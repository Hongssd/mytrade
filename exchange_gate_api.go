package mytrade

import (
	"strconv"
	"strings"

	"github.com/Hongssd/mygateapi"
	"github.com/shopspring/decimal"
)

func (g *GateTradeEngine) apiSpotOrderCreate(req *OrderParam) *mygateapi.PrivateRestSpotOrdersPostAPI {
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().
		NewPrivateRestSpotOrdersPost().
		// Account(req.AccountType).
		CurrencyPair(req.Symbol).
		Type(g.gateConverter.ToGateOrderType(req.OrderType)).
		Side(g.gateConverter.ToGateOrderSide(req.OrderSide)).
		Amount(req.Quantity)

	account := g.gateConverter.ToOrderSpotAccountType(GateAccountType(req.AccountType), req.IsMargin, req.IsIsolated)
	api.Account(account.String())

	if req.OrderType == ORDER_TYPE_LIMIT {
		if !req.Price.IsZero() {
			api.Price(req.Price)
		}
	}

	// 自定义订单id
	if req.ClientOrderId != "" {
		api.Text(req.ClientOrderId)
	}

	if req.TimeInForce != "" {
		api.TimeInForce(g.gateConverter.ToGateTimeInForce(req.TimeInForce))
	}

	return api
}
func (g *GateTradeEngine) apiSpotPriceOrderCreate(req *OrderParam) *mygateapi.PrivateRestSpotPriceOrdersPostAPI {
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestSpotPriceOrdersPost().
		Market(req.Symbol)

	var expiration int
	if req.Expiration != 0 {
		expiration = req.Expiration
	}
	api.Trigger(mygateapi.PrivateRestSpotPriceOrdersPostTriggerReq{
		Price:      GetPointer(req.TriggerPrice.String()),
		Rule:       GetPointer(g.gateConverter.ToGateTriggerRule(req.TriggerType, req.OrderSide)),
		Expiration: GetPointer(expiration),
	})

	var timeInForce string
	if req.TimeInForce != "" {
		timeInForce = g.gateConverter.ToGateTimeInForce(req.TimeInForce)
	}

	var text string
	if req.ClientOrderId != "" {
		text = req.ClientOrderId
	}
	api.Put(mygateapi.PrivateRestSpotPriceOrdersPostPutReq{
		Account:     GetPointer(g.gateConverter.ToGateSpotPriceOrderAccount(GateAccountType(req.AccountType))),
		Type:        GetPointer(g.gateConverter.ToGateOrderType(req.OrderType)),
		Side:        GetPointer(g.gateConverter.ToGateOrderSide(req.OrderSide)),
		Amount:      GetPointer(req.Quantity.String()),
		Price:       GetPointer(req.Price.String()),
		TimeInForce: GetPointer(timeInForce),
		Text:        GetPointer(text),
	})

	return api
}

