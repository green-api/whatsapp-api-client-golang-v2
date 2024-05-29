package greenapi

import "encoding/json"

type ReadMarkCategory struct {
	GreenAPI GreenAPIInterface
}

// ------------------------------------------------------------------ ReadChat block

type RequestReadChat struct {
	ChatId    string `json:"chatId"`
	IdMessage string `json:"idMessage,omitempty"`
}

type ReadChatOption func(*RequestReadChat) error

func OptionalIdMessage(idMessage string) ReadChatOption {
	return func(r *RequestReadChat) error {
		r.IdMessage = idMessage
		return nil
	}
}

func (c ReadMarkCategory) ReadChat(chatId string, options ...ReadChatOption) (*APIResponse, error) {
	r := &RequestReadChat{
		ChatId: chatId,
	}

	for _, o := range options {
		o(r)
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "readChat", jsonData)
}
