package mytrade

type BybitTradeEngine struct {
	exchangeBase

	apiKey    string
	secretKey string
}

func (b BybitTradeEngine) NewOrderReq() *OrderParam {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) NewQueryOrderReq() *QueryOrderParam {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) NewQueryTradeReq() *QueryTradeParam {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) QueryOpenOrders(param *QueryOrderParam) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) QueryOrder(param *QueryOrderParam) (*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) QueryTrades(param *QueryTradeParam) ([]*Trade, error) {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) CreateOrder(param *OrderParam) (*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) AmendOrder(param *OrderParam) (*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) CancelOrder(param *OrderParam) (*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) CreateOrders(params []*OrderParam) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) AmendOrders(params []*OrderParam) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) CancelOrders(params []*OrderParam) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) NewSubscribeOrderReq() *SubscribeOrderParam {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) SubscribeOrder(param *SubscribeOrderParam) (TradeSubscribe[Order], error) {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) WsCreateOrder(param *OrderParam) (*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) WsAmendOrder(param *OrderParam) (*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) WsCancelOrder(param *OrderParam) (*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) WsCreateOrders(params []*OrderParam) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) WsAmendOrders(params []*OrderParam) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) WsCancelOrders(params []*OrderParam) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}
