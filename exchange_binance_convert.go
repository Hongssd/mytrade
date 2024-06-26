package mytrade

import "strconv"

// 枚举转换器
type BinanceEnumConverter struct{}

// 订单方向
func (c *BinanceEnumConverter) FromBNOrderSide(t string) OrderSide {
	switch t {
	case BN_ORDER_SIDE_BUY:
		return ORDER_SIDE_BUY
	case BN_ORDER_SIDE_SELL:
		return ORDER_SIDE_SELL
	default:
		return ORDER_SIDE_UNKNOWN
	}
}
func (c *BinanceEnumConverter) ToBNOrderSide(t OrderSide) string {
	switch t {
	case ORDER_SIDE_BUY:
		return BN_ORDER_SIDE_BUY
	case ORDER_SIDE_SELL:
		return BN_ORDER_SIDE_SELL
	default:
		return ""
	}
}

// 订单类型
func (c *BinanceEnumConverter) FromBNOrderType(t string) OrderType {
	switch t {
	case BN_ORDER_TYPE_LIMIT:
		return ORDER_TYPE_LIMIT
	case BN_ORDER_TYPE_MARKET:
		return ORDER_TYPE_MARKET
	default:
		return ORDER_TYPE_UNKNOWN
	}
}
func (c *BinanceEnumConverter) ToBNOrderType(t OrderType) string {
	switch t {
	case ORDER_TYPE_LIMIT:
		return BN_ORDER_TYPE_LIMIT
	case ORDER_TYPE_MARKET:
		return BN_ORDER_TYPE_MARKET
	default:
		return ""
	}
}

// 仓位方向
func (c *BinanceEnumConverter) FromBNPositionSide(t string) PositionSide {
	switch t {
	case BN_POSITION_SIDE_LONG:
		return POSITION_SIDE_LONG
	case BN_POSITION_SIDE_SHORT:
		return POSITION_SIDE_SHORT
	case BN_POSITION_SIDE_BOTH:
		return POSITION_SIDE_BOTH
	default:
		return POSITION_SIDE_UNKNOWN
	}
}
func (c *BinanceEnumConverter) ToBNPositionSide(t PositionSide) string {
	switch t {
	case POSITION_SIDE_LONG:
		return BN_POSITION_SIDE_LONG
	case POSITION_SIDE_SHORT:
		return BN_POSITION_SIDE_SHORT
	case POSITION_SIDE_BOTH:
		return BN_POSITION_SIDE_BOTH
	default:
		return ""
	}
}

// 有效方式
func (c *BinanceEnumConverter) FromBNTimeInForce(t string) TimeInForce {
	switch t {
	case BN_TIME_IN_FORCE_GTC:
		return TIME_IN_FORCE_GTC
	case BN_TIME_IN_FORCE_IOC:
		return TIME_IN_FORCE_IOC
	case BN_TIME_IN_FORCE_FOK:
		return TIME_IN_FORCE_FOK
	case BN_TIME_IN_FORCE_POST_ONLY:
		return TIME_IN_FORCE_POST_ONLY
	default:
		return TIME_IN_FORCE_UNKNOWN
	}
}
func (c *BinanceEnumConverter) ToBNTimeInForce(t TimeInForce) string {
	switch t {
	case TIME_IN_FORCE_GTC:
		return BN_TIME_IN_FORCE_GTC
	case TIME_IN_FORCE_IOC:
		return BN_TIME_IN_FORCE_IOC
	case TIME_IN_FORCE_FOK:
		return BN_TIME_IN_FORCE_FOK
	case TIME_IN_FORCE_POST_ONLY:
		return BN_TIME_IN_FORCE_POST_ONLY
	default:
		return ""
	}
}

