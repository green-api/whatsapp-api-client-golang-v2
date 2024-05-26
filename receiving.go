package greenapi

import (
	"encoding/json"
	"fmt"
)

type ReceivingCategory struct {
	GreenAPI GreenAPIInterface
}

// ------------------------------------------------------------------ ReceiveNotification block

type RequestReceiveNotification struct {
	ReceiveTimeout int    `json:"receiveTimeout"`
	AddUrl         string `json:"addUrl"`
}

type ReceiveNotificationOption func(*RequestReceiveNotification) error

func OptionReceiveTimeout(seconds int) ReceiveNotificationOption {
	return func(r *RequestReceiveNotification) error {
		r.ReceiveTimeout = seconds
		return nil
	}
}

func (c ReceivingCategory) ReceiveNotification(options ...ReceiveNotificationOption) (*APIResponse, error) {

	r := &RequestReceiveNotification{}

	for _, o := range options {
		o(r)
	}

	if r.ReceiveTimeout != 0 {
		r.AddUrl = fmt.Sprintf("?receiveTimeout=%v", r.ReceiveTimeout)
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("GET", "receiveNotification", jsonData, WithGetParams(r.AddUrl))
}

// ------------------------------------------------------------------ DeleteNotification block

type RequestDeleteNotification struct {
	ReceiptId int    `json:"receiptId"`
	AddUrl    string `json:"addUrl"`
}

func (c ReceivingCategory) DeleteNotification(receiptId int) (*APIResponse, error) {
	r := &RequestDeleteNotification{
		ReceiptId: receiptId,
	}

	r.AddUrl = fmt.Sprintf("/%v", receiptId)

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("DELETE", "deleteNotification", jsonData, WithGetParams(r.AddUrl))
}

// ------------------------------------------------------------------ DownloadFile block

type RequestDownloadFile struct {
	ChatId    string `json:"chatId"`
	IdMessage string `json:"idMessage"`
}

func (c ReceivingCategory) DownloadFile(chatId, idMessage string) (*APIResponse, error) {
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
