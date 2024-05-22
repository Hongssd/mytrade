package mytrade

import (
	"github.com/shopspring/decimal"
)

type symbolInfoStruct struct {
	//基础规范
	Exchange      string `json:"exchange"`      //交易所名字
	AccountType   string `json:"accountType"`   //账户类型
	Symbol        string `json:"symbol"`        //交易对名字
	BaseCoin      string `json:"baseCoin"`      //交易币种
	QuoteCoin     string `json:"quoteCoin"`     //计价币种
	IsTrading     bool   `json:"isTrading"`     //是否交易中
	IsContract    bool   `json:"isContract"`    //是否合约
	IsContractAmt bool   `json:"isContractAmt"` //是否合约张数计量
	ContractSize  string `json:"contractSize"`  //合约面值
	ContractCoin  string `json:"contractCoin"`  //合约面值计价币种
	ContractType  string `json:"contractType"`  //合约类型

	//精度
	PricePrecision int `json:"pricePrecision"` //价格精度
	AmtPrecision   int `json:"amtPrecision"`   //数量精度

	//价格规范
	TickSize string `json:"tickSize"` //下单价格精度
	MinPrice string `json:"minPrice"` //最小下单价格
	MaxPrice string `json:"maxPrice"` //最大下单价格

	//数量规范
	LotSize   string `json:"lotSize"`   //下单数量精度
	MinAmt    string `json:"minAmt"`    //最小下单数量
	MaxLmtAmt string `json:"maxLmtAmt"` //最大限价单下单数量
	MaxMktAmt string `json:"maxMktAmt"` //最大市价单下单数量

	//其他规范
	MaxLeverage  string `json:"maxLeverage"`  //最大杠杆
	MinLeverage  string `json:"minLeverage"`  //最小杠杆
	StepLeverage string `json:"stepLeverage"` //修改杠杆的步长
	MaxOrderNum  int    `json:"maxOrderNum"`  //最大订单数
	MinNotional  string `json:"minNotional"`  //最小名义价值
}

type symbolInfo struct {
	symbolInfoStruct
}

// 基础规范
// 交易所名字
func (info *symbolInfo) Exchange() string {
	return info.symbolInfoStruct.Exchange
}

// 账户类型
func (info *symbolInfo) AccountType() string {
	return info.symbolInfoStruct.AccountType
}

// 交易对名字
func (info *symbolInfo) Symbol() string {
	return info.symbolInfoStruct.Symbol
}

// 交易币种
func (info *symbolInfo) BaseCoin() string {
	return info.symbolInfoStruct.BaseCoin
}

// 计价币种
func (info *symbolInfo) QuoteCoin() string {
	return info.symbolInfoStruct.QuoteCoin
}

// 是否交易中
func (info *symbolInfo) IsTrading() bool {
	return info.symbolInfoStruct.IsTrading
}

// 是否合约
func (info *symbolInfo) IsContract() bool {
	return info.symbolInfoStruct.IsContract
}

// 是否合约张数计量
func (info *symbolInfo) IsContractAmt() bool {
	return info.symbolInfoStruct.IsContractAmt
}

// 合约面值
func (info *symbolInfo) ContractSize() decimal.Decimal {
	return decimal.RequireFromString(info.symbolInfoStruct.ContractSize)
}

// 合约面值计价币种
func (info *symbolInfo) ContractCoin() string {
	return info.symbolInfoStruct.ContractCoin
}

// 合约类型
func (info *symbolInfo) ContractType() string {
	return info.symbolInfoStruct.ContractType
}

// 精度
// 价格精度
func (info *symbolInfo) PricePrecision() int {
	return info.symbolInfoStruct.PricePrecision
}

// 数量精度
func (info *symbolInfo) AmtPrecision() int {
	return info.symbolInfoStruct.AmtPrecision
}

// 价格规范
// 下单价格精度
func (info *symbolInfo) TickSize() decimal.Decimal {
	return decimal.RequireFromString(info.symbolInfoStruct.TickSize)
}

// 最小下单价格
func (info *symbolInfo) MinPrice() decimal.Decimal {
	return decimal.RequireFromString(info.symbolInfoStruct.MinPrice)
}

// 最大下单价格
func (info *symbolInfo) MaxPrice() decimal.Decimal {
	return decimal.RequireFromString(info.symbolInfoStruct.MaxPrice)
}

// 数量规范
// 当合约张数计量时为张的数量
// 下单数量精度
func (info *symbolInfo) LotSize() decimal.Decimal {
	return decimal.RequireFromString(info.symbolInfoStruct.LotSize)
}

// 最小下单数量
func (info *symbolInfo) MinAmt() decimal.Decimal {
	return decimal.RequireFromString(info.symbolInfoStruct.MinAmt)
}

// 最大限价单下单数量
func (info *symbolInfo) MaxLmtAmt() decimal.Decimal {
	return decimal.RequireFromString(info.symbolInfoStruct.MaxLmtAmt)
}

// 最大市价单下单数量
func (info *symbolInfo) MaxMktAmt() decimal.Decimal {
	return decimal.RequireFromString(info.symbolInfoStruct.MaxMktAmt)
}

// 其他规范
// 最大杠杆
func (info *symbolInfo) MaxLeverage() decimal.Decimal {
	return decimal.RequireFromString(info.symbolInfoStruct.MaxLeverage)
}

// 最小杠杆
func (info *symbolInfo) MinLeverage() decimal.Decimal {
	return decimal.RequireFromString(info.symbolInfoStruct.MinLeverage)
}

// 修改杠杆的步长
func (info *symbolInfo) StepLeverage() decimal.Decimal {
	return decimal.RequireFromString(info.symbolInfoStruct.StepLeverage)
}

// 最大订单数
func (info *symbolInfo) MaxOrderNum() int {
	return info.symbolInfoStruct.MaxOrderNum
}

// 最小名义价值
func (info *symbolInfo) MinNotional() decimal.Decimal {
	return decimal.RequireFromString(info.symbolInfoStruct.MinNotional)
}

func (info *symbolInfo) MarshalJsonIndent(prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(info.symbolInfoStruct, prefix, indent)
}

func (info *symbolInfo) MarshalJson() ([]byte, error) {
	return json.Marshal(info.symbolInfoStruct)
}
