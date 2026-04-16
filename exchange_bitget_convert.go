package mytrade

import "strings"

type BitgetEnumConverter struct{}

func (c BitgetEnumConverter) FromBitgetAccountMode(m string) AccountMode {
	switch strings.ToLower(strings.TrimSpace(m)) {
	case BITGET_AM_SINGLE:
		return ACCOUNT_MODE_SINGLE_MARGIN
	case BITGET_AM_UNION:
		return ACCOUNT_MODE_MULTI_MARGIN
	}
	return ACCOUNT_MODE_UNKNOWN
}

func (c BitgetEnumConverter) FromBitgetMarginMode(m string) MarginMode {
	switch strings.ToLower(strings.TrimSpace(m)) {
	case BITGET_MARGIN_MODE_CROSSED:
		return MARGIN_MODE_CROSSED
	case BITGET_MARGIN_MODE_ISOLATED:
		return MARGIN_MODE_ISOLATED
	}
	return MARGIN_MODE_UNKNOWN
}

func (c BitgetEnumConverter) ToBitgetMarginMode(m MarginMode) string {
	switch m {
	case MARGIN_MODE_CROSSED:
		return BITGET_MARGIN_MODE_CROSSED
	case MARGIN_MODE_ISOLATED:
		return BITGET_MARGIN_MODE_ISOLATED
	}
	return ""
}

func (c BitgetEnumConverter) FromBitgetPositionMode(m string) PositionMode {
	switch strings.ToLower(strings.TrimSpace(m)) {
	case BITGET_POSITION_MODE_ONE_WAY:
		return POSITION_MODE_ONEWAY
	case BITGET_POSITION_MODE_HEDGE:
		return POSITION_MODE_HEDGE
	}
	return POSITION_MODE_UNKNOWN
}

func (c BitgetEnumConverter) ToBitgetPositionMode(m PositionMode) string {
	switch m {
	case POSITION_MODE_ONEWAY:
		return BITGET_POSITION_MODE_ONE_WAY
	case POSITION_MODE_HEDGE:
		return BITGET_POSITION_MODE_HEDGE
	}
	return ""
}

func (c BitgetEnumConverter) FromBitgetOrderSide(s string) OrderSide {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case BITGET_ORDER_SIDE_BUY:
		return ORDER_SIDE_BUY
	case BITGET_ORDER_SIDE_SELL:
		return ORDER_SIDE_SELL
	default:
		return ORDER_SIDE_UNKNOWN
	}
}

func (c BitgetEnumConverter) ToBitgetOrderSide(s OrderSide) string {
	switch s {
	case ORDER_SIDE_BUY:
		return BITGET_ORDER_SIDE_BUY
	case ORDER_SIDE_SELL:
		return BITGET_ORDER_SIDE_SELL
	default:
		return ""
	}
}

func (c BitgetEnumConverter) FromBitgetOrderType(t string) (OrderType, TimeInForce) {
	switch strings.ToLower(strings.TrimSpace(t)) {
	case BITGET_ORDER_TYPE_LIMIT:
		return ORDER_TYPE_LIMIT, TIME_IN_FORCE_GTC
	case BITGET_ORDER_TYPE_MARKET:
		return ORDER_TYPE_MARKET, TIME_IN_FORCE_GTC
	default:
		return ORDER_TYPE_UNKNOWN, TIME_IN_FORCE_UNKNOWN
	}
}

func (c BitgetEnumConverter) FromBitgetOrderTypeWithTIF(orderType, tif string) (OrderType, TimeInForce) {
	ot, _ := c.FromBitgetOrderType(orderType)
	switch strings.ToLower(strings.TrimSpace(tif)) {
	case BITGET_TIF_GTC:
		return ot, TIME_IN_FORCE_GTC
	case BITGET_TIF_IOC:
		return ORDER_TYPE_LIMIT, TIME_IN_FORCE_IOC
	case BITGET_TIF_FOK:
		return ORDER_TYPE_LIMIT, TIME_IN_FORCE_FOK
	case BITGET_TIF_POST_ONLY:
		return ORDER_TYPE_LIMIT, TIME_IN_FORCE_POST_ONLY
	default:
		if ot == ORDER_TYPE_MARKET {
			return ORDER_TYPE_MARKET, TIME_IN_FORCE_GTC
		}
		return ot, TIME_IN_FORCE_GTC
	}
}

func (c BitgetEnumConverter) ToBitgetOrderType(t OrderType) string {
	switch t {
	case ORDER_TYPE_LIMIT:
		return BITGET_ORDER_TYPE_LIMIT
	case ORDER_TYPE_MARKET:
		return BITGET_ORDER_TYPE_MARKET
	default:
		return ""
	}
}

