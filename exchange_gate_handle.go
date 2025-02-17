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
func (g *GateTradeEngine) handleOrderFromFuturesOrderCancel(req *OrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestFuturesSettleOrdersOrderIdDeleteRes]) *Order {
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
		ReduceOnly:           false,
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
func (g *GateTradeEngine) handleOrderFromDeliveryOrderCancel(req *OrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestDeliverySettleOrdersOrderIdDeleteRes]) *Order {
	executedQty := decimal.NewFromInt(res.Data.Size).Sub(decimal.NewFromInt(res.Data.Left))
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
		ReduceOnly:           false,
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
func (g *GateTradeEngine) handleOrdersFromFuturesOpenOrders(req *QueryOrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestFuturesSettleOrdersGetRes]) []*Order {
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
			ReduceOnly:           false,
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
func (g *GateTradeEngine) handleOrdersFromDeliveryOpenOrders(req *QueryOrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestDeliverySettleOrdersGetRes]) []*Order {
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
			ReduceOnly:           false,
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
func (g *GateTradeEngine) handleOrderFromFuturesOrderQuery(req *QueryOrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestFuturesSettleOrdersOrderIdGetRes]) *Order {
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
		TimeInForce:          "",
		FeeAmount:            "",
		FeeCcy:               "",
		ReduceOnly:           false,
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
func (g *GateTradeEngine) handleOrderFromDeliveryOrderQuery(req *QueryOrderParam, res *mygateapi.GateRestRes[mygateapi.PrivateRestDeliverySettleOrdersOrderIdGetRes]) *Order {
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
		TimeInForce:          "",
		FeeAmount:            "",
		FeeCcy:               "",
		ReduceOnly:           false,
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
