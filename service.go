package greenapi

import "encoding/json"

type ServiceCategory struct {
	GreenAPI GreenAPIInterface
}

// ------------------------------------------------------------------ CheckWhatsapp block

type RequestCheckWhatsapp struct {
	PhoneNumber int `json:"phoneNumber"`
}

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

// ------------------------------------------------------------------ GetAvatar block

type RequestGetAvatar struct {
	ChatId string `json:"chatId"`
}

func (c ServiceCategory) GetAvatar(chatId string) (*APIResponse, error) {
	r := &RequestGetAvatar{
		ChatId: chatId,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "getAvatar", jsonData)
}

// ------------------------------------------------------------------ GetContacts block

func (c ServiceCategory) GetContacts() (*APIResponse, error) {
	return c.GreenAPI.Request("GET", "getContacts", nil)
}

// ------------------------------------------------------------------ GetContactInfo block

type RequestGetContactInfo struct {
	ChatId string `json:"chatId"`
}

func (c ServiceCategory) GetContactInfo(chatId string) (*APIResponse, error) {
	r := &RequestGetContactInfo{
		ChatId: chatId,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "getContactInfo", jsonData)
}

// ------------------------------------------------------------------ DeleteMessage block

type RequestDeleteMessage struct {
	ChatId    string `json:"chatId"`
	IdMessage string `json:"idMessage"`
}

func (c ServiceCategory) DeleteMessage(chatId, idMessage string) (*APIResponse, error) {
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

// ------------------------------------------------------------------ ArchiveChat block

type RequestArchiveChat struct {
	ChatId string `json:"chatId"`
}

func (c ServiceCategory) ArchiveChat(chatId string) (*APIResponse, error) {
	r := &RequestArchiveChat{
		ChatId: chatId,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "archiveChat", jsonData)
}

func (c ServiceCategory) UnarchiveChat(chatId string) (*APIResponse, error) {
	r := &RequestArchiveChat{
		ChatId: chatId,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "unarchiveChat", jsonData)
}

// ------------------------------------------------------------------ SetDisappearingChat block

type RequestSetDisappearingChat struct {
	ChatId              string `json:"chatId"`
	EphemeralExpiration int    `json:"ephemeralExpiration"`
}

func (c ServiceCategory) SetDisappearingChat(chatId string, ephemeralExpiration int) (*APIResponse, error) {
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
