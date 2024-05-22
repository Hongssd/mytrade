package mytrade

import (
	"errors"
	"github.com/shopspring/decimal"
	"sort"
)

type Kline struct {
	Exchange             string  `json:"exchange" bson:"exchange"`                             //交易所
	AccountType          string  `json:"account_type" bson:"account_type"`                     //K线类型
	Symbol               string  `json:"symbol" bson:"symbol"`                                 //交易对
	Interval             string  `json:"interval" bson:"interval"`                             //K线间隔
	StartTime            int64   `json:"start_time" bson:"start_time" gorm:"primaryKey"`       //开盘时间
	Open                 float64 `json:"open" bson:"open"`                                     //开盘价
	High                 float64 `json:"high" bson:"high"`                                     //最高价
	Low                  float64 `json:"low" bson:"low"`                                       //最低价
	Close                float64 `json:"close" bson:"close"`                                   //收盘价
	Volume               float64 `json:"volume" bson:"volume"`                                 //成交量
	CloseTime            int64   `json:"close_time" bson:"close_time"`                         //收盘时间
	TransactionVolume    float64 `json:"transaction_volume" bson:"transaction_volume"`         //成交额
	TransactionNumber    int64   `json:"transaction_number" bson:"transaction_number"`         //成交笔数
	BuyTransactionVolume float64 `json:"buy_transaction_volume" bson:"buy_transaction_volume"` //主动买入成交量
	BuyTransactionAmount float64 `json:"buy_transaction_amount" bson:"buy_transaction_amount"` //主动买入成交额
}

type KlinesDesc []Kline
type KlinesAsc []Kline

func (k KlinesDesc) Len() int {
	return len(k)
}
func (k KlinesDesc) Less(i, j int) bool {
	return k[i].StartTime > k[j].StartTime
}
func (k KlinesDesc) Swap(i, j int) {
	k[i], k[j] = k[j], k[i]
}

func (k KlinesAsc) Len() int {
	return len(k)
}
func (k KlinesAsc) Less(i, j int) bool {
	return k[i].StartTime < k[j].StartTime
}
func (k KlinesAsc) Swap(i, j int) {
	k[i], k[j] = k[j], k[i]
}

func SortKline(klines []Kline, desc bool) []Kline {
	if desc {
		kd := KlinesDesc(klines)
		sort.Sort(kd)
		return kd
	}
	ka := KlinesAsc(klines)
	sort.Sort(ka)
	return ka
}

type OrderBook struct {
	Exchange    string `json:"exchange" bson:"exchange"`         //交易所
	AccountType string `json:"account_type" bson:"account_type"` //K线类型
	Symbol      string `json:"symbol" bson:"symbol"`             //交易对
	Timestamp   int64  `json:"timestamp" bson:"timestamp"`       //时间戳
	Asks        []Book `json:"asks" bson:"asks"`                 //卖单
	Bids        []Book `json:"bids" bson:"bids"`                 //买单
}

type Book struct {
	Price    string
	Quantity string
}

func (b *Book) Float64() (float64, float64) {
	return stringToFloat64(b.Price), stringToFloat64(b.Quantity)
}
func (b *Book) Decimal() (decimal.Decimal, decimal.Decimal) {
	return stringToDecimal(b.Price), stringToDecimal(b.Quantity)
}

func (ob *OrderBook) WeightAvgBidsPrice(level int) (decimal.Decimal, decimal.Decimal, error) {
	return ob.weightAvgPrice(ob.Bids, level)
}

func (ob *OrderBook) WeightAvgAsksPrice(level int) (decimal.Decimal, decimal.Decimal, error) {
	return ob.weightAvgPrice(ob.Asks, level)
}

func (ob *OrderBook) weightAvgPrice(books []Book, level int) (decimal.Decimal, decimal.Decimal, error) {
	if level > len(books) {
		return decimal.Zero, decimal.Zero, errors.New("level is too large")
	}
	var total decimal.Decimal
	var totalQuantity decimal.Decimal
	for i := 0; i < level; i++ {
		price, quantity := books[i].Decimal()
		total = total.Add(price.Mul(quantity))
		totalQuantity = totalQuantity.Add(quantity)
	}
	return total.Div(totalQuantity), totalQuantity, nil
}
