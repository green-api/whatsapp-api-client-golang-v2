package voip_test

import (
	"context"
	"encoding/json"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/green-api/whatsapp-api-client-golang-v2/opus"
	"github.com/green-api/whatsapp-api-client-golang-v2/voip"
)

type tone struct{ pts int64 }

func (t *tone) Recv(ctx context.Context) (voip.AudioFrame, error) {
	if err := ctx.Err(); err != nil {
		return voip.AudioFrame{}, err
	}

	samples := make([]int16, 960)

	for i := range samples {
		samples[i] = int16(10000 * math.Sin(2*math.Pi*440*float64(t.pts+int64(i))/48000))
	}

	frame := voip.AudioFrame{Samples: samples, SampleRate: 48000, Channels: 1, PTS: t.pts}
	t.pts += 960

	return frame, nil
}

func (*tone) Close() error { return nil }

func requireCodec(t *testing.T) {
	t.Helper()

	if _, err := opus.New(1); err != nil {
		t.Skip(err)
	}
}

// Exercise the public API against real REST, WebSocket, ICE, DTLS, SRTP and Opus.
func TestCallsConnectionWithRealServerAndAudio(t *testing.T) {
	requireCodec(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	var mu sync.Mutex

	var timeline []string

	record := func(s string) { mu.Lock(); timeline = append(timeline, s); mu.Unlock() }
	serverAudio := make(chan struct{}, 1)
	clientAudio := make(chan struct{}, 1)
	stopReceived := make(chan struct{}, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "callsDial"):
			record("dial")
			w.WriteHeader(204)

		case strings.Contains(r.URL.Path, "callsGetIceServers"):
			record("ice")
			io.WriteString(w, "[]")

		case strings.Contains(r.URL.Path, "callsRtc"):
			socket, e := websocket.Accept(w, r, nil)

			if e != nil {
				t.Error(e)

				return
			}

			defer socket.CloseNow()

			peer, e := voip.NewAiortcBridge(nil)

			if e != nil {
				t.Error(e)

				return
			}

			defer peer.Close()

			sink := voip.NewFrameAudioSink(func(context.Context, voip.AudioFrame) error {
				select {
				case serverAudio <- struct{}{}:

				default:
				}

				return nil
			})

			defer sink.CloseAsync()

			peer.OnTrack(func(track voip.RemoteAudioTrack) { _ = sink.AttachAsync(ctx, track) })

			if e = peer.AddTrack(voip.NewPCMTrack(&tone{}, opus.New, 1)); e != nil {
				t.Error(e)

				return
			}

			send := func(v any) error { b, _ := json.Marshal(v); return socket.Write(ctx, websocket.MessageText, b) }
			_ = send(map[string]any{"type": "state", "state": map[string]string{"state": "out-call"}})

			for {
				_, raw, e := socket.Read(ctx)

				if e != nil {
					return
				}

				var m struct {
					Type      string                  `json:"type"`
					Offer     voip.SessionDescription `json:"offer"`
					Candidate *voip.IceCandidate      `json:"candidate"`
				}

				if e = json.Unmarshal(raw, &m); e != nil {
					t.Error(e)

					return
				}

				switch m.Type {
				case "offer":
					record("offer")

					if e = peer.SetRemoteDescription(m.Offer); e != nil {
						t.Error(e)

						return
					}

					answer, e := peer.CreateAnswer()

					if e != nil {
						t.Error(e)

						return
					}

					if e = peer.SetLocalDescription(ctx, answer); e != nil {
						t.Error(e)

						return
					}

					candidates := peer.LocalCandidates()

					if len(candidates) > 0 {
						_ = send(map[string]any{"type": "ice-candidate", "candidate": candidates[0]})
					}

					if e = send(map[string]any{"type": "answer", "answer": peer.LocalDescription()}); e != nil {
						return
					}

					_ = send(map[string]any{"type": "state", "state": map[string]string{"state": "on-call"}})

				case "ice-candidate":
					if e = peer.AddIceCandidate(m.Candidate); e != nil {
						t.Error(e)

						return
					}

				case "stop":
					stopReceived <- struct{}{}

					_ = send(map[string]any{"type": "state", "state": map[string]string{"state": "idle"}})
				}
			}

		default:
			t.Errorf("unexpected endpoint %s", r.URL.Path)
			w.WriteHeader(404)
		}
	}))

	defer server.Close()

	api := voip.New(voip.Config{Host: server.URL, IDInstance: "123", APITokenInstance: "test"})

	calls := api.Connect(func(context.Context) (*voip.CallAudio, error) {
		source := &tone{}

		sink := voip.NewFrameAudioSink(func(context.Context, voip.AudioFrame) error {
			select {
			case clientAudio <- struct{}{}:

			default:
			}

			return nil
		})

		return &voip.CallAudio{LocalTrack: voip.NewPCMTrack(source, opus.New, 1), OnRemoteTrack: sink.AttachAsync, Close: func() error { _ = source.Close(); return sink.CloseAsync() }}, nil
	})

	defer calls.CloseAsync()

	if e := calls.OpenAsync(ctx); e != nil {
		t.Fatal(e)
	}

	if e := api.DialAsync(ctx, "123"); e != nil {
		t.Fatal(e)
	}

	if e := calls.StartAudioAsync(ctx); e != nil {
		t.Fatal(e)
	}

	for _, ch := range []<-chan struct{}{serverAudio, clientAudio} {
		select {
		case <-ch:

		case <-ctx.Done():
			t.Fatal("no bidirectional audio", ctx.Err())
		}
	}

	if e := calls.StopAudioAsync(); e != nil {
		t.Fatal(e)
	}

	select {
	case <-stopReceived:

	case <-ctx.Done():
		t.Fatal("stop not received")
	}

	mu.Lock()

	defer mu.Unlock()

	if len(timeline) < 3 || timeline[0] != "dial" || timeline[1] != "ice" || timeline[2] != "offer" {
		t.Fatal(timeline)
	}
}
