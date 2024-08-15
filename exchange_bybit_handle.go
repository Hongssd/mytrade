package mytrade

import (
	"errors"
	"github.com/Hongssd/mybybitapi"
	"github.com/shopspring/decimal"
)

// 查询订单结果处理
func (b *BybitTradeEngine) handleOrdersFromQueryOpenOrders(req *QueryOrderParam, res mybybitapi.OrderRealtimeRes) []*Order {
	var orders []*Order
	for _, order := range res.List {
		var isMargin bool
		if order.IsLeverage == "1" {
			isMargin = true
		}
		orders = append(orders, &Order{
			Exchange:      BYBIT_NAME.String(),
			AccountType:   req.AccountType,
			Symbol:        order.Symbol,
			IsMargin:      isMargin,
			OrderId:       order.OrderId,
			ClientOrderId: order.OrderLinkId,
			Price:         order.Price,
			Quantity:      order.Qty,
			ExecutedQty:   order.CumExecQty,
			CumQuoteQty:   order.CumExecValue,
			AvgPrice:      order.AvgPrice,
			Status:        b.bybitConverter.FromBYBITOrderStatus(order.OrderStatus),
			Type:          b.bybitConverter.FromBYBITOrderType(order.OrderType),
			Side:          b.bybitConverter.FromBYBITOrderSide(order.Side),
			PositionSide:  b.bybitConverter.FromBYBITPositionSide(order.PositionIdx),
			TimeInForce:   b.bybitConverter.FromBYBITTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    stringToInt64(order.CreatedTime),
			UpdateTime:    stringToInt64(order.UpdatedTime),

			FeeAmount: order.CumExecFee,
			FeeCcy:    "",

			TriggerPrice:         order.TriggerPrice,
			TriggerType:          b.bybitConverter.FromBYBITTriggerConditionForTriggerType(order.TriggerDirection, order.Side),
			TriggerConditionType: b.bybitConverter.FromBYBITTriggerCondition(order.TriggerDirection),
		})

	}
	return orders
}
func (b *BybitTradeEngine) handleOrdersFromQueryOrders(req *QueryOrderParam, res mybybitapi.OrderHistoryRes) []*Order {
	var orders []*Order
	for _, order := range res.List {
		orders = append(orders, &Order{
			Exchange:      BYBIT_NAME.String(),
			AccountType:   req.AccountType,
			Symbol:        order.Symbol,
			IsMargin:      req.IsMargin,
			IsIsolated:    req.IsIsolated,
			OrderId:       order.OrderId,
			ClientOrderId: order.OrderLinkId,
			Price:         order.Price,
			Quantity:      order.Qty,
			ExecutedQty:   order.CumExecQty,
			CumQuoteQty:   order.CumExecValue,
			AvgPrice:      order.AvgPrice,
			Status:        b.bybitConverter.FromBYBITOrderStatus(order.OrderStatus),
			Type:          b.bybitConverter.FromBYBITOrderType(order.OrderType),
			Side:          b.bybitConverter.FromBYBITOrderSide(order.Side),
			PositionSide:  b.bybitConverter.FromBYBITPositionSide(order.PositionIdx),
			TimeInForce:   b.bybitConverter.FromBYBITTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    stringToInt64(order.CreatedTime),
			UpdateTime:    stringToInt64(order.UpdatedTime),

			FeeAmount: order.CumExecFee,
			FeeCcy:    "",

			TriggerPrice:         order.TriggerPrice,
			TriggerType:          b.bybitConverter.FromBYBITTriggerConditionForTriggerType(order.TriggerDirection, order.Side),
			TriggerConditionType: b.bybitConverter.FromBYBITTriggerCondition(order.TriggerDirection),
		})
	}
	return orders
}
func (b *BybitTradeEngine) handleTradesFromQueryTrades(req *QueryTradeParam, res mybybitapi.OrderExecutionListRes) []*Trade {
	var trades []*Trade
	for _, r := range res.List {
		quoteQty := decimal.RequireFromString(r.ExecPrice).Mul(decimal.RequireFromString(r.ExecQty))
		trades = append(trades, &Trade{
			Exchange:     BYBIT_NAME.String(),
			AccountType:  req.AccountType,
			Symbol:       r.Symbol,
			TradeId:      r.ExecId,
			OrderId:      r.OrderId,
			Price:        r.ExecPrice,
			Quantity:     r.ExecQty,
			QuoteQty:     quoteQty.String(),
			Side:         b.bybitConverter.FromBYBITOrderSide(r.Side),
			PositionSide: "",
			FeeAmount:    r.ExecFee,
			FeeCcy:       r.FeeCurrency,
			RealizedPnl:  "",
			IsMaker:      r.IsMaker,
			Timestamp:    stringToInt64(r.ExecTime),
		})
	}
	return trades
}

