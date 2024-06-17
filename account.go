package greenapi

import (
	"encoding/json"
)

type AccountCategory struct {
	GreenAPI GreenAPIInterface
}

// ------------------------------------------------------------------ GetSettings block

// Getting settings of an instance.
//
// https://green-api.com/en/docs/api/account/GetSettings/
func (c AccountCategory) GetSettings() (*APIResponse, error) {
	return c.GreenAPI.Request("GET", "getSettings", nil)
}

// ------------------------------------------------------------------ SetSettings block

type RequestSetSettings struct {
	WebhookUrl                        *string `json:"webhookUrl,omitempty"`
	WebhookUrlToken                   *string `json:"webhookUrlToken,omitempty"`
	DelaySendMessagesMilliseconds     *uint    `json:"delaySendMessagesMilliseconds,omitempty"`
	MarkIncomingMessagesReaded        string  `json:"markIncomingMessagesReaded,omitempty"`
	MarkIncomingMessagesReadedOnReply string  `json:"markIncomingMessagesReadedOnReply,omitempty"`
	OutgoingWebhook                   string  `json:"outgoingWebhook,omitempty"`
	OutgoingMessageWebhook            string  `json:"outgoingMessageWebhook,omitempty"`
	OutgoingAPIMessageWebhook         string  `json:"outgoingAPIMessageWebhook,omitempty"`
	StateWebhook                      string  `json:"stateWebhook,omitempty"`
	IncomingWebhook                   string  `json:"incomingWebhook,omitempty"`
	DeviceWebhook                     string  `json:"deviceWebhook,omitempty"`
	KeepOnlineStatus                  string  `json:"keepOnlineStatus,omitempty"`
	PollMessageWebhook                string  `json:"pollMessageWebhook,omitempty"`
	IncomingBlockWebhook              string  `json:"incomingBlockWebhook,omitempty"`
	IncomingCallWebhook               string  `json:"incomingCallWebhook,omitempty"`
}

type SetSettingsOption func(*RequestSetSettings) error

