package voip

import (
	"encoding/json"
	"testing"
)

func TestCandidateWireAndMultiSectionSDP(t *testing.T) {
	sdp := "v=0\r\nm=audio 9 UDP/TLS/RTP/SAVPF 96\r\na=candidate:1 1 udp 1 192.0.2.1 5000 typ host\r\na=mid:audio\r\nm=audio 9 UDP/TLS/RTP/SAVPF 96\r\na=mid:other\r\na=candidate:2 1 udp 2 192.0.2.2 5001 typ host\r\n"
	c := candidatesFromSDP(sdp)

	if len(c) != 2 || *c[0].SDPMid != "audio" || *c[1].SDPMLineIndex != 1 {
		t.Fatal(c)
	}

	for _, v := range c {
		b, e := json.Marshal(v)

		if e != nil {
			t.Fatal(e)
		}

		var round IceCandidate

		if e = json.Unmarshal(b, &round); e != nil || validateCandidate(&round) != nil {
			t.Fatal(e)
		}
	}

	if validateCandidate(&IceCandidate{Candidate: "bad"}) == nil {
		t.Fatal("missing mid accepted")
	}

	if validateCandidate(nil) != nil {
		t.Fatal("end-of-candidates rejected")
	}
}
