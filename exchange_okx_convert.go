package mytrade

// 枚举转换器
type OkxEnumConverter struct{}

// 订单方向
func (c *OkxEnumConverter) FromOKXOrderSide(t string) OrderSide {
	switch t {
	case OKX_ORDER_SIDE_BUY:
		return ORDER_SIDE_BUY
	case OKX_ORDER_SIDE_SELL:
		return ORDER_SIDE_SELL
	default:
		return ORDER_SIDE_UNKNOWN
	}
}
func (c *OkxEnumConverter) ToOKXOrderSide(t OrderSide) string {
	switch t {
	case ORDER_SIDE_BUY:
		return OKX_ORDER_SIDE_BUY
	case ORDER_SIDE_SELL:
		return OKX_ORDER_SIDE_SELL
	default:
		return ""
	}
}

// 订单类型
func (c *OkxEnumConverter) FromOKXOrderType(t string) (OrderType, TimeInForce) {
	switch t {
	case OKX_ORDER_TYPE_LIMIT:
		return ORDER_TYPE_LIMIT, TIME_IN_FORCE_GTC
	case OKX_ORDER_TYPE_MARKET:
		return ORDER_TYPE_MARKET, TIME_IN_FORCE_GTC
	case OKX_ORDER_TYPE_FOK:
		return ORDER_TYPE_LIMIT, TIME_IN_FORCE_FOK
	case OKX_ORDER_TYPE_IOC:
		return ORDER_TYPE_LIMIT, TIME_IN_FORCE_IOC
	case OKX_ORDER_TYPE_POST_ONLY:
		return ORDER_TYPE_LIMIT, TIME_IN_FORCE_POST_ONLY
	default:
		return ORDER_TYPE_UNKNOWN, TIME_IN_FORCE_UNKNOWN
	}
}
func (c *OkxEnumConverter) ToOKXOrderType(t OrderType, t2 TimeInForce) string {
	if t2 == TIME_IN_FORCE_GTC {
		switch t {
		case ORDER_TYPE_LIMIT:
			return OKX_ORDER_TYPE_LIMIT
		case ORDER_TYPE_MARKET:
			return OKX_ORDER_TYPE_MARKET
		default:
			return ""
		}
	} else {
		switch t2 {
		case TIME_IN_FORCE_FOK:
			return OKX_ORDER_TYPE_FOK
		case TIME_IN_FORCE_IOC:
			return OKX_ORDER_TYPE_IOC
		case TIME_IN_FORCE_POST_ONLY:
			return OKX_ORDER_TYPE_POST_ONLY
		default:
			return ""
		}
	}
}

// 策略订单类型

// 仓位方向
func (c *OkxEnumConverter) FromOKXPositionSide(t string) PositionSide {
	switch t {
	case OKX_POSITION_SIDE_LONG:
		return POSITION_SIDE_LONG
	case OKX_POSITION_SIDE_SHORT:
		return POSITION_SIDE_SHORT
	case OKX_POSITION_SIDE_BOTH:
		return POSITION_SIDE_BOTH
	default:
		return POSITION_SIDE_UNKNOWN
	}
}
func (c *OkxEnumConverter) ToOKXPositionSide(t PositionSide) string {
	switch t {
	case POSITION_SIDE_LONG:
		return OKX_POSITION_SIDE_LONG
	case POSITION_SIDE_SHORT:
		return OKX_POSITION_SIDE_SHORT
	case POSITION_SIDE_BOTH:
		return OKX_POSITION_SIDE_BOTH
	default:
		return ""
	}
}

// 订单状态
func (c *OkxEnumConverter) FromOKXOrderStatus(t string) OrderStatus {
	switch t {
	case OKX_ORDER_STATUS_NEW:
		return ORDER_STATUS_NEW
	case OKX_ORDER_STATUS_PARTIALLY_FILLED:
		return ORDER_STATUS_PARTIALLY_FILLED
	case OKX_ORDER_STATUS_FILLED:
		return ORDER_STATUS_FILLED
	case OKX_ORDER_STATUS_CANCELED:
		return ORDER_STATUS_CANCELED
	case OKX_ORDER_STATUS_REJECTED:
		return ORDER_STATUS_REJECTED
	default:
		return ORDER_STATUS_UNKNOWN
	}
}
func (c *OkxEnumConverter) ToOKXOrderStatus(t OrderStatus) string {
	switch t {
	case ORDER_STATUS_NEW:
		return OKX_ORDER_STATUS_NEW
	case ORDER_STATUS_PARTIALLY_FILLED:
		return OKX_ORDER_STATUS_PARTIALLY_FILLED
	case ORDER_STATUS_FILLED:
		return OKX_ORDER_STATUS_FILLED
	case ORDER_STATUS_CANCELED:
		return OKX_ORDER_STATUS_CANCELED
	case ORDER_STATUS_REJECTED:
		return OKX_ORDER_STATUS_REJECTED
	default:
		return ""
	}
}

