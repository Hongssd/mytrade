package mytrade

import (
	"github.com/Hongssd/mybinanceapi"
	"github.com/shopspring/decimal"
	"strconv"
	"time"
)

// UM
func (b *BinanceTradeEngine) handlePortfolioMarginUmOrderCreate(req *OrderParam, res *mybinanceapi.PortfolioMarginUmOrderPostRes) *Order {
	return &Order{
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
		CumQuoteQty:   res.CumQty,
		AvgPrice:      res.AvgPrice,
		Status:        b.bnConverter.FromBNOrderStatus(res.Status, res.Type),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		CreateTime:    res.UpdateTime,
		UpdateTime:    res.UpdateTime,

		TriggerPrice:         req.TriggerPrice.String(),
		TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(res.Type),
		TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(res.Side, res.Type),
	}
}
func (b *BinanceTradeEngine) handlePortfolioMarginUmOrderAmend(req *OrderParam, res *mybinanceapi.PortfolioMarginUmOrderPutRes) *Order {
	return &Order{
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
		CumQuoteQty:   res.CumQty,
		AvgPrice:      res.AvgPrice,
		Status:        b.bnConverter.FromBNOrderStatus(res.Status, res.Type),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		CreateTime:    res.UpdateTime,
		UpdateTime:    res.UpdateTime,

		TriggerPrice:         req.TriggerPrice.String(),
		TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(res.Type),
		TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(res.Side, res.Type),
	}
}
func (b *BinanceTradeEngine) handlePortfolioMarginUmOrderCancel(req *OrderParam, res *mybinanceapi.PortfolioMarginUmOrderDeleteRes) *Order {
	return &Order{
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
		CumQuoteQty:   res.CumQty,
		AvgPrice:      res.AvgPrice,
		Status:        b.bnConverter.FromBNOrderStatus(res.Status, res.Type),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		CreateTime:    res.UpdateTime,
		UpdateTime:    res.UpdateTime,

		TriggerPrice:         req.TriggerPrice.String(),
		TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(res.Type),
		TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(res.Side, res.Type),
	}
}

