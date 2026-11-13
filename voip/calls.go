// Package voip ports the GREEN-API Python tools/voip protocol and lifecycle.
package voip

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	StateIdle         = "idle"
	StateIncomingCall = "inc-call"
	StateOutgoingCall = "out-call"
	StateOnCall       = "on-call"
)

type CallInfo struct {
	ID   string  `json:"id"`
	WID  string  `json:"wid"`
	Name *string `json:"name"`
}

type CallState struct {
	State  string    `json:"state"`
	Info   *CallInfo `json:"info,omitempty"`
	Reason string    `json:"reason,omitempty"`
}

func isCall(s string) bool { return s == "inc-call" || s == "out-call" || s == "on-call" }

type Config struct {
	Host, IDInstance, APITokenInstance string
	HTTPClient                         *http.Client
	Logger                             *slog.Logger
}

type Voip struct {
	config Config
	client *http.Client
}

// readDeadlineConn implements aiohttp's idle socket-read timeout, not a total deadline.
type readDeadlineConn struct{ net.Conn }

func (c readDeadlineConn) Read(p []byte) (int, error) {
	_ = c.SetReadDeadline(time.Now().Add(5 * time.Second))

	return c.Conn.Read(p)
}

func New(config Config) *Voip {
	if config.Host == "" {
		config.Host = "https://api.green-api.com"
	}

	if config.Logger == nil {
		config.Logger = slog.Default()
	}

	client := config.HTTPClient

	if client == nil {
		dialer := net.Dialer{Timeout: 5 * time.Second}

		client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyFromEnvironment, DisableKeepAlives: true, TLSHandshakeTimeout: 5 * time.Second, ResponseHeaderTimeout: 5 * time.Second, DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			c, e := dialer.DialContext(ctx, network, addr)

			if e != nil {
				return nil, e
			}

			return readDeadlineConn{c}, nil
		}}}
	}

	copyClient := *client
	copyClient.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }

	return &Voip{config: config, client: &copyClient}
}

func (v *Voip) endpoint(method string) string {
	return strings.TrimRight(v.config.Host, "/") + "/waInstance" + v.config.IDInstance + "/" + method + "/" + v.config.APITokenInstance
}

func (v *Voip) WebsocketUrl() (string, error) {
	u, e := url.Parse(v.config.Host)

	if e != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") {
		return "", errors.New("VoIP host must use http or https")
	}

	if u.Scheme == "https" {
		u.Scheme = "wss"
	} else {
		u.Scheme = "ws"
	}

	u.Path = strings.TrimRight(u.Path, "/") + "/waInstance" + v.config.IDInstance + "/callsRtc/" + v.config.APITokenInstance
	u.RawPath = ""
	u.RawQuery = ""
	u.Fragment = ""

	return u.String(), nil
}

func (v *Voip) request(ctx context.Context, method, endpoint string, payload any, out any) error {
	var body io.Reader

	if payload != nil {
		b, e := json.Marshal(payload)

		if e != nil {
			return e
		}

		body = bytes.NewReader(b)
	}

	req, e := http.NewRequestWithContext(ctx, method, v.endpoint(endpoint), body)

	if e != nil {
		return fmt.Errorf("%s: invalid request", endpoint)
	}

	req.Header.Set("User-Agent", "GREEN-API_SDK_GO/1.0")

	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, e := v.client.Do(req)

	if e != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		return fmt.Errorf("%s: connection failed", endpoint)
	}

	defer resp.Body.Close()

	b, e := io.ReadAll(resp.Body)

	if e != nil {
		return fmt.Errorf("%s: reading response failed", endpoint)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("%s failed: %d %s", endpoint, resp.StatusCode, b)
	}

	if resp.StatusCode == 204 || len(b) == 0 {
		return nil
	}

	if out == nil {
		var discard any

		out = &discard
	}

	return json.Unmarshal(b, out)
}

func (v *Voip) GetStateAsync(ctx context.Context) (CallState, error) {
	var s CallState

	e := v.request(ctx, "GET", "callsState", nil, &s)

	if e == nil && s.State == "" {
		e = errors.New("callsState response is missing state")
	}

	return s, e
}

func (v *Voip) GetIceServersAsync(ctx context.Context) ([]IceServer, error) {
	var s []IceServer

	e := v.request(ctx, "GET", "callsGetIceServers", nil, &s)

	return s, e
}

func (v *Voip) DialAsync(ctx context.Context, target string) error {
	if !strings.Contains(target, "@") {
		target += "@c.us"
	}

	return v.request(ctx, "POST", "callsDial", map[string]string{"chatId": target}, nil)
}

