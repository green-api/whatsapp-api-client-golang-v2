package greenapi

type QueuesCategory struct {
	GreenAPI GreenAPIInterface
}

// ------------------------------------------------------------------ ShowMessagesQueue block

// Getting a list of messages in the queue to be sent.
//
// https://green-api.com/en/docs/api/queues/ShowMessagesQueue/
func (c QueuesCategory) ShowMessagesQueue() (*APIResponse, error) {
	return c.GreenAPI.Request("GET", "showMessagesQueue", nil)
}

// ------------------------------------------------------------------ ClearMessagesQueue block

// Clearing the queue of messages to be sent.
//
// https://green-api.com/en/docs/api/queues/ClearMessagesQueue/
func (c QueuesCategory) ClearMessagesQueue() (*APIResponse, error) {
	return c.GreenAPI.Request("GET", "clearMessagesQueue", nil)
}
