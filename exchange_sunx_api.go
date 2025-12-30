package mytrade

import (
	"strconv"

	"github.com/Hongssd/mysunxapi"
	"github.com/shopspring/decimal"
)

func (e *SunxTradeEngine) apiQueryOpenOrders(req *QueryOrderParam) *mysunxapi.PrivateRestTradeOrderOpensAPI {
	api := sunx.NewPrivateRestClient(e.accessKey, e.secretKey).NewPrivateRestTradeOrderOpens().MarginMode(SUNX_MARGIN_MODE_CROSSED)
	if req.Symbol != "" {
		api.ContractCode(req.Symbol)
	}
	if req.From != 0 {
		api.From(req.From)
	}
	if req.Direct != "" {
		api.Direct(req.Direct)
	}
	if req.Limit != 0 {
		api.Limit(req.Limit)
	}
	if req.OrderId != "" {
		api.OrderId(req.OrderId)
	}
	if req.ClientOrderId != "" {
		api.ClientOrderId(req.ClientOrderId)
	}
	return api
}

func (e *SunxTradeEngine) apiQueryOrder(req *QueryOrderParam) *mysunxapi.PrivateRestTradeOrderGetAPI {
	api := sunx.NewPrivateRestClient(e.accessKey, e.secretKey).NewPrivateRestTradeOrderGet().MarginMode(SUNX_MARGIN_MODE_CROSSED)
	if req.Symbol != "" {
		api.ContractCode(req.Symbol)
	}
	if req.OrderId != "" {
		api.OrderId(req.OrderId)
	}
	if req.ClientOrderId != "" {
		api.ClientOrderId(req.ClientOrderId)
	}
	return api
}

func (e *SunxTradeEngine) apiQueryOrders(req *QueryOrderParam) *mysunxapi.PrivateRestTradeOrderHistoryAPI {
	if req.Symbol == "" {
		log.Error("symbol is required")
		return nil
	}
	api := sunx.NewPrivateRestClient(e.accessKey, e.secretKey).NewPrivateRestTradeOrderHistory().
		ContractCode(req.Symbol).MarginMode(SUNX_MARGIN_MODE_CROSSED)
	if req.States != "" {
		api.States(req.States)
	}
	if req.StartTime != 0 {
		api.StartTime(strconv.FormatInt(req.StartTime, 10))
	}
	if req.EndTime != 0 {
		api.EndTime(strconv.FormatInt(req.EndTime, 10))
	}
	if req.Limit != 0 {
		api.Limit(req.Limit)
	}
	if req.Direct != "" {
		api.Direct(req.Direct)
	}
	return api
}

func (e *SunxTradeEngine) apiQueryTrades(req *QueryTradeParam) *mysunxapi.PrivateRestTradeOrderDetailsAPI {
	api := sunx.NewPrivateRestClient(e.accessKey, e.secretKey).NewPrivateRestTradeOrderDetails()
	if req.Symbol != "" {
		api.ContractCode(req.Symbol)
	}
	if req.OrderId != "" {
		api.OrderId(req.OrderId)
	}
	if req.StartTime != 0 {
		api.StartTime(strconv.FormatInt(req.StartTime, 10))
	}
	if req.EndTime != 0 {
		api.EndTime(strconv.FormatInt(req.EndTime, 10))
	}
	if req.From != 0 {
		api.From(req.From)
	}
	if req.Limit != 0 {
		api.Limit(req.Limit)
	}
	if req.Direct != "" {
		api.Direct(req.Direct)
	}
	return api
}

