package voip

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/pion/webrtc/v4"
)

type SessionDescription = webrtc.SessionDescription
type IceCandidate = webrtc.ICECandidateInit

// IceServer accepts the same string-or-list URLs as aiortc.
type IceServer struct {
	URLs       []string `json:"urls"`
	Username   string   `json:"username,omitempty"`
	Credential string   `json:"credential,omitempty"`
}

func (s *IceServer) UnmarshalJSON(b []byte) error {
	var v struct {
		URLs       json.RawMessage `json:"urls"`
		Username   string          `json:"username"`
		Credential string          `json:"credential"`
	}

	if e := json.Unmarshal(b, &v); e != nil {
		return e
	}

	if e := json.Unmarshal(v.URLs, &s.URLs); e != nil {
		var one string

		if e = json.Unmarshal(v.URLs, &one); e != nil {
			return e
		}

		s.URLs = []string{one}
	}

	s.Username = v.Username
	s.Credential = v.Credential

	return nil
}

type Bridge interface {
	AddTrack(LocalAudioTrack) error
	OnTrack(func(RemoteAudioTrack))
	CreateOffer() (SessionDescription, error)
	SetLocalDescription(context.Context, SessionDescription) error
	LocalDescription() *SessionDescription
	LocalCandidates() []IceCandidate
	SetRemoteDescription(SessionDescription) error
	AddIceCandidate(*IceCandidate) error
	Close() error
}

// AiortcBridge corresponds to Python rtc.AiortcBridge; Pion details live in rtc_adapter.go.
type AiortcBridge struct{ *RTCAdapter }

func NewAiortcBridge(servers []IceServer) (*AiortcBridge, error) {
	a, e := newRTCAdapter(servers)

	if e != nil {
		return nil, e
	}

	return &AiortcBridge{a}, nil
}

func (b *AiortcBridge) LocalCandidates() []IceCandidate {
	d := b.LocalDescription()

	if d == nil {
		return nil
	}

	return candidatesFromSDP(d.SDP)
}

func candidatesFromSDP(sdp string) []IceCandidate {
	var result []IceCandidate

	var candidates []string

	var mid *string

	index := -1

	flush := func() {
		for _, text := range candidates {
			i := uint16(index)
			result = append(result, IceCandidate{Candidate: text, SDPMid: mid, SDPMLineIndex: &i})
		}
	}

	for _, line := range strings.Split(strings.ReplaceAll(sdp, "\r\n", "\n"), "\n") {
		switch {
		case strings.HasPrefix(line, "m="):
			if index >= 0 {
				flush()
			}

			index++

			mid = nil
			candidates = nil

		case strings.HasPrefix(line, "a=mid:"):
			v := strings.TrimPrefix(line, "a=mid:")
			mid = &v

		case strings.HasPrefix(line, "a=candidate:"):
			candidates = append(candidates, strings.TrimPrefix(line, "a="))
		}
	}

	if index >= 0 {
		flush()
	}

	return result
}

func validateCandidate(c *IceCandidate) error {
	if c != nil && c.Candidate != "" && c.SDPMid == nil && c.SDPMLineIndex == nil {
		return errors.New("ICE candidate needs sdpMid or sdpMLineIndex")
	}

	return nil
}
