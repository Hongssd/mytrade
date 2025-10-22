package mytrade

import (
	"strconv"
	"time"

	"github.com/Hongssd/mybinanceapi"
	"github.com/shopspring/decimal"
)

// UM
func (b *BinanceTradeEngine) handlePortfolioMarginUmOrderCreate(req *OrderParam, order *mybinanceapi.PortfolioMarginUmOrderPostRes) *Order {
	return &Order{
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
		CumQuoteQty:   order.CumQty,
		AvgPrice:      order.AvgPrice,
		Status:        b.bnConverter.FromBNOrderStatus(order.Status, order.Type),
		Type:          b.bnConverter.FromBNOrderType(order.Type),
		Side:          b.bnConverter.FromBNOrderSide(order.Side),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
		PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
		CreateTime:    order.UpdateTime,
		UpdateTime:    order.UpdateTime,

		TriggerPrice:         req.TriggerPrice.String(),
		TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(order.Type),
		TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(order.Side, order.Type),
	}
}
func (b *BinanceTradeEngine) handlePortfolioMarginUmOrderAmend(req *OrderParam, order *mybinanceapi.PortfolioMarginUmOrderPutRes) *Order {
	return &Order{
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
		CumQuoteQty:   order.CumQty,
		AvgPrice:      order.AvgPrice,
		Status:        b.bnConverter.FromBNOrderStatus(order.Status, order.Type),
		Type:          b.bnConverter.FromBNOrderType(order.Type),
		Side:          b.bnConverter.FromBNOrderSide(order.Side),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
		PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
		CreateTime:    order.UpdateTime,
		UpdateTime:    order.UpdateTime,

		TriggerPrice:         req.TriggerPrice.String(),
		TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(order.Type),
		TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(order.Side, order.Type),
	}
}
func (b *BinanceTradeEngine) handlePortfolioMarginUmOrderCancel(req *OrderParam, order *mybinanceapi.PortfolioMarginUmOrderDeleteRes) *Order {
	return &Order{
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
		CumQuoteQty:   order.CumQty,
		AvgPrice:      order.AvgPrice,
		Status:        b.bnConverter.FromBNOrderStatus(order.Status, order.Type),
		Type:          b.bnConverter.FromBNOrderType(order.Type),
		Side:          b.bnConverter.FromBNOrderSide(order.Side),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
		PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
		CreateTime:    order.UpdateTime,
		UpdateTime:    order.UpdateTime,

		TriggerPrice:         req.TriggerPrice.String(),
		TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(order.Type),
		TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(order.Side, order.Type),
	}
}

