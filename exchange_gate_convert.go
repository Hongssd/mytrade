package mytrade

// 枚举转换器
type GateEnumConverter struct{}

func (c *GateEnumConverter) FromGateAccountMode(t string) AccountMode {
	switch t {
	case GATE_ACCOUNT_MODE_CLASSIC:
		return ACCOUNT_MODE_FREE_MARGIN
	case GATE_ACCOUNT_MODE_SINGLE_MARGIN:
		return ACCOUNT_MODE_SINGLE_MARGIN
	case GATE_ACCOUNT_MODE_MULTI_MARGIN:
		return ACCOUNT_MODE_MULTI_MARGIN
	case GATE_ACCOUNT_MODE_PORTFOLIO:
		return ACCOUNT_MODE_PORTFOLIO
	}
	return ACCOUNT_MODE_UNKNOWN
}

func (c *GateEnumConverter) ToGateAccountMode(t AccountMode) (string, bool, bool) {
	switch t {
	case ACCOUNT_MODE_FREE_MARGIN:
		return GATE_ACCOUNT_MODE_CLASSIC, false, false
	case ACCOUNT_MODE_SINGLE_MARGIN:
		//单币种保证金
		return GATE_ACCOUNT_MODE_SINGLE_MARGIN, false, false
	case ACCOUNT_MODE_MULTI_MARGIN:
		//跨币种保证金
		return GATE_ACCOUNT_MODE_MULTI_MARGIN, true, true
	case ACCOUNT_MODE_PORTFOLIO:
		//组合保证金
		return GATE_ACCOUNT_MODE_PORTFOLIO, true, true
	}
	return "", false, false
}

func (c *GateEnumConverter) FromGatePositionMode(t bool) PositionMode {
	switch t {
	case GATE_POSITION_MODE_ONEWAY:

		return POSITION_MODE_ONEWAY
	case GATE_POSITION_MODE_HEDGE:
		return POSITION_MODE_HEDGE
	}
	return POSITION_MODE_UNKNOWN
}

func (c *GateEnumConverter) ToGatePositionMode(t PositionMode) bool {
	switch t {
	case POSITION_MODE_ONEWAY:
		return GATE_POSITION_MODE_ONEWAY
	case POSITION_MODE_HEDGE:
		return GATE_POSITION_MODE_HEDGE
	}
	return GATE_POSITION_MODE_ONEWAY
}

func (c *GateEnumConverter) ToGateAssetType(t AssetType) string {
	switch t {
	case ASSET_TYPE_FUND:
		return GATE_ASSET_TYPE_SPOT
	case ASSET_TYPE_MARGIN_ISOLATED:
		return GATE_ASSET_TYPE_ISOLATED_MARGIN
	case ASSET_TYPE_MARGIN_CROSSED:
		return GATE_ASSET_TYPE_CROSS_MARGIN
	case ASSET_TYPE_UMFUTURE:
		return GATE_ASSET_TYPE_FUTURES
	case ASSET_TYPE_DELIVERY:
		return GATE_ASSET_TYPE_DELIVERY
	case ASSET_TYPE_UNIFIED:
		return GATE_ASSET_TYPE_UNFIED
	}
	return ""
}

func (c *GateEnumConverter) FromGateAssetType(t string) AssetType {
	switch t {
	case GATE_ASSET_TYPE_SPOT:
		return ASSET_TYPE_FUND
	case GATE_ASSET_TYPE_ISOLATED_MARGIN:
		return ASSET_TYPE_MARGIN_ISOLATED
	case GATE_ASSET_TYPE_CROSS_MARGIN:
		return ASSET_TYPE_MARGIN_CROSSED
	case GATE_ASSET_TYPE_FUTURES:
		return ASSET_TYPE_UMFUTURE
	case GATE_ASSET_TYPE_DELIVERY:
		return ASSET_TYPE_DELIVERY
	case GATE_ASSET_TYPE_UNFIED:
		return ASSET_TYPE_UNIFIED
	}
	return ""
}
func (c *GateEnumConverter) ToOrderSpotAccountType(t GateAccountType, isMargin, isIsolated bool) GateAccountType {
	if t == GATE_ACCOUNT_TYPE_SPOT {
		if isMargin {
			if isIsolated {
				return GATE_ACCOUNT_TYPE_MARGIN //逐仓杠杆
			} else {
				return GATE_ACCOUNT_TYPE_CROSS_MARGIN //全倉杠杆
			}
		} else {
			return GATE_ACCOUNT_TYPE_SPOT //现货
		}
	}
	return GATE_ACCOUNT_TYPE_UNKNOWN
}
func (c *GateEnumConverter) FromOrderSpotAccountType(t GateAccountType) (GateAccountType, bool, bool) {
	//return accountType, isMargin, isIsolated
	switch t {
	case GATE_ACCOUNT_TYPE_SPOT, GATE_ACCOUNT_TYPE_UNIFIED:
		return GATE_ACCOUNT_TYPE_SPOT, false, false
	case GATE_ACCOUNT_TYPE_MARGIN:
		return GATE_ACCOUNT_TYPE_SPOT, true, true
	case GATE_ACCOUNT_TYPE_CROSS_MARGIN:
		return GATE_ACCOUNT_TYPE_SPOT, true, false
	}
	return GATE_ACCOUNT_TYPE_UNKNOWN, false, false
}

