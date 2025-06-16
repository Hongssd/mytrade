package mytrade

import (
	"strconv"

	"github.com/Hongssd/mybinanceapi"
)

// 现货订单API接口
func (b *BinanceTradeEngine) apiSpotOpenOrders(req *QueryOrderParam) *mybinanceapi.SpotOpenOrdersApi {
	api := binance.NewSpotRestClient(b.apiKey, b.secretKey).NewOpenOrders()
	if req.Symbol != "" {
		api.Symbol(req.Symbol)
	}
	return api
}
func (b *BinanceTradeEngine) apiSpotOrderQuery(req *QueryOrderParam) *mybinanceapi.SpotOrderGetApi {
	api := binance.NewSpotRestClient(b.apiKey, b.secretKey).NewSpotOrderGet().Symbol(req.Symbol)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api = api.OrderId(orderId)
	} else {
		if req.ClientOrderId != "" {
			api = api.OrigClientOrderId(req.ClientOrderId)
		}
	}
	return api
}
func (b *BinanceTradeEngine) apiSpotOrdersQuery(req *QueryOrderParam) *mybinanceapi.SpotAllOrdersApi {
	api := binance.NewSpotRestClient(b.apiKey, b.secretKey).NewAllOrders().Symbol(req.Symbol)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api.OrderId(orderId)
	}
	if req.StartTime != 0 {
		api = api.StartTime(req.StartTime)
	}
	if req.EndTime != 0 {
		api = api.EndTime(req.EndTime)
	}
	if req.Limit != 0 {
		api = api.Limit(req.Limit)
	}
	return api
}
func (b *BinanceTradeEngine) apiSpotTradeQuery(req *QueryTradeParam) *mybinanceapi.SpotMyTradesApi {
	api := binance.NewSpotRestClient(b.apiKey, b.secretKey).NewSpotMyTrades().
		Symbol(req.Symbol)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api = api.OrderId(orderId)
	}
	if req.StartTime != 0 {
		api = api.StartTime(req.StartTime)
	}
	if req.EndTime != 0 {
		api = api.EndTime(req.EndTime)
	}
	if req.Limit != 0 {
		api = api.Limit(req.Limit)
	}
	return api
}

