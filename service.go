package greenapi

import "encoding/json"

type ServiceCategory struct {
	GreenAPI GreenAPIInterface
}

// ------------------------------------------------------------------ CheckAccount

type RequestCheckAccount struct {
	PhoneNumber int `json:"phoneNumber"`
}

// Checking a MAX account availability on a phone number.
//
// https://green-api.com/v3/docs/api/service/CheckAccount/
func (c ServiceCategory) CheckAccount(phoneNumber int) (*APIResponse, error) {
	r := &RequestCheckAccount{
		PhoneNumber: phoneNumber,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "CheckAccount", jsonData)
}

// ------------------------------------------------------------------ GetAvatar

type RequestGetAvatar struct {
	ChatId string `json:"chatId"`
}

// Getting a user or a group chat avatar.
//
// https://green-api.com/v3/docs/api/service/GetAvatar/
func (c ServiceCategory) GetAvatar(chatId string) (*APIResponse, error) {
	err := ValidateChatId(chatId)
	if err != nil {
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
// https://green-api.com/v3/docs/api/service/GetContacts/
func (c ServiceCategory) GetContacts() (*APIResponse, error) {
	return c.GreenAPI.Request("GET", "getContacts", nil)
}

// ------------------------------------------------------------------ GetContactInfo

type RequestGetContactInfo struct {
	ChatId string `json:"chatId"`
}

// Getting information about a contact.
//
// https://green-api.com/v3/docs/api/service/GetContactInfo/
func (c ServiceCategory) GetContactInfo(chatId string) (*APIResponse, error) {
	err := ValidateChatId(chatId)
	if err != nil {
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
