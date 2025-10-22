package mytrade

import (
	"strconv"
	"time"

	"github.com/Hongssd/mybinanceapi"
	"github.com/shopspring/decimal"
)

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
			IsMargin:      req.IsMargin,
			IsIsolated:    req.IsIsolated,
			OrderId:       strconv.FormatInt(order.OrderId, 10),
			ClientOrderId: order.ClientOrderId,
			Price:         order.Price,
			Quantity:      order.OrigQty,
			ExecutedQty:   order.ExecutedQty,
			CumQuoteQty:   order.CummulativeQuoteQty,
			AvgPrice:      avgPrice.String(),
			Status:        b.bnConverter.FromBNOrderStatus(order.Status, order.Type),
			Type:          b.bnConverter.FromBNOrderType(order.Type),
			Side:          b.bnConverter.FromBNOrderSide(order.Side),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
			CreateTime:    order.Time,
			UpdateTime:    order.UpdateTime,

			TriggerPrice:         order.StopPrice,
			TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(order.Type),
			TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(order.Side, order.Type),
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
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.ClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CummulativeQuoteQty,
		AvgPrice:      avgPrice.String(),
		Status:        b.bnConverter.FromBNOrderStatus(res.Status, res.Type),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		CreateTime:    res.Time,
		UpdateTime:    res.UpdateTime,

		TriggerPrice:         res.StopPrice,
		TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(res.Type),
		TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(res.Side, res.Type),
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
			IsMargin:      req.IsMargin,
			IsIsolated:    req.IsIsolated,
			OrderId:       strconv.FormatInt(order.OrderId, 10),
			ClientOrderId: order.ClientOrderId,
			Price:         order.Price,
			Quantity:      order.OrigQty,
			ExecutedQty:   order.ExecutedQty,
			CumQuoteQty:   order.CummulativeQuoteQty,
			AvgPrice:      avgPrice.String(),
			Status:        b.bnConverter.FromBNOrderStatus(order.Status, order.Type),
			Type:          b.bnConverter.FromBNOrderType(order.Type),
			Side:          b.bnConverter.FromBNOrderSide(order.Side),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
			CreateTime:    order.Time,
			UpdateTime:    order.UpdateTime,

			TriggerPrice:         order.StopPrice,
			TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(order.Type),
			TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(order.Side, order.Type),
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

	// d, _ := json.MarshalIndent(res, "", "  ")
	// log.Info(string(d))
	order := &Order{
		Exchange:      BINANCE_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.ClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CummulativeQuoteQty,
		AvgPrice:      avgPrice.String(),
		Status:        b.bnConverter.FromBNOrderStatus(res.Status, res.Type),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		CreateTime:    res.WorkingTime,
		UpdateTime:    res.WorkingTime,

		TriggerPrice:         req.TriggerPrice.String(),
		TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(res.Type),
		TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(res.Side, res.Type),
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
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
		OrderId:       strconv.FormatInt(res.NewOrderResponse.OrderId, 10),
		ClientOrderId: res.NewOrderResponse.ClientOrderId,
		Price:         res.NewOrderResponse.Price,
		Quantity:      res.NewOrderResponse.OrigQty,
		ExecutedQty:   res.NewOrderResponse.ExecutedQty,
		CumQuoteQty:   res.NewOrderResponse.CummulativeQuoteQty,
		AvgPrice:      avgPrice.String(),
		Status:        b.bnConverter.FromBNOrderStatus(res.NewOrderResponse.Status, res.NewOrderResponse.Type),
		Type:          b.bnConverter.FromBNOrderType(res.NewOrderResponse.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.NewOrderResponse.Side),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.NewOrderResponse.TimeInForce),
		CreateTime:    res.NewOrderResponse.WorkingTime,
		UpdateTime:    res.NewOrderResponse.WorkingTime,

		TriggerPrice:         req.TriggerPrice.String(),
		TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(res.NewOrderResponse.Type),
		TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(res.NewOrderResponse.Side, res.NewOrderResponse.Type),
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
	//d, _ := json.MarshalIndent(res, "", "  ")
	//log.Info(string(d))
	order := &Order{
		Exchange:      BINANCE_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.OrigClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CummulativeQuoteQty,
		AvgPrice:      avgPrice.String(),
		Status:        b.bnConverter.FromBNOrderStatus(res.Status, res.Type),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		CreateTime:    res.TransactTime,
		UpdateTime:    nowTimestamp,

		TriggerPrice:         req.TriggerPrice.String(),
		TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(res.Type),
		TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(res.Side, res.Type),
	}
	return order
}

