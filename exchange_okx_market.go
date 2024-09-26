package mytrade

import (
	"github.com/Hongssd/myokxapi"
	"strconv"
)

type OkxMarketData struct {
	ExchangeBase
}

func (m *OkxMarketData) NewKlineReq() *KlineParam {
	return &KlineParam{}
}
func (m *OkxMarketData) NewBookReq() *BookParam {
	return &BookParam{}
}

func (m *OkxMarketData) GetKline(req *KlineParam) (*[]Kline, error) {
	if req == nil || req.AccountType == "" || req.Symbol == "" || req.Interval == "" {
		return nil, ErrorInvalidParam
	}

	client := okx.NewRestClient("", "", "").PublicRestClient()
	switch OkxAccountType(req.AccountType) {
	case OKX_AC_SPOT, OKX_AC_SWAP, OKX_AC_FUTURES:
		api := client.NewPublicRestMarketCandles().InstId(req.Symbol).Bar(req.Interval)
		if req.StartTime != 0 {
			api.Before(req.StartTime)
		}
		if req.EndTime != 0 {
			api.After(req.EndTime)
		}
		if req.Limit != 0 {
			api.Limit(req.Limit)
		}
		data, err := api.Do()
		if err != nil {
			return nil, err
		}
		return m.convertToKline(req.AccountType, req.Symbol, req.Interval, &data.Data), nil
	default:
		return nil, ErrorInvalidParam
	}
}
func (m *OkxMarketData) GetBook(req *BookParam) (*OrderBook, error) {
	if req == nil || req.AccountType == "" || req.Symbol == "" {
		return nil, ErrorInvalidParam
	}

	client := okx.NewRestClient("", "", "").PublicRestClient()
	switch OkxAccountType(req.AccountType) {
	case OKX_AC_SPOT, OKX_AC_SWAP, OKX_AC_FUTURES:
		api := client.NewPublicRestMarketBooks().InstId(req.Symbol)
		if req.Level != 0 {
			api.Sz(strconv.Itoa(req.Level))
		}
		data, err := api.Do()
		if err != nil {
			return nil, err
		}
		return m.convertToOrderBook(req.AccountType, req.Symbol, &data.Data), nil
	default:
		return nil, ErrorAccountType
	}

}

func (m *OkxMarketData) convertToKline(accountType, symbol, interval string, data *myokxapi.PublicRestMarketCandlesRes) *[]Kline {
	var list []Kline
	if data == nil || len(*data) == 0 {
		return &list
	}
	for _, v := range *data {
		startTime, err := strconv.ParseInt(v.Ts, 10, 64)
		if err != nil {
			continue
		}
		volume := 0.0
		if OkxAccountType(accountType) == OKX_AC_SPOT {
			volume = stringToFloat64(v.Vol)
		} else {
			volume = stringToFloat64(v.VolCcy)
		}
		kline := Kline{
			Exchange:             OKX_NAME.String(),
			AccountType:          accountType,
			Symbol:               symbol,
			Interval:             interval,
			StartTime:            startTime,
			Open:                 stringToFloat64(v.O),
			High:                 stringToFloat64(v.H),
			Low:                  stringToFloat64(v.L),
			Close:                stringToFloat64(v.C),
			Volume:               volume,
			CloseTime:            okxGetKlineCloseTime(startTime, interval),
			TransactionVolume:    stringToFloat64(v.VolCcyQuote),
			TransactionNumber:    0,
			BuyTransactionVolume: 0,
			BuyTransactionAmount: 0,
		}
		list = append(list, kline)
	}
	return &list
}

func (m *OkxMarketData) convertToOrderBook(accountType, symbol string, dataList *myokxapi.PublicRestMarketBooksRes) *OrderBook {
	var ob OrderBook
	ob.Exchange = OKX_NAME.String()
	ob.AccountType = accountType
	ob.Symbol = symbol
	if dataList == nil || len(*dataList) != 1 {
		return &ob
	}

	data := (*dataList)[0]

	ob.Timestamp = stringToInt64(data.Ts)
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
