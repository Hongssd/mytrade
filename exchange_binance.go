package mytrade

type BinanceExchange struct {
	exchangeBase
}

// 获取交易规范
func (b *BinanceExchange) NewExchangeInfo() TradeExchangeInfo {
	return &BinanceExchangeInfo{
		isLoaded: false,
	}
}

// 获取市场数据
func (b *BinanceExchange) NewMarketData() TradeMarketData {
	return &BinanceMarketData{}
}

// 获取交易引擎
func (b *BinanceExchange) NewTradeEngine(apiKey, secretKey, passphrase string) TradeEngine {
	return &BinanceTradeEngine{
		apiKey:    apiKey,
		secretKey: secretKey,
	}
}
