package mytrade

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestOkxSymbolSettleCcy(t *testing.T) {
	t.Parallel()
	tests := []struct {
		symbol string
		want   string
	}{
		{"SNDK-USDT-SWAP", "USDT"},
		{"BTC-USD-SWAP", "USD"},
		{"BTC-USDC-SWAP", "USDC"},
		{"BTC-USDT", "USDT"},
		{"BTC-USDT-250328", "USDT"},
		{"ETH-USD-250328", "USD"},
		{"BTC", ""},
		{"", ""},
	}
	for _, tc := range tests {
		if got := okxSymbolSettleCcy(tc.symbol); got != tc.want {
			t.Fatalf("okxSymbolSettleCcy(%q)=%q want %q", tc.symbol, got, tc.want)
		}
	}
}

func TestOkxTradeFeeInstFamily(t *testing.T) {
	t.Parallel()
	tests := []struct {
		accountType string
		symbol      string
		want        string
	}{
		{OKX_AC_SWAP.String(), "SNDK-USDT-SWAP", "SNDK-USDT"},
		{OKX_AC_SWAP.String(), "BTC-USD-SWAP", "BTC-USD"},
		{OKX_AC_FUTURES.String(), "BTC-USDT-250328", "BTC-USDT"},
		{OKX_AC_OPTION.String(), "BTC-USD-241227-100000-C", "BTC-USD"},
		{OKX_AC_SPOT.String(), "BTC-USDT", ""},
		{OKX_AC_MARGIN.String(), "BTC-USDT", ""},
		{OKX_AC_SWAP.String(), "BTC", ""},
		{OKX_AC_SWAP.String(), "", ""},
	}
	for _, tc := range tests {
		if got := okxTradeFeeInstFamily(tc.accountType, tc.symbol); got != tc.want {
			t.Fatalf("okxTradeFeeInstFamily(%q,%q)=%q want %q",
				tc.accountType, tc.symbol, got, tc.want)
		}
	}
}

func TestOkxFeeRateFromTradeFeeQuote(t *testing.T) {
	t.Parallel()
	// VIP2 coin-margined vs actual USDT-margined (dashboard: maker rebate 0.0020%, taker 0.0200%).
	q := okxAccountTradeFeeQuote{
		Maker:     "-0.00015",
		Taker:     "-0.00036",
		MakerU:    "0.00002",
		TakerU:    "-0.0002",
		MakerUSDC: "-0.0001",
		TakerUSDC: "-0.00025",
	}
	mustDec := func(s string) decimal.Decimal {
		d, err := decimal.NewFromString(s)
		if err != nil {
			t.Fatalf("decimal %q: %v", s, err)
		}
		return d
	}
	tests := []struct {
		name        string
		accountType string
		symbol      string
		wantMaker   string
		wantTaker   string
	}{
		{
			name:        "usdt_swap_uses_makerU_not_vip_coin_margin",
			accountType: OKX_AC_SWAP.String(),
			symbol:      "SNDK-USDT-SWAP",
			wantMaker:   "-0.00002",
			wantTaker:   "0.0002",
		},
		{
			name:        "usdt_futures_uses_makerU",
			accountType: OKX_AC_FUTURES.String(),
			symbol:      "BTC-USDT-250328",
			wantMaker:   "-0.00002",
			wantTaker:   "0.0002",
		},
		{
			name:        "coin_margin_swap_uses_maker_taker",
			accountType: OKX_AC_SWAP.String(),
			symbol:      "BTC-USD-SWAP",
			wantMaker:   "0.00015",
			wantTaker:   "0.00036",
		},
		{
			name:        "usdc_swap_uses_makerUSDC",
			accountType: OKX_AC_SWAP.String(),
			symbol:      "BTC-USDC-SWAP",
			wantMaker:   "0.0001",
			wantTaker:   "0.00025",
		},
		{
			name:        "spot_usdt_uses_maker_taker",
			accountType: OKX_AC_SPOT.String(),
			symbol:      "BTC-USDT",
			wantMaker:   "0.00015",
			wantTaker:   "0.00036",
		},
		{
			name:        "spot_usdc_uses_makerUSDC",
			accountType: OKX_AC_SPOT.String(),
			symbol:      "BTC-USDC",
			wantMaker:   "0.0001",
			wantTaker:   "0.00025",
		},
		{
			name:        "option_uses_maker_taker",
			accountType: OKX_AC_OPTION.String(),
			symbol:      "BTC-USD-241227-100000-C",
			wantMaker:   "0.00015",
			wantTaker:   "0.00036",
		},
		{
			name:        "empty_symbol_swap_falls_to_coin_margin_fields",
			accountType: OKX_AC_SWAP.String(),
			symbol:      "",
			wantMaker:   "0.00015",
			wantTaker:   "0.00036",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := okxFeeRateFromTradeFeeQuote(tc.accountType, tc.symbol, q)
			if !got.Maker.Equal(mustDec(tc.wantMaker)) || !got.Taker.Equal(mustDec(tc.wantTaker)) {
				t.Fatalf("fee=%s/%s want %s/%s",
					got.Maker.String(), got.Taker.String(), tc.wantMaker, tc.wantTaker)
			}
		})
	}
}

func TestOkxPickTradeFeeStrings_EmptyQuote(t *testing.T) {
	t.Parallel()
	maker, taker := okxPickTradeFeeStrings(OKX_AC_SWAP.String(), "SNDK-USDT-SWAP", okxAccountTradeFeeQuote{})
	if maker != "" || taker != "" {
		t.Fatalf("empty quote got %q/%q", maker, taker)
	}
	got := okxFeeRateFromTradeFeeQuote(OKX_AC_SWAP.String(), "SNDK-USDT-SWAP", okxAccountTradeFeeQuote{})
	if !got.Maker.IsZero() || !got.Taker.IsZero() {
		t.Fatalf("empty quote parsed fee=%s/%s", got.Maker.String(), got.Taker.String())
	}
}
