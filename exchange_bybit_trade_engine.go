package mytrade

type BybitTradeEngine struct {
	exchangeBase

	bybitConverter BybitEnumConverter

	broadcasterSpot    *bybitOrderBroadcaster
	broadcasterLinear  *bybitOrderBroadcaster
	broadcasterInverse *bybitOrderBroadcaster

	apiKey    string
	secretKey string
}

func (o *BybitTradeEngine) NewOrderReq() *OrderParam {
	//TODO implement me
	panic("implement me")
}

func (o *BybitTradeEngine) NewQueryOrderReq() *QueryOrderParam {
	//TODO implement me
	panic("implement me")
}

func (o *BybitTradeEngine) NewQueryTradeReq() *QueryTradeParam {
	//TODO implement me
	panic("implement me")
}

func (o *BybitTradeEngine) QueryOpenOrders(param *QueryOrderParam) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *BybitTradeEngine) QueryOrder(param *QueryOrderParam) (*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *BybitTradeEngine) QueryTrades(param *QueryTradeParam) ([]*Trade, error) {
	//TODO implement me
	panic("implement me")
}

func (o *BybitTradeEngine) CreateOrder(param *OrderParam) (*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *BybitTradeEngine) AmendOrder(param *OrderParam) (*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *BybitTradeEngine) CancelOrder(param *OrderParam) (*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *BybitTradeEngine) CreateOrders(params []*OrderParam) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *BybitTradeEngine) AmendOrders(params []*OrderParam) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *BybitTradeEngine) CancelOrders(params []*OrderParam) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *BybitTradeEngine) NewSubscribeOrderReq() *SubscribeOrderParam {
	return &SubscribeOrderParam{}
}

func (o *BybitTradeEngine) SubscribeOrder(req *SubscribeOrderParam) (TradeSubscribe[Order], error) {
	switch BybitAccountType(req.AccountType) {
	case BYBIT_AC_SPOT, BYBIT_AC_LINEAR, BYBIT_AC_INVERSE:
	default:
		return nil, ErrorAccountType
	}
	b := o.getBroadcastFromAccountType(req.AccountType)

	sub, err := o.newOrderSubscriber(b, "", req.AccountType, "")
	if err != nil {
		return nil, err
	}

	middleSub := &subscription[Order]{
		resultChan: make(chan Order, 100),
		errChan:    make(chan error, 10),
		closeChan:  make(chan struct{}, 10),
	}

	//循环将订单数据中转到目标订阅
	go func() {
		for {
			select {
			case <-sub.ch.CloseChan():
				middleSub.closeChan <- struct{}{}
				return
			case err := <-sub.ch.ErrChan():
				middleSub.errChan <- err
			case order := <-sub.ch.ResultChan():
				middleSub.resultChan <- order
			}
		}
	}()

	return middleSub, nil
}

func (o *BybitTradeEngine) WsCreateOrder(param *OrderParam) (*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *BybitTradeEngine) WsAmendOrder(param *OrderParam) (*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *BybitTradeEngine) WsCancelOrder(param *OrderParam) (*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *BybitTradeEngine) WsCreateOrders(params []*OrderParam) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *BybitTradeEngine) WsAmendOrders(params []*OrderParam) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *BybitTradeEngine) WsCancelOrders(params []*OrderParam) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}
