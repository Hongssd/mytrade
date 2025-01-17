package mytrade

// 枚举转换器
type GateEnumConverter struct{}

func (c *GateEnumConverter) FromGateAccountMode(t int) AccountMode {
	switch t {
	case GATE_ACCOUNT_MODE_FREE_MARGIN:
		return ACCOUNT_MODE_FREE_MARGIN
	case GATE_ACCOUNT_MODE_MULTI_MARGIN:
		return ACCOUNT_MODE_MULTI_MARGIN
	}
	return ACCOUNT_MODE_UNKNOWN
}

func (c *GateEnumConverter) FromGatePositionMode(t string) PositionMode {
	switch t {
	case GATE_POSITION_MODE_ONEWAY:
		return POSITION_MODE_ONEWAY
	case GATE_POSITION_MODE_HEDGE_LONG, GATE_POSITION_MODE_HEDGE_SHORT:
		return POSITION_MODE_HEDGE
	}
	return POSITION_MODE_UNKNOWN
}

func (c *GateEnumConverter) ToGateAssetType(t AssetType) string {
	switch t {
	case ASSET_TYPE_FUND:
		return GATE_ASSET_TYPE_SPOT
	case ASSET_TYPE_MARGIN:
		return GATE_ASSET_TYPE_MARGIN
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
	case GATE_ASSET_TYPE_MARGIN:
		return ASSET_TYPE_MARGIN
	case GATE_ASSET_TYPE_FUTURES:
		return ASSET_TYPE_UMFUTURE
	case GATE_ASSET_TYPE_DELIVERY:
		return ASSET_TYPE_DELIVERY
	case GATE_ASSET_TYPE_UNFIED:
		return ASSET_TYPE_UNIFIED
	}
	return ""
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
	}
	return TIME_IN_FORCE_UNKNOWN
}

func (c *GateEnumConverter) FromGateOrderStatus(t string) OrderStatus {
	switch t {
	case GATE_ORDER_STATUS_NEW:
		return ORDER_STATUS_NEW
	case GATE_ORDER_STATUS_FILLED:
		return ORDER_STATUS_FILLED
	case GATE_ORDER_STATUS_CANCELLED:
		return ORDER_STATUS_CANCELED
	}
	return ORDER_STATUS_UNKNOWN
}

func (c *GateEnumConverter) ToGateOrderStatus(t OrderStatus) string {
	switch t {
	case ORDER_STATUS_NEW:
		return GATE_ORDER_STATUS_NEW
	case ORDER_STATUS_FILLED:
		return GATE_ORDER_STATUS_FILLED
	case ORDER_STATUS_CANCELED:
		return GATE_ORDER_STATUS_CANCELLED
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
