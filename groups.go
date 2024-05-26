package greenapi

import "encoding/json"

type GroupsCategory struct {
	GreenAPI GreenAPIInterface
}

// ------------------------------------------------------------------ CreateGroup block

type RequestCreateGroup struct {
	GroupName string   `json:"groupName"`
	ChatIds   []string `json:"chatIds"`
}

func (c GroupsCategory) CreateGroup(groupName string, chatIds []string) (*APIResponse, error) {
	r := &RequestCreateGroup{
		GroupName: groupName,
		ChatIds:   chatIds,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	var payload map[string]interface{}

	if err := json.Unmarshal(jsonData, &payload); err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("GET", "createGroup", payload)
}
