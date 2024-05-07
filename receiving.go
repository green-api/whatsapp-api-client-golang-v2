package greenapi

import (
	"encoding/json"
	"fmt"
)

type ReceivingCategory struct {
	GreenAPI GreenAPIInterface
}

// ------------------------------------------------------------------ ReceiveNotification block

type requestReceiveNotification struct {
	ReceiveTimeout int    `json:"receiveTimeout"`
	AddUrl         string `json:"addUrl"`
}

type receiveNotificationOption func(*requestReceiveNotification)

func WithReceiveTimeout(timeout int) receiveNotificationOption {
	return func(r *requestReceiveNotification) {
		r.ReceiveTimeout = timeout
	}
}

func (c ReceivingCategory) ReceiveNotification(options ...receiveNotificationOption) (interface{}, error) {

	r := &requestReceiveNotification{}

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
