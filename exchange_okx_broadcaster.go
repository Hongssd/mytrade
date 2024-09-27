package mytrade

import (
	"context"
	"errors"
	"github.com/Hongssd/myokxapi"
	"sync"
	"time"
)

// 广播 用于异步接收订单操作的返回结果
type okxOrderBroadcaster struct {
	accountType      string
	okxWsAccount     *myokxapi.PrivateWsStreamClient
	currentSubscribe *myokxapi.Subscription[myokxapi.WsOrders]
	subscribers      *MySyncMap[string, *okxOrderSubscriber]
	keys             *MySyncMap[*okxOrderSubscriber, string]
	mu               sync.RWMutex
}

// 订阅者 用于异步接收订单操作的返回结果
type okxOrderSubscriber struct {
	symbol        string
	clientOrderId string
	ch            *subscription[Order]
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
	if *ob == nil {
		newOb, err := o.newOrderBroadcaster(accountType)
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

// 等待广播消息，超时返回
func (o *OkxTradeEngine) waitSubscribeReturn(sub *okxOrderSubscriber, timeout time.Duration) (*Order, error) {
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

// 新建广播者，本质上是一条ws链接接收所有的订单推送结果进行广播
func (o *OkxTradeEngine) newOrderBroadcaster(accountType string) (*okxOrderBroadcaster, error) {
	broadcaster := &okxOrderBroadcaster{
		accountType:  accountType,
		okxWsAccount: okx.NewPrivateWsStreamClient(),
		subscribers:  GetPointer(NewMySyncMap[string, *okxOrderSubscriber]()),
		keys:         GetPointer(NewMySyncMap[*okxOrderSubscriber, string]()),
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

	sub, err := broadcaster.okxWsAccount.SubscribeOrders(accountType, "", "")
	if err != nil {
		return nil, err
	}

	broadcaster.currentSubscribe = sub

	go func() {
		for {
			sub := broadcaster.currentSubscribe
			select {
			case err := <-sub.ErrChan():
				broadcaster.subscribers.Range(func(key string, value *okxOrderSubscriber) bool {
					value.ch.ErrChan() <- err

					return true
				})
			case result := <-sub.ResultChan():
				//log.Infof("订单频道订阅接收到消息：%s", result)
				order := o.handleOrderFromWsOrder(result)
				order.AccountType = broadcaster.accountType
				broadcaster.subscribers.Range(func(key string, value *okxOrderSubscriber) bool {
					if value.clientOrderId == "" || order.ClientOrderId == value.clientOrderId {
						value.ch.ResultChan() <- *order
					}
					return true
				})
			case <-sub.CloseChan():
				broadcaster.subscribers.Range(func(key string, value *okxOrderSubscriber) bool {
					value.ch.CloseChan() <- struct{}{}
					return true
				})
				return
			}
		}
	}()

	return broadcaster, nil
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
