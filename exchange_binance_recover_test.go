package mytrade

import (
	"errors"
	"testing"
	"time"

	"github.com/Hongssd/mybinanceapi"
)

func TestBinanceOrderStreamsPicksStreamByAccount(t *testing.T) {
	bn := mybinanceapi.MyBinance{}

	pm := &BinanceTradeEngine{
		isPortfolioMargin:   true,
		wsPMContractAccount: bn.NewPMContractStreamClient(),
	}
	for _, ac := range []BinanceAccountType{BN_AC_FUTURE, BN_AC_SWAP} {
		if n := len(pm.binanceOrderStreams(ac.String())); n != 1 {
			t.Fatalf("portfolio margin %s should use the PM contract stream, got %d", ac, n)
		}
	}

	plain := &BinanceTradeEngine{wsFutureAccount: bn.NewFutureWsStreamClient()}
	if n := len(plain.binanceOrderStreams(BN_AC_FUTURE.String())); n != 1 {
		t.Fatalf("future stream missing, got %d", n)
	}
	if n := len(plain.binanceOrderStreams(BN_AC_SWAP.String())); n != 0 {
		t.Fatalf("an unopened swap stream must not be created, got %d", n)
	}
}

func TestBinanceRecoverOrderStreamNotOpen(t *testing.T) {
	engine := &BinanceTradeEngine{}
	if err := engine.RecoverOrderStream(BN_AC_FUTURE.String(), OrderStreamReconnect, "test"); err == nil {
		t.Fatal("recovering a stream that was never opened should fail")
	}
}

func TestPushBinanceOrderSubDropsWhenStalled(t *testing.T) {
	old := binanceOrderSubStallLimit
	binanceOrderSubStallLimit = 30 * time.Millisecond
	defer func() { binanceOrderSubStallLimit = old }()

	sub := &subscription[Order]{resultChan: make(chan Order, 1), errChan: make(chan error, 1)}
	pushBinanceOrderSub(sub, Order{ClientOrderId: "a"})

	start := time.Now()
	pushBinanceOrderSub(sub, Order{ClientOrderId: "b"})
	if time.Since(start) < 20*time.Millisecond {
		t.Fatal("a full channel should be waited on before dropping")
	}
	if got := (<-sub.resultChan).ClientOrderId; got != "a" {
		t.Fatalf("first push must survive, got %s", got)
	}
	select {
	case o := <-sub.resultChan:
		t.Fatalf("stalled push should be dropped, got %s", o.ClientOrderId)
	default:
	}

	pushBinanceOrderSubErr(sub, errors.New("boom"))
	pushBinanceOrderSubErr(sub, errors.New("boom")) // full: must not block
}
