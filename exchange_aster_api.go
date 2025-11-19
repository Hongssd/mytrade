package mytrade

import (
	"strconv"

	"github.com/Hongssd/myasterapi"
)

// 现货订单API接口
func (b *AsterTradeEngine) apiSpotOpenOrders(req *QueryOrderParam) *myasterapi.SpotOpenOrdersApi {
	api := aster.NewSpotRestClient(b.apiKey, b.secretKey).NewOpenOrders()
	if req.Symbol != "" {
		api.Symbol(req.Symbol)
	}
	return api
}
func (b *AsterTradeEngine) apiSpotOrderQuery(req *QueryOrderParam) *myasterapi.SpotOrderGetApi {
	api := aster.NewSpotRestClient(b.apiKey, b.secretKey).NewSpotOrderGet().Symbol(req.Symbol)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api = api.OrderId(orderId)
	} else {
		if req.ClientOrderId != "" {
			api = api.OrigClientOrderId(req.ClientOrderId)
		}
	}
	return api
}
func (b *AsterTradeEngine) apiSpotOrdersQuery(req *QueryOrderParam) *myasterapi.SpotAllOrdersApi {
	api := aster.NewSpotRestClient(b.apiKey, b.secretKey).NewAllOrders().Symbol(req.Symbol)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api.OrderId(orderId)
	}
	if req.StartTime != 0 {
		api = api.StartTime(req.StartTime)
	}
	if req.EndTime != 0 {
		api = api.EndTime(req.EndTime)
	}
	if req.Limit != 0 {
		api = api.Limit(req.Limit)
	}
	return api
}
func (b *AsterTradeEngine) apiSpotTradeQuery(req *QueryTradeParam) *myasterapi.SpotMyTradesApi {
	api := aster.NewSpotRestClient(b.apiKey, b.secretKey).NewSpotMyTrades().
		Symbol(req.Symbol)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api = api.OrderId(orderId)
	} else {
		if req.StartTime != 0 {
			api = api.StartTime(req.StartTime)
		}
		if req.EndTime != 0 {
			api = api.EndTime(req.EndTime)
		}
	}

	if req.Limit != 0 {
		api = api.Limit(req.Limit)
	}
	return api
}

func (b *AsterTradeEngine) apiSpotOrderCreate(req *OrderParam) *myasterapi.SpotOrderPostApi {
	api := aster.NewSpotRestClient(b.apiKey, b.secretKey).NewSpotOrderPost().
		Symbol(req.Symbol).
		Side(b.asterConverter.ToAsterOrderSide(req.OrderSide)).
		Quantity(req.Quantity)

	if req.NewOrderRespType != "" {
		api.NewOrderRespType(req.NewOrderRespType)
	} else {
		api.NewOrderRespType("RESULT")
	}

	api.Type(b.asterConverter.ToTriggerBnOrderType(AsterAccountType(req.AccountType), req.OrderType, req.TriggerType))

	if !req.TriggerPrice.IsZero() {
		api.StopPrice(req.TriggerPrice)
	}

	// log.Info(req)

	if !req.Price.IsZero() {
		api = api.Price(req.Price)
	}
	if req.ClientOrderId != "" {
		api = api.NewClientOrderId(req.ClientOrderId)
	}
	if req.TimeInForce != "" {
		if req.TimeInForce == TIME_IN_FORCE_POST_ONLY {
			// 现货POSTONLY下单 不传timeInforce并且将订单类型为LIMIT_MAKER
			// api.TimeInForce(b.bnConverter.ToAsterTimeInForce(TIME_IN_FORCE_GTC))
			api.Type(ASTER_ORDER_TYPE_LIMIT_MAKER)
		} else {
			api = api.TimeInForce(b.asterConverter.ToAsterTimeInForce(req.TimeInForce))
		}
	}
	return api
}
func (b *AsterTradeEngine) apiSpotOrderAmend(req *OrderParam) *myasterapi.SpotOrderCancelReplaceApi {
	api := aster.NewSpotRestClient(b.apiKey, b.secretKey).NewSpotOrderCancelReplace().
		Symbol(req.Symbol).CancelReplaceMode("STOP_ON_FAILURE").
		Side(b.asterConverter.ToAsterOrderSide(req.OrderSide)).
		Quantity(req.Quantity)

	api.Type(b.asterConverter.ToTriggerBnOrderType(AsterAccountType(req.AccountType), req.OrderType, req.TriggerType))
	if !req.TriggerPrice.IsZero() {
		api.StopPrice(req.TriggerPrice)
	}
	if !req.Price.IsZero() {
		api = api.Price(req.Price)
	}
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api = api.CancelOrderId(orderId)
	}
	if req.ClientOrderId != "" {
		api = api.CancelOrigClientOrderId(req.ClientOrderId)
	}
	if req.NewClientOrderId != "" {
		api = api.NewClientOrderId(req.NewClientOrderId)
	}
	if req.TimeInForce != "" {
		api = api.TimeInForce(b.asterConverter.ToAsterTimeInForce(req.TimeInForce))
	}
	return api
}
func (b *AsterTradeEngine) apiSpotOrderCancel(req *OrderParam) *myasterapi.SpotOrderDeleteApi {
	api := aster.NewSpotRestClient(b.apiKey, b.secretKey).NewSpotOrderDelete().
		Symbol(req.Symbol)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api = api.OrderId(orderId)
	} else {
		api = api.OrigClientOrderId(req.ClientOrderId)
	}

	return api
}

