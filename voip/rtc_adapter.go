package voip

import (
	"context"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/pion/interceptor"
	"github.com/pion/rtp/codecs"
	"github.com/pion/webrtc/v4"
	"github.com/pion/webrtc/v4/pkg/media"
	"github.com/pion/webrtc/v4/pkg/media/samplebuilder"
)

// RTCAdapter implements only the aiortc boundary used by rtc.py. Audio remains
// application-owned; PCM codecs are injected to keep this package CGo-free.
type RTCAdapter struct {
	peer         *webrtc.PeerConnection
	mu           sync.Mutex
	onTrack      func(RemoteAudioTrack)
	onError      func(error)
	closed       bool
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	codecFactory CodecFactory
	channels     int
	once         sync.Once
	closeErr     error
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

func newRTCAdapter(servers []IceServer) (*RTCAdapter, error) {
	var ice []webrtc.ICEServer

	for _, s := range servers {
		ice = append(ice, webrtc.ICEServer{URLs: s.URLs, Username: s.Username, Credential: s.Credential})
	}
	// aiortc offers Opus first, then G722/PCMU/PCMA. This PCM adapter negotiates
	// Opus only, the codec used by the GREEN-API audio bridge.

	engine := &webrtc.MediaEngine{}

	if e := engine.RegisterCodec(webrtc.RTPCodecParameters{RTPCodecCapability: webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeOpus, ClockRate: 48000, Channels: 2, SDPFmtpLine: "minptime=10;useinbandfec=1"}, PayloadType: 96}, webrtc.RTPCodecTypeAudio); e != nil {
		return nil, e
	}

	registry := &interceptor.Registry{}

	if e := webrtc.RegisterDefaultInterceptors(engine, registry); e != nil {
		return nil, e
	}

	api := webrtc.NewAPI(webrtc.WithMediaEngine(engine), webrtc.WithInterceptorRegistry(registry))
	peer, e := api.NewPeerConnection(webrtc.Configuration{ICEServers: ice})

	if e != nil {
		return nil, e
	}

	ctx, cancel := context.WithCancel(context.Background())
	a := &RTCAdapter{peer: peer, ctx: ctx, cancel: cancel}

	peer.OnTrack(func(t *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
		a.mu.Lock()

		if a.closed {
			a.mu.Unlock()

			return
		}

		a.wg.Add(1)
		factory, channels, fn := a.codecFactory, a.channels, a.onTrack

		defer a.wg.Done()

		a.mu.Unlock()

		if factory == nil || fn == nil {
			return
		}

		codec, e := factory(channels)

		if e != nil {
			a.reportError(e)

			return
		}

		remote := &decodedTrack{frames: make(chan AudioFrame, 32), done: make(chan struct{})}
		a.wg.Add(1)

		go func() { defer a.wg.Done(); remote.consume(a.ctx, t, codec, channels) }()

		fn(remote)
	})

	return a, nil
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

	codec, e := pcm.CodecFactory(pcm.Channels)

	if e != nil {
		return e
	}

	a.mu.Lock()
	a.codecFactory = pcm.CodecFactory
	a.channels = pcm.Channels
	a.mu.Unlock()
	track, e := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeOpus, ClockRate: 48000, Channels: 2}, "audio", "green-api")

	if e != nil {
		return e
	}

	sender, e := a.peer.AddTrack(track)

	if e != nil {
		return e
	}

	a.wg.Add(2)

	go func() {
		defer a.wg.Done()

		b := make([]byte, 1500)

		for {
			if _, _, e := sender.Read(b); e != nil {
				return
			}
		}
	}()

	go func() {
		defer a.wg.Done()

		ticker := time.NewTicker(20 * time.Millisecond)

		defer ticker.Stop()

		if waitForPeer(a.ctx, a.peer) != nil {
			return
		}

		var pending []int16

		for {
			frame, e := source.Recv(a.ctx)

			if e != nil {
				if a.ctx.Err() == nil && !errors.Is(e, io.EOF) {
					a.reportError(e)
				}

				return
			}

			if frame.SampleRate != 48000 || frame.Channels != pcm.Channels {
				a.reportError(errors.New("PCM frames must use 48000 Hz and the configured channel count"))

				return
			}

			pending = append(pending, frame.Samples...)
			size := 960 * pcm.Channels

			for len(pending) >= size {
				packet, e := codec.Encode(pending[:size], 960)

				if e != nil {
					a.reportError(e)

					return
				}

				pending = pending[size:]

				select {
				case <-a.ctx.Done():
					return

				case <-ticker.C:
				}

				if e = track.WriteSample(media.Sample{Data: packet, Duration: 20 * time.Millisecond}); e != nil {
					return
				}
			}
		}
	}()

	return nil
}

func waitForPeer(ctx context.Context, peer *webrtc.PeerConnection) error {
	ticker := time.NewTicker(5 * time.Millisecond)

	defer ticker.Stop()

	for {
		if peer.ConnectionState() == webrtc.PeerConnectionStateConnected {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-ticker.C:
		}
	}
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
