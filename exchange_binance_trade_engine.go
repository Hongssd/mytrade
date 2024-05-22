package mytrade

type BinanceTradeEngine struct {
	exchangeBase

	apiKey    string
	secretKey string
}

func (b BinanceTradeEngine) IsConnectedWs() bool {
	//TODO implement me
	panic("implement me")
}

func (b BinanceTradeEngine) NewOrderReq() *OrderParam {
	//TODO implement me
	panic("implement me")
}

func (b BinanceTradeEngine) NewQueryOrderReq() *QueryHistoryParam {
	//TODO implement me
	panic("implement me")
}

func (b BinanceTradeEngine) NewQueryTradeReq() *QueryTradeParam {
	//TODO implement me
	panic("implement me")
}

func (b BinanceTradeEngine) QueryOpenOrders(req *QueryHistoryParam) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b BinanceTradeEngine) QueryOrder(req *QueryHistoryParam) (*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b BinanceTradeEngine) QueryTrades(req *QueryTradeParam) ([]*Trade, error) {
	//TODO implement me
	panic("implement me")
}

func (b BinanceTradeEngine) CreateOrder(req *OrderParam) (*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b BinanceTradeEngine) AmendOrder(req *OrderParam) (*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b BinanceTradeEngine) CancelOrder(req *OrderParam) error {
	//TODO implement me
	panic("implement me")
}

func (b BinanceTradeEngine) CreateOrders(reqs []*OrderParam) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b BinanceTradeEngine) AmendOrders(reqs []*OrderParam) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b BinanceTradeEngine) CancelOrders(reqs []*OrderParam) error {
	//TODO implement me
	panic("implement me")
}

func (b BinanceTradeEngine) OpenOrderWs() error {
	//TODO implement me
	panic("implement me")
}

func (b BinanceTradeEngine) CloseOrderWs() error {
	//TODO implement me
	panic("implement me")
}

func (b BinanceTradeEngine) NewSubscribeOrderReq() *SubscribeOrderParam {
	//TODO implement me
	panic("implement me")
}

func (b BinanceTradeEngine) SubscribeOrder(req *SubscribeOrderParam) (TradeSubscribe[Order], error) {
	//TODO implement me
	panic("implement me")
}

func (b BinanceTradeEngine) WsCreateOrder(req *OrderParam) (*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b BinanceTradeEngine) WsAmendOrder(req *OrderParam) (*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b BinanceTradeEngine) WsCancelOrder(req *OrderParam) error {
	//TODO implement me
	panic("implement me")
}

func (b BinanceTradeEngine) WsCreateOrders(reqs []*OrderParam) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b BinanceTradeEngine) WsAmendOrders(reqs []*OrderParam) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (b BinanceTradeEngine) WsCancelOrders(reqs []*OrderParam) error {
	//TODO implement me
	panic("implement me")
}