func (b *BinanceTradeEngine) handleOrderFromSpotBatchErr(req *OrderParam, err error) *Order {
	return &Order{
		Exchange:      BINANCE_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
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

// 现货杠杆订单处理
func (b *BinanceTradeEngine) handleOrdersFromSpotMarginOpenOrders(req *QueryOrderParam, res *mybinanceapi.MarginOpenOrdersRes) []*Order {
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
			IsMargin:      req.IsMargin,
			IsIsolated:    req.IsIsolated,
			OrderId:       strconv.FormatInt(order.OrderId, 10),
			ClientOrderId: order.ClientOrderId,
			Price:         order.Price,
			Quantity:      order.OrigQty,
			ExecutedQty:   order.ExecutedQty,
			CumQuoteQty:   order.CummulativeQuoteQty,
			AvgPrice:      avgPrice.String(),
			Status:        b.bnConverter.FromBNOrderStatus(order.Status, order.Type),
			Type:          b.bnConverter.FromBNOrderType(order.Type),
			Side:          b.bnConverter.FromBNOrderSide(order.Side),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
			CreateTime:    order.Time,
			UpdateTime:    order.UpdateTime,

			TriggerPrice:         order.StopPrice,
			TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(order.Type),
			TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(order.Side, order.Type),
		})
	}
	return orders
}
func (b *BinanceTradeEngine) handleOrderFromSpotMarginOrderQuery(req *QueryOrderParam, res *mybinanceapi.MarginOrderGetRes) *Order {
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
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.ClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CummulativeQuoteQty,
		AvgPrice:      avgPrice.String(),
		Status:        b.bnConverter.FromBNOrderStatus(res.Status, res.Type),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		CreateTime:    res.Time,
		UpdateTime:    res.UpdateTime,

		TriggerPrice:         res.StopPrice,
		TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(res.Type),
		TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(res.Side, res.Type),
	}
	return order
}
func (b *BinanceTradeEngine) handleOrderFromSpotMarginOrdersQuery(req *QueryOrderParam, res *mybinanceapi.MarginAllOrdersRes) []*Order {
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
			IsMargin:      req.IsMargin,
			IsIsolated:    req.IsIsolated,
			OrderId:       strconv.FormatInt(order.OrderId, 10),
			ClientOrderId: order.ClientOrderId,
			Price:         order.Price,
			Quantity:      order.OrigQty,
			ExecutedQty:   order.ExecutedQty,
			CumQuoteQty:   order.CummulativeQuoteQty,
			AvgPrice:      avgPrice.String(),
			Status:        b.bnConverter.FromBNOrderStatus(order.Status, order.Type),
			Type:          b.bnConverter.FromBNOrderType(order.Type),
			Side:          b.bnConverter.FromBNOrderSide(order.Side),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
			CreateTime:    order.Time,
			UpdateTime:    order.UpdateTime,

			TriggerPrice:         order.StopPrice,
			TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(order.Type),
			TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(order.Side, order.Type),
		})
	}
	return orders
}

