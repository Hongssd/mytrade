package mytrade

type exchangeBase struct {
	exchangeType ExchangeType
}

func (e *exchangeBase) ExchangeType() ExchangeType {
	return e.exchangeType
}

func NewBinanceExchange() TradeExchange {
	return &BinanceExchange{
		exchangeBase: exchangeBase{
			exchangeType: BINANCE_NAME,
		},
	}
}
func NewOkxExchange() TradeExchange {
	return &OkxExchange{
		exchangeBase: exchangeBase{
			exchangeType: OKX_NAME,
		},
	}
}
func NewBybitExchange() TradeExchange {
	return &BybitExchange{
		exchangeBase: exchangeBase{
			exchangeType: BYBIT_NAME,
		},
	}
}
