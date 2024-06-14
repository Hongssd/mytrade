package mytrade

import (
	"errors"
	"github.com/Hongssd/mybybitapi"
	"github.com/shopspring/decimal"
)

// 查询订单接口获取
func (b *BybitTradeEngine) apiQueryOpenOrders(req *QueryOrderParam, pageCursor string) *mybybitapi.OrderRealtimeAPI {
	client := mybybitapi.NewRestClient(b.apiKey, b.secretKey).PrivateRestClient()
	api := client.NewOrderRealtime().Category(req.AccountType).Symbol(req.Symbol).Limit(50)
	if pageCursor != "" {
		api.Cursor(pageCursor)
	}
	return api
}
func (b *BybitTradeEngine) apiQueryOrder(req *QueryOrderParam) *mybybitapi.OrderRealtimeAPI {
	client := mybybitapi.NewRestClient(b.apiKey, b.secretKey).PrivateRestClient()
	api := client.NewOrderRealtime().Category(req.AccountType).Symbol(req.Symbol)
	if req.OrderId != "" {
		api.OrderId(req.OrderId)
	}
	if req.ClientOrderId != "" {
		api.OrderLinkId(req.ClientOrderId)
	}
	return api
}
func (b *BybitTradeEngine) apiQueryTrades(req *QueryTradeParam, pageCursor string) *mybybitapi.OrderExecutionListAPI {
	client := mybybitapi.NewRestClient(b.apiKey, b.secretKey).PrivateRestClient()
	api := client.NewOrderExecutionList().Category(req.AccountType).Symbol(req.Symbol)
	if req.OrderId != "" {
		api.OrderId(req.OrderId)
	}
	if req.StartTime != 0 {
		api.StartTime(req.StartTime)
	}
	if req.EndTime != 0 {
		api.EndTime(req.EndTime)
	}
	if req.Limit != 0 {
		api.Limit(req.Limit)
	}

	if pageCursor != "" {
		api.Cursor(pageCursor)
	}

	return api
}

// 单订单接口获取
func (b *BybitTradeEngine) apiOrderCreate(req *OrderParam) *mybybitapi.OrderCreateAPI {
	client := mybybitapi.NewRestClient(b.apiKey, b.secretKey).PrivateRestClient()
	api := client.NewOrderCreate().
		Category(req.AccountType).
		Symbol(req.Symbol).
		Side(b.bybitConverter.ToBYBITOrderSide(req.OrderSide)).
		OrderType(b.bybitConverter.ToBYBITOrderType(req.OrderType)).
		Price(req.Price.String()).
		Qty(req.Quantity.String())

	if req.PositionSide != "" {
		api.PositionIdx(b.bybitConverter.ToBYBITPositionSide(req.OrderSide, req.PositionSide))
	}
	if req.ClientOrderId != "" {
		api.OrderLinkId(req.ClientOrderId)
	}
	if req.TimeInForce != "" {
		api.TimeInForce(b.bybitConverter.ToBYBITTimeInForce(req.TimeInForce))
	}
	if req.ReduceOnly {
		api.ReduceOnly(req.ReduceOnly)
	}

	return api
}
func (b *BybitTradeEngine) apiOrderAmend(req *OrderParam) *mybybitapi.OrderAmendAPI {
	client := mybybitapi.NewRestClient(b.apiKey, b.secretKey).PrivateRestClient()
	api := client.NewOrderAmend().
		Category(req.AccountType).
		Symbol(req.Symbol)

	if req.OrderId != "" {
		api.OrderId(req.OrderId)
	}
	if req.ClientOrderId != "" {
		api.OrderLinkId(req.ClientOrderId)
	}

	if req.Price != decimal.Zero {
		api.Price(req.Price.String())
	}
	if req.Quantity != decimal.Zero {
		api.Qty(req.Quantity.String())
	}
	return api
}
func (b *BybitTradeEngine) apiOrderCancel(req *OrderParam) *mybybitapi.OrderCancelAPI {
	client := mybybitapi.NewRestClient(b.apiKey, b.secretKey).PrivateRestClient()
	api := client.NewOrderCancel().
		Category(req.AccountType).
		Symbol(req.Symbol).
		OrderId(req.OrderId)

	if req.ClientOrderId != "" {
		api.OrderLinkId(req.ClientOrderId)

	}

	return api
}

