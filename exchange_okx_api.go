package mytrade

import (
	"errors"
	"fmt"
	"github.com/Hongssd/myokxapi"
	"github.com/shopspring/decimal"
)

// 查询订单接口获取
func (o *OkxTradeEngine) apiQueryOpenOrders(req *QueryOrderParam) *myokxapi.PrivateRestTradeOrdersPendingAPI {
	client := okx.NewRestClient(o.apiKey, o.secretKey, o.passphrase).PrivateRestClient()
	api := client.NewPrivateRestTradeOrdersPending().InstType(req.AccountType).InstId(req.Symbol)
	return api
}
func (o *OkxTradeEngine) apiQueryOrder(req *QueryOrderParam) *myokxapi.PrivateRestTradeOrderGetAPI {
	client := okx.NewRestClient(o.apiKey, o.secretKey, o.passphrase).PrivateRestClient()
	api := client.NewPrivateRestTradeOrderGet().InstId(req.Symbol)
	if req.OrderId != "" {
		api.OrdId(req.OrderId)
	}
	if req.ClientOrderId != "" {
		api.ClOrdId(req.ClientOrderId)
	}
	return api
}
func (o *OkxTradeEngine) apiQueryTrades(req *QueryTradeParam) *myokxapi.PrivateRestTradeFillsAPI {
	client := okx.NewRestClient(o.apiKey, o.secretKey, o.passphrase).PrivateRestClient()
	api := client.NewPrivateRestTradeFills().InstType(req.AccountType).InstId(req.Symbol)
	if req.OrderId != "" {
		api.OrdId(req.OrderId)
	}
	return api
}

// 单订单接口获取
func (o *OkxTradeEngine) apiOrderCreate(req *OrderParam) *myokxapi.PrivateRestTradeOrderPostAPI {
	client := okx.NewRestClient(o.apiKey, o.secretKey, o.passphrase).PrivateRestClient()

	tdMode := o.okxConverter.getTdModeFromAccountType(OkxAccountType(req.AccountType), req.IsIsolated)

	api := client.NewPrivateRestTradeOrderPost().
		InstId(req.Symbol).TdMode(tdMode).
		Side(o.okxConverter.ToOKXOrderSide(req.OrderSide)).
		OrdType(o.okxConverter.ToOKXOrderType(req.OrderType, req.TimeInForce)).
		Px(req.Price.String()).
		Sz(req.Quantity.String())

	if req.ClientOrderId != "" {
		api.ClOrdId(req.ClientOrderId)
	}

	return api
}
func (o *OkxTradeEngine) apiOrderAmend(req *OrderParam) *myokxapi.PrivateRestTradeAmendOrderAPI {
	client := okx.NewRestClient(o.apiKey, o.secretKey, o.passphrase).PrivateRestClient()

	api := client.NewPrivateRestTradeAmendOrder().
		InstId(req.Symbol)
	if req.OrderId != "" {
		api.OrdId(req.OrderId)
	}
	if req.ClientOrderId != "" {
		api.ClOrdId(req.ClientOrderId)
	}
	if !req.Price.IsZero() {
		api.NewPx(req.Price.String())
	}
	if !req.Quantity.IsZero() {
		api.NewSz(req.Quantity.String())
	}

	return api
}
func (o *OkxTradeEngine) apiOrderCancel(req *OrderParam) *myokxapi.PrivateRestTradeCancelOrderAPI {
	client := okx.NewRestClient(o.apiKey, o.secretKey, o.passphrase).PrivateRestClient()

	api := client.NewPrivateRestTradeCancelOrder().
		InstId(req.Symbol)
	if req.OrderId != "" {
		api.OrdId(req.OrderId)
	}
	if req.ClientOrderId != "" {
		api.ClOrdId(req.ClientOrderId)
	}

	return api
}

