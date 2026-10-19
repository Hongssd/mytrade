package mytrade

import "time"

// binanceOrderSubStallLimit bounds how long one stalled consumer may hold the
// shared stream reader. Past it the push is dropped, so every other order of the
// account keeps flowing.
var binanceOrderSubStallLimit = 2 * time.Second

// pushBinanceOrderSub never blocks the stream reader on a stalled consumer.
func pushBinanceOrderSub(sub *subscription[Order], order Order) {
	select {
	case sub.resultChan <- order:
		return
	default:
	}
	timer := time.NewTimer(binanceOrderSubStallLimit)
	defer timer.Stop()
	select {
	case sub.resultChan <- order:
	case <-timer.C:
		log.Errorf("[binance-orders] subscriber stalled, push dropped: clientOrderId=%s status=%v", order.ClientOrderId, order.Status)
	}
}

func pushBinanceOrderSubErr(sub *subscription[Order], err error) {
	select {
	case sub.errChan <- err:
	default:
		log.Warnf("[binance-orders] subscriber error chan full, dropped: %v", err)
	}
}
