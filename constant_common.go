package mytrade

type ExchangeType string

func (e ExchangeType) String() string {
	return string(e)
}

const (
	BINANCE_NAME ExchangeType = "BINANCE"
	OKX_NAME     ExchangeType = "OKX"
	BYBIT_NAME   ExchangeType = "BYBIT"
)

type OrderStatus string

func (o OrderStatus) String() string {
	return string(o)
}

const (
	ORDER_STATUS_UNKNOWN          OrderStatus = ""
	ORDER_STATUS_NEW              OrderStatus = "NEW"              //新订单
	ORDER_STATUS_PARTIALLY_FILLED OrderStatus = "PARTIALLY_FILLED" //部分成交
	ORDER_STATUS_FILLED           OrderStatus = "FILLED"           //完全成交
	ORDER_STATUS_CANCELED         OrderStatus = "CANCELED"         //已撤销
	ORDER_STATUS_REJECTED         OrderStatus = "REJECTED"         //已拒绝
	ORDER_STATUS_UN_TRIGGERED     OrderStatus = "UN_TRIGGERED"     //未触发
	ORDER_STATUS_TRIGGERED        OrderStatus = "TRIGGERED"        //已触发
)

type OrderType string

func (o OrderType) String() string {
	return string(o)
}

const (
	ORDER_TYPE_UNKNOWN OrderType = ""
	ORDER_TYPE_LIMIT   OrderType = "LIMIT"  //限价单
	ORDER_TYPE_MARKET  OrderType = "MARKET" //市价单
)

type OrderSide string

func (o OrderSide) String() string {
	return string(o)
}

const (
	ORDER_SIDE_UNKNOWN OrderSide = ""
	ORDER_SIDE_BUY     OrderSide = "BUY"  //买
	ORDER_SIDE_SELL    OrderSide = "SELL" //卖
)

type PositionSide string

func (p PositionSide) String() string {
	return string(p)
}

const (
	POSITION_SIDE_UNKNOWN PositionSide = ""
	POSITION_SIDE_LONG    PositionSide = "LONG"  //多头
	POSITION_SIDE_SHORT   PositionSide = "SHORT" //空头
	POSITION_SIDE_BOTH    PositionSide = "BOTH"  //双向
)

type TimeInForce string

func (t TimeInForce) String() string {
	return string(t)
}

const (
	TIME_IN_FORCE_UNKNOWN   TimeInForce = ""
	TIME_IN_FORCE_GTC       TimeInForce = "GTC"       //成交为止, 一直有效
	TIME_IN_FORCE_IOC       TimeInForce = "IOC"       //立即成交或取消
	TIME_IN_FORCE_FOK       TimeInForce = "FOK"       //全部成交或立即取消
	TIME_IN_FORCE_POST_ONLY TimeInForce = "POST_ONLY" //只做maker
)

// 账户模式 无保证金/单币种保证金/多币种保证金/组合保证金
type AccountMode string

func (a AccountMode) String() string {
	return string(a)
}

const (
	ACCOUNT_MODE_UNKNOWN       AccountMode = ""          //未知
	ACCOUNT_MODE_FREE_MARGIN   AccountMode = "FREE"      //无保证金
	ACCOUNT_MODE_SINGLE_MARGIN AccountMode = "SINGLE"    //单币种保证金
	ACCOUNT_MODE_MULTI_MARGIN  AccountMode = "MULTI"     //多币种保证金
	ACCOUNT_MODE_PORTFOLIO     AccountMode = "PORTFOLIO" //组合保证金
)

// 仓位保证金模式 全仓 逐仓
type MarginMode string

func (m MarginMode) String() string {
	return string(m)
}

const (
	MARGIN_MODE_UNKNOWN  MarginMode = ""         //未知
	MARGIN_MODE_CROSSED  MarginMode = "CROSSED"  //全仓
	MARGIN_MODE_ISOLATED MarginMode = "ISOLATED" //逐仓
)

// 仓位模式 双向持仓/单向持仓
type PositionMode string

func (p PositionMode) String() string {
	return string(p)
}

const (
	POSITION_MODE_UNKNOWN PositionMode = ""        //未知
	POSITION_MODE_HEDGE   PositionMode = "HEDGE"   //双向持仓
	POSITION_MODE_ONEWAY  PositionMode = "ONE_WAY" //单向持仓
)

// 止盈止损触发类型
type OrderTriggerType string

func (p OrderTriggerType) String() string {
	return string(p)
}

const (
	ORDER_TRIGGER_TYPE_UNKNOWN     OrderTriggerType = ""
	ORDER_TRIGGER_TYPE_STOP_LOSS   OrderTriggerType = "STOP_LOSS"   //止损
	ORDER_TRIGGER_TYPE_TAKE_PROFIT OrderTriggerType = "TAKE_PROFIT" //止盈
)

// 触发条件类型 上穿 下穿
type OrderTriggerConditionType string

func (p OrderTriggerConditionType) String() string {
	return string(p)
}

const (
	ORDER_TRIGGER_CONDITION_TYPE_UNKNOWN      OrderTriggerConditionType = ""
	ORDER_TRIGGER_CONDITION_TYPE_THROUGH_UP   OrderTriggerConditionType = "THROUGH_UP"   //上穿
	ORDER_TRIGGER_CONDITION_TYPE_THROUGH_DOWN OrderTriggerConditionType = "THROUGH_DOWN" //下穿
)

// 账户类型
type AssetType string

func (p AssetType) String() string {
	return string(p)
}

const (
	ASSET_TYPE_FUND     AssetType = "FUND"     //资金账户
	ASSET_TYPE_UNIFIED  AssetType = "UNIFIED"  // 统一账户
	ASSET_TYPE_CONTRACT AssetType = "CONTRACT" // 合约账户
	ASSET_TYPE_UMFUTURE AssetType = "UMFUTURE" // U本位合约账户
	ASSET_TYPE_CMFUTURE AssetType = "CMFUTURE" // 币本位合约账户
)

// 划转状态类型
type TransferStatusType string

func (p TransferStatusType) String() string {
	return string(p)
}

const (
	TRANSFER_STATUS_TYPE_UNKNOWN TransferStatusType = "UNKNOWN"
	TRANSFER_STATUS_TYPE_PENDING TransferStatusType = "PENDING"
	TRANSFER_STATUS_TYPE_SUCCESS TransferStatusType = "SUCCESS"
	TRANSFER_STATUS_TYPE_FAILED  TransferStatusType = "FAILED"
)
