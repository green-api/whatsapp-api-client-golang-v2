// Package example contains the shared WAV audio factory for the two call examples.
package example

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/green-api/whatsapp-api-client-golang-v2/opus"
	"github.com/green-api/whatsapp-api-client-golang-v2/voip"
)

// wavTrack streams PCM WAV data, rather than loading the entire file.
type wavTrack struct {
	mu        sync.Mutex
	file      *os.File
	remaining uint32
	channels  int
	pts       int64
	closed    bool
}

func openWAV(path string) (*wavTrack, error) {
	f, e := os.Open(path)

	if e != nil {
		return nil, e
	}

	fail := func(e error) (*wavTrack, error) { f.Close(); return nil, e }

	var header [12]byte

	if _, e = io.ReadFull(f, header[:]); e != nil {
		return fail(e)
	}

	if string(header[:4]) != "RIFF" || string(header[8:]) != "WAVE" {
		return fail(errors.New("expected RIFF/WAVE"))
	}

	channels := 0

	for {
		var chunk [8]byte

		if _, e = io.ReadFull(f, chunk[:]); e != nil {
			return fail(e)
		}

		size := binary.LittleEndian.Uint32(chunk[4:])

		switch string(chunk[:4]) {
		case "fmt ":
			if size < 16 {
				return fail(errors.New("invalid WAV fmt chunk"))
			}

			var format [16]byte

			if _, e = io.ReadFull(f, format[:]); e != nil {
				return fail(e)
			}

			channels = int(binary.LittleEndian.Uint16(format[2:]))

			if binary.LittleEndian.Uint16(format[:]) != 1 || binary.LittleEndian.Uint32(format[4:]) != 48000 || binary.LittleEndian.Uint16(format[14:]) != 16 || (channels != 1 && channels != 2) {
				return fail(errors.New("WAV must be 48000 Hz, signed 16-bit PCM, mono or stereo"))
			}

			if _, e = f.Seek(int64(size-16)+int64(size%2), io.SeekCurrent); e != nil {
				return fail(e)
			}

		case "data":
			if channels == 0 {
				return fail(errors.New("WAV data precedes fmt"))
			}

			return &wavTrack{file: f, remaining: size, channels: channels}, nil

		default:
			if _, e = f.Seek(int64(size)+int64(size%2), io.SeekCurrent); e != nil {
				return fail(e)
			}
		}
	}
}

func (t *wavTrack) Recv(ctx context.Context) (voip.AudioFrame, error) {
	if e := ctx.Err(); e != nil {
		return voip.AudioFrame{}, e
	}

	t.mu.Lock()

	defer t.mu.Unlock()

	if t.closed || t.remaining == 0 {
		return voip.AudioFrame{}, io.EOF
	}

	size := uint32(960 * t.channels * 2)

	if t.remaining < size {
		size = t.remaining
	}

	raw := make([]byte, size)

	if _, e := io.ReadFull(t.file, raw); e != nil {
		return voip.AudioFrame{}, e
	}

	t.remaining -= size
	pcm := make([]int16, 960*t.channels)

	for i := 0; i+1 < len(raw); i += 2 {
		pcm[i/2] = int16(binary.LittleEndian.Uint16(raw[i:]))
	}

	frame := voip.AudioFrame{Samples: pcm, SampleRate: 48000, Channels: t.channels, PTS: t.pts}
	t.pts += 960

	return frame, nil
}

func (t *wavTrack) Close() error {
	t.mu.Lock()

	defer t.mu.Unlock()

	if t.closed {
		return nil
	}

	t.closed = true

	return t.file.Close()
}

func AudioFactory(path string) voip.AudioFactory {
	return func(context.Context) (*voip.CallAudio, error) {
		source, e := openWAV(path)

		if e != nil {
			return nil, e
		}

		if _, e = opus.New(source.channels); e != nil {
			source.Close()

			return nil, e
		}

		sink := voip.NewFrameAudioSink(func(_ context.Context, f voip.AudioFrame) error {
			fmt.Printf("received %d samples at %d Hz\n", len(f.Samples)/f.Channels, f.SampleRate)

			return nil
		})

		return &voip.CallAudio{LocalTrack: voip.NewPCMTrack(source, opus.New, source.channels), OnRemoteTrack: sink.AttachAsync, Close: func() error { return errors.Join(sink.CloseAsync(), source.Close()) }}, nil
	}
}

func Config() (voip.Config, error) {
	id, token := os.Getenv("GREEN_API_ID"), os.Getenv("GREEN_API_TOKEN")

	if id == "" || token == "" {
		return voip.Config{}, errors.New("set GREEN_API_ID and GREEN_API_TOKEN")
	}

	return voip.Config{Host: os.Getenv("GREEN_API_HOST"), IDInstance: id, APITokenInstance: token}, nil
}
