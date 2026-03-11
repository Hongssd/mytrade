package mytrade

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Hongssd/myxcoinapi"
	"github.com/shopspring/decimal"
)

func getXcoinCumQuoteQty(quoteQty, avgPrice, totalFillQty string) string {
	// 优先使用交易所返回的 QuoteQty
	if quoteQty != "" {
		return quoteQty
	}
	avg, err := decimal.NewFromString(avgPrice)
	if err != nil {
		return "0"
	}
	filled, err := decimal.NewFromString(totalFillQty)
	if err != nil {
		return "0"
	}
	return avg.Mul(filled).String()
}

func (e *XcoinTradeEngine) handleOrdersFromQueryOpenOrders(req *QueryOrderParam, res *myxcoinapi.XcoinRestRes[myxcoinapi.PrivateRestTradeOpenOrdersRes]) []*Order {
	var orders []*Order
	for _, order := range res.Data {
		orderType, timeInForce := e.xcoinConverter.FromXcoinOrderType(order.OrderType, order.TimeInForce)
		symbolSplit := strings.Split(order.Symbol, "-")
		if len(symbolSplit) < 2 {
			log.Errorf("symbol split error: %s", order.Symbol)
			continue
		}
		feeCcy := symbolSplit[1]
		orders = append(orders, &Order{
			Exchange:      e.ExchangeType().String(),
			AccountType:   req.AccountType,
			Symbol:        order.Symbol,
			IsMargin:      false,
			IsIsolated:    false, // sunx 仅全仓
			OrderId:       order.OrderId,
			ClientOrderId: order.ClientOrderId,
			Price:         order.Price,
			Quantity:      order.Qty,
			ExecutedQty:   order.TotalFillQty,
			CumQuoteQty:   getXcoinCumQuoteQty(order.QuoteQty, order.AvgPrice, order.TotalFillQty),
			AvgPrice:      order.AvgPrice,
			Status:        e.xcoinConverter.FromXcoinOrderStatus(order.Status),
			Type:          orderType,
			Side:          e.xcoinConverter.FromXcoinOrderSide(order.Side),
			PositionSide:  e.xcoinConverter.FromXcoinPositionSide(order.PosSide),
			TimeInForce:   timeInForce,
			ReduceOnly:    order.ReduceOnly,
			FeeAmount:     order.QuoteFee,
			FeeCcy:        feeCcy,
			CreateTime:    stringToInt64(order.CreateTime),
			UpdateTime:    stringToInt64(order.UpdateTime),
		})
	}
	return orders
}

func (e *XcoinTradeEngine) handleOrderFromQueryOrder(req *QueryOrderParam, res *myxcoinapi.XcoinRestRes[myxcoinapi.PrivateRestTradeOrderInfoRes]) (*Order, error) {
	r := res.Data
	symbolSplit := strings.Split(r.Symbol, "-")
	if len(symbolSplit) < 2 {
		log.Errorf("symbol split error: %s", r.Symbol)
		return nil, errors.New("symbol split error")
	}
	feeCcy := symbolSplit[1]
	orderType, timeInForce := e.xcoinConverter.FromXcoinOrderType(r.OrderType, r.TimeInForce)
	ocoTpOrdType, _ := e.xcoinConverter.FromXcoinOrderType(r.TpslOrder.TpOrderType, "")
	ocoSlOrdType, _ := e.xcoinConverter.FromXcoinOrderType(r.TpslOrder.SlOrderType, "")
	if feeCcy == "" {
		log.Errorf("fee ccy is empty: %s", r.Symbol)
		return nil, errors.New("fee ccy is empty")
	}
	return &Order{
		Exchange:          e.ExchangeType().String(),
		AccountType:       req.AccountType,
		Symbol:            r.Symbol,
		IsMargin:          false,
		IsIsolated:        false, // sunx 仅全仓
		OrderId:           r.OrderId,
		ClientOrderId:     r.ClientOrderId,
		Price:             r.Price,
		Quantity:          r.Qty,
		ExecutedQty:       r.TotalFillQty,
		CumQuoteQty:       getXcoinCumQuoteQty(r.QuoteQty, r.AvgPrice, r.TotalFillQty),
		AvgPrice:          r.AvgPrice,
		Status:            e.xcoinConverter.FromXcoinOrderStatus(r.Status),
		Type:              orderType,
		Side:              e.xcoinConverter.FromXcoinOrderSide(r.Side),
		PositionSide:      e.xcoinConverter.FromXcoinPositionSide(r.PosSide),
		TimeInForce:       timeInForce,
		ReduceOnly:        r.ReduceOnly,
		FeeAmount:         r.QuoteFee,
		FeeCcy:            feeCcy,
		CreateTime:        stringToInt64(r.CreateTime),
		UpdateTime:        stringToInt64(r.UpdateTime),
		IsAlgo:            false,
		OcoTpTriggerPrice: r.TpslOrder.TakeProfit,
		OcoTpOrdType:      ocoTpOrdType,
		OcoTpOrdPrice:     r.TpslOrder.TpLimitPrice,
		OcoSlTriggerPrice: r.TpslOrder.StopLoss,
		OcoSlOrdType:      ocoSlOrdType,
		OcoSlOrdPrice:     r.TpslOrder.SlLimitPrice,
	}, nil
}

