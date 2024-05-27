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

	return c.GreenAPI.Request("POST", "createGroup", jsonData)
}

// ------------------------------------------------------------------ UpdateGroupName block

type RequestUpdateGroupName struct {
	GroupId   string `json:"groupId"`
	GroupName string `json:"groupName"`
}

func (c GroupsCategory) UpdateGroupName(groupId, groupName string) (*APIResponse, error) {
	r := &RequestUpdateGroupName{
		GroupId:   groupId,
		GroupName: groupName,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "updateGroupName", jsonData)
}

// ------------------------------------------------------------------ GetGroupData block

type RequestGetGroupData struct {
	GroupId string `json:"groupId"`
}

func (c GroupsCategory) GetGroupData(groupId string) (*APIResponse, error) {
	r := &RequestGetGroupData{
		GroupId: groupId,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "getGroupData", jsonData)
}

// ------------------------------------------------------------------ GroupParticipant block

type RequestModifyGroupParticipant struct {
	GroupId           string `json:"groupId"`
	ParticipantChatId string `json:"participantChatId"`
}

func (c GroupsCategory) AddGroupParticipant(groupId, participantChatId string) (*APIResponse, error) {
	r := &RequestModifyGroupParticipant{
		GroupId:           groupId,
		ParticipantChatId: participantChatId,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "addGroupParticipant", jsonData)
}

func (c GroupsCategory) RemoveGroupParticipant(groupId, participantChatId string) (*APIResponse, error) {
	r := &RequestModifyGroupParticipant{
		GroupId:           groupId,
		ParticipantChatId: participantChatId,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "removeGroupParticipant", jsonData)
}

func (c GroupsCategory) SetGroupAdmin(groupId, participantChatId string) (*APIResponse, error) {
	r := &RequestModifyGroupParticipant{
		GroupId:           groupId,
		ParticipantChatId: participantChatId,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "setGroupAdmin", jsonData)
}

func (c GroupsCategory) RemoveAdmin(groupId, participantChatId string) (*APIResponse, error) {
	r := &RequestModifyGroupParticipant{
		GroupId:           groupId,
		ParticipantChatId: participantChatId,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "removeAdmin", jsonData)
}

// ------------------------------------------------------------------ SetGroupPicture block

type RequestSetGroupPicture struct {
	File    string `json:"file"`
	GroupId string `json:"groupId"`
}

func (c GroupsCategory) SetGroupPicture(filepath, groupId string) (*APIResponse, error) {
	r := &RequestSetGroupPicture{
		File:    filepath,
		GroupId: groupId,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "setGroupPicture", jsonData, WithFormData(true))
}

// ------------------------------------------------------------------ LeaveGroup block

type RequestLeaveGroup struct {
	GroupId string `json:"groupId"`
}

func (c GroupsCategory) LeaveGroup(groupId string) (*APIResponse, error) {
	r := &RequestLeaveGroup{
		GroupId: groupId,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "leaveGroup", jsonData)
}
