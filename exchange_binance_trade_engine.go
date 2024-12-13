package mytrade

import (
	"errors"
	"fmt"
	"github.com/Hongssd/mybinanceapi"
	"sync"
)

type BinanceTradeEngine struct {
	ExchangeBase

	bnConverter       BinanceEnumConverter
	apiKey            string
	secretKey         string
	isPortfolioMargin bool

	wsSpotAccount               *mybinanceapi.SpotWsStreamClient
	wsSpotMarginAccount         *mybinanceapi.SpotWsStreamClient
	wsSpotIsolatedMarginAccount MySyncMap[string, *mybinanceapi.SpotWsStreamClient]
	wsFutureAccount             *mybinanceapi.FutureWsStreamClient
	wsSwapAccount               *mybinanceapi.SwapWsStreamClient
	wsPMMarginAccount           *mybinanceapi.PMMarginStreamClient
	wsPMContractAccount         *mybinanceapi.PMContractStreamClient

	wsSpotWsApi   *mybinanceapi.SpotWsStreamClient
	wsFutureWsApi *mybinanceapi.FutureWsStreamClient
}

func (b *BinanceTradeEngine) NewOrderReq() *OrderParam {
	return &OrderParam{}
}
func (b *BinanceTradeEngine) NewQueryOrderReq() *QueryOrderParam {
	return &QueryOrderParam{}
}
func (b *BinanceTradeEngine) NewQueryTradeReq() *QueryTradeParam {
	return &QueryTradeParam{}
}

