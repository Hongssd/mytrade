package mytrade

import (
	"sync"
	"time"

	"github.com/Hongssd/mysunxapi"
)

type SunxTradeEngine struct {
	ExchangeBase

	sunxConverter SunxEnumConverter
	accessKey     string
	secretKey     string

	wsForSwapOrder *mysunxapi.PrivateWsStreamClient

	// 广播器管理 - 按交易对管理广播器
	broadcasters   *MySyncMap[string, *sunxOrderBroadcaster]
	broadcastersMu sync.RWMutex
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
		o, err := s.handleOrderFromQueryOrder(req, res)
		log.Infof("query order: %+v", o)
		return o, nil
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
		// 获取或创建该交易对的广播器（同一交易对只会创建一次）
		broadcaster, err := s.getBroadcasterFromSymbol(req.Symbol)
		if err != nil {
			return nil, err
		}

		// 创建订阅者
		sub, err := s.newOrderSubscriber(broadcaster, req.ClientOrderId, "", req.Symbol)
		if err != nil {
			return nil, err
		}
		defer s.closeSubscribe(broadcaster, sub)

		// 发起 API 请求
		api := s.apiOrderCreate(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}

		//处理API返回值
		_, err = s.handleOrderFromOrderCreate(req, res)
		if err != nil {
			return nil, err
		}

		// 等待 WS 推送返回完整订单信息
		return s.waitSubscribeReturn(sub, 10*time.Second)
	default:
		return nil, ErrorAccountType
	}
}

func (s *SunxTradeEngine) AmendOrder(req *OrderParam) (*Order, error) {
	return nil, ErrorNotSupport
}

func (s *SunxTradeEngine) CancelOrder(req *OrderParam) (*Order, error) {
	if err := s.accountTypePreCheck(req.AccountType); err != nil {
		return nil, err
	}
	switch SunxAccountType(req.AccountType) {
	case SUNX_ACCOUNT_TYPE_SWAP:
		// 获取或创建该交易对的广播器（同一交易对只会创建一次）
		broadcaster, err := s.getBroadcasterFromSymbol(req.Symbol)
		if err != nil {
			return nil, err
		}

		// 创建订阅者
		sub, err := s.newOrderSubscriber(broadcaster, req.ClientOrderId, req.OrderId, req.Symbol)
		if err != nil {
			return nil, err
		}
		defer s.closeSubscribe(broadcaster, sub)

		// 发起 API 请求
		api := s.apiOrderCancel(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}

		//处理API返回值
		_, err = s.handleOrderFromOrderCancel(req, res)
		if err != nil {
			return nil, err
		}

		// 等待 WS 推送返回完整订单信息
		return s.waitSubscribeReturn(sub, 10*time.Second)
	default:
		return nil, ErrorAccountType
	}
}

func (s *SunxTradeEngine) CreateOrders(reqs []*OrderParam) ([]*Order, error) {
	if err := s.restBatchPreCheck(reqs); err != nil {
		return nil, err
	}

	// 按交易对分组创建订阅者（同一交易对只会创建一次广播器）
	symbolBroadcasters := make(map[string]*sunxOrderBroadcaster)
	subscribers := make([]*sunxOrderSubscriber, 0, len(reqs))

	for _, req := range reqs {
		// 获取或创建该交易对的广播器
		var broadcaster *sunxOrderBroadcaster
		if sb, ok := symbolBroadcasters[req.Symbol]; ok {
			broadcaster = sb
		} else {
			var err error
			broadcaster, err = s.getBroadcasterFromSymbol(req.Symbol)
			if err != nil {
				return nil, err
			}
			symbolBroadcasters[req.Symbol] = broadcaster
		}

		// 创建订阅者
		sub, err := s.newOrderSubscriber(broadcaster, req.ClientOrderId, "", req.Symbol)
		if err != nil {
			return nil, err
		}
		subscribers = append(subscribers, sub)
	}

	// 确保在返回前关闭所有订阅者
	defer func() {
		for i, sub := range subscribers {
			broadcaster := symbolBroadcasters[reqs[i].Symbol]
			s.closeSubscribe(broadcaster, sub)
		}
	}()

	// 发起批量下单请求
	api := s.apiBatchOrderCreate(reqs)
	res, err := api.Do()
	if err != nil {
		return nil, err
	}

	// 等待所有订单的 WS 推送
	orders := make([]*Order, 0, len(reqs))
	for _, sub := range subscribers {
		order, err := s.waitSubscribeReturn(sub, 10*time.Second)
		if err != nil {
			log.Warnf("wait order ws return failed: %v", err)
			continue
		}
		orders = append(orders, order)
	}

	// 如果没有收到任何 WS 推送，使用 API 响应
	if len(orders) == 0 {
		return s.handleOrdersFromBatchOrderCreate(reqs, res)
	}

	return orders, nil
}

