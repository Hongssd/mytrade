package mytrade

import (
	"sync"
	"time"

	mybitgetapi "github.com/Hongssd/mybitgetapi"
)

type BitgetTradeEngine struct {
	ExchangeBase

	converter  BitgetEnumConverter
	apiKey     string
	secretKey  string
	passphrase string

	privateClient *mybitgetapi.PrivateRestClient
	isClassic     bool
	posModeHedge  bool
	modeDetectErr error

	broadcasterSpot           *bitgetOrderBroadcaster
	broadcasterMarginCrossed  *bitgetOrderBroadcaster
	broadcasterMarginIsolated *bitgetOrderBroadcaster
	broadcasterFutures        *bitgetOrderBroadcaster
}

func (e *BitgetTradeEngine) checkMode() error {
	if e.modeDetectErr != nil {
		return e.modeDetectErr
	}
	return nil
}

func (e *BitgetTradeEngine) NewOrderReq() *OrderParam {
	return &OrderParam{}
}

func (e *BitgetTradeEngine) NewQueryOrderReq() *QueryOrderParam {
	return &QueryOrderParam{}
}

func (e *BitgetTradeEngine) NewQueryTradeReq() *QueryTradeParam {
	return &QueryTradeParam{}
}

func (e *BitgetTradeEngine) NewSubscribeOrderReq() *SubscribeOrderParam {
	return &SubscribeOrderParam{}
}

