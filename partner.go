package greenapi

import (
	"encoding/json"
	"fmt"
)

type PartnerCategory struct {
	GreenAPIPartner GreenAPIPartnerInterface
}

type RequestCreateInstance struct {
	RequestSetSettings
	Name *string `json:"name,omitempty"`
}

type CreateInstanceOption func(*RequestCreateInstance) error

// Instance name
func OptionalName(name string) CreateInstanceOption {
	return func(r *RequestCreateInstance) error {
		r.Name = &name
		return nil
	}
}

// ------------------------------------------------------------------ GetInstances

// Getting all the account instances created by the partner.
//
// https://green-api.com/v3/docs/partners/getInstances/
func (c PartnerCategory) GetInstances() (*APIResponse, error) {
	return c.GreenAPIPartner.PartnerRequest("GET", "getInstances", nil)
}

// ------------------------------------------------------------------ CreateInstance

// Creating an instance.
//
// https://green-api.com/v3/docs/partners/createInstance/
//
// Add optional arguments by passing these functions:
//
//	OptionalName(name string) <- Name for instance.
//	OptionalWebhookUrl(webhookUrl string) <- URL for sending notifications.
//	OptionalWebhookUrlToken(webhookUrlToken string) <- Token to access your notification server.
//	OptionalDelaySendMesssages(delaySendMessagesMilliseconds int) <- Message sending delay.
//	OptionalMarkIncomingMessagesRead(markIncomingMessagesReaded bool) <- Mark incoming messages as read or not.
//	OptionalMarkIncomingMessagesReadOnReply(markIncomingMessagesReadedOnReply bool) <- Mark incoming messages as read when posting a message to the chat via API.
//	OptionalOutgoingWebhook(outgoingWebhook bool) <- Get notifications about outgoing messages sending/delivering/reading statuses.
//	OptionalOutgoingMessageWebhook(outgoingMessageWebhook bool) <- Get notifications about messages sent from the phone.
//	OptionalOutgoingAPIMessageWebhook(outgoingAPIMessageWebhook bool) <- Get notifications about messages sent from API.
//	OptionalStateWebhook(stateWebhook bool) <- Get notifications about the instance authorization state change.
//	OptionalIncomingWebhook(incomingWebhook bool) <- Get notifications about incoming messages and files.

func (c PartnerCategory) CreateInstance(options ...any) (*APIResponse, error) {
	rCreateInstance := &RequestCreateInstance{}

	for _, o := range options {
		switch v := o.(type) {
		case SetSettingsOption:
			err := v(&rCreateInstance.RequestSetSettings)
			if err != nil {
				return nil, err
			}
		case CreateInstanceOption:
			err := v(rCreateInstance)
			if err != nil {
				return nil, err
			}
		default:
			err := fmt.Errorf("greenapi.GreenAPIPartner.CreateInstance: wrong type in argument: %T", o)
			return nil, err
		}
	}

	jsonData, err := json.Marshal(rCreateInstance)
	if err != nil {
		return nil, err
	}
	return c.GreenAPIPartner.PartnerRequest("POST", "createInstance", jsonData)
}

// ------------------------------------------------------------------ DeleteInstanceAccount

type RequestDeleteInstanceAccount struct {
	IdInstance uint `json:"idInstance"`
}

// Deleting an instance.
//
// https://green-api.com/v3/docs/partners/deleteInstanceAccount/
func (c PartnerCategory) DeleteInstanceAccount(idInstance uint) (*APIResponse, error) {
	r := &RequestDeleteInstanceAccount{
		IdInstance: idInstance,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPIPartner.PartnerRequest("POST", "deleteInstanceAccount", jsonData)
}
