package mytrade

import (
	"context"
	"errors"
	"github.com/Hongssd/myokxapi"
	"sync"
	"time"
)

type OkxTradeEngine struct {
	exchangeBase

	okxConverter OkxEnumConverter
	apiKey       string
	secretKey    string
	passphrase   string

	broadcasterSpot   *okxOrderBroadcaster
	broadcasterSwap   *okxOrderBroadcaster
	broadcasterFuture *okxOrderBroadcaster

	parent *OkxExchange
}

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

// 关闭订阅者
func (o *OkxTradeEngine) closeSubscribe(b **okxOrderBroadcaster, sub *okxOrderSubscriber) {
	(*b).mu.Lock()
	defer (*b).mu.Unlock()
	key, _ := (*b).keys.Load(sub)
	(*b).subscribers.Delete(key)
	(*b).keys.Delete(sub)
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

func (o *OkxTradeEngine) NewOrderReq() *OrderParam {
	return &OrderParam{}
}

func (o *OkxTradeEngine) NewQueryOrderReq() *QueryOrderParam {
	return &QueryOrderParam{}
}

func (o *OkxTradeEngine) NewQueryTradeReq() *QueryTradeParam {
	return &QueryTradeParam{}
}

func (o *OkxTradeEngine) QueryOpenOrders(req *QueryOrderParam) ([]*Order, error) {
	api := o.apiQueryOpenOrders(req)
	res, err := api.Do()
	if err != nil {
		return nil, err
	}

	orders, err := o.handleOrdersFromQueryOpenOrders(req, res)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *OkxTradeEngine) QueryOrder(req *QueryOrderParam) (*Order, error) {
	api := o.apiQueryOrder(req)
	res, err := api.Do()
	if err != nil {
		return nil, err
	}

	order, err := o.handleOrderFromQueryOrderGet(req, res)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (o *OkxTradeEngine) QueryTrades(req *QueryTradeParam) ([]*Trade, error) {
	api := o.apiQueryTrades(req)
	res, err := api.Do()
	if err != nil {
		return nil, err
	}

	trades, err := o.handleTradesFromQueryTrades(req, res)
	if err != nil {
		return nil, err
	}

	return trades, nil
}

func (o *OkxTradeEngine) CreateOrder(req *OrderParam) (*Order, error) {
	//获取API
	api := o.apiOrderCreate(req)

	b := o.getBoardcastFromAccountType(req.AccountType)

	//创建订阅
	sub, err := o.newOrderSubscriber(b, req.ClientOrderId, req.AccountType, req.Symbol)
	if err != nil {
		return nil, err
	}
	defer o.closeSubscribe(b, sub)

	//执行API
	res, err := api.Do()
	if err != nil {
		return nil, err
	}

	//处理API返回值
	_, err = o.handleOrderFromOrderCreate(req, res)
	if err != nil {
		return nil, err
	}

	//异步接收ws结果，10秒超时
	return o.waitSubscribeReturn(sub, 10*time.Second)
}
func (o *OkxTradeEngine) AmendOrder(req *OrderParam) (*Order, error) {
	//获取API
	api := o.apiOrderAmend(req)

	b := o.getBoardcastFromAccountType(req.AccountType)

	//创建订阅
	sub, err := o.newOrderSubscriber(b, req.ClientOrderId, req.AccountType, req.Symbol)
	if err != nil {
		return nil, err
	}
	defer o.closeSubscribe(b, sub)

	//执行API
	res, err := api.Do()
	if err != nil {
		return nil, err
	}

	//处理API返回值
	_, err = o.handleOrderFromOrderAmend(req, res)
	if err != nil {
		return nil, err
	}

	//异步接收ws结果，10秒超时
	return o.waitSubscribeReturn(sub, 10*time.Second)
}
func (o *OkxTradeEngine) CancelOrder(req *OrderParam) (*Order, error) {
	//获取API
	api := o.apiOrderCancel(req)

	b := o.getBoardcastFromAccountType(req.AccountType)

	//创建订阅
	sub, err := o.newOrderSubscriber(b, req.ClientOrderId, req.AccountType, req.Symbol)
	if err != nil {
		return nil, err
	}
	defer o.closeSubscribe(b, sub)
	//执行API
	res, err := api.Do()
	if err != nil {
		return nil, err
	}

	//处理API返回值
	_, err = o.handleOrderFromOrderCancel(req, res)
	if err != nil {
		return nil, err
	}

	//异步接收ws结果，10秒超时
	return o.waitSubscribeReturn(sub, 10*time.Second)
}

func (o *OkxTradeEngine) CreateOrders(reqs []*OrderParam) ([]*Order, error) {
	//获取API
	api := o.apiBatchOrderCreate(reqs)

	var defers []func()

	subs := make([]*okxOrderSubscriber, 0, len(reqs))
	//批量创建订阅
	for _, req := range reqs {
		b := o.getBoardcastFromAccountType(req.AccountType)

		sub, err := o.newOrderSubscriber(b, req.ClientOrderId, req.AccountType, req.Symbol)
		if err != nil {
			return nil, err
		}

		subs = append(subs, sub)
		defers = append(defers, func() {
			o.closeSubscribe(b, sub)
		})
	}
	defer func() {
		for _, d := range defers {
			d()
		}
	}()
	//执行API
	res, err := api.Do()
	if err != nil {
		return nil, err
	}
	//处理API返回值
	_, err = o.handleOrderFromBatchOrderCreate(reqs, res)
	if err != nil {
		return nil, err
	}

	var orders []*Order
	//批量异步接收ws结果，10秒超时
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, sub := range subs {
		wg.Add(1)
		go func(sub *okxOrderSubscriber) {
			defer wg.Done()
			order, err := o.waitSubscribeReturn(sub, 10*time.Second)
			if err != nil {
				log.Error(err)
			}
			mu.Lock()
			orders = append(orders, order)
			mu.Unlock()
		}(sub)
	}

	wg.Wait()
	return orders, nil

}
func (o *OkxTradeEngine) AmendOrders(reqs []*OrderParam) ([]*Order, error) {
	//获取API
	api := o.apiBatchOrderAmend(reqs)

	var defers []func()

	subs := make([]*okxOrderSubscriber, 0, len(reqs))
	//批量创建订阅
	for _, req := range reqs {
		b := o.getBoardcastFromAccountType(req.AccountType)

		sub, err := o.newOrderSubscriber(b, req.ClientOrderId, req.AccountType, req.Symbol)
		if err != nil {
			return nil, err
		}
		subs = append(subs, sub)

		defers = append(defers, func() {
			o.closeSubscribe(b, sub)
		})
	}
	defer func() {
		for _, d := range defers {
			d()
		}
	}()
	//执行API
	res, err := api.Do()
	if err != nil {
		return nil, err
	}

	//处理API返回值
	_, err = o.handleOrderFromBatchOrderAmend(reqs, res)
	if err != nil {
		return nil, err
	}

	var orders []*Order

	//批量异步接收ws结果，10秒超时
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, sub := range subs {
		wg.Add(1)
		go func(sub *okxOrderSubscriber) {
			defer wg.Done()
			order, err := o.waitSubscribeReturn(sub, 10*time.Second)
			if err != nil {
				log.Error(err)
			}
			mu.Lock()
			orders = append(orders, order)
			mu.Unlock()
		}(sub)

	}
	wg.Wait()
	return orders, nil
}
func (o *OkxTradeEngine) CancelOrders(reqs []*OrderParam) ([]*Order, error) {
	//获取API
	api := o.apiBatchOrderCancel(reqs)

	var defers []func()

	subs := make([]*okxOrderSubscriber, 0, len(reqs))
	//批量创建订阅
	for _, req := range reqs {
		b := o.getBoardcastFromAccountType(req.AccountType)

		sub, err := o.newOrderSubscriber(b, req.ClientOrderId, req.AccountType, req.Symbol)
		if err != nil {
			return nil, err
		}
		subs = append(subs, sub)

		defers = append(defers, func() {
			o.closeSubscribe(b, sub)
		})
	}
	defer func() {
		for _, d := range defers {
			d()
		}
	}()
	//执行API
	res, err := api.Do()
	if err != nil {
		return nil, err
	}

	//处理API返回值
	_, err = o.handleOrderFromBatchOrderCancel(reqs, res)
	if err != nil {
		return nil, err
	}

	var orders []*Order

	//批量异步接收ws结果，10秒超时
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, sub := range subs {
		wg.Add(1)
		go func(sub *okxOrderSubscriber) {
			defer wg.Done()
			order, err := o.waitSubscribeReturn(sub, 10*time.Second)
			if err != nil {
				log.Error(err)
			}
			mu.Lock()
			orders = append(orders, order)
			mu.Unlock()
		}(sub)

	}

	wg.Wait()
	return orders, nil
}

func (o *OkxTradeEngine) NewSubscribeOrderReq() *SubscribeOrderParam {
	return &SubscribeOrderParam{}
}

func (o *OkxTradeEngine) SubscribeOrder(req *SubscribeOrderParam) (TradeSubscribe[Order], error) {

	switch OkxAccountType(req.AccountType) {
	case OKX_AC_SPOT, OKX_AC_SWAP, OKX_AC_FUTURES:
	default:
		return nil, ErrorAccountType
	}
	b := o.getBoardcastFromAccountType(req.AccountType)

	sub, err := o.newOrderSubscriber(b, "", req.AccountType, "")
	if err != nil {
		return nil, err
	}

	middleSub := &subscription[Order]{
		resultChan: make(chan Order, 100),
		errChan:    make(chan error, 10),
		closeChan:  make(chan struct{}, 10),
	}

	//循环将订单数据中转到目标订阅
	go func() {
		for {
			select {
			case <-sub.ch.CloseChan():
				middleSub.closeChan <- struct{}{}
				return
			case err := <-sub.ch.ErrChan():
				middleSub.errChan <- err
			case order := <-sub.ch.ResultChan():
				middleSub.resultChan <- order
			}
		}
	}()

	return middleSub, nil

}

func (o *OkxTradeEngine) WsCreateOrder(req *OrderParam) (*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *OkxTradeEngine) WsAmendOrder(req *OrderParam) (*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *OkxTradeEngine) WsCancelOrder(req *OrderParam) (*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *OkxTradeEngine) WsCreateOrders(reqs []*OrderParam) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *OkxTradeEngine) WsAmendOrders(reqs []*OrderParam) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *OkxTradeEngine) WsCancelOrders(reqs []*OrderParam) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}
