package voip

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/pion/webrtc/v4"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

type fakeSocket struct {
	mu        sync.Mutex
	listeners map[string]func(any)
	sent      []json.RawMessage
	fail      bool
}

func (s *fakeSocket) On(n string, f func(any))        { s.listeners[n] = f }
func (s *fakeSocket) OpenAsync(context.Context) error { return nil }
func (s *fakeSocket) CloseAsync() error               { return nil }
func (s *fakeSocket) Send(_ context.Context, v any) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	if s.fail {
		return errors.New("send failed")
	}

	b, _ := json.Marshal(v)
	s.sent = append(s.sent, b)

	return nil
}

func (s *fakeSocket) count() int       { s.mu.Lock(); defer s.mu.Unlock(); return len(s.sent) }
func (s *fakeSocket) message(v string) { s.listeners["message"](json.RawMessage(v)) }

type fakeTrack struct{}

func (fakeTrack) Recv(ctx context.Context) (AudioFrame, error) {
	<-ctx.Done()

	return AudioFrame{}, ctx.Err()
}

func (fakeTrack) Close() error { return nil }

type fakeBridge struct {
	mu           sync.Mutex
	calls        []string
	closed       bool
	onTrack      func(RemoteAudioTrack)
	badCandidate bool
	onError      func(error)
}

func (b *fakeBridge) record(s string)                                               { b.mu.Lock(); b.calls = append(b.calls, s); b.mu.Unlock() }
func (b *fakeBridge) AddTrack(LocalAudioTrack) error                                { return nil }
func (b *fakeBridge) OnTrack(f func(RemoteAudioTrack))                              { b.onTrack = f }
func (b *fakeBridge) CreateOffer() (SessionDescription, error)                      { return SessionDescription{}, nil }
func (b *fakeBridge) SetLocalDescription(context.Context, SessionDescription) error { return nil }
func (b *fakeBridge) LocalDescription() *SessionDescription                         { return nil }
func (b *fakeBridge) LocalCandidates() []IceCandidate                               { return nil }

func (b *fakeBridge) SetRemoteDescription(SessionDescription) error { b.record("answer"); return nil }
func (b *fakeBridge) AddIceCandidate(*IceCandidate) error {
	b.record("candidate")

	if b.badCandidate {
		return errors.New("bad candidate")
	}

	return nil
}

func (b *fakeBridge) Close() error           { b.mu.Lock(); b.closed = true; b.mu.Unlock(); return nil }
func (b *fakeBridge) OnError(fn func(error)) { b.mu.Lock(); b.onError = fn; b.mu.Unlock() }
func (b *fakeBridge) fail(err error) {
	b.mu.Lock()
	fn := b.onError
	b.mu.Unlock()
	if fn != nil {
		fn(err)
	}
}

type harness struct {
	c                       *CallsConnection
	s                       *fakeSocket
	mu                      sync.Mutex
	bridges                 []*fakeBridge
	audioCount, closedAudio int
}

func newHarness(t *testing.T) *harness {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, "[]") }))
	t.Cleanup(server.Close)
	h := &harness{s: &fakeSocket{listeners: map[string]func(any){}}}
	v := New(Config{Host: server.URL, IDInstance: "123", APITokenInstance: "secret"})

	h.c = newCallsConnection(v, func(context.Context) (*CallAudio, error) {
		h.mu.Lock()

		h.audioCount++

		h.mu.Unlock()

		return &CallAudio{LocalTrack: fakeTrack{}, OnRemoteTrack: func(context.Context, RemoteAudioTrack) error { return nil }, Close: func() error { h.mu.Lock(); h.closedAudio++; h.mu.Unlock(); return nil }}, nil
	}, h.s, func([]IceServer) (Bridge, error) {
		b := &fakeBridge{}
		h.mu.Lock()
		h.bridges = append(h.bridges, b)
		h.mu.Unlock()

		return b, nil
	})

	t.Cleanup(func() { _ = h.c.CloseAsync() })

	return h
}