func (b *BinanceTradeEngine) QueryOpenOrders(req *QueryOrderParam) ([]*Order, error) {
	var orders []*Order
	switch BinanceAccountType(req.AccountType) {
	case BN_AC_SPOT:
		if req.IsMargin {
			if b.isPortfolioMargin && !req.IsIsolated {
				api := b.apiPortfolioMarginMarginOpenOrdersQuery(req)
				res, err := api.Do()
				if err != nil {
					return nil, err
				}
				orders = b.handlePortfolioMarginMarginOpenOrders(req, res)
			} else {
				api := b.apiSpotMarginOpenOrders(req)
				res, err := api.Do()
				if err != nil {
					return nil, err
				}
				orders = b.handleOrdersFromSpotMarginOpenOrders(req, res)
			}
		} else {
			api := b.apiSpotOpenOrders(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			orders = b.handleOrdersFromSpotOpenOrders(req, res)
		}
	case BN_AC_FUTURE:
		if b.isPortfolioMargin {
			api := b.apiPortfolioMarginUmOpenOrdersQuery(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			orders = b.handlePortfolioMarginUmOpenOrders(req, res)
		} else {
			api := b.apiFutureOpenOrders(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			orders = b.handleOrdersFromFutureOpenOrders(req, res)
		}
	case BN_AC_SWAP:
		if b.isPortfolioMargin {
			api := b.apiPortfolioMarginCmOpenOrdersQuery(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			orders = b.handlePortfolioMarginCmOpenOrders(req, res)
		} else {
			api := b.apiSwapOpenOrders(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			orders = b.handleOrdersFromSwapOpenOrders(req, res)
		}
	default:
		return nil, ErrorAccountType
	}
	return orders, nil
}
func (b *BinanceTradeEngine) QueryOrder(req *QueryOrderParam) (*Order, error) {
	var order *Order

	switch BinanceAccountType(req.AccountType) {
	case BN_AC_SPOT:
		if req.IsMargin {
			if b.isPortfolioMargin && !req.IsIsolated {
				api := b.apiPortfolioMarginMarginOrderQuery(req)
				res, err := api.Do()
				if err != nil {
					return nil, err
				}
				order = b.handlePortfolioMarginMarginOrderQuery(req, res)
			} else {
				api := b.apiSpotMarginOrderQuery(req)
				res, err := api.Do()
				if err != nil {
					return nil, err
				}
				order = b.handleOrderFromSpotMarginOrderQuery(req, res)
			}
		} else {
			api := b.apiSpotOrderQuery(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			order = b.handleOrderFromSpotOrderQuery(req, res)
		}
	case BN_AC_FUTURE:
		if b.isPortfolioMargin {
			api := b.apiPortfolioMarginUmOrderQuery(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			order = b.handlePortfolioMarginUmOrderQuery(req, res)
		} else {
			api := b.apiFutureOrderQuery(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			order = b.handleOrderFromFutureOrderQuery(req, res)
		}
	case BN_AC_SWAP:
		if b.isPortfolioMargin {
			api := b.apiPortfolioMarginCmOrderQuery(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			order = b.handlePortfolioMarginCmOrderQuery(req, res)
		} else {
			api := b.apiSwapOrderQuery(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			order = b.handleOrderFromSwapOrderQuery(req, res)
		}
	default:
		return nil, ErrorAccountType
	}

	return order, nil
}
func (b *BinanceTradeEngine) QueryOrders(req *QueryOrderParam) ([]*Order, error) {
	var orders []*Order

	switch BinanceAccountType(req.AccountType) {
	case BN_AC_SPOT:
		if req.IsMargin {
			if b.isPortfolioMargin && !req.IsIsolated {
				api := b.apiPortfolioMarginMarginOrdersQuery(req)
				res, err := api.Do()
				if err != nil {
					return nil, err
				}
				orders = b.handlePortfolioMarginMarginOrdersQuery(req, res)
			} else {
				api := b.apiSpotMarginOrdersQuery(req)
				res, err := api.Do()
				if err != nil {
					return nil, err
				}
				orders = b.handleOrderFromSpotMarginOrdersQuery(req, res)
			}
		} else {
			api := b.apiSpotOrdersQuery(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			orders = b.handleOrderFromSpotOrdersQuery(req, res)
		}
	case BN_AC_FUTURE:
		if b.isPortfolioMargin {
			api := b.apiPortfolioMarginUmOrdersQuery(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			orders = b.handlePortfolioMarginUmOrdersQuery(req, res)
		} else {
			api := b.apiFutureOrdersQuery(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			orders = b.handleOrderFromFutureOrdersQuery(req, res)
		}
	case BN_AC_SWAP:
		if b.isPortfolioMargin {
			api := b.apiPortfolioMarginCmOrdersQuery(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			orders = b.handlePortfolioMarginCmOrdersQuery(req, res)
		} else {
			api := b.apiSwapOrdersQuery(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			orders = b.handleOrderFromSwapOrdersQuery(req, res)
		}
	default:
		return nil, ErrorAccountType
	}

	return orders, nil
}
func (b *BinanceTradeEngine) QueryTrades(req *QueryTradeParam) ([]*Trade, error) {
	var trades []*Trade

	switch BinanceAccountType(req.AccountType) {
	case BN_AC_SPOT:
		if b.isPortfolioMargin && req.IsMargin && !req.IsIsolated {
			api := b.apiPortfolioMarginMarginTradesQuery(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			trades = b.handlePortfolioMarginMarginTradesQuery(req, res)
		} else {
			api := b.apiSpotTradeQuery(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			trades = b.handleTradesFromSpotTradeQuery(req, res)
		}
	case BN_AC_FUTURE:
		if b.isPortfolioMargin {
			api := b.apiPortfolioMarginUmTradesQuery(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			trades = b.handlePortfolioMarginUmTradesQuery(req, res)
		} else {
			api := b.apiFutureTradeQuery(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			trades = b.handleTradesFromFutureTradeQuery(req, res)
		}
	case BN_AC_SWAP:
		if b.isPortfolioMargin {
			api := b.apiPortfolioMarginCmTradesQuery(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			trades = b.handlePortfolioMarginCmTradesQuery(req, res)
		} else {
			api := b.apiSwapTradeQuery(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			trades = b.handleTradesFromSwapTradeQuery(req, res)
		}
	default:
		return nil, ErrorAccountType
	}

	return trades, nil
}

func (b *BinanceTradeEngine) CreateOrder(req *OrderParam) (*Order, error) {
	var order *Order
	switch BinanceAccountType(req.AccountType) {
	case BN_AC_SPOT:
		if req.IsMargin {
			if b.isPortfolioMargin && !req.IsIsolated {
				api := b.apiPortfolioMarginMarginOrderCreate(req)
				res, err := api.Do()
				if err != nil {
					return nil, err
				}
				order = b.handlePortfolioMarginMarginOrderCreate(req, res)
			} else {
				api := b.apiSpotMarginOrderCreate(req)
				res, err := api.Do()
				if err != nil {
					return nil, err
				}
				order = b.handleOrderFromSpotMarginOrderCreate(req, res)
			}
		} else {
			api := b.apiSpotOrderCreate(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			order = b.handleOrderFromSpotOrderCreate(req, res)
		}
	case BN_AC_FUTURE:
		if b.isPortfolioMargin {
			api := b.apiPortfolioMarginUmOrderCreate(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			order = b.handlePortfolioMarginUmOrderCreate(req, res)
		} else {
			api := b.apiFutureOrderCreate(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			order = b.handleOrderFromFutureOrderCreate(req, res)
		}
	case BN_AC_SWAP:
		if b.isPortfolioMargin {
			api := b.apiPortfolioMarginCmOrderCreate(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			order = b.handlePortfolioMarginCmOrderCreate(req, res)
		} else {
			api := b.apiSwapOrderCreate(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			order = b.handleOrderFromSwapOrderCreate(req, res)
		}
	default:
		return nil, ErrorAccountType
	}
	return order, nil
}
func (b *BinanceTradeEngine) AmendOrder(req *OrderParam) (*Order, error) {
	var order *Order

	switch BinanceAccountType(req.AccountType) {
	case BN_AC_SPOT:
		if req.IsMargin {
			return nil, errors.New("margin not support amend order")
		}
		api := b.apiSpotOrderAmend(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		if res.CancelResult != "SUCCESS" {
			return nil, errors.New("cancel order failed")
		}
		if res.NewOrderResult != "SUCCESS" {
			return nil, errors.New("cancel order success and amend order failed")
		}
		order = b.handleOrderFromSpotOrderAmend(req, res)
	case BN_AC_FUTURE:
		if b.isPortfolioMargin {
			api := b.apiPortfolioMarginUmOrderAmend(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			order = b.handlePortfolioMarginUmOrderAmend(req, res)
		} else {
			api := b.apiFutureOrderAmend(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			order = b.handleOrderFromFutureOrderAmend(req, res)
		}
	case BN_AC_SWAP:
		if b.isPortfolioMargin {
			api := b.apiPortfolioMarginCmOrderAmend(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			order = b.handlePortfolioMarginCmOrderAmend(req, res)
		} else {
			api := b.apiSwapOrderAmend(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			order = b.handleOrderFromSwapOrderAmend(req, res)
		}
	default:
		return nil, ErrorAccountType
	}

	return order, nil
}
func (b *BinanceTradeEngine) CancelOrder(req *OrderParam) (*Order, error) {
	var order *Order
	switch BinanceAccountType(req.AccountType) {
	case BN_AC_SPOT:
		if req.IsMargin {
			if b.isPortfolioMargin && !req.IsIsolated {
				api := b.apiPortfolioMarginMarginOrderCancel(req)
				res, err := api.Do()
				if err != nil {
					return nil, err
				}
				order = b.handlePortfolioMarginMarginOrderCancel(req, res)
			} else {
				api := b.apiSpotMarginOrderCancel(req)
				res, err := api.Do()
				if err != nil {
					return nil, err
				}
				order = b.handleOrderFromSpotMarginOrderCancel(req, res)
			}
		} else {
			api := b.apiSpotOrderCancel(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			order = b.handleOrderFromSpotOrderCancel(req, res)
		}
	case BN_AC_FUTURE:
		if b.isPortfolioMargin {
			api := b.apiPortfolioMarginUmOrderCancel(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			order = b.handlePortfolioMarginUmOrderCancel(req, res)
		} else {
			api := b.apiFutureOrderCancel(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			order = b.handleOrderFromFutureOrderCancel(req, res)
		}
	case BN_AC_SWAP:
		if b.isPortfolioMargin {
			api := b.apiPortfolioMarginCmOrderCancel(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			order = b.handlePortfolioMarginCmOrderCancel(req, res)
		} else {
			api := b.apiSwapOrderCancel(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			order = b.handleOrderFromSwapOrderCancel(req, res)
		}
	default:
		return nil, ErrorAccountType
	}

	return order, nil
}

func (b *BinanceTradeEngine) CreateOrders(reqs []*OrderParam) ([]*Order, error) {
	var orders []*Order
	err := b.restBatchPreCheck(reqs)
	if err != nil {
		return nil, err
	}
	switch BinanceAccountType(reqs[0].AccountType) {
	case BN_AC_SPOT:
		//现货无批量接口，直接并发下单
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
					orders = append(orders, b.handleOrderFromSpotBatchErr(req, err))
					mu.Unlock()
				}
				mu.Lock()
				orders = append(orders, order)
				mu.Unlock()
			}()
		}
		wg.Wait()
	case BN_AC_FUTURE:
		api := b.apiFutureBatchOrderCreate(reqs)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		orders = b.handleOrdersFromFutureBatchOrderCreate(reqs, res)
	case BN_AC_SWAP:
		api := b.apiSwapBatchOrderCreate(reqs)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		orders = b.handleOrdersFromSwapBatchOrderCreate(reqs, res)
	default:
		return nil, ErrorAccountType
	}

	return orders, nil
}
func (b *BinanceTradeEngine) AmendOrders(reqs []*OrderParam) ([]*Order, error) {
	var orders []*Order
	err := b.restBatchPreCheck(reqs)
	if err != nil {
		return nil, err
	}

	switch BinanceAccountType(reqs[0].AccountType) {
	case BN_AC_SPOT:
		//现货无批量接口，直接并发改单
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
					orders = append(orders, b.handleOrderFromSpotBatchErr(req, err))
					mu.Unlock()
				}
				mu.Lock()
				orders = append(orders, order)
				mu.Unlock()
			}()
		}
		wg.Wait()
	case BN_AC_FUTURE:
		api := b.apiFutureBatchOrderAmend(reqs)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		orders = b.handleOrdersFromFutureBatchOrderAmend(reqs, res)
	case BN_AC_SWAP:
		api := b.apiSwapBatchOrderAmend(reqs)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		orders = b.handleOrdersFromSwapBatchOrderAmend(reqs, res)
	default:
		return nil, ErrorAccountType
	}

	return orders, nil
}
func (b *BinanceTradeEngine) CancelOrders(reqs []*OrderParam) ([]*Order, error) {
	var orders []*Order

	err := b.restBatchPreCheck(reqs)
	if err != nil {
		return nil, err
	}

	switch BinanceAccountType(reqs[0].AccountType) {
	case BN_AC_SPOT:
		//现货无批量接口，直接并发撤单
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
					orders = append(orders, b.handleOrderFromSpotBatchErr(req, err))
					mu.Unlock()
				}
				mu.Lock()
				orders = append(orders, order)
				mu.Unlock()
			}()
		}
		wg.Wait()
	case BN_AC_FUTURE:
		api, err := b.apiFutureBatchOrderCancel(reqs)
		if err != nil {
			return nil, err
		}
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		orders = b.handleOrdersFromFutureBatchOrderCancel(reqs, res)
	case BN_AC_SWAP:
		api, err := b.apiSwapBatchOrderCancel(reqs)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		orders = b.handleOrdersFromSwapBatchOrderCancel(reqs, res)
	default:
		return nil, ErrorAccountType
	}

	return orders, nil
}

func (b *BinanceTradeEngine) NewSubscribeOrderReq() *SubscribeOrderParam {
	return &SubscribeOrderParam{}
}
func (b *BinanceTradeEngine) SubscribeOrder(r *SubscribeOrderParam) (TradeSubscribe[Order], error) {
	req := *r
	var err error
	//构建一个推送订单数据的中转订阅
	newSub := &subscription[Order]{
		resultChan: make(chan Order, 100),
		errChan:    make(chan error, 10),
		closeChan:  make(chan struct{}, 10),
	}
	switch BinanceAccountType(req.AccountType) {
	case BN_AC_SPOT:
		var targetWs *mybinanceapi.SpotWsStreamClient
		if r.IsMargin {
			if !r.IsIsolated {
				//全仓杠杆
				if !b.isPortfolioMargin {
					err = b.checkWsSpotMarginAccount()
					if err != nil {
						return nil, err
					}
					targetWs = b.wsSpotMarginAccount
				} else {
					err = b.checkWsPmMarginAccount()
					if err != nil {
						return nil, err
					}
					newPayload, err := b.wsPMMarginAccount.CreatePayload()
					if err != nil {
						return nil, err
					}
					b.handleSubscribeOrderFromPMMarginPayload(req, newPayload, newSub)
					return newSub, nil
				}
			} else {
				//逐仓杠杆
				err = b.checkWsSpotIsolatedMarginAccount(r.IsolatedSymbol)
				if err != nil {
					return nil, err
				}
				var ok bool
				targetWs, ok = b.wsSpotIsolatedMarginAccount.Load(r.IsolatedSymbol)
				if !ok {
					return nil, fmt.Errorf("isolated symbol %s not found", r.IsolatedSymbol)
				}
			}
		} else {
			//现货
			err = b.checkWsSpotAccount()
			if err != nil {
				return nil, err
			}
			targetWs = b.wsSpotAccount

		}
		newPayload, err := targetWs.CreatePayload()
		if err != nil {
			return nil, err
		}
		b.handleSubscribeOrderFromSpotPayload(req, newPayload, newSub)
		return newSub, nil
	case BN_AC_FUTURE:
		if !b.isPortfolioMargin {
			err := b.checkWsFutureAccount()
			if err != nil {
				return nil, err
			}
			newPayload, err := b.wsFutureAccount.CreatePayload()
			if err != nil {
				return nil, err
			}
			b.handleSubscribeOrderFromFuturePayload(req, newPayload, newSub)
		}
		return newSub, nil
	case BN_AC_SWAP:
		if !b.isPortfolioMargin {

			err := b.checkWsSwapAccount()
			if err != nil {
				return nil, err
			}
			newPayload, err := b.wsSwapAccount.CreatePayload()
			if err != nil {
				return nil, err
			}

			b.handleSubscribeOrderFromSwapPayload(req, newPayload, newSub)
			return newSub, nil
		}
	default:
		return nil, ErrorAccountType
	}

	if b.isPortfolioMargin {
		switch BinanceAccountType(req.AccountType) {
		case BN_AC_FUTURE, BN_AC_SWAP:
			err := b.checkWsPmContractAccount()
			if err != nil {
				return nil, err
			}
			newPayload, err := b.wsPMContractAccount.CreatePayload()
			if err != nil {
				return nil, err
			}
			b.handleSubscribeOrderFromPMContractPayload(req, newPayload, newSub)
		}
	}
	return newSub, nil
}

func (b *BinanceTradeEngine) WsCreateOrder(req *OrderParam) (*Order, error) {
	var order *Order
	var err error
	switch BinanceAccountType(req.AccountType) {
	case BN_AC_SPOT:
		if b.wsSpotWsApi == nil {
			wsSpotWsApi, err := binance.NewSpotWsStreamClient().ConvertToWsApi(b.apiKey, b.secretKey)
			if err != nil {
				return nil, err
			}
			b.wsSpotWsApi = wsSpotWsApi
			err = b.wsSpotWsApi.OpenConn()
			if err != nil {
				return nil, err
			}
		}

		res, err := b.wsSpotWsApi.CreateOrder(b.apiSpotOrderCreate(req))
		if err != nil {
			return nil, err
		}
		if res.Error.Msg != "" {
			return nil, fmt.Errorf("[%d]%s", res.Error.Code, res.Error.Msg)
		}
		order = b.handleOrderFromSpotOrderCreate(req, &res.Result)
	case BN_AC_FUTURE:
		if b.wsFutureWsApi == nil {
			b.wsFutureWsApi, err = binance.NewFutureWsStreamClient().ConvertToWsApi(b.apiKey, b.secretKey)
			if err != nil {
				return nil, err
			}
			err := b.wsFutureWsApi.OpenConn()
			if err != nil {
				return nil, err
			}
		}

		res, err := b.wsFutureWsApi.CreateOrder(b.apiFutureOrderCreate(req))
		if err != nil {
			return nil, err
		}
		if res.Error.Msg != "" {
			return nil, fmt.Errorf("[%d]%s", res.Error.Code, res.Error.Msg)
		}
		order = b.handleOrderFromFutureOrderCreate(req, &res.Result)
	case BN_AC_SWAP:
		//币合约无WS API接口，直接调用REST
		return b.CreateOrder(req)
	default:
		return nil, ErrorAccountType
	}

	return order, nil
}
func (b *BinanceTradeEngine) WsAmendOrder(req *OrderParam) (*Order, error) {
	var order *Order

	var err error
	switch BinanceAccountType(req.AccountType) {
	case BN_AC_SPOT:
		if b.wsSpotWsApi == nil {
			b.wsSpotWsApi, err = binance.NewSpotWsStreamClient().ConvertToWsApi(b.apiKey, b.secretKey)
			if err != nil {
				return nil, err
			}
			err := b.wsSpotWsApi.OpenConn()
			if err != nil {
				return nil, err
			}
		}

		res, err := b.wsSpotWsApi.CancelReplaceOrder(b.apiSpotOrderAmend(req))
		if err != nil {
			return nil, err
		}
		if res.Error.Msg != "" {
			return nil, fmt.Errorf("[%d]%s", res.Error.Code, res.Error.Msg)
		}
		order = b.handleOrderFromSpotOrderAmend(req, &res.Result)
	case BN_AC_FUTURE:
		if b.wsFutureWsApi == nil {
			b.wsFutureWsApi, err = binance.NewFutureWsStreamClient().ConvertToWsApi(b.apiKey, b.secretKey)
			if err != nil {
				return nil, err
			}
			err := b.wsFutureWsApi.OpenConn()
			if err != nil {
				return nil, err
			}
		}

		res, err := b.wsFutureWsApi.AmendOrder(b.apiFutureOrderAmend(req))
		if err != nil {
			return nil, err
		}

		if res.Error.Msg != "" {
			return nil, fmt.Errorf("[%d]%s", res.Error.Code, res.Error.Msg)
		}
		order = b.handleOrderFromFutureOrderAmend(req, &res.Result)
	case BN_AC_SWAP:
		//币合约无WS API接口，直接调用REST
		return b.AmendOrder(req)
	default:
		return nil, ErrorAccountType
	}
	return order, nil
}
func (b *BinanceTradeEngine) WsCancelOrder(req *OrderParam) (*Order, error) {
	var order *Order
	var err error
	switch BinanceAccountType(req.AccountType) {
	case BN_AC_SPOT:
		if b.wsSpotWsApi == nil {
			b.wsSpotWsApi, err = binance.NewSpotWsStreamClient().ConvertToWsApi(b.apiKey, b.secretKey)
			if err != nil {
				return nil, err
			}
			err := b.wsSpotWsApi.OpenConn()
			if err != nil {
				return nil, err
			}
		}

		res, err := b.wsSpotWsApi.CancelOrder(b.apiSpotOrderCancel(req))
		if err != nil {
			return nil, err
		}
		if res.Error.Msg != "" {
			return nil, fmt.Errorf("[%d]%s", res.Error.Code, res.Error.Msg)
		}
		order = b.handleOrderFromSpotOrderCancel(req, &res.Result)
	case BN_AC_FUTURE:
		if b.wsFutureWsApi == nil {
			b.wsFutureWsApi, err = binance.NewFutureWsStreamClient().ConvertToWsApi(b.apiKey, b.secretKey)
			if err != nil {
				return nil, err
			}
			err := b.wsFutureWsApi.OpenConn()
			if err != nil {
				return nil, err
			}
		}

		res, err := b.wsFutureWsApi.CancelOrder(b.apiFutureOrderCancel(req))
		if err != nil {
			return nil, err
		}
		if res.Error.Msg != "" {
			return nil, fmt.Errorf("[%d]%s", res.Error.Code, res.Error.Msg)
		}
		order = b.handleOrderFromFutureOrderCancel(req, &res.Result)
	case BN_AC_SWAP:
		//币合约无WS API接口，直接调用REST
		return b.CancelOrder(req)
	default:
		return nil, ErrorAccountType
	}
	return order, nil
}

func (b *BinanceTradeEngine) WsCreateOrders(reqs []*OrderParam) ([]*Order, error) {

	err := b.restBatchPreCheck(reqs)
	if err != nil {
		return nil, err
	}

	switch BinanceAccountType(reqs[0].AccountType) {
	case BN_AC_SPOT, BN_AC_FUTURE, BN_AC_SWAP:
		//合约WS无批量接口库，直接调用REST批量接口

		return b.CreateOrders(reqs)
	default:
		return nil, ErrorAccountType
	}

}
func (b *BinanceTradeEngine) WsAmendOrders(reqs []*OrderParam) ([]*Order, error) {
	err := b.restBatchPreCheck(reqs)
	if err != nil {
		return nil, err
	}

	switch BinanceAccountType(reqs[0].AccountType) {
	case BN_AC_SPOT, BN_AC_FUTURE, BN_AC_SWAP:
		//合约WS无批量接口库，直接调用REST批量接口
		return b.AmendOrders(reqs)
	default:
		return nil, ErrorAccountType
	}
}
func (b *BinanceTradeEngine) WsCancelOrders(reqs []*OrderParam) ([]*Order, error) {
	err := b.restBatchPreCheck(reqs)
	if err != nil {
		return nil, err
	}

	switch BinanceAccountType(reqs[0].AccountType) {
	case BN_AC_SPOT, BN_AC_FUTURE, BN_AC_SWAP:
		//合约WS无批量接口库，直接调用REST批量接口

		return b.CancelOrders(reqs)
	default:
		return nil, ErrorAccountType
	}
}
