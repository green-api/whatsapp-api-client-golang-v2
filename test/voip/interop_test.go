package voip_test

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/green-api/whatsapp-api-client-golang-v2/opus"
	"github.com/green-api/whatsapp-api-client-golang-v2/voip"
)

// Set GREEN_VOIP_PYTHON to a Python interpreter with the reference dependencies.
func TestPythonAiortcInterop(t *testing.T) {
	python := os.Getenv("GREEN_VOIP_PYTHON")

	if python == "" {
		t.Skip("set GREEN_VOIP_PYTHON to enable original Python interoperability check")
	}

	requireCodec(t)
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)

	defer cancel()

	cmd := exec.CommandContext(ctx, python, "../../test_voip/python_peer.py")
	cmd.Env = append(os.Environ(), "PYTHONDONTWRITEBYTECODE=1")
	input, e := cmd.StdinPipe()

	if e != nil {
		t.Fatal(e)
	}

	output, e := cmd.StdoutPipe()

	if e != nil {
		t.Fatal(e)
	}

	var stderr bytes.Buffer

	cmd.Stderr = &stderr

	if e = cmd.Start(); e != nil {
		t.Fatal(e)
	}

	defer func() {
		input.Close()
		_ = cmd.Process.Kill()
		_ = cmd.Wait()

		if t.Failed() {
			t.Log(stderr.String())
		}
	}()

	peer, e := voip.NewAiortcBridge(nil)

	if e != nil {
		t.Fatal(e)
	}

	defer peer.Close()

	received := make(chan struct{}, 1)

	sink := voip.NewFrameAudioSink(func(_ context.Context, f voip.AudioFrame) error {
		for _, s := range f.Samples {
			if s > 500 || s < -500 {
				select {
				case received <- struct{}{}:

				default:
				}

				break
			}
		}

		return nil
	})

	defer sink.CloseAsync()

	peer.OnTrack(func(track voip.RemoteAudioTrack) { _ = sink.AttachAsync(ctx, track) })

	if e = peer.AddTrack(voip.NewPCMTrack(&tone{}, opus.New, 1)); e != nil {
		t.Fatal(e)
	}

	offer, e := peer.CreateOffer()

	if e != nil {
		t.Fatal(e)
	}

	if e = peer.SetLocalDescription(ctx, offer); e != nil {
		t.Fatal(e)
	}

	if e = json.NewEncoder(input).Encode(peer.LocalDescription()); e != nil {
		t.Fatal(e)
	}

	scanner := bufio.NewScanner(output)
	scanner.Buffer(make([]byte, 4096), 1<<20)

	if !scanner.Scan() {
		t.Fatal("no Python answer", scanner.Err())
	}

	var answer voip.SessionDescription

	if e = json.Unmarshal(scanner.Bytes(), &answer); e != nil {
		t.Fatal(e)
	}

	if e = peer.SetRemoteDescription(answer); e != nil {
		t.Fatal(e)
	}

	if !scanner.Scan() {
		t.Fatal("Python did not decode Go audio", scanner.Err())
	}

	var report struct {
		Received bool `json:"received_audio"`
	}

	if e = json.Unmarshal(scanner.Bytes(), &report); e != nil || !report.Received {
		t.Fatal("no Python audio confirmation", e)
	}

	select {
	case <-received:

	case <-ctx.Done():
		t.Fatal("Go did not decode Python audio")
	}

	_, _ = input.Write([]byte("stop\n"))
	t.Log("original aiortc and Go exchanged non-silent PCM audio in both directions")
}
