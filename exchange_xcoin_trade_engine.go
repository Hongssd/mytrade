package mytrade

import (
	"fmt"
	"sync"
	"time"
)

type XcoinTradeEngine struct {
	ExchangeBase

	xcoinConverter XcoinEnumConverter
	apiKey         string
	apiSecret      string

	broadcasterSpot            *xcoinOrderBroadcaster
	broadcasterLinearPerpetual *xcoinOrderBroadcaster
	broadcasterLinearFutures   *xcoinOrderBroadcaster
}

func (x *XcoinTradeEngine) NewOrderReq() *OrderParam {
	return &OrderParam{}
}

func (x *XcoinTradeEngine) NewQueryOrderReq() *QueryOrderParam {
	return &QueryOrderParam{}
}

func (x *XcoinTradeEngine) NewQueryTradeReq() *QueryTradeParam {
	return &QueryTradeParam{}
}

func (x *XcoinTradeEngine) QueryOpenOrders(req *QueryOrderParam) ([]*Order, error) {
	if err := x.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}

	api := x.apiQueryOpenOrders(req)
	res, err := api.Do()
	if err != nil {
		return nil, err
	}
	return x.handleOrdersFromQueryOpenOrders(req, res), nil
}

func (x *XcoinTradeEngine) QueryOrder(req *QueryOrderParam) (*Order, error) {
	if err := x.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}

	api := x.apiQueryOrder(req)
	res, err := api.Do()
	if err != nil {
		return nil, err
	}
	o, err := x.handleOrderFromQueryOrder(req, res)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (x *XcoinTradeEngine) QueryOrders(req *QueryOrderParam) ([]*Order, error) {
	if err := x.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}

	api := x.apiQueryOrders(req)
	res, err := api.Do()
	if err != nil {
		return nil, err
	}
	return x.handleOrdersFromQueryOrders(req, res), nil
}

func (x *XcoinTradeEngine) QueryTrades(req *QueryTradeParam) ([]*Trade, error) {
	if err := x.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}

	api := x.apiQueryTrades(req)
	res, err := api.Do()
	if err != nil {
		return nil, err
	}
	return x.handleTradesFromQueryTrades(req, res), nil
}

func (x *XcoinTradeEngine) CreateOrder(req *OrderParam) (*Order, error) {
	if err := x.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}

	api := x.apiOrderCreate(req)

	b := x.getBroadcastFromAccountType(req.AccountType)
	// 创建订阅
	sub, err := x.newOrderSubscriber(b, req.ClientOrderId, req.OrderId, req.AccountType, req.Symbol)
	if err != nil {
		return nil, err
	}
	defer x.closeSubscribe(b, sub)

	// 执行API
	res, err := api.Do()
	if err != nil {
		return nil, err
	}
	// 处理API返回值
	_, err = x.handleOrderFromOrderCreate(req, res)
	if err != nil {
		return nil, err
	}
	// 异步接收ws结果，1秒超时
	return x.waitSubscribeReturn(sub, 1*time.Second)
}

func (x *XcoinTradeEngine) AmendOrder(req *OrderParam) (*Order, error) {
	return nil, ErrorNotSupport
}

func (x *XcoinTradeEngine) CancelOrder(req *OrderParam) (*Order, error) {
	if err := x.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}

	api := xcoin.NewRestClient(x.apiKey, x.apiSecret).PrivateRestClient().
		NewPrivateRestTradeCancelOrder().
		Symbol(req.Symbol)
	if req.OrderId != "" {
		api.OrderId(req.OrderId)
	}
	if req.ClientOrderId != "" {
		api.ClientOrderId(req.ClientOrderId)
	}

	b := x.getBroadcastFromAccountType(req.AccountType)
	// 创建订阅
	sub, err := x.newOrderSubscriber(b, req.ClientOrderId, req.OrderId, req.AccountType, req.Symbol)
	if err != nil {
		return nil, err
	}
	defer x.closeSubscribe(b, sub)

	// 执行API
	res, err := api.Do()
	if err != nil {
		return nil, err
	}
	// 处理API返回值
	if res.Code != "0" {
		return nil, fmt.Errorf("[%s]:%s", res.Code, res.Msg)
	}
	// 异步接收ws结果，1秒超时
	return x.waitSubscribeReturn(sub, 1*time.Second)
}

