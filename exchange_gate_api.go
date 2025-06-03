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
		Amount(req.Quantity).Price(req.Price)

	account := g.gateConverter.ToOrderSpotAccountType(GateAccountType(req.AccountType), req.IsMargin, req.IsIsolated)
	api.Account(account.String())

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
		Rule:       GetPointer(g.gateConverter.ToGateSpotPriceOrderTriggerRule(req.TriggerType, req.OrderSide)),
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
	account := g.gateConverter.ToOrderSpotAccountType(GateAccountType(req.AccountType), req.IsMargin, req.IsIsolated)
	api.Put(mygateapi.PrivateRestSpotPriceOrdersPostPutReq{
		Account:     GetPointer(g.gateConverter.ToGateSpotPriceOrderAccount(account)),
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
	settle := "usdt"
	split := strings.Split(req.Symbol, "_")
	if len(split) == 2 {
		settle = strings.ToLower(split[1])
	}
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
func (g *GateTradeEngine) apiFuturesPriceOrderCreate(req *OrderParam) *mygateapi.PrivateRestFuturesSettlePriceOrdersPostAPI {
	settle := "usdt"
	split := strings.Split(req.Symbol, "_")
	if len(split) == 2 {
		settle = strings.ToLower(split[1])
	}

	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestFuturesSettlePriceOrdersPost().
		Settle(settle)

	var size int64
	if req.OrderSide == ORDER_SIDE_BUY {
		size = req.Quantity.Abs().IntPart()
	} else if req.OrderSide == ORDER_SIDE_SELL {
		size = req.Quantity.Abs().Neg().IntPart()
	}

	var price string
	var tif string
	if req.OrderType == ORDER_TYPE_LIMIT {
		price = req.Price.String()
		tif = g.gateConverter.ToGateTimeInForce(req.TimeInForce)
	} else {
		price = decimal.Zero.String()
		tif = GATE_TIME_IN_FORCE_IOC
	}

	api.Initial(mygateapi.PrivateRestFuturesSettlePriceOrdersPostInitialReq{
		Contract:   GetPointer(req.Symbol),
		Size:       GetPointer(size),
		Price:      GetPointer(price),
		Tif:        GetPointer(tif),
		Text:       GetPointer(req.ClientOrderId),
		ReduceOnly: GetPointer(req.ReduceOnly),
		AutoSize:   GetPointer(req.GateAutoSize),
	})

	api.Trigger(mygateapi.PrivateRestFuturesSettlePriceOrdersPostTriggerReq{
		StrategyType: GetPointer(int32(0)),
		Rule:         GetPointer(g.gateConverter.ToGateFuturesPriceOrderTriggerRule(req.TriggerType, req.OrderSide)),
		Price:        GetPointer(req.TriggerPrice.String()),
		PriceType:    GetPointer(int32(0)),
		Expiration:   GetPointer(req.Expiration),
	})

	api.OrderType(req.GatePriceOrderType)

	return api
}

func (g *GateTradeEngine) apiDeliveryOrderCreate(req *OrderParam) *mygateapi.PrivateRestDeliverySettleOrdersPostAPI {
	settle := "usdt"
	split := strings.Split(req.Symbol, "_")
	if len(split) == 3 {
		settle = strings.ToLower(split[1])
	}
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
func (g *GateTradeEngine) apiDeliveryPriceOrderCreate(req *OrderParam) *mygateapi.PrivateRestDeliverySettlePriceOrdersPostAPI {
	settle := "usdt"
	split := strings.Split(req.Symbol, "_")
	if len(split) == 3 {
		settle = strings.ToLower(split[1])
	}
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestDeliverySettlePriceOrdersPost().
		Settle(settle)

	var size int64
	if req.OrderSide == ORDER_SIDE_BUY {
		size = req.Quantity.Abs().IntPart()
	} else if req.OrderSide == ORDER_SIDE_SELL {
		size = req.Quantity.Abs().Neg().IntPart()
	}

	var price string
	var tif string
	if req.OrderType == ORDER_TYPE_LIMIT {
		price = req.Price.String()
		tif = g.gateConverter.ToGateTimeInForce(req.TimeInForce)
	} else {
		price = decimal.Zero.String()
		tif = GATE_TIME_IN_FORCE_IOC
	}

	api.Initial(mygateapi.PrivateRestFuturesSettlePriceOrdersPostInitialReq{
		Contract:   GetPointer(req.Symbol),
		Size:       GetPointer(size),
		Price:      GetPointer(price),
		Tif:        GetPointer(tif),
		Text:       GetPointer(req.ClientOrderId),
		ReduceOnly: GetPointer(req.ReduceOnly),
		AutoSize:   GetPointer(req.GateAutoSize),
	})

	api.Trigger(mygateapi.PrivateRestFuturesSettlePriceOrdersPostTriggerReq{
		StrategyType: GetPointer(int32(0)),
		Rule:         GetPointer(g.gateConverter.ToGateFuturesPriceOrderTriggerRule(req.TriggerType, req.OrderSide)),
		Price:        GetPointer(req.TriggerPrice.String()),
		PriceType:    GetPointer(int32(0)),
		Expiration:   GetPointer(req.Expiration),
	})

	api.OrderType(req.GatePriceOrderType)
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
	settle := "usdt"
	split := strings.Split(req.Symbol, "_")
	if len(split) == 2 {
		settle = strings.ToLower(split[1])
	}
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
	settle := "usdt"
	split := strings.Split(req.Symbol, "_")
	if len(split) == 2 {
		settle = strings.ToLower(split[1])
	}
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
func (g *GateTradeEngine) apiFuturesPriceOrderCancel(req *OrderParam) *mygateapi.PrivateRestFuturesSettlePriceOrdersOrderIdDeleteAPI {
	settle := "usdt"
	split := strings.Split(req.Symbol, "_")
	if len(split) == 2 {
		settle = strings.ToLower(split[1])
	}

	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestFuturesSettlePriceOrdersOrderIdDelete().
		Settle(settle).OrderId(req.OrderId)

	return api
}
func (g *GateTradeEngine) apiDeliveryOrderCancel(req *OrderParam) *mygateapi.PrivateRestDeliverySettleOrdersOrderIdDeleteAPI {
	settle := "usdt"
	split := strings.Split(req.Symbol, "_")
	if len(split) == 3 {
		settle = strings.ToLower(split[1])
	}
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
func (g *GateTradeEngine) apiDeliveryPriceOrderCancel(req *OrderParam) *mygateapi.PrivateRestDeliverySettlePriceOrdersOrderIdDeleteAPI {
	settle := "usdt"
	split := strings.Split(req.Symbol, "_")
	if len(split) == 3 {
		settle = strings.ToLower(split[1])
	}
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestDeliverySettlePriceOrdersOrderIdDelete().
		Settle(settle).OrderId(req.OrderId)

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
func (g *GateTradeEngine) apiSpotPriceOpenOrders(req *QueryOrderParam) *mygateapi.PrivateRestSpotPriceOrdersGetAPI {
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestSpotPriceOrdersGet().
		Status(GATE_ORDER_SPOT_PRICE_STATUS_OPEN)

	account := g.gateConverter.ToOrderSpotAccountType(GateAccountType(req.AccountType), req.IsMargin, req.IsIsolated)
	if account != GATE_ACCOUNT_TYPE_SPOT {
		api.Account(g.gateConverter.ToGateSpotPriceOrderAccount(account))
	}

	if req.Symbol != "" {
		api.Market(req.Symbol)
	}

	if req.Limit != 0 {
		api.Limit(req.Limit)
	}

	return api
}

func (g *GateTradeEngine) apiFuturesOpenOrders(req *QueryOrderParam) *mygateapi.PrivateRestFuturesSettleOrdersGetAPI {
	settle := "usdt"
	split := strings.Split(req.Symbol, "_")
	if len(split) == 2 {
		settle = strings.ToLower(split[1])
	}

	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestFuturesSettleOrdersGet().
		Settle(settle).Status(GATE_ORDER_CONTRACT_STATUS_OPEN)

	if req.Limit != 0 {
		api.Limit(req.Limit)
	}

	return api
}
func (g *GateTradeEngine) apiFuturesPriceOpenOrders(req *QueryOrderParam) *mygateapi.PrivateRestFuturesSettlePriceOrdersGetAPI {
	settle := "usdt"
	split := strings.Split(req.Symbol, "_")
	if len(split) == 2 {
		settle = strings.ToLower(split[1])
	}
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestFuturesSettlePriceOrdersGet().
		Settle(settle).Contract(req.Symbol).Status(GATE_ORDER_SPOT_PRICE_STATUS_OPEN)

	if req.Limit != 0 {
		api.Limit(req.Limit)
	}

	return api
}
func (g *GateTradeEngine) apiDeliveryOpenOrders(req *QueryOrderParam) *mygateapi.PrivateRestDeliverySettleOrdersGetAPI {
	settle := "usdt"
	split := strings.Split(req.Symbol, "_")
	if len(split) == 3 {
		settle = strings.ToLower(split[1])
	}
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestDeliverySettleOrdersGet().
		Settle(settle).Status(GATE_ORDER_CONTRACT_STATUS_OPEN)

	if req.Limit != 0 {
		api.Limit(req.Limit)
	}

	return api
}
func (g *GateTradeEngine) apiDeliveryPriceOpenOrders(req *QueryOrderParam) *mygateapi.PrivateRestDeliverySettlePriceOrdersGetAPI {
	settle := "usdt"
	split := strings.Split(req.Symbol, "_")
	if len(split) == 3 {
		settle = strings.ToLower(split[1])
	}
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestDeliverySettlePriceOrdersGet().
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
func (g *GateTradeEngine) apiSpotPriceOrderQuery(req *QueryOrderParam) *mygateapi.PrivateRestSpotPriceOrdersOrderIdGetAPI {
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestSpotPriceOrdersOrderIdGet().
		OrderId(req.OrderId)

	return api
}
func (g *GateTradeEngine) apiFuturesOrderQuery(req *QueryOrderParam) *mygateapi.PrivateRestFuturesSettleOrdersOrderIdGetAPI {
	settle := "usdt"
	split := strings.Split(req.Symbol, "_")
	if len(split) == 2 {
		settle = strings.ToLower(split[1])
	}
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
func (g *GateTradeEngine) apiFuturesPriceOrderQuery(req *QueryOrderParam) *mygateapi.PrivateRestFuturesSettlePriceOrdersOrderIdGetAPI {
	settle := "usdt"
	split := strings.Split(req.Symbol, "_")
	if len(split) == 2 {
		settle = strings.ToLower(split[1])
	}
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestFuturesSettlePriceOrdersOrderIdGet().
		Settle(settle).OrderId(req.OrderId)

	return api
}

func (g *GateTradeEngine) apiDeliveryOrderQuery(req *QueryOrderParam) *mygateapi.PrivateRestDeliverySettleOrdersOrderIdGetAPI {
	settle := "usdt"
	split := strings.Split(req.Symbol, "_")
	if len(split) == 3 {
		settle = strings.ToLower(split[1])
	}
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
func (g *GateTradeEngine) apiDeliveryPriceOrderQuery(req *QueryOrderParam) *mygateapi.PrivateRestDeliverySettlePriceOrdersOrderIdGetAPI {
	settle := "usdt"
	split := strings.Split(req.Symbol, "_")
	if len(split) == 3 {
		settle = strings.ToLower(split[1])
	}
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestDeliverySettlePriceOrdersOrderIdGet().
		Settle(settle).OrderId(req.OrderId)

	return api
}

func (g *GateTradeEngine) apiSpotOrdersQuery(req *QueryOrderParam) *mygateapi.PrivateRestSpotOrdersGetAPI {
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestSpotOrdersGet()

	api.Status(GATE_ORDER_CONTRACT_STATUS_FINISHED)

	account := g.gateConverter.ToOrderSpotAccountType(GateAccountType(req.AccountType), req.IsMargin, req.IsIsolated)
	if account != GATE_ACCOUNT_TYPE_SPOT {
		api.Account(account.String())
	}

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
func (g *GateTradeEngine) apiSpotPriceOrdersQuery(req *QueryOrderParam) *mygateapi.PrivateRestSpotPriceOrdersGetAPI {
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().
		NewPrivateRestSpotPriceOrdersGet().
		Status(GATE_ORDER_CONTRACT_STATUS_FINISHED)

	account := g.gateConverter.ToOrderSpotAccountType(GateAccountType(req.AccountType), req.IsMargin, req.IsIsolated)
	if account != GATE_ACCOUNT_TYPE_SPOT {
		api.Account(g.gateConverter.ToGateSpotPriceOrderAccount(account))
	}

	if req.Symbol != "" {
		api.Market(req.Symbol)
	}

	if req.Limit != 0 {
		api.Limit(req.Limit)
	}

	if req.Symbol != "" {
		api.Market(req.Symbol)
	}

	return api
}
func (g *GateTradeEngine) apiFuturesOrdersQuery(req *QueryOrderParam) *mygateapi.PrivateRestFuturesSettleOrdersGetAPI {
	settle := "usdt"
	if req.Symbol != "" {
		split := strings.Split(req.Symbol, "_")
		if len(split) == 2 {
			settle = strings.ToLower(split[1])
		}
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
func (g *GateTradeEngine) apiFuturesPriceOrdersQuery(req *QueryOrderParam) *mygateapi.PrivateRestFuturesSettlePriceOrdersGetAPI {
	settle := "usdt"
	split := strings.Split(req.Symbol, "_")
	if len(split) == 2 {
		settle = strings.ToLower(split[1])
	}
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestFuturesSettlePriceOrdersGet().
		Settle(settle).Contract(req.Symbol).
		Status(GATE_ORDER_CONTRACT_STATUS_FINISHED)

	if req.Limit != 0 {
		api.Limit(req.Limit)
	}

	return api
}
func (g *GateTradeEngine) apiDeliveryOrdersQuery(req *QueryOrderParam) *mygateapi.PrivateRestDeliverySettleOrdersGetAPI {
	settle := "usdt"
	if req.Symbol != "" {
		split := strings.Split(req.Symbol, "_")
		if len(split) == 3 {
			settle = strings.ToLower(split[1])
		}
	}
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().
		NewPrivateRestDeliverySettleOrdersGet().
		Settle(settle).
		Status(GATE_ORDER_CONTRACT_STATUS_FINISHED)

	if req.Limit != 0 {
		api.Limit(req.Limit)
	} else {
		api.Limit(100)
	}
	return api
}
func (g *GateTradeEngine) apiDeliveryPriceOrdersQuery(req *QueryOrderParam) *mygateapi.PrivateRestDeliverySettlePriceOrdersGetAPI {
	settle := "usdt"
	split := strings.Split(req.Symbol, "_")
	if len(split) == 3 {
		settle = strings.ToLower(split[1])
	}
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestDeliverySettlePriceOrdersGet().
		Settle(settle).Contract(req.Symbol).Status(GATE_ORDER_CONTRACT_STATUS_FINISHED)

	if req.Limit != 0 {
		api.Limit(req.Limit)
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
	settle := "usdt"
	split := strings.Split(req.Symbol, "_")
	if len(split) == 2 {
		settle = strings.ToLower(split[1])
	}
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestFuturesSettleMyTrades().
		Settle(settle)

	if req.Symbol != "" {
		api.Contract(req.Symbol)
	}

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
	settle := "usdt"
	split := strings.Split(req.Symbol, "_")
	if len(split) == 3 {
		settle = strings.ToLower(split[1])
	}
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestDeliverySettleMyTrades().
		Settle(settle)

	if req.Symbol != "" {
		api.Contract(req.Symbol)
	}

	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api.Order(orderId)
	}

	if req.Limit != 0 {
		api.Limit(req.Limit)
	}

	return api
}

func (g *GateTradeEngine) checkWsForSpotOrder() error {
	if g.wsForSpotOrder == nil {
		g.wsForSpotOrder = mygateapi.NewSpotWsStreamClient(mygateapi.NewRestClient(g.apiKey, g.secretKey))
		err := g.wsForSpotOrder.OpenConn()
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *GateTradeEngine) checkWsForFuturesOrder() error {
	if g.wsForUSDTFuturesOrder == nil {
		g.wsForUSDTFuturesOrder = mygateapi.NewFuturesWsStreamClient(mygateapi.NewRestClient(g.apiKey, g.secretKey), mygateapi.USDT_CONTRACT)
		err := g.wsForUSDTFuturesOrder.OpenConn()
		if err != nil {
			return err
		}
	}
	if g.wsForBTCFuturesOrder == nil {
		g.wsForBTCFuturesOrder = mygateapi.NewFuturesWsStreamClient(mygateapi.NewRestClient(g.apiKey, g.secretKey), mygateapi.BTC_CONTRACT)
		err := g.wsForBTCFuturesOrder.OpenConn()
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *GateTradeEngine) checkWsForDeliveryOrder() error {
	if g.wsForUSDTDeliveryOrder == nil {
		g.wsForUSDTDeliveryOrder = mygateapi.NewDeliveryWsStreamClient(mygateapi.NewRestClient(g.apiKey, g.secretKey), mygateapi.USDT_CONTRACT)
		err := g.wsForUSDTDeliveryOrder.OpenConn()
		if err != nil {
			return err
		}
	}
	if g.wsForBTCFuturesDeliveryOrder == nil {
		g.wsForBTCFuturesDeliveryOrder = mygateapi.NewDeliveryWsStreamClient(mygateapi.NewRestClient(g.apiKey, g.secretKey), mygateapi.BTC_CONTRACT)
		err := g.wsForBTCFuturesDeliveryOrder.OpenConn()
		if err != nil {
			return err
		}
	}
	return nil
}
