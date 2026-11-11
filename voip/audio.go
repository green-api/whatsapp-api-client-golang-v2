package voip

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"sync"
)

// AudioFrame is interleaved signed 16-bit PCM. PTS is in samples per channel.
// The adapter supports the aiortc audio clock: 48 kHz, mono or stereo.
type AudioFrame struct {
	Samples              []int16
	SampleRate, Channels int
	PTS                  int64
}

type LocalAudioTrack interface {
	Recv(context.Context) (AudioFrame, error)
	Close() error
}

type RemoteAudioTrack interface {
	Recv(context.Context) (AudioFrame, error)
}

type CallAudio struct {
	LocalTrack    LocalAudioTrack
	OnRemoteTrack func(context.Context, RemoteAudioTrack) error
	Close         func() error
}

type AudioFactory func(context.Context) (*CallAudio, error)

// FrameAudioSink waits for each callback before receiving the next frame.
type FrameAudioSink struct {
	mu      sync.Mutex
	closed  bool
	cancel  []context.CancelFunc
	wg      sync.WaitGroup
	onFrame func(context.Context, AudioFrame) error
	onError func(error)
}

func NewFrameAudioSink(onFrame func(context.Context, AudioFrame) error) *FrameAudioSink {
	return &FrameAudioSink{onFrame: onFrame, onError: func(error) { slog.Error("Remote audio frame consumer failed") }}
}

func (s *FrameAudioSink) AttachAsync(ctx context.Context, track RemoteAudioTrack) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	if s.closed {
		return errors.New("Audio sink is closed")
	}

	if track == nil || s.onFrame == nil {
		return errors.New("audio track and callback are required")
	}

	ctx, cancel := context.WithCancel(ctx)
	s.cancel = append(s.cancel, cancel)
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()

		defer func() {
			if recover() != nil {
				slog.Error("Remote audio frame callback panicked")
			}
		}()

		defer cancel()

		for {
			f, e := track.Recv(ctx)

			if e == nil {
				e = s.onFrame(ctx, f)
			}

			if e != nil {
				if !errors.Is(e, io.EOF) && ctx.Err() == nil && s.onError != nil {
					s.onError(e)
				}

				return
			}
		}
	}()

	return nil
}

func (s *FrameAudioSink) CloseAsync() error {
	s.mu.Lock()
	s.closed = true

	for _, cancel := range s.cancel {
		cancel()
	}

	s.cancel = nil
	s.mu.Unlock()
	s.wg.Wait()

	return nil
}

// Codec instances are per direction and per media bridge; they must not be shared.
// opus supplies libopus without adding CGo to the core package.
type AudioCodec interface {
	Encode([]int16, int) ([]byte, error)
	Decode([]byte, int) ([]int16, error)
}

type CodecFactory func(channels int) (AudioCodec, error)

// attachmentGroup mirrors Python's tracked remote attachment tasks. No new task
// may enter once teardown starts; callbacks must honor their context.
type attachmentGroup struct {
	mu     sync.Mutex
	closed bool
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func newAttachmentGroup(parent context.Context) *attachmentGroup {
	ctx, cancel := context.WithCancel(parent)

	return &attachmentGroup{ctx: ctx, cancel: cancel}
}

func (g *attachmentGroup) start(fn func(context.Context)) {
	g.mu.Lock()

	defer g.mu.Unlock()

	if g.closed {
		return
	}

	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		defer func() {
			if recover() != nil {
				slog.Error("Remote audio attachment panicked")
			}
		}()

		fn(g.ctx)
	}()
}

func (g *attachmentGroup) close() {
	g.mu.Lock()
	g.closed = true
	g.cancel()
	g.mu.Unlock()
	g.wg.Wait()
}