func (v *Voip) AcceptAsync(ctx context.Context) error {
	return v.request(ctx, "POST", "callsAccept", nil, nil)
}

func (v *Voip) RejectAsync(ctx context.Context) error {
	return v.request(ctx, "POST", "callsReject", nil, nil)
}

func (v *Voip) HangUpAsync(ctx context.Context) error {
	return v.request(ctx, "POST", "callsHangUp", nil, nil)
}

type Callback func(any)
type callbackEntry struct {
	id    uint64
	fn    Callback
	async bool
}

type pendingBridge struct {
	done              chan error
	bridgeGeneration  uint64
	startContext      context.Context
	mediaContext      context.Context
	stopStartupCancel func() bool
}

// CallsConnection owns one callsRtc socket and its current audio bridge.
// op serializes protocol mutations; slow startup work runs outside it, as Python awaits do.
type CallsConnection struct {
	attachments                                                  *attachmentGroup
	voip                                                         *Voip
	socket                                                       signalingSocket
	audioFactory                                                 AudioFactory
	bridgeFactory                                                func([]IceServer) (Bridge, error)
	op                                                           sync.Mutex
	mu                                                           sync.RWMutex
	state                                                        *CallState
	peerConnection                                               Bridge
	audio                                                        *CallAudio
	pendingBridge                                                *pendingBridge
	pendingRemoteCandidates                                      []*IceCandidate
	remoteDescriptionSet, resumePending, refusalReported, closed bool
	bridgeGeneration                                             uint64
	bridgeCancel                                                 context.CancelFunc
	callbacks                                                    map[string][]callbackEntry
	nextCallback                                                 uint64
}

type signalingSocket interface {
	OpenAsync(context.Context) error
	Send(context.Context, any) error
	CloseAsync() error
	On(string, func(any))
}

func (v *Voip) Connect(factory AudioFactory) *CallsConnection {
	u, e := v.WebsocketUrl()
	s := NewReconnectingSocket(u)
	s.urlError = e

	return newCallsConnection(v, factory, s, func(servers []IceServer) (Bridge, error) { return NewAiortcBridge(servers) })
}

func newCallsConnection(v *Voip, f AudioFactory, s signalingSocket, b func([]IceServer) (Bridge, error)) *CallsConnection {
	c := &CallsConnection{voip: v, audioFactory: f, socket: s, bridgeFactory: b, callbacks: make(map[string][]callbackEntry)}
	s.On("connect", func(d any) { c.emit("connect", d) })
	s.On("disconnect", func(d any) { c.onSocketDisconnect(d.(DisconnectDetail)) })

	s.On("message", func(d any) {
		if m, ok := d.(json.RawMessage); ok {
			c.onMessage(m)
		}
	})

	return c
}

// On matches a synchronous Python callback. It must not block or call lifecycle methods;
// use OnAsync for callbacks that start/stop audio or issue other blocking operations.
// The returned function unregisters this exact callback (Go functions aren't comparable).
func (c *CallsConnection) On(name string, fn Callback) func() { return c.addCallback(name, fn, false) }
func (c *CallsConnection) OnAsync(name string, fn Callback) func() {
	return c.addCallback(name, fn, true)
}

func (c *CallsConnection) Off(unsubscribe func()) {
	if unsubscribe != nil {
		unsubscribe()
	}
}

func (c *CallsConnection) addCallback(name string, fn Callback, async bool) func() {
	c.mu.Lock()

	c.nextCallback++

	id := c.nextCallback
	c.callbacks[name] = append(c.callbacks[name], callbackEntry{id, fn, async})
	c.mu.Unlock()

	return func() {
		c.mu.Lock()

		defer c.mu.Unlock()

		for i, e := range c.callbacks[name] {
			if e.id == id {
				c.callbacks[name] = append(c.callbacks[name][:i], c.callbacks[name][i+1:]...)

				break
			}
		}
	}
}

func (c *CallsConnection) emit(name string, d any) {
	c.mu.RLock()
	entries := append([]callbackEntry(nil), c.callbacks[name]...)
	c.mu.RUnlock()

	for _, e := range entries {
		invoke := func() {
			defer func() {
				if recover() != nil {
					c.voip.config.Logger.Error("VoIP callback panicked", "event", name)
				}
			}()

			e.fn(d)
		}

		if e.async {
			go invoke()
		} else {
			invoke()
		}
	}
}

