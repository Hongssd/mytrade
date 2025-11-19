package mytrade

import "strconv"

// 枚举转换器
type AsterEnumConverter struct{}

// 订单方向
func (c *AsterEnumConverter) FromAsterOrderSide(t string) OrderSide {
	switch t {
	case ASTER_ORDER_SIDE_BUY:
		return ORDER_SIDE_BUY
	case ASTER_ORDER_SIDE_SELL:
		return ORDER_SIDE_SELL
	default:
		return ORDER_SIDE_UNKNOWN
	}
}
func (c *AsterEnumConverter) ToAsterOrderSide(t OrderSide) string {
	switch t {
	case ORDER_SIDE_BUY:
		return ASTER_ORDER_SIDE_BUY
	case ORDER_SIDE_SELL:
		return ASTER_ORDER_SIDE_SELL
	default:
		return ""
	}
}

// 订单类型
func (c *AsterEnumConverter) FromAsterOrderType(t string) OrderType {
	switch t {
	//币安 标准限价 限价做市 现货止盈止损限价 合约止盈止损限价 转为限价单
	case ASTER_ORDER_TYPE_LIMIT,
		ASTER_ORDER_TYPE_LIMIT_MAKER,
		ASTER_ORDER_TYPE_SPOT_STOP_LOSS_LIMIT,
		ASTER_ORDER_TYPE_SPOT_TAKE_PROFIT_LIMIT,
		ASTER_ORDER_TYPE_FUTURE_STOP,
		ASTER_ORDER_TYPE_FUTURE_TAKE_PROFIT:
		return ORDER_TYPE_LIMIT
	//币安 市价 合约止盈止损市价 转为市价单
	case ASTER_ORDER_TYPE_MARKET,
		ASTER_ORDER_TYPE_FUTURE_STOP_MARKET,
		ASTER_ORDER_TYPE_FUTURE_TAKE_PROFIT_MARKET:
		return ORDER_TYPE_MARKET
	default:
		return OrderType(t)
	}
}
func (c *AsterEnumConverter) ToAsterOrderType(t OrderType) string {
	switch t {
	case ORDER_TYPE_LIMIT:
		return ASTER_ORDER_TYPE_LIMIT
	case ORDER_TYPE_MARKET:
		return ASTER_ORDER_TYPE_MARKET
	default:
		return ""
	}
}

// 触发类型
func (c *AsterEnumConverter) FromAsterOrderTypeForTriggerType(t string) OrderTriggerType {
	switch t {
	case ASTER_ORDER_TYPE_LIMIT,
		ASTER_ORDER_TYPE_MARKET:
		return ORDER_TRIGGER_TYPE_UNKNOWN
	case ASTER_ORDER_TYPE_SPOT_STOP_LOSS_LIMIT,
		ASTER_ORDER_TYPE_FUTURE_STOP,
		ASTER_ORDER_TYPE_FUTURE_STOP_MARKET:
		return ORDER_TRIGGER_TYPE_STOP_LOSS
	case ASTER_ORDER_TYPE_SPOT_TAKE_PROFIT_LIMIT,
		ASTER_ORDER_TYPE_FUTURE_TAKE_PROFIT,
		ASTER_ORDER_TYPE_FUTURE_TAKE_PROFIT_MARKET:
		return ORDER_TRIGGER_TYPE_TAKE_PROFIT
	default:
		return ORDER_TRIGGER_TYPE_UNKNOWN
	}

}
func (c *AsterEnumConverter) ToTriggerBnOrderType(accountType AsterAccountType, ot OrderType, tt OrderTriggerType) string {
	switch ot {
	case ORDER_TYPE_LIMIT:
		switch tt {
		case ORDER_TRIGGER_TYPE_STOP_LOSS:
			switch accountType {
			case ASTER_AC_SPOT:
				return ASTER_ORDER_TYPE_SPOT_STOP_LOSS_LIMIT
			case ASTER_AC_FUTURE:
				return ASTER_ORDER_TYPE_FUTURE_STOP
			}
		case ORDER_TRIGGER_TYPE_TAKE_PROFIT:
			switch accountType {
			case ASTER_AC_SPOT:
				return ASTER_ORDER_TYPE_SPOT_TAKE_PROFIT_LIMIT
			case ASTER_AC_FUTURE:
				return ASTER_ORDER_TYPE_FUTURE_TAKE_PROFIT
			}
		default:
			return ASTER_ORDER_TYPE_LIMIT
		}
	case ORDER_TYPE_MARKET:
		switch tt {
		case ORDER_TRIGGER_TYPE_STOP_LOSS:
			switch accountType {
			case ASTER_AC_SPOT:
				return ASTER_ORDER_TYPE_SPOT_STOP_LOSS_LIMIT
			case ASTER_AC_FUTURE:
				return ASTER_ORDER_TYPE_FUTURE_STOP_MARKET
			}
		case ORDER_TRIGGER_TYPE_TAKE_PROFIT:
			switch accountType {
			case ASTER_AC_SPOT:
				return ASTER_ORDER_TYPE_SPOT_TAKE_PROFIT_LIMIT
			case ASTER_AC_FUTURE:
				return ASTER_ORDER_TYPE_FUTURE_TAKE_PROFIT_MARKET
			}
		default:
			return ASTER_ORDER_TYPE_MARKET
		}
	}
	return ""
}

