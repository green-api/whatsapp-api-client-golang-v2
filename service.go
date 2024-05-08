package greenapi

import "encoding/json"

type ServiceCategory struct {
	GreenAPI GreenAPIInterface
}

// ------------------------------------------------------------------ DeleteMessage block

type requestDeleteMessage struct {
	ChatId    string `json:"chatId"`
	IdMessage string `json:"idMessage"`
}

func (c ServiceCategory) DeleteMessage(chatId, idMessage string) (interface{}, error) {
	r := &requestDeleteMessage{
		ChatId:    chatId,
		IdMessage: idMessage,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	var payload map[string]interface{}

	if err := json.Unmarshal(jsonData, &payload); err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "deleteMessage", payload)
}
