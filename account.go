package greenapi

import (
	"encoding/json"
	"fmt"
)

type AccountCategory struct {
	GreenAPI GreenAPIInterface
}

// ------------------------------------------------------------------ GetSettings block

func (c AccountCategory) GetSettings() (*APIResponse, error) {
	return c.GreenAPI.Request("GET", "getSettings", nil)
}

// ------------------------------------------------------------------ SetSettings block

type RequestSetSettings struct {
	WebhookUrl                        *string `json:"webhookUrl,omitempty"`
	WebhookUrlToken                   *string `json:"webhookUrlToken,omitempty"`
	DelaySendMessagesMilliseconds     *int    `json:"delaySendMessagesMilliseconds,omitempty"`
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

func OptionWebhookUrl(webhookUrl string) SetSettingsOption {
	return func(r *RequestSetSettings) error {
		r.WebhookUrl = &webhookUrl
		return nil
	}
}

func OptionWebhookUrlToken(webhookUrlToken string) SetSettingsOption {
	return func(r *RequestSetSettings) error {
		r.WebhookUrlToken = &webhookUrlToken
		return nil
	}
}

func OptionDelaySendMesssages(delaySendMessagesMilliseconds int) SetSettingsOption {
	return func(r *RequestSetSettings) error {
		r.DelaySendMessagesMilliseconds = &delaySendMessagesMilliseconds
		return nil
	}
}

func OptionMarkIncomingMessagesRead(markIncomingMessagesReaded bool) SetSettingsOption {
	return func(r *RequestSetSettings) error {
		if markIncomingMessagesReaded {
			r.MarkIncomingMessagesReaded = "yes"
		} else {
			r.MarkIncomingMessagesReaded = "no"
		}
		return nil
	}
}

func OptionMarkIncomingMessagesReadOnReply(markIncomingMessagesReadedOnReply bool) SetSettingsOption {
	return func(r *RequestSetSettings) error {
		if markIncomingMessagesReadedOnReply {
			r.MarkIncomingMessagesReadedOnReply = "yes"
		} else {
			r.MarkIncomingMessagesReadedOnReply = "no"
		}
		return nil
	}
}

func OptionOutgoingWebhook(outgoingWebhook bool) SetSettingsOption {
	return func(r *RequestSetSettings) error {
		if outgoingWebhook {
			r.OutgoingWebhook = "yes"
		} else {
			r.OutgoingWebhook = "no"
		}
		return nil
	}
}

func OptionOutgoingMessageWebhook(outgoingMessageWebhook bool) SetSettingsOption {
	return func(r *RequestSetSettings) error {
		if outgoingMessageWebhook {
			r.OutgoingMessageWebhook = "yes"
		} else {
			r.OutgoingMessageWebhook = "no"
		}
		return nil
	}
}

func OptionOutgoingAPIMessageWebhook(outgoingAPIMessageWebhook bool) SetSettingsOption {
	return func(r *RequestSetSettings) error {
		if outgoingAPIMessageWebhook {
			r.OutgoingAPIMessageWebhook = "yes"
		} else {
			r.OutgoingAPIMessageWebhook = "no"
		}
		return nil
	}
}

func OptionStateWebhook(stateWebhook bool) SetSettingsOption {
	return func(r *RequestSetSettings) error {
		if stateWebhook {
			r.StateWebhook = "yes"
		} else {
			r.StateWebhook = "no"
		}
		return nil
	}
}

func OptionIncomingWebhook(incomingWebhook bool) SetSettingsOption {
	return func(r *RequestSetSettings) error {
		if incomingWebhook {
			r.IncomingWebhook = "yes"
		} else {
			r.IncomingWebhook = "no"
		}
		return nil
	}
}

func OptionDeviceWebhook(deviceWebhook bool) SetSettingsOption {
	return func(r *RequestSetSettings) error {
		if deviceWebhook {
			r.DeviceWebhook = "yes"
		} else {
			r.DeviceWebhook = "no"
		}
		return nil
	}
}

func OptionKeepOnlineStatus(keepOnlineStatus bool) SetSettingsOption {
	return func(r *RequestSetSettings) error {
		if keepOnlineStatus {
			r.KeepOnlineStatus = "yes"
		} else {
			r.KeepOnlineStatus = "no"
		}
		return nil
	}
}

func OptionPollMessageWebhook(pollMessageWebhook bool) SetSettingsOption {
	return func(r *RequestSetSettings) error {
		if pollMessageWebhook {
			r.PollMessageWebhook = "yes"
		} else {
			r.PollMessageWebhook = "no"
		}
		return nil
	}
}

func OptionIncomingBlockWebhook(incomingBlockWebhook bool) SetSettingsOption {
	return func(r *RequestSetSettings) error {
		if incomingBlockWebhook {
			r.IncomingBlockWebhook = "yes"
		} else {
			r.IncomingBlockWebhook = "no"
		}
		return nil
	}
}

func OptionIncomingCallWebhook(incomingCallWebhook bool) SetSettingsOption {
	return func(r *RequestSetSettings) error {
		if incomingCallWebhook {
			r.IncomingCallWebhook = "yes"
		} else {
			r.IncomingCallWebhook = "no"
		}
		return nil
	}
}

func (c AccountCategory) SetSettings(options ...SetSettingsOption) (*APIResponse, error) {

	r := &RequestSetSettings{}
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

	fmt.Println(payload)
	return c.GreenAPI.Request("POST", "setSettings", payload)
}