func (b *BinanceTradeEngine) apiSpotOrderCreate(req *OrderParam) *mybinanceapi.SpotOrderPostApi {
	api := binance.NewSpotRestClient(b.apiKey, b.secretKey).NewSpotOrderPost().
		Symbol(req.Symbol).
		Side(b.bnConverter.ToBNOrderSide(req.OrderSide)).
		Quantity(req.Quantity).NewOrderRespType("FULL")

	api.Type(b.bnConverter.ToTriggerBnOrderType(BinanceAccountType(req.AccountType), req.OrderType, req.TriggerType))

	if !req.TriggerPrice.IsZero() {
		api.StopPrice(req.TriggerPrice)
	}

	log.Info(req)

	if !req.Price.IsZero() {
		api = api.Price(req.Price)
	}
	if req.ClientOrderId != "" {
		api = api.NewClientOrderId(req.ClientOrderId)
	}
	if req.TimeInForce != "" {
		if req.TimeInForce == TIME_IN_FORCE_POST_ONLY {
			// 现货POSTONLY下单 不传timeInforce并且将订单类型为LIMIT_MAKER
			// api.TimeInForce(b.bnConverter.ToBNTimeInForce(TIME_IN_FORCE_GTC))
			api.Type(BN_ORDER_TYPE_LIMIT_MAKER)
		} else {
			api = api.TimeInForce(b.bnConverter.ToBNTimeInForce(req.TimeInForce))
		}
	}
	return api
}
func (b *BinanceTradeEngine) apiSpotOrderAmend(req *OrderParam) *mybinanceapi.SpotOrderCancelReplaceApi {
	api := binance.NewSpotRestClient(b.apiKey, b.secretKey).NewSpotOrderCancelReplace().
		Symbol(req.Symbol).CancelReplaceMode("STOP_ON_FAILURE").
		Side(b.bnConverter.ToBNOrderSide(req.OrderSide)).
		Quantity(req.Quantity)

	api.Type(b.bnConverter.ToTriggerBnOrderType(BinanceAccountType(req.AccountType), req.OrderType, req.TriggerType))
	if !req.TriggerPrice.IsZero() {
		api.StopPrice(req.TriggerPrice)
	}
	if !req.Price.IsZero() {
		api = api.Price(req.Price)
	}
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
func (b *BinanceTradeEngine) apiSpotOrderCancel(req *OrderParam) *mybinanceapi.SpotOrderDeleteApi {
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

// 现货杠杆订单API接口
func (b *BinanceTradeEngine) apiSpotMarginOpenOrders(req *QueryOrderParam) *mybinanceapi.SpotMarginOpenOrdersApi {
	api := binance.NewSpotRestClient(b.apiKey, b.secretKey).NewSpotMarginOpenOrders()
	if req.Symbol != "" {
		api.Symbol(req.Symbol)
	}
	// 如果isIsolated=“TRUE”，symbol为必填
	if req.IsIsolated {
		api.IsIsolated("TRUE")
	}
	return api
}
func (b *BinanceTradeEngine) apiSpotMarginOrderQuery(req *QueryOrderParam) *mybinanceapi.SpotMarginOrderGetApi {
	api := binance.NewSpotRestClient(b.apiKey, b.secretKey).NewSpotMarginOrderGet().Symbol(req.Symbol)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api = api.OrderId(orderId)
	} else {
		api = api.OrigClientOrderId(req.ClientOrderId)
	}
	if req.IsIsolated {
		api.IsIsolated("TRUE")
	}
	return api
}
func (b *BinanceTradeEngine) apiSpotMarginOrdersQuery(req *QueryOrderParam) *mybinanceapi.SpotMarginAllOrdersApi {
	api := binance.NewSpotRestClient(b.apiKey, b.secretKey).NewSpotMarginAllOrders().Symbol(req.Symbol)
	if req.IsIsolated {
		api = api.IsIsolated("TRUE")
	}
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api.OrderId(orderId)
	}
	if req.StartTime != 0 {
		api = api.StartTime(req.StartTime)
	}
	if req.EndTime != 0 {
		api = api.EndTime(req.EndTime)
	}
	if req.Limit != 0 {
		api = api.Limit(int64(req.Limit))
	}
	return api
}

func (b *BinanceTradeEngine) apiSpotMarginOrderCreate(req *OrderParam) *mybinanceapi.SpotMarginOrderPostApi {
	api := binance.NewSpotRestClient(b.apiKey, b.secretKey).NewSpotMarginOrderPost().
		Symbol(req.Symbol).
		Side(b.bnConverter.ToBNOrderSide(req.OrderSide)).
		Quantity(req.Quantity).NewOrderRespType("FULL")

	api.Type(b.bnConverter.ToTriggerBnOrderType(BinanceAccountType(req.AccountType), req.OrderType, req.TriggerType))

	// 默认为全仓模式 FALSE
	if req.IsIsolated {
		api.IsIsolated("TRUE")
	}

	if !req.TriggerPrice.IsZero() {
		api.StopPrice(req.TriggerPrice)
	}

	log.Info(req)

	if !req.Price.IsZero() {
		api = api.Price(req.Price)
	}
	if req.ClientOrderId != "" {
		api = api.NewClientOrderId(req.ClientOrderId)
	}
	if req.TimeInForce != "" {
		if req.TimeInForce == TIME_IN_FORCE_POST_ONLY {
			// 现货杠杆POSTONLY下单 不传timeInforce并且将订单类型为LIMIT_MAKER
			// api.TimeInForce(b.bnConverter.ToBNTimeInForce(TIME_IN_FORCE_GTC))
			api.Type(BN_ORDER_TYPE_LIMIT_MAKER)
		} else {
			api = api.TimeInForce(b.bnConverter.ToBNTimeInForce(req.TimeInForce))
		}
	}

	return api
}
func (b *BinanceTradeEngine) apiSpotMarginOrderCancel(req *OrderParam) *mybinanceapi.SpotMarginOrderDeleteApi {
	api := binance.NewSpotRestClient(b.apiKey, b.secretKey).NewSpotMarginOrderDelete().
		Symbol(req.Symbol)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api = api.OrderId(orderId)
	} else {
		api = api.OrigClientOrderId(req.ClientOrderId)
	}
	if req.IsIsolated {
		api.IsIsolated("TRUE")
	}
	return api
}

