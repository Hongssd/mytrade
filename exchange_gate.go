package mytrade

type GateExchange struct {
	ExchangeBase
}

// 获取交易规范
func (o *GateExchange) NewExchangeInfo() TradeExchangeInfo {
	return &GateExchangeInfo{
		isLoaded: false,
	}
}

// 获取市场数据
func (o *GateExchange) NewMarketData() TradeMarketData {
	return &GateMarketData{}
}

// 获取交易引擎
func (o *GateExchange) NewTradeEngine(apiKey, secretKey, passphrase string) TradeEngine {
	return &GateTradeEngine{
		ExchangeBase: o.ExchangeBase,
		apiKey:       apiKey,
		secretKey:    secretKey,
	}
}

// 获取交易账户
func (o *GateExchange) NewTradeAccount(apiKey, secretKey, passphrase string) TradeAccount {
	return &GateTradeAccount{
		ExchangeBase: o.ExchangeBase,
		apiKey:       apiKey,
		secretKey:    secretKey,
	}
}