func (e *SunxTradeEngine) apiOrderCreate(req *OrderParam) *mysunxapi.PrivateRestTradeOrderPostAPI {
	api := sunx.NewPrivateRestClient(e.accessKey, e.secretKey).NewPrivateRestTradeOrderPost().MarginMode(SUNX_MARGIN_MODE_CROSSED)
	if req.Symbol != "" {
		api.ContractCode(req.Symbol)
	}
	if req.PositionSide != "" {
		api.PositionSide(e.sunxConverter.ToSunxPositionSide(req.PositionSide))
	}
	if req.OrderSide != "" {
		api.Side(e.sunxConverter.ToSunxOrderSide(req.OrderSide))
	}
	if req.OrderType != "" {
		api.Type(e.sunxConverter.ToSunxOrderType(req.OrderType))
	}
	if !req.Price.IsZero() {
		api.Price(req.Price.InexactFloat64())
	}
	if !req.Quantity.IsZero() {
		api.Volume(req.Quantity.InexactFloat64())
	}
	if req.ClientOrderId != "" {
		api.ClientOrderId(stringToInt64(req.ClientOrderId))
	}
	if req.TimeInForce != "" {
		api.TimeInForce(e.sunxConverter.ToSunxTimeInForce(req.TimeInForce))
	}
	if req.ReduceOnly {
		api.ReduceOnly(1)
	} else {
		api.ReduceOnly(0)
	}
	// 止盈止损
	if !req.OcoSlTriggerPx.IsZero() {
		api.SlTriggerPrice(req.OcoSlTriggerPx.String())
	}
	if req.OcoSlOrdType != "" {
		api.SlType(e.sunxConverter.ToSunxOrderType(req.OcoSlOrdType))
	}
	if !req.OcoSlOrdPx.IsZero() {
		api.SlOrderPrice(req.OcoSlOrdPx.String())
	}
	if !req.OcoTpTriggerPx.IsZero() {
		api.TpTriggerPrice(req.OcoTpTriggerPx.String())
	}
	if req.OcoTpOrdType != "" {
		api.TpType(e.sunxConverter.ToSunxOrderType(req.OcoTpOrdType))
	}
	if !req.OcoTpOrdPx.IsZero() {
		api.TpOrderPrice(req.OcoTpOrdPx.String())
	}

	return api
}

func (e *SunxTradeEngine) apiOrderCancel(req *OrderParam) *mysunxapi.PrivateRestTradeCancelOrderAPI {
	api := sunx.NewPrivateRestClient(e.accessKey, e.secretKey).NewPrivateRestTradeCancelOrder()
	if req.Symbol != "" {
		api.ContractCode(req.Symbol)
	}
	if req.OrderId != "" {
		api.OrderId(req.OrderId)
	}
	if req.ClientOrderId != "" {
		api.ClientOrderIds(req.ClientOrderId)
	}
	return api
}

// 批量订单
func (e *SunxTradeEngine) apiBatchOrderCreate(reqs []*OrderParam) *mysunxapi.PrivateRestTradeBatchOrdersAPI {
	api := sunx.NewPrivateRestClient(e.accessKey, e.secretKey).NewPrivateRestTradeBatchOrders()
	for _, req := range reqs {
		reduceOnly := 0
		if req.ReduceOnly {
			reduceOnly = 1
		}
		api.AddOrder(mysunxapi.PrivateRestTradeOrderPostReq{
			ContractCode:   GetPointer(req.Symbol),
			MarginMode:     GetPointer(SUNX_MARGIN_MODE_CROSSED),
			Side:           GetPointer(e.sunxConverter.ToSunxOrderSide(req.OrderSide)),
			Type:           GetPointer(e.sunxConverter.ToSunxOrderType(req.OrderType)),
			Price:          GetPointer(req.Price.InexactFloat64()),
			Volume:         GetPointer(req.Quantity.InexactFloat64()),
			ClientOrderId:  GetPointer(stringToInt64(req.ClientOrderId)),
			TimeInForce:    GetPointer(e.sunxConverter.ToSunxTimeInForce(req.TimeInForce)),
			ReduceOnly:     GetPointer(reduceOnly),
			PositionSide:   GetPointer(e.sunxConverter.ToSunxPositionSide(req.PositionSide)),
			SlTriggerPrice: GetPointer(req.OcoSlTriggerPx.String()),
			SlType:         GetPointer(e.sunxConverter.ToSunxOrderType(req.OcoSlOrdType)),
			SlOrderPrice:   GetPointer(req.OcoSlOrdPx.String()),
			TpTriggerPrice: GetPointer(req.OcoTpTriggerPx.String()),
			TpType:         GetPointer(e.sunxConverter.ToSunxOrderType(req.OcoTpOrdType)),
			TpOrderPrice:   GetPointer(req.OcoTpOrdPx.String()),
		})
	}
	return api
}

