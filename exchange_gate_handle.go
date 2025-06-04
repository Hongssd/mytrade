package mytrade

import (
	"strconv"
	"strings"
	"time"

	"github.com/Hongssd/mygateapi"
	"github.com/shopspring/decimal"
)

var gateTimeMul = decimal.NewFromInt(1000)

// 现货订单查询
func (g *GateTradeEngine) handleOrdersFromSpotOpenOrders(req *QueryOrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestSpotOpenOrdersRes]) []*Order {
	var orders []*Order

	for _, symbol := range res.Data {
		for _, order := range symbol.Orders {
			accountType, isMargin, isIsolated := g.gateConverter.FromOrderSpotAccountType(GateAccountType(order.Account))
			amt, _ := decimal.NewFromString(order.Amount)
			orders = append(orders, &Order{
				Exchange:      g.ExchangeType().String(),
				AccountType:   accountType.String(),
				Symbol:        order.CurrencyPair,
				IsMargin:      isMargin,
				IsIsolated:    isIsolated,
				OrderId:       order.ID,
				ClientOrderId: order.Text,
				Price:         order.Price,
				Quantity:      amt.Abs().String(),
				ExecutedQty:   order.FilledAmount,
				CumQuoteQty:   order.FilledTotal,
				AvgPrice:      order.AvgDealPrice,
				Status:        g.gateConverter.FromGateSpotOrderStatus(order.Status),
				Type:          g.gateConverter.FromGateOrderType(order.Type),
				Side:          g.gateConverter.FromGateOrderSide(order.Side),
				TimeInForce:   g.gateConverter.FromGateTimeInForce(order.TimeInForce),
				FeeAmount:     order.Fee,

				FeeCcy:     order.FeeCurrency,
				CreateTime: order.CreateTimeMs,
				UpdateTime: order.UpdateTimeMs,
			})
		}
	}
	return orders
}
func (g *GateTradeEngine) handleOrdersFromSpotPriceOpenOrders(req *QueryOrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestSpotPriceOrdersGetRes]) []*Order {
	var orders []*Order
	for _, order := range res.Data {
		accountType, isMargin, isIsolated := g.gateConverter.FromOrderSpotPriceAccountType(order.Put.Account)
		updateTime := decimal.NewFromInt(order.Ftime)
		if updateTime.IsZero() {
			updateTime = decimal.NewFromInt(order.Ctime)
		}
		amt, _ := decimal.NewFromString(order.Put.Amount)

		triggerType := g.gateConverter.FromGateSpotPriceOrderTriggerRule(order.Trigger.Rule, OrderSide(g.gateConverter.FromGateOrderSide(order.Put.Side)))
		triggerConditionType := g.gateConverter.FromGateTriggerCondition(OrderSide(g.gateConverter.FromGateOrderSide(order.Put.Side)), triggerType)

		orders = append(orders, &Order{
			Exchange:             g.ExchangeType().String(),
			AccountType:          accountType.String(),
			Symbol:               order.Market,
			IsMargin:             isMargin,
			IsIsolated:           isIsolated,
			OrderId:              strconv.FormatInt(order.ID, 10),
			ClientOrderId:        order.Put.Text,
			Price:                order.Put.Price,
			Quantity:             amt.Abs().String(),
			ExecutedQty:          "0",
			CumQuoteQty:          "0",
			AvgPrice:             "0",
			Status:               g.gateConverter.FromGateSpotPriceOrderStatus(order.Status),
			Type:                 OrderType(g.gateConverter.FromGateOrderType(order.Put.Type)),
			Side:                 OrderSide(g.gateConverter.FromGateOrderSide(order.Put.Side)),
			PositionSide:         "",
			TimeInForce:          g.gateConverter.FromGateTimeInForce(order.Put.TimeInForce),
			FeeAmount:            "",
			FeeCcy:               "",
			ReduceOnly:           false,
			CreateTime:           decimal.NewFromInt(order.Ctime).Mul(gateTimeMul).IntPart(),
			UpdateTime:           updateTime.Mul(gateTimeMul).IntPart(),
			RealizedPnl:          "",
			AttachTpTriggerPrice: "",
			AttachTpOrdPrice:     "",
			AttachSlTriggerPrice: "",
			IsAlgo:               true,
			TriggerPrice:         order.Trigger.Price,
			TriggerType:          triggerType,
			TriggerConditionType: triggerConditionType,
			Expiration:           order.Trigger.Expiration,
			OcoTpTriggerPrice:    "",
			OcoTpOrdType:         "",
			OcoTpOrdPrice:        "",
			OcoSlTriggerPrice:    "",
		})
	}
	return orders
}
func (g *GateTradeEngine) handleOrderFromSpotOrderQuery(req *QueryOrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestSpotOrdersOrderIdGetRes]) *Order {
	accountType, isMargin, isIsolated := g.gateConverter.FromOrderSpotAccountType(GateAccountType(res.Data.Account))
	amt, _ := decimal.NewFromString(res.Data.Amount)
	fillAmt, _ := decimal.NewFromString(res.Data.FilledAmount)
	fillTotal, _ := decimal.NewFromString(res.Data.FilledTotal)

	return &Order{
		Exchange:      g.ExchangeType().String(),
		AccountType:   accountType.String(),
		Symbol:        res.Data.CurrencyPair,
		IsMargin:      isMargin,
		IsIsolated:    isIsolated,
		OrderId:       req.OrderId,
		ClientOrderId: req.ClientOrderId,
		Price:         res.Data.Price,
		Quantity:      amt.Abs().String(),
		ExecutedQty:   fillAmt.Abs().String(),
		CumQuoteQty:   fillTotal.Abs().String(),
		AvgPrice:      res.Data.AvgDealPrice,
		Status:        g.gateConverter.FromGateSpotOrderStatus(res.Data.Status),
		Type:          g.gateConverter.FromGateOrderType(res.Data.Type),
		Side:          g.gateConverter.FromGateOrderSide(res.Data.Side),
		PositionSide:  "",

		TimeInForce:           g.gateConverter.FromGateTimeInForce(res.Data.TimeInForce),
		FeeAmount:             res.Data.Fee,
		FeeCcy:                res.Data.FeeCurrency,
		ReduceOnly:            false,
		CreateTime:            res.Data.CreateTimeMs,
		UpdateTime:            res.Data.UpdateTimeMs,
		RealizedPnl:           "",
		AttachTpTriggerPrice:  "",
		AttachTpOrdPrice:      "",
		AttachSlTriggerPrice:  "",
		AttachSlOrdPrice:      "",
		IsAlgo:                req.IsAlgo,
		OrderAlgoType:         "",
		TriggerPrice:          "",
		TriggerType:           "",
		TriggerConditionType:  "",
		OcoTpTriggerPrice:     "",
		OcoTpOrdType:          "",
		OcoTpOrdPrice:         "",
		OcoSlTriggerPrice:     "",
		OcoSlOrdType:          "",
		OcoSlOrdPrice:         "",
		MarginBuyBorrowAmount: "",
		MarginBuyBorrowAsset:  "",
		ErrorCode:             "",
		ErrorMsg:              "",
	}
}
func (g *GateTradeEngine) handleOrderFromSpotPriceOrderQuery(req *QueryOrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestSpotPriceOrdersOrderIdGetRes]) *Order {
	accountType, isMargin, isIsolated := g.gateConverter.FromOrderSpotPriceAccountType(res.Data.Put.Account)

	amt, _ := decimal.NewFromString(res.Data.Put.Amount)

	updateTime := decimal.NewFromInt(res.Data.Ftime).Mul(gateTimeMul).IntPart()
	if updateTime == 0 {
		updateTime = decimal.NewFromInt(res.Data.Ctime).Mul(gateTimeMul).IntPart()
	}

	triggerType := g.gateConverter.FromGateSpotPriceOrderTriggerRule(res.Data.Trigger.Rule, OrderSide(g.gateConverter.FromGateOrderSide(res.Data.Put.Side)))
	triggerConditionType := g.gateConverter.FromGateTriggerCondition(OrderSide(g.gateConverter.FromGateOrderSide(res.Data.Put.Side)), triggerType)

	return &Order{
		Exchange:             g.ExchangeType().String(),
		AccountType:          accountType.String(),
		Symbol:               res.Data.Market,
		IsMargin:             isMargin,
		IsIsolated:           isIsolated,
		OrderId:              strconv.FormatInt(res.Data.ID, 10),
		ClientOrderId:        res.Data.Put.Text,
		Price:                res.Data.Put.Price,
		Quantity:             amt.Abs().String(),
		ExecutedQty:          "0",
		CumQuoteQty:          "0",
		AvgPrice:             "0",
		Status:               g.gateConverter.FromGateSpotPriceOrderStatus(res.Data.Status),
		Type:                 OrderType(g.gateConverter.FromGateOrderType(res.Data.Put.Type)),
		Side:                 OrderSide(g.gateConverter.FromGateOrderSide(res.Data.Put.Side)),
		PositionSide:         "",
		TimeInForce:          g.gateConverter.FromGateTimeInForce(res.Data.Put.TimeInForce),
		FeeAmount:            "",
		FeeCcy:               "",
		ReduceOnly:           false,
		CreateTime:           decimal.NewFromInt(res.Data.Ctime).Mul(gateTimeMul).IntPart(),
		UpdateTime:           updateTime,
		IsAlgo:               true,
		RealizedPnl:          "",
		AttachTpTriggerPrice: "",
		AttachTpOrdPrice:     "",
		AttachSlTriggerPrice: "",
		TriggerPrice:         res.Data.Trigger.Price,
		TriggerType:          triggerType,
		TriggerConditionType: triggerConditionType,
		Expiration:           res.Data.Trigger.Expiration,
		OcoTpTriggerPrice:    "",
		OcoTpOrdType:         "",
		OcoTpOrdPrice:        "",
		OcoSlTriggerPrice:    "",
	}
}
func (g *GateTradeEngine) handleOrdersFromSpotOrdersQuery(req *QueryOrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestSpotOrdersGetRes]) []*Order {
	var orders []*Order
	for _, order := range res.Data {
		accountType, isMargin, isIsolated := g.gateConverter.FromOrderSpotAccountType(GateAccountType(order.Account))
		amt, _ := decimal.NewFromString(order.Amount)
		fillAmt, _ := decimal.NewFromString(order.FilledAmount)
		fillTotal, _ := decimal.NewFromString(order.FilledTotal)
		orders = append(orders, &Order{
			Exchange:      g.ExchangeType().String(),
			AccountType:   accountType.String(),
			Symbol:        order.CurrencyPair,
			IsMargin:      isMargin,
			IsIsolated:    isIsolated,
			OrderId:       order.ID,
			ClientOrderId: order.Text,
			Price:         order.Price,
			Quantity:      amt.Abs().String(),
			ExecutedQty:   fillAmt.Abs().String(),
			CumQuoteQty:   fillTotal.Abs().String(),
			AvgPrice:      order.AvgDealPrice,
			Status:        g.gateConverter.FromGateSpotOrderStatus(order.Status),
			Type:          g.gateConverter.FromGateOrderType(order.Type),
			Side:          g.gateConverter.FromGateOrderSide(order.Side),
			PositionSide:  "",

			TimeInForce:          g.gateConverter.FromGateTimeInForce(order.TimeInForce),
			FeeAmount:            order.Fee,
			FeeCcy:               order.FeeCurrency,
			ReduceOnly:           false,
			CreateTime:           int64(order.CreateTimeMs),
			UpdateTime:           int64(order.UpdateTimeMs),
			RealizedPnl:          "",
			AttachTpTriggerPrice: "",
			AttachTpOrdPrice:     "",
			AttachSlTriggerPrice: "",

			AttachSlOrdPrice:      "",
			IsAlgo:                req.IsAlgo,
			OrderAlgoType:         "",
			TriggerPrice:          "",
			TriggerType:           "",
			TriggerConditionType:  "",
			OcoTpTriggerPrice:     "",
			OcoTpOrdType:          "",
			OcoTpOrdPrice:         "",
			OcoSlTriggerPrice:     "",
			OcoSlOrdType:          "",
			OcoSlOrdPrice:         "",
			MarginBuyBorrowAmount: "",
			MarginBuyBorrowAsset:  "",
			ErrorCode:             "",
			ErrorMsg:              "",
		})
	}

	return orders
}
func (g *GateTradeEngine) handleOrdersFromSpotPriceOrdersQuery(req *QueryOrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestSpotPriceOrdersGetRes]) []*Order {
	var orders []*Order
	for _, order := range res.Data {
		accountType, isMargin, isIsolated := g.gateConverter.FromGateSpotPriceOrderAccount(order.Put.Account)

		updateTime := decimal.NewFromInt(order.Ftime).Mul(gateTimeMul).IntPart()
		if updateTime == 0 {
			updateTime = decimal.NewFromInt(order.Ctime).Mul(gateTimeMul).IntPart()
		}
		amt, _ := decimal.NewFromString(order.Put.Amount)

		triggerType := g.gateConverter.FromGateSpotPriceOrderTriggerRule(order.Trigger.Rule, OrderSide(g.gateConverter.FromGateOrderSide(order.Put.Side)))
		triggerConditionType := g.gateConverter.FromGateTriggerCondition(OrderSide(g.gateConverter.FromGateOrderSide(order.Put.Side)), triggerType)

		orders = append(orders, &Order{
			Exchange:              g.ExchangeType().String(),
			AccountType:           accountType.String(),
			Symbol:                order.Market,
			IsMargin:              isMargin,
			IsIsolated:            isIsolated,
			OrderId:               strconv.FormatInt(order.ID, 10),
			ClientOrderId:         order.Put.Text,
			Price:                 order.Put.Price,
			Quantity:              amt.Abs().String(),
			ExecutedQty:           "0",
			CumQuoteQty:           "0",
			AvgPrice:              "0",
			Status:                g.gateConverter.FromGateSpotPriceOrderStatus(order.Status),
			Type:                  OrderType(g.gateConverter.FromGateOrderType(order.Put.Type)),
			Side:                  OrderSide(g.gateConverter.FromGateOrderSide(order.Put.Side)),
			PositionSide:          "",
			TimeInForce:           g.gateConverter.FromGateTimeInForce(order.Put.TimeInForce),
			FeeAmount:             "",
			FeeCcy:                "",
			ReduceOnly:            false,
			CreateTime:            decimal.NewFromInt(order.Ctime).Mul(gateTimeMul).IntPart(),
			UpdateTime:            updateTime,
			RealizedPnl:           "",
			AttachTpTriggerPrice:  "",
			AttachTpOrdPrice:      "",
			AttachSlTriggerPrice:  "",
			IsAlgo:                true,
			TriggerPrice:          order.Trigger.Price,
			TriggerType:           triggerType,
			TriggerConditionType:  triggerConditionType,
			Expiration:            order.Trigger.Expiration,
			OcoTpTriggerPrice:     "",
			OcoTpOrdType:          "",
			OcoTpOrdPrice:         "",
			OcoSlTriggerPrice:     "",
			OcoSlOrdType:          "",
			OcoSlOrdPrice:         "",
			MarginBuyBorrowAmount: "",
			MarginBuyBorrowAsset:  "",
			ErrorCode:             "",
			ErrorMsg:              "",
		})
	}
	return orders
}
func (g *GateTradeEngine) handleTradesFromSpotTradesQuery(req *QueryTradeParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestSpotMyTradesRes]) []*Trade {
	var trades []*Trade
	for _, trade := range res.Data {
		price, _ := decimal.NewFromString(trade.Price)
		amt, _ := decimal.NewFromString(trade.Amount)
		quoteQty := price.Mul(amt)

		isMaker := false
		if trade.Role == "maker" {
			isMaker = true
		}

		trades = append(trades, &Trade{
			Exchange:      g.ExchangeType().String(),
			AccountType:   req.AccountType,
			Symbol:        req.Symbol,
			TradeId:       trade.ID,
			OrderId:       trade.OrderID,
			ClientOrderId: trade.Text,
			Price:         trade.Price,
			Quantity:      trade.Amount,
			QuoteQty:      quoteQty.String(),
			Side:          g.gateConverter.FromGateOrderSide(trade.Side),
			PositionSide:  "",
			FeeAmount:     trade.Fee,
			FeeCcy:        trade.FeeCurrency,
			RealizedPnl:   "",
			IsMaker:       isMaker,
			Timestamp:     stringToInt64(trade.CreateTimeMs),
		})
	}

	return trades
}

