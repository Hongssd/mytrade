package mytrade

import (
	"github.com/shopspring/decimal"
)

type TradeExchangeType interface {
	//交易所类型
	ExchangeType() ExchangeType
}

type TradeExchange interface {
	TradeExchangeType
	//获取交易规范
	NewExchangeInfo() TradeExchangeInfo
	//获取市场数据
	NewMarketData() TradeMarketData
	//获取交易引擎
	NewTradeEngine(apiKey, secretKey, passphrase string) TradeEngine
	//获取账户信息
	NewTradeAccount(apiKey, secretKey, passphrase string) TradeAccount
}

// 交易规范
type TradeExchangeInfo interface {
	TradeExchangeType
	//获取交易对规范
	GetSymbolInfo(accountType, symbol string) (TradeSymbolInfo, error)

	//获取全部交易对规范
	GetAllSymbolInfo(accountType string) ([]TradeSymbolInfo, error)

	//刷新交易规范
	Refresh() error
}

// 交易对规范接口
type TradeSymbolInfo interface {
	//基础规范
	Exchange() string
	//账户类型
	AccountType() string
	//交易对名字
	Symbol() string
	//交易币种
	BaseCoin() string
	//计价币种
	QuoteCoin() string
	//是否交易中
	IsTrading() bool
	//是否合约
	IsContract() bool
	//是否合约张数计量
	IsContractAmt() bool
	//合约面值
	ContractSize() decimal.Decimal
	//合约面值计价币种
	ContractCoin() string
	//合约类型
	ContractType() string

	//精度
	//价格精度
	PricePrecision() int
	//数量精度
	AmtPrecision() int

	//价格规范
	//下单价格精度
	TickSize() decimal.Decimal
	//最小下单价格
	MinPrice() decimal.Decimal
	//最大下单价格
	MaxPrice() decimal.Decimal

	//数量规范
	//当合约张数计量时为张的数量
	//下单数量精度
	LotSize() decimal.Decimal
	//最小下单数量
	MinAmt() decimal.Decimal
	//最大限价单下单数量
	MaxLmtAmt() decimal.Decimal
	//最大市价单下单数量
	MaxMktAmt() decimal.Decimal

	//其他规范
	//最大杠杆
	MaxLeverage() decimal.Decimal
	//最小杠杆
	MinLeverage() decimal.Decimal
	//杠杆步长
	StepLeverage() decimal.Decimal
	//最大订单数
	MaxOrderNum() int
	//最小名义价值
	MinNotional() decimal.Decimal

	MarshalJson() ([]byte, error)
	MarshalJsonIndent(prefix, indent string) ([]byte, error)
}

// 市场数据接口
type TradeMarketData interface {
	TradeExchangeType
	//新建K线请求参数
	NewKlineReq() *KlineParam
	//新建深度请求参数
	NewBookReq() *BookParam
	//查询K线
	GetKline(*KlineParam) (*[]Kline, error)
	//查询深度
	GetBook(*BookParam) (*OrderBook, error)
}

// 交易引擎接口
type TradeEngine interface {
	TradeExchangeType

	//新建订单请求参数
	NewOrderReq() *OrderParam
	//新建订单查询请求参数
	NewQueryOrderReq() *QueryOrderParam
	//新建成交查询请求参数
	NewQueryTradeReq() *QueryTradeParam

	//查挂单
	QueryOpenOrders(*QueryOrderParam) ([]*Order, error)
	//查指定单
	QueryOrder(*QueryOrderParam) (*Order, error)

	//查订单列表
	QueryOrders(*QueryOrderParam) ([]*Order, error)

	//查成交
	QueryTrades(*QueryTradeParam) ([]*Trade, error)

	//下单
	CreateOrder(*OrderParam) (*Order, error)
	//修改订单
	AmendOrder(*OrderParam) (*Order, error)
	//撤单
	CancelOrder(*OrderParam) (*Order, error)

	//批量下单
	CreateOrders([]*OrderParam) ([]*Order, error)
	//批量修改订单
	AmendOrders([]*OrderParam) ([]*Order, error)
	//批量撤单
	CancelOrders([]*OrderParam) ([]*Order, error)

	//新建交易订单订阅请求参数
	NewSubscribeOrderReq() *SubscribeOrderParam

	//开启订单交易订阅
	SubscribeOrder(*SubscribeOrderParam) (TradeSubscribe[Order], error)

	//websocket下单
	WsCreateOrder(*OrderParam) (*Order, error)
	//websocket修改订单
	WsAmendOrder(*OrderParam) (*Order, error)
	//websocket撤单
	WsCancelOrder(*OrderParam) (*Order, error)

	//websocket批量下单
	WsCreateOrders([]*OrderParam) ([]*Order, error)
	//websocket批量修改订单
	WsAmendOrders([]*OrderParam) ([]*Order, error)
	//websocket批量撤单
	WsCancelOrders([]*OrderParam) ([]*Order, error)
}

// 交易账户接口
type TradeAccount interface {
	//获取账户模式 无保证金/单币种保证金/多币种保证金/组合保证金
	GetAccountMode() (AccountMode, error)
	//获取保证金模式 全仓/逐仓
	GetMarginMode(accountType, symbol string, positionSide PositionSide) (MarginMode, error)
	//获取持仓模式 单向/多向
	GetPositionMode(accountType, symbol string) (PositionMode, error)
	//获取杠杆
	GetLeverage(accountType, symbol string,
		marginMode MarginMode, positionSide PositionSide) (decimal.Decimal, error)

	//获取手续费率,taker maker
	GetFeeRate(accountType, symbol string) (*FeeRate, error)
	//获取持仓
	GetPositions(accountType string, symbols ...string) ([]*Position, error)
	//获取资产
	GetAssets(accountType string, currencies ...string) ([]*Asset, error)

	//设置账户模式
	SetAccountMode(mode AccountMode) error
	//设置保证金模式
	SetMarginMode(accountType, symbol string, mode MarginMode) error
	//设置持仓模式
	SetPositionMode(accountType, symbol string, mode PositionMode) error
	//设置杠杆
	SetLeverage(accountType, symbol string,
		marginMode MarginMode, positionSide PositionSide,
		leverage decimal.Decimal) error

	//资金划转
	AssetTransfer(AssetTransferParams *AssetTransferParams) ([]*AssetTransfer, error)
	//查询资金划转历史
	QueryAssetTransfer(req *QueryAssetTransferParams) ([]*QueryAssetTransfer, error)
}

type TradeSubscribe[T any] interface {
	ErrChan() chan error
	ResultChan() chan T
	CloseChan() chan struct{}
}
