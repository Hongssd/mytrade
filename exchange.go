package mytrade

type exchangeBase struct {
	exchangeType ExchangeType
}

func (e *exchangeBase) ExchangeType() ExchangeType {
	return e.exchangeType
}

func NewBinanceExchange() TradeExchange {
	return &BinanceExchange{}
}
func NewOkxExchange() TradeExchange {
	return &OkxExchange{}
}
func NewBybitExchange() TradeExchange {
	return &BybitExchange{}
}
