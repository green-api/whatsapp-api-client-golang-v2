package greenapi

import "encoding/json"

type GreenAPI struct {
	APIURL           string
	MediaURL         string
	IDInstance       string
	APITokenInstance string
	PartnerToken     string
}

type APIResponse struct {
	StatusCode int             `json:"status_code"`
	Body       json.RawMessage `json:"body"`
	Timestamp  string          `json:"timestamp"`
}