// U本位合约订单API接口
func (b *AsterTradeEngine) apiFutureOpenOrders(req *QueryOrderParam) *myasterapi.FutureOpenOrdersApi {
	api := aster.NewFutureRestClient(b.apiKey, b.secretKey).NewOpenOrders()
	if req.Symbol != "" {
		api.Symbol(req.Symbol)
	}
	return api
}
func (b *AsterTradeEngine) apiFutureOrderQuery(req *QueryOrderParam) *myasterapi.FutureOrderGetApi {
	api := aster.NewFutureRestClient(b.apiKey, b.secretKey).NewFutureOrderGet().Symbol(req.Symbol)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api = api.OrderId(orderId)
	} else {
		api = api.OrigClientOrderId(req.ClientOrderId)
	}
	return api
}
func (b *AsterTradeEngine) apiFutureOrdersQuery(req *QueryOrderParam) *myasterapi.FutureAllOrdersApi {
	api := aster.NewFutureRestClient(b.apiKey, b.secretKey).NewAllOrders().Symbol(req.Symbol)

	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api = api.OrderId(orderId)
	}

	if req.StartTime != 0 {
		api = api.StartTime(req.StartTime)
	}
	if req.EndTime != 0 {
		api = api.EndTime(req.EndTime)
	}
	if req.Limit != 0 {
		api = api.Limit(int64(req.Limit))
	}

	return api
}
func (b *AsterTradeEngine) apiFutureTradeQuery(req *QueryTradeParam) *myasterapi.FutureUserTradesApi {
	api := aster.NewFutureRestClient(b.apiKey, b.secretKey).NewFutureUserTrades().
		Symbol(req.Symbol)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api = api.OrderId(orderId)
	}
	if req.StartTime != 0 {
		api = api.StartTime(req.StartTime)
	}
	if req.EndTime != 0 {
		api = api.EndTime(req.EndTime)
	}
	if req.Limit != 0 {
		api = api.Limit(int64(req.Limit))
	}
	return api
}

func (b *AsterTradeEngine) apiFutureOrderCreate(req *OrderParam) *myasterapi.FutureOrderPostApi {
	api := aster.NewFutureRestClient(b.apiKey, b.secretKey).NewFutureOrderPost().
		Symbol(req.Symbol).
		Side(b.asterConverter.ToAsterOrderSide(req.OrderSide)).
		PositionSide(b.asterConverter.ToAsterPositionSide(req.PositionSide)).
		Quantity(req.Quantity)

	if req.NewOrderRespType != "" {
		api.NewOrderRespType(req.NewOrderRespType)
	} else {
		api.NewOrderRespType("RESULT")
	}

	api.Type(b.asterConverter.ToTriggerBnOrderType(AsterAccountType(req.AccountType), req.OrderType, req.TriggerType))
	if !req.TriggerPrice.IsZero() {
		api.StopPrice(req.TriggerPrice)
	}
	if !req.Price.IsZero() {
		api = api.Price(req.Price)
	}
	if req.ClientOrderId != "" {
		api = api.NewClientOrderId(req.ClientOrderId)
	}
	if req.TimeInForce != "" {
		api = api.TimeInForce(b.asterConverter.ToAsterTimeInForce(req.TimeInForce))
	}
	return api
}
func (b *AsterTradeEngine) apiFutureOrderAmend(req *OrderParam) *myasterapi.FutureOrderPutApi {
	api := aster.NewFutureRestClient(b.apiKey, b.secretKey).NewFutureOrderPut().
		Symbol(req.Symbol).
		Side(b.asterConverter.ToAsterOrderSide(req.OrderSide)).
		Quantity(req.Quantity)

	if !req.Price.IsZero() {
		api = api.Price(req.Price)
	}
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api = api.OrderId(orderId)
	}
	if req.ClientOrderId != "" {
		api = api.OrigClientOrderId(req.ClientOrderId)
	}
	return api
}
func (b *AsterTradeEngine) apiFutureOrderCancel(req *OrderParam) *myasterapi.FutureOrderDeleteApi {
	api := aster.NewFutureRestClient(b.apiKey, b.secretKey).NewFutureOrderDelete().
		Symbol(req.Symbol)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api = api.OrderId(orderId)
	} else {
		api = api.OrigClientOrderId(req.ClientOrderId)
	}
	return api
}