func waitFor(t *testing.T, f func() bool) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)

	for !f() {
		if time.Now().After(deadline) {
			t.Fatal("condition timed out")
		}

		time.Sleep(time.Millisecond)
	}
}

func (h *harness) start(t *testing.T) chan error {
	t.Helper()
	before := h.s.count()
	done := make(chan error, 1)

	go func() { done <- h.c.StartAudioAsync(context.Background()) }()

	waitFor(t, func() bool { return h.s.count() > before })

	return done
}

func answer(h *harness) { h.s.message(`{"type":"answer","answer":{"type":"answer","sdp":"v=0"}}`) }
func result(t *testing.T, ch <-chan error, wantError bool) {
	t.Helper()

	select {
	case e := <-ch:
		if (e != nil) != wantError {
			t.Fatalf("error = %v", e)
		}

	case <-time.After(5 * time.Second):
		t.Fatal("operation hung")
	}
}

func TestRESTContract(t *testing.T) {
	var paths []string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.Method+" "+r.URL.Path)

		switch {
		case strings.Contains(r.URL.Path, "callsDial"):
			var p map[string]string

			if e := json.NewDecoder(r.Body).Decode(&p); e != nil {
				t.Error(e)
			}

			if p["chatId"] != "123@c.us" {
				t.Error(p)
			}

		case strings.Contains(r.URL.Path, "callsState"):
			io.WriteString(w, `{"state":"inc-call","info":{"id":"a","wid":"b","name":""}}`)

			return

		case strings.Contains(r.URL.Path, "callsGetIceServers"):
			io.WriteString(w, `[{"urls":"stun:example.org"}]`)

			return

		default:
			b, _ := io.ReadAll(r.Body)

			if len(b) != 0 || r.Header.Get("Content-Type") != "" {
				t.Error("unexpected POST body or content type")
			}
		}

		w.WriteHeader(204)
	}))

	defer server.Close()

	v := New(Config{Host: server.URL + "/prefix/", IDInstance: "1", APITokenInstance: "t"})
	ctx := context.Background()

	for _, e := range []error{v.DialAsync(ctx, "123"), v.AcceptAsync(ctx), v.RejectAsync(ctx), v.HangUpAsync(ctx)} {
		if e != nil {
			t.Fatal(e)
		}
	}

	s, e := v.GetStateAsync(ctx)

	if e != nil || s.Info.Name == nil || *s.Info.Name != "" {
		t.Fatal(s, e)
	}

	ice, e := v.GetIceServersAsync(ctx)

	if e != nil || len(ice) != 1 || ice[0].URLs[0] != "stun:example.org" {
		t.Fatal(ice, e)
	}

	if len(paths) != 6 {
		t.Fatal(paths)
	}

	u, e := v.WebsocketUrl()

	if e != nil || !strings.Contains(u, "/prefix/waInstance1/callsRtc/t") {
		t.Fatal(u, e)
	}
}

func TestRESTErrorsAndNoRetry(t *testing.T) {
	for _, code := range []int{200, 204, 302, 429, 500} {
		t.Run(http.StatusText(code), func(t *testing.T) {
			count := 0

			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				count++

				w.Header().Set("Location", "/redirect")
				w.WriteHeader(code)

				if code == 200 {
					io.WriteString(w, "invalid json")
				}
			}))

			defer s.Close()

			v := New(Config{Host: s.URL})
			e := v.AcceptAsync(context.Background())

			if (e == nil) != (code == 204) {
				t.Fatal(code, e)
			}

			if count != 1 {
				t.Fatal("retried", count)
			}
		})
	}
}

