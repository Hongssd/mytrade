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
	ExchangeInfo() TradeExchangeInfo

	//获取市场数据
	MarketData() TradeMarketData

	//获取交易引擎
	TradeEngine(apiKey, secretKey string, options ...TradeEngineOption) TradeEngine
}

type TradeEngineOption func(e TradeEngine)

// 交易规范
type TradeExchangeInfo interface {
	TradeExchangeType
	//获取交易对规范
	GetSymbolInfo(accountType, symbol string) (TradeSymbolInfo, error)

	//刷新交易规范
	Refresh() error
}

// 交易对规范接口
type TradeSymbolInfo interface {
	//基础规范
	Exchange() string
	AccountType() string           //账户类型
	Symbol() string                //交易对名字
	BaseCoin() string              //交易币种
	QuoteCoin() string             //计价币种
	IsTrading() bool               //是否交易中
	IsContract() bool              //是否合约
	IsContractAmt() bool           //是否合约张数计量
	ContractSize() decimal.Decimal //合约面值
	ContractCoin() string          //合约面值计价币种
	ContractType() string          //合约类型

	//精度
	PricePrecision() int //价格精度
	AmtPrecision() int   //数量精度

	//价格规范
	TickSize() decimal.Decimal //下单价格精度
	MinPrice() decimal.Decimal //最小下单价格
	MaxPrice() decimal.Decimal //最大下单价格

	//数量规范
	//当合约张数计量时为张的数量
	LotSize() decimal.Decimal   //下单数量精度
	MinAmt() decimal.Decimal    //最小下单数量
	MaxLmtAmt() decimal.Decimal //最大限价单下单数量
	MaxMktAmt() decimal.Decimal //最大市价单下单数量

	//其他规范
	MaxLeverage() decimal.Decimal  //最大杠杆
	MinLeverage() decimal.Decimal  //最大杠杆
	StepLeverage() decimal.Decimal //最大杠杆
	MaxOrderNum() int              //最大订单数
	MinNotional() decimal.Decimal  //最小名义价值

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
	WsCancelOrder(*OrderParam) error

	//websocket批量下单
	WsCreateOrders([]*OrderParam) ([]*Order, error)
	//websocket批量修改订单
	WsAmendOrders([]*OrderParam) ([]*Order, error)
	//websocket批量撤单
	WsCancelOrders([]*OrderParam) error
}

type TradeSubscribe[T any] interface {
	ErrChan() chan error
	ResultChan() chan T
	CloseChan() chan struct{}
}
