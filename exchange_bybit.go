package mytrade

type BybitExchange struct {
	exchangeBase
}

func (b *BybitExchange) NewExchangeInfo() TradeExchangeInfo {
	return &BybitExchangeInfo{
		isLoaded: false,
	}
}

// 获取市场数据
func (b *BybitExchange) NewMarketData() TradeMarketData {
	return &BybitMarketData{}
}

// 获取交易引擎
func (b *BybitExchange) NewTradeEngine(apiKey, secretKey, passphrase string) TradeEngine {
	return &BybitTradeEngine{
		exchangeBase: b.exchangeBase,
		apiKey:       apiKey,
		secretKey:    secretKey,
	}
}

// 获取交易账户
func (b *BybitExchange) NewTradeAccount(apiKey, secretKey, passphrase string) TradeAccount {
	return &BybitTradeAccount{
		exchangeBase: b.exchangeBase,
		apiKey:       apiKey,
		secretKey:    secretKey,
	}
}
