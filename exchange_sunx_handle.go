package mytrade

import (
	"errors"
	"fmt"

	"github.com/Hongssd/mysunxapi"
)

func (e *SunxTradeEngine) handleOrdersFromQueryOpenOrders(req *QueryOrderParam, res *mysunxapi.SunxRestRes[mysunxapi.PrivateRestTradeOrderOpensRes]) []*Order {
	var orders []*Order

	for _, order := range res.Data {
		orders = append(orders, &Order{
			Exchange:      e.ExchangeType().String(),
			AccountType:   req.AccountType,
			Symbol:        order.ContractCode,
			IsMargin:      true,
			IsIsolated:    false, // sunx 仅全仓
			OrderId:       order.OrderId,
			ClientOrderId: order.ClientOrderId,
			Price:         order.Price,
			Quantity:      order.Volume,
			ExecutedQty:   order.TradeVolume,
			CumQuoteQty:   order.TradeTurnover,
			AvgPrice:      order.TradeAvgPrice,
			Status:        e.sunxConverter.FromSunxOrderStatus(order.State),
			Type:          e.sunxConverter.FromSunxOrderType(order.Type),
			Side:          e.sunxConverter.FromSunxOrderSide(order.Side),
			PositionSide:  e.sunxConverter.FromSunxPositionSide(order.PositionSide),
			TimeInForce:   e.sunxConverter.FromSunxTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			FeeAmount:     order.Fee,
			FeeCcy:        order.FeeCurrency,
			CreateTime:    stringToInt64(order.CreatedTime),
			UpdateTime:    stringToInt64(order.UpdatedTime),
		})
	}
	return orders
}

func (e *SunxTradeEngine) handleOrderFromQueryOrder(req *QueryOrderParam, res *mysunxapi.SunxRestRes[mysunxapi.PrivateRestTradeOrderGetRes]) (*Order, error) {
	r := res.Data
	return &Order{
		Exchange:          e.ExchangeType().String(),
		AccountType:       req.AccountType,
		Symbol:            r.ContractCode,
		IsMargin:          true,
		IsIsolated:        false, // sunx 仅全仓
		OrderId:           r.OrderId,
		Price:             r.Price,
		Quantity:          r.Volume,
		ExecutedQty:       r.TradeVolume,
		CumQuoteQty:       r.TradeTurnover,
		AvgPrice:          r.TradeAvgPrice,
		Status:            e.sunxConverter.FromSunxOrderStatus(r.State),
		Type:              e.sunxConverter.FromSunxOrderType(r.Type),
		Side:              e.sunxConverter.FromSunxOrderSide(r.Side),
		PositionSide:      e.sunxConverter.FromSunxPositionSide(r.PositionSide),
		TimeInForce:       e.sunxConverter.FromSunxTimeInForce(r.TimeInForce),
		ReduceOnly:        r.ReduceOnly,
		FeeAmount:         r.Fee,
		FeeCcy:            r.FeeCurrency,
		CreateTime:        stringToInt64(r.CreatedTime),
		UpdateTime:        stringToInt64(r.UpdatedTime),
		IsAlgo:            false,
		OcoTpTriggerPrice: r.TpTriggerPrice,
		OcoTpOrdType:      e.sunxConverter.FromSunxOrderType(r.TpType),
		OcoTpOrdPrice:     r.TpOrderPrice,
		OcoSlTriggerPrice: r.SlTriggerPrice,
		OcoSlOrdType:      e.sunxConverter.FromSunxOrderType(r.SlType),
		OcoSlOrdPrice:     r.SlOrderPrice,
	}, nil
}

func (e *SunxTradeEngine) handleOrdersFromQueryOrders(req *QueryOrderParam, res *mysunxapi.SunxRestRes[mysunxapi.PrivateRestTradeOrderHistoryRes]) []*Order {
	var orders []*Order

	for _, order := range res.Data {
		orders = append(orders, &Order{
			Exchange:      e.ExchangeType().String(),
			AccountType:   req.AccountType,
			Symbol:        order.ContractCode,
			IsMargin:      true,
			IsIsolated:    false, // sunx 仅全仓
			OrderId:       order.OrderId,
			ClientOrderId: order.ClientOrderId,
			Price:         order.Price,
			Quantity:      order.Volume,
			ExecutedQty:   order.TradeVolume,
			CumQuoteQty:   order.TradeTurnover,
			AvgPrice:      order.TradeAvgPrice,
			Status:        e.sunxConverter.FromSunxOrderStatus(order.State),
			Type:          e.sunxConverter.FromSunxOrderType(order.Type),
			Side:          e.sunxConverter.FromSunxOrderSide(order.Side),
			PositionSide:  e.sunxConverter.FromSunxPositionSide(order.PositionSide),
			TimeInForce:   e.sunxConverter.FromSunxTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			FeeAmount:     order.Fee,
			FeeCcy:        order.FeeCurrency,
			CreateTime:    stringToInt64(order.CreatedTime),
			UpdateTime:    stringToInt64(order.UpdatedTime),
		})
	}
	return orders
}

