package mytrade

import "github.com/Hongssd/mysunxapi"

type SunxTradeEngine struct {
	ExchangeBase

	sunxConverter SunxEnumConverter
	accessKey     string
	secretKey     string

	wsForSwapOrder *mysunxapi.PrivateWsStreamClient
}

func (s *SunxTradeEngine) NewOrderReq() *OrderParam {
	return &OrderParam{}
}
func (s *SunxTradeEngine) NewQueryOrderReq() *QueryOrderParam {
	return &QueryOrderParam{}
}
func (s *SunxTradeEngine) NewQueryTradeReq() *QueryTradeParam {
	return &QueryTradeParam{}
}

func (s *SunxTradeEngine) QueryOpenOrders(req *QueryOrderParam) ([]*Order, error) {
	if err := s.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}

	switch SunxAccountType(req.AccountType) {
	case SUNX_ACCOUNT_TYPE_SWAP:
		api := s.apiQueryOpenOrders(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		return s.handleOrdersFromQueryOpenOrders(req, res), nil
	default:
		return nil, ErrorAccountType
	}
}

func (s *SunxTradeEngine) QueryOrder(req *QueryOrderParam) (*Order, error) {
	if err := s.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}

	switch SunxAccountType(req.AccountType) {
	case SUNX_ACCOUNT_TYPE_SWAP:
		api := s.apiQueryOrder(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		return s.handleOrderFromQueryOrder(req, res)
	default:
		return nil, ErrorAccountType
	}
}

func (s *SunxTradeEngine) QueryOrders(req *QueryOrderParam) ([]*Order, error) {
	if err := s.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}

	switch SunxAccountType(req.AccountType) {
	case SUNX_ACCOUNT_TYPE_SWAP:
		api := s.apiQueryOrders(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		return s.handleOrdersFromQueryOrders(req, res), nil
	default:
		return nil, ErrorAccountType
	}
}

func (s *SunxTradeEngine) QueryTrades(req *QueryTradeParam) ([]*Trade, error) {
	if err := s.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}

	switch SunxAccountType(req.AccountType) {
	case SUNX_ACCOUNT_TYPE_SWAP:
		api := s.apiQueryTrades(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		return s.handleTradesFromQueryTrades(req, res), nil
	default:
		return nil, ErrorAccountType
	}
}

func (s *SunxTradeEngine) CreateOrder(req *OrderParam) (*Order, error) {
	if err := s.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}

	switch SunxAccountType(req.AccountType) {
	case SUNX_ACCOUNT_TYPE_SWAP:
		api := s.apiOrderCreate(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		return s.handleOrderFromOrderCreate(req, res)
	default:
		return nil, ErrorAccountType
	}
}

func (s *SunxTradeEngine) AmendOrder(req *OrderParam) (*Order, error) {
	if err := s.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}
	switch SunxAccountType(req.AccountType) {
	case SUNX_ACCOUNT_TYPE_SWAP:
		// 查单
		queryReq := &QueryOrderParam{
			AccountType:   req.AccountType,
			Symbol:        req.Symbol,
			OrderId:       req.OrderId,
			ClientOrderId: req.ClientOrderId,
		}
		queryApi := s.apiQueryOrder(queryReq)
		queryRes, err := queryApi.Do()
		if err != nil {
			return nil, err
		}
		currOrder, err := s.handleOrderFromQueryOrder(queryReq, queryRes)
		if err != nil {
			return nil, err
		}

		// 撤单
		cancelReq := &OrderParam{
			AccountType:   req.AccountType,
			Symbol:        queryRes.Data.ContractCode,
			OrderId:       queryRes.Data.OrderId,
			ClientOrderId: queryRes.Data.ClientOrderId,
		}
		cancelApi := s.apiOrderCancel(cancelReq)
		_, err = cancelApi.Do()
		if err != nil {
			return nil, err
		}
		log.Warn("cancel success")
		// 下单
		amendReq, amendApi := s.apiAmendOrderCreate(currOrder, req) // 改单
		res, err := amendApi.Do()
		if err != nil {
			return nil, err
		}
		return s.handleOrderFromOrderCreate(amendReq, res)
	default:
		return nil, ErrorAccountType
	}
}

func (s *SunxTradeEngine) CancelOrder(req *OrderParam) (*Order, error) {
	if err := s.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}
	switch SunxAccountType(req.AccountType) {
	case SUNX_ACCOUNT_TYPE_SWAP:
		api := s.apiOrderCancel(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		return s.handleOrderFromOrderCancel(req, res)
	default:
		return nil, ErrorAccountType
	}
}

func (s *SunxTradeEngine) CreateOrders(reqs []*OrderParam) ([]*Order, error) {
	if err := s.restBatchPreCheck(reqs); err != nil {
		return nil, err
	}
	api := s.apiBatchOrderCreate(reqs)
	res, err := api.Do()
	if err != nil {
		return nil, err
	}
	return s.handleOrdersFromBatchOrderCreate(reqs, res)
}

func (s *SunxTradeEngine) AmendOrders(reqs []*OrderParam) ([]*Order, error) {
	return nil, ErrorNotSupport
}

func (s *SunxTradeEngine) CancelOrders(reqs []*OrderParam) ([]*Order, error) {
	if err := s.restBatchPreCheck(reqs); err != nil {
		return nil, err
	}
	api := s.apiBatchOrderCancel(reqs)
	res, err := api.Do()
	if err != nil {
		return nil, err
	}
	return s.handleOrdersFromBatchOrderCancel(reqs, res)
}

func (s *SunxTradeEngine) NewSubscribeOrderReq() *SubscribeOrderParam {
	return &SubscribeOrderParam{}
}

func (s *SunxTradeEngine) SubscribeOrder(r *SubscribeOrderParam) (TradeSubscribe[Order], error) {
	req := *r
	// 构建一个推送订单数据的中转订阅
	newSub := &subscription[Order]{
		resultChan: make(chan Order, 100),
		errChan:    make(chan error, 10),
		closeChan:  make(chan struct{}, 10),
	}
	switch SunxAccountType(req.AccountType) {
	case SUNX_ACCOUNT_TYPE_SWAP:
		err := s.checkWsForSwapOrder()
		if err != nil {
			return nil, err
		}
		swapSub, err := s.wsForSwapOrder.SubscribeOrders(req.ContractCodes, true)
		if err != nil {
			return nil, err
		}
		s.handleSubscribeOrderFromSwapSub(&req, swapSub, newSub)
		return newSub, nil
	default:
		return nil, ErrorAccountType
	}
}

func (s *SunxTradeEngine) WsCreateOrder(req *OrderParam) (*Order, error) {
	return nil, ErrorNotSupport
}

func (s *SunxTradeEngine) WsAmendOrder(req *OrderParam) (*Order, error) {
	return nil, ErrorNotSupport
}

func (s *SunxTradeEngine) WsCancelOrder(req *OrderParam) (*Order, error) {
	return nil, ErrorNotSupport
}

func (s *SunxTradeEngine) WsCreateOrders(reqs []*OrderParam) ([]*Order, error) {
	return nil, ErrorNotSupport
}

func (s *SunxTradeEngine) WsAmendOrders(reqs []*OrderParam) ([]*Order, error) {
	return nil, ErrorNotSupport
}

func (s *SunxTradeEngine) WsCancelOrders(reqs []*OrderParam) ([]*Order, error) {
	return nil, ErrorNotSupport
}
