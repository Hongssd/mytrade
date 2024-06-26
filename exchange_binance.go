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
		exchangeBase: b.exchangeBase,
		apiKey:       apiKey,
		secretKey:    secretKey,
	}
}

// 获取账户信息
func (b *BinanceExchange) NewTradeAccount(apiKey, secretKey, passphrase string) TradeAccount {
	return &BinanceTradeAccount{
		exchangeBase: b.exchangeBase,
		apiKey:       apiKey,
		secretKey:    secretKey,
	}
}
