package mytrade

import "github.com/shopspring/decimal"

type AssetTransferParams struct {
	// required
	Asset  string          //币种 All required
	Amount decimal.Decimal //数量 All required
	From   AssetType       //从哪个账户划转 All required
	To     AssetType       //划转到哪个账户 All required
}