// U本位合约订单API接口
func (b *BinanceTradeEngine) apiFutureOpenOrders(req *QueryOrderParam) *mybinanceapi.FutureOpenOrdersApi {
	api := binance.NewFutureRestClient(b.apiKey, b.secretKey).NewOpenOrders()
	if req.Symbol != "" {
		api.Symbol(req.Symbol)
	}
	return api
}
func (b *BinanceTradeEngine) apiFutureOrderQuery(req *QueryOrderParam) *mybinanceapi.FutureOrderGetApi {
	api := binance.NewFutureRestClient(b.apiKey, b.secretKey).NewFutureOrderGet().Symbol(req.Symbol)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api = api.OrderId(orderId)
	} else {
		api = api.OrigClientOrderId(req.ClientOrderId)
	}
	return api
}
func (b *BinanceTradeEngine) apiFutureOrdersQuery(req *QueryOrderParam) *mybinanceapi.FutureAllOrdersApi {
	api := binance.NewFutureRestClient(b.apiKey, b.secretKey).NewAllOrders().Symbol(req.Symbol)

	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api = api.OrderId(orderId)
	}

	if req.StartTime != 0 {
		api = api.StartTime(req.StartTime)
	}
	if req.EndTime != 0 {
		api = api.EndTime(req.EndTime)
	}
	if req.Limit != 0 {
		api = api.Limit(int64(req.Limit))
	}

	return api
}
func (b *BinanceTradeEngine) apiFutureTradeQuery(req *QueryTradeParam) *mybinanceapi.FutureUserTradesApi {
	api := binance.NewFutureRestClient(b.apiKey, b.secretKey).NewFutureUserTrades().
		Symbol(req.Symbol)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api = api.OrderId(orderId)
	}
	if req.StartTime != 0 {
		api = api.StartTime(req.StartTime)
	}
	if req.EndTime != 0 {
		api = api.EndTime(req.EndTime)
	}
	if req.Limit != 0 {
		api = api.Limit(int64(req.Limit))
	}
	return api
}

