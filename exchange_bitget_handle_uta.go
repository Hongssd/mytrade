package mytrade

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	mybitgetapi "github.com/Hongssd/mybitgetapi"
)

func (b *BitgetTradeEngine) handleUtaTradeFeeDetailsAgg(details []mybitgetapi.PrivateRestUtaTradeFeeDetail) (amount, ccy string) {
	if len(details) == 0 {
		return "0", ""
	}
	var sum float64
	ccy = details[0].FeeCoin
	for _, d := range details {
		if d.FeeCoin != "" {
			ccy = d.FeeCoin
		}
		sum += stringToFloat64(d.Fee)
	}
	return strconv.FormatFloat(sum, 'f', -1, 64), ccy
}

func (b *BitgetTradeEngine) handleOrderFromUtaTradeOrderInfo(c BitgetEnumConverter, accountType string, d *mybitgetapi.PrivateRestUtaTradeOrderInfoRes) *Order {
	if d == nil {
		return nil
	}
	ot, tif := c.FromBitgetOrderTypeWithTIF(d.OrderType, d.TimeInForce)
	feeAmt, feeCcy := b.handleUtaTradeFeeDetailsAgg(d.FeeDetail)
	return &Order{
		Exchange:      BITGET_NAME.String(),
		AccountType:   accountType,
		Symbol:        d.Symbol,
		OrderId:       d.OrderId,
		ClientOrderId: d.ClientOid,
		Price:         d.Price,
		Quantity:      d.Qty,
		ExecutedQty:   d.CumExecQty,
		CumQuoteQty:   d.CumExecValue,
		AvgPrice:      d.AvgPrice,
		Status:        c.FromBitgetOrderStatusUTA(d.OrderStatus),
		Type:          ot,
		Side:          c.FromBitgetOrderSide(d.Side),
		PositionSide:  c.FromBitgetPositionSide(d.PosSide),
		TimeInForce:   tif,
		FeeAmount:     feeAmt,
		FeeCcy:        feeCcy,
		ReduceOnly:    c.ReduceOnlyFromString(d.ReduceOnly),
		CreateTime:    stringToInt64(d.CreatedTime),
		UpdateTime:    stringToInt64(d.UpdatedTime),
	}
}

func (b *BitgetTradeEngine) handleOrderFromUtaTradeOrderListItem(c BitgetEnumConverter, accountType string, d *mybitgetapi.PrivateRestUtaTradeOrderListItem) *Order {
	if d == nil {
		return nil
	}
	ot, tif := c.FromBitgetOrderTypeWithTIF(d.OrderType, d.TimeInForce)
	feeAmt, feeCcy := b.handleUtaTradeFeeDetailsAgg(d.FeeDetail)
	return &Order{
		Exchange:      BITGET_NAME.String(),
		AccountType:   accountType,
		Symbol:        d.Symbol,
		OrderId:       d.OrderId,
		ClientOrderId: d.ClientOid,
		Price:         d.Price,
		Quantity:      d.Qty,
		ExecutedQty:   d.CumExecQty,
		CumQuoteQty:   d.CumExecValue,
		AvgPrice:      d.AvgPrice,
		Status:        c.FromBitgetOrderStatusUTA(d.OrderStatus),
		Type:          ot,
		Side:          c.FromBitgetOrderSide(d.Side),
		PositionSide:  c.FromBitgetPositionSide(d.PosSide),
		TimeInForce:   tif,
		FeeAmount:     feeAmt,
		FeeCcy:        feeCcy,
		ReduceOnly:    c.ReduceOnlyFromString(d.ReduceOnly),
		CreateTime:    stringToInt64(d.CreatedTime),
		UpdateTime:    stringToInt64(d.UpdatedTime),
	}
}

func (b *BitgetTradeEngine) handleTradesFromUtaTradeFillList(c BitgetEnumConverter, accountType string, list []mybitgetapi.PrivateRestUtaTradeFillItem) []*Trade {
	out := make([]*Trade, 0, len(list))
	for i := range list {
		d := &list[i]
		feeAmt, feeCcy := b.handleUtaTradeFeeDetailsAgg(d.FeeDetail)
		out = append(out, &Trade{
			Exchange:      BITGET_NAME.String(),
			AccountType:   accountType,
			Symbol:        d.Symbol,
			TradeId:       d.ExecId,
			OrderId:       d.OrderId,
			ClientOrderId: d.ClientOid,
			Price:         d.ExecPrice,
			Quantity:      d.ExecQty,
			QuoteQty:      d.ExecValue,
			Side:          c.FromBitgetOrderSide(d.Side),
			PositionSide:  POSITION_SIDE_BOTH,
			FeeAmount:     feeAmt,
			FeeCcy:        feeCcy,
			RealizedPnl:   d.ExecPnl,
			IsMaker:       strings.EqualFold(d.TradeScope, "maker"),
			Timestamp:     stringToInt64(d.CreatedTime),
		})
	}
	return out
}

