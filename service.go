package greenapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type ServiceCategory struct {
	GreenAPI GreenAPIInterface
}

// ------------------------------------------------------------------ CheckWhatsapp

type RequestCheckWhatsapp struct {
	ChatId      *string `json:"chatId,omitempty"`
	PhoneNumber *int    `json:"phoneNumber,omitempty"`
	Force       *bool   `json:"force,omitempty"`
}

type CheckWhatsappOption func(*RequestCheckWhatsapp) error

// Deprecated: Use chatId parameter
func OptionalChatID(chatID string) CheckWhatsappOption {
	return func(r *RequestCheckWhatsapp) error {
		r.ChatId = &chatID
		return nil
	}
}

// Force check without cache. Default is false.
func OptionalForce(force bool) CheckWhatsappOption {
	return func(r *RequestCheckWhatsapp) error {
		r.Force = &force
		return nil
	}
}

// Checking a WhatsApp account availability on a phone number or whatsapp chat ip.
//
// https://green-api.com/en/docs/api/service/CheckWhatsapp/
// The `phone` function argument is optional and is retained for backward compatibility.
// If `phone == 0`, then `WhatsappChatId` must be specified.
// Add optional arguments by passing these functions:
//
//	OptionalChatID(chatID string) <- Specified if the phone == 0.
//	OptionalForce(force bool) <- Force check without cache. Default is false.
func (c ServiceCategory) CheckWhatsapp(phoneNumber int, options ...CheckWhatsappOption) (*APIResponse, error) {
	r := &RequestCheckWhatsapp{
		PhoneNumber: &phoneNumber,
	}

	for _, o := range options {
		err := o(r)
		if err != nil {
			return nil, err
		}
	}
	if (r.PhoneNumber == nil || *r.PhoneNumber == 0) && r.ChatId == nil {
		return nil, errors.New("CheckWhatsapp: phone and whatsappChatID is nil")
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

type RequestGetContacts struct {
	Group *bool `json:"group,omitempty"`
	Count int   `json:"count,omitempty"`
}

type GetContactsOption func(*RequestGetContacts) error

// Filter contacts by type (groups or personal chats).
func OptionalGetContactsGroup(group bool) GetContactsOption {
	return func(r *RequestGetContacts) error {
		r.Group = &group
		return nil
	}
}

// Limit the number of contacts returned.
func OptionalGetContactsCount(count int) GetContactsOption {
	return func(r *RequestGetContacts) error {
		r.Count = count
		return nil
	}
}

// Getting a list of the current account contacts.
//
// https://green-api.com/en/docs/api/service/GetContacts/
//
// Add optional arguments by passing these functions:
//
//	OptionalGetContactsGroup(group bool) <- Filter contacts by type (groups or personal chats).
//	OptionalGetContactsCount(count int) <- Limit the number of contacts returned.
func (c ServiceCategory) GetContacts(options ...GetContactsOption) (*APIResponse, error) {
	r := &RequestGetContacts{}

	for _, o := range options {
		err := o(r)
		if err != nil {
			return nil, err
		}
	}

	var addUrl string
	var params []string
	if r.Group != nil {
		if *r.Group {
			params = append(params, "group=true")
		} else {
			params = append(params, "group=false")
		}
	}
	if r.Count != 0 {
		params = append(params, fmt.Sprintf("count=%d", r.Count))
	}
	if len(params) > 0 {
		addUrl = "?" + strings.Join(params, "&")
	}

	return c.GreenAPI.Request("GET", "getContacts", nil, WithGetParams(addUrl))
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

// ------------------------------------------------------------------ DeleteMessage

type RequestDeleteMessage struct {
	ChatId           string `json:"chatId"`
	IdMessage        string `json:"idMessage"`
	OnlySenderDelete *bool  `json:"onlySenderDelete,omitempty"`
}

type DeleteMessageOption func(*RequestDeleteMessage) error

func OptionalOnlySenderDelete(OnlySenderDelete bool) DeleteMessageOption {
	return func(r *RequestDeleteMessage) error {
		r.OnlySenderDelete = &OnlySenderDelete
		return nil
	}
}

// Deleting a message from a chat.
//
// https://green-api.com/en/docs/api/service/deleteMessage/
func (c ServiceCategory) DeleteMessage(chatId, idMessage string, options ...DeleteMessageOption) (*APIResponse, error) {
	err := ValidateChatId(chatId)
	if err != nil {
		return nil, err
	}

	r := &RequestDeleteMessage{
		ChatId:    chatId,
		IdMessage: idMessage,
	}

	for _, o := range options {
		err := o(r)
		if err != nil {
			return nil, err
		}
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

// ------------------------------------------------------------------ EditMessage

type RequestEditMessage struct {
	ChatId    string `json:"chatId"`
	IdMessage string `json:"idMessage"`
	Message   string `json:"message"`
}

// Edit a message in chat.
//
// https://green-api.com/en/docs/api/service/editMessage/
func (c ServiceCategory) EditMessage(chatId, idMessage, message string) (*APIResponse, error) {
	err := ValidateChatId(chatId)
	if err != nil {
		return nil, err
	}

	r := &RequestEditMessage{
		ChatId:    chatId,
		IdMessage: idMessage,
		Message:   message,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "editMessage", jsonData)
}

// Archiving a chat.
//
// https://green-api.com/en/docs/api/service/archiveChat/
func (c ServiceCategory) ArchiveChat(chatId string) (*APIResponse, error) {
	err := ValidateChatId(chatId)
	if err != nil {
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
	if err != nil {
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
//
//	0 (off), 86400 (24 hours), 604800 (7 days), 7776000 (90 days).
func (c ServiceCategory) SetDisappearingChat(chatId string, ephemeralExpiration int) (*APIResponse, error) {
	err := ValidateChatId(chatId)
	if err != nil {
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

// ------------------------------------------------------------------ SendTyping

type RequestSendTyping struct {
	ChatId     string `json:"chatId"`
	TypingTime int    `json:"typingTime,omitempty"`
	TypingType string `json:"typingType,omitempty"`
}

type SendTypingOption func(*RequestSendTyping) error

// Duration of typing indication in milliseconds. Must be between 1000 and 20000 milliseconds.
func OptionalSendTypingTime(typingTime int) SendTypingOption {
	return func(r *RequestSendTyping) error {
		if typingTime < 1000 || typingTime > 20000 {
			return fmt.Errorf("typingTime must be between 1000 and 20000 milliseconds, got: %d", typingTime)
		}
		r.TypingTime = typingTime
		return nil
	}
}

// Type of typing indication. Available values: "typing", "recording". Default: "typing"
func OptionalSendTypingType(typingType string) SendTypingOption {
	return func(r *RequestSendTyping) error {
		allowedTypes := map[string]bool{
			"typing":    true,
			"recording": true,
		}
		if !allowedTypes[typingType] {
			return fmt.Errorf("invalid typingType: %s. Allowed values: typing, recording", typingType)
		}
		r.TypingType = typingType
		return nil
	}
}

// Sending a typing indicator.
//
// This method allows you to show the user that you are typing a message or recording a voice message.
//
// https://green-api.com/en/docs/api/sending/SendTyping/
//
// Add optional arguments by passing these functions:
//
//	OptionalSendTypingTime(typingTime int) <- Duration of typing indication in milliseconds. Must be between 1000 and 20000 milliseconds.
//	OptionalSendTypingType(typingType string) <- Type of typing indication. Available values: "typing", "recording". Default: "typing"
func (c ServiceCategory) SendTyping(chatId string, options ...SendTypingOption) (*APIResponse, error) {
	err := ValidateChatId(chatId)
	if err != nil {
		return nil, err
	}

	r := &RequestSendTyping{
		ChatId: chatId,
	}

	for _, o := range options {
		err := o(r)
		if err != nil {
			return nil, err
		}
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "sendTyping", jsonData)
}

// ------------------------------------------------------------------ GetChats

type GetChatsOption func(*int) error

// Limit the number of chats returned.
func OptionalChatsCount(count int) GetChatsOption {
	return func(r *int) error {
		*r = count
		return nil
	}
}

// Getting a list of chats.
//
// https://green-api.com/en/docs/api/service/GetChats/
//
// Add optional arguments by passing these functions:
//
//	OptionalChatsCount(count int) <- Limit the number of chats returned.
func (c ServiceCategory) GetChats(options ...GetChatsOption) (*APIResponse, error) {
	var count int

	for _, o := range options {
		err := o(&count)
		if err != nil {
			return nil, err
		}
	}

	var addUrl string
	if count != 0 {
		addUrl = fmt.Sprintf("?count=%d", count)
	}

	return c.GreenAPI.Request("GET", "getChats", nil, WithGetParams(addUrl))
}