// 批量订单返回结果处理
func (b *BybitTradeEngine) handleOrderFromBatchOrderCreate(reqs []*OrderParam, res *mybybitapi.BybitRestRes[mybybitapi.OrderCreateBatchRes]) ([]*Order, error) {
	if len(res.Result.List) != len(reqs) {
		return nil, errors.New("api return invalid data")
	}
	orders := make([]*Order, 0, len(reqs))
	for i, r := range res.Result.List {
		order := &Order{
			Exchange:      BYBIT_NAME.String(),
			OrderId:       r.OrderId,
			ClientOrderId: r.OrderLinkId,
			AccountType:   r.Category,
			Symbol:        r.Symbol,
			IsMargin:      reqs[i].IsMargin,
			IsIsolated:    reqs[i].IsIsolated,
			CreateTime:    stringToInt64(r.CreateAt),
		}
		orders = append(orders, order)
	}
	return orders, nil
}
func (b *BybitTradeEngine) handleOrderFromBatchOrderAmend(reqs []*OrderParam, res *mybybitapi.BybitRestRes[mybybitapi.OrderAmendBatchRes]) ([]*Order, error) {
	if len(res.Result.List) != len(reqs) {
		return nil, errors.New("api return invalid data")
	}
	orders := make([]*Order, 0, len(reqs))
	for i, r := range res.Result.List {
		order := &Order{
			Exchange:      BYBIT_NAME.String(),
			OrderId:       r.OrderId,
			ClientOrderId: r.OrderLinkId,
			AccountType:   r.Category,
			Symbol:        r.Symbol,
			IsMargin:      reqs[i].IsMargin,
			IsIsolated:    reqs[i].IsIsolated,
		}
		orders = append(orders, order)
	}
	return orders, nil
}
func (b *BybitTradeEngine) handleOrderFromBatchOrderCancel(reqs []*OrderParam, res *mybybitapi.BybitRestRes[mybybitapi.OrderCancelBatchRes]) ([]*Order, error) {
	if len(res.Result.List) != len(reqs) {
		return nil, errors.New("api return invalid data")
	}
	orders := make([]*Order, 0, len(reqs))
	for i, r := range res.Result.List {
		order := &Order{
			Exchange:      BYBIT_NAME.String(),
			OrderId:       r.OrderId,
			ClientOrderId: r.OrderLinkId,
			AccountType:   r.Category,
			Symbol:        r.Symbol,
			IsMargin:      reqs[i].IsMargin,
			IsIsolated:    reqs[i].IsIsolated,
		}
		orders = append(orders, order)
	}
	return orders, nil
}

// 订单推送处理
func (b *BybitTradeEngine) handleOrderFromWsOrder(orders mybybitapi.WsOrder) []*Order {

	// 从ws订单信息转换为本地订单信息
	var res []*Order
	for _, order := range orders.Data {
		var isMargin, isIsolated bool
		if order.IsLeverage == "1" {
			isMargin = true
			isIsolated = false
		}
		order := &Order{
			Exchange:      BYBIT_NAME.String(),
			AccountType:   order.Category,
			Symbol:        order.Symbol,
			OrderId:       order.OrderId,
			ClientOrderId: order.OrderLinkId,
			Price:         order.Price,
			Quantity:      order.Qty,
			ExecutedQty:   order.CumExecQty,
			CumQuoteQty:   order.CumExecValue,
			AvgPrice:      order.AvgPrice,
			Status:        b.bybitConverter.FromBYBITOrderStatus(order.OrderStatus),
			Type:          b.bybitConverter.FromBYBITOrderType(order.OrderType),
			Side:          b.bybitConverter.FromBYBITOrderSide(order.Side),
			PositionSide:  b.bybitConverter.FromBYBITPositionSide(order.PositionIdx),
			TimeInForce:   b.bybitConverter.FromBYBITTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    stringToInt64(order.CreatedTime),
			UpdateTime:    stringToInt64(order.UpdatedTime),

			FeeAmount: order.CumExecFee,
			FeeCcy:    order.FeeCurrency,

			TriggerPrice:         order.TriggerPrice,
			TriggerType:          b.bybitConverter.FromBYBITTriggerConditionForTriggerType(order.TriggerDirection, order.Side),
			TriggerConditionType: b.bybitConverter.FromBYBITTriggerCondition(order.TriggerDirection),

			IsIsolated: isIsolated,
			IsMargin:   isMargin,
		}
		res = append(res, order)
	}
	return res
}

func (b *BybitTradeEngine) handleOrderFromInverseBatchErr(req *OrderParam, err error) *Order {
	return &Order{
		Exchange:      BYBIT_NAME.String(),
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