func TestEarlyCandidateAndNoAnswerTimeout(t *testing.T) {
	h := newHarness(t)
	done := h.start(t)
	h.s.message(`{"type":"ice-candidate","candidate":{"candidate":"candidate:test","sdpMid":"0"}}`)

	select {
	case <-done:
		t.Fatal("completed before answer")

	case <-time.After(25 * time.Millisecond):
	}

	answer(h)
	result(t, done, false)
	h.mu.Lock()
	b := h.bridges[0]
	h.mu.Unlock()
	b.mu.Lock()

	defer b.mu.Unlock()

	if !reflect.DeepEqual(b.calls, []string{"answer", "candidate"}) {
		t.Fatal(b.calls)
	}
}

func TestStopAndCloseDoNotHangUp(t *testing.T) {
	h := newHarness(t)
	h.s.message(`{"type":"state","state":{"state":"on-call"}}`)
	done := h.start(t)
	answer(h)
	result(t, done, false)

	if e := h.c.StopAudioAsync(); e != nil {
		t.Fatal(e)
	}

	if h.c.State().State != "on-call" || h.c.HasAudioBridge() {
		t.Fatal("stop changed call state")
	}

	h.s.mu.Lock()
	last := string(h.s.sent[len(h.s.sent)-1])
	h.s.mu.Unlock()

	if last != `{"type":"stop"}` {
		t.Fatal(last)
	}

	if e := h.c.CloseAsync(); e != nil {
		t.Fatal(e)
	}
}

func TestReconnectFreshAudioAndOldTrack(t *testing.T) {
	h := newHarness(t)
	h.s.message(`{"type":"state","state":{"state":"on-call"}}`)
	done := h.start(t)
	answer(h)
	result(t, done, false)
	h.mu.Lock()
	old := h.bridges[0]
	h.mu.Unlock()
	h.s.listeners["disconnect"](DisconnectDetail{Code: 1006})
	h.s.message(`{"type":"state","state":{"state":"on-call"}}`)
	waitFor(t, func() bool { return h.s.count() == 2 })
	old.onTrack(nil)
	answer(h)
	h.mu.Lock()

	defer h.mu.Unlock()

	if h.audioCount != 2 || h.closedAudio != 1 {
		t.Fatal(h.audioCount, h.closedAudio)
	}
}

func TestRepeatedDisconnectKeepsPendingResume(t *testing.T) {
	h := newHarness(t)
	h.s.message(`{"type":"state","state":{"state":"on-call"}}`)
	first := h.start(t)
	answer(h)
	result(t, first, false)
	ended := 0
	errors := make(chan ErrorDetail, 2)
	h.c.On("end_call", func(any) { ended++ })
	h.c.On("error", func(value any) { errors <- value.(ErrorDetail) })
	h.s.listeners["disconnect"](DisconnectDetail{Code: 1006})
	h.s.message(`{"type":"state","state":{"state":"on-call"}}`)
	waitFor(t, func() bool { return h.s.count() == 2 })
	h.s.listeners["disconnect"](DisconnectDetail{Code: 1006})
	h.c.op.Lock()
	pending := h.c.resumePending
	h.c.op.Unlock()
	if ended != 0 || h.c.State() == nil || !pending {
		t.Fatal("second disconnect ended a resumable call")
	}
	h.s.message(`{"type":"state","state":{"state":"on-call"}}`)
	waitFor(t, func() bool { return h.s.count() == 3 })
	answer(h)
	waitFor(t, func() bool {
		h.c.op.Lock()
		defer h.c.op.Unlock()
		return h.c.HasAudioBridge() && !h.c.resumePending
	})
	h.mu.Lock()
	count := h.audioCount
	h.mu.Unlock()
	if ended != 0 || count != 3 {
		t.Fatal("call did not resume after repeated disconnect")
	}
	select {
	case event := <-errors:
		t.Fatalf("transient disconnect became an error: %+v", event)
	case <-time.After(20 * time.Millisecond):
	}
}