// 批量接口订单获取
func (b *BybitTradeEngine) apiBatchOrderCreate(reqs []*OrderParam) *mybybitapi.OrderCreateBatchAPI {
	client := mybybitapi.NewRestClient(b.apiKey, b.secretKey).PrivateRestClient()
	api := client.NewOrderCreateBatch()
	for _, req := range reqs {
		api.AddNewOrderCreateReq(b.apiOrderCreate(req))
	}
	return api
}
func (b *BybitTradeEngine) apiBatchOrderAmend(reqs []*OrderParam) *mybybitapi.OrderAmendBatchAPI {
	client := mybybitapi.NewRestClient(b.apiKey, b.secretKey).PrivateRestClient()
	api := client.NewOrderAmendBatch()
	for _, req := range reqs {
		api.AddNewOrderAmendReq(b.apiOrderAmend(req))
	}
	return api
}
func (b *BybitTradeEngine) apiBatchOrderCancel(reqs []*OrderParam) *mybybitapi.OrderCancelBatchAPI {
	client := mybybitapi.NewRestClient(b.apiKey, b.secretKey).PrivateRestClient()
	api := client.NewOrderCancelBatch()
	for _, req := range reqs {
		api.AddNewOrderCancelReq(b.apiOrderCancel(req))
	}
	return api
}

// 查询订单结果处理
func (b *BybitTradeEngine) handleOrdersFromQueryOpenOrders(req *QueryOrderParam, res mybybitapi.OrderRealtimeRes) []*Order {
	var orders []*Order
	for _, order := range res.List {
		orders = append(orders, &Order{
			Exchange:      BYBIT_NAME.String(),
			AccountType:   req.AccountType,
			Symbol:        req.Symbol,
			OrderId:       order.OrderId,
			ClientOrderId: order.OrderLinkId,
			Price:         order.Price,
			Quantity:      order.Qty,
			ExecutedQty:   order.CumExecQty,
			CumQuoteQty:   order.CumExecValue,
			AvgPrice:      order.AvgPrice,
			Status:        b.bybitConverter.FromBYBITOrderStatus(order.OrderStatus),
			Type:          b.bybitConverter.FromBYBITOrderType(order.OrderType),
			Side:          b.bybitConverter.FromBYBITOrderSide(order.Side),
			PositionSide:  b.bybitConverter.FromBYBITPositionSide(order.PositionIdx),
			TimeInForce:   b.bybitConverter.FromBYBITTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    stringToInt64(order.CreatedTime),
			UpdateTime:    stringToInt64(order.UpdatedTime),

			FeeAmount: order.CumExecFee,
			FeeCcy:    "",
		})

	}
	return orders
}
func (b *BybitTradeEngine) handleTradesFromQueryTrades(req *QueryTradeParam, res mybybitapi.OrderExecutionListRes) []*Trade {
	var trades []*Trade
	for _, r := range res.List {
		quoteQty := decimal.RequireFromString(r.ExecPrice).Mul(decimal.RequireFromString(r.ExecQty))
		trades = append(trades, &Trade{
			Exchange:     BYBIT_NAME.String(),
			AccountType:  req.AccountType,
			Symbol:       req.Symbol,
			TradeId:      r.ExecId,
			OrderId:      r.OrderId,
			Price:        r.ExecPrice,
			Quantity:     r.ExecQty,
			QuoteQty:     quoteQty.String(),
			Side:         b.bybitConverter.FromBYBITOrderSide(r.Side),
			PositionSide: "",
			FeeAmount:    r.ExecFee,
			FeeCcy:       r.FeeCurrency,
			RealizedPnl:  "",
			IsMaker:      r.IsMaker,
			Timestamp:    stringToInt64(r.ExecTime),
		})
	}
	return trades
}

