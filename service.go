package greenapi

import "encoding/json"

type ServiceCategory struct {
	GreenAPI GreenAPIInterface
}

// ------------------------------------------------------------------ DeleteMessage block

type RequestDeleteMessage struct {
	ChatId    string `json:"chatId"`
	IdMessage string `json:"idMessage"`
}

func (c ServiceCategory) DeleteMessage(chatId, idMessage string) (*APIResponse, error) {
	r := &RequestDeleteMessage{
		ChatId:    chatId,
		IdMessage: idMessage,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "deleteMessage", jsonData)
}
