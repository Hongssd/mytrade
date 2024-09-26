package mytrade

import (
	"github.com/Hongssd/mybybitapi"
	"strconv"
)

type BybitMarketData struct {
	ExchangeBase
}

func (m *BybitMarketData) NewKlineReq() *KlineParam {
	return &KlineParam{}
}
func (m *BybitMarketData) NewBookReq() *BookParam {
	return &BookParam{}
}

func (m *BybitMarketData) GetKline(req *KlineParam) (*[]Kline, error) {
	if req == nil || req.AccountType == "" || req.Symbol == "" || req.Interval == "" {
		return nil, ErrorInvalidParam
	}

	client := mybybitapi.NewRestClient("", "").PublicRestClient()
	switch BybitAccountType(req.AccountType) {
	case BYBIT_AC_SPOT, BYBIT_AC_LINEAR, BYBIT_AC_INVERSE:
		api := client.NewMarketKline().Category(req.AccountType).Symbol(req.Symbol).Interval(req.Interval)
		if req.StartTime != 0 {
			api.Start(req.StartTime)
		}
		if req.EndTime != 0 {
			api.End(req.EndTime)
		}
		if req.Limit != 0 {
			api.Limit(req.Limit)
		}
		data, err := api.Do()
		if err != nil {
			return nil, err
		}
		return m.convertToKline(req.AccountType, req.Symbol, req.Interval, &data.Result), nil
	default:
		return nil, ErrorInvalidParam
	}
}
func (m *BybitMarketData) GetBook(req *BookParam) (*OrderBook, error) {
	if req == nil || req.AccountType == "" || req.Symbol == "" {
		return nil, ErrorInvalidParam
	}

	client := mybybitapi.NewRestClient("", "").PublicRestClient()
	switch BybitAccountType(req.AccountType) {
	case BYBIT_AC_SPOT, BYBIT_AC_LINEAR, BYBIT_AC_INVERSE:
		api := client.NewMarketOrderBook().Category(req.AccountType).Symbol(req.Symbol)
		if req.Level != 0 {
			api.Limit(req.Level)
		}
		data, err := api.Do()
		if err != nil {
			return nil, err
		}
		return m.convertToOrderBook(req.AccountType, req.Symbol, &data.Result), nil
	default:
		return nil, ErrorAccountType
	}
}

func (m *BybitMarketData) convertToKline(accountType, symbol, interval string, data *mybybitapi.MarketKlineRes) *[]Kline {
	var list []Kline
	if data == nil || len(data.List) == 0 {
		return &list
	}
	for _, v := range data.List {
		startTime, err := strconv.ParseInt(v.StartTime, 10, 64)
		if err != nil {
			continue
		}

		kline := Kline{
			Exchange:             BYBIT_NAME.String(),
			AccountType:          accountType,
			Symbol:               symbol,
			Interval:             interval,
			StartTime:            startTime,
			Open:                 stringToFloat64(v.OpenPrice),
			High:                 stringToFloat64(v.HighPrice),
			Low:                  stringToFloat64(v.LowPrice),
			Close:                stringToFloat64(v.ClosePrice),
			Volume:               stringToFloat64(v.Volume),
			CloseTime:            bybitGetKlineCloseTime(startTime, interval),
			TransactionVolume:    stringToFloat64(v.Turnover),
			TransactionNumber:    0,
			BuyTransactionVolume: 0,
			BuyTransactionAmount: 0,
		}
		list = append(list, kline)
	}
	return &list
}

func (m *BybitMarketData) convertToOrderBook(accountType, symbol string, data *mybybitapi.MarketOrderBookRes) *OrderBook {
	var ob OrderBook
	ob.Exchange = BYBIT_NAME.String()
	ob.AccountType = accountType
	ob.Symbol = symbol
	if data == nil {
		return &ob
	}

	ob.Timestamp = data.Ts
	ob.Asks = make([]Book, len(data.Asks))
	ob.Bids = make([]Book, len(data.Bids))
	for i, v := range data.Asks {
		ob.Asks[i] = Book{
			Price:    v.Price,
			Quantity: v.Quantity,
		}
	}
	for i, v := range data.Bids {
		ob.Bids[i] = Book{
			Price:    v.Price,
			Quantity: v.Quantity,
		}
	}
	return &ob
}
