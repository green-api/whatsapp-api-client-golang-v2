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

func OptionReceiveTimeout(timeout int) ReceiveNotificationOption {
	return func(r *RequestReceiveNotification) error {
		r.ReceiveTimeout = timeout
		return nil
	}
}

func (c ReceivingCategory) ReceiveNotification(options ...ReceiveNotificationOption) (interface{}, error) {

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

	var payload map[string]interface{}

	if err := json.Unmarshal(jsonData, &payload); err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("GET", "receiveNotification", payload, WithGetParams(true))
}

// ------------------------------------------------------------------ DeleteNotification block

type RequestDeleteNotification struct {
	ReceiptId int    `json:"receiptId"`
	AddUrl    string `json:"addUrl"`
}

func (c ReceivingCategory) DeleteNotification(receiptId int) (interface{}, error) {
	r := &RequestDeleteNotification{
		ReceiptId: receiptId,
	}

	r.AddUrl = fmt.Sprintf("/%v", receiptId)

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	var payload map[string]interface{}

	if err := json.Unmarshal(jsonData, &payload); err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("DELETE", "deleteNotification", payload, WithGetParams(true))
}
