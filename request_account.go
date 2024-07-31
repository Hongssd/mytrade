package mytrade

import "github.com/shopspring/decimal"

type AssetTransferParams struct {
	// required
	Exchange   string          //交易所名称 All required
	TransferId string          //划转ID ByBit required
	Asset      string          //币种 All required
	Amount     decimal.Decimal //数量 All required
	From       string          //从哪个账户划转 All required
	To         string          //划转到哪个账户 All required
	FromSymbol string          //从哪个合约划转 BN optional
	ToSymbol   string          //划转到哪个合约 BN optional
	Timestamp  int64           //时间戳 BN required

	// optional
	RecvWindow  int64  // 超时时间（应小于60s，单位ms）BN optional
	SubAcct     string // 子账户名 OKX optional
	LoanTrans   bool   // 是否支持跨币种保证金模式或组合保证金模式下的借币转出，默认false不借出 OKX optional
	OmitPosRisk string // 是否忽略仓位风险，默认为false，仅适用于组合保证金模式 OKX optional
	ClientId    string // 客户自定义ID 字母（区分大小写）与数字的组合，可以是纯字母、纯数字且长度要在1-32位之间。OKX optional
}