// 触发条件类型
// 订单方向为买入时 止盈为下穿 止损为上穿
// 订单方向为卖出时 止盈为上穿 止损为下穿
func (c *AsterEnumConverter) FromAsterOrderSideForTriggerConditionType(bnOrderSide, bnOrderType string) OrderTriggerConditionType {
	tt := c.FromAsterOrderTypeForTriggerType(bnOrderType)
	switch bnOrderSide {
	case ASTER_ORDER_SIDE_BUY:
		switch tt {
		case ORDER_TRIGGER_TYPE_TAKE_PROFIT:
			//买入止盈 价格下穿触发
			return ORDER_TRIGGER_CONDITION_TYPE_THROUGH_DOWN
		case ORDER_TRIGGER_TYPE_STOP_LOSS:
			//买入止损 价格上穿触发
			return ORDER_TRIGGER_CONDITION_TYPE_THROUGH_UP
		}
	case ASTER_ORDER_SIDE_SELL:
		switch tt {
		case ORDER_TRIGGER_TYPE_TAKE_PROFIT:
			//卖出止盈 价格上穿触发
			return ORDER_TRIGGER_CONDITION_TYPE_THROUGH_UP
		case ORDER_TRIGGER_TYPE_STOP_LOSS:
			//卖出止损 价格下穿触发
			return ORDER_TRIGGER_CONDITION_TYPE_THROUGH_DOWN
		}
	}
	return ORDER_TRIGGER_CONDITION_TYPE_UNKNOWN
}

// 仓位方向
func (c *AsterEnumConverter) FromAsterPositionSide(t string) PositionSide {
	switch t {
	case ASTER_POSITION_SIDE_LONG:
		return POSITION_SIDE_LONG
	case ASTER_POSITION_SIDE_SHORT:
		return POSITION_SIDE_SHORT
	case ASTER_POSITION_SIDE_BOTH:
		return POSITION_SIDE_BOTH
	default:
		return POSITION_SIDE_UNKNOWN
	}
}
func (c *AsterEnumConverter) ToAsterPositionSide(t PositionSide) string {
	switch t {
	case POSITION_SIDE_LONG:
		return ASTER_POSITION_SIDE_LONG
	case POSITION_SIDE_SHORT:
		return ASTER_POSITION_SIDE_SHORT
	case POSITION_SIDE_BOTH:
		return ASTER_POSITION_SIDE_BOTH
	default:
		return ""
	}
}