// 订单状态
func (c *BinanceEnumConverter) FromBNOrderStatus(t string) OrderStatus {
	switch t {
	case BN_ORDER_STATUS_NEW:
		return ORDER_STATUS_NEW
	case BN_ORDER_STATUS_PARTIALLY_FILLED:
		return ORDER_STATUS_PARTIALLY_FILLED
	case BN_ORDER_STATUS_FILLED:
		return ORDER_STATUS_FILLED
	case BN_ORDER_STATUS_CANCELED, BN_ORDER_STATUS_EXPIRED:
		return ORDER_STATUS_CANCELED
	case BN_ORDER_STATUS_REJECTED:
		return ORDER_STATUS_REJECTED
	default:
		return ORDER_STATUS_UNKNOWN
	}
}
func (c *BinanceEnumConverter) ToBNOrderStatus(t OrderStatus) string {
	switch t {
	case ORDER_STATUS_NEW:
		return BN_ORDER_STATUS_NEW
	case ORDER_STATUS_PARTIALLY_FILLED:
		return BN_ORDER_STATUS_PARTIALLY_FILLED
	case ORDER_STATUS_FILLED:
		return BN_ORDER_STATUS_FILLED
	case ORDER_STATUS_CANCELED:
		return BN_ORDER_STATUS_CANCELED
	case ORDER_STATUS_REJECTED:
		return BN_ORDER_STATUS_REJECTED
	default:
		return ""
	}
}

// 账户模式
func (c *BinanceEnumConverter) FromBNAccountMode(t bool) AccountMode {
	switch t {
	case BN_ACCOUNT_MOED_MULTI_CURRENCY_MARGIN:
		return ACCOUNT_MODE_MULTI_MARGIN
	case BN_ACCOUNT_MODE_SINGLE_CURRENCY_MARGIN:
		return ACCOUNT_MODE_SINGLE_MARGIN
	default:
		return ACCOUNT_MODE_UNKNOWN
	}
}
func (c *BinanceEnumConverter) ToBNAccountMode(t AccountMode) string {
	switch t {
	case ACCOUNT_MODE_MULTI_MARGIN, ACCOUNT_MODE_PORTFOLIO:
		return strconv.FormatBool(BN_ACCOUNT_MOED_MULTI_CURRENCY_MARGIN)
	case ACCOUNT_MODE_FREE_MARGIN, ACCOUNT_MODE_SINGLE_MARGIN:
		return strconv.FormatBool(BN_ACCOUNT_MODE_SINGLE_CURRENCY_MARGIN)
	default:
		return ""
	}
}

// 保证金模式
func (c *BinanceEnumConverter) FromBNMarginMode(t bool) MarginMode {
	switch t {
	case BN_MARGIN_MODE_ISOLATED:
		return MARGIN_MODE_ISOLATED
	case BN_MARGIN_MODE_CROSSED:
		return MARGIN_MODE_CROSSED
	default:
		return MARGIN_MODE_UNKNOWN
	}
}
func (c *BinanceEnumConverter) ToBNMarginMode(t MarginMode) bool {
	switch t {
	case MARGIN_MODE_ISOLATED:
		return BN_MARGIN_MODE_ISOLATED
	case MARGIN_MODE_CROSSED:
		return BN_MARGIN_MODE_CROSSED
	default:
		return false
	}
}
func (c *BinanceEnumConverter) ToBNMarginModeStr(t MarginMode) string {
	switch t {
	case MARGIN_MODE_ISOLATED:
		return BN_MARGIN_MODE_ISOLATED_STR
	case MARGIN_MODE_CROSSED:
		return BN_MARGIN_MODE_CROSSED_STR
	default:
		return ""
	}
}

// 仓位模式
func (c *BinanceEnumConverter) FromBNPositionMode(t bool) PositionMode {
	switch t {
	case BN_POSITION_MODE_HEDGE:
		return POSITION_MODE_HEDGE
	case BN_POSITION_MODE_ONEWAY:
		return POSITION_MODE_ONEWAY
	default:
		return POSITION_MODE_UNKNOWN
	}
}
func (c *BinanceEnumConverter) ToBNPositionMode(t PositionMode) string {
	switch t {
	case POSITION_MODE_HEDGE:
		return strconv.FormatBool(BN_POSITION_MODE_HEDGE)
	case POSITION_MODE_ONEWAY:
		return strconv.FormatBool(BN_POSITION_MODE_ONEWAY)
	default:
		return ""
	}
}
