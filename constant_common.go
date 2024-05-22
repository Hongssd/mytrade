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

type OrderType string

func (o OrderType) String() string {
	return string(o)
}

const (
	ORDER_TYPE_LIMIT  OrderType = "LIMIT"  //限价单
	ORDER_TYPE_MARKET OrderType = "MARKET" //市价单
)

type OrderSide string

func (o OrderSide) String() string {
	return string(o)
}

const (
	ORDER_SIDE_BUY  OrderSide = "BUY"  //买
	ORDER_SIDE_SELL OrderSide = "SELL" //卖
)

type PositionSide string

func (p PositionSide) String() string {
	return string(p)
}

const (
	POSITION_SIDE_LONG  PositionSide = "LONG"  //多头
	POSITION_SIDE_SHORT PositionSide = "SHORT" //空头
)

type TimeInForce string

func (t TimeInForce) String() string {
	return string(t)
}

const (
	TIME_IN_FORCE_GTC       TimeInForce = "GTC"       //成交为止, 一直有效
	TIME_IN_FORCE_IOC       TimeInForce = "IOC"       //立即成交或取消
	TIME_IN_FORCE_FOK       TimeInForce = "FOK"       //全部成交或立即取消
	TIME_IN_FORCE_POST_ONLY TimeInForce = "POST_ONLY" //只做maker
)
