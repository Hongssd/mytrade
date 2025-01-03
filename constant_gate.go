package mytrade

type GateAccountType string

func (b GateAccountType) String() string {
	return string(b)
}

const (
	GATE_AC_SPOT     GateAccountType = "SPOT"     //现货
	GATE_AC_FUTURES  GateAccountType = "FUTURES"  //合约
	GATE_AC_DELIVERY GateAccountType = "DELIVERY" //交割
)

// 参数	值
// interval	10s
// interval	1m
// interval	5m
// interval	15m
// interval	30m
// interval	1h
// interval	4h
// interval	8h
// interval	1d
// interval	7d
// interval	30d
const (
	GATE_KLINE_INTERVAL_1m  = "1m"
	GATE_KLINE_INTERVAL_3m  = "3m"
	GATE_KLINE_INTERVAL_5m  = "5m"
	GATE_KLINE_INTERVAL_15m = "15m"
	GATE_KLINE_INTERVAL_30m = "30m"
	GATE_KLINE_INTERVAL_1H  = "1h"
	GATE_KLINE_INTERVAL_2H  = "2h"
	GATE_KLINE_INTERVAL_4H  = "4h"
	GATE_KLINE_INTERVAL_6H  = "6h"
	GATE_KLINE_INTERVAL_8H  = "8h"
	GATE_KLINE_INTERVAL_12H = "12h"
	GATE_KLINE_INTERVAL_1D  = "1d"
	GATE_KLINE_INTERVAL_2D  = "2d"
	GATE_KLINE_INTERVAL_3D  = "3d"
	GATE_KLINE_INTERVAL_5D  = "5d"
	GATE_KLINE_INTERVAL_7D  = "7d"
	GATE_KLINE_INTERVAL_30D = "30d"
)

func gateGetMillisecondFromInterval(interval string) int64 {
	switch interval {
	case GATE_KLINE_INTERVAL_1m:
		return 60 * 1000
	case GATE_KLINE_INTERVAL_3m:
		return 3 * 60 * 1000
	case GATE_KLINE_INTERVAL_5m:
		return 5 * 60 * 1000
	case GATE_KLINE_INTERVAL_15m:
		return 15 * 60 * 1000
	case GATE_KLINE_INTERVAL_30m:
		return 30 * 60 * 1000
	case GATE_KLINE_INTERVAL_1H:
		return 60 * 60 * 1000
	case GATE_KLINE_INTERVAL_2H:
		return 2 * 60 * 60 * 1000
	case GATE_KLINE_INTERVAL_4H:
		return 4 * 60 * 60 * 1000
	case GATE_KLINE_INTERVAL_6H:
		return 6 * 60 * 60 * 1000
	case GATE_KLINE_INTERVAL_8H:
		return 8 * 60 * 60 * 1000
	case GATE_KLINE_INTERVAL_12H:
		return 12 * 60 * 60 * 1000
	case GATE_KLINE_INTERVAL_1D:
		return 24 * 60 * 60 * 1000
	case GATE_KLINE_INTERVAL_2D:
		return 2 * 24 * 60 * 60 * 1000
	case GATE_KLINE_INTERVAL_3D:
		return 3 * 24 * 60 * 60 * 1000
	case GATE_KLINE_INTERVAL_5D:
		return 5 * 24 * 60 * 60 * 1000
	case GATE_KLINE_INTERVAL_7D:
		return 7 * 24 * 60 * 60 * 1000
	case GATE_KLINE_INTERVAL_30D:
		return 30 * 24 * 60 * 60 * 1000

	default:
		return 60 * 1000
	}
}

func gateGetKlineCloseTime(ts int64, interval string) int64 {
	return ts + gateGetMillisecondFromInterval(interval) - 1
}
