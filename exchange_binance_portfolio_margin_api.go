package mytrade

import (
	"github.com/Hongssd/mybinanceapi"
	"strconv"
)

// UM
func (b *BinanceTradeEngine) apiPortfolioMarginUmOrderCreate(req *OrderParam) *mybinanceapi.PortfolioMarginUmOrderPostApi {
	api := binance.NewPortfolioMarginClient(b.apiKey, b.secretKey).NewUmOrderPost().
		Symbol(req.Symbol).
		Side(b.bnConverter.ToBNOrderSide(req.OrderSide)).
		PositionSide(b.bnConverter.ToBNPositionSide(req.PositionSide)).
		Type(b.bnConverter.ToBNOrderType(req.OrderType)).
		Quantity(req.Quantity)

	if !req.Price.IsZero() {
		api.Price(req.Price)
	}

	if req.TimeInForce != "" {
		api.TimeInForce(b.bnConverter.ToBNTimeInForce(req.TimeInForce))
	}

	if req.ClientOrderId != "" {
		api.NewClientOrderId(req.ClientOrderId)
	}

	return api
}
func (b *BinanceTradeEngine) apiPortfolioMarginUmOrderAmend(req *OrderParam) *mybinanceapi.PortfolioMarginUmOrderPutApi {
	api := binance.NewPortfolioMarginClient(b.apiKey, b.secretKey).NewUmOrderPut().
		Symbol(req.Symbol).
		Side(b.bnConverter.ToBNOrderSide(req.OrderSide)).
		Price(req.Price).
		Quantity(req.Quantity)

	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api.OrderId(orderId)
	}
	if req.ClientOrderId != "" {
		api.OrigClientOrderId(req.ClientOrderId)
	}

	return api
}
func (b *BinanceTradeEngine) apiPortfolioMarginUmOrderCancel(req *OrderParam) *mybinanceapi.PortfolioMarginUmOrderDeleteApi {
	api := binance.NewPortfolioMarginClient(b.apiKey, b.secretKey).NewUmOrderDelete().Symbol(req.Symbol)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api.OrderId(orderId)
	}
	if req.ClientOrderId != "" {
		api.OrigClientOrderId(req.ClientOrderId)
	}

	return api
}

func (b *BinanceTradeEngine) apiPortfolioMarginUmOrderQuery(req *QueryOrderParam) *mybinanceapi.PortfolioMarginUmOrderGetApi {
	api := binance.NewPortfolioMarginClient(b.apiKey, b.secretKey).NewUmOrderGet().
		Symbol(req.Symbol)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api.OrderId(orderId)
	}
	if req.ClientOrderId != "" {
		api.OrigClientOrderId(req.ClientOrderId)
	}

	return api
}
func (b *BinanceTradeEngine) apiPortfolioMarginUmOrdersQuery(req *QueryOrderParam) *mybinanceapi.PortfolioMarginUmAllOrdersGetApi {
	api := binance.NewPortfolioMarginClient(b.apiKey, b.secretKey).NewUmAllOrdersGet().Symbol(req.Symbol)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api.OrderId(orderId)
	}
	if req.StartTime != 0 {
		api.StartTime(req.StartTime)
	}
	if req.EndTime != 0 {
		api.EndTime(req.EndTime)
	}
	if req.Limit != 0 {
		api.Limit(int32(req.Limit))
	}

	return api
}
func (b *BinanceTradeEngine) apiPortfolioMarginUmOpenOrderQuery(req *QueryOrderParam) *mybinanceapi.PortfolioMarginUmOpenOrderGetApi {
	api := binance.NewPortfolioMarginClient(b.apiKey, b.secretKey).NewUmOpenOrderGet().Symbol(req.Symbol)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api.OrderId(orderId)
	}
	if req.ClientOrderId != "" {
		api.OrigClientOrderId(req.ClientOrderId)
	}

	return api
}
func (b *BinanceTradeEngine) apiPortfolioMarginUmOpenOrdersQuery(req *QueryOrderParam) *mybinanceapi.PortfolioMarginUmOpenOrdersGetApi {
	api := binance.NewPortfolioMarginClient(b.apiKey, b.secretKey).NewUmOpenOrdersGet()
	if req.Symbol != "" {
		api.Symbol(req.Symbol)
	}

	return api
}
func (b *BinanceTradeEngine) apiPortfolioMarginUmTradesQuery(req *QueryTradeParam) *mybinanceapi.PortfolioMarginUmUserTradesApi {
	api := binance.NewPortfolioMarginClient(b.apiKey, b.secretKey).NewUmUserTrades().Symbol(req.Symbol)
	if req.StartTime != 0 {
		api.StartTime(req.StartTime)
	}
	if req.EndTime != 0 {
		api.EndTime(req.EndTime)
	}
	if req.Limit != 0 {
		api.Limit(int32(req.Limit))
	}

	return api
}

