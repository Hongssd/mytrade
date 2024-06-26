package mytrade

type OkxExchange struct {
	exchangeBase
}

// 获取交易规范
func (o *OkxExchange) NewExchangeInfo() TradeExchangeInfo {
	return &OkxExchangeInfo{
		isLoaded: false,
	}
}

// 获取市场数据
func (o *OkxExchange) NewMarketData() TradeMarketData {
	return &OkxMarketData{}
}

// 获取交易引擎
func (o *OkxExchange) NewTradeEngine(apiKey, secretKey, passphrase string) TradeEngine {
	return &OkxTradeEngine{
		exchangeBase: o.exchangeBase,
		apiKey:       apiKey,
		secretKey:    secretKey,
		passphrase:   passphrase,
	}
}

// 获取交易账户
func (o *OkxExchange) NewTradeAccount(apiKey, secretKey, passphrase string) TradeAccount {
	return &OkxTradeAccount{
		exchangeBase: o.exchangeBase,
		apiKey:       apiKey,
		secretKey:    secretKey,
		passphrase:   passphrase,
	}
}