// 有效方式
func (c *AsterEnumConverter) FromAsterTimeInForce(t string) TimeInForce {
	switch t {
	case ASTER_TIME_IN_FORCE_GTC:
		return TIME_IN_FORCE_GTC
	case ASTER_TIME_IN_FORCE_IOC:
		return TIME_IN_FORCE_IOC
	case ASTER_TIME_IN_FORCE_FOK:
		return TIME_IN_FORCE_FOK
	case ASTER_TIME_IN_FORCE_POST_ONLY:
		return TIME_IN_FORCE_POST_ONLY
	default:
		return TIME_IN_FORCE_UNKNOWN
	}
}
func (c *AsterEnumConverter) ToAsterTimeInForce(t TimeInForce) string {
	switch t {
	case TIME_IN_FORCE_GTC:
		return ASTER_TIME_IN_FORCE_GTC
	case TIME_IN_FORCE_IOC:
		return ASTER_TIME_IN_FORCE_IOC
	case TIME_IN_FORCE_FOK:
		return ASTER_TIME_IN_FORCE_FOK
	case TIME_IN_FORCE_POST_ONLY:
		return ASTER_TIME_IN_FORCE_POST_ONLY
	default:
		return ""
	}
}

// 订单状态
func (c *AsterEnumConverter) FromAsterOrderStatus(t string, orderType string) OrderStatus {
	switch t {
	case ASTER_ORDER_STATUS_NEW:
		switch orderType {
		case ASTER_ORDER_TYPE_SPOT_STOP_LOSS_LIMIT,
			ASTER_ORDER_TYPE_SPOT_TAKE_PROFIT_LIMIT,
			ASTER_ORDER_TYPE_FUTURE_STOP,
			ASTER_ORDER_TYPE_FUTURE_TAKE_PROFIT,
			ASTER_ORDER_TYPE_FUTURE_STOP_MARKET,
			ASTER_ORDER_TYPE_FUTURE_TAKE_PROFIT_MARKET:
			return ORDER_STATUS_UN_TRIGGERED
		default:
			return ORDER_STATUS_NEW
		}
	case ASTER_ORDER_STATUS_PARTIALLY_FILLED:
		return ORDER_STATUS_PARTIALLY_FILLED
	case ASTER_ORDER_STATUS_FILLED:
		return ORDER_STATUS_FILLED
	case ASTER_ORDER_STATUS_CANCELED, ASTER_ORDER_STATUS_EXPIRED:
		return ORDER_STATUS_CANCELED
	case ASTER_ORDER_STATUS_REJECTED:
		return ORDER_STATUS_REJECTED
	default:
		return ORDER_STATUS_UNKNOWN
	}
}
func (c *AsterEnumConverter) ToAsterOrderStatus(t OrderStatus) string {
	switch t {
	case ORDER_STATUS_NEW:
		return ASTER_ORDER_STATUS_NEW
	case ORDER_STATUS_PARTIALLY_FILLED:
		return ASTER_ORDER_STATUS_PARTIALLY_FILLED
	case ORDER_STATUS_FILLED:
		return ASTER_ORDER_STATUS_FILLED
	case ORDER_STATUS_CANCELED:
		return ASTER_ORDER_STATUS_CANCELED
	case ORDER_STATUS_REJECTED:
		return ASTER_ORDER_STATUS_REJECTED
	default:
		return ""
	}
}