func (b *BinanceTradeEngine) handlePortfolioMarginUmOpenOrders(req *QueryOrderParam, res *mybinanceapi.PortfolioMarginUmOpenOrdersGetRes) []*Order {
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
			TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
			PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
			CreateTime:    order.Time,
			UpdateTime:    order.UpdateTime,

			TriggerPrice:         order.Price, // 返回值无相关参数，price代替
			TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(order.Type),
			TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(order.Side, order.Type),
		})
	}
	return orders
}
func (b *BinanceTradeEngine) handlePortfolioMarginUmOrderQuery(req *QueryOrderParam, order *mybinanceapi.PortfolioMarginUmOrderGetRes) *Order {
	return &Order{
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
		TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
		PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
		CreateTime:    order.Time,
		UpdateTime:    order.UpdateTime,

		TriggerPrice:         order.Price, // 返回值无相关参数，price代替
		TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(order.Type),
		TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(order.Side, order.Type),
	}
}
func (b *BinanceTradeEngine) handlePortfolioMarginUmOrdersQuery(req *QueryOrderParam, res *mybinanceapi.PortfolioMarginUmAllOrdersGetRes) []*Order {
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
			TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
			PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
			CreateTime:    order.Time,
			UpdateTime:    order.UpdateTime,

			TriggerPrice:         order.Price, // 返回值无相关参数，price代替
			TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(order.Type),
			TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(order.Side, order.Type),
		})
	}
	return orders
}
func (b *BinanceTradeEngine) handlePortfolioMarginUmTradesQuery(req *QueryTradeParam, res *mybinanceapi.PortfolioMarginUmUserTradesRes) []*Trade {
	var trades []*Trade
	for _, trade := range *res {
		if req.OrderId != "" && strconv.FormatInt(trade.OrderId, 10) != req.OrderId {
			continue
		}
		trades = append(trades, &Trade{
			Exchange:     BINANCE_NAME.String(),
			AccountType:  req.AccountType,
			Symbol:       trade.Symbol,
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

// CM
func (b *BinanceTradeEngine) handlePortfolioMarginCmOrderCreate(req *OrderParam, order *mybinanceapi.PortfolioMarginCmOrderPostRes) *Order {
	return &Order{
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
		CumQuoteQty:   order.CumQty,
		AvgPrice:      order.AvgPrice,
		Status:        b.bnConverter.FromBNOrderStatus(order.Status, order.Type),
		Type:          b.bnConverter.FromBNOrderType(order.Type),
		Side:          b.bnConverter.FromBNOrderSide(order.Side),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
		PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
		CreateTime:    order.UpdateTime,
		UpdateTime:    order.UpdateTime,

		TriggerPrice:         req.TriggerPrice.String(),
		TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(order.Type),
		TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(order.Side, order.Type),
	}
}
func (b *BinanceTradeEngine) handlePortfolioMarginCmOrderAmend(req *OrderParam, order *mybinanceapi.PortfolioMarginCmOrderPutRes) *Order {
	return &Order{
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
		CumQuoteQty:   order.CumQty,
		AvgPrice:      order.AvgPrice,
		Status:        b.bnConverter.FromBNOrderStatus(order.Status, order.Type),
		Type:          b.bnConverter.FromBNOrderType(order.Type),
		Side:          b.bnConverter.FromBNOrderSide(order.Side),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
		PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
		CreateTime:    order.UpdateTime,
		UpdateTime:    order.UpdateTime,

		TriggerPrice:         req.TriggerPrice.String(),
		TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(order.Type),
		TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(order.Side, order.Type),
	}
}
func (b *BinanceTradeEngine) handlePortfolioMarginCmOrderCancel(req *OrderParam, order *mybinanceapi.PortfolioMarginCmOrderDeleteRes) *Order {
	return &Order{
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
		CumQuoteQty:   order.CumQty,
		AvgPrice:      order.AvgPrice,
		Status:        b.bnConverter.FromBNOrderStatus(order.Status, order.Type),
		Type:          b.bnConverter.FromBNOrderType(order.Type),
		Side:          b.bnConverter.FromBNOrderSide(order.Side),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
		PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
		CreateTime:    order.UpdateTime,
		UpdateTime:    order.UpdateTime,

		TriggerPrice:         req.TriggerPrice.String(),
		TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(order.Type),
		TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(order.Side, order.Type),
	}
}

func (b *BinanceTradeEngine) handlePortfolioMarginCmOpenOrders(req *QueryOrderParam, res *mybinanceapi.PortfolioMarginCmOpenOrdersRes) []*Order {
	var orders []*Order
	for _, order := range *res {
		cumQuoteQty := decimal.RequireFromString(order.AvgPrice).Mul(decimal.RequireFromString(order.CumBase))
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
			CumQuoteQty:   cumQuoteQty.String(),
			AvgPrice:      order.AvgPrice,
			Status:        b.bnConverter.FromBNOrderStatus(order.Status, order.Type),
			Type:          b.bnConverter.FromBNOrderType(order.Type),
			Side:          b.bnConverter.FromBNOrderSide(order.Side),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
			PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
			CreateTime:    order.Time,
			UpdateTime:    order.UpdateTime,

			TriggerPrice:         order.Price, // 返回值无相关参数，price代替
			TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(order.Type),
			TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(order.Side, order.Type),
		})
	}
	return orders
}
func (b *BinanceTradeEngine) handlePortfolioMarginCmOrderQuery(req *QueryOrderParam, order *mybinanceapi.PortfolioMarginCmOrderGetRes) *Order {
	cumQuoteQty := decimal.RequireFromString(order.AvgPrice).Mul(decimal.RequireFromString(order.CumBase))
	return &Order{
		Exchange:      BINANCE_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        order.Symbol,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
		OrderId:       strconv.FormatInt(order.OrderId, 10),
		ClientOrderId: order.ClientOrderId,
		Price:         order.Price,
		Quantity:      order.OrigQty,
		ExecutedQty:   order.ExecuteQty,
		CumQuoteQty:   cumQuoteQty.String(),
		AvgPrice:      order.AvgPrice,
		Status:        b.bnConverter.FromBNOrderStatus(order.Status, order.Type),
		Type:          b.bnConverter.FromBNOrderType(order.Type),
		Side:          b.bnConverter.FromBNOrderSide(order.Side),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
		PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
		CreateTime:    order.Time,
		UpdateTime:    order.UpdateTime,

		TriggerPrice:         order.Price, // 返回值无相关参数，price代替
		TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(order.Type),
		TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(order.Side, order.Type),
	}
}
func (b *BinanceTradeEngine) handlePortfolioMarginCmOrdersQuery(req *QueryOrderParam, res *mybinanceapi.PortfolioMarginCmAllOrdersRes) []*Order {
	var orders []*Order
	for _, order := range *res {
		cumQuoteQty := decimal.RequireFromString(order.AvgPrice).Mul(decimal.RequireFromString(order.CumBase))
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
			ExecutedQty:   order.ExecuteQty,
			CumQuoteQty:   cumQuoteQty.String(),
			AvgPrice:      order.AvgPrice,
			Status:        b.bnConverter.FromBNOrderStatus(order.Status, order.Type),
			Type:          b.bnConverter.FromBNOrderType(order.Type),
			Side:          b.bnConverter.FromBNOrderSide(order.Side),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
			PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
			CreateTime:    order.Time,
			UpdateTime:    order.UpdateTime,

			TriggerPrice:         order.Price, // 返回值无相关参数，price代替
			TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(order.Type),
			TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(order.Side, order.Type),
		})
	}
	return orders
}
func (b *BinanceTradeEngine) handlePortfolioMarginCmTradesQuery(req *QueryTradeParam, res *mybinanceapi.PortfolioMarginCmUserTradesRes) []*Trade {
	var trades []*Trade
	for _, trade := range *res {
		quoteQty := decimal.RequireFromString(trade.Price).Mul(decimal.RequireFromString(trade.BaseQty))
		trades = append(trades, &Trade{
			Exchange:     BINANCE_NAME.String(),
			AccountType:  req.AccountType,
			Symbol:       trade.Symbol,
			TradeId:      strconv.FormatInt(trade.Id, 10),
			OrderId:      strconv.FormatInt(trade.OrderId, 10),
			Price:        trade.Price,
			Quantity:     trade.Qty,
			QuoteQty:     quoteQty.String(),
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

// Margin
func (b *BinanceTradeEngine) handlePortfolioMarginMarginOrderCreate(req *OrderParam, order *mybinanceapi.PortfolioMarginMarginOrderPostRes) *Order {
	avgPrice := decimal.Zero
	if order.ExecutedQty != "" && order.CummulativeQuoteQty != "" {
		executedQty, _ := decimal.NewFromString(order.ExecutedQty)
		cumQuoteQty, _ := decimal.NewFromString(order.CummulativeQuoteQty)
		if !executedQty.IsZero() {
			avgPrice = cumQuoteQty.Div(executedQty)
		}
	}
	return &Order{
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
		CreateTime:    order.TransactTime,
		UpdateTime:    order.TransactTime,

		MarginBuyBorrowAmount: order.MarginBuyBorrowAmount,
		MarginBuyBorrowAsset:  order.MarginBuyBorrowAsset,

		TriggerPrice:         req.TriggerPrice.String(),
		TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(order.Type),
		TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(order.Side, order.Type),
	}
}
func (b *BinanceTradeEngine) handlePortfolioMarginMarginOrderCancel(req *OrderParam, order *mybinanceapi.PortfolioMarginMarginOrderDeleteRes) *Order {
	avgPrice := decimal.Zero
	if order.ExecutedQty != "" && order.CummulativeQuoteQty != "" {
		executedQty, _ := decimal.NewFromString(order.ExecutedQty)
		cumQuoteQty, _ := decimal.NewFromString(order.CummulativeQuoteQty)
		if !executedQty.IsZero() {
			avgPrice = cumQuoteQty.Div(executedQty)
		}
	}
	timestamp := time.Now().UnixMilli()
	return &Order{
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
		CreateTime:    timestamp, // 返回值无相关参数，使用当前系统时间代替
		UpdateTime:    timestamp,

		TriggerPrice:         req.TriggerPrice.String(),
		TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(order.Type),
		TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(order.Side, order.Type),
	}
}

func (b *BinanceTradeEngine) handlePortfolioMarginMarginOpenOrders(req *QueryOrderParam, res *mybinanceapi.PortfolioMarginMarginOpenOrdersRes) []*Order {
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
func (b *BinanceTradeEngine) handlePortfolioMarginMarginOrderQuery(req *QueryOrderParam, order *mybinanceapi.PortfolioMarginMarginOrderGetRes) *Order {
	avgPrice := decimal.Zero
	if order.ExecutedQty != "" && order.CummulativeQuoteQty != "" {
		executedQty, _ := decimal.NewFromString(order.ExecutedQty)
		cumQuoteQty, _ := decimal.NewFromString(order.CummulativeQuoteQty)
		if !executedQty.IsZero() {
			avgPrice = cumQuoteQty.Div(executedQty)
		}
	}
	return &Order{
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
	}
}
func (b *BinanceTradeEngine) handlePortfolioMarginMarginOrdersQuery(req *QueryOrderParam, res *mybinanceapi.PortfolioMarginMarginAllOrdersRes) []*Order {
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
func (b *BinanceTradeEngine) handlePortfolioMarginMarginTradesQuery(req *QueryTradeParam, res *mybinanceapi.PortfolioMarginMarginMyTradesRes) []*Trade {
	var trades []*Trade
	for _, trade := range *res {
		var orderSide OrderSide
		if trade.IsBuyer {
			orderSide = ORDER_SIDE_BUY
		} else {
			orderSide = ORDER_SIDE_SELL
		}
		quoteQty := decimal.RequireFromString(trade.Price).Mul(decimal.RequireFromString(trade.Qty))
		trades = append(trades, &Trade{
			Exchange:    BINANCE_NAME.String(),
			AccountType: req.AccountType,
			Symbol:      trade.Symbol,
			TradeId:     strconv.FormatInt(trade.Id, 10),
			OrderId:     strconv.FormatInt(trade.OrderId, 10),
			Price:       trade.Price,
			Quantity:    trade.Qty,
			QuoteQty:    quoteQty.String(),
			Side:        orderSide,
			FeeAmount:   trade.Commission,
			FeeCcy:      trade.CommissionAsset,
			IsMaker:     trade.IsMaker,
			Timestamp:   trade.Time,
		})
	}
	return trades
}

// handle ws
func (b *BinanceTradeEngine) handleSubscribeOrderFromPMMarginPayload(req SubscribeOrderParam, newPayload *mybinanceapi.WsPMMarginPayload, newSub *subscription[Order]) {
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
					UpdateTime:    r.TradeTime,
					IsMargin:      true,

					TriggerPrice:         r.StopPrice,
					TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(r.Type),
					TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(r.Side, r.Type),
				}
				newSub.resultChan <- order
			}
		}
	}()
}
func (b *BinanceTradeEngine) handleSubscribeOrderFromPMContractPayload(req SubscribeOrderParam, newPayload *mybinanceapi.WsPMContractPayload, newSub *subscription[Order]) {
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
