package mytrade

type AsterExchange struct {
	ExchangeBase
}

// 获取交易规范
func (b *AsterExchange) NewExchangeInfo() TradeExchangeInfo {
	return &AsterExchangeInfo{
		isLoaded: false,
	}
}

// 获取市场数据
func (b *AsterExchange) NewMarketData() TradeMarketData {
	return &AsterMarketData{}
}

// 获取交易引擎
func (b *AsterExchange) NewTradeEngine(apiKey, secretKey, passphrase string) TradeEngine {
	return &AsterTradeEngine{
		ExchangeBase: b.ExchangeBase,
		apiKey:       apiKey,
		secretKey:    secretKey,
	}
}

// 获取账户信息
func (b *AsterExchange) NewTradeAccount(apiKey, secretKey, passphrase string) TradeAccount {
	return &AsterTradeAccount{
		ExchangeBase: b.ExchangeBase,
		apiKey:       apiKey,
		secretKey:    secretKey,
	}
}