func (g *GateTradeEngine) apiFuturesOrderCreate(req *OrderParam) *mygateapi.PrivateRestFuturesSettleOrdersPostAPI {
	split := strings.Split(req.Symbol, "_")
	if len(split) != 2 {
		log.Error("symbol error")
		return nil
	}
	settle := strings.ToLower(split[1])
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestFuturesSettleOrdersPost().
		Settle(settle).Contract(req.Symbol).Price(req.Price.String())

	if req.OrderSide == ORDER_SIDE_BUY {
		api.Size(req.Quantity.Abs().IntPart())
		if req.PositionSide == POSITION_SIDE_SHORT {
			//BUY SHORT 只平仓
			api.ReduceOnly(true)
		}
	} else {
		api.Size(req.Quantity.Abs().Neg().IntPart())
		if req.PositionSide == POSITION_SIDE_LONG {
			//SELL LONG 只平仓
			api.ReduceOnly(true)
		}
	}

	if req.TimeInForce != "" {
		api.Tif(g.gateConverter.ToGateTimeInForce(req.TimeInForce))
	}

	// 市价委托指定tif为ioc
	if req.OrderType == ORDER_TYPE_MARKET {
		api.Tif(GATE_TIME_IN_FORCE_IOC).Price(decimal.Zero.String())
	}

	// 平仓时 size = 0, close = true
	if req.Quantity.IsZero() {
		api.Close(true)
	}

	if req.ReduceOnly {
		api.ReduceOnly(req.ReduceOnly)
	}
	if req.ClientOrderId != "" {
		api.Text(req.ClientOrderId)
	}

	return api
}
func (g *GateTradeEngine) apiDeliveryOrderCreate(req *OrderParam) *mygateapi.PrivateRestDeliverySettleOrdersPostAPI {
	split := strings.Split(req.Symbol, "_")
	if len(split) != 3 {
		log.Error("symbol error")
		return nil
	}
	settle := strings.ToLower(split[1])
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestDeliverySettleOrdersPost().
		Settle(settle).Contract(req.Symbol).Size(req.Quantity.IntPart()).Price(req.Price.String())
	if req.OrderSide == ORDER_SIDE_BUY {
		api.Size(req.Quantity.Abs().IntPart())
		if req.PositionSide == POSITION_SIDE_SHORT {
			//BUY SHORT 只平仓
			api.ReduceOnly(true)
		}
	} else {
		api.Size(req.Quantity.Abs().Neg().IntPart())
		if req.PositionSide == POSITION_SIDE_LONG {
			//SELL LONG 只平仓
			api.ReduceOnly(true)
		}
	}

	if req.TimeInForce != "" {
		api.Tif(g.gateConverter.ToGateTimeInForce(req.TimeInForce))
	}

	// 市价委托指定tif为ioc
	if req.OrderType == ORDER_TYPE_MARKET {
		api.Tif(GATE_TIME_IN_FORCE_IOC).Price(decimal.Zero.String())
	}

	// 平仓时 size = 0, close = true
	if req.Quantity.IsZero() {
		api.Close(true)
	}

	if req.ReduceOnly {
		api.ReduceOnly(req.ReduceOnly)
	}
	if req.ClientOrderId != "" {
		api.Text(req.ClientOrderId)
	}

	return api
}

func (g *GateTradeEngine) apiSpotOrderAmend(req *OrderParam) *mygateapi.PrivateRestSpotOrdersOrderIdPatchAPI {
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().
		NewPrivateRestSpotOrdersOrderIdPatch().
		OrderId(req.OrderId)

	if req.Symbol != "" {
		api.CurrencyPair(req.Symbol)
	}

	account := g.gateConverter.ToOrderSpotAccountType(GateAccountType(req.AccountType), req.IsMargin, req.IsIsolated)
	api.Account(account.String())

	if !req.Quantity.IsZero() {
		api.Amount(req.Quantity)
	}

	if !req.Price.IsZero() {
		api.Price(req.Price)
	}

	if req.ClientOrderId != "" {
		api.AmendText(req.ClientOrderId)
	}

	return api
}
func (g *GateTradeEngine) apiFuturesOrderAmend(req *OrderParam) *mygateapi.PrivateRestFuturesSettleOrdersOrderIdPutAPI {
	split := strings.Split(req.Symbol, "_")
	if len(split) != 2 {
		log.Error("symbol error")
		return nil
	}
	settle := strings.ToLower(split[1])
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestFuturesSettleOrdersOrderIdPut().
		Settle(settle).OrderId(req.OrderId)

	if req.Quantity != decimal.Zero {
		api.Size(req.Quantity.IntPart())
	}

	if req.Price != decimal.Zero {
		api.Price(req.Price.String())
	}

	if req.ClientOrderId != "" {
		api.AmendText(req.ClientOrderId)
	}

	return api
}

