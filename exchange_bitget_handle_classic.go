package mytrade

import (
	"fmt"
	"strconv"
	"strings"

	mybitgetapi "github.com/Hongssd/mybitgetapi"
)

func bitgetClassicRespPreCheck[T any](res *mybitgetapi.BitgetRestRes[T]) error {
	if res == nil {
		return ErrorInvalidParam
	}
	if res.Code != "00000" {
		return fmt.Errorf("[%s]:%s", res.Code, res.Msg)
	}
	return nil
}

func (b *BitgetTradeEngine) handleOrdersFromClassicSpotQueryOpenOrders(req *QueryOrderParam, res *mybitgetapi.BitgetRestRes[mybitgetapi.PrivateRestClassicSpotTradeUnfilledOrdersRes]) ([]*Order, error) {
	if err := bitgetClassicRespPreCheck(res); err != nil {
		return nil, err
	}
	orders := make([]*Order, 0, len(res.Data))
	for _, r := range res.Data {
		ot, tif := b.converter.FromBitgetOrderType(r.OrderType)
		order := &Order{
			Exchange:      BITGET_NAME.String(),
			OrderId:       r.OrderId,
			ClientOrderId: r.ClientOid,
			AccountType:   req.AccountType,
			Symbol:        req.Symbol,
			IsMargin:      false,
			IsIsolated:    false,
			IsAlgo:        req.IsAlgo,
			Price:         r.Price,
			Quantity:      r.Size,
			ExecutedQty:   r.BaseVolume,
			AvgPrice:      r.PriceAvg,
			Status:        b.converter.FromBitgetOrderStatusClassicSpot(r.Status),
			Type:          ot,
			Side:          b.converter.FromBitgetOrderSide(r.Side),
			TimeInForce:   tif,
			CreateTime:    stringToInt64(r.CTime),
			UpdateTime:    stringToInt64(r.UTime),
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (b *BitgetTradeEngine) handleOrdersFromClassicMarginIsolatedQueryOpenOrders(req *QueryOrderParam, res *mybitgetapi.BitgetRestRes[mybitgetapi.PrivateRestClassicMarginIsolatedTradeOpenOrdersRes]) ([]*Order, error) {
	if err := bitgetClassicRespPreCheck(res); err != nil {
		return nil, err
	}
	orders := make([]*Order, 0, len(res.Data.OrderList))
	for _, r := range res.Data.OrderList {
		ot, tif := b.converter.FromBitgetOrderType(r.OrderType)
		order := &Order{
			Exchange:      BITGET_NAME.String(),
			OrderId:       r.OrderId,
			ClientOrderId: r.ClientOid,
			AccountType:   req.AccountType,
			Symbol:        req.Symbol,
			IsMargin:      false,
			IsIsolated:    false,
			IsAlgo:        req.IsAlgo,
			Price:         r.Price,
			Quantity:      r.BaseSize,
			ExecutedQty:   r.Size,
			AvgPrice:      r.PriceAvg,
			Status:        b.converter.FromBitgetOrderStatusClassicSpot(r.Status),
			Type:          ot,
			Side:          b.converter.FromBitgetOrderSide(r.Side),
			TimeInForce:   tif,
			CreateTime:    stringToInt64(r.CTime),
			UpdateTime:    stringToInt64(r.UTime),
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (b *BitgetTradeEngine) handleOrdersFromClassicMarginCrossQueryOpenOrders(req *QueryOrderParam, res *mybitgetapi.BitgetRestRes[mybitgetapi.PrivateRestClassicMarginCrossTradeOpenOrdersRes]) ([]*Order, error) {
	if err := bitgetClassicRespPreCheck(res); err != nil {
		return nil, err
	}
	orders := make([]*Order, 0, len(res.Data.OrderList))
	for _, r := range res.Data.OrderList {
		ot, tif := b.converter.FromBitgetOrderType(r.OrderType)
		order := &Order{
			Exchange:      BITGET_NAME.String(),
			OrderId:       r.OrderId,
			ClientOrderId: r.ClientOid,
			AccountType:   req.AccountType,
			Symbol:        req.Symbol,
			IsMargin:      false,
			IsIsolated:    false,
			IsAlgo:        req.IsAlgo,
			Price:         r.Price,
			Quantity:      r.BaseSize,
			ExecutedQty:   r.Size,
			AvgPrice:      r.PriceAvg,
			Status:        b.converter.FromBitgetOrderStatusClassicSpot(r.Status),
			Type:          ot,
			Side:          b.converter.FromBitgetOrderSide(r.Side),
			TimeInForce:   tif,
			CreateTime:    stringToInt64(r.CTime),
			UpdateTime:    stringToInt64(r.UTime),
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (b *BitgetTradeEngine) handleOrderFromClassicSpotOrderCreate(req *OrderParam, res *mybitgetapi.PrivateRestClassicSpotTradePlaceOrderRes) (*Order, error) {
	if res == nil {
		return nil, ErrorOrderNotFound
	}
	order := &Order{
		Exchange:      BITGET_NAME.String(),
		OrderId:       res.OrderId,
		ClientOrderId: res.ClientOid,
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
	}
	return order, nil
}

func (b *BitgetTradeEngine) handleOrderFromClassicFuturesOrderCreate(req *OrderParam, res *mybitgetapi.PrivateRestClassicFuturesTradePlaceOrderRes) (*Order, error) {
	if res == nil {
		return nil, ErrorOrderNotFound
	}
	order := &Order{
		Exchange:      BITGET_NAME.String(),
		OrderId:       res.OrderId,
		ClientOrderId: res.ClientOid,
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
	}
	return order, nil
}

func (b *BitgetTradeEngine) handleOrderFromClassicSpotCancelOrder(req *OrderParam, res *mybitgetapi.PrivateRestClassicSpotTradeCancelOrderRes) (*Order, error) {
	if res == nil {
		return nil, ErrorOrderNotFound
	}
	order := &Order{
		Exchange:      BITGET_NAME.String(),
		OrderId:       res.OrderId,
		ClientOrderId: res.ClientOid,
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
	}
	return order, nil
}

func (b *BitgetTradeEngine) handleOrderFromClassicMarginIsolatedCancelOrder(req *OrderParam, res *mybitgetapi.PrivateRestClassicMarginIsolatedTradeCancelOrderRes) (*Order, error) {
	if res == nil {
		return nil, ErrorOrderNotFound
	}
	order := &Order{
		Exchange:      BITGET_NAME.String(),
		OrderId:       res.OrderId,
		ClientOrderId: res.ClientOid,
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
	}
	return order, nil
}

func (b *BitgetTradeEngine) handleOrderFromClassicMarginCrossCancelOrder(req *OrderParam, res *mybitgetapi.PrivateRestClassicMarginCrossTradeCancelOrderRes) (*Order, error) {
	if res == nil {
		return nil, ErrorOrderNotFound
	}
	order := &Order{
		Exchange:      BITGET_NAME.String(),
		OrderId:       res.OrderId,
		ClientOrderId: res.ClientOid,
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
	}
	return order, nil
}

func (b *BitgetTradeEngine) handleOrderFromClassicFuturesCancelOrder(req *OrderParam, res *mybitgetapi.PrivateRestClassicFuturesTradeCancelOrderRes) (*Order, error) {
	if res == nil {
		return nil, ErrorOrderNotFound
	}
	order := &Order{
		Exchange:      BITGET_NAME.String(),
		OrderId:       res.OrderId,
		ClientOrderId: res.ClientOid,
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
	}
	return order, nil
}

func (b *BitgetTradeEngine) handleOrderFromClassicSpotQueryOrderInfoRow(c BitgetEnumConverter, accountType string, d *mybitgetapi.PrivateRestClassicSpotTradeOrderInfoRow) *Order {
	if d == nil {
		return nil
	}
	ot, _ := c.FromBitgetOrderType(d.OrderType)
	return &Order{
		Exchange:      BITGET_NAME.String(),
		AccountType:   accountType,
		Symbol:        d.Symbol,
		OrderId:       d.OrderId,
		ClientOrderId: d.ClientOid,
		Price:         d.Price,
		Quantity:      d.Size,
		ExecutedQty:   d.BaseVolume,
		CumQuoteQty:   d.QuoteVolume,
		AvgPrice:      d.PriceAvg,
		Status:        c.FromBitgetOrderStatusClassicSpot(d.Status),
		Type:          ot,
		Side:          c.FromBitgetOrderSide(d.Side),
		TimeInForce:   TIME_IN_FORCE_GTC,
		CreateTime:    stringToInt64(d.CTime),
		UpdateTime:    stringToInt64(d.UTime),
	}
}

func (b *BitgetTradeEngine) handleOrderFromClassicFuturesQueryOrder(req *QueryOrderParam, res *mybitgetapi.BitgetRestRes[mybitgetapi.PrivateRestClassicFuturesTradeOrderDetailRes]) (*Order, error) {
	if res == nil {
		return nil, ErrorOrderNotFound
	}
	if err := bitgetClassicRespPreCheck(res); err != nil {
		return nil, err
	}
	d := &res.Data
	ot, _ := b.converter.FromBitgetOrderType(res.Data.OrderType)
	order := &Order{
		Exchange:      BITGET_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        d.Symbol,
		OrderId:       d.OrderId,
		ClientOrderId: d.ClientOid,
		Price:         d.Price,
		Quantity:      d.Size,
		ExecutedQty:   d.BaseVolume,
		CumQuoteQty:   d.QuoteVolume,
		AvgPrice:      d.PriceAvg,
		Status:        b.converter.FromBitgetOrderStatusClassicFutures(d.State),
		Type:          ot,
		Side:          b.converter.FromBitgetOrderSide(d.Side),
		PositionSide:  b.converter.FromBitgetPositionSide(d.PosSide),
		TimeInForce:   b.converter.FromClassicForce(d.Force),
		FeeAmount:     d.Fee,
		FeeCcy:        d.MarginCoin,
		ReduceOnly:    strings.EqualFold(d.ReduceOnly, "yes") || d.ReduceOnly == "YES",
		CreateTime:    stringToInt64(d.CTime),
		UpdateTime:    stringToInt64(d.UTime),
	}
	return order, nil
}

func (b *BitgetTradeEngine) handleTradesFromClassicSpotQueryTrades(req *QueryTradeParam, res *mybitgetapi.BitgetRestRes[mybitgetapi.PrivateRestClassicSpotTradeFillsRes]) ([]*Trade, error) {
	if req == nil {
		return nil, ErrorInvalidParam
	}
	if err := bitgetClassicRespPreCheck(res); err != nil {
		return nil, err
	}
	trades := make([]*Trade, 0, len(res.Data))
	for _, r := range res.Data {
		d := &r
		fee := d.FeeDetail.TotalFee
		if fee == "" {
			fee = "0"
		}
		trades = append(trades, &Trade{
			Exchange:    BITGET_NAME.String(),
			AccountType: req.AccountType,
			Symbol:      d.Symbol,
			TradeId:     d.TradeId,
			OrderId:     d.OrderId,
			Price:       d.PriceAvg,
			Quantity:    d.Size,
			QuoteQty:    d.Amount,
			Side:        b.converter.FromBitgetOrderSide(d.Side),
			FeeAmount:   fee,
			FeeCcy:      d.FeeDetail.FeeCoin,
			IsMaker:     strings.EqualFold(d.TradeScope, "maker"),
			Timestamp:   stringToInt64(d.CTime),
		})
	}
	return trades, nil
}

func bitgetClassicFuturesFeeDetailsAgg(details []mybitgetapi.PrivateRestClassicFuturesTradeFeeDetail) (amount, ccy string) {
	if len(details) == 0 {
		return "0", ""
	}
	var sum float64
	ccy = details[0].FeeCoin
	for _, d := range details {
		if d.FeeCoin != "" {
			ccy = d.FeeCoin
		}
		sum += stringToFloat64(d.TotalFee)
	}
	return strconv.FormatFloat(sum, 'f', -1, 64), ccy
}

func (b *BitgetTradeEngine) handleTradesFromClassicFuturesQueryTrades(req *QueryTradeParam, res *mybitgetapi.BitgetRestRes[mybitgetapi.PrivateRestClassicFuturesTradeFillHistoryRes]) ([]*Trade, error) {
	if req == nil {
		return nil, ErrorInvalidParam
	}
	if err := bitgetClassicRespPreCheck(res); err != nil {
		return nil, err
	}
	trades := make([]*Trade, 0, len(res.Data.FillList))
	for _, r := range res.Data.FillList {
		d := &r
		feeAmt, feeCcy := bitgetClassicFuturesFeeDetailsAgg(d.FeeDetail)
		posSide := POSITION_SIDE_BOTH
		if strings.EqualFold(d.PosMode, BITGET_POSITION_MODE_HEDGE) {
			switch {
			case strings.EqualFold(d.Side, BITGET_ORDER_SIDE_BUY) && strings.EqualFold(d.TradeSide, BITGET_TRADE_SIDE_OPEN):
				posSide = POSITION_SIDE_LONG
			case strings.EqualFold(d.Side, BITGET_ORDER_SIDE_SELL) && strings.EqualFold(d.TradeSide, BITGET_TRADE_SIDE_OPEN):
				posSide = POSITION_SIDE_SHORT
			case strings.EqualFold(d.Side, BITGET_ORDER_SIDE_BUY) && strings.EqualFold(d.TradeSide, BITGET_TRADE_SIDE_CLOSE):
				posSide = POSITION_SIDE_SHORT
			case strings.EqualFold(d.Side, BITGET_ORDER_SIDE_SELL) && strings.EqualFold(d.TradeSide, BITGET_TRADE_SIDE_CLOSE):
				posSide = POSITION_SIDE_LONG
			}
		}
		trades = append(trades, &Trade{
			Exchange:     BITGET_NAME.String(),
			AccountType:  req.AccountType,
			Symbol:       d.Symbol,
			TradeId:      d.TradeId,
			OrderId:      d.OrderId,
			Price:        d.Price,
			Quantity:     d.BaseVolume,
			QuoteQty:     d.QuoteVolume,
			Side:         b.converter.FromBitgetOrderSide(d.Side),
			PositionSide: posSide,
			FeeAmount:    feeAmt,
			FeeCcy:       feeCcy,
			RealizedPnl:  d.Profit,
			IsMaker:      strings.EqualFold(d.TradeScope, "maker"),
			Timestamp:    stringToInt64(d.CTime),
		})
	}
	return trades, nil
}

func (b *BitgetTradeEngine) handleOrdersFromClassicSpotQueryOrders(req *QueryOrderParam, res *mybitgetapi.BitgetRestRes[mybitgetapi.PrivateRestClassicSpotTradeHistoryOrdersRes]) ([]*Order, error) {
	if req == nil {
		return nil, ErrorInvalidParam
	}
	if err := bitgetClassicRespPreCheck(res); err != nil {
		return nil, err
	}
	orders := make([]*Order, 0, len(res.Data))
	for _, r := range res.Data {
		ot, _ := b.converter.FromBitgetOrderType(r.OrderType)
		order := &Order{
			Exchange:      BITGET_NAME.String(),
			AccountType:   req.AccountType,
			Symbol:        r.Symbol,
			OrderId:       r.OrderId,
			ClientOrderId: r.ClientOid,
			IsMargin:      req.IsMargin,
			IsIsolated:    req.IsIsolated,
			IsAlgo:        req.IsAlgo,
			OrderAlgoType: req.OrderAlgoType,
			Price:         r.Price,
			Quantity:      r.Size,
			ExecutedQty:   r.BaseVolume,
			CumQuoteQty:   r.QuoteVolume,
			AvgPrice:      r.PriceAvg,
			Status:        b.converter.FromBitgetOrderStatusClassicSpot(r.Status),
			Type:          ot,
			Side:          b.converter.FromBitgetOrderSide(r.Side),
			TimeInForce:   TIME_IN_FORCE_GTC,
			CreateTime:    stringToInt64(r.CTime),
			UpdateTime:    stringToInt64(r.UTime),
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (b *BitgetTradeEngine) handleOrdersFromClassicMarginIsolatedQueryOrders(req *QueryOrderParam, res *mybitgetapi.BitgetRestRes[mybitgetapi.PrivateRestClassicMarginIsolatedTradeHistoryOrdersRes]) ([]*Order, error) {
	if req == nil {
		return nil, ErrorInvalidParam
	}
	if err := bitgetClassicRespPreCheck(res); err != nil {
		return nil, err
	}
	orders := make([]*Order, 0, len(res.Data.OrderList))
	for _, r := range res.Data.OrderList {
		ot, _ := b.converter.FromBitgetOrderType(r.OrderType)
		order := &Order{
			Exchange:      BITGET_NAME.String(),
			AccountType:   req.AccountType,
			Symbol:        r.Symbol,
			OrderId:       r.OrderId,
			ClientOrderId: r.ClientOid,
			IsMargin:      true,
			IsIsolated:    true,
			Price:         r.Price,
			Quantity:      r.BaseSize,
			ExecutedQty:   r.Size,
			CumQuoteQty:   r.Amount,
			AvgPrice:      r.PriceAvg,
			Status:        b.converter.FromBitgetOrderStatusClassicSpot(r.Status),
			Type:          ot,
			Side:          b.converter.FromBitgetOrderSide(r.Side),
			TimeInForce:   b.converter.FromClassicForce(r.Force),
			CreateTime:    stringToInt64(r.CTime),
			UpdateTime:    stringToInt64(r.UTime),
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (b *BitgetTradeEngine) handleOrdersFromClassicMarginCrossedQueryOrders(req *QueryOrderParam, res *mybitgetapi.BitgetRestRes[mybitgetapi.PrivateRestClassicMarginCrossTradeHistoryOrdersRes]) ([]*Order, error) {
	if req == nil {
		return nil, ErrorInvalidParam
	}
	if err := bitgetClassicRespPreCheck(res); err != nil {
		return nil, err
	}
	orders := make([]*Order, 0, len(res.Data.OrderList))
	for _, r := range res.Data.OrderList {
		ot, _ := b.converter.FromBitgetOrderType(r.OrderType)
		order := &Order{
			Exchange:      BITGET_NAME.String(),
			AccountType:   req.AccountType,
			Symbol:        r.Symbol,
			OrderId:       r.OrderId,
			ClientOrderId: r.ClientOid,
			IsMargin:      req.IsMargin,
			IsIsolated:    req.IsIsolated,
			IsAlgo:        req.IsAlgo,
			OrderAlgoType: req.OrderAlgoType,
			Price:         r.Price,
			Quantity:      r.BaseSize,
			ExecutedQty:   r.Size,
			CumQuoteQty:   r.Amount,
			AvgPrice:      r.PriceAvg,
			Status:        b.converter.FromBitgetOrderStatusClassicSpot(r.Status),
			Type:          ot,
			Side:          b.converter.FromBitgetOrderSide(r.Side),
			TimeInForce:   b.converter.FromClassicForce(r.Force),
			CreateTime:    stringToInt64(r.CTime),
			UpdateTime:    stringToInt64(r.UTime),
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (b *BitgetTradeEngine) handleOrderFromClassicSpotQueryOrder(req *QueryOrderParam, res *mybitgetapi.BitgetRestRes[mybitgetapi.PrivateRestClassicSpotTradeOrderInfoRes]) (*Order, error) {
	if req == nil || res == nil || len(res.Data) == 0 {
		return nil, ErrorOrderNotFound
	}
	match := func(d *mybitgetapi.PrivateRestClassicSpotTradeOrderInfoRow) bool {
		if d == nil {
			return false
		}
		if req.OrderId != "" && d.OrderId == req.OrderId {
			return true
		}
		if req.ClientOrderId != "" && d.ClientOid == req.ClientOrderId {
			return true
		}
		return false
	}
	for i := range res.Data {
		if match(&res.Data[i]) {
			return b.handleOrderFromClassicSpotQueryOrderInfoRow(b.converter, req.AccountType, &res.Data[i]), nil
		}
	}
	return nil, ErrorOrderNotFound
}

func (b *BitgetTradeEngine) handleTradesFromClassicMarginIsolatedQueryTrades(req *QueryTradeParam, res *mybitgetapi.BitgetRestRes[mybitgetapi.PrivateRestClassicMarginIsolatedTradeFillsRes]) ([]*Trade, error) {
	if req == nil {
		return nil, ErrorInvalidParam
	}
	if err := bitgetClassicRespPreCheck(res); err != nil {
		return nil, err
	}
	trades := make([]*Trade, 0, len(res.Data.Fills))
	for _, r := range res.Data.Fills {
		d := &r
		fee := d.FeeDetail.TotalFee
		if fee == "" {
			fee = "0"
		}
		trades = append(trades, &Trade{
			Exchange:    BITGET_NAME.String(),
			AccountType: req.AccountType,
			TradeId:     d.TradeId,
			OrderId:     d.OrderId,
			Price:       d.PriceAvg,
			Quantity:    d.Size,
			QuoteQty:    d.Amount,
			Side:        b.converter.FromBitgetOrderSide(d.Side),
			FeeAmount:   fee,
			FeeCcy:      d.FeeDetail.FeeCoin,
			IsMaker:     strings.EqualFold(d.TradeScope, "maker"),
			Timestamp:   stringToInt64(d.CTime),
		})
	}
	return trades, nil
}

func (b *BitgetTradeEngine) handleTradesFromClassicMarginCrossedQueryTrades(req *QueryTradeParam, res *mybitgetapi.BitgetRestRes[mybitgetapi.PrivateRestClassicMarginCrossTradeFillsRes]) ([]*Trade, error) {
	if req == nil {
		return nil, ErrorInvalidParam
	}
	if err := bitgetClassicRespPreCheck(res); err != nil {
		return nil, err
	}
	trades := make([]*Trade, 0, len(res.Data.Fills))
	for _, r := range res.Data.Fills {
		d := &r
		fee := d.FeeDetail.TotalFee
		if fee == "" {
			fee = "0"
		}
		trades = append(trades, &Trade{
			Exchange:    BITGET_NAME.String(),
			AccountType: req.AccountType,
			TradeId:     d.TradeId,
			OrderId:     d.OrderId,
			Price:       d.PriceAvg,
			Quantity:    d.Size,
			QuoteQty:    d.Amount,
			Side:        b.converter.FromBitgetOrderSide(d.Side),
			FeeAmount:   fee,
			FeeCcy:      d.FeeDetail.FeeCoin,
			IsMaker:     strings.EqualFold(d.TradeScope, "maker"),
			Timestamp:   stringToInt64(d.CTime),
		})
	}
	return trades, nil
}

func (b *BitgetTradeEngine) handleOrderFromClassicMarginIsolatedOrderCreate(req *OrderParam, res *mybitgetapi.PrivateRestClassicMarginIsolatedTradePlaceOrderRes) (*Order, error) {
	if res == nil {
		return nil, ErrorOrderNotFound
	}
	return &Order{
		Exchange:      BITGET_NAME.String(),
		OrderId:       res.OrderId,
		ClientOrderId: res.ClientOid,
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
	}, nil
}

func (b *BitgetTradeEngine) handleOrderFromClassicMarginCrossedOrderCreate(req *OrderParam, res *mybitgetapi.PrivateRestClassicMarginCrossTradePlaceOrderRes) (*Order, error) {
	if res == nil {
		return nil, ErrorOrderNotFound
	}
	return &Order{
		Exchange:      BITGET_NAME.String(),
		OrderId:       res.OrderId,
		ClientOrderId: res.ClientOid,
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		IsMargin:      req.IsMargin,
		IsIsolated:    req.IsIsolated,
	}, nil
}

func (b *BitgetTradeEngine) handleOrdersFromClassicFuturesQueryOrders(req *QueryOrderParam, res *mybitgetapi.BitgetRestRes[mybitgetapi.PrivateRestClassicFuturesTradeOrdersHistoryRes]) ([]*Order, error) {
	if req == nil {
		return nil, ErrorInvalidParam
	}
	if err := bitgetClassicRespPreCheck(res); err != nil {
		return nil, err
	}
	orders := make([]*Order, 0, len(res.Data.EntrustedList))
	for i := range res.Data.EntrustedList {
		r := &res.Data.EntrustedList[i]
		ot, _ := b.converter.FromBitgetOrderType(r.OrderType)
		orders = append(orders, &Order{
			Exchange:      BITGET_NAME.String(),
			AccountType:   req.AccountType,
			Symbol:        r.Symbol,
			OrderId:       r.OrderId,
			ClientOrderId: r.ClientOid,
			IsMargin:      req.IsMargin,
			IsIsolated:    req.IsIsolated,
			IsAlgo:        req.IsAlgo,
			OrderAlgoType: req.OrderAlgoType,
			Price:         r.Price,
			Quantity:      r.Size,
			ExecutedQty:   r.BaseVolume,
			CumQuoteQty:   r.QuoteVolume,
			AvgPrice:      r.PriceAvg,
			Status:        b.converter.FromBitgetOrderStatusClassicFutures(r.Status),
			Type:          ot,
			Side:          b.converter.FromBitgetOrderSide(r.Side),
			PositionSide:  b.converter.FromBitgetPositionSide(r.PosSide),
			TimeInForce:   b.converter.FromClassicForce(r.Force),
			FeeAmount:     r.Fee,
			FeeCcy:        r.MarginCoin,
			ReduceOnly:    b.converter.ReduceOnlyFromString(r.ReduceOnly),
			CreateTime:    stringToInt64(r.CTime),
			UpdateTime:    stringToInt64(r.UTime),
		})
	}
	return orders, nil
}

func (b *BitgetTradeEngine) handleOrdersFromClassicSpotBatchOrders(reqs []*OrderParam, res *mybitgetapi.PrivateRestClassicSpotTradeBatchResultRes) ([]*Order, error) {
	orderMap := make(map[string]struct {
		orderId, clientOid, errCode, errMsg string
	})
	for _, s := range res.SuccessList {
		k := s.ClientOid
		if k == "" {
			k = s.OrderId
		}
		orderMap[k] = struct {
			orderId, clientOid, errCode, errMsg string
		}{s.OrderId, s.ClientOid, "", ""}
	}
	for _, f := range res.FailureList {
		k := f.ClientOid
		if k == "" {
			k = f.OrderId
		}
		orderMap[k] = struct {
			orderId, clientOid, errCode, errMsg string
		}{f.OrderId, f.ClientOid, f.ErrorCode, f.ErrorMsg}
	}
	errStr := ""
	orders := make([]*Order, 0, len(reqs))
	for _, req := range reqs {
		k := req.ClientOrderId
		row := orderMap[k]
		if row.errMsg != "" || row.errCode != "" {
			errStr += fmt.Sprintf("{[%s][%s]:%s}", k, row.errCode, row.errMsg)
		}
		orders = append(orders, &Order{
			Exchange:      BITGET_NAME.String(),
			OrderId:       row.orderId,
			ClientOrderId: row.clientOid,
			AccountType:   req.AccountType,
			Symbol:        req.Symbol,
			IsMargin:      req.IsMargin,
			IsIsolated:    req.IsIsolated,
		})
	}
	if errStr != "" {
		return orders, fmt.Errorf("[batch]%s", errStr)
	}
	return orders, nil
}

func (b *BitgetTradeEngine) handleOrdersFromClassicFuturesBatchOrders(reqs []*OrderParam, res *mybitgetapi.PrivateRestClassicFuturesTradeBatchRes) ([]*Order, error) {
	orderMap := make(map[string]struct {
		orderId, clientOid, errCode, errMsg string
	})
	for _, s := range res.SuccessList {
		k := s.ClientOid
		if k == "" {
			k = s.OrderId
		}
		orderMap[k] = struct {
			orderId, clientOid, errCode, errMsg string
		}{s.OrderId, s.ClientOid, "", ""}
	}
	for _, f := range res.FailureList {
		k := f.ClientOid
		if k == "" {
			k = f.OrderId
		}
		orderMap[k] = struct {
			orderId, clientOid, errCode, errMsg string
		}{f.OrderId, f.ClientOid, f.ErrorCode, f.ErrorMsg}
	}
	errStr := ""
	orders := make([]*Order, 0, len(reqs))
	for _, r := range reqs {
		k := r.ClientOrderId
		row := orderMap[k]
		if row.errMsg != "" || row.errCode != "" {
			errStr += fmt.Sprintf("{[%s][%s]:%s}", k, row.errCode, row.errMsg)
		}
		orders = append(orders, &Order{
			Exchange:      BITGET_NAME.String(),
			OrderId:       row.orderId,
			ClientOrderId: row.clientOid,
			AccountType:   r.AccountType,
			Symbol:        r.Symbol,
			IsMargin:      r.IsMargin,
			IsIsolated:    r.IsIsolated,
		})
	}
	if errStr != "" {
		return orders, fmt.Errorf("[batch]%s", errStr)
	}
	return orders, nil
}

func (b *BitgetTradeEngine) handleOrdersFromClassicSpotBatchCancelOrders(reqs []*OrderParam, res *mybitgetapi.PrivateRestClassicSpotTradeBatchCancelOrderRes) ([]*Order, error) {
	orderMap := make(map[string]struct {
		orderId, clientOid, errCode, errMsg string
	})
	for _, s := range res.SuccessList {
		k := s.ClientOid
		if k == "" {
			k = s.OrderId
		}
		orderMap[k] = struct {
			orderId, clientOid, errCode, errMsg string
		}{s.OrderId, s.ClientOid, "", ""}
	}
	for _, f := range res.FailureList {
		k := f.ClientOid
		if k == "" {
			k = f.OrderId
		}
		orderMap[k] = struct {
			orderId, clientOid, errCode, errMsg string
		}{f.OrderId, f.ClientOid, f.ErrorCode, f.ErrorMsg}
	}
	errStr := ""
	orders := make([]*Order, 0, len(reqs))
	for _, req := range reqs {
		k := req.ClientOrderId
		row := orderMap[k]
		if row.errMsg != "" || row.errCode != "" {
			errStr += fmt.Sprintf("{[%s][%s]:%s}", k, row.errCode, row.errMsg)
		}
		orders = append(orders, &Order{
			Exchange:      BITGET_NAME.String(),
			OrderId:       row.orderId,
			ClientOrderId: row.clientOid,
			AccountType:   req.AccountType,
			Symbol:        req.Symbol,
			IsMargin:      req.IsMargin,
			IsIsolated:    req.IsIsolated,
		})
	}
	if errStr != "" {
		return orders, fmt.Errorf("[batch]%s", errStr)
	}
	return orders, nil
}

func (b *BitgetTradeEngine) handleOrdersFromClassicMarginIsolatedBatchCancelOrders(reqs []*OrderParam, res *mybitgetapi.PrivateRestClassicMarginIsolatedTradeBatchCancelOrderRes) ([]*Order, error) {
	orderMap := make(map[string]struct {
		orderId, clientOid, errCode, errMsg string
	})
	for _, s := range res.SuccessList {
		k := s.ClientOid
		if k == "" {
			k = s.OrderId
		}
		orderMap[k] = struct {
			orderId, clientOid, errCode, errMsg string
		}{s.OrderId, s.ClientOid, "", ""}
	}
	for _, f := range res.FailureList {
		k := f.ClientOid
		if k == "" {
			k = f.OrderId
		}
		orderMap[k] = struct {
			orderId, clientOid, errCode, errMsg string
		}{f.OrderId, f.ClientOid, f.ErrorCode, f.ErrorMsg}
	}
	errStr := ""
	orders := make([]*Order, 0, len(reqs))
	for _, req := range reqs {
		k := req.ClientOrderId
		row := orderMap[k]
		if row.errMsg != "" || row.errCode != "" {
			errStr += fmt.Sprintf("{[%s][%s]:%s}", k, row.errCode, row.errMsg)
		}
		orders = append(orders, &Order{
			Exchange:      BITGET_NAME.String(),
			OrderId:       row.orderId,
			ClientOrderId: row.clientOid,
			AccountType:   req.AccountType,
			Symbol:        req.Symbol,
			IsMargin:      req.IsMargin,
			IsIsolated:    req.IsIsolated,
		})
	}
	if errStr != "" {
		return orders, fmt.Errorf("[batch]%s", errStr)
	}
	return orders, nil
}

func (b *BitgetTradeEngine) handleOrdersFromClassicMarginCrossBatchCancelOrders(reqs []*OrderParam, res *mybitgetapi.PrivateRestClassicMarginCrossTradeBatchCancelOrderRes) ([]*Order, error) {
	orderMap := make(map[string]struct {
		orderId, clientOid, errCode, errMsg string
	})
	for _, s := range res.SuccessList {
		k := s.ClientOid
		if k == "" {
			k = s.OrderId
		}
		orderMap[k] = struct {
			orderId, clientOid, errCode, errMsg string
		}{s.OrderId, s.ClientOid, "", ""}
	}
	for _, f := range res.FailureList {
		k := f.ClientOid
		if k == "" {
			k = f.OrderId
		}
		orderMap[k] = struct {
			orderId, clientOid, errCode, errMsg string
		}{f.OrderId, f.ClientOid, f.ErrorCode, f.ErrorMsg}
	}
	errStr := ""
	orders := make([]*Order, 0, len(reqs))
	for _, req := range reqs {
		k := req.ClientOrderId
		row := orderMap[k]
		if row.errMsg != "" || row.errCode != "" {
			errStr += fmt.Sprintf("{[%s][%s]:%s}", k, row.errCode, row.errMsg)
		}
		orders = append(orders, &Order{
			Exchange:      BITGET_NAME.String(),
			OrderId:       row.orderId,
			ClientOrderId: row.clientOid,
			AccountType:   req.AccountType,
			Symbol:        req.Symbol,
			IsMargin:      req.IsMargin,
			IsIsolated:    req.IsIsolated,
		})
	}
	if errStr != "" {
		return orders, fmt.Errorf("[batch]%s", errStr)
	}
	return orders, nil
}

func (b *BitgetTradeEngine) handleOrdersFromClassicFuturesBatchCancelOrders(reqs []*OrderParam, res *mybitgetapi.PrivateRestClassicFuturesTradeBatchCancelOrderRes) ([]*Order, error) {
	orderMap := make(map[string]struct {
		orderId, clientOid, errCode, errMsg string
	})
	for _, s := range res.SuccessList {
		k := s.ClientOid
		if k == "" {
			k = s.OrderId
		}
		orderMap[k] = struct {
			orderId, clientOid, errCode, errMsg string
		}{s.OrderId, s.ClientOid, "", ""}
	}
	for _, f := range res.FailureList {
		k := f.ClientOid
		if k == "" {
			k = f.OrderId
		}
		orderMap[k] = struct {
			orderId, clientOid, errCode, errMsg string
		}{f.OrderId, f.ClientOid, f.ErrorCode, f.ErrorMsg}
	}
	errStr := ""
	orders := make([]*Order, 0, len(reqs))
	for _, req := range reqs {
		k := req.ClientOrderId
		row := orderMap[k]
		if row.errMsg != "" || row.errCode != "" {
			errStr += fmt.Sprintf("{[%s][%s]:%s}", k, row.errCode, row.errMsg)
		}
		orders = append(orders, &Order{
			Exchange:      BITGET_NAME.String(),
			OrderId:       row.orderId,
			ClientOrderId: row.clientOid,
			AccountType:   req.AccountType,
			Symbol:        req.Symbol,
			IsMargin:      req.IsMargin,
			IsIsolated:    req.IsIsolated,
		})
	}
	if errStr != "" {
		return orders, fmt.Errorf("[batch]%s", errStr)
	}
	return orders, nil
}

func (b *BitgetTradeEngine) handleOrdersFromClassicMarginCrossBatchOrders(reqs []*OrderParam, data *mybitgetapi.PrivateRestClassicMarginCrossTradeBatchRes) ([]*Order, error) {
	orderMap := make(map[string]struct {
		orderId, clientOid, errCode, errMsg string
	})
	for _, s := range data.SuccessList {
		k := s.ClientOid
		if k == "" {
			k = s.OrderId
		}
		orderMap[k] = struct {
			orderId, clientOid, errCode, errMsg string
		}{s.OrderId, s.ClientOid, "", ""}
	}
	for _, f := range data.FailureList {
		k := f.ClientOid
		if k == "" {
			k = f.OrderId
		}
		orderMap[k] = struct {
			orderId, clientOid, errCode, errMsg string
		}{f.OrderId, f.ClientOid, f.ErrorCode, f.ErrorMsg}
	}
	errStr := ""
	orders := make([]*Order, 0, len(reqs))
	for _, r := range reqs {
		k := r.ClientOrderId
		row := orderMap[k]
		if row.errMsg != "" || row.errCode != "" {
			errStr += fmt.Sprintf("{[%s][%s]:%s}", k, row.errCode, row.errMsg)
		}
		orders = append(orders, &Order{
			Exchange:      BITGET_NAME.String(),
			OrderId:       row.orderId,
			ClientOrderId: row.clientOid,
			AccountType:   r.AccountType,
			Symbol:        r.Symbol,
			IsMargin:      r.IsMargin,
			IsIsolated:    r.IsIsolated,
		})
	}
	if errStr != "" {
		return orders, fmt.Errorf("[batch]%s", errStr)
	}
	return orders, nil
}

func (b *BitgetTradeEngine) handleOrdersFromClassicMarginIsolatedBatchOrders(reqs []*OrderParam, data *mybitgetapi.PrivateRestClassicMarginIsolatedTradeBatchRes) ([]*Order, error) {
	orderMap := make(map[string]struct {
		orderId, clientOid, errCode, errMsg string
	})
	for _, s := range data.SuccessList {
		k := s.ClientOid
		if k == "" {
			k = s.OrderId
		}
		orderMap[k] = struct {
			orderId, clientOid, errCode, errMsg string
		}{s.OrderId, s.ClientOid, "", ""}
	}
	for _, f := range data.FailureList {
		k := f.ClientOid
		if k == "" {
			k = f.OrderId
		}
		orderMap[k] = struct {
			orderId, clientOid, errCode, errMsg string
		}{f.OrderId, f.ClientOid, f.ErrorCode, f.ErrorMsg}
	}
	errStr := ""
	orders := make([]*Order, 0, len(reqs))
	for _, r := range reqs {
		k := r.ClientOrderId
		row := orderMap[k]
		if row.errMsg != "" || row.errCode != "" {
			errStr += fmt.Sprintf("{[%s][%s]:%s}", k, row.errCode, row.errMsg)
		}
		orders = append(orders, &Order{
			Exchange:      BITGET_NAME.String(),
			OrderId:       row.orderId,
			ClientOrderId: row.clientOid,
			AccountType:   r.AccountType,
			Symbol:        r.Symbol,
			IsMargin:      r.IsMargin,
			IsIsolated:    r.IsIsolated,
		})
	}
	if errStr != "" {
		return orders, fmt.Errorf("[batch]%s", errStr)
	}
	return orders, nil
}
