package greenapi

import (
	"encoding/json"
)

type AccountCategory struct {
	GreenAPI GreenAPIInterface
}

// ------------------------------------------------------------------ GetSettings

// Getting settings of an instance.
//
// https://green-api.com/v3/docs/api/account/GetSettings/
func (c AccountCategory) GetSettings() (*APIResponse, error) {
	return c.GreenAPI.Request("GET", "getSettings", nil)
}

// ------------------------------------------------------------------ SetSettings

type RequestSetSettings struct {
	WebhookUrl                        *string `json:"webhookUrl,omitempty"`
	WebhookUrlToken                   *string `json:"webhookUrlToken,omitempty"`
	DelaySendMessagesMilliseconds     *uint   `json:"delaySendMessagesMilliseconds,omitempty"`
	MarkIncomingMessagesReaded        string  `json:"markIncomingMessagesReaded,omitempty"`
	MarkIncomingMessagesReadedOnReply string  `json:"markIncomingMessagesReadedOnReply,omitempty"`
	OutgoingWebhook                   string  `json:"outgoingWebhook,omitempty"`
	OutgoingMessageWebhook            string  `json:"outgoingMessageWebhook,omitempty"`
	OutgoingAPIMessageWebhook         string  `json:"outgoingAPIMessageWebhook,omitempty"`
	StateWebhook                      string  `json:"stateWebhook,omitempty"`
	IncomingWebhook                   string  `json:"incomingWebhook,omitempty"`
}

type SetSettingsOption func(*RequestSetSettings) error

// URL for sending notifications.
func OptionalWebhookUrl(webhookUrl string) SetSettingsOption {
	return func(r *RequestSetSettings) error {
		err := ValidateURL(webhookUrl)
		if err != nil {
			return err
		}
		r.WebhookUrl = &webhookUrl
		return nil
	}
}

// Token to access your notification server.
func OptionalWebhookUrlToken(webhookUrlToken string) SetSettingsOption {
	return func(r *RequestSetSettings) error {
		r.WebhookUrlToken = &webhookUrlToken
		return nil
	}
}

// Message sending delay.
func OptionalDelaySendMessages(delaySendMessagesMilliseconds uint) SetSettingsOption {
	return func(r *RequestSetSettings) error {
		r.DelaySendMessagesMilliseconds = &delaySendMessagesMilliseconds
		return nil
	}
}

// Mark incoming messages as read or not.
func OptionalMarkIncomingMessagesRead(markIncomingMessagesReaded bool) SetSettingsOption {
	return func(r *RequestSetSettings) error {
		if markIncomingMessagesReaded {
			r.MarkIncomingMessagesReaded = "yes"
		} else {
			r.MarkIncomingMessagesReaded = "no"
		}
		return nil
	}
}

// Mark incoming messages as read when posting a message to the chat via API.
func OptionalMarkIncomingMessagesReadOnReply(markIncomingMessagesReadedOnReply bool) SetSettingsOption {
	return func(r *RequestSetSettings) error {
		if markIncomingMessagesReadedOnReply {
			r.MarkIncomingMessagesReadedOnReply = "yes"
		} else {
			r.MarkIncomingMessagesReadedOnReply = "no"
		}
		return nil
	}
}

// Get notifications about outgoing messages sending/delivering/reading statuses
func OptionalOutgoingWebhook(outgoingWebhook bool) SetSettingsOption {
	return func(r *RequestSetSettings) error {
		if outgoingWebhook {
			r.OutgoingWebhook = "yes"
		} else {
			r.OutgoingWebhook = "no"
		}
		return nil
	}
}

// Get notifications about messages sent from the phone.
func OptionalOutgoingMessageWebhook(outgoingMessageWebhook bool) SetSettingsOption {
	return func(r *RequestSetSettings) error {
		if outgoingMessageWebhook {
			r.OutgoingMessageWebhook = "yes"
		} else {
			r.OutgoingMessageWebhook = "no"
		}
		return nil
	}
}

// Get notifications about messages sent from API.
func OptionalOutgoingAPIMessageWebhook(outgoingAPIMessageWebhook bool) SetSettingsOption {
	return func(r *RequestSetSettings) error {
		if outgoingAPIMessageWebhook {
			r.OutgoingAPIMessageWebhook = "yes"
		} else {
			r.OutgoingAPIMessageWebhook = "no"
		}
		return nil
	}
}

// Get notifications about the instance authorization state change.
func OptionalStateWebhook(stateWebhook bool) SetSettingsOption {
	return func(r *RequestSetSettings) error {
		if stateWebhook {
			r.StateWebhook = "yes"
		} else {
			r.StateWebhook = "no"
		}
		return nil
	}
}

