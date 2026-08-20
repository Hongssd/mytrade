package mytrade

import (
	"strings"

	"github.com/shopspring/decimal"
)

// okxAccountTradeFeeQuote is the fee-rate payload from GET /api/v5/account/trade-fee.
// Maker/Taker are coin-margined (or SPOT USDT zone); MakerU/TakerU are USDT-margined;
// MakerUSDC/TakerUSDC are USDC-margined (or SPOT USDⓈ&Crypto zone).
type okxAccountTradeFeeQuote struct {
	Maker     string
	Taker     string
	MakerU    string
	TakerU    string
	MakerUSDC string
	TakerUSDC string
}

func okxSymbolSettleCcy(symbol string) string {
	parts := strings.Split(symbol, "-")
	if len(parts) < 2 {
		return ""
	}
	return strings.ToUpper(parts[1])
}

func okxTradeFeeInstFamily(accountType, symbol string) string {
	switch OkxAccountType(accountType) {
	case OKX_AC_SWAP, OKX_AC_FUTURES, OKX_AC_OPTION:
		parts := strings.Split(symbol, "-")
		if len(parts) < 2 {
			return ""
		}
		return parts[0] + "-" + parts[1]
	default:
		return ""
	}
}

func okxPickTradeFeeStrings(accountType, symbol string, q okxAccountTradeFeeQuote) (maker, taker string) {
	settle := okxSymbolSettleCcy(symbol)
	switch OkxAccountType(accountType) {
	case OKX_AC_SPOT, OKX_AC_MARGIN:
		if settle == "USDC" {
			return q.MakerUSDC, q.TakerUSDC
		}
		return q.Maker, q.Taker
	case OKX_AC_SWAP, OKX_AC_FUTURES:
		switch settle {
		case "USDT":
			return q.MakerU, q.TakerU
		case "USDC":
			return q.MakerUSDC, q.TakerUSDC
		default:
			return q.Maker, q.Taker
		}
	default:
		return q.Maker, q.Taker
	}
}

func okxFeeRateFromTradeFeeQuote(accountType, symbol string, q okxAccountTradeFeeQuote) FeeRate {
	makerStr, takerStr := okxPickTradeFeeStrings(accountType, symbol, q)
	maker, _ := decimal.NewFromString(makerStr)
	taker, _ := decimal.NewFromString(takerStr)
	// OKX: charge is negative, rebate is positive. mytrade: pay is positive, rebate is negative.
	return FeeRate{
		Maker: maker.Mul(decimal.NewFromInt(-1)),
		Taker: taker.Mul(decimal.NewFromInt(-1)),
	}
}