// 现货订单操作
func (g *GateTradeEngine) handleOrderFromSpotOrderCreate(req *OrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestSpotOrdersPostRes]) *Order {
	accountType, isMargin, isIsolated := g.gateConverter.FromOrderSpotAccountType(GateAccountType(res.Data.Account))
	amt, _ := decimal.NewFromString(res.Data.Amount)
	fillAmt, _ := decimal.NewFromString(res.Data.FilledAmount)
	fillTotal, _ := decimal.NewFromString(res.Data.FilledTotal)
	return &Order{
		Exchange:      g.ExchangeType().String(),
		AccountType:   accountType.String(),
		Symbol:        res.Data.CurrencyPair,
		IsMargin:      isMargin,
		IsIsolated:    isIsolated,
		OrderId:       res.Data.ID,
		ClientOrderId: res.Data.Text,
		Price:         res.Data.Price,
		Quantity:      amt.Abs().String(),
		ExecutedQty:   fillAmt.Abs().String(),
		CumQuoteQty:   fillTotal.Abs().String(),
		AvgPrice:      res.Data.AvgDealPrice,
		Status:        g.gateConverter.FromGateSpotOrderStatus(res.Data.Status),
		Type:          g.gateConverter.FromGateOrderType(res.Data.Type),
		Side:          g.gateConverter.FromGateOrderSide(res.Data.Side),
		TimeInForce:   g.gateConverter.FromGateTimeInForce(res.Data.TimeInForce),
		FeeAmount:     res.Data.Fee,
		FeeCcy:        res.Data.FeeCurrency,
		CreateTime:    int64(res.Data.CreateTimeMs),
		UpdateTime:    int64(res.Data.UpdateTimeMs),
	}
}
func (g *GateTradeEngine) handleOrderFromSpotPriceOrderCreate(req *OrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestSpotPriceOrdersPostRes]) *Order {
	amt := req.Quantity
	return &Order{
		Exchange:             g.ExchangeType().String(),
		AccountType:          req.AccountType,
		Symbol:               req.Symbol,
		IsMargin:             req.IsMargin,
		IsIsolated:           req.IsIsolated,
		OrderId:              strconv.FormatInt(res.Data.ID, 10),
		ClientOrderId:        req.ClientOrderId,
		Price:                req.Price.String(),
		Quantity:             amt.Abs().String(),
		ExecutedQty:          "0",
		CumQuoteQty:          "0",
		AvgPrice:             "0",
		Status:               ORDER_STATUS_UN_TRIGGERED,
		Type:                 req.OrderType,
		Side:                 req.OrderSide,
		TimeInForce:          req.TimeInForce,
		FeeAmount:            "",
		FeeCcy:               "",
		CreateTime:           stringToInt64(res.GateTimeRes.OutTime) / 1000,
		UpdateTime:           stringToInt64(res.GateTimeRes.OutTime) / 1000,
		RealizedPnl:          "",
		AttachTpTriggerPrice: "",
		AttachTpOrdPrice:     "",
		AttachSlTriggerPrice: "",
		AttachSlOrdPrice:     "",
		IsAlgo:               req.IsAlgo,
		OrderAlgoType:        "", // 仅支持单向止盈止损
		TriggerPrice:         req.TriggerPrice.String(),
		TriggerType:          req.TriggerType,
		TriggerConditionType: g.gateConverter.FromGateTriggerCondition(req.OrderSide, req.TriggerType),
		OcoTpTriggerPrice:    "",
		OcoTpOrdType:         "",
		OcoTpOrdPrice:        "",
		OcoSlTriggerPrice:    "",
		OcoSlOrdType:         "",
		OcoSlOrdPrice:        "",
	}
}
func (g *GateTradeEngine) handleOrderFromSpotOrderAmend(req *OrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestSpotOrdersOrderIdPatchRes]) *Order {
	accountType, isMargin, isIsolated := g.gateConverter.FromOrderSpotAccountType(GateAccountType(res.Data.Account))
	amt, _ := decimal.NewFromString(res.Data.Amount)
	fillAmt, _ := decimal.NewFromString(res.Data.FilledAmount)
	fillTotal, _ := decimal.NewFromString(res.Data.FilledTotal)
	return &Order{
		Exchange:      g.ExchangeType().String(),
		AccountType:   accountType.String(),
		Symbol:        res.Data.CurrencyPair,
		IsMargin:      isMargin,
		IsIsolated:    isIsolated,
		OrderId:       res.Data.ID,
		ClientOrderId: res.Data.Text,
		Price:         res.Data.Price,
		Quantity:      amt.Abs().String(),
		ExecutedQty:   fillAmt.Abs().String(),
		CumQuoteQty:   fillTotal.Abs().String(),
		AvgPrice:      res.Data.AvgDealPrice,
		Status:        g.gateConverter.FromGateSpotOrderStatus(res.Data.Status),
		Type:          g.gateConverter.FromGateOrderType(res.Data.Type),
		Side:          g.gateConverter.FromGateOrderSide(res.Data.Side),
		TimeInForce:   g.gateConverter.FromGateTimeInForce(res.Data.TimeInForce),

		FeeAmount:  res.Data.Fee,
		FeeCcy:     res.Data.FeeCurrency,
		CreateTime: res.Data.CreateTimeMs,
		UpdateTime: res.Data.UpdateTimeMs,
	}
}
func (g *GateTradeEngine) handleOrderFromSpotOrderCancel(req *OrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestSpotOrdersOrderIdDeleteRes]) *Order {
	accountType, isMargin, isIsolated := g.gateConverter.FromOrderSpotAccountType(GateAccountType(res.Data.Account))
	amt, _ := decimal.NewFromString(res.Data.Amount)
	fillAmt, _ := decimal.NewFromString(res.Data.FilledAmount)
	fillTotal, _ := decimal.NewFromString(res.Data.FilledTotal)
	return &Order{
		Exchange:      g.ExchangeType().String(),
		AccountType:   accountType.String(),
		Symbol:        res.Data.CurrencyPair,
		IsMargin:      isMargin,
		IsIsolated:    isIsolated,
		OrderId:       res.Data.ID,
		ClientOrderId: res.Data.Text,
		Price:         res.Data.Price,
		Quantity:      amt.Abs().String(),
		ExecutedQty:   fillAmt.Abs().String(),
		CumQuoteQty:   fillTotal.Abs().String(),
		AvgPrice:      res.Data.AvgDealPrice,
		Status:        g.gateConverter.FromGateSpotOrderStatus(res.Data.Status),
		Type:          g.gateConverter.FromGateOrderType(res.Data.Type),
		Side:          g.gateConverter.FromGateOrderSide(res.Data.Side),
		PositionSide:  "",

		TimeInForce:           g.gateConverter.FromGateTimeInForce(res.Data.TimeInForce),
		FeeAmount:             res.Data.Fee,
		FeeCcy:                res.Data.FeeCurrency,
		ReduceOnly:            false,
		CreateTime:            res.Data.CreateTimeMs,
		UpdateTime:            res.Data.UpdateTimeMs,
		RealizedPnl:           "",
		AttachTpTriggerPrice:  "",
		AttachTpOrdPrice:      "",
		AttachSlTriggerPrice:  "",
		AttachSlOrdPrice:      "",
		IsAlgo:                req.IsAlgo,
		OrderAlgoType:         "",
		TriggerPrice:          "",
		TriggerType:           "",
		TriggerConditionType:  "",
		OcoTpTriggerPrice:     "",
		OcoTpOrdType:          "",
		OcoTpOrdPrice:         "",
		OcoSlTriggerPrice:     "",
		OcoSlOrdType:          "",
		OcoSlOrdPrice:         "",
		MarginBuyBorrowAmount: "",
		MarginBuyBorrowAsset:  "",
		ErrorCode:             "",
		ErrorMsg:              "",
	}
}
func (g *GateTradeEngine) handleOrderFromSpotPriceOrderCancel(req *OrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestSpotPriceOrdersOrderIdDeleteRes]) *Order {
	accountType, isMargin, isIsolated := g.gateConverter.FromOrderSpotPriceAccountType(res.Data.Put.Account)
	amt, _ := decimal.NewFromString(res.Data.Put.Amount)

	triggerType := g.gateConverter.FromGateSpotPriceOrderTriggerRule(res.Data.Trigger.Rule, OrderSide(g.gateConverter.FromGateOrderSide(res.Data.Put.Side)))
	triggerConditionType := g.gateConverter.FromGateTriggerCondition(OrderSide(g.gateConverter.FromGateOrderSide(res.Data.Put.Side)), triggerType)

	return &Order{
		Exchange:             g.ExchangeType().String(),
		AccountType:          accountType.String(),
		Symbol:               res.Data.Market,
		IsMargin:             isMargin,
		IsIsolated:           isIsolated,
		OrderId:              strconv.FormatInt(res.Data.ID, 10),
		ClientOrderId:        res.Data.Put.Text,
		Price:                res.Data.Put.Price,
		Quantity:             amt.Abs().String(),
		ExecutedQty:          "0",
		CumQuoteQty:          "0",
		AvgPrice:             "0",
		Status:               g.gateConverter.FromGateSpotPriceOrderStatus(res.Data.Status),
		Type:                 OrderType(g.gateConverter.FromGateOrderType(res.Data.Put.Type)),
		Side:                 OrderSide(g.gateConverter.FromGateOrderSide(res.Data.Put.Side)),
		PositionSide:         "",
		TimeInForce:          g.gateConverter.FromGateTimeInForce(res.Data.Put.TimeInForce),
		FeeAmount:            "",
		FeeCcy:               "",
		ReduceOnly:           false,
		CreateTime:           decimal.NewFromInt((res.Data.Ctime)).Mul(gateTimeMul).IntPart(),
		UpdateTime:           decimal.NewFromInt((res.Data.Ftime)).Mul(gateTimeMul).IntPart(),
		RealizedPnl:          "",
		AttachTpTriggerPrice: "",
		AttachTpOrdPrice:     "",
		AttachSlTriggerPrice: "",
		IsAlgo:               true,
		TriggerPrice:         res.Data.Trigger.Price,
		TriggerType:          triggerType,
		TriggerConditionType: triggerConditionType,
		Expiration:           res.Data.Trigger.Expiration,
		OcoTpTriggerPrice:    "",
		OcoTpOrdType:         "",
		OcoTpOrdPrice:        "",
		OcoSlTriggerPrice:    "",
	}
}
func (g *GateTradeEngine) getPositionSide(accountType, contract string, size int64, reduceOnly bool) (PositionSide, error) {
	apiParam := ExchangeApiParam{
		Exchange:  GATE_NAME.String(),
		ApiKey:    g.apiKey,
		ApiSecret: g.secretKey,
	}
	positions, err := InnerExchangeManager.GetPositions(apiParam, accountType, contract)
	if err != nil {
		log.Error(err)
		return POSITION_SIDE_UNKNOWN, err
	}
	if len(positions) == 1 {
		return POSITION_SIDE_BOTH, nil
	} else {
		if reduceOnly {
			//只减仓
			if size > 0 {
				//只减仓买入 平空买入
				return POSITION_SIDE_LONG, nil
			} else {
				//只减仓卖出 平多卖出
				return POSITION_SIDE_SHORT, nil
			}
		} else {
			if size > 0 {
				//开仓买入 开多
				return POSITION_SIDE_LONG, nil
			} else {
				//开仓卖出 开空
				return POSITION_SIDE_SHORT, nil
			}
		}
	}
}

