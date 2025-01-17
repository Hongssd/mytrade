package mytrade

type GateAccountType string

func (b GateAccountType) String() string {
	return string(b)
}

const (
	GATE_AC_SPOT     GateAccountType = "SPOT"     //现货
	GATE_AC_MARGIN   GateAccountType = "MARGIN"   //现货杠杆
	GATE_AC_UNIFIED  GateAccountType = "UNIFIED"  //统一账户
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

const (
	GATE_ACCOUNT_MODE_FREE_MARGIN  = 1 //现货模式
	GATE_ACCOUNT_MODE_MULTI_MARGIN = 2 //跨币种保证金模式
)

const (
	GATE_ACCOUNT_TYPE_SPOT     = "spot"     //现货
	GATE_ACCOUNT_TYPE_MARGIN   = "margin"   //现货杠杆
	GATE_ACCOUNT_TYPE_FUTURES  = "futures"  //合约
	GATE_ACCOUNT_TYPE_DELIVERY = "delivery" //交割
	GATE_ACCOUNT_TYPE_UNIFIED  = "unified"  //统一账户
	GATE_ACCOUNT_TYPE_UNKNOWN  = ""         //未知
)

const (
	GATE_ASSET_TYPE_SPOT     = "spot"     //现货
	GATE_ASSET_TYPE_MARGIN   = "margin"   //现货杠杆
	GATE_ASSET_TYPE_FUTURES  = "futures"  //合约
	GATE_ASSET_TYPE_DELIVERY = "delivery" //交割
	GATE_ASSET_TYPE_UNFIED   = "unified"  //统一账户
)

const (
	GATE_POSITION_MODE_ONEWAY      = "single"     //单向持仓
	GATE_POSITION_MODE_HEDGE_LONG  = "dual_long"  //双向持仓多头
	GATE_POSITION_MODE_HEDGE_SHORT = "dual_short" //双向持仓空头
)

const (
	GATE_ORDER_SIDE_BUY  = "buy"
	GATE_ORDER_SIDE_SELL = "sell"
)

const (
	GATE_ORDER_TYPE_LIMIT  = "limit"
	GATE_ORDER_TYPE_MARKET = "market"
)

const (
	GATE_TIME_IN_FORCE_GTC = "gtc" // - GoodTillCancelled
	GATE_TIME_IN_FORCE_IOC = "ioc" // - ioc: ImmediateOrCancelled ，立即成交或者取消，只吃单不挂单
	GATE_TIME_IN_FORCE_POC = "poc" // - poc: PendingOrCancelled，被动委托，只挂单不吃单
)

const (
	GATE_ORDER_STATUS_NEW       = "open"
	GATE_ORDER_STATUS_FILLED    = "closed"
	GATE_ORDER_STATUS_CANCELLED = "cancelled"
)
