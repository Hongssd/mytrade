package mytrade

type SunxAccountType string

func (a SunxAccountType) String() string {
	return string(a)
}

const (
	SUNX_AC_SWAP SunxAccountType = "swap"
)

const (
	SUNX_KLINE_INTERVAL_1m  = "1min"
	SUNX_KLINE_INTERVAL_5m  = "5min"
	SUNX_KLINE_INTERVAL_15m = "15min"
	SUNX_KLINE_INTERVAL_30m = "30min"
	SUNX_KLINE_INTERVAL_1H  = "60min"
	SUNX_KLINE_INTERVAL_4H  = "4hour"
	SUNX_KLINE_INTERVAL_1D  = "1day"
	SUNX_KLINE_INTERVAL_7D  = "1week"
	SUNX_KLINE_INTERVAL_30D = "1mon"
)

func sunxGetMillisecondFromInterval(interval string) int64 {
	switch interval {
	case SUNX_KLINE_INTERVAL_1m:
		return 60 * 1000
	case SUNX_KLINE_INTERVAL_5m:
		return 5 * 60 * 1000
	case SUNX_KLINE_INTERVAL_15m:
		return 15 * 60 * 1000
	case SUNX_KLINE_INTERVAL_30m:
		return 30 * 60 * 1000
	case SUNX_KLINE_INTERVAL_1H:
		return 60 * 60 * 1000
	case SUNX_KLINE_INTERVAL_1D:
		return 24 * 60 * 60 * 1000
	case SUNX_KLINE_INTERVAL_7D:
		return 7 * 24 * 60 * 60 * 1000
	case SUNX_KLINE_INTERVAL_30D:
		return 30 * 24 * 60 * 60 * 1000
	default:
		return 60 * 1000
	}
}

func sunxGetKlineCloseTime(ts int64, interval string) int64 {
	return ts + sunxGetMillisecondFromInterval(interval) - 1
}

// 仅统一账户模式
const (
	SUNX_ACCOUNT_MODE_UNIFIED = "unified" // 统一账户模式
)

const (
	SUNX_ACCOUNT_TYPE_SWAP     SunxAccountType = "swap"     // 永续合约
	SUNX_ACCOUNT_TYPE_DELIVERY SunxAccountType = "delivery" // 交割合约
)

const (
	SUNX_ASSET_TYPE_SWAP     = "swap"     // 永续合约
	SUNX_ASSET_TYPE_DELIVERY = "delivery" // 交割合约
)

const (
	SUNX_POSITION_MODE_SINGLE = "single_side" // 单向持仓
	SUNX_POSITION_MODE_HEDGE  = "dual_side"   // 双向持仓
)

const (
	SUNX_MARGIN_MODE_CROSSED = "cross" // 全仓
)

const (
	SUNX_ORDER_STATUS_NEW                = "new"
	SUNX_ORDER_STATUS_PARTIALLY_FILLED   = "partially_filled"
	SUNX_ORDER_STATUS_FILLED             = "filled"
	SUNX_ORDER_STATUS_PARTIALLY_CANCELED = "partially_canceled"
	SUNX_ORDER_STATUS_CANCELED           = "canceled"
	SUNX_ORDER_STATUS_REJECTED           = "rejected"
)

const (
	SUNX_ORDER_TYPE_LIMIT  = "limit"
	SUNX_ORDER_TYPE_MARKET = "market"
)

const (
	SUNX_ORDER_SIDE_BUY  = "buy"
	SUNX_ORDER_SIDE_SELL = "sell"
)

const (
	SUNX_POSITION_SIDE_LONG  = "long"
	SUNX_POSITION_SIDE_SHORT = "short"
	SUNX_POSITION_SIDE_BOTH  = "both"
)

const (
	SUNX_TIME_IN_FORCE_GTC = "gtc"
	SUNX_TIME_IN_FORCE_IOC = "ioc"
	SUNX_TIME_INFORCE_FOK  = "fok"
)

const (
	SUNX_ORDER_BOOK_STEP0  = "step0"
	SUNX_ORDER_BOOK_STEP1  = "step1"
	SUNX_ORDER_BOOK_STEP2  = "step2"
	SUNX_ORDER_BOOK_STEP3  = "step3"
	SUNX_ORDER_BOOK_STEP4  = "step4"
	SUNX_ORDER_BOOK_STEP5  = "step5"
	SUNX_ORDER_BOOK_STEP14 = "step14"
	SUNX_ORDER_BOOK_STEP15 = "step15"
	SUNX_ORDER_BOOK_STEP16 = "step16"
	SUNX_ORDER_BOOK_STEP17 = "step17"
)