func (c *GateEnumConverter) FromOrderSpotPriceAccountType(t string) (GateAccountType, bool, bool) {
	//return accountType, isMargin, isIsolated
	switch t {
	case GATE_SPOT_PRICE_ORDER_ACCOUNT_NORMAL:
		return GATE_ACCOUNT_TYPE_SPOT, false, false
	case GATE_SPOT_PRICE_ORDER_ACCOUNT_MARGIN:
		return GATE_ACCOUNT_TYPE_SPOT, true, true
	case GATE_SPOT_PRICE_ORDER_ACCOUNT_CROSSED_MARGIN:
		return GATE_ACCOUNT_TYPE_SPOT, true, false
	}
	return GATE_ACCOUNT_TYPE_UNKNOWN, false, false
}

func (c *GateEnumConverter) ToGateFuturesPriceOrderTriggerRule(t OrderTriggerType, s OrderSide) int32 {
	switch t {
	case ORDER_TRIGGER_TYPE_TAKE_PROFIT:
		switch s {
		case ORDER_SIDE_BUY:
			return GATE_FUTURES_PRICE_ORDER_TRIGGER_RULE_LTE
		case ORDER_SIDE_SELL:
			return GATE_FUTURES_PRICE_ORDER_TRIGGER_RULE_GTE
		}
	case ORDER_TRIGGER_TYPE_STOP_LOSS:
		switch s {
		case ORDER_SIDE_BUY:
			return GATE_FUTURES_PRICE_ORDER_TRIGGER_RULE_GTE
		case ORDER_SIDE_SELL:
			return GATE_FUTURES_PRICE_ORDER_TRIGGER_RULE_LTE
		}
	}
	return 0
}

func (c *GateEnumConverter) FromGateFuturesPriceOrderTriggerRule(t int64, s OrderSide) OrderTriggerType {
	switch t {
	case GATE_FUTURES_PRICE_ORDER_TRIGGER_RULE_GTE:
		switch s {
		case ORDER_SIDE_BUY:
			//价格上穿时买入，止损
			return ORDER_TRIGGER_TYPE_STOP_LOSS
		case ORDER_SIDE_SELL:
			//价格上穿时卖出，止盈
			return ORDER_TRIGGER_TYPE_TAKE_PROFIT
		}
	case GATE_FUTURES_PRICE_ORDER_TRIGGER_RULE_LTE:
		switch s {
		case ORDER_SIDE_BUY:
			//价格下穿时买入，止盈
			return ORDER_TRIGGER_TYPE_TAKE_PROFIT
		case ORDER_SIDE_SELL:
			//价格下穿时卖出，止损
			return ORDER_TRIGGER_TYPE_STOP_LOSS
		}
	}
	return ORDER_TRIGGER_TYPE_UNKNOWN
}

func (c *GateEnumConverter) ToGateOrderType(t OrderType) string {
	switch t {
	case ORDER_TYPE_LIMIT:
		return GATE_ORDER_TYPE_LIMIT
	case ORDER_TYPE_MARKET:
		return GATE_ORDER_TYPE_MARKET
	}
	return ""
}