func (c *CallsConnection) State() *CallState {
	c.mu.RLock()

	defer c.mu.RUnlock()

	if c.state == nil {
		return nil
	}

	s := *c.state

	if s.Info != nil {
		i := *s.Info

		if i.Name != nil {
			n := *i.Name
			i.Name = &n
		}

		s.Info = &i
	}

	return &s
}

func (c *CallsConnection) HasAudioBridge() bool {
	c.mu.RLock()

	defer c.mu.RUnlock()

	return c.peerConnection != nil
}

func (c *CallsConnection) OpenAsync(ctx context.Context) error {
	c.mu.RLock()
	closed := c.closed
	c.mu.RUnlock()

	if closed {
		return errors.New("CallsConnection is closed")
	}

	return c.socket.OpenAsync(ctx)
}

func (c *CallsConnection) StartAudioAsync(ctx context.Context) error {
	return c.startAudioAsync(ctx, nil)
}

func (c *CallsConnection) startAudioAsync(ctx context.Context, expected *uint64) error {
	c.op.Lock()
	c.mu.Lock()

	if c.closed || c.peerConnection != nil || c.pendingBridge != nil || (expected != nil && c.bridgeGeneration != *expected) {
		c.mu.Unlock()
		c.op.Unlock()

		return errors.New("CallsConnection closed or audio bridge already starting or active")
	}

	bctx, cancel := context.WithCancel(context.WithoutCancel(ctx))
	p := &pendingBridge{done: make(chan error, 1), bridgeGeneration: c.bridgeGeneration, startContext: ctx, mediaContext: bctx}
	p.stopStartupCancel = context.AfterFunc(ctx, cancel)
	c.pendingBridge = p
	c.bridgeCancel = cancel
	c.attachments = newAttachmentGroup(bctx)
	c.mu.Unlock()
	c.op.Unlock()
	err := c.startAudioInternal(ctx, p)

	if err != nil {
		c.cancelPendingStart(p, err)
	}

	return c.awaitPendingStart(ctx, p)
}

func (c *CallsConnection) cancelPendingStart(p *pendingBridge, err error) bool {
	c.op.Lock()

	defer c.op.Unlock()

	c.mu.RLock()
	current := c.pendingBridge == p
	c.mu.RUnlock()

	if current {
		c.rejectPending(err)
		_ = c.teardownBridge(false)
	}

	return current
}

func (c *CallsConnection) awaitPendingStart(ctx context.Context, p *pendingBridge) error {
	select {
	case result := <-p.done:
		return result

	case <-ctx.Done():
		if c.cancelPendingStart(p, ctx.Err()) {
			return ctx.Err()
		}

		return <-p.done
	}
}

func (c *CallsConnection) current(p *pendingBridge) bool {
	c.mu.RLock()

	defer c.mu.RUnlock()

	return c.bridgeGeneration == p.bridgeGeneration && !c.closed
}

func (c *CallsConnection) startAudioInternal(ctx context.Context, p *pendingBridge) error {
	audio, err := c.prepareAudio(p)

	if err != nil {
		return err
	}

	peer, err := c.preparePeer(ctx, p, audio)

	if err != nil {
		return err
	}

	return c.sendOffer(ctx, p, peer)
}

func (c *CallsConnection) prepareAudio(p *pendingBridge) (*CallAudio, error) {
	if c.audioFactory == nil {
		return nil, errors.New("audio_factory is required")
	}

	audio, err := c.audioFactory(p.mediaContext)

	if err != nil {
		return nil, err
	}

	if audio == nil || audio.LocalTrack == nil || audio.Close == nil || audio.OnRemoteTrack == nil {
		if audio != nil && audio.Close != nil {
			_ = audio.Close()
		}

		return nil, errors.New("audio_factory returned incomplete CallAudio")
	}

	c.op.Lock()

	if !c.current(p) {
		c.op.Unlock()
		_ = audio.Close()

		return nil, errors.New("audio bridge superseded")
	}

	c.mu.Lock()
	c.audio = audio
	c.mu.Unlock()
	c.op.Unlock()

	return audio, nil
}

