package mytrade

type XcoinAccountType string

func (a XcoinAccountType) String() string {
	return string(a)
}

const (
	XCOIN_ACCOUNT_TYPE_SPOT             XcoinAccountType = "spot"             // 现货
	XCOIN_ACCOUNT_TYPE_LINEAR_PERPETUAL XcoinAccountType = "linear_perpetual" // 永续合约
	XCOIN_ACCOUNT_TYPE_LINEAR_FUTURES   XcoinAccountType = "linear_futures"   // 交割合约
)

//1s，1m，3m，5m，15m，30m，1h，2h，4h，6h，8h，12h，1d，3d，1w，1M

const (
	XCOIN_KLINE_INTERVAL_1s  = "1s"
	XCOIN_KLINE_INTERVAL_1m  = "1m"
	XCOIN_KLINE_INTERVAL_3m  = "3m"
	XCOIN_KLINE_INTERVAL_5m  = "5m"
	XCOIN_KLINE_INTERVAL_15m = "15m"
	XCOIN_KLINE_INTERVAL_30m = "30m"
	XCOIN_KLINE_INTERVAL_1h  = "1h"
	XCOIN_KLINE_INTERVAL_2h  = "2h"
	XCOIN_KLINE_INTERVAL_4h  = "4h"
	XCOIN_KLINE_INTERVAL_6h  = "6h"
	XCOIN_KLINE_INTERVAL_8h  = "8h"
	XCOIN_KLINE_INTERVAL_12h = "12h"
	XCOIN_KLINE_INTERVAL_1d  = "1d"
	XCOIN_KLINE_INTERVAL_3d  = "3d"
	XCOIN_KLINE_INTERVAL_1w  = "1w"
	XCOIN_KLINE_INTERVAL_1M  = "1M"
)

func xcoinGetMillisecondFromInterval(interval string) int64 {
	switch interval {
	case XCOIN_KLINE_INTERVAL_1s:
		return 1000
	case XCOIN_KLINE_INTERVAL_1m:
		return 60 * 1000
	case XCOIN_KLINE_INTERVAL_3m:
		return 3 * 60 * 1000
	case XCOIN_KLINE_INTERVAL_5m:
		return 5 * 60 * 1000
	case XCOIN_KLINE_INTERVAL_15m:
		return 15 * 60 * 1000
	case XCOIN_KLINE_INTERVAL_30m:
		return 30 * 60 * 1000
	case XCOIN_KLINE_INTERVAL_1h:
		return 60 * 60 * 1000
	case XCOIN_KLINE_INTERVAL_2h:
		return 2 * 60 * 60 * 1000
	case XCOIN_KLINE_INTERVAL_4h:
		return 4 * 60 * 60 * 1000
	case XCOIN_KLINE_INTERVAL_6h:
		return 6 * 60 * 60 * 1000
	case XCOIN_KLINE_INTERVAL_8h:
		return 8 * 60 * 60 * 1000
	case XCOIN_KLINE_INTERVAL_12h:
		return 12 * 60 * 60 * 1000
	case XCOIN_KLINE_INTERVAL_1d:
		return 24 * 60 * 60 * 1000
	case XCOIN_KLINE_INTERVAL_3d:
		return 3 * 24 * 60 * 60 * 1000
	case XCOIN_KLINE_INTERVAL_1w:
		return 7 * 24 * 60 * 60 * 1000
	case XCOIN_KLINE_INTERVAL_1M:
		return 30 * 24 * 60 * 60 * 1000
	default:
		return 60 * 1000
	}
}

func xcoinGetKlineCloseTime(ts int64, interval string) int64 {
	return ts + xcoinGetMillisecondFromInterval(interval) - 1
}

const (
	XCOIN_ASSET_TYPE_FUNDING    = "funding"
	XCOIN_ASSET_TYPE_TRADING    = "trading"
	XCOIN_ASSET_TYPE_SECURITIES = "securities"
)

// 资金划转类型
const (
	XCOIN_TRANSFER_STATUS_TYPE_SUCCESS = "success"
	XCOIN_TRANSFER_STATUS_TYPE_PENDING = "pending"
	XCOIN_TRANSFER_STATUS_TYPE_FAILED  = "failed"
)

// 订单类型
const (
	XCOIN_ORDER_TYPE_LIMIT     = "limit"
	XCOIN_ORDER_TYPE_MARKET    = "market"
	XCOIN_ORDER_TYPE_POST_ONLY = "post_only"
)

// 订单方向
const (
	XCOIN_ORDER_SIDE_BUY  = "buy"
	XCOIN_ORDER_SIDE_SELL = "sell"
)

// 时间类型
const (
	XCOIN_TIME_IN_FORCE_GTC       = "gtc"
	XCOIN_TIME_IN_FORCE_IOC       = "ioc"
	XCOIN_TIME_IN_FORCE_FOK       = "fok"
	XCOIN_TIME_IN_FORCE_POST_ONLY = "post_only"
)

// 订单状态
const (
	XCOIN_ORDER_STATUS_NEW                = "new"
	XCOIN_ORDER_STATUS_PARTIALLY_FILLED   = "partially_filled"
	XCOIN_ORDER_STATUS_FILLED             = "filled"
	XCOIN_ORDER_STATUS_PARTIALLY_CANCELED = "partially_canceled"
	XCOIN_ORDER_STATUS_CANCELED           = "canceled"
	XCOIN_ORDER_STATUS_REJECTED           = "rejected"
	XCOIN_ORDER_STATUS_UNTRIGGERED        = "untriggered"
	XCOIN_ORDER_STATUS_TRIGGERED          = "triggered"
)
