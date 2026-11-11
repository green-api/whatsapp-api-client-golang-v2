package voip_test

import (
	"context"
	"math"
	"sync/atomic"
	"testing"
	"time"

	"github.com/green-api/whatsapp-api-client-golang-v2/opus"
	"github.com/green-api/whatsapp-api-client-golang-v2/voip"
)

type tone struct{ pts int64 }

func (t *tone) Recv(ctx context.Context) (voip.AudioFrame, error) {
	if e := ctx.Err(); e != nil {
		return voip.AudioFrame{}, e
	}

	samples := make([]int16, 960)

	for i := range samples {
		samples[i] = int16(10000 * math.Sin(2*math.Pi*440*float64(t.pts+int64(i))/48000))
	}

	f := voip.AudioFrame{Samples: samples, SampleRate: 48000, Channels: 1, PTS: t.pts}
	t.pts += 960

	return f, nil
}

func (*tone) Close() error { return nil }
func requireCodec(t *testing.T) {
	t.Helper()

	if _, e := opus.New(1); e != nil {
		t.Skip(e)
	}
}

func TestTwoRealPeersExchangeDecodedAudio(t *testing.T) {
	requireCodec(t)
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)

	defer cancel()

	caller, e := voip.NewAiortcBridge(nil)

	if e != nil {
		t.Fatal(e)
	}

	defer caller.Close()

	callee, e := voip.NewAiortcBridge(nil)

	if e != nil {
		t.Fatal(e)
	}

	defer callee.Close()

	var count atomic.Int32

	sink := func(_ context.Context, frame voip.AudioFrame) error {
		if frame.SampleRate != 48000 || frame.Channels != 1 {
			t.Error("wrong frame format")
		}

		for _, s := range frame.Samples {
			if s > 1000 || s < -1000 {
				count.Add(1)

				break
			}
		}

		return nil
	}

	left, right := voip.NewFrameAudioSink(sink), voip.NewFrameAudioSink(sink)

	defer left.CloseAsync()

	defer right.CloseAsync()

	caller.OnTrack(func(track voip.RemoteAudioTrack) {
		if e := left.AttachAsync(ctx, track); e != nil {
			t.Error(e)
		}
	})

	callee.OnTrack(func(track voip.RemoteAudioTrack) {
		if e := right.AttachAsync(ctx, track); e != nil {
			t.Error(e)
		}
	})

	for _, peer := range []*voip.AiortcBridge{caller, callee} {
		if e = peer.AddTrack(voip.NewPCMTrack(&tone{}, opus.New, 1)); e != nil {
			t.Fatal(e)
		}
	}

	offer, e := caller.CreateOffer()

	if e != nil {
		t.Fatal(e)
	}

	if e = caller.SetLocalDescription(ctx, offer); e != nil {
		t.Fatal(e)
	}

	if len(caller.LocalCandidates()) == 0 {
		t.Fatal("no gathered candidates")
	}

	if e = callee.SetRemoteDescription(*caller.LocalDescription()); e != nil {
		t.Fatal(e)
	}

	answer, e := callee.CreateAnswer()

	if e != nil {
		t.Fatal(e)
	}

	if e = callee.SetLocalDescription(ctx, answer); e != nil {
		t.Fatal(e)
	}

	if e = caller.SetRemoteDescription(*callee.LocalDescription()); e != nil {
		t.Fatal(e)
	}

	ticker := time.NewTicker(10 * time.Millisecond)

	defer ticker.Stop()

	for count.Load() < 10 {
		select {
		case <-ctx.Done():
			t.Fatalf("audio not exchanged: %d frames; %v", count.Load(), ctx.Err())

		case <-ticker.C:
		}
	}

	t.Logf("received %d non-silent decoded frames", count.Load())
}
