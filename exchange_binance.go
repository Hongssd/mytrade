package mytrade

type BinanceExchange struct {
	exchangeBase
	*BinanceExchangeInfo
	*BinanceMarketData
	*BinanceTradeEngine
}

func (b *BinanceExchange) ExchangeInfo() TradeExchangeInfo {
	if b.BinanceExchangeInfo == nil {
		b.BinanceExchangeInfo = &BinanceExchangeInfo{
			isLoaded: false,
		}
	}
	return b.BinanceExchangeInfo
}

// 获取市场数据
func (b *BinanceExchange) MarketData() TradeMarketData {
	if b.BinanceMarketData == nil {
		b.BinanceMarketData = &BinanceMarketData{}
	}
	return b.BinanceMarketData
}

// 获取交易引擎
func (b *BinanceExchange) TradeEngine(apiKey, secretKey string, options ...TradeEngineOption) TradeEngine {
	if b.BinanceTradeEngine == nil {
		b.BinanceTradeEngine = &BinanceTradeEngine{
			apiKey:    apiKey,
			secretKey: secretKey,
		}
		if len(options) > 0 {
			for _, option := range options {
				option(b.BinanceTradeEngine)
			}
		}
	}
	return b.BinanceTradeEngine
}
