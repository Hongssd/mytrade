package mytrade

import (
	"context"
	"sync"

	"github.com/Hongssd/mygateapi"
	"golang.org/x/sync/errgroup"
)

type GateTradeEngine struct {
	ExchangeBase

	gateConverter GateEnumConverter
	apiKey        string
	secretKey     string
	passphrase    string

	wsForSpotOrder       *mygateapi.SpotWsStreamClient
	wsForSpotOrderMu     sync.Mutex
	wsForFuturesOrder    *mygateapi.FuturesWsStreamClient
	wsForFuturesOrderMu  sync.Mutex
	wsForDeliveryOrder   *mygateapi.DeliveryWsStreamClient
	wsForDeliveryOrderMu sync.Mutex
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
		}
		api := g.apiSpotOpenOrders(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		return g.handleOrdersFromSpotOpenOrders(req, res), nil
	case GATE_ACCOUNT_TYPE_FUTURES:
		if req.IsAlgo {
			api := g.apiFuturesPriceOpenOrders(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return g.handleOrdersFromFuturesPriceOpenOrders(req, res), nil
		}
		api := g.apiFuturesOpenOrders(req)
		if api == nil {
			return nil, ErrorSymbolNotFound
		}
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		return g.handleOrdersFromFuturesOpenOrders(req, res), nil
	case GATE_ACCOUNT_TYPE_DELIVERY:
		if req.IsAlgo {
			api := g.apiDeliveryPriceOpenOrders(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return g.handleOrdersFromDeliveryPriceOpenOrders(req, res), nil
		}
		api := g.apiDeliveryOpenOrders(req)
		if api == nil {
			return nil, ErrorSymbolNotFound
		}
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		return g.handleOrdersFromDeliveryOpenOrders(req, res), nil
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
		}
		api := g.apiSpotOrderQuery(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		return g.handleOrderFromSpotOrderQuery(req, res), nil
	case GATE_ACCOUNT_TYPE_FUTURES:
		if req.IsAlgo {
			api := g.apiFuturesPriceOrderQuery(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return g.handleOrderFromFuturesPriceOrderQuery(req, res), nil
		}
		api := g.apiFuturesOrderQuery(req)
		if api == nil {
			return nil, ErrorSymbolNotFound
		}
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		return g.handleOrderFromFuturesOrderQuery(req, res), nil
	case GATE_ACCOUNT_TYPE_DELIVERY:
		if req.IsAlgo {
			api := g.apiDeliveryPriceOrderQuery(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return g.handleOrderFromDeliveryPriceOrderQuery(req, res), nil
		}
		api := g.apiDeliveryOrderQuery(req)
		if api == nil {
			return nil, ErrorSymbolNotFound
		}
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		return g.handleOrderFromDeliveryOrderQuery(req, res), nil
	default:
		return nil, ErrorAccountType
	}
}
func (g *GateTradeEngine) QueryOrders(req *QueryOrderParam) ([]*Order, error) {
	switch GateAccountType(req.AccountType) {
	case GATE_ACCOUNT_TYPE_SPOT:
		if req.IsAlgo {
			api := g.apiSpotPriceOrdersQuery(req)

			var orders []*Order
			var statuses []string = []string{GATE_ORDER_CONTRACT_STATUS_OPEN, GATE_ORDER_CONTRACT_STATUS_FINISHED}
			// errGroup
			errG, _ := errgroup.WithContext(context.Background())
			for _, status := range statuses {
				status := status
				errG.Go(func() error {
					api.Status(status)
					res, err := api.Do()
					if err != nil {
						return err
					}
					orders = append(orders, g.handleOrdersFromSpotPriceOrdersQuery(req, res)...)
					return nil
				})
			}
			err := errG.Wait()
			if err != nil {
				return nil, err
			}
		}
		api := g.apiSpotOrdersQuery(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		return g.handleOrdersFromSpotOrdersQuery(req, res), nil
	case GATE_ACCOUNT_TYPE_FUTURES:
		if req.IsAlgo {
			api := g.apiFuturesPriceOrdersQuery(req)

			var orders []*Order
			var statuses []string = []string{GATE_ORDER_CONTRACT_STATUS_OPEN, GATE_ORDER_CONTRACT_STATUS_FINISHED}
			// errGroup
			errG, _ := errgroup.WithContext(context.Background())
			for _, status := range statuses {
				status := status
				errG.Go(func() error {
					api.Status(status)
					res, err := api.Do()
					if err != nil {
						return err
					}
					orders = append(orders, g.handleOrdersFromFuturesPriceOrdersQuery(req, res)...)
					return nil
				})
			}

			err := errG.Wait()
			if err != nil {
				return nil, err
			}

			return orders, nil
		}
		api := g.apiFuturesOrdersQuery(req)
		if api == nil {
			return nil, ErrorSymbolNotFound
		}
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		return g.handleOrdersFromFuturesOrdersQuery(req, res), nil
	case GATE_ACCOUNT_TYPE_DELIVERY:
		if req.IsAlgo {
			api := g.apiDeliveryPriceOrdersQuery(req)

			var orders []*Order
			var statuses []string = []string{GATE_ORDER_CONTRACT_STATUS_OPEN, GATE_ORDER_CONTRACT_STATUS_FINISHED}
			// errGroup
			errG, _ := errgroup.WithContext(context.Background())
			for _, status := range statuses {
				status := status
				errG.Go(func() error {
					api.Status(status)
					res, err := api.Do()
					if err != nil {
						return err
					}
					orders = append(orders, g.handleOrdersFromDeliveryPriceOrdersQuery(req, res)...)
					return nil
				})
			}

			err := errG.Wait()
			if err != nil {
				return nil, err
			}

			return orders, nil
		}
		api := g.apiDeliveryOrdersQuery(req)
		if api == nil {
			return nil, ErrorSymbolNotFound
		}
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		return g.handleOrdersFromDeliveryOrdersQuery(req, res), nil
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
		}
		api := g.apiSpotOrderCreate(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		return g.handleOrderFromSpotOrderCreate(req, res), nil
	case GATE_ACCOUNT_TYPE_FUTURES:
		if req.IsAlgo {
			api := g.apiFuturesPriceOrderCreate(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return g.handleOrderFromFuturesPriceOrderCreate(req, res), nil
		}
		api := g.apiFuturesOrderCreate(req)
		if api == nil {
			return nil, ErrorSymbolNotFound
		}
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		return g.handleOrderFromFuturesOrderCreate(req, res), nil
	case GATE_ACCOUNT_TYPE_DELIVERY:
		if req.IsAlgo {
			api := g.apiDeliveryPriceOrderCreate(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return g.handleOrderFromDeliveryPriceOrderCreate(req, res), nil
		}
		api := g.apiDeliveryOrderCreate(req)
		if api == nil {
			return nil, ErrorSymbolNotFound
		}
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		return g.handleOrderFromDeliveryOrderCreate(req, res), nil
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
		if api == nil {
			return nil, ErrorSymbolNotFound
		}
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
		}
		api := g.apiSpotOrderCancel(req)
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		order = g.handleOrderFromSpotOrderCancel(req, res)
	case GATE_ACCOUNT_TYPE_FUTURES:
		if req.IsAlgo {
			api := g.apiFuturesPriceOrderCancel(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return g.handleOrderFromFuturesPriceOrderCancel(req, res), nil
		}
		api := g.apiFuturesOrderCancel(req)
		if api == nil {
			return nil, ErrorSymbolNotFound
		}
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		order = g.handleOrderFromFuturesOrderCancel(req, res)
	case GATE_ACCOUNT_TYPE_DELIVERY:
		if req.IsAlgo {
			api := g.apiDeliveryPriceOrderCancel(req)
			res, err := api.Do()
			if err != nil {
				return nil, err
			}
			return g.handleOrderFromDeliveryPriceOrderCancel(req, res), nil
		}
		api := g.apiDeliveryOrderCancel(req)
		if api == nil {
			return nil, ErrorSymbolNotFound
		}
		res, err := api.Do()
		if err != nil {
			return nil, err
		}
		order = g.handleOrderFromDeliveryOrderCancel(req, res)
	default:
		return nil, ErrorAccountType
	}
	return order, nil
}

func (g *GateTradeEngine) CreateOrders(reqs []*OrderParam) ([]*Order, error) {
	return nil, nil
}
func (g *GateTradeEngine) AmendOrders(reqs []*OrderParam) ([]*Order, error) {
	return nil, nil
}
func (g *GateTradeEngine) CancelOrders(reqs []*OrderParam) ([]*Order, error) {
	return nil, nil
}

func (g *GateTradeEngine) NewSubscribeOrderReq() *SubscribeOrderParam {
	return &SubscribeOrderParam{}
}

func (g *GateTradeEngine) SubscribeOrder(req *SubscribeOrderParam) (TradeSubscribe[Order], error) {
	return nil, nil
}

func (g *GateTradeEngine) WsCreateOrder(req *OrderParam) (*Order, error) {
	return nil, nil
}
func (g *GateTradeEngine) WsAmendOrder(req *OrderParam) (*Order, error) {
	return nil, nil
}
func (g *GateTradeEngine) WsCancelOrder(req *OrderParam) (*Order, error) {
	return nil, nil
}

func (g *GateTradeEngine) WsCreateOrders(reqs []*OrderParam) ([]*Order, error) {
	return nil, nil
}
func (g *GateTradeEngine) WsAmendOrders(reqs []*OrderParam) ([]*Order, error) {
	return nil, nil
}
func (g *GateTradeEngine) WsCancelOrders(reqs []*OrderParam) ([]*Order, error) {
	return nil, nil
}
