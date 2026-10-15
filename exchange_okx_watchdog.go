package mytrade

import (
	"sync"
	"time"
)

type okxRecoverLevel int

const (
	okxRecoverNone okxRecoverLevel = iota
	// Drop the connection and let myokxapi reconnect, login and resubscribe.
	okxRecoverReconnect
	// Build a brand-new ws client and subscription and swap it in.
	okxRecoverRebuild
)

func (l okxRecoverLevel) String() string {
	switch l {
	case okxRecoverReconnect:
		return "reconnect"
	case okxRecoverRebuild:
		return "rebuild"
	default:
		return "none"
	}
}

const (
	// REST accepted the request but no push arrived for this many requests in a row.
	okxAckTimeoutThreshold = 3
	// Minimum gap between two recoveries, so a burst of timed-out waiters
	// (a batch times out together) triggers one recovery.
	okxRecoverCooldown = 5 * time.Second
	// A new recovery request within this window means the previous level did not help.
	okxRecoverEscalateWindow = 3 * time.Minute
	// No frame at all (not even a pong) for this long means myokxapi's own
	// reconnect is stuck.
	okxWsSilentLimit = 60 * time.Second
	// How often the health ticker inspects the connection.
	okxHealthInterval = 10 * time.Second
)

// okxAckHealth decides when the orders push channel needs recovery. It holds no
// connection state so the decision logic can be tested on its own.
type okxAckHealth struct {
	mu            sync.Mutex
	strikes       int
	lastRecoverAt time.Time
	lastLevel     okxRecoverLevel
}

// onSuccess records a push that reached a waiter. Any success proves the
// channel works, so strikes and the escalation level start over.
func (h *okxAckHealth) onSuccess() {
	h.mu.Lock()
	h.strikes = 0
	h.lastLevel = okxRecoverNone
	h.mu.Unlock()
}

// onTimeout records a waiter whose push never arrived and reports whether the
// caller should start a recovery now.
func (h *okxAckHealth) onTimeout(now time.Time) (okxRecoverLevel, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.strikes++
	if h.strikes < okxAckTimeoutThreshold {
		return okxRecoverNone, false
	}
	return h.nextLevelLocked(now)
}

// onSilent records a connection that delivered no frame for too long.
func (h *okxAckHealth) onSilent(now time.Time) (okxRecoverLevel, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.nextLevelLocked(now)
}

func (h *okxAckHealth) nextLevelLocked(now time.Time) (okxRecoverLevel, bool) {
	if !h.lastRecoverAt.IsZero() && now.Sub(h.lastRecoverAt) < okxRecoverCooldown {
		return okxRecoverNone, false
	}
	level := okxRecoverReconnect
	if h.lastLevel != okxRecoverNone && now.Sub(h.lastRecoverAt) < okxRecoverEscalateWindow {
		level = okxRecoverRebuild
	}
	h.strikes = 0
	h.lastRecoverAt = now
	h.lastLevel = level
	return level, true
}

func (h *okxAckHealth) strikeCount() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.strikes
}
