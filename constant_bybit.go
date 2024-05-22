package mytrade

type BybitAccountType string

func (b BybitAccountType) String() string {
	return string(b)
}

const (
	BYBIT_AC_SPOT    BybitAccountType = "spot"    //现货
	BYBIT_AC_LINEAR  BybitAccountType = "linear"  //永续合约
	BYBIT_AC_INVERSE BybitAccountType = "inverse" //交割合约
	BYBIT_AC_OPTION  BybitAccountType = "option"  //期权
)

// 時間粒度. 1,3,5,15,30,60,120,240,360,720,D,M,W
const (
	BYBIT_KLINE_INTERVAL_1m  = "1"
	BYBIT_KLINE_INTERVAL_3m  = "3"
	BYBIT_KLINE_INTERVAL_5m  = "5"
	BYBIT_KLINE_INTERVAL_15m = "15"
	BYBIT_KLINE_INTERVAL_30m = "30"
	BYBIT_KLINE_INTERVAL_1H  = "60"
	BYBIT_KLINE_INTERVAL_2H  = "120"
	BYBIT_KLINE_INTERVAL_4H  = "240"
	BYBIT_KLINE_INTERVAL_6H  = "360"
	BYBIT_KLINE_INTERVAL_12H = "720"
	BYBIT_KLINE_INTERVAL_1D  = "D"
	BYBIT_KLINE_INTERVAL_1W  = "W"
	BYBIT_KLINE_INTERVAL_1M  = "M"
)

func bybitGetMillisecondFromInterval(interval string) int64 {
	switch interval {
	case BYBIT_KLINE_INTERVAL_1m:
		return 60 * 1000
	case BYBIT_KLINE_INTERVAL_3m:
		return 3 * 60 * 1000
	case BYBIT_KLINE_INTERVAL_5m:
		return 5 * 60 * 1000
	case BYBIT_KLINE_INTERVAL_15m:
		return 15 * 60 * 1000
	case BYBIT_KLINE_INTERVAL_30m:
		return 30 * 60 * 1000
	case BYBIT_KLINE_INTERVAL_1H:
		return 60 * 60 * 1000
	case BYBIT_KLINE_INTERVAL_2H:
		return 2 * 60 * 60 * 1000
	case BYBIT_KLINE_INTERVAL_4H:
		return 4 * 60 * 60 * 1000
	case BYBIT_KLINE_INTERVAL_6H:
		return 6 * 60 * 60 * 1000
	case BYBIT_KLINE_INTERVAL_12H:
		return 12 * 60 * 60 * 1000
	case BYBIT_KLINE_INTERVAL_1D:
		return 24 * 60 * 60 * 1000
	case BYBIT_KLINE_INTERVAL_1W:
		return 7 * 24 * 60 * 60 * 1000
	case BYBIT_KLINE_INTERVAL_1M:
		return 30 * 24 * 60 * 60 * 1000
	default:
		return 60 * 1000
	}
}
func bybitGetKlineCloseTime(ts int64, interval string) int64 {
	return ts + bybitGetMillisecondFromInterval(interval) - 1
}
