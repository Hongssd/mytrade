package mytrade

import (
	"strconv"
	"time"

	"github.com/Hongssd/myasterapi"
	"github.com/shopspring/decimal"
)

// 现货订单处理
func (b *AsterTradeEngine) handleOrdersFromSpotOpenOrders(req *QueryOrderParam, res *myasterapi.SpotOpenOrdersRes) []*Order {
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
			Exchange:      ASTER_NAME.String(),
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
			Status:        b.asterConverter.FromAsterOrderStatus(order.Status, order.Type),
			Type:          b.asterConverter.FromAsterOrderType(order.Type),
			Side:          b.asterConverter.FromAsterOrderSide(order.Side),
			TimeInForce:   b.asterConverter.FromAsterTimeInForce(order.TimeInForce),
			CreateTime:    order.Time,
			UpdateTime:    order.UpdateTime,

			TriggerPrice:         order.StopPrice,
			TriggerType:          b.asterConverter.FromAsterOrderTypeForTriggerType(order.Type),
			TriggerConditionType: b.asterConverter.FromAsterOrderSideForTriggerConditionType(order.Side, order.Type),
		})
	}
	return orders
}
func (b *AsterTradeEngine) handleOrderFromSpotOrderQuery(req *QueryOrderParam, res *myasterapi.SpotOrderGetRes) *Order {
	avgPrice := decimal.Zero
	if res.ExecutedQty != "" && res.CummulativeQuoteQty != "" {
		executedQty, _ := decimal.NewFromString(res.ExecutedQty)
		cumQuoteQty, _ := decimal.NewFromString(res.CummulativeQuoteQty)
		if !executedQty.IsZero() {
			avgPrice = cumQuoteQty.Div(executedQty)
		}
	}
	order := &Order{
		Exchange:      ASTER_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        res.Symbol,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.ClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CummulativeQuoteQty,
		AvgPrice:      avgPrice.String(),
		Status:        b.asterConverter.FromAsterOrderStatus(res.Status, res.Type),
		Type:          b.asterConverter.FromAsterOrderType(res.Type),
		Side:          b.asterConverter.FromAsterOrderSide(res.Side),
		TimeInForce:   b.asterConverter.FromAsterTimeInForce(res.TimeInForce),
		CreateTime:    res.Time,
		UpdateTime:    res.UpdateTime,

		TriggerPrice:         res.StopPrice,
		TriggerType:          b.asterConverter.FromAsterOrderTypeForTriggerType(res.Type),
		TriggerConditionType: b.asterConverter.FromAsterOrderSideForTriggerConditionType(res.Side, res.Type),
	}
	return order
}
func (b *AsterTradeEngine) handleOrderFromSpotOrdersQuery(req *QueryOrderParam, res *myasterapi.SpotAllOrdersRes) []*Order {
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
			Exchange:      ASTER_NAME.String(),
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
			Status:        b.asterConverter.FromAsterOrderStatus(order.Status, order.Type),
			Type:          b.asterConverter.FromAsterOrderType(order.Type),
			Side:          b.asterConverter.FromAsterOrderSide(order.Side),
			TimeInForce:   b.asterConverter.FromAsterTimeInForce(order.TimeInForce),
			CreateTime:    order.Time,
			UpdateTime:    order.UpdateTime,

			TriggerPrice:         order.StopPrice,
			TriggerType:          b.asterConverter.FromAsterOrderTypeForTriggerType(order.Type),
			TriggerConditionType: b.asterConverter.FromAsterOrderSideForTriggerConditionType(order.Side, order.Type),
		})
	}
	return orders
}
func (b *AsterTradeEngine) handleTradesFromSpotTradeQuery(req *QueryTradeParam, res *myasterapi.SpotMyTradesRes) []*Trade {
	var trades []*Trade
	for _, trade := range *res {
		var orderSide OrderSide
		if trade.IsBuyer {
			orderSide = ORDER_SIDE_BUY
		} else {
			orderSide = ORDER_SIDE_SELL
		}
		trades = append(trades, &Trade{
			Exchange:    ASTER_NAME.String(),
			AccountType: req.AccountType,
			Symbol:      trade.Symbol,
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

func (b *AsterTradeEngine) handleOrderFromSpotOrderCreate(req *OrderParam, res *myasterapi.SpotOrderPostRes) *Order {
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
		Exchange:      ASTER_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        res.Symbol,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.ClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CummulativeQuoteQty,
		AvgPrice:      avgPrice.String(),
		Status:        b.asterConverter.FromAsterOrderStatus(res.Status, res.Type),
		Type:          b.asterConverter.FromAsterOrderType(res.Type),
		Side:          b.asterConverter.FromAsterOrderSide(res.Side),
		TimeInForce:   b.asterConverter.FromAsterTimeInForce(res.TimeInForce),
		CreateTime:    res.WorkingTime,
		UpdateTime:    res.WorkingTime,

		TriggerPrice:         req.TriggerPrice.String(),
		TriggerType:          b.asterConverter.FromAsterOrderTypeForTriggerType(res.Type),
		TriggerConditionType: b.asterConverter.FromAsterOrderSideForTriggerConditionType(res.Side, res.Type),
	}
	return order
}
func (b *AsterTradeEngine) handleOrderFromSpotOrderAmend(req *OrderParam, res *myasterapi.SpotOrderCancelReplaceRes) *Order {
	avgPrice := decimal.Zero
	if res.NewOrderResponse.ExecutedQty != "" && res.NewOrderResponse.CummulativeQuoteQty != "" {
		executedQty, _ := decimal.NewFromString(res.NewOrderResponse.ExecutedQty)
		cumQuoteQty, _ := decimal.NewFromString(res.NewOrderResponse.CummulativeQuoteQty)
		if !executedQty.IsZero() {
			avgPrice = cumQuoteQty.Div(executedQty)
		}
	}
	order := &Order{
		Exchange:      ASTER_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        res.NewOrderResponse.Symbol,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
		OrderId:       strconv.FormatInt(res.NewOrderResponse.OrderId, 10),
		ClientOrderId: res.NewOrderResponse.ClientOrderId,
		Price:         res.NewOrderResponse.Price,
		Quantity:      res.NewOrderResponse.OrigQty,
		ExecutedQty:   res.NewOrderResponse.ExecutedQty,
		CumQuoteQty:   res.NewOrderResponse.CummulativeQuoteQty,
		AvgPrice:      avgPrice.String(),
		Status:        b.asterConverter.FromAsterOrderStatus(res.NewOrderResponse.Status, res.NewOrderResponse.Type),
		Type:          b.asterConverter.FromAsterOrderType(res.NewOrderResponse.Type),
		Side:          b.asterConverter.FromAsterOrderSide(res.NewOrderResponse.Side),
		TimeInForce:   b.asterConverter.FromAsterTimeInForce(res.NewOrderResponse.TimeInForce),
		CreateTime:    res.NewOrderResponse.WorkingTime,
		UpdateTime:    res.NewOrderResponse.WorkingTime,

		TriggerPrice:         req.TriggerPrice.String(),
		TriggerType:          b.asterConverter.FromAsterOrderTypeForTriggerType(res.NewOrderResponse.Type),
		TriggerConditionType: b.asterConverter.FromAsterOrderSideForTriggerConditionType(res.NewOrderResponse.Side, res.NewOrderResponse.Type),
	}
	return order
}
func (b *AsterTradeEngine) handleOrderFromSpotOrderCancel(req *OrderParam, res *myasterapi.SpotOrderDeleteRes) *Order {
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
		Exchange:      ASTER_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        res.Symbol,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.OrigClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CummulativeQuoteQty,
		AvgPrice:      avgPrice.String(),
		Status:        b.asterConverter.FromAsterOrderStatus(res.Status, res.Type),
		Type:          b.asterConverter.FromAsterOrderType(res.Type),
		Side:          b.asterConverter.FromAsterOrderSide(res.Side),
		TimeInForce:   b.asterConverter.FromAsterTimeInForce(res.TimeInForce),
		CreateTime:    res.TransactTime,
		UpdateTime:    nowTimestamp,

		TriggerPrice:         req.TriggerPrice.String(),
		TriggerType:          b.asterConverter.FromAsterOrderTypeForTriggerType(res.Type),
		TriggerConditionType: b.asterConverter.FromAsterOrderSideForTriggerConditionType(res.Side, res.Type),
	}
	return order
}

func (b *AsterTradeEngine) handleOrderFromSpotBatchErr(req *OrderParam, err error) *Order {
	return &Order{
		Exchange:      ASTER_NAME.String(),
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

// U合约订单处理
func (b *AsterTradeEngine) handleOrdersFromFutureOpenOrders(req *QueryOrderParam, res *myasterapi.FutureOpenOrdersRes) []*Order {
	var orders []*Order
	for _, order := range *res {
		orders = append(orders, &Order{
			Exchange:      ASTER_NAME.String(),
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
			Status:        b.asterConverter.FromAsterOrderStatus(order.Status, order.Type),
			Type:          b.asterConverter.FromAsterOrderType(order.Type),
			Side:          b.asterConverter.FromAsterOrderSide(order.Side),
			PositionSide:  b.asterConverter.FromAsterPositionSide(order.PositionSide),
			TimeInForce:   b.asterConverter.FromAsterTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    order.Time,
			UpdateTime:    order.UpdateTime,

			TriggerPrice:         order.StopPrice,
			TriggerType:          b.asterConverter.FromAsterOrderTypeForTriggerType(order.Type),
			TriggerConditionType: b.asterConverter.FromAsterOrderSideForTriggerConditionType(order.Side, order.Type),
		})
	}
	return orders
}
func (b *AsterTradeEngine) handleOrderFromFutureOrderQuery(req *QueryOrderParam, res *myasterapi.FutureOrderGetRes) *Order {
	order := &Order{
		Exchange:      ASTER_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        res.Symbol,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.ClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CumQuote,
		AvgPrice:      res.AvgPrice,
		Status:        b.asterConverter.FromAsterOrderStatus(res.Status, res.Type),
		Type:          b.asterConverter.FromAsterOrderType(res.Type),
		Side:          b.asterConverter.FromAsterOrderSide(res.Side),
		PositionSide:  b.asterConverter.FromAsterPositionSide(res.PositionSide),
		TimeInForce:   b.asterConverter.FromAsterTimeInForce(res.TimeInForce),
		ReduceOnly:    res.ReduceOnly,
		CreateTime:    res.Time,
		UpdateTime:    res.UpdateTime,

		TriggerPrice:         res.StopPrice,
		TriggerType:          b.asterConverter.FromAsterOrderTypeForTriggerType(res.Type),
		TriggerConditionType: b.asterConverter.FromAsterOrderSideForTriggerConditionType(res.Side, res.Type),
	}
	return order
}
func (b *AsterTradeEngine) handleOrderFromFutureOrdersQuery(req *QueryOrderParam, res *myasterapi.FutureAllOrdersRes) []*Order {
	var orders []*Order
	for _, order := range *res {
		orders = append(orders, &Order{
			Exchange:      ASTER_NAME.String(),
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
			Status:        b.asterConverter.FromAsterOrderStatus(order.Status, order.Type),
			Type:          b.asterConverter.FromAsterOrderType(order.Type),
			Side:          b.asterConverter.FromAsterOrderSide(order.Side),
			PositionSide:  b.asterConverter.FromAsterPositionSide(order.PositionSide),
			TimeInForce:   b.asterConverter.FromAsterTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    order.Time,
			UpdateTime:    order.UpdateTime,

			TriggerPrice:         order.StopPrice,
			TriggerType:          b.asterConverter.FromAsterOrderTypeForTriggerType(order.Type),
			TriggerConditionType: b.asterConverter.FromAsterOrderSideForTriggerConditionType(order.Side, order.Type),
		})
	}
	return orders
}
func (b *AsterTradeEngine) handleTradesFromFutureTradeQuery(req *QueryTradeParam, res *myasterapi.FutureUserTradesRes) []*Trade {
	var trades []*Trade
	for _, trade := range *res {
		trades = append(trades, &Trade{
			Exchange:     ASTER_NAME.String(),
			AccountType:  req.AccountType,
			Symbol:       trade.Symbol,
			TradeId:      strconv.FormatInt(trade.Id, 10),
			OrderId:      strconv.FormatInt(trade.OrderId, 10),
			Price:        trade.Price,
			Quantity:     trade.Qty,
			QuoteQty:     trade.QuoteQty,
			Side:         b.asterConverter.FromAsterOrderSide(trade.Side),
			PositionSide: b.asterConverter.FromAsterPositionSide(trade.PositionSide),
			FeeAmount:    trade.Commission,
			FeeCcy:       trade.CommissionAsset,
			RealizedPnl:  trade.RealizedPnl,
			IsMaker:      trade.Maker,
			Timestamp:    trade.Time,
		})
	}
	return trades
}

func (b *AsterTradeEngine) handleOrderFromFutureOrderCreate(req *OrderParam, res *myasterapi.FutureOrderPostRes) *Order {
	nowTimestamp := time.Now().UnixMilli()
	order := &Order{
		Exchange:      ASTER_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        res.Symbol,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.ClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CumQuote,
		AvgPrice:      res.AvgPrice,
		Status:        b.asterConverter.FromAsterOrderStatus(res.Status, res.Type),
		Type:          b.asterConverter.FromAsterOrderType(res.Type),
		Side:          b.asterConverter.FromAsterOrderSide(res.Side),
		PositionSide:  b.asterConverter.FromAsterPositionSide(res.PositionSide),
		TimeInForce:   b.asterConverter.FromAsterTimeInForce(res.TimeInForce),
		ReduceOnly:    res.ReduceOnly,
		CreateTime:    nowTimestamp,
		UpdateTime:    nowTimestamp,

		TriggerPrice:         res.StopPrice,
		TriggerType:          b.asterConverter.FromAsterOrderTypeForTriggerType(res.Type),
		TriggerConditionType: b.asterConverter.FromAsterOrderSideForTriggerConditionType(res.Side, res.Type),
	}
	return order
}
func (b *AsterTradeEngine) handleOrderFromFutureOrderAmend(req *OrderParam, res *myasterapi.FutureOrderPutRes) *Order {
	nowTimestamp := time.Now().UnixMilli()
	order := &Order{
		Exchange:      ASTER_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        res.Symbol,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.ClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CumQuote,
		AvgPrice:      res.AvgPrice,
		Status:        b.asterConverter.FromAsterOrderStatus(res.Status, res.Type),
		Type:          b.asterConverter.FromAsterOrderType(res.Type),
		Side:          b.asterConverter.FromAsterOrderSide(res.Side),
		PositionSide:  b.asterConverter.FromAsterPositionSide(res.PositionSide),
		TimeInForce:   b.asterConverter.FromAsterTimeInForce(res.TimeInForce),
		ReduceOnly:    res.ReduceOnly,
		CreateTime:    nowTimestamp,
		UpdateTime:    nowTimestamp,

		TriggerPrice:         res.StopPrice,
		TriggerType:          b.asterConverter.FromAsterOrderTypeForTriggerType(res.Type),
		TriggerConditionType: b.asterConverter.FromAsterOrderSideForTriggerConditionType(res.Side, res.Type),
	}
	return order
}
func (b *AsterTradeEngine) handleOrderFromFutureOrderCancel(req *OrderParam, res *myasterapi.FutureOrderDeleteRes) *Order {
	order := &Order{
		Exchange:      ASTER_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        res.Symbol,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.ClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CumQuote,
		AvgPrice:      res.AvgPrice,
		Status:        b.asterConverter.FromAsterOrderStatus(res.Status, res.Type),
		Type:          b.asterConverter.FromAsterOrderType(res.Type),
		Side:          b.asterConverter.FromAsterOrderSide(res.Side),
		PositionSide:  b.asterConverter.FromAsterPositionSide(res.PositionSide),
		TimeInForce:   b.asterConverter.FromAsterTimeInForce(res.TimeInForce),
		ReduceOnly:    res.ReduceOnly,
		CreateTime:    res.UpdateTime,
		UpdateTime:    res.UpdateTime,

		TriggerPrice:         res.StopPrice,
		TriggerType:          b.asterConverter.FromAsterOrderTypeForTriggerType(res.Type),
		TriggerConditionType: b.asterConverter.FromAsterOrderSideForTriggerConditionType(res.Side, res.Type),
	}
	return order
}

func (b *AsterTradeEngine) handleOrdersFromFutureBatchOrderCreate(reqs []*OrderParam, res *myasterapi.FutureBatchOrdersPostRes) []*Order {
	var orders []*Order
	nowTimestamp := time.Now().UnixMilli()
	for _, order := range *res {
		code, ok := order.Code.(float64)
		if !ok {
			code = 0
		}
		codeInt := int(code)
		orders = append(orders, &Order{
			Exchange:      ASTER_NAME.String(),
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
			Status:        b.asterConverter.FromAsterOrderStatus(order.Status, order.Type),
			Type:          b.asterConverter.FromAsterOrderType(order.Type),
			Side:          b.asterConverter.FromAsterOrderSide(order.Side),
			PositionSide:  b.asterConverter.FromAsterPositionSide(order.PositionSide),
			TimeInForce:   b.asterConverter.FromAsterTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    nowTimestamp,
			UpdateTime:    nowTimestamp,
			ErrorCode:     strconv.Itoa(codeInt),
			ErrorMsg:      order.Msg,

			TriggerPrice:         order.StopPrice,
			TriggerType:          b.asterConverter.FromAsterOrderTypeForTriggerType(order.Type),
			TriggerConditionType: b.asterConverter.FromAsterOrderSideForTriggerConditionType(order.Side, order.Type),
		})
	}
	return orders
}
func (b *AsterTradeEngine) handleOrdersFromFutureBatchOrderAmend(reqs []*OrderParam, res *myasterapi.FutureBatchOrdersPutRes) []*Order {
	var orders []*Order
	nowTimestamp := time.Now().UnixMilli()
	for _, order := range *res {
		code, ok := order.Code.(float64)
		if !ok {
			code = 0
		}
		codeInt := int(code)
		orders = append(orders, &Order{
			Exchange:      ASTER_NAME.String(),
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
			Status:        b.asterConverter.FromAsterOrderStatus(order.Status, order.Type),
			Type:          b.asterConverter.FromAsterOrderType(order.Type),
			Side:          b.asterConverter.FromAsterOrderSide(order.Side),
			PositionSide:  b.asterConverter.FromAsterPositionSide(order.PositionSide),
			TimeInForce:   b.asterConverter.FromAsterTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    nowTimestamp,
			UpdateTime:    nowTimestamp,
			ErrorCode:     strconv.Itoa(codeInt),
			ErrorMsg:      order.Msg,

			TriggerPrice:         order.StopPrice,
			TriggerType:          b.asterConverter.FromAsterOrderTypeForTriggerType(order.Type),
			TriggerConditionType: b.asterConverter.FromAsterOrderSideForTriggerConditionType(order.Side, order.Type),
		})
	}
	return orders
}
func (b *AsterTradeEngine) handleOrdersFromFutureBatchOrderCancel(reqs []*OrderParam, res *myasterapi.FutureBatchOrdersDeleteRes) []*Order {
	var orders []*Order
	nowTimestamp := time.Now().UnixMilli()
	for _, order := range *res {
		code, ok := order.Code.(float64)
		if !ok {
			code = 0
		}
		codeInt := int(code)
		orders = append(orders, &Order{
			Exchange:      ASTER_NAME.String(),
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
			Status:        b.asterConverter.FromAsterOrderStatus(order.Status, order.Type),
			Type:          b.asterConverter.FromAsterOrderType(order.Type),
			Side:          b.asterConverter.FromAsterOrderSide(order.Side),
			PositionSide:  b.asterConverter.FromAsterPositionSide(order.PositionSide),
			TimeInForce:   b.asterConverter.FromAsterTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    nowTimestamp,
			UpdateTime:    nowTimestamp,
			ErrorCode:     strconv.Itoa(codeInt),
			ErrorMsg:      order.Msg,

			TriggerPrice:         order.StopPrice,
			TriggerType:          b.asterConverter.FromAsterOrderTypeForTriggerType(order.Type),
			TriggerConditionType: b.asterConverter.FromAsterOrderSideForTriggerConditionType(order.Side, order.Type),
		})
	}
	return orders
}

// handle ws
func (b *AsterTradeEngine) handleSubscribeOrderFromSpotPayload(req SubscribeOrderParam, newPayload *myasterapi.WsSpotPayload, newSub *subscription[Order]) {
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
					Exchange:      ASTER_NAME.String(),
					AccountType:   req.AccountType,
					Symbol:        r.Symbol,
					OrderId:       strconv.FormatInt(r.OrderId, 10),
					ClientOrderId: r.ClientOrderId,
					Price:         r.Price,
					Quantity:      r.OrigQty,
					ExecutedQty:   r.ExecutedQty,
					CumQuoteQty:   r.CummulativeQuoteQty,
					AvgPrice:      avgPrice.String(),
					Status:        b.asterConverter.FromAsterOrderStatus(r.Status, r.Type),
					Type:          b.asterConverter.FromAsterOrderType(r.Type),
					Side:          b.asterConverter.FromAsterOrderSide(r.Side),
					TimeInForce:   b.asterConverter.FromAsterTimeInForce(r.TimeInForce),
					FeeAmount:     r.FeeQty,
					FeeCcy:        r.FeeAsset,
					CreateTime:    r.OrderCreateTime,
					UpdateTime:    r.Timestamp,

					TriggerPrice:         r.StopPrice,
					TriggerType:          b.asterConverter.FromAsterOrderTypeForTriggerType(r.Type),
					TriggerConditionType: b.asterConverter.FromAsterOrderSideForTriggerConditionType(r.Side, r.Type),
				}
				newSub.resultChan <- order
			}
		}
	}()
}
func (b *AsterTradeEngine) handleSubscribeOrderFromFuturePayload(req SubscribeOrderParam, newPayload *myasterapi.WsFuturePayload, newSub *subscription[Order]) {
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
					Exchange:      ASTER_NAME.String(),
					AccountType:   req.AccountType,
					Symbol:        r.Symbol,
					OrderId:       strconv.FormatInt(r.OrderId, 10),
					ClientOrderId: r.ClientOrderId,
					Price:         r.Price,
					Quantity:      r.OrigQty,
					ExecutedQty:   r.ExecutedQty,
					CumQuoteQty:   CumQuoteQty.String(),
					AvgPrice:      r.AvgPrice,
					Status:        b.asterConverter.FromAsterOrderStatus(r.Status, r.Type),
					Type:          b.asterConverter.FromAsterOrderType(r.Type),
					Side:          b.asterConverter.FromAsterOrderSide(r.Side),
					PositionSide:  b.asterConverter.FromAsterPositionSide(r.PositionSide),
					TimeInForce:   b.asterConverter.FromAsterTimeInForce(r.TimeInForce),
					FeeAmount:     r.FeeQty,
					FeeCcy:        r.FeeAsset,
					ReduceOnly:    r.IsReduceOnly,
					CreateTime:    result.TradeTime,
					UpdateTime:    r.TradeTime,

					TriggerPrice:         r.StopPrice,
					TriggerType:          b.asterConverter.FromAsterOrderTypeForTriggerType(r.Type),
					TriggerConditionType: b.asterConverter.FromAsterOrderSideForTriggerConditionType(r.Side, r.Type),
				}
				newSub.resultChan <- order
			}
		}
	}()
}