func (g *GateTradeEngine) apiSpotOrderCancel(req *OrderParam) *mygateapi.PrivateRestSpotOrdersOrderIdDeleteAPI {
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestSpotOrdersOrderIdDelete().CurrencyPair(req.Symbol)

	if req.OrderId != "" {
		api.OrderId(req.OrderId)
	}

	if req.ClientOrderId != "" {
		api.OrderId(req.ClientOrderId)
	}

	account := g.gateConverter.ToOrderSpotAccountType(GateAccountType(req.AccountType), req.IsMargin, req.IsIsolated)
	api.Account(account.String())

	return api
}
func (g *GateTradeEngine) apiSpotPriceOrderCancel(req *OrderParam) *mygateapi.PrivateRestSpotPriceOrdersOrderIdDeleteAPI {
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestSpotPriceOrdersOrderIdDelete().OrderId(req.OrderId)
	return api
}
func (g *GateTradeEngine) apiFuturesOrderCancel(req *OrderParam) *mygateapi.PrivateRestFuturesSettleOrdersOrderIdDeleteAPI {
	split := strings.Split(req.Symbol, "_")
	if len(split) != 2 {
		log.Error("symbol error")
		return nil
	}
	settle := strings.ToLower(split[1])
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestFuturesSettleOrdersOrderIdDelete().
		Settle(settle)

	if req.OrderId != "" {
		api.OrderId(req.OrderId)
	}

	if req.ClientOrderId != "" {
		api.OrderId(req.ClientOrderId)
	}

	return api
}
func (g *GateTradeEngine) apiDeliveryOrderCancel(req *OrderParam) *mygateapi.PrivateRestDeliverySettleOrdersOrderIdDeleteAPI {
	split := strings.Split(req.Symbol, "_")
	if len(split) != 3 {
		log.Error("symbol error")
		return nil
	}
	settle := strings.ToLower(split[1])
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestDeliverySettleOrdersOrderIdDelete().
		Settle(settle)

	if req.OrderId != "" {
		api.OrderId(req.OrderId)
	}

	if req.ClientOrderId != "" {
		api.OrderId(req.ClientOrderId)
	}

	return api
}

func (g *GateTradeEngine) apiSpotOpenOrders(req *QueryOrderParam) *mygateapi.PrivateRestSpotOpenOrdersAPI {
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().
		NewPrivateRestSpotOpenOrders()

	account := g.gateConverter.ToOrderSpotAccountType(GateAccountType(req.AccountType), req.IsMargin, req.IsIsolated)
	if account != GATE_ACCOUNT_TYPE_SPOT {
		api.Account(account.String())
	}

	if req.Limit != 0 {
		api.Limit(req.Limit)
	}

	return api
}
func (g *GateTradeEngine) apiFuturesOpenOrders(req *QueryOrderParam) *mygateapi.PrivateRestFuturesSettleOrdersGetAPI {
	split := strings.Split(req.Symbol, "_")
	if len(split) != 2 {
		log.Error("symbol error")
		return nil
	}
	settle := strings.ToLower(split[1])
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestFuturesSettleOrdersGet().
		Settle(settle).Status(GATE_ORDER_CONTRACT_STATUS_OPEN)

	if req.Limit != 0 {
		api.Limit(req.Limit)
	}

	return api
}
func (g *GateTradeEngine) apiDeliveryOpenOrders(req *QueryOrderParam) *mygateapi.PrivateRestDeliverySettleOrdersGetAPI {
	split := strings.Split(req.Symbol, "_")
	if len(split) != 3 {
		log.Error("symbol error")
		return nil
	}
	settle := strings.ToLower(split[1])
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestDeliverySettleOrdersGet().
		Settle(settle).Status(GATE_ORDER_CONTRACT_STATUS_OPEN)

	if req.Limit != 0 {
		api.Limit(req.Limit)
	}

	return api
}

func (g *GateTradeEngine) apiSpotOrderQuery(req *QueryOrderParam) *mygateapi.PrivateRestSpotOrdersOrderIdGetAPI {
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestSpotOrdersOrderIdGet()

	account := g.gateConverter.ToOrderSpotAccountType(GateAccountType(req.AccountType), req.IsMargin, req.IsIsolated)
	api.Account(account.String())

	if req.Symbol != "" {
		api.CurrencyPair(req.Symbol)
	}

	if req.OrderId != "" {
		api.OrderId(req.OrderId)
	}

	if req.ClientOrderId != "" {
		api.OrderId(req.ClientOrderId)
	}

	return api
}
func (g *GateTradeEngine) apiFuturesOrderQuery(req *QueryOrderParam) *mygateapi.PrivateRestFuturesSettleOrdersOrderIdGetAPI {
	split := strings.Split(req.Symbol, "_")
	if len(split) != 2 {
		log.Error("symbol error")
		return nil
	}
	settle := strings.ToLower(split[1])
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestFuturesSettleOrdersOrderIdGet().
		Settle(settle)

	if req.OrderId != "" {
		api.OrderId(req.OrderId)
	}

	if req.ClientOrderId != "" {
		api.OrderId(req.ClientOrderId)
	}

	return api
}
func (g *GateTradeEngine) apiDeliveryOrderQuery(req *QueryOrderParam) *mygateapi.PrivateRestDeliverySettleOrdersOrderIdGetAPI {
	split := strings.Split(req.Symbol, "_")
	if len(split) != 3 {
		log.Error("symbol error")
		return nil
	}
	settle := strings.ToLower(split[1])
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestDeliverySettleOrdersOrderIdGet().
		Settle(settle)

	if req.OrderId != "" {
		api.OrderId(req.OrderId)
	}

	if req.ClientOrderId != "" {
		api.OrderId(req.ClientOrderId)
	}

	return api
}

