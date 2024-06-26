package mytrade

// 枚举转换器
type BybitEnumConverter struct{}

// 订单方向
func (c *BybitEnumConverter) FromBYBITOrderSide(t string) OrderSide {
	switch t {
	case BYBIT_ORDER_SIDE_BUY:
		return ORDER_SIDE_BUY
	case BYBIT_ORDER_SIDE_SELL:
		return ORDER_SIDE_SELL
	default:
		return ORDER_SIDE_UNKNOWN
	}
}
func (c *BybitEnumConverter) ToBYBITOrderSide(t OrderSide) string {
	switch t {
	case ORDER_SIDE_BUY:
		return BYBIT_ORDER_SIDE_BUY
	case ORDER_SIDE_SELL:
		return BYBIT_ORDER_SIDE_SELL
	default:
		return ""
	}
}

// 订单类型
func (c *BybitEnumConverter) FromBYBITOrderType(t string) OrderType {
	switch t {
	case BYBIT_ORDER_TYPE_LIMIT:
		return ORDER_TYPE_LIMIT
	case BYBIT_ORDER_TYPE_MARKET:
		return ORDER_TYPE_MARKET
	default:
		return ORDER_TYPE_UNKNOWN
	}
}
func (c *BybitEnumConverter) ToBYBITOrderType(t OrderType) string {
	switch t {
	case ORDER_TYPE_LIMIT:
		return BYBIT_ORDER_TYPE_LIMIT
	case ORDER_TYPE_MARKET:
		return BYBIT_ORDER_TYPE_MARKET
	default:
		return ""
	}
}

// 仓位方向
func (c *BybitEnumConverter) FromBYBITPositionSide(t int) PositionSide {
	switch t {
	case BYBIT_POSITION_SIDE_LONG:
		return POSITION_SIDE_LONG
	case BYBIT_POSITION_SIDE_SHORT:
		return POSITION_SIDE_SHORT
	case BYBIT_POSITION_SIDE_BOTH:
		return POSITION_SIDE_BOTH
	default:
		return POSITION_SIDE_UNKNOWN
	}
}
func (c *BybitEnumConverter) ToBYBITPositionSide(t PositionSide) int {
	switch t {
	case POSITION_SIDE_LONG:
		return BYBIT_POSITION_SIDE_LONG
	case POSITION_SIDE_SHORT:
		return BYBIT_POSITION_SIDE_SHORT
	case POSITION_SIDE_BOTH:
		return BYBIT_POSITION_SIDE_BOTH
	}
	return 0
}

// 有效方式
func (c *BybitEnumConverter) FromBYBITTimeInForce(t string) TimeInForce {
	switch t {
	case BYBIT_TIME_IN_FORCE_GTC:
		return TIME_IN_FORCE_GTC
	case BYBIT_TIME_IN_FORCE_IOC:
		return TIME_IN_FORCE_IOC
	case BYBIT_TIME_IN_FORCE_FOK:
		return TIME_IN_FORCE_FOK
	case BYBIT_TIME_IN_FORCE_POST_ONLY:
		return TIME_IN_FORCE_POST_ONLY
	default:
		return TIME_IN_FORCE_UNKNOWN
	}
}
func (c *BybitEnumConverter) ToBYBITTimeInForce(t TimeInForce) string {
	switch t {
	case TIME_IN_FORCE_GTC:
		return BYBIT_TIME_IN_FORCE_GTC
	case TIME_IN_FORCE_IOC:
		return BYBIT_TIME_IN_FORCE_IOC
	case TIME_IN_FORCE_FOK:
		return BYBIT_TIME_IN_FORCE_FOK
	case TIME_IN_FORCE_POST_ONLY:
		return BYBIT_TIME_IN_FORCE_POST_ONLY
	default:
		return ""
	}
}

