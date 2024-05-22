package mytrade

import "github.com/shopspring/decimal"

type OrderParam struct {

	// 以下是必填参数
	AccountType string          //账户类型
	Symbol      string          //交易对
	Price       decimal.Decimal //价格
	Quantity    decimal.Decimal //数量
	OrderType   OrderType       //订单类型
	OrderSide   OrderSide       //买卖方向

	// 以下是可选参数
	OrderId       string       //交易所自动生成的订单ID
	ClientOrderId string       //用户自己生成的订单ID  不填则自动生成
	PositionSide  PositionSide //仓位方向 合约交易时必填
	TimeInForce   TimeInForce  //有效方式
}

func (o *OrderParam) SetAccountType(accountType string) *OrderParam {
	o.AccountType = accountType
	return o
}

func (o *OrderParam) SetSymbol(symbol string) *OrderParam {
	o.Symbol = symbol
	return o
}

func (o *OrderParam) SetPrice(price decimal.Decimal) *OrderParam {
	o.Price = price
	return o
}

func (o *OrderParam) SetQuantity(quantity decimal.Decimal) *OrderParam {
	o.Quantity = quantity
	return o
}

func (o *OrderParam) SetOrderType(orderType OrderType) *OrderParam {
	o.OrderType = orderType
	return o
}

func (o *OrderParam) SetOrderSide(orderSide OrderSide) *OrderParam {
	o.OrderSide = orderSide
	return o
}

func (o *OrderParam) SetOrderId(orderId string) *OrderParam {
	o.OrderId = orderId
	return o
}

func (o *OrderParam) SetClientOrderId(clientOrderId string) *OrderParam {
	o.ClientOrderId = clientOrderId
	return o
}

func (o *OrderParam) SetPositionSide(positionSide PositionSide) *OrderParam {
	o.PositionSide = positionSide
	return o
}

func (o *OrderParam) SetTimeInForce(timeInForce TimeInForce) *OrderParam {
	o.TimeInForce = timeInForce
	return o
}

type QueryHistoryParam struct {
	AccountType   string //账户类型
	Symbol        string //交易对
	OrderId       string //交易所自动生成的订单ID 选填
	ClientOrderId string //用户自己生成的订单ID 选填
	StartTime     int64  //开始时间 选填 默认返回7天内订单
	EndTime       int64  //结束时间 选填 默认返回7天内订单
	Limit         int    //限制返回的订单数量 选填 默认返回100条 最大100
}

func (q *QueryHistoryParam) SetAccountType(accountType string) *QueryHistoryParam {
	q.AccountType = accountType
	return q
}

func (q *QueryHistoryParam) SetSymbol(symbol string) *QueryHistoryParam {
	q.Symbol = symbol
	return q
}

func (q *QueryHistoryParam) SetOrderId(orderId string) *QueryHistoryParam {
	q.OrderId = orderId
	return q
}

func (q *QueryHistoryParam) SetClientOrderId(clientOrderId string) *QueryHistoryParam {
	q.ClientOrderId = clientOrderId
	return q
}

type QueryTradeParam struct {
}

type SubscribeOrderParam struct {
}
