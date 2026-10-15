package mytrade

import (
	"context"
	"errors"
	"github.com/Hongssd/myokxapi"
	"sync"
	"sync/atomic"
	"time"
)

// 广播 用于异步接收订单操作的返回结果
type okxOrderBroadcaster struct {
	accountType string
	subscribers *MySyncMap[string, *okxOrderSubscriber]
	keys        *MySyncMap[*okxOrderSubscriber, string]
	mu          sync.RWMutex

	// wsMu guards the fields below, which change when the connection is rebuilt.
	wsMu             sync.RWMutex
	okxWsAccount     *myokxapi.PrivateWsStreamClient
	currentSubscribe *myokxapi.Subscription[myokxapi.WsOrders]
	stopRun          chan struct{}

	health     okxAckHealth
	recovering atomic.Bool
}

// 订阅者 用于异步接收订单操作的返回结果
type okxOrderSubscriber struct {
	symbol        string
	clientOrderId string
	ch            *subscription[Order]
	owner         *okxOrderBroadcaster
}

// 新建订阅者
func (o *OkxTradeEngine) newOrderSubscriber(ob **okxOrderBroadcaster, clientOrderId, accountType, symbol string) (*okxOrderSubscriber, error) {

	sub := &okxOrderSubscriber{
		symbol:        symbol,
		clientOrderId: clientOrderId,
		ch: &subscription[Order]{
			resultChan: make(chan Order, 100),
			errChan:    make(chan error, 10),
			closeChan:  make(chan struct{}, 10),
		},
	}

	// Without the lock two callers could both see nil and open two connections.
	o.broadcasterInitMu.Lock()
	if *ob == nil {
		newOb, err := o.newOrderBroadcaster(accountType)
		if err != nil {
			o.broadcasterInitMu.Unlock()
			return nil, err
		}
		*ob = newOb
	}
	b := *ob
	o.broadcasterInitMu.Unlock()

	sub.owner = b
	b.mu.Lock()
	defer b.mu.Unlock()

	key := clientOrderId
	b.subscribers.Store(key, sub)
	b.keys.Store(sub, key)
	return sub, nil
}

// 等待广播消息，超时返回
func (o *OkxTradeEngine) waitSubscribeReturn(sub *okxOrderSubscriber, timeout time.Duration) (*Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	select {
	//超时返回
	case <-ctx.Done():
		if sub.owner != nil {
			sub.owner.onAckTimeout(o, sub)
		}
		return nil, errors.New("api msg timeout")
	case <-sub.ch.CloseChan():
		//链接关闭的情况，使用接口查询
		return o.QueryOrder(o.NewQueryOrderReq().SetSymbol(sub.symbol).SetClientOrderId(sub.clientOrderId))
	case order := <-sub.ch.ResultChan():
		if sub.owner != nil {
			sub.owner.health.onSuccess()
		}
		return &order, nil
	}
}

// connectOrdersWs opens a private ws, logs in and subscribes to orders.
func (o *OkxTradeEngine) connectOrdersWs(accountType string) (*myokxapi.PrivateWsStreamClient, *myokxapi.Subscription[myokxapi.WsOrders], error) {
	ws := okx.NewPrivateWsStreamClient()
	if err := ws.OpenConn(); err != nil {
		return nil, nil, err
	}

	err := ws.Login(okx.NewRestClient(o.apiKey, o.secretKey, o.passphrase))
	if err != nil {
		closeOkxWsQuietly(ws)
		return nil, nil, err
	}

	if accountType == "SPOT" {
		accountType = "ANY"
	}

	sub, err := ws.SubscribeOrders(accountType, "", "")
	if err != nil {
		closeOkxWsQuietly(ws)
		return nil, nil, err
	}
	return ws, sub, nil
}

// closeOkxWsQuietly closes a ws that is being discarded. Close notifies the
// client's own subscriptions and may block on them, so it runs detached.
func closeOkxWsQuietly(ws *myokxapi.PrivateWsStreamClient) {
	go func() {
		defer func() { _ = recover() }()
		if err := ws.Close(); err != nil {
			log.Warnf("[okx-orders] close discarded ws failed: %v", err)
		}
	}()
}

