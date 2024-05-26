package greenapi

type QueuesCategory struct {
	GreenAPI GreenAPIInterface
}

// ------------------------------------------------------------------ ShowMessagesQueue block

func (c QueuesCategory) ShowMessagesQueue() (*APIResponse, error) {
	return c.GreenAPI.Request("GET", "showMessagesQueue", nil)
}

// ------------------------------------------------------------------ ClearMessagesQueue block

func (c QueuesCategory) ClearMessagesQueue() (*APIResponse, error) {
	return c.GreenAPI.Request("GET", "clearMessagesQueue", nil)
}
