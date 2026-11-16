package voip

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/pion/interceptor"
	"github.com/pion/rtp/codecs"
	"github.com/pion/webrtc/v4"
	"github.com/pion/webrtc/v4/pkg/media"
	"github.com/pion/webrtc/v4/pkg/media/samplebuilder"
)

// RTCAdapter implements the media transport. Audio remains
// application-owned; PCM codecs are injected to keep this package CGo-free.
type RTCAdapter struct {
	peer          *webrtc.PeerConnection
	mu            sync.Mutex
	onTrack       func(RemoteAudioTrack)
	onError       func(error)
	closed        bool
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	codecFactory  CodecFactory
	channels      int
	once          sync.Once
	connectedOnce sync.Once
	closeErr      error
	connected     chan struct{}
}

// PCMTrack associates an application PCM source with its codec configuration.
type PCMTrack struct {
	LocalAudioTrack
	CodecFactory CodecFactory
	Channels     int
}

func NewPCMTrack(source LocalAudioTrack, factory CodecFactory, channels int) *PCMTrack {
	return &PCMTrack{source, factory, channels}
}

// NewRTCAdapter creates a Pion implementation of Bridge.
func NewRTCAdapter(servers []IceServer) (*RTCAdapter, error) {
	peer, err := newPeer(servers)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	a := &RTCAdapter{peer: peer, ctx: ctx, cancel: cancel, connected: make(chan struct{})}
	peer.OnConnectionStateChange(a.onConnectionState)
	peer.OnTrack(a.onRemoteTrack)

	return a, nil
}

func newPeer(servers []IceServer) (*webrtc.PeerConnection, error) {
	var ice []webrtc.ICEServer

	for _, s := range servers {
		ice = append(ice, webrtc.ICEServer{URLs: s.URLs, Username: s.Username, Credential: s.Credential})
	}
	// The GREEN-API audio bridge negotiates Opus.

	engine := &webrtc.MediaEngine{}

	if e := engine.RegisterCodec(webrtc.RTPCodecParameters{RTPCodecCapability: webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeOpus, ClockRate: 48000, Channels: 2, SDPFmtpLine: "minptime=10;useinbandfec=1"}, PayloadType: 96}, webrtc.RTPCodecTypeAudio); e != nil {
		return nil, e
	}

	registry := &interceptor.Registry{}

	if e := webrtc.RegisterDefaultInterceptors(engine, registry); e != nil {
		return nil, e
	}

	api := webrtc.NewAPI(webrtc.WithMediaEngine(engine), webrtc.WithInterceptorRegistry(registry))
	return api.NewPeerConnection(webrtc.Configuration{ICEServers: ice})
}

func (a *RTCAdapter) onConnectionState(state webrtc.PeerConnectionState) {
	switch state {
	case webrtc.PeerConnectionStateConnected:
		a.connectedOnce.Do(func() { close(a.connected) })
	case webrtc.PeerConnectionStateFailed:
		a.reportError(fmt.Errorf("peer connection failed (ICE state: %s)", a.peer.ICEConnectionState()))
		a.cancel()
	}
}

func (a *RTCAdapter) onRemoteTrack(track *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
	a.mu.Lock()

	if a.closed {
		a.mu.Unlock()
		return
	}

	a.wg.Add(1)
	factory, channels, fn := a.codecFactory, a.channels, a.onTrack
	a.mu.Unlock()

	defer a.wg.Done()

	if factory == nil || fn == nil {
		return
	}

	codec, err := factory(channels)

	if err != nil {
		a.reportError(err)
		return
	}

	remote := &decodedTrack{frames: make(chan AudioFrame, 32), done: make(chan struct{})}
	a.wg.Add(1)

	go func() { defer a.wg.Done(); remote.consume(a.ctx, track, codec, channels) }()

	fn(remote)
}

func (a *RTCAdapter) OnTrack(fn func(RemoteAudioTrack)) { a.mu.Lock(); a.onTrack = fn; a.mu.Unlock() }
func (a *RTCAdapter) AddTrack(source LocalAudioTrack) error {
	a.mu.Lock()

	if a.closed {
		a.mu.Unlock()
		return errors.New("peer closed")
	}

	a.wg.Add(1)
	a.mu.Unlock()

	defer a.wg.Done()
	pcm, ok := source.(*PCMTrack)

	if !ok || pcm.CodecFactory == nil || pcm.LocalAudioTrack == nil || (pcm.Channels != 1 && pcm.Channels != 2) {
		return errors.New("local track must be NewPCMTrack(source, codecFactory, 1 or 2)")
	}

	codec, err := pcm.CodecFactory(pcm.Channels)

	if err != nil {
		return err
	}

	a.mu.Lock()
	a.codecFactory = pcm.CodecFactory
	a.channels = pcm.Channels
	a.mu.Unlock()

	track, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeOpus, ClockRate: 48000, Channels: 2}, "audio", "green-api")

	if err != nil {
		return err
	}

	sender, err := a.peer.AddTrack(track)

	if err != nil {
		return err
	}

	a.wg.Add(2)

	go func() { defer a.wg.Done(); a.readSender(sender) }()
	go func() { defer a.wg.Done(); a.sendAudio(source, codec, track, pcm.Channels) }()

	return nil
}