func (c *GateEnumConverter) FromGateOrderType(t string) OrderType {
	switch t {
	case GATE_ORDER_TYPE_LIMIT:
		return ORDER_TYPE_LIMIT
	case GATE_ORDER_TYPE_MARKET:
		return ORDER_TYPE_MARKET
	}
	return ORDER_TYPE_UNKNOWN
}
func (c *GateEnumConverter) ToGateOrderSide(t OrderSide) string {
	switch t {
	case ORDER_SIDE_BUY:
		return GATE_ORDER_SIDE_BUY
	case ORDER_SIDE_SELL:
		return GATE_ORDER_SIDE_SELL
	}
	return ""
}

func (c *GateEnumConverter) FromGateOrderSide(t string) OrderSide {
	switch t {
	case GATE_ORDER_SIDE_BUY:
		return ORDER_SIDE_BUY
	case GATE_ORDER_SIDE_SELL:
		return ORDER_SIDE_SELL
	}
	return ORDER_SIDE_UNKNOWN
}

func (c *GateEnumConverter) ToGateTimeInForce(t TimeInForce) string {
	switch t {
	case TIME_IN_FORCE_GTC:
		return GATE_TIME_IN_FORCE_GTC
	case TIME_IN_FORCE_IOC:
		return GATE_TIME_IN_FORCE_IOC
	case TIME_IN_FORCE_POST_ONLY:
		return GATE_TIME_IN_FORCE_POC
	case TIME_IN_FORCE_FOK:
		return GATE_TIME_IN_FORCE_FOK
	}
	return ""
}

func (c *GateEnumConverter) FromGateTimeInForce(t string) TimeInForce {
	switch t {
	case GATE_TIME_IN_FORCE_GTC:
		return TIME_IN_FORCE_GTC
	case GATE_TIME_IN_FORCE_IOC:
		return TIME_IN_FORCE_IOC
	case GATE_TIME_IN_FORCE_POC:
		return TIME_IN_FORCE_POST_ONLY
	case GATE_TIME_IN_FORCE_FOK:
		return TIME_IN_FORCE_FOK
	}
	return TIME_IN_FORCE_UNKNOWN
}
func (c *GateEnumConverter) ToGateSpotOrderStatus(t OrderStatus) string {
	switch t {
	case ORDER_STATUS_NEW:
		return GATE_ORDER_SPOT_STATUS_NEW
	case ORDER_STATUS_FILLED:
		return GATE_ORDER_SPOT_STATUS_FILLED
	case ORDER_STATUS_CANCELED:
		return GATE_ORDER_SPOT_STATUS_CANCELLED
	}
	return ""
}
func (c *GateEnumConverter) FromGateSpotOrderStatus(t string) OrderStatus {
	switch t {
	case GATE_ORDER_SPOT_STATUS_NEW:
		return ORDER_STATUS_NEW
	case GATE_ORDER_SPOT_STATUS_FILLED:
		return ORDER_STATUS_FILLED
	case GATE_ORDER_SPOT_STATUS_CANCELLED:
		return ORDER_STATUS_CANCELED
	}
	return ORDER_STATUS_UNKNOWN
}
func (c *GateEnumConverter) FromGateWsSportOrderStatus(event, finishAs string) OrderStatus {
	switch event {
	case GATE_ORDER_WS_SPOT_EVENT_PUT:
		return ORDER_STATUS_NEW
	case GATE_ORDER_WS_SPOT_EVENT_UPDATE, GATE_ORDER_WS_SPOT_EVENT_FINISH:
		switch finishAs {
		case GATE_ORDER_WS_SPOT_FINISH_AS_OPEN:
			return ORDER_STATUS_NEW
		case GATE_ORDER_WS_SPOT_FINISH_AS_FILLED:
			return ORDER_STATUS_FILLED
		case GATE_ORDER_WS_SPOT_FINISH_AS_CANCELLED,
			GATE_ORDER_WS_SPOT_FINISH_AS_IOC,
			GATE_ORDER_WS_SPOT_FINISH_AS_STP,
			GATE_ORDER_WS_SPOT_FINISH_AS_POC,
			GATE_ORDER_WS_SPOT_FINISH_AS_FOK,
			GATE_ORDER_WS_SPOT_FINISH_AS_TRADER_NOT_ENOUGH,
			GATE_ORDER_WS_SPOT_FINISH_AS_DEPTH_NOT_ENOUGH,
			GATE_ORDER_WS_SPOT_FINISH_AS_LIQUIDATE_CANCELLED,
			GATE_ORDER_WS_SPOT_FINISH_AS_SMALL:
			return ORDER_STATUS_CANCELED
		default:
			return ORDER_STATUS_UNKNOWN
		}
	}
	return ORDER_STATUS_UNKNOWN
}