// 批量订单返回结果处理
func (b *BybitTradeEngine) handleOrderFromBatchOrderCreate(reqs []*OrderParam, res *mybybitapi.BybitRestRes[mybybitapi.OrderCreateBatchRes]) ([]*Order, error) {
	if len(res.Result.List) != len(reqs) {
		return nil, errors.New("api return invalid data")
	}
	orders := make([]*Order, 0, len(reqs))
	for _, r := range res.Result.List {
		order := &Order{
			Exchange:      BYBIT_NAME.String(),
			OrderId:       r.OrderId,
			ClientOrderId: r.OrderLinkId,
			AccountType:   r.Category,
			Symbol:        r.Symbol,
			CreateTime:    stringToInt64(r.CreateAt),
		}
		orders = append(orders, order)
	}
	return orders, nil
}
func (b *BybitTradeEngine) handleOrderFromBatchOrderAmend(reqs []*OrderParam, res *mybybitapi.BybitRestRes[mybybitapi.OrderAmendBatchRes]) ([]*Order, error) {
	if len(res.Result.List) != len(reqs) {
		return nil, errors.New("api return invalid data")
	}
	orders := make([]*Order, 0, len(reqs))
	for _, r := range res.Result.List {
		order := &Order{
			Exchange:      BYBIT_NAME.String(),
			OrderId:       r.OrderId,
			ClientOrderId: r.OrderLinkId,
			AccountType:   r.Category,
			Symbol:        r.Symbol,
		}
		orders = append(orders, order)
	}
	return orders, nil
}
func (b *BybitTradeEngine) handleOrderFromBatchOrderCancel(reqs []*OrderParam, res *mybybitapi.BybitRestRes[mybybitapi.OrderCancelBatchRes]) ([]*Order, error) {
	if len(res.Result.List) != len(reqs) {
		return nil, errors.New("api return invalid data")
	}
	orders := make([]*Order, 0, len(reqs))
	for _, r := range res.Result.List {
		order := &Order{
			Exchange:      BYBIT_NAME.String(),
			OrderId:       r.OrderId,
			ClientOrderId: r.OrderLinkId,
			AccountType:   r.Category,
			Symbol:        r.Symbol,
		}
		orders = append(orders, order)
	}
	return orders, nil
}

// 订单推送处理
func (b *BybitTradeEngine) handleOrderFromWsOrder(orders mybybitapi.WsOrder) []*Order {
	// 从ws订单信息转换为本地订单信息
	var res []*Order
	for _, order := range orders.Data {
		order := &Order{
			Exchange:      BYBIT_NAME.String(),
			AccountType:   order.Category,
			Symbol:        order.Symbol,
			OrderId:       order.OrderId,
			ClientOrderId: order.OrderLinkId,
			Price:         order.Price,
			Quantity:      order.Qty,
			ExecutedQty:   order.CumExecQty,
			CumQuoteQty:   order.CumExecValue,
			AvgPrice:      order.AvgPrice,
			Status:        b.bybitConverter.FromBYBITOrderStatus(order.OrderStatus),
			Type:          b.bybitConverter.FromBYBITOrderType(order.OrderType),
			Side:          b.bybitConverter.FromBYBITOrderSide(order.Side),
			PositionSide:  b.bybitConverter.FromBYBITPositionSide(order.PositionIdx),
			TimeInForce:   b.bybitConverter.FromBYBITTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    stringToInt64(order.CreatedTime),
			UpdateTime:    stringToInt64(order.UpdatedTime),

			FeeAmount: order.CumExecFee,
			FeeCcy:    order.FeeCurrency,
		}
		res = append(res, order)
	}
	return res
}

func (b *BybitTradeEngine) handleOrderFromInverseBatchErr(req *OrderParam, err error) *Order {
	return &Order{
		Exchange:      BYBIT_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
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

func (b *BybitTradeEngine) restBatchPreCheck(reqs []*OrderParam) error {
	//检测长度，BYBIT最多批量下10个订单
	if len(reqs) > 10 {
		return ErrorInvalid("bybit order param length require less than 10")
	}

	//检测类型是否相同
	for _, req := range reqs {
		if err := b.accountTypePreCheck(req.AccountType); err != nil {
			return err
		}
		if req.AccountType != reqs[0].AccountType {
			return ErrorInvalid("order param account type require same")
		}
	}
	return nil
}

func (b *BybitTradeEngine) accountTypePreCheck(accountType string) error {
	switch BybitAccountType(accountType) {
	case BYBIT_AC_SPOT, BYBIT_AC_LINEAR, BYBIT_AC_INVERSE:
	default:
		return ErrorAccountType
	}
	return nil
}

// ws订单请求前置检查
func (b *BybitTradeEngine) wsOrderPreCheck() (bool, error) {
	b.wsForOrderMu.Lock()
	defer b.wsForOrderMu.Unlock()

	if b.wsForOrder == nil {
		newWs := mybybitapi.NewTradeWsStreamClient()
		err := newWs.OpenConn()
		if err != nil {
			return false, err
		}

		err = newWs.Auth(mybybitapi.NewRestClient(b.apiKey, b.secretKey))
		if err != nil {
			return false, err
		}
		b.wsForOrder = newWs
	}
	return true, nil
}