// 永续合约订单查询
func (g *GateTradeEngine) handleOrdersFromFuturesOpenOrders(req *QueryOrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestFuturesSettleOrdersGetRes]) []*Order {
	var orders []*Order
	for _, order := range res.Data {
		amt := decimal.NewFromInt(order.Size)
		executedQty := decimal.NewFromInt(order.Size).Sub(decimal.NewFromInt(order.Left)).Abs()
		fillPrice, _ := decimal.NewFromString(order.FillPrice)
		symbolInfo, err := InnerExchangeManager.GetSymbolInfo(GATE_NAME.String(), req.AccountType, order.Contract)
		if err != nil {
			log.Error(err)
			continue
		}
		cumQuoteQty := executedQty.Mul(symbolInfo.ContractSize()).Mul(fillPrice)

		var updateTime int64
		if decimal.NewFromFloat(order.FinishTime).IsZero() {
			updateTime = decimal.NewFromFloat(order.CreateTime).Mul(gateTimeMul).IntPart()
		} else {
			updateTime = decimal.NewFromFloat(order.FinishTime).Mul(gateTimeMul).IntPart()
		}

		var orderType OrderType
		if order.Price == "0" && order.Tif == GATE_TIME_IN_FORCE_IOC {
			orderType = ORDER_TYPE_MARKET
		} else {
			orderType = ORDER_TYPE_LIMIT
		}

		var orderSide OrderSide
		if order.Size > 0 {
			orderSide = ORDER_SIDE_BUY
		} else {
			orderSide = ORDER_SIDE_SELL
		}

		positionSide, err := g.getPositionSide(req.AccountType, order.Contract, order.Size, order.ReduceOnly)
		if err != nil {
			log.Error(err)
		}
		orders = append(orders, &Order{
			Exchange:             g.ExchangeType().String(),
			AccountType:          req.AccountType,
			Symbol:               order.Contract,
			IsMargin:             req.IsMargin,
			IsIsolated:           req.IsIsolated,
			OrderId:              strconv.FormatInt(order.Id, 10),
			ClientOrderId:        order.Text,
			Price:                order.Price,
			Quantity:             amt.Abs().String(),
			ExecutedQty:          executedQty.String(),
			CumQuoteQty:          cumQuoteQty.String(),
			AvgPrice:             order.FillPrice,
			Status:               g.gateConverter.FromGateContractOrderStatus(order.Status, order.FinishAs),
			Type:                 orderType,
			Side:                 orderSide,
			PositionSide:         positionSide,
			TimeInForce:          g.gateConverter.FromGateTimeInForce(order.Tif),
			FeeAmount:            "",
			FeeCcy:               "",
			ReduceOnly:           order.ReduceOnly,
			CreateTime:           decimal.NewFromFloat(order.CreateTime).Mul(gateTimeMul).IntPart(),
			UpdateTime:           updateTime,
			RealizedPnl:          "",
			AttachTpTriggerPrice: "",
			AttachTpOrdPrice:     "",
			AttachSlTriggerPrice: "",

			AttachSlOrdPrice:      "",
			IsAlgo:                req.IsAlgo,
			OrderAlgoType:         "",
			TriggerPrice:          "",
			TriggerType:           "",
			TriggerConditionType:  "",
			OcoTpTriggerPrice:     "",
			OcoTpOrdType:          "",
			OcoTpOrdPrice:         "",
			OcoSlTriggerPrice:     "",
			OcoSlOrdType:          "",
			OcoSlOrdPrice:         "",
			MarginBuyBorrowAmount: "",
			MarginBuyBorrowAsset:  "",
			ErrorCode:             "",
			ErrorMsg:              "",
		})
	}

	return orders
}
func (g *GateTradeEngine) handleOrdersFromFuturesPriceOpenOrders(req *QueryOrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestFuturesSettlePriceOrdersGetRes]) []*Order {
	var orders []*Order
	for _, order := range res.Data {

		amt := decimal.NewFromInt(order.Initial.Size)
		var updateTime int64
		if decimal.NewFromFloat(order.FinishTime).IsZero() {
			updateTime = decimal.NewFromFloat(order.CreateTime).Mul(gateTimeMul).IntPart()
		} else {
			updateTime = decimal.NewFromFloat(order.FinishTime).Mul(gateTimeMul).IntPart()
		}

		var orderType OrderType
		if order.Initial.Price == "0" && order.Initial.Tif == GATE_TIME_IN_FORCE_IOC {
			orderType = ORDER_TYPE_MARKET
		} else {
			orderType = ORDER_TYPE_LIMIT
		}

		var orderSide OrderSide
		if order.Initial.Size > 0 {
			orderSide = ORDER_SIDE_BUY
		} else {
			orderSide = ORDER_SIDE_SELL
		}

		triggerType := g.gateConverter.FromGateFuturesPriceOrderTriggerRule(order.Trigger.Rule, orderSide)
		triggerConditionType := g.gateConverter.FromGateTriggerCondition(orderSide, triggerType)
		positionSide, err := g.getPositionSide(req.AccountType, order.Initial.Contract, order.Initial.Size, order.Initial.IsReduceOnly)
		if err != nil {
			log.Error(err)
		}
		orders = append(orders, &Order{
			Exchange:             g.ExchangeType().String(),
			AccountType:          req.AccountType,
			Symbol:               order.Initial.Contract,
			IsMargin:             req.IsMargin,
			IsIsolated:           req.IsIsolated,
			OrderId:              strconv.FormatInt(order.Id, 10),
			ClientOrderId:        order.Initial.Text,
			Price:                order.Initial.Price,
			Quantity:             amt.Abs().String(),
			ExecutedQty:          "0",
			CumQuoteQty:          "0",
			AvgPrice:             "0",
			Status:               g.gateConverter.FromGateContractPriceOrderStatus(order.Status, order.FinishAs),
			Type:                 orderType,
			Side:                 orderSide,
			PositionSide:         positionSide,
			TimeInForce:          g.gateConverter.FromGateTimeInForce(order.Initial.Tif),
			FeeAmount:            "",
			FeeCcy:               "",
			ReduceOnly:           order.Initial.IsReduceOnly,
			CreateTime:           decimal.NewFromFloat(order.CreateTime).Mul(gateTimeMul).IntPart(),
			UpdateTime:           updateTime,
			IsAlgo:               true,
			RealizedPnl:          "",
			AttachTpTriggerPrice: "",
			AttachTpOrdPrice:     "",
			AttachSlTriggerPrice: "",
			TriggerPrice:         order.Trigger.Price,
			TriggerType:          triggerType,
			TriggerConditionType: triggerConditionType,
			OcoTpTriggerPrice:    "",
			OcoTpOrdType:         "",
			OcoTpOrdPrice:        "",
			OcoSlTriggerPrice:    "",
			OcoSlOrdType:         "",
			OcoSlOrdPrice:        "",
		})
	}

	return orders
}
func (g *GateTradeEngine) handleOrderFromFuturesOrderQuery(req *QueryOrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestFuturesSettleOrdersOrderIdGetRes]) *Order {
	amt := decimal.NewFromInt(res.Data.Size)
	executedQty := decimal.NewFromInt(res.Data.Size).Sub(decimal.NewFromInt(res.Data.Left)).Abs()
	fillPrice, _ := decimal.NewFromString(res.Data.FillPrice)
	symbolInfo, err := InnerExchangeManager.GetSymbolInfo(GATE_NAME.String(), req.AccountType, res.Data.Contract)
	if err != nil {
		log.Error(err)
		return nil
	}
	cumQuoteQty := executedQty.Mul(symbolInfo.ContractSize()).Mul(fillPrice)
	var updateTime int64

	if decimal.NewFromFloat(res.Data.FinishTime).IsZero() {
		updateTime = decimal.NewFromFloat(res.Data.CreateTime).Mul(gateTimeMul).IntPart()
	} else {
		updateTime = decimal.NewFromFloat(res.Data.FinishTime).Mul(gateTimeMul).IntPart()
	}

	var orderType OrderType
	if res.Data.Price == "0" && res.Data.Tif == GATE_TIME_IN_FORCE_IOC {
		orderType = ORDER_TYPE_MARKET
	} else {
		orderType = ORDER_TYPE_LIMIT
	}

	var orderSide OrderSide
	if res.Data.Size > 0 {
		orderSide = ORDER_SIDE_BUY
	} else {
		orderSide = ORDER_SIDE_SELL
	}

	positionSide, err := g.getPositionSide(req.AccountType, res.Data.Contract, res.Data.Size, res.Data.ReduceOnly)
	if err != nil {
		log.Error(err)
	}

	return &Order{
		Exchange:      g.ExchangeType().String(),
		AccountType:   req.AccountType,
		Symbol:        res.Data.Contract,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
		OrderId:       strconv.FormatInt(res.Data.Id, 10),
		ClientOrderId: res.Data.Text,
		Price:         res.Data.Price,
		Quantity:      amt.Abs().String(),

		ExecutedQty:          executedQty.String(),
		CumQuoteQty:          cumQuoteQty.String(),
		AvgPrice:             res.Data.FillPrice,
		Status:               g.gateConverter.FromGateContractOrderStatus(res.Data.Status, res.Data.FinishAs),
		Type:                 orderType,
		Side:                 orderSide,
		PositionSide:         positionSide,
		TimeInForce:          g.gateConverter.FromGateTimeInForce(res.Data.Tif),
		FeeAmount:            "",
		FeeCcy:               "",
		ReduceOnly:           res.Data.ReduceOnly,
		CreateTime:           decimal.NewFromFloat(res.Data.CreateTime).Mul(gateTimeMul).IntPart(),
		UpdateTime:           updateTime,
		RealizedPnl:          "",
		AttachTpTriggerPrice: "",
		AttachTpOrdPrice:     "",
		AttachSlTriggerPrice: "",

		AttachSlOrdPrice:      "",
		IsAlgo:                req.IsAlgo,
		OrderAlgoType:         "",
		TriggerPrice:          "",
		TriggerType:           "",
		TriggerConditionType:  "",
		OcoTpTriggerPrice:     "",
		OcoTpOrdType:          "",
		OcoTpOrdPrice:         "",
		OcoSlTriggerPrice:     "",
		OcoSlOrdType:          "",
		OcoSlOrdPrice:         "",
		MarginBuyBorrowAmount: "",
		MarginBuyBorrowAsset:  "",
		ErrorCode:             "",
		ErrorMsg:              "",
	}
}
func (g *GateTradeEngine) handleOrderFromFuturesPriceOrderQuery(req *QueryOrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestFuturesSettlePriceOrdersOrderIdGetRes]) *Order {

	amt := decimal.NewFromInt(res.Data.Initial.Size)

	updateTime := decimal.NewFromFloat(res.Data.FinishTime).Mul(gateTimeMul).IntPart()
	if updateTime == 0 {
		updateTime = decimal.NewFromFloat(res.Data.CreateTime).Mul(gateTimeMul).IntPart()
	}

	var orderType OrderType
	if res.Data.Initial.Price == "0" && res.Data.Initial.Tif == GATE_TIME_IN_FORCE_IOC {
		orderType = ORDER_TYPE_MARKET
	} else {
		orderType = ORDER_TYPE_LIMIT
	}

	var orderSide OrderSide
	if res.Data.Initial.Size > 0 {
		orderSide = ORDER_SIDE_BUY
	} else {
		orderSide = ORDER_SIDE_SELL
	}

	triggerType := g.gateConverter.FromGateFuturesPriceOrderTriggerRule(res.Data.Trigger.Rule, orderSide)
	triggerConditionType := g.gateConverter.FromGateTriggerCondition(orderSide, triggerType)

	positionSide, err := g.getPositionSide(req.AccountType, res.Data.Initial.Contract, res.Data.Initial.Size, res.Data.Initial.IsReduceOnly)
	if err != nil {
		log.Error(err)
	}

	return &Order{
		Exchange:             g.ExchangeType().String(),
		AccountType:          req.AccountType,
		Symbol:               res.Data.Initial.Contract,
		IsMargin:             req.IsMargin,
		IsIsolated:           req.IsIsolated,
		OrderId:              strconv.FormatInt(res.Data.Id, 10),
		ClientOrderId:        res.Data.Initial.Text,
		Price:                res.Data.Initial.Price,
		Quantity:             amt.Abs().String(),
		ExecutedQty:          "0",
		CumQuoteQty:          "0",
		AvgPrice:             "0",
		Status:               g.gateConverter.FromGateContractPriceOrderStatus(res.Data.Status, res.Data.FinishAs),
		Type:                 orderType,
		Side:                 orderSide,
		PositionSide:         positionSide,
		TimeInForce:          g.gateConverter.FromGateTimeInForce(res.Data.Initial.Tif),
		FeeAmount:            "",
		FeeCcy:               "",
		ReduceOnly:           res.Data.Initial.IsReduceOnly,
		CreateTime:           decimal.NewFromFloat(res.Data.CreateTime).Mul(gateTimeMul).IntPart(),
		UpdateTime:           updateTime,
		IsAlgo:               true,
		RealizedPnl:          "",
		AttachTpTriggerPrice: "",
		AttachTpOrdPrice:     "",
		AttachSlTriggerPrice: "",
		TriggerPrice:         res.Data.Trigger.Price,
		TriggerType:          triggerType,
		TriggerConditionType: triggerConditionType,
		OcoTpTriggerPrice:    "",
		OcoTpOrdType:         "",
		OcoTpOrdPrice:        "",
		OcoSlTriggerPrice:    "",
		OcoSlOrdType:         "",
		OcoSlOrdPrice:        "",
	}
}
func (g *GateTradeEngine) handleOrdersFromFuturesOrdersQuery(req *QueryOrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestFuturesSettleOrdersGetRes]) []*Order {
	var orders []*Order
	for _, order := range res.Data {
		amt := decimal.NewFromInt(order.Size)
		executedQty := decimal.NewFromInt(order.Size).Sub(decimal.NewFromInt(order.Left)).Abs()
		fillPrice, _ := decimal.NewFromString(order.FillPrice)
		cumQuoteQty := decimal.Zero
		symbolInfo, err := InnerExchangeManager.GetSymbolInfo(GATE_NAME.String(), req.AccountType, order.Contract)
		if err != nil {
			log.Error(err)
			continue
		}

		// log.Info(order.Id, ":", order.IsClose)

		if symbolInfo != nil && symbolInfo.IsContract() && symbolInfo.IsContractAmt() {
			cumQuoteQty = executedQty.Mul(symbolInfo.ContractSize()).Mul(fillPrice)
		}

		var updateTime int64

		finishTime := decimal.NewFromFloat(order.FinishTime)
		if finishTime.IsZero() {
			updateTime = decimal.NewFromFloat(order.CreateTime).Mul(gateTimeMul).IntPart()
		} else {
			updateTime = finishTime.Mul(gateTimeMul).IntPart()
		}

		var orderType OrderType
		if order.Price == "0" && order.Tif == GATE_TIME_IN_FORCE_IOC {
			orderType = ORDER_TYPE_MARKET
		} else {
			orderType = ORDER_TYPE_LIMIT
		}

		var orderSide OrderSide
		if order.Size > 0 {
			orderSide = ORDER_SIDE_BUY
		} else {
			orderSide = ORDER_SIDE_SELL
		}

		positionSide, err := g.getPositionSide(req.AccountType, order.Contract, order.Size, order.ReduceOnly)
		if err != nil {
			log.Error(err)
		}

		orders = append(orders, &Order{
			Exchange:              g.ExchangeType().String(),
			AccountType:           req.AccountType,
			Symbol:                order.Contract,
			IsMargin:              req.IsMargin,
			IsIsolated:            req.IsIsolated,
			OrderId:               strconv.FormatInt(order.Id, 10),
			ClientOrderId:         order.Text,
			Price:                 order.Price,
			Quantity:              amt.Abs().String(),
			ExecutedQty:           executedQty.String(),
			CumQuoteQty:           cumQuoteQty.String(),
			AvgPrice:              order.FillPrice,
			Status:                g.gateConverter.FromGateContractOrderStatus(order.Status, order.FinishAs),
			Type:                  orderType,
			Side:                  orderSide,
			PositionSide:          positionSide,
			TimeInForce:           g.gateConverter.FromGateTimeInForce(order.Tif),
			FeeAmount:             "",
			FeeCcy:                "",
			ReduceOnly:            order.IsReduceOnly,
			CreateTime:            decimal.NewFromFloat(order.CreateTime).Mul(gateTimeMul).IntPart(),
			UpdateTime:            updateTime,
			RealizedPnl:           "",
			AttachTpTriggerPrice:  "",
			AttachTpOrdPrice:      "",
			AttachSlTriggerPrice:  "",
			AttachSlOrdPrice:      "",
			IsAlgo:                req.IsAlgo,
			OrderAlgoType:         "",
			TriggerPrice:          "",
			TriggerType:           "",
			TriggerConditionType:  "",
			OcoTpTriggerPrice:     "",
			OcoTpOrdType:          "",
			OcoTpOrdPrice:         "",
			OcoSlTriggerPrice:     "",
			OcoSlOrdType:          "",
			OcoSlOrdPrice:         "",
			MarginBuyBorrowAmount: "",
			MarginBuyBorrowAsset:  "",
			ErrorCode:             "",
			ErrorMsg:              "",
		})
	}

	return orders
}
func (g *GateTradeEngine) handleOrdersFromFuturesPriceOrdersQuery(req *QueryOrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestFuturesSettlePriceOrdersGetRes]) []*Order {
	var orders []*Order
	for _, order := range res.Data {

		amt := decimal.NewFromInt(order.Initial.Size)
		var updateTime int64
		if decimal.NewFromFloat(order.FinishTime).IsZero() {
			updateTime = decimal.NewFromFloat(order.CreateTime).Mul(gateTimeMul).IntPart()
		} else {
			updateTime = decimal.NewFromFloat(order.FinishTime).Mul(gateTimeMul).IntPart()
		}

		var orderType OrderType
		if order.Initial.Price == "0" && order.Initial.Tif == GATE_TIME_IN_FORCE_IOC {
			orderType = ORDER_TYPE_MARKET
		} else {
			orderType = ORDER_TYPE_LIMIT
		}

		var orderSide OrderSide
		if order.Initial.Size > 0 {
			orderSide = ORDER_SIDE_BUY
		} else {
			orderSide = ORDER_SIDE_SELL
		}

		triggerType := g.gateConverter.FromGateFuturesPriceOrderTriggerRule(order.Trigger.Rule, orderSide)
		triggerConditionType := g.gateConverter.FromGateTriggerCondition(orderSide, triggerType)

		positionSide, err := g.getPositionSide(req.AccountType, order.Initial.Contract, order.Initial.Size, order.Initial.IsReduceOnly)
		if err != nil {
			log.Error(err)
		}

		orders = append(orders, &Order{
			Exchange:             g.ExchangeType().String(),
			AccountType:          req.AccountType,
			Symbol:               order.Initial.Contract,
			IsMargin:             req.IsMargin,
			IsIsolated:           req.IsIsolated,
			OrderId:              strconv.FormatInt(order.Id, 10),
			ClientOrderId:        order.Initial.Text,
			Price:                order.Initial.Price,
			Quantity:             amt.Abs().String(),
			ExecutedQty:          "0",
			CumQuoteQty:          "0",
			AvgPrice:             "0",
			Status:               g.gateConverter.FromGateContractPriceOrderStatus(order.Status, order.FinishAs),
			Type:                 orderType,
			Side:                 orderSide,
			PositionSide:         positionSide,
			TimeInForce:          g.gateConverter.FromGateTimeInForce(order.Initial.Tif),
			FeeAmount:            "",
			FeeCcy:               "",
			ReduceOnly:           order.Initial.IsReduceOnly,
			CreateTime:           decimal.NewFromFloat(order.CreateTime).Mul(gateTimeMul).IntPart(),
			UpdateTime:           updateTime,
			IsAlgo:               true,
			RealizedPnl:          "",
			AttachTpTriggerPrice: "",
			AttachTpOrdPrice:     "",
			AttachSlTriggerPrice: "",
			TriggerPrice:         order.Trigger.Price,
			TriggerType:          triggerType,
			TriggerConditionType: triggerConditionType,
			OcoTpTriggerPrice:    "",
			OcoTpOrdType:         "",
			OcoTpOrdPrice:        "",
			OcoSlTriggerPrice:    "",
			OcoSlOrdType:         "",
			OcoSlOrdPrice:        "",
		})
	}
	return orders
}
func (g *GateTradeEngine) handleTradesFromFuturesTradesQuery(req *QueryTradeParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestFuturesSettleMyTradesRes]) []*Trade {
	var trades []*Trade
	var order *Order
	var positionSide PositionSide
	var err error

	order, err = g.QueryOrder(&QueryOrderParam{
		AccountType: req.AccountType,
		Symbol:      req.Symbol,
		OrderId:     req.OrderId,
	})
	if err != nil {
		log.Error(err)
		return nil
	}
	if order != nil {
		positionSide = order.PositionSide
	}

	for _, trade := range res.Data {
		price, _ := decimal.NewFromString(trade.Price)
		amt := decimal.NewFromInt(trade.Size)

		quoteQty := decimal.Zero
		symbolInfo, err := InnerExchangeManager.GetSymbolInfo(GATE_NAME.String(), req.AccountType, trade.Contract)
		if err != nil {
			log.Error(err)
			continue
		}
		if symbolInfo != nil && symbolInfo.IsContract() && symbolInfo.IsContractAmt() {
			quoteQty = amt.Abs().Mul(symbolInfo.ContractSize()).Mul(price)
		}

		isMaker := false
		if trade.Role == "maker" {
			isMaker = true
		}

		var side OrderSide
		if trade.Size > 0 {
			side = ORDER_SIDE_BUY
		} else {
			side = ORDER_SIDE_SELL
		}
		feeCcy := strings.Split(trade.Contract, "_")[1]

		trades = append(trades, &Trade{
			Exchange:      g.ExchangeType().String(),
			AccountType:   req.AccountType,
			Symbol:        req.Symbol,
			TradeId:       strconv.FormatInt(trade.Id, 10),
			OrderId:       trade.OrderId,
			ClientOrderId: trade.Text,
			Price:         trade.Price,
			Quantity:      amt.Abs().String(),
			QuoteQty:      quoteQty.String(),
			Side:          side,
			PositionSide:  positionSide,
			FeeAmount:     trade.Fee,
			FeeCcy:        feeCcy,
			RealizedPnl:   "",
			IsMaker:       isMaker,
			Timestamp:     decimal.NewFromFloat(trade.CreateTime).Mul(gateTimeMul).IntPart(),
		})
	}
	return trades
}

