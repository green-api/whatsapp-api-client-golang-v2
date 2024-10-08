package greenapi

import "encoding/json"

type ServiceCategory struct {
	GreenAPI GreenAPIInterface
}

// ------------------------------------------------------------------ CheckWhatsapp

type RequestCheckWhatsapp struct {
	PhoneNumber int `json:"phoneNumber"`
}

// Checking a WhatsApp account availability on a phone number.
//
// https://green-api.com/en/docs/api/service/CheckWhatsapp/
func (c ServiceCategory) CheckWhatsapp(phoneNumber int) (*APIResponse, error) {
	r := &RequestCheckWhatsapp{
		PhoneNumber: phoneNumber,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "checkWhatsapp", jsonData)
}

// ------------------------------------------------------------------ GetAvatar

type RequestGetAvatar struct {
	ChatId string `json:"chatId"`
}

// Getting a user or a group chat avatar.
// 
// https://green-api.com/en/docs/api/service/GetAvatar/
func (c ServiceCategory) GetAvatar(chatId string) (*APIResponse, error) {
	err := ValidateChatId(chatId)
	if err!=nil {
		return nil, err
	}

	r := &RequestGetAvatar{
		ChatId: chatId,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "getAvatar", jsonData)
}

// ------------------------------------------------------------------ GetContacts

// Getting a list of the current account contacts.
//
// https://green-api.com/en/docs/api/service/GetContacts/
func (c ServiceCategory) GetContacts() (*APIResponse, error) {
	return c.GreenAPI.Request("GET", "getContacts", nil)
}

// ------------------------------------------------------------------ GetContactInfo

type RequestGetContactInfo struct {
	ChatId string `json:"chatId"`
}

// Getting information about a contact.
// 
// https://green-api.com/en/docs/api/service/GetContactInfo/
func (c ServiceCategory) GetContactInfo(chatId string) (*APIResponse, error) {
	err := ValidateChatId(chatId)
	if err!=nil {
		return nil, err
	}

	r := &RequestGetContactInfo{
		ChatId: chatId,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "getContactInfo", jsonData)
}

// ------------------------------------------------------------------ DeleteMessage

type RequestDeleteMessage struct {
	ChatId    string `json:"chatId"`
	IdMessage string `json:"idMessage"`
}

// Deleting a message from a chat.
// 
// https://green-api.com/en/docs/api/service/deleteMessage/
func (c ServiceCategory) DeleteMessage(chatId, idMessage string) (*APIResponse, error) {
	err := ValidateChatId(chatId)
	if err!=nil {
		return nil, err
	}

	r := &RequestDeleteMessage{
		ChatId:    chatId,
		IdMessage: idMessage,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "deleteMessage", jsonData)
}

// ------------------------------------------------------------------ ArchiveChat

type RequestArchiveChat struct {
	ChatId string `json:"chatId"`
}

// Archiving a chat.
// 
// https://green-api.com/en/docs/api/service/archiveChat/
func (c ServiceCategory) ArchiveChat(chatId string) (*APIResponse, error) {
	err := ValidateChatId(chatId)
	if err!=nil {
		return nil, err
	}

	r := &RequestArchiveChat{
		ChatId: chatId,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "archiveChat", jsonData)
}

// Unarchiving a chat.
//
// https://green-api.com/en/docs/api/service/unarchiveChat/
func (c ServiceCategory) UnarchiveChat(chatId string) (*APIResponse, error) {
	err := ValidateChatId(chatId)
	if err!=nil {
		return nil, err
	}

	r := &RequestArchiveChat{
		ChatId: chatId,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "unarchiveChat", jsonData)
}

// ------------------------------------------------------------------ SetDisappearingChat

type RequestSetDisappearingChat struct {
	ChatId              string `json:"chatId"`
	EphemeralExpiration int    `json:"ephemeralExpiration"`
}

// Changing settings of disappearing messages in chats.
// 
// https://green-api.com/en/docs/api/service/SetDisappearingChat/
// 
// The standard settings of the application are to be used: 
//  0 (off), 86400 (24 hours), 604800 (7 days), 7776000 (90 days).
func (c ServiceCategory) SetDisappearingChat(chatId string, ephemeralExpiration int) (*APIResponse, error) {
	err := ValidateChatId(chatId)
	if err!=nil {
		return nil, err
	}

	r := &RequestSetDisappearingChat{
		ChatId:              chatId,
		EphemeralExpiration: ephemeralExpiration,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "setDisappearingChat", jsonData)
}