// 新建广播者，本质上是一条ws链接接收所有的订单推送结果进行广播
func (o *OkxTradeEngine) newOrderBroadcaster(accountType string) (*okxOrderBroadcaster, error) {
	broadcaster := &okxOrderBroadcaster{
		accountType: accountType,
		subscribers: GetPointer(NewMySyncMap[string, *okxOrderSubscriber]()),
		keys:        GetPointer(NewMySyncMap[*okxOrderSubscriber, string]()),
		mu:          sync.RWMutex{},
	}
	ws, sub, err := o.connectOrdersWs(accountType)
	if err != nil {
		return nil, err
	}
	broadcaster.okxWsAccount = ws
	broadcaster.currentSubscribe = sub
	broadcaster.stopRun = make(chan struct{})

	go o.runOrderBroadcaster(broadcaster, ws, sub, broadcaster.stopRun)
	go o.watchOrderBroadcaster(broadcaster)

	return broadcaster, nil
}

// runOrderBroadcaster fans pushes of one ws generation out to the subscribers.
// It ends when its generation is replaced (stop) or its ws is closed.
func (o *OkxTradeEngine) runOrderBroadcaster(b *okxOrderBroadcaster, ws *myokxapi.PrivateWsStreamClient,
	sub *myokxapi.Subscription[myokxapi.WsOrders], stop chan struct{}) {
	for {
		select {
		case <-stop:
			return
		case err := <-sub.ErrChan():
			log.Errorf("[okx-orders] %s subscription error: %v", b.accountType, err)
			b.subscribers.Range(func(key string, value *okxOrderSubscriber) bool {
				pushToOkxSubscriber(value, nil, err)
				return true
			})
		case result := <-sub.ResultChan():
			order := o.handleOrderFromWsOrder(result)
			order.AccountType = b.accountType
			b.subscribers.Range(func(key string, value *okxOrderSubscriber) bool {
				if value.clientOrderId == "" || order.ClientOrderId == value.clientOrderId {
					pushToOkxSubscriber(value, order, nil)
				}
				return true
			})
		case <-sub.CloseChan():
			// Not caused by a rebuild (that stops this loop first), so the
			// connection was closed from outside. Waiters fall back to REST and the
			// channel is rebuilt rather than staying dead.
			log.Warnf("[okx-orders] %s subscription closed unexpectedly, rebuilding", b.accountType)
			b.subscribers.Range(func(key string, value *okxOrderSubscriber) bool {
				if value.clientOrderId != "" {
					select {
					case value.ch.CloseChan() <- struct{}{}:
					default:
					}
				}
				return true
			})
			go o.recoverOrderBroadcaster(b, okxRecoverRebuild, "subscription closed")
			return
		}
	}
}

// pushToOkxSubscriber never blocks the shared reader on a stalled consumer.
func pushToOkxSubscriber(s *okxOrderSubscriber, order *Order, err error) {
	if err != nil {
		select {
		case s.ch.ErrChan() <- err:
		default:
			log.Warnf("[okx-orders] subscriber error chan full, dropped: clOrdId=%s", s.clientOrderId)
		}
		return
	}
	select {
	case s.ch.ResultChan() <- *order:
		return
	default:
	}
	timer := time.NewTimer(2 * time.Second)
	defer timer.Stop()
	select {
	case s.ch.ResultChan() <- *order:
	case <-timer.C:
		log.Errorf("[okx-orders] subscriber stalled, push dropped: subClOrdId=%q orderClOrdId=%s status=%v",
			s.clientOrderId, order.ClientOrderId, order.Status)
	}
}

// onAckTimeout counts a request that REST accepted but whose push never came.
func (b *okxOrderBroadcaster) onAckTimeout(o *OkxTradeEngine, sub *okxOrderSubscriber) {
	level, ok := b.health.onTimeout(time.Now())
	log.Warnf("[okx-orders] %s order push timeout: clOrdId=%s strikes=%d recover=%v",
		b.accountType, sub.clientOrderId, b.health.strikeCount(), ok)
	if ok {
		go o.recoverOrderBroadcaster(b, level, "consecutive order push timeouts")
	}
}

