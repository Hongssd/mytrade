package mytrade

import (
	"errors"
	"fmt"
	"github.com/Hongssd/mybinanceapi"
	"sync"
)

type BinanceTradeEngine struct {
	exchangeBase

	bnConverter BinanceEnumConverter
	apiKey      string
	secretKey   string

	wsSpotAccount   *mybinanceapi.SpotWsStreamClient
	wsFutureAccount *mybinanceapi.FutureWsStreamClient
	wsSwapAccount   *mybinanceapi.SwapWsStreamClient

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
			api := b.apiSpotMarginOpenOrders(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			orders = b.handleOrdersFromSpotMarginOpenOrders(req, res)
		} else {
			api := b.apiSpotOpenOrders(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			orders = b.handleOrdersFromSpotOpenOrders(req, res)
		}
	case BN_AC_FUTURE:
		api := b.apiFutureOpenOrders(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		orders = b.handleOrdersFromFutureOpenOrders(req, res)
	case BN_AC_SWAP:
		api := b.apiSwapOpenOrders(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		orders = b.handleOrdersFromSwapOpenOrders(req, res)
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
			api := b.apiSpotMarginOrderQuery(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			order = b.handleOrderFromSpotMarginOrderQuery(req, res)
		} else {
			api := b.apiSpotOrderQuery(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			order = b.handleOrderFromSpotOrderQuery(req, res)
		}
	case BN_AC_FUTURE:
		api := b.apiFutureOrderQuery(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		order = b.handleOrderFromFutureOrderQuery(req, res)
	case BN_AC_SWAP:
		api := b.apiSwapOrderQuery(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		order = b.handleOrderFromSwapOrderQuery(req, res)
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
			api := b.apiSpotMarginOrdersQuery(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			orders = b.handleOrderFromSpotMarginOrdersQuery(req, res)
		} else {
			api := b.apiSpotOrdersQuery(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			orders = b.handleOrderFromSpotOrdersQuery(req, res)
		}
	case BN_AC_FUTURE:
		api := b.apiFutureOrdersQuery(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		orders = b.handleOrderFromFutureOrdersQuery(req, res)
	case BN_AC_SWAP:
		api := b.apiSwapOrdersQuery(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		orders = b.handleOrderFromSwapOrdersQuery(req, res)
	default:
		return nil, ErrorAccountType
	}

	return orders, nil
}
func (b *BinanceTradeEngine) QueryTrades(req *QueryTradeParam) ([]*Trade, error) {
	var trades []*Trade

	switch BinanceAccountType(req.AccountType) {
	case BN_AC_SPOT:
		api := b.apiSpotTradeQuery(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		trades = b.handleTradesFromSpotTradeQuery(req, res)
	case BN_AC_FUTURE:
		api := b.apiFutureTradeQuery(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		trades = b.handleTradesFromFutureTradeQuery(req, res)
	case BN_AC_SWAP:
		api := b.apiSwapTradeQuery(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		trades = b.handleTradesFromSwapTradeQuery(req, res)
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
			api := b.apiSpotMarginOrderCreate(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			order = b.handleOrderFromSpotMarginOrderCreate(req, res)
		} else {
			api := b.apiSpotOrderCreate(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			order = b.handleOrderFromSpotOrderCreate(req, res)
		}
	case BN_AC_FUTURE:
		api := b.apiFutureOrderCreate(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		order = b.handleOrderFromFutureOrderCreate(req, res)
	case BN_AC_SWAP:
		api := b.apiSwapOrderCreate(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		order = b.handleOrderFromSwapOrderCreate(req, res)
	default:
		return nil, ErrorAccountType
	}
	return order, nil
}
func (b *BinanceTradeEngine) AmendOrder(req *OrderParam) (*Order, error) {
	var order *Order

	switch BinanceAccountType(req.AccountType) {
	case BN_AC_SPOT:
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
		api := b.apiFutureOrderAmend(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		order = b.handleOrderFromFutureOrderAmend(req, res)
	case BN_AC_SWAP:
		api := b.apiSwapOrderAmend(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		order = b.handleOrderFromSwapOrderAmend(req, res)
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
			api := b.apiSpotMarginOrderCancel(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			order = b.handleOrderFromSpotMarginOrderCancel(req, res)
		} else {
			api := b.apiSpotOrderCancel(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			order = b.handleOrderFromSpotOrderCancel(req, res)
		}
	case BN_AC_FUTURE:
		api := b.apiFutureOrderCancel(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		order = b.handleOrderFromFutureOrderCancel(req, res)
	case BN_AC_SWAP:
		api := b.apiSwapOrderCancel(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		order = b.handleOrderFromSwapOrderCancel(req, res)
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
	binance := &mybinanceapi.MyBinance{}
	var err error
	switch BinanceAccountType(req.AccountType) {
	case BN_AC_SPOT:
		if b.wsSpotAccount == nil {
			b.wsSpotAccount, err = binance.NewSpotWsStreamClient().ConvertToAccountWs(b.apiKey, b.secretKey, mybinanceapi.SPOT_WS_TYPE)
			if err != nil {
				return nil, err
			}
			err := b.wsSpotAccount.OpenConn()
			if err != nil {
				return nil, err
			}
		}

		newPayload, err := b.wsSpotAccount.CreatePayload()
		if err != nil {
			return nil, err
		}
		//构建一个推送订单数据的中转订阅
		newSub := &subscription[Order]{
			resultChan: make(chan Order, 100),
			errChan:    make(chan error, 10),
			closeChan:  make(chan struct{}, 10),
		}
		b.handleSubscribeOrderFromSpotPayload(req, newPayload, newSub)
		return newSub, nil
	case BN_AC_FUTURE:
		if b.wsFutureAccount == nil {
			b.wsFutureAccount, err = binance.NewFutureWsStreamClient().ConvertToAccountWs(b.apiKey, b.secretKey)
			if err != nil {
				return nil, err
			}
			err := b.wsFutureAccount.OpenConn()
			if err != nil {
				return nil, err
			}
		}

		newPayload, err := b.wsFutureAccount.CreatePayload()
		if err != nil {
			return nil, err
		}

		//构建一个推送订单数据的中转订阅
		newSub := &subscription[Order]{
			resultChan: make(chan Order, 100),
			errChan:    make(chan error, 10),
			closeChan:  make(chan struct{}, 10),
		}

		b.handleSubscribeOrderFromFuturePayload(req, newPayload, newSub)

		return newSub, nil
	case BN_AC_SWAP:
		if b.wsSwapAccount == nil {
			b.wsSwapAccount, err = binance.NewSwapWsStreamClient().ConvertToAccountWs(b.apiKey, b.secretKey)
			if err != nil {
				return nil, err
			}
			err := b.wsSwapAccount.OpenConn()
			if err != nil {
				return nil, err
			}
		}

		newPayload, err := b.wsSwapAccount.CreatePayload()
		if err != nil {
			return nil, err
		}

		//构建一个推送订单数据的中转订阅
		newSub := &subscription[Order]{
			resultChan: make(chan Order, 100),
			errChan:    make(chan error, 10),
			closeChan:  make(chan struct{}, 10),
		}

		b.handleSubscribeOrderFromSwapPayload(req, newPayload, newSub)
		return newSub, nil
	default:
		return nil, ErrorAccountType
	}
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