func (b *BinanceTradeEngine) handlePortfolioMarginUmOpenOrders(req *QueryOrderParam, res *mybinanceapi.PortfolioMarginUmOpenOrdersGetRes) []*Order {
	var orders []*Order
	for _, o := range *res {
		orders = append(orders, &Order{
			Exchange:      BINANCE_NAME.String(),
			AccountType:   req.AccountType,
			Symbol:        req.Symbol,
			IsMargin:      req.IsMargin,
			IsIsolated:    req.IsIsolated,
			OrderId:       strconv.FormatInt(o.OrderId, 10),
			ClientOrderId: o.ClientOrderId,
			Price:         o.Price,
			Quantity:      o.OrigQty,
			ExecutedQty:   o.ExecutedQty,
			CumQuoteQty:   o.CumQuote,
			AvgPrice:      o.AvgPrice,
			Status:        b.bnConverter.FromBNOrderStatus(o.Status, o.Type),
			Type:          b.bnConverter.FromBNOrderType(o.Type),
			Side:          b.bnConverter.FromBNOrderSide(o.Side),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(o.TimeInForce),
			CreateTime:    o.Time,
			UpdateTime:    o.UpdateTime,

			TriggerPrice:         o.Price, // 返回值无相关参数，price代替
			TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(o.Type),
			TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(o.Side, o.Type),
		})
	}
	return orders
}
func (b *BinanceTradeEngine) handlePortfolioMarginUmOrderQuery(req *QueryOrderParam, res *mybinanceapi.PortfolioMarginUmOrderGetRes) *Order {
	return &Order{
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
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		CreateTime:    res.UpdateTime,
		UpdateTime:    res.UpdateTime,

		TriggerPrice:         res.Price, // 返回值无相关参数，price代替
		TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(res.Type),
		TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(res.Side, res.Type),
	}
}
func (b *BinanceTradeEngine) handlePortfolioMarginUmOrdersQuery(req *QueryOrderParam, res *mybinanceapi.PortfolioMarginUmAllOrdersGetRes) []*Order {
	var orders []*Order
	for _, o := range *res {
		orders = append(orders, &Order{
			Exchange:      BINANCE_NAME.String(),
			AccountType:   req.AccountType,
			Symbol:        req.Symbol,
			IsMargin:      req.IsMargin,
			IsIsolated:    req.IsIsolated,
			OrderId:       strconv.FormatInt(o.OrderId, 10),
			ClientOrderId: o.ClientOrderId,
			Price:         o.Price,
			Quantity:      o.OrigQty,
			ExecutedQty:   o.ExecutedQty,
			CumQuoteQty:   o.CumQuote,
			AvgPrice:      o.AvgPrice,
			Status:        b.bnConverter.FromBNOrderStatus(o.Status, o.Type),
			Type:          b.bnConverter.FromBNOrderType(o.Type),
			Side:          b.bnConverter.FromBNOrderSide(o.Side),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(o.TimeInForce),
			CreateTime:    o.UpdateTime,
			UpdateTime:    o.UpdateTime,

			TriggerPrice:         o.Price, // 返回值无相关参数，price代替
			TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(o.Type),
			TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(o.Side, o.Type),
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

// CM
func (b *BinanceTradeEngine) handlePortfolioMarginCmOrderCreate(req *OrderParam, res *mybinanceapi.PortfolioMarginCmOrderPostRes) *Order {
	return &Order{
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
		CumQuoteQty:   res.CumQty,
		AvgPrice:      res.AvgPrice,
		Status:        b.bnConverter.FromBNOrderStatus(res.Status, res.Type),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		CreateTime:    res.UpdateTime,
		UpdateTime:    res.UpdateTime,

		TriggerPrice:         req.TriggerPrice.String(),
		TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(res.Type),
		TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(res.Side, res.Type),
	}
}
func (b *BinanceTradeEngine) handlePortfolioMarginCmOrderAmend(req *OrderParam, res *mybinanceapi.PortfolioMarginCmOrderPutRes) *Order {
	return &Order{
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
		CumQuoteQty:   res.CumQty,
		AvgPrice:      res.AvgPrice,
		Status:        b.bnConverter.FromBNOrderStatus(res.Status, res.Type),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		CreateTime:    res.UpdateTime,
		UpdateTime:    res.UpdateTime,

		TriggerPrice:         req.TriggerPrice.String(),
		TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(res.Type),
		TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(res.Side, res.Type),
	}
}
func (b *BinanceTradeEngine) handlePortfolioMarginCmOrderCancel(req *OrderParam, res *mybinanceapi.PortfolioMarginCmOrderDeleteRes) *Order {
	return &Order{
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
		CumQuoteQty:   res.CumQty,
		AvgPrice:      res.AvgPrice,
		Status:        b.bnConverter.FromBNOrderStatus(res.Status, res.Type),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		CreateTime:    res.UpdateTime,
		UpdateTime:    res.UpdateTime,

		TriggerPrice:         req.TriggerPrice.String(),
		TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(res.Type),
		TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(res.Side, res.Type),
	}
}

func (b *BinanceTradeEngine) handlePortfolioMarginCmOpenOrders(req *QueryOrderParam, res *mybinanceapi.PortfolioMarginCmOpenOrdersRes) []*Order {
	var orders []*Order
	for _, o := range *res {
		cumQuoteQty := decimal.RequireFromString(o.Price).Mul(decimal.RequireFromString(o.ExecutedQty))
		orders = append(orders, &Order{
			Exchange:      BINANCE_NAME.String(),
			AccountType:   req.AccountType,
			Symbol:        req.Symbol,
			IsMargin:      req.IsMargin,
			IsIsolated:    req.IsIsolated,
			OrderId:       strconv.FormatInt(o.OrderId, 10),
			ClientOrderId: o.ClientOrderId,
			Price:         o.Price,
			Quantity:      o.OrigQty,
			ExecutedQty:   o.ExecutedQty,
			CumQuoteQty:   cumQuoteQty.String(),
			AvgPrice:      o.AvgPrice,
			Status:        b.bnConverter.FromBNOrderStatus(o.Status, o.Type),
			Type:          b.bnConverter.FromBNOrderType(o.Type),
			Side:          b.bnConverter.FromBNOrderSide(o.Side),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(o.TimeInForce),
			CreateTime:    o.UpdateTime,
			UpdateTime:    o.UpdateTime,

			TriggerPrice:         o.Price, // 返回值无相关参数，price代替
			TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(o.Type),
			TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(o.Side, o.Type),
		})
	}
	return orders
}
func (b *BinanceTradeEngine) handlePortfolioMarginCmOrderQuery(req *QueryOrderParam, res *mybinanceapi.PortfolioMarginCmOrderGetRes) *Order {
	cumQuoteQty := decimal.RequireFromString(res.Price).Mul(decimal.RequireFromString(res.ExecuteQty))
	return &Order{
		Exchange:      BINANCE_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.ClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecuteQty,
		CumQuoteQty:   cumQuoteQty.String(),
		AvgPrice:      res.AvgPrice,
		Status:        b.bnConverter.FromBNOrderStatus(res.Status, res.Type),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		CreateTime:    res.UpdateTime,
		UpdateTime:    res.UpdateTime,

		TriggerPrice:         res.Price, // 返回值无相关参数，price代替
		TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(res.Type),
		TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(res.Side, res.Type),
	}
}
func (b *BinanceTradeEngine) handlePortfolioMarginCmOrdersQuery(req *QueryOrderParam, res *mybinanceapi.PortfolioMarginCmAllOrdersRes) []*Order {
	var orders []*Order
	for _, o := range *res {
		cumQuoteQty := decimal.RequireFromString(o.Price).Mul(decimal.RequireFromString(o.ExecuteQty))
		orders = append(orders, &Order{
			Exchange:      BINANCE_NAME.String(),
			AccountType:   req.AccountType,
			Symbol:        req.Symbol,
			IsMargin:      req.IsMargin,
			IsIsolated:    req.IsIsolated,
			OrderId:       strconv.FormatInt(o.OrderId, 10),
			ClientOrderId: o.ClientOrderId,
			Price:         o.Price,
			Quantity:      o.OrigQty,
			ExecutedQty:   o.ExecuteQty,
			CumQuoteQty:   cumQuoteQty.String(),
			AvgPrice:      o.AvgPrice,
			Status:        b.bnConverter.FromBNOrderStatus(o.Status, o.Type),
			Type:          b.bnConverter.FromBNOrderType(o.Type),
			Side:          b.bnConverter.FromBNOrderSide(o.Side),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(o.TimeInForce),
			CreateTime:    o.UpdateTime,
			UpdateTime:    o.UpdateTime,

			TriggerPrice:         o.Price, // 返回值无相关参数，price代替
			TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(o.Type),
			TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(o.Side, o.Type),
		})
	}
	return orders
}
func (b *BinanceTradeEngine) handlePortfolioMarginCmTradesQuery(req *QueryTradeParam, res *mybinanceapi.PortfolioMarginCmUserTradesRes) []*Trade {
	var trades []*Trade
	for _, trade := range *res {
		quoteQty := decimal.RequireFromString(trade.Price).Mul(decimal.RequireFromString(trade.Qty))
		trades = append(trades, &Trade{
			Exchange:     BINANCE_NAME.String(),
			AccountType:  req.AccountType,
			Symbol:       req.Symbol,
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
func (b *BinanceTradeEngine) handlePortfolioMarginMarginOrderCreate(req *OrderParam, res *mybinanceapi.PortfolioMarginMarginOrderPostRes) *Order {
	avgPrice := decimal.Zero
	if res.ExecutedQty != "" && res.CummulativeQuoteQty != "" {
		executedQty, _ := decimal.NewFromString(res.ExecutedQty)
		cumQuoteQty, _ := decimal.NewFromString(res.CummulativeQuoteQty)
		if !executedQty.IsZero() {
			avgPrice = cumQuoteQty.Div(executedQty)
		}
	}
	return &Order{
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
		CreateTime:    res.TransactTime,
		UpdateTime:    res.TransactTime,

		TriggerPrice:         req.TriggerPrice.String(),
		TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(res.Type),
		TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(res.Side, res.Type),
	}
}
func (b *BinanceTradeEngine) handlePortfolioMarginMarginOrderCancel(req *OrderParam, res *mybinanceapi.PortfolioMarginMarginOrderDeleteRes) *Order {
	avgPrice := decimal.Zero
	if res.ExecutedQty != "" && res.CummulativeQuoteQty != "" {
		executedQty, _ := decimal.NewFromString(res.ExecutedQty)
		cumQuoteQty, _ := decimal.NewFromString(res.CummulativeQuoteQty)
		if !executedQty.IsZero() {
			avgPrice = cumQuoteQty.Div(executedQty)
		}
	}
	timestamp := time.Now().UnixMilli()
	return &Order{
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
		CreateTime:    timestamp, // 返回值无相关参数，使用当前系统时间代替
		UpdateTime:    timestamp,

		TriggerPrice:         req.TriggerPrice.String(),
		TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(res.Type),
		TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(res.Side, res.Type),
	}
}

func (b *BinanceTradeEngine) handlePortfolioMarginMarginOpenOrders(req *QueryOrderParam, res *mybinanceapi.PortfolioMarginMarginOpenOrdersRes) []*Order {
	var orders []*Order
	for _, o := range *res {
		avgPrice := decimal.Zero
		if o.ExecutedQty != "" && o.CummulativeQuoteQty != "" {
			executedQty, _ := decimal.NewFromString(o.ExecutedQty)
			cumQuoteQty, _ := decimal.NewFromString(o.CummulativeQuoteQty)
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
			OrderId:       strconv.FormatInt(o.OrderId, 10),
			ClientOrderId: o.ClientOrderId,
			Price:         o.Price,
			Quantity:      o.OrigQty,
			ExecutedQty:   o.ExecutedQty,
			CumQuoteQty:   o.CummulativeQuoteQty,
			AvgPrice:      avgPrice.String(),
			Status:        b.bnConverter.FromBNOrderStatus(o.Status, o.Type),
			Type:          b.bnConverter.FromBNOrderType(o.Type),
			Side:          b.bnConverter.FromBNOrderSide(o.Side),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(o.TimeInForce),
			CreateTime:    o.Time,
			UpdateTime:    o.UpdateTime,

			TriggerPrice:         o.StopPrice,
			TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(o.Type),
			TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(o.Side, o.Type),
		})
	}
	return orders
}
func (b *BinanceTradeEngine) handlePortfolioMarginMarginOrderQuery(req *QueryOrderParam, res *mybinanceapi.PortfolioMarginMarginOrderGetRes) *Order {
	avgPrice := decimal.Zero
	if res.ExecutedQty != "" && res.CummulativeQuoteQty != "" {
		executedQty, _ := decimal.NewFromString(res.ExecutedQty)
		cumQuoteQty, _ := decimal.NewFromString(res.CummulativeQuoteQty)
		if !executedQty.IsZero() {
			avgPrice = cumQuoteQty.Div(executedQty)
		}
	}
	return &Order{
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
}
func (b *BinanceTradeEngine) handlePortfolioMarginMarginOrdersQuery(req *QueryOrderParam, res *mybinanceapi.PortfolioMarginMarginAllOrdersRes) []*Order {
	var orders []*Order
	for _, o := range *res {
		avgPrice := decimal.Zero
		if o.ExecutedQty != "" && o.CummulativeQuoteQty != "" {
			executedQty, _ := decimal.NewFromString(o.ExecutedQty)
			cumQuoteQty, _ := decimal.NewFromString(o.CummulativeQuoteQty)
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
			OrderId:       strconv.FormatInt(o.OrderId, 10),
			ClientOrderId: o.ClientOrderId,
			Price:         o.Price,
			Quantity:      o.OrigQty,
			ExecutedQty:   o.ExecutedQty,
			CumQuoteQty:   o.CummulativeQuoteQty,
			AvgPrice:      avgPrice.String(),
			Status:        b.bnConverter.FromBNOrderStatus(o.Status, o.Type),
			Type:          b.bnConverter.FromBNOrderType(o.Type),
			Side:          b.bnConverter.FromBNOrderSide(o.Side),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(o.TimeInForce),
			CreateTime:    o.Time,
			UpdateTime:    o.UpdateTime,

			TriggerPrice:         o.StopPrice,
			TriggerType:          b.bnConverter.FromBNOrderTypeForTriggerType(o.Type),
			TriggerConditionType: b.bnConverter.FromBNOrderSideForTriggerConditionType(o.Side, o.Type),
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
			Symbol:      req.Symbol,
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
