package mytrade

import "strings"

type XcoinEnumConverter struct{}

// 资金划转类型
func (c *XcoinEnumConverter) FromXcoinAssetType(t string) AssetType {
	switch t {
	case XCOIN_ASSET_TYPE_FUNDING:
		return ASSET_TYPE_FUND
	case XCOIN_ASSET_TYPE_TRADING:
		return ASSET_TYPE_UNIFIED
	case XCOIN_ASSET_TYPE_SECURITIES:
		return ASSET_TYPE_SECURITIES
	default:
		return ""
	}
}

func (c *XcoinEnumConverter) ToXcoinAssetType(t AssetType) string {
	switch t {
	case ASSET_TYPE_FUND:
		return XCOIN_ASSET_TYPE_FUNDING
	case ASSET_TYPE_UNIFIED:
		return XCOIN_ASSET_TYPE_TRADING
	case ASSET_TYPE_SECURITIES:
		return XCOIN_ASSET_TYPE_SECURITIES
	default:
		return ""
	}
}

// 划转状态类型 success, failed, pending
func (c *XcoinEnumConverter) FromXcoinTransferStatus(t string) TransferStatusType {
	switch t {
	case XCOIN_TRANSFER_STATUS_TYPE_SUCCESS:
		return TRANSFER_STATUS_TYPE_SUCCESS
	case XCOIN_TRANSFER_STATUS_TYPE_PENDING:
		return TRANSFER_STATUS_TYPE_PENDING
	case XCOIN_TRANSFER_STATUS_TYPE_FAILED:
		return TRANSFER_STATUS_TYPE_FAILED
	default:
		return TRANSFER_STATUS_TYPE_UNKNOWN
	}
}

func (c *XcoinEnumConverter) FromXcoinOrderType(orderType string, timeInForce string) (OrderType, TimeInForce) {
	if strings.ToLower(timeInForce) == XCOIN_TIME_IN_FORCE_GTC || timeInForce == "" {
		switch strings.ToLower(orderType) {
		case XCOIN_ORDER_TYPE_LIMIT:
			return ORDER_TYPE_LIMIT, TIME_IN_FORCE_GTC
		case XCOIN_ORDER_TYPE_MARKET:
			return ORDER_TYPE_MARKET, TIME_IN_FORCE_GTC
		case XCOIN_ORDER_TYPE_POST_ONLY:
			return ORDER_TYPE_LIMIT, TIME_IN_FORCE_POST_ONLY
		default:
			return ORDER_TYPE_UNKNOWN, TIME_IN_FORCE_UNKNOWN
		}
	}
	switch strings.ToLower(timeInForce) {
	case XCOIN_TIME_IN_FORCE_FOK:
		return ORDER_TYPE_LIMIT, TIME_IN_FORCE_FOK
	case XCOIN_TIME_IN_FORCE_IOC:
		return ORDER_TYPE_LIMIT, TIME_IN_FORCE_IOC
	case XCOIN_TIME_IN_FORCE_POST_ONLY:
		return ORDER_TYPE_LIMIT, TIME_IN_FORCE_POST_ONLY
	default:
		return ORDER_TYPE_UNKNOWN, TIME_IN_FORCE_UNKNOWN
	}
}

func (c *XcoinEnumConverter) ToXcoinOrderType(orderType OrderType, timeInForce TimeInForce) string {
	if timeInForce == TIME_IN_FORCE_GTC || timeInForce == "" {
		switch orderType {
		case ORDER_TYPE_MARKET:
			return XCOIN_ORDER_TYPE_MARKET
		case ORDER_TYPE_LIMIT:
			return XCOIN_ORDER_TYPE_LIMIT
		default:
			return ""
		}
	}
	switch timeInForce {
	case TIME_IN_FORCE_POST_ONLY:
		return XCOIN_ORDER_TYPE_POST_ONLY
	case TIME_IN_FORCE_FOK, TIME_IN_FORCE_IOC:
		return XCOIN_ORDER_TYPE_LIMIT
	default:
		return ""
	}
}

func (c *XcoinEnumConverter) FromXcoinOrderSide(side string) OrderSide {
	switch strings.ToLower(side) {
	case XCOIN_ORDER_SIDE_BUY:
		return ORDER_SIDE_BUY
	case XCOIN_ORDER_SIDE_SELL:
		return ORDER_SIDE_SELL
	default:
		return ORDER_SIDE_UNKNOWN
	}
}

func (c *XcoinEnumConverter) ToXcoinOrderSide(side OrderSide) string {
	switch side {
	case ORDER_SIDE_BUY:
		return XCOIN_ORDER_SIDE_BUY
	case ORDER_SIDE_SELL:
		return XCOIN_ORDER_SIDE_SELL
	default:
		return ""
	}
}

func (c *XcoinEnumConverter) FromXcoinPositionSide(posSide string) PositionSide {
	switch strings.ToLower(posSide) {
	default:
		return POSITION_SIDE_UNKNOWN
	}
}

func (c *XcoinEnumConverter) FromXcoinTimeInForce(timeInForce string) TimeInForce {
	switch strings.ToLower(timeInForce) {
	case XCOIN_TIME_IN_FORCE_GTC:
		return TIME_IN_FORCE_GTC
	case XCOIN_TIME_IN_FORCE_IOC:
		return TIME_IN_FORCE_IOC
	case XCOIN_TIME_IN_FORCE_FOK:
		return TIME_IN_FORCE_FOK
	case XCOIN_TIME_IN_FORCE_POST_ONLY:
		return TIME_IN_FORCE_POST_ONLY
	default:
		return TIME_IN_FORCE_UNKNOWN
	}
}

func (c *XcoinEnumConverter) ToXcoinTimeInForce(timeInForce TimeInForce) string {
	switch timeInForce {
	case TIME_IN_FORCE_GTC:
		return XCOIN_TIME_IN_FORCE_GTC
	case TIME_IN_FORCE_IOC:
		return XCOIN_TIME_IN_FORCE_IOC
	case TIME_IN_FORCE_FOK:
		return XCOIN_TIME_IN_FORCE_FOK
	case TIME_IN_FORCE_POST_ONLY:
		return XCOIN_TIME_IN_FORCE_POST_ONLY
	default:
		return ""
	}
}

func (c *XcoinEnumConverter) FromXcoinOrderStatus(status string) OrderStatus {
	switch strings.ToLower(status) {
	case XCOIN_ORDER_STATUS_NEW:
		return ORDER_STATUS_NEW
	case XCOIN_ORDER_STATUS_PARTIALLY_FILLED:
		return ORDER_STATUS_PARTIALLY_FILLED
	case XCOIN_ORDER_STATUS_FILLED:
		return ORDER_STATUS_FILLED
	case XCOIN_ORDER_STATUS_CANCELED, XCOIN_ORDER_STATUS_PARTIALLY_CANCELED:
		return ORDER_STATUS_CANCELED
	case XCOIN_ORDER_STATUS_REJECTED:
		return ORDER_STATUS_REJECTED
	case XCOIN_ORDER_STATUS_UNTRIGGERED:
		return ORDER_STATUS_UN_TRIGGERED
	case XCOIN_ORDER_STATUS_TRIGGERED:
		return ORDER_STATUS_TRIGGERED
	default:
		return ORDER_STATUS_UNKNOWN
	}
}
