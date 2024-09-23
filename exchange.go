package mytrade

type ExchangeBase struct {
	exchangeType ExchangeType
}

func (e *ExchangeBase) ExchangeType() ExchangeType {
	return e.exchangeType
}

func NewBinanceExchange() TradeExchange {
	return &BinanceExchange{
		ExchangeBase: ExchangeBase{
			exchangeType: BINANCE_NAME,
		},
	}
}
func NewOkxExchange() TradeExchange {
	return &OkxExchange{
		ExchangeBase: ExchangeBase{
			exchangeType: OKX_NAME,
		},
	}
}
func NewBybitExchange() TradeExchange {
	return &BybitExchange{
		ExchangeBase: ExchangeBase{
			exchangeType: BYBIT_NAME,
		},
	}
}
