package mytrade

type OkxAccountType string

func (b OkxAccountType) String() string {
	return string(b)
}

const (
	OKX_AC_SPOT    OkxAccountType = "SPOT"    //现货
	OKX_AC_MARGIN  OkxAccountType = "MARGIN"  //币币
	OKX_AC_SWAP    OkxAccountType = "SWAP"    //永续合约
	OKX_AC_FUTURES OkxAccountType = "FUTURES" //交割合约
	OKX_AC_OPTION  OkxAccountType = "OPTION"  //期权
)

// 订单类型
// market：市价单
// limit：限价单
// post_only：只做maker单
// fok：全部成交或立即取消
// ioc：立即成交并取消剩余
// optimal_limit_ioc：市价委托立即成交并取消剩余（仅适用交割、永续）
// mmp：做市商保护(仅适用于组合保证金账户模式下的期权订单)
// mmp_and_post_only：做市商保护且只做maker单(仅适用于组合保证金账户模式下的期权订单)
const (
	OKX_ORDER_TYPE_MARKET    = "market"
	OKX_ORDER_TYPE_LIMIT     = "limit"
	OKX_ORDER_TYPE_POST_ONLY = "post_only"
	OKX_ORDER_TYPE_FOK       = "fok"
	OKX_ORDER_TYPE_IOC       = "ioc"
)

const (
	OKX_ORDER_SIDE_BUY  = "buy"
	OKX_ORDER_SIDE_SELL = "sell"
)

const (
	OKX_POSITION_SIDE_LONG  = "long"
	OKX_POSITION_SIDE_SHORT = "short"
	OKX_POSITION_SIDE_BOTH  = "net"
)

// 订单状态
// canceled：撤单成功
// live：等待成交
// partially_filled：部分成交
// filled：完全成交
// mmp_canceled：做市商保护机制导致的自动撤单
const (
	OKX_ORDER_STATUS_NEW              = "live"
	OKX_ORDER_STATUS_PARTIALLY_FILLED = "partially_filled"
	OKX_ORDER_STATUS_FILLED           = "filled"
	OKX_ORDER_STATUS_CANCELED         = "canceled"
	OKX_ORDER_STATUS_REJECTED         = "mmp_canceled"
)

// 时间粒度，默认值1m
// 如 [1m/3m/5m/15m/30m/1H/2H/4H]
// 香港时间开盘价k线：[6H/12H/1D/2D/3D/1W/1M/3M]
// UTC时间开盘价k线：[/6Hutc/12Hutc/1Dutc/2Dutc/3Dutc/1Wutc/1Mutc/3Mutc]
const (
	OKX_KLINE_INTERVAL_1m     = "1m"
	OKX_KLINE_INTERVAL_3m     = "3m"
	OKX_KLINE_INTERVAL_5m     = "5m"
	OKX_KLINE_INTERVAL_15m    = "15m"
	OKX_KLINE_INTERVAL_30m    = "30m"
	OKX_KLINE_INTERVAL_1H     = "1H"
	OKX_KLINE_INTERVAL_2H     = "2H"
	OKX_KLINE_INTERVAL_4H     = "4H"
	OKX_KLINE_INTERVAL_6H     = "6H"
	OKX_KLINE_INTERVAL_12H    = "12H"
	OKX_KLINE_INTERVAL_1D     = "1D"
	OKX_KLINE_INTERVAL_2D     = "2D"
	OKX_KLINE_INTERVAL_3D     = "3D"
	OKX_KLINE_INTERVAL_1W     = "1W"
	OKX_KLINE_INTERVAL_1M     = "1M"
	OKX_KLINE_INTERVAL_3M     = "3M"
	OKX_KLINE_INTERVAL_6Hutc  = "6Hutc"
	OKX_KLINE_INTERVAL_12Hutc = "12Hutc"
	OKX_KLINE_INTERVAL_1Dutc  = "1Dutc"
	OKX_KLINE_INTERVAL_2Dutc  = "2Dutc"
	OKX_KLINE_INTERVAL_3Dutc  = "3Dutc"
	OKX_KLINE_INTERVAL_1Wutc  = "1Wutc"
	OKX_KLINE_INTERVAL_1Mutc  = "1Mutc"
	OKX_KLINE_INTERVAL_3Mutc  = "3Mutc"
)

func okxGetMillisecondFromInterval(interval string) int64 {
	switch interval {
	case OKX_KLINE_INTERVAL_1m:
		return 60 * 1000
	case OKX_KLINE_INTERVAL_3m:
		return 3 * 60 * 1000
	case OKX_KLINE_INTERVAL_5m:
		return 5 * 60 * 1000
	case OKX_KLINE_INTERVAL_15m:
		return 15 * 60 * 1000
	case OKX_KLINE_INTERVAL_30m:
		return 30 * 60 * 1000
	case OKX_KLINE_INTERVAL_1H:
		return 60 * 60 * 1000
	case OKX_KLINE_INTERVAL_2H:
		return 2 * 60 * 60 * 1000
	case OKX_KLINE_INTERVAL_4H:
		return 4 * 60 * 60 * 1000
	case OKX_KLINE_INTERVAL_6H:
		return 6 * 60 * 60 * 1000
	case OKX_KLINE_INTERVAL_12H:
		return 12 * 60 * 60 * 1000
	case OKX_KLINE_INTERVAL_1D:
		return 24 * 60 * 60 * 1000
	case OKX_KLINE_INTERVAL_2D:
		return 2 * 24 * 60 * 60 * 1000
	case OKX_KLINE_INTERVAL_3D:
		return 3 * 24 * 60 * 60 * 1000
	case OKX_KLINE_INTERVAL_1W:
		return 7 * 24 * 60 * 60 * 1000
	case OKX_KLINE_INTERVAL_1M:
		return 30 * 24 * 60 * 60 * 1000
	case OKX_KLINE_INTERVAL_3M:
		return 3 * 30 * 24 * 60 * 60 * 1000
	case OKX_KLINE_INTERVAL_6Hutc:
		return 6 * 60 * 60 * 1000
	case OKX_KLINE_INTERVAL_12Hutc:
		return 12 * 60 * 60 * 1000
	case OKX_KLINE_INTERVAL_1Dutc:
		return 24 * 60 * 60 * 1000
	case OKX_KLINE_INTERVAL_2Dutc:
		return 2 * 24 * 60 * 1000
	case OKX_KLINE_INTERVAL_3Dutc:
		return 3 * 24 * 60 * 60 * 1000
	case OKX_KLINE_INTERVAL_1Wutc:
		return 7 * 24 * 60 * 60 * 1000
	case OKX_KLINE_INTERVAL_1Mutc:
		return 30 * 24 * 60 * 60 * 1000
	case OKX_KLINE_INTERVAL_3Mutc:
		return 3 * 30 * 24 * 60 * 60 * 1000
	default:
		return 60 * 1000
	}
}

func okxGetKlineCloseTime(ts int64, interval string) int64 {
	return ts + okxGetMillisecondFromInterval(interval) - 1
}
