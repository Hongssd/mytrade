package mytrade

import (
	"errors"
	"fmt"
	"time"
)

type ExchangeManager struct {
	//交易所管理器
	ExchangeMap     *MySyncMap[string, TradeExchange]
	ExchangeInfoMap *MySyncMap[string, TradeExchangeInfo]
	MarketDataMap   *MySyncMap[string, TradeMarketData]
	TradeEngineMap  *MySyncMap[string, TradeEngine]
}

type ExchangeApiParam struct {
	Exchange   string
	ApiKey     string
	ApiSecret  string
	Passphrase string
}

// 创建交易所管理器
func NewExchangeManager() *ExchangeManager {
	e := &ExchangeManager{
		ExchangeMap:     GetPointer(NewMySyncMap[string, TradeExchange]()),
		ExchangeInfoMap: GetPointer(NewMySyncMap[string, TradeExchangeInfo]()),
		MarketDataMap:   GetPointer(NewMySyncMap[string, TradeMarketData]()),
		TradeEngineMap:  GetPointer(NewMySyncMap[string, TradeEngine]()),
	}
	e.ExchangeMap.Store(BINANCE_NAME.String(), NewBinanceExchange())
	e.ExchangeMap.Store(OKX_NAME.String(), NewOkxExchange())
	e.ExchangeMap.Store(BYBIT_NAME.String(), NewBybitExchange())
	return e
}

// 获取交易所规范
func (e *ExchangeManager) getTradeExchangeInfo(exchange string) (TradeExchangeInfo, error) {
	tradeExchange, ok := e.ExchangeMap.Load(exchange)
	if !ok {
		return nil, errors.Join(errors.New(exchange), ErrorNotSupport)
	}
	tradeExchangeInfo, ok := e.ExchangeInfoMap.Load(exchange)
	if !ok {
		tradeExchangeInfo = tradeExchange.NewExchangeInfo()
		e.ExchangeInfoMap.Store(exchange, tradeExchangeInfo)
	}
	return tradeExchangeInfo, nil
}

// 获取交易所市场数据
func (e *ExchangeManager) getTradeMarketData(exchange string) (TradeMarketData, error) {
	tradeExchange, ok := e.ExchangeMap.Load(exchange)
	if !ok {
		return nil, errors.Join(errors.New(exchange), ErrorNotSupport)
	}
	marketData, ok := e.MarketDataMap.Load(exchange)
	if !ok {
		marketData = tradeExchange.NewMarketData()
		e.MarketDataMap.Store(exchange, marketData)
	}
	return marketData, nil
}

// 获取交易引擎
func (e *ExchangeManager) getTradeEngine(exchange string, apikey, apiSecret, passphrase string) (TradeEngine, error) {
	tradeExchange, ok := e.ExchangeMap.Load(exchange)
	if !ok {
		return nil, errors.Join(errors.New(exchange), ErrorNotSupport)
	}

	key := fmt.Sprintf("%s_%s", exchange, apikey)

	tradeEngine, ok := e.TradeEngineMap.Load(key)
	if !ok {
		tradeEngine = tradeExchange.NewTradeEngine(apikey, apiSecret, passphrase)
		e.TradeEngineMap.Store(key, tradeEngine)
	}
	return tradeEngine, nil
}

// 刷新交易所规范
func (e *ExchangeManager) RefreshExchangeInfo(exchange string) error {
	tradeExchangeInfo, err := e.getTradeExchangeInfo(exchange)
	if err != nil {
		return err
	}

	//刷新交易规范
	retry := 0
	err = tradeExchangeInfo.Refresh()
	for err != nil && retry < 5 {
		retry++
		log.Errorf("刷新交易规范失败，重试第%d次:%s", retry, err.Error())
		time.Sleep(5 * time.Second)
		err = tradeExchangeInfo.Refresh()
	}
	log.Infof("%s刷新交易规范成功", exchange)
	return nil
}

