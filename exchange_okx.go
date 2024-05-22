package mytrade

type OkxExchange struct {
	exchangeBase
	*OkxExchangeInfo
	*OkxMarketData
	*OkxTradeEngine
}

func (o *OkxExchange) ExchangeInfo() TradeExchangeInfo {
	if o.OkxExchangeInfo == nil {
		o.OkxExchangeInfo = &OkxExchangeInfo{
			isLoaded: false,
		}
	}
	return o.OkxExchangeInfo
}

// 获取市场数据
func (o *OkxExchange) MarketData() TradeMarketData {
	if o.OkxMarketData == nil {
		o.OkxMarketData = &OkxMarketData{}
	}
	return o.OkxMarketData
}

// 获取交易引擎
func (o *OkxExchange) TradeEngine(apiKey, secretKey string, options ...TradeEngineOption) TradeEngine {

	if o.OkxTradeEngine == nil {
		o.OkxTradeEngine = &OkxTradeEngine{
			apiKey:    apiKey,
			secretKey: secretKey,
		}
		if len(options) > 0 {
			for _, option := range options {
				option(o.OkxTradeEngine)
			}
		}
	}
	return o.OkxTradeEngine
}

func WithOkxPassphrase(passphrase string) TradeEngineOption {
	return func(engine TradeEngine) {
		if okxEngine, ok := engine.(*OkxTradeEngine); ok {
			okxEngine.passphrase = passphrase
		}
	}
}
