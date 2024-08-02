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
}

type GreenAPIInterface interface {
	Request(httpMethod, APImethod string, requestBody []byte, options ...requestOptions) (*APIResponse, error)
}

type GreenAPIPartner struct {
	PartnerToken string
	Email        string
}

type GreenAPIPartnerInterface interface {
	PartnerRequest(HTTPMethod, APIMethod string, requestBody []byte) (*APIResponse, error)
}

type APIResponse struct {
	StatusCode    int             `json:"status_code"`
	StatusMessage []byte          `json:"status_message"`
	Body          json.RawMessage `json:"body"`
	Timestamp     time.Time       `json:"timestamp"`
}
