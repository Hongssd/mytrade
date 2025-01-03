package mytrade

import (
	"github.com/Hongssd/mygateapi"
	"sync"
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

func (o *GateTradeEngine) NewOrderReq() *OrderParam {
	return &OrderParam{}
}
func (o *GateTradeEngine) NewQueryOrderReq() *QueryOrderParam {
	return &QueryOrderParam{}
}
func (o *GateTradeEngine) NewQueryTradeReq() *QueryTradeParam {
	return &QueryTradeParam{}
}

func (o *GateTradeEngine) QueryOpenOrders(req *QueryOrderParam) ([]*Order, error) {

	return nil, nil
}
func (o *GateTradeEngine) QueryOrder(req *QueryOrderParam) (*Order, error) {
	return nil, nil
}
func (o *GateTradeEngine) QueryOrders(req *QueryOrderParam) ([]*Order, error) {
	return nil, nil
}

func (o *GateTradeEngine) QueryTrades(req *QueryTradeParam) ([]*Trade, error) {
	return nil, nil
}

func (o *GateTradeEngine) CreateOrder(req *OrderParam) (*Order, error) {
	return nil, nil

}
func (o *GateTradeEngine) AmendOrder(req *OrderParam) (*Order, error) {
	return nil, nil
}
func (o *GateTradeEngine) CancelOrder(req *OrderParam) (*Order, error) {
	return nil, nil
}

func (o *GateTradeEngine) CreateOrders(reqs []*OrderParam) ([]*Order, error) {
	return nil, nil
}
func (o *GateTradeEngine) AmendOrders(reqs []*OrderParam) ([]*Order, error) {
	return nil, nil
}
func (o *GateTradeEngine) CancelOrders(reqs []*OrderParam) ([]*Order, error) {
	return nil, nil
}

func (o *GateTradeEngine) NewSubscribeOrderReq() *SubscribeOrderParam {
	return &SubscribeOrderParam{}
}
func (o *GateTradeEngine) SubscribeOrder(req *SubscribeOrderParam) (TradeSubscribe[Order], error) {
	return nil, nil
}

func (o *GateTradeEngine) WsCreateOrder(req *OrderParam) (*Order, error) {
	return nil, nil
}
func (o *GateTradeEngine) WsAmendOrder(req *OrderParam) (*Order, error) {
	return nil, nil
}
func (o *GateTradeEngine) WsCancelOrder(req *OrderParam) (*Order, error) {
	return nil, nil
}

func (o *GateTradeEngine) WsCreateOrders(reqs []*OrderParam) ([]*Order, error) {
	return nil, nil
}
func (o *GateTradeEngine) WsAmendOrders(reqs []*OrderParam) ([]*Order, error) {
	return nil, nil
}
func (o *GateTradeEngine) WsCancelOrders(reqs []*OrderParam) ([]*Order, error) {
	return nil, nil
}
