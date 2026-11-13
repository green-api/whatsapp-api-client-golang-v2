//go:build cgo

package opus

import (
	"math"
	"testing"
)

func TestOpusRoundTripAndFrameValidation(t *testing.T) {
	for _, channels := range []int{1, 2} {
		codec, e := New(channels)

		if e != nil {
			t.Fatal(e)
		}

		pcm := make([]int16, 960*channels)

		for i := range pcm {
			pcm[i] = int16(10000 * math.Sin(2*math.Pi*440*float64(i/channels)/48000))
		}

		packet, e := codec.Encode(pcm, 960)

		if e != nil {
			t.Fatal(e)
		}

		decoded, e := codec.Decode(packet, 5760)

		if e != nil || len(decoded) != len(pcm) {
			t.Fatal(len(decoded), e)
		}

		if _, e = codec.Encode(pcm[:1], 960); e == nil {
			t.Fatal("bad frame accepted")
		}
	}

	if _, e := New(3); e == nil {
		t.Fatal("bad channels accepted")
	}
}
