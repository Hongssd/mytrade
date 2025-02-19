package mytrade

type GateAccountType string

func (b GateAccountType) String() string {
	return string(b)
}

// const (
// 	GATE_ACCOUNT_TYPE_SPOT     GateAccountType = "SPOT"     //现货
// 	GATE_ACCOUNT_TYPE_MARGIN   GateAccountType = "MARGIN"   //现货杠杆
// 	GATE_ACCOUNT_TYPE_UNIFIED  GateAccountType = "UNIFIED"  //统一账户
// 	GATE_ACCOUNT_TYPE_FUTURES  GateAccountType = "FUTURES"  //合约
// 	GATE_ACCOUNT_TYPE_DELIVERY GateAccountType = "DELIVERY" //交割
// )

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

// 统一账户模式：
// - classic: 经典账户模式
// - multi_currency: 跨币种保证金模式
// - portfolio: 组合保证金模式
// - single_currency: 单币种保证金模式
const (
	GATE_ACCOUNT_MODE_CLASSIC       = "classic"         //经典保证金模式
	GATE_ACCOUNT_MODE_MULTI_MARGIN  = "multi_currency"  //跨币种保证金模式
	GATE_ACCOUNT_MODE_PORTFOLIO     = "portfolio"       //组合保证金模式
	GATE_ACCOUNT_MODE_SINGLE_MARGIN = "single_currency" //单币种保证金模式
)

const (
	GATE_ACCOUNT_TYPE_SPOT         GateAccountType = "spot"         //现货
	GATE_ACCOUNT_TYPE_MARGIN       GateAccountType = "margin"       //现货逐仓杠杆
	GATE_ACCOUNT_TYPE_CROSS_MARGIN GateAccountType = "cross_margin" //现货全仓杠杆
	GATE_ACCOUNT_TYPE_FUTURES      GateAccountType = "futures"      //合约
	GATE_ACCOUNT_TYPE_DELIVERY     GateAccountType = "delivery"     //交割
	GATE_ACCOUNT_TYPE_UNIFIED      GateAccountType = "unified"      //统一账户
	GATE_ACCOUNT_TYPE_UNKNOWN      GateAccountType = ""             //未知

)

const (
	GATE_ASSET_TYPE_SPOT            = "spot"         //现货
	GATE_ASSET_TYPE_ISOLATED_MARGIN = "margin"       //现货逐仓杠杆
	GATE_ASSET_TYPE_CROSS_MARGIN    = "cross_margin" //现货全仓杠杆
	GATE_ASSET_TYPE_FUTURES         = "futures"      //合约
	GATE_ASSET_TYPE_DELIVERY        = "delivery"     //交割
	GATE_ASSET_TYPE_UNFIED          = "unified"      //统一账户
)

const (
	GATE_POSITION_MODE_ONEWAY = false //单向持仓
	GATE_POSITION_MODE_HEDGE  = true  //双向持仓
)

const (
	GATE_POSITION_SIDE_BOTH  = "single"     //单向持仓
	GATE_POSITION_SIDE_LONG  = "dual_long"  //双向持仓多头
	GATE_POSITION_SIDE_SHORT = "dual_short" //双向持仓空头
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
	GATE_ORDER_SPOT_STATUS_NEW       = "open"
	GATE_ORDER_SPOT_STATUS_FILLED    = "closed"
	GATE_ORDER_SPOT_STATUS_CANCELLED = "cancelled"
)

const (
	GATE_ORDER_CONTRACT_STATUS_OPEN     = "open"
	GATE_ORDER_CONTRACT_STATUS_FINISHED = "finished"
)

const (
	GATE_ORDER_SPOT_PRICE_STATUS_OPEN     = "open"     // 正在运行
	GATE_ORDER_SPOT_PRICE_STATUS_CANCELED = "canceled" // 被取消
	GATE_ORDER_SPOT_PRICE_STATUS_FINISHED = "finish"   // 成功结束
	GATE_ORDER_SPOT_PRICE_STATUS_FAILED   = "failed"   // 失败
	GATE_ORDER_SPOT_PRICE_STATUS_EXPIRED  = "expired"  // 过期
)