// CM
func (b *BinanceTradeEngine) apiPortfolioMarginCmOrderCreate(req *OrderParam) *mybinanceapi.PortfolioMarginCmOrderPostApi {
	api := binance.NewPortfolioMarginClient(b.apiKey, b.secretKey).NewCmOrderPost().
		Symbol(req.Symbol).
		Side(b.bnConverter.ToBNOrderSide(req.OrderSide)).
		PositionSide(b.bnConverter.ToBNPositionSide(req.PositionSide)).
		Type(b.bnConverter.ToBNOrderType(req.OrderType)).
		Quantity(req.Quantity)

	if !req.Price.IsZero() {
		api.Price(req.Price)
	}

	if req.TimeInForce != "" {
		api.TimeInForce(b.bnConverter.ToBNTimeInForce(req.TimeInForce))
	}

	if req.ClientOrderId != "" {
		api.NewClientOrderId(req.ClientOrderId)
	}

	return api
}
func (b *BinanceTradeEngine) apiPortfolioMarginCmOrderAmend(req *OrderParam) *mybinanceapi.PortfolioMarginCmOrderPutApi {
	api := binance.NewPortfolioMarginClient(b.apiKey, b.secretKey).NewCmOrderPut().
		Symbol(req.Symbol).
		Side(b.bnConverter.ToBNOrderSide(req.OrderSide)).
		Price(req.Price).
		Quantity(req.Quantity)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api.OrderId(orderId)
	}
	if req.ClientOrderId != "" {
		api.OrigClientOrderId(req.ClientOrderId)
	}

	return api
}
func (b *BinanceTradeEngine) apiPortfolioMarginCmOrderCancel(req *OrderParam) *mybinanceapi.PortfolioMarginCmOrderDeleteApi {
	api := binance.NewPortfolioMarginClient(b.apiKey, b.secretKey).NewCmOrderDelete().Symbol(req.Symbol)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api.OrderId(orderId)
	}
	if req.ClientOrderId != "" {
		api.OrigClientOrderId(req.ClientOrderId)
	}

	return api
}

func (b *BinanceTradeEngine) apiPortfolioMarginCmOrderQuery(req *QueryOrderParam) *mybinanceapi.PortfolioMarginCmOrderGetApi {
	api := binance.NewPortfolioMarginClient(b.apiKey, b.secretKey).NewCmOrderGet().
		Symbol(req.Symbol)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api.OrderId(orderId)
	}
	if req.ClientOrderId != "" {
		api.OrigClientOrderId(req.ClientOrderId)
	}

	return api
}
func (b *BinanceTradeEngine) apiPortfolioMarginCmOrdersQuery(req *QueryOrderParam) *mybinanceapi.PortfolioMarginCmAllOrdersApi {
	api := binance.NewPortfolioMarginClient(b.apiKey, b.secretKey).NewCmAllOrders().Symbol(req.Symbol)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api.OrderId(orderId)
	}
	if req.StartTime != 0 {
		api.StartTime(req.StartTime)
	}
	if req.EndTime != 0 {
		api.EndTime(req.EndTime)
	}
	if req.Limit != 0 {
		api.Limit(int32(req.Limit))
	}

	return api
}
func (b *BinanceTradeEngine) apiPortfolioMarginCmOpenOrdersQuery(req *QueryOrderParam) *mybinanceapi.PortfolioMarginCmOpenOrdersApi {
	api := binance.NewPortfolioMarginClient(b.apiKey, b.secretKey).NewCmOpenOrders()
	if req.Symbol != "" {
		api.Symbol(req.Symbol)
	}
	//api.Pair()

	return api
}
func (b *BinanceTradeEngine) apiPortfolioMarginCmTradesQuery(req *QueryTradeParam) *mybinanceapi.PortfolioMarginCmUserTradesApi {
	api := binance.NewPortfolioMarginClient(b.apiKey, b.secretKey).NewCmUserTrades().Symbol(req.Symbol)
	if req.StartTime != 0 {
		api.StartTime(req.StartTime)
	}
	if req.EndTime != 0 {
		api.EndTime(req.EndTime)
	}
	if req.Limit != 0 {
		api.Limit(int32(req.Limit))
	}

	return api
}

