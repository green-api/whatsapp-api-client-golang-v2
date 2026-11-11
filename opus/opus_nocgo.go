//go:build !cgo

package opus

import (
	"errors"

	"github.com/green-api/whatsapp-api-client-golang-v2/voip"
)

func New(int) (voip.AudioCodec, error) {
	return nil, errors.New("opus needs CGO_ENABLED=1 and a C compiler; alternatively supply your own CodecFactory")
}