func (e *SunxTradeEngine) handleTradesFromQueryTrades(req *QueryTradeParam, res *mysunxapi.SunxRestRes[mysunxapi.PrivateRestTradeOrderDetailsRes]) []*Trade {
	var trades []*Trade

	for _, trade := range res.Data {
		trades = append(trades, &Trade{
			Exchange:     e.ExchangeType().String(),
			AccountType:  req.AccountType,
			Symbol:       trade.ContractCode,
			TradeId:      trade.TradeId,
			OrderId:      trade.OrderId,
			Price:        trade.TradePrice,
			Quantity:     trade.TradeVolume,
			QuoteQty:     trade.TradeTurnover,
			Side:         e.sunxConverter.FromSunxOrderSide(trade.Side),
			PositionSide: e.sunxConverter.FromSunxPositionSide(trade.PositionSide),
			FeeAmount:    trade.TradeFee,
			FeeCcy:       trade.FeeCurrency,
			RealizedPnl:  trade.Profit,
			IsMaker:      trade.Role == "MAKER",
			Timestamp:    stringToInt64(trade.UpdatedTime),
		})
	}
	return trades
}

func (e *SunxTradeEngine) handleOrderFromOrderCreate(req *OrderParam, res *mysunxapi.SunxRestRes[mysunxapi.PrivateRestTradeOrderResCommon]) (*Order, error) {
	if res.Message != "Success" {
		return nil, errors.New(res.Message)
	}
	order := &Order{
		Exchange:      e.ExchangeType().String(),
		OrderId:       res.Data.OrderId,
		ClientOrderId: res.Data.ClientOrderId,
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		IsMargin:      true,
		IsIsolated:    false, // sunx 仅全仓
		Status:        ORDER_STATUS_NEW,
		UpdateTime:    res.Ts,
	}
	return order, nil
}

func (e *SunxTradeEngine) handleOrdersFromBatchOrderCreate(reqs []*OrderParam, res *mysunxapi.SunxRestRes[mysunxapi.PrivateRestTradeBatchOrdersRes]) ([]*Order, error) {
	if len(res.Data) != len(reqs) {
		return nil, errors.New("api return invalid data")
	}
	orders := make([]*Order, 0, len(reqs))
	for i, r := range res.Data {
		if r.Code != 200 {
			log.Errorf("{[%s][%d]:%s}", reqs[i].ClientOrderId, r.Code, r.Message)
			continue
		}
		orders = append(orders, &Order{
			Exchange:      e.ExchangeType().String(),
			OrderId:       r.OrderId,
			ClientOrderId: reqs[i].ClientOrderId,
			AccountType:   reqs[i].AccountType,
			Symbol:        reqs[i].Symbol,
			IsMargin:      reqs[i].IsMargin,
			IsIsolated:    reqs[i].IsIsolated,
			Status:        ORDER_STATUS_NEW,
			UpdateTime:    res.Ts,
		})
	}
	return orders, nil
}

func (e *SunxTradeEngine) handleOrderFromOrderCancel(req *OrderParam, res *mysunxapi.SunxRestRes[mysunxapi.PrivateRestTradeCancelOrderRes]) (*Order, error) {
	if res.Message != "Success" {
		return nil, errors.New(res.Message)
	}
	order := &Order{
		Exchange:      e.ExchangeType().String(),
		OrderId:       res.Data.OrderId,
		ClientOrderId: res.Data.ClientOrderId,
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		IsMargin:      true,
		IsIsolated:    false, // sunx 仅全仓
		Status:        ORDER_STATUS_CANCELED,
		UpdateTime:    res.Ts,
	}
	return order, nil
}

func (e *SunxTradeEngine) handleOrdersFromBatchOrderCancel(reqs []*OrderParam, res *mysunxapi.SunxRestRes[mysunxapi.PrivateRestTradeCancelBatchOrdersRes]) ([]*Order, error) {
	if len(res.Data) != len(reqs) {
		return nil, errors.New("api return invalid data")
	}
	errStr := ""
	orders := make([]*Order, 0, len(reqs))
	for i, r := range res.Data {
		if r.Code != 200 {
			errStr += fmt.Sprintf("{[%s][%d]:%s}", reqs[i].ClientOrderId, r.Code, r.Message)
		}
		order := &Order{
			Exchange:      e.ExchangeType().String(),
			OrderId:       r.OrderId,
			ClientOrderId: reqs[i].ClientOrderId,
			AccountType:   reqs[i].AccountType,
			Symbol:        reqs[i].Symbol,
			IsMargin:      reqs[i].IsMargin,
			IsIsolated:    reqs[i].IsIsolated,
			Status:        ORDER_STATUS_CANCELED,
			UpdateTime:    res.Ts,
		}
		orders = append(orders, order)
	}
	if errStr != "" {
		return orders, fmt.Errorf("[%d]%s: [%s]", res.Code, res.Message, errStr)
	}
	return orders, nil
}
