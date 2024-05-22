package mytrade

type BybitTradeEngine struct {
	exchangeBase

	apiKey    string
	secretKey string
}

func (b BybitTradeEngine) IsConnectedWs() bool {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) NewOrderReq() *OrderParam {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) NewQueryOrderReq() *QueryHistoryParam {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) NewQueryTradeReq() *QueryTradeParam {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) QueryOpenOrders(req *QueryHistoryParam) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) QueryOrder(req *QueryHistoryParam) (*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) QueryTrades(req *QueryTradeParam) ([]*Trade, error) {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) CreateOrder(req *OrderParam) (*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) AmendOrder(req *OrderParam) (*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) CancelOrder(req *OrderParam) error {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) CreateOrders(reqs []*OrderParam) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) AmendOrders(reqs []*OrderParam) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) CancelOrders(reqs []*OrderParam) error {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) OpenOrderWs() error {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) CloseOrderWs() error {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) NewSubscribeOrderReq() *SubscribeOrderParam {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) SubscribeOrder(req *SubscribeOrderParam) (TradeSubscribe[Order], error) {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) WsCreateOrder(req *OrderParam) (*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) WsAmendOrder(req *OrderParam) (*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) WsCancelOrder(req *OrderParam) error {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) WsCreateOrders(reqs []*OrderParam) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) WsAmendOrders(reqs []*OrderParam) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b BybitTradeEngine) WsCancelOrders(reqs []*OrderParam) error {
	//TODO implement me
	panic("implement me")
}
