package greenapi

type AccountCategory struct {
	GreenAPI GreenAPIInterface
}

func (c AccountCategory) GetSettings() (interface{}, error) {
	return c.GreenAPI.Request("GET", "getSettings", nil)
}

func (c AccountCategory) GetInstances() (interface{}, error) {
	return c.GreenAPI.Request("GET", "getInstances", map[string]interface{}{}) //обязательно передавать пустой интерфейс вместо nil
}

// OptionsSetSettings contains available parameters for SetSettings
type requestSetSettings struct {
	WebhookUrl                        string
	WebhookUrlToken                   string
	DelaySendMessagesMilliseconds     int
	MarkIncomingMessagesReaded        string
	MarkIncomingMessagesReadedOnReply string
	OutgoingWebhook                   string
	OutgoingMessageWebhook            string
	OutgoingAPIMessageWebhook         string
	StateWebhook                      string
	IncomingWebhook                   string
	DeviceWebhook                     string
	KeepOnlineStatus                  string
	PollMessageWebhook                string
	IncomingBlockWebhook              string
	IncomingCallWebhook               string
}

type AccountOption func(*requestSetSettings)

func WithWebhookUrl(webhookUrl string) AccountOption {
	return func(r *requestSetSettings) {
		r.WebhookUrl = webhookUrl
	}
}

// func (c AccountCategory) SetSettings(options ...AccountOption) (interface{}, error) {

// 	r := &requestSetSettings{}
// 	for _, o := range options {
// 		o(r)
// 	}

// 	return c.GreenAPI.Request("POST", "setSettings", r)
// }