// 批量订单接口获取
func (o *OkxTradeEngine) apiBatchOrderCreate(reqs []*OrderParam) *myokxapi.PrivateRestTradeBatchOrdersAPI {
	client := okx.NewRestClient(o.apiKey, o.secretKey, o.passphrase).PrivateRestClient()
	api := client.NewPrivateRestTradeBatchOrders()
	for _, req := range reqs {
		api.AddNewOrderReq(o.apiOrderCreate(req))
	}
	return api
}
func (o *OkxTradeEngine) apiBatchOrderAmend(reqs []*OrderParam) *myokxapi.PrivateRestTradeAmendBatchOrdersAPI {
	client := okx.NewRestClient(o.apiKey, o.secretKey, o.passphrase).PrivateRestClient()
	api := client.NewPrivateRestTradeAmendBatchOrders()
	for _, req := range reqs {
		api.AddNewOrderReq(o.apiOrderAmend(req))
	}
	return api
}
func (o *OkxTradeEngine) apiBatchOrderCancel(reqs []*OrderParam) *myokxapi.PrivateRestTradeCancelBatchOrdersAPI {
	client := okx.NewRestClient(o.apiKey, o.secretKey, o.passphrase).PrivateRestClient()
	api := client.NewPrivateRestTradeCancelBatchOrders()
	for _, req := range reqs {
		api.AddNewOrderReq(o.apiOrderCancel(req))
	}
	return api
}

// 查询订单返回结果处理
func (o *OkxTradeEngine) handleOrdersFromQueryOpenOrders(req *QueryOrderParam, res *myokxapi.OkxRestRes[myokxapi.PrivateRestTradeOrdersPendingRes]) ([]*Order, error) {
	if res.Code != "0" {
		return nil, fmt.Errorf("[%s]:%s", res.Code, res.Msg)
	}
	orders := make([]*Order, 0, len(res.Data))
	for _, r := range res.Data {
		orderType, timeInForce := o.okxConverter.FromOKXOrderType(r.OrdType)
		order := &Order{
			Exchange:      OKX_NAME.String(),
			OrderId:       r.OrdId,
			ClientOrderId: r.ClOrdId,
			AccountType:   req.AccountType,
			Symbol:        req.Symbol,
			Price:         r.Px,
			Quantity:      r.Sz,
			ExecutedQty:   r.FillSz,
			Status:        o.okxConverter.FromOKXOrderStatus(r.State),
			Type:          orderType,
			Side:          o.okxConverter.FromOKXOrderSide(r.Side),
			PositionSide:  o.okxConverter.FromOKXPositionSide(r.PosSide),
			TimeInForce:   timeInForce,
			ReduceOnly:    stringToBool(r.ReduceOnly),
			CreateTime:    stringToInt64(r.CTime),
			UpdateTime:    stringToInt64(r.UTime),
			RealizedPnl:   r.Pnl,
		}
		orders = append(orders, order)
	}
	return orders, nil
}
func (o *OkxTradeEngine) handleOrderFromQueryOrderGet(req *QueryOrderParam, res *myokxapi.OkxRestRes[myokxapi.PrivateRestTradeOrderGetRes]) (*Order, error) {
	if res.Code != "0" {
		return nil, fmt.Errorf("[%s]:%s", res.Code, res.Msg)
	}
	if len(res.Data) != 1 {
		return nil, errors.New("api return invalid data")
	}
	r := res.Data[0]

	orderType, timeInForce := o.okxConverter.FromOKXOrderType(r.OrdType)

	order := &Order{
		Exchange:      OKX_NAME.String(),
		OrderId:       r.OrdId,
		ClientOrderId: r.ClOrdId,
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		Price:         r.Px,
		Quantity:      r.Sz,
		ExecutedQty:   r.FillSz,
		Status:        o.okxConverter.FromOKXOrderStatus(r.State),
		Type:          orderType,
		Side:          o.okxConverter.FromOKXOrderSide(r.Side),
		PositionSide:  o.okxConverter.FromOKXPositionSide(r.PosSide),
		TimeInForce:   timeInForce,
		ReduceOnly:    stringToBool(r.ReduceOnly),
		CreateTime:    stringToInt64(r.CTime),
		UpdateTime:    stringToInt64(r.UTime),
		RealizedPnl:   r.Pnl,
	}
	return order, nil
}
func (o *OkxTradeEngine) handleTradesFromQueryTrades(req *QueryTradeParam, res *myokxapi.OkxRestRes[myokxapi.PrivateRestTradeFillsRes]) ([]*Trade, error) {
	if res.Code != "0" {
		return nil, fmt.Errorf("[%s]:%s", res.Code, res.Msg)
	}
	trades := make([]*Trade, 0, len(res.Data))

	for _, r := range res.Data {
		quoteQty := decimal.RequireFromString(r.FillPx).Mul(decimal.RequireFromString(r.FillSz))
		isMaker := r.ExecType == "M"
		trade := &Trade{
			Exchange:     OKX_NAME.String(),
			AccountType:  req.AccountType,
			Symbol:       req.Symbol,
			TradeId:      r.TradeId,
			OrderId:      r.OrdId,
			Price:        r.FillPx,
			Quantity:     r.FillSz,
			QuoteQty:     quoteQty.String(),
			Side:         o.okxConverter.FromOKXOrderSide(r.Side),
			PositionSide: o.okxConverter.FromOKXPositionSide(r.PosSide),
			FeeAmount:    r.Fee,
			FeeCcy:       r.FeeCcy,
			RealizedPnl:  r.FillPnl,
			IsMaker:      isMaker,
			Timestamp:    stringToInt64(r.FillTime),
		}
		trades = append(trades, trade)
	}
	return trades, nil
}

