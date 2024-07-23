package mytrade

type BinanceAccountType string

func (b BinanceAccountType) String() string {
	return string(b)
}

const (
	BN_AC_FUNDING BinanceAccountType = "FUNDING"
	BN_AC_SPOT    BinanceAccountType = "SPOT"
	BN_AC_FUTURE  BinanceAccountType = "FUTURE"
	BN_AC_SWAP    BinanceAccountType = "SWAP"
)

const (
	BN_ORDER_TYPE_LIMIT                  = "LIMIT"
	BN_ORDER_TYPE_MARKET                 = "MARKET"
	BN_ORDER_TYPE_SPOT_STOP_LOSS_LIMIT   = "STOP_LOSS_LIMIT"
	BN_ORDER_TYPE_SPOT_TAKE_PROFIT_LIMIT = "TAKE_PROFIT_LIMIT"

	BN_ORDER_TYPE_FUTURE_STOP               = "STOP"
	BN_ORDER_TYPE_FUTURE_TAKE_PROFIT        = "TAKE_PROFIT"
	BN_ORDER_TYPE_FUTURE_STOP_MARKET        = "STOP_MARKET"
	BN_ORDER_TYPE_FUTURE_TAKE_PROFIT_MARKET = "TAKE_PROFIT_MARKET"
)

const (
	BN_ORDER_SIDE_BUY  = "BUY"
	BN_ORDER_SIDE_SELL = "SELL"
)

const (
	BN_POSITION_SIDE_LONG  = "LONG"
	BN_POSITION_SIDE_SHORT = "SHORT"
	BN_POSITION_SIDE_BOTH  = "BOTH"
)

const (
	BN_ORDER_STATUS_NEW              = "NEW"
	BN_ORDER_STATUS_PARTIALLY_FILLED = "PARTIALLY_FILLED"
	BN_ORDER_STATUS_FILLED           = "FILLED"
	BN_ORDER_STATUS_CANCELED         = "CANCELED"
	BN_ORDER_STATUS_REJECTED         = "REJECTED"
	BN_ORDER_STATUS_EXPIRED          = "EXPIRED"
)

const (
	BN_TIME_IN_FORCE_GTC       = "GTC"
	BN_TIME_IN_FORCE_IOC       = "IOC"
	BN_TIME_IN_FORCE_FOK       = "FOK"
	BN_TIME_IN_FORCE_POST_ONLY = "GTX"
)

const (
	BN_WORKING_TYPE_MARK_PRICE     = "MARK_PRICE"
	BN_WORKING_TYPE_CONTRACT_PRICE = "CONTRACT_PRICE"
)

// {
// "multiAssetsMargin": true // "true": 联合保证金模式开启；"false": 联合保证金模式关闭
// }
const (
	BN_ACCOUNT_MOED_MULTI_CURRENCY_MARGIN  = true
	BN_ACCOUNT_MODE_SINGLE_CURRENCY_MARGIN = false
)

// "isolated": true,  // 是否是逐仓模式
const (
	BN_MARGIN_MODE_ISOLATED = true
	BN_MARGIN_MODE_CROSSED  = false
)

const (
	BN_MARGIN_MODE_ISOLATED_STR = "ISOLATED"
	BN_MARGIN_MODE_CROSSED_STR  = "CROSSED"
)

// {
// "dualSidePosition": true // "true": 双向持仓模式；"false": 单向持仓模式
// }
const (
	BN_POSITION_MODE_HEDGE  = true
	BN_POSITION_MODE_ONEWAY = false
)
