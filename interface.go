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
	GetKline(req *KlineParam) (*[]Kline, error)
	//查询深度
	GetBook(req *BookParam) (*OrderBook, error)
}

// 交易引擎接口
type TradeEngine interface {
	TradeExchangeType

	//新建订单请求参数
	NewOrderReq() *OrderParam
	//新建订单查询请求参数
	NewQueryOrderReq() *QueryHistoryParam
	//新建成交查询请求参数
	NewQueryTradeReq() *QueryTradeParam

	//查挂单
	QueryOpenOrders(req *QueryHistoryParam) ([]*Order, error)
	//查指定单
	QueryOrder(req *QueryHistoryParam) (*Order, error)
	//查成交
	QueryTrades(req *QueryTradeParam) ([]*Trade, error)

	//下单
	CreateOrder(req *OrderParam) (*Order, error)
	//修改订单
	AmendOrder(req *OrderParam) (*Order, error)
	//撤单
	CancelOrder(req *OrderParam) error

	//批量下单
	CreateOrders(reqs []*OrderParam) ([]*Order, error)
	//批量修改订单
	AmendOrders(reqs []*OrderParam) ([]*Order, error)
	//批量撤单
	CancelOrders(reqs []*OrderParam) error

	//websocket相关
	//开启交易websocket
	OpenOrderWs() error
	//是否已连接交易websocket
	IsConnectedWs() bool
	//关闭交易websocket
	CloseOrderWs() error

	//新建交易订单订阅请求参数
	NewSubscribeOrderReq() *SubscribeOrderParam

	//开启订单交易订阅
	SubscribeOrder(req *SubscribeOrderParam) (TradeSubscribe[Order], error)

	//websocket下单
	WsCreateOrder(req *OrderParam) (*Order, error)
	//websocket修改订单
	WsAmendOrder(req *OrderParam) (*Order, error)
	//websocket撤单
	WsCancelOrder(req *OrderParam) error

	//websocket批量下单
	WsCreateOrders(reqs []*OrderParam) ([]*Order, error)
	//websocket批量修改订单
	WsAmendOrders(reqs []*OrderParam) ([]*Order, error)
	//websocket批量撤单
	WsCancelOrders(reqs []*OrderParam) error
}

type TradeSubscribe[T any] interface {
	ErrChan() chan error
	ResultChan() chan T
	CloseChan() chan struct{}
}