// Get notifications about incoming messages and files.
func OptionalIncomingWebhook(incomingWebhook bool) SetSettingsOption {
	return func(r *RequestSetSettings) error {
		if incomingWebhook {
			r.IncomingWebhook = "yes"
		} else {
			r.IncomingWebhook = "no"
		}
		return nil
	}
}

// Applying settings for an instance.
//
// https://green-api.com/v3/docs/api/account/SetSettings/
//
// Add optional arguments by passing these functions:
//
//	OptionalWebhookUrl(webhookUrl string) <- URL for sending notifications.
//	OptionalWebhookUrlToken(webhookUrlToken string) <- Token to access your notification server.
//	OptionalDelaySendMesssages(delaySendMessagesMilliseconds int) <- Message sending delay.
//	OptionalMarkIncomingMessagesRead(markIncomingMessagesReaded bool) <- Mark incoming messages as read or not.
//	OptionalMarkIncomingMessagesReadOnReply(markIncomingMessagesReadedOnReply bool) <- Mark incoming messages as read when posting a message to the chat via API.
//	OptionalOutgoingWessebhook(outgoingWebhook bool) <- Get notifications about outgoing messages sending/delivering/reading statuses.
//	OptionalOutgoingMageWebhook(outgoingMessageWebhook bool) <- Get notifications about messages sent from the phone.
//	OptionalOutgoingAPIMessageWebhook(outgoingAPIMessageWebhook bool) <- Get notifications about messages sent from API.
//	OptionalStateWebhook(stateWebhook bool) <- Get notifications about the instance authorization state change.
//	OptionalIncomingWebhook(incomingWebhook bool) <- Get notifications about incoming messages and files.
func (c AccountCategory) SetSettings(options ...SetSettingsOption) (*APIResponse, error) {

	r := &RequestSetSettings{}
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

	return c.GreenAPI.Request("POST", "setSettings", jsonData)
}

// ------------------------------------------------------------------ GetStateInstance

// Getting state of an instance.
//
// https://green-api.com/v3/docs/api/account/GetStateInstance/
func (c AccountCategory) GetStateInstance() (*APIResponse, error) {
	return c.GreenAPI.Request("GET", "getStateInstance", nil)
}

// ------------------------------------------------------------------ GetStatusInstance

// Getting the status of an instance socket connection with MAX.
//
// https://green-api.com/v3/docs/api/account/GetStatusInstance/
func (c AccountCategory) GetStatusInstance() (*APIResponse, error) {
	return c.GreenAPI.Request("GET", "getStatusInstance", nil)
}

// ------------------------------------------------------------------ Reboot

// Rebooting an instance.
//
// https://green-api.com/v3/docs/api/account/Reboot/
func (c AccountCategory) Reboot() (*APIResponse, error) {
	return c.GreenAPI.Request("GET", "reboot", nil)
}

// ------------------------------------------------------------------ Logout

// Logging out an instance.
//
// https://green-api.com/v3/docs/api/account/Logout/
func (c AccountCategory) Logout() (*APIResponse, error) {
	return c.GreenAPI.Request("GET", "logout", nil)
}

// ------------------------------------------------------------------ StartAuthorization

type RequestStartAuthorization struct {
	PhoneNumber int `json:"phoneNumber"`
}

// Start instance authorization
//
// https://green-api.com/v3/en/docs/api/account/StartAuthorization/
func (c AccountCategory) StartAuthorization(phoneNumber int) (*APIResponse, error) {
	r := &RequestStartAuthorization{
		PhoneNumber: phoneNumber,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "startAuthorization", jsonData)
}

// ------------------------------------------------------------------ SendAuthorizationCode

type RequestSendAuthorizationCode struct {
	Code string `json:"code"`
}

// Start instance authorization
//
// https://green-api.com/v3/en/docs/api/account/StartAuthorization/
func (c AccountCategory) SendAuthorizationCode(code string) (*APIResponse, error) {
	r := &RequestSendAuthorizationCode{
		Code: code,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "sendAuthorizationCode", jsonData)
}

// ------------------------------------------------------------------ SetProfilePicture

type RequestSetProfilePicture struct {
	File string `json:"file"`
}

// Setting a profile picture.
//
// https://green-api.com/v3/docs/api/account/SetProfilePicture/
func (c AccountCategory) SetProfilePicture(filepath string) (*APIResponse, error) {
	r := &RequestSetProfilePicture{
		File: filepath,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "setProfilePicture", jsonData, WithFormData(true))
}

// ------------------------------------------------------------------ GetAccountSettings

// Getting information about the MAX account
//
// https://green-api.com/v3/docs/api/account/GetAccountSettings/
func (c AccountCategory) GetAccountSettings() (*APIResponse, error) {
	return c.GreenAPI.Request("GET", "getAccountSettings", nil)
}