func (e *SunxTradeEngine) accountTypePreCheck(accountType string) error {
	switch SunxAccountType(accountType) {
	case SUNX_ACCOUNT_TYPE_SWAP:
	default:
		return ErrorAccountType
	}
	return nil
}

func (e *SunxTradeEngine) restBatchPreCheck(reqs []*OrderParam) error {
	if len(reqs) == 0 {
		return ErrorInvalid("sunx order param length require greater than 0")
	}
	//检测长度，OKX最多批量下20个订单
	if len(reqs) > 20 {
		return ErrorInvalid("okx order param length require less than 20")
	}

	//检测类型是否相同
	for _, req := range reqs {
		if err := e.accountTypePreCheck(req.AccountType); err != nil {
			return err
		}
	}

	return nil
}

func (e *SunxTradeEngine) apiBatchOrderCancel(reqs []*OrderParam) *mysunxapi.PrivateRestTradeCancelBatchOrdersAPI {
	api := sunx.NewPrivateRestClient(e.accessKey, e.secretKey).NewPrivateRestTradeCancelBatchOrders()
	if reqs[0].Symbol != "" {
		api.ContractCode(reqs[0].Symbol)
	}
	orderids := []string{}
	clientorderids := []string{}
	for _, req := range reqs {
		if req.OrderId != "" {
			orderids = append(orderids, req.OrderId)
		} else if req.ClientOrderId != "" {
			clientorderids = append(clientorderids, req.ClientOrderId)
		}
	}
	if len(orderids) > 0 {
		api.OrderId(orderids)
	}
	if len(clientorderids) > 0 {
		api.ClientOrderId(clientorderids)
	}
	return api
}

