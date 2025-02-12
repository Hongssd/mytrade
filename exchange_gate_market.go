package mytrade

import (
	"strconv"
	"time"

	"github.com/Hongssd/mygateapi"
)

type GateMarketData struct {
	ExchangeBase
}

func (m *GateMarketData) NewKlineReq() *KlineParam {
	return &KlineParam{}
}
func (m *GateMarketData) NewBookReq() *BookParam {
	return &BookParam{}
}

func (m *GateMarketData) GetKline(req *KlineParam) (*[]Kline, error) {
	if req == nil || req.AccountType == "" || req.Symbol == "" || req.Interval == "" {
		return nil, ErrorInvalidParam
	}

	client := mygateapi.NewRestClient("", "").PublicRestClient()
	switch GateAccountType(req.AccountType) {
	case GATE_ACCOUNT_TYPE_SPOT:
		api := client.NewPublicRestSpotCandlesticks().CurrencyPair(req.Symbol).Interval(req.Interval)

		if req.StartTime != 0 {
			api.From(req.StartTime)
		}
		if req.EndTime != 0 {
			api.To(req.EndTime)
		}
		if req.Limit != 0 {
			api.Limit(req.Limit)
		}
		data, err := api.Do()
		if err != nil {
			return nil, err
		}
		return m.convertToSpotKline(req.AccountType, req.Symbol, req.Interval, &data.Data), nil
	case GATE_ACCOUNT_TYPE_FUTURES:
		api := client.NewPublicRestFuturesSettleCandlesticks().Settle("usdt").Contract(req.Symbol).Interval(req.Interval)
		if req.StartTime != 0 {
			api.From(req.StartTime)
		}
		if req.EndTime != 0 {
			api.To(req.EndTime)
		}
		if req.Limit != 0 {
			api.Limit(req.Limit)
		}
		data, err := api.Do()
		if err != nil {
			return nil, err
		}
		return m.convertToFuturesKline(req.AccountType, req.Symbol, req.Interval, &data.Data), nil

	case GATE_ACCOUNT_TYPE_DELIVERY:
		api := client.NewPublicRestDeliverySettleCandlesticks().Settle("usdt").Contract(req.Symbol).Interval(req.Interval)
		if req.StartTime != 0 {
			api.From(req.StartTime)
		}
		if req.EndTime != 0 {
			api.To(req.EndTime)
		}
		if req.Limit != 0 {
			api.Limit(req.Limit)
		}
		data, err := api.Do()
		if err != nil {
			return nil, err
		}
		return m.convertToDeliveryKline(req.AccountType, req.Symbol, req.Interval, &data.Data), nil
	default:
		return nil, ErrorInvalidParam
	}
}
func (m *GateMarketData) GetBook(req *BookParam) (*OrderBook, error) {
	if req == nil || req.AccountType == "" || req.Symbol == "" {
		return nil, ErrorInvalidParam
	}

	client := mygateapi.NewRestClient("", "").PublicRestClient()
	switch GateAccountType(req.AccountType) {
	case GATE_ACCOUNT_TYPE_SPOT:
		api := client.NewPublicRestSpotOrderBook().CurrencyPair(req.Symbol)
		if req.Level != 0 {
			api.Limit(req.Level)
		}
		data, err := api.Do()
		if err != nil {
			return nil, err
		}
		return m.convertToSpotOrderBook(req.AccountType, req.Symbol, &data.Data), nil
	case GATE_ACCOUNT_TYPE_FUTURES:
		api := client.NewPublicRestFuturesSettleOrderBook().Settle("usdt").Contract(req.Symbol)
		if req.Level != 0 {
			api.Limit(req.Level)
		}
		data, err := api.Do()
		if err != nil {
			return nil, err
		}
		return m.convertToFuturesOrderBook(req.AccountType, req.Symbol, &data.Data), nil
	case GATE_ACCOUNT_TYPE_DELIVERY:
		api := client.NewPublicRestDeliverySettleOrderBook().Settle("usdt").Contract(req.Symbol)
		if req.Level != 0 {
			api.Limit(req.Level)
		}
		data, err := api.Do()
		if err != nil {
			return nil, err
		}
		return m.convertToDeliveryOrderBook(req.AccountType, req.Symbol, &data.Data), nil
	default:
		return nil, ErrorAccountType
	}
}

func (m *GateMarketData) convertToSpotKline(accountType, symbol, interval string, data *mygateapi.PublicRestSpotCandlesticksRes) *[]Kline {
	var list []Kline
	if data == nil || len(*data) == 0 {
		return &list
	}
	for _, v := range *data {
		startTime, err := strconv.ParseInt(v.Ts, 10, 64)
		if err != nil {
			continue
		}
		kline := Kline{
			Exchange:             GATE_NAME.String(),
			AccountType:          accountType,
			Symbol:               symbol,
			Interval:             interval,
			StartTime:            startTime * 1000,
			Open:                 stringToFloat64(v.O),
			High:                 stringToFloat64(v.H),
			Low:                  stringToFloat64(v.L),
			Close:                stringToFloat64(v.C),
			Volume:               stringToFloat64(v.VolCcy),
			CloseTime:            gateGetKlineCloseTime(startTime, interval),
			TransactionVolume:    stringToFloat64(v.VolCcyQuote),
			TransactionNumber:    0,
			BuyTransactionVolume: 0,
			BuyTransactionAmount: 0,
		}
		list = append(list, kline)
	}
	return &list
}

func (m *GateMarketData) convertToFuturesKline(accountType, symbol, interval string, data *mygateapi.PublicRestFuturesSettleCandlesticksRes) *[]Kline {
	var list []Kline
	if data == nil || len(*data) == 0 {
		return &list
	}
	for _, v := range *data {
		startTime := time.Unix(v.Ts, 0).UnixMilli()
		kline := Kline{
			Exchange:             GATE_NAME.String(),
			AccountType:          accountType,
			Symbol:               symbol,
			Interval:             interval,
			StartTime:            startTime,
			Open:                 stringToFloat64(v.O),
			High:                 stringToFloat64(v.H),
			Low:                  stringToFloat64(v.L),
			Close:                stringToFloat64(v.C),
			Volume:               float64(v.Vol),
			CloseTime:            gateGetKlineCloseTime(startTime, interval),
			TransactionVolume:    stringToFloat64(v.VolCcyQuote),
			TransactionNumber:    0,
			BuyTransactionVolume: 0,
			BuyTransactionAmount: 0,
		}
		list = append(list, kline)
	}
	return &list
}

func (m *GateMarketData) convertToDeliveryKline(accountType, symbol, interval string, data *mygateapi.PublicRestDeliverySettleCandlesticksRes) *[]Kline {
	var list []Kline
	if data == nil || len(*data) == 0 {
		return &list
	}
	for _, v := range *data {
		startTime := time.Unix(v.Timestamp, 0).UnixMilli()
		kline := Kline{
			Exchange:             GATE_NAME.String(),
			AccountType:          accountType,
			Symbol:               symbol,
			Interval:             interval,
			StartTime:            startTime,
			Open:                 stringToFloat64(v.Open),
			High:                 stringToFloat64(v.High),
			Low:                  stringToFloat64(v.Low),
			Close:                stringToFloat64(v.Close),
			Volume:               float64(v.Volume),
			CloseTime:            gateGetKlineCloseTime(startTime, interval),
			TransactionVolume:    0,
			TransactionNumber:    0,
			BuyTransactionVolume: 0,
			BuyTransactionAmount: 0,
		}
		list = append(list, kline)
	}
	return &list
}

func (m *GateMarketData) convertToSpotOrderBook(accountType, symbol string, data *mygateapi.PublicRestSpotOrderBookRes) *OrderBook {
	var ob OrderBook
	ob.Exchange = GATE_NAME.String()
	ob.AccountType = accountType
	ob.Symbol = symbol

	ob.Timestamp = data.Update
	ob.Asks = make([]Book, len(data.Asks))
	ob.Bids = make([]Book, len(data.Bids))
	for i, v := range data.Asks {
		ob.Asks[i] = Book{
			Price:    v[0],
			Quantity: v[1],
		}
	}
	for i, v := range data.Bids {
		ob.Bids[i] = Book{
			Price:    v[0],
			Quantity: v[1],
		}
	}
	return &ob
}

func (m *GateMarketData) convertToFuturesOrderBook(accountType, symbol string, data *mygateapi.PublicRestFuturesSettleOrderBookRes) *OrderBook {
	var ob OrderBook
	ob.Exchange = GATE_NAME.String()
	ob.AccountType = accountType
	ob.Symbol = symbol

	ob.Timestamp = int64(data.Update * 1000)
	ob.Asks = make([]Book, len(data.Asks))
	ob.Bids = make([]Book, len(data.Bids))
	for i, v := range data.Asks {
		ob.Asks[i] = Book{
			Price:    v.P,
			Quantity: strconv.FormatInt(v.S, 10),
		}
	}
	for i, v := range data.Bids {
		ob.Bids[i] = Book{
			Price:    v.P,
			Quantity: strconv.FormatInt(v.S, 10),
		}
	}
	return &ob
}

func (m *GateMarketData) convertToDeliveryOrderBook(accountType, symbol string, data *mygateapi.PublicRestDeliverySettleOrderBookRes) *OrderBook {
	var ob OrderBook
	ob.Exchange = GATE_NAME.String()
	ob.AccountType = accountType
	ob.Symbol = symbol

	ob.Timestamp = int64(data.Update * 1000)
	ob.Asks = make([]Book, len(data.Asks))
	ob.Bids = make([]Book, len(data.Bids))
	for i, v := range data.Asks {
		ob.Asks[i] = Book{
			Price:    v.Price,
			Quantity: strconv.FormatInt(v.Quantity, 10),
		}
	}
	for i, v := range data.Bids {
		ob.Bids[i] = Book{
			Price:    v.Price,
			Quantity: strconv.FormatInt(v.Quantity, 10),
		}
	}
	return &ob
}