func (c *CallsConnection) preparePeer(ctx context.Context, p *pendingBridge, audio *CallAudio) (Bridge, error) {
	servers, err := c.voip.GetIceServersAsync(ctx)

	if err != nil {
		return nil, err
	}

	if !c.current(p) {
		return nil, errors.New("audio bridge superseded")
	}

	peer, err := c.bridgeFactory(servers)

	if err != nil {
		return nil, err
	}

	c.op.Lock()

	if !c.current(p) {
		c.op.Unlock()
		_ = peer.Close()

		return nil, errors.New("audio bridge superseded")
	}

	c.mu.Lock()
	c.peerConnection = peer
	c.remoteDescriptionSet = false
	c.pendingRemoteCandidates = nil
	attachments := c.attachments
	c.mu.Unlock()

	c.bindPeer(p, audio, peer, attachments)
	err = peer.AddTrack(audio.LocalTrack)
	c.op.Unlock()

	return peer, err
}

func (c *CallsConnection) bindPeer(p *pendingBridge, audio *CallAudio, peer Bridge, attachments *attachmentGroup) {
	peer.OnTrack(func(track RemoteAudioTrack) {
		if !c.current(p) {
			return
		}

		attachments.start(func(attachCtx context.Context) {
			if !c.current(p) {
				return
			}

			c.emit("remote_track", track)

			if err := audio.OnRemoteTrack(attachCtx, track); err != nil && attachCtx.Err() == nil {
				c.emit("error", ErrorDetail{err.Error()})
			}
		})
	})

	if reporter, ok := peer.(interface{ OnError(func(error)) }); ok {
		reporter.OnError(func(err error) {
			if c.current(p) {
				c.emit("error", ErrorDetail{err.Error()})
			}
		})
	}
}

func (c *CallsConnection) sendOffer(ctx context.Context, p *pendingBridge, peer Bridge) error {
	offer, err := peer.CreateOffer()

	if err != nil {
		return err
	}

	if !c.current(p) {
		return errors.New("audio bridge superseded")
	}

	if err = peer.SetLocalDescription(ctx, offer); err != nil {
		return err
	}

	c.op.Lock()

	defer c.op.Unlock()

	if !c.current(p) {
		return errors.New("audio bridge superseded")
	}

	if local := peer.LocalDescription(); local != nil {
		offer = *local
	}

	if err = c.socket.Send(ctx, map[string]any{"type": "offer", "offer": offer}); err != nil {
		return err
	}

	for _, candidate := range peer.LocalCandidates() {
		if err = c.socket.Send(ctx, map[string]any{"type": "ice-candidate", "candidate": candidate}); err != nil {
			return err
		}
	}

	return nil
}

func (c *CallsConnection) rejectPending(e error) {
	c.mu.Lock()
	p := c.pendingBridge
	c.pendingBridge = nil
	c.mu.Unlock()

	if p != nil {
		p.stopStartupCancel()
		p.startContext = nil
		p.stopStartupCancel = nil
		p.done <- e
	}
}

func (c *CallsConnection) teardownBridge(notify bool) error {
	c.mu.Lock()

	c.bridgeGeneration++

	peer, a, cancel := c.peerConnection, c.audio, c.bridgeCancel
	attachments := c.attachments
	c.attachments = nil
	c.peerConnection = nil
	c.audio = nil
	c.bridgeCancel = nil
	c.pendingRemoteCandidates = nil
	c.remoteDescriptionSet = false
	c.mu.Unlock()

	if cancel != nil {
		cancel()
	}

	if attachments != nil {
		attachments.close()
	}

	var err error

	if notify && peer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		err = c.socket.Send(ctx, map[string]string{"type": "stop"})
		cancel()
	}

	if peer != nil {
		err = errors.Join(err, peer.Close())
	}

	if a != nil {
		err = errors.Join(err, a.Close())
	}

	return err
}

func (c *CallsConnection) StopAudioAsync() error {
	c.op.Lock()

	defer c.op.Unlock()

	c.resumePending = false
	c.rejectPending(errors.New("Audio bridge stopped"))

	return c.teardownBridge(true)
}

func (c *CallsConnection) CloseAsync() error {
	c.op.Lock()
	c.mu.Lock()

	if c.closed {
		c.mu.Unlock()
		c.op.Unlock()

		return nil
	}

	c.closed = true
	c.mu.Unlock()
	c.resumePending = false
	c.rejectPending(errors.New("CallsConnection closed"))
	e := c.teardownBridge(false)
	c.op.Unlock()

	return errors.Join(e, c.socket.CloseAsync())
}

type ErrorDetail struct {
	Message string `json:"message"`
}

type EndCallDetail struct {
	Reason string `json:"reason"`
	Cause  string `json:"cause,omitempty"`
}

