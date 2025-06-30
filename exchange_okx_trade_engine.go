package mytrade

import (
	"sync"
	"time"

	"github.com/Hongssd/myokxapi"
)

type OkxTradeEngine struct {
	ExchangeBase

	okxConverter OkxEnumConverter
	apiKey       string
	secretKey    string
	passphrase   string

	broadcasterSpot   *okxOrderBroadcaster
	broadcasterMargin *okxOrderBroadcaster
	broadcasterSwap   *okxOrderBroadcaster
	broadcasterFuture *okxOrderBroadcaster

	orderAlgoBroadcasterSpot   *okxOrderAlgoBroadcaster
	orderAlgoBroadcasterMargin *okxOrderAlgoBroadcaster
	orderAlgoBroadcasterSwap   *okxOrderAlgoBroadcaster
	orderAlgoBroadcasterFuture *okxOrderAlgoBroadcaster

	wsForOrder   *myokxapi.PrivateWsStreamClient
	wsForOrderMu sync.Mutex

	wsForOrderAlgo   *myokxapi.BusinessWsStreamClient
	wsForOrderAlgoMu sync.Mutex
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
	if err := o.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}

	if req.IsAlgo {
		algoApi := o.apiQueryOpenOrderAlgo(req)
		res, err := algoApi.Do()
		if err != nil {
			return nil, err
		}
		orders, err := o.handleOrdersFromQueryOpenOrderAlgo(req, res)
		if err != nil {
			return nil, err
		}
		return orders, nil
	}
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
	if err := o.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}

	if req.IsAlgo {
		algoApi := o.apiQueryOrderAlgo(req)
		res, err := algoApi.Do()
		if err != nil {
			return nil, err
		}
		order, err := o.handleOrderFromQueryOrderAlgo(req, res)
		if err != nil {
			return nil, err
		}
		return order, nil
	}
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
func (o *OkxTradeEngine) QueryOrders(req *QueryOrderParam) ([]*Order, error) {
	if err := o.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}

	if req.IsAlgo {
		api := o.apiQueryOrdersAlgo(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		orders, err := o.handleOrdersFromQueryOrderAlgo(req, res)
		if err != nil {
			return nil, err
		}
		return orders, nil
	}

	api := o.apiQueryOrders(req)
	res, err := api.Do()
	if err != nil {
		return nil, err
	}
	orders, err := o.handleOrdersFromQueryOrderGet(req, res)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (o *OkxTradeEngine) QueryTrades(req *QueryTradeParam) ([]*Trade, error) {
	if err := o.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}
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
	if err := o.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}

	if req.IsAlgo {
		//获取API
		algoApi := o.apiOrderAlgoCreate(req)
		b := o.getOrderAlgoBroadcastFromAccountType(req.AccountType)
		//创建订阅
		sub, err := o.newOrderAlgoSubscriber(b, req.ClientOrderId, req.AccountType, req.Symbol)
		if err != nil {
			return nil, err
		}
		defer o.closeSubscribeAlgo(b, sub)
		//执行API
		res, err := algoApi.Do()
		if err != nil {
			return nil, err
		}
		//处理API返回值
		_, err = o.handleOrderFromOrderAlgoCreate(req, res)
		if err != nil {
			return nil, err
		}
		//异步接收ws结果，1秒超时
		return o.waitOrderAlgoSubscribeReturn(sub, 1*time.Second)
	}

	api := o.apiOrderCreate(req)

	b := o.getBroadcastFromAccountType(req.AccountType)
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
	//异步接收ws结果，1秒超时
	return o.waitSubscribeReturn(sub, 1*time.Second)

}
func (o *OkxTradeEngine) AmendOrder(req *OrderParam) (*Order, error) {
	if err := o.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}

	if req.IsAlgo {
		//获取API
		algoApi := o.apiOrderAlgoAmend(req)
		b := o.getOrderAlgoBroadcastFromAccountType(req.AccountType)
		//创建订阅
		sub, err := o.newOrderAlgoSubscriber(b, req.ClientOrderId, req.AccountType, req.Symbol)
		if err != nil {
			return nil, err
		}
		defer o.closeSubscribeAlgo(b, sub)
		//执行API
		res, err := algoApi.Do()
		if err != nil {
			return nil, err
		}
		//处理API返回值
		_, err = o.handleOrderFromOrderAlgoAmend(req, res)
		if err != nil {
			return nil, err
		}
		//异步接收ws结果，1秒超时
		return o.waitOrderAlgoSubscribeReturn(sub, 1*time.Second)
	}

	//获取API
	api := o.apiOrderAmend(req)
	b := o.getBroadcastFromAccountType(req.AccountType)
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
	//异步接收ws结果，1秒超时
	return o.waitSubscribeReturn(sub, 1*time.Second)
}
func (o *OkxTradeEngine) CancelOrder(req *OrderParam) (*Order, error) {
	if err := o.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}

	if req.IsAlgo {
		//获取API
		algoApi := o.apiOrderAlgoCancel(req)
		b := o.getOrderAlgoBroadcastFromAccountType(req.AccountType)
		//创建订阅
		sub, err := o.newOrderAlgoSubscriber(b, req.ClientOrderId, req.AccountType, req.Symbol)
		if err != nil {
			return nil, err
		}
		defer o.closeSubscribeAlgo(b, sub)
		//执行API
		res, err := algoApi.Do()
		if err != nil {
			return nil, err
		}
		//处理API返回值
		_, err = o.handleOrderFromOrderAlgoCancel(req, res)
		if err != nil {
			return nil, err
		}
		//异步接收ws结果，1秒超时
		return o.waitOrderAlgoSubscribeReturn(sub, 1*time.Second)
	}
	//获取API
	api := o.apiOrderCancel(req)
	b := o.getBroadcastFromAccountType(req.AccountType)

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

	//异步接收ws结果，1秒超时
	return o.waitSubscribeReturn(sub, 1*time.Second)
}

