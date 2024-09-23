package mytrade

import (
	"github.com/Hongssd/mybybitapi"
	"sync"
	"time"
)

type BybitTradeEngine struct {
	ExchangeBase

	bybitConverter BybitEnumConverter

	broadcasterSpot    *bybitOrderBroadcaster
	broadcasterLinear  *bybitOrderBroadcaster
	broadcasterInverse *bybitOrderBroadcaster

	wsForOrder   *mybybitapi.TradeWsStreamClient
	wsForOrderMu sync.Mutex

	apiKey    string
	secretKey string
}

func (b *BybitTradeEngine) NewOrderReq() *OrderParam {
	return &OrderParam{}
}
func (b *BybitTradeEngine) NewQueryOrderReq() *QueryOrderParam {
	return &QueryOrderParam{}
}
func (b *BybitTradeEngine) NewQueryTradeReq() *QueryTradeParam {
	return &QueryTradeParam{}
}

func (b *BybitTradeEngine) QueryOpenOrders(req *QueryOrderParam) ([]*Order, error) {
	if err := b.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}

	api := b.apiQueryOpenOrders(req, "")

	res, err := api.Do()
	if err != nil {
		return nil, err
	}
	orders := b.handleOrdersFromQueryOpenOrders(req, res.Result)
	for res.Result.NextPageCursor != "" {
		api = b.apiQueryOpenOrders(req, res.Result.NextPageCursor)
		res, err = api.Do()
		if err != nil {
			return nil, err
		}
		orders = append(orders, b.handleOrdersFromQueryOpenOrders(req, res.Result)...)
	}

	return orders, nil
}
func (b *BybitTradeEngine) QueryOrder(req *QueryOrderParam) (*Order, error) {
	if err := b.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}

	api := b.apiQueryOrder(req)

	res, err := api.Do()
	if err != nil {
		return nil, err
	}
	d, _ := json.Marshal(res)
	log.Warn(string(d))
	orders := b.handleOrdersFromQueryOpenOrders(req, res.Result)
	if len(orders) != 1 {
		return nil, ErrorOrderNotFound
	}

	return orders[0], nil
}
func (b *BybitTradeEngine) QueryOrders(req *QueryOrderParam) ([]*Order, error) {
	if err := b.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}

	api := b.apiQueryOrders(req)

	res, err := api.Do()
	if err != nil {
		return nil, err
	}

	orders := b.handleOrdersFromQueryOrders(req, res.Result)
	return orders, nil
}

func (b *BybitTradeEngine) QueryTrades(req *QueryTradeParam) ([]*Trade, error) {
	if err := b.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}

	api := b.apiQueryTrades(req, "")

	res, err := api.Do()
	if err != nil {
		return nil, err
	}

	trades := b.handleTradesFromQueryTrades(req, res.Result)
	return trades, nil
}

func (b *BybitTradeEngine) CreateOrder(req *OrderParam) (*Order, error) {
	if err := b.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}

	//获取API
	api := b.apiOrderCreate(req)

	bb := b.getBroadcastFromAccountType(req.AccountType)

	//创建订阅
	sub, err := b.newOrderSubscriber(bb, req.ClientOrderId, req.AccountType, req.Symbol)
	if err != nil {
		return nil, err
	}
	defer b.closeSubscribe(bb, sub)

	//执行API
	_, err = api.Do()
	if err != nil {
		return nil, err
	}

	//异步接收ws结果，1秒超时
	return b.waitSubscribeReturn(sub, 1*time.Second)
}
func (b *BybitTradeEngine) AmendOrder(req *OrderParam) (*Order, error) {
	if err := b.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}

	//获取API
	api := b.apiOrderAmend(req)

	bb := b.getBroadcastFromAccountType(req.AccountType)

	//创建订阅
	sub, err := b.newOrderSubscriber(bb, req.ClientOrderId, req.AccountType, req.Symbol)
	if err != nil {
		return nil, err
	}
	defer b.closeSubscribe(bb, sub)

	//执行API
	_, err = api.Do()
	if err != nil {
		return nil, err
	}

	//异步接收ws结果，1秒超时
	return b.waitSubscribeReturn(sub, 1*time.Second)
}
func (b *BybitTradeEngine) CancelOrder(req *OrderParam) (*Order, error) {
	if err := b.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}

	//获取API
	api := b.apiOrderCancel(req)

	bb := b.getBroadcastFromAccountType(req.AccountType)

	//创建订阅
	sub, err := b.newOrderSubscriber(bb, req.ClientOrderId, req.AccountType, req.Symbol)
	if err != nil {
		return nil, err
	}
	defer b.closeSubscribe(bb, sub)

	//执行API
	_, err = api.Do()
	if err != nil {
		return nil, err
	}

	//异步接收ws结果，1秒超时
	return b.waitSubscribeReturn(sub, 1*time.Second)
}