func (e *SunxTradeEngine) apiAmendOrderCreate(currOrder *Order, amendReq *OrderParam) (*OrderParam, *mysunxapi.PrivateRestTradeOrderPostAPI) {
	api := sunx.NewPrivateRestClient(e.accessKey, e.secretKey).NewPrivateRestTradeOrderPost().MarginMode(SUNX_MARGIN_MODE_CROSSED)
	orderParam := &OrderParam{
		Symbol:         currOrder.Symbol,
		PositionSide:   currOrder.PositionSide,
		OrderSide:      currOrder.Side,
		OrderType:      currOrder.Type,
		Price:          decimal.RequireFromString(currOrder.Price),
		Quantity:       decimal.RequireFromString(currOrder.Quantity),
		ClientOrderId:  currOrder.ClientOrderId,
		TimeInForce:    currOrder.TimeInForce,
		ReduceOnly:     currOrder.ReduceOnly,
		OcoSlTriggerPx: decimal.RequireFromString(currOrder.OcoSlTriggerPrice),
		OcoSlOrdType:   currOrder.OcoSlOrdType,
		OcoSlOrdPx:     decimal.RequireFromString(currOrder.OcoSlOrdPrice),
		OcoTpTriggerPx: decimal.RequireFromString(currOrder.OcoTpTriggerPrice),
		OcoTpOrdType:   currOrder.OcoTpOrdType,
		OcoTpOrdPx:     decimal.RequireFromString(currOrder.OcoTpOrdPrice),
	}
	api.ContractCode(orderParam.Symbol)
	if amendReq.Symbol != "" {
		orderParam.Symbol = amendReq.Symbol
		api.ContractCode(amendReq.Symbol)
	}

	api.PositionSide(e.sunxConverter.ToSunxPositionSide(orderParam.PositionSide))
	if amendReq.PositionSide != "" {
		orderParam.PositionSide = amendReq.PositionSide
		api.PositionSide(e.sunxConverter.ToSunxPositionSide(orderParam.PositionSide))
	}

	api.Side(e.sunxConverter.ToSunxOrderSide(orderParam.OrderSide))
	if amendReq.OrderSide != "" {
		orderParam.OrderSide = amendReq.OrderSide
		api.Side(e.sunxConverter.ToSunxOrderSide(orderParam.OrderSide))
	}

	api.Type(e.sunxConverter.ToSunxOrderType(orderParam.OrderType))
	if amendReq.OrderType != "" {
		orderParam.OrderType = amendReq.OrderType
		api.Type(e.sunxConverter.ToSunxOrderType(orderParam.OrderType))
	}

	api.Price(orderParam.Price.InexactFloat64())
	if !amendReq.Price.IsZero() {
		orderParam.Price = amendReq.Price
		api.Price(amendReq.Price.InexactFloat64())
	}

	api.Volume(orderParam.Quantity.InexactFloat64())
	if !amendReq.Quantity.IsZero() {
		orderParam.Quantity = amendReq.Quantity
		api.Volume(amendReq.Quantity.InexactFloat64())
	}

	api.ClientOrderId(stringToInt64(orderParam.ClientOrderId))
	if amendReq.ClientOrderId != "" {
		orderParam.ClientOrderId = amendReq.ClientOrderId
		api.ClientOrderId(stringToInt64(orderParam.ClientOrderId))
	}

	api.TimeInForce(e.sunxConverter.ToSunxTimeInForce(orderParam.TimeInForce))
	if amendReq.TimeInForce != "" {
		orderParam.TimeInForce = amendReq.TimeInForce
		api.TimeInForce(e.sunxConverter.ToSunxTimeInForce(orderParam.TimeInForce))
	}

	if amendReq.ReduceOnly {
		api.ReduceOnly(1)
	} else {
		api.ReduceOnly(0)
	}
	orderParam.ReduceOnly = amendReq.ReduceOnly

	api.SlTriggerPrice(orderParam.OcoSlTriggerPx.String())
	if !amendReq.OcoSlTriggerPx.IsZero() {
		orderParam.OcoSlTriggerPx = amendReq.OcoSlTriggerPx
		api.SlTriggerPrice(amendReq.OcoSlTriggerPx.String())
	}

	api.SlType(e.sunxConverter.ToSunxOrderType(orderParam.OcoSlOrdType))
	if amendReq.OcoSlOrdType != "" {
		orderParam.OcoSlOrdType = amendReq.OcoSlOrdType
		api.SlType(e.sunxConverter.ToSunxOrderType(orderParam.OcoSlOrdType))
	}

	api.SlOrderPrice(orderParam.OcoSlOrdPx.String())
	if !amendReq.OcoSlOrdPx.IsZero() {
		orderParam.OcoSlOrdPx = amendReq.OcoSlOrdPx
		api.SlOrderPrice(amendReq.OcoSlOrdPx.String())
	}

	api.TpTriggerPrice(orderParam.OcoTpTriggerPx.String())
	if !amendReq.OcoTpTriggerPx.IsZero() {
		orderParam.OcoTpTriggerPx = amendReq.OcoTpTriggerPx
		api.TpTriggerPrice(amendReq.OcoTpTriggerPx.String())
	}

	api.TpType(e.sunxConverter.ToSunxOrderType(orderParam.OcoTpOrdType))
	if amendReq.OcoTpOrdType != "" {
		orderParam.OcoTpOrdType = amendReq.OcoTpOrdType
		api.TpType(e.sunxConverter.ToSunxOrderType(orderParam.OcoTpOrdType))
	}

	api.TpOrderPrice(orderParam.OcoTpOrdPx.String())
	if !amendReq.OcoTpOrdPx.IsZero() {
		orderParam.OcoTpOrdPx = amendReq.OcoTpOrdPx
		api.TpOrderPrice(amendReq.OcoTpOrdPx.String())
	}

	return orderParam, api
}

func (s *SunxTradeEngine) checkWsForSwapOrder() error {
	if s.wsForSwapOrder == nil {
		s.wsForSwapOrder = sunx.NewPrivateWsStreamClient(s.accessKey, s.secretKey, mysunxapi.WsAPITypeNotification)
		err := s.wsForSwapOrder.OpenConn()
		if err != nil {
			return err
		}
	}
	return nil
}
