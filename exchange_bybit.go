package mytrade

type BybitExchange struct {
	exchangeBase
	*BybitExchangeInfo
	*BybitMarketData
	*BybitTradeEngine
}

func (b *BybitExchange) ExchangeInfo() TradeExchangeInfo {
	if b.BybitExchangeInfo == nil {
		b.BybitExchangeInfo = &BybitExchangeInfo{
			isLoaded: false,
		}
	}
	return b.BybitExchangeInfo
}

// 获取市场数据
func (b *BybitExchange) MarketData() TradeMarketData {
	if b.BybitMarketData == nil {
		b.BybitMarketData = &BybitMarketData{}
	}
	return b.BybitMarketData
}

// 获取交易引擎
func (b *BybitExchange) TradeEngine(apiKey, secretKey string, options ...TradeEngineOption) TradeEngine {

	if b.BybitTradeEngine == nil {
		b.BybitTradeEngine = &BybitTradeEngine{
			apiKey:    apiKey,
			secretKey: secretKey,
		}
		if len(options) > 0 {
			for _, option := range options {
				option(b.BybitTradeEngine)
			}
		}
	}
	return b.BybitTradeEngine
}
