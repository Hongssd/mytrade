package mytrade

import (
	"github.com/Hongssd/mybinanceapi"
	"strconv"
	"time"
)

// 现货订单API接口
func (b BinanceTradeEngine) apiSpotOrderCreate(req *OrderParam) *mybinanceapi.SpotOrderPostApi {
	api := binance.NewSpotRestClient(b.apiKey, b.secretKey).NewSpotOrderPost().
		Symbol(req.Symbol).
		Type(b.bnConverter.ToBNOrderType(req.OrderType)).
		Side(b.bnConverter.ToBNOrderSide(req.OrderSide)).
		Price(req.Price).
		Quantity(req.Quantity)
	if req.ClientOrderId != "" {
		api = api.NewClientOrderId(req.ClientOrderId)
	}
	if req.TimeInForce != "" {
		api = api.TimeInForce(b.bnConverter.ToBNTimeInForce(req.TimeInForce))
	}
	return api
}
func (b BinanceTradeEngine) apiSpotOrderAmend(req *OrderParam) *mybinanceapi.SpotOrderCancelReplaceApi {
	api := binance.NewSpotRestClient(b.apiKey, b.secretKey).NewSpotOrderCancelReplace().
		Symbol(req.Symbol).CancelReplaceMode("STOP_ON_FAILURE").
		Side(b.bnConverter.ToBNOrderSide(req.OrderSide)).
		Type(b.bnConverter.ToBNOrderType(req.OrderType)).
		Price(req.Price).
		Quantity(req.Quantity)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api = api.CancelOrderId(orderId)
	}
	if req.ClientOrderId != "" {
		api = api.CancelOrigClientOrderId(req.ClientOrderId)
	}
	if req.NewClientOrderId != "" {
		api = api.NewClientOrderId(req.NewClientOrderId)
	}
	if req.TimeInForce != "" {
		api = api.TimeInForce(b.bnConverter.ToBNTimeInForce(req.TimeInForce))
	}
	return api
}
func (b BinanceTradeEngine) apiSpotOrderCancel(req *OrderParam) *mybinanceapi.SpotOrderDeleteApi {
	api := binance.NewSpotRestClient(b.apiKey, b.secretKey).NewSpotOrderDelete().
		Symbol(req.Symbol)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api = api.OrderId(orderId)
	} else {
		api = api.OrigClientOrderId(req.ClientOrderId)
	}

	return api
}

// U本位合约订单API接口
func (b BinanceTradeEngine) apiFutureOrderCreate(req *OrderParam) *mybinanceapi.FutureOrderPostApi {
	api := binance.NewFutureRestClient(b.apiKey, b.secretKey).NewFutureOrderPost().
		Symbol(req.Symbol).
		Type(b.bnConverter.ToBNOrderType(req.OrderType)).
		Side(b.bnConverter.ToBNOrderSide(req.OrderSide)).
		PositionSide(b.bnConverter.ToBNPositionSide(req.PositionSide)).
		Price(req.Price).
		Quantity(req.Quantity)
	if req.ClientOrderId != "" {
		api = api.NewClientOrderId(req.ClientOrderId)
	}
	if req.TimeInForce != "" {
		api = api.TimeInForce(b.bnConverter.ToBNTimeInForce(req.TimeInForce))
	}
	return api
}
func (b BinanceTradeEngine) apiFutureOrderAmend(req *OrderParam) *mybinanceapi.FutureOrderPutApi {
	api := binance.NewFutureRestClient(b.apiKey, b.secretKey).NewFutureOrderPut().
		Symbol(req.Symbol).
		Side(b.bnConverter.ToBNOrderSide(req.OrderSide)).
		Price(req.Price).
		Quantity(req.Quantity)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api = api.OrderId(orderId)
	}
	if req.ClientOrderId != "" {
		api = api.OrigClientOrderId(req.ClientOrderId)
	}
	return api
}
func (b BinanceTradeEngine) apiFutureOrderCancel(req *OrderParam) *mybinanceapi.FutureOrderDeleteApi {
	api := binance.NewFutureRestClient(b.apiKey, b.secretKey).NewFutureOrderDelete().
		Symbol(req.Symbol)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api = api.OrderId(orderId)
	} else {
		api = api.OrigClientOrderId(req.ClientOrderId)
	}
	return api
}

