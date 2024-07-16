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
	OrderId          string       //交易所自动生成的订单ID
	ClientOrderId    string       //用户自己生成的订单ID  不填则自动生成
	PositionSide     PositionSide //仓位方向 合约交易时必填
	TimeInForce      TimeInForce  //有效方式
	ReduceOnly       bool         //是否只减仓
	NewClientOrderId string       //新的用户自己生成的订单ID 改单时可用
	IsIsolated       bool         //是否是逐仓模式

	//止盈止损
	TriggerPrice decimal.Decimal  //止盈止损触发价
	TriggerType  OrderTriggerType //触发类型
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
func (o *OrderParam) SetReduceOnly(reduceOnly bool) *OrderParam {
	o.ReduceOnly = reduceOnly
	return o
}
func (o *OrderParam) SetNewClientOrderId(newClientOrderId string) *OrderParam {
	o.NewClientOrderId = newClientOrderId
	return o
}
func (o *OrderParam) SetIsIsolated(isIsolated bool) *OrderParam {
	o.IsIsolated = isIsolated
	return o
}
func (o *OrderParam) SetTriggerPrice(triggerPrice decimal.Decimal) *OrderParam {
	o.TriggerPrice = triggerPrice
	return o
}
func (o *OrderParam) SetOrderTriggerType(orderTriggerType OrderTriggerType) *OrderParam {
	o.TriggerType = orderTriggerType
	return o
}

type QueryOrderParam struct {
	AccountType   string //账户类型
	Symbol        string //交易对
	BaseCoin      string //交易幣種
	SettleCoin    string //結算幣種
	OrderId       string //交易所自动生成的订单ID 选填
	ClientOrderId string //用户自己生成的订单ID 选填
	StartTime     int64  //开始时间 选填 默认返回7天内订单
	EndTime       int64  //结束时间 选填 默认返回7天内订单
	Limit         int    //限制返回的订单数量 选填 默认返回100条 最大100

}

func (q *QueryOrderParam) SetAccountType(accountType string) *QueryOrderParam {
	q.AccountType = accountType
	return q
}

func (q *QueryOrderParam) SetSymbol(symbol string) *QueryOrderParam {
	q.Symbol = symbol
	return q
}

func (q *QueryOrderParam) SetBaseCoin(baseCoin string) *QueryOrderParam {
	q.BaseCoin = baseCoin
	return q
}

func (q *QueryOrderParam) SetSettleCoin(settleCoin string) *QueryOrderParam {
	q.SettleCoin = settleCoin
	return q
}

func (q *QueryOrderParam) SetOrderId(orderId string) *QueryOrderParam {
	q.OrderId = orderId
	return q
}

func (q *QueryOrderParam) SetClientOrderId(clientOrderId string) *QueryOrderParam {
	q.ClientOrderId = clientOrderId
	return q
}

func (q *QueryOrderParam) SetStartTime(startTime int64) *QueryOrderParam {
	q.StartTime = startTime
	return q
}

func (q *QueryOrderParam) SetEndTime(endTime int64) *QueryOrderParam {
	q.EndTime = endTime
	return q
}

func (q *QueryOrderParam) SetLimit(limit int) *QueryOrderParam {
	q.Limit = limit
	return q
}

type QueryTradeParam struct {
	AccountType string //账户类型
	Symbol      string //交易对
	OrderId     string //订单ID
	StartTime   int64  //开始时间
	EndTime     int64  //结束时间
	Limit       int    //限制返回的成交明细数量
}

func (q *QueryTradeParam) SetAccountType(accountType string) *QueryTradeParam {
	q.AccountType = accountType
	return q
}
func (q *QueryTradeParam) SetSymbol(symbol string) *QueryTradeParam {
	q.Symbol = symbol
	return q
}
func (q *QueryTradeParam) SetOrderId(orderId string) *QueryTradeParam {
	q.OrderId = orderId
	return q
}
func (q *QueryTradeParam) SetStartTime(startTime int64) *QueryTradeParam {
	q.StartTime = startTime
	return q
}
func (q *QueryTradeParam) SetEndTime(endTime int64) *QueryTradeParam {
	q.EndTime = endTime
	return q
}
func (q *QueryTradeParam) SetLimit(limit int) *QueryTradeParam {
	q.Limit = limit
	return q
}

type SubscribeOrderParam struct {
	AccountType string //账户类型
}

func (s *SubscribeOrderParam) SetAccountType(accountType string) *SubscribeOrderParam {
	s.AccountType = accountType
	return s
}