func (c *GateEnumConverter) FromGateContractOrderStatus(t, fas string) OrderStatus {
	switch t {
	case GATE_ORDER_CONTRACT_STATUS_OPEN:
		return ORDER_STATUS_NEW
	case GATE_ORDER_CONTRACT_STATUS_FINISHED:
		switch fas {
		case GATE_ORDER_CONTRACT_FINISH_AS_NEW,
			GATE_ORDER_CONTRACT_FINISH_AS_UPDATE,
			GATE_ORDER_CONTRACT_FINISH_AS_REDUCE_OUT:
			return ORDER_STATUS_NEW
		case GATE_ORDER_CONTRACT_FINISH_AS_FILLED:
			return ORDER_STATUS_FILLED
		case GATE_ORDER_CONTRACT_FINISH_AS_CANCELLED,
			GATE_ORDER_CONTRACT_FINISH_AS_LIQUIDATED,
			GATE_ORDER_CONTRACT_FINISH_AS_IOC,
			GATE_ORDER_CONTRACT_FINISH_AS_AUTO_DELEVERAGED,
			GATE_ORDER_CONTRACT_FINISH_AS_REDUCE_ONLY,
			GATE_ORDER_CONTRACT_FINISH_AS_POSITION_CLOSE,
			GATE_ORDER_CONTRACT_FINISH_AS_STP:
			return ORDER_STATUS_CANCELED
		}
		return ORDER_STATUS_FILLED
	}
	return ORDER_STATUS_UNKNOWN
}

func (c *GateEnumConverter) ToGatePositionSide(t PositionSide) string {
	switch t {
	case POSITION_SIDE_LONG:
		return GATE_POSITION_SIDE_LONG
	case POSITION_SIDE_SHORT:
		return GATE_POSITION_SIDE_SHORT
	case POSITION_SIDE_BOTH:
		return GATE_POSITION_SIDE_BOTH
	}
	return ""
}

func (c *GateEnumConverter) FromGatePositionSide(t string) PositionSide {
	switch t {
	case GATE_POSITION_SIDE_BOTH:
		return POSITION_SIDE_BOTH
	case GATE_POSITION_SIDE_LONG:
		return POSITION_SIDE_LONG
	case GATE_POSITION_SIDE_SHORT:
		return POSITION_SIDE_SHORT
	}
	return POSITION_SIDE_UNKNOWN
}

func (c *GateEnumConverter) ToGateSpotPriceOrderTriggerRule(t OrderTriggerType, s OrderSide) string {
	switch t {
	case ORDER_TRIGGER_TYPE_TAKE_PROFIT:
		switch s {
		case ORDER_SIDE_BUY:
			return "<="
		case ORDER_SIDE_SELL:
			return ">="
		}
	case ORDER_TRIGGER_TYPE_STOP_LOSS:
		switch s {
		case ORDER_SIDE_BUY:
			return ">="
		case ORDER_SIDE_SELL:
			return "<="
		}
	}
	return ""
}

func (c *GateEnumConverter) FromGateSpotPriceOrderTriggerRule(rule string, s OrderSide) OrderTriggerType {
	switch rule {
	case GATE_SPOT_PRICE_ORDER_TRIGGER_RULE_LTE:
		switch s {
		case ORDER_SIDE_BUY:
			return ORDER_TRIGGER_TYPE_TAKE_PROFIT
		case ORDER_SIDE_SELL:
			return ORDER_TRIGGER_TYPE_STOP_LOSS
		}
	case GATE_SPOT_PRICE_ORDER_TRIGGER_RULE_GTE:
		switch s {
		case ORDER_SIDE_BUY:
			return ORDER_TRIGGER_TYPE_STOP_LOSS
		case ORDER_SIDE_SELL:
			return ORDER_TRIGGER_TYPE_TAKE_PROFIT
		}
	}
	return ""
}

