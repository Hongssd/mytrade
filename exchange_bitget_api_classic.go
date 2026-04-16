package mytrade

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	mybitgetapi "github.com/Hongssd/mybitgetapi"
)

const (
	BITGET_BATCH_CLASSIC_ORDER_MAX = 20
)

func bitgetLimitString(limit int) string {
	if limit <= 0 {
		return ""
	}
	return strconv.Itoa(limit)
}

func (e *BitgetTradeEngine) bitgetPrivateClient() *mybitgetapi.PrivateRestClient {
	return e.privateClient
}

func (e *BitgetTradeEngine) accountTypePreCheck(accountType string) error {
	switch accountType {
	case BITGET_AC_SPOT, BITGET_AC_MARGIN, BITGET_AC_USDT_FUTURES, BITGET_AC_COIN_FUTURES, BITGET_AC_USDC_FUTURES:
		return nil
	default:
		return ErrorAccountType
	}
}

func (e *BitgetTradeEngine) restBatchPreCheck(reqs []*OrderParam) error {
	if len(reqs) == 0 {
		return ErrorInvalidParam
	}
	if e.isClassic {
		if len(reqs) > BITGET_BATCH_CLASSIC_ORDER_MAX {
			return ErrorInvalid(fmt.Sprintf("bitget batch order count must be <= %d", BITGET_BATCH_CLASSIC_ORDER_MAX))
		}
	} else {
		if len(reqs) > BITGET_BATCH_UTA_ORDER_MAX {
			return ErrorInvalid(fmt.Sprintf("bitget batch order count must be <= %d", BITGET_BATCH_UTA_ORDER_MAX))
		}
	}
	for _, r := range reqs {
		if r == nil {
			return ErrorInvalidParam
		}
		if err := e.accountTypePreCheck(r.AccountType); err != nil {
			return err
		}
	}
	return nil
}

func bitgetMarginCoinFromSymbol(symbol, ccy string) string {
	if ccy != "" {
		return ccy
	}
	for _, suf := range []string{"USDT", "USDC", "USD"} {
		if strings.HasSuffix(symbol, suf) {
			return suf
		}
	}
	return ""
}

func bitgetLoanTypeFromReq(req *OrderParam) string {
	if req.SideEffectType != "" {
		return req.SideEffectType
	}
	return BITGET_LOAN_TYPE_NORMAL
}

func bitgetClassicFuturesMarginMode(req *OrderParam) string {
	if req.IsIsolated {
		return BITGET_MARGIN_MODE_ISOLATED
	}
	return BITGET_MARGIN_MODE_CROSSED
}

