package mytrade

import (
	"strconv"
	"sync"

	"github.com/Hongssd/mygateapi"
)

type GateTradeEngine struct {
	ExchangeBase

	gateConverter GateEnumConverter
	apiKey        string
	secretKey     string
	passphrase    string

	wsForSpotOrder               *mygateapi.SpotWsStreamClient
	wsForUSDTFuturesOrder        *mygateapi.FuturesWsStreamClient
	wsForBTCFuturesOrder         *mygateapi.FuturesWsStreamClient
	wsForUSDTDeliveryOrder       *mygateapi.DeliveryWsStreamClient
	wsForBTCFuturesDeliveryOrder *mygateapi.DeliveryWsStreamClient
}

func (g *GateTradeEngine) NewOrderReq() *OrderParam {
	return &OrderParam{}
}
func (g *GateTradeEngine) NewQueryOrderReq() *QueryOrderParam {
	return &QueryOrderParam{}
}
func (g *GateTradeEngine) NewQueryTradeReq() *QueryTradeParam {
	return &QueryTradeParam{}
}

func (g *GateTradeEngine) QueryOpenOrders(req *QueryOrderParam) ([]*Order, error) {
	switch GateAccountType(req.AccountType) {
	case GATE_ACCOUNT_TYPE_SPOT:
		if req.IsAlgo {
			api := g.apiSpotPriceOpenOrders(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return g.handleOrdersFromSpotPriceOpenOrders(req, res), nil
		} else {
			api := g.apiSpotOpenOrders(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return g.handleOrdersFromSpotOpenOrders(req, res), nil
		}
	case GATE_ACCOUNT_TYPE_FUTURES:
		if req.IsAlgo {
			api := g.apiFuturesPriceOpenOrders(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return g.handleOrdersFromFuturesPriceOpenOrders(req, res), nil
		} else {
			api := g.apiFuturesOpenOrders(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return g.handleOrdersFromFuturesOpenOrders(req, res), nil
		}
	case GATE_ACCOUNT_TYPE_DELIVERY:
		if req.IsAlgo {
			api := g.apiDeliveryPriceOpenOrders(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return g.handleOrdersFromDeliveryPriceOpenOrders(req, res), nil
		} else {
			api := g.apiDeliveryOpenOrders(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return g.handleOrdersFromDeliveryOpenOrders(req, res), nil
		}

	default:
		return nil, ErrorAccountType
	}
}
func (g *GateTradeEngine) QueryOrder(req *QueryOrderParam) (*Order, error) {
	switch GateAccountType(req.AccountType) {
	case GATE_ACCOUNT_TYPE_SPOT:
		if req.IsAlgo {
			api := g.apiSpotPriceOrderQuery(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return g.handleOrderFromSpotPriceOrderQuery(req, res), nil
		} else {
			api := g.apiSpotOrderQuery(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return g.handleOrderFromSpotOrderQuery(req, res), nil
		}
	case GATE_ACCOUNT_TYPE_FUTURES:
		if req.IsAlgo {
			api := g.apiFuturesPriceOrderQuery(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return g.handleOrderFromFuturesPriceOrderQuery(req, res), nil
		} else {
			api := g.apiFuturesOrderQuery(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return g.handleOrderFromFuturesOrderQuery(req, res), nil
		}
	case GATE_ACCOUNT_TYPE_DELIVERY:
		if req.IsAlgo {
			api := g.apiDeliveryPriceOrderQuery(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return g.handleOrderFromDeliveryPriceOrderQuery(req, res), nil
		} else {
			api := g.apiDeliveryOrderQuery(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return g.handleOrderFromDeliveryOrderQuery(req, res), nil
		}
	default:
		return nil, ErrorAccountType
	}
}
func (g *GateTradeEngine) QueryOrders(req *QueryOrderParam) ([]*Order, error) {
	switch GateAccountType(req.AccountType) {
	case GATE_ACCOUNT_TYPE_SPOT:
		if req.IsAlgo {
			api := g.apiSpotPriceOrdersQuery(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return g.handleOrdersFromSpotPriceOrdersQuery(req, res), nil
		} else {
			api := g.apiSpotOrdersQuery(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return g.handleOrdersFromSpotOrdersQuery(req, res), nil
		}
	case GATE_ACCOUNT_TYPE_FUTURES:
		if req.IsAlgo {
			api := g.apiFuturesPriceOrdersQuery(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return g.handleOrdersFromFuturesPriceOrdersQuery(req, res), nil
		} else {
			api := g.apiFuturesOrdersQuery(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return g.handleOrdersFromFuturesOrdersQuery(req, res), nil
		}

	case GATE_ACCOUNT_TYPE_DELIVERY:
		if req.IsAlgo {
			api := g.apiDeliveryPriceOrdersQuery(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return g.handleOrdersFromDeliveryPriceOrdersQuery(req, res), nil
		} else {
			api := g.apiDeliveryOrdersQuery(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return g.handleOrdersFromDeliveryOrdersQuery(req, res), nil
		}
	default:
		return nil, ErrorAccountType
	}
}

func (g *GateTradeEngine) QueryTrades(req *QueryTradeParam) ([]*Trade, error) {
	switch GateAccountType(req.AccountType) {
	case GATE_ACCOUNT_TYPE_SPOT:
		api := g.apiSpotTradesQuery(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		return g.handleTradesFromSpotTradesQuery(req, res), nil
	case GATE_ACCOUNT_TYPE_FUTURES:
		api := g.apiFuturesTradesQuery(req)
		if api == nil {
			return nil, ErrorSymbolNotFound
		}
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		return g.handleTradesFromFuturesTradesQuery(req, res), nil
	case GATE_ACCOUNT_TYPE_DELIVERY:
		api := g.apiDeliveryTradesQuery(req)
		if api == nil {
			return nil, ErrorSymbolNotFound
		}
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		return g.handleTradesFromDeliveryTradesQuery(req, res), nil
	default:
		return nil, ErrorAccountType
	}
}

func (g *GateTradeEngine) CreateOrder(req *OrderParam) (*Order, error) {
	switch GateAccountType(req.AccountType) {
	case GATE_ACCOUNT_TYPE_SPOT:
		if req.IsAlgo {
			api := g.apiSpotPriceOrderCreate(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return g.handleOrderFromSpotPriceOrderCreate(req, res), nil
		} else {
			api := g.apiSpotOrderCreate(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return g.handleOrderFromSpotOrderCreate(req, res), nil
		}
	case GATE_ACCOUNT_TYPE_FUTURES:
		if req.IsAlgo {
			api := g.apiFuturesPriceOrderCreate(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return g.handleOrderFromFuturesPriceOrderCreate(req, res), nil
		} else {
			api := g.apiFuturesOrderCreate(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return g.handleOrderFromFuturesOrderCreate(req, res), nil
		}
	case GATE_ACCOUNT_TYPE_DELIVERY:
		if req.IsAlgo {
			api := g.apiDeliveryPriceOrderCreate(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return g.handleOrderFromDeliveryPriceOrderCreate(req, res), nil
		} else {
			api := g.apiDeliveryOrderCreate(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return g.handleOrderFromDeliveryOrderCreate(req, res), nil
		}
	default:
		return nil, ErrorAccountType
	}
}
func (g *GateTradeEngine) AmendOrder(req *OrderParam) (*Order, error) {
	switch GateAccountType(req.AccountType) {
	case GATE_ACCOUNT_TYPE_SPOT:
		api := g.apiSpotOrderAmend(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		return g.handleOrderFromSpotOrderAmend(req, res), nil
	case GATE_ACCOUNT_TYPE_FUTURES:
		api := g.apiFuturesOrderAmend(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		return g.handleOrderFromFuturesOrderAmend(req, res), nil
	case GATE_ACCOUNT_TYPE_DELIVERY:
		return nil, ErrorNotSupport
	default:
		return nil, ErrorAccountType
	}
}
func (g *GateTradeEngine) CancelOrder(req *OrderParam) (*Order, error) {
	var order *Order
	switch GateAccountType(req.AccountType) {
	case GATE_ACCOUNT_TYPE_SPOT:
		if req.IsAlgo {
			api := g.apiSpotPriceOrderCancel(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return g.handleOrderFromSpotPriceOrderCancel(req, res), nil
		} else {
			api := g.apiSpotOrderCancel(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			order = g.handleOrderFromSpotOrderCancel(req, res)
		}
	case GATE_ACCOUNT_TYPE_FUTURES:
		if req.IsAlgo {
			api := g.apiFuturesPriceOrderCancel(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return g.handleOrderFromFuturesPriceOrderCancel(req, res), nil
		} else {
			api := g.apiFuturesOrderCancel(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			order = g.handleOrderFromFuturesOrderCancel(req, res)
		}
	case GATE_ACCOUNT_TYPE_DELIVERY:
		if req.IsAlgo {
			api := g.apiDeliveryPriceOrderCancel(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return g.handleOrderFromDeliveryPriceOrderCancel(req, res), nil
		} else {
			api := g.apiDeliveryOrderCancel(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			order = g.handleOrderFromDeliveryOrderCancel(req, res)
		}
	default:
		return nil, ErrorAccountType
	}
	return order, nil
}

func (g *GateTradeEngine) CreateOrders(reqs []*OrderParam) ([]*Order, error) {
	//使用并发下单
	orders := []*Order{}
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, req := range reqs {
		req := req
		wg.Add(1)
		go func() {
			defer wg.Done()
			order, err := g.CreateOrder(req)
			if err != nil {
				mu.Lock()
				orders = append(orders, g.handleOrderFromBatchErr(req, err))
				mu.Unlock()
			}
			mu.Lock()
			orders = append(orders, order)
			mu.Unlock()
		}()
	}
	wg.Wait()
	return orders, nil
}
func (g *GateTradeEngine) AmendOrders(reqs []*OrderParam) ([]*Order, error) {
	//使用并发改单
	orders := []*Order{}
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, req := range reqs {
		req := req
		wg.Add(1)
		go func() {
			defer wg.Done()
			order, err := g.AmendOrder(req)
			if err != nil {
				mu.Lock()
				orders = append(orders, g.handleOrderFromBatchErr(req, err))
				mu.Unlock()
			}
			mu.Lock()
			orders = append(orders, order)
			mu.Unlock()
		}()
	}
	wg.Wait()
	return orders, nil
}
func (g *GateTradeEngine) CancelOrders(reqs []*OrderParam) ([]*Order, error) {
	//使用并发撤单
	orders := []*Order{}
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, req := range reqs {
		req := req
		wg.Add(1)
		go func() {
			defer wg.Done()
			order, err := g.CancelOrder(req)
			if err != nil {
				mu.Lock()
				orders = append(orders, g.handleOrderFromBatchErr(req, err))
				mu.Unlock()
			}
			mu.Lock()
			orders = append(orders, order)
			mu.Unlock()
		}()
	}
	wg.Wait()
	return orders, nil
}

func (g *GateTradeEngine) NewSubscribeOrderReq() *SubscribeOrderParam {
	return &SubscribeOrderParam{}
}

func (g *GateTradeEngine) SubscribeOrder(r *SubscribeOrderParam) (TradeSubscribe[Order], error) {
	req := *r
	//构建一个推送订单数据的中转订阅
	newSub := &subscription[Order]{
		resultChan: make(chan Order, 100),
		errChan:    make(chan error, 10),
		closeChan:  make(chan struct{}, 10),
	}
	switch GateAccountType(req.AccountType) {
	case GATE_ACCOUNT_TYPE_SPOT:

		err := g.checkWsForSpotOrder()
		if err != nil {
			return nil, err
		}
		spotSub, err := g.wsForSpotOrder.SubscribeOrders()
		if err != nil {
			return nil, err
		}
		g.handleSubscribeOrderFromSpotSub(req, spotSub, newSub)
		return newSub, nil
	case GATE_ACCOUNT_TYPE_FUTURES:

		err := g.checkWsForFuturesOrder()
		if err != nil {
			return nil, err
		}
		accountDetail, err := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestAccountDetail().Do()
		if err != nil {
			return nil, err
		}
		userId := strconv.FormatInt(accountDetail.Data.UserId, 10)

		futuresUsdtSub, err := g.wsForUSDTFuturesOrder.SubscribeOrder(userId, string(mygateapi.USDT_CONTRACT))
		if err != nil {
			return nil, err
		}
		futuresBtcSub, err := g.wsForBTCFuturesOrder.SubscribeOrder(userId, string(mygateapi.BTC_CONTRACT))
		if err != nil {
			return nil, err
		}
		g.handleSubscribeOrderFromFuturesOrDeliverySub(req, futuresUsdtSub, newSub)
		g.handleSubscribeOrderFromFuturesOrDeliverySub(req, futuresBtcSub, newSub)
		return newSub, nil
	case GATE_ACCOUNT_TYPE_DELIVERY:
		err := g.checkWsForDeliveryOrder()
		if err != nil {
			return nil, err
		}
		accountDetail, err := mygateapi.NewRestClient(g.apiKey, g.secretKey).PrivateRestClient().NewPrivateRestAccountDetail().Do()
		if err != nil {
			return nil, err
		}
		userId := strconv.FormatInt(accountDetail.Data.UserId, 10)
		deliveryUsdtSub, err := g.wsForUSDTDeliveryOrder.SubscribeOrder(userId, string(mygateapi.USDT_CONTRACT))
		if err != nil {
			return nil, err
		}
		deliveryBtcSub, err := g.wsForBTCFuturesDeliveryOrder.SubscribeOrder(userId, string(mygateapi.BTC_CONTRACT))
		if err != nil {
			return nil, err
		}
		g.handleSubscribeOrderFromFuturesOrDeliverySub(req, deliveryUsdtSub, newSub)
		g.handleSubscribeOrderFromFuturesOrDeliverySub(req, deliveryBtcSub, newSub)
		return newSub, nil
	default:
		return nil, ErrorAccountType
	}
}

func (g *GateTradeEngine) WsCreateOrder(req *OrderParam) (*Order, error) {
	return nil, ErrorNotSupport
}
func (g *GateTradeEngine) WsAmendOrder(req *OrderParam) (*Order, error) {
	return nil, ErrorNotSupport
}
func (g *GateTradeEngine) WsCancelOrder(req *OrderParam) (*Order, error) {
	return nil, ErrorNotSupport
}

func (g *GateTradeEngine) WsCreateOrders(reqs []*OrderParam) ([]*Order, error) {
	return nil, ErrorNotSupport
}
func (g *GateTradeEngine) WsAmendOrders(reqs []*OrderParam) ([]*Order, error) {
	return nil, ErrorNotSupport
}
func (g *GateTradeEngine) WsCancelOrders(reqs []*OrderParam) ([]*Order, error) {
	return nil, ErrorNotSupport
}