// Margin
func (b *BinanceTradeEngine) apiPortfolioMarginMarginOrderCreate(req *OrderParam) *mybinanceapi.PortfolioMarginMarginOrderPostApi {
	api := binance.NewPortfolioMarginClient(b.apiKey, b.secretKey).NewMarginOrderPost().
		Symbol(req.Symbol).
		Side(b.bnConverter.ToBNOrderSide(req.OrderSide)).
		Type(b.bnConverter.ToBNOrderType(req.OrderType))

	if !req.Quantity.IsZero() {
		api.Quantity(req.Quantity)
	}

	if !req.Price.IsZero() {
		api.Price(req.Price)
	}

	if !req.TriggerPrice.IsZero() {
		api.StopPrice(req.TriggerPrice)
	}

	if req.ClientOrderId != "" {
		api.NewClientOrderId(req.ClientOrderId)
	}

	if req.TimeInForce != "" {
		api.TimeInForce(b.bnConverter.ToBNTimeInForce(req.TimeInForce))
	}

	return api
}
func (b *BinanceTradeEngine) apiPortfolioMarginMarginOrderCancel(req *OrderParam) *mybinanceapi.PortfolioMarginMarginOrderDeleteApi {
	api := binance.NewPortfolioMarginClient(b.apiKey, b.secretKey).NewMarginOrderDelete().Symbol(req.Symbol)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api.OrderId(orderId)
	}
	if req.ClientOrderId != "" {
		api.OrigClientOrderId(req.ClientOrderId)
	}

	return api
}

func (b *BinanceTradeEngine) apiPortfolioMarginMarginOrderQuery(req *QueryOrderParam) *mybinanceapi.PortfolioMarginMarginOrderGetApi {
	api := binance.NewPortfolioMarginClient(b.apiKey, b.secretKey).NewMarginOrderGet().Symbol(req.Symbol)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api.OrderId(orderId)
	}
	if req.ClientOrderId != "" {
		api.OrigClientOrderId(req.ClientOrderId)
	}

	return api
}
func (b *BinanceTradeEngine) apiPortfolioMarginMarginOrdersQuery(req *QueryOrderParam) *mybinanceapi.PortfolioMarginMarginAllOrdersApi {
	api := binance.NewPortfolioMarginClient(b.apiKey, b.secretKey).NewMarginAllOrders().Symbol(req.Symbol)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api.OrderId(orderId)
	}
	if req.StartTime != 0 {
		api.StartTime(req.StartTime)
	}
	if req.EndTime != 0 {
		api.EndTime(req.EndTime)
	}
	if req.Limit != 0 {
		api.Limit(int32(req.Limit))
	}

	return api
}
func (b *BinanceTradeEngine) apiPortfolioMarginMarginOpenOrdersQuery(req *QueryOrderParam) *mybinanceapi.PortfolioMarginMarginOpenOrdersApi {
	api := binance.NewPortfolioMarginClient(b.apiKey, b.secretKey).NewMarginOpenOrders()
	if req.Symbol != "" {
		api.Symbol(req.Symbol)
	}
	return api
}
func (b *BinanceTradeEngine) apiPortfolioMarginMarginTradesQuery(req *QueryTradeParam) *mybinanceapi.PortfolioMarginMarginMyTradesApi {
	api := binance.NewPortfolioMarginClient(b.apiKey, b.secretKey).NewMarginMyTrades().Symbol(req.Symbol)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api.OrderId(orderId)
	}
	if req.StartTime != 0 {
		api.StartTime(req.StartTime)
	}
	if req.EndTime != 0 {
		api.EndTime(req.EndTime)
	}
	if req.Limit != 0 {
		api.Limit(int32(req.Limit))
	}

	return api
}
