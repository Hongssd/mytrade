package mytrade

import (
	"strconv"

	"github.com/Hongssd/myokxapi"
)

// 查询订单接口获取
func (o *OkxTradeEngine) apiQueryOpenOrders(req *QueryOrderParam) *myokxapi.PrivateRestTradeOrdersPendingAPI {
	client := okx.NewRestClient(o.apiKey, o.secretKey, o.passphrase).PrivateRestClient()
	api := client.NewPrivateRestTradeOrdersPending()
	//if req.AccountType != "" {
	//	api.InstType(req.AccountType)
	//}
	if req.Symbol != "" {
		api.InstId(req.Symbol)
	}
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
func (o *OkxTradeEngine) apiQueryOrders(req *QueryOrderParam) *myokxapi.PrivateRestTradeOrderHistoryAPI {
	client := okx.NewRestClient(o.apiKey, o.secretKey, o.passphrase).PrivateRestClient()
	api := client.NewPrivateRestTradeOrderHistory()
	if req.AccountType == OKX_AC_SPOT.String() && req.IsMargin {
		api.InstType(OKX_AC_MARGIN.String())
	} else {
		api.InstType(req.AccountType)
	}

	if req.Symbol != "" {
		api.InstId(req.Symbol)
	}
	if req.OrderId != "" {
		api.Before(req.OrderId)
	}
	if req.StartTime != 0 {
		api.Begin(strconv.FormatInt(req.StartTime, 10))
	}
	if req.EndTime != 0 {
		api.End(strconv.FormatInt(req.EndTime, 10))
	}
	if req.Limit != 0 {
		api.Limit(strconv.Itoa(req.Limit))
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
	tdMode := o.okxConverter.getTdModeFromAccountType(OkxAccountType(req.AccountType), o.okxConverter.ToOKXAccountMode(req.AccountMode), req.IsIsolated, req.IsMargin)

	api := client.NewPrivateRestTradeOrderPost().
		InstId(req.Symbol).TdMode(tdMode).
		Side(o.okxConverter.ToOKXOrderSide(req.OrderSide)).
		OrdType(o.okxConverter.ToOKXOrderType(req.OrderType, req.TimeInForce)).
		Px(req.Price.String()).
		Sz(req.Quantity.String())
	if req.IsMargin && !req.IsIsolated {
		if req.Ccy == "" {
			api.Ccy("USDT")
		} else {
			api.Ccy(req.Ccy)
		}
	}

	if !req.AttachOcoTpTriggerPrice.IsZero() || !req.AttachOcoSlTriggerPrice.IsZero() {
		api.AttachAlgoOrds([]myokxapi.PrivateRestTradeOrderPostReqAttachAlgoOrd{
			{
				TpTriggerPx: GetPointer(req.AttachOcoTpTriggerPrice.String()),
				TpOrdPx:     GetPointer(req.AttachOcoTpOrderPrice.String()),
				SlTriggerPx: GetPointer(req.AttachOcoSlTriggerPrice.String()),
				SlOrdPx:     GetPointer(req.AttachOcoSlOrderPrice.String()),
			},
		})
	}

	if OkxAccountType(req.AccountType) != "SPOT" {
		api.PosSide(o.okxConverter.ToOKXPositionSide(req.PositionSide))
	} else {
		api.TgtCcy("base_ccy")
	}
	if OkxAccountType(req.AccountType) != "SPOT" && !req.IsMargin {
		if req.ReduceOnly {
			api.ReduceOnly(req.ReduceOnly)
		}
	}
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

// 策略委托订单接口获取
func (o *OkxTradeEngine) apiOrderAlgoCreate(req *OrderParam) *myokxapi.PrivateRestTradeOrderAlgoPostAPI {
	client := okx.NewRestClient(o.apiKey, o.secretKey, o.passphrase).PrivateRestClient()

	tdMode := o.okxConverter.getTdModeFromAccountType(OkxAccountType(req.AccountType), o.okxConverter.ToOKXAccountMode(req.AccountMode), req.IsIsolated, req.IsMargin)

	api := client.NewPrivateRestTradeOrderAlgoPost().
		InstId(req.Symbol).
		TdMode(tdMode).
		Side(o.okxConverter.ToOKXOrderSide(req.OrderSide)).
		Sz(req.Quantity.String())

	if req.OrderAlgoType != "" {
		api.OrdType(req.OrderAlgoType.String())
		switch req.OrderAlgoType {
		case OKX_ORDER_ALGO_TYPE_CONDITIONAL:
			if req.TriggerType != ORDER_TRIGGER_TYPE_UNKNOWN && !req.TriggerPrice.IsZero() {
				switch req.TriggerType {
				case ORDER_TRIGGER_TYPE_STOP_LOSS:
					api.ConditionalSlTriggerPx(req.TriggerPrice.String())
					switch req.OrderType {
					case ORDER_TYPE_LIMIT:
						api.ConditionalSlOrdPx(req.Price.String())
					case ORDER_TYPE_MARKET:
						api.ConditionalSlOrdPx("-1")
					}
				case ORDER_TRIGGER_TYPE_TAKE_PROFIT:
					api.ConditionalTpTriggerPx(req.TriggerPrice.String())
					switch req.OrderType {
					case ORDER_TYPE_LIMIT:
						api.ConditionalTpOrdPx(req.Price.String())
					case ORDER_TYPE_MARKET:
						api.ConditionalTpOrdPx("-1")
					}
				}
			}
		case OKX_ORDER_ALGO_TYPE_OCO:
			// 止盈
			api.ConditionalTpTriggerPx(req.OcoTpTriggerPx.String())
			switch req.OcoTpOrdType {
			case ORDER_TYPE_LIMIT:
				api.ConditionalTpOrdPx(req.OcoTpOrdPx.String())
			case ORDER_TYPE_MARKET:

				api.ConditionalTpOrdPx("-1")
			}
			// 止损
			api.ConditionalSlTriggerPx(req.OcoSlTriggerPx.String())
			switch req.OcoSlOrdType {
			case ORDER_TYPE_LIMIT:
				api.ConditionalSlOrdPx(req.OcoSlOrdPx.String())
			case ORDER_TYPE_MARKET:
				api.ConditionalSlOrdPx("-1")
			}
		}

	}

	if OkxAccountType(req.AccountType) != "SPOT" {
		api.PosSide(o.okxConverter.ToOKXPositionSide(req.PositionSide))
	} else {
		api.TgtCcy("base_ccy")
	}

	if req.ReduceOnly {
		api.ConditionalReduceOnly(req.ReduceOnly)
	}
	if req.ClientOrderId != "" {
		api.AlgoClOrdId(req.ClientOrderId)
	}

	return api
}

// 修改策略委托订单接口
func (o *OkxTradeEngine) apiOrderAlgoAmend(req *OrderParam) *myokxapi.PrivateRestTradeAmendOrderAlgoAPI {
	client := okx.NewRestClient(o.apiKey, o.secretKey, o.passphrase).PrivateRestClient()

	api := client.NewPrivateRestTradeAmendOrderAlgo().
		InstId(req.Symbol)
	if req.OrderId != "" {
		api.AlgoId(req.OrderId)
	}
	if req.ClientOrderId != "" {
		api.AlgoClOrdId(req.ClientOrderId)
	}
	if !req.Quantity.IsZero() {
		api.NewSz(req.Quantity.String())
	}

	switch req.OrderAlgoType {
	case OKX_ORDER_ALGO_TYPE_CONDITIONAL:
		if req.TriggerType != ORDER_TRIGGER_TYPE_UNKNOWN {
			switch req.TriggerType {
			case ORDER_TRIGGER_TYPE_STOP_LOSS:
				api.NewSlTriggerPx(req.TriggerPrice.String())
				switch req.OrderType {
				case ORDER_TYPE_LIMIT:
					api.NewSlOrdPx(req.Price.String())
				case ORDER_TYPE_MARKET:
					api.NewSlOrdPx("-1")
				}
			case ORDER_TRIGGER_TYPE_TAKE_PROFIT:
				api.NewTpTriggerPx(req.TriggerPrice.String())
				switch req.OrderType {
				case ORDER_TYPE_LIMIT:
					api.NewTpOrdPx(req.Price.String())
				case ORDER_TYPE_MARKET:
					api.NewTpOrdPx("-1")
				}
			}
		}
	case OKX_ORDER_ALGO_TYPE_OCO:
		// 止盈
		api.NewTpTriggerPx(req.OcoTpTriggerPx.String())
		switch req.OcoTpOrdType {
		case ORDER_TYPE_LIMIT:
			api.NewTpOrdPx(req.OcoTpOrdPx.String())
		case ORDER_TYPE_MARKET:
			api.NewTpOrdPx("-1")
		}

		// 止损
		api.NewSlTriggerPx(req.OcoSlTriggerPx.String())
		switch req.OcoSlOrdType {
		case ORDER_TYPE_LIMIT:
			api.NewSlOrdPx(req.OcoSlOrdPx.String())
		case ORDER_TYPE_MARKET:
			api.NewSlOrdPx("-1")
		}

	}

	if req.ClientOrderId != "" {
		api.AlgoClOrdId(req.ClientOrderId)
	}

	return api
}

// 取消策略委托订单接口
func (o *OkxTradeEngine) apiOrderAlgoCancel(req *OrderParam) *myokxapi.PrivateRestTradeCancelOrderAlgoAPI {
	client := okx.NewRestClient(o.apiKey, o.secretKey, o.passphrase).PrivateRestClient()
	api := client.NewPrivateRestTradeCancelOrderAlgo()
	cancelOrderAlgo := api.NewCancelOrderAlgo()
	cancelOrderAlgo.SetInstId(req.Symbol).SetAlgoId(req.OrderId)
	api.AddCancelOrderAlgo(cancelOrderAlgo)

	return api
}

// 查询策略委托订单接口
func (o *OkxTradeEngine) apiQueryOpenOrderAlgo(req *QueryOrderParam) *myokxapi.PrivateRestTradeOrderAlgoPendingAPI {
	client := okx.NewRestClient(o.apiKey, o.secretKey, o.passphrase).PrivateRestClient()
	api := client.NewPrivateRestTradeOrderAlgoPending()
	api.OrdType(req.OrderAlgoType.String()) // 策略委托类型

	if req.Symbol != "" {
		api.InstId(req.Symbol)
	}
	if req.OrderId != "" {
		api.AlgoId(req.OrderId)
	}
	if req.ClientOrderId != "" {
		api.AlgoClOrdId(req.ClientOrderId)
	}
	if req.StartTime != 0 {
		api.After(strconv.FormatInt(req.StartTime, 10))
	}
	if req.EndTime != 0 {
		api.Before(strconv.FormatInt(req.EndTime, 10))
	}
	if req.Limit != 0 {
		api.Limit(strconv.Itoa(req.Limit))
	}
	return api
}
func (o *OkxTradeEngine) apiQueryOrderAlgo(req *QueryOrderParam) *myokxapi.PrivateRestTradeOrderAlgoGetAPI {
	client := okx.NewRestClient(o.apiKey, o.secretKey, o.passphrase).PrivateRestClient()
	api := client.NewPrivateRestTradeOrderAlgoGet() // 暂仅支持查询单向/双向止盈止损订单
	if req.OrderId != "" {
		api.AlgoId(req.OrderId)
	}
	if req.ClientOrderId != "" {
		api.AlgoClOrdId(req.ClientOrderId)
	}
	return api
}
func (o *OkxTradeEngine) apiQueryOrdersAlgo(req *QueryOrderParam) *myokxapi.PrivateRestTradeOrderAlgoHistoryAPI {
	client := okx.NewRestClient(o.apiKey, o.secretKey, o.passphrase).PrivateRestClient()
	api := client.NewPrivateRestTradeOrderAlgoHistory()
	api.OrdType(req.OrderAlgoType.String()) // 策略委托类型

	api.State(o.okxConverter.ToOKXOrderStatus(req.AlgoState, true))

	if req.Symbol != "" {
		api.InstId(req.Symbol)
	}
	if req.OrderId != "" {
		api.AlgoId(req.OrderId)
	}
	if req.StartTime != 0 {
		api.After(strconv.FormatInt(req.StartTime, 10))
	}
	if req.EndTime != 0 {
		api.Before(strconv.FormatInt(req.EndTime, 10))
	}
	if req.Limit != 0 {
		api.Limit(strconv.Itoa(req.Limit))
	}
	return api
}

// ws订单请求前置检查
func (o *OkxTradeEngine) wsOrderPreCheck() (bool, error) {
	o.wsForOrderMu.Lock()
	defer o.wsForOrderMu.Unlock()

	if o.wsForOrder == nil {
		newWs := okx.NewPrivateWsStreamClient()
		err := newWs.OpenConn()
		if err != nil {
			return false, err
		}

		err = newWs.Login(okx.NewRestClient(o.apiKey, o.secretKey, o.passphrase))
		if err != nil {
			return false, err
		}
		o.wsForOrder = newWs
	}
	return true, nil
}

func (o *OkxTradeEngine) restBatchPreCheck(reqs []*OrderParam) error {
	//检测长度，OKX最多批量下20个订单
	if len(reqs) > 20 {
		return ErrorInvalid("okx order param length require less than 20")

	}

	//检测类型是否相同
	for _, req := range reqs {
		if err := o.accountTypePreCheck(req.AccountType); err != nil {
			return err
		}
	}

	return nil
}

func (o *OkxTradeEngine) accountTypePreCheck(accountType string) error {
	switch OkxAccountType(accountType) {
	case OKX_AC_SPOT, OKX_AC_MARGIN, OKX_AC_SWAP, OKX_AC_FUTURES:
	default:
		return ErrorAccountType
	}
	return nil
}
