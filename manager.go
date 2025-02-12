package mytrade

import (
	"errors"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type ExchangeManager struct {
	//交易所管理器
	ExchangeMap     *MySyncMap[string, TradeExchange]
	ExchangeInfoMap *MySyncMap[string, TradeExchangeInfo]
	MarketDataMap   *MySyncMap[string, TradeMarketData]
	TradeEngineMap  *MySyncMap[string, TradeEngine]
	TradeAccountMap *MySyncMap[string, TradeAccount]
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
		TradeAccountMap: GetPointer(NewMySyncMap[string, TradeAccount]()),
	}
	e.ExchangeMap.Store(BINANCE_NAME.String(), NewBinanceExchange())
	e.ExchangeMap.Store(OKX_NAME.String(), NewOkxExchange())
	e.ExchangeMap.Store(BYBIT_NAME.String(), NewBybitExchange())
	e.ExchangeMap.Store(GATE_NAME.String(), NewGateExchange())
	return e
}

var InnerExchangeManager *ExchangeManager = NewExchangeManager()

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

// 获取交易账号
func (e *ExchangeManager) getTradeAccount(exchange string, apikey, apiSecret, passphrase string) (TradeAccount, error) {
	tradeExchange, ok := e.ExchangeMap.Load(exchange)
	if !ok {
		return nil, errors.Join(errors.New(exchange), ErrorNotSupport)
	}

	key := fmt.Sprintf("%s_%s", exchange, apikey)

	tradeAccount, ok := e.TradeAccountMap.Load(key)
	if !ok {
		tradeAccount = tradeExchange.NewTradeAccount(apikey, apiSecret, passphrase)
		e.TradeAccountMap.Store(key, tradeAccount)
	}
	return tradeAccount, nil
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

// 获取全部交易对规范
func (e *ExchangeManager) GetAllSymbolInfo(exchange, accountType string) ([]TradeSymbolInfo, error) {
	tradeExchangeInfo, err := e.getTradeExchangeInfo(exchange)
	if err != nil {
		return nil, err
	}
	return tradeExchangeInfo.GetAllSymbolInfo(accountType)
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

// 获取账户模式
func (e *ExchangeManager) GetAccountMode(api ExchangeApiParam) (AccountMode, error) {
	tradeAccount, err := e.getTradeAccount(api.Exchange, api.ApiKey, api.ApiSecret, api.Passphrase)
	if err != nil {
		return ACCOUNT_MODE_UNKNOWN, err
	}
	return tradeAccount.GetAccountMode()
}

// 获取保证金模式
func (e *ExchangeManager) GetMarginMode(api ExchangeApiParam, accountType, symbol string, positionSide PositionSide) (MarginMode, error) {
	tradeAccount, err := e.getTradeAccount(api.Exchange, api.ApiKey, api.ApiSecret, api.Passphrase)
	if err != nil {
		return MARGIN_MODE_UNKNOWN, err
	}
	return tradeAccount.GetMarginMode(accountType, symbol, positionSide)
}

// 获取仓位模式
func (e *ExchangeManager) GetPositionMode(api ExchangeApiParam, accountType, symbol string) (PositionMode, error) {
	tradeAccount, err := e.getTradeAccount(api.Exchange, api.ApiKey, api.ApiSecret, api.Passphrase)
	if err != nil {
		return POSITION_MODE_UNKNOWN, err
	}
	return tradeAccount.GetPositionMode(accountType, symbol)
}

// 获取杠杆倍数
func (e *ExchangeManager) GetLeverage(api ExchangeApiParam, accountType, symbol string, marginMode MarginMode, positionSide PositionSide) (decimal.Decimal, error) {
	tradeAccount, err := e.getTradeAccount(api.Exchange, api.ApiKey, api.ApiSecret, api.Passphrase)
	if err != nil {
		return decimal.Zero, err
	}
	return tradeAccount.GetLeverage(accountType, symbol, marginMode, positionSide)
}

// 获取手续费率
func (e *ExchangeManager) GetFeeRate(api ExchangeApiParam, accountType, symbol string) (*FeeRate, error) {
	tradeAccount, err := e.getTradeAccount(api.Exchange, api.ApiKey, api.ApiSecret, api.Passphrase)
	if err != nil {
		return nil, err
	}
	return tradeAccount.GetFeeRate(accountType, symbol)
}

// 获取持仓
func (e *ExchangeManager) GetPositions(api ExchangeApiParam, accountType string, symbols ...string) ([]*Position, error) {
	tradeAccount, err := e.getTradeAccount(api.Exchange, api.ApiKey, api.ApiSecret, api.Passphrase)
	if err != nil {
		return nil, err
	}
	return tradeAccount.GetPositions(accountType, symbols...)
}

// 获取资产
func (e *ExchangeManager) GetAssets(api ExchangeApiParam, accountType string, currencies ...string) ([]*Asset, error) {
	tradeAccount, err := e.getTradeAccount(api.Exchange, api.ApiKey, api.ApiSecret, api.Passphrase)
	if err != nil {
		return nil, err
	}
	return tradeAccount.GetAssets(accountType, currencies...)
}

// 设置账户模式
func (e *ExchangeManager) SetAccountMode(api ExchangeApiParam, mode AccountMode) error {
	tradeAccount, err := e.getTradeAccount(api.Exchange, api.ApiKey, api.ApiSecret, api.Passphrase)
	if err != nil {
		return err
	}
	return tradeAccount.SetAccountMode(mode)
}

// 设置保证金模式
func (e *ExchangeManager) SetMarginMode(api ExchangeApiParam, accountType, symbol string, mode MarginMode) error {
	tradeAccount, err := e.getTradeAccount(api.Exchange, api.ApiKey, api.ApiSecret, api.Passphrase)
	if err != nil {
		return err
	}
	return tradeAccount.SetMarginMode(accountType, symbol, mode)
}

// 设置仓位模式
func (e *ExchangeManager) SetPositionMode(api ExchangeApiParam, accountType, symbol string, mode PositionMode) error {
	tradeAccount, err := e.getTradeAccount(api.Exchange, api.ApiKey, api.ApiSecret, api.Passphrase)
	if err != nil {
		return err
	}
	return tradeAccount.SetPositionMode(accountType, symbol, mode)
}

// 设置杠杆倍数
func (e *ExchangeManager) SetLeverage(api ExchangeApiParam, accountType, symbol string, marginMode MarginMode, positionSide PositionSide, leverage decimal.Decimal) error {
	tradeAccount, err := e.getTradeAccount(api.Exchange, api.ApiKey, api.ApiSecret, api.Passphrase)
	if err != nil {
		return err
	}
	return tradeAccount.SetLeverage(accountType, symbol, marginMode, positionSide, leverage)
}

// 资金划转（账户内）
func (e *ExchangeManager) AssetTransfer(api ExchangeApiParam, req *AssetTransferParams) ([]*AssetTransfer, error) {
	tradeAccount, err := e.getTradeAccount(api.Exchange, api.ApiKey, api.ApiSecret, api.Passphrase)
	if err != nil {
		return nil, err
	}
	return tradeAccount.AssetTransfer(req)
}

// 资金划转历史记录
func (e *ExchangeManager) QueryAssetTransfer(api ExchangeApiParam, req *QueryAssetTransferParams) ([]*QueryAssetTransfer, error) {
	tradeAccount, err := e.getTradeAccount(api.Exchange, api.ApiKey, api.ApiSecret, api.Passphrase)
	if err != nil {
		return nil, err
	}
	return tradeAccount.QueryAssetTransfer(req)
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
