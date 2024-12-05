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
	//币安 标准限价 现货止盈止损限价 合约止盈止损限价 转为限价单
	case BN_ORDER_TYPE_LIMIT,
		BN_ORDER_TYPE_SPOT_STOP_LOSS_LIMIT,
		BN_ORDER_TYPE_SPOT_TAKE_PROFIT_LIMIT,
		BN_ORDER_TYPE_FUTURE_STOP,
		BN_ORDER_TYPE_FUTURE_TAKE_PROFIT:
		return ORDER_TYPE_LIMIT
	//币安 市价 合约止盈止损市价 转为市价单
	case BN_ORDER_TYPE_MARKET,
		BN_ORDER_TYPE_FUTURE_STOP_MARKET,
		BN_ORDER_TYPE_FUTURE_TAKE_PROFIT_MARKET:
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

// 触发类型
func (c *BinanceEnumConverter) FromBNOrderTypeForTriggerType(t string) OrderTriggerType {
	switch t {
	case BN_ORDER_TYPE_LIMIT,
		BN_ORDER_TYPE_MARKET:
		return ORDER_TRIGGER_TYPE_UNKNOWN
	case BN_ORDER_TYPE_SPOT_STOP_LOSS_LIMIT,
		BN_ORDER_TYPE_FUTURE_STOP,
		BN_ORDER_TYPE_FUTURE_STOP_MARKET:
		return ORDER_TRIGGER_TYPE_STOP_LOSS
	case BN_ORDER_TYPE_SPOT_TAKE_PROFIT_LIMIT,
		BN_ORDER_TYPE_FUTURE_TAKE_PROFIT,
		BN_ORDER_TYPE_FUTURE_TAKE_PROFIT_MARKET:
		return ORDER_TRIGGER_TYPE_TAKE_PROFIT
	default:
		return ORDER_TRIGGER_TYPE_UNKNOWN
	}

}
func (c *BinanceEnumConverter) ToTriggerBnOrderType(accountType BinanceAccountType, ot OrderType, tt OrderTriggerType) string {
	switch ot {
	case ORDER_TYPE_LIMIT:
		switch tt {
		case ORDER_TRIGGER_TYPE_STOP_LOSS:
			switch accountType {
			case BN_AC_SPOT:
				return BN_ORDER_TYPE_SPOT_STOP_LOSS_LIMIT
			case BN_AC_FUTURE, BN_AC_SWAP:
				return BN_ORDER_TYPE_FUTURE_STOP
			}
		case ORDER_TRIGGER_TYPE_TAKE_PROFIT:
			switch accountType {
			case BN_AC_SPOT:
				return BN_ORDER_TYPE_SPOT_TAKE_PROFIT_LIMIT
			case BN_AC_FUTURE, BN_AC_SWAP:
				return BN_ORDER_TYPE_FUTURE_TAKE_PROFIT
			}
		default:
			return BN_ORDER_TYPE_LIMIT
		}
	case ORDER_TYPE_MARKET:
		switch tt {
		case ORDER_TRIGGER_TYPE_STOP_LOSS:
			switch accountType {
			case BN_AC_SPOT:
				return BN_ORDER_TYPE_SPOT_STOP_LOSS_LIMIT
			case BN_AC_FUTURE, BN_AC_SWAP:
				return BN_ORDER_TYPE_FUTURE_STOP_MARKET
			}
		case ORDER_TRIGGER_TYPE_TAKE_PROFIT:
			switch accountType {
			case BN_AC_SPOT:
				return BN_ORDER_TYPE_SPOT_TAKE_PROFIT_LIMIT
			case BN_AC_FUTURE, BN_AC_SWAP:
				return BN_ORDER_TYPE_FUTURE_TAKE_PROFIT_MARKET
			}
		default:
			return BN_ORDER_TYPE_MARKET
		}
	}
	return ""
}

