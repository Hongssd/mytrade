package mytrade

import (
	"github.com/Hongssd/mysunxapi"
	"github.com/shopspring/decimal"
)

type SunxMarketData struct {
	ExchangeBase
}

func (m *SunxMarketData) NewKlineReq() *KlineParam {
	return &KlineParam{}
}
func (m *SunxMarketData) NewBookReq() *BookParam {
	return &BookParam{}
}

func (m *SunxMarketData) GetKline(req *KlineParam) (*[]Kline, error) {
	if req == nil || req.AccountType == "" || req.Symbol == "" || req.Interval == "" {
		return nil, ErrorInvalidParam
	}

	client := sunx.NewPublicRestClient()
	switch SunxAccountType(req.AccountType) {
	case SUNX_ACCOUNT_TYPE_SWAP:
		api := client.NewPublicRestMarketHistoryKline().ContractCode(req.Symbol).Period(req.Interval)
		if req.StartTime != 0 {
			api.From(req.StartTime)
		}
		if req.EndTime != 0 {
			api.To(req.EndTime)
		}
		if req.Limit != 0 {
			api.Size(req.Limit)
		}
		data, err := api.Do()
		if err != nil {
			return nil, err
		}
		return m.convertToSwapKline(req.AccountType, req.Symbol, req.Interval, &data.Data), nil
	default:
		return nil, ErrorAccountType
	}
}

func (m *SunxMarketData) convertToSwapKline(accountType, symbol, interval string, data *mysunxapi.PublicRestMarketHistoryKlineRes) *[]Kline {
	var list []Kline
	if data == nil || len(*data) == 0 {
		return &list
	}
	for _, v := range *data {
		startTime := decimal.NewFromInt(v.Id).IntPart()
		kline := Kline{
			Exchange:             SUNX_NAME.String(),
			AccountType:          accountType,
			Symbol:               symbol,
			Interval:             interval,
			StartTime:            startTime,
			Open:                 v.Open,
			High:                 v.High,
			Low:                  v.Low,
			Close:                v.Close,
			Volume:               v.Vol,
			CloseTime:            sunxGetKlineCloseTime(startTime, interval),
			TransactionVolume:    v.Amount,
			TransactionNumber:    decimal.NewFromFloat(v.Count).IntPart(),
			BuyTransactionVolume: 0,
			BuyTransactionAmount: 0,
		}
		list = append(list, kline)
	}
	return &list
}

func (m *SunxMarketData) GetBook(req *BookParam) (*OrderBook, error) {
	if req == nil || req.AccountType == "" || req.Symbol == "" {
		return nil, ErrorInvalidParam
	}

	client := sunx.NewPublicRestClient()
	switch SunxAccountType(req.AccountType) {
	case SUNX_ACCOUNT_TYPE_SWAP:
		api := client.NewPublicRestMarketDepth().ContractCode(req.Symbol)
		if req.Step != "" {
			api.Type(req.Step)
		} else {
			if req.Level <= 20 {
				api.Type("step6") //20档
			} else {
				api.Type("step0") //150档
			}

		}
		data, err := api.Do()
		if err != nil {
			return nil, err
		}
		if req.Level > 0 {
			askEnd := min(req.Level, len(data.Data.Asks))
			bidEnd := min(req.Level, len(data.Data.Bids))
			data.Data.Asks = data.Data.Asks[:askEnd]
			data.Data.Bids = data.Data.Bids[:bidEnd]
		}
		return m.convertToSwapOrderBook(req.AccountType, req.Symbol, &data.Data), nil
	default:
		return nil, ErrorAccountType
	}
}

func (m *SunxMarketData) convertToSwapOrderBook(accountType, symbol string, data *mysunxapi.PublicRestMarketDepthRes) *OrderBook {
	var ob OrderBook
	ob.Exchange = SUNX_NAME.String()
	ob.AccountType = accountType
	ob.Symbol = symbol
	ob.Timestamp = data.Ts
	ob.Asks = make([]Book, len(data.Asks))
	ob.Bids = make([]Book, len(data.Bids))
	for i, v := range data.Asks {
		ob.Asks[i] = Book{
			Price:    decimal.NewFromFloat(v.Price).String(),
			Quantity: decimal.NewFromFloat(v.Volume).String(),
		}
	}
	for i, v := range data.Bids {
		ob.Bids[i] = Book{
			Price:    decimal.NewFromFloat(v.Price).String(),
			Quantity: decimal.NewFromFloat(v.Volume).String(),
		}
	}
	return &ob
}