func (a *RTCAdapter) readSender(sender *webrtc.RTPSender) {
	data := make([]byte, 1500)

	for {
		if _, _, err := sender.Read(data); err != nil {
			return
		}
	}
}

func (a *RTCAdapter) sendAudio(source LocalAudioTrack, codec AudioCodec, track *webrtc.TrackLocalStaticSample, channels int) {
	select {
	case <-a.connected:
	case <-a.ctx.Done():
		return
	}
	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()
	var pending []int16
	for {
		frame, err := source.Recv(a.ctx)

		if err != nil {
			if a.ctx.Err() == nil && !errors.Is(err, io.EOF) {
				a.reportError(err)
			}

			return
		}

		if frame.SampleRate != 48000 || frame.Channels != channels {
			a.reportError(errors.New("PCM frames must use 48000 Hz and the configured channel count"))
			return
		}

		pending = append(pending, frame.Samples...)
		pending, err = a.writeAudio(pending, codec, track, channels, ticker.C)

		if err != nil {
			if a.ctx.Err() == nil {
				a.reportError(err)
			}

			return
		}
	}
}

func (a *RTCAdapter) writeAudio(pending []int16, codec AudioCodec, track *webrtc.TrackLocalStaticSample, channels int, tick <-chan time.Time) ([]int16, error) {
	size := 960 * channels

	for len(pending) >= size {
		packet, err := codec.Encode(pending[:size], 960)

		if err != nil {
			return pending, err
		}

		pending = pending[size:]

		select {
		case <-a.ctx.Done():
			return pending, a.ctx.Err()
		case <-tick:
		}

		if err = track.WriteSample(media.Sample{Data: packet, Duration: 20 * time.Millisecond}); err != nil {
			return pending, err
		}
	}

	return pending, nil
}

func (a *RTCAdapter) CreateOffer() (SessionDescription, error)  { return a.peer.CreateOffer(nil) }
func (a *RTCAdapter) CreateAnswer() (SessionDescription, error) { return a.peer.CreateAnswer(nil) }
func (a *RTCAdapter) SetLocalDescription(ctx context.Context, d SessionDescription) error {
	done := webrtc.GatheringCompletePromise(a.peer)

	if e := a.peer.SetLocalDescription(d); e != nil {
		return e
	}

	select {
	case <-done:
		return nil

	case <-ctx.Done():
		return ctx.Err()

	case <-a.ctx.Done():
		return errors.New("peer closed")
	}
}

func (a *RTCAdapter) LocalDescription() *SessionDescription { return a.peer.LocalDescription() }
func (a *RTCAdapter) LocalCandidates() []IceCandidate {
	description := a.LocalDescription()
	if description == nil {
		return nil
	}
	return candidatesFromSDP(description.SDP)
}

func (a *RTCAdapter) SetRemoteDescription(d SessionDescription) error {
	return a.peer.SetRemoteDescription(d)
}

func (a *RTCAdapter) AddIceCandidate(c *IceCandidate) error {
	if e := validateCandidate(c); e != nil {
		return e
	}

	if c == nil {
		return a.peer.AddICECandidate(IceCandidate{})
	}

	return a.peer.AddICECandidate(*c)
}

func (a *RTCAdapter) GetStats() webrtc.StatsReport { return a.peer.GetStats() }
func (a *RTCAdapter) Close() error {
	a.once.Do(func() {
		a.mu.Lock()
		a.closed = true
		a.mu.Unlock()
		a.cancel()
		a.closeErr = a.peer.Close()
		a.wg.Wait()
	})

	return a.closeErr
}

type decodedTrack struct {
	frames chan AudioFrame
	done   chan struct{}
	mu     sync.Mutex
	err    error
}

func (t *decodedTrack) Recv(ctx context.Context) (AudioFrame, error) {
	select {
	case f, ok := <-t.frames:
		if ok {
			return f, nil
		}

		t.mu.Lock()
		e := t.err
		t.mu.Unlock()

		if e == nil {
			e = io.EOF
		}

		return AudioFrame{}, e

	case <-ctx.Done():
		return AudioFrame{}, ctx.Err()
	}
}

func (t *decodedTrack) consume(ctx context.Context, remote *webrtc.TrackRemote, codec AudioCodec, channels int) {
	defer close(t.frames)

	defer close(t.done)

	builder := samplebuilder.New(10, &codecs.OpusPacket{}, 48000)

	var pts int64

	for {
		packet, _, e := remote.ReadRTP()

		if e != nil {
			t.mu.Lock()
			t.err = e
			t.mu.Unlock()

			return
		}

		builder.Push(packet)

		for {
			sample := builder.Pop()

			if sample == nil {
				break
			}

			pcm, e := codec.Decode(sample.Data, 5760)

			if e != nil {
				t.mu.Lock()
				t.err = e
				t.mu.Unlock()

				return
			}

			frame := AudioFrame{Samples: pcm, SampleRate: 48000, Channels: channels, PTS: pts}
			pts += int64(len(pcm) / channels)

			select {
			case t.frames <- frame:

			case <-ctx.Done():
				return
			}
		}
	}
}

func (a *RTCAdapter) OnError(fn func(error)) { a.mu.Lock(); a.onError = fn; a.mu.Unlock() }
func (a *RTCAdapter) reportError(err error) {
	a.mu.Lock()
	fn, closed := a.onError, a.closed
	a.mu.Unlock()

	if fn != nil && !closed {
		go fn(err)
	}
}