// 订单状态
func (c *BybitEnumConverter) FromBYBITOrderStatus(t string) OrderStatus {
	switch t {
	case BYBIT_ORDER_STATUS_NEW:
		return ORDER_STATUS_NEW
	case BYBIT_ORDER_STATUS_PARTIALLY_FILLED:
		return ORDER_STATUS_PARTIALLY_FILLED
	case BYBIT_ORDER_STATUS_FILLED:
		return ORDER_STATUS_FILLED
	case BYBIT_ORDER_STATUS_CANCELED, BYBIT_ORDER_STATUS_PARTIALLY_FILLED_CANCELED:
		return ORDER_STATUS_CANCELED
	case BYBIT_ORDER_STATUS_REJECTED:
		return ORDER_STATUS_REJECTED
	default:
		return ORDER_STATUS_UNKNOWN
	}
}
func (c *BybitEnumConverter) ToBYBITOrderStatus(t OrderStatus) string {
	switch t {
	case ORDER_STATUS_NEW:
		return BYBIT_ORDER_STATUS_NEW
	case ORDER_STATUS_PARTIALLY_FILLED:
		return BYBIT_ORDER_STATUS_PARTIALLY_FILLED
	case ORDER_STATUS_FILLED:
		return BYBIT_ORDER_STATUS_FILLED
	case ORDER_STATUS_CANCELED:
		return BYBIT_ORDER_STATUS_CANCELED
	case ORDER_STATUS_REJECTED:
		return BYBIT_ORDER_STATUS_REJECTED
	default:
		return ""
	}
}

// 账户模式
func (c *BybitEnumConverter) FromBYBITAccountMode(t string) AccountMode {
	switch t {
	case BYBIT_ACCOUNT_MOED_ISOLATED_MARGIN:
		return ACCOUNT_MODE_SINGLE_MARGIN
	case BYBIT_ACCOUNT_MODE_REGULAR_MARGIN:
		return ACCOUNT_MODE_MULTI_MARGIN
	case BYBIT_ACCOUNT_MODE_PORTFOLIO_MARGIN:
		return ACCOUNT_MODE_PORTFOLIO
	default:
		return ACCOUNT_MODE_UNKNOWN
	}
}
func (c *BybitEnumConverter) ToBYBITAccountMode(t AccountMode) string {
	switch t {
	case ACCOUNT_MODE_SINGLE_MARGIN:
		return BYBIT_ACCOUNT_MOED_ISOLATED_MARGIN
	case ACCOUNT_MODE_FREE_MARGIN, ACCOUNT_MODE_MULTI_MARGIN:
		return BYBIT_ACCOUNT_MODE_REGULAR_MARGIN
	case ACCOUNT_MODE_PORTFOLIO:
		return BYBIT_ACCOUNT_MODE_PORTFOLIO_MARGIN
	default:
		return ""
	}
}

// 保证金模式
func (c *BybitEnumConverter) FromBYBITMarginMode(t int) MarginMode {
	switch t {
	case BYBIT_MARGIN_MODE_ISOLATED:
		return MARGIN_MODE_ISOLATED
	case BYBIT_MARGIN_MODE_CROSSED:
		return MARGIN_MODE_CROSSED
	default:
		return MARGIN_MODE_UNKNOWN
	}
}
func (c *BybitEnumConverter) ToBYBITMarginMode(t MarginMode) int {
	switch t {
	case MARGIN_MODE_ISOLATED:
		return BYBIT_MARGIN_MODE_ISOLATED
	case MARGIN_MODE_CROSSED:
		return BYBIT_MARGIN_MODE_CROSSED
	default:
		return 0
	}
}

// 仓位模式
func (c *BybitEnumConverter) FromBYBITPositionMode(t int) PositionMode {
	switch t {
	case BYBIT_POSITION_MODE_HEDGE:
		return POSITION_MODE_HEDGE
	case BYBIT_POSITION_MODE_ONEWAY:
		return POSITION_MODE_ONEWAY
	default:
		return POSITION_MODE_UNKNOWN
	}
}
func (c *BybitEnumConverter) ToBYBITPositionMode(t PositionMode) int {
	switch t {
	case POSITION_MODE_HEDGE:
		return BYBIT_POSITION_MODE_HEDGE
	case POSITION_MODE_ONEWAY:
		return BYBIT_POSITION_MODE_ONEWAY
	default:
		return 0
	}
}
