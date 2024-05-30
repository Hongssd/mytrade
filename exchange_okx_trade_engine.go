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

	broadcaster *okxOrderBroadcaster
}

// 广播 用于异步接收订单操作的返回结果
type okxOrderBroadcaster struct {
	okxWsAccount     *myokxapi.PrivateWsStreamClient
	currentSubscribe *myokxapi.Subscription[myokxapi.WsOrders]
	subscribers      *MySyncMap[string, *okxOrderSubscriber]
	keys             *MySyncMap[*okxOrderSubscriber, string]
	mu               sync.RWMutex
}

// 订阅者 用于异步接收订单操作的返回结果
type okxOrderSubscriber struct {
	clientOrderId string
	accountType   string
	ch            chan Order
}

// 新建订阅者
func (o *OkxTradeEngine) newOrderSubscriber(clientOrderId, accountType string) (*okxOrderSubscriber, error) {
	var err error
	sub := &okxOrderSubscriber{
		clientOrderId: clientOrderId,
		accountType:   accountType,
		ch:            make(chan Order, 10),
	}
	if o.broadcaster == nil {
		o.broadcaster, err = o.newOrderBroadcaster()
		if err != nil {
			return nil, err
		}
	}

	o.broadcaster.mu.Lock()
	defer o.broadcaster.mu.Unlock()

	key := clientOrderId
	o.broadcaster.subscribers.Store(key, sub)
	o.broadcaster.keys.Store(sub, key)
	return sub, nil
}

// 关闭订阅者
func (o *OkxTradeEngine) closeSubscribe(sub *okxOrderSubscriber) {
	o.broadcaster.mu.Lock()
	defer o.broadcaster.mu.Unlock()

	key, _ := o.broadcaster.keys.Load(sub)
	o.broadcaster.subscribers.Delete(key)
	o.broadcaster.keys.Delete(sub)
}

// 等待广播消息，超时返回
func (o *OkxTradeEngine) waitSubscribeReturn(sub *okxOrderSubscriber, timeout time.Duration) (*Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	select {
	case <-ctx.Done():
		return nil, errors.New("api msg timeout")
	case order := <-sub.ch:
		order.AccountType = sub.accountType
		return &order, nil
	}
}

// 新建广播者，本质上是一条ws链接接收所有的订单推送结果进行广播
func (o *OkxTradeEngine) newOrderBroadcaster() (*okxOrderBroadcaster, error) {
	broadcaster := &okxOrderBroadcaster{
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

	sub, err := broadcaster.okxWsAccount.SubscribeOrders("ANY", "", "")
	if err != nil {
		return nil, err
	}

	broadcaster.currentSubscribe = sub

	go func() {
		for {
			sub := o.broadcaster.currentSubscribe
			select {
			case err := <-sub.ErrChan():
				log.Error(err)
			case result := <-sub.ResultChan():
				//log.Infof("订单频道订阅接收到消息：%s", result)
				order := o.handleOrderFromWsOrder(result)
				broadcaster.subscribers.Range(func(key string, value *okxOrderSubscriber) bool {
					value.ch <- *order
					return true
				})
			case <-sub.CloseChan():
				log.Info("订阅已关闭: ", sub.Args)
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
	//TODO implement me
	panic("implement me")
}

func (o *OkxTradeEngine) QueryOrder(req *QueryOrderParam) (*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *OkxTradeEngine) QueryTrades(req *QueryTradeParam) ([]*Trade, error) {
	//TODO implement me
	panic("implement me")
}

func (o *OkxTradeEngine) CreateOrder(req *OrderParam) (*Order, error) {
	//获取API
	api, err := o.apiOrderCreate(req)
	if err != nil {
		return nil, err
	}

	//创建订阅
	sub, err := o.newOrderSubscriber(req.ClientOrderId, req.AccountType)
	if err != nil {
		return nil, err
	}
	defer o.closeSubscribe(sub)

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
	api, err := o.apiOrderAmend(req)
	if err != nil {
		return nil, err
	}

	//创建订阅
	sub, err := o.newOrderSubscriber(req.ClientOrderId, req.AccountType)
	if err != nil {
		return nil, err
	}
	defer o.closeSubscribe(sub)

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
	api, err := o.apiOrderCancel(req)
	if err != nil {
		return nil, err
	}
	//创建订阅
	sub, err := o.newOrderSubscriber(req.ClientOrderId, req.AccountType)
	if err != nil {
		return nil, err
	}
	defer o.closeSubscribe(sub)
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
	api, err := o.apiBatchOrderCreate(reqs)
	if err != nil {
		return nil, err
	}

	subs := make([]*okxOrderSubscriber, 0, len(reqs))
	//批量创建订阅
	for _, req := range reqs {
		sub, err := o.newOrderSubscriber(req.ClientOrderId, req.AccountType)
		if err != nil {
			return nil, err
		}

		subs = append(subs, sub)
	}
	defer func() {
		for _, sub := range subs {
			o.closeSubscribe(sub)
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
	api, err := o.apiBatchOrderAmend(reqs)
	if err != nil {
		return nil, err
	}

	subs := make([]*okxOrderSubscriber, 0, len(reqs))
	//批量创建订阅
	for _, req := range reqs {
		sub, err := o.newOrderSubscriber(req.ClientOrderId, req.AccountType)
		if err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}
	defer func() {
		for _, sub := range subs {
			o.closeSubscribe(sub)
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
	api, err := o.apiBatchOrderCancel(reqs)
	if err != nil {
		return nil, err
	}

	subs := make([]*okxOrderSubscriber, 0, len(reqs))
	//批量创建订阅
	for _, req := range reqs {
		sub, err := o.newOrderSubscriber(req.ClientOrderId, req.AccountType)
		if err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}
	defer func() {
		for _, sub := range subs {
			o.closeSubscribe(sub)
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
	//TODO implement me
	panic("implement me")
}

func (o *OkxTradeEngine) SubscribeOrder(req *SubscribeOrderParam) (TradeSubscribe[Order], error) {
	//TODO implement me
	panic("implement me")
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
