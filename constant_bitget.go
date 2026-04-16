package mytrade

type BitgetAccountType string

func (b BitgetAccountType) String() string {
	return string(b)
}

// Account types aligned with mybitgetapi@v0.0.1 InstType.
const (
	BITGET_AC_SPOT            = "SPOT"
	BITGET_AC_MARGIN          = "MARGIN"
	BITGET_AC_MARGIN_CROSSED  = "MARGIN_CROSSED"  // 非交易所枚举类型 仅用于区分账号资产类型接口
	BITGET_AC_MARGIN_ISOLATED = "MARGIN_ISOLATED" // 非交易所枚举类型 仅用于区分账号资产类型接口
	BITGET_AC_USDT_FUTURES    = "USDT-FUTURES"
	BITGET_AC_COIN_FUTURES    = "COIN-FUTURES"
	BITGET_AC_USDC_FUTURES    = "USDC-FUTURES"
	BITGET_AC_UTA             = "UTA"
)

// Account Mode
const (
	BITGET_AM_SINGLE = "single"
	BITGET_AM_UNION  = "union"
)

// Order side.
const (
	BITGET_ORDER_SIDE_BUY  = "buy"
	BITGET_ORDER_SIDE_SELL = "sell"
)

// Order type.
const (
	BITGET_ORDER_TYPE_LIMIT  = "limit"
	BITGET_ORDER_TYPE_MARKET = "market"
)

// Time in force.
const (
	BITGET_TIF_GTC       = "gtc"
	BITGET_TIF_IOC       = "ioc"
	BITGET_TIF_FOK       = "fok"
	BITGET_TIF_POST_ONLY = "post_only"
)

// Position side.
const (
	BITGET_POS_SIDE_LONG  = "long"
	BITGET_POS_SIDE_SHORT = "short"
	BITGET_POS_SIDE_NET   = "net"
)

// reduceOnly
const (
	BITGET_REDUCE_ONLY_YES = "YES"
	BITGET_REDUCE_ONLY_NO  = "NO"
)

// Classic 杠杆下单 loanType（与 Bitget 文档一致；可用 OrderParam.SideEffectType 覆盖）
const (
	BITGET_LOAN_TYPE_NORMAL = "normal"
)

// Classic 合约 tradeSide
const (
	BITGET_TRADE_SIDE_OPEN  = "open"
	BITGET_TRADE_SIDE_CLOSE = "close"
)

// UTA / WS common order statuses.
const (
	BITGET_ORDER_STATUS_NEW              = "new"
	BITGET_ORDER_STATUS_LIVE             = "live"
	BITGET_ORDER_STATUS_PARTIALLY_FILLED = "partially_filled"
	BITGET_ORDER_STATUS_FILLED           = "filled"
	BITGET_ORDER_STATUS_CANCELED         = "canceled"
	BITGET_ORDER_STATUS_CANCELLED        = "cancelled"
	BITGET_ORDER_STATUS_REJECTED         = "rejected"
)

// Classic common order statuses.
const (
	BITGET_CLASSIC_ORDER_STATUS_NEW          = "new"
	BITGET_CLASSIC_ORDER_STATUS_PARTIAL_FILL = "partial_fill"
	BITGET_CLASSIC_ORDER_STATUS_FULL_FILL    = "full_fill"
)

// Margin / position modes aligned with SDK.
const (
	BITGET_MARGIN_MODE_ISOLATED = "isolated"
	BITGET_MARGIN_MODE_CROSSED  = "crossed"

	BITGET_POSITION_MODE_ONE_WAY = "one_way_mode"
	BITGET_POSITION_MODE_HEDGE   = "hedge_mode"
)

// Account switch/upgrade statuses.
const (
	BITGET_ACCOUNT_SWITCH_STATUS_PROCESS = "process"
	BITGET_ACCOUNT_SWITCH_STATUS_SUCCESS = "success"
	BITGET_ACCOUNT_SWITCH_STATUS_FAIL    = "fail"
)