func (b *BinanceTradeEngine) apiFutureOrderCreate(req *OrderParam) *mybinanceapi.FutureOrderPostApi {
	api := binance.NewFutureRestClient(b.apiKey, b.secretKey).NewFutureOrderPost().
		Symbol(req.Symbol).
		Side(b.bnConverter.ToBNOrderSide(req.OrderSide)).
		PositionSide(b.bnConverter.ToBNPositionSide(req.PositionSide)).
		Quantity(req.Quantity).NewOrderRespType("FULL")

	api.Type(b.bnConverter.ToTriggerBnOrderType(BinanceAccountType(req.AccountType), req.OrderType, req.TriggerType))
	if !req.TriggerPrice.IsZero() {
		api.StopPrice(req.TriggerPrice)
	}
	if !req.Price.IsZero() {
		api = api.Price(req.Price)
	}
	if req.ClientOrderId != "" {
		api = api.NewClientOrderId(req.ClientOrderId)
	}
	if req.TimeInForce != "" {
		api = api.TimeInForce(b.bnConverter.ToBNTimeInForce(req.TimeInForce))
	}
	return api
}
func (b *BinanceTradeEngine) apiFutureOrderAmend(req *OrderParam) *mybinanceapi.FutureOrderPutApi {
	api := binance.NewFutureRestClient(b.apiKey, b.secretKey).NewFutureOrderPut().
		Symbol(req.Symbol).
		Side(b.bnConverter.ToBNOrderSide(req.OrderSide)).
		Quantity(req.Quantity)

	if !req.Price.IsZero() {
		api = api.Price(req.Price)
	}
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api = api.OrderId(orderId)
	}
	if req.ClientOrderId != "" {
		api = api.OrigClientOrderId(req.ClientOrderId)
	}
	return api
}
func (b *BinanceTradeEngine) apiFutureOrderCancel(req *OrderParam) *mybinanceapi.FutureOrderDeleteApi {
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

func (b *BinanceTradeEngine) apiFutureBatchOrderCreate(reqs []*OrderParam) *mybinanceapi.FutureBatchOrdersPostApi {
	client := binance.NewFutureRestClient(b.apiKey, b.secretKey)
	api := client.NewFutureBatchOrdersPost()
	for _, req := range reqs {
		thisApi := client.NewFutureOrderPost().Symbol(req.Symbol).
			Side(b.bnConverter.ToBNOrderSide(req.OrderSide)).
			PositionSide(b.bnConverter.ToBNPositionSide(req.PositionSide)).
			Quantity(req.Quantity).NewOrderRespType("FULL")

		thisApi.Type(b.bnConverter.ToTriggerBnOrderType(BinanceAccountType(req.AccountType), req.OrderType, req.TriggerType))
		if !req.TriggerPrice.IsZero() {
			thisApi.StopPrice(req.TriggerPrice)
		}
		if !req.Price.IsZero() {
			thisApi = thisApi.Price(req.Price)
		}
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
func (b *BinanceTradeEngine) apiFutureBatchOrderAmend(reqs []*OrderParam) *mybinanceapi.FutureBatchOrdersPutApi {
	client := binance.NewFutureRestClient(b.apiKey, b.secretKey)
	api := client.NewFutureBatchOrdersPut()
	for _, req := range reqs {
		thisApi := client.NewFutureOrderPut().Symbol(req.Symbol).
			Side(b.bnConverter.ToBNOrderSide(req.OrderSide)).
			Quantity(req.Quantity)
		if !req.Price.IsZero() {
			thisApi = thisApi.Price(req.Price)
		}
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
func (b *BinanceTradeEngine) apiFutureBatchOrderCancel(reqs []*OrderParam) (*mybinanceapi.FutureBatchOrdersDeleteApi, error) {
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
func (b *BinanceTradeEngine) apiSwapOpenOrders(req *QueryOrderParam) *mybinanceapi.SwapOpenOrdersApi {
	api := binance.NewSwapRestClient(b.apiKey, b.secretKey).NewOpenOrders()
	if req.Symbol != "" {
		api.Symbol(req.Symbol)
	}
	return api
}
func (b *BinanceTradeEngine) apiSwapOrderQuery(req *QueryOrderParam) *mybinanceapi.SwapOrderGetApi {
	api := binance.NewSwapRestClient(b.apiKey, b.secretKey).NewSwapOrderGet().Symbol(req.Symbol)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api = api.OrderId(orderId)
	} else {
		api = api.OrigClientOrderId(req.ClientOrderId)
	}
	return api
}
func (b *BinanceTradeEngine) apiSwapOrdersQuery(req *QueryOrderParam) *mybinanceapi.SwapAllOrdersApi {
	api := binance.NewSwapRestClient(b.apiKey, b.secretKey).NewAllOrders().Symbol(req.Symbol)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api = api.OrderId(orderId)
	}
	if req.StartTime != 0 {
		api = api.StartTime(req.StartTime)
	}
	if req.EndTime != 0 {
		api = api.EndTime(req.EndTime)
	}
	if req.Limit != 0 {
		api = api.Limit(int64(req.Limit))
	}
	return api
}
func (b *BinanceTradeEngine) apiSwapTradeQuery(req *QueryTradeParam) *mybinanceapi.SwapUserTradesApi {
	api := binance.NewSwapRestClient(b.apiKey, b.secretKey).NewSwapUserTrades().
		Symbol(req.Symbol)
	if req.OrderId != "" {
		orderId, _ := strconv.ParseInt(req.OrderId, 10, 64)
		api = api.OrderId(orderId)
	}
	if req.StartTime != 0 {
		api = api.StartTime(req.StartTime)
	}
	if req.EndTime != 0 {
		api = api.EndTime(req.EndTime)
	}
	if req.Limit != 0 {
		api = api.Limit(int64(req.Limit))
	}
	return api
}

func (b *BinanceTradeEngine) apiSwapOrderCreate(req *OrderParam) *mybinanceapi.SwapOrderPostApi {
	api := binance.NewSwapRestClient(b.apiKey, b.secretKey).NewSwapOrderPost().
		Symbol(req.Symbol).
		Side(b.bnConverter.ToBNOrderSide(req.OrderSide)).
		PositionSide(b.bnConverter.ToBNPositionSide(req.PositionSide)).
		Quantity(req.Quantity).NewOrderRespType("FULL")

	api.Type(b.bnConverter.ToTriggerBnOrderType(BinanceAccountType(req.AccountType), req.OrderType, req.TriggerType))
	if !req.TriggerPrice.IsZero() {
		api.StopPrice(req.TriggerPrice)
	}
	if !req.Price.IsZero() {
		api = api.Price(req.Price)
	}
	if req.ClientOrderId != "" {
		api = api.NewClientOrderId(req.ClientOrderId)
	}
	if req.TimeInForce != "" {
		api = api.TimeInForce(b.bnConverter.ToBNTimeInForce(req.TimeInForce))
	}
	return api
}
func (b *BinanceTradeEngine) apiSwapOrderAmend(req *OrderParam) *mybinanceapi.SwapOrderPutApi {
	api := binance.NewSwapRestClient(b.apiKey, b.secretKey).NewSwapOrderPut().
		Symbol(req.Symbol).
		Side(b.bnConverter.ToBNOrderSide(req.OrderSide)).
		Quantity(req.Quantity)
	if !req.Price.IsZero() {
		api = api.Price(req.Price)
	}
	if req.OrderId != "" {
		api = api.OrderId(req.OrderId)
	}
	if req.ClientOrderId != "" {
		api = api.OrigClientOrderId(req.ClientOrderId)
	}
	return api
}
func (b *BinanceTradeEngine) apiSwapOrderCancel(req *OrderParam) *mybinanceapi.SwapOrderDeleteApi {
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

func (b *BinanceTradeEngine) apiSwapBatchOrderCreate(reqs []*OrderParam) *mybinanceapi.SwapBatchOrdersPostApi {
	client := binance.NewSwapRestClient(b.apiKey, b.secretKey)
	api := client.NewSwapBatchOrdersPost()
	for _, req := range reqs {
		thisApi := client.NewSwapOrderPost().Symbol(req.Symbol).
			Side(b.bnConverter.ToBNOrderSide(req.OrderSide)).
			PositionSide(b.bnConverter.ToBNPositionSide(req.PositionSide)).
			Quantity(req.Quantity).NewOrderRespType("FULL")

		thisApi.Type(b.bnConverter.ToTriggerBnOrderType(BinanceAccountType(req.AccountType), req.OrderType, req.TriggerType))
		if !req.TriggerPrice.IsZero() {
			thisApi.StopPrice(req.TriggerPrice)
		}
		if !req.Price.IsZero() {
			thisApi = thisApi.Price(req.Price)
		}
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
func (b *BinanceTradeEngine) apiSwapBatchOrderAmend(reqs []*OrderParam) *mybinanceapi.SwapBatchOrdersPutApi {
	client := binance.NewSwapRestClient(b.apiKey, b.secretKey)
	api := client.NewSwapBatchOrdersPut()
	for _, req := range reqs {
		thisApi := client.NewSwapOrderPut().Symbol(req.Symbol).
			Side(b.bnConverter.ToBNOrderSide(req.OrderSide)).
			Quantity(req.Quantity)
		if !req.Price.IsZero() {
			thisApi = thisApi.Price(req.Price)
		}
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
func (b *BinanceTradeEngine) apiSwapBatchOrderCancel(reqs []*OrderParam) (*mybinanceapi.SwapBatchOrdersDeleteApi, error) {
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

func (b *BinanceTradeEngine) restBatchPreCheck(reqs []*OrderParam) error {
	//检测长度，BINANCE最多批量下5个订单
	if len(reqs) > 5 {
		return ErrorInvalid("binance order param length require less than 5")
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

func (b *BinanceTradeEngine) accountTypePreCheck(accountType string) error {
	switch BinanceAccountType(accountType) {
	case BN_AC_SPOT, BN_AC_FUTURE, BN_AC_SWAP:
		return nil
	default:
		return ErrorInvalid("binance account type invalid")
	}
}
func (b *BinanceTradeEngine) checkWsSpotAccount() error {
	var err error
	if b.wsSpotAccount == nil {
		b.wsSpotAccount, err = binance.NewSpotWsStreamClient().ConvertToAccountWs(b.apiKey, b.secretKey, mybinanceapi.SPOT_WS_TYPE)
		if err != nil {
			return err
		}
		err := b.wsSpotAccount.OpenConn()
		if err != nil {
			return err
		}
	}
	return nil
}
func (b *BinanceTradeEngine) checkWsSpotMarginAccount() error {
	var err error
	if b.wsSpotMarginAccount == nil {
		b.wsSpotMarginAccount, err = binance.NewSpotWsStreamClient().ConvertToAccountWs(b.apiKey, b.secretKey, mybinanceapi.SPOT_MARGIN_WS_TYPE)
		if err != nil {
			return err
		}
		err := b.wsSpotMarginAccount.OpenConn()
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *BinanceTradeEngine) checkWsSpotIsolatedMarginAccount(symbol string) error {
	var err error
	if _, ok := b.wsSpotIsolatedMarginAccount.Load(symbol); !ok {
		var targetWs *mybinanceapi.SpotWsStreamClient
		targetWs, err = binance.NewSpotWsStreamClient().ConvertToAccountWs(b.apiKey, b.secretKey, mybinanceapi.SPOT_ISOLATED_MARGIN_WS_TYPE, symbol)
		if err != nil {
			return err
		}
		err := targetWs.OpenConn()
		if err != nil {
			return err
		}
		b.wsSpotIsolatedMarginAccount.Store(symbol, targetWs)
	}
	return nil
}

func (b *BinanceTradeEngine) checkWsFutureAccount() error {
	var err error
	if b.wsFutureAccount == nil {
		b.wsFutureAccount, err = binance.NewFutureWsStreamClient().ConvertToAccountWs(b.apiKey, b.secretKey)
		if err != nil {
			return err
		}
		err := b.wsFutureAccount.OpenConn()
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *BinanceTradeEngine) checkWsSwapAccount() error {
	var err error
	if b.wsSwapAccount == nil {
		b.wsSwapAccount, err = binance.NewSwapWsStreamClient().ConvertToAccountWs(b.apiKey, b.secretKey)
		if err != nil {
			return err
		}
		err := b.wsSwapAccount.OpenConn()
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *BinanceTradeEngine) checkWsPmContractAccount() error {
	var err error
	if b.wsPMContractAccount == nil {
		b.wsPMContractAccount, err = binance.NewPMContractStreamClient().ConvertToAccountWs(b.apiKey, b.secretKey)
		if err != nil {
			return err
		}
		err := b.wsPMContractAccount.OpenConn()
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *BinanceTradeEngine) checkWsPmMarginAccount() error {
	var err error
	if b.wsPMMarginAccount == nil {
		b.wsPMMarginAccount, err = binance.NewPMMarginStreamClient().ConvertToAccountWs(b.apiKey, b.secretKey)
		if err != nil {
			return err
		}
		err := b.wsPMMarginAccount.OpenConn()
		if err != nil {
			return err
		}
	}
	return nil
}