func (b *BinanceTradeEngine) handleOrderFromSpotMarginOrderCreate(req *OrderParam, res *mybinanceapi.SpotMarginOrderPostRes) *Order {
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
		IsMargin:      req.IsMargin,
		IsIsolated:    res.IsIsolated,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.ClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CummulativeQuoteQty,
		AvgPrice:      avgPrice.String(),
		Status:        b.bnConverter.FromBNOrderStatus(res.Status, res.Type),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		//杠杆下单没有以下返回值
		//CreateTime:    res.WorkingTime,
		//UpdateTime:    res.WorkingTime,
		//更改为
		CreateTime: res.TransactTime,
		UpdateTime: res.TransactTime,

		TriggerPrice:         req.TriggerPrice.String(),
		TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(res.Type),
		TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(res.Side, res.Type),
	}
	return order
}
func (b *BinanceTradeEngine) handleOrderFromSpotMarginOrderCancel(req *OrderParam, res *mybinanceapi.SpotMarginOrderDeleteRes) *Order {
	nowTimestamp := time.Now().UnixMilli()
	avgPrice := decimal.Zero
	if res.ExecutedQty != "" && res.CummulativeQuoteQty != "" {
		executedQty, _ := decimal.NewFromString(res.ExecutedQty)
		cumQuoteQty, _ := decimal.NewFromString(res.CummulativeQuoteQty)
		if !executedQty.IsZero() {
			avgPrice = cumQuoteQty.Div(executedQty)
		}
	}
	//d, _ := json.MarshalIndent(res, "", "  ")
	//log.Info(string(d))
	order := &Order{
		Exchange:      BINANCE_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
		OrderId:       res.OrderId,
		ClientOrderId: res.OrigClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CummulativeQuoteQty,
		AvgPrice:      avgPrice.String(),
		Status:        b.bnConverter.FromBNOrderStatus(res.Status, res.Type),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		CreateTime:    nowTimestamp,
		UpdateTime:    nowTimestamp,

		TriggerPrice:         req.TriggerPrice.String(),
		TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(res.Type),
		TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(res.Side, res.Type),
	}
	return order
}

