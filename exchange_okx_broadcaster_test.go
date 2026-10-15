package mytrade

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// fakeOkxOrdersServer answers login, subscribe and text ping like OKX's
// private endpoint, and lets a test push order messages on a chosen connection.
type fakeOkxOrdersServer struct {
	srv *httptest.Server

	mu    sync.Mutex
	conns []*websocket.Conn

	logins atomic.Int32
	subs   atomic.Int32
}

func newFakeOkxOrdersServer(t *testing.T) *fakeOkxOrdersServer {
	t.Helper()
	f := &fakeOkxOrdersServer{}
	up := websocket.Upgrader{}
	f.srv = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := up.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		f.mu.Lock()
		f.conns = append(f.conns, c)
		f.mu.Unlock()
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				return
			}
			text := string(msg)
			switch {
			case text == "ping":
				_ = c.WriteMessage(websocket.TextMessage, []byte("pong"))
			case strings.Contains(text, `"op":"login"`):
				f.logins.Add(1)
				_ = c.WriteMessage(websocket.TextMessage, []byte(`{"event":"login","code":"0","msg":"","connId":"fake"}`))
			case strings.Contains(text, `"op":"subscribe"`):
				f.subs.Add(1)
				_ = c.WriteMessage(websocket.TextMessage,
					[]byte(`{"event":"subscribe","arg":{"channel":"orders","instType":"SWAP"},"connId":"fake"}`))
			}
		}
	}))
	t.Cleanup(f.srv.Close)

	// The client always dials wss://ws.okx.com:8443, so redirect the dial.
	d := websocket.DefaultDialer
	oldDial, oldTLS, oldProxy := d.NetDialContext, d.TLSClientConfig, d.Proxy
	addr := f.srv.Listener.Addr().String()
	d.NetDialContext = func(ctx context.Context, network, _ string) (net.Conn, error) {
		return (&net.Dialer{}).DialContext(ctx, network, addr)
	}
	d.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	d.Proxy = nil
	t.Cleanup(func() { d.NetDialContext, d.TLSClientConfig, d.Proxy = oldDial, oldTLS, oldProxy })
	return f
}

func (f *fakeOkxOrdersServer) connCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.conns)
}

func (f *fakeOkxOrdersServer) push(t *testing.T, connIdx int, clOrdId string) {
	t.Helper()
	f.mu.Lock()
	var c *websocket.Conn
	if connIdx < len(f.conns) {
		c = f.conns[connIdx]
	}
	f.mu.Unlock()
	if c == nil {
		t.Fatalf("server conn %d not established", connIdx)
	}
	msg := `{"arg":{"channel":"orders","instType":"SWAP"},"data":[{"instType":"SWAP","instId":"BTC-USDT-SWAP","clOrdId":"` +
		clOrdId + `","ordId":"1","state":"live","accFillSz":"0"}]}`
	if err := c.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		t.Fatalf("push on conn %d: %v", connIdx, err)
	}
}

func waitUntil(t *testing.T, what string, d time.Duration, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(d)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("timeout waiting for %s", what)
}

func newTestOkxEngine() *OkxTradeEngine {
	return &OkxTradeEngine{apiKey: "k", secretKey: "s", passphrase: "p"}
}

func mustSubscriber(t *testing.T, o *OkxTradeEngine, clOrdId string) *okxOrderSubscriber {
	t.Helper()
	sub, err := o.newOrderSubscriber(&o.broadcasterSwap, clOrdId, "SWAP", "BTC-USDT-SWAP")
	if err != nil {
		t.Fatalf("newOrderSubscriber(%q): %v", clOrdId, err)
	}
	return sub
}

func expectOrder(t *testing.T, sub *okxOrderSubscriber, clOrdId string) {
	t.Helper()
	select {
	case o := <-sub.ch.ResultChan():
		if o.ClientOrderId != clOrdId {
			t.Fatalf("got clientOrderId %q, want %q", o.ClientOrderId, clOrdId)
		}
	case <-time.After(3 * time.Second):
		t.Fatalf("no push for %q", clOrdId)
	}
}

func expectNoOrder(t *testing.T, sub *okxOrderSubscriber) {
	t.Helper()
	select {
	case o := <-sub.ch.ResultChan():
		t.Fatalf("unexpected push: %+v", o)
	case <-time.After(200 * time.Millisecond):
	}
}

func TestOkxBroadcaster_RoutesByClientOrderId(t *testing.T) {
	f := newFakeOkxOrdersServer(t)
	o := newTestOkxEngine()
	waiter := mustSubscriber(t, o, "c1")
	other := mustSubscriber(t, o, "c2")
	longLived := mustSubscriber(t, o, "")

	f.push(t, 0, "c1")

	order, err := o.waitSubscribeReturn(waiter, 3*time.Second)
	if err != nil || order.ClientOrderId != "c1" {
		t.Fatalf("waiter got order=%+v err=%v", order, err)
	}
	expectOrder(t, longLived, "c1")
	expectNoOrder(t, other)
}