func (b *BybitTradeEngine) CreateOrders(reqs []*OrderParam) ([]*Order, error) {
	if err := b.restBatchPreCheck(reqs); err != nil {
		return nil, err
	}

	var orders []*Order
	switch BybitAccountType(reqs[0].AccountType) {
	case BYBIT_AC_INVERSE:
		//BYBIT批量接口不支持反向合约,使用并发模式调用单订单
		var wg sync.WaitGroup
		var mu sync.Mutex
		for _, req := range reqs {
			req := req
			wg.Add(1)
			go func() {
				defer wg.Done()
				order, err := b.CreateOrder(req)
				if err != nil {
					mu.Lock()
					orders = append(orders, b.handleOrderFromInverseBatchErr(req, err))
					mu.Unlock()
				}
				mu.Lock()
				orders = append(orders, order)
				mu.Unlock()
			}()
		}
		wg.Wait()
	case BYBIT_AC_LINEAR, BYBIT_AC_SPOT:
		//获取API
		api := b.apiBatchOrderCreate(reqs)

		var defers []func()

		subs := make([]*bybitOrderSubscriber, 0, len(reqs))
		//批量创建订阅
		for _, req := range reqs {
			bb := b.getBroadcastFromAccountType(req.AccountType)

			sub, err := b.newOrderSubscriber(bb, req.ClientOrderId, req.AccountType, req.Symbol)
			if err != nil {
				return nil, err
			}

			subs = append(subs, sub)
			defers = append(defers, func() {
				b.closeSubscribe(bb, sub)
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
		ords, err := b.handleOrderFromBatchOrderCreate(reqs, res)
		if err != nil {
			return ords, err
		}

		//批量异步接收ws结果，1秒超时
		var wg sync.WaitGroup
		var mu sync.Mutex
		for _, sub := range subs {
			wg.Add(1)
			go func(sub *bybitOrderSubscriber) {
				defer wg.Done()
				order, err := b.waitSubscribeReturn(sub, 1*time.Second)
				if err != nil {
					log.Error(err)
				}
				mu.Lock()
				orders = append(orders, order)
				mu.Unlock()
			}(sub)
		}

		wg.Wait()
	default:
		return nil, ErrorAccountType
	}

	return orders, nil
}
func (b *BybitTradeEngine) AmendOrders(reqs []*OrderParam) ([]*Order, error) {
	if err := b.restBatchPreCheck(reqs); err != nil {
		return nil, err
	}

	var orders []*Order
	switch BybitAccountType(reqs[0].AccountType) {
	case BYBIT_AC_INVERSE:
		//BYBIT批量接口不支持反向合约,使用并发模式调用单订单
		var wg sync.WaitGroup
		var mu sync.Mutex
		for _, req := range reqs {
			req := req
			wg.Add(1)
			go func() {
				defer wg.Done()
				order, err := b.AmendOrder(req)
				if err != nil {
					mu.Lock()
					orders = append(orders, b.handleOrderFromInverseBatchErr(req, err))
					mu.Unlock()
				}
				mu.Lock()
				orders = append(orders, order)
				mu.Unlock()
			}()
		}

		wg.Wait()
	case BYBIT_AC_LINEAR, BYBIT_AC_SPOT:
		//获取API
		api := b.apiBatchOrderAmend(reqs)

		var defers []func()

		subs := make([]*bybitOrderSubscriber, 0, len(reqs))
		//批量创建订阅
		for _, req := range reqs {
			bb := b.getBroadcastFromAccountType(req.AccountType)

			sub, err := b.newOrderSubscriber(bb, req.ClientOrderId, req.AccountType, req.Symbol)
			if err != nil {
				return nil, err
			}

			subs = append(subs, sub)
			defers = append(defers, func() {
				b.closeSubscribe(bb, sub)
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
		ords, err := b.handleOrderFromBatchOrderAmend(reqs, res)
		if err != nil {
			return ords, err
		}

		//批量异步接收ws结果，1秒超时
		var wg sync.WaitGroup
		var mu sync.Mutex
		for _, sub := range subs {
			wg.Add(1)
			go func(sub *bybitOrderSubscriber) {
				defer wg.Done()
				order, err := b.waitSubscribeReturn(sub, 1*time.Second)
				if err != nil {
					log.Error(err)
				}
				mu.Lock()
				orders = append(orders, order)
				mu.Unlock()
			}(sub)

		}

		wg.Wait()
	default:
		return nil, ErrorAccountType
	}

	return orders, nil
}
func (b *BybitTradeEngine) CancelOrders(reqs []*OrderParam) ([]*Order, error) {
	if err := b.restBatchPreCheck(reqs); err != nil {
		return nil, err
	}
	var orders []*Order
	switch BybitAccountType(reqs[0].AccountType) {
	case BYBIT_AC_INVERSE:
		//BYBIT批量接口不支持反向合约,使用并发模式调用单订单
		var wg sync.WaitGroup
		var mu sync.Mutex
		for _, req := range reqs {
			req := req
			wg.Add(1)
			go func() {
				defer wg.Done()
				order, err := b.CancelOrder(req)
				if err != nil {
					mu.Lock()
					orders = append(orders, b.handleOrderFromInverseBatchErr(req, err))
					mu.Unlock()
				}
				mu.Lock()
				orders = append(orders, order)
				mu.Unlock()
			}()
		}
		wg.Wait()
	case BYBIT_AC_LINEAR, BYBIT_AC_SPOT:
		//获取API
		api := b.apiBatchOrderCancel(reqs)

		var defers []func()

		subs := make([]*bybitOrderSubscriber, 0, len(reqs))
		//批量创建订阅
		for _, req := range reqs {
			bb := b.getBroadcastFromAccountType(req.AccountType)

			sub, err := b.newOrderSubscriber(bb, req.ClientOrderId, req.AccountType, req.Symbol)
			if err != nil {
				return nil, err
			}

			subs = append(subs, sub)
			defers = append(defers, func() {
				b.closeSubscribe(bb, sub)
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
		ords, err := b.handleOrderFromBatchOrderCancel(reqs, res)
		if err != nil {
			return ords, err
		}

		//批量异步接收ws结果，1秒超时
		var wg sync.WaitGroup
		var mu sync.Mutex
		for _, sub := range subs {
			wg.Add(1)
			go func(sub *bybitOrderSubscriber) {
				defer wg.Done()
				order, err := b.waitSubscribeReturn(sub, 1*time.Second)
				if err != nil {
					log.Error(err)
				}
				mu.Lock()
				orders = append(orders, order)
				mu.Unlock()
			}(sub)

		}

		wg.Wait()
	default:
		return nil, ErrorAccountType
	}
	return orders, nil
}

func (b *BybitTradeEngine) NewSubscribeOrderReq() *SubscribeOrderParam {
	return &SubscribeOrderParam{}
}
func (b *BybitTradeEngine) SubscribeOrder(req *SubscribeOrderParam) (TradeSubscribe[Order], error) {

	if err := b.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}

	bb := b.getBroadcastFromAccountType(req.AccountType)

	sub, err := b.newOrderSubscriber(bb, "", req.AccountType, "")
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

func (b *BybitTradeEngine) WsCreateOrder(req *OrderParam) (*Order, error) {
	if err := b.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}
	if ok, err := b.wsOrderPreCheck(); !ok {
		return nil, err
	}

	switch BybitAccountType(req.AccountType) {
	case BYBIT_AC_INVERSE:
		//BYBIT Ws接口不支持反向合约，直接调用rest接口
		return b.CreateOrder(req)
	case BYBIT_AC_LINEAR, BYBIT_AC_SPOT:

		bb := b.getBroadcastFromAccountType(req.AccountType)

		//创建订阅
		sub, err := b.newOrderSubscriber(bb, req.ClientOrderId, req.AccountType, req.Symbol)
		if err != nil {
			return nil, err
		}
		defer b.closeSubscribe(bb, sub)

		_, err = b.wsForOrder.CreateOrder(b.apiOrderCreate(req))
		if err != nil {
			return nil, err
		}

		//异步接收ws结果，1秒超时
		return b.waitSubscribeReturn(sub, 1*time.Second)
	default:
		return nil, ErrorAccountType
	}
}
func (b *BybitTradeEngine) WsAmendOrder(req *OrderParam) (*Order, error) {
	if err := b.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}
	if ok, err := b.wsOrderPreCheck(); !ok {
		return nil, err
	}

	switch BybitAccountType(req.AccountType) {
	case BYBIT_AC_INVERSE:
		//BYBIT Ws接口不支持反向合约，直接调用rest接口
		return b.AmendOrder(req)
	case BYBIT_AC_LINEAR, BYBIT_AC_SPOT:
		bb := b.getBroadcastFromAccountType(req.AccountType)

		//创建订阅
		sub, err := b.newOrderSubscriber(bb, req.ClientOrderId, req.AccountType, req.Symbol)
		if err != nil {
			return nil, err
		}
		defer b.closeSubscribe(bb, sub)

		_, err = b.wsForOrder.AmendOrder(b.apiOrderAmend(req))
		if err != nil {
			return nil, err
		}

		//异步接收ws结果，1秒超时
		return b.waitSubscribeReturn(sub, 1*time.Second)
	default:
		return nil, ErrorAccountType
	}
}
func (b *BybitTradeEngine) WsCancelOrder(req *OrderParam) (*Order, error) {
	if err := b.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}
	if ok, err := b.wsOrderPreCheck(); !ok {
		return nil, err
	}

	switch BybitAccountType(req.AccountType) {
	case BYBIT_AC_INVERSE:
		//BYBIT Ws接口不支持反向合约，直接调用rest接口
		return b.CancelOrder(req)
	case BYBIT_AC_LINEAR, BYBIT_AC_SPOT:
		bb := b.getBroadcastFromAccountType(req.AccountType)

		//创建订阅
		sub, err := b.newOrderSubscriber(bb, req.ClientOrderId, req.AccountType, req.Symbol)
		if err != nil {
			return nil, err
		}
		defer b.closeSubscribe(bb, sub)

		_, err = b.wsForOrder.CancelOrder(b.apiOrderCancel(req))
		if err != nil {
			return nil, err
		}

		//异步接收ws结果，1秒超时
		return b.waitSubscribeReturn(sub, 1*time.Second)
	default:
		return nil, ErrorAccountType
	}
}

func (b *BybitTradeEngine) WsCreateOrders(reqs []*OrderParam) ([]*Order, error) {
	//bybit ws接口不支持批量接口，直接调用批量rest接口
	return b.CreateOrders(reqs)
}
func (b *BybitTradeEngine) WsAmendOrders(reqs []*OrderParam) ([]*Order, error) {
	//bybit ws接口不支持批量接口，使用并发模式调用ws单订单接口
	return b.AmendOrders(reqs)
}
func (b *BybitTradeEngine) WsCancelOrders(reqs []*OrderParam) ([]*Order, error) {
	//bybit ws接口不支持批量接口，使用并发模式调用ws单订单接口
	return b.CancelOrders(reqs)
}
