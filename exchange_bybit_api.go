package mytrade

import (
	"github.com/Hongssd/mybybitapi"
	"github.com/shopspring/decimal"
)

// 查询订单接口获取
func (b *BybitTradeEngine) apiQueryOpenOrders(req *QueryOrderParam, pageCursor string) *mybybitapi.OrderRealtimeAPI {
	client := mybybitapi.NewRestClient(b.apiKey, b.secretKey).PrivateRestClient()
	api := client.NewOrderRealtime().Category(req.AccountType).Limit(50)

	if req.Symbol != "" {
		api.Symbol(req.Symbol)
	} else {
		if req.AccountType == BYBIT_AC_LINEAR.String() {
			//合约必传交易对或币种名
			if req.SettleCoin != "" {
				api.SettleCoin(req.SettleCoin)
			}
			if req.BaseCoin != "" {
				api.BaseCoin(req.BaseCoin)
			}
		}
	}

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
func (b *BybitTradeEngine) apiQueryOrders(req *QueryOrderParam) *mybybitapi.OrderHistoryAPI {
	client := mybybitapi.NewRestClient(b.apiKey, b.secretKey).PrivateRestClient()
	api := client.NewOrderHistory().Category(req.AccountType)

	if req.Symbol != "" {
		api.Symbol(req.Symbol)
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

	if !req.TriggerPrice.IsZero() && req.TriggerType != ORDER_TRIGGER_TYPE_UNKNOWN {
		api.TriggerPrice(req.TriggerPrice.String())
		api.TriggerDirection(b.bybitConverter.ToBYBITTriggerCondition(req.TriggerType, req.OrderSide))
	}

	if req.AccountType == BYBIT_AC_SPOT.String() {
		api.MarketUnit("baseCoin")
	}
	if req.IsMargin {
		api.IsLeverage(1)
	}

	if req.PositionSide != "" {
		api.PositionIdx(b.bybitConverter.ToBYBITPositionSide(req.PositionSide))
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
