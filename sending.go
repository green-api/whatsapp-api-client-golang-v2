package greenapi

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
)

type SendingCategory struct {
	GreenAPI GreenAPIInterface
}

// ------------------------------------------------------------------ SendMessage block

type RequestSendMessage struct {
	ChatId          string `json:"chatId"`
	Message         string `json:"message"`
	QuotedMessageId string `json:"quotedMessageId,omitempty"`
	LinkPreview     *bool   `json:"linkPreview,omitempty"`
}

type SendMessageOption func(*RequestSendMessage) error

// Quoted message ID. If present, the message will be sent quoting the specified chat message.
func OptionalQuotedMessageId(quotedMessageId string) SendMessageOption {
	return func(r *RequestSendMessage) error {
		r.QuotedMessageId = quotedMessageId
		return nil
	}
}

// The parameter includes displaying a preview and a description of the link. Enabled by default.
func OptionalLinkPreview(linkPreview bool) SendMessageOption {
	return func(r *RequestSendMessage) error {
		r.LinkPreview = &linkPreview
		return nil
	}
}

// Sending a text message.
//
// https://green-api.com/en/docs/api/sending/SendMessage/
//
// Add optional arguments by passing these functions:
//
//  OptionalQuotedMessageId(quotedMessageId string) <- Quoted message ID. If present, the message will be sent quoting the specified chat message.
//  OptionalLinkPreview(linkPreview bool) <- The parameter includes displaying a preview and a description of the link. Enabled by default.
func (c SendingCategory) SendMessage(chatId, message string, options ...SendMessageOption) (*APIResponse, error) {
	err := ValidateChatId(chatId)
	if err!=nil {
		return nil, err
	}

	err = ValidateMessageLength(message, 20000)
	if err!=nil {
		return nil, err
	}

	r := &RequestSendMessage{
		ChatId:  chatId,
		Message: message,
	}

	for _, o := range options {
		o(r)
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "sendMessage", jsonData)
}

// ------------------------------------------------------------------ SendPoll block

type PollOption struct {
	OptionName string `json:"optionName"`
}

type RequestSendPoll struct {
	ChatId          string       `json:"chatId"`
	Message         string       `json:"message"`
	PollOptions     []PollOption `json:"options"`
	MultipleAnswers *bool         `json:"multipleAnswers,omitempty"`
	QuotedMessageId string       `json:"quotedMessageId,omitempty"`
}

type SendPollOption func(*RequestSendPoll) error

// Allow multiple answers. Disabled by default.
func OptionalMultipleAnswers(multipleAnswers bool) SendPollOption {
	return func(r *RequestSendPoll) error {
		r.MultipleAnswers = &multipleAnswers
		return nil
	}
}

// If specified, the message will be sent quoting the specified chat message.
func OptionalPollQuotedMessageId(quotedMessageId string) SendPollOption {
	return func(r *RequestSendPoll) error {
		r.QuotedMessageId = quotedMessageId
		return nil
	}
}

// Sending messages with a poll.
//
// https://green-api.com/en/docs/api/sending/SendPoll/
//
// Add optional arguments by passing these functions:
//
//  OptionalMultipleAnswers(multipleAnswers bool) <- Allow multiple answers. Disabled by default.
//  OptionalPollQuotedMessageId(quotedMessageId string) <- If specified, the message will be sent quoting the specified chat message.
func (c SendingCategory) SendPoll(chatId, message string, pollOptions []string, options ...SendPollOption) (*APIResponse, error) {
	err := ValidateChatId(chatId)
	if err!=nil {
		return nil, err
	}

	err = ValidateMessageLength(message, 255)
	if err!=nil {
		return nil, err
	}

	if len(pollOptions) < 2 {
		return nil, fmt.Errorf("cannot create less than 2 poll options")
	} else if len(pollOptions) > 12 {
		return nil, fmt.Errorf("cannot create more than 12 poll options")
	}

	//map to check for duplicates in pollOptions 
	seen := make(map[string]bool)

	for _, pollOption := range pollOptions {
		if len(pollOption) > 100 {
			return nil, fmt.Errorf("poll option should not exceed 100 characters")
		}
		if seen[pollOption] {
			return nil, fmt.Errorf("poll options cannot have duplicates: %s", pollOption)
		}
		seen[pollOption] = true
	}

	r := &RequestSendPoll{
		ChatId:  chatId,
		Message: message,
	}

	for _, v := range pollOptions {
		r.PollOptions = append(r.PollOptions, PollOption{OptionName: v})
	}

	for _, o := range options {
		o(r)
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "sendPoll", jsonData)
}