// 永续合约订单操作
func (g *GateTradeEngine) handleOrderFromFuturesOrderCreate(req *OrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestFuturesSettleOrdersPostRes]) *Order {
	amt := decimal.NewFromInt(res.Data.Size)
	executedQty := decimal.NewFromInt(res.Data.Size).Sub(decimal.NewFromInt(res.Data.Left)).Abs()
	fillPrice, _ := decimal.NewFromString(res.Data.FillPrice)

	symbolInfo, err := InnerExchangeManager.GetSymbolInfo(GATE_NAME.String(), req.AccountType, res.Data.Contract)
	if err != nil {
		log.Error(err)
		return nil
	}
	cumQuoteQty := executedQty.Mul(symbolInfo.ContractSize()).Mul(fillPrice)

	var updateTime int64
	if decimal.NewFromFloat(res.Data.FinishTime).IsZero() {
		updateTime = decimal.NewFromFloat(res.Data.CreateTime).Mul(gateTimeMul).IntPart()
	} else {
		updateTime = decimal.NewFromFloat(res.Data.FinishTime).Mul(gateTimeMul).IntPart()
	}

	return &Order{
		Exchange:      g.ExchangeType().String(),
		AccountType:   req.AccountType,
		Symbol:        res.Data.Contract,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
		OrderId:       decimal.NewFromInt(res.Data.Id).String(),
		ClientOrderId: res.Data.Text,
		Price:         res.Data.Price,
		Quantity:      amt.Abs().String(),
		ExecutedQty:   executedQty.String(),
		CumQuoteQty:   cumQuoteQty.String(),
		AvgPrice:      res.Data.FillPrice,
		Status:        g.gateConverter.FromGateContractOrderStatus(res.Data.Status, res.Data.FinishAs),
		Type:          req.OrderType,
		Side:          req.OrderSide,
		PositionSide:  req.PositionSide,

		TimeInForce:           g.gateConverter.FromGateTimeInForce(res.Data.Tif),
		FeeAmount:             "",
		FeeCcy:                "",
		ReduceOnly:            res.Data.IsReduceOnly,
		CreateTime:            decimal.NewFromFloat(res.Data.CreateTime).Mul(gateTimeMul).IntPart(),
		UpdateTime:            updateTime,
		RealizedPnl:           "",
		AttachTpTriggerPrice:  "",
		AttachTpOrdPrice:      "",
		AttachSlTriggerPrice:  "",
		AttachSlOrdPrice:      "",
		IsAlgo:                req.IsAlgo,
		OrderAlgoType:         "",
		TriggerPrice:          "",
		TriggerType:           "",
		TriggerConditionType:  "",
		OcoTpTriggerPrice:     "",
		OcoTpOrdType:          "",
		OcoTpOrdPrice:         "",
		OcoSlTriggerPrice:     "",
		OcoSlOrdType:          "",
		OcoSlOrdPrice:         "",
		MarginBuyBorrowAmount: "",
		MarginBuyBorrowAsset:  "",
		ErrorCode:             "",
		ErrorMsg:              "",
	}
}
func (g *GateTradeEngine) handleOrderFromFuturesPriceOrderCreate(req *OrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestFuturesSettlePriceOrdersPostRes]) *Order {
	amt := req.Quantity
	return &Order{
		Exchange:             g.ExchangeType().String(),
		AccountType:          req.AccountType,
		Symbol:               req.Symbol,
		IsMargin:             req.IsMargin,
		IsIsolated:           req.IsIsolated,
		OrderId:              strconv.FormatInt(res.Data.Id, 10),
		ClientOrderId:        req.ClientOrderId,
		Price:                req.Price.String(),
		Quantity:             amt.Abs().String(),
		ExecutedQty:          "0",
		CumQuoteQty:          "0",
		AvgPrice:             "0",
		Status:               ORDER_STATUS_UN_TRIGGERED,
		Type:                 req.OrderType,
		Side:                 req.OrderSide,
		TimeInForce:          req.TimeInForce,
		FeeAmount:            "",
		FeeCcy:               "",
		CreateTime:           time.Now().UnixMilli(),
		UpdateTime:           time.Now().UnixMilli(),
		RealizedPnl:          "",
		AttachTpTriggerPrice: "",
		AttachTpOrdPrice:     "",
		AttachSlTriggerPrice: "",
		AttachSlOrdPrice:     "",
		IsAlgo:               req.IsAlgo,
		OrderAlgoType:        "", // 仅支持单向止盈止损
		TriggerPrice:         req.TriggerPrice.String(),
		TriggerType:          req.TriggerType,
		TriggerConditionType: g.gateConverter.FromGateTriggerCondition(req.OrderSide, req.TriggerType),
		OcoTpTriggerPrice:    "",
		OcoTpOrdType:         "",
		OcoTpOrdPrice:        "",
		OcoSlTriggerPrice:    "",
		OcoSlOrdType:         "",
		OcoSlOrdPrice:        "",
	}
}
func (g *GateTradeEngine) handleOrderFromFuturesOrderAmend(req *OrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestFuturesSettleOrdersOrderIdPutRes]) *Order {
	amt := decimal.NewFromInt(res.Data.Size)
	executedQty := decimal.NewFromInt(res.Data.Size).Sub(decimal.NewFromInt(res.Data.Left))
	fillPrice, _ := decimal.NewFromString(res.Data.FillPrice)

	symbolInfo, err := InnerExchangeManager.GetSymbolInfo(GATE_NAME.String(), req.AccountType, res.Data.Contract)
	if err != nil {
		log.Error(err)
		return nil
	}
	cumQuoteQty := executedQty.Mul(symbolInfo.ContractSize()).Mul(fillPrice)

	var updateTime int64
	if decimal.NewFromFloat(res.Data.FinishTime).IsZero() {
		updateTime = decimal.NewFromFloat(res.Data.CreateTime).Mul(gateTimeMul).IntPart()
	} else {
		updateTime = decimal.NewFromFloat(res.Data.FinishTime).Mul(gateTimeMul).IntPart()
	}

	orderType := ORDER_TYPE_LIMIT
	price, _ := decimal.NewFromString(res.Data.Price)
	if price.IsZero() && res.Data.Tif == GATE_TIME_IN_FORCE_IOC {
		orderType = ORDER_TYPE_MARKET
	}

	orderSide := ORDER_SIDE_BUY
	if res.Data.Size < 0 {
		orderSide = ORDER_SIDE_SELL
	}

	return &Order{
		Exchange:             g.ExchangeType().String(),
		AccountType:          req.AccountType,
		Symbol:               res.Data.Contract,
		IsMargin:             req.IsMargin,
		IsIsolated:           req.IsIsolated,
		OrderId:              decimal.NewFromInt(res.Data.Id).String(),
		ClientOrderId:        res.Data.Text,
		Price:                res.Data.Price,
		Quantity:             amt.Abs().String(),
		ExecutedQty:          executedQty.String(),
		CumQuoteQty:          cumQuoteQty.String(),
		AvgPrice:             res.Data.FillPrice,
		Status:               g.gateConverter.FromGateContractOrderStatus(res.Data.Status, res.Data.FinishAs),
		Type:                 orderType,
		Side:                 orderSide,
		PositionSide:         req.PositionSide,
		TimeInForce:          g.gateConverter.FromGateTimeInForce(res.Data.Tif),
		FeeAmount:            "",
		FeeCcy:               "",
		ReduceOnly:           res.Data.IsReduceOnly,
		CreateTime:           decimal.NewFromFloat(res.Data.CreateTime).Mul(gateTimeMul).IntPart(),
		UpdateTime:           updateTime,
		RealizedPnl:          "",
		AttachTpTriggerPrice: "",
		AttachTpOrdPrice:     "",
		AttachSlTriggerPrice: "",

		AttachSlOrdPrice:      "",
		IsAlgo:                req.IsAlgo,
		OrderAlgoType:         "",
		TriggerPrice:          "",
		TriggerType:           "",
		TriggerConditionType:  "",
		OcoTpTriggerPrice:     "",
		OcoTpOrdType:          "",
		OcoTpOrdPrice:         "",
		OcoSlTriggerPrice:     "",
		OcoSlOrdType:          "",
		OcoSlOrdPrice:         "",
		MarginBuyBorrowAmount: "",
		MarginBuyBorrowAsset:  "",
		ErrorCode:             "",
		ErrorMsg:              "",
	}
}
func (g *GateTradeEngine) handleOrderFromFuturesOrderCancel(req *OrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestFuturesSettleOrdersOrderIdDeleteRes]) *Order {
	amt := decimal.NewFromInt(res.Data.Size)
	executedQty := decimal.NewFromInt(res.Data.Size).Sub(decimal.NewFromInt(res.Data.Left)).Abs()
	fillPrice, _ := decimal.NewFromString(res.Data.FillPrice)

	symbolInfo, err := InnerExchangeManager.GetSymbolInfo(GATE_NAME.String(), req.AccountType, res.Data.Contract)
	if err != nil {
		log.Error(err)
		return nil
	}
	cumQuoteQty := executedQty.Mul(symbolInfo.ContractSize()).Mul(fillPrice)
	var updateTime int64
	if decimal.NewFromFloat(res.Data.FinishTime).IsZero() {
		updateTime = decimal.NewFromFloat(res.Data.CreateTime).Mul(gateTimeMul).IntPart()
	} else {
		updateTime = decimal.NewFromFloat(res.Data.FinishTime).Mul(gateTimeMul).IntPart()
	}

	orderType := ORDER_TYPE_LIMIT
	price, _ := decimal.NewFromString(res.Data.Price)
	if price.IsZero() && res.Data.Tif == GATE_TIME_IN_FORCE_IOC {
		orderType = ORDER_TYPE_MARKET
	}

	orderSide := ORDER_SIDE_BUY
	if res.Data.Size < 0 {
		orderSide = ORDER_SIDE_SELL
	}

	positionSide, err := g.getPositionSide(req.AccountType, res.Data.Contract, res.Data.Size, res.Data.ReduceOnly)
	if err != nil {
		log.Error(err)
	}

	return &Order{
		Exchange:             g.ExchangeType().String(),
		AccountType:          req.AccountType,
		Symbol:               res.Data.Contract,
		IsMargin:             req.IsMargin,
		IsIsolated:           req.IsIsolated,
		OrderId:              decimal.NewFromInt(res.Data.Id).String(),
		ClientOrderId:        res.Data.Text,
		Price:                res.Data.Price,
		Quantity:             amt.Abs().String(),
		ExecutedQty:          executedQty.String(),
		CumQuoteQty:          cumQuoteQty.String(),
		AvgPrice:             res.Data.FillPrice,
		Status:               g.gateConverter.FromGateContractOrderStatus(res.Data.Status, res.Data.FinishAs),
		Type:                 orderType,
		Side:                 orderSide,
		PositionSide:         positionSide,
		TimeInForce:          g.gateConverter.FromGateTimeInForce(res.Data.Tif),
		FeeAmount:            "",
		FeeCcy:               "",
		ReduceOnly:           res.Data.ReduceOnly,
		CreateTime:           decimal.NewFromFloat(res.Data.CreateTime).Mul(gateTimeMul).IntPart(),
		UpdateTime:           updateTime,
		RealizedPnl:          "",
		AttachTpTriggerPrice: "",
		AttachTpOrdPrice:     "",
		AttachSlTriggerPrice: "",

		AttachSlOrdPrice:      "",
		IsAlgo:                req.IsAlgo,
		OrderAlgoType:         "",
		TriggerPrice:          "",
		TriggerType:           "",
		TriggerConditionType:  "",
		OcoTpTriggerPrice:     "",
		OcoTpOrdType:          "",
		OcoTpOrdPrice:         "",
		OcoSlTriggerPrice:     "",
		OcoSlOrdType:          "",
		OcoSlOrdPrice:         "",
		MarginBuyBorrowAmount: "",
		MarginBuyBorrowAsset:  "",
		ErrorCode:             "",
		ErrorMsg:              "",
	}
}
func (g *GateTradeEngine) handleOrderFromFuturesPriceOrderCancel(req *OrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestFuturesSettlePriceOrdersOrderIdDeleteRes]) *Order {
	amt := decimal.NewFromInt(res.Data.Initial.Size)

	var updateTime int64
	if decimal.NewFromFloat(res.Data.FinishTime).IsZero() {
		updateTime = decimal.NewFromFloat(res.Data.CreateTime).Mul(gateTimeMul).IntPart()
	} else {
		updateTime = decimal.NewFromFloat(res.Data.FinishTime).Mul(gateTimeMul).IntPart()
	}

	orderType := ORDER_TYPE_LIMIT
	price, _ := decimal.NewFromString(res.Data.Initial.Price)
	if price.IsZero() && res.Data.Initial.Tif == GATE_TIME_IN_FORCE_IOC {
		orderType = ORDER_TYPE_MARKET
	}

	orderSide := ORDER_SIDE_BUY
	if res.Data.Initial.Size < 0 {
		orderSide = ORDER_SIDE_SELL
	}

	triggerType := g.gateConverter.FromGateFuturesPriceOrderTriggerRule(res.Data.Trigger.Rule, orderSide)
	triggerConditionType := g.gateConverter.FromGateTriggerCondition(orderSide, triggerType)

	positionSide, err := g.getPositionSide(req.AccountType, res.Data.Initial.Contract, res.Data.Initial.Size, res.Data.Initial.IsReduceOnly)
	if err != nil {
		log.Error(err)
	}

	return &Order{
		Exchange:             g.ExchangeType().String(),
		AccountType:          req.AccountType,
		Symbol:               res.Data.Initial.Contract,
		IsMargin:             req.IsMargin,
		IsIsolated:           req.IsIsolated,
		OrderId:              strconv.FormatInt(res.Data.Id, 10),
		ClientOrderId:        res.Data.Initial.Text,
		Price:                res.Data.Initial.Price,
		Quantity:             amt.Abs().String(),
		ExecutedQty:          "0",
		CumQuoteQty:          "0",
		AvgPrice:             "0",
		Status:               g.gateConverter.FromGateContractPriceOrderStatus(res.Data.Status, res.Data.FinishAs),
		Type:                 orderType,
		Side:                 orderSide,
		PositionSide:         positionSide,
		TimeInForce:          g.gateConverter.FromGateTimeInForce(res.Data.Initial.Tif),
		FeeAmount:            "",
		FeeCcy:               "",
		ReduceOnly:           res.Data.Initial.ReduceOnly,
		CreateTime:           decimal.NewFromFloat(res.Data.CreateTime).Mul(gateTimeMul).IntPart(),
		UpdateTime:           updateTime,
		IsAlgo:               true,
		RealizedPnl:          "",
		AttachTpTriggerPrice: "",
		AttachTpOrdPrice:     "",
		AttachSlTriggerPrice: "",
		TriggerPrice:         res.Data.Trigger.Price,
		TriggerType:          triggerType,
		TriggerConditionType: triggerConditionType,
		OcoTpTriggerPrice:    "",
		OcoTpOrdType:         "",
		OcoTpOrdPrice:        "",
		OcoSlTriggerPrice:    "",
		OcoSlOrdType:         "",
		OcoSlOrdPrice:        "",
	}
}