func TestDuplicateAnswerKeepsCurrentBridge(t *testing.T) {
	h := newHarness(t)
	done := h.start(t)
	answer(h)
	result(t, done, false)
	answer(h)
	h.mu.Lock()
	peer := h.bridges[0]
	h.mu.Unlock()
	peer.mu.Lock()
	closed, calls := peer.closed, append([]string(nil), peer.calls...)
	peer.mu.Unlock()
	if closed || !h.c.HasAudioBridge() || !reflect.DeepEqual(calls, []string{"answer"}) {
		t.Fatalf("duplicate answer changed bridge: closed=%v calls=%v", closed, calls)
	}
}

func TestPeerFailureClosesOnlyItsOwnBridge(t *testing.T) {
	h := newHarness(t)
	done := h.start(t)
	answer(h)
	result(t, done, false)
	h.mu.Lock()
	old := h.bridges[0]
	h.mu.Unlock()
	events := make(chan ErrorDetail, 2)
	h.c.On("error", func(value any) { events <- value.(ErrorDetail) })
	old.fail(errors.New("ICE failed"))
	waitFor(t, func() bool { return !h.c.HasAudioBridge() })
	if (<-events).Message != "ICE failed" {
		t.Fatal("peer failure was not reported")
	}
	done = h.start(t)
	answer(h)
	result(t, done, false)
	old.fail(errors.New("late failure"))
	if !h.c.HasAudioBridge() {
		t.Fatal("old peer error closed the new bridge")
	}
}

func TestRTCFailureReportsErrorAndStopsMedia(t *testing.T) {
	adapter, err := NewRTCAdapter(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer adapter.Close()
	reported := make(chan error, 1)
	adapter.OnError(func(err error) { reported <- err })
	adapter.onConnectionState(webrtc.PeerConnectionStateFailed)
	select {
	case err := <-reported:
		if !strings.Contains(err.Error(), "peer connection failed") {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("RTC failure was not reported")
	}
	select {
	case <-adapter.ctx.Done():
	default:
		t.Fatal("RTC media did not stop")
	}
}

func TestRequestErrorKeepsOriginalURLCause(t *testing.T) {
	token := "private-token"
	original := &url.Error{Op: "Get", URL: "https://host/" + token, Err: errors.New("dial refused")}
	client := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) { return nil, original })}
	v := New(Config{Host: "https://host", IDInstance: "1", APITokenInstance: token, HTTPClient: client})
	_, err := v.GetStateAsync(context.Background())
	var request *RequestError
	var gotURL *url.Error
	if !errors.As(err, &request) || !errors.As(err, &gotURL) || !errors.Is(err, original) {
		t.Fatalf("original error lost: %v", err)
	}
	if strings.Contains(err.Error(), token) || request.Endpoint != "callsState" {
		t.Fatalf("unsafe or unhelpful error: %v", err)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestPermanentRefusalAndPendingDisconnect(t *testing.T) {
	for _, permanent := range []bool{false, true} {
		h := newHarness(t)
		h.s.message(`{"type":"state","state":{"state":"on-call"}}`)
		done := h.start(t)
		h.s.listeners["disconnect"](DisconnectDetail{Code: 4001, Permanent: permanent})
		result(t, done, true)

		if permanent && h.c.State() != nil {
			t.Fatal("permanent refusal retained call")
		}

		_ = h.c.CloseAsync()
	}
}

func TestErrorBeforeAndAfterAnswer(t *testing.T) {
	h := newHarness(t)
	done := h.start(t)
	h.s.message(`{"type":"error","message":"refused"}`)
	result(t, done, true)

	if h.c.HasAudioBridge() {
		t.Fatal("bridge retained")
	}

	done = h.start(t)
	answer(h)
	result(t, done, false)
	h.s.message(`{"type":"error","message":"notice"}`)

	if !h.c.HasAudioBridge() {
		t.Fatal("post-answer error closed bridge")
	}
}

func TestBadCandidateClosesBridge(t *testing.T) {
	h := newHarness(t)
	done := h.start(t)
	answer(h)
	result(t, done, false)
	h.mu.Lock()
	h.bridges[0].badCandidate = true
	h.mu.Unlock()
	h.s.message(`{"type":"ice-candidate","candidate":{"candidate":"bad"}}`)

	if h.c.HasAudioBridge() {
		t.Fatal("bridge retained")
	}
}

func TestCloseWhileFactoryPending(t *testing.T) {
	h := newHarness(t)
	entered, release := make(chan struct{}), make(chan struct{})
	closed := make(chan struct{})

	h.c.audioFactory = func(context.Context) (*CallAudio, error) {
		close(entered)

		<-release

		return &CallAudio{LocalTrack: fakeTrack{}, OnRemoteTrack: func(context.Context, RemoteAudioTrack) error { return nil }, Close: func() error { close(closed); return nil }}, nil
	}

	done := make(chan error, 1)

	go func() { done <- h.c.StartAudioAsync(context.Background()) }()

	<-entered

	if e := h.c.CloseAsync(); e != nil {
		t.Fatal(e)
	}

	close(release)
	result(t, done, true)

	<-closed
}

func TestIdleEndsEvenIfStopFails(t *testing.T) {
	h := newHarness(t)
	h.s.message(`{"type":"state","state":{"state":"on-call"}}`)
	done := h.start(t)
	answer(h)
	result(t, done, false)
	ended := false
	h.c.On("end_call", func(d any) { ended = d.(EndCallDetail).Reason == "call-ended" })
	h.s.mu.Lock()
	h.s.fail = true
	h.s.mu.Unlock()
	h.s.message(`{"type":"state","state":{"state":"idle","reason":"hangup"}}`)

	if !ended || h.c.HasAudioBridge() {
		t.Fatal("end not reported")
	}
}

func TestCancelNegotiationAndDuplicateStart(t *testing.T) {
	h := newHarness(t)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)

	go func() { done <- h.c.StartAudioAsync(ctx) }()

	waitFor(t, func() bool { return h.s.count() == 1 })

	if e := h.c.StartAudioAsync(ctx); e == nil {
		t.Fatal("duplicate start accepted")
	}

	cancel()
	result(t, done, true)

	if h.c.HasAudioBridge() {
		t.Fatal("bridge retained")
	}
}