// 获取交易对规范
func (e *ExchangeManager) GetSymbolInfo(exchange, accountType, symbol string) (TradeSymbolInfo, error) {
	tradeExchangeInfo, err := e.getTradeExchangeInfo(exchange)
	if err != nil {
		return nil, err
	}
	return tradeExchangeInfo.GetSymbolInfo(accountType, symbol)
}

// 获取K线数据
func (e *ExchangeManager) GetKlineData(symbolInfo TradeSymbolInfo, interval string, start, end int64, limit int) (*[]Kline, error) {
	marketData, err := e.getTradeMarketData(symbolInfo.Exchange())
	if err != nil {
		return nil, err
	}
	req := marketData.NewKlineReq().
		SetAccountType(symbolInfo.AccountType()).SetSymbol(symbolInfo.Symbol()).SetInterval(interval)

	if start != 0 {
		req.SetStartTime(start)
	}
	if end != 0 {
		req.SetEndTime(end)
	}
	if limit != 0 {
		req.SetLimit(limit)
	}

	return marketData.GetKline(req)
}

// 获取深度数据
func (e *ExchangeManager) GetBookData(symbolInfo TradeSymbolInfo, level int) (*OrderBook, error) {
	marketData, err := e.getTradeMarketData(symbolInfo.Exchange())
	if err != nil {
		return nil, err
	}
	req := marketData.NewBookReq().
		SetAccountType(symbolInfo.AccountType()).SetSymbol(symbolInfo.Symbol())

	if level != 0 {
		req.SetLevel(level)
	}

	return marketData.GetBook(req)
}

// 查询单个订单
func (e *ExchangeManager) QueryOrder(api ExchangeApiParam, req *QueryOrderParam) (*Order, error) {
	tradeEngine, err := e.getTradeEngine(api.Exchange, api.ApiKey, api.ApiSecret, api.Passphrase)
	if err != nil {
		return nil, err
	}
	return tradeEngine.QueryOrder(req)
}

// 查询最近订单
func (e *ExchangeManager) QueryOrders(api ExchangeApiParam, req *QueryOrderParam) ([]*Order, error) {
	tradeEngine, err := e.getTradeEngine(api.Exchange, api.ApiKey, api.ApiSecret, api.Passphrase)
	if err != nil {
		return nil, err
	}
	return tradeEngine.QueryOrders(req)
}

// 查询挂单
func (e *ExchangeManager) QueryOpenOrders(api ExchangeApiParam, req *QueryOrderParam) ([]*Order, error) {
	tradeEngine, err := e.getTradeEngine(api.Exchange, api.ApiKey, api.ApiSecret, api.Passphrase)
	if err != nil {
		return nil, err
	}
	return tradeEngine.QueryOpenOrders(req)
}

// 查询成交
func (e *ExchangeManager) QueryTrades(api ExchangeApiParam, req *QueryTradeParam) ([]*Trade, error) {
	tradeEngine, err := e.getTradeEngine(api.Exchange, api.ApiKey, api.ApiSecret, api.Passphrase)
	if err != nil {
		return nil, err
	}
	return tradeEngine.QueryTrades(req)
}

// 下单
func (e *ExchangeManager) CreateOrder(api ExchangeApiParam, req *OrderParam) (*Order, error) {
	tradeEngine, err := e.getTradeEngine(api.Exchange, api.ApiKey, api.ApiSecret, api.Passphrase)
	if err != nil {
		return nil, err
	}
	return tradeEngine.CreateOrder(req)
}

// 改单
func (e *ExchangeManager) AmendOrder(api ExchangeApiParam, req *OrderParam) (*Order, error) {
	tradeEngine, err := e.getTradeEngine(api.Exchange, api.ApiKey, api.ApiSecret, api.Passphrase)
	if err != nil {
		return nil, err
	}
	return tradeEngine.AmendOrder(req)
}

// 撤单
func (e *ExchangeManager) CancelOrder(api ExchangeApiParam, req *OrderParam) (*Order, error) {
	tradeEngine, err := e.getTradeEngine(api.Exchange, api.ApiKey, api.ApiSecret, api.Passphrase)
	if err != nil {
		return nil, err
	}
	return tradeEngine.CancelOrder(req)
}

