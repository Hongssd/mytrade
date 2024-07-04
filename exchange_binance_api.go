package mytrade

import (
	"github.com/Hongssd/mybinanceapi"
	"github.com/shopspring/decimal"
	"strconv"
	"time"
)

// 现货订单API接口
func (b *BinanceTradeEngine) apiSpotOpenOrders(req *QueryOrderParam) *mybinanceapi.SpotOpenOrdersApi {
	api := binance.NewSpotRestClient(b.apiKey, b.secretKey).NewOpenOrders()
	if req.Symbol != "" {
		api.Symbol(req.Symbol)
	}
	return api
}
func (b *BinanceTradeEngine) apiSpotOrderQuery(req *QueryOrderParam) *mybinanceapi.SpotOrderGetApi {
	api := binance.NewSpotRestClient(b.apiKey, b.secretKey).NewSpotOrderGet().Symbol(req.Symbol)
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
func (b *BinanceTradeEngine) apiSpotOrdersQuery(req *QueryOrderParam) *mybinanceapi.SpotAllOrdersApi {
	api := binance.NewSpotRestClient(b.apiKey, b.secretKey).NewAllOrders().Symbol(req.Symbol)
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
func (b *BinanceTradeEngine) apiSpotTradeQuery(req *QueryTradeParam) *mybinanceapi.SpotMyTradesApi {
	api := binance.NewSpotRestClient(b.apiKey, b.secretKey).NewSpotMyTrades().
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
		api = api.Limit(req.Limit)
	}
	return api
}

func (b *BinanceTradeEngine) apiSpotOrderCreate(req *OrderParam) *mybinanceapi.SpotOrderPostApi {
	api := binance.NewSpotRestClient(b.apiKey, b.secretKey).NewSpotOrderPost().
		Symbol(req.Symbol).
		Type(b.bnConverter.ToBNOrderType(req.OrderType)).
		Side(b.bnConverter.ToBNOrderSide(req.OrderSide)).
		Quantity(req.Quantity)
	if !req.Price.IsZero() {
		api = api.Price(req.Price)
	}
	if req.ClientOrderId != "" {
		api = api.NewClientOrderId(req.ClientOrderId)
	}
	if req.TimeInForce != "" {
		api = api.TimeInForce(b.bnConverter.ToBNTimeInForce(req.TimeInForce))
	}
	return api
}
func (b *BinanceTradeEngine) apiSpotOrderAmend(req *OrderParam) *mybinanceapi.SpotOrderCancelReplaceApi {
	api := binance.NewSpotRestClient(b.apiKey, b.secretKey).NewSpotOrderCancelReplace().
		Symbol(req.Symbol).CancelReplaceMode("STOP_ON_FAILURE").
		Side(b.bnConverter.ToBNOrderSide(req.OrderSide)).
		Type(b.bnConverter.ToBNOrderType(req.OrderType)).
		Quantity(req.Quantity)
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
		api = api.TimeInForce(b.bnConverter.ToBNTimeInForce(req.TimeInForce))
	}
	return api
}
func (b *BinanceTradeEngine) apiSpotOrderCancel(req *OrderParam) *mybinanceapi.SpotOrderDeleteApi {
	api := binance.NewSpotRestClient(b.apiKey, b.secretKey).NewSpotOrderDelete().
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
func (b *BinanceTradeEngine) apiFutureOpenOrders(req *QueryOrderParam) *mybinanceapi.FutureOpenOrdersApi {
	api := binance.NewFutureRestClient(b.apiKey, b.secretKey).NewOpenOrders()
	if req.Symbol != "" {
		api.Symbol(req.Symbol)
	}
	return api
}
func (b *BinanceTradeEngine) apiFutureOrderQuery(req *QueryOrderParam) *mybinanceapi.FutureOrderGetApi {
	api := binance.NewFutureRestClient(b.apiKey, b.secretKey).NewFutureOrderGet().Symbol(req.Symbol)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api = api.OrderId(orderId)
	} else {
		api = api.OrigClientOrderId(req.ClientOrderId)
	}
	return api
}
func (b *BinanceTradeEngine) apiFutureOrdersQuery(req *QueryOrderParam) *mybinanceapi.FutureAllOrdersApi {
	api := binance.NewFutureRestClient(b.apiKey, b.secretKey).NewAllOrders().Symbol(req.Symbol)

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
func (b *BinanceTradeEngine) apiFutureTradeQuery(req *QueryTradeParam) *mybinanceapi.FutureUserTradesApi {
	api := binance.NewFutureRestClient(b.apiKey, b.secretKey).NewFutureUserTrades().
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

func (b *BinanceTradeEngine) apiFutureOrderCreate(req *OrderParam) *mybinanceapi.FutureOrderPostApi {
	api := binance.NewFutureRestClient(b.apiKey, b.secretKey).NewFutureOrderPost().
		Symbol(req.Symbol).
		Type(b.bnConverter.ToBNOrderType(req.OrderType)).
		Side(b.bnConverter.ToBNOrderSide(req.OrderSide)).
		PositionSide(b.bnConverter.ToBNPositionSide(req.PositionSide)).
		Quantity(req.Quantity)
	if !req.Price.IsZero() {
		api = api.Price(req.Price)
	}
	if req.ClientOrderId != "" {
		api = api.NewClientOrderId(req.ClientOrderId)
	}
	if req.TimeInForce != "" {
		api = api.TimeInForce(b.bnConverter.ToBNTimeInForce(req.TimeInForce))
	}
	return api
}
func (b *BinanceTradeEngine) apiFutureOrderAmend(req *OrderParam) *mybinanceapi.FutureOrderPutApi {
	api := binance.NewFutureRestClient(b.apiKey, b.secretKey).NewFutureOrderPut().
		Symbol(req.Symbol).
		Side(b.bnConverter.ToBNOrderSide(req.OrderSide)).
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
func (b *BinanceTradeEngine) apiFutureOrderCancel(req *OrderParam) *mybinanceapi.FutureOrderDeleteApi {
	api := binance.NewFutureRestClient(b.apiKey, b.secretKey).NewFutureOrderDelete().
		Symbol(req.Symbol)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api = api.OrderId(orderId)
	} else {
		api = api.OrigClientOrderId(req.ClientOrderId)
	}
	return api
}

func (b *BinanceTradeEngine) apiFutureBatchOrderCreate(reqs []*OrderParam) *mybinanceapi.FutureBatchOrdersPostApi {
	client := binance.NewFutureRestClient(b.apiKey, b.secretKey)
	api := client.NewFutureBatchOrdersPost()
	for _, req := range reqs {
		thisApi := client.NewFutureOrderPost().Symbol(req.Symbol).
			Side(b.bnConverter.ToBNOrderSide(req.OrderSide)).
			Type(b.bnConverter.ToBNOrderType(req.OrderType)).
			PositionSide(b.bnConverter.ToBNPositionSide(req.PositionSide)).
			Quantity(req.Quantity)
		if !req.Price.IsZero() {
			thisApi = thisApi.Price(req.Price)
		}
		if req.ClientOrderId != "" {
			thisApi = thisApi.NewClientOrderId(req.ClientOrderId)
		}
		if req.TimeInForce != "" {
			thisApi = thisApi.TimeInForce(b.bnConverter.ToBNTimeInForce(req.TimeInForce))
		}
		api = api.AddOrders(thisApi)
	}
	return api
}
func (b *BinanceTradeEngine) apiFutureBatchOrderAmend(reqs []*OrderParam) *mybinanceapi.FutureBatchOrdersPutApi {
	client := binance.NewFutureRestClient(b.apiKey, b.secretKey)
	api := client.NewFutureBatchOrdersPut()
	for _, req := range reqs {
		thisApi := client.NewFutureOrderPut().Symbol(req.Symbol).
			Side(b.bnConverter.ToBNOrderSide(req.OrderSide)).
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
func (b *BinanceTradeEngine) apiFutureBatchOrderCancel(reqs []*OrderParam) (*mybinanceapi.FutureBatchOrdersDeleteApi, error) {
	client := binance.NewFutureRestClient(b.apiKey, b.secretKey)
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

// 币本位合约订单API接口
func (b *BinanceTradeEngine) apiSwapOpenOrders(req *QueryOrderParam) *mybinanceapi.SwapOpenOrdersApi {
	api := binance.NewSwapRestClient(b.apiKey, b.secretKey).NewOpenOrders()
	if req.Symbol != "" {
		api.Symbol(req.Symbol)
	}
	return api
}
func (b *BinanceTradeEngine) apiSwapOrderQuery(req *QueryOrderParam) *mybinanceapi.SwapOrderGetApi {
	api := binance.NewSwapRestClient(b.apiKey, b.secretKey).NewSwapOrderGet().Symbol(req.Symbol)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api = api.OrderId(orderId)
	} else {
		api = api.OrigClientOrderId(req.ClientOrderId)
	}
	return api
}
func (b *BinanceTradeEngine) apiSwapOrdersQuery(req *QueryOrderParam) *mybinanceapi.SwapAllOrdersApi {
	api := binance.NewSwapRestClient(b.apiKey, b.secretKey).NewAllOrders().Symbol(req.Symbol)
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
func (b *BinanceTradeEngine) apiSwapTradeQuery(req *QueryTradeParam) *mybinanceapi.SwapUserTradesApi {
	api := binance.NewSwapRestClient(b.apiKey, b.secretKey).NewSwapUserTrades().
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

func (b *BinanceTradeEngine) apiSwapOrderCreate(req *OrderParam) *mybinanceapi.SwapOrderPostApi {
	api := binance.NewSwapRestClient(b.apiKey, b.secretKey).NewSwapOrderPost().
		Symbol(req.Symbol).
		Type(b.bnConverter.ToBNOrderType(req.OrderType)).
		Side(b.bnConverter.ToBNOrderSide(req.OrderSide)).
		PositionSide(b.bnConverter.ToBNPositionSide(req.PositionSide)).
		Quantity(req.Quantity)
	if !req.Price.IsZero() {
		api = api.Price(req.Price)
	}
	if req.ClientOrderId != "" {
		api = api.NewClientOrderId(req.ClientOrderId)
	}
	if req.TimeInForce != "" {
		api = api.TimeInForce(b.bnConverter.ToBNTimeInForce(req.TimeInForce))
	}
	return api
}
func (b *BinanceTradeEngine) apiSwapOrderAmend(req *OrderParam) *mybinanceapi.SwapOrderPutApi {
	api := binance.NewSwapRestClient(b.apiKey, b.secretKey).NewSwapOrderPut().
		Symbol(req.Symbol).
		Side(b.bnConverter.ToBNOrderSide(req.OrderSide)).
		Quantity(req.Quantity)
	if !req.Price.IsZero() {
		api = api.Price(req.Price)
	}
	if req.OrderId != "" {
		api = api.OrderId(req.OrderId)
	}
	if req.ClientOrderId != "" {
		api = api.OrigClientOrderId(req.ClientOrderId)
	}
	return api
}
func (b *BinanceTradeEngine) apiSwapOrderCancel(req *OrderParam) *mybinanceapi.SwapOrderDeleteApi {
	api := binance.NewSwapRestClient(b.apiKey, b.secretKey).NewSwapOrderDelete().
		Symbol(req.Symbol)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api = api.OrderId(orderId)
	} else {
		api = api.OrigClientOrderId(req.ClientOrderId)
	}
	return api
}

func (b *BinanceTradeEngine) apiSwapBatchOrderCreate(reqs []*OrderParam) *mybinanceapi.SwapBatchOrdersPostApi {
	client := binance.NewSwapRestClient(b.apiKey, b.secretKey)
	api := client.NewSwapBatchOrdersPost()
	for _, req := range reqs {
		thisApi := client.NewSwapOrderPost().Symbol(req.Symbol).
			Side(b.bnConverter.ToBNOrderSide(req.OrderSide)).
			Type(b.bnConverter.ToBNOrderType(req.OrderType)).
			PositionSide(b.bnConverter.ToBNPositionSide(req.PositionSide)).
			Quantity(req.Quantity)
		if !req.Price.IsZero() {
			thisApi = thisApi.Price(req.Price)
		}
		if req.ClientOrderId != "" {
			thisApi = thisApi.NewClientOrderId(req.ClientOrderId)
		}
		if req.TimeInForce != "" {
			thisApi = thisApi.TimeInForce(b.bnConverter.ToBNTimeInForce(req.TimeInForce))
		}
		api = api.AddOrders(thisApi)
	}
	return api
}
func (b *BinanceTradeEngine) apiSwapBatchOrderAmend(reqs []*OrderParam) *mybinanceapi.SwapBatchOrdersPutApi {
	client := binance.NewSwapRestClient(b.apiKey, b.secretKey)
	api := client.NewSwapBatchOrdersPut()
	for _, req := range reqs {
		thisApi := client.NewSwapOrderPut().Symbol(req.Symbol).
			Side(b.bnConverter.ToBNOrderSide(req.OrderSide)).
			Quantity(req.Quantity)
		if !req.Price.IsZero() {
			thisApi = thisApi.Price(req.Price)
		}
		if req.OrderId != "" {
			thisApi = thisApi.OrderId(req.OrderId)
		}
		if req.ClientOrderId != "" {
			thisApi = thisApi.OrigClientOrderId(req.ClientOrderId)
		}
		api = api.AddOrders(thisApi)
	}
	return api
}
func (b *BinanceTradeEngine) apiSwapBatchOrderCancel(reqs []*OrderParam) (*mybinanceapi.SwapBatchOrdersDeleteApi, error) {
	client := binance.NewSwapRestClient(b.apiKey, b.secretKey)
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
	api := client.NewSwapBatchOrdersDelete().
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

// 现货订单处理
func (b *BinanceTradeEngine) handleOrdersFromSpotOpenOrders(req *QueryOrderParam, res *mybinanceapi.SpotOpenOrdersRes) []*Order {
	var orders []*Order

	for _, order := range *res {
		avgPrice := decimal.Zero
		if order.ExecutedQty != "" && order.CummulativeQuoteQty != "" {
			executedQty, _ := decimal.NewFromString(order.ExecutedQty)
			cumQuoteQty, _ := decimal.NewFromString(order.CummulativeQuoteQty)
			if !executedQty.IsZero() {
				avgPrice = cumQuoteQty.Div(executedQty)
			}
		}
		orders = append(orders, &Order{
			Exchange:      BINANCE_NAME.String(),
			AccountType:   req.AccountType,
			Symbol:        order.Symbol,
			OrderId:       strconv.FormatInt(order.OrderId, 10),
			ClientOrderId: order.ClientOrderId,
			Price:         order.Price,
			Quantity:      order.OrigQty,
			ExecutedQty:   order.ExecutedQty,
			CumQuoteQty:   order.CummulativeQuoteQty,
			AvgPrice:      avgPrice.String(),
			Status:        b.bnConverter.FromBNOrderStatus(order.Status),
			Type:          b.bnConverter.FromBNOrderType(order.Type),
			Side:          b.bnConverter.FromBNOrderSide(order.Side),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
			CreateTime:    order.Time,
			UpdateTime:    order.UpdateTime,
		})
	}
	return orders
}
func (b *BinanceTradeEngine) handleOrderFromSpotOrderQuery(req *QueryOrderParam, res *mybinanceapi.SpotOrderGetRes) *Order {
	avgPrice := decimal.Zero
	if res.ExecutedQty != "" && res.CummulativeQuoteQty != "" {
		executedQty, _ := decimal.NewFromString(res.ExecutedQty)
		cumQuoteQty, _ := decimal.NewFromString(res.CummulativeQuoteQty)
		if !executedQty.IsZero() {
			avgPrice = cumQuoteQty.Div(executedQty)
		}
	}
	order := &Order{
		Exchange:      BINANCE_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.ClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CummulativeQuoteQty,
		AvgPrice:      avgPrice.String(),
		Status:        b.bnConverter.FromBNOrderStatus(res.Status),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		CreateTime:    res.Time,
		UpdateTime:    res.UpdateTime,
	}
	return order
}
func (b *BinanceTradeEngine) handleOrderFromSpotOrdersQuery(req *QueryOrderParam, res *mybinanceapi.SpotAllOrdersRes) []*Order {
	var orders []*Order
	for _, order := range *res {
		avgPrice := decimal.Zero
		if order.ExecutedQty != "" && order.CummulativeQuoteQty != "" {
			executedQty, _ := decimal.NewFromString(order.ExecutedQty)
			cumQuoteQty, _ := decimal.NewFromString(order.CummulativeQuoteQty)
			if !executedQty.IsZero() {
				avgPrice = cumQuoteQty.Div(executedQty)
			}
		}
		orders = append(orders, &Order{
			Exchange:      BINANCE_NAME.String(),
			AccountType:   req.AccountType,
			Symbol:        req.Symbol,
			OrderId:       strconv.FormatInt(order.OrderId, 10),
			ClientOrderId: order.ClientOrderId,
			Price:         order.Price,
			Quantity:      order.OrigQty,
			ExecutedQty:   order.ExecutedQty,
			CumQuoteQty:   order.CummulativeQuoteQty,
			AvgPrice:      avgPrice.String(),
			Status:        b.bnConverter.FromBNOrderStatus(order.Status),
			Type:          b.bnConverter.FromBNOrderType(order.Type),
			Side:          b.bnConverter.FromBNOrderSide(order.Side),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
			CreateTime:    order.Time,
			UpdateTime:    order.UpdateTime,
		})
	}
	return orders
}
func (b *BinanceTradeEngine) handleTradesFromSpotTradeQuery(req *QueryTradeParam, res *mybinanceapi.SpotMyTradesRes) []*Trade {
	var trades []*Trade
	for _, trade := range *res {
		var orderSide OrderSide
		if trade.IsBuyer {
			orderSide = ORDER_SIDE_BUY
		} else {
			orderSide = ORDER_SIDE_SELL
		}
		trades = append(trades, &Trade{
			Exchange:    BINANCE_NAME.String(),
			AccountType: req.AccountType,
			Symbol:      req.Symbol,
			TradeId:     strconv.FormatInt(trade.Id, 10),
			OrderId:     strconv.FormatInt(trade.OrderId, 10),
			Price:       trade.Price,
			Quantity:    trade.Qty,
			QuoteQty:    trade.QuoteQty,
			Side:        orderSide,
			FeeAmount:   trade.Commission,
			FeeCcy:      trade.CommissionAsset,
			IsMaker:     trade.IsMaker,
			Timestamp:   trade.Time,
		})
	}
	return trades
}

func (b *BinanceTradeEngine) handleOrderFromSpotOrderCreate(req *OrderParam, res *mybinanceapi.SpotOrderPostRes) *Order {
	avgPrice := decimal.Zero
	if res.ExecutedQty != "" && res.CummulativeQuoteQty != "" {
		executedQty, _ := decimal.NewFromString(res.ExecutedQty)
		cumQuoteQty, _ := decimal.NewFromString(res.CummulativeQuoteQty)
		if !executedQty.IsZero() {
			avgPrice = cumQuoteQty.Div(executedQty)
		}
	}
	order := &Order{
		Exchange:      BINANCE_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.ClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CummulativeQuoteQty,
		AvgPrice:      avgPrice.String(),
		Status:        b.bnConverter.FromBNOrderStatus(res.Status),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		CreateTime:    res.WorkingTime,
		UpdateTime:    res.WorkingTime,
	}
	return order
}
func (b *BinanceTradeEngine) handleOrderFromSpotOrderAmend(req *OrderParam, res *mybinanceapi.SpotOrderCancelReplaceRes) *Order {
	avgPrice := decimal.Zero
	if res.NewOrderResponse.ExecutedQty != "" && res.NewOrderResponse.CummulativeQuoteQty != "" {
		executedQty, _ := decimal.NewFromString(res.NewOrderResponse.ExecutedQty)
		cumQuoteQty, _ := decimal.NewFromString(res.NewOrderResponse.CummulativeQuoteQty)
		if !executedQty.IsZero() {
			avgPrice = cumQuoteQty.Div(executedQty)
		}
	}
	order := &Order{
		Exchange:      BINANCE_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		OrderId:       strconv.FormatInt(res.NewOrderResponse.OrderId, 10),
		ClientOrderId: res.NewOrderResponse.ClientOrderId,
		Price:         res.NewOrderResponse.Price,
		Quantity:      res.NewOrderResponse.OrigQty,
		ExecutedQty:   res.NewOrderResponse.ExecutedQty,
		CumQuoteQty:   res.NewOrderResponse.CummulativeQuoteQty,
		AvgPrice:      avgPrice.String(),
		Status:        b.bnConverter.FromBNOrderStatus(res.NewOrderResponse.Status),
		Type:          b.bnConverter.FromBNOrderType(res.NewOrderResponse.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.NewOrderResponse.Side),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.NewOrderResponse.TimeInForce),
		CreateTime:    res.NewOrderResponse.WorkingTime,
		UpdateTime:    res.NewOrderResponse.WorkingTime,
	}
	return order
}
func (b *BinanceTradeEngine) handleOrderFromSpotOrderCancel(req *OrderParam, res *mybinanceapi.SpotOrderDeleteRes) *Order {
	nowTimestamp := time.Now().UnixMilli()
	avgPrice := decimal.Zero
	if res.ExecutedQty != "" && res.CummulativeQuoteQty != "" {
		executedQty, _ := decimal.NewFromString(res.ExecutedQty)
		cumQuoteQty, _ := decimal.NewFromString(res.CummulativeQuoteQty)
		if !executedQty.IsZero() {
			avgPrice = cumQuoteQty.Div(executedQty)
		}
	}
	order := &Order{
		Exchange:      BINANCE_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.OrigClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CummulativeQuoteQty,
		AvgPrice:      avgPrice.String(),
		Status:        b.bnConverter.FromBNOrderStatus(res.Status),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		CreateTime:    res.TransactTime,
		UpdateTime:    nowTimestamp,
	}
	return order
}

func (b *BinanceTradeEngine) handleOrderFromSpotBatchErr(req *OrderParam, err error) *Order {
	return &Order{
		Exchange:      BINANCE_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		OrderId:       req.OrderId,
		ClientOrderId: req.ClientOrderId,
		Price:         req.Price.String(),
		Quantity:      req.Quantity.String(),
		Type:          req.OrderType,
		Side:          req.OrderSide,
		PositionSide:  req.PositionSide,
		TimeInForce:   req.TimeInForce,
		Status:        ORDER_STATUS_REJECTED,
		ErrorMsg:      err.Error(),
	}
}

// U合约订单处理
func (b *BinanceTradeEngine) handleOrdersFromFutureOpenOrders(req *QueryOrderParam, res *mybinanceapi.FutureOpenOrdersRes) []*Order {
	var orders []*Order
	for _, order := range *res {
		orders = append(orders, &Order{
			Exchange:      BINANCE_NAME.String(),
			AccountType:   req.AccountType,
			Symbol:        order.Symbol,
			OrderId:       strconv.FormatInt(order.OrderId, 10),
			ClientOrderId: order.ClientOrderId,
			Price:         order.Price,
			Quantity:      order.OrigQty,
			ExecutedQty:   order.ExecutedQty,
			CumQuoteQty:   order.CumQuote,
			AvgPrice:      order.AvgPrice,
			Status:        b.bnConverter.FromBNOrderStatus(order.Status),
			Type:          b.bnConverter.FromBNOrderType(order.Type),
			Side:          b.bnConverter.FromBNOrderSide(order.Side),
			PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    order.Time,
			UpdateTime:    order.UpdateTime,
		})
	}
	return orders
}
func (b *BinanceTradeEngine) handleOrderFromFutureOrderQuery(req *QueryOrderParam, res *mybinanceapi.FutureOrderGetRes) *Order {
	order := &Order{
		Exchange:      BINANCE_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.ClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CumQuote,
		AvgPrice:      res.AvgPrice,
		Status:        b.bnConverter.FromBNOrderStatus(res.Status),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		PositionSide:  b.bnConverter.FromBNPositionSide(res.PositionSide),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		ReduceOnly:    res.ReduceOnly,
		CreateTime:    res.Time,
		UpdateTime:    res.UpdateTime,
	}
	return order
}
func (b *BinanceTradeEngine) handleOrderFromFutureOrdersQuery(req *QueryOrderParam, res *mybinanceapi.FutureAllOrdersRes) []*Order {
	var orders []*Order
	for _, order := range *res {
		orders = append(orders, &Order{
			Exchange:      BINANCE_NAME.String(),
			AccountType:   req.AccountType,
			Symbol:        order.Symbol,
			OrderId:       strconv.FormatInt(order.OrderId, 10),
			ClientOrderId: order.ClientOrderId,
			Price:         order.Price,
			Quantity:      order.OrigQty,
			ExecutedQty:   order.ExecutedQty,
			CumQuoteQty:   order.CumQuote,
			AvgPrice:      order.AvgPrice,
			Status:        b.bnConverter.FromBNOrderStatus(order.Status),
			Type:          b.bnConverter.FromBNOrderType(order.Type),
			Side:          b.bnConverter.FromBNOrderSide(order.Side),
			PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    order.Time,
			UpdateTime:    order.UpdateTime,
		})
	}
	return orders
}
func (b *BinanceTradeEngine) handleTradesFromFutureTradeQuery(req *QueryTradeParam, res *mybinanceapi.FutureUserTradesRes) []*Trade {
	var trades []*Trade
	for _, trade := range *res {
		trades = append(trades, &Trade{
			Exchange:     BINANCE_NAME.String(),
			AccountType:  req.AccountType,
			Symbol:       req.Symbol,
			TradeId:      strconv.FormatInt(trade.Id, 10),
			OrderId:      strconv.FormatInt(trade.OrderId, 10),
			Price:        trade.Price,
			Quantity:     trade.Qty,
			QuoteQty:     trade.QuoteQty,
			Side:         b.bnConverter.FromBNOrderSide(trade.Side),
			PositionSide: b.bnConverter.FromBNPositionSide(trade.PositionSide),
			FeeAmount:    trade.Commission,
			FeeCcy:       trade.CommissionAsset,
			RealizedPnl:  trade.RealizedPnl,
			IsMaker:      trade.Maker,
			Timestamp:    trade.Time,
		})
	}
	return trades
}

func (b *BinanceTradeEngine) handleOrderFromFutureOrderCreate(req *OrderParam, res *mybinanceapi.FutureOrderPostRes) *Order {
	nowTimestamp := time.Now().UnixMilli()
	order := &Order{
		Exchange:      BINANCE_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.ClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CumQuote,
		AvgPrice:      res.AvgPrice,
		Status:        b.bnConverter.FromBNOrderStatus(res.Status),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		PositionSide:  b.bnConverter.FromBNPositionSide(res.PositionSide),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		ReduceOnly:    res.ReduceOnly,
		CreateTime:    nowTimestamp,
		UpdateTime:    nowTimestamp,
	}
	return order
}
func (b *BinanceTradeEngine) handleOrderFromFutureOrderAmend(req *OrderParam, res *mybinanceapi.FutureOrderPutRes) *Order {
	nowTimestamp := time.Now().UnixMilli()
	order := &Order{
		Exchange:      BINANCE_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.ClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CumQuote,
		AvgPrice:      res.AvgPrice,
		Status:        b.bnConverter.FromBNOrderStatus(res.Status),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		PositionSide:  b.bnConverter.FromBNPositionSide(res.PositionSide),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		ReduceOnly:    res.ReduceOnly,
		CreateTime:    nowTimestamp,
		UpdateTime:    nowTimestamp,
	}
	return order
}
func (b *BinanceTradeEngine) handleOrderFromFutureOrderCancel(req *OrderParam, res *mybinanceapi.FutureOrderDeleteRes) *Order {
	order := &Order{
		Exchange:      BINANCE_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.ClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CumQuote,
		AvgPrice:      res.AvgPrice,
		Status:        b.bnConverter.FromBNOrderStatus(res.Status),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		PositionSide:  b.bnConverter.FromBNPositionSide(res.PositionSide),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		ReduceOnly:    res.ReduceOnly,
		CreateTime:    res.UpdateTime,
		UpdateTime:    res.UpdateTime,
	}
	return order
}

func (b *BinanceTradeEngine) handleOrdersFromFutureBatchOrderCreate(reqs []*OrderParam, res *mybinanceapi.FutureBatchOrdersPostRes) []*Order {
	var orders []*Order
	nowTimestamp := time.Now().UnixMilli()
	for _, order := range *res {
		orders = append(orders, &Order{
			Exchange:      BINANCE_NAME.String(),
			AccountType:   reqs[0].AccountType,
			Symbol:        order.Symbol,
			OrderId:       strconv.FormatInt(order.OrderId, 10),
			ClientOrderId: order.ClientOrderId,
			Price:         order.Price,
			Quantity:      order.OrigQty,
			ExecutedQty:   order.ExecutedQty,
			CumQuoteQty:   order.CumQuote,
			AvgPrice:      order.AvgPrice,
			Status:        b.bnConverter.FromBNOrderStatus(order.Status),
			Type:          b.bnConverter.FromBNOrderType(order.Type),
			Side:          b.bnConverter.FromBNOrderSide(order.Side),
			PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    nowTimestamp,
			UpdateTime:    nowTimestamp,
			ErrorCode:     strconv.Itoa(order.Code),
			ErrorMsg:      order.Msg,
		})
	}
	return orders
}
func (b *BinanceTradeEngine) handleOrdersFromFutureBatchOrderAmend(reqs []*OrderParam, res *mybinanceapi.FutureBatchOrdersPutRes) []*Order {
	var orders []*Order
	nowTimestamp := time.Now().UnixMilli()
	for _, order := range *res {
		orders = append(orders, &Order{
			Exchange:      BINANCE_NAME.String(),
			AccountType:   reqs[0].AccountType,
			Symbol:        order.Symbol,
			OrderId:       strconv.FormatInt(order.OrderId, 10),
			ClientOrderId: order.ClientOrderId,
			Price:         order.Price,
			Quantity:      order.OrigQty,
			ExecutedQty:   order.ExecutedQty,
			CumQuoteQty:   order.CumQuote,
			AvgPrice:      order.AvgPrice,
			Status:        b.bnConverter.FromBNOrderStatus(order.Status),
			Type:          b.bnConverter.FromBNOrderType(order.Type),
			Side:          b.bnConverter.FromBNOrderSide(order.Side),
			PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    nowTimestamp,
			UpdateTime:    nowTimestamp,
			ErrorCode:     strconv.Itoa(order.Code),
			ErrorMsg:      order.Msg,
		})
	}
	return orders
}
func (b *BinanceTradeEngine) handleOrdersFromFutureBatchOrderCancel(reqs []*OrderParam, res *mybinanceapi.FutureBatchOrdersDeleteRes) []*Order {
	var orders []*Order
	nowTimestamp := time.Now().UnixMilli()
	for _, order := range *res {
		orders = append(orders, &Order{
			Exchange:      BINANCE_NAME.String(),
			AccountType:   reqs[0].AccountType,
			Symbol:        order.Symbol,
			OrderId:       strconv.FormatInt(order.OrderId, 10),
			ClientOrderId: order.ClientOrderId,
			Price:         order.Price,
			Quantity:      order.OrigQty,
			ExecutedQty:   order.ExecutedQty,
			CumQuoteQty:   order.CumQuote,
			AvgPrice:      order.AvgPrice,
			Status:        b.bnConverter.FromBNOrderStatus(order.Status),
			Type:          b.bnConverter.FromBNOrderType(order.Type),
			Side:          b.bnConverter.FromBNOrderSide(order.Side),
			PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    nowTimestamp,
			UpdateTime:    nowTimestamp,
			ErrorCode:     strconv.Itoa(order.Code),
			ErrorMsg:      order.Msg,
		})
	}
	return orders
}

// 币本位合约订单处理
func (b *BinanceTradeEngine) handleOrdersFromSwapOpenOrders(req *QueryOrderParam, res *mybinanceapi.SwapOpenOrdersRes) []*Order {
	var orders []*Order
	for _, order := range *res {
		orders = append(orders, &Order{
			Exchange:      BINANCE_NAME.String(),
			AccountType:   req.AccountType,
			Symbol:        order.Symbol,
			OrderId:       strconv.FormatInt(order.OrderId, 10),
			ClientOrderId: order.ClientOrderId,
			Price:         order.Price,
			Quantity:      order.OrigQty,
			ExecutedQty:   order.ExecutedQty,
			CumQuoteQty:   order.CumQuote,
			AvgPrice:      order.AvgPrice,
			Status:        b.bnConverter.FromBNOrderStatus(order.Status),
			Type:          b.bnConverter.FromBNOrderType(order.Type),
			Side:          b.bnConverter.FromBNOrderSide(order.Side),
			PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    order.Time,
			UpdateTime:    order.UpdateTime,
		})
	}
	return orders
}
func (b *BinanceTradeEngine) handleOrderFromSwapOrderQuery(req *QueryOrderParam, res *mybinanceapi.SwapOrderGetRes) *Order {
	order := &Order{
		Exchange:      BINANCE_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.ClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CumQuote,
		AvgPrice:      res.AvgPrice,
		Status:        b.bnConverter.FromBNOrderStatus(res.Status),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		PositionSide:  b.bnConverter.FromBNPositionSide(res.PositionSide),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		ReduceOnly:    res.ReduceOnly,
		CreateTime:    res.Time,
		UpdateTime:    res.UpdateTime,
	}
	return order
}
func (b *BinanceTradeEngine) handleOrderFromSwapOrdersQuery(req *QueryOrderParam, res *mybinanceapi.SwapAllOrdersRes) []*Order {
	var orders []*Order
	for _, order := range *res {
		orders = append(orders, &Order{
			Exchange:      BINANCE_NAME.String(),
			AccountType:   req.AccountType,
			Symbol:        order.Symbol,
			OrderId:       strconv.FormatInt(order.OrderId, 10),
			ClientOrderId: order.ClientOrderId,
			Price:         order.Price,
			Quantity:      order.OrigQty,
			ExecutedQty:   order.ExecutedQty,
			CumQuoteQty:   order.CumQuote,
			AvgPrice:      order.AvgPrice,
			Status:        b.bnConverter.FromBNOrderStatus(order.Status),
			Type:          b.bnConverter.FromBNOrderType(order.Type),
			Side:          b.bnConverter.FromBNOrderSide(order.Side),
			PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    order.Time,
			UpdateTime:    order.UpdateTime,
		})
	}
	return orders
}
func (b *BinanceTradeEngine) handleTradesFromSwapTradeQuery(req *QueryTradeParam, res *mybinanceapi.SwapUserTradesRes) []*Trade {
	var trades []*Trade
	for _, trade := range *res {
		trades = append(trades, &Trade{
			Exchange:     BINANCE_NAME.String(),
			AccountType:  req.AccountType,
			Symbol:       req.Symbol,
			TradeId:      strconv.FormatInt(trade.Id, 10),
			OrderId:      strconv.FormatInt(trade.OrderId, 10),
			Price:        trade.Price,
			Quantity:     trade.Qty,
			QuoteQty:     trade.BaseQty,
			Side:         b.bnConverter.FromBNOrderSide(trade.Side),
			PositionSide: b.bnConverter.FromBNPositionSide(trade.PositionSide),
			FeeAmount:    trade.Commission,
			FeeCcy:       trade.CommissionAsset,
			RealizedPnl:  trade.RealizedPnl,
			IsMaker:      trade.Maker,
			Timestamp:    trade.Time,
		})
	}
	return trades
}

func (b *BinanceTradeEngine) handleOrderFromSwapOrderCreate(req *OrderParam, res *mybinanceapi.SwapOrderPostRes) *Order {
	nowTimestamp := time.Now().UnixMilli()
	order := Order{
		Exchange:      BINANCE_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.ClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CumQuote,
		AvgPrice:      res.AvgPrice,
		Status:        b.bnConverter.FromBNOrderStatus(res.Status),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		PositionSide:  b.bnConverter.FromBNPositionSide(res.PositionSide),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		ReduceOnly:    res.ReduceOnly,
		CreateTime:    nowTimestamp,
		UpdateTime:    nowTimestamp,
	}
	return &order
}
func (b *BinanceTradeEngine) handleOrderFromSwapOrderAmend(req *OrderParam, res *mybinanceapi.SwapOrderPutRes) *Order {
	nowTimestamp := time.Now().UnixMilli()
	order := &Order{
		Exchange:      BINANCE_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.ClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CumQuote,
		AvgPrice:      res.AvgPrice,
		Status:        b.bnConverter.FromBNOrderStatus(res.Status),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		PositionSide:  b.bnConverter.FromBNPositionSide(res.PositionSide),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		ReduceOnly:    res.ReduceOnly,
		CreateTime:    nowTimestamp,
		UpdateTime:    nowTimestamp,
	}
	return order
}
func (b *BinanceTradeEngine) handleOrderFromSwapOrderCancel(req *OrderParam, res *mybinanceapi.SwapOrderDeleteRes) *Order {
	order := &Order{
		Exchange:      BINANCE_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.ClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CumQuote,
		AvgPrice:      res.AvgPrice,
		Status:        b.bnConverter.FromBNOrderStatus(res.Status),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		PositionSide:  b.bnConverter.FromBNPositionSide(res.PositionSide),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		ReduceOnly:    res.ReduceOnly,
		CreateTime:    res.UpdateTime,
		UpdateTime:    res.UpdateTime,
	}
	return order
}

func (b *BinanceTradeEngine) handleOrdersFromSwapBatchOrderCreate(reqs []*OrderParam, res *mybinanceapi.SwapBatchOrdersPostRes) []*Order {
	var orders []*Order
	nowTimestamp := time.Now().UnixMilli()
	for _, order := range *res {
		orders = append(orders, &Order{
			Exchange:      BINANCE_NAME.String(),
			AccountType:   reqs[0].AccountType,
			Symbol:        order.Symbol,
			OrderId:       strconv.FormatInt(order.OrderId, 10),
			ClientOrderId: order.ClientOrderId,
			Price:         order.Price,
			Quantity:      order.OrigQty,
			ExecutedQty:   order.ExecutedQty,
			CumQuoteQty:   order.CumQuote,
			AvgPrice:      order.AvgPrice,
			Status:        b.bnConverter.FromBNOrderStatus(order.Status),
			Type:          b.bnConverter.FromBNOrderType(order.Type),
			Side:          b.bnConverter.FromBNOrderSide(order.Side),
			PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    nowTimestamp,
			UpdateTime:    nowTimestamp,
			ErrorCode:     strconv.Itoa(order.Code),
			ErrorMsg:      order.Msg,
		})
	}
	return orders
}
func (b *BinanceTradeEngine) handleOrdersFromSwapBatchOrderAmend(reqs []*OrderParam, res *mybinanceapi.SwapBatchOrdersPutRes) []*Order {
	var orders []*Order
	nowTimestamp := time.Now().UnixMilli()
	for _, order := range *res {
		orders = append(orders, &Order{
			Exchange:      BINANCE_NAME.String(),
			AccountType:   reqs[0].AccountType,
			Symbol:        order.Symbol,
			OrderId:       strconv.FormatInt(order.OrderId, 10),
			ClientOrderId: order.ClientOrderId,
			Price:         order.Price,
			Quantity:      order.OrigQty,
			ExecutedQty:   order.ExecutedQty,
			CumQuoteQty:   order.CumQuote,
			AvgPrice:      order.AvgPrice,
			Status:        b.bnConverter.FromBNOrderStatus(order.Status),
			Type:          b.bnConverter.FromBNOrderType(order.Type),
			Side:          b.bnConverter.FromBNOrderSide(order.Side),
			PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    nowTimestamp,
			UpdateTime:    nowTimestamp,
			ErrorCode:     strconv.Itoa(order.Code),
			ErrorMsg:      order.Msg,
		})
	}
	return orders
}
func (b *BinanceTradeEngine) handleOrdersFromSwapBatchOrderCancel(reqs []*OrderParam, res *mybinanceapi.SwapBatchOrdersDeleteRes) []*Order {
	var orders []*Order
	nowTimestamp := time.Now().UnixMilli()
	for _, order := range *res {
		orders = append(orders, &Order{
			Exchange:      BINANCE_NAME.String(),
			AccountType:   reqs[0].AccountType,
			Symbol:        order.Symbol,
			OrderId:       strconv.FormatInt(order.OrderId, 10),
			ClientOrderId: order.ClientOrderId,
			Price:         order.Price,
			Quantity:      order.OrigQty,
			ExecutedQty:   order.ExecutedQty,
			CumQuoteQty:   order.CumQuote,
			AvgPrice:      order.AvgPrice,
			Status:        b.bnConverter.FromBNOrderStatus(order.Status),
			Type:          b.bnConverter.FromBNOrderType(order.Type),
			Side:          b.bnConverter.FromBNOrderSide(order.Side),
			PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    nowTimestamp,
			UpdateTime:    nowTimestamp,
			ErrorCode:     strconv.Itoa(order.Code),
			ErrorMsg:      order.Msg,
		})
	}
	return orders
}

// handle ws
func (b *BinanceTradeEngine) handleSubscribeOrderFromSpotPayload(req SubscribeOrderParam, newPayload *mybinanceapi.WsSpotPayload, newSub *subscription[Order]) {
	//处理不需要的订阅数据
	go func() {
		for {
			select {
			case <-newPayload.BalanceUpdatePayload.ErrChan():
				continue
			case <-newSub.closeChan:
				return
			case r := <-newPayload.BalanceUpdatePayload.ResultChan():
				_ = r
			}
		}
	}()
	go func() {
		for {
			select {
			case <-newPayload.OutboundAccountPositionPayload.ErrChan():
				continue
			case <-newSub.closeChan:
				return
			case r := <-newPayload.OutboundAccountPositionPayload.ResultChan():
				_ = r
			}
		}
	}()

	//处理订单推送订阅
	go func() {
		for {
			select {
			case err := <-newPayload.ExecutionReportPayload.ErrChan():
				newSub.errChan <- err
			case <-newSub.closeChan:
				newSub.CloseChan() <- struct{}{}
				return
			case r := <-newPayload.ExecutionReportPayload.ResultChan():
				avgPrice := decimal.Zero
				if r.ExecutedQty != "" && r.CummulativeQuoteQty != "" {
					executedQty, _ := decimal.NewFromString(r.ExecutedQty)
					cumQuoteQty, _ := decimal.NewFromString(r.CummulativeQuoteQty)
					if !executedQty.IsZero() {
						avgPrice = cumQuoteQty.Div(executedQty)
					}
				}
				order := Order{
					Exchange:      BINANCE_NAME.String(),
					AccountType:   req.AccountType,
					Symbol:        r.Symbol,
					OrderId:       strconv.FormatInt(r.OrderId, 10),
					ClientOrderId: r.ClientOrderId,
					Price:         r.Price,
					Quantity:      r.OrigQty,
					ExecutedQty:   r.ExecutedQty,
					CumQuoteQty:   r.CummulativeQuoteQty,
					AvgPrice:      avgPrice.String(),
					Status:        b.bnConverter.FromBNOrderStatus(r.Status),
					Type:          b.bnConverter.FromBNOrderType(r.Type),
					Side:          b.bnConverter.FromBNOrderSide(r.Side),
					TimeInForce:   b.bnConverter.FromBNTimeInForce(r.TimeInForce),
					FeeAmount:     r.FeeQty,
					FeeCcy:        r.FeeAsset,
					CreateTime:    r.OrderCreateTime,
					UpdateTime:    r.Timestamp,
				}
				newSub.resultChan <- order
			}
		}
	}()
}
func (b *BinanceTradeEngine) handleSubscribeOrderFromFuturePayload(req SubscribeOrderParam, newPayload *mybinanceapi.WsFuturePayload, newSub *subscription[Order]) {
	//处理不需要的订阅数据
	go func() {
		for {
			select {
			case <-newPayload.AccountUpdatePayload.ErrChan():
				continue
			case <-newSub.closeChan:
				return
			case r := <-newPayload.AccountUpdatePayload.ResultChan():
				_ = r
			}
		}
	}()

	//处理订单推送订阅
	go func() {
		for {
			select {
			case err := <-newPayload.OrderTradeUpdatePayload.ErrChan():
				newSub.errChan <- err
			case <-newSub.closeChan:
				newSub.CloseChan() <- struct{}{}
				return
			case result := <-newPayload.OrderTradeUpdatePayload.ResultChan():
				r := result.Order
				CumQuoteQty := decimal.Zero
				avgPrice, err := decimal.NewFromString(r.AvgPrice)
				if err != nil {
					newSub.ErrChan() <- err
				}
				CumQuoteQty = avgPrice.Mul(decimal.RequireFromString(r.ExecutedQty))
				order := Order{
					Exchange:      BINANCE_NAME.String(),
					AccountType:   req.AccountType,
					Symbol:        r.Symbol,
					OrderId:       strconv.FormatInt(r.OrderId, 10),
					ClientOrderId: r.ClientOrderId,
					Price:         r.Price,
					Quantity:      r.OrigQty,
					ExecutedQty:   r.ExecutedQty,
					CumQuoteQty:   CumQuoteQty.String(),
					AvgPrice:      r.AvgPrice,
					Status:        b.bnConverter.FromBNOrderStatus(r.Status),
					Type:          b.bnConverter.FromBNOrderType(r.Type),
					Side:          b.bnConverter.FromBNOrderSide(r.Side),
					PositionSide:  b.bnConverter.FromBNPositionSide(r.PositionSide),
					TimeInForce:   b.bnConverter.FromBNTimeInForce(r.TimeInForce),
					FeeAmount:     r.FeeQty,
					FeeCcy:        r.FeeAsset,
					ReduceOnly:    r.IsReduceOnly,
					CreateTime:    result.Timestamp,
					UpdateTime:    result.Timestamp,
				}
				newSub.resultChan <- order
			}
		}
	}()
}
func (b *BinanceTradeEngine) handleSubscribeOrderFromSwapPayload(req SubscribeOrderParam, newPayload *mybinanceapi.WsSwapPayload, newSub *subscription[Order]) {
	//处理不需要的订阅数据
	go func() {
		for {
			select {
			case <-newPayload.AccountUpdatePayload.ErrChan():
				continue
			case <-newSub.closeChan:
				return
			case r := <-newPayload.AccountUpdatePayload.ResultChan():
				_ = r
			}
		}
	}()

	//处理订单推送订阅
	go func() {
		for {
			select {
			case err := <-newPayload.OrderTradeUpdatePayload.ErrChan():
				newSub.errChan <- err
			case <-newSub.closeChan:
				newSub.CloseChan() <- struct{}{}
				return
			case result := <-newPayload.OrderTradeUpdatePayload.ResultChan():
				r := result.Order
				CumQuoteQty := decimal.Zero
				avgPrice, err := decimal.NewFromString(r.AvgPrice)
				if err != nil {
					newSub.ErrChan() <- err
				}
				CumQuoteQty = avgPrice.Mul(decimal.RequireFromString(r.ExecutedQty))
				order := Order{
					Exchange:      BINANCE_NAME.String(),
					AccountType:   req.AccountType,
					Symbol:        r.Symbol,
					OrderId:       strconv.FormatInt(r.OrderId, 10),
					ClientOrderId: r.ClientOrderId,
					Price:         r.Price,
					Quantity:      r.OrigQty,
					ExecutedQty:   r.ExecutedQty,
					CumQuoteQty:   CumQuoteQty.String(),
					AvgPrice:      r.AvgPrice,
					Status:        b.bnConverter.FromBNOrderStatus(r.Status),
					Type:          b.bnConverter.FromBNOrderType(r.Type),
					Side:          b.bnConverter.FromBNOrderSide(r.Side),
					PositionSide:  b.bnConverter.FromBNPositionSide(r.PositionSide),
					TimeInForce:   b.bnConverter.FromBNTimeInForce(r.TimeInForce),
					FeeAmount:     r.FeeQty,
					FeeCcy:        r.FeeAsset,
					ReduceOnly:    r.IsReduceOnly,
					CreateTime:    result.Timestamp,
					UpdateTime:    result.Timestamp,
				}
				newSub.resultChan <- order
			}
		}
	}()
}

func (b *BinanceTradeEngine) restBatchPreCheck(reqs []*OrderParam) error {
	//检测长度，BINANCE最多批量下5个订单
	if len(reqs) > 5 {
		return ErrorInvalid("binance order param length require less than 5")
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

func (b *BinanceTradeEngine) accountTypePreCheck(accountType string) error {
	switch BinanceAccountType(accountType) {
	case BN_AC_SPOT, BN_AC_FUTURE, BN_AC_SWAP:
		return nil
	default:
		return ErrorInvalid("binance account type invalid")
	}
}