func TestCancelAfterAnswerKeepsBridgeAlive(t *testing.T) {
	h := newHarness(t)
	startCtx, cancel := context.WithCancel(context.Background())
	mediaContext := make(chan context.Context, 1)
	attachmentContext := make(chan context.Context, 1)

	h.c.audioFactory = func(ctx context.Context) (*CallAudio, error) {
		mediaContext <- ctx

		return &CallAudio{
			LocalTrack: fakeTrack{},
			OnRemoteTrack: func(ctx context.Context, _ RemoteAudioTrack) error {
				attachmentContext <- ctx
				<-ctx.Done()

				return ctx.Err()
			},
			Close: func() error { return nil },
		}, nil
	}

	done := make(chan error, 1)
	go func() { done <- h.c.StartAudioAsync(startCtx) }()

	waitFor(t, func() bool { return h.s.count() == 1 })
	answer(h)
	result(t, done, false)
	cancel()

	if ctx := <-mediaContext; ctx.Err() != nil {
		t.Fatal("audio factory context ended after successful answer", ctx.Err())
	}

	h.mu.Lock()
	bridge := h.bridges[0]
	h.mu.Unlock()
	bridge.onTrack(&oneFrame{})

	ctx := <-attachmentContext
	if ctx.Err() != nil || !h.c.HasAudioBridge() {
		t.Fatal("bridge stopped when start context was canceled")
	}

	if err := h.c.StopAudioAsync(); err != nil {
		t.Fatal(err)
	}

	if ctx.Err() == nil {
		t.Fatal("attachment context survived bridge stop")
	}
}