// 交割合约订单查询
func (g *GateTradeEngine) handleOrdersFromDeliveryOpenOrders(req *QueryOrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestDeliverySettleOrdersGetRes]) []*Order {
	var orders []*Order
	for _, order := range res.Data {
		amt := decimal.NewFromInt(order.Size)
		executedQty := decimal.NewFromInt(order.Size).Sub(decimal.NewFromInt(order.Left)).Abs()
		fillPrice, _ := decimal.NewFromString(order.FillPrice)
		symbolInfo, err := InnerExchangeManager.GetSymbolInfo(GATE_NAME.String(), req.AccountType, order.Contract)
		if err != nil {
			log.Error(err)
			continue
		}
		cumQuoteQty := executedQty.Mul(symbolInfo.ContractSize()).Mul(fillPrice)
		var updateTime int64

		if decimal.NewFromFloat(order.FinishTime).IsZero() {
			updateTime = decimal.NewFromFloat(order.CreateTime).Mul(gateTimeMul).IntPart()
		} else {
			updateTime = decimal.NewFromFloat(order.FinishTime).Mul(gateTimeMul).IntPart()
		}

		var orderType OrderType
		if order.Price == "0" && order.Tif == GATE_TIME_IN_FORCE_IOC {
			orderType = ORDER_TYPE_MARKET
		} else {
			orderType = ORDER_TYPE_LIMIT
		}

		var orderSide OrderSide
		if order.Size > 0 {
			orderSide = ORDER_SIDE_BUY
		} else {
			orderSide = ORDER_SIDE_SELL
		}

		positionSide, err := g.getPositionSide(req.AccountType, order.Contract, order.Size, order.ReduceOnly)
		if err != nil {
			log.Error(err)
		}

		orders = append(orders, &Order{
			Exchange:             g.ExchangeType().String(),
			AccountType:          req.AccountType,
			Symbol:               order.Contract,
			IsMargin:             req.IsMargin,
			IsIsolated:           req.IsIsolated,
			OrderId:              strconv.FormatInt(order.Id, 10),
			ClientOrderId:        order.Text,
			Price:                order.Price,
			Quantity:             amt.Abs().String(),
			ExecutedQty:          executedQty.String(),
			CumQuoteQty:          cumQuoteQty.String(),
			AvgPrice:             order.FillPrice,
			Status:               g.gateConverter.FromGateContractOrderStatus(order.Status, order.FinishAs),
			Type:                 orderType,
			Side:                 orderSide,
			PositionSide:         positionSide,
			TimeInForce:          g.gateConverter.FromGateTimeInForce(order.Tif),
			FeeAmount:            "",
			FeeCcy:               "",
			ReduceOnly:           order.ReduceOnly,
			CreateTime:           decimal.NewFromFloat(order.CreateTime).Mul(gateTimeMul).IntPart(),
			UpdateTime:           updateTime,
			RealizedPnl:          "",
			AttachTpTriggerPrice: "",
			AttachTpOrdPrice:     "",
			AttachSlTriggerPrice: "",

			AttachSlOrdPrice:      "",
			IsAlgo:                req.IsAlgo,
			OrderAlgoType:         "",
			TriggerPrice:          "",
			TriggerType:           "",
			TriggerConditionType:  "",
			OcoTpTriggerPrice:     "",
			OcoTpOrdType:          "",
			OcoTpOrdPrice:         "",
			OcoSlTriggerPrice:     "",
			OcoSlOrdType:          "",
			OcoSlOrdPrice:         "",
			MarginBuyBorrowAmount: "",
			MarginBuyBorrowAsset:  "",
			ErrorCode:             "",
			ErrorMsg:              "",
		})
	}

	return orders
}
func (g *GateTradeEngine) handleOrdersFromDeliveryPriceOpenOrders(req *QueryOrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestDeliverySettlePriceOrdersGetRes]) []*Order {
	var orders []*Order
	for _, order := range res.Data {

		amt := decimal.NewFromInt(order.Initial.Size)

		var updateTime int64
		if decimal.NewFromFloat(order.FinishTime).IsZero() {
			updateTime = decimal.NewFromFloat(order.CreateTime).Mul(gateTimeMul).IntPart()
		} else {
			updateTime = decimal.NewFromFloat(order.FinishTime).Mul(gateTimeMul).IntPart()
		}

		var orderType OrderType
		if order.Initial.Price == "0" && order.Initial.Tif == GATE_TIME_IN_FORCE_IOC {
			orderType = ORDER_TYPE_MARKET
		} else {
			orderType = ORDER_TYPE_LIMIT
		}

		var orderSide OrderSide
		if order.Initial.Size > 0 {
			orderSide = ORDER_SIDE_BUY
		} else {
			orderSide = ORDER_SIDE_SELL
		}

		triggerType := g.gateConverter.FromGateFuturesPriceOrderTriggerRule(order.Trigger.Rule, orderSide)
		triggerConditionType := g.gateConverter.FromGateTriggerCondition(orderSide, triggerType)

		positionSide, err := g.getPositionSide(req.AccountType, order.Initial.Contract, order.Initial.Size, order.Initial.IsReduceOnly)
		if err != nil {
			log.Error(err)
		}

		orders = append(orders, &Order{
			Exchange:             g.ExchangeType().String(),
			AccountType:          req.AccountType,
			Symbol:               order.Initial.Contract,
			IsMargin:             req.IsMargin,
			IsIsolated:           req.IsIsolated,
			OrderId:              strconv.FormatInt(order.Id, 10),
			ClientOrderId:        order.Initial.Text,
			Price:                order.Initial.Price,
			Quantity:             amt.Abs().String(),
			ExecutedQty:          "0",
			CumQuoteQty:          "0",
			AvgPrice:             "0",
			Status:               g.gateConverter.FromGateContractPriceOrderStatus(order.Status, order.FinishAs),
			Type:                 orderType,
			Side:                 orderSide,
			PositionSide:         positionSide,
			TimeInForce:          g.gateConverter.FromGateTimeInForce(order.Initial.Tif),
			FeeAmount:            "",
			FeeCcy:               "",
			ReduceOnly:           order.Initial.IsReduceOnly,
			CreateTime:           decimal.NewFromFloat(order.CreateTime).Mul(gateTimeMul).IntPart(),
			UpdateTime:           updateTime,
			IsAlgo:               true,
			RealizedPnl:          "",
			AttachTpTriggerPrice: "",
			AttachTpOrdPrice:     "",
			AttachSlTriggerPrice: "",
			TriggerPrice:         order.Trigger.Price,
			TriggerType:          triggerType,
			TriggerConditionType: triggerConditionType,
			OcoTpTriggerPrice:    "",
			OcoTpOrdType:         "",
			OcoTpOrdPrice:        "",
			OcoSlTriggerPrice:    "",
			OcoSlOrdType:         "",
			OcoSlOrdPrice:        "",
		})
	}

	return orders
}
func (g *GateTradeEngine) handleOrderFromDeliveryOrderQuery(req *QueryOrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestDeliverySettleOrdersOrderIdGetRes]) *Order {
	amt := decimal.NewFromInt(res.Data.Size)
	executedQty := decimal.NewFromInt(res.Data.Size).Sub(decimal.NewFromInt(res.Data.Left)).Abs()
	fillPrice, _ := decimal.NewFromString(res.Data.FillPrice)
	symbolInfo, err := InnerExchangeManager.GetSymbolInfo(GATE_NAME.String(), req.AccountType, res.Data.Contract)
	if err != nil {
		log.Error(err)
		return nil
	}
	cumQuoteQty := executedQty.Mul(symbolInfo.ContractSize()).Mul(fillPrice)
	var updateTime int64

	if decimal.NewFromFloat(res.Data.FinishTime).IsZero() {
		updateTime = decimal.NewFromFloat(res.Data.CreateTime).Mul(gateTimeMul).IntPart()
	} else {
		updateTime = decimal.NewFromFloat(res.Data.FinishTime).Mul(gateTimeMul).IntPart()
	}

	var orderType OrderType
	if res.Data.Price == "0" && res.Data.Tif == GATE_TIME_IN_FORCE_IOC {
		orderType = ORDER_TYPE_MARKET
	} else {
		orderType = ORDER_TYPE_LIMIT
	}

	var orderSide OrderSide
	if res.Data.Size > 0 {
		orderSide = ORDER_SIDE_BUY
	} else {
		orderSide = ORDER_SIDE_SELL
	}

	positionSide, err := g.getPositionSide(req.AccountType, res.Data.Contract, res.Data.Size, res.Data.ReduceOnly)
	if err != nil {
		log.Error(err)
	}

	return &Order{
		Exchange:      g.ExchangeType().String(),
		AccountType:   req.AccountType,
		Symbol:        res.Data.Contract,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
		OrderId:       strconv.FormatInt(res.Data.Id, 10),
		ClientOrderId: res.Data.Text,
		Price:         res.Data.Price,
		Quantity:      amt.Abs().String(),

		ExecutedQty:          executedQty.String(),
		CumQuoteQty:          cumQuoteQty.String(),
		AvgPrice:             res.Data.FillPrice,
		Status:               g.gateConverter.FromGateContractOrderStatus(res.Data.Status, res.Data.FinishAs),
		Type:                 orderType,
		Side:                 orderSide,
		PositionSide:         positionSide,
		TimeInForce:          g.gateConverter.FromGateTimeInForce(res.Data.Tif),
		FeeAmount:            "",
		FeeCcy:               "",
		ReduceOnly:           res.Data.ReduceOnly,
		CreateTime:           decimal.NewFromFloat(res.Data.CreateTime).Mul(gateTimeMul).IntPart(),
		UpdateTime:           updateTime,
		RealizedPnl:          "",
		AttachTpTriggerPrice: "",
		AttachTpOrdPrice:     "",
		AttachSlTriggerPrice: "",

		AttachSlOrdPrice:      "",
		IsAlgo:                req.IsAlgo,
		OrderAlgoType:         "",
		TriggerPrice:          "",
		TriggerType:           "",
		TriggerConditionType:  "",
		OcoTpTriggerPrice:     "",
		OcoTpOrdType:          "",
		OcoTpOrdPrice:         "",
		OcoSlTriggerPrice:     "",
		OcoSlOrdType:          "",
		OcoSlOrdPrice:         "",
		MarginBuyBorrowAmount: "",
		MarginBuyBorrowAsset:  "",
		ErrorCode:             "",
		ErrorMsg:              "",
	}
}
func (g *GateTradeEngine) handleOrderFromDeliveryPriceOrderQuery(req *QueryOrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestDeliverySettlePriceOrdersOrderIdGetRes]) *Order {
	amt := decimal.NewFromInt(res.Data.Initial.Size)
	updateTime := decimal.NewFromFloat(res.Data.CreateTime).Mul(gateTimeMul).IntPart()
	if updateTime == 0 {
		updateTime = decimal.NewFromFloat(res.Data.FinishTime).Mul(gateTimeMul).IntPart()
	}

	var orderType OrderType
	if res.Data.Initial.Price == "0" && res.Data.Initial.Tif == GATE_TIME_IN_FORCE_IOC {
		orderType = ORDER_TYPE_MARKET
	} else {
		orderType = ORDER_TYPE_LIMIT
	}

	var orderSide OrderSide
	if res.Data.Initial.Size > 0 {
		orderSide = ORDER_SIDE_BUY
	} else {
		orderSide = ORDER_SIDE_SELL
	}

	triggerType := g.gateConverter.FromGateFuturesPriceOrderTriggerRule(res.Data.Trigger.Rule, orderSide)
	triggerConditionType := g.gateConverter.FromGateTriggerCondition(orderSide, triggerType)

	positionSide, err := g.getPositionSide(req.AccountType, res.Data.Initial.Contract, res.Data.Initial.Size, res.Data.Initial.IsReduceOnly)
	if err != nil {
		log.Error(err)
	}

	return &Order{
		Exchange:             g.ExchangeType().String(),
		AccountType:          req.AccountType,
		Symbol:               res.Data.Initial.Contract,
		IsMargin:             req.IsMargin,
		IsIsolated:           req.IsIsolated,
		OrderId:              strconv.FormatInt(res.Data.Id, 10),
		ClientOrderId:        res.Data.Initial.Text,
		Price:                res.Data.Initial.Price,
		Quantity:             amt.Abs().String(),
		ExecutedQty:          "0",
		CumQuoteQty:          "0",
		AvgPrice:             "0",
		Status:               g.gateConverter.FromGateContractPriceOrderStatus(res.Data.Status, res.Data.FinishAs),
		Type:                 orderType,
		Side:                 orderSide,
		PositionSide:         positionSide,
		TimeInForce:          g.gateConverter.FromGateTimeInForce(res.Data.Initial.Tif),
		FeeAmount:            "",
		FeeCcy:               "",
		ReduceOnly:           res.Data.Initial.IsReduceOnly,
		CreateTime:           decimal.NewFromFloat(res.Data.CreateTime).Mul(gateTimeMul).IntPart(),
		UpdateTime:           updateTime,
		IsAlgo:               true,
		RealizedPnl:          "",
		AttachTpTriggerPrice: "",
		AttachTpOrdPrice:     "",
		AttachSlTriggerPrice: "",
		TriggerPrice:         res.Data.Trigger.Price,
		TriggerType:          triggerType,
		TriggerConditionType: triggerConditionType,
		OcoTpTriggerPrice:    "",
		OcoTpOrdType:         "",
		OcoTpOrdPrice:        "",
		OcoSlTriggerPrice:    "",
		OcoSlOrdType:         "",
		OcoSlOrdPrice:        "",
	}
}
func (g *GateTradeEngine) handleOrdersFromDeliveryOrdersQuery(req *QueryOrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestDeliverySettleOrdersGetRes]) []*Order {
	var orders []*Order
	for _, order := range res.Data {
		amt := decimal.NewFromInt(order.Size)
		executedQty := decimal.NewFromInt(order.Size).Sub(decimal.NewFromInt(order.Left))
		fillPrice, _ := decimal.NewFromString(order.FillPrice)
		symbolInfo, err := InnerExchangeManager.GetSymbolInfo(GATE_NAME.String(), req.AccountType, order.Contract)
		if err != nil {
			log.Error(err)
			continue
		}
		cumQuoteQty := executedQty.Mul(symbolInfo.ContractSize()).Mul(fillPrice)
		var updateTime int64
		if decimal.NewFromFloat(order.FinishTime).IsZero() {
			updateTime = decimal.NewFromFloat(order.CreateTime).Mul(gateTimeMul).IntPart()
		} else {
			updateTime = decimal.NewFromFloat(order.FinishTime).Mul(gateTimeMul).IntPart()
		}

		var orderType OrderType
		if order.Price == "0" && order.Tif == GATE_TIME_IN_FORCE_IOC {
			orderType = ORDER_TYPE_MARKET
		} else {
			orderType = ORDER_TYPE_LIMIT
		}

		var orderSide OrderSide
		if order.Size > 0 {
			orderSide = ORDER_SIDE_BUY
		} else {
			orderSide = ORDER_SIDE_SELL
		}

		positionSide, err := g.getPositionSide(req.AccountType, order.Contract, order.Size, order.IsReduceOnly)
		if err != nil {
			log.Error(err)
		}

		orders = append(orders, &Order{
			Exchange:             g.ExchangeType().String(),
			AccountType:          req.AccountType,
			Symbol:               order.Contract,
			IsMargin:             req.IsMargin,
			IsIsolated:           req.IsIsolated,
			OrderId:              strconv.FormatInt(order.Id, 10),
			ClientOrderId:        order.Text,
			Price:                order.Price,
			Quantity:             amt.Abs().String(),
			ExecutedQty:          executedQty.Abs().String(),
			CumQuoteQty:          cumQuoteQty.Abs().String(),
			AvgPrice:             order.FillPrice,
			Status:               g.gateConverter.FromGateContractOrderStatus(order.Status, order.FinishAs),
			Type:                 orderType,
			Side:                 orderSide,
			PositionSide:         positionSide,
			TimeInForce:          g.gateConverter.FromGateTimeInForce(order.Tif),
			FeeAmount:            "",
			FeeCcy:               "",
			ReduceOnly:           order.IsReduceOnly,
			CreateTime:           decimal.NewFromFloat(order.CreateTime).Mul(gateTimeMul).IntPart(),
			UpdateTime:           updateTime,
			RealizedPnl:          "",
			AttachTpTriggerPrice: "",
			AttachTpOrdPrice:     "",
			AttachSlTriggerPrice: "",

			AttachSlOrdPrice:      "",
			IsAlgo:                req.IsAlgo,
			OrderAlgoType:         "",
			TriggerPrice:          "",
			TriggerType:           "",
			TriggerConditionType:  "",
			OcoTpTriggerPrice:     "",
			OcoTpOrdType:          "",
			OcoTpOrdPrice:         "",
			OcoSlTriggerPrice:     "",
			OcoSlOrdType:          "",
			OcoSlOrdPrice:         "",
			MarginBuyBorrowAmount: "",
			MarginBuyBorrowAsset:  "",
			ErrorCode:             "",
			ErrorMsg:              "",
		})
	}

	return orders
}
func (g *GateTradeEngine) handleOrdersFromDeliveryPriceOrdersQuery(req *QueryOrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestDeliverySettlePriceOrdersGetRes]) []*Order {
	var orders []*Order
	for _, order := range res.Data {
		amt := decimal.NewFromInt(order.Initial.Size)
		var updateTime int64
		if decimal.NewFromFloat(order.FinishTime).IsZero() {
			updateTime = decimal.NewFromFloat(order.CreateTime).Mul(gateTimeMul).IntPart()
		} else {
			updateTime = decimal.NewFromFloat(order.FinishTime).Mul(gateTimeMul).IntPart()
		}

		var orderType OrderType
		if order.Initial.Price == "0" && order.Initial.Tif == GATE_TIME_IN_FORCE_IOC {
			orderType = ORDER_TYPE_MARKET
		} else {
			orderType = ORDER_TYPE_LIMIT
		}

		var orderSide OrderSide
		if order.Initial.Size > 0 {
			orderSide = ORDER_SIDE_BUY
		} else {
			orderSide = ORDER_SIDE_SELL
		}

		triggerType := g.gateConverter.FromGateFuturesPriceOrderTriggerRule(order.Trigger.Rule, orderSide)
		triggerConditionType := g.gateConverter.FromGateTriggerCondition(orderSide, triggerType)

		positionSide, err := g.getPositionSide(req.AccountType, order.Initial.Contract, order.Initial.Size, order.Initial.IsReduceOnly)
		if err != nil {
			log.Error(err)
		}

		orders = append(orders, &Order{
			Exchange:             g.ExchangeType().String(),
			AccountType:          req.AccountType,
			Symbol:               order.Initial.Contract,
			IsMargin:             req.IsMargin,
			IsIsolated:           req.IsIsolated,
			OrderId:              strconv.FormatInt(order.Id, 10),
			ClientOrderId:        order.Initial.Text,
			Price:                order.Initial.Price,
			Quantity:             amt.Abs().String(),
			ExecutedQty:          "0",
			CumQuoteQty:          "0",
			AvgPrice:             "0",
			Status:               g.gateConverter.FromGateContractPriceOrderStatus(order.Status, order.FinishAs),
			Type:                 orderType,
			Side:                 orderSide,
			PositionSide:         positionSide,
			TimeInForce:          g.gateConverter.FromGateTimeInForce(order.Initial.Tif),
			FeeAmount:            "",
			FeeCcy:               "",
			ReduceOnly:           order.Initial.IsReduceOnly,
			CreateTime:           decimal.NewFromFloat(order.CreateTime).Mul(gateTimeMul).IntPart(),
			UpdateTime:           updateTime,
			IsAlgo:               true,
			RealizedPnl:          "",
			AttachTpTriggerPrice: "",
			AttachTpOrdPrice:     "",
			AttachSlTriggerPrice: "",
			TriggerPrice:         order.Trigger.Price,
			TriggerType:          triggerType,
			TriggerConditionType: triggerConditionType,
			OcoTpTriggerPrice:    "",
			OcoTpOrdType:         "",
			OcoTpOrdPrice:        "",
			OcoSlTriggerPrice:    "",
			OcoSlOrdType:         "",
			OcoSlOrdPrice:        "",
		})
	}
	return orders
}
func (g *GateTradeEngine) handleTradesFromDeliveryTradesQuery(req *QueryTradeParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestDeliverySettleMyTradesRes]) []*Trade {
	var trades []*Trade
	var err error
	var order *Order
	var positionSide PositionSide

	order, err = g.QueryOrder(&QueryOrderParam{
		AccountType: req.AccountType,
		Symbol:      req.Symbol,
		OrderId:     req.OrderId,
	})
	if err != nil {
		log.Error(err)
	}

	if order != nil {
		positionSide = order.PositionSide
	}

	for _, trade := range res.Data {
		price, _ := decimal.NewFromString(trade.Price)
		amt := decimal.NewFromInt(trade.Size)

		symbolInfo, err := InnerExchangeManager.GetSymbolInfo(GATE_NAME.String(), req.AccountType, trade.Contract)
		if err != nil {
			log.Error(err)
			continue
		}
		quoteQty := amt.Abs().Mul(symbolInfo.ContractSize()).Mul(price)

		isMaker := false
		if trade.Role == "maker" {
			isMaker = true
		}

		var side OrderSide
		if trade.Size > 0 {
			side = ORDER_SIDE_BUY
		} else {
			side = ORDER_SIDE_SELL
		}
		feeCcy := strings.Split(trade.Contract, "_")[1]

		trades = append(trades, &Trade{
			Exchange:      g.ExchangeType().String(),
			AccountType:   req.AccountType,
			Symbol:        req.Symbol,
			TradeId:       strconv.FormatInt(trade.ID, 10),
			OrderId:       trade.OrderID,
			ClientOrderId: trade.Text,
			Price:         trade.Price,
			Quantity:      amt.Abs().String(),
			QuoteQty:      quoteQty.String(),
			Side:          side,
			PositionSide:  positionSide,
			FeeAmount:     trade.Fee,
			FeeCcy:        feeCcy,
			RealizedPnl:   "",
			IsMaker:       isMaker,
			Timestamp:     decimal.NewFromFloat(trade.CreateTime).Mul(gateTimeMul).IntPart(),
		})
	}
	return trades
}

