package greenapi

import (
	"encoding/json"
	"fmt"
)

type JournalsCategory struct {
	GreenAPI GreenAPIInterface
}

// ------------------------------------------------------------------ GetChatHistory block

type RequestGetChatHistory struct {
	ChatId string `json:"chatId"`
	Count  int    `json:"count,omitempty"`
}

type GetChatHistoryOption func(*RequestGetChatHistory) error

func OptionCount(count int) GetChatHistoryOption {
	return func(r *RequestGetChatHistory) error {
		r.Count = count
		return nil
	}
}

func (c JournalsCategory) GetChatHistory(chatId string, options ...GetChatHistoryOption) (*APIResponse, error) {
	r := &RequestGetChatHistory{
		ChatId: chatId,
	}

	for _, o := range options {
		o(r)
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	var payload map[string]interface{}

	if err := json.Unmarshal(jsonData, &payload); err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "getChatHistory", payload)
}

// ------------------------------------------------------------------ GetMessage block

type RequestGetMessage struct {
	ChatId    string `json:"chatId"`
	IdMessage string `json:"idMessage"`
}

func (c JournalsCategory) GetMessage(chatId, idMessage string) (*APIResponse, error) {
	r := &RequestGetMessage{
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

	return c.GreenAPI.Request("POST", "getMessage", payload)
}

// ------------------------------------------------------------------ LastIncomingMessages + LastOutgoingMessages block

type RequestLastMessages struct {
	Minutes int `json:"minutes,omitempty"`
}

type LastMessagesOption func(*RequestLastMessages) error

func OptionMinutes(minutes int) LastMessagesOption {
	return func(r *RequestLastMessages) error {
		r.Minutes = minutes
		return nil
	}
}

func (c JournalsCategory) LastIncomingMessages(options ...LastMessagesOption) (*APIResponse, error) {
	r := &RequestLastMessages{}

	for _, o := range options {
		o(r)
	}

	var addUrl string
	if r.Minutes != 0 {
		addUrl = fmt.Sprintf("?minutes=%v", r.Minutes)
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	var payload map[string]interface{}

	if err := json.Unmarshal(jsonData, &payload); err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("GET", "lastIncomingMessages", payload, WithGetParams(addUrl))
}

func (c JournalsCategory) LastOutgoingMessages(options ...LastMessagesOption) (*APIResponse, error) {
	r := &RequestLastMessages{}

	for _, o := range options {
		o(r)
	}

	var addUrl string
	if r.Minutes != 0 {
		addUrl = fmt.Sprintf("?minutes=%v", r.Minutes)
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	var payload map[string]interface{}

	if err := json.Unmarshal(jsonData, &payload); err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("GET", "lastOutgoingMessages", payload, WithGetParams(addUrl))
}