func (g *GateTradeEngine) apiSpotOrdersQuery(req *QueryOrderParam) *mygateapi.PrivateRestSpotOrdersGetAPI {
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestSpotOrdersGet()

	api.Status(GATE_ORDER_CONTRACT_STATUS_FINISHED)

	// account := g.gateConverter.ToOrderSpotAccountType(GateAccountType(req.AccountType), req.IsMargin, req.IsIsolated)
	// api.Account("unified")

	if req.Symbol != "" {
		api.CurrencyPair(req.Symbol)
	}

	if req.Limit != 0 {
		api.Limit(req.Limit)
	} else {
		api.Limit(100)
	}
	return api
}
func (g *GateTradeEngine) apiFuturesOrdersQuery(req *QueryOrderParam) *mygateapi.PrivateRestFuturesSettleOrdersGetAPI {
	settle := "usdt"
	if req.Symbol != "" {
		split := strings.Split(req.Symbol, "_")
		if len(split) != 2 {
			log.Error("symbol error")
			return nil
		}
		settle = strings.ToLower(split[1])
	}

	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().
		NewPrivateRestFuturesSettleOrdersGet().
		Settle(settle).Status(GATE_ORDER_CONTRACT_STATUS_FINISHED)

	if req.Limit != 0 {
		api.Limit(req.Limit)
	} else {
		api.Limit(100)
	}
	return api
}
func (g *GateTradeEngine) apiDeliveryOrdersQuery(req *QueryOrderParam) *mygateapi.PrivateRestDeliverySettleOrdersGetAPI {
	settle := "usdt"
	if req.Symbol != "" {
		split := strings.Split(req.Symbol, "_")
		if len(split) != 3 {
			log.Error("symbol error")
			return nil
		}
		settle = strings.ToLower(split[1])
	}
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestDeliverySettleOrdersGet().Settle(settle).
		Status(GATE_ORDER_CONTRACT_STATUS_FINISHED)

	if req.Limit != 0 {
		api.Limit(req.Limit)
	} else {
		api.Limit(100)
	}
	return api
}

func (g *GateTradeEngine) apiSpotTradesQuery(req *QueryTradeParam) *mygateapi.PrivateRestSpotMyTradesAPI {
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestSpotMyTrades()

	account := g.gateConverter.ToOrderSpotAccountType(GateAccountType(req.AccountType), req.IsMargin, req.IsIsolated)
	api.Account(account.String())

	if req.Symbol != "" {
		api.CurrencyPair(req.Symbol)
	}

	if req.OrderId != "" {
		api.OrderId(req.OrderId)
	}

	if req.Limit != 0 {
		api.Limit(req.Limit)
	}

	if req.StartTime != 0 {
		api.From(req.StartTime)
	}

	if req.EndTime != 0 {
		api.To(req.EndTime)
	}

	if req.Limit != 0 {
		api.Limit(req.Limit)
	}

	return api
}
func (g *GateTradeEngine) apiFuturesTradesQuery(req *QueryTradeParam) *mygateapi.PrivateRestFuturesSettleMyTradesAPI {
	split := strings.Split(req.Symbol, "_")
	if len(split) != 2 {
		log.Error("symbol error")
		return nil
	}
	settle := strings.ToLower(split[1])
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestFuturesSettleMyTrades().
		Settle(settle).Contract(req.Symbol)

	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api.Order(orderId)
	}

	if req.Limit != 0 {
		api.Limit(req.Limit)
	}

	return api
}
func (g *GateTradeEngine) apiDeliveryTradesQuery(req *QueryTradeParam) *mygateapi.PrivateRestDeliverySettleMyTradesAPI {
	split := strings.Split(req.Symbol, "_")
	if len(split) != 3 {
		log.Error("symbol error")
		return nil
	}
	settle := strings.ToLower(split[1])
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestDeliverySettleMyTrades().
		Settle(settle).Contract(req.Symbol)

	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api.Order(orderId)
	}

	if req.Limit != 0 {
		api.Limit(req.Limit)
	}

	return api
}