func (o *OkxTradeEngine) CreateOrders(reqs []*OrderParam) ([]*Order, error) {
	if err := o.restBatchPreCheck(reqs); err != nil {
		return nil, err
	}
	//获取API
	api := o.apiBatchOrderCreate(reqs)

	var defers []func()

	subs := make([]*okxOrderSubscriber, 0, len(reqs))
	//批量创建订阅
	for _, req := range reqs {
		b := o.getBroadcastFromAccountType(req.AccountType)

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
	if err != nil && res == nil {
		return nil, err
	}
	//处理API返回值
	ords, err := o.handleOrderFromBatchOrderCreate(reqs, res)
	if err != nil {
		return ords, err
	}

	var orders []*Order
	//批量异步接收ws结果，1秒超时
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, sub := range subs {
		wg.Add(1)
		go func(sub *okxOrderSubscriber) {
			defer wg.Done()
			order, err := o.waitSubscribeReturn(sub, 1*time.Second)
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
	if err := o.restBatchPreCheck(reqs); err != nil {
		return nil, err
	}
	//获取API
	api := o.apiBatchOrderAmend(reqs)

	var defers []func()

	subs := make([]*okxOrderSubscriber, 0, len(reqs))
	//批量创建订阅
	for _, req := range reqs {
		b := o.getBroadcastFromAccountType(req.AccountType)

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
	ords, err := o.handleOrderFromBatchOrderAmend(reqs, res)
	if err != nil {
		return ords, err
	}

	var orders []*Order

	//批量异步接收ws结果，1秒超时
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, sub := range subs {
		wg.Add(1)
		go func(sub *okxOrderSubscriber) {
			defer wg.Done()
			order, err := o.waitSubscribeReturn(sub, 1*time.Second)
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
	if err := o.restBatchPreCheck(reqs); err != nil {
		return nil, err
	}

	//获取API
	api := o.apiBatchOrderCancel(reqs)

	var defers []func()

	subs := make([]*okxOrderSubscriber, 0, len(reqs))
	//批量创建订阅
	for _, req := range reqs {
		b := o.getBroadcastFromAccountType(req.AccountType)

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
	ords, err := o.handleOrderFromBatchOrderCancel(reqs, res)
	if err != nil {
		return ords, err
	}

	var orders []*Order

	//批量异步接收ws结果，1秒超时
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, sub := range subs {
		wg.Add(1)
		go func(sub *okxOrderSubscriber) {
			defer wg.Done()
			order, err := o.waitSubscribeReturn(sub, 1*time.Second)
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
	if err := o.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}
	if ok, err := o.wsOrderPreCheck(); !ok {
		return nil, err
	}

	b := o.getBroadcastFromAccountType(req.AccountType)

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
	if err := o.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}
	if ok, err := o.wsOrderPreCheck(); !ok {
		return nil, err
	}

	b := o.getBroadcastFromAccountType(req.AccountType)

	//创建订阅
	sub, err := o.newOrderSubscriber(b, req.ClientOrderId, req.AccountType, req.Symbol)
	if err != nil {
		return nil, err
	}
	defer o.closeSubscribe(b, sub)

	res, err := o.wsForOrder.Order(o.handleWsOrderCreateFromOrderParam(req), time.Now().UnixMilli()+5000)
	if err != nil {
		return nil, err
	}

	//处理API返回值
	_, err = o.handleOrderFromWsOrderResult(req, res)
	if err != nil {
		return nil, err
	}

	//异步接收ws结果，1秒超时
	return o.waitSubscribeReturn(sub, 1*time.Second)
}
func (o *OkxTradeEngine) WsAmendOrder(req *OrderParam) (*Order, error) {
	if err := o.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}
	if ok, err := o.wsOrderPreCheck(); !ok {
		return nil, err
	}

	b := o.getBroadcastFromAccountType(req.AccountType)

	//创建订阅
	sub, err := o.newOrderSubscriber(b, req.ClientOrderId, req.AccountType, req.Symbol)
	if err != nil {
		return nil, err
	}
	defer o.closeSubscribe(b, sub)

	res, err := o.wsForOrder.AmendOrder(o.handleWsOrderAmendFromOrderParam(req), time.Now().UnixMilli()+5000)
	if err != nil {
		return nil, err
	}

	//处理API返回值
	_, err = o.handleOrderFromWsOrderResult(req, res)
	if err != nil {
		return nil, err
	}

	//异步接收ws结果，1秒超时
	return o.waitSubscribeReturn(sub, 1*time.Second)
}
func (o *OkxTradeEngine) WsCancelOrder(req *OrderParam) (*Order, error) {
	if err := o.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}
	if ok, err := o.wsOrderPreCheck(); !ok {
		return nil, err
	}

	b := o.getBroadcastFromAccountType(req.AccountType)

	//创建订阅
	sub, err := o.newOrderSubscriber(b, req.ClientOrderId, req.AccountType, req.Symbol)
	if err != nil {
		return nil, err
	}
	defer o.closeSubscribe(b, sub)

	res, err := o.wsForOrder.CancelOrder(o.handleWsOrderCancelFromOrderParam(req), time.Now().UnixMilli()+5000)
	if err != nil {
		return nil, err
	}

	//处理API返回值
	_, err = o.handleOrderFromWsOrderResult(req, res)
	if err != nil {
		return nil, err
	}

	//异步接收ws结果，1秒超时
	return o.waitSubscribeReturn(sub, 1*time.Second)
}

func (o *OkxTradeEngine) WsCreateOrders(reqs []*OrderParam) ([]*Order, error) {
	if err := o.restBatchPreCheck(reqs); err != nil {
		return nil, err
	}
	if ok, err := o.wsOrderPreCheck(); !ok {
		return nil, err
	}

	var defers []func()

	subs := make([]*okxOrderSubscriber, 0, len(reqs))
	//批量创建订阅
	for _, req := range reqs {
		b := o.getBroadcastFromAccountType(req.AccountType)

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
	res, err := o.wsForOrder.BatchOrder(o.handleBatchWsOrderCreateFromOrderParams(reqs), time.Now().UnixMilli()+5000)
	if err != nil && res == nil {
		return nil, err
	}
	//处理API返回值
	ords, err := o.handleOrdersFromWsBatchOrderResult(reqs, res)
	if err != nil {
		return ords, err
	}

	var orders []*Order
	//批量异步接收ws结果，1秒超时
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, sub := range subs {
		wg.Add(1)
		go func(sub *okxOrderSubscriber) {
			defer wg.Done()
			order, err := o.waitSubscribeReturn(sub, 1*time.Second)
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
func (o *OkxTradeEngine) WsAmendOrders(reqs []*OrderParam) ([]*Order, error) {
	if err := o.restBatchPreCheck(reqs); err != nil {
		return nil, err
	}
	if ok, err := o.wsOrderPreCheck(); !ok {
		return nil, err
	}

	var defers []func()

	subs := make([]*okxOrderSubscriber, 0, len(reqs))
	//批量创建订阅
	for _, req := range reqs {
		b := o.getBroadcastFromAccountType(req.AccountType)

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
	res, err := o.wsForOrder.BatchAmendOrder(o.handleBatchWsOrderAmendFromOrderParams(reqs), time.Now().UnixMilli()+5000)
	if err != nil {
		return nil, err
	}

	//处理API返回值
	ords, err := o.handleOrdersFromWsBatchOrderResult(reqs, res)
	if err != nil {
		return ords, err
	}

	var orders []*Order

	//批量异步接收ws结果，1秒超时
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, sub := range subs {
		wg.Add(1)
		go func(sub *okxOrderSubscriber) {
			defer wg.Done()
			order, err := o.waitSubscribeReturn(sub, 1*time.Second)
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
func (o *OkxTradeEngine) WsCancelOrders(reqs []*OrderParam) ([]*Order, error) {
	if err := o.restBatchPreCheck(reqs); err != nil {
		return nil, err
	}
	if ok, err := o.wsOrderPreCheck(); !ok {
		return nil, err
	}

	var defers []func()

	subs := make([]*okxOrderSubscriber, 0, len(reqs))
	//批量创建订阅
	for _, req := range reqs {
		b := o.getBroadcastFromAccountType(req.AccountType)

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
	res, err := o.wsForOrder.BatchCancelOrder(o.handleBatchWsOrderCancelFromOrderParams(reqs), time.Now().UnixMilli()+5000)
	if err != nil {
		return nil, err
	}

	//处理API返回值
	ords, err := o.handleOrdersFromWsBatchOrderResult(reqs, res)
	if err != nil {
		return ords, err
	}

	var orders []*Order

	//批量异步接收ws结果，1秒超时
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, sub := range subs {
		wg.Add(1)
		go func(sub *okxOrderSubscriber) {
			defer wg.Done()
			order, err := o.waitSubscribeReturn(sub, 1*time.Second)
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