func (b BinanceTradeEngine) apiFutureBatchOrderCreate(reqs []*OrderParam) *mybinanceapi.FutureBatchOrdersPostApi {
	client := binance.NewFutureRestClient(b.apiKey, b.secretKey)
	api := client.NewFutureBatchOrdersPost()
	for _, req := range reqs {
		thisApi := client.NewFutureOrderPost().Symbol(req.Symbol).
			Side(b.bnConverter.ToBNOrderSide(req.OrderSide)).
			Type(b.bnConverter.ToBNOrderType(req.OrderType)).
			PositionSide(b.bnConverter.ToBNPositionSide(req.PositionSide)).
			Price(req.Price).
			Quantity(req.Quantity)
		if req.ClientOrderId != "" {
			thisApi = thisApi.NewClientOrderId(req.ClientOrderId)
		}
		if req.TimeInForce != "" {
			thisApi = thisApi.TimeInForce(b.bnConverter.ToBNTimeInForce(req.TimeInForce))
		}
		api = api.AddOrders(thisApi)
	}
	return api
}
func (b BinanceTradeEngine) apiFutureBatchOrderAmend(reqs []*OrderParam) *mybinanceapi.FutureBatchOrdersPutApi {
	client := binance.NewFutureRestClient(b.apiKey, b.secretKey)
	api := client.NewFutureBatchOrdersPut()
	for _, req := range reqs {
		thisApi := client.NewFutureOrderPut().Symbol(req.Symbol).
			Side(b.bnConverter.ToBNOrderSide(req.OrderSide)).
			Price(req.Price).
			Quantity(req.Quantity)
		if req.OrderId != "" {
			orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
			thisApi = thisApi.OrderId(orderId)
		}
		if req.ClientOrderId != "" {
			thisApi = thisApi.OrigClientOrderId(req.ClientOrderId)
		}
		api = api.AddOrders(thisApi)
	}
	return api
}
func (b BinanceTradeEngine) apiFutureBatchOrderCancel(reqs []*OrderParam) (*mybinanceapi.FutureBatchOrdersDeleteApi, error) {
	client := binance.NewFutureRestClient(b.apiKey, b.secretKey)
	orderIds := []int64{}
	clientOrderIds := []string{}
	for _, req := range reqs {
		if req.OrderId != "" {
			orderId, err := strconv.ParseInt(req.OrderId, 10, 64)
			if err != nil {
				return nil, ErrorInvalid("order id")
			}
			orderIds = append(orderIds, orderId)
		} else if req.ClientOrderId != "" {
			clientOrderIds = append(clientOrderIds, req.ClientOrderId)
		} else {
			return nil, ErrorInvalid("order id or client order id is required")
		}
	}
	api := client.NewFutureBatchOrdersDelete().
		Symbol(reqs[0].Symbol)
	if len(orderIds) > 0 {
		api = api.OrderIdList(orderIds)
	} else if len(clientOrderIds) > 0 {
		api = api.OrigClientOrderIdList(clientOrderIds)
	} else {
		return nil, ErrorInvalid("order id or client order id is required")
	}
	return api, nil
}

// 币本位合约订单API接口
func (b BinanceTradeEngine) apiSwapOrderCreate(req *OrderParam) *mybinanceapi.SwapOrderPostApi {
	api := binance.NewSwapRestClient(b.apiKey, b.secretKey).NewSwapOrderPost().
		Symbol(req.Symbol).
		Type(b.bnConverter.ToBNOrderType(req.OrderType)).
		Side(b.bnConverter.ToBNOrderSide(req.OrderSide)).
		PositionSide(b.bnConverter.ToBNPositionSide(req.PositionSide)).
		Price(req.Price).
		Quantity(req.Quantity)
	if req.ClientOrderId != "" {
		api = api.NewClientOrderId(req.ClientOrderId)
	}
	if req.TimeInForce != "" {
		api = api.TimeInForce(b.bnConverter.ToBNTimeInForce(req.TimeInForce))
	}
	return api
}
func (b BinanceTradeEngine) apiSwapOrderAmend(req *OrderParam) *mybinanceapi.SwapOrderPutApi {
	api := binance.NewSwapRestClient(b.apiKey, b.secretKey).NewSwapOrderPut().
		Symbol(req.Symbol).
		Side(b.bnConverter.ToBNOrderSide(req.OrderSide)).
		Price(req.Price).
		Quantity(req.Quantity)
	if req.OrderId != "" {
		api = api.OrderId(req.OrderId)
	}
	if req.ClientOrderId != "" {
		api = api.OrigClientOrderId(req.ClientOrderId)
	}
	return api
}
func (b BinanceTradeEngine) apiSwapOrderCancel(req *OrderParam) *mybinanceapi.SwapOrderDeleteApi {
	api := binance.NewSwapRestClient(b.apiKey, b.secretKey).NewSwapOrderDelete().
		Symbol(req.Symbol)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api = api.OrderId(orderId)
	} else {
		api = api.OrigClientOrderId(req.ClientOrderId)
	}
	return api
}