// Classic 现货 — 当前委托
func (e *BitgetTradeEngine) apiClassicSpotQueryOpenOrders(req *QueryOrderParam) *mybitgetapi.PrivateRestClassicSpotTradeUnfilledOrdersAPI {
	api := e.bitgetPrivateClient().NewPrivateRestClassicSpotTradeUnfilledOrders()
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

// Classic 杠杆 - 全仓 - 当前委托
func (e *BitgetTradeEngine) apiClassicMarginCrossQueryOpenOrders(req *QueryOrderParam) *mybitgetapi.PrivateRestClassicMarginCrossTradeOpenOrdersAPI {
	api := e.bitgetPrivateClient().NewPrivateRestClassicMarginCrossTradeOpenOrders()
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

// Classic 杠杆 - 逐仓 - 当前委托
func (e *BitgetTradeEngine) apiClassicMarginIsolatedQueryOpenOrders(req *QueryOrderParam) *mybitgetapi.PrivateRestClassicMarginIsolatedTradeOpenOrdersAPI {
	api := e.bitgetPrivateClient().NewPrivateRestClassicMarginIsolatedTradeOpenOrders()
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

// Classic 现货 — 订单详情
func (e *BitgetTradeEngine) apiClassicSpotQueryOrder(req *QueryOrderParam) *mybitgetapi.PrivateRestClassicSpotTradeOrderInfoAPI {
	api := e.bitgetPrivateClient().NewPrivateRestClassicSpotTradeOrderInfo()
	if req.OrderId != "" {
		api.OrderId(req.OrderId)
	}
	if req.ClientOrderId != "" {
		api.ClientOid(req.ClientOrderId)
	}
	return api
}

// Classic 现货 — 历史委托
func (e *BitgetTradeEngine) apiClassicSpotQueryOrders(req *QueryOrderParam) *mybitgetapi.PrivateRestClassicSpotTradeHistoryOrdersAPI {
	api := e.bitgetPrivateClient().NewPrivateRestClassicSpotTradeHistoryOrders()
	if req.Symbol != "" {
		api.Symbol(req.Symbol)
	}
	if req.OrderId != "" {
		api.OrderId(req.OrderId)
	}
	if req.ClientOrderId != "" {
		api.OrderId(req.ClientOrderId)
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

// Classic 杠杆 - 逐仓 - 历史委托
func (e *BitgetTradeEngine) apiClassicMarginIsolatedQueryOrders(req *QueryOrderParam) *mybitgetapi.PrivateRestClassicMarginIsolatedTradeHistoryOrdersAPI {
	api := e.bitgetPrivateClient().NewPrivateRestClassicMarginIsolatedTradeHistoryOrders()
	if req.Symbol != "" {
		api.Symbol(req.Symbol)
	}
	if req.OrderId != "" {
		api.OrderId(req.OrderId)
	}
	if req.ClientOrderId != "" {
		api.ClientOid(req.ClientOrderId)
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

// Classic 杠杆 - 全仓 - 历史委托
func (e *BitgetTradeEngine) apiClassicMarginCrossedQueryOrders(req *QueryOrderParam) *mybitgetapi.PrivateRestClassicMarginCrossTradeHistoryOrdersAPI {
	api := e.bitgetPrivateClient().NewPrivateRestClassicMarginCrossTradeHistoryOrders()
	if req.Symbol != "" {
		api.Symbol(req.Symbol)
	}
	if req.OrderId != "" {
		api.OrderId(req.OrderId)
	}
	if req.ClientOrderId != "" {
		api.ClientOid(req.ClientOrderId)
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

// Classic 现货 — 成交明细
func (e *BitgetTradeEngine) apiClassicSpotQueryTrades(req *QueryTradeParam) *mybitgetapi.PrivateRestClassicSpotTradeFillsAPI {
	api := e.bitgetPrivateClient().NewPrivateRestClassicSpotTradeFills()
	if req.Symbol != "" {
		api.Symbol(req.Symbol)
	}
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

// Classic 杠杆 - 逐仓 - 成交明细
func (e *BitgetTradeEngine) apiClassicMarginIsolatedQueryTrades(req *QueryTradeParam) *mybitgetapi.PrivateRestClassicMarginIsolatedTradeFillsAPI {
	api := e.bitgetPrivateClient().NewPrivateRestClassicMarginIsolatedTradeFills()
	if req.Symbol != "" {
		api.Symbol(req.Symbol)
	}
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

// Classic 杠杆 - 全仓 - 成交明细
func (e *BitgetTradeEngine) apiClassicMarginCrossedQueryTrades(req *QueryTradeParam) *mybitgetapi.PrivateRestClassicMarginCrossTradeFillsAPI {
	api := e.bitgetPrivateClient().NewPrivateRestClassicMarginCrossTradeFills()
	if req.Symbol != "" {
		api.Symbol(req.Symbol)
	}
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

// Classic 合约 — 当前委托
func (e *BitgetTradeEngine) apiClassicFuturesQueryOpenOrders(productType string, req *QueryOrderParam) *mybitgetapi.PrivateRestClassicFuturesTradeOrdersPendingAPI {
	api := e.bitgetPrivateClient().NewPrivateRestClassicFuturesTradeOrdersPending().ProductType(productType)
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

// Classic 合约 — 订单详情
func (e *BitgetTradeEngine) apiClassicFuturesQueryOrder(symbol, productType string, req *QueryOrderParam) *mybitgetapi.PrivateRestClassicFuturesTradeOrderDetailAPI {
	api := e.bitgetPrivateClient().NewPrivateRestClassicFuturesTradeOrderDetail().Symbol(symbol).ProductType(productType)
	if req.OrderId != "" {
		api.OrderId(req.OrderId)
	}
	if req.ClientOrderId != "" {
		api.ClientOid(req.ClientOrderId)
	}
	return api
}

// Classic 合约 — 历史委托
func (e *BitgetTradeEngine) apiClassicFuturesQueryOrders(productType string, req *QueryOrderParam) *mybitgetapi.PrivateRestClassicFuturesTradeOrdersHistoryAPI {
	api := e.bitgetPrivateClient().NewPrivateRestClassicFuturesTradeOrdersHistory().ProductType(productType)
	if req.Symbol != "" {
		api.Symbol(req.Symbol)
	}
	if req.OrderId != "" {
		api.OrderId(req.OrderId)
	}
	if req.ClientOrderId != "" {
		api.ClientOid(req.ClientOrderId)
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

// Classic 合约 — 历史成交
func (e *BitgetTradeEngine) apiClassicFuturesQueryTrades(productType string, req *QueryTradeParam) *mybitgetapi.PrivateRestClassicFuturesTradeFillHistoryAPI {
	api := e.bitgetPrivateClient().NewPrivateRestClassicFuturesTradeFillHistory().ProductType(productType)
	if req.Symbol != "" {
		api.Symbol(req.Symbol)
	}
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

// ——— Classic 现货 ———

func (e *BitgetTradeEngine) apiClassicSpotOrderCreate(req *OrderParam) *mybitgetapi.PrivateRestClassicSpotTradePlaceOrderAPI {
	api := e.bitgetPrivateClient().NewPrivateRestClassicSpotTradePlaceOrder().
		Symbol(req.Symbol).
		Side(e.converter.ToBitgetOrderSide(req.OrderSide)).
		OrderType(e.converter.ToBitgetOrderType(req.OrderType)).
		Force(e.converter.ToBitgetTimeInForce(req.TimeInForce)).
		Size(req.Quantity.String())

	if req.ClientOrderId != "" {
		api.ClientOid(req.ClientOrderId)
	}

	if req.OrderType == ORDER_TYPE_LIMIT {
		api.Price(req.Price.String())
	}
	return api
}

func (e *BitgetTradeEngine) apiClassicSpotAmendOrder(req *OrderParam) *mybitgetapi.PrivateRestClassicSpotTradeCancelReplaceOrderAPI {
	api := e.bitgetPrivateClient().NewPrivateRestClassicSpotTradeCancelReplaceOrder().
		Symbol(req.Symbol).
		Price(req.Price.String()).
		Size(req.Quantity.String())
	if req.OrderId != "" {
		api.OrderId(req.OrderId)
	}
	if req.ClientOrderId != "" {
		api.ClientOid(req.ClientOrderId)
	}
	if req.NewClientOrderId != "" {
		api.NewClientOid(req.NewClientOrderId)
	}
	return api
}

func (e *BitgetTradeEngine) apiClassicSpotCancelOrder(req *OrderParam) *mybitgetapi.PrivateRestClassicSpotTradeCancelOrderAPI {
	api := e.bitgetPrivateClient().NewPrivateRestClassicSpotTradeCancelOrder().Symbol(req.Symbol)
	if req.OrderId != "" {
		api.OrderId(req.OrderId)
	}
	if req.ClientOrderId != "" {
		api.ClientOid(req.ClientOrderId)
	}
	return api
}

func (e *BitgetTradeEngine) apiClassicSpotBatchCreateOrders(reqs []*OrderParam, multiple bool) *mybitgetapi.PrivateRestClassicSpotTradeBatchOrdersAPI {
	api := e.bitgetPrivateClient().NewPrivateRestClassicSpotTradeBatchOrders()
	if multiple {
		api.BatchMode("multiple")
	} else if len(reqs) > 0 {
		api.Symbol(reqs[0].Symbol).BatchMode("single")
	}
	for _, req := range reqs {
		item := mybitgetapi.PrivateRestClassicSpotTradeBatchOrdersReqItem{
			Side:      GetPointer(e.converter.ToBitgetOrderSide(req.OrderSide)),
			OrderType: GetPointer(e.converter.ToBitgetOrderType(req.OrderType)),
			Force:     GetPointer(e.converter.ToBitgetTimeInForce(req.TimeInForce)),
			Size:      GetPointer(req.Quantity.String()),
			ClientOid: GetPointer(req.ClientOrderId),
		}
		if multiple {
			item.Symbol = GetPointer(req.Symbol)
		}
		if req.OrderType == ORDER_TYPE_LIMIT {
			item.Price = GetPointer(req.Price.String())
		}
		api.AddOrder(item)
	}
	return api
}

func (e *BitgetTradeEngine) apiClassicSpotBatchCancelOrders(reqs []*OrderParam) *mybitgetapi.PrivateRestClassicSpotTradeBatchCancelOrderAPI {
	api := e.bitgetPrivateClient().NewPrivateRestClassicSpotTradeBatchCancelOrder()
	multiple := false
	if len(reqs) > 0 {
		for _, r := range reqs[1:] {
			if r.Symbol != reqs[0].Symbol {
				multiple = true
				break
			}
		}
	}
	if multiple {
		api.BatchMode("multiple")
	} else if len(reqs) > 0 {
		api.Symbol(reqs[0].Symbol).BatchMode("single")
	}
	for _, r := range reqs {
		item := mybitgetapi.PrivateRestClassicSpotTradeBatchCancelOrderItem{}
		if multiple {
			item.Symbol = GetPointer(r.Symbol)
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

// ——— Classic 合约 ———

func (e *BitgetTradeEngine) apiClassicFuturesOrderCreate(req *OrderParam, marginCoin string) *mybitgetapi.PrivateRestClassicFuturesTradePlaceOrderAPI {
	api := e.bitgetPrivateClient().NewPrivateRestClassicFuturesTradePlaceOrder().
		Symbol(req.Symbol).
		ProductType(req.AccountType).
		MarginMode(bitgetClassicFuturesMarginMode(req)).
		MarginCoin(marginCoin).
		Size(req.Quantity.String()).
		Side(e.converter.ToBitgetOrderSide(req.OrderSide)).
		OrderType(e.converter.ToBitgetOrderType(req.OrderType)).
		Force(e.converter.ToBitgetTimeInForce(req.TimeInForce))

	if req.ClientOrderId != "" {
		api.ClientOid(req.ClientOrderId)
	}
	if req.TradeSide != "" {
		api.TradeSide(req.TradeSide)
	}
	if req.OrderType == ORDER_TYPE_LIMIT {
		api.Price(req.Price.String())
	}
	if req.ReduceOnly {
		api.ReduceOnly("YES")
	} else {
		api.ReduceOnly("NO")
	}
	return api
}

func (e *BitgetTradeEngine) apiClassicFuturesAmendOrder(req *OrderParam) *mybitgetapi.PrivateRestClassicFuturesTradeModifyOrderAPI {
	api := e.bitgetPrivateClient().NewPrivateRestClassicFuturesTradeModifyOrder().
		Symbol(req.Symbol).
		ProductType(req.AccountType)

	if req.OrderId != "" {
		api.OrderId(req.OrderId)
	}
	if req.ClientOrderId != "" {
		api.ClientOid(req.ClientOrderId)
	}
	if req.NewClientOrderId != "" {
		api.NewClientOid(req.NewClientOrderId)
	}
	api.NewSize(req.Quantity.String())
	if req.OrderType == ORDER_TYPE_LIMIT {
		api.NewPrice(req.Price.String())
	}
	return api
}

func (e *BitgetTradeEngine) apiClassicFuturesCancelOrder(req *OrderParam) *mybitgetapi.PrivateRestClassicFuturesTradeCancelOrderAPI {
	api := e.bitgetPrivateClient().NewPrivateRestClassicFuturesTradeCancelOrder().
		Symbol(req.Symbol).
		ProductType(req.AccountType)
	if req.OrderId != "" {
		api.OrderId(req.OrderId)
	}
	if req.ClientOrderId != "" {
		api.ClientOid(req.ClientOrderId)
	}
	return api
}

func (e *BitgetTradeEngine) apiClassicFuturesBatchCreateOrders(reqs []*OrderParam, marginCoin string) *mybitgetapi.PrivateRestClassicFuturesTradeBatchOrderAPI {
	api := e.bitgetPrivateClient().NewPrivateRestClassicFuturesTradeBatchOrder()
	if len(reqs) == 0 {
		return api
	}

	r0 := reqs[0]
	api.Symbol(r0.Symbol).
		ProductType(r0.AccountType).
		MarginCoin(marginCoin).
		MarginMode(bitgetClassicFuturesMarginMode(r0))

	for _, r := range reqs {
		item := mybitgetapi.PrivateRestClassicFuturesTradeBatchOrderItem{
			Size:      GetPointer(r.Quantity.String()),
			Side:      GetPointer(e.converter.ToBitgetOrderSide(r.OrderSide)),
			OrderType: GetPointer(e.converter.ToBitgetOrderType(r.OrderType)),
			Force:     GetPointer(e.converter.ToBitgetTimeInForce(r.TimeInForce)),
			ClientOid: GetPointer(r.ClientOrderId),
		}
		if r.TradeSide != "" {
			item.TradeSide = GetPointer(r.TradeSide)
		}
		if r.OrderType == ORDER_TYPE_LIMIT {
			item.Price = GetPointer(r.Price.String())
		}
		if r.ReduceOnly {
			item.ReduceOnly = GetPointer("YES")
		} else {
			item.ReduceOnly = GetPointer("NO")
		}
		api.AddOrder(item)
	}
	return api
}

func (e *BitgetTradeEngine) apiClassicFuturesBatchCancelOrders(reqs []*OrderParam, marginCoin string) *mybitgetapi.PrivateRestClassicFuturesTradeBatchCancelOrderAPI {
	api := e.bitgetPrivateClient().NewPrivateRestClassicFuturesTradeBatchCancelOrder()
	if len(reqs) == 0 {
		return api
	}

	r0 := reqs[0]
	api.Symbol(r0.Symbol).
		ProductType(r0.AccountType).
		MarginCoin(marginCoin)

	for _, r := range reqs {
		item := mybitgetapi.PrivateRestClassicFuturesTradeCancelOrderIdItem{}
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

// ——— Classic 全仓杠杆 ———

func (e *BitgetTradeEngine) apiClassicMarginCrossCreateOrder(req *OrderParam) *mybitgetapi.PrivateRestClassicMarginCrossTradePlaceOrderAPI {
	lt := bitgetLoanTypeFromReq(req)
	api := e.bitgetPrivateClient().NewPrivateRestClassicMarginCrossTradePlaceOrder().
		Symbol(req.Symbol).
		OrderType(e.converter.ToBitgetOrderType(req.OrderType)).
		LoanType(lt).
		Force(e.converter.ToBitgetTimeInForce(req.TimeInForce)).
		Side(e.converter.ToBitgetOrderSide(req.OrderSide))

	if req.ClientOrderId != "" {
		api.ClientOid(req.ClientOrderId)
	}
	if req.OrderType == ORDER_TYPE_LIMIT {
		api.Price(req.Price.String())
		api.BaseSize(req.Quantity.String())
	} else {
		if req.OrderSide == ORDER_SIDE_BUY {
			api.QuoteSize(req.Quantity.String())
		} else {
			api.BaseSize(req.Quantity.String())
		}
	}
	return api
}

func (e *BitgetTradeEngine) apiClassicMarginCrossedOrderCreate(req *OrderParam) *mybitgetapi.PrivateRestClassicMarginCrossTradePlaceOrderAPI {
	return e.apiClassicMarginCrossCreateOrder(req)
}

func (e *BitgetTradeEngine) apiClassicMarginCrossCancelOrder(req *OrderParam) *mybitgetapi.PrivateRestClassicMarginCrossTradeCancelOrderAPI {
	api := e.bitgetPrivateClient().NewPrivateRestClassicMarginCrossTradeCancelOrder().Symbol(req.Symbol)
	if req.OrderId != "" {
		api.OrderId(req.OrderId)
	}
	if req.ClientOrderId != "" {
		api.ClientOid(req.ClientOrderId)
	}
	return api
}

func (e *BitgetTradeEngine) apiClassicMarginCrossBatchCreateOrders(reqs []*OrderParam) *mybitgetapi.PrivateRestClassicMarginCrossTradeBatchPlaceOrderAPI {
	api := e.bitgetPrivateClient().NewPrivateRestClassicMarginCrossTradeBatchPlaceOrder().Symbol(reqs[0].Symbol)
	lt := bitgetLoanTypeFromReq(reqs[0])
	for _, r := range reqs {
		item := mybitgetapi.PrivateRestClassicMarginCrossTradeBatchOrderItem{
			OrderType: GetPointer(e.converter.ToBitgetOrderType(r.OrderType)),
			LoanType:  GetPointer(lt),
			Force:     GetPointer(e.converter.ToBitgetTimeInForce(r.TimeInForce)),
			Side:      GetPointer(e.converter.ToBitgetOrderSide(r.OrderSide)),
			ClientOid: GetPointer(r.ClientOrderId),
		}
		if r.OrderType == ORDER_TYPE_LIMIT {
			item.Price = GetPointer(r.Price.String())
			item.BaseSize = GetPointer(r.Quantity.String())
		} else {
			if r.OrderSide == ORDER_SIDE_BUY {
				item.QuoteSize = GetPointer(r.Quantity.String())
			} else {
				item.BaseSize = GetPointer(r.Quantity.String())
			}
		}
		api.AddOrder(item)
	}
	return api
}

func (e *BitgetTradeEngine) apiClassicMarginCrossBatchCancelOrders(reqs []*OrderParam) *mybitgetapi.PrivateRestClassicMarginCrossTradeBatchCancelOrderAPI {
	api := e.bitgetPrivateClient().NewPrivateRestClassicMarginCrossTradeBatchCancelOrder().Symbol(reqs[0].Symbol)
	for _, r := range reqs {
		item := mybitgetapi.PrivateRestClassicMarginCrossTradeCancelOrderItem{}
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

// ——— Classic 逐仓杠杆 ———

func (e *BitgetTradeEngine) apiClassicMarginIsolatedCreateOrder(req *OrderParam) *mybitgetapi.PrivateRestClassicMarginIsolatedTradePlaceOrderAPI {
	lt := bitgetLoanTypeFromReq(req)
	api := e.bitgetPrivateClient().NewPrivateRestClassicMarginIsolatedTradePlaceOrder().
		Symbol(req.Symbol).
		OrderType(e.converter.ToBitgetOrderType(req.OrderType)).
		LoanType(lt).
		Force(e.converter.ToBitgetTimeInForce(req.TimeInForce)).
		Side(e.converter.ToBitgetOrderSide(req.OrderSide))

	if req.ClientOrderId != "" {
		api.ClientOid(req.ClientOrderId)
	}
	if req.OrderType == ORDER_TYPE_LIMIT {
		api.Price(req.Price.String())
		api.BaseSize(req.Quantity.String())
	} else {
		if req.OrderSide == ORDER_SIDE_BUY {
			api.QuoteSize(req.Quantity.String())
		} else {
			api.BaseSize(req.Quantity.String())
		}
	}
	return api
}

func (e *BitgetTradeEngine) apiClassicMarginIsolatedOrderCreate(req *OrderParam) *mybitgetapi.PrivateRestClassicMarginIsolatedTradePlaceOrderAPI {
	return e.apiClassicMarginIsolatedCreateOrder(req)
}

func (e *BitgetTradeEngine) apiClassicMarginIsolatedCancelOrder(req *OrderParam) *mybitgetapi.PrivateRestClassicMarginIsolatedTradeCancelOrderAPI {
	api := e.bitgetPrivateClient().NewPrivateRestClassicMarginIsolatedTradeCancelOrder().Symbol(req.Symbol)
	if req.OrderId != "" {
		api.OrderId(req.OrderId)
	}
	if req.ClientOrderId != "" {
		api.ClientOid(req.ClientOrderId)
	}
	return api
}

func (e *BitgetTradeEngine) apiClassicMarginIsolatedBatchCreateOrders(reqs []*OrderParam) *mybitgetapi.PrivateRestClassicMarginIsolatedTradeBatchPlaceOrderAPI {
	api := e.bitgetPrivateClient().NewPrivateRestClassicMarginIsolatedTradeBatchPlaceOrder().Symbol(reqs[0].Symbol)
	lt := bitgetLoanTypeFromReq(reqs[0])
	for _, r := range reqs {
		item := mybitgetapi.PrivateRestClassicMarginIsolatedTradeBatchOrderItem{
			OrderType: GetPointer(e.converter.ToBitgetOrderType(r.OrderType)),
			LoanType:  GetPointer(lt),
			Force:     GetPointer(e.converter.ToBitgetTimeInForce(r.TimeInForce)),
			Side:      GetPointer(e.converter.ToBitgetOrderSide(r.OrderSide)),
			ClientOid: GetPointer(r.ClientOrderId),
		}
		if r.OrderType == ORDER_TYPE_LIMIT {
			item.Price = GetPointer(r.Price.String())
			item.BaseSize = GetPointer(r.Quantity.String())
		} else {
			if r.OrderSide == ORDER_SIDE_BUY {
				item.QuoteSize = GetPointer(r.Quantity.String())
			} else {
				item.BaseSize = GetPointer(r.Quantity.String())
			}
		}
		api.AddOrder(item)
	}
	return api
}

func (e *BitgetTradeEngine) apiClassicMarginIsolatedBatchCancelOrders(reqs []*OrderParam) *mybitgetapi.PrivateRestClassicMarginIsolatedTradeBatchCancelOrderAPI {
	api := e.bitgetPrivateClient().NewPrivateRestClassicMarginIsolatedTradeBatchCancelOrder().Symbol(reqs[0].Symbol)
	for _, r := range reqs {
		item := mybitgetapi.PrivateRestClassicMarginIsolatedTradeCancelOrderItem{}
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

func bitgetClassicFuturesBatchPreCheck(reqs []*OrderParam) error {
	if len(reqs) == 0 {
		return errors.New("empty batch")
	}
	r0 := reqs[0]
	for _, r := range reqs {
		if r.AccountType != r0.AccountType || r.Symbol != r0.Symbol || r.IsIsolated != r0.IsIsolated {
			return ErrorInvalid("classic futures batch requires same accountType, symbol and margin mode")
		}
	}
	mc0 := bitgetMarginCoinFromSymbol(r0.Symbol, r0.Ccy)
	for _, r := range reqs {
		if bitgetMarginCoinFromSymbol(r.Symbol, r.Ccy) != mc0 {
			return ErrorInvalid("classic futures batch requires same margin coin")
		}
	}
	if mc0 == "" {
		return ErrorInvalid("margin coin required: set OrderParam.Ccy or use symbol suffix USDT/USDC")
	}
	return nil
}

func bitgetClassicMarginBatchPreCheck(reqs []*OrderParam) error {
	if len(reqs) == 0 {
		return errors.New("empty batch")
	}
	sym := reqs[0].Symbol
	for _, r := range reqs {
		if r.Symbol != sym {
			return ErrorInvalid("classic margin batch requires same symbol")
		}
	}
	return nil
}
