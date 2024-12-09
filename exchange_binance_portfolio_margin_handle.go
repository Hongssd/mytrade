package mytrade

import (
	"github.com/Hongssd/mybinanceapi"
	"github.com/shopspring/decimal"
	"strconv"
	"time"
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
		CreateTime:    order.UpdateTime,
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
			CreateTime:    order.UpdateTime,
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
			CreateTime:    order.UpdateTime,
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
		CreateTime:    order.UpdateTime,
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
			CreateTime:    order.UpdateTime,
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