// 触发条件类型
// 订单方向为买入时 止盈为下穿 止损为上穿
// 订单方向为卖出时 止盈为上穿 止损为下穿
func (c *BinanceEnumConverter) FromBNOrderSideForTriggerConditionType(bnOrderSide, bnOrderType string) OrderTriggerConditionType {
	tt := c.FromBNOrderTypeForTriggerType(bnOrderType)
	switch bnOrderSide {
	case BN_ORDER_SIDE_BUY:
		switch tt {
		case ORDER_TRIGGER_TYPE_TAKE_PROFIT:
			//买入止盈 价格下穿触发
			return ORDER_TRIGGER_CONDITION_TYPE_THROUGH_DOWN
		case ORDER_TRIGGER_TYPE_STOP_LOSS:
			//买入止损 价格上穿触发
			return ORDER_TRIGGER_CONDITION_TYPE_THROUGH_UP
		}
	case BN_ORDER_SIDE_SELL:
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
func (c *BinanceEnumConverter) FromBNOrderStatus(t string, orderType string) OrderStatus {
	switch t {
	case BN_ORDER_STATUS_NEW:
		switch orderType {
		case BN_ORDER_TYPE_SPOT_STOP_LOSS_LIMIT,
			BN_ORDER_TYPE_SPOT_TAKE_PROFIT_LIMIT,
			BN_ORDER_TYPE_FUTURE_STOP,
			BN_ORDER_TYPE_FUTURE_TAKE_PROFIT,
			BN_ORDER_TYPE_FUTURE_STOP_MARKET,
			BN_ORDER_TYPE_FUTURE_TAKE_PROFIT_MARKET:
			return ORDER_STATUS_UN_TRIGGERED
		default:
			return ORDER_STATUS_NEW
		}
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

// 划转类型转换
func (c *BinanceEnumConverter) FromBNAssetType(t string) AssetType {
	switch t {
	case BN_ASSET_TYPE_FUND:
		return ASSET_TYPE_FUND
	case BN_ASSET_TYPE_UNIFIED:
		return ASSET_TYPE_UNIFIED
	case BN_ASSET_TYPE_UMFUTURE:
		return ASSET_TYPE_UMFUTURE
	case BN_ASSET_TYPE_CMFUTURE:
		return ASSET_TYPE_CMFUTURE
	case BN_ASSET_TYPE_PORTFOLIO_MARGIN:
		return ASSET_TYPE_PORTFOLIO_MARGIN
	case BN_ASSET_TYPE_MARGIN_ISOLATED:
		return ASSET_TYPE_MARGIN_ISOLATED
	case BN_ASSET_TYPE_MARGIN_CROSSED:
		return ASSET_TYPE_MARGIN_CROSSED
	default:
		return ""
	}
}
func (c *BinanceEnumConverter) ToBNAssetType(t AssetType) string {
	switch t {
	case ASSET_TYPE_FUND:
		return BN_ASSET_TYPE_FUND
	case ASSET_TYPE_UNIFIED:
		return BN_ASSET_TYPE_UNIFIED
	case ASSET_TYPE_UMFUTURE:
		return BN_ASSET_TYPE_UMFUTURE
	case ASSET_TYPE_CMFUTURE:
		return BN_ASSET_TYPE_CMFUTURE
	case ASSET_TYPE_PORTFOLIO_MARGIN:
		return BN_ASSET_TYPE_PORTFOLIO_MARGIN
	case ASSET_TYPE_MARGIN_CROSSED:
		return BN_ASSET_TYPE_MARGIN_CROSSED
	case ASSET_TYPE_MARGIN_ISOLATED:
		return BN_ASSET_TYPE_MARGIN_ISOLATED

	default:
		return ""
	}
}

// 划转状态类型
func (c *BinanceEnumConverter) FromBinanceTransferStatus(t string) TransferStatusType {
	switch t {
	case BN_TRANSFER_STATUS_TYPE_SUCCESS:
		return TRANSFER_STATUS_TYPE_SUCCESS
	case BN_TRANSFER_STATUS_TYPE_PENDING:
		return TRANSFER_STATUS_TYPE_PENDING
	case BN_TRANSFER_STATUS_TYPE_FAILED:
		return TRANSFER_STATUS_TYPE_FAILED
	default:
		return TRANSFER_STATUS_TYPE_UNKNOWN
	}
}