func (b BinanceTradeEngine) apiSwapBatchOrderCreate(reqs []*OrderParam) *mybinanceapi.SwapBatchOrdersPostApi {
	client := binance.NewSwapRestClient(b.apiKey, b.secretKey)
	api := client.NewSwapBatchOrdersPost()
	for _, req := range reqs {
		thisApi := client.NewSwapOrderPost().Symbol(req.Symbol).
			Side(b.bnConverter.ToBNOrderSide(req.OrderSide)).
			Type(b.bnConverter.ToBNOrderType(req.OrderType)).
			PositionSide(b.bnConverter.ToBNPositionSide(req.PositionSide)).
			Price(req.Price).
			Quantity(req.Quantity)
		if req.ClientOrderId != "" {
			thisApi = thisApi.NewClientOrderId(req.ClientOrderId)
		}
		if req.TimeInForce != "" {
			thisApi = thisApi.TimeInForce(b.bnConverter.ToBNTimeInForce(req.TimeInForce))
		}
		api = api.AddOrders(thisApi)
	}
	return api
}
func (b BinanceTradeEngine) apiSwapBatchOrderAmend(reqs []*OrderParam) *mybinanceapi.SwapBatchOrdersPutApi {
	client := binance.NewSwapRestClient(b.apiKey, b.secretKey)
	api := client.NewSwapBatchOrdersPut()
	for _, req := range reqs {
		thisApi := client.NewSwapOrderPut().Symbol(req.Symbol).
			Side(b.bnConverter.ToBNOrderSide(req.OrderSide)).
			Price(req.Price).
			Quantity(req.Quantity)
		if req.OrderId != "" {
			thisApi = thisApi.OrderId(req.OrderId)
		}
		if req.ClientOrderId != "" {
			thisApi = thisApi.OrigClientOrderId(req.ClientOrderId)
		}
		api = api.AddOrders(thisApi)
	}
	return api
}
func (b BinanceTradeEngine) apiSwapBatchOrderCancel(reqs []*OrderParam) (*mybinanceapi.SwapBatchOrdersDeleteApi, error) {
	client := binance.NewSwapRestClient(b.apiKey, b.secretKey)
	orderIds := []int64{}
	clientOrderIds := []string{}
	for _, req := range reqs {
		if req.OrderId != "" {
			orderId, err := strconv.ParseInt(req.OrderId, 10, 64)
			if err != nil {
				return nil, ErrorInvalid("order id")
			}
			orderIds = append(orderIds, orderId)
		} else if req.ClientOrderId != "" {
			clientOrderIds = append(clientOrderIds, req.ClientOrderId)
		} else {
			return nil, ErrorInvalid("order id or client order id is required")
		}
	}
	api := client.NewSwapBatchOrdersDelete().
		Symbol(reqs[0].Symbol)
	if len(orderIds) > 0 {
		api = api.OrderIdList(orderIds)
	} else if len(clientOrderIds) > 0 {
		api = api.OrigClientOrderIdList(clientOrderIds)
	} else {
		return nil, ErrorInvalid("order id or client order id is required")
	}
	return api, nil
}