// URL for sending notifications.
func OptionalWebhookUrl(webhookUrl string) SetSettingsOption {
	return func(r *RequestSetSettings) error {
		err := ValidateURL(webhookUrl)
		if err!=nil {
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
func OptionalDelaySendMesssages(delaySendMessagesMilliseconds uint) SetSettingsOption {
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

// Get notifications about the device (phone) and battery level.
func OptionalDeviceWebhook(deviceWebhook bool) SetSettingsOption {
	return func(r *RequestSetSettings) error {
		if deviceWebhook {
			r.DeviceWebhook = "yes"
		} else {
			r.DeviceWebhook = "no"
		}
		return nil
	}
}

// Sets the 'Online' status for your Whatsapp account.
func OptionalKeepOnlineStatus(keepOnlineStatus bool) SetSettingsOption {
	return func(r *RequestSetSettings) error {
		if keepOnlineStatus {
			r.KeepOnlineStatus = "yes"
		} else {
			r.KeepOnlineStatus = "no"
		}
		return nil
	}
}

// Get notifications about the creation of a poll and voting in the poll.
func OptionalPollMessageWebhook(pollMessageWebhook bool) SetSettingsOption {
	return func(r *RequestSetSettings) error {
		if pollMessageWebhook {
			r.PollMessageWebhook = "yes"
		} else {
			r.PollMessageWebhook = "no"
		}
		return nil
	}
}

// Get notifications about adding a chat to the list of blocked contacts.
func OptionalIncomingBlockWebhook(incomingBlockWebhook bool) SetSettingsOption {
	return func(r *RequestSetSettings) error {
		if incomingBlockWebhook {
			r.IncomingBlockWebhook = "yes"
		} else {
			r.IncomingBlockWebhook = "no"
		}
		return nil
	}
}

// Get notifications about incoming call statuses.
func OptionalIncomingCallWebhook(incomingCallWebhook bool) SetSettingsOption {
	return func(r *RequestSetSettings) error {
		if incomingCallWebhook {
			r.IncomingCallWebhook = "yes"
		} else {
			r.IncomingCallWebhook = "no"
		}
		return nil
	}
}

// Applying settings for an instance.
// 
// https://green-api.com/en/docs/api/account/SetSettings/
//
// Add optional arguments by passing these functions:
//  OptionalWebhookUrl(webhookUrl string) <- URL for sending notifications.
//  OptionalWebhookUrlToken(webhookUrlToken string) <- Token to access your notification server.
//  OptionalDelaySendMesssages(delaySendMessagesMilliseconds int) <- Message sending delay. 
//  OptionalMarkIncomingMessagesRead(markIncomingMessagesReaded bool) <- Mark incoming messages as read or not.
//  OptionalMarkIncomingMessagesReadOnReply(markIncomingMessagesReadedOnReply bool) <- Mark incoming messages as read when posting a message to the chat via API.
//  OptionalOutgoingWessebhook(outgoingWebhook bool) <- Get notifications about outgoing messages sending/delivering/reading statuses.
//  OptionalOutgoingMageWebhook(outgoingMessageWebhook bool) <- Get notifications about messages sent from the phone.
//  OptionalOutgoingAPIMessageWebhook(outgoingAPIMessageWebhook bool) <- Get notifications about messages sent from API.
//  OptionalStateWebhook(stateWebhook bool) <- Get notifications about the instance authorization state change.
//  OptionalIncomingWebhook(incomingWebhook bool) <- Get notifications about incoming messages and files.
//  OptionalDeviceWebhook(deviceWebhook bool) <- Get notifications about the device (phone) and battery level.
//  OptionalKeepOnlineStatus(keepOnlineStatus bool) <- Sets the 'Online' status for your Whatsapp account.
//  OptionalPollMessageWebhook(pollMessageWebhook bool) <- Get notifications about the creation of a poll and voting in the poll.
//  OptionalIncomingBlockWebhook(incomingBlockWebhook bool) <- Get notifications about adding a chat to the list of blocked contacts.
//  OptionalIncomingCallWebhook(incomingCallWebhook bool) <- Get notifications about incoming call statuses.
func (c AccountCategory) SetSettings(options ...SetSettingsOption) (*APIResponse, error) {

	r := &RequestSetSettings{}
	for _, o := range options {
		o(r)
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "setSettings", jsonData)
}

// ------------------------------------------------------------------ GetStateInstance block

// Getting state of an instance.
//
// https://green-api.com/en/docs/api/account/GetStateInstance/
func (c AccountCategory) GetStateInstance() (*APIResponse, error) {
	return c.GreenAPI.Request("GET", "getStateInstance", nil)
}

// ------------------------------------------------------------------ GetStatusInstance block

// Getting the status of an instance socket connection with WhatsApp.
//
// https://green-api.com/en/docs/api/account/GetStatusInstance/
func (c AccountCategory) GetStatusInstance() (*APIResponse, error) {
	return c.GreenAPI.Request("GET", "getStatusInstance", nil)
}

// ------------------------------------------------------------------ Reboot block

// Rebooting an instance.
// 
// https://green-api.com/en/docs/api/account/Reboot/
func (c AccountCategory) Reboot() (*APIResponse, error) {
	return c.GreenAPI.Request("GET", "reboot", nil)
}

// ------------------------------------------------------------------ Logout block

// Logging out an instance.
// 
// https://green-api.com/docs/api/account/Logout/
func (c AccountCategory) Logout() (*APIResponse, error) {
	return c.GreenAPI.Request("GET", "logout", nil)
}

// ------------------------------------------------------------------ QR block

// Getting QR code for authorization.
// 
// https://green-api.com/en/docs/api/account/QR/
func (c AccountCategory) QR() (*APIResponse, error) {
	return c.GreenAPI.Request("GET", "qr", nil)
}

// ------------------------------------------------------------------ GetAuthorizationCode block

type RequestGetAuthorizationCode struct {
	PhoneNumber int `json:"phoneNumber"`
}

// Authorize an instance by phone number.
//
// https://green-api.com/en/docs/api/account/GetAuthorizationCode/
func (c AccountCategory) GetAuthorizationCode(phoneNumber int) (*APIResponse, error) {
	r := &RequestGetAuthorizationCode{
		PhoneNumber: phoneNumber,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "getAuthorizationCode", jsonData)
}

// ------------------------------------------------------------------ SetProfilePicture block

type RequestSetProfilePicture struct {
	File string `json:"file"`
}

// Setting a profile picture.
// 
// https://green-api.com/en/docs/api/account/SetProfilePicture/
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

// ------------------------------------------------------------------ GetWaSettings block

// Getting information about the WhatsApp account
//
// https://green-api.com/en/docs/api/account/GetWaSettings/
func (c AccountCategory) GetWaSettings() (*APIResponse, error) {
	return c.GreenAPI.Request("GET", "getWaSettings", nil)
}
