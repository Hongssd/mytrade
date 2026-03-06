package mytrade

import (
	"strconv"

	"github.com/Hongssd/myxcoinapi"
)

type XcoinMarketData struct {
	ExchangeBase
}

func (m *XcoinMarketData) NewKlineReq() *KlineParam {
	return &KlineParam{}
}

func (m *XcoinMarketData) NewBookReq() *BookParam {
	return &BookParam{}
}

func (m *XcoinMarketData) GetKline(req *KlineParam) (*[]Kline, error) {
	if req == nil || req.AccountType == "" || req.Symbol == "" || req.Interval == "" {
		return nil, ErrorInvalidParam
	}

	client := xcoin.NewRestClient("", "").PublicRestClient()
	switch XcoinAccountType(req.AccountType) {
	case XCOIN_ACCOUNT_TYPE_SPOT, XCOIN_ACCOUNT_TYPE_LINEAR_PERPETUAL, XCOIN_ACCOUNT_TYPE_LINEAR_FUTURES:
		api := client.NewPublicRestMarketKline().
			BusinessType(req.AccountType).
			Symbol(req.Symbol).
			Period(req.Interval)
		if req.StartTime != 0 {
			api.StartTime(strconv.FormatInt(req.StartTime, 10))
		}
		if req.EndTime != 0 {
			api.EndTime(strconv.FormatInt(req.EndTime, 10))
		}
		if req.Limit != 0 {
			api.Limit(strconv.Itoa(req.Limit))
		}
		data, err := api.Do()
		if err != nil {
			return nil, err
		}
		return m.convertToKline(req.AccountType, req.Symbol, req.Interval, &data.Data), nil
	default:
		return nil, ErrorAccountType
	}
}

func (m *XcoinMarketData) GetBook(req *BookParam) (*OrderBook, error) {
	if req == nil || req.AccountType == "" || req.Symbol == "" {
		return nil, ErrorInvalidParam
	}

	client := xcoin.NewRestClient("", "").PublicRestClient()
	switch XcoinAccountType(req.AccountType) {
	case XCOIN_ACCOUNT_TYPE_SPOT, XCOIN_ACCOUNT_TYPE_LINEAR_PERPETUAL, XCOIN_ACCOUNT_TYPE_LINEAR_FUTURES:
		api := client.NewPublicRestMarketDepth().
			BusinessType(req.AccountType).
			Symbol(req.Symbol)
		if req.Level != 0 {
			api.Limit(strconv.Itoa(req.Level))
		}
		data, err := api.Do()
		if err != nil {
			return nil, err
		}
		return m.convertToOrderBook(req.AccountType, req.Symbol, data.Ts, &data.Data), nil
	default:
		return nil, ErrorAccountType
	}
}

func (m *XcoinMarketData) convertToKline(accountType, symbol, interval string, data *myxcoinapi.PublicRestMarketKlineRes) *[]Kline {
	var list []Kline
	if data == nil || len(*data) == 0 {
		return &list
	}
	for _, v := range *data {
		kline := Kline{
			Exchange:             XCOIN_NAME.String(),
			AccountType:          accountType,
			Symbol:               symbol,
			Interval:             interval,
			StartTime:            stringToInt64(v.OpenTime),
			Open:                 stringToFloat64(v.OpenPrice),
			High:                 stringToFloat64(v.HighPrice),
			Low:                  stringToFloat64(v.LowPrice),
			Close:                stringToFloat64(v.ClosePrice),
			Volume:               stringToFloat64(v.FillQty),
			CloseTime:            stringToInt64(v.CloseTime),
			TransactionVolume:    stringToFloat64(v.FillAmount),
			TransactionNumber:    stringToInt64(v.Count),
			BuyTransactionVolume: 0,
			BuyTransactionAmount: 0,
		}
		list = append(list, kline)
	}
	return &list
}

func (m *XcoinMarketData) convertToOrderBook(accountType, symbol, ts string, data *myxcoinapi.PublicRestMarketDepthRes) *OrderBook {
	var ob OrderBook
	ob.Exchange = XCOIN_NAME.String()
	ob.AccountType = accountType
	ob.Symbol = symbol
	ob.Timestamp = stringToInt64(ts)
	if data == nil {
		return &ob
	}

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