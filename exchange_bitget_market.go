package mytrade

import (
	"strconv"

	mybitgetapi "github.com/Hongssd/mybitgetapi"
)

type BitgetMarketData struct {
	ExchangeBase
}

func (m *BitgetMarketData) preCheckAccountType(accountType string) bool {
	switch accountType {
	case BITGET_AC_SPOT, BITGET_AC_MARGIN, BITGET_AC_USDT_FUTURES, BITGET_AC_COIN_FUTURES, BITGET_AC_USDC_FUTURES:
		return true
	default:
		return false
	}
}

func (m *BitgetMarketData) NewKlineReq() *KlineParam {
	return &KlineParam{}
}

func (m *BitgetMarketData) NewBookReq() *BookParam {
	return &BookParam{}
}

// GetKline / GetBook 统一使用 UTA 公共行情，与 BitgetExchangeInfo 的 instruments 数据源一致（无密钥的 MarketData 无法感知 isClassic）。
func (m *BitgetMarketData) GetKline(req *KlineParam) (*[]Kline, error) {
	if req == nil || req.AccountType == "" || req.Symbol == "" || req.Interval == "" {
		return nil, ErrorInvalidParam
	}
	client := mybitgetapi.NewRestClient("", "", "").PublicRestClient()
	switch req.AccountType {
	case BITGET_AC_SPOT, BITGET_AC_MARGIN, BITGET_AC_USDT_FUTURES, BITGET_AC_COIN_FUTURES, BITGET_AC_USDC_FUTURES:
		api := client.NewPublicRestUtaMarketCandles().
			Category(req.AccountType).
			Symbol(req.Symbol).
			Interval(req.Interval)
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

func (m *BitgetMarketData) convertToKline(accountType, symbol, interval string, rows *mybitgetapi.PublicRestUtaMarketCandlesRes) *[]Kline {
	var list []Kline
	if rows == nil || len(*rows) == 0 {
		return &list
	}
	for _, v := range *rows {
		startTime := stringToInt64(v.OpenTime)
		k := Kline{
			Exchange:             BITGET_NAME.String(),
			AccountType:          accountType,
			Symbol:               symbol,
			Interval:             interval,
			StartTime:            startTime,
			Open:                 stringToFloat64(v.OpenPrice),
			High:                 stringToFloat64(v.HighPrice),
			Low:                  stringToFloat64(v.LowPrice),
			Close:                stringToFloat64(v.ClosePrice),
			Volume:               stringToFloat64(v.FillQty),
			CloseTime:            0,
			TransactionVolume:    stringToFloat64(v.FillAmount),
			TransactionNumber:    0,
			BuyTransactionVolume: 0,
			BuyTransactionAmount: 0,
		}
		list = append(list, k)
	}
	return &list
}

func (m *BitgetMarketData) GetBook(req *BookParam) (*OrderBook, error) {
	if req == nil || req.AccountType == "" || req.Symbol == "" {
		return nil, ErrorInvalidParam
	}
	client := mybitgetapi.NewRestClient("", "", "").PublicRestClient()
	switch req.AccountType {
	case BITGET_AC_SPOT, BITGET_AC_MARGIN, BITGET_AC_USDT_FUTURES, BITGET_AC_COIN_FUTURES, BITGET_AC_USDC_FUTURES:
		api := client.NewPublicRestUtaMarketOrderBook().Category(req.AccountType).Symbol(req.Symbol)
		if req.Level != 0 {
			api.Limit(strconv.Itoa(req.Level))
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

func (m *BitgetMarketData) convertToOrderBook(accountType, symbol string, data *mybitgetapi.PublicRestUtaMarketOrderBookRes) *OrderBook {
	ob := OrderBook{
		Exchange:    BITGET_NAME.String(),
		AccountType: accountType,
		Symbol:      symbol,
	}
	if data == nil {
		return &ob
	}

	ob.Timestamp = stringToInt64(data.Ts)
	ob.Asks = make([]Book, len(data.Asks))
	for i, a := range data.Asks {
		ob.Asks[i] = Book{Price: a.Price, Quantity: a.Quantity}
	}
	ob.Bids = make([]Book, len(data.Bids))
	for i, b := range data.Bids {
		ob.Bids[i] = Book{Price: b.Price, Quantity: b.Quantity}
	}
	return &ob
}