func (c *OkxEnumConverter) getTdModeFromAccountType(accountType OkxAccountType, isIsolated bool) string {
	tdMode := ""
	switch accountType {
	case OKX_AC_SPOT, OKX_AC_MARGIN:
		if !isIsolated {
			tdMode = "cash"
		} else {
			tdMode = "spot_isolated"
		}
	case OKX_AC_SWAP, OKX_AC_FUTURES:
		if !isIsolated {
			tdMode = "cross"
		} else {
			tdMode = "isolated"
		}
	default:
		return tdMode
	}
	return tdMode
}

// 账户模式
func (c *OkxEnumConverter) FromOKXAccountMode(t string) AccountMode {
	switch t {
	case OKX_ACCOUNT_MODE_FREE_MARGIN:
		return ACCOUNT_MODE_FREE_MARGIN
	case OKX_ACCOUNT_MODE_SINGLE_CURRENCY_MARGIN:
		return ACCOUNT_MODE_SINGLE_MARGIN
	case OKX_ACCOUNT_MODE_MULTI_CURRENCY_MARGIN:
		return ACCOUNT_MODE_MULTI_MARGIN
	case OKX_ACCOUNT_MODE_PORTFOLIO_MARGIN:
		return ACCOUNT_MODE_PORTFOLIO
	default:
		return ACCOUNT_MODE_UNKNOWN
	}
}
func (c *OkxEnumConverter) ToOKXAccountMode(t AccountMode) string {
	switch t {
	case ACCOUNT_MODE_FREE_MARGIN:
		return OKX_ACCOUNT_MODE_FREE_MARGIN
	case ACCOUNT_MODE_SINGLE_MARGIN:
		return OKX_ACCOUNT_MODE_SINGLE_CURRENCY_MARGIN
	case ACCOUNT_MODE_MULTI_MARGIN:
		return OKX_ACCOUNT_MODE_MULTI_CURRENCY_MARGIN
	case ACCOUNT_MODE_PORTFOLIO:
		return OKX_ACCOUNT_MODE_PORTFOLIO_MARGIN
	default:
		return ""
	}
}

// 保证金模式
func (c *OkxEnumConverter) FromOKXMarginMode(t string) MarginMode {
	switch t {
	case OKX_MARGIN_MODE_ISOLATED:
		return MARGIN_MODE_ISOLATED
	case OKX_MARGIN_MODE_CROSSED:
		return MARGIN_MODE_CROSSED
	default:
		return MARGIN_MODE_UNKNOWN
	}
}
func (c *OkxEnumConverter) ToOKXMarginMode(t MarginMode) string {
	switch t {
	case MARGIN_MODE_ISOLATED:
		return OKX_MARGIN_MODE_ISOLATED
	case MARGIN_MODE_CROSSED:
		return OKX_MARGIN_MODE_CROSSED
	default:
		return ""
	}
}

// 仓位模式
func (c *OkxEnumConverter) FromOKXPositionMode(t string) PositionMode {
	switch t {
	case OKX_POSITION_MODE_HEDGE:
		return POSITION_MODE_HEDGE
	case OKX_POSITION_MODE_ONEWAY:
		return POSITION_MODE_ONEWAY
	default:
		return POSITION_MODE_UNKNOWN
	}
}
func (c *OkxEnumConverter) ToOKXPositionMode(t PositionMode) string {
	switch t {
	case POSITION_MODE_HEDGE:
		return OKX_POSITION_MODE_HEDGE
	case POSITION_MODE_ONEWAY:
		return OKX_POSITION_MODE_ONEWAY
	default:
		return ""
	}
}