// U合约订单处理
func (b *BinanceTradeEngine) handleOrdersFromFutureOpenOrders(req *QueryOrderParam, res *mybinanceapi.FutureOpenOrdersRes) []*Order {
	var orders []*Order
	for _, order := range *res {
		orders = append(orders, &Order{
			Exchange:      BINANCE_NAME.String(),
			AccountType:   req.AccountType,
			Symbol:        order.Symbol,
			IsMargin:      req.IsMargin,
			IsIsolated:    req.IsIsolated,
			OrderId:       strconv.FormatInt(order.OrderId, 10),
			ClientOrderId: order.ClientOrderId,
			Price:         order.Price,
			Quantity:      order.OrigQty,
			ExecutedQty:   order.ExecutedQty,
			CumQuoteQty:   order.CumQuote,
			AvgPrice:      order.AvgPrice,
			Status:        b.bnConverter.FromBNOrderStatus(order.Status, order.Type),
			Type:          b.bnConverter.FromBNOrderType(order.Type),
			Side:          b.bnConverter.FromBNOrderSide(order.Side),
			PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    order.Time,
			UpdateTime:    order.UpdateTime,

			TriggerPrice:         order.StopPrice,
			TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(order.Type),
			TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(order.Side, order.Type),
		})
	}
	return orders
}
func (b *BinanceTradeEngine) handleOrderFromFutureOrderQuery(req *QueryOrderParam, res *mybinanceapi.FutureOrderGetRes) *Order {
	order := &Order{
		Exchange:      BINANCE_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.ClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CumQuote,
		AvgPrice:      res.AvgPrice,
		Status:        b.bnConverter.FromBNOrderStatus(res.Status, res.Type),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		PositionSide:  b.bnConverter.FromBNPositionSide(res.PositionSide),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		ReduceOnly:    res.ReduceOnly,
		CreateTime:    res.Time,
		UpdateTime:    res.UpdateTime,

		TriggerPrice:         res.StopPrice,
		TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(res.Type),
		TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(res.Side, res.Type),
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
			IsMargin:      req.IsMargin,
			IsIsolated:    req.IsIsolated,
			OrderId:       strconv.FormatInt(order.OrderId, 10),
			ClientOrderId: order.ClientOrderId,
			Price:         order.Price,
			Quantity:      order.OrigQty,
			ExecutedQty:   order.ExecutedQty,
			CumQuoteQty:   order.CumQuote,
			AvgPrice:      order.AvgPrice,
			Status:        b.bnConverter.FromBNOrderStatus(order.Status, order.Type),
			Type:          b.bnConverter.FromBNOrderType(order.Type),
			Side:          b.bnConverter.FromBNOrderSide(order.Side),
			PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    order.Time,
			UpdateTime:    order.UpdateTime,

			TriggerPrice:         order.StopPrice,
			TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(order.Type),
			TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(order.Side, order.Type),
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
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.ClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CumQuote,
		AvgPrice:      res.AvgPrice,
		Status:        b.bnConverter.FromBNOrderStatus(res.Status, res.Type),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		PositionSide:  b.bnConverter.FromBNPositionSide(res.PositionSide),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		ReduceOnly:    res.ReduceOnly,
		CreateTime:    nowTimestamp,
		UpdateTime:    nowTimestamp,

		TriggerPrice:         res.StopPrice,
		TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(res.Type),
		TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(res.Side, res.Type),
	}
	return order
}
func (b *BinanceTradeEngine) handleOrderFromFutureOrderAmend(req *OrderParam, res *mybinanceapi.FutureOrderPutRes) *Order {
	nowTimestamp := time.Now().UnixMilli()
	order := &Order{
		Exchange:      BINANCE_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.ClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CumQuote,
		AvgPrice:      res.AvgPrice,
		Status:        b.bnConverter.FromBNOrderStatus(res.Status, res.Type),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		PositionSide:  b.bnConverter.FromBNPositionSide(res.PositionSide),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		ReduceOnly:    res.ReduceOnly,
		CreateTime:    nowTimestamp,
		UpdateTime:    nowTimestamp,

		TriggerPrice:         res.StopPrice,
		TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(res.Type),
		TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(res.Side, res.Type),
	}
	return order
}
func (b *BinanceTradeEngine) handleOrderFromFutureOrderCancel(req *OrderParam, res *mybinanceapi.FutureOrderDeleteRes) *Order {
	order := &Order{
		Exchange:      BINANCE_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.ClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CumQuote,
		AvgPrice:      res.AvgPrice,
		Status:        b.bnConverter.FromBNOrderStatus(res.Status, res.Type),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		PositionSide:  b.bnConverter.FromBNPositionSide(res.PositionSide),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		ReduceOnly:    res.ReduceOnly,
		CreateTime:    res.UpdateTime,
		UpdateTime:    res.UpdateTime,

		TriggerPrice:         res.StopPrice,
		TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(res.Type),
		TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(res.Side, res.Type),
	}
	return order
}

