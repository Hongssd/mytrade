package mytrade

type Order struct {
	Exchange      string       `json:"exchange"`      //交易所名字
	AccountType   string       `json:"accountType"`   //账户类型
	Symbol        string       `json:"symbol"`        //交易对名字
	IsMargin      bool         `json:"isMargin"`      //是否是杠杆订单
	IsIsolated    bool         `json:"isIsolated"`    //是否是逐仓symbol交易（杠杆）
	OrderId       string       `json:"orderId"`       //交易所自动生成的订单ID
	ClientOrderId string       `json:"clientOrderId"` //用户自己生成的订单ID  不填则自动生成
	Price         string       `json:"price"`         //下单价格
	Quantity      string       `json:"quantity"`      //下单数量
	ExecutedQty   string       `json:"executedQty"`   //成交数量
	CumQuoteQty   string       `json:"cumQuoteQty"`   //成交金额
	AvgPrice      string       `json:"avgPrice"`      //成交均价
	Status        OrderStatus  `json:"status"`        //订单状态
	Type          OrderType    `json:"orderType"`     //订单类型
	Side          OrderSide    `json:"orderSide"`     //买卖方向
	PositionSide  PositionSide `json:"positionSide"`  //仓位方向 合约交易时存在
	TimeInForce   TimeInForce  `json:"timeInForce"`   //有效方式
	FeeAmount     string       `json:"feeAmount"`     //手续费数量
	FeeCcy        string       `json:"feeCcy"`        //手续费币种
	ReduceOnly    bool         `json:"reduceOnly"`    //是否只减仓
	CreateTime    int64        `json:"createTime"`    //创建时间
	UpdateTime    int64        `json:"updateTime"`    //更新时间
	RealizedPnl   string       `json:"fillPnl"`       //成交盈亏

	IsAlgo               bool                      `json:"isAlgo"`               //是否是策略订单
	TriggerPrice         string                    `json:"triggerPrice"`         //触发价格
	TriggerType          OrderTriggerType          `json:"triggerType"`          //触发类型
	TriggerConditionType OrderTriggerConditionType `json:"triggerConditionType"` //触发条件类型

	MarginBuyBorrowAmount string `json:"marginBuyBorrowAmount"` //下单后没有发生借款则不返回该字段（杠杆）
	MarginBuyBorrowAsset  string `json:"marginBuyBorrowAsset"`  //下单后没有发生借款则不返回该字段（杠杆）

	ErrorCode string `json:"errorCode"` //错误码
	ErrorMsg  string `json:"errorMsg"`  //错误信息
}

type Trade struct {
	Exchange      string       `json:"exchange"`      //交易所名字
	AccountType   string       `json:"accountType"`   //账户类型
	Symbol        string       `json:"symbol"`        //交易对名字
	TradeId       string       `json:"tradeId"`       //交易所自动生成的交易ID
	OrderId       string       `json:"orderId"`       //交易所自动生成的订单ID
	ClientOrderId string       `json:"clientOrderId"` //用户自己生成的订单ID
	Price         string       `json:"price"`         //成交价格
	Quantity      string       `json:"quantity"`      //成交数量
	QuoteQty      string       `json:"quoteQty"`      //成交金额
	Side          OrderSide    `json:"side"`          //买卖方向
	PositionSide  PositionSide `json:"positionSide"`  //仓位方向 合约交易时存在
	FeeAmount     string       `json:"feeAmount"`     //手续费数量
	FeeCcy        string       `json:"feeCcy"`        //手续费币种
	RealizedPnl   string       `json:"realizedPnl"`   //实现盈亏
	IsMaker       bool         `json:"isMaker"`       //是否maker成交
	Timestamp     int64        `json:"timestamp"`     //成交时间
}
