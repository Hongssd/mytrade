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

const (
	BN_ORDER_TYPE_LIMIT  = "LIMIT"
	BN_ORDER_TYPE_MARKET = "MARKET"
)

const (
	BN_ORDER_SIDE_BUY  = "BUY"
	BN_ORDER_SIDE_SELL = "SELL"
)

const (
	BN_POSITION_SIDE_LONG  = "LONG"
	BN_POSITION_SIDE_SHORT = "SHORT"
)

const (
	BN_ORDER_STATUS_NEW              = "NEW"
	BN_ORDER_STATUS_PARTIALLY_FILLED = "PARTIALLY_FILLED"
	BN_ORDER_STATUS_FILLED           = "FILLED"
	BN_ORDER_STATUS_CANCELED         = "CANCELED"
	BN_ORDER_STATUS_REJECTED         = "REJECTED"
)

const (
	BN_TIME_IN_FORCE_GTC       = "GTC"
	BN_TIME_IN_FORCE_IOC       = "IOC"
	BN_TIME_IN_FORCE_FOK       = "FOK"
	BN_TIME_IN_FORCE_POST_ONLY = "GTX"
)
