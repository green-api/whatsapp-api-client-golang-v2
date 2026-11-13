package voip

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/coder/websocket"
	"sync"
	"testing"
	"time"
)

type rawSocket struct {
	incoming chan any
	closed   chan struct{}
	once     sync.Once
}

func newRawSocket() *rawSocket {
	return &rawSocket{make(chan any, 10), make(chan struct{}), sync.Once{}}
}

func (s *rawSocket) Read(ctx context.Context) (websocket.MessageType, []byte, error) {
	select {
	case v := <-s.incoming:
		if e, ok := v.(error); ok {
			return 0, nil, e
		}

		return websocket.MessageText, []byte(v.(string)), nil

	case <-s.closed:
		return 0, nil, errors.New("closed")

	case <-ctx.Done():
		return 0, nil, ctx.Err()
	}
}

func (s *rawSocket) Write(context.Context, websocket.MessageType, []byte) error { return nil }
func (s *rawSocket) Close(websocket.StatusCode, string) error {
	s.once.Do(func() { close(s.closed) })

	return nil
}

func (s *rawSocket) Ping(context.Context) error { return nil }
func TestSocketReconnectPolicy(t *testing.T) {
	s := NewReconnectingSocket("ws://example.test")

	var mu sync.Mutex

	var sockets []*rawSocket

	var delays []time.Duration

	s.factory = func(context.Context, string) (socketConn, error) {
		raw := newRawSocket()
		mu.Lock()
		sockets = append(sockets, raw)
		mu.Unlock()

		return raw, nil
	}

	s.sleep = func(ctx context.Context, d time.Duration) error {
		mu.Lock()
		delays = append(delays, d)
		mu.Unlock()

		return ctx.Err()
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	if e := s.OpenAsync(ctx); e != nil {
		t.Fatal(e)
	}

	mu.Lock()
	first := sockets[0]
	mu.Unlock()

	first.incoming <- websocket.CloseError{Code: 1006}

	waitFor(t, func() bool { mu.Lock(); defer mu.Unlock(); return len(sockets) == 2 })
	mu.Lock()
	second := sockets[1]
	mu.Unlock()

	second.incoming <- websocket.CloseError{Code: 4001}

	waitFor(t, s.Refused)
	_ = s.CloseAsync()
	mu.Lock()

	defer mu.Unlock()

	if len(delays) != 1 || delays[0] != InitialBackoff {
		t.Fatal(delays)
	}
}

func TestBackoffCapAndOpenTimeoutKeepsRetrying(t *testing.T) {
	s := NewReconnectingSocket("ws://example.test")
	attempts := make(chan time.Duration, 20)
	s.factory = func(context.Context, string) (socketConn, error) { return nil, errors.New("failed") }

	s.sleep = func(ctx context.Context, d time.Duration) error {
		select {
		case attempts <- d:

		case <-ctx.Done():
			return ctx.Err()
		}

		select {
		case <-time.After(time.Millisecond):
			return nil

		case <-ctx.Done():
			return ctx.Err()
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)

	defer cancel()

	if s.OpenAsync(ctx) == nil {
		t.Fatal("expected timeout")
	}

	for _, want := range []time.Duration{500 * time.Millisecond, time.Second, 2 * time.Second, 4 * time.Second, 8 * time.Second, 10 * time.Second, 10 * time.Second} {
		select {
		case d := <-attempts:
			if d != want {
				t.Fatal(d, want)
			}

		case <-time.After(time.Second):
			t.Fatal("reconnect stopped")
		}
	}

	_ = s.CloseAsync()
}

func TestOrderedFrames(t *testing.T) {
	s := NewReconnectingSocket("ws://example.test")
	raw := newRawSocket()
	s.factory = func(context.Context, string) (socketConn, error) { return raw, nil }
	entered, release, second := make(chan struct{}), make(chan struct{}), make(chan struct{})

	s.On("message", func(d any) {
		if string(d.(json.RawMessage)) == `{"n":1}` {
			close(entered)

			<-release
		} else {
			close(second)
		}
	})

	if e := s.OpenAsync(context.Background()); e != nil {
		t.Fatal(e)
	}

	raw.incoming <- `{"n":1}`

	raw.incoming <- `{"n":2}`

	<-entered

	select {
	case <-second:
		t.Fatal("frames ran concurrently")

	case <-time.After(10 * time.Millisecond):
	}

	close(release)

	<-second

	_ = s.CloseAsync()
}