func (e *XcoinTradeEngine) handleOrdersFromQueryOrders(req *QueryOrderParam, res *myxcoinapi.XcoinRestRes[myxcoinapi.PrivateRestTradeHistoryOrdersRes]) []*Order {
	var orders []*Order

	for _, order := range res.Data {
		orderType, timeInForce := e.xcoinConverter.FromXcoinOrderType(order.OrderType, order.TimeInForce)
		symbolSplit := strings.Split(order.Symbol, "-")
		if len(symbolSplit) < 2 {
			log.Errorf("symbol split error: %s", order.Symbol)
			continue
		}
		feeCcy := symbolSplit[1]
		orders = append(orders, &Order{
			Exchange:      e.ExchangeType().String(),
			AccountType:   req.AccountType,
			Symbol:        order.Symbol,
			IsMargin:      false,
			IsIsolated:    false, // xcoin 仅全仓
			OrderId:       order.OrderId,
			ClientOrderId: order.ClientOrderId,
			Price:         order.Price,
			Quantity:      order.Qty,
			ExecutedQty:   order.TotalFillQty,
			CumQuoteQty:   getXcoinCumQuoteQty(order.QuoteQty, order.AvgPrice, order.TotalFillQty),
			AvgPrice:      order.AvgPrice,
			Status:        e.xcoinConverter.FromXcoinOrderStatus(order.Status),
			Type:          orderType,
			Side:          e.xcoinConverter.FromXcoinOrderSide(order.Side),
			PositionSide:  e.xcoinConverter.FromXcoinPositionSide(order.PosSide),
			TimeInForce:   timeInForce,
			ReduceOnly:    order.ReduceOnly,
			FeeAmount:     order.QuoteFee,
			FeeCcy:        feeCcy,
			CreateTime:    stringToInt64(order.CreateTime),
			UpdateTime:    stringToInt64(order.UpdateTime),
		})
	}
	return orders
}

func (e *XcoinTradeEngine) handleTradesFromQueryTrades(req *QueryTradeParam, res *myxcoinapi.XcoinRestRes[myxcoinapi.PrivateRestTradeHistoryTradesRes]) []*Trade {
	var trades []*Trade

	for _, trade := range res.Data {
		quoteQty := "0"
		price, err := decimal.NewFromString(trade.FillPrice)
		if err == nil {
			qty, err := decimal.NewFromString(trade.FillQty)
			if err == nil {
				quoteQty = price.Mul(qty).String()
			}
		}
		trades = append(trades, &Trade{
			Exchange:     e.ExchangeType().String(),
			AccountType:  trade.BusinessType,
			Symbol:       trade.Symbol,
			TradeId:      trade.TradeId,
			OrderId:      trade.OrderId,
			Price:        trade.FillPrice,
			Quantity:     trade.FillQty,
			QuoteQty:     quoteQty,
			Side:         e.xcoinConverter.FromXcoinOrderSide(trade.Side),
			PositionSide: POSITION_SIDE_UNKNOWN,
			FeeAmount:    trade.Fee,
			FeeCcy:       trade.FeeCurrency,
			RealizedPnl:  trade.Pnl,
			IsMaker:      strings.EqualFold(trade.Role, "maker"),
			Timestamp:    stringToInt64(trade.FillTime),
		})
	}
	return trades
}

// 单订单返回结果处理
func (e *XcoinTradeEngine) handleOrderFromOrderCreate(req *OrderParam, res *myxcoinapi.XcoinRestRes[myxcoinapi.PrivateRestTradeOrderRes]) (*Order, error) {
	if res.Code != "0" {
		return nil, fmt.Errorf("[%s]:%s", res.Code, res.Msg)
	}

	r := res.Data
	order := &Order{
		Exchange:      e.ExchangeType().String(),
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		IsMargin:      false,
		IsIsolated:    false, // xcoin 仅全仓
		OrderId:       r.OrderId,
		ClientOrderId: r.ClientOrderId,
	}
	return order, nil
}

func (e *XcoinTradeEngine) handleOrderFromBatchOrderCreate(reqs []*OrderParam, res *myxcoinapi.XcoinRestRes[myxcoinapi.PrivateRestTradeBatchOrderRes]) ([]*Order, error) {
	if res.Code != "0" {
		return nil, fmt.Errorf("[%s]:%s", res.Code, res.Msg)
	}
	if len(res.Data) != len(reqs) {
		return nil, errors.New("api return invalid data")
	}
	orders := make([]*Order, 0, len(reqs))
	for i, r := range res.Data {
		if r.Code != "0" {
			return nil, fmt.Errorf("[%s][%s]:%s", r.ClientOrderId, r.Code, r.Msg)
		}
		order := &Order{
			Exchange:      e.ExchangeType().String(),
			AccountType:   reqs[i].AccountType,
			Symbol:        reqs[i].Symbol,
			IsMargin:      false,
			IsIsolated:    false,
			OrderId:       r.OrderId,
			ClientOrderId: r.ClientOrderId,
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (e *XcoinTradeEngine) handleOrderFromBatchOrderCancel(reqs []*OrderParam, res *myxcoinapi.XcoinRestRes[myxcoinapi.PrivateRestTradeBatchCancelOrderRes]) ([]*Order, error) {
	if res.Code != "0" {
		return nil, fmt.Errorf("[%s]:%s", res.Code, res.Msg)
	}
	if len(res.Data) != len(reqs) {
		return nil, errors.New("api return invalid data")
	}
	orders := make([]*Order, 0, len(reqs))
	for i, r := range res.Data {
		if r.Code != 0 {
			return nil, fmt.Errorf("[%s][%d]:%s", r.OrderId, r.Code, r.Msg)
		}
		order := &Order{
			Exchange:      e.ExchangeType().String(),
			AccountType:   reqs[i].AccountType,
			Symbol:        reqs[i].Symbol,
			IsMargin:      false,
			IsIsolated:    false,
			OrderId:       r.OrderId,
			ClientOrderId: reqs[i].ClientOrderId,
		}
		orders = append(orders, order)
	}
	return orders, nil
}