// 现货订单处理
func (b BinanceTradeEngine) handleOrderFromSpotOrderCreate(req *OrderParam, res *mybinanceapi.SpotOrderPostRes) *Order {
	order := &Order{
		Exchange:      BINANCE_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.ClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CummulativeQuoteQty,
		Status:        b.bnConverter.FromBNOrderStatus(res.Status),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		CreateTime:    res.WorkingTime,
		UpdateTime:    res.WorkingTime,
	}
	return order
}
func (b BinanceTradeEngine) handleOrderFromSpotOrderAmend(req *OrderParam, res *mybinanceapi.SpotOrderCancelReplaceRes) *Order {
	order := &Order{
		Exchange:      BINANCE_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		OrderId:       strconv.FormatInt(res.NewOrderResponse.OrderId, 10),
		ClientOrderId: res.NewOrderResponse.ClientOrderId,
		Price:         res.NewOrderResponse.Price,
		Quantity:      res.NewOrderResponse.OrigQty,
		ExecutedQty:   res.NewOrderResponse.ExecutedQty,
		CumQuoteQty:   res.NewOrderResponse.CummulativeQuoteQty,
		Status:        b.bnConverter.FromBNOrderStatus(res.NewOrderResponse.Status),
		Type:          b.bnConverter.FromBNOrderType(res.NewOrderResponse.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.NewOrderResponse.Side),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.NewOrderResponse.TimeInForce),
		CreateTime:    res.NewOrderResponse.WorkingTime,
		UpdateTime:    res.NewOrderResponse.WorkingTime,
	}
	return order
}
func (b BinanceTradeEngine) handleOrderFromSpotOrderCancel(req *OrderParam, res *mybinanceapi.SpotOrderDeleteRes) *Order {
	nowTimestamp := time.Now().UnixMilli()
	order := &Order{
		Exchange:      BINANCE_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.OrigClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CummulativeQuoteQty,
		Status:        b.bnConverter.FromBNOrderStatus(res.Status),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		CreateTime:    res.TransactTime,
		UpdateTime:    nowTimestamp,
	}
	return order
}

