package voip

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/coder/websocket"
)

const InitialBackoff = 500 * time.Millisecond
const MaxBackoff = 10 * time.Second
const PermanentCloseMin = 4000
const PermanentCloseMax = 4999

type DisconnectDetail struct {
	Reason    string `json:"reason"`
	Code      int    `json:"code"`
	Permanent bool   `json:"permanent"`
}

type socketConn interface {
	Read(context.Context) (websocket.MessageType, []byte, error)
	Write(context.Context, websocket.MessageType, []byte) error
	Close(websocket.StatusCode, string) error
	Ping(context.Context) error
}

type ReconnectingSocket struct {
	url                      string
	urlError                 error
	mu                       sync.Mutex
	socket                   socketConn
	closed, refused, started bool
	ctx                      context.Context
	cancel                   context.CancelFunc
	done                     chan struct{}
	changed                  chan struct{}
	listeners                map[string][]func(any)
	factory                  func(context.Context, string) (socketConn, error)
	sleep                    func(context.Context, time.Duration) error
}

func NewReconnectingSocket(url string) *ReconnectingSocket {
	ctx, cancel := context.WithCancel(context.Background())

	return &ReconnectingSocket{url: url, ctx: ctx, cancel: cancel, done: make(chan struct{}), changed: make(chan struct{}), listeners: make(map[string][]func(any)), factory: func(ctx context.Context, url string) (socketConn, error) {
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)

		defer cancel()

		c, _, e := websocket.Dial(ctx, url, &websocket.DialOptions{CompressionMode: websocket.CompressionContextTakeover})

		if e == nil {
			c.SetReadLimit(1 << 20)
		}

		return c, e
	}, sleep: func(ctx context.Context, d time.Duration) error {
		t := time.NewTimer(d)

		defer t.Stop()

		select {
		case <-t.C:
			return nil

		case <-ctx.Done():
			return ctx.Err()
		}
	}}
}

func (s *ReconnectingSocket) On(name string, fn func(any)) {
	s.mu.Lock()

	defer s.mu.Unlock()

	s.listeners[name] = append(s.listeners[name], fn)
}

func (s *ReconnectingSocket) dispatch(name string, d any) {
	s.mu.Lock()
	listeners := append([]func(any){}, s.listeners[name]...)
	s.mu.Unlock()

	for _, fn := range listeners {
		func() {
			defer func() {
				if recover() != nil {
					slog.Error("callsRtc handler panicked", "event", name)
				}
			}()

			fn(d)
		}()
	}
}

func (s *ReconnectingSocket) Refused() bool { s.mu.Lock(); defer s.mu.Unlock(); return s.refused }
func (s *ReconnectingSocket) signal()       { close(s.changed); s.changed = make(chan struct{}) }
func (s *ReconnectingSocket) OpenAsync(ctx context.Context) error {
	s.mu.Lock()

	if s.urlError != nil {
		s.mu.Unlock()

		return s.urlError
	}

	if s.closed {
		s.mu.Unlock()

		return errors.New("ReconnectingSocket is closed")
	}

	if !s.started {
		s.started = true

		go s.connect()
	}

	s.mu.Unlock()

	for {
		s.mu.Lock()
		connected, closed := s.socket != nil, s.closed
		changed := s.changed
		s.mu.Unlock()

		if connected {
			return nil
		}

		if closed {
			return errors.New("ReconnectingSocket is closed")
		}

		select {
		case <-changed:

		case <-s.done:
			return errors.New("callsRtc connection was refused or closed")

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (s *ReconnectingSocket) Send(ctx context.Context, data any) error {
	s.mu.Lock()
	socket := s.socket
	s.mu.Unlock()

	if socket == nil {
		return errors.New("ReconnectingSocket is not connected")
	}

	b, e := json.Marshal(data)

	if e != nil {
		return e
	}

	return socket.Write(ctx, websocket.MessageText, b)
}

func (s *ReconnectingSocket) CloseAsync() error {
	s.mu.Lock()

	if s.closed {
		s.mu.Unlock()

		return nil
	}

	s.closed = true
	s.cancel()
	socket := s.socket
	s.socket = nil
	s.signal()
	started := s.started
	s.mu.Unlock()

	if socket != nil {
		_ = socket.Close(websocket.StatusNormalClosure, "")
	}

	if started {
		<-s.done
	}

	return nil
}

func (s *ReconnectingSocket) connect() {
	defer close(s.done)

	backoff := InitialBackoff

	for s.ctx.Err() == nil {
		socket, e := s.factory(s.ctx, s.url)

		if e == nil {
			s.mu.Lock()

			if s.closed {
				s.mu.Unlock()
				_ = socket.Close(websocket.StatusNormalClosure, "")

				return
			}

			s.socket = socket
			s.signal()
			s.mu.Unlock()
			backoff = InitialBackoff
			life, cancel := context.WithCancel(s.ctx)
			heartbeatDone := make(chan struct{})

			go func() {
				defer close(heartbeatDone)

				ticker := time.NewTicker(20 * time.Second)

				defer ticker.Stop()

				for {
					select {
					case <-life.Done():
						return

					case <-ticker.C:
						pingCtx, stop := context.WithTimeout(life, 20*time.Second)
						err := socket.Ping(pingCtx)
						stop()

						if err != nil {
							_ = socket.Close(websocket.StatusInternalError, "keepalive timeout")

							return
						}
					}
				}
			}()

			s.dispatch("connect", nil)

			for s.ctx.Err() == nil {
				_, b, err := socket.Read(s.ctx)

				if err != nil {
					e = err

					break
				}

				if json.Valid(b) {
					s.dispatch("message", json.RawMessage(append([]byte(nil), b...)))
				}
			}

			cancel()

			<-heartbeatDone

			s.mu.Lock()
			s.socket = nil
			s.signal()
			closed := s.closed
			s.mu.Unlock()
			_ = socket.Close(websocket.StatusNormalClosure, "")

			if closed {
				return
			}

			code := int(websocket.CloseStatus(e))

			if code < 0 {
				code = 1006
			}

			reason := "connection closed"

			var closeErr websocket.CloseError

			if errors.As(e, &closeErr) && closeErr.Reason != "" {
				reason = closeErr.Reason
			}

			permanent := code >= PermanentCloseMin && code <= PermanentCloseMax
			s.mu.Lock()
			s.refused = s.refused || permanent
			s.mu.Unlock()
			s.dispatch("disconnect", DisconnectDetail{reason, code, permanent})

			if permanent {
				return
			}
		} else if s.ctx.Err() == nil {
			s.dispatch("disconnect", DisconnectDetail{"connection failed", 1006, false})
		}

		if s.sleep(s.ctx, backoff) != nil {
			return
		}

		backoff *= 2

		if backoff > MaxBackoff {
			backoff = MaxBackoff
		}
	}
}
