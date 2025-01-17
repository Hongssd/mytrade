package mytrade

import (
	"github.com/Hongssd/mygateapi"
	"github.com/shopspring/decimal"
	"strconv"
	"strings"
)

func (g *GateTradeEngine) apiSpotOrderCreate(req *OrderParam) *mygateapi.PrivateRestSpotOrdersPostAPI {
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestSpotOrdersPost().
		Account(g.gateConverter.ToGateAssetType(AssetType(req.AccountType))).
		CurrencyPair(req.Symbol).Side(g.gateConverter.ToGateOrderSide(req.OrderSide)).Amount(req.Quantity)
	if req.OrderType == ORDER_TYPE_LIMIT {
		api.Type(GATE_ORDER_TYPE_LIMIT).Price(req.Price)
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
func (g *GateTradeEngine) apiFuturesOrderCreate(req *OrderParam) *mygateapi.PrivateRestFuturesSettleOrdersPostAPI {
	split := strings.Split(req.Symbol, "_")
	if len(split) != 2 {
		log.Error("symbol error")
		return nil
	}
	settle := strings.ToLower(split[1])
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestFuturesSettleOrdersPost().
		Settle(settle).Contract(req.Symbol).Size(req.Quantity.IntPart()).Price(req.Price.String())

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

	if req.ClientOrderId != "" {
		api.Text(req.ClientOrderId)
	}

	return api
}

func (g *GateTradeEngine) apiSpotOrderAmend(req *OrderParam) *mygateapi.PrivateRestSpotOrdersOrderIdPatchAPI {
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestSpotOrdersOrderIdPatch().
		OrderId(req.OrderId)

	if req.Symbol != "" {
		api.CurrencyPair(req.Symbol)
	}

	if req.AccountType != "" {
		api.Account(g.gateConverter.ToGateAssetType(AssetType(req.AccountType)))
	}

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

	if req.AccountType != "" {
		api.Account(g.gateConverter.ToGateAssetType(AssetType(req.AccountType)))
	}

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
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestSpotOpenOrders()

	if req.AccountType != "" {
		api.Account(g.gateConverter.ToGateAssetType(AssetType(req.AccountType)))
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
		Settle(settle).Status(GATE_ORDER_STATUS_NEW)

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
		Settle(settle).Status(GATE_ORDER_STATUS_NEW)

	if req.Limit != 0 {
		api.Limit(req.Limit)
	}

	return api
}

func (g *GateTradeEngine) apiSpotOrderQuery(req *QueryOrderParam) *mygateapi.PrivateRestSpotOrdersOrderIdGetAPI {
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestSpotOrdersOrderIdGet()

	if req.AccountType != "" {
		api.Account(g.gateConverter.ToGateAssetType(AssetType(req.AccountType)))
	}

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
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestSpotOrdersGet().
		Status(g.gateConverter.ToGateOrderStatus(req.Status))

	if req.AccountType != "" {
		api.Account(g.gateConverter.ToGateAssetType(AssetType(req.AccountType)))
	}

	if req.Symbol != "" {
		api.CurrencyPair(req.Symbol)
	}

	return api
}
func (g *GateTradeEngine) apiFuturesOrdersQuery(req *QueryOrderParam) *mygateapi.PrivateRestFuturesSettleOrdersGetAPI {
	split := strings.Split(req.Symbol, "_")
	if len(split) != 2 {
		log.Error("symbol error")
		return nil
	}
	settle := strings.ToLower(split[1])
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestFuturesSettleOrdersGet().Settle(settle).
		Status(g.gateConverter.ToGateOrderStatus(req.Status))

	return api
}
func (g *GateTradeEngine) apiDeliveryOrdersQuery(req *QueryOrderParam) *mygateapi.PrivateRestDeliverySettleOrdersGetAPI {
	split := strings.Split(req.Symbol, "_")
	if len(split) != 3 {
		log.Error("symbol error")
		return nil
	}
	settle := strings.ToLower(split[1])
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestDeliverySettleOrdersGet().Settle(settle).
		Status(g.gateConverter.ToGateOrderStatus(req.Status))

	return api
}

func (g *GateTradeEngine) apiSpotTradesQuery(req *QueryTradeParam) *mygateapi.PrivateRestSpotMyTradesAPI {
	api := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestSpotMyTrades()

	if req.AccountType != "" {
		api.Account(g.gateConverter.ToGateAssetType(AssetType(req.AccountType)))
	}

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
