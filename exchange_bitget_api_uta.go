package mytrade

import (
	"strconv"

	mybitgetapi "github.com/Hongssd/mybitgetapi"
)

const (
	BITGET_BATCH_UTA_ORDER_MAX = 50
)

// UTA 查询当前委托
func (e *BitgetTradeEngine) apiUtaTradeUnfilledOrders(req *QueryOrderParam) *mybitgetapi.PrivateRestUtaTradeUnfilledOrdersAPI {
	api := e.bitgetPrivateClient().NewPrivateRestUtaTradeUnfilledOrders().Category(req.AccountType)
	if req.Symbol != "" {
		api.Symbol(req.Symbol)
	}
	if req.StartTime != 0 {
		api.StartTime(strconv.FormatInt(req.StartTime, 10))
	}
	if req.EndTime != 0 {
		api.EndTime(strconv.FormatInt(req.EndTime, 10))
	}
	if ls := bitgetLimitString(req.Limit); ls != "" {
		api.Limit(ls)
	}
	return api
}

// UTA 查询订单详情
func (e *BitgetTradeEngine) apiUtaTradeOrderInfo(req *QueryOrderParam) *mybitgetapi.PrivateRestUtaTradeOrderInfoAPI {
	api := e.bitgetPrivateClient().NewPrivateRestUtaTradeOrderInfo()
	if req.OrderId != "" {
		api.OrderId(req.OrderId)
	}
	if req.ClientOrderId != "" {
		api.ClientOid(req.ClientOrderId)
	}
	return api
}

// UTA 查询历史委托
func (e *BitgetTradeEngine) apiUtaTradeHistoryOrders(req *QueryOrderParam) *mybitgetapi.PrivateRestUtaTradeHistoryOrdersAPI {
	api := e.bitgetPrivateClient().NewPrivateRestUtaTradeHistoryOrders().Category(req.AccountType)
	if req.Symbol != "" {
		api.Symbol(req.Symbol)
	}
	if req.StartTime != 0 {
		api.StartTime(strconv.FormatInt(req.StartTime, 10))
	}
	if req.EndTime != 0 {
		api.EndTime(strconv.FormatInt(req.EndTime, 10))
	}
	if ls := bitgetLimitString(req.Limit); ls != "" {
		api.Limit(ls)
	}
	return api
}

// UTA 查询成交明细
func (e *BitgetTradeEngine) apiUtaTradeFills(req *QueryTradeParam) *mybitgetapi.PrivateRestUtaTradeFillsAPI {
	api := e.bitgetPrivateClient().NewPrivateRestUtaTradeFills().Category(req.AccountType)
	if req.OrderId != "" {
		api.OrderId(req.OrderId)
	}
	if req.StartTime != 0 {
		api.StartTime(strconv.FormatInt(req.StartTime, 10))
	}
	if req.EndTime != 0 {
		api.EndTime(strconv.FormatInt(req.EndTime, 10))
	}
	if ls := bitgetLimitString(req.Limit); ls != "" {
		api.Limit(ls)
	}
	return api
}

// ——— UTA 下单 / 改单 / 撤单 ———

func (e *BitgetTradeEngine) apiUtaTradePlaceOrder(req *OrderParam, clientOid string) *mybitgetapi.PrivateRestUtaTradePlaceOrderAPI {
	api := e.bitgetPrivateClient().NewPrivateRestUtaTradePlaceOrder().
		Category(req.AccountType).Symbol(req.Symbol).
		Qty(req.Quantity.String()).
		Side(e.converter.ToBitgetOrderSide(req.OrderSide)).
		OrderType(e.converter.ToBitgetOrderType(req.OrderType)).
		ClientOid(clientOid)
	if req.OrderType == ORDER_TYPE_LIMIT {
		api.Price(req.Price.String())
		api.TimeInForce(e.converter.ToBitgetTimeInForce(req.TimeInForce))
	}
	switch req.AccountType {
	case BITGET_AC_USDT_FUTURES, BITGET_AC_COIN_FUTURES, BITGET_AC_USDC_FUTURES:
		if req.PositionSide != POSITION_SIDE_BOTH && req.PositionSide != POSITION_SIDE_UNKNOWN {
			api.PosSide(e.converter.ToBitgetPositionSide(req.PositionSide))
		}
		if req.ReduceOnly {
			api.ReduceOnly(BITGET_REDUCE_ONLY_YES)
		} else {
			api.ReduceOnly(BITGET_REDUCE_ONLY_NO)
		}
	}
	return api
}

