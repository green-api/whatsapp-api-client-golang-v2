package greenapi

import (
	"encoding/json"
	"time"
)

type GreenAPI struct {
	APIURL           string
	MediaURL         string
	IDInstance       string
	APITokenInstance string
	PartnerToken     string
}

type APIResponse struct {
	StatusCode    int             `json:"status_code"`
	StatusMessage []byte          `json:"status_message"`
	Body          json.RawMessage `json:"body"`
	Timestamp     time.Time       `json:"timestamp"`
}
