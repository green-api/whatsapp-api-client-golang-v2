package greenapi

import (
	"encoding/json"
	"fmt"
)

type JournalsCategory struct {
	GreenAPI GreenAPIInterface
}

// ------------------------------------------------------------------ GetChatHistory

type RequestGetChatHistory struct {
	ChatId string `json:"chatId"`
	Count  int    `json:"count,omitempty"`
}

type GetChatHistoryOption func(*RequestGetChatHistory) error

// The number of messages to get. The default is 100
func OptionalCount(count int) GetChatHistoryOption {
	return func(r *RequestGetChatHistory) error {
		r.Count = count
		return nil
	}
}

// Getting a chat messages history.
//
// https://green-api.com/v3/docs/api/journals/GetChatHistory/
//
// Add optional arguments by passing these functions:
//
//	OptionalCount(count int) <- The number of messages to get. The default is 100
func (c JournalsCategory) GetChatHistory(chatId string, options ...GetChatHistoryOption) (*APIResponse, error) {
	err := ValidateChatId(chatId)
	if err != nil {
		return nil, err
	}

	r := &RequestGetChatHistory{
		ChatId: chatId,
	}

	for _, o := range options {
		err := o(r)
		if err != nil {
			return nil, err
		}
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "getChatHistory", jsonData)
}

// ------------------------------------------------------------------ GetMessage

type RequestGetMessage struct {
	ChatId    string `json:"chatId"`
	IdMessage string `json:"idMessage"`
}

// Getting a message information.
//
// https://green-api.com/v3/docs/api/journals/GetMessage/
func (c JournalsCategory) GetMessage(chatId, idMessage string) (*APIResponse, error) {
	err := ValidateChatId(chatId)
	if err != nil {
		return nil, err
	}

	r := &RequestGetMessage{
		ChatId:    chatId,
		IdMessage: idMessage,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "getMessage", jsonData)
}

// ------------------------------------------------------------------ LastIncomingMessages + LastOutgoingMessages

type RequestLastMessages struct {
	Minutes int `json:"minutes,omitempty"`
}

type LastMessagesOption func(*RequestLastMessages) error

// Time in minutes for which the messages should be displayed (default is 1440 minutes)
func OptionalMinutes(minutes int) LastMessagesOption {
	return func(r *RequestLastMessages) error {
		r.Minutes = minutes
		return nil
	}
}

// Getting the last incoming messages of the account.
//
// https://green-api.com/v3/docs/api/journals/LastIncomingMessages/
//
// Add optional arguments by passing these functions:
//
//	OptionalMinutes(minutes int) <- Time in minutes for which the messages should be displayed (default is 1440 minutes)
func (c JournalsCategory) LastIncomingMessages(options ...LastMessagesOption) (*APIResponse, error) {
	r := &RequestLastMessages{}

	for _, o := range options {
		err := o(r)
		if err != nil {
			return nil, err
		}
	}

	var addUrl string
	if r.Minutes != 0 {
		addUrl = fmt.Sprintf("?minutes=%v", r.Minutes)
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("GET", "lastIncomingMessages", jsonData, WithGetParams(addUrl))
}

// Getting the last outgoung messages of the account.
//
// https://green-api.com/v3/docs/api/journals/LastOutgoingMessages/
//
//	OptionalMinutes(minutes int) <- Time in minutes for which the messages should be displayed (default is 1440 minutes)
func (c JournalsCategory) LastOutgoingMessages(options ...LastMessagesOption) (*APIResponse, error) {
	r := &RequestLastMessages{}

	for _, o := range options {
		err := o(r)
		if err != nil {
			return nil, err
		}
	}

	var addUrl string
	if r.Minutes != 0 {
		addUrl = fmt.Sprintf("?minutes=%v", r.Minutes)
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("GET", "lastOutgoingMessages", jsonData, WithGetParams(addUrl))
}
