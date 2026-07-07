package greenapi

type QueuesCategory struct {
	GreenAPI GreenAPIInterface
}

// ------------------------------------------------------------------ ShowMessagesQueue

// Getting a list of messages in the queue to be sent.
//
// https://green-api.com/en/docs/api/queues/ShowMessagesQueue/
func (c QueuesCategory) ShowMessagesQueue() (*APIResponse, error) {
	return c.GreenAPI.Request("GET", "showMessagesQueue", nil)
}

// ------------------------------------------------------------------ ClearMessagesQueue

// Clearing the queue of messages to be sent.
//
// https://green-api.com/en/docs/api/queues/ClearMessagesQueue/
func (c QueuesCategory) ClearMessagesQueue() (*APIResponse, error) {
	return c.GreenAPI.Request("GET", "clearMessagesQueue", nil)
}

// ------------------------------------------------------------------ GetWebhooksCount

// Getting the number of webhooks in the queue.
//
// https://green-api.com/en/docs/api/queues/GetWebhooksCount/
func (c QueuesCategory) GetWebhooksCount() (*APIResponse, error) {
	return c.GreenAPI.Request("GET", "getWebhooksCount", nil)
}

// ------------------------------------------------------------------ ClearWebhooksQueue

// Clearing the webhooks queue. Rate limited to once per 60 seconds.
//
// https://green-api.com/en/docs/api/queues/ClearWebhooksQueue/
func (c QueuesCategory) ClearWebhooksQueue() (*APIResponse, error) {
	return c.GreenAPI.Request("DELETE", "clearWebhooksQueue", nil)
}