func (c *CallsConnection) onMessage(raw json.RawMessage) {
	var m struct {
		Type      string             `json:"type"`
		State     CallState          `json:"state"`
		Answer    SessionDescription `json:"answer"`
		Candidate *IceCandidate      `json:"candidate"`
		Message   string             `json:"message"`
	}

	if json.Unmarshal(raw, &m) != nil {
		return
	}

	c.op.Lock()

	defer c.op.Unlock()

	c.mu.RLock()
	closed, peer, remote := c.closed, c.peerConnection, c.remoteDescriptionSet
	c.mu.RUnlock()

	if closed {
		return
	}

	switch m.Type {
	case "state":
		c.onState(m.State)

	case "answer":
		c.onAnswer(peer, m.Answer)

	case "ice-candidate":
		if remote && peer != nil {
			if e := peer.AddIceCandidate(m.Candidate); e != nil {
				c.rejectPending(e)
				_ = c.teardownBridge(false)
				c.emit("error", ErrorDetail{e.Error()})
			}
		} else if !remote {
			c.mu.Lock()
			c.pendingRemoteCandidates = append(c.pendingRemoteCandidates, m.Candidate)
			c.mu.Unlock()
		}

	case "error":
		c.mu.RLock()
		pending := c.pendingBridge != nil
		c.mu.RUnlock()

		if pending {
			c.rejectPending(errors.New(m.Message))
			_ = c.teardownBridge(false)
		} else {
			c.refusalReported = true
			c.emit("error", ErrorDetail{m.Message})
		}
	}
}

func (c *CallsConnection) onAnswer(peer Bridge, answer SessionDescription) {
	if peer == nil {
		return
	}

	err := peer.SetRemoteDescription(answer)

	if err == nil {
		c.mu.Lock()
		c.remoteDescriptionSet = true
		candidates := c.pendingRemoteCandidates
		c.pendingRemoteCandidates = nil
		c.mu.Unlock()

		for _, candidate := range candidates {
			if err = peer.AddIceCandidate(candidate); err != nil {
				break
			}
		}
	}

	if err == nil {
		c.mu.RLock()
		pending := c.pendingBridge
		c.mu.RUnlock()

		if pending != nil {
			err = pending.startContext.Err()

			if err == nil && !pending.stopStartupCancel() {
				err = context.Canceled
			}
		}
	}

	if err != nil {
		c.rejectPending(err)
		_ = c.teardownBridge(false)

		return
	}

	c.rejectPending(nil)
}

func (c *CallsConnection) onState(s CallState) {
	c.mu.Lock()
	previous := ""

	if c.state != nil {
		previous = c.state.State
	}

	c.state = &s
	c.mu.Unlock()
	c.emit("state", s)

	if s.State == "inc-call" && previous != "inc-call" && s.Info != nil {
		c.emit("incoming_call", *s.Info)

		return
	}

	if isCall(previous) && !isCall(s.State) {
		c.resumePending = false
		c.rejectPending(errors.New("Call ended during negotiation"))

		if e := c.teardownBridge(true); e != nil {
			c.voip.config.Logger.Warn("Could not send callsRtc stop")
		}

		c.emit("end_call", EndCallDetail{"call-ended", s.Reason})

		return
	}

	c.mu.RLock()
	canResume := c.peerConnection == nil && c.pendingBridge == nil
	c.mu.RUnlock()

	if c.resumePending && isCall(s.State) && canResume {
		c.resumePending = false
		generation := c.bridgeGeneration

		go func() {
			if e := c.startAudioAsync(context.Background(), &generation); e != nil {
				c.emit("error", ErrorDetail{e.Error()})
			}
		}()
	}
}

func (c *CallsConnection) onSocketDisconnect(d DisconnectDetail) {
	c.op.Lock()

	defer c.op.Unlock()

	c.mu.RLock()
	previous := ""

	if c.state != nil {
		previous = c.state.State
	}

	hasPeer := c.peerConnection != nil
	closed := c.closed
	c.mu.RUnlock()

	if closed {
		return
	}

	if isCall(previous) {
		c.resumePending = hasPeer
		c.rejectPending(errors.New("Socket disconnected during negotiation"))
		_ = c.teardownBridge(false)

		if d.Permanent {
			c.resumePending = false
		}

		if !c.resumePending {
			c.mu.Lock()
			c.state = nil
			c.mu.Unlock()
			c.emit("end_call", EndCallDetail{Reason: "connection-lost"})
		}
	} else {
		c.mu.Lock()
		c.state = nil
		c.mu.Unlock()
	}

	if d.Permanent && !c.refusalReported {
		c.emit("error", ErrorDetail{d.Reason})
	}

	c.refusalReported = false
	c.emit("disconnect", d)
}
