package voip

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
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

// HandshakeRefusedError keeps the HTTP response and the WebSocket dial error.
// Its message omits the request URL because that URL contains the API token.
type HandshakeRefusedError struct {
	Response *http.Response
	Err      error
}

func (e *HandshakeRefusedError) Error() string {
	return fmt.Sprintf("callsRtc handshake refused: HTTP %d", e.StatusCode())
}

func (e *HandshakeRefusedError) Unwrap() error { return e.Err }
func (e *HandshakeRefusedError) StatusCode() int {
	if e.Response == nil {
		return 0
	}
	return e.Response.StatusCode
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
	stopError                error
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

		c, response, e := websocket.Dial(ctx, url, &websocket.DialOptions{CompressionMode: websocket.CompressionContextTakeover})

		if e != nil && response != nil && (response.StatusCode == http.StatusUnauthorized || response.StatusCode == http.StatusForbidden) {
			return nil, &HandshakeRefusedError{Response: response, Err: e}
		}

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
				if value := recover(); value != nil {
					slog.Error("callsRtc handler panicked", "event", name, "panic_type", fmt.Sprintf("%T", value), "stack", string(debug.Stack()))
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
	return s.waitForConnection(ctx)
}

func (s *ReconnectingSocket) waitForConnection(ctx context.Context) error {
	for {
		s.mu.Lock()
		connected, closed, stopError := s.socket != nil, s.closed, s.stopError
		changed := s.changed
		s.mu.Unlock()

		if connected {
			return nil
		}

		if closed {
			return errors.New("ReconnectingSocket is closed")
		}

		if stopError != nil {
			return stopError
		}

		select {
		case <-changed:

		case <-s.done:
			return s.connectionStopped()

		case <-ctx.Done():
			s.mu.Lock()
			stopError := s.stopError
			s.mu.Unlock()

			if stopError != nil {
				return stopError
			}

			return ctx.Err()
		}
	}
}

func (s *ReconnectingSocket) connectionStopped() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.stopError != nil {
		return s.stopError
	}

	return errors.New("callsRtc connection was refused or closed")
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
		socket, err := s.factory(s.ctx, s.url)
		var refusal *HandshakeRefusedError

		if errors.As(err, &refusal) {
			s.reportHandshakeRefusal(refusal, err)
			return
		}

		if err == nil {
			if !s.openSocket(socket) {
				return
			}

			backoff = InitialBackoff
			err = s.readSocket(socket)

			if s.closeSocket(socket) || s.reportSocketClose(err) {
				return
			}
		} else if s.ctx.Err() == nil {
			s.dispatch("disconnect", DisconnectDetail{"connection failed", 1006, false})
		}

		if s.sleep(s.ctx, backoff) != nil {
			return
		}

		backoff = min(backoff*2, MaxBackoff)
	}
}

func (s *ReconnectingSocket) reportHandshakeRefusal(refusal *HandshakeRefusedError, err error) {
	s.mu.Lock()
	s.refused = true
	s.stopError = err
	s.signal()
	s.mu.Unlock()
	s.dispatch("disconnect", DisconnectDetail{Reason: refusal.Error(), Code: refusal.StatusCode(), Permanent: true})
}

func (s *ReconnectingSocket) openSocket(socket socketConn) bool {
	s.mu.Lock()

	if s.closed {
		s.mu.Unlock()
		_ = socket.Close(websocket.StatusNormalClosure, "")
		return false
	}

	s.socket = socket
	s.signal()
	s.mu.Unlock()
	s.dispatch("connect", nil)

	return true
}

func (s *ReconnectingSocket) readSocket(socket socketConn) error {
	life, cancel := context.WithCancel(s.ctx)
	heartbeatDone := make(chan struct{})

	go func() { defer close(heartbeatDone); s.heartbeat(life, socket) }()
	defer func() { cancel(); <-heartbeatDone }()

	for s.ctx.Err() == nil {
		_, data, err := socket.Read(s.ctx)

		if err != nil {
			return err
		}

		if json.Valid(data) {
			s.dispatch("message", json.RawMessage(append([]byte(nil), data...)))
		}
	}

	return s.ctx.Err()
}

func (s *ReconnectingSocket) heartbeat(ctx context.Context, socket socketConn) {
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			pingCtx, stop := context.WithTimeout(ctx, 20*time.Second)
			err := socket.Ping(pingCtx)
			stop()

			if err != nil {
				_ = socket.Close(websocket.StatusInternalError, "keepalive timeout")
				return
			}
		}
	}
}

func (s *ReconnectingSocket) closeSocket(socket socketConn) bool {
	s.mu.Lock()
	s.socket = nil
	s.signal()
	closed := s.closed
	s.mu.Unlock()
	_ = socket.Close(websocket.StatusNormalClosure, "")

	return closed
}

func (s *ReconnectingSocket) reportSocketClose(err error) bool {
	code := int(websocket.CloseStatus(err))

	if code < 0 {
		code = 1006
	}

	reason := "connection closed"
	var closeErr websocket.CloseError

	if errors.As(err, &closeErr) && closeErr.Reason != "" {
		reason = closeErr.Reason
	}

	permanent := code >= PermanentCloseMin && code <= PermanentCloseMax
	s.mu.Lock()
	s.refused = s.refused || permanent
	s.mu.Unlock()
	s.dispatch("disconnect", DisconnectDetail{reason, code, permanent})

	return permanent
}