func (e *BitgetTradeEngine) apiUtaTradeModifyOrder(req *OrderParam) *mybitgetapi.PrivateRestUtaTradeModifyOrderAPI {
	api := e.bitgetPrivateClient().NewPrivateRestUtaTradeModifyOrder().
		Symbol(req.Symbol).Category(req.AccountType)
	if req.OrderId != "" {
		api.OrderId(req.OrderId)
	}
	if req.ClientOrderId != "" {
		api.ClientOid(req.ClientOrderId)
	}
	api.Qty(req.Quantity.String())
	if req.OrderType == ORDER_TYPE_LIMIT {
		api.Price(req.Price.String())
	}
	return api
}

func (e *BitgetTradeEngine) apiUtaTradeCancelOrder(req *OrderParam) *mybitgetapi.PrivateRestUtaTradeCancelOrderAPI {
	api := e.bitgetPrivateClient().NewPrivateRestUtaTradeCancelOrder()
	if req.OrderId != "" {
		api.OrderId(req.OrderId)
	}
	if req.ClientOrderId != "" {
		api.ClientOid(req.ClientOrderId)
	}
	if req.AccountType != "" {
		api.Category(req.AccountType)
	}
	return api
}

func (e *BitgetTradeEngine) apiUtaTradePlaceBatch(reqs []*OrderParam, clientOids []string) *mybitgetapi.PrivateRestUtaTradePlaceBatchAPI {
	api := e.bitgetPrivateClient().NewPrivateRestUtaTradePlaceBatch()
	for i, r := range reqs {
		item := mybitgetapi.PrivateRestUtaTradePlaceBatchReqItem{
			Category:  r.AccountType,
			Symbol:    r.Symbol,
			Qty:       r.Quantity.String(),
			Side:      e.converter.ToBitgetOrderSide(r.OrderSide),
			OrderType: e.converter.ToBitgetOrderType(r.OrderType),
		}
		if r.OrderType == ORDER_TYPE_LIMIT {
			p := r.Price.String()
			item.Price = &p
			tif := e.converter.ToBitgetTimeInForce(r.TimeInForce)
			item.TimeInForce = &tif
		}
		co := clientOids[i]
		item.ClientOid = &co
		switch r.AccountType {
		case BITGET_AC_USDT_FUTURES, BITGET_AC_COIN_FUTURES, BITGET_AC_USDC_FUTURES:
			if r.PositionSide != POSITION_SIDE_BOTH && r.PositionSide != POSITION_SIDE_UNKNOWN && r.PositionSide != "" {
				ps := e.converter.ToBitgetPositionSide(r.PositionSide)
				item.PosSide = &ps
			}
		}
		api.AddOrder(item)
	}
	return api
}

func (e *BitgetTradeEngine) apiUtaTradeBatchModifyOrder(reqs []*OrderParam) *mybitgetapi.PrivateRestUtaTradeBatchModifyOrderAPI {
	api := e.bitgetPrivateClient().NewPrivateRestUtaTradeBatchModifyOrder()
	for _, r := range reqs {
		item := mybitgetapi.PrivateRestUtaTradeBatchModifyOrderReqItem{
			Symbol:   GetPointer(r.Symbol),
			Category: GetPointer(r.AccountType),
		}
		if r.OrderId != "" {
			item.OrderId = GetPointer(r.OrderId)
		}
		if r.ClientOrderId != "" {
			item.ClientOid = GetPointer(r.ClientOrderId)
		}
		q := r.Quantity.String()
		item.Qty = &q
		if r.OrderType == ORDER_TYPE_LIMIT {
			p := r.Price.String()
			item.Price = &p
		}
		api.AddOrder(item)
	}
	return api
}

func (e *BitgetTradeEngine) apiUtaTradeCancelBatch(reqs []*OrderParam) *mybitgetapi.PrivateRestUtaTradeCancelBatchAPI {
	api := e.bitgetPrivateClient().NewPrivateRestUtaTradeCancelBatch()
	for _, r := range reqs {
		item := mybitgetapi.PrivateRestUtaTradeCancelBatchReqOrder{
			Category: GetPointer(r.AccountType),
			Symbol:   GetPointer(r.Symbol),
		}
		if r.OrderId != "" {
			item.OrderId = GetPointer(r.OrderId)
		}
		if r.ClientOrderId != "" {
			item.ClientOid = GetPointer(r.ClientOrderId)
		}
		api.AddOrder(item)
	}
	return api
}