func (x *XcoinTradeEngine) CreateOrders(reqs []*OrderParam) ([]*Order, error) {
	if err := x.restBatchPreCheck(reqs); err != nil {
		return nil, err
	}

	api := x.apiBatchOrderCreate(reqs)

	var defers []func()
	subs := make([]*xcoinOrderSubscriber, 0, len(reqs))
	// 批量创建订阅
	for _, req := range reqs {
		b := x.getBroadcastFromAccountType(req.AccountType)
		sub, err := x.newOrderSubscriber(b, req.ClientOrderId, req.OrderId, req.AccountType, req.Symbol)
		if err != nil {
			return nil, err
		}
		subs = append(subs, sub)
		defers = append(defers, func() {
			x.closeSubscribe(b, sub)
		})
	}
	defer func() {
		for _, d := range defers {
			d()
		}
	}()

	// 执行API
	res, err := api.Do()
	if err != nil && res == nil {
		return nil, err
	}
	// 处理API返回值
	ords, err := x.handleOrderFromBatchOrderCreate(reqs, res)
	if err != nil {
		return ords, err
	}

	var orders []*Order
	// 批量异步接收ws结果，1秒超时
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, sub := range subs {
		wg.Add(1)
		go func(sub *xcoinOrderSubscriber) {
			defer wg.Done()
			order, err := x.waitSubscribeReturn(sub, 1*time.Second)
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

func (x *XcoinTradeEngine) AmendOrders(reqs []*OrderParam) ([]*Order, error) {
	return nil, ErrorNotSupport
}

func (x *XcoinTradeEngine) CancelOrders(reqs []*OrderParam) ([]*Order, error) {
	if err := x.restBatchPreCheck(reqs); err != nil {
		return nil, err
	}

	api := x.apiBatchOrderCancel(reqs)

	var defers []func()
	subs := make([]*xcoinOrderSubscriber, 0, len(reqs))
	// 批量创建订阅
	for _, req := range reqs {
		b := x.getBroadcastFromAccountType(req.AccountType)
		sub, err := x.newOrderSubscriber(b, req.ClientOrderId, req.OrderId, req.AccountType, req.Symbol)
		if err != nil {
			return nil, err
		}
		subs = append(subs, sub)
		defers = append(defers, func() {
			x.closeSubscribe(b, sub)
		})
	}
	defer func() {
		for _, d := range defers {
			d()
		}
	}()

	// 执行API
	res, err := api.Do()
	if err != nil && res == nil {
		return nil, err
	}
	// 处理API返回值
	ords, err := x.handleOrderFromBatchOrderCancel(reqs, res)
	if err != nil {
		return ords, err
	}

	var orders []*Order
	// 批量异步接收ws结果，1秒超时
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, sub := range subs {
		wg.Add(1)
		go func(sub *xcoinOrderSubscriber) {
			defer wg.Done()
			order, err := x.waitSubscribeReturn(sub, 1*time.Second)
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

func (x *XcoinTradeEngine) NewSubscribeOrderReq() *SubscribeOrderParam {
	return &SubscribeOrderParam{}
}

func (x *XcoinTradeEngine) SubscribeOrder(req *SubscribeOrderParam) (TradeSubscribe[Order], error) {
	return nil, ErrorNotSupport
}

func (x *XcoinTradeEngine) WsCreateOrder(req *OrderParam) (*Order, error) {
	return nil, ErrorNotSupport
}

func (x *XcoinTradeEngine) WsAmendOrder(req *OrderParam) (*Order, error) {
	return nil, ErrorNotSupport
}

func (x *XcoinTradeEngine) WsCancelOrder(req *OrderParam) (*Order, error) {
	return nil, ErrorNotSupport
}

func (x *XcoinTradeEngine) WsCreateOrders(reqs []*OrderParam) ([]*Order, error) {
	return nil, ErrorNotSupport
}

func (x *XcoinTradeEngine) WsAmendOrders(reqs []*OrderParam) ([]*Order, error) {
	return nil, ErrorNotSupport
}

func (x *XcoinTradeEngine) WsCancelOrders(reqs []*OrderParam) ([]*Order, error) {
	return nil, ErrorNotSupport
}
