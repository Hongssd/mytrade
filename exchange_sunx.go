package mytrade

type SunxExchange struct {
	ExchangeBase
}

// 获取交易规范
func (b *SunxExchange) NewExchangeInfo() TradeExchangeInfo {
	return &SunxExchangeInfo{
		isLoaded: false,
	}
}

// 获取市场数据
func (b *SunxExchange) NewMarketData() TradeMarketData {
	return &SunxMarketData{}
}

// 获取交易引擎
func (b *SunxExchange) NewTradeEngine(accessKey, secretKey, passphrase string) TradeEngine {
	return &SunxTradeEngine{
		ExchangeBase: b.ExchangeBase,
		accessKey:    accessKey,
		secretKey:    secretKey,
	}
}

// 获取账户信息
func (b *SunxExchange) NewTradeAccount(accessKey, secretKey, passphrase string) TradeAccount {
	return &SunxTradeAccount{
		ExchangeBase: b.ExchangeBase,
		accessKey:    accessKey,
		secretKey:    secretKey,
	}
}