// watchOrderBroadcaster is the active probe. myokxapi already pings and closes a
// silent socket, so this only catches a connection whose own reconnect is stuck.
func (o *OkxTradeEngine) watchOrderBroadcaster(b *okxOrderBroadcaster) {
	ticker := time.NewTicker(okxHealthInterval)
	defer ticker.Stop()
	for range ticker.C {
		b.wsMu.RLock()
		ws := b.okxWsAccount
		b.wsMu.RUnlock()
		if ws == nil {
			continue
		}
		idle := time.Since(ws.LastMessageTime())
		if idle < okxWsSilentLimit {
			continue
		}
		if level, ok := b.health.onSilent(time.Now()); ok {
			log.Warnf("[okx-orders] %s ws silent for %s", b.accountType, idle.Round(time.Second))
			go o.recoverOrderBroadcaster(b, level, "ws silent")
		}
	}
}

// recoverOrderBroadcaster runs one recovery at a time.
func (o *OkxTradeEngine) recoverOrderBroadcaster(b *okxOrderBroadcaster, level okxRecoverLevel, reason string) {
	if !b.recovering.CompareAndSwap(false, true) {
		return
	}
	defer b.recovering.Store(false)

	log.Warnf("[okx-orders] %s recover start: level=%s reason=%s", b.accountType, level, reason)

	if level == okxRecoverReconnect {
		b.wsMu.RLock()
		ws := b.okxWsAccount
		b.wsMu.RUnlock()
		if err := ws.ForceReconnect(reason); err == nil {
			return
		} else {
			log.Warnf("[okx-orders] %s force reconnect failed, rebuilding: %v", b.accountType, err)
		}
	}

	ws, sub, err := o.connectOrdersWs(b.accountType)
	if err != nil {
		log.Errorf("[okx-orders] %s rebuild failed: %v", b.accountType, err)
		return
	}

	b.wsMu.Lock()
	oldWs, oldStop := b.okxWsAccount, b.stopRun
	b.okxWsAccount = ws
	b.currentSubscribe = sub
	b.stopRun = make(chan struct{})
	newStop := b.stopRun
	b.wsMu.Unlock()

	// Stop the old reader first so closing the old ws is not seen as a failure.
	if oldStop != nil {
		close(oldStop)
	}
	go o.runOrderBroadcaster(b, ws, sub, newStop)
	if oldWs != nil {
		closeOkxWsQuietly(oldWs)
	}
	log.Warnf("[okx-orders] %s rebuild success, subscribers kept", b.accountType)
}

func (o *OkxTradeEngine) getBroadcastFromAccountType(accountType string) **okxOrderBroadcaster {
	switch OkxAccountType(accountType) {
	case OKX_AC_SPOT:
		return &o.broadcasterSpot
	case OKX_AC_MARGIN:
		return &o.broadcasterMargin
	case OKX_AC_SWAP:
		return &o.broadcasterSwap
	case OKX_AC_FUTURES:
		return &o.broadcasterFuture
	default:
		return nil
	}
}

// 订阅者 用于异步接收订单操作的返回结果
type okxOrderAlgoSubscriber struct {
	symbol        string
	clientOrderId string
	ch            *subscription[Order]
}

// 广播 用于异步接收策略订单操作的返回结果
type okxOrderAlgoBroadcaster struct {
	accountType      string
	okxWsAccount     *myokxapi.BusinessWsStreamClient
	currentSubscribe *myokxapi.Subscription[myokxapi.WsOrdersAlgo]
	subscribers      *MySyncMap[string, *okxOrderAlgoSubscriber]
	keys             *MySyncMap[*okxOrderAlgoSubscriber, string]
	mu               sync.RWMutex
}

