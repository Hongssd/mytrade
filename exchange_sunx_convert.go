package mytrade

type SunxEnumConverter struct{}

// 账户模式
func (c *SunxEnumConverter) ToSunxAccountMode() string {
	return SUNX_ACCOUNT_MODE_UNIFIED
}

func (c *SunxEnumConverter) FromSunxAccountMode() AccountMode {
	return ACCOUNT_MODE_UNIFIED
}

// 保证金模式（仅支持全仓 cross)
func (c *SunxEnumConverter) FromSunxMarginMode(t string) MarginMode {
	switch t {
	case SUNX_MARGIN_MODE_CROSSED:
		return MARGIN_MODE_CROSSED
	default:
		return MARGIN_MODE_UNKNOWN
	}
}
func (c *SunxEnumConverter) ToSunxMarginMode(marginMode MarginMode) string {
	switch marginMode {
	case MARGIN_MODE_CROSSED:
		return SUNX_MARGIN_MODE_CROSSED
	default:
		return ""
	}
}

// 仓位模式
func (c *SunxEnumConverter) FromSunxPositionMode(t string) PositionMode {
	switch t {
	case SUNX_POSITION_MODE_SINGLE:
		return POSITION_MODE_ONEWAY
	case SUNX_POSITION_MODE_HEDGE:
		return POSITION_MODE_HEDGE
	default:
		return POSITION_MODE_UNKNOWN
	}
}

func (c *SunxEnumConverter) ToSunxPositionMode(t PositionMode) string {
	switch t {
	case POSITION_MODE_ONEWAY:
		return SUNX_POSITION_MODE_SINGLE
	case POSITION_MODE_HEDGE:
		return SUNX_POSITION_MODE_HEDGE
	default:
		return ""
	}
}

// 订单状态
func (c *SunxEnumConverter) FromSunxOrderStatus(t string) OrderStatus {
	switch t {
	case SUNX_ORDER_STATUS_NEW:
		return ORDER_STATUS_NEW
	case SUNX_ORDER_STATUS_PARTIALLY_FILLED:
		return ORDER_STATUS_PARTIALLY_FILLED
	case SUNX_ORDER_STATUS_FILLED:
		return ORDER_STATUS_FILLED
	case SUNX_ORDER_STATUS_PARTIALLY_CANCELED:
		return ORDER_STATUS_CANCELED
	case SUNX_ORDER_STATUS_CANCELED:
		return ORDER_STATUS_CANCELED
	case SUNX_ORDER_STATUS_REJECTED:
		return ORDER_STATUS_REJECTED
	default:
		return ORDER_STATUS_UNKNOWN
	}
}

// 订单类型
func (c *SunxEnumConverter) FromSunxOrderType(t string) OrderType {
	switch t {
	case SUNX_ORDER_TYPE_LIMIT:
		return ORDER_TYPE_LIMIT
	case SUNX_ORDER_TYPE_MARKET:
		return ORDER_TYPE_MARKET
	default:
		return ORDER_TYPE_UNKNOWN
	}
}

func (c *SunxEnumConverter) ToSunxOrderType(t OrderType, tif TimeInForce) string {
	if tif == TIME_IN_FORCE_POST_ONLY {
		return SUNX_ORDER_TYPE_POST_ONLY
	}

	switch t {
	case ORDER_TYPE_LIMIT:
		return SUNX_ORDER_TYPE_LIMIT
	case ORDER_TYPE_MARKET:
		return SUNX_ORDER_TYPE_MARKET
	default:
		return ""
	}
}

// 订单方向
func (c *SunxEnumConverter) FromSunxOrderSide(t string) OrderSide {
	switch t {
	case SUNX_ORDER_SIDE_BUY:
		return ORDER_SIDE_BUY
	case SUNX_ORDER_SIDE_SELL:
		return ORDER_SIDE_SELL
	default:
		return ORDER_SIDE_UNKNOWN
	}
}

func (c *SunxEnumConverter) ToSunxOrderSide(t OrderSide) string {
	switch t {
	case ORDER_SIDE_BUY:
		return SUNX_ORDER_SIDE_BUY
	case ORDER_SIDE_SELL:
		return SUNX_ORDER_SIDE_SELL
	default:
		return ""
	}
}

// 仓位方向
func (c *SunxEnumConverter) FromSunxPositionSide(t string) PositionSide {
	switch t {
	case SUNX_POSITION_SIDE_LONG:
		return POSITION_SIDE_LONG
	case SUNX_POSITION_SIDE_SHORT:
		return POSITION_SIDE_SHORT
	case SUNX_POSITION_SIDE_BOTH:
		return POSITION_SIDE_BOTH
	default:
		return POSITION_SIDE_UNKNOWN
	}
}

func (c *SunxEnumConverter) ToSunxPositionSide(t PositionSide) string {
	switch t {
	case POSITION_SIDE_LONG:
		return SUNX_POSITION_SIDE_LONG
	case POSITION_SIDE_SHORT:
		return SUNX_POSITION_SIDE_SHORT
	case POSITION_SIDE_BOTH:
		return SUNX_POSITION_SIDE_BOTH
	default:
		return ""
	}
}

// 时间InForce
func (c *SunxEnumConverter) FromSunxTimeInForce(ot string, tif string) (OrderType, TimeInForce) {
	if ot == SUNX_ORDER_TYPE_POST_ONLY {
		return ORDER_TYPE_LIMIT, TIME_IN_FORCE_POST_ONLY
	}
	orderType := c.FromSunxOrderType(ot)
	switch tif {
	case SUNX_TIME_IN_FORCE_GTC:
		return orderType, TIME_IN_FORCE_GTC
	case SUNX_TIME_IN_FORCE_IOC:
		return orderType, TIME_IN_FORCE_IOC
	case SUNX_TIME_INFORCE_FOK:
		return orderType, TIME_IN_FORCE_FOK
	}
	return ORDER_TYPE_UNKNOWN, TIME_IN_FORCE_UNKNOWN
}

func (c *SunxEnumConverter) ToSunxTimeInForce(t TimeInForce) string {
	switch t {
	case TIME_IN_FORCE_GTC:
		return SUNX_TIME_IN_FORCE_GTC
	case TIME_IN_FORCE_IOC:
		return SUNX_TIME_IN_FORCE_IOC
	case TIME_IN_FORCE_FOK:
		return SUNX_TIME_INFORCE_FOK
	case TIME_IN_FORCE_POST_ONLY:
		return ""
	default:
		return ""
	}
}