func (c BitgetEnumConverter) ToBitgetTimeInForce(tif TimeInForce) string {
	if tif == TIME_IN_FORCE_UNKNOWN || tif == "" {
		return BITGET_TIF_GTC
	}
	switch tif {
	case TIME_IN_FORCE_GTC:
		return BITGET_TIF_GTC
	case TIME_IN_FORCE_IOC:
		return BITGET_TIF_IOC
	case TIME_IN_FORCE_FOK:
		return BITGET_TIF_FOK
	case TIME_IN_FORCE_POST_ONLY:
		return BITGET_TIF_POST_ONLY
	default:
		return BITGET_TIF_GTC
	}
}

func (c BitgetEnumConverter) ToBitgetPositionSide(p PositionSide) string {
	switch p {
	case POSITION_SIDE_LONG:
		return BITGET_POS_SIDE_LONG
	case POSITION_SIDE_SHORT:
		return BITGET_POS_SIDE_SHORT
	default:
		return BITGET_POS_SIDE_NET
	}
}

func (c BitgetEnumConverter) FromBitgetPositionSide(s string) PositionSide {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case BITGET_POS_SIDE_LONG:
		return POSITION_SIDE_LONG
	case BITGET_POS_SIDE_SHORT:
		return POSITION_SIDE_SHORT
	case BITGET_POS_SIDE_NET, "":
		return POSITION_SIDE_BOTH
	default:
		return POSITION_SIDE_UNKNOWN
	}
}

func (c BitgetEnumConverter) ReduceOnlyFromString(s string) bool {
	return strings.EqualFold(strings.TrimSpace(s), BITGET_REDUCE_ONLY_YES)
}

func (c BitgetEnumConverter) FromBitgetOrderStatusUTA(s string) OrderStatus {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case BITGET_ORDER_STATUS_LIVE, BITGET_ORDER_STATUS_NEW:
		return ORDER_STATUS_NEW
	case BITGET_ORDER_STATUS_PARTIALLY_FILLED:
		return ORDER_STATUS_PARTIALLY_FILLED
	case BITGET_ORDER_STATUS_FILLED:
		return ORDER_STATUS_FILLED
	case BITGET_ORDER_STATUS_CANCELLED, BITGET_ORDER_STATUS_CANCELED:
		return ORDER_STATUS_CANCELED
	case BITGET_ORDER_STATUS_REJECTED:
		return ORDER_STATUS_REJECTED
	default:
		return ORDER_STATUS_UNKNOWN
	}
}

func (c BitgetEnumConverter) FromBitgetOrderStatusClassicSpot(s string) OrderStatus {
	switch strings.ToLower(s) {
	case BITGET_CLASSIC_ORDER_STATUS_NEW:
		return ORDER_STATUS_NEW
	case BITGET_CLASSIC_ORDER_STATUS_PARTIAL_FILL, "partial-fill", BITGET_ORDER_STATUS_PARTIALLY_FILLED:
		return ORDER_STATUS_PARTIALLY_FILLED
	case BITGET_CLASSIC_ORDER_STATUS_FULL_FILL, "full-fill", BITGET_ORDER_STATUS_FILLED:
		return ORDER_STATUS_FILLED
	case BITGET_ORDER_STATUS_CANCELLED, BITGET_ORDER_STATUS_CANCELED:
		return ORDER_STATUS_CANCELED
	default:
		return ORDER_STATUS_UNKNOWN
	}
}

func (c BitgetEnumConverter) FromBitgetOrderStatusClassicFutures(s string) OrderStatus {
	switch strings.ToLower(s) {
	case BITGET_CLASSIC_ORDER_STATUS_NEW:
		return ORDER_STATUS_NEW
	case BITGET_CLASSIC_ORDER_STATUS_PARTIAL_FILL:
		return ORDER_STATUS_PARTIALLY_FILLED
	case BITGET_ORDER_STATUS_FILLED, BITGET_CLASSIC_ORDER_STATUS_FULL_FILL:
		return ORDER_STATUS_FILLED
	case BITGET_ORDER_STATUS_CANCELED, BITGET_ORDER_STATUS_CANCELLED:
		return ORDER_STATUS_CANCELED
	default:
		return ORDER_STATUS_UNKNOWN
	}
}

func (c BitgetEnumConverter) FromClassicForce(f string) TimeInForce {
	switch strings.ToLower(strings.TrimSpace(f)) {
	case BITGET_TIF_GTC, "normal", "":
		return TIME_IN_FORCE_GTC
	case BITGET_TIF_IOC:
		return TIME_IN_FORCE_IOC
	case BITGET_TIF_FOK:
		return TIME_IN_FORCE_FOK
	case BITGET_TIF_POST_ONLY:
		return TIME_IN_FORCE_POST_ONLY
	default:
		return TIME_IN_FORCE_GTC
	}
}