func TestOkxBroadcaster_ConcurrentCreationOpensOneConnection(t *testing.T) {
	f := newFakeOkxOrdersServer(t)
	o := newTestOkxEngine()

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := o.newOrderSubscriber(&o.broadcasterSwap, "", "SWAP", ""); err != nil {
				t.Errorf("newOrderSubscriber: %v", err)
			}
		}()
	}
	wg.Wait()

	if n := f.connCount(); n != 1 {
		t.Fatalf("opened %d connections, want 1", n)
	}
}

func TestOkxBroadcaster_TimeoutsTriggerReconnect(t *testing.T) {
	f := newFakeOkxOrdersServer(t)
	o := newTestOkxEngine()
	longLived := mustSubscriber(t, o, "")

	for i := 0; i < okxAckTimeoutThreshold; i++ {
		w := mustSubscriber(t, o, "no-push")
		if _, err := o.waitSubscribeReturn(w, 50*time.Millisecond); err == nil {
			t.Fatal("expected api msg timeout")
		}
	}

	waitUntil(t, "relogin and resubscribe after timeouts", 10*time.Second, func() bool {
		return f.logins.Load() >= 2 && f.subs.Load() >= 2
	})

	f.push(t, 1, "after-reconnect")
	expectOrder(t, longLived, "after-reconnect")
}

func TestOkxBroadcaster_SuccessResetsStrikes(t *testing.T) {
	f := newFakeOkxOrdersServer(t)
	o := newTestOkxEngine()
	mustSubscriber(t, o, "")
	b := o.broadcasterSwap

	for i := 0; i < okxAckTimeoutThreshold-1; i++ {
		w := mustSubscriber(t, o, "t")
		_, _ = o.waitSubscribeReturn(w, 30*time.Millisecond)
	}
	ok := mustSubscriber(t, o, "ok")
	f.push(t, 0, "ok")
	if _, err := o.waitSubscribeReturn(ok, 3*time.Second); err != nil {
		t.Fatalf("waitSubscribeReturn: %v", err)
	}
	if n := b.health.strikeCount(); n != 0 {
		t.Fatalf("strikes = %d after a successful push, want 0", n)
	}
}

func TestOkxBroadcaster_RebuildKeepsSubscribers(t *testing.T) {
	f := newFakeOkxOrdersServer(t)
	o := newTestOkxEngine()
	longLived := mustSubscriber(t, o, "")
	b := o.broadcasterSwap

	b.wsMu.RLock()
	oldWs := b.okxWsAccount
	b.wsMu.RUnlock()

	o.recoverOrderBroadcaster(b, okxRecoverRebuild, "test")

	b.wsMu.RLock()
	newWs := b.okxWsAccount
	b.wsMu.RUnlock()
	if newWs == oldWs {
		t.Fatal("ws client was not replaced")
	}
	if n := f.connCount(); n != 2 {
		t.Fatalf("server saw %d connections, want 2", n)
	}

	f.push(t, 1, "from-new-conn")
	expectOrder(t, longLived, "from-new-conn")
	select {
	case <-longLived.ch.CloseChan():
		t.Fatal("long-lived subscriber must not receive a close signal on rebuild")
	default:
	}
}

func TestOkxBroadcaster_SingleFlightRecovery(t *testing.T) {
	f := newFakeOkxOrdersServer(t)
	o := newTestOkxEngine()
	mustSubscriber(t, o, "")
	b := o.broadcasterSwap

	b.recovering.Store(true)
	o.recoverOrderBroadcaster(b, okxRecoverRebuild, "test")
	if n := f.connCount(); n != 1 {
		t.Fatalf("recovery ran while another was in flight: %d connections", n)
	}
}

func TestPushToOkxSubscriber_StalledConsumerDoesNotBlockForever(t *testing.T) {
	sub := &okxOrderSubscriber{
		clientOrderId: "",
		ch: &subscription[Order]{
			resultChan: make(chan Order, 1),
			errChan:    make(chan error, 1),
			closeChan:  make(chan struct{}, 1),
		},
	}
	pushToOkxSubscriber(sub, &Order{ClientOrderId: "fills-buffer"}, nil)

	done := make(chan struct{})
	go func() {
		pushToOkxSubscriber(sub, &Order{ClientOrderId: "dropped"}, nil)
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("push to a stalled subscriber blocked")
	}

	pushToOkxSubscriber(sub, nil, context.Canceled)
	pushToOkxSubscriber(sub, nil, context.Canceled) // error chan full: dropped, must not block
}