func (e *BitgetTradeEngine) QueryOpenOrders(req *QueryOrderParam) ([]*Order, error) {
	if err := e.checkMode(); err != nil {
		return nil, err
	}
	if req == nil || req.AccountType == "" {
		return nil, ErrorInvalidParam
	}
	if err := e.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}
	if e.isClassic {
		switch req.AccountType {
		case BITGET_AC_SPOT:
			api := e.apiClassicSpotQueryOpenOrders(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			orders, err := e.handleOrdersFromClassicSpotQueryOpenOrders(req, res)
			if err != nil {
				log.Error(err)
				return nil, err
			}
			return orders, nil
		case BITGET_AC_MARGIN:
			if req.IsIsolated {
				api := e.apiClassicMarginIsolatedQueryOpenOrders(req)
				res, err := api.Do()
				if err != nil {
					return nil, err
				}
				orders, err := e.handleOrdersFromClassicMarginIsolatedQueryOpenOrders(req, res)
				if err != nil {
					log.Error(err)
					return nil, err
				}
				return orders, nil
			} else {
				api := e.apiClassicMarginCrossQueryOpenOrders(req)
				res, err := api.Do()
				if err != nil {
					return nil, err
				}
				orders, err := e.handleOrdersFromClassicMarginCrossQueryOpenOrders(req, res)
				if err != nil {
					log.Error(err)
					return nil, err
				}
				return orders, nil
			}
		case BITGET_AC_USDT_FUTURES, BITGET_AC_COIN_FUTURES, BITGET_AC_USDC_FUTURES:
			api := e.apiClassicFuturesQueryOpenOrders(req.AccountType, req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return e.handleOrdersFromClassicFuturesQueryOrders(req, res)
		default:
			return nil, ErrorAccountType
		}
	} else {
		// TODO UTA:
		return nil, ErrorNotSupport
		// res, err := e.apiUtaTradeUnfilledOrders(req).Do()
		// if err != nil {
		// 	return nil, err
		// }
		// return e.handleOrdersFromUtaTradeOrderList(e.converter, req.AccountType, res.Data.List), nil
	}
}

func (e *BitgetTradeEngine) QueryOrder(req *QueryOrderParam) (*Order, error) {
	if err := e.checkMode(); err != nil {
		return nil, err
	}
	if req == nil || req.AccountType == "" {
		return nil, ErrorInvalidParam
	}
	if err := e.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}
	if req.OrderId == "" && req.ClientOrderId == "" {
		return nil, ErrorInvalidParam
	}
	if e.isClassic {
		switch req.AccountType {
		case BITGET_AC_SPOT:
			api := e.apiClassicSpotQueryOrder(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			if len(res.Data) == 0 {
				return nil, ErrorOrderNotFound
			}
			order, err := e.handleOrderFromClassicSpotQueryOrder(req, res)
			if err != nil {
				return nil, err
			}
			return order, nil
		case BITGET_AC_MARGIN:
			return nil, ErrorNotSupport
		case BITGET_AC_USDT_FUTURES, BITGET_AC_COIN_FUTURES, BITGET_AC_USDC_FUTURES:
			if req.Symbol == "" {
				return nil, ErrorInvalidParam
			}
			api := e.apiClassicFuturesQueryOrder(req.Symbol, req.AccountType, req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			order, err := e.handleOrderFromClassicFuturesQueryOrder(req, res)
			if err != nil {
				return nil, err
			}
			return order, nil
		default:
			return nil, ErrorAccountType
		}
	} else {
		// TODO UTA:
		return nil, ErrorNotSupport
	}

	// res, err := e.apiUtaTradeOrderInfo(req).Do()
	// if err != nil {
	// 	return nil, err
	// }
	// return e.handleOrderFromUtaTradeOrderInfo(e.converter, req.AccountType, &res.Data), nil
}

func (e *BitgetTradeEngine) QueryOrders(req *QueryOrderParam) ([]*Order, error) {
	if err := e.checkMode(); err != nil {
		return nil, err
	}
	if req == nil || req.AccountType == "" {
		return nil, ErrorInvalidParam
	}
	if err := e.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}
	if e.isClassic {
		switch req.AccountType {
		case BITGET_AC_SPOT:
			api := e.apiClassicSpotQueryOrders(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return e.handleOrdersFromClassicSpotQueryOrders(req, res)
		case BITGET_AC_MARGIN:
			if req.IsIsolated {
				api := e.apiClassicMarginIsolatedQueryOrders(req)
				res, err := api.Do()
				if err != nil {
					return nil, err
				}
				return e.handleOrdersFromClassicMarginIsolatedQueryOrders(req, res)
			} else {
				api := e.apiClassicMarginCrossedQueryOrders(req)
				res, err := api.Do()
				if err != nil {
					return nil, err
				}
				return e.handleOrdersFromClassicMarginCrossedQueryOrders(req, res)
			}
		case BITGET_AC_USDT_FUTURES, BITGET_AC_COIN_FUTURES, BITGET_AC_USDC_FUTURES:
			api := e.apiClassicFuturesQueryOrders(req.AccountType, req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return e.handleOrdersFromClassicFuturesQueryOrders(req, res)
		default:
			return nil, ErrorAccountType
		}
	}
	// TODO UTA
	return nil, ErrorNotSupport
}

func (e *BitgetTradeEngine) QueryTrades(req *QueryTradeParam) ([]*Trade, error) {
	if err := e.checkMode(); err != nil {
		return nil, err
	}
	if req == nil || req.AccountType == "" {
		return nil, ErrorInvalidParam
	}
	if err := e.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}
	if e.isClassic {
		switch req.AccountType {
		case BITGET_AC_SPOT:
			api := e.apiClassicSpotQueryTrades(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return e.handleTradesFromClassicSpotQueryTrades(req, res)
		case BITGET_AC_MARGIN:
			if req.IsIsolated {
				api := e.apiClassicMarginIsolatedQueryTrades(req)
				res, err := api.Do()
				if err != nil {
					return nil, err
				}
				return e.handleTradesFromClassicMarginIsolatedQueryTrades(req, res)
			} else {
				api := e.apiClassicMarginCrossedQueryTrades(req)
				res, err := api.Do()
				if err != nil {
					return nil, err
				}
				return e.handleTradesFromClassicMarginCrossedQueryTrades(req, res)
			}
		case BITGET_AC_USDT_FUTURES, BITGET_AC_COIN_FUTURES, BITGET_AC_USDC_FUTURES:
			api := e.apiClassicFuturesQueryTrades(req.AccountType, req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return e.handleTradesFromClassicFuturesQueryTrades(req, res)
		default:
			return nil, ErrorAccountType
		}
	}
	// TODO UTA
	return nil, ErrorNotSupport
}

func (e *BitgetTradeEngine) CreateOrder(req *OrderParam) (*Order, error) {
	if err := e.checkMode(); err != nil {
		return nil, err
	}
	if req == nil || req.AccountType == "" {
		return nil, ErrorInvalidParam
	}
	if err := e.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}

	if e.isClassic {
		b := e.getBroadcastFromAccountType(req.AccountType)
		sub, err := e.newOrderSubscriber(b, req.ClientOrderId, req.AccountType, req.Symbol)
		if err != nil {
			return nil, err
		}
		defer e.closeSubscribe(b, sub)

		switch req.AccountType {
		case BITGET_AC_SPOT:
			api := e.apiClassicSpotOrderCreate(req)
			// 执行API
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			// 处理API返回值
			_, err = e.handleOrderFromClassicSpotOrderCreate(req, &res.Data)
			if err != nil {
				return nil, err
			}
		case BITGET_AC_MARGIN:
			if req.IsIsolated {
				api := e.apiClassicMarginIsolatedOrderCreate(req)
				res, err := api.Do()
				if err != nil {
					return nil, err
				}
				_, err = e.handleOrderFromClassicMarginIsolatedOrderCreate(req, &res.Data)
				if err != nil {
					return nil, err
				}
			} else {
				api := e.apiClassicMarginCrossedOrderCreate(req)
				res, err := api.Do()
				if err != nil {
					return nil, err
				}
				_, err = e.handleOrderFromClassicMarginCrossedOrderCreate(req, &res.Data)
				if err != nil {
					return nil, err
				}
			}
		case BITGET_AC_USDT_FUTURES, BITGET_AC_COIN_FUTURES, BITGET_AC_USDC_FUTURES:
			marginCode := bitgetMarginCoinFromSymbol(req.Symbol, req.Ccy)
			api := e.apiClassicFuturesOrderCreate(req, marginCode)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			_, err = e.handleOrderFromClassicFuturesOrderCreate(req, &res.Data)
			if err != nil {
				return nil, err
			}
		}

		return e.waitSubscribeReturn(sub, 1*time.Second)
	}

	return nil, ErrorNotSupport
}

func (e *BitgetTradeEngine) AmendOrder(req *OrderParam) (*Order, error) {
	if err := e.checkMode(); err != nil {
		return nil, err
	}
	if req == nil || req.AccountType == "" {
		return nil, ErrorInvalidParam
	}
	if err := e.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}
	if req.IsAlgo {
		return nil, ErrorNotSupport
	}
	if e.isClassic {
		b := e.getBroadcastFromAccountType(req.AccountType)
		sub, err := e.newOrderSubscriber(b, req.NewClientOrderId, req.AccountType, req.Symbol)
		if err != nil {
			return nil, err
		}
		defer e.closeSubscribe(b, sub)
		switch req.AccountType {
		case BITGET_AC_SPOT:
			api := e.apiClassicSpotAmendOrder(req)
			_, err = api.Do()
			if err != nil {
				return nil, err
			}
			// return e.handleOrderFromClassicSpotAmendOrder(req, &res.Data)
		case BITGET_AC_MARGIN:
			return nil, ErrorNotSupport
		case BITGET_AC_USDT_FUTURES, BITGET_AC_COIN_FUTURES, BITGET_AC_USDC_FUTURES:
			api := e.apiClassicFuturesAmendOrder(req)
			_, err := api.Do()
			if err != nil {
				return nil, err
			}
			// return e.handleOrderFromClassicFuturesAmendOrder(req, &res.Data)
		}

		return e.waitSubscribeReturn(sub, 1*time.Second)
	}

	// TODO UTA
	return nil, ErrorNotSupport
}

func (e *BitgetTradeEngine) CancelOrder(req *OrderParam) (*Order, error) {
	if err := e.checkMode(); err != nil {
		return nil, err
	}
	if req == nil || req.AccountType == "" {
		return nil, ErrorInvalidParam
	}
	if err := e.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}
	if req.IsAlgo {
		return nil, ErrorNotSupport
	}

	if e.isClassic {
		b := e.getBroadcastFromAccountType(req.AccountType)
		sub, err := e.newOrderSubscriber(b, req.NewClientOrderId, req.AccountType, req.Symbol)
		if err != nil {
			return nil, err
		}
		defer e.closeSubscribe(b, sub)

		switch req.AccountType {
		case BITGET_AC_SPOT:
			api := e.apiClassicSpotCancelOrder(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			_, err = e.handleOrderFromClassicSpotCancelOrder(req, &res.Data)
			if err != nil {
				return nil, err
			}
		case BITGET_AC_MARGIN:
			if req.IsIsolated {
				api := e.apiClassicMarginIsolatedCancelOrder(req)
				res, err := api.Do()
				if err != nil {
					return nil, err
				}
				_, err = e.handleOrderFromClassicMarginIsolatedCancelOrder(req, &res.Data)
				if err != nil {
					return nil, err
				}
			} else {
				api := e.apiClassicMarginCrossCancelOrder(req)
				res, err := api.Do()
				if err != nil {
					return nil, err
				}
				_, err = e.handleOrderFromClassicMarginCrossCancelOrder(req, &res.Data)
				if err != nil {
					return nil, err
				}
			}
		case BITGET_AC_USDT_FUTURES, BITGET_AC_COIN_FUTURES, BITGET_AC_USDC_FUTURES:
			api := e.apiClassicFuturesCancelOrder(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			_, err = e.handleOrderFromClassicFuturesCancelOrder(req, &res.Data)
			if err != nil {
				return nil, err
			}
		}
		return e.waitSubscribeReturn(sub, 1*time.Second)
	}

	// TODO UTA
	return nil, ErrorNotSupport
}

func (e *BitgetTradeEngine) CreateOrders(reqs []*OrderParam) ([]*Order, error) {
	if err := e.checkMode(); err != nil {
		return nil, err
	}
	if err := e.restBatchPreCheck(reqs); err != nil {
		return nil, err
	}
	for _, r := range reqs {
		if r.IsAlgo {
			return nil, ErrorNotSupport
		}
	}

	var defers []func()
	subs := make([]*bitgetOrderSubscriber, 0, len(reqs))
	for _, req := range reqs {
		b := e.getBroadcastFromAccountType(req.AccountType)
		sub, err := e.newOrderSubscriber(b, req.ClientOrderId, req.AccountType, req.Symbol)
		if err != nil {
			return nil, err
		}
		subs = append(subs, sub)
		br := b
		su := sub
		defers = append(defers, func() {
			e.closeSubscribe(br, su)
		})
	}
	defer func() {
		for _, d := range defers {
			d()
		}
	}()

	var ords []*Order
	var batchErr error

	if e.isClassic {
		at := reqs[0].AccountType
		for _, r := range reqs[1:] {
			if r.AccountType != at {
				return nil, ErrorInvalid("classic batch requires same AccountType")
			}
		}
		switch at {
		case BITGET_AC_SPOT:
			multiple := false
			for _, r := range reqs[1:] {
				if r.Symbol != reqs[0].Symbol {
					multiple = true
					break
				}
			}
			// 获取API
			api := e.apiClassicSpotBatchCreateOrders(reqs, multiple)
			// 执行API
			res, err := api.Do()
			if err != nil && res == nil {
				return nil, err
			}
			// 处理API返回值
			ords, batchErr = e.handleOrdersFromClassicSpotBatchOrders(reqs, &res.Data)
		case BITGET_AC_MARGIN:
			if err := bitgetClassicMarginBatchPreCheck(reqs); err != nil {
				return nil, err
			}
			if reqs[0].IsIsolated {
				api := e.apiClassicMarginIsolatedBatchCreateOrders(reqs)
				res, err := api.Do()
				if err != nil && res == nil {
					return nil, err
				}
				ords, batchErr = e.handleOrdersFromClassicMarginIsolatedBatchOrders(reqs, &res.Data)
			} else {
				api := e.apiClassicMarginCrossBatchCreateOrders(reqs)
				res, err := api.Do()
				if err != nil && res == nil {
					return nil, err
				}
				ords, batchErr = e.handleOrdersFromClassicMarginCrossBatchOrders(reqs, &res.Data)
			}
		case BITGET_AC_USDT_FUTURES, BITGET_AC_COIN_FUTURES, BITGET_AC_USDC_FUTURES:
			if err := bitgetClassicFuturesBatchPreCheck(reqs); err != nil {
				return nil, err
			}
			marginCode := bitgetMarginCoinFromSymbol(reqs[0].Symbol, reqs[0].Ccy)
			res, err := e.apiClassicFuturesBatchCreateOrders(reqs, marginCode).Do()
			if err != nil && res == nil {
				return nil, err
			}
			ords, batchErr = e.handleOrdersFromClassicFuturesBatchOrders(reqs, &res.Data)
		default:
			return nil, ErrorAccountType
		}
	} else {
		// TODO UTA
		return nil, ErrorNotSupport
	}

	if batchErr != nil {
		return ords, batchErr
	}

	var orders []*Order
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, sub := range subs {
		wg.Add(1)
		go func(sub *bitgetOrderSubscriber) {
			defer wg.Done()
			order, err := e.waitSubscribeReturn(sub, 1*time.Second)
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

func (e *BitgetTradeEngine) AmendOrders(reqs []*OrderParam) ([]*Order, error) {
	return nil, ErrorNotSupport
}

func (e *BitgetTradeEngine) CancelOrders(reqs []*OrderParam) ([]*Order, error) {
	if err := e.checkMode(); err != nil {
		return nil, err
	}
	if err := e.restBatchPreCheck(reqs); err != nil {
		return nil, err
	}
	for _, r := range reqs {
		if r.IsAlgo {
			return nil, ErrorNotSupport
		}
	}

	var defers []func()
	subs := make([]*bitgetOrderSubscriber, 0, len(reqs))
	for _, req := range reqs {
		b := e.getBroadcastFromAccountType(req.AccountType)
		sub, err := e.newOrderSubscriber(b, req.ClientOrderId, req.AccountType, req.Symbol)
		if err != nil {
			return nil, err
		}
		subs = append(subs, sub)
		br := b
		su := sub
		defers = append(defers, func() {
			e.closeSubscribe(br, su)
		})
	}
	defer func() {
		for _, d := range defers {
			d()
		}
	}()

	var ords []*Order
	var batchErr error

	if e.isClassic {
		at := reqs[0].AccountType
		for _, r := range reqs[1:] {
			if r.AccountType != at {
				return nil, ErrorInvalid("classic batch requires same AccountType")
			}
		}
		switch at {
		case BITGET_AC_SPOT:
			api := e.apiClassicSpotBatchCancelOrders(reqs)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			ords, batchErr = e.handleOrdersFromClassicSpotBatchCancelOrders(reqs, &res.Data)
		case BITGET_AC_MARGIN:
			if err := bitgetClassicMarginBatchPreCheck(reqs); err != nil {
				return nil, err
			}
			if reqs[0].IsIsolated {
				api := e.apiClassicMarginIsolatedBatchCancelOrders(reqs)
				res, err := api.Do()
				if err != nil {
					return nil, err
				}
				ords, batchErr = e.handleOrdersFromClassicMarginIsolatedBatchCancelOrders(reqs, &res.Data)
			} else {
				res, err := e.apiClassicMarginCrossBatchCancelOrders(reqs).Do()
				if err != nil {
					return nil, err
				}
				ords, batchErr = e.handleOrdersFromClassicMarginCrossBatchCancelOrders(reqs, &res.Data)
			}
		case BITGET_AC_USDT_FUTURES, BITGET_AC_COIN_FUTURES, BITGET_AC_USDC_FUTURES:
			if err := bitgetClassicFuturesBatchPreCheck(reqs); err != nil {
				return nil, err
			}
			marginCode := bitgetMarginCoinFromSymbol(reqs[0].Symbol, reqs[0].Ccy)
			res, err := e.apiClassicFuturesBatchCancelOrders(reqs, marginCode).Do()
			if err != nil {
				return nil, err
			}
			ords, batchErr = e.handleOrdersFromClassicFuturesBatchCancelOrders(reqs, &res.Data)
		default:
			return nil, ErrorAccountType
		}
	} else {
		// TODO UTA
		return nil, ErrorNotSupport
	}

	if batchErr != nil {
		return ords, batchErr
	}

	var orders []*Order
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, sub := range subs {
		wg.Add(1)
		go func(sub *bitgetOrderSubscriber) {
			defer wg.Done()
			order, err := e.waitSubscribeReturn(sub, time.Second)
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

func (e *BitgetTradeEngine) SubscribeOrder(req *SubscribeOrderParam) (TradeSubscribe[Order], error) {
	if err := e.checkMode(); err != nil {
		return nil, err
	}
	if err := e.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}
	b := e.getBroadcastFromAccountType(req.AccountType)
	sub, err := e.newOrderSubscriber(b, "", req.AccountType, "")
	if err != nil {
		return nil, err
	}
	middleSub := &subscription[Order]{
		resultChan: make(chan Order, 100),
		errChan:    make(chan error, 10),
		closeChan:  make(chan struct{}, 10),
	}
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

func (e *BitgetTradeEngine) WsCreateOrder(req *OrderParam) (*Order, error) {
	return e.CreateOrder(req)
}

func (e *BitgetTradeEngine) WsAmendOrder(req *OrderParam) (*Order, error) {
	return e.AmendOrder(req)
}

func (e *BitgetTradeEngine) WsCancelOrder(req *OrderParam) (*Order, error) {
	return e.CancelOrder(req)
}

func (e *BitgetTradeEngine) WsCreateOrders(reqs []*OrderParam) ([]*Order, error) {
	return e.CreateOrders(reqs)
}

func (e *BitgetTradeEngine) WsAmendOrders(reqs []*OrderParam) ([]*Order, error) {
	return e.AmendOrders(reqs)
}

func (e *BitgetTradeEngine) WsCancelOrders(reqs []*OrderParam) ([]*Order, error) {
	return e.CancelOrders(reqs)
}