// ------------------------------------------------------------------ SendFileByUpload block

type RequestSendFileByUpload struct {
	ChatId          string `json:"chatId"`
	File            string `json:"file"`
	FileName        string `json:"fileName"`
	Caption         string `json:"caption,omitempty"`
	QuotedMessageId string `json:"quotedMessageId,omitempty"`
}

type SendFileByUploadOption func(*RequestSendFileByUpload) error

// File caption. Caption added to video, images. The maximum field length is 20000 characters.
func OptionalCaptionSendUpload(caption string) SendFileByUploadOption {
	return func(r *RequestSendFileByUpload) error {
		r.Caption = caption
		return nil
	}
}

// If specified, the message will be sent quoting the specified chat message.
func OptionalQuotedMessageIdSendUpload(quotedMessageId string) SendFileByUploadOption {
	return func(r *RequestSendFileByUpload) error {
		r.QuotedMessageId = quotedMessageId
		return nil
	}
}
// Uploading and sending a file.
// 
// https://green-api.com/en/docs/api/sending/SendFileByUpload/
//
// Add optional arguments by passing these functions:
//
//  OptionalCaptionSendUpload(caption string) <- File caption. Caption added to video, images. The maximum field length is 20000 characters.
//  OptionalQuotedMessageIdSendUpload(quotedMessageId string) <- If specified, the message will be sent quoting the specified chat message.
func (c SendingCategory) SendFileByUpload(chatId, filePath, fileName string, options ...SendFileByUploadOption) (*APIResponse, error) {
	err := ValidateChatId(chatId)
	if err!=nil {
		return nil, err
	}

	r := &RequestSendFileByUpload{
		ChatId:   chatId,
		FileName: fileName,
		File:     filePath,
	}

	for _, o := range options {
		o(r)
	}

	if r.Caption != ""{
		err = ValidateMessageLength(r.Caption, 20000)
		if err!=nil {
			return nil, err
		}
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "sendFileByUpload", jsonData, WithFormData(true), WithMediaHost(true))
}

// ------------------------------------------------------------------ SendFileByUrl block

type RequestSendFileByUrl struct {
	ChatId          string `json:"chatId"`
	UrlFile         string `json:"urlFile"`
	FileName        string `json:"fileName"`
	Caption         string `json:"caption,omitempty"`
	QuotedMessageId string `json:"quotedMessageId,omitempty"`
}

type SendFileByUrlOption func(*RequestSendFileByUrl) error

// File caption. Caption added to video, images. The maximum field length is 20000 characters.
func OptionalCaptionSendUrl(caption string) SendFileByUrlOption {
	return func(r *RequestSendFileByUrl) error {
		r.Caption = caption
		return nil
	}
}

// If specified, the message will be sent quoting the specified chat message.
func OptionalQuotedMessageIdSendUrl(quotedMessageId string) SendFileByUrlOption {
	return func(r *RequestSendFileByUrl) error {
		r.QuotedMessageId = quotedMessageId
		return nil
	}
}

// Sending a file by URL.
//
// https://green-api.com/en/docs/api/sending/SendFileByUrl/
//
// Add optional arguments by passing these functions:
//
//  OptionalCaptionSendUrl(caption string) <- File caption. Caption added to video, images. The maximum field length is 20000 characters.
//  OptionalQuotedMessageIdSendUrl(quotedMessageId string) <- If specified, the message will be sent quoting the specified chat message.
func (c SendingCategory) SendFileByUrl(chatId, urlFile, fileName string, options ...SendFileByUrlOption) (*APIResponse, error) {
	err := ValidateChatId(chatId)
	if err!=nil {
		return nil, err
	}

	err = ValidateURL(urlFile)
	if err!=nil {
		return nil, err
	}
	
	r := &RequestSendFileByUrl{
		ChatId:   chatId,
		UrlFile:  urlFile,
		FileName: fileName,
	}

	for _, o := range options {
		o(r)
	}

	if r.Caption != "" {
		err = ValidateMessageLength(r.Caption, 20000)
		if err!=nil {
			return nil, err
		}
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "sendFileByUrl", jsonData)
}

// ------------------------------------------------------------------ UploadFile block

type RequestUploadFile struct {
	File []byte `json:"file"`
	FileName string `json:"fileName"`
}