func (c *GateEnumConverter) FromGateSpotPriceOrderStatus(t string) OrderStatus {
	switch t {
	case GATE_ORDER_SPOT_PRICE_STATUS_OPEN:
		return ORDER_STATUS_UN_TRIGGERED
	case GATE_ORDER_SPOT_PRICE_STATUS_CANCELED:
		return ORDER_STATUS_CANCELED
	case GATE_ORDER_SPOT_PRICE_STATUS_FINISHED:
		return ORDER_STATUS_TRIGGERED
	case GATE_ORDER_SPOT_PRICE_STATUS_FAILED:
		return ORDER_STATUS_CANCELED
	case GATE_ORDER_SPOT_PRICE_STATUS_EXPIRED:
		return ORDER_STATUS_CANCELED
	}
	return ORDER_STATUS_UNKNOWN
}
func (c *GateEnumConverter) FromGateContractPriceOrderStatus(t string, fas string) OrderStatus {
	switch t {
	case GATE_ORDER_FUTURES_PRICE_STATUS_OPEN, GATE_ORDER_FUTURES_PRICE_STATUS_INACTIVE:
		return ORDER_STATUS_UN_TRIGGERED
	case GATE_ORDER_FUTURES_PRICE_STATUS_FINISHED:
		switch fas {
		case GATE_ORDER_FUTURES_PRICE_FINISH_AS_CANCELLED,
			GATE_ORDER_FUTURES_PRICE_FINISH_AS_FAILED,
			GATE_ORDER_FUTURES_PRICE_FINISH_AS_EXPIRED:
			return ORDER_STATUS_CANCELED
		case GATE_ORDER_FUTURES_PRICE_FINISH_AS_SUCCEEDED:
			return ORDER_STATUS_TRIGGERED
		}
	case GATE_ORDER_FUTURES_PRICE_STATUS_INVALID:
		return ORDER_STATUS_REJECTED
	}
	return ORDER_STATUS_UNKNOWN
}

func (c *GateEnumConverter) ToGateSpotPriceOrderAccount(t GateAccountType) string {
	switch t {
	case GATE_ACCOUNT_TYPE_SPOT:
		return GATE_SPOT_PRICE_ORDER_ACCOUNT_NORMAL
	case GATE_ACCOUNT_TYPE_MARGIN:
		return GATE_SPOT_PRICE_ORDER_ACCOUNT_MARGIN
	case GATE_ACCOUNT_TYPE_CROSS_MARGIN:
		return GATE_SPOT_PRICE_ORDER_ACCOUNT_CROSSED_MARGIN
	}
	return ""
}
func (c *GateEnumConverter) FromGateSpotPriceOrderAccount(t string) (GateAccountType, bool, bool) {
	switch t {
	case GATE_SPOT_PRICE_ORDER_ACCOUNT_NORMAL:
		return GATE_ACCOUNT_TYPE_SPOT, false, false
	case GATE_SPOT_PRICE_ORDER_ACCOUNT_MARGIN:
		return GATE_ACCOUNT_TYPE_SPOT, true, true
	case GATE_SPOT_PRICE_ORDER_ACCOUNT_CROSSED_MARGIN:
		return GATE_ACCOUNT_TYPE_SPOT, true, false
	}
	return GATE_ACCOUNT_TYPE_UNKNOWN, false, false
}

func (c *GateEnumConverter) FromGateTriggerCondition(OrderSide OrderSide, triggerType OrderTriggerType) OrderTriggerConditionType {
	switch OrderSide {
	case ORDER_SIDE_BUY:
		switch triggerType {
		case ORDER_TRIGGER_TYPE_STOP_LOSS:
			//止损买入时高价 价格上穿触发
			return ORDER_TRIGGER_CONDITION_TYPE_THROUGH_UP
		case ORDER_TRIGGER_TYPE_TAKE_PROFIT:
			//止盈买入时低价 价格下穿触发
			return ORDER_TRIGGER_CONDITION_TYPE_THROUGH_DOWN
		}
	case ORDER_SIDE_SELL:
		switch triggerType {
		case ORDER_TRIGGER_TYPE_STOP_LOSS:
			//止损卖出时低价 价格下穿触发
			return ORDER_TRIGGER_CONDITION_TYPE_THROUGH_DOWN
		case ORDER_TRIGGER_TYPE_TAKE_PROFIT:
			//止盈卖出时高价 价格上穿触发
			return ORDER_TRIGGER_CONDITION_TYPE_THROUGH_UP
		}
	}
	return ""
}
