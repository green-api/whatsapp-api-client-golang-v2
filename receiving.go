package greenapi

import (
	"encoding/json"
	"fmt"
)

type ReceivingCategory struct {
	GreenAPI GreenAPIInterface
}

// ------------------------------------------------------------------ ReceiveNotification

type RequestReceiveNotification struct {
	ReceiveTimeout int `json:"receiveTimeout,omitempty"`
}

type ReceiveNotificationOption func(*RequestReceiveNotification) error

// Notification waiting timeout, takes a value from 5 to 60 seconds (5 seconds by default)
func OptionalReceiveTimeout(seconds int) ReceiveNotificationOption {
	return func(r *RequestReceiveNotification) error {
		r.ReceiveTimeout = seconds
		return nil
	}
}

// Receiving one incoming notification from the notifications queue.
//
// https://green-api.com/v3/docs/api/receiving/technology-http-api/ReceiveNotification/
//
// Add optional arguments by passing these functions:
//
//	OptionalReceiveTimeout(seconds int) <- Notification waiting timeout, takes a value from 5 to 60 seconds (5 seconds by default)
func (c ReceivingCategory) ReceiveNotification(options ...ReceiveNotificationOption) (*APIResponse, error) {

	r := &RequestReceiveNotification{}

	for _, o := range options {
		err := o(r)
		if err != nil {
			return nil, err
		}
	}

	var addUrl string
	if r.ReceiveTimeout != 0 {
		addUrl = fmt.Sprintf("?receiveTimeout=%v", r.ReceiveTimeout)
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("GET", "receiveNotification", jsonData, WithGetParams(addUrl))
}

// ------------------------------------------------------------------ DeleteNotification

type RequestDeleteNotification struct {
	ReceiptId int `json:"receiptId"`
}

// Deleting an incoming notification from the notification queue.
//
// https://green-api.com/v3/docs/api/receiving/technology-http-api/DeleteNotification/
func (c ReceivingCategory) DeleteNotification(receiptId int) (*APIResponse, error) {
	addUrl := fmt.Sprintf("/%v", receiptId)

	return c.GreenAPI.Request("DELETE", "deleteNotification", nil, WithGetParams(addUrl))
}

// ------------------------------------------------------------------ DownloadFile

type RequestDownloadFile struct {
	ChatId    string `json:"chatId"`
	IdMessage string `json:"idMessage"`
}

// Downloading incoming and outgoing files from a chat.
//
// https://green-api.com/v3/docs/api/receiving/files/DownloadFile/
func (c ReceivingCategory) DownloadFile(chatId, idMessage string) (*APIResponse, error) {
	err := ValidateChatId(chatId)
	if err != nil {
		return nil, err
	}

	r := &RequestDownloadFile{
		ChatId:    chatId,
		IdMessage: idMessage,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "downloadFile", jsonData)
}