// Uploading a file to the cloud storage. 
//
// https://green-api.com/en/docs/api/sending/UploadFile/
func (c SendingCategory) UploadFile(filePath string) (*APIResponse, error) {

	binary, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "uploadFile", binary, WithSetMimetype(mtype{
		Mimetype: mimetype.Detect(binary).String(),
		FileName: filepath.Base(filePath),
	}), WithMediaHost(true))
}

// ------------------------------------------------------------------ SendLocation block

type RequestSendLocation struct {
	ChatId          string  `json:"chatId"`
	NameLocation    string  `json:"nameLocation,omitempty"`
	Address         string  `json:"address,omitempty"`
	Latitude        float32 `json:"latitude"`
	Longitude       float32 `json:"longitude"`
	QuotedMessageId string  `json:"quotedMessageId,omitempty"`
}

type SendLocationOption func(*RequestSendLocation) error

// Location name.
func OptionalNameLocation(nameLocation string) SendLocationOption {
	return func(r *RequestSendLocation) error {
		r.NameLocation = nameLocation
		return nil
	}
}

// Location address.
func OptionalAddress(address string) SendLocationOption {
	return func(r *RequestSendLocation) error {
		r.Address = address
		return nil
	}
}

// If specified, the message will be sent quoting the specified chat message.
func OptionalQuotedMessageIdLocation(quotedMessageId string) SendLocationOption {
	return func(r *RequestSendLocation) error {
		r.QuotedMessageId = quotedMessageId
		return nil
	}
}

// Sending a location message.
//
// https://green-api.com/en/docs/api/sending/SendLocation/
//
// Add optional arguments by passing these functions:
//
//  OptionalNameLocation(nameLocation string) <- Location name.
//  OptionalAddress(address string) <- Location address.
//  OptionalQuotedMessageIdLocation(quotedMessageId string) <- If specified, the message will be sent quoting the specified chat message.
func (c SendingCategory) SendLocation(chatId string, latitude, longitude float32, options ...SendLocationOption) (*APIResponse, error) {
	err := ValidateChatId(chatId)
	if err!=nil {
		return nil, err
	}
	
	r := &RequestSendLocation{
		ChatId:    chatId,
		Latitude:  latitude,
		Longitude: longitude,
	}

	for _, o := range options {
		o(r)
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "sendLocation", jsonData)
}

// ------------------------------------------------------------------ SendContact block

type Contact struct {
	PhoneContact int    `json:"phoneContact"` //phoneContact comment
	FirstName    string `json:"firstName,omitempty"`
	MiddleName   string `json:"middleName,omitempty"`
	LastName     string `json:"lastName,omitempty"`
	Company      string `json:"company,omitempty"`
}

type RequestSendContact struct {
	ChatId          string  `json:"chatId"`
	Contact         Contact `json:"contact"`
	QuotedMessageId string  `json:"quotedMessageId,omitempty"`
}

type SendContactOption func(*RequestSendContact) error

// If specified, the message will be sent quoting the specified chat message.
func OptionalQuotedMessageIdContact(quotedMessageId string) SendContactOption {
	return func(r *RequestSendContact) error {
		r.QuotedMessageId = quotedMessageId
		return nil
	}
}

// Sending a contact message.
// 
// https://green-api.com/en/docs/api/sending/SendContact/
//
// Add optional arguments by passing these functions:
//
//  OptionalQuotedMessageIdContact(quotedMessageId string) <- If specified, the message will be sent quoting the specified chat message.
func (c SendingCategory) SendContact(chatId string, contact Contact, options ...SendContactOption) (*APIResponse, error) {
	err := ValidateChatId(chatId)
	if err!=nil {
		return nil, err
	}
	
	r := &RequestSendContact{
		ChatId:  chatId,
		Contact: contact,
	}

	for _, o := range options {
		o(r)
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "sendContact", jsonData)
}

// ------------------------------------------------------------------ ForwardMessages block

type RequestForwardMessages struct {
	ChatId     string   `json:"chatId"`
	ChatIdFrom string   `json:"chatIdFrom"`
	Messages   []string `json:"messages"`
}

// Forwarding messages from one chat to another.
//
// https://green-api.com/en/docs/api/sending/ForwardMessages/
func (c SendingCategory) ForwardMessages(chatId, chatIdFrom string, messages []string) (*APIResponse, error) {
	err := ValidateChatId(chatId)
	if err!=nil {
		return nil, err
	}
	
	r := &RequestForwardMessages{
		ChatId:     chatId,
		ChatIdFrom: chatIdFrom,
		Messages:   messages,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "forwardMessages", jsonData)
}
