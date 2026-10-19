package mytrade

import "errors"

// OrderStreamRecoverLevel is how hard a recovery tries.
type OrderStreamRecoverLevel int

const (
	// Replace the connection. A listenKey stream also re-POSTs the key first.
	OrderStreamReconnect OrderStreamRecoverLevel = iota + 1
	// Delete the listenKey too so the exchange issues a brand-new one. Use it
	// only after a reconnect did not bring pushes back, because a key that
	// belongs to the account also drops other connections sharing it.
	OrderStreamRenew
)

func (l OrderStreamRecoverLevel) String() string {
	switch l {
	case OrderStreamReconnect:
		return "reconnect"
	case OrderStreamRenew:
		return "renew"
	default:
		return "unknown"
	}
}

// OrderStreamRecoverer is implemented by engines whose orders push can go
// silent without the connection closing. A caller that sees REST confirm orders
// the push never reported uses it to repair the stream.
type OrderStreamRecoverer interface {
	RecoverOrderStream(accountType string, level OrderStreamRecoverLevel, reason string) error
}

// RecoverOrderStream repairs the orders push of one account. It returns
// ErrorNotSupport for an exchange that cannot be repaired this way.
func (e *ExchangeManager) RecoverOrderStream(api ExchangeApiParam, accountType string, level OrderStreamRecoverLevel, reason string) error {
	tradeEngine, err := e.getTradeEngine(api.Exchange, api.ApiKey, api.ApiSecret, api.Passphrase)
	if err != nil {
		return err
	}
	recoverer, ok := tradeEngine.(OrderStreamRecoverer)
	if !ok {
		return errors.Join(errors.New(api.Exchange), ErrorNotSupport)
	}
	return recoverer.RecoverOrderStream(accountType, level, reason)
}