// 批量下单
func (e *ExchangeManager) CreateOrders(api ExchangeApiParam, req []*OrderParam) ([]*Order, error) {
	tradeEngine, err := e.getTradeEngine(api.Exchange, api.ApiKey, api.ApiSecret, api.Passphrase)
	if err != nil {
		return nil, err
	}
	return tradeEngine.CreateOrders(req)
}

// 批量改单
func (e *ExchangeManager) AmendOrders(api ExchangeApiParam, req []*OrderParam) ([]*Order, error) {
	tradeEngine, err := e.getTradeEngine(api.Exchange, api.ApiKey, api.ApiSecret, api.Passphrase)
	if err != nil {
		return nil, err
	}
	return tradeEngine.AmendOrders(req)
}

// 批量撤单
func (e *ExchangeManager) CancelOrders(api ExchangeApiParam, req []*OrderParam) ([]*Order, error) {
	tradeEngine, err := e.getTradeEngine(api.Exchange, api.ApiKey, api.ApiSecret, api.Passphrase)
	if err != nil {
		return nil, err
	}
	return tradeEngine.CancelOrders(req)
}

// 订阅订单
func (e *ExchangeManager) SubscribeOrder(api ExchangeApiParam, req *SubscribeOrderParam) (TradeSubscribe[Order], error) {
	tradeEngine, err := e.getTradeEngine(api.Exchange, api.ApiKey, api.ApiSecret, api.Passphrase)
	if err != nil {
		return nil, err
	}
	return tradeEngine.SubscribeOrder(req)
}

// websocket下单
func (e *ExchangeManager) WsCreateOrder(api ExchangeApiParam, req *OrderParam) (*Order, error) {
	tradeEngine, err := e.getTradeEngine(api.Exchange, api.ApiKey, api.ApiSecret, api.Passphrase)
	if err != nil {
		return nil, err
	}
	return tradeEngine.WsCreateOrder(req)
}

// websocket改单
func (e *ExchangeManager) WsAmendOrder(api ExchangeApiParam, req *OrderParam) (*Order, error) {
	tradeEngine, err := e.getTradeEngine(api.Exchange, api.ApiKey, api.ApiSecret, api.Passphrase)
	if err != nil {
		return nil, err
	}
	return tradeEngine.WsAmendOrder(req)
}

// websocket撤单
func (e *ExchangeManager) WsCancelOrder(api ExchangeApiParam, req *OrderParam) (*Order, error) {
	tradeEngine, err := e.getTradeEngine(api.Exchange, api.ApiKey, api.ApiSecret, api.Passphrase)
	if err != nil {
		return nil, err
	}
	return tradeEngine.WsCancelOrder(req)
}

// websocket批量下单
func (e *ExchangeManager) WsCreateOrders(api ExchangeApiParam, req []*OrderParam) ([]*Order, error) {
	tradeEngine, err := e.getTradeEngine(api.Exchange, api.ApiKey, api.ApiSecret, api.Passphrase)
	if err != nil {
		return nil, err
	}
	return tradeEngine.WsCreateOrders(req)
}

// websocket批量改单
func (e *ExchangeManager) WsAmendOrders(api ExchangeApiParam, req []*OrderParam) ([]*Order, error) {
	tradeEngine, err := e.getTradeEngine(api.Exchange, api.ApiKey, api.ApiSecret, api.Passphrase)
	if err != nil {
		return nil, err
	}
	return tradeEngine.WsAmendOrders(req)
}

// websocket批量撤单
func (e *ExchangeManager) WsCancelOrders(api ExchangeApiParam, req []*OrderParam) ([]*Order, error) {
	tradeEngine, err := e.getTradeEngine(api.Exchange, api.ApiKey, api.ApiSecret, api.Passphrase)
	if err != nil {
		return nil, err
	}
	return tradeEngine.WsCancelOrders(req)
}