func (s *SunxTradeEngine) AmendOrders(reqs []*OrderParam) ([]*Order, error) {
	return nil, ErrorNotSupport
}

func (s *SunxTradeEngine) CancelOrders(reqs []*OrderParam) ([]*Order, error) {
	if err := s.restBatchPreCheck(reqs); err != nil {
		return nil, err
	}

	// 按交易对分组创建订阅者（同一交易对只会创建一次广播器）
	symbolBroadcasters := make(map[string]*sunxOrderBroadcaster)
	subscribers := make([]*sunxOrderSubscriber, 0, len(reqs))

	for _, req := range reqs {
		// 获取或创建该交易对的广播器
		var broadcaster *sunxOrderBroadcaster
		if sb, ok := symbolBroadcasters[req.Symbol]; ok {
			broadcaster = sb
		} else {
			var err error
			broadcaster, err = s.getBroadcasterFromSymbol(req.Symbol)
			if err != nil {
				return nil, err
			}
			symbolBroadcasters[req.Symbol] = broadcaster
		}

		// 创建订阅者
		sub, err := s.newOrderSubscriber(broadcaster, req.ClientOrderId, req.OrderId, req.Symbol)
		if err != nil {
			return nil, err
		}
		subscribers = append(subscribers, sub)
	}

	// 确保在返回前关闭所有订阅者
	defer func() {
		for i, sub := range subscribers {
			broadcaster := symbolBroadcasters[reqs[i].Symbol]
			s.closeSubscribe(broadcaster, sub)
		}
	}()

	// 发起批量撤单请求
	api := s.apiBatchOrderCancel(reqs)
	res, err := api.Do()
	if err != nil {
		return nil, err
	}

	// 等待所有订单的 WS 推送
	orders := make([]*Order, 0, len(reqs))
	for _, sub := range subscribers {
		order, err := s.waitSubscribeReturn(sub, 10*time.Second)
		if err != nil {
			log.Warnf("wait order ws return failed: %v", err)
			continue
		}
		orders = append(orders, order)
	}

	// 如果没有收到任何 WS 推送，使用 API 响应
	if len(orders) == 0 {
		return s.handleOrdersFromBatchOrderCancel(reqs, res)
	}

	return orders, nil
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

// 获取指定交易对的广播器，如果不存在则创建
func (s *SunxTradeEngine) getBroadcasterFromSymbol(symbol string) (*sunxOrderBroadcaster, error) {
	s.broadcastersMu.Lock()
	defer s.broadcastersMu.Unlock()

	// 初始化 broadcasters map
	if s.broadcasters == nil {
		s.broadcasters = GetPointer(NewMySyncMap[string, *sunxOrderBroadcaster]())
	}

	// 尝试获取已存在的广播器
	if broadcaster, ok := s.broadcasters.Load(symbol); ok {
		return broadcaster, nil
	}

	// 创建新的广播器
	newBroadcaster, err := s.newOrderBroadcaster(symbol)
	if err != nil {
		return nil, err
	}

	// 存储广播器
	s.broadcasters.Store(symbol, newBroadcaster)
	log.Infof("创建新的订单广播器，交易对: %s", symbol)

	return newBroadcaster, nil
}
