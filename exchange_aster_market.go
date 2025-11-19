package mytrade

import (
	"time"

	"github.com/Hongssd/myasterapi"
)

type AsterMarketData struct {
	ExchangeBase
}

func (m *AsterMarketData) NewKlineReq() *KlineParam {
	return &KlineParam{}
}
func (m *AsterMarketData) NewBookReq() *BookParam {
	return &BookParam{}
}

func (m *AsterMarketData) GetKline(req *KlineParam) (*[]Kline, error) {
	if req == nil || req.AccountType == "" || req.Symbol == "" || req.Interval == "" {
		return nil, ErrorInvalidParam
	}
	client := aster.NewSpotRestClient("", "")
	switch AsterAccountType(req.AccountType) {
	case ASTER_AC_SPOT:
		api := client.NewSpotKlines().
			Symbol(req.Symbol).Interval(req.Interval)
		if req.StartTime != 0 {
			api.StartTime(req.StartTime)
		}
		if req.EndTime != 0 {
			api.EndTime(req.EndTime)
		}
		if req.Limit != 0 {
			api.Limit(req.Limit)
		}
		data, err := api.Do()
		if err != nil {
			return nil, err
		}
		return m.convertToKline(req.AccountType, req.Symbol, req.Interval, data), nil
	case ASTER_AC_FUTURE:
		api := aster.NewFutureRestClient("", "").NewFutureKlines().
			Symbol(req.Symbol).Interval(req.Interval)
		if req.StartTime != 0 {
			api.StartTime(req.StartTime)
		}
		if req.EndTime != 0 {
			api.EndTime(req.EndTime)
		}
		if req.Limit != 0 {
			api.Limit(req.Limit)
		}
		data, err := api.Do()
		if err != nil {
			return nil, err
		}
		return m.convertToKline(req.AccountType, req.Symbol, req.Interval, data), nil
	default:
		return nil, ErrorInvalidParam
	}
}
func (m *AsterMarketData) GetBook(req *BookParam) (*OrderBook, error) {
	if req == nil || req.AccountType == "" || req.Symbol == "" {
		return nil, ErrorInvalidParam
	}
	client := aster.NewSpotRestClient("", "")
	switch AsterAccountType(req.AccountType) {
	case ASTER_AC_SPOT:
		api := client.NewSpotDepth().Symbol(req.Symbol)
		if req.Level != 0 {
			api.Limit(req.Level)
		}
		data, err := api.Do()
		if err != nil {
			return nil, err
		}
		return m.convertToOrderBook(req.AccountType, req.Symbol, data.Bids, data.Asks), nil
	case ASTER_AC_FUTURE:
		api := aster.NewFutureRestClient("", "").NewFutureDepth().Symbol(req.Symbol)
		if req.Level != 0 {
			api.Limit(req.Level)
		}
		data, err := api.Do()
		if err != nil {
			return nil, err
		}
		return m.convertToOrderBook(req.AccountType, req.Symbol, data.Bids, data.Asks), nil
	default:
		return nil, ErrorAccountType
	}
}

func (m *AsterMarketData) convertToKline(accountType, symbol, interval string, data *myasterapi.KlinesRes) *[]Kline {
	var list []Kline
	if data == nil || len(*data) == 0 {
		return &list
	}
	for _, v := range *data {
		kline := Kline{
			Exchange:             ASTER_NAME.String(),
			AccountType:          accountType,
			Symbol:               symbol,
			Interval:             interval,
			StartTime:            v.StartTime,
			Open:                 v.Open,
			High:                 v.High,
			Low:                  v.Low,
			Close:                v.Close,
			Volume:               v.Volume,
			CloseTime:            v.CloseTime,
			TransactionVolume:    v.TransactionVolume,
			TransactionNumber:    v.TransactionNumber,
			BuyTransactionVolume: v.BuyTransactionVolume,
			BuyTransactionAmount: v.BuyTransactionAmount,
		}
		list = append(list, kline)
	}
	return &list

}

func (m *AsterMarketData) convertToOrderBook(accountType, symbol string, bids, asks []myasterapi.DepthGear) *OrderBook {
	var ob OrderBook
	ob.Exchange = ASTER_NAME.String()
	ob.AccountType = accountType
	ob.Symbol = symbol
	ob.Timestamp = time.Now().UnixMilli()
	ob.Asks = make([]Book, len(asks))
	ob.Bids = make([]Book, len(bids))
	for i, v := range asks {
		ob.Asks[i] = Book{
			Price:    v.Price,
			Quantity: v.Quantity,
		}
	}
	for i, v := range bids {
		ob.Bids[i] = Book{
			Price:    v.Price,
			Quantity: v.Quantity,
		}
	}
	return &ob
}