func (b *BinanceTradeEngine) handleOrdersFromFutureBatchOrderCreate(reqs []*OrderParam, res *mybinanceapi.FutureBatchOrdersPostRes) []*Order {
	var orders []*Order
	nowTimestamp := time.Now().UnixMilli()
	for _, order := range *res {
		code, ok := order.Code.(float64)
		if !ok {
			code = 0
		}
		codeInt := int(code)
		orders = append(orders, &Order{
			Exchange:      BINANCE_NAME.String(),
			AccountType:   reqs[0].AccountType,
			Symbol:        order.Symbol,
			IsMargin:      reqs[0].IsMargin,
			IsIsolated:    reqs[0].IsIsolated,
			OrderId:       strconv.FormatInt(order.OrderId, 10),
			ClientOrderId: order.ClientOrderId,
			Price:         order.Price,
			Quantity:      order.OrigQty,
			ExecutedQty:   order.ExecutedQty,
			CumQuoteQty:   order.CumQuote,
			AvgPrice:      order.AvgPrice,
			Status:        b.bnConverter.FromBNOrderStatus(order.Status, order.Type),
			Type:          b.bnConverter.FromBNOrderType(order.Type),
			Side:          b.bnConverter.FromBNOrderSide(order.Side),
			PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    nowTimestamp,
			UpdateTime:    nowTimestamp,
			ErrorCode:     strconv.Itoa(codeInt),
			ErrorMsg:      order.Msg,

			TriggerPrice:         order.StopPrice,
			TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(order.Type),
			TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(order.Side, order.Type),
		})
	}
	return orders
}
func (b *BinanceTradeEngine) handleOrdersFromFutureBatchOrderAmend(reqs []*OrderParam, res *mybinanceapi.FutureBatchOrdersPutRes) []*Order {
	var orders []*Order
	nowTimestamp := time.Now().UnixMilli()
	for _, order := range *res {
		code, ok := order.Code.(float64)
		if !ok {
			code = 0
		}
		codeInt := int(code)
		orders = append(orders, &Order{
			Exchange:      BINANCE_NAME.String(),
			AccountType:   reqs[0].AccountType,
			Symbol:        order.Symbol,
			IsMargin:      reqs[0].IsMargin,
			IsIsolated:    reqs[0].IsIsolated,
			OrderId:       strconv.FormatInt(order.OrderId, 10),
			ClientOrderId: order.ClientOrderId,
			Price:         order.Price,
			Quantity:      order.OrigQty,
			ExecutedQty:   order.ExecutedQty,
			CumQuoteQty:   order.CumQuote,
			AvgPrice:      order.AvgPrice,
			Status:        b.bnConverter.FromBNOrderStatus(order.Status, order.Type),
			Type:          b.bnConverter.FromBNOrderType(order.Type),
			Side:          b.bnConverter.FromBNOrderSide(order.Side),
			PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    nowTimestamp,
			UpdateTime:    nowTimestamp,
			ErrorCode:     strconv.Itoa(codeInt),
			ErrorMsg:      order.Msg,

			TriggerPrice:         order.StopPrice,
			TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(order.Type),
			TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(order.Side, order.Type),
		})
	}
	return orders
}
func (b *BinanceTradeEngine) handleOrdersFromFutureBatchOrderCancel(reqs []*OrderParam, res *mybinanceapi.FutureBatchOrdersDeleteRes) []*Order {
	var orders []*Order
	nowTimestamp := time.Now().UnixMilli()
	for _, order := range *res {
		code, ok := order.Code.(float64)
		if !ok {
			code = 0
		}
		codeInt := int(code)
		orders = append(orders, &Order{
			Exchange:      BINANCE_NAME.String(),
			AccountType:   reqs[0].AccountType,
			Symbol:        order.Symbol,
			IsMargin:      reqs[0].IsMargin,
			IsIsolated:    reqs[0].IsIsolated,
			OrderId:       strconv.FormatInt(order.OrderId, 10),
			ClientOrderId: order.ClientOrderId,
			Price:         order.Price,
			Quantity:      order.OrigQty,
			ExecutedQty:   order.ExecutedQty,
			CumQuoteQty:   order.CumQuote,
			AvgPrice:      order.AvgPrice,
			Status:        b.bnConverter.FromBNOrderStatus(order.Status, order.Type),
			Type:          b.bnConverter.FromBNOrderType(order.Type),
			Side:          b.bnConverter.FromBNOrderSide(order.Side),
			PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    nowTimestamp,
			UpdateTime:    nowTimestamp,
			ErrorCode:     strconv.Itoa(codeInt),
			ErrorMsg:      order.Msg,

			TriggerPrice:         order.StopPrice,
			TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(order.Type),
			TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(order.Side, order.Type),
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
			IsMargin:      req.IsMargin,
			IsIsolated:    req.IsIsolated,
			OrderId:       strconv.FormatInt(order.OrderId, 10),
			ClientOrderId: order.ClientOrderId,
			Price:         order.Price,
			Quantity:      order.OrigQty,
			ExecutedQty:   order.ExecutedQty,
			CumQuoteQty:   order.CumQuote,
			AvgPrice:      order.AvgPrice,
			Status:        b.bnConverter.FromBNOrderStatus(order.Status, order.Type),
			Type:          b.bnConverter.FromBNOrderType(order.Type),
			Side:          b.bnConverter.FromBNOrderSide(order.Side),
			PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    order.Time,
			UpdateTime:    order.UpdateTime,

			TriggerPrice:         order.StopPrice,
			TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(order.Type),
			TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(order.Side, order.Type),
		})
	}
	return orders
}
func (b *BinanceTradeEngine) handleOrderFromSwapOrderQuery(req *QueryOrderParam, res *mybinanceapi.SwapOrderGetRes) *Order {
	order := &Order{
		Exchange:      BINANCE_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.ClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CumQuote,
		AvgPrice:      res.AvgPrice,
		Status:        b.bnConverter.FromBNOrderStatus(res.Status, res.Type),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		PositionSide:  b.bnConverter.FromBNPositionSide(res.PositionSide),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		ReduceOnly:    res.ReduceOnly,
		CreateTime:    res.Time,
		UpdateTime:    res.UpdateTime,

		TriggerPrice:         res.StopPrice,
		TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(res.Type),
		TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(res.Side, res.Type),
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
			IsMargin:      req.IsMargin,
			IsIsolated:    req.IsIsolated,
			OrderId:       strconv.FormatInt(order.OrderId, 10),
			ClientOrderId: order.ClientOrderId,
			Price:         order.Price,
			Quantity:      order.OrigQty,
			ExecutedQty:   order.ExecutedQty,
			CumQuoteQty:   order.CumQuote,
			AvgPrice:      order.AvgPrice,
			Status:        b.bnConverter.FromBNOrderStatus(order.Status, order.Type),
			Type:          b.bnConverter.FromBNOrderType(order.Type),
			Side:          b.bnConverter.FromBNOrderSide(order.Side),
			PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    order.Time,
			UpdateTime:    order.UpdateTime,

			TriggerPrice:         order.StopPrice,
			TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(order.Type),
			TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(order.Side, order.Type),
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
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.ClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CumQuote,
		AvgPrice:      res.AvgPrice,
		Status:        b.bnConverter.FromBNOrderStatus(res.Status, res.Type),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		PositionSide:  b.bnConverter.FromBNPositionSide(res.PositionSide),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		ReduceOnly:    res.ReduceOnly,
		CreateTime:    nowTimestamp,
		UpdateTime:    nowTimestamp,

		TriggerPrice:         res.StopPrice,
		TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(res.Type),
		TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(res.Side, res.Type),
	}
	return &order
}
func (b *BinanceTradeEngine) handleOrderFromSwapOrderAmend(req *OrderParam, res *mybinanceapi.SwapOrderPutRes) *Order {
	nowTimestamp := time.Now().UnixMilli()
	order := &Order{
		Exchange:      BINANCE_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.ClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CumQuote,
		AvgPrice:      res.AvgPrice,
		Status:        b.bnConverter.FromBNOrderStatus(res.Status, res.Type),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		PositionSide:  b.bnConverter.FromBNPositionSide(res.PositionSide),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		ReduceOnly:    res.ReduceOnly,
		CreateTime:    nowTimestamp,
		UpdateTime:    nowTimestamp,

		TriggerPrice:         res.StopPrice,
		TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(res.Type),
		TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(res.Side, res.Type),
	}
	return order
}
func (b *BinanceTradeEngine) handleOrderFromSwapOrderCancel(req *OrderParam, res *mybinanceapi.SwapOrderDeleteRes) *Order {
	order := &Order{
		Exchange:      BINANCE_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.ClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CumQuote,
		AvgPrice:      res.AvgPrice,
		Status:        b.bnConverter.FromBNOrderStatus(res.Status, res.Type),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		PositionSide:  b.bnConverter.FromBNPositionSide(res.PositionSide),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		ReduceOnly:    res.ReduceOnly,
		CreateTime:    res.UpdateTime,
		UpdateTime:    res.UpdateTime,

		TriggerPrice:         res.StopPrice,
		TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(res.Type),
		TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(res.Side, res.Type),
	}
	return order
}

func (b *BinanceTradeEngine) handleOrdersFromSwapBatchOrderCreate(reqs []*OrderParam, res *mybinanceapi.SwapBatchOrdersPostRes) []*Order {
	var orders []*Order
	nowTimestamp := time.Now().UnixMilli()
	for _, order := range *res {
		code, ok := order.Code.(float64)
		if !ok {
			code = 0
		}
		codeInt := int(code)
		orders = append(orders, &Order{
			Exchange:      BINANCE_NAME.String(),
			AccountType:   reqs[0].AccountType,
			Symbol:        order.Symbol,
			IsMargin:      reqs[0].IsMargin,
			IsIsolated:    reqs[0].IsIsolated,
			OrderId:       strconv.FormatInt(order.OrderId, 10),
			ClientOrderId: order.ClientOrderId,
			Price:         order.Price,
			Quantity:      order.OrigQty,
			ExecutedQty:   order.ExecutedQty,
			CumQuoteQty:   order.CumQuote,
			AvgPrice:      order.AvgPrice,
			Status:        b.bnConverter.FromBNOrderStatus(order.Status, order.Type),
			Type:          b.bnConverter.FromBNOrderType(order.Type),
			Side:          b.bnConverter.FromBNOrderSide(order.Side),
			PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    nowTimestamp,
			UpdateTime:    nowTimestamp,
			ErrorCode:     strconv.Itoa(codeInt),
			ErrorMsg:      order.Msg,

			TriggerPrice:         order.StopPrice,
			TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(order.Type),
			TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(order.Side, order.Type),
		})
	}
	return orders
}
func (b *BinanceTradeEngine) handleOrdersFromSwapBatchOrderAmend(reqs []*OrderParam, res *mybinanceapi.SwapBatchOrdersPutRes) []*Order {
	var orders []*Order
	nowTimestamp := time.Now().UnixMilli()
	for _, order := range *res {
		code, ok := order.Code.(float64)
		if !ok {
			code = 0
		}
		codeInt := int(code)
		orders = append(orders, &Order{
			Exchange:      BINANCE_NAME.String(),
			AccountType:   reqs[0].AccountType,
			Symbol:        order.Symbol,
			IsMargin:      reqs[0].IsMargin,
			IsIsolated:    reqs[0].IsIsolated,
			OrderId:       strconv.FormatInt(order.OrderId, 10),
			ClientOrderId: order.ClientOrderId,
			Price:         order.Price,
			Quantity:      order.OrigQty,
			ExecutedQty:   order.ExecutedQty,
			CumQuoteQty:   order.CumQuote,
			AvgPrice:      order.AvgPrice,
			Status:        b.bnConverter.FromBNOrderStatus(order.Status, order.Type),
			Type:          b.bnConverter.FromBNOrderType(order.Type),
			Side:          b.bnConverter.FromBNOrderSide(order.Side),
			PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    nowTimestamp,
			UpdateTime:    nowTimestamp,
			ErrorCode:     strconv.Itoa(codeInt),
			ErrorMsg:      order.Msg,

			TriggerPrice:         order.StopPrice,
			TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(order.Type),
			TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(order.Side, order.Type),
		})
	}
	return orders
}
func (b *BinanceTradeEngine) handleOrdersFromSwapBatchOrderCancel(reqs []*OrderParam, res *mybinanceapi.SwapBatchOrdersDeleteRes) []*Order {
	var orders []*Order
	nowTimestamp := time.Now().UnixMilli()
	for _, order := range *res {
		code, ok := order.Code.(float64)
		if !ok {
			code = 0
		}
		codeInt := int(code)
		orders = append(orders, &Order{
			Exchange:      BINANCE_NAME.String(),
			AccountType:   reqs[0].AccountType,
			Symbol:        order.Symbol,
			IsMargin:      reqs[0].IsMargin,
			IsIsolated:    reqs[0].IsIsolated,
			OrderId:       strconv.FormatInt(order.OrderId, 10),
			ClientOrderId: order.ClientOrderId,
			Price:         order.Price,
			Quantity:      order.OrigQty,
			ExecutedQty:   order.ExecutedQty,
			CumQuoteQty:   order.CumQuote,
			AvgPrice:      order.AvgPrice,
			Status:        b.bnConverter.FromBNOrderStatus(order.Status, order.Type),
			Type:          b.bnConverter.FromBNOrderType(order.Type),
			Side:          b.bnConverter.FromBNOrderSide(order.Side),
			PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    nowTimestamp,
			UpdateTime:    nowTimestamp,
			ErrorCode:     strconv.Itoa(codeInt),
			ErrorMsg:      order.Msg,

			TriggerPrice:         order.StopPrice,
			TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(order.Type),
			TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(order.Side, order.Type),
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
					Status:        b.bnConverter.FromBNOrderStatus(r.Status, r.Type),
					Type:          b.bnConverter.FromBNOrderType(r.Type),
					Side:          b.bnConverter.FromBNOrderSide(r.Side),
					TimeInForce:   b.bnConverter.FromBNTimeInForce(r.TimeInForce),
					FeeAmount:     r.FeeQty,
					FeeCcy:        r.FeeAsset,
					CreateTime:    r.OrderCreateTime,
					UpdateTime:    r.Timestamp,

					TriggerPrice:         r.StopPrice,
					TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(r.Type),
					TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(r.Side, r.Type),
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
					Status:        b.bnConverter.FromBNOrderStatus(r.Status, r.Type),
					Type:          b.bnConverter.FromBNOrderType(r.Type),
					Side:          b.bnConverter.FromBNOrderSide(r.Side),
					PositionSide:  b.bnConverter.FromBNPositionSide(r.PositionSide),
					TimeInForce:   b.bnConverter.FromBNTimeInForce(r.TimeInForce),
					FeeAmount:     r.FeeQty,
					FeeCcy:        r.FeeAsset,
					ReduceOnly:    r.IsReduceOnly,
					CreateTime:    result.TradeTime,
					UpdateTime:    r.TradeTime,

					TriggerPrice:         r.StopPrice,
					TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(r.Type),
					TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(r.Side, r.Type),
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
					Status:        b.bnConverter.FromBNOrderStatus(r.Status, r.Type),
					Type:          b.bnConverter.FromBNOrderType(r.Type),
					Side:          b.bnConverter.FromBNOrderSide(r.Side),
					PositionSide:  b.bnConverter.FromBNPositionSide(r.PositionSide),
					TimeInForce:   b.bnConverter.FromBNTimeInForce(r.TimeInForce),
					FeeAmount:     r.FeeQty,
					FeeCcy:        r.FeeAsset,
					ReduceOnly:    r.IsReduceOnly,
					CreateTime:    result.Timestamp,
					UpdateTime:    result.Timestamp,

					TriggerPrice:         r.StopPrice,
					TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(r.Type),
					TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(r.Side, r.Type),
				}
				newSub.resultChan <- order
			}
		}
	}()
}
