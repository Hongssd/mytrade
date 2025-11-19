package mytrade

import (
	"errors"
	"fmt"
	"sync"

	"github.com/Hongssd/myasterapi"
)

type AsterTradeEngine struct {
	ExchangeBase

	asterConverter AsterEnumConverter
	apiKey         string
	secretKey      string

	wsSpotAccount   *myasterapi.SpotWsStreamClient
	wsFutureAccount *myasterapi.FutureWsStreamClient

	wsSpotWsApi   *myasterapi.SpotWsStreamClient
	wsFutureWsApi *myasterapi.FutureWsStreamClient
}

func (b *AsterTradeEngine) NewOrderReq() *OrderParam {
	return &OrderParam{}
}
func (b *AsterTradeEngine) NewQueryOrderReq() *QueryOrderParam {
	return &QueryOrderParam{}
}
func (b *AsterTradeEngine) NewQueryTradeReq() *QueryTradeParam {
	return &QueryTradeParam{}
}

func (b *AsterTradeEngine) QueryOpenOrders(req *QueryOrderParam) ([]*Order, error) {
	var orders []*Order
	switch AsterAccountType(req.AccountType) {
	case ASTER_AC_SPOT:

		api := b.apiSpotOpenOrders(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		orders = b.handleOrdersFromSpotOpenOrders(req, res)
	case ASTER_AC_FUTURE:

		api := b.apiFutureOpenOrders(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		orders = b.handleOrdersFromFutureOpenOrders(req, res)
	default:
		return nil, ErrorAccountType
	}
	return orders, nil
}
func (b *AsterTradeEngine) QueryOrder(req *QueryOrderParam) (*Order, error) {
	var order *Order

	switch AsterAccountType(req.AccountType) {
	case ASTER_AC_SPOT:
		api := b.apiSpotOrderQuery(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		order = b.handleOrderFromSpotOrderQuery(req, res)
	case ASTER_AC_FUTURE:
		api := b.apiFutureOrderQuery(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		order = b.handleOrderFromFutureOrderQuery(req, res)

	default:
		return nil, ErrorAccountType
	}

	return order, nil
}
func (b *AsterTradeEngine) QueryOrders(req *QueryOrderParam) ([]*Order, error) {
	var orders []*Order

	switch AsterAccountType(req.AccountType) {
	case ASTER_AC_SPOT:
		api := b.apiSpotOrdersQuery(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		orders = b.handleOrderFromSpotOrdersQuery(req, res)
	case ASTER_AC_FUTURE:
		api := b.apiFutureOrdersQuery(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		orders = b.handleOrderFromFutureOrdersQuery(req, res)
	default:
		return nil, ErrorAccountType
	}

	return orders, nil
}
func (b *AsterTradeEngine) QueryTrades(req *QueryTradeParam) ([]*Trade, error) {
	var trades []*Trade

	switch AsterAccountType(req.AccountType) {
	case ASTER_AC_SPOT:

		api := b.apiSpotTradeQuery(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		trades = b.handleTradesFromSpotTradeQuery(req, res)

	case ASTER_AC_FUTURE:

		api := b.apiFutureTradeQuery(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		trades = b.handleTradesFromFutureTradeQuery(req, res)

	default:
		return nil, ErrorAccountType
	}

	return trades, nil
}

func (b *AsterTradeEngine) CreateOrder(req *OrderParam) (*Order, error) {
	var order *Order
	switch AsterAccountType(req.AccountType) {
	case ASTER_AC_SPOT:

		api := b.apiSpotOrderCreate(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		order = b.handleOrderFromSpotOrderCreate(req, res)

	case ASTER_AC_FUTURE:

		api := b.apiFutureOrderCreate(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		order = b.handleOrderFromFutureOrderCreate(req, res)

	default:
		return nil, ErrorAccountType
	}
	return order, nil
}
func (b *AsterTradeEngine) AmendOrder(req *OrderParam) (*Order, error) {
	var order *Order

	switch AsterAccountType(req.AccountType) {
	case ASTER_AC_SPOT:
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
	case ASTER_AC_FUTURE:

		api := b.apiFutureOrderAmend(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		order = b.handleOrderFromFutureOrderAmend(req, res)

	default:
		return nil, ErrorAccountType
	}

	return order, nil
}
func (b *AsterTradeEngine) CancelOrder(req *OrderParam) (*Order, error) {
	var order *Order
	switch AsterAccountType(req.AccountType) {
	case ASTER_AC_SPOT:

		api := b.apiSpotOrderCancel(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		order = b.handleOrderFromSpotOrderCancel(req, res)

	case ASTER_AC_FUTURE:

		api := b.apiFutureOrderCancel(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		order = b.handleOrderFromFutureOrderCancel(req, res)

	default:
		return nil, ErrorAccountType
	}

	return order, nil
}

func (b *AsterTradeEngine) CreateOrders(reqs []*OrderParam) ([]*Order, error) {
	var orders []*Order
	err := b.restBatchPreCheck(reqs)
	if err != nil {
		return nil, err
	}
	switch AsterAccountType(reqs[0].AccountType) {
	case ASTER_AC_SPOT:
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
	case ASTER_AC_FUTURE:

		api := b.apiFutureBatchOrderCreate(reqs)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		orders = b.handleOrdersFromFutureBatchOrderCreate(reqs, res)

	default:
		return nil, ErrorAccountType
	}

	return orders, nil
}
func (b *AsterTradeEngine) AmendOrders(reqs []*OrderParam) ([]*Order, error) {
	var orders []*Order
	err := b.restBatchPreCheck(reqs)
	if err != nil {
		return nil, err
	}

	switch AsterAccountType(reqs[0].AccountType) {
	case ASTER_AC_SPOT:
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
	case ASTER_AC_FUTURE:

		api := b.apiFutureBatchOrderAmend(reqs)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		orders = b.handleOrdersFromFutureBatchOrderAmend(reqs, res)

	default:
		return nil, ErrorAccountType
	}

	return orders, nil
}
func (b *AsterTradeEngine) CancelOrders(reqs []*OrderParam) ([]*Order, error) {
	var orders []*Order

	err := b.restBatchPreCheck(reqs)
	if err != nil {
		return nil, err
	}

	switch AsterAccountType(reqs[0].AccountType) {
	case ASTER_AC_SPOT:
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
	case ASTER_AC_FUTURE:

		api, err := b.apiFutureBatchOrderCancel(reqs)
		if err != nil {
			return nil, err
		}
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		orders = b.handleOrdersFromFutureBatchOrderCancel(reqs, res)

	default:
		return nil, ErrorAccountType
	}

	return orders, nil
}

func (b *AsterTradeEngine) NewSubscribeOrderReq() *SubscribeOrderParam {
	return &SubscribeOrderParam{}
}
func (b *AsterTradeEngine) SubscribeOrder(r *SubscribeOrderParam) (TradeSubscribe[Order], error) {
	req := *r
	var err error
	//构建一个推送订单数据的中转订阅
	newSub := &subscription[Order]{
		resultChan: make(chan Order, 100),
		errChan:    make(chan error, 10),
		closeChan:  make(chan struct{}, 10),
	}
	switch AsterAccountType(req.AccountType) {
	case ASTER_AC_SPOT:
		var targetWs *myasterapi.SpotWsStreamClient
		//现货
		err = b.checkWsSpotAccount()
		if err != nil {
			return nil, err
		}
		targetWs = b.wsSpotAccount

		newPayload, err := targetWs.CreatePayload()
		if err != nil {
			return nil, err
		}
		b.handleSubscribeOrderFromSpotPayload(req, newPayload, newSub)
		return newSub, nil
	case ASTER_AC_FUTURE:
		err := b.checkWsFutureAccount()
		if err != nil {
			return nil, err
		}
		newPayload, err := b.wsFutureAccount.CreatePayload()
		if err != nil {
			return nil, err
		}
		b.handleSubscribeOrderFromFuturePayload(req, newPayload, newSub)
		return newSub, nil
	default:
		return nil, ErrorAccountType
	}
}

func (b *AsterTradeEngine) WsCreateOrder(req *OrderParam) (*Order, error) {
	var order *Order
	var err error
	switch AsterAccountType(req.AccountType) {
	case ASTER_AC_SPOT:
		if b.wsSpotWsApi == nil {
			wsSpotWsApi, err := aster.NewSpotWsStreamClient().ConvertToWsApi(b.apiKey, b.secretKey)
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
	case ASTER_AC_FUTURE:
		if b.wsFutureWsApi == nil {
			b.wsFutureWsApi, err = aster.NewFutureWsStreamClient().ConvertToWsApi(b.apiKey, b.secretKey)
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
	default:
		return nil, ErrorAccountType
	}

	return order, nil
}
func (b *AsterTradeEngine) WsAmendOrder(req *OrderParam) (*Order, error) {
	var order *Order

	var err error
	switch AsterAccountType(req.AccountType) {
	case ASTER_AC_SPOT:
		if b.wsSpotWsApi == nil {
			b.wsSpotWsApi, err = aster.NewSpotWsStreamClient().ConvertToWsApi(b.apiKey, b.secretKey)
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
	case ASTER_AC_FUTURE:
		if b.wsFutureWsApi == nil {
			b.wsFutureWsApi, err = aster.NewFutureWsStreamClient().ConvertToWsApi(b.apiKey, b.secretKey)
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
	default:
		return nil, ErrorAccountType
	}
	return order, nil
}
func (b *AsterTradeEngine) WsCancelOrder(req *OrderParam) (*Order, error) {
	var order *Order
	var err error
	switch AsterAccountType(req.AccountType) {
	case ASTER_AC_SPOT:
		if b.wsSpotWsApi == nil {
			b.wsSpotWsApi, err = aster.NewSpotWsStreamClient().ConvertToWsApi(b.apiKey, b.secretKey)
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
	case ASTER_AC_FUTURE:
		if b.wsFutureWsApi == nil {
			b.wsFutureWsApi, err = aster.NewFutureWsStreamClient().ConvertToWsApi(b.apiKey, b.secretKey)
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
	default:
		return nil, ErrorAccountType
	}
	return order, nil
}

func (b *AsterTradeEngine) WsCreateOrders(reqs []*OrderParam) ([]*Order, error) {

	err := b.restBatchPreCheck(reqs)
	if err != nil {
		return nil, err
	}

	switch AsterAccountType(reqs[0].AccountType) {
	case ASTER_AC_SPOT, ASTER_AC_FUTURE:
		//合约WS无批量接口库，直接调用REST批量接口

		return b.CreateOrders(reqs)
	default:
		return nil, ErrorAccountType
	}

}
func (b *AsterTradeEngine) WsAmendOrders(reqs []*OrderParam) ([]*Order, error) {
	err := b.restBatchPreCheck(reqs)
	if err != nil {
		return nil, err
	}

	switch AsterAccountType(reqs[0].AccountType) {
	case ASTER_AC_SPOT, ASTER_AC_FUTURE:
		//合约WS无批量接口库，直接调用REST批量接口
		return b.AmendOrders(reqs)
	default:
		return nil, ErrorAccountType
	}
}
func (b *AsterTradeEngine) WsCancelOrders(reqs []*OrderParam) ([]*Order, error) {
	err := b.restBatchPreCheck(reqs)
	if err != nil {
		return nil, err
	}

	switch AsterAccountType(reqs[0].AccountType) {
	case ASTER_AC_SPOT, ASTER_AC_FUTURE:
		//合约WS无批量接口库，直接调用REST批量接口

		return b.CancelOrders(reqs)
	default:
		return nil, ErrorAccountType
	}
}
