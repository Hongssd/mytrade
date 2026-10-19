package mytrade

import (
	"errors"
	"fmt"

	"github.com/Hongssd/mybinanceapi"
)

// binanceStreamRecoverer is what every user stream client exposes.
type binanceStreamRecoverer interface {
	ForceReconnect(reason string) error
	RenewListenKeyAndReconnect(reason string) error
}

// binanceOrderStreams lists the already opened user streams that carry orders
// for accountType. It never opens a new one.
func (b *BinanceTradeEngine) binanceOrderStreams(accountType string) []binanceStreamRecoverer {
	var streams []binanceStreamRecoverer
	switch BinanceAccountType(accountType) {
	case BN_AC_SPOT:
		if b.wsSpotAccount != nil {
			streams = append(streams, b.wsSpotAccount)
		}
		if b.wsSpotMarginAccount != nil {
			streams = append(streams, b.wsSpotMarginAccount)
		}
		b.wsSpotIsolatedMarginAccount.Range(func(_ string, ws *mybinanceapi.SpotWsStreamClient) bool {
			streams = append(streams, ws)
			return true
		})
		if b.wsPMMarginAccount != nil {
			streams = append(streams, b.wsPMMarginAccount)
		}
	case BN_AC_FUTURE:
		if b.isPortfolioMargin {
			if b.wsPMContractAccount != nil {
				streams = append(streams, b.wsPMContractAccount)
			}
		} else if b.wsFutureAccount != nil {
			streams = append(streams, b.wsFutureAccount)
		}
	case BN_AC_SWAP:
		if b.isPortfolioMargin {
			if b.wsPMContractAccount != nil {
				streams = append(streams, b.wsPMContractAccount)
			}
		} else if b.wsSwapAccount != nil {
			streams = append(streams, b.wsSwapAccount)
		}
	}
	return streams
}

// RecoverOrderStream repairs the user stream behind the orders push. The stream
// client keeps its subscriptions, so existing SubscribeOrder results keep flowing.
func (b *BinanceTradeEngine) RecoverOrderStream(accountType string, level OrderStreamRecoverLevel, reason string) error {
	streams := b.binanceOrderStreams(accountType)
	if len(streams) == 0 {
		return fmt.Errorf("binance %s orders stream is not open", accountType)
	}
	var errs []error
	for _, ws := range streams {
		var err error
		if level == OrderStreamRenew {
			err = ws.RenewListenKeyAndReconnect(reason)
		} else {
			err = ws.ForceReconnect(reason)
		}
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
