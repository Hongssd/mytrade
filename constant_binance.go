package mytrade

type BinanceAccountType string

func (b BinanceAccountType) String() string {
	return string(b)
}

const (
	BN_AC_SPOT   BinanceAccountType = "SPOT"
	BN_AC_FUTURE BinanceAccountType = "FUTURE"
	BN_AC_SWAP   BinanceAccountType = "SWAP"
)
