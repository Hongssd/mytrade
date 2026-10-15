package mytrade

import (
	"testing"
	"time"
)

func TestOkxAckHealthOnTimeout(t *testing.T) {
	base := time.Unix(1_700_000_000, 0)

	t.Run("below threshold does not recover", func(t *testing.T) {
		h := &okxAckHealth{}
		for i := 0; i < okxAckTimeoutThreshold-1; i++ {
			if _, ok := h.onTimeout(base); ok {
				t.Fatalf("recovered after %d timeouts", i+1)
			}
		}
	})

	t.Run("threshold triggers reconnect and resets strikes", func(t *testing.T) {
		h := &okxAckHealth{}
		var level okxRecoverLevel
		var ok bool
		for i := 0; i < okxAckTimeoutThreshold; i++ {
			level, ok = h.onTimeout(base)
		}
		if !ok || level != okxRecoverReconnect {
			t.Fatalf("got level=%v ok=%v, want reconnect", level, ok)
		}
		if h.strikeCount() != 0 {
			t.Fatalf("strikes not reset: %d", h.strikeCount())
		}
	})

	t.Run("success resets strikes", func(t *testing.T) {
		h := &okxAckHealth{}
		h.onTimeout(base)
		h.onTimeout(base)
		h.onSuccess()
		if _, ok := h.onTimeout(base); ok {
			t.Fatal("recovered although a success interrupted the streak")
		}
	})

	t.Run("cooldown suppresses a second recovery", func(t *testing.T) {
		h := &okxAckHealth{}
		for i := 0; i < okxAckTimeoutThreshold; i++ {
			h.onTimeout(base)
		}
		soon := base.Add(okxRecoverCooldown / 2)
		for i := 0; i < okxAckTimeoutThreshold+2; i++ {
			if _, ok := h.onTimeout(soon); ok {
				t.Fatal("recovered inside cooldown")
			}
		}
	})

	t.Run("repeat failure after cooldown escalates to rebuild", func(t *testing.T) {
		h := &okxAckHealth{}
		for i := 0; i < okxAckTimeoutThreshold; i++ {
			h.onTimeout(base)
		}
		later := base.Add(okxRecoverCooldown + time.Second)
		var level okxRecoverLevel
		var ok bool
		for i := 0; i < okxAckTimeoutThreshold; i++ {
			level, ok = h.onTimeout(later)
		}
		if !ok || level != okxRecoverRebuild {
			t.Fatalf("got level=%v ok=%v, want rebuild", level, ok)
		}
	})

	t.Run("failure after the escalate window starts over at reconnect", func(t *testing.T) {
		h := &okxAckHealth{}
		for i := 0; i < okxAckTimeoutThreshold; i++ {
			h.onTimeout(base)
		}
		later := base.Add(okxRecoverEscalateWindow + time.Second)
		var level okxRecoverLevel
		var ok bool
		for i := 0; i < okxAckTimeoutThreshold; i++ {
			level, ok = h.onTimeout(later)
		}
		if !ok || level != okxRecoverReconnect {
			t.Fatalf("got level=%v ok=%v, want reconnect", level, ok)
		}
	})

	t.Run("success clears escalation", func(t *testing.T) {
		h := &okxAckHealth{}
		for i := 0; i < okxAckTimeoutThreshold; i++ {
			h.onTimeout(base)
		}
		h.onSuccess()
		later := base.Add(okxRecoverCooldown + time.Second)
		var level okxRecoverLevel
		for i := 0; i < okxAckTimeoutThreshold; i++ {
			level, _ = h.onTimeout(later)
		}
		if level != okxRecoverReconnect {
			t.Fatalf("got level=%v, want reconnect", level)
		}
	})
}

func TestOkxAckHealthOnSilent(t *testing.T) {
	base := time.Unix(1_700_000_000, 0)
	h := &okxAckHealth{}

	level, ok := h.onSilent(base)
	if !ok || level != okxRecoverReconnect {
		t.Fatalf("first silent: level=%v ok=%v, want reconnect", level, ok)
	}
	if _, ok := h.onSilent(base.Add(time.Second)); ok {
		t.Fatal("silent recovery inside cooldown")
	}
	level, ok = h.onSilent(base.Add(okxHealthInterval + okxRecoverCooldown))
	if !ok || level != okxRecoverRebuild {
		t.Fatalf("second silent: level=%v ok=%v, want rebuild", level, ok)
	}
}