// 账户模式
func (c *AsterEnumConverter) FromAsterAccountMode(t bool) AccountMode {
	switch t {
	case ASTER_ACCOUNT_MOED_MULTI_CURRENCY_MARGIN:
		return ACCOUNT_MODE_MULTI_MARGIN
	case ASTER_ACCOUNT_MODE_SINGLE_CURRENCY_MARGIN:
		return ACCOUNT_MODE_SINGLE_MARGIN
	default:
		return ACCOUNT_MODE_UNKNOWN
	}
}
func (c *AsterEnumConverter) ToAsterAccountMode(t AccountMode) string {
	switch t {
	case ACCOUNT_MODE_MULTI_MARGIN, ACCOUNT_MODE_PORTFOLIO:
		return strconv.FormatBool(ASTER_ACCOUNT_MOED_MULTI_CURRENCY_MARGIN)
	case ACCOUNT_MODE_FREE_MARGIN, ACCOUNT_MODE_SINGLE_MARGIN:
		return strconv.FormatBool(ASTER_ACCOUNT_MODE_SINGLE_CURRENCY_MARGIN)
	default:
		return ""
	}
}

// 保证金模式
func (c *AsterEnumConverter) FromAsterMarginMode(t bool) MarginMode {
	switch t {
	case ASTER_MARGIN_MODE_ISOLATED:
		return MARGIN_MODE_ISOLATED
	case ASTER_MARGIN_MODE_CROSSED:
		return MARGIN_MODE_CROSSED
	default:
		return MARGIN_MODE_UNKNOWN
	}
}
func (c *AsterEnumConverter) ToAsterMarginMode(t MarginMode) bool {
	switch t {
	case MARGIN_MODE_ISOLATED:
		return ASTER_MARGIN_MODE_ISOLATED
	case MARGIN_MODE_CROSSED:
		return ASTER_MARGIN_MODE_CROSSED
	default:
		return false
	}
}
func (c *AsterEnumConverter) ToAsterMarginModeStr(t MarginMode) string {
	switch t {
	case MARGIN_MODE_ISOLATED:
		return ASTER_MARGIN_MODE_ISOLATED_STR
	case MARGIN_MODE_CROSSED:
		return ASTER_MARGIN_MODE_CROSSED_STR
	default:
		return ""
	}
}

// 仓位模式
func (c *AsterEnumConverter) FromAsterPositionMode(t bool) PositionMode {
	switch t {
	case ASTER_POSITION_MODE_HEDGE:
		return POSITION_MODE_HEDGE
	case ASTER_POSITION_MODE_ONEWAY:
		return POSITION_MODE_ONEWAY
	default:
		return POSITION_MODE_UNKNOWN
	}
}
func (c *AsterEnumConverter) ToAsterPositionMode(t PositionMode) string {
	switch t {
	case POSITION_MODE_HEDGE:
		return strconv.FormatBool(ASTER_POSITION_MODE_HEDGE)
	case POSITION_MODE_ONEWAY:
		return strconv.FormatBool(ASTER_POSITION_MODE_ONEWAY)
	default:
		return ""
	}
}

// 划转类型转换
func (c *AsterEnumConverter) FromAsterAssetType(t string) AssetType {
	switch t {
	case ASTER_ASSET_TYPE_UNIFIED:
		return ASSET_TYPE_UNIFIED
	case ASTER_ASSET_TYPE_UMFUTURE:
		return ASSET_TYPE_UMFUTURE
	default:
		return ""
	}
}
func (c *AsterEnumConverter) ToAsterAssetType(t AssetType) string {
	switch t {
	case ASSET_TYPE_UNIFIED:
		return ASTER_ASSET_TYPE_UNIFIED
	case ASSET_TYPE_UMFUTURE:
		return ASTER_ASSET_TYPE_UMFUTURE
	default:
		return ""
	}
}

// 划转状态类型
func (c *AsterEnumConverter) FromAsterTransferStatus(t string) TransferStatusType {
	switch t {
	case ASTER_TRANSFER_STATUS_TYPE_SUCCESS:
		return TRANSFER_STATUS_TYPE_SUCCESS
	case ASTER_TRANSFER_STATUS_TYPE_PENDING:
		return TRANSFER_STATUS_TYPE_PENDING
	case ASTER_TRANSFER_STATUS_TYPE_FAILED:
		return TRANSFER_STATUS_TYPE_FAILED
	default:
		return TRANSFER_STATUS_TYPE_UNKNOWN
	}
}
