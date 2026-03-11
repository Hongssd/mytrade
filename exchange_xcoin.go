package mytrade

type XcoinExchange struct {
	ExchangeBase
}

// 获取交易规范
func (x *XcoinExchange) NewExchangeInfo() TradeExchangeInfo {
	return &XcoinExchangeInfo{
		isLoaded: false,
	}
}

// 获取市场数据
func (x *XcoinExchange) NewMarketData() TradeMarketData {
	return &XcoinMarketData{}
}

// 获取交易引擎
func (x *XcoinExchange) NewTradeEngine(apiKey, secretKey, passphrase string) TradeEngine {
	return &XcoinTradeEngine{
		ExchangeBase: x.ExchangeBase,
		apiKey:       apiKey,
		apiSecret:    secretKey,
	}
}

// 获取交易账户
func (x *XcoinExchange) NewTradeAccount(apiKey, secretKey, passphrase string) TradeAccount {
	return &XcoinTradeAccount{
		ExchangeBase: x.ExchangeBase,
		apiKey:       apiKey,
		apiSecret:    secretKey,
	}
}
