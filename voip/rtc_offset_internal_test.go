package voip

import "testing"

func TestRTPOffsetKeepsGapsAndWraps(t *testing.T) {
	cases := []struct {
		first, timestamp uint32
		want             int64
	}{
		{1000, 1000, 0},
		{1000, 1960, 960},
		{1000, 1000 + 960*11, 960 * 11},
		{4294966336, 640, 1600},
	}

	for _, c := range cases {
		if got := rtpOffset(c.first, c.timestamp); got != c.want {
			t.Errorf("rtpOffset(%d, %d) = %d, want %d", c.first, c.timestamp, got, c.want)
		}
	}
}