func (b *BitgetTradeEngine) handleOrdersFromUtaTradeOrderList(c BitgetEnumConverter, accountType string, list []mybitgetapi.PrivateRestUtaTradeOrderListItem) []*Order {
	out := make([]*Order, 0, len(list))
	for i := range list {
		out = append(out, b.handleOrderFromUtaTradeOrderListItem(c, accountType, &list[i]))
	}
	return out
}

func bitgetUtaBatchRowOK(code string) bool {
	return code == "" || code == "00000"
}

func handleTradesFilterBySymbol(trades []*Trade, symbol string) []*Trade {
	if symbol == "" {
		return trades
	}
	filtered := make([]*Trade, 0, len(trades))
	for _, t := range trades {
		if t.Symbol == symbol {
			filtered = append(filtered, t)
		}
	}
	return filtered
}

func (b *BitgetTradeEngine) handleOrdersFromUtaTradePlaceBatch(reqs []*OrderParam, res mybitgetapi.PrivateRestUtaTradePlaceBatchRes) ([]*Order, error) {
	if len(res) != len(reqs) {
		return nil, errors.New("api return invalid data")
	}
	errStr := ""
	orders := make([]*Order, 0, len(reqs))
	for i, r := range res {
		if !bitgetUtaBatchRowOK(r.Code) {
			errStr += fmt.Sprintf("{[%s][%s]:%s}", r.ClientOid, r.Code, r.Msg)
		}
		orders = append(orders, &Order{
			Exchange:      BITGET_NAME.String(),
			OrderId:       r.OrderId,
			ClientOrderId: r.ClientOid,
			AccountType:   reqs[i].AccountType,
			Symbol:        reqs[i].Symbol,
			IsMargin:      reqs[i].IsMargin,
			IsIsolated:    reqs[i].IsIsolated,
		})
	}
	if errStr != "" {
		return orders, fmt.Errorf("[batch]%s", errStr)
	}
	return orders, nil
}

func (b *BitgetTradeEngine) handleOrdersFromUtaTradeBatchModify(reqs []*OrderParam, res mybitgetapi.PrivateRestUtaTradeBatchModifyOrderRes) ([]*Order, error) {
	if len(res) != len(reqs) {
		return nil, errors.New("api return invalid data")
	}
	errStr := ""
	orders := make([]*Order, 0, len(reqs))
	for i, r := range res {
		if !bitgetUtaBatchRowOK(r.Code) {
			errStr += fmt.Sprintf("{[%s][%s]:%s}", r.ClientOid, r.Code, r.Msg)
		}
		orders = append(orders, &Order{
			Exchange:      BITGET_NAME.String(),
			OrderId:       r.OrderId,
			ClientOrderId: r.ClientOid,
			AccountType:   reqs[i].AccountType,
			Symbol:        reqs[i].Symbol,
			IsMargin:      reqs[i].IsMargin,
			IsIsolated:    reqs[i].IsIsolated,
		})
	}
	if errStr != "" {
		return orders, fmt.Errorf("[batch]%s", errStr)
	}
	return orders, nil
}

func (b *BitgetTradeEngine) handleOrdersFromUtaTradeCancelBatch(reqs []*OrderParam, res mybitgetapi.PrivateRestUtaTradeCancelBatchRes) ([]*Order, error) {
	if len(res) != len(reqs) {
		return nil, errors.New("api return invalid data")
	}
	errStr := ""
	orders := make([]*Order, 0, len(reqs))
	for i, r := range res {
		if !bitgetUtaBatchRowOK(r.Code) {
			errStr += fmt.Sprintf("{[%s][%s]:%s}", r.ClientOid, r.Code, r.Msg)
		}
		orders = append(orders, &Order{
			Exchange:      BITGET_NAME.String(),
			OrderId:       r.OrderId,
			ClientOrderId: r.ClientOid,
			AccountType:   reqs[i].AccountType,
			Symbol:        reqs[i].Symbol,
			IsMargin:      reqs[i].IsMargin,
			IsIsolated:    reqs[i].IsIsolated,
		})
	}
	if errStr != "" {
		return orders, fmt.Errorf("[batch]%s", errStr)
	}
	return orders, nil
}