// 单订单返回结果处理
func (o *OkxTradeEngine) handleOrderFromOrderCreate(req *OrderParam, res *myokxapi.OkxRestRes[myokxapi.PrivateRestTradeOrderPostRes]) (*Order, error) {

	if res.Code != "0" {
		return nil, fmt.Errorf("[%s]:%s", res.Code, res.Msg)
	}
	if len(res.Data) != 1 {
		return nil, errors.New("api return invalid data")
	}
	if res.Data[0].SCode != "0" {
		return nil, fmt.Errorf("[%s]:%s", res.Data[0].SCode, res.Data[0].SMsg)
	}
	r := res.Data[0]
	order := &Order{
		Exchange:      OKX_NAME.String(),
		OrderId:       r.OrdId,
		ClientOrderId: r.ClOrdId,
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
	}
	return order, nil
}
func (o *OkxTradeEngine) handleOrderFromOrderAmend(req *OrderParam, res *myokxapi.OkxRestRes[myokxapi.PrivateRestTradeAmendOrderRes]) (*Order, error) {

	if res.Code != "0" {
		return nil, fmt.Errorf("[%s]:%s", res.Code, res.Msg)
	}
	if len(res.Data) != 1 {
		return nil, errors.New("api return invalid data")
	}
	if res.Data[0].SCode != "0" {
		return nil, fmt.Errorf("[%s]:%s", res.Data[0].SCode, res.Data[0].SMsg)
	}
	r := res.Data[0]
	order := &Order{
		Exchange:      OKX_NAME.String(),
		OrderId:       r.OrdId,
		ClientOrderId: r.ClOrdId,
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
	}
	return order, nil
}
func (o *OkxTradeEngine) handleOrderFromOrderCancel(req *OrderParam, res *myokxapi.OkxRestRes[myokxapi.PrivateRestTradeCancelOrderRes]) (*Order, error) {

	if res.Code != "0" {
		return nil, fmt.Errorf("[%s]:%s", res.Code, res.Msg)
	}
	if len(res.Data) != 1 {
		return nil, errors.New("api return invalid data")
	}
	if res.Data[0].SCode != "0" {
		return nil, fmt.Errorf("[%s]:%s", res.Data[0].SCode, res.Data[0].SMsg)
	}
	r := res.Data[0]
	order := &Order{
		Exchange:      OKX_NAME.String(),
		OrderId:       r.OrdId,
		ClientOrderId: r.ClOrdId,
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
	}
	return order, nil
}

// 批量订单返回结果处理
func (o *OkxTradeEngine) handleOrderFromBatchOrderCreate(reqs []*OrderParam, res *myokxapi.OkxRestRes[myokxapi.PrivateRestTradeBatchOrdersRes]) ([]*Order, error) {
	if res.Code != "0" {
		return nil, fmt.Errorf("API ERRO[%s]:%s", res.Code, res.Msg)
	}
	if len(res.Data) != len(reqs) {
		return nil, errors.New("api return invalid data")
	}
	orders := make([]*Order, 0, len(reqs))
	for i, r := range res.Data {
		if r.SCode != "0" {
			return nil, fmt.Errorf("ORDE ERRO[%s]:%s", r.SCode, r.SMsg)
		}
		order := &Order{
			Exchange:      OKX_NAME.String(),
			OrderId:       r.OrdId,
			ClientOrderId: r.ClOrdId,
			AccountType:   reqs[i].AccountType,
			Symbol:        reqs[i].Symbol,
		}
		orders = append(orders, order)
	}
	return orders, nil
}
func (o *OkxTradeEngine) handleOrderFromBatchOrderAmend(reqs []*OrderParam, res *myokxapi.OkxRestRes[myokxapi.PrivateRestTradeAmendBatchOrdersRes]) ([]*Order, error) {
	if res.Code != "0" {
		return nil, fmt.Errorf("[%s]:%s", res.Code, res.Msg)
	}
	if len(res.Data) != len(reqs) {
		return nil, errors.New("api return invalid data")
	}
	orders := make([]*Order, 0, len(reqs))
	for i, r := range res.Data {
		if r.SCode != "0" {
			return nil, fmt.Errorf("[%s]:%s", r.SCode, r.SMsg)
		}
		order := &Order{
			Exchange:      OKX_NAME.String(),
			OrderId:       r.OrdId,
			ClientOrderId: r.ClOrdId,
			AccountType:   reqs[i].AccountType,
			Symbol:        reqs[i].Symbol,
		}
		orders = append(orders, order)
	}
	return orders, nil
}
func (o *OkxTradeEngine) handleOrderFromBatchOrderCancel(reqs []*OrderParam, res *myokxapi.OkxRestRes[myokxapi.PrivateRestTradeCancelBatchOrdersRes]) ([]*Order, error) {
	if res.Code != "0" {
		return nil, fmt.Errorf("[%s]:%s", res.Code, res.Msg)
	}
	if len(res.Data) != len(reqs) {
		return nil, errors.New("api return invalid data")
	}
	orders := make([]*Order, 0, len(reqs))
	for i, r := range res.Data {
		if r.SCode != "0" {
			return nil, fmt.Errorf("[%s]:%s", r.SCode, r.SMsg)
		}
		order := &Order{
			Exchange:      OKX_NAME.String(),
			OrderId:       r.OrdId,
			ClientOrderId: r.ClOrdId,
			AccountType:   reqs[i].AccountType,
			Symbol:        reqs[i].Symbol,
		}
		orders = append(orders, order)
	}
	return orders, nil
}

// 订单推送处理
func (o *OkxTradeEngine) handleOrderFromWsOrder(order myokxapi.WsOrders) *Order {

	orderType, timeInForce := o.okxConverter.FromOKXOrderType(order.OrdType)

	avgPx := decimal.RequireFromString(order.AvgPx)
	cumQuoteQty := decimal.Zero
	if !avgPx.IsZero() {
		cumQuoteQty = avgPx.Mul(decimal.RequireFromString(order.FillSz))
	}

	return &Order{
		Exchange:      OKX_NAME.String(),
		Symbol:        order.Orders.InstId,
		OrderId:       order.OrdId,
		ClientOrderId: order.ClOrdId,
		Price:         order.Px,
		Quantity:      order.Sz,
		ExecutedQty:   order.FillSz,
		CumQuoteQty:   cumQuoteQty.String(),
		Status:        o.okxConverter.FromOKXOrderStatus(order.State),
		Type:          orderType,
		Side:          o.okxConverter.FromOKXOrderSide(order.Side),
		PositionSide:  o.okxConverter.FromOKXPositionSide(order.PosSide),
		TimeInForce:   timeInForce,
		ReduceOnly:    stringToBool(order.ReduceOnly),
		CreateTime:    stringToInt64(order.CTime),
		UpdateTime:    stringToInt64(order.UTime),
		RealizedPnl:   order.FillPnl,

		ErrorMsg:  order.Msg,
		ErrorCode: order.Code,
	}

}
func (o *OkxTradeEngine) getBoardcastFromAccountType(accountType string) **okxOrderBroadcaster {
	switch OkxAccountType(accountType) {
	case OKX_AC_SPOT:
		return &o.broadcasterSpot
	case OKX_AC_SWAP:
		return &o.broadcasterSwap
	case OKX_AC_FUTURES:
		return &o.broadcasterFuture
	default:
		return nil
	}
}
