package greenapi

type GreenAPI struct {
	Host             string
	MediaHost        string
	IDInstance       string
	APITokenInstance string
	PartnerToken     string
}

type ApiResponse struct {
	StatusCode int         `json:"status_code"`
	Body       interface{} `json:"body"`
	Timestamp  string      `json:"timestamp"`
}
