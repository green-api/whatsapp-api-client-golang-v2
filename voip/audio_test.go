package voip

import (
	"context"
	"io"
	"sync/atomic"
	"testing"
	"time"
)

type oneFrame struct{ done bool }

func (t *oneFrame) Recv(context.Context) (AudioFrame, error) {
	if t.done {
		return AudioFrame{}, io.EOF
	}

	t.done = true

	return AudioFrame{Samples: make([]int16, 960), SampleRate: 48000, Channels: 1}, nil
}

func TestFrameSinkDeliversAndCloses(t *testing.T) {
	var frames atomic.Int32

	sink := NewFrameAudioSink(func(context.Context, AudioFrame) error { frames.Add(1); return nil })

	if e := sink.AttachAsync(context.Background(), &oneFrame{}); e != nil {
		t.Fatal(e)
	}

	waitFor(t, func() bool { return frames.Load() == 1 })

	if e := sink.CloseAsync(); e != nil {
		t.Fatal(e)
	}

	if e := sink.AttachAsync(context.Background(), &oneFrame{}); e == nil {
		t.Fatal("closed sink accepted track")
	}

	_ = sink.CloseAsync()
}

func TestFrameSinkWaitsForCallbackCancellation(t *testing.T) {
	entered := make(chan struct{})

	sink := NewFrameAudioSink(func(ctx context.Context, _ AudioFrame) error { close(entered); <-ctx.Done(); return ctx.Err() })
	_ = sink.AttachAsync(context.Background(), &oneFrame{})

	<-entered

	done := make(chan struct{})

	go func() { _ = sink.CloseAsync(); close(done) }()

	select {
	case <-done:

	case <-time.After(time.Second):
		t.Fatal("sink close hung")
	}
}

func TestCallbacksUnsubscribeAndAsyncReentry(t *testing.T) {
	h := newHarness(t)
	count := 0
	off := h.c.On("state", func(any) { count++ })
	h.s.message(`{"type":"state","state":{"state":"idle"}}`)
	h.c.Off(off)
	h.s.message(`{"type":"state","state":{"state":"idle"}}`)

	if count != 1 {
		t.Fatal(count)
	}

	done := make(chan struct{})
	h.c.OnAsync("state", func(any) { _ = h.c.StopAudioAsync(); close(done) })
	h.s.message(`{"type":"state","state":{"state":"idle"}}`)

	select {
	case <-done:

	case <-time.After(time.Second):
		t.Fatal("async callback deadlocked")
	}
}
