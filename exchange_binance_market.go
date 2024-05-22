package mytrade

import (
	"github.com/Hongssd/mybinanceapi"
	"time"
)

type BinanceMarketData struct {
	exchangeBase
}

func (m *BinanceMarketData) NewKlineReq() *KlineParam {
	return &KlineParam{}
}
func (m *BinanceMarketData) NewBookReq() *BookParam {
	return &BookParam{}
}

func (m *BinanceMarketData) GetKline(req *KlineParam) (*[]Kline, error) {
	if req == nil || req.AccountType == "" || req.Symbol == "" || req.Interval == "" {
		return nil, ErrorInvalidParam
	}
	client := binance.NewSpotRestClient("", "")
	switch BinanceAccountType(req.AccountType) {
	case BN_AC_SPOT:
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
	case BN_AC_FUTURE:
		api := binance.NewFutureRestClient("", "").NewFutureKlines().
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
	case BN_AC_SWAP:
		api := binance.NewSwapRestClient("", "").NewSwapKlines().
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
func (m *BinanceMarketData) GetBook(req *BookParam) (*OrderBook, error) {
	if req == nil || req.AccountType == "" || req.Symbol == "" {
		return nil, ErrorInvalidParam
	}
	client := binance.NewSpotRestClient("", "")
	switch BinanceAccountType(req.AccountType) {
	case BN_AC_SPOT:
		api := client.NewSpotDepth().Symbol(req.Symbol)
		if req.Level != 0 {
			api.Limit(req.Level)
		}
		data, err := api.Do()
		if err != nil {
			return nil, err
		}
		return m.convertToOrderBook(req.AccountType, req.Symbol, data.Bids, data.Asks), nil
	case BN_AC_FUTURE:
		api := binance.NewFutureRestClient("", "").NewFutureDepth().Symbol(req.Symbol)
		if req.Level != 0 {
			api.Limit(req.Level)
		}
		data, err := api.Do()
		if err != nil {
			return nil, err
		}
		return m.convertToOrderBook(req.AccountType, req.Symbol, data.Bids, data.Asks), nil
	case BN_AC_SWAP:
		api := binance.NewSwapRestClient("", "").NewSwapDepth().Symbol(req.Symbol)
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

func (m *BinanceMarketData) convertToKline(accountType, symbol, interval string, data *mybinanceapi.KlinesRes) *[]Kline {
	var list []Kline
	if data == nil || len(*data) == 0 {
		return &list
	}
	for _, v := range *data {
		kline := Kline{
			Exchange:             BINANCE_NAME.String(),
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

func (m *BinanceMarketData) convertToOrderBook(accountType, symbol string, bids, asks []mybinanceapi.DepthGear) *OrderBook {
	var ob OrderBook
	ob.Exchange = BINANCE_NAME.String()
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
