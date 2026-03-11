package mytrade

import (
	"strconv"

	"github.com/Hongssd/myxcoinapi"
)

func (e *XcoinTradeEngine) accountTypePreCheck(accountType string) error {
	switch XcoinAccountType(accountType) {
	case XCOIN_ACCOUNT_TYPE_SPOT, XCOIN_ACCOUNT_TYPE_LINEAR_PERPETUAL, XCOIN_ACCOUNT_TYPE_LINEAR_FUTURES:
	default:
		return ErrorAccountType
	}
	return nil
}

func (e *XcoinTradeEngine) restBatchPreCheck(reqs []*OrderParam) error {
	if len(reqs) > 20 {
		return ErrorInvalid("xcoin order param length require less than 20")
	}
	for _, req := range reqs {
		if err := e.accountTypePreCheck(req.AccountType); err != nil {
			return err
		}
	}
	return nil
}

func (e *XcoinTradeEngine) wsOrderPreCheck() (bool, error) {
	e.wsForOrderMu.Lock()
	defer e.wsForOrderMu.Unlock()
	if e.wsForOrder == nil {
		restClient := xcoin.NewRestClient(e.apiKey, e.apiSecret)
		newWs := xcoin.NewPrivateWsStreamClient(restClient)
		err := newWs.OpenConn()
		if err != nil {
			return false, err
		}
		err = newWs.Auth()
		if err != nil {
			return false, err
		}
		e.wsForOrder = newWs
	}
	return true, nil
}

func (e *XcoinTradeEngine) apiQueryOpenOrders(req *QueryOrderParam) *myxcoinapi.PrivateRestTradeOpenOrdersAPI {
	api := xcoin.NewRestClient(e.apiKey, e.apiSecret).PrivateRestClient().NewPrivateRestTradeOpenOrders()
	if req.Symbol != "" {
		api.Symbol(req.Symbol)
	}
	if req.OrderFilter != "" {
		api.OrderFilter(req.OrderFilter)
	}
	if req.AccountType != "" {
		api.BusinessType(req.AccountType)
	}

	return api
}

func (e *XcoinTradeEngine) apiQueryOrder(req *QueryOrderParam) *myxcoinapi.PrivateRestTradeOrderInfoAPI {
	api := xcoin.NewRestClient(e.apiKey, e.apiSecret).PrivateRestClient().NewPrivateRestTradeOrderInfo()
	if req.OrderFilter != "" {
		api.OrderFilter(req.OrderFilter)
	}
	if req.OrderId != "" {
		api.OrderId(req.OrderId)
	}
	if req.ClientOrderId != "" {
		api.ClientOrderId(req.ClientOrderId)
	}
	return api
}

func (e *XcoinTradeEngine) apiQueryOrders(req *QueryOrderParam) *myxcoinapi.PrivateRestTradeHistoryOrdersAPI {
	api := xcoin.NewRestClient(e.apiKey, e.apiSecret).PrivateRestClient().NewPrivateRestTradeHistoryOrders()
	if req.AccountType != "" {
		api.BusinessType(req.AccountType)
	}
	if req.Symbol != "" {
		api.Symbol(req.Symbol)
	}
	if req.OrderFilter != "" {
		api.OrderFilter(req.OrderFilter)
	}
	if req.StartTime != 0 {
		api.BeginTime(strconv.FormatInt(req.StartTime, 10))
	}
	if req.EndTime != 0 {
		api.EndTime(strconv.FormatInt(req.EndTime, 10))
	}
	if req.Limit != 0 {
		api.Limit(strconv.Itoa(req.Limit))
	}
	if req.OrderId != "" {
		api.OrderId(req.OrderId)
	}
	return api
}

func (e *XcoinTradeEngine) apiQueryTrades(req *QueryTradeParam) *myxcoinapi.PrivateRestTradeHistoryTradesAPI {
	api := xcoin.NewRestClient(e.apiKey, e.apiSecret).PrivateRestClient().NewPrivateRestTradeHistoryTrades()
	if req.AccountType != "" {
		api.BusinessType(req.AccountType)
	}
	if req.Symbol != "" {
		api.Symbol(req.Symbol)
	}
	if req.OrderType != "" {
		api.OrderType(req.OrderType)
	}
	if req.StartTime != 0 {
		api.BeginTime(strconv.FormatInt(req.StartTime, 10))
	}
	if req.EndTime != 0 {
		api.EndTime(strconv.FormatInt(req.EndTime, 10))
	}
	if req.Limit != 0 {
		api.Limit(strconv.Itoa(req.Limit))
	}
	return api
}

// 单订单接口获取
func (e *XcoinTradeEngine) apiOrderCreate(req *OrderParam) *myxcoinapi.PrivateRestTradeOrderAPI {
	api := xcoin.NewRestClient(e.apiKey, e.apiSecret).PrivateRestClient().NewPrivateRestTradeOrder()
	orderType := e.xcoinConverter.ToXcoinOrderType(req.OrderType, req.TimeInForce)

	// required
	api.Symbol(req.Symbol).
		Side(e.xcoinConverter.ToXcoinOrderSide(req.OrderSide)).
		OrderType(orderType).
		Price(req.Price.String()).
		Qty(req.Quantity.String())

	// optional
	if tif := e.xcoinConverter.ToXcoinTimeInForce(req.TimeInForce); tif != "" {
		api.TimeInForce(tif)
	}
	if req.ClientOrderId != "" {
		api.ClientOrderId(req.ClientOrderId)
	}
	if req.ReduceOnly {
		api.ReduceOnly(req.ReduceOnly)
	}
	if !req.OcoTpTriggerPx.IsZero() || !req.OcoTpOrdPx.IsZero() {
		api.TakeProfit(req.OcoTpTriggerPx.String())
		api.TpOrderType(e.xcoinConverter.ToXcoinOrderType(req.OcoTpOrdType, ""))
		api.TpLimitPrice(req.OcoTpOrdPx.String())
		api.StopLoss(req.OcoSlTriggerPx.String())
		api.SlOrderType(e.xcoinConverter.ToXcoinOrderType(req.OcoSlOrdType, ""))
		api.SlLimitPrice(req.OcoSlOrdPx.String())
	}
	return api
}

func (e *XcoinTradeEngine) apiBatchOrderCreate(reqs []*OrderParam) *myxcoinapi.PrivateRestTradeBatchOrderAPI {
	api := xcoin.NewRestClient(e.apiKey, e.apiSecret).PrivateRestClient().NewPrivateRestTradeBatchOrder()
	for _, req := range reqs {
		orderReq := &myxcoinapi.PrivateRestTradeOrderReq{
			Symbol:    GetPointer(req.Symbol),
			Side:      GetPointer(e.xcoinConverter.ToXcoinOrderSide(req.OrderSide)),
			OrderType: GetPointer(e.xcoinConverter.ToXcoinOrderType(req.OrderType, req.TimeInForce)),
			Qty:       GetPointer(req.Quantity.String()),
			Price:     GetPointer(req.Price.String()),
		}
		if tif := e.xcoinConverter.ToXcoinTimeInForce(req.TimeInForce); tif != "" {
			orderReq.TimeInForce = GetPointer(tif)
		}
		if req.ClientOrderId != "" {
			orderReq.ClientOrderId = GetPointer(req.ClientOrderId)
		}
		if req.ReduceOnly {
			orderReq.ReduceOnly = GetPointer(req.ReduceOnly)
		}
		if !req.OcoTpTriggerPx.IsZero() || !req.OcoTpOrdPx.IsZero() {
			orderReq.TpslOrder = &myxcoinapi.TpslOrder{
				TakeProfit:   GetPointer(req.OcoTpTriggerPx.String()),
				TpOrderType:  GetPointer(e.xcoinConverter.ToXcoinOrderType(req.OcoTpOrdType, "")),
				TpLimitPrice: GetPointer(req.OcoTpOrdPx.String()),
				StopLoss:     GetPointer(req.OcoSlTriggerPx.String()),
				SlOrderType:  GetPointer(e.xcoinConverter.ToXcoinOrderType(req.OcoSlOrdType, "")),
				SlLimitPrice: GetPointer(req.OcoSlOrdPx.String()),
			}
		}
		api.AddOrderReq(orderReq)
	}
	return api
}

func (e *XcoinTradeEngine) apiBatchOrderCancel(reqs []*OrderParam) *myxcoinapi.PrivateRestTradeBatchCancelOrderAPI {
	api := xcoin.NewRestClient(e.apiKey, e.apiSecret).PrivateRestClient().NewPrivateRestTradeBatchCancelOrder()
	for _, req := range reqs {
		cancelReq := &myxcoinapi.PrivateRestTradeCancelOrderReq{
			Symbol: GetPointer(req.Symbol),
		}
		if req.OrderId != "" {
			cancelReq.OrderId = GetPointer(req.OrderId)
		}
		if req.ClientOrderId != "" {
			cancelReq.ClientOrderId = GetPointer(req.ClientOrderId)
		}
		api.AddOrderReq(cancelReq)
	}
	return api
}
