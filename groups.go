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

// Creating a group chat.
// 
// https://green-api.com/en/docs/api/groups/CreateGroup/
func (c GroupsCategory) CreateGroup(groupName string, chatIds []string) (*APIResponse, error) {
	for _, chatId := range chatIds {
		err := ValidateChatId(chatId)
		if err!=nil {
			return nil, err
		}
	} 

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

// Change a group chat name.
// 
// https://green-api.com/en/docs/api/groups/UpdateGroupName/
func (c GroupsCategory) UpdateGroupName(groupId, groupName string) (*APIResponse, error) {
	err := ValidateChatId(groupId)
	if err!=nil {
		return nil, err
	}

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

// Getting a group chat data
// 
// https://green-api.com/en/docs/api/groups/GetGroupData/
func (c GroupsCategory) GetGroupData(groupId string) (*APIResponse, error) {
	err := ValidateChatId(groupId)
	if err!=nil {
		return nil, err
	}

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

// Adding a participant to a group chat.
// 
// https://green-api.com/en/docs/api/groups/AddGroupParticipant/
func (c GroupsCategory) AddGroupParticipant(groupId, participantChatId string) (*APIResponse, error) {
	err := ValidateChatId(groupId, participantChatId)
	if err!=nil {
		return nil, err
	}

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

// Removing a participant from a group chat.
// 
// https://green-api.com/en/docs/api/groups/RemoveGroupParticipant/
func (c GroupsCategory) RemoveGroupParticipant(groupId, participantChatId string) (*APIResponse, error) {
	err := ValidateChatId(groupId, participantChatId)
	if err!=nil {
		return nil, err
	}

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

// Setting a group chat participant as an administrator.
// 
// https://green-api.com/en/docs/api/groups/SetGroupAdmin/
func (c GroupsCategory) SetGroupAdmin(groupId, participantChatId string) (*APIResponse, error) {
	err := ValidateChatId(groupId, participantChatId)
	if err!=nil {
		return nil, err
	}
	
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

// Removing a participant from the group chat administration rights.
// 
// https://green-api.com/en/docs/api/groups/RemoveAdmin/
func (c GroupsCategory) RemoveAdmin(groupId, participantChatId string) (*APIResponse, error) {
	err := ValidateChatId(groupId, participantChatId)
	if err!=nil {
		return nil, err
	}

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

// Setting a group picture.
// 
// https://green-api.com/en/docs/api/groups/SetGroupPicture/
func (c GroupsCategory) SetGroupPicture(filepath, groupId string) (*APIResponse, error) {
	err := ValidateChatId(groupId)
	if err!=nil {
		return nil, err
	}

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

// Leaving a group chat.
// 
// https://green-api.com/en/docs/api/groups/LeaveGroup/
func (c GroupsCategory) LeaveGroup(groupId string) (*APIResponse, error) {
	err := ValidateChatId(groupId)
	if err!=nil {
		return nil, err
	}
	
	r := &RequestLeaveGroup{
		GroupId: groupId,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "leaveGroup", jsonData)
}
