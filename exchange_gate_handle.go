package mytrade

import (
	"strconv"
	"strings"

	"github.com/Hongssd/mygateapi"
	"github.com/shopspring/decimal"
)

var gateTimeMul = decimal.NewFromInt(1000)

func (g *GateTradeEngine) handleOrderFromSpotOrderCreate(req *OrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestSpotOrdersPostRes]) *Order {
	accountType, isMargin, isIsolated := g.gateConverter.FromOrderSpotAccountType(GateAccountType(res.Data.Account))
	return &Order{
		Exchange:      g.ExchangeType().String(),
		AccountType:   accountType.String(),
		Symbol:        res.Data.CurrencyPair,
		IsMargin:      isMargin,
		IsIsolated:    isIsolated,
		OrderId:       res.Data.ID,
		ClientOrderId: res.Data.Text,
		Price:         res.Data.Price,
		Quantity:      res.Data.Amount,
		ExecutedQty:   res.Data.FilledAmount,
		CumQuoteQty:   res.Data.FilledTotal,
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
	accountType, isMargin, isIsolated := g.gateConverter.FromOrderSpotPriceAccountType(req.AccountType)
	return &Order{
		Exchange:             g.ExchangeType().String(),
		AccountType:          accountType.String(),
		Symbol:               req.Symbol,
		IsMargin:             isMargin,
		IsIsolated:           isIsolated,
		OrderId:              strconv.FormatInt(res.Data.ID, 10),
		ClientOrderId:        req.ClientOrderId,
		Price:                req.Price.String(),
		Quantity:             req.Quantity.String(),
		ExecutedQty:          "",
		CumQuoteQty:          "",
		AvgPrice:             "",
		Status:               "",
		Type:                 req.OrderType,
		Side:                 req.OrderSide,
		TimeInForce:          req.TimeInForce,
		FeeAmount:            "",
		FeeCcy:               "",
		CreateTime:           0,
		UpdateTime:           0,
		RealizedPnl:          "",
		AttachTpTriggerPrice: "",
		AttachTpOrdPrice:     "",
		AttachSlTriggerPrice: "",
		AttachSlOrdPrice:     "",
		IsAlgo:               req.IsAlgo,
		OrderAlgoType:        OKX_ORDER_ALGO_TYPE_CONDITIONAL, // 仅支持单向止盈止损
		TriggerPrice:         req.TriggerPrice.String(),
		TriggerType:          req.TriggerType,
		TriggerConditionType: "",
		OcoTpTriggerPrice:    "",
		OcoTpOrdType:         "",
		OcoTpOrdPrice:        "",
		OcoSlTriggerPrice:    "",
		OcoSlOrdType:         "",
		OcoSlOrdPrice:        "",
	}
}
func (g *GateTradeEngine) handleOrderFromFuturesOrderCreate(req *OrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestFuturesSettleOrdersPostRes]) *Order {
	executedQty := decimal.NewFromInt(res.Data.Size).Sub(decimal.NewFromInt(res.Data.Left)).Abs()
	fillPrice, _ := decimal.NewFromString(res.Data.FillPrice)

	symbolInfo, err := InnerExchangeManager.GetSymbolInfo(GATE_NAME.String(), req.AccountType, res.Data.Contract)
	if err != nil {
		log.Error(err)
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
		Quantity:      decimal.NewFromInt(res.Data.Size).Abs().String(),
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
	return &Order{
		Exchange:             g.ExchangeType().String(),
		AccountType:          req.AccountType,
		Symbol:               req.Symbol,
		IsMargin:             req.IsMargin,
		IsIsolated:           req.IsIsolated,
		OrderId:              strconv.FormatInt(res.Data.Id, 10),
		ClientOrderId:        req.ClientOrderId,
		Price:                req.Price.String(),
		Quantity:             req.Quantity.String(),
		ExecutedQty:          "",
		CumQuoteQty:          "",
		AvgPrice:             "",
		Status:               "",
		Type:                 req.OrderType,
		Side:                 req.OrderSide,
		TimeInForce:          req.TimeInForce,
		FeeAmount:            "",
		FeeCcy:               "",
		CreateTime:           0,
		UpdateTime:           0,
		RealizedPnl:          "",
		AttachTpTriggerPrice: "",
		AttachTpOrdPrice:     "",
		AttachSlTriggerPrice: "",
		AttachSlOrdPrice:     "",
		IsAlgo:               req.IsAlgo,
		OrderAlgoType:        OKX_ORDER_ALGO_TYPE_CONDITIONAL, // 仅支持单向止盈止损
		TriggerPrice:         req.TriggerPrice.String(),
		TriggerType:          req.TriggerType,
		TriggerConditionType: "",
		OcoTpTriggerPrice:    "",
		OcoTpOrdType:         "",
		OcoTpOrdPrice:        "",
		OcoSlTriggerPrice:    "",
		OcoSlOrdType:         "",
		OcoSlOrdPrice:        "",
	}
}
func (g *GateTradeEngine) handleOrderFromDeliveryOrderCreate(req *OrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestDeliverySettleOrdersPostRes]) *Order {
	executedQty := decimal.NewFromInt(res.Data.Size).Sub(decimal.NewFromInt(res.Data.Left))
	fillPrice, _ := decimal.NewFromString(res.Data.FillPrice)

	symbolInfo, err := InnerExchangeManager.GetSymbolInfo(GATE_NAME.String(), req.AccountType, res.Data.Contract)
	if err != nil {
		log.Error(err)
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
		Quantity:             decimal.NewFromInt(res.Data.Size).String(),
		ExecutedQty:          executedQty.Abs().String(),
		CumQuoteQty:          cumQuoteQty.String(),
		AvgPrice:             res.Data.FillPrice,
		Status:               g.gateConverter.FromGateContractOrderStatus(res.Data.Status, res.Data.FinishAs),
		Type:                 req.OrderType,
		Side:                 req.OrderSide,
		PositionSide:         "",
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
	return &Order{
		Exchange:             g.ExchangeType().String(),
		AccountType:          req.AccountType,
		Symbol:               req.Symbol,
		IsMargin:             req.IsMargin,
		IsIsolated:           req.IsIsolated,
		OrderId:              strconv.FormatInt(res.Data.ID, 10),
		ClientOrderId:        req.ClientOrderId,
		Price:                req.Price.String(),
		Quantity:             req.Quantity.String(),
		ExecutedQty:          "",
		CumQuoteQty:          "",
		AvgPrice:             "",
		Status:               "",
		Type:                 req.OrderType,
		Side:                 req.OrderSide,
		TimeInForce:          req.TimeInForce,
		FeeAmount:            "",
		FeeCcy:               "",
		CreateTime:           0,
		UpdateTime:           0,
		RealizedPnl:          "",
		AttachTpTriggerPrice: "",
		AttachTpOrdPrice:     "",
		AttachSlTriggerPrice: "",
		AttachSlOrdPrice:     "",
		IsAlgo:               req.IsAlgo,
		OrderAlgoType:        OKX_ORDER_ALGO_TYPE_CONDITIONAL, // 仅支持单向止盈止损
		TriggerPrice:         req.TriggerPrice.String(),
		TriggerType:          req.TriggerType,
		TriggerConditionType: "",
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
	return &Order{
		Exchange:      g.ExchangeType().String(),
		AccountType:   accountType.String(),
		Symbol:        res.Data.CurrencyPair,
		IsMargin:      isMargin,
		IsIsolated:    isIsolated,
		OrderId:       res.Data.ID,
		ClientOrderId: res.Data.Text,
		Price:         res.Data.Price,
		Quantity:      res.Data.Amount,
		ExecutedQty:   res.Data.FilledAmount,
		CumQuoteQty:   res.Data.FilledTotal,
		AvgPrice:      res.Data.AvgDealPrice,
		Status:        g.gateConverter.FromGateSpotOrderStatus(res.Data.Status),
		Type:          g.gateConverter.FromGateOrderType(res.Data.Type),
		Side:          g.gateConverter.FromGateOrderSide(res.Data.Side),
		TimeInForce:   g.gateConverter.FromGateTimeInForce(res.Data.TimeInForce),

		FeeAmount:  res.Data.Fee,
		FeeCcy:     res.Data.FeeCurrency,
		CreateTime: stringToInt64(res.Data.CreateTime),
		UpdateTime: stringToInt64(res.Data.UpdateTime),
	}
}
func (g *GateTradeEngine) handleOrderFromFuturesOrderAmend(req *OrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestFuturesSettleOrdersOrderIdPutRes]) *Order {
	executedQty := decimal.NewFromInt(res.Data.Size).Sub(decimal.NewFromInt(res.Data.Left))
	fillPrice, _ := decimal.NewFromString(res.Data.FillPrice)

	symbolInfo, err := InnerExchangeManager.GetSymbolInfo(GATE_NAME.String(), req.AccountType, res.Data.Contract)
	if err != nil {
		log.Error(err)
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
		Quantity:             decimal.NewFromInt(res.Data.Size).String(),
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

func (g *GateTradeEngine) handleOrderFromSpotOrderCancel(req *OrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestSpotOrdersOrderIdDeleteRes]) *Order {
	accountType, isMargin, isIsolated := g.gateConverter.FromOrderSpotAccountType(GateAccountType(res.Data.Account))
	return &Order{
		Exchange:      g.ExchangeType().String(),
		AccountType:   accountType.String(),
		Symbol:        res.Data.CurrencyPair,
		IsMargin:      isMargin,
		IsIsolated:    isIsolated,
		OrderId:       res.Data.ID,
		ClientOrderId: res.Data.Text,
		Price:         res.Data.Price,
		Quantity:      res.Data.Amount,
		ExecutedQty:   res.Data.FilledAmount,
		CumQuoteQty:   res.Data.FilledTotal,
		AvgPrice:      res.Data.AvgDealPrice,
		Status:        g.gateConverter.FromGateSpotOrderStatus(res.Data.Status),
		Type:          g.gateConverter.FromGateOrderType(res.Data.Type),
		Side:          g.gateConverter.FromGateOrderSide(res.Data.Side),
		PositionSide:  "",

		TimeInForce:           g.gateConverter.FromGateTimeInForce(res.Data.TimeInForce),
		FeeAmount:             res.Data.Fee,
		FeeCcy:                res.Data.FeeCurrency,
		ReduceOnly:            false,
		CreateTime:            decimal.NewFromInt(stringToInt64(res.Data.CreateTime)).Mul(gateTimeMul).IntPart(),
		UpdateTime:            decimal.NewFromInt(stringToInt64(res.Data.UpdateTime)).Mul(gateTimeMul).IntPart(),
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

	return &Order{
		Exchange:             g.ExchangeType().String(),
		AccountType:          accountType.String(),
		Symbol:               res.Data.Market,
		IsMargin:             isMargin,
		IsIsolated:           isIsolated,
		OrderId:              strconv.FormatInt(res.Data.ID, 10),
		ClientOrderId:        "",
		Price:                res.Data.Put.Price,
		Quantity:             res.Data.Put.Amount,
		ExecutedQty:          "",
		CumQuoteQty:          "",
		AvgPrice:             res.Data.Put.Price,
		Status:               g.gateConverter.FromGateSpotPriceOrderStatus(res.Data.Status),
		Type:                 OrderType(g.gateConverter.FromGateOrderType(res.Data.Put.Type)),
		Side:                 OrderSide(g.gateConverter.FromGateOrderSide(res.Data.Put.Side)),
		PositionSide:         "",
		TimeInForce:          g.gateConverter.FromGateTimeInForce(res.Data.Put.TimeInForce),
		FeeAmount:            "",
		FeeCcy:               "",
		ReduceOnly:           false,
		CreateTime:           res.Data.Ctime,
		UpdateTime:           res.Data.Ftime,
		RealizedPnl:          "",
		AttachTpTriggerPrice: "",
		AttachTpOrdPrice:     "",
		AttachSlTriggerPrice: "",
		TriggerPrice:         res.Data.Trigger.Price,
		TriggerType:          g.gateConverter.FromGateSpotPriceOrderTriggerRule(res.Data.Trigger.Rule, OrderSide(g.gateConverter.FromGateOrderSide(res.Data.Put.Side))),
		TriggerConditionType: "",
		Expiration:           res.Data.Trigger.Expiration,
		OcoTpTriggerPrice:    "",
		OcoTpOrdType:         "",
		OcoTpOrdPrice:        "",
		OcoSlTriggerPrice:    "",
	}
}
func (g *GateTradeEngine) handleOrderFromFuturesOrderCancel(req *OrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestFuturesSettleOrdersOrderIdDeleteRes]) *Order {
	executedQty := decimal.NewFromInt(res.Data.Size).Sub(decimal.NewFromInt(res.Data.Left)).Abs()
	fillPrice, _ := decimal.NewFromString(res.Data.FillPrice)

	symbolInfo, err := InnerExchangeManager.GetSymbolInfo(GATE_NAME.String(), req.AccountType, res.Data.Contract)
	if err != nil {
		log.Error(err)
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
		Quantity:             decimal.NewFromInt(res.Data.Size).String(),
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
		ReduceOnly:           res.Data.ReduceOnly,
		CreateTime:           decimal.NewFromFloat(res.Data.CreateTime).Mul(gateTimeMul).IntPart(),
		UpdateTime:           updateTime,
		RealizedPnl:          "",
		AttachTpTriggerPrice: "",
		AttachTpOrdPrice:     "",
		AttachSlTriggerPrice: "",

		AttachSlOrdPrice:      "",
		IsAlgo:                req.IsAlgo,
		OrderAlgoType:         OKX_ORDER_ALGO_TYPE_CONDITIONAL, // 仅支持单向止盈止损
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
	executedQty := decimal.NewFromInt(res.Data.Initial.Size).Sub(decimal.NewFromInt(res.Data.Initial.Size))
	fillPrice, _ := decimal.NewFromString(res.Data.Initial.Price)

	symbolInfo, err := InnerExchangeManager.GetSymbolInfo(GATE_NAME.String(), req.AccountType, res.Data.Initial.Contract)
	if err != nil {
		log.Error(err)
	}
	cumQuoteQty := executedQty.Mul(symbolInfo.ContractSize()).Mul(fillPrice)

	return &Order{
		Exchange:             g.ExchangeType().String(),
		AccountType:          req.AccountType,
		Symbol:               res.Data.Initial.Contract,
		IsMargin:             req.IsMargin,
		IsIsolated:           req.IsIsolated,
		OrderId:              strconv.FormatInt(res.Data.Id, 10),
		ClientOrderId:        res.Data.Initial.Text,
		Price:                res.Data.Initial.Price,
		Quantity:             strconv.FormatInt(res.Data.Initial.Size, 10),
		ExecutedQty:          executedQty.String(),
		CumQuoteQty:          cumQuoteQty.String(),
		AvgPrice:             res.Data.Initial.Price,
		Status:               g.gateConverter.FromGateContractOrderStatus(res.Data.Status, res.Data.FinishAs),
		Type:                 ORDER_TYPE_LIMIT,
		Side:                 ORDER_SIDE_BUY,
		PositionSide:         "",
		TimeInForce:          g.gateConverter.FromGateTimeInForce(res.Data.Initial.Tif),
		FeeAmount:            "",
		FeeCcy:               "",
		ReduceOnly:           false,
		CreateTime:           decimal.NewFromFloat(res.Data.CreateTime).Mul(gateTimeMul).IntPart(),
		UpdateTime:           decimal.NewFromFloat(res.Data.FinishTime).Mul(gateTimeMul).IntPart(),
		RealizedPnl:          "",
		AttachTpTriggerPrice: "",
		AttachTpOrdPrice:     "",
		AttachSlTriggerPrice: "",
		TriggerPrice:         res.Data.Trigger.Price,
		TriggerType:          g.gateConverter.FromGateFuturesPriceOrderTriggerRule(res.Data.Trigger.Rule, ORDER_SIDE_BUY),
		TriggerConditionType: "",
		OcoTpTriggerPrice:    "",
		OcoTpOrdType:         "",
		OcoTpOrdPrice:        "",
		OcoSlTriggerPrice:    "",
		OcoSlOrdType:         "",
		OcoSlOrdPrice:        "",
	}
}
func (g *GateTradeEngine) handleOrderFromDeliveryOrderCancel(req *OrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestDeliverySettleOrdersOrderIdDeleteRes]) *Order {
	executedQty := decimal.NewFromInt(res.Data.Size).Sub(decimal.NewFromInt(res.Data.Left)).Abs()
	fillPrice, _ := decimal.NewFromString(res.Data.FillPrice)

	symbolInfo, err := InnerExchangeManager.GetSymbolInfo(GATE_NAME.String(), req.AccountType, res.Data.Contract)
	if err != nil {
		log.Error(err)
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

	return &Order{
		Exchange:             g.ExchangeType().String(),
		AccountType:          req.AccountType,
		Symbol:               res.Data.Contract,
		IsMargin:             req.IsMargin,
		IsIsolated:           req.IsIsolated,
		OrderId:              decimal.NewFromInt(res.Data.Id).String(),
		ClientOrderId:        res.Data.Text,
		Price:                res.Data.Price,
		Quantity:             decimal.NewFromInt(res.Data.Size).String(),
		ExecutedQty:          executedQty.Abs().String(),
		CumQuoteQty:          cumQuoteQty.Abs().String(),
		AvgPrice:             res.Data.FillPrice,
		Status:               g.gateConverter.FromGateContractOrderStatus(res.Data.Status, res.Data.FinishAs),
		Type:                 orderType,
		Side:                 orderSide,
		PositionSide:         "",
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
	executedQty := decimal.NewFromInt(res.Data.Initial.Size).Sub(decimal.NewFromInt(res.Data.Initial.Size))
	fillPrice, _ := decimal.NewFromString(res.Data.Initial.Price)

	symbolInfo, err := InnerExchangeManager.GetSymbolInfo(GATE_NAME.String(), req.AccountType, res.Data.Initial.Contract)
	if err != nil {
		log.Error(err)
	}
	cumQuoteQty := executedQty.Mul(symbolInfo.ContractSize()).Mul(fillPrice)

	return &Order{
		Exchange:             g.ExchangeType().String(),
		AccountType:          req.AccountType,
		Symbol:               res.Data.Initial.Contract,
		IsMargin:             req.IsMargin,
		IsIsolated:           req.IsIsolated,
		OrderId:              strconv.FormatInt(res.Data.Id, 10),
		ClientOrderId:        res.Data.Initial.Text,
		Price:                res.Data.Initial.Price,
		Quantity:             strconv.FormatInt(res.Data.Initial.Size, 10),
		ExecutedQty:          executedQty.String(),
		CumQuoteQty:          cumQuoteQty.String(),
		AvgPrice:             res.Data.Initial.Price,
		Status:               g.gateConverter.FromGateContractOrderStatus(res.Data.Status, res.Data.FinishAs),
		Type:                 ORDER_TYPE_LIMIT,
		Side:                 ORDER_SIDE_BUY,
		PositionSide:         "",
		TimeInForce:          g.gateConverter.FromGateTimeInForce(res.Data.Initial.Tif),
		FeeAmount:            "",
		FeeCcy:               "",
		ReduceOnly:           false,
		CreateTime:           decimal.NewFromFloat(res.Data.CreateTime).Mul(gateTimeMul).IntPart(),
		UpdateTime:           decimal.NewFromFloat(res.Data.FinishTime).Mul(gateTimeMul).IntPart(),
		RealizedPnl:          "",
		AttachTpTriggerPrice: "",
		AttachTpOrdPrice:     "",
		AttachSlTriggerPrice: "",
		TriggerPrice:         res.Data.Trigger.Price,
		TriggerType:          g.gateConverter.FromGateFuturesPriceOrderTriggerRule(res.Data.Trigger.Rule, ORDER_SIDE_BUY),
		TriggerConditionType: "",
		OcoTpTriggerPrice:    "",
		OcoTpOrdType:         "",
		OcoTpOrdPrice:        "",
		OcoSlTriggerPrice:    "",
		OcoSlOrdType:         "",
		OcoSlOrdPrice:        "",
	}
}

func (g *GateTradeEngine) handleOrdersFromSpotOpenOrders(req *QueryOrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestSpotOpenOrdersRes]) []*Order {
	var orders []*Order

	for _, symbol := range res.Data {
		for _, order := range symbol.Orders {
			accountType, isMargin, isIsolated := g.gateConverter.FromOrderSpotAccountType(GateAccountType(order.Account))
			orders = append(orders, &Order{
				Exchange:      g.ExchangeType().String(),
				AccountType:   accountType.String(),
				Symbol:        order.CurrencyPair,
				IsMargin:      isMargin,
				IsIsolated:    isIsolated,
				OrderId:       order.ID,
				ClientOrderId: order.Text,
				Price:         order.Price,
				Quantity:      order.Amount,
				ExecutedQty:   order.FilledAmount,
				CumQuoteQty:   order.FilledTotal,
				AvgPrice:      order.AvgDealPrice,
				Status:        g.gateConverter.FromGateSpotOrderStatus(order.Status),
				Type:          g.gateConverter.FromGateOrderType(order.Type),
				Side:          g.gateConverter.FromGateOrderSide(order.Side),
				TimeInForce:   g.gateConverter.FromGateTimeInForce(order.TimeInForce),
				FeeAmount:     order.Fee,

				FeeCcy:     order.FeeCurrency,
				CreateTime: stringToInt64(order.CreateTime),
				UpdateTime: stringToInt64(order.UpdateTime),
			})
		}
	}
	return orders
}
func (g *GateTradeEngine) handleOrdersFromSpotPriceOpenOrders(req *QueryOrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestSpotPriceOrdersGetRes]) []*Order {
	var orders []*Order
	for _, order := range res.Data {
		accountType, isMargin, isIsolated := g.gateConverter.FromOrderSpotPriceAccountType(order.Put.Account)
		orders = append(orders, &Order{
			Exchange:             g.ExchangeType().String(),
			AccountType:          accountType.String(),
			Symbol:               order.Market,
			IsMargin:             isMargin,
			IsIsolated:           isIsolated,
			OrderId:              strconv.FormatInt(order.ID, 10),
			ClientOrderId:        "",
			Price:                order.Put.Price,
			Quantity:             order.Put.Amount,
			ExecutedQty:          "",
			CumQuoteQty:          "",
			AvgPrice:             order.Put.Price,
			Status:               g.gateConverter.FromGateSpotPriceOrderStatus(order.Status),
			Type:                 OrderType(g.gateConverter.FromGateOrderType(order.Put.Type)),
			Side:                 OrderSide(g.gateConverter.FromGateOrderSide(order.Put.Side)),
			PositionSide:         "",
			TimeInForce:          g.gateConverter.FromGateTimeInForce(order.Put.TimeInForce),
			FeeAmount:            "",
			FeeCcy:               "",
			ReduceOnly:           false,
			CreateTime:           order.Ctime,
			UpdateTime:           order.Ftime,
			RealizedPnl:          "",
			AttachTpTriggerPrice: "",
			AttachTpOrdPrice:     "",
			AttachSlTriggerPrice: "",
			TriggerPrice:         order.Trigger.Price,
			TriggerType:          g.gateConverter.FromGateSpotPriceOrderTriggerRule(order.Trigger.Rule, OrderSide(g.gateConverter.FromGateOrderSide(order.Put.Side))),
			TriggerConditionType: "",
			Expiration:           order.Trigger.Expiration,
			OcoTpTriggerPrice:    "",
			OcoTpOrdType:         "",
			OcoTpOrdPrice:        "",
			OcoSlTriggerPrice:    "",
		})
	}
	return orders
}
func (g *GateTradeEngine) handleOrdersFromFuturesOpenOrders(req *QueryOrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestFuturesSettleOrdersGetRes]) []*Order {
	var orders []*Order
	for _, order := range res.Data {
		executedQty := decimal.NewFromInt(order.Size).Sub(decimal.NewFromInt(order.Left)).Abs()
		fillPrice, _ := decimal.NewFromString(order.FillPrice)
		symbolInfo, err := InnerExchangeManager.GetSymbolInfo(GATE_NAME.String(), req.AccountType, order.Contract)
		if err != nil {
			log.Error(err)
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

		orders = append(orders, &Order{
			Exchange:             g.ExchangeType().String(),
			AccountType:          req.AccountType,
			Symbol:               order.Contract,
			IsMargin:             req.IsMargin,
			IsIsolated:           req.IsIsolated,
			OrderId:              strconv.FormatInt(order.Id, 10),
			ClientOrderId:        order.Text,
			Price:                order.Price,
			Quantity:             strconv.FormatInt(order.Size, 10),
			ExecutedQty:          executedQty.String(),
			CumQuoteQty:          cumQuoteQty.String(),
			AvgPrice:             order.FillPrice,
			Status:               g.gateConverter.FromGateContractOrderStatus(order.Status, order.FinishAs),
			Type:                 orderType,
			Side:                 orderSide,
			PositionSide:         "",
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
		executedQty := decimal.NewFromInt(order.Initial.Size).Sub(decimal.NewFromInt(order.Initial.Size))
		fillPrice, _ := decimal.NewFromString(order.Initial.Price)

		symbolInfo, err := InnerExchangeManager.GetSymbolInfo(GATE_NAME.String(), req.AccountType, order.Initial.Contract)
		if err != nil {
			log.Error(err)
		}
		cumQuoteQty := executedQty.Mul(symbolInfo.ContractSize()).Mul(fillPrice)

		orders = append(orders, &Order{
			Exchange:             g.ExchangeType().String(),
			AccountType:          req.AccountType,
			Symbol:               order.Initial.Contract,
			IsMargin:             req.IsMargin,
			IsIsolated:           req.IsIsolated,
			OrderId:              strconv.FormatInt(order.Id, 10),
			ClientOrderId:        order.Initial.Text,
			Price:                order.Initial.Price,
			Quantity:             strconv.FormatInt(order.Initial.Size, 10),
			ExecutedQty:          executedQty.String(),
			CumQuoteQty:          cumQuoteQty.String(),
			AvgPrice:             order.Initial.Price,
			Status:               g.gateConverter.FromGateContractOrderStatus(order.Status, order.FinishAs),
			Type:                 ORDER_TYPE_LIMIT,
			Side:                 ORDER_SIDE_BUY,
			PositionSide:         "",
			TimeInForce:          g.gateConverter.FromGateTimeInForce(order.Initial.Tif),
			FeeAmount:            "",
			FeeCcy:               "",
			ReduceOnly:           false,
			CreateTime:           decimal.NewFromFloat(order.CreateTime).Mul(gateTimeMul).IntPart(),
			UpdateTime:           decimal.NewFromFloat(order.FinishTime).Mul(gateTimeMul).IntPart(),
			RealizedPnl:          "",
			AttachTpTriggerPrice: "",
			AttachTpOrdPrice:     "",
			AttachSlTriggerPrice: "",
			TriggerPrice:         order.Trigger.Price,
			TriggerType:          g.gateConverter.FromGateFuturesPriceOrderTriggerRule(order.Trigger.Rule, ORDER_SIDE_BUY),
			TriggerConditionType: "",
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
func (g *GateTradeEngine) handleOrdersFromDeliveryOpenOrders(req *QueryOrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestDeliverySettleOrdersGetRes]) []*Order {
	var orders []*Order
	for _, order := range res.Data {
		executedQty := decimal.NewFromInt(order.Size).Sub(decimal.NewFromInt(order.Left)).Abs()
		fillPrice, _ := decimal.NewFromString(order.FillPrice)
		symbolInfo, err := InnerExchangeManager.GetSymbolInfo(GATE_NAME.String(), req.AccountType, order.Contract)
		if err != nil {
			log.Error(err)
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

		orders = append(orders, &Order{
			Exchange:             g.ExchangeType().String(),
			AccountType:          req.AccountType,
			Symbol:               order.Contract,
			IsMargin:             req.IsMargin,
			IsIsolated:           req.IsIsolated,
			OrderId:              strconv.FormatInt(order.Id, 10),
			ClientOrderId:        order.Text,
			Price:                order.Price,
			Quantity:             strconv.FormatInt(order.Size, 10),
			ExecutedQty:          executedQty.String(),
			CumQuoteQty:          cumQuoteQty.String(),
			AvgPrice:             order.FillPrice,
			Status:               g.gateConverter.FromGateContractOrderStatus(order.Status, order.FinishAs),
			Type:                 orderType,
			Side:                 orderSide,
			PositionSide:         "",
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
		executedQty := decimal.NewFromInt(order.Initial.Size).Sub(decimal.NewFromInt(order.Initial.Size))
		fillPrice, _ := decimal.NewFromString(order.Initial.Price)

		symbolInfo, err := InnerExchangeManager.GetSymbolInfo(GATE_NAME.String(), req.AccountType, order.Initial.Contract)
		if err != nil {
			log.Error(err)
		}
		cumQuoteQty := executedQty.Mul(symbolInfo.ContractSize()).Mul(fillPrice)

		orders = append(orders, &Order{
			Exchange:             g.ExchangeType().String(),
			AccountType:          req.AccountType,
			Symbol:               order.Initial.Contract,
			IsMargin:             req.IsMargin,
			IsIsolated:           req.IsIsolated,
			OrderId:              strconv.FormatInt(order.Id, 10),
			ClientOrderId:        order.Initial.Text,
			Price:                order.Initial.Price,
			Quantity:             strconv.FormatInt(order.Initial.Size, 10),
			ExecutedQty:          executedQty.String(),
			CumQuoteQty:          cumQuoteQty.String(),
			AvgPrice:             order.Initial.Price,
			Status:               g.gateConverter.FromGateContractOrderStatus(order.Status, order.FinishAs),
			Type:                 ORDER_TYPE_LIMIT,
			Side:                 ORDER_SIDE_BUY,
			PositionSide:         "",
			TimeInForce:          g.gateConverter.FromGateTimeInForce(order.Initial.Tif),
			FeeAmount:            "",
			FeeCcy:               "",
			ReduceOnly:           false,
			CreateTime:           decimal.NewFromFloat(order.CreateTime).Mul(gateTimeMul).IntPart(),
			UpdateTime:           decimal.NewFromFloat(order.FinishTime).Mul(gateTimeMul).IntPart(),
			RealizedPnl:          "",
			AttachTpTriggerPrice: "",
			AttachTpOrdPrice:     "",
			AttachSlTriggerPrice: "",
			TriggerPrice:         order.Trigger.Price,
			TriggerType:          g.gateConverter.FromGateFuturesPriceOrderTriggerRule(order.Trigger.Rule, ORDER_SIDE_BUY),
			TriggerConditionType: "",
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

func (g *GateTradeEngine) handleOrderFromSpotOrderQuery(req *QueryOrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestSpotOrdersOrderIdGetRes]) *Order {
	accountType, isMargin, isIsolated := g.gateConverter.FromOrderSpotAccountType(GateAccountType(res.Data.Account))
	return &Order{
		Exchange:      g.ExchangeType().String(),
		AccountType:   accountType.String(),
		Symbol:        res.Data.CurrencyPair,
		IsMargin:      isMargin,
		IsIsolated:    isIsolated,
		OrderId:       req.OrderId,
		ClientOrderId: req.ClientOrderId,
		Price:         res.Data.Price,
		Quantity:      res.Data.Amount,
		ExecutedQty:   res.Data.FilledAmount,
		CumQuoteQty:   res.Data.FilledTotal,
		AvgPrice:      res.Data.AvgDealPrice,
		Status:        g.gateConverter.FromGateSpotOrderStatus(res.Data.Status),
		Type:          g.gateConverter.FromGateOrderType(res.Data.Type),
		Side:          g.gateConverter.FromGateOrderSide(res.Data.Side),
		PositionSide:  "",

		TimeInForce:           g.gateConverter.FromGateTimeInForce(res.Data.TimeInForce),
		FeeAmount:             res.Data.Fee,
		FeeCcy:                res.Data.FeeCurrency,
		ReduceOnly:            false,
		CreateTime:            stringToInt64(res.Data.CreateTime),
		UpdateTime:            stringToInt64(res.Data.UpdateTime),
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
	return &Order{
		Exchange:             g.ExchangeType().String(),
		AccountType:          accountType.String(),
		Symbol:               res.Data.Market,
		IsMargin:             isMargin,
		IsIsolated:           isIsolated,
		OrderId:              strconv.FormatInt(res.Data.ID, 10),
		ClientOrderId:        "",
		Price:                res.Data.Put.Price,
		Quantity:             res.Data.Put.Amount,
		ExecutedQty:          "",
		CumQuoteQty:          "",
		AvgPrice:             res.Data.Put.Price,
		Status:               g.gateConverter.FromGateSpotPriceOrderStatus(res.Data.Status),
		Type:                 OrderType(g.gateConverter.FromGateOrderType(res.Data.Put.Type)),
		Side:                 OrderSide(g.gateConverter.FromGateOrderSide(res.Data.Put.Side)),
		PositionSide:         "",
		TimeInForce:          g.gateConverter.FromGateTimeInForce(res.Data.Put.TimeInForce),
		FeeAmount:            "",
		FeeCcy:               "",
		ReduceOnly:           false,
		CreateTime:           res.Data.Ctime,
		UpdateTime:           res.Data.Ftime,
		RealizedPnl:          "",
		AttachTpTriggerPrice: "",
		AttachTpOrdPrice:     "",
		AttachSlTriggerPrice: "",
		TriggerPrice:         res.Data.Trigger.Price,
		TriggerType:          g.gateConverter.FromGateSpotPriceOrderTriggerRule(res.Data.Trigger.Rule, OrderSide(g.gateConverter.FromGateOrderSide(res.Data.Put.Side))),
		TriggerConditionType: "",
		Expiration:           res.Data.Trigger.Expiration,
		OcoTpTriggerPrice:    "",
		OcoTpOrdType:         "",
		OcoTpOrdPrice:        "",
		OcoSlTriggerPrice:    "",
	}
}
func (g *GateTradeEngine) handleOrderFromFuturesOrderQuery(req *QueryOrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestFuturesSettleOrdersOrderIdGetRes]) *Order {
	executedQty := decimal.NewFromInt(res.Data.Size).Sub(decimal.NewFromInt(res.Data.Left)).Abs()
	fillPrice, _ := decimal.NewFromString(res.Data.FillPrice)
	symbolInfo, err := InnerExchangeManager.GetSymbolInfo(GATE_NAME.String(), req.AccountType, res.Data.Contract)
	if err != nil {
		log.Error(err)
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

	return &Order{
		Exchange:      g.ExchangeType().String(),
		AccountType:   req.AccountType,
		Symbol:        res.Data.Contract,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
		OrderId:       strconv.FormatInt(res.Data.Id, 10),
		ClientOrderId: res.Data.Text,
		Price:         res.Data.Price,
		Quantity:      decimal.NewFromInt(res.Data.Size).String(),

		ExecutedQty:          executedQty.String(),
		CumQuoteQty:          cumQuoteQty.String(),
		AvgPrice:             res.Data.FillPrice,
		Status:               g.gateConverter.FromGateContractOrderStatus(res.Data.Status, res.Data.FinishAs),
		Type:                 orderType,
		Side:                 orderSide,
		PositionSide:         "",
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
	executedQty := decimal.NewFromInt(res.Data.Initial.Size).Sub(decimal.NewFromInt(res.Data.Initial.Size))
	fillPrice, _ := decimal.NewFromString(res.Data.Initial.Price)

	symbolInfo, err := InnerExchangeManager.GetSymbolInfo(GATE_NAME.String(), req.AccountType, res.Data.Initial.Contract)
	if err != nil {
		log.Error(err)
	}
	cumQuoteQty := executedQty.Mul(symbolInfo.ContractSize()).Mul(fillPrice)

	return &Order{
		Exchange:             g.ExchangeType().String(),
		AccountType:          req.AccountType,
		Symbol:               res.Data.Initial.Contract,
		IsMargin:             req.IsMargin,
		IsIsolated:           req.IsIsolated,
		OrderId:              strconv.FormatInt(res.Data.Id, 10),
		ClientOrderId:        res.Data.Initial.Text,
		Price:                res.Data.Initial.Price,
		Quantity:             strconv.FormatInt(res.Data.Initial.Size, 10),
		ExecutedQty:          executedQty.String(),
		CumQuoteQty:          cumQuoteQty.String(),
		AvgPrice:             res.Data.Initial.Price,
		Status:               g.gateConverter.FromGateContractOrderStatus(res.Data.Status, res.Data.FinishAs),
		Type:                 ORDER_TYPE_LIMIT,
		Side:                 ORDER_SIDE_BUY,
		PositionSide:         "",
		TimeInForce:          g.gateConverter.FromGateTimeInForce(res.Data.Initial.Tif),
		FeeAmount:            "",
		FeeCcy:               "",
		ReduceOnly:           false,
		CreateTime:           decimal.NewFromFloat(res.Data.CreateTime).Mul(gateTimeMul).IntPart(),
		UpdateTime:           decimal.NewFromFloat(res.Data.FinishTime).Mul(gateTimeMul).IntPart(),
		RealizedPnl:          "",
		AttachTpTriggerPrice: "",
		AttachTpOrdPrice:     "",
		AttachSlTriggerPrice: "",
		TriggerPrice:         res.Data.Trigger.Price,
		TriggerType:          g.gateConverter.FromGateFuturesPriceOrderTriggerRule(res.Data.Trigger.Rule, ORDER_SIDE_BUY),
		TriggerConditionType: "",
		OcoTpTriggerPrice:    "",
		OcoTpOrdType:         "",
		OcoTpOrdPrice:        "",
		OcoSlTriggerPrice:    "",
		OcoSlOrdType:         "",
		OcoSlOrdPrice:        "",
	}
}
func (g *GateTradeEngine) handleOrderFromDeliveryOrderQuery(req *QueryOrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestDeliverySettleOrdersOrderIdGetRes]) *Order {
	executedQty := decimal.NewFromInt(res.Data.Size).Sub(decimal.NewFromInt(res.Data.Left)).Abs()
	fillPrice, _ := decimal.NewFromString(res.Data.FillPrice)
	symbolInfo, err := InnerExchangeManager.GetSymbolInfo(GATE_NAME.String(), req.AccountType, res.Data.Contract)
	if err != nil {
		log.Error(err)
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

	return &Order{
		Exchange:      g.ExchangeType().String(),
		AccountType:   req.AccountType,
		Symbol:        res.Data.Contract,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
		OrderId:       strconv.FormatInt(res.Data.Id, 10),
		ClientOrderId: res.Data.Text,
		Price:         res.Data.Price,
		Quantity:      decimal.NewFromInt(res.Data.Size).String(),

		ExecutedQty:          executedQty.String(),
		CumQuoteQty:          cumQuoteQty.String(),
		AvgPrice:             res.Data.FillPrice,
		Status:               g.gateConverter.FromGateContractOrderStatus(res.Data.Status, res.Data.FinishAs),
		Type:                 orderType,
		Side:                 orderSide,
		PositionSide:         "",
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
	executedQty := decimal.NewFromInt(res.Data.Initial.Size).Sub(decimal.NewFromInt(res.Data.Initial.Size))
	fillPrice, _ := decimal.NewFromString(res.Data.Initial.Price)

	symbolInfo, err := InnerExchangeManager.GetSymbolInfo(GATE_NAME.String(), req.AccountType, res.Data.Initial.Contract)
	if err != nil {
		log.Error(err)
	}
	cumQuoteQty := executedQty.Mul(symbolInfo.ContractSize()).Mul(fillPrice)

	return &Order{
		Exchange:             g.ExchangeType().String(),
		AccountType:          req.AccountType,
		Symbol:               res.Data.Initial.Contract,
		IsMargin:             req.IsMargin,
		IsIsolated:           req.IsIsolated,
		OrderId:              strconv.FormatInt(res.Data.Id, 10),
		ClientOrderId:        res.Data.Initial.Text,
		Price:                res.Data.Initial.Price,
		Quantity:             strconv.FormatInt(res.Data.Initial.Size, 10),
		ExecutedQty:          executedQty.String(),
		CumQuoteQty:          cumQuoteQty.String(),
		AvgPrice:             res.Data.Initial.Price,
		Status:               g.gateConverter.FromGateContractOrderStatus(res.Data.Status, res.Data.FinishAs),
		Type:                 ORDER_TYPE_LIMIT,
		Side:                 ORDER_SIDE_BUY,
		PositionSide:         "",
		TimeInForce:          g.gateConverter.FromGateTimeInForce(res.Data.Initial.Tif),
		FeeAmount:            "",
		FeeCcy:               "",
		ReduceOnly:           false,
		CreateTime:           decimal.NewFromFloat(res.Data.CreateTime).Mul(gateTimeMul).IntPart(),
		UpdateTime:           decimal.NewFromFloat(res.Data.FinishTime).Mul(gateTimeMul).IntPart(),
		RealizedPnl:          "",
		AttachTpTriggerPrice: "",
		AttachTpOrdPrice:     "",
		AttachSlTriggerPrice: "",
		TriggerPrice:         res.Data.Trigger.Price,
		TriggerType:          g.gateConverter.FromGateFuturesPriceOrderTriggerRule(res.Data.Trigger.Rule, ORDER_SIDE_BUY),
		TriggerConditionType: "",
		OcoTpTriggerPrice:    "",
		OcoTpOrdType:         "",
		OcoTpOrdPrice:        "",
		OcoSlTriggerPrice:    "",
		OcoSlOrdType:         "",
		OcoSlOrdPrice:        "",
	}
}

func (g *GateTradeEngine) handleOrdersFromSpotOrdersQuery(req *QueryOrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestSpotOrdersGetRes]) []*Order {
	var orders []*Order
	for _, order := range res.Data {
		accountType, isMargin, isIsolated := g.gateConverter.FromOrderSpotAccountType(GateAccountType(order.Account))
		orders = append(orders, &Order{
			Exchange:      g.ExchangeType().String(),
			AccountType:   accountType.String(),
			Symbol:        order.CurrencyPair,
			IsMargin:      isMargin,
			IsIsolated:    isIsolated,
			OrderId:       order.ID,
			ClientOrderId: order.Text,
			Price:         order.Price,
			Quantity:      order.Amount,
			ExecutedQty:   order.FilledAmount,
			CumQuoteQty:   order.FilledTotal,
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
		accountType, isMargin, isIsolated := g.gateConverter.FromOrderSpotAccountType(GateAccountType(order.Put.Account))
		orders = append(orders, &Order{
			Exchange:              g.ExchangeType().String(),
			AccountType:           accountType.String(),
			Symbol:                order.Market,
			IsMargin:              isMargin,
			IsIsolated:            isIsolated,
			OrderId:               strconv.FormatInt(order.ID, 10),
			ClientOrderId:         "",
			Price:                 order.Put.Price,
			Quantity:              order.Put.Amount,
			ExecutedQty:           "",
			CumQuoteQty:           "",
			AvgPrice:              order.Put.Price,
			Status:                g.gateConverter.FromGateSpotPriceOrderStatus(order.Status),
			Type:                  OrderType(g.gateConverter.FromGateOrderType(order.Put.Type)),
			Side:                  OrderSide(g.gateConverter.FromGateOrderSide(order.Put.Side)),
			PositionSide:          "",
			TimeInForce:           g.gateConverter.FromGateTimeInForce(order.Put.TimeInForce),
			FeeAmount:             "",
			FeeCcy:                "",
			ReduceOnly:            false,
			CreateTime:            order.Ctime,
			UpdateTime:            order.Ftime,
			RealizedPnl:           "",
			AttachTpTriggerPrice:  "",
			AttachTpOrdPrice:      "",
			AttachSlTriggerPrice:  "",
			TriggerPrice:          order.Trigger.Price,
			TriggerType:           g.gateConverter.FromGateSpotPriceOrderTriggerRule(order.Trigger.Rule, OrderSide(g.gateConverter.FromGateOrderSide(order.Put.Side))),
			TriggerConditionType:  "",
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
func (g *GateTradeEngine) handleOrdersFromFuturesOrdersQuery(req *QueryOrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestFuturesSettleOrdersGetRes]) []*Order {
	var orders []*Order
	for _, order := range res.Data {
		executedQty := decimal.NewFromInt(order.Size).Sub(decimal.NewFromInt(order.Left)).Abs()
		fillPrice, _ := decimal.NewFromString(order.FillPrice)
		cumQuoteQty := decimal.Zero
		symbolInfo, err := InnerExchangeManager.GetSymbolInfo(GATE_NAME.String(), req.AccountType, order.Contract)
		if err != nil {
			log.Error(err)
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

		orders = append(orders, &Order{
			Exchange:              g.ExchangeType().String(),
			AccountType:           req.AccountType,
			Symbol:                order.Contract,
			IsMargin:              req.IsMargin,
			IsIsolated:            req.IsIsolated,
			OrderId:               strconv.FormatInt(order.Id, 10),
			ClientOrderId:         order.Text,
			Price:                 order.Price,
			Quantity:              decimal.NewFromInt(order.Size).Abs().String(),
			ExecutedQty:           executedQty.String(),
			CumQuoteQty:           cumQuoteQty.String(),
			AvgPrice:              order.FillPrice,
			Status:                g.gateConverter.FromGateContractOrderStatus(order.Status, order.FinishAs),
			Type:                  orderType,
			Side:                  orderSide,
			PositionSide:          "",
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
		executedQty := decimal.NewFromInt(order.Initial.Size).Sub(decimal.NewFromInt(order.Initial.Size))
		fillPrice, _ := decimal.NewFromString(order.Initial.Price)

		symbolInfo, err := InnerExchangeManager.GetSymbolInfo(GATE_NAME.String(), req.AccountType, order.Initial.Contract)
		if err != nil {
			log.Error(err)
		}
		cumQuoteQty := executedQty.Mul(symbolInfo.ContractSize()).Mul(fillPrice)

		orders = append(orders, &Order{
			Exchange:             g.ExchangeType().String(),
			AccountType:          req.AccountType,
			Symbol:               order.Initial.Contract,
			IsMargin:             req.IsMargin,
			IsIsolated:           req.IsIsolated,
			OrderId:              strconv.FormatInt(order.Id, 10),
			ClientOrderId:        order.Initial.Text,
			Price:                order.Initial.Price,
			Quantity:             strconv.FormatInt(order.Initial.Size, 10),
			ExecutedQty:          executedQty.String(),
			CumQuoteQty:          cumQuoteQty.String(),
			AvgPrice:             order.Initial.Price,
			Status:               g.gateConverter.FromGateContractOrderStatus(order.Status, order.FinishAs),
			Type:                 ORDER_TYPE_LIMIT,
			Side:                 ORDER_SIDE_BUY,
			PositionSide:         "",
			TimeInForce:          g.gateConverter.FromGateTimeInForce(order.Initial.Tif),
			FeeAmount:            "",
			FeeCcy:               "",
			ReduceOnly:           false,
			CreateTime:           decimal.NewFromFloat(order.CreateTime).Mul(gateTimeMul).IntPart(),
			UpdateTime:           decimal.NewFromFloat(order.FinishTime).Mul(gateTimeMul).IntPart(),
			RealizedPnl:          "",
			AttachTpTriggerPrice: "",
			AttachTpOrdPrice:     "",
			AttachSlTriggerPrice: "",
			TriggerPrice:         order.Trigger.Price,
			TriggerType:          g.gateConverter.FromGateFuturesPriceOrderTriggerRule(order.Trigger.Rule, ORDER_SIDE_BUY),
			TriggerConditionType: "",
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
func (g *GateTradeEngine) handleOrdersFromDeliveryOrdersQuery(req *QueryOrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestDeliverySettleOrdersGetRes]) []*Order {
	var orders []*Order
	for _, order := range res.Data {
		executedQty := decimal.NewFromInt(order.Size).Sub(decimal.NewFromInt(order.Left))
		fillPrice, _ := decimal.NewFromString(order.FillPrice)
		symbolInfo, err := InnerExchangeManager.GetSymbolInfo(GATE_NAME.String(), req.AccountType, order.Contract)
		if err != nil {
			log.Error(err)
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

		orders = append(orders, &Order{
			Exchange:             g.ExchangeType().String(),
			AccountType:          req.AccountType,
			Symbol:               order.Contract,
			IsMargin:             req.IsMargin,
			IsIsolated:           req.IsIsolated,
			OrderId:              strconv.FormatInt(order.Id, 10),
			ClientOrderId:        order.Text,
			Price:                order.Price,
			Quantity:             decimal.NewFromInt(order.Size).Abs().String(),
			ExecutedQty:          executedQty.Abs().String(),
			CumQuoteQty:          cumQuoteQty.Abs().String(),
			AvgPrice:             order.Price,
			Status:               g.gateConverter.FromGateContractOrderStatus(order.Status, order.FinishAs),
			Type:                 orderType,
			Side:                 orderSide,
			PositionSide:         "",
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
		executedQty := decimal.NewFromInt(order.Initial.Size).Sub(decimal.NewFromInt(order.Initial.Size))
		fillPrice, _ := decimal.NewFromString(order.Initial.Price)

		symbolInfo, err := InnerExchangeManager.GetSymbolInfo(GATE_NAME.String(), req.AccountType, order.Initial.Contract)
		if err != nil {
			log.Error(err)
		}
		cumQuoteQty := executedQty.Mul(symbolInfo.ContractSize()).Mul(fillPrice)

		orders = append(orders, &Order{
			Exchange:             g.ExchangeType().String(),
			AccountType:          req.AccountType,
			Symbol:               order.Initial.Contract,
			IsMargin:             req.IsMargin,
			IsIsolated:           req.IsIsolated,
			OrderId:              strconv.FormatInt(order.Id, 10),
			ClientOrderId:        order.Initial.Text,
			Price:                order.Initial.Price,
			Quantity:             strconv.FormatInt(order.Initial.Size, 10),
			ExecutedQty:          executedQty.String(),
			CumQuoteQty:          cumQuoteQty.String(),
			AvgPrice:             order.Initial.Price,
			Status:               g.gateConverter.FromGateContractOrderStatus(order.Status, order.FinishAs),
			Type:                 ORDER_TYPE_LIMIT,
			Side:                 ORDER_SIDE_BUY,
			PositionSide:         "",
			TimeInForce:          g.gateConverter.FromGateTimeInForce(order.Initial.Tif),
			FeeAmount:            "",
			FeeCcy:               "",
			ReduceOnly:           false,
			CreateTime:           decimal.NewFromFloat(order.CreateTime).Mul(gateTimeMul).IntPart(),
			UpdateTime:           decimal.NewFromFloat(order.FinishTime).Mul(gateTimeMul).IntPart(),
			RealizedPnl:          "",
			AttachTpTriggerPrice: "",
			AttachTpOrdPrice:     "",
			AttachSlTriggerPrice: "",
			TriggerPrice:         order.Trigger.Price,
			TriggerType:          g.gateConverter.FromGateFuturesPriceOrderTriggerRule(order.Trigger.Rule, ORDER_SIDE_BUY),
			TriggerConditionType: "",
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

func (g *GateTradeEngine) handleTradesFromFuturesTradesQuery(req *QueryTradeParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestFuturesSettleMyTradesRes]) []*Trade {
	var trades []*Trade
	for _, trade := range res.Data {
		price, _ := decimal.NewFromString(trade.Price)
		amt := decimal.NewFromInt(trade.Size)

		quoteQty := decimal.Zero
		symbolInfo, err := InnerExchangeManager.GetSymbolInfo(GATE_NAME.String(), req.AccountType, trade.Contract)
		if err != nil {
			log.Error(err)
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
			PositionSide:  "",
			FeeAmount:     trade.Fee,
			FeeCcy:        feeCcy,
			RealizedPnl:   "",
			IsMaker:       isMaker,
			Timestamp:     decimal.NewFromFloat(trade.CreateTime).Mul(gateTimeMul).IntPart(),
		})
	}
	return trades
}
func (g *GateTradeEngine) handleTradesFromDeliveryTradesQuery(req *QueryTradeParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestDeliverySettleMyTradesRes]) []*Trade {
	var trades []*Trade
	for _, trade := range res.Data {
		price, _ := decimal.NewFromString(trade.Price)
		amt := decimal.NewFromInt(trade.Size)

		symbolInfo, err := InnerExchangeManager.GetSymbolInfo(GATE_NAME.String(), req.AccountType, trade.Contract)
		if err != nil {
			log.Error(err)
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
			PositionSide:  "",
			FeeAmount:     trade.Fee,
			FeeCcy:        feeCcy,
			RealizedPnl:   "",
			IsMaker:       isMaker,
			Timestamp:     decimal.NewFromFloat(trade.CreateTime).Mul(gateTimeMul).IntPart(),
		})
	}
	return trades
}
