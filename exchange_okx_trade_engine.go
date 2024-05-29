package mytrade

type OkxTradeEngine struct {
	exchangeBase

	apiKey     string
	secretKey  string
	passphrase string
}

func (o OkxTradeEngine) NewOrderReq() *OrderParam {
	//TODO implement me
	panic("implement me")
}

func (o OkxTradeEngine) NewQueryOrderReq() *QueryOrderParam {
	//TODO implement me
	panic("implement me")
}

func (o OkxTradeEngine) NewQueryTradeReq() *QueryTradeParam {
	//TODO implement me
	panic("implement me")
}

func (o OkxTradeEngine) QueryOpenOrders(req *QueryOrderParam) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o OkxTradeEngine) QueryOrder(req *QueryOrderParam) (*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o OkxTradeEngine) QueryTrades(req *QueryTradeParam) ([]*Trade, error) {
	//TODO implement me
	panic("implement me")
}

func (o OkxTradeEngine) CreateOrder(req *OrderParam) (*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o OkxTradeEngine) AmendOrder(req *OrderParam) (*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o OkxTradeEngine) CancelOrder(req *OrderParam) (*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o OkxTradeEngine) CreateOrders(reqs []*OrderParam) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o OkxTradeEngine) AmendOrders(reqs []*OrderParam) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o OkxTradeEngine) CancelOrders(reqs []*OrderParam) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o OkxTradeEngine) NewSubscribeOrderReq() *SubscribeOrderParam {
	//TODO implement me
	panic("implement me")
}

func (o OkxTradeEngine) SubscribeOrder(req *SubscribeOrderParam) (TradeSubscribe[Order], error) {
	//TODO implement me
	panic("implement me")
}

func (o OkxTradeEngine) WsCreateOrder(req *OrderParam) (*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o OkxTradeEngine) WsAmendOrder(req *OrderParam) (*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o OkxTradeEngine) WsCancelOrder(req *OrderParam) error {
	//TODO implement me
	panic("implement me")
}

func (o OkxTradeEngine) WsCreateOrders(reqs []*OrderParam) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o OkxTradeEngine) WsAmendOrders(reqs []*OrderParam) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o OkxTradeEngine) WsCancelOrders(reqs []*OrderParam) error {
	//TODO implement me
	panic("implement me")
}