const (
	GATE_ORDER_FUTURES_PRICE_STATUS_OPEN     = "open"     // 活跃中
	GATE_ORDER_FUTURES_PRICE_STATUS_FINISHED = "finished" // 已结束
	GATE_ORDER_FUTURES_PRICE_STATUS_INACTIVE = "inactive" // 未生效
	GATE_ORDER_FUTURES_PRICE_STATUS_INVALID  = "invalid"  // 无效
)

// 结束状态，cancelled - 被取消；succeeded - 成功；failed - 失败；expired - 过期
const (
	GATE_ORDER_FUTURES_PRICE_FINISH_AS_CANCELLED = "cancelled" // 被取消
	GATE_ORDER_FUTURES_PRICE_FINISH_AS_SUCCEEDED = "succeeded" // 成功
	GATE_ORDER_FUTURES_PRICE_FINISH_AS_FAILED    = "failed"    // 失败
	GATE_ORDER_FUTURES_PRICE_FINISH_AS_EXPIRED   = "expired"   // 过期
)

const (
	GATE_SPOT_PRICE_ORDER_ACCOUNT_NORMAL         = "normal"       // 现货交易
	GATE_SPOT_PRICE_ORDER_ACCOUNT_MARGIN         = "margin"       // 杠杆交易
	GATE_SPOT_PRICE_ORDER_ACCOUNT_CROSSED_MARGIN = "cross_margin" // 全仓杠杆交易
)

// gate合约订单
// 结束方式，包括：
// - filled: 完全成交
// - cancelled: 用户撤销
// - liquidated: 强制平仓撤销
// - ioc: 未立即完全成交，因为tif设置为ioc
// - auto_deleveraged: 自动减仓撤销
// - reduce_only: 增持仓位撤销，因为设置reduce_only或平仓
// - position_closed: 因为仓位平掉了，所以挂单被撤掉
// - reduce_out: 只减仓被排除的不容易成交的挂单
// - stp: 订单发生自成交限制而被撤销
const (
	GATE_ORDER_CONTRACT_FINISH_AS_FILLED           = "filled"
	GATE_ORDER_CONTRACT_FINISH_AS_CANCELLED        = "cancelled"
	GATE_ORDER_CONTRACT_FINISH_AS_LIQUIDATED       = "liquidated"
	GATE_ORDER_CONTRACT_FINISH_AS_IOC              = "ioc"
	GATE_ORDER_CONTRACT_FINISH_AS_AUTO_DELEVERAGED = "auto_deleveraged"
	GATE_ORDER_CONTRACT_FINISH_AS_REDUCE_ONLY      = "reduce_only"
)

const (
	GATE_SPOT_PRICE_ORDER_TRIGGER_RULE_LTE = "<="
	GATE_SPOT_PRICE_ORDER_TRIGGER_RULE_GTE = ">="
)

const (
	GATE_FUTURES_PRICE_ORDER_TRIGGER_RULE_GTE = 1
	GATE_FUTURES_PRICE_ORDER_TRIGGER_RULE_LTE = 2
)

const (
	GATE_PRICE_ORDER_TYPE_CLOSE_LONG_ORDER          = "close-long-order"
	GATE_PRICE_ORDER_TYPE_CLOSE_SHORT_ORDER         = "close-short-order"
	GATE_PRICE_ORDER_TYPE_CLOSE_LONG_POSITION       = "close-long-position"
	GATE_PRICE_ORDER_TYPE_CLOSE_SHORT_POSITION      = "close-short-position"
	GATE_PRICE_ORDER_TYPE_PLAN_CLOSE_LONG_POSITION  = "plan-close-long-position"
	GATE_PRICE_ORDER_TYPE_PLAN_CLOSE_SHORT_POSITION = "plan-close-short-position"
)