func (b *AsterTradeEngine) apiFutureBatchOrderCreate(reqs []*OrderParam) *myasterapi.FutureBatchOrdersPostApi {
	client := aster.NewFutureRestClient(b.apiKey, b.secretKey)
	api := client.NewFutureBatchOrdersPost()
	for _, req := range reqs {
		thisApi := client.NewFutureOrderPost().Symbol(req.Symbol).
			Side(b.asterConverter.ToAsterOrderSide(req.OrderSide)).
			PositionSide(b.asterConverter.ToAsterPositionSide(req.PositionSide)).
			Quantity(req.Quantity)

		if req.NewOrderRespType != "" {
			thisApi.NewOrderRespType(req.NewOrderRespType)
		} else {
			thisApi.NewOrderRespType("RESULT")
		}
		thisApi.Type(b.asterConverter.ToTriggerBnOrderType(AsterAccountType(req.AccountType), req.OrderType, req.TriggerType))
		if !req.TriggerPrice.IsZero() {
			thisApi.StopPrice(req.TriggerPrice)
		}
		if !req.Price.IsZero() {
			thisApi = thisApi.Price(req.Price)
		}
		if req.ClientOrderId != "" {
			thisApi = thisApi.NewClientOrderId(req.ClientOrderId)
		}
		if req.TimeInForce != "" {
			thisApi = thisApi.TimeInForce(b.asterConverter.ToAsterTimeInForce(req.TimeInForce))
		}
		api = api.AddOrders(thisApi)
	}
	return api
}
func (b *AsterTradeEngine) apiFutureBatchOrderAmend(reqs []*OrderParam) *myasterapi.FutureBatchOrdersPutApi {
	client := aster.NewFutureRestClient(b.apiKey, b.secretKey)
	api := client.NewFutureBatchOrdersPut()
	for _, req := range reqs {
		thisApi := client.NewFutureOrderPut().Symbol(req.Symbol).
			Side(b.asterConverter.ToAsterOrderSide(req.OrderSide)).
			Quantity(req.Quantity)
		if !req.Price.IsZero() {
			thisApi = thisApi.Price(req.Price)
		}
		if req.OrderId != "" {
			orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
			thisApi = thisApi.OrderId(orderId)
		}
		if req.ClientOrderId != "" {
			thisApi = thisApi.OrigClientOrderId(req.ClientOrderId)
		}
		api = api.AddOrders(thisApi)
	}
	return api
}
func (b *AsterTradeEngine) apiFutureBatchOrderCancel(reqs []*OrderParam) (*myasterapi.FutureBatchOrdersDeleteApi, error) {
	client := aster.NewFutureRestClient(b.apiKey, b.secretKey)
	orderIds := []int64{}
	clientOrderIds := []string{}
	for _, req := range reqs {
		if req.OrderId != "" {
			orderId, err := strconv.ParseInt(req.OrderId, 10, 64)
			if err != nil {
				return nil, ErrorInvalid("order id")
			}
			orderIds = append(orderIds, orderId)
		} else if req.ClientOrderId != "" {
			clientOrderIds = append(clientOrderIds, req.ClientOrderId)
		} else {
			return nil, ErrorInvalid("order id or client order id is required")
		}
	}
	api := client.NewFutureBatchOrdersDelete().
		Symbol(reqs[0].Symbol)
	if len(orderIds) > 0 {
		api = api.OrderIdList(orderIds)
	} else if len(clientOrderIds) > 0 {
		api = api.OrigClientOrderIdList(clientOrderIds)
	} else {
		return nil, ErrorInvalid("order id or client order id is required")
	}
	return api, nil
}

func (b *AsterTradeEngine) restBatchPreCheck(reqs []*OrderParam) error {
	//检测长度，ASTER最多批量下5个订单
	if len(reqs) > 5 {
		return ErrorInvalid("aster order param length require less than 5")
	}

	//检测类型是否相同
	for _, req := range reqs {
		if err := b.accountTypePreCheck(req.AccountType); err != nil {
			return err
		}
		if req.AccountType != reqs[0].AccountType {
			return ErrorInvalid("order param account type require same")
		}
	}

	return nil
}

func (b *AsterTradeEngine) accountTypePreCheck(accountType string) error {
	switch AsterAccountType(accountType) {
	case ASTER_AC_SPOT, ASTER_AC_FUTURE:
		return nil
	default:
		return ErrorInvalid("aster account type invalid")
	}
}
func (b *AsterTradeEngine) checkWsSpotAccount() error {
	var err error
	if b.wsSpotAccount == nil {
		b.wsSpotAccount, err = aster.NewSpotWsStreamClient().ConvertToAccountWs(b.apiKey, b.secretKey, myasterapi.SPOT_WS_TYPE)
		if err != nil {
			return err
		}
		err := b.wsSpotAccount.OpenConn()
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *AsterTradeEngine) checkWsFutureAccount() error {
	var err error
	if b.wsFutureAccount == nil {
		b.wsFutureAccount, err = aster.NewFutureWsStreamClient().ConvertToAccountWs(b.apiKey, b.secretKey)
		if err != nil {
			return err
		}
		err := b.wsFutureAccount.OpenConn()
		if err != nil {
			return err
		}
	}
	return nil
}