// 交割合约订单操作
func (g *GateTradeEngine) handleOrderFromDeliveryOrderCreate(req *OrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestDeliverySettleOrdersPostRes]) *Order {
	amt := decimal.NewFromInt(res.Data.Size)
	executedQty := decimal.NewFromInt(res.Data.Size).Sub(decimal.NewFromInt(res.Data.Left))
	fillPrice, _ := decimal.NewFromString(res.Data.FillPrice)

	symbolInfo, err := InnerExchangeManager.GetSymbolInfo(GATE_NAME.String(), req.AccountType, res.Data.Contract)
	if err != nil {
		log.Error(err)
		return nil
	}
	cumQuoteQty := executedQty.Abs().Mul(symbolInfo.ContractSize()).Mul(fillPrice)

	// var positionSide PositionSide
	// if res.Data.Size > 0 {
	// 	positionSide = POSITION_SIDE_LONG
	// } else if res.Data.Size < 0 {
	// 	positionSide = POSITION_SIDE_SHORT
	// }
	var updateTime int64
	if decimal.NewFromFloat(res.Data.FinishTime).IsZero() {
		updateTime = decimal.NewFromFloat(res.Data.CreateTime).Mul(gateTimeMul).IntPart()
	} else {
		updateTime = decimal.NewFromFloat(res.Data.FinishTime).Mul(gateTimeMul).IntPart()
	}

	return &Order{
		Exchange:             g.ExchangeType().String(),
		AccountType:          req.AccountType,
		Symbol:               res.Data.Contract,
		IsMargin:             req.IsMargin,
		IsIsolated:           req.IsIsolated,
		OrderId:              decimal.NewFromInt(res.Data.Id).String(),
		ClientOrderId:        res.Data.Text,
		Price:                res.Data.Price,
		Quantity:             amt.Abs().String(),
		ExecutedQty:          executedQty.Abs().String(),
		CumQuoteQty:          cumQuoteQty.String(),
		AvgPrice:             res.Data.FillPrice,
		Status:               g.gateConverter.FromGateContractOrderStatus(res.Data.Status, res.Data.FinishAs),
		Type:                 req.OrderType,
		Side:                 req.OrderSide,
		PositionSide:         req.PositionSide,
		TimeInForce:          g.gateConverter.FromGateTimeInForce(res.Data.Tif),
		FeeAmount:            "",
		FeeCcy:               "",
		ReduceOnly:           res.Data.IsReduceOnly,
		CreateTime:           decimal.NewFromFloat(res.Data.CreateTime).Mul(gateTimeMul).IntPart(),
		UpdateTime:           updateTime,
		RealizedPnl:          "",
		AttachTpTriggerPrice: "",
		AttachTpOrdPrice:     "",
		AttachSlTriggerPrice: "",

		AttachSlOrdPrice:      "",
		IsAlgo:                req.IsAlgo,
		OrderAlgoType:         "",
		TriggerPrice:          "",
		TriggerType:           "",
		TriggerConditionType:  "",
		OcoTpTriggerPrice:     "",
		OcoTpOrdType:          "",
		OcoTpOrdPrice:         "",
		OcoSlTriggerPrice:     "",
		OcoSlOrdType:          "",
		OcoSlOrdPrice:         "",
		MarginBuyBorrowAmount: "",
		MarginBuyBorrowAsset:  "",
		ErrorCode:             "",
		ErrorMsg:              "",
	}
}
func (g *GateTradeEngine) handleOrderFromDeliveryPriceOrderCreate(req *OrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestDeliverySettlePriceOrdersPostRes]) *Order {
	amt := req.Quantity
	return &Order{
		Exchange:             g.ExchangeType().String(),
		AccountType:          req.AccountType,
		Symbol:               req.Symbol,
		IsMargin:             req.IsMargin,
		IsIsolated:           req.IsIsolated,
		OrderId:              strconv.FormatInt(res.Data.ID, 10),
		ClientOrderId:        req.ClientOrderId,
		Price:                req.Price.String(),
		Quantity:             amt.Abs().String(),
		ExecutedQty:          "0",
		CumQuoteQty:          "0",
		AvgPrice:             "0",
		Status:               ORDER_STATUS_UN_TRIGGERED,
		Type:                 req.OrderType,
		Side:                 req.OrderSide,
		TimeInForce:          req.TimeInForce,
		FeeAmount:            "",
		FeeCcy:               "",
		CreateTime:           time.Now().UnixMilli(),
		UpdateTime:           time.Now().UnixMilli(),
		RealizedPnl:          "",
		AttachTpTriggerPrice: "",
		AttachTpOrdPrice:     "",
		AttachSlTriggerPrice: "",
		AttachSlOrdPrice:     "",
		IsAlgo:               req.IsAlgo,
		OrderAlgoType:        "",
		TriggerPrice:         req.TriggerPrice.String(),
		TriggerType:          req.TriggerType,
		TriggerConditionType: g.gateConverter.FromGateTriggerCondition(req.OrderSide, req.TriggerType),
		OcoTpTriggerPrice:    "",
		OcoTpOrdType:         "",
		OcoTpOrdPrice:        "",
		OcoSlTriggerPrice:    "",
		OcoSlOrdType:         "",
		OcoSlOrdPrice:        "",
	}
}
func (g *GateTradeEngine) handleOrderFromDeliveryOrderCancel(req *OrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestDeliverySettleOrdersOrderIdDeleteRes]) *Order {
	amt := decimal.NewFromInt(res.Data.Size)
	executedQty := decimal.NewFromInt(res.Data.Size).Sub(decimal.NewFromInt(res.Data.Left)).Abs()
	fillPrice, _ := decimal.NewFromString(res.Data.FillPrice)

	symbolInfo, err := InnerExchangeManager.GetSymbolInfo(GATE_NAME.String(), req.AccountType, res.Data.Contract)
	if err != nil {
		log.Error(err)
		return nil
	}
	cumQuoteQty := executedQty.Mul(symbolInfo.ContractSize()).Mul(fillPrice)
	// var positionSide PositionSide
	// if res.Data.Size > 0 {
	// 	positionSide = POSITION_SIDE_LONG
	// } else if res.Data.Size < 0 {
	// 	positionSide = POSITION_SIDE_SHORT
	// }
	var updateTime int64
	if decimal.NewFromFloat(res.Data.FinishTime).IsZero() {
		updateTime = decimal.NewFromFloat(res.Data.CreateTime).Mul(gateTimeMul).IntPart()
	} else {
		updateTime = decimal.NewFromFloat(res.Data.FinishTime).Mul(gateTimeMul).IntPart()
	}

	orderType := ORDER_TYPE_LIMIT
	price, _ := decimal.NewFromString(res.Data.Price)
	if price.IsZero() && res.Data.Tif == GATE_TIME_IN_FORCE_IOC {
		orderType = ORDER_TYPE_MARKET
	}

	orderSide := ORDER_SIDE_BUY
	if res.Data.Size < 0 {
		orderSide = ORDER_SIDE_SELL
	}

	positionSide, err := g.getPositionSide(req.AccountType, res.Data.Contract, res.Data.Size, res.Data.ReduceOnly)
	if err != nil {
		log.Error(err)
	}

	return &Order{
		Exchange:             g.ExchangeType().String(),
		AccountType:          req.AccountType,
		Symbol:               res.Data.Contract,
		IsMargin:             req.IsMargin,
		IsIsolated:           req.IsIsolated,
		OrderId:              decimal.NewFromInt(res.Data.Id).String(),
		ClientOrderId:        res.Data.Text,
		Price:                res.Data.Price,
		Quantity:             amt.Abs().String(),
		ExecutedQty:          executedQty.Abs().String(),
		CumQuoteQty:          cumQuoteQty.Abs().String(),
		AvgPrice:             res.Data.FillPrice,
		Status:               g.gateConverter.FromGateContractOrderStatus(res.Data.Status, res.Data.FinishAs),
		Type:                 orderType,
		Side:                 orderSide,
		PositionSide:         positionSide,
		TimeInForce:          g.gateConverter.FromGateTimeInForce(res.Data.Tif),
		FeeAmount:            "",
		FeeCcy:               "",
		ReduceOnly:           res.Data.ReduceOnly,
		CreateTime:           decimal.NewFromFloat(res.Data.CreateTime).Mul(gateTimeMul).IntPart(),
		UpdateTime:           updateTime,
		RealizedPnl:          "",
		AttachTpTriggerPrice: "",
		AttachTpOrdPrice:     "",
		AttachSlTriggerPrice: "",

		AttachSlOrdPrice:      "",
		IsAlgo:                req.IsAlgo,
		OrderAlgoType:         "",
		TriggerPrice:          "",
		TriggerType:           "",
		TriggerConditionType:  "",
		OcoTpTriggerPrice:     "",
		OcoTpOrdType:          "",
		OcoTpOrdPrice:         "",
		OcoSlTriggerPrice:     "",
		OcoSlOrdType:          "",
		OcoSlOrdPrice:         "",
		MarginBuyBorrowAmount: "",
		MarginBuyBorrowAsset:  "",
		ErrorCode:             "",
		ErrorMsg:              "",
	}
}
func (g *GateTradeEngine) handleOrderFromDeliveryPriceOrderCancel(req *OrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestDeliverySettlePriceOrdersOrderIdDeleteRes]) *Order {
	amt := decimal.NewFromInt(res.Data.Initial.Size)
	var updateTime int64
	if decimal.NewFromFloat(res.Data.FinishTime).IsZero() {
		updateTime = decimal.NewFromFloat(res.Data.CreateTime).Mul(gateTimeMul).IntPart()
	} else {
		updateTime = decimal.NewFromFloat(res.Data.FinishTime).Mul(gateTimeMul).IntPart()
	}

	orderType := ORDER_TYPE_LIMIT
	price, _ := decimal.NewFromString(res.Data.Initial.Price)
	if price.IsZero() && res.Data.Initial.Tif == GATE_TIME_IN_FORCE_IOC {
		orderType = ORDER_TYPE_MARKET
	}

	orderSide := ORDER_SIDE_BUY
	if res.Data.Initial.Size < 0 {
		orderSide = ORDER_SIDE_SELL
	}

	triggerType := g.gateConverter.FromGateFuturesPriceOrderTriggerRule(res.Data.Trigger.Rule, orderSide)
	triggerConditionType := g.gateConverter.FromGateTriggerCondition(orderSide, triggerType)

	positionSide, err := g.getPositionSide(req.AccountType, res.Data.Initial.Contract, res.Data.Initial.Size, res.Data.Initial.IsReduceOnly)
	if err != nil {
		log.Error(err)
	}

	return &Order{
		Exchange:             g.ExchangeType().String(),
		AccountType:          req.AccountType,
		Symbol:               res.Data.Initial.Contract,
		IsMargin:             req.IsMargin,
		IsIsolated:           req.IsIsolated,
		OrderId:              strconv.FormatInt(res.Data.Id, 10),
		ClientOrderId:        res.Data.Initial.Text,
		Price:                res.Data.Initial.Price,
		Quantity:             amt.Abs().String(),
		ExecutedQty:          "0",
		CumQuoteQty:          "0",
		AvgPrice:             "0",
		Status:               g.gateConverter.FromGateContractPriceOrderStatus(res.Data.Status, res.Data.FinishAs),
		Type:                 orderType,
		Side:                 orderSide,
		PositionSide:         positionSide,
		TimeInForce:          g.gateConverter.FromGateTimeInForce(res.Data.Initial.Tif),
		FeeAmount:            "",
		FeeCcy:               "",
		ReduceOnly:           res.Data.Initial.IsReduceOnly,
		CreateTime:           decimal.NewFromFloat(res.Data.CreateTime).Mul(gateTimeMul).IntPart(),
		UpdateTime:           updateTime,
		IsAlgo:               true,
		RealizedPnl:          "",
		AttachTpTriggerPrice: "",
		AttachTpOrdPrice:     "",
		AttachSlTriggerPrice: "",
		TriggerPrice:         res.Data.Trigger.Price,
		TriggerType:          triggerType,
		TriggerConditionType: triggerConditionType,
		OcoTpTriggerPrice:    "",
		OcoTpOrdType:         "",
		OcoTpOrdPrice:        "",
		OcoSlTriggerPrice:    "",
		OcoSlOrdType:         "",
		OcoSlOrdPrice:        "",
	}
}
func (g *GateTradeEngine) handleOrderFromBatchErr(req *OrderParam, err error) *Order {
	return &Order{
		Exchange:      GATE_NAME.String(),
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

// handle ws
func (g *GateTradeEngine) handleSubscribeOrderFromSpotSub(req SubscribeOrderParam,
	spotSub *mygateapi.MultipleSubscription[mygateapi.WsSubscribeResult[mygateapi.WsSpotOrder]], newSub *subscription[Order]) {
	//处理订单推送订阅
	go func() {
		for {
			select {
			case err := <-spotSub.ErrChan():
				newSub.errChan <- err
			case <-spotSub.CloseChan():
				newSub.CloseChan() <- struct{}{}
				return
			case result := <-spotSub.ResultChan():
				r := result.Result
				avgPrice := decimal.Zero
				filledAmount, _ := decimal.NewFromString(r.FilledAmount)
				filledTotal, _ := decimal.NewFromString(r.FilledTotal)
				if !filledAmount.IsZero() && !filledTotal.IsZero() {
					avgPrice, _ = decimal.NewFromString(r.AvgDealPrice)
					if avgPrice.IsZero() {
						avgPrice = filledTotal.Div(filledAmount)
					}
				}
				createTime, _ := strconv.ParseInt(r.CreateTimeMs, 10, 64)
				updateTime, _ := strconv.ParseInt(r.UpdateTimeMs, 10, 64)
				order := Order{
					Exchange:      GATE_NAME.String(),
					AccountType:   req.AccountType,
					Symbol:        r.CurrencyPair,
					OrderId:       r.Id,
					ClientOrderId: r.Text,
					Price:         r.Price,
					Quantity:      r.Amount,
					ExecutedQty:   r.FilledAmount,
					CumQuoteQty:   r.FilledTotal,
					AvgPrice:      avgPrice.String(),
					Status:        g.gateConverter.FromGateWsSportOrderStatus(r.Event, r.FinishAs),
					Type:          g.gateConverter.FromGateOrderType(r.Type),
					Side:          g.gateConverter.FromGateOrderSide(r.Side),
					TimeInForce:   g.gateConverter.FromGateTimeInForce(r.TimeInForce),
					FeeAmount:     r.Fee,
					FeeCcy:        r.FeeCurrency,
					CreateTime:    createTime,
					UpdateTime:    updateTime,
				}
				newSub.resultChan <- order
			}
		}
	}()
}

// handle ws
func (g *GateTradeEngine) handleSubscribeOrderFromFuturesOrDeliverySub(req SubscribeOrderParam,
	targetSub *mygateapi.MultipleSubscription[mygateapi.WsSubscribeResult[[]mygateapi.WsFuturesOrder]], newSub *subscription[Order]) {
	//处理订单推送订阅
	go func() {
		for {
			select {
			case err := <-targetSub.ErrChan():
				newSub.errChan <- err
			case <-targetSub.CloseChan():
				newSub.CloseChan() <- struct{}{}
				return
			case result := <-targetSub.ResultChan():
				if result.Result == nil {
					continue
				}
				for _, r := range *result.Result {
					price := decimal.NewFromFloat(r.Price)
					avgPrice := decimal.NewFromFloat(r.FillPrice)
					//size为正买入 size为负卖出
					quantity := decimal.NewFromInt(r.Size).Abs()
					quantityLeft := decimal.NewFromInt(r.Left).Abs()
					filledAmount := quantity.Sub(quantityLeft)
					filledTotal := filledAmount.Mul(avgPrice)

					createTime := r.CreateTimeMs
					updateTime := r.UpdateTime

					var orderType OrderType
					if price.IsZero() {
						//市价单
						orderType = ORDER_TYPE_MARKET
					} else {
						//限价单
						orderType = ORDER_TYPE_LIMIT
					}

					var orderSide OrderSide
					if r.Size > 0 {
						orderSide = ORDER_SIDE_BUY
					} else {
						orderSide = ORDER_SIDE_SELL
					}

					mkfr := decimal.NewFromFloat(r.Mkfr)
					tkfr := decimal.NewFromFloat(r.Tkfr)
					fee := mkfr.Add(tkfr)

					positionSide, err := g.getPositionSide(req.AccountType, r.Contract, r.Size, r.IsReduceOnly)
					if err != nil {
						log.Error(err)
					}

					order := Order{
						Exchange:      GATE_NAME.String(),
						AccountType:   req.AccountType,
						Symbol:        r.Contract,
						OrderId:       strconv.FormatInt(r.Id, 10),
						ClientOrderId: r.Text,
						Price:         price.String(),
						Quantity:      quantity.String(),
						ExecutedQty:   filledAmount.String(),
						CumQuoteQty:   filledTotal.String(),
						AvgPrice:      avgPrice.String(),
						Status:        g.gateConverter.FromGateContractOrderStatus(r.Status, r.FinishAs),
						Type:          orderType,
						Side:          orderSide,
						PositionSide:  positionSide,
						TimeInForce:   g.gateConverter.FromGateTimeInForce(r.Tif),
						FeeAmount:     fee.String(),
						FeeCcy:        "",
						ReduceOnly:    r.IsReduceOnly,
						CreateTime:    createTime,
						UpdateTime:    updateTime,
					}
					newSub.resultChan <- order
				}
			}
		}
	}()
}
