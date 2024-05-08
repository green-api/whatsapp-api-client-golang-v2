package greenapi

import (
	"encoding/json"
	"fmt"
)

type AccountCategory struct {
	GreenAPI GreenAPIInterface
}

// ------------------------------------------------------------------ GetSettings block

func (c AccountCategory) GetSettings() (interface{}, error) {
	return c.GreenAPI.Request("GET", "getSettings", nil)
}

// ------------------------------------------------------------------ SetSettings block

type requestSetSettings struct {
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

type setSettingsOption func(*requestSetSettings)

func WithWebhookUrl(webhookUrl string) setSettingsOption {
	return func(r *requestSetSettings) {
		r.WebhookUrl = &webhookUrl
	}
}

func WithWebhookUrlToken(webhookUrlToken string) setSettingsOption {
	return func(r *requestSetSettings) {
		r.WebhookUrlToken = &webhookUrlToken
	}
}

func WithDelaySendMesssages(delaySendMessagesMilliseconds int) setSettingsOption {
	return func(r *requestSetSettings) {
		r.DelaySendMessagesMilliseconds = &delaySendMessagesMilliseconds
	}
}

func WithMarkIncomingMessagesRead(markIncomingMessagesReaded bool) setSettingsOption {
	return func(r *requestSetSettings) {
		if markIncomingMessagesReaded {
			r.MarkIncomingMessagesReaded = "yes"
		} else {
			r.MarkIncomingMessagesReaded = "no"
		}
	}
}

func WithMarkIncomingMessagesReadOnReply(markIncomingMessagesReadedOnReply bool) setSettingsOption {
	return func(r *requestSetSettings) {
		if markIncomingMessagesReadedOnReply {
			r.MarkIncomingMessagesReadedOnReply = "yes"
		} else {
			r.MarkIncomingMessagesReadedOnReply = "no"
		}
	}
}

func WithOutgoingWebhook(outgoingWebhook bool) setSettingsOption {
	return func(r *requestSetSettings) {
		if outgoingWebhook {
			r.OutgoingWebhook = "yes"
		} else {
			r.OutgoingWebhook = "no"
		}
	}
}

func WithOutgoingMessageWebhook(outgoingMessageWebhook bool) setSettingsOption {
	return func(r *requestSetSettings) {
		if outgoingMessageWebhook {
			r.OutgoingMessageWebhook = "yes"
		} else {
			r.OutgoingMessageWebhook = "no"
		}
	}
}

func WithOutgoingAPIMessageWebhook(outgoingAPIMessageWebhook bool) setSettingsOption {
	return func(r *requestSetSettings) {
		if outgoingAPIMessageWebhook {
			r.OutgoingAPIMessageWebhook = "yes"
		} else {
			r.OutgoingAPIMessageWebhook = "no"
		}
	}
}

func WithStateWebhook(stateWebhook bool) setSettingsOption {
	return func(r *requestSetSettings) {
		if stateWebhook {
			r.StateWebhook = "yes"
		} else {
			r.StateWebhook = "no"
		}
	}
}

func WithIncomingWebhook(incomingWebhook bool) setSettingsOption {
	return func(r *requestSetSettings) {
		if incomingWebhook {
			r.IncomingWebhook = "yes"
		} else {
			r.IncomingWebhook = "no"
		}
	}
}

func WithDeviceWebhook(deviceWebhook bool) setSettingsOption {
	return func(r *requestSetSettings) {
		if deviceWebhook {
			r.DeviceWebhook = "yes"
		} else {
			r.DeviceWebhook = "no"
		}
	}
}

func WithKeepOnlineStatus(keepOnlineStatus bool) setSettingsOption {
	return func(r *requestSetSettings) {
		if keepOnlineStatus {
			r.KeepOnlineStatus = "yes"
		} else {
			r.KeepOnlineStatus = "no"
		}
	}
}

func WithPollMessageWebhook(pollMessageWebhook bool) setSettingsOption {
	return func(r *requestSetSettings) {
		if pollMessageWebhook {
			r.PollMessageWebhook = "yes"
		} else {
			r.PollMessageWebhook = "no"
		}
	}
}

func WithIncomingBlockWebhook(incomingBlockWebhook bool) setSettingsOption {
	return func(r *requestSetSettings) {
		if incomingBlockWebhook {
			r.IncomingBlockWebhook = "yes"
		} else {
			r.IncomingBlockWebhook = "no"
		}
	}
}

func WithIncomingCallWebhook(incomingCallWebhook bool) setSettingsOption {
	return func(r *requestSetSettings) {
		if incomingCallWebhook {
			r.IncomingCallWebhook = "yes"
		} else {
			r.IncomingCallWebhook = "no"
		}
	}
}

func (c AccountCategory) SetSettings(options ...setSettingsOption) (interface{}, error) {

	r := &requestSetSettings{}
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