func TestAnswerAndCancellationRace(t *testing.T) {
	for i := 0; i < 12; i++ {
		h := newHarness(t)
		ctx, cancel := context.WithCancel(context.Background())
		done := make(chan error, 1)

		go func() { done <- h.c.StartAudioAsync(ctx) }()
		waitFor(t, func() bool { return h.s.count() == 1 })

		var actions sync.WaitGroup
		actions.Add(2)

		go func() { defer actions.Done(); answer(h) }()
		go func() { defer actions.Done(); cancel() }()

		actions.Wait()
		err := <-done

		if err == nil && !h.c.HasAudioBridge() {
			t.Fatal("successful start lost its bridge")
		}

		if err != nil && h.c.HasAudioBridge() {
			t.Fatal("canceled start retained its bridge", err)
		}

		_ = h.c.CloseAsync()
	}
}

func TestCloseCancelsPendingRemoteAttachment(t *testing.T) {
	h := newHarness(t)
	entered := make(chan struct{})
	left := make(chan struct{})
	closed := make(chan struct{})

	h.c.audioFactory = func(context.Context) (*CallAudio, error) {
		return &CallAudio{LocalTrack: fakeTrack{}, OnRemoteTrack: func(ctx context.Context, _ RemoteAudioTrack) error {
			close(entered)

			<-ctx.Done()

			close(left)

			return ctx.Err()
		}, Close: func() error {
			select {
			case <-left:

			default:
				t.Error("audio closed before attachment exited")
			}

			close(closed)

			return nil
		}}, nil
	}

	done := h.start(t)
	answer(h)
	result(t, done, false)
	h.mu.Lock()
	bridge := h.bridges[0]
	h.mu.Unlock()
	bridge.onTrack(&oneFrame{})

	<-entered

	if e := h.c.CloseAsync(); e != nil {
		t.Fatal(e)
	}

	<-closed
}

func TestStateAndIncomingEvents(t *testing.T) {
	h := newHarness(t)

	var events []string

	var incoming CallInfo

	var end EndCallDetail

	h.c.On("state", func(d any) { events = append(events, d.(CallState).State) })
	h.c.On("incoming_call", func(d any) { events = append(events, "incoming"); incoming = d.(CallInfo) })
	h.c.On("end_call", func(d any) { events = append(events, "end"); end = d.(EndCallDetail) })
	h.s.message(`{"type":"state","state":{"state":"inc-call","info":{"id":"a","wid":"100@lid","name":""}}}`)

	if incoming.Name == nil || *incoming.Name != "" || incoming.WID != "100@lid" {
		t.Fatal(incoming)
	}

	h.s.message(`{"type":"state","state":{"state":"inc-call","info":{"id":"a","wid":"100@lid"}}}`)
	h.s.message(`{"type":"state","state":{"state":"idle","reason":"remote-hangup"}}`)

	if !reflect.DeepEqual(events, []string{"inc-call", "incoming", "inc-call", "idle", "end"}) || end.Cause != "remote-hangup" {
		t.Fatal(events, end)
	}
}

func TestStartSendFailureDisposesAudio(t *testing.T) {
	h := newHarness(t)
	h.s.fail = true

	if e := h.c.StartAudioAsync(context.Background()); e == nil {
		t.Fatal("send failure accepted")
	}

	if h.c.HasAudioBridge() {
		t.Fatal("bridge retained")
	}

	h.mu.Lock()

	defer h.mu.Unlock()

	if h.closedAudio != 1 {
		t.Fatal(h.closedAudio)
	}
}

func TestClosePendingNegotiationAndPostCloseOpen(t *testing.T) {
	h := newHarness(t)
	done := h.start(t)

	if e := h.c.CloseAsync(); e != nil {
		t.Fatal(e)
	}

	result(t, done, true)

	if e := h.c.OpenAsync(context.Background()); e == nil {
		t.Fatal("closed connection reopened")
	}

	if e := h.c.CloseAsync(); e != nil {
		t.Fatal(e)
	}
}