// 新建广播者
func (o *OkxTradeEngine) newOrderAlgoBroadcaster(accountType string) (*okxOrderAlgoBroadcaster, error) {
	broadcaster := &okxOrderAlgoBroadcaster{
		accountType:  accountType,
		okxWsAccount: okx.NewBusinessWsStreamClient(),
		subscribers:  GetPointer(NewMySyncMap[string, *okxOrderAlgoSubscriber]()),
		keys:         GetPointer(NewMySyncMap[*okxOrderAlgoSubscriber, string]()),
		mu:           sync.RWMutex{},
	}
	err := broadcaster.okxWsAccount.OpenConn()
	if err != nil {
		return nil, err
	}

	err = broadcaster.okxWsAccount.Login(okx.NewRestClient(o.apiKey, o.secretKey, o.passphrase))
	if err != nil {
		return nil, err
	}

	if accountType == "SPOT" {
		accountType = "ANY"
	}

	sub, err := broadcaster.okxWsAccount.SubscribeOrdersAlgo(accountType, "", "")
	if err != nil {
		return nil, err
	}

	broadcaster.currentSubscribe = sub

	go func() {
		for {
			sub := broadcaster.currentSubscribe
			select {
			case err := <-sub.ErrChan():
				broadcaster.subscribers.Range(func(key string, value *okxOrderAlgoSubscriber) bool {
					value.ch.ErrChan() <- err

					return true
				})
			case result := <-sub.ResultChan():
				//log.Infof("订单频道订阅接收到消息：%s", result)

				order := o.handleOrderFromWsOrderAlgo(result)
				order.AccountType = broadcaster.accountType
				if order.AccountType != OKX_AC_SPOT.String() {
					order.IsMargin = false
				}
				broadcaster.subscribers.Range(func(key string, value *okxOrderAlgoSubscriber) bool {
					if value.clientOrderId == "" || order.ClientOrderId == value.clientOrderId {
						value.ch.ResultChan() <- *order
					}
					return true
				})
			case <-sub.CloseChan():
				broadcaster.subscribers.Range(func(key string, value *okxOrderAlgoSubscriber) bool {
					value.ch.CloseChan() <- struct{}{}
					return true
				})
				return
			}
		}
	}()

	return broadcaster, nil
}

// 新建策略委托订单订阅者
func (o *OkxTradeEngine) newOrderAlgoSubscriber(ob **okxOrderAlgoBroadcaster, clientOrderId, accountType, symbol string) (*okxOrderAlgoSubscriber, error) {

	sub := &okxOrderAlgoSubscriber{
		symbol:        symbol,
		clientOrderId: clientOrderId,
		ch: &subscription[Order]{
			resultChan: make(chan Order, 100),
			errChan:    make(chan error, 10),
			closeChan:  make(chan struct{}, 10),
		},
	}
	if *ob == nil {
		newOb, err := o.newOrderAlgoBroadcaster(accountType)
		if err != nil {
			return nil, err
		}
		*ob = newOb
	}

	(*ob).mu.Lock()
	defer (*ob).mu.Unlock()

	key := clientOrderId
	(*ob).subscribers.Store(key, sub)
	(*ob).keys.Store(sub, key)
	return sub, nil
}

// 关闭策略委托订单订阅者
func (o *OkxTradeEngine) closeSubscribeAlgo(b **okxOrderAlgoBroadcaster, sub *okxOrderAlgoSubscriber) {
	(*b).mu.Lock()
	defer (*b).mu.Unlock()
	key, _ := (*b).keys.Load(sub)
	(*b).subscribers.Delete(key)
	(*b).keys.Delete(sub)
}

// 等待策略委托订单广播消息，超时返回
func (o *OkxTradeEngine) waitOrderAlgoSubscribeReturn(sub *okxOrderAlgoSubscriber, timeout time.Duration) (*Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	select {
	//超时返回
	case <-ctx.Done():
		return nil, errors.New("api msg timeout")
	case <-sub.ch.CloseChan():
		//链接关闭的情况，使用接口查询
		return o.QueryOrder(o.NewQueryOrderReq().SetSymbol(sub.symbol).SetClientOrderId(sub.clientOrderId))
	case order := <-sub.ch.ResultChan():
		return &order, nil
	}
}

// 关闭订阅者
func (o *OkxTradeEngine) closeSubscribe(b **okxOrderBroadcaster, sub *okxOrderSubscriber) {
	(*b).mu.Lock()
	defer (*b).mu.Unlock()
	key, _ := (*b).keys.Load(sub)
	(*b).subscribers.Delete(key)
	(*b).keys.Delete(sub)
}

func (o *OkxTradeEngine) getOrderAlgoBroadcastFromAccountType(accountType string) **okxOrderAlgoBroadcaster {
	switch OkxAccountType(accountType) {
	case OKX_AC_SPOT:
		return &o.orderAlgoBroadcasterSpot
	case OKX_AC_MARGIN:
		return &o.orderAlgoBroadcasterMargin
	case OKX_AC_SWAP:
		return &o.orderAlgoBroadcasterSwap
	case OKX_AC_FUTURES:
		return &o.orderAlgoBroadcasterFuture
	default:
		return nil
	}
}
