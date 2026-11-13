//go:build cgo

// Package opus supplies the PCM codec for voip.NewPCMTrack. Gopus embeds libopus
// on supported platforms, so a C compiler is needed, not a system opus install.
package opus

import (
	"fmt"

	"github.com/green-api/whatsapp-api-client-golang-v2/voip"
	"layeh.com/gopus"
)

type codec struct {
	encoder  *gopus.Encoder
	decoder  *gopus.Decoder
	channels int
}

func New(channels int) (voip.AudioCodec, error) {
	if channels != 1 && channels != 2 {
		return nil, fmt.Errorf("Opus channels must be 1 or 2")
	}

	e, err := gopus.NewEncoder(48000, channels, gopus.Audio)

	if err != nil {
		return nil, err
	}

	e.SetBitrate(96000)
	d, err := gopus.NewDecoder(48000, channels)

	if err != nil {
		return nil, err
	}

	return &codec{
		encoder:  e,
		decoder:  d,
		channels: channels,
	}, nil
}

func (c *codec) Encode(pcm []int16, frameSize int) ([]byte, error) {
	if len(pcm) != frameSize*c.channels {
		return nil, fmt.Errorf("invalid PCM frame size")
	}

	return c.encoder.Encode(pcm, frameSize, 4000)
}

func (c *codec) Decode(packet []byte, frameSize int) ([]int16, error) {
	return c.decoder.Decode(packet, frameSize, false)
}