func (b BinanceTradeEngine) handleOrderFromSpotBatchErr(req *OrderParam, err error) *Order {
	return &Order{
		Exchange:      BINANCE_NAME.String(),
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

// U合约订单处理
func (b BinanceTradeEngine) handleOrderFromFutureOrderCreate(req *OrderParam, res *mybinanceapi.FutureOrderPostRes) *Order {
	nowTimestamp := time.Now().UnixMilli()
	order := &Order{
		Exchange:      BINANCE_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.ClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CumQuote,
		Status:        b.bnConverter.FromBNOrderStatus(res.Status),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		PositionSide:  b.bnConverter.FromBNPositionSide(res.PositionSide),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		ReduceOnly:    res.ReduceOnly,
		CreateTime:    nowTimestamp,
		UpdateTime:    nowTimestamp,
	}
	return order
}
func (b BinanceTradeEngine) handleOrderFromFutureOrderAmend(req *OrderParam, res *mybinanceapi.FutureOrderPutRes) *Order {
	nowTimestamp := time.Now().UnixMilli()
	order := &Order{
		Exchange:      BINANCE_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.ClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CumQuote,
		Status:        b.bnConverter.FromBNOrderStatus(res.Status),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		PositionSide:  b.bnConverter.FromBNPositionSide(res.PositionSide),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		ReduceOnly:    res.ReduceOnly,
		CreateTime:    nowTimestamp,
		UpdateTime:    nowTimestamp,
	}
	return order
}
func (b BinanceTradeEngine) handleOrderFromFutureOrderCancel(req *OrderParam, res *mybinanceapi.FutureOrderDeleteRes) *Order {
	order := &Order{
		Exchange:      BINANCE_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.ClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CumQuote,
		Status:        b.bnConverter.FromBNOrderStatus(res.Status),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		PositionSide:  b.bnConverter.FromBNPositionSide(res.PositionSide),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		ReduceOnly:    res.ReduceOnly,
		CreateTime:    res.UpdateTime,
		UpdateTime:    res.UpdateTime,
	}
	return order
}

func (b BinanceTradeEngine) handleOrdersFromFutureBatchOrderCreate(reqs []*OrderParam, res *mybinanceapi.FutureBatchOrdersPostRes) []*Order {
	var orders []*Order
	nowTimestamp := time.Now().UnixMilli()
	for _, order := range *res {
		orders = append(orders, &Order{
			Exchange:      BINANCE_NAME.String(),
			AccountType:   reqs[0].AccountType,
			Symbol:        order.Symbol,
			OrderId:       strconv.FormatInt(order.OrderId, 10),
			ClientOrderId: order.ClientOrderId,
			Price:         order.Price,
			Quantity:      order.OrigQty,
			ExecutedQty:   order.ExecutedQty,
			CumQuoteQty:   order.CumQuote,
			Status:        b.bnConverter.FromBNOrderStatus(order.Status),
			Type:          b.bnConverter.FromBNOrderType(order.Type),
			Side:          b.bnConverter.FromBNOrderSide(order.Side),
			PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    nowTimestamp,
			UpdateTime:    nowTimestamp,
			ErrorCode:     strconv.Itoa(order.Code),
			ErrorMsg:      order.Msg,
		})
	}
	return orders
}
func (b BinanceTradeEngine) handleOrdersFromFutureBatchOrderAmend(reqs []*OrderParam, res *mybinanceapi.FutureBatchOrdersPutRes) []*Order {
	var orders []*Order
	nowTimestamp := time.Now().UnixMilli()
	for _, order := range *res {
		orders = append(orders, &Order{
			Exchange:      BINANCE_NAME.String(),
			AccountType:   reqs[0].AccountType,
			Symbol:        order.Symbol,
			OrderId:       strconv.FormatInt(order.OrderId, 10),
			ClientOrderId: order.ClientOrderId,
			Price:         order.Price,
			Quantity:      order.OrigQty,
			ExecutedQty:   order.ExecutedQty,
			CumQuoteQty:   order.CumQuote,
			Status:        b.bnConverter.FromBNOrderStatus(order.Status),
			Type:          b.bnConverter.FromBNOrderType(order.Type),
			Side:          b.bnConverter.FromBNOrderSide(order.Side),
			PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    nowTimestamp,
			UpdateTime:    nowTimestamp,
			ErrorCode:     strconv.Itoa(order.Code),
			ErrorMsg:      order.Msg,
		})
	}
	return orders
}
func (b BinanceTradeEngine) handleOrdersFromFutureBatchOrderCancel(reqs []*OrderParam, res *mybinanceapi.FutureBatchOrdersDeleteRes) []*Order {
	var orders []*Order
	nowTimestamp := time.Now().UnixMilli()
	for _, order := range *res {
		orders = append(orders, &Order{
			Exchange:      BINANCE_NAME.String(),
			AccountType:   reqs[0].AccountType,
			Symbol:        order.Symbol,
			OrderId:       strconv.FormatInt(order.OrderId, 10),
			ClientOrderId: order.ClientOrderId,
			Price:         order.Price,
			Quantity:      order.OrigQty,
			ExecutedQty:   order.ExecutedQty,
			CumQuoteQty:   order.CumQuote,
			Status:        b.bnConverter.FromBNOrderStatus(order.Status),
			Type:          b.bnConverter.FromBNOrderType(order.Type),
			Side:          b.bnConverter.FromBNOrderSide(order.Side),
			PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    nowTimestamp,
			UpdateTime:    nowTimestamp,
			ErrorCode:     strconv.Itoa(order.Code),
			ErrorMsg:      order.Msg,
		})
	}
	return orders
}

// 币本位合约订单处理
func (b BinanceTradeEngine) handleOrderFromSwapOrderCreate(req *OrderParam, res *mybinanceapi.SwapOrderPostRes) *Order {
	nowTimestamp := time.Now().UnixMilli()
	order := Order{
		Exchange:      BINANCE_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.ClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CumQuote,
		Status:        b.bnConverter.FromBNOrderStatus(res.Status),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		PositionSide:  b.bnConverter.FromBNPositionSide(res.PositionSide),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		ReduceOnly:    res.ReduceOnly,
		CreateTime:    nowTimestamp,
		UpdateTime:    nowTimestamp,
	}
	return &order
}
func (b BinanceTradeEngine) handleOrderFromSwapOrderAmend(req *OrderParam, res *mybinanceapi.SwapOrderPutRes) *Order {
	nowTimestamp := time.Now().UnixMilli()
	order := &Order{
		Exchange:      BINANCE_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.ClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CumQuote,
		Status:        b.bnConverter.FromBNOrderStatus(res.Status),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		PositionSide:  b.bnConverter.FromBNPositionSide(res.PositionSide),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		ReduceOnly:    res.ReduceOnly,
		CreateTime:    nowTimestamp,
		UpdateTime:    nowTimestamp,
	}
	return order
}
func (b BinanceTradeEngine) handleOrderFromSwapOrderCancel(req *OrderParam, res *mybinanceapi.SwapOrderDeleteRes) *Order {
	order := &Order{
		Exchange:      BINANCE_NAME.String(),
		AccountType:   req.AccountType,
		Symbol:        req.Symbol,
		OrderId:       strconv.FormatInt(res.OrderId, 10),
		ClientOrderId: res.ClientOrderId,
		Price:         res.Price,
		Quantity:      res.OrigQty,
		ExecutedQty:   res.ExecutedQty,
		CumQuoteQty:   res.CumQuote,
		Status:        b.bnConverter.FromBNOrderStatus(res.Status),
		Type:          b.bnConverter.FromBNOrderType(res.Type),
		Side:          b.bnConverter.FromBNOrderSide(res.Side),
		PositionSide:  b.bnConverter.FromBNPositionSide(res.PositionSide),
		TimeInForce:   b.bnConverter.FromBNTimeInForce(res.TimeInForce),
		ReduceOnly:    res.ReduceOnly,
		CreateTime:    res.UpdateTime,
		UpdateTime:    res.UpdateTime,
	}
	return order
}

func (b BinanceTradeEngine) handleOrdersFromSwapBatchOrderCreate(reqs []*OrderParam, res *mybinanceapi.SwapBatchOrdersPostRes) []*Order {
	var orders []*Order
	nowTimestamp := time.Now().UnixMilli()
	for _, order := range *res {
		orders = append(orders, &Order{
			Exchange:      BINANCE_NAME.String(),
			AccountType:   reqs[0].AccountType,
			Symbol:        order.Symbol,
			OrderId:       strconv.FormatInt(order.OrderId, 10),
			ClientOrderId: order.ClientOrderId,
			Price:         order.Price,
			Quantity:      order.OrigQty,
			ExecutedQty:   order.ExecutedQty,
			CumQuoteQty:   order.CumQuote,
			Status:        b.bnConverter.FromBNOrderStatus(order.Status),
			Type:          b.bnConverter.FromBNOrderType(order.Type),
			Side:          b.bnConverter.FromBNOrderSide(order.Side),
			PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    nowTimestamp,
			UpdateTime:    nowTimestamp,
			ErrorCode:     strconv.Itoa(order.Code),
			ErrorMsg:      order.Msg,
		})
	}
	return orders
}
func (b BinanceTradeEngine) handleOrdersFromSwapBatchOrderAmend(reqs []*OrderParam, res *mybinanceapi.SwapBatchOrdersPutRes) []*Order {
	var orders []*Order
	nowTimestamp := time.Now().UnixMilli()
	for _, order := range *res {
		orders = append(orders, &Order{
			Exchange:      BINANCE_NAME.String(),
			AccountType:   reqs[0].AccountType,
			Symbol:        order.Symbol,
			OrderId:       strconv.FormatInt(order.OrderId, 10),
			ClientOrderId: order.ClientOrderId,
			Price:         order.Price,
			Quantity:      order.OrigQty,
			ExecutedQty:   order.ExecutedQty,
			CumQuoteQty:   order.CumQuote,
			Status:        b.bnConverter.FromBNOrderStatus(order.Status),
			Type:          b.bnConverter.FromBNOrderType(order.Type),
			Side:          b.bnConverter.FromBNOrderSide(order.Side),
			PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    nowTimestamp,
			UpdateTime:    nowTimestamp,
			ErrorCode:     strconv.Itoa(order.Code),
			ErrorMsg:      order.Msg,
		})
	}
	return orders
}
func (b BinanceTradeEngine) handleOrdersFromSwapBatchOrderCancel(reqs []*OrderParam, res *mybinanceapi.SwapBatchOrdersDeleteRes) []*Order {
	var orders []*Order
	nowTimestamp := time.Now().UnixMilli()
	for _, order := range *res {
		orders = append(orders, &Order{
			Exchange:      BINANCE_NAME.String(),
			AccountType:   reqs[0].AccountType,
			Symbol:        order.Symbol,
			OrderId:       strconv.FormatInt(order.OrderId, 10),
			ClientOrderId: order.ClientOrderId,
			Price:         order.Price,
			Quantity:      order.OrigQty,
			ExecutedQty:   order.ExecutedQty,
			CumQuoteQty:   order.CumQuote,
			Status:        b.bnConverter.FromBNOrderStatus(order.Status),
			Type:          b.bnConverter.FromBNOrderType(order.Type),
			Side:          b.bnConverter.FromBNOrderSide(order.Side),
			PositionSide:  b.bnConverter.FromBNPositionSide(order.PositionSide),
			TimeInForce:   b.bnConverter.FromBNTimeInForce(order.TimeInForce),
			ReduceOnly:    order.ReduceOnly,
			CreateTime:    nowTimestamp,
			UpdateTime:    nowTimestamp,
			ErrorCode:     strconv.Itoa(order.Code),
			ErrorMsg:      order.Msg,
		})
	}
	return orders
}
