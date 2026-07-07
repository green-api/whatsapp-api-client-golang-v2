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

// ------------------------------------------------------------------ SendMessage

type CustomPreview struct {
	Title         string `json:"title,omitempty"`
	Description   string `json:"description,omitempty"`
	Link          string `json:"link,omitempty"`
	UrlFile       string `json:"urlFile,omitempty"`
	JpegThumbnail string `json:"jpegThumbnail,omitempty"`
}

type RequestSendMessage struct {
	ChatId          string         `json:"chatId"`
	Message         string         `json:"message"`
	QuotedMessageId string         `json:"quotedMessageId,omitempty"`
	LinkPreview     *bool          `json:"linkPreview,omitempty"`
	TypePreview     string         `json:"typePreview,omitempty"`
	CustomPreview   *CustomPreview `json:"customPreview,omitempty"`
	TypingTime      int            `json:"typingTime,omitempty"`
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

// Link preview type. Available values: "large", "small".
func OptionalTypePreview(typePreview string) SendMessageOption {
	return func(r *RequestSendMessage) error {
		if typePreview != "large" && typePreview != "small" {
			return fmt.Errorf("invalid typePreview: %s. Allowed values: large, small", typePreview)
		}
		r.TypePreview = typePreview
		return nil
	}
}

// Custom link preview.
func OptionalCustomPreview(customPreview CustomPreview) SendMessageOption {
	return func(r *RequestSendMessage) error {
		r.CustomPreview = &customPreview
		return nil
	}
}

// Duration of typing indication in milliseconds. Must be between 1000 and 20000 milliseconds.
func OptionalMessageTypingTime(typingTime int) SendMessageOption {
	return func(r *RequestSendMessage) error {
		if typingTime != 0 && (typingTime < 1000 || typingTime > 20000) {
			return fmt.Errorf("typingTime must be between 1000 and 20000 milliseconds, got: %d", typingTime)
		}
		r.TypingTime = typingTime
		return nil
	}
}

// Sending a text message.
//
// https://green-api.com/en/docs/api/sending/SendMessage/
//
// Add optional arguments by passing these functions:
//
//	OptionalQuotedMessageId(quotedMessageId string) <- Quoted message ID. If present, the message will be sent quoting the specified chat message.
//	OptionalLinkPreview(linkPreview bool) <- The parameter includes displaying a preview and a description of the link. Enabled by default.
//	OptionalTypePreview(typePreview string) <- Link preview type. Available values: "large", "small".
//	OptionalCustomPreview(customPreview CustomPreview) <- Custom link preview.
//	OptionalMessageTypingTime(typingTime int) <- Duration of typing indication in milliseconds. Must be between 1000 and 20000 milliseconds.
func (c SendingCategory) SendMessage(chatId, message string, options ...SendMessageOption) (*APIResponse, error) {
	err := ValidateChatId(chatId)
	if err != nil {
		return nil, err
	}

	err = ValidateMessageLength(message, 20000)
	if err != nil {
		return nil, err
	}

	r := &RequestSendMessage{
		ChatId:  chatId,
		Message: message,
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

	return c.GreenAPI.Request("POST", "sendMessage", jsonData)
}

// ------------------------------------------------------------------ SendPoll

type PollOption struct {
	OptionName string `json:"optionName"`
}

type RequestSendPoll struct {
	ChatId          string       `json:"chatId"`
	Message         string       `json:"message"`
	PollOptions     []PollOption `json:"options"`
	MultipleAnswers *bool        `json:"multipleAnswers,omitempty"`
	QuotedMessageId string       `json:"quotedMessageId,omitempty"`
	TypingTime      int          `json:"typingTime,omitempty"`
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

// Duration of typing indication in milliseconds. Must be between 1000 and 20000 milliseconds.
func OptionalPollTypingTime(typingTime int) SendPollOption {
	return func(r *RequestSendPoll) error {
		if typingTime != 0 && (typingTime < 1000 || typingTime > 20000) {
			return fmt.Errorf("typingTime must be between 1000 and 20000 milliseconds, got: %d", typingTime)
		}
		r.TypingTime = typingTime
		return nil
	}
}

// Sending messages with a poll.
//
// https://green-api.com/en/docs/api/sending/SendPoll/
//
// Add optional arguments by passing these functions:
//
//	OptionalMultipleAnswers(multipleAnswers bool) <- Allow multiple answers. Disabled by default.
//	OptionalPollQuotedMessageId(quotedMessageId string) <- If specified, the message will be sent quoting the specified chat message.
//	OptionalPollTypingTime(typingTime int) <- Duration of typing indication in milliseconds. Must be between 1000 and 20000 milliseconds.
func (c SendingCategory) SendPoll(chatId, message string, pollOptions []string, options ...SendPollOption) (*APIResponse, error) {
	err := ValidateChatId(chatId)
	if err != nil {
		return nil, err
	}

	err = ValidateMessageLength(message, 255)
	if err != nil {
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
		err := o(r)
		if err != nil {
			return nil, err
		}
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "sendPoll", jsonData)
}

// ------------------------------------------------------------------ SendFileByUpload

type RequestSendFileByUpload struct {
	ChatId          string `json:"chatId"`
	File            string `json:"file"`
	FileName        string `json:"fileName"`
	Caption         string `json:"caption,omitempty"`
	QuotedMessageId string `json:"quotedMessageId,omitempty"`
	TypingTime      int    `json:"typingTime,omitempty"`
	TypingType      string `json:"typingType,omitempty"`
}

type SendFileByUploadOption func(*RequestSendFileByUpload) error

// File caption. Caption added to video, images. The maximum field length is 20000 characters.
func OptionalCaptionSendUpload(caption string) SendFileByUploadOption {
	return func(r *RequestSendFileByUpload) error {
		err := ValidateMessageLength(caption, 20000)
		if err != nil {
			return err
		}
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

// Duration of typing indication in milliseconds. Must be between 1000 and 20000 milliseconds.
func OptionalUploadTypingTime(typingTime int) SendFileByUploadOption {
	return func(r *RequestSendFileByUpload) error {
		if typingTime != 0 && (typingTime < 1000 || typingTime > 20000) {
			return fmt.Errorf("typingTime must be between 1000 and 20000 milliseconds, got: %d", typingTime)
		}
		r.TypingTime = typingTime
		return nil
	}
}

// Type of typing indication. Available values: "typing", "recording". Default: "typing"
func OptionalUploadTypingType(typingType string) SendFileByUploadOption {
	return func(r *RequestSendFileByUpload) error {
		if typingType != "" {
			allowedTypes := map[string]bool{
				"typing":    true,
				"recording": true,
			}
			if !allowedTypes[typingType] {
				return fmt.Errorf("invalid typingType: %s. Allowed values: typing, recording", typingType)
			}
		}
		r.TypingType = typingType
		return nil
	}
}

// Uploading and sending a file.
//
// https://green-api.com/en/docs/api/sending/SendFileByUpload/
//
// Add optional arguments by passing these functions:
//
//	OptionalCaptionSendUpload(caption string) <- File caption. Caption added to video, images. The maximum field length is 20000 characters.
//	OptionalQuotedMessageIdSendUpload(quotedMessageId string) <- If specified, the message will be sent quoting the specified chat message.
//	OptionalUploadTypingTime(typingTime int) <- Duration of typing indication in milliseconds. Must be between 1000 and 20000 milliseconds.
//	OptionalUploadTypingType(typingType string) <- Type of typing indication. Available values: "typing", "recording". Default: "typing"
func (c SendingCategory) SendFileByUpload(chatId, filePath, fileName string, options ...SendFileByUploadOption) (*APIResponse, error) {
	err := ValidateChatId(chatId)
	if err != nil {
		return nil, err
	}

	r := &RequestSendFileByUpload{
		ChatId:   chatId,
		FileName: fileName,
		File:     filePath,
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

	return c.GreenAPI.Request("POST", "sendFileByUpload", jsonData, WithFormData(true), WithMediaHost(true))
}

// ------------------------------------------------------------------ SendFileByUrl

type RequestSendFileByUrl struct {
	ChatId          string `json:"chatId"`
	UrlFile         string `json:"urlFile"`
	FileName        string `json:"fileName"`
	Caption         string `json:"caption,omitempty"`
	QuotedMessageId string `json:"quotedMessageId,omitempty"`
	TypingTime      int    `json:"typingTime,omitempty"`
	TypingType      string `json:"typingType,omitempty"`
}

type SendFileByUrlOption func(*RequestSendFileByUrl) error

// File caption. Caption added to video, images. The maximum field length is 20000 characters.
func OptionalCaptionSendUrl(caption string) SendFileByUrlOption {
	return func(r *RequestSendFileByUrl) error {
		err := ValidateMessageLength(caption, 20000)
		if err != nil {
			return err
		}
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

// Duration of typing indication in milliseconds. Must be between 1000 and 20000 milliseconds.
func OptionalUrlTypingTime(typingTime int) SendFileByUrlOption {
	return func(r *RequestSendFileByUrl) error {
		if typingTime != 0 && (typingTime < 1000 || typingTime > 20000) {
			return fmt.Errorf("typingTime must be between 1000 and 20000 milliseconds, got: %d", typingTime)
		}
		r.TypingTime = typingTime
		return nil
	}
}

// Type of typing indication. Available values: "typing", "recording". Default: "typing"
func OptionalUrlTypingType(typingType string) SendFileByUrlOption {
	return func(r *RequestSendFileByUrl) error {
		if typingType != "" {
			allowedTypes := map[string]bool{
				"typing":    true,
				"recording": true,
			}
			if !allowedTypes[typingType] {
				return fmt.Errorf("invalid typingType: %s. Allowed values: typing, recording", typingType)
			}
		}
		r.TypingType = typingType
		return nil
	}
}

// Sending a file by URL.
//
// https://green-api.com/en/docs/api/sending/SendFileByUrl/
//
// Add optional arguments by passing these functions:
//
//	OptionalCaptionSendUrl(caption string) <- File caption. Caption added to video, images. The maximum field length is 20000 characters.
//	OptionalQuotedMessageIdSendUrl(quotedMessageId string) <- If specified, the message will be sent quoting the specified chat message.
//	OptionalUrlTypingTime(typingTime int) <- Duration of typing indication in milliseconds. Must be between 1000 and 20000 milliseconds.
//	OptionalUrlTypingType(typingType string) <- Type of typing indication. Available values: "typing", "recording". Default: "typing"
func (c SendingCategory) SendFileByUrl(chatId, urlFile, fileName string, options ...SendFileByUrlOption) (*APIResponse, error) {
	err := ValidateChatId(chatId)
	if err != nil {
		return nil, err
	}

	err = ValidateURL(urlFile)
	if err != nil {
		return nil, err
	}

	r := &RequestSendFileByUrl{
		ChatId:   chatId,
		UrlFile:  urlFile,
		FileName: fileName,
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

	return c.GreenAPI.Request("POST", "sendFileByUrl", jsonData)
}

// ------------------------------------------------------------------ UploadFile

type RequestUploadFile struct {
	File     []byte `json:"file"`
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

// ------------------------------------------------------------------ SendLocation

type RequestSendLocation struct {
	ChatId          string  `json:"chatId"`
	NameLocation    string  `json:"nameLocation,omitempty"`
	Address         string  `json:"address,omitempty"`
	Latitude        float32 `json:"latitude"`
	Longitude       float32 `json:"longitude"`
	QuotedMessageId string  `json:"quotedMessageId,omitempty"`
	TypingTime      int     `json:"typingTime,omitempty"`
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

// Duration of typing indication in milliseconds. Must be between 1000 and 20000 milliseconds.
func OptionalLocationTypingTime(typingTime int) SendLocationOption {
	return func(r *RequestSendLocation) error {
		if typingTime != 0 && (typingTime < 1000 || typingTime > 20000) {
			return fmt.Errorf("typingTime must be between 1000 and 20000 milliseconds, got: %d", typingTime)
		}
		r.TypingTime = typingTime
		return nil
	}
}

// Sending a location message.
//
// https://green-api.com/en/docs/api/sending/SendLocation/
//
// Add optional arguments by passing these functions:
//
//	OptionalNameLocation(nameLocation string) <- Location name.
//	OptionalAddress(address string) <- Location address.
//	OptionalQuotedMessageIdLocation(quotedMessageId string) <- If specified, the message will be sent quoting the specified chat message.
//	OptionalLocationTypingTime(typingTime int) <- Duration of typing indication in milliseconds. Must be between 1000 and 20000 milliseconds.
func (c SendingCategory) SendLocation(chatId string, latitude, longitude float32, options ...SendLocationOption) (*APIResponse, error) {
	err := ValidateChatId(chatId)
	if err != nil {
		return nil, err
	}

	r := &RequestSendLocation{
		ChatId:    chatId,
		Latitude:  latitude,
		Longitude: longitude,
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

	return c.GreenAPI.Request("POST", "sendLocation", jsonData)
}

// ------------------------------------------------------------------ SendContact

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
	TypingTime      int     `json:"typingTime,omitempty"`
}

type SendContactOption func(*RequestSendContact) error

// If specified, the message will be sent quoting the specified chat message.
func OptionalQuotedMessageIdContact(quotedMessageId string) SendContactOption {
	return func(r *RequestSendContact) error {
		r.QuotedMessageId = quotedMessageId
		return nil
	}
}

// Duration of typing indication in milliseconds. Must be between 1000 and 20000 milliseconds.
func OptionalContactTypingTime(typingTime int) SendContactOption {
	return func(r *RequestSendContact) error {
		if typingTime != 0 && (typingTime < 1000 || typingTime > 20000) {
			return fmt.Errorf("typingTime must be between 1000 and 20000 milliseconds, got: %d", typingTime)
		}
		r.TypingTime = typingTime
		return nil
	}
}

// Sending a contact message.
//
// https://green-api.com/en/docs/api/sending/SendContact/
//
// Add optional arguments by passing these functions:
//
//	OptionalQuotedMessageIdContact(quotedMessageId string) <- If specified, the message will be sent quoting the specified chat message.
//	OptionalContactTypingTime(typingTime int) <- Duration of typing indication in milliseconds. Must be between 1000 and 20000 milliseconds.
func (c SendingCategory) SendContact(chatId string, contact Contact, options ...SendContactOption) (*APIResponse, error) {
	err := ValidateChatId(chatId)
	if err != nil {
		return nil, err
	}

	r := &RequestSendContact{
		ChatId:  chatId,
		Contact: contact,
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

	return c.GreenAPI.Request("POST", "sendContact", jsonData)
}

// ------------------------------------------------------------------ ForwardMessages

type RequestForwardMessages struct {
	ChatId     string   `json:"chatId"`
	ChatIdFrom string   `json:"chatIdFrom"`
	Messages   []string `json:"messages"`
	TypingTime int      `json:"typingTime,omitempty"`
}

type ForwardMessagesOption func(*RequestForwardMessages) error

// Duration of typing indication in milliseconds. Must be between 1000 and 20000 milliseconds.
func OptionalForwardTypingTime(typingTime int) ForwardMessagesOption {
	return func(r *RequestForwardMessages) error {
		if typingTime != 0 && (typingTime < 1000 || typingTime > 20000) {
			return fmt.Errorf("typingTime must be between 1000 and 20000 milliseconds, got: %d", typingTime)
		}
		r.TypingTime = typingTime
		return nil
	}
}

// Forwarding messages from one chat to another.
//
// https://green-api.com/en/docs/api/sending/ForwardMessages/
//
// Add optional arguments by passing these functions:
//
//	OptionalForwardTypingTime(typingTime int) <- Duration of typing indication in milliseconds. Must be between 1000 and 20000 milliseconds.
func (c SendingCategory) ForwardMessages(chatId, chatIdFrom string, messages []string, options ...ForwardMessagesOption) (*APIResponse, error) {
	err := ValidateChatId(chatId)
	if err != nil {
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

// ------------------------------------------------------------------ SendInteractiveButtons

type InteractiveButton struct {
	Type        string `json:"type"`
	ButtonId    string `json:"buttonId"`
	ButtonText  string `json:"buttonText"`
	CopyCode    string `json:"copyCode,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
	URL         string `json:"url,omitempty"`
}

type RequestSendInteractiveButtons struct {
	ChatId          string              `json:"chatId"`
	Header          string              `json:"header,omitempty"`
	Body            string              `json:"body"`
	Footer          string              `json:"footer,omitempty"`
	Buttons         []InteractiveButton `json:"buttons"`
	QuotedMessageId string              `json:"quotedMessageId,omitempty"`
	TypingTime      int                 `json:"typingTime,omitempty"`
}

type SendInteractiveButtonsOption func(*RequestSendInteractiveButtons) error

// Message header.
func OptionalInteractiveHeader(header string) SendInteractiveButtonsOption {
	return func(r *RequestSendInteractiveButtons) error {
		r.Header = header
		return nil
	}
}

// Message footer.
func OptionalInteractiveFooter(footer string) SendInteractiveButtonsOption {
	return func(r *RequestSendInteractiveButtons) error {
		r.Footer = footer
		return nil
	}
}

// If specified, the message will be sent quoting the specified chat message.
func OptionalInteractiveQuotedMessageId(quotedMessageId string) SendInteractiveButtonsOption {
	return func(r *RequestSendInteractiveButtons) error {
		r.QuotedMessageId = quotedMessageId
		return nil
	}
}

// Duration of typing indication in milliseconds. Must be between 1000 and 20000 milliseconds.
func OptionalInteractiveTypingTime(typingTime int) SendInteractiveButtonsOption {
	return func(r *RequestSendInteractiveButtons) error {
		if typingTime != 0 && (typingTime < 1000 || typingTime > 20000) {
			return fmt.Errorf("typingTime must be between 1000 and 20000 milliseconds, got: %d", typingTime)
		}
		r.TypingTime = typingTime
		return nil
	}
}

// Sending interactive buttons.
//
// https://green-api.com/en/docs/api/sending/SendInteractiveButtons/
//
// Add optional arguments by passing these functions:
//
//	OptionalInteractiveHeader(header string) <- Message header.
//	OptionalInteractiveFooter(footer string) <- Message footer.
//	OptionalInteractiveQuotedMessageId(quotedMessageId string) <- If specified, the message will be sent quoting the specified chat message.
//	OptionalInteractiveTypingTime(typingTime int) <- Duration of typing indication in milliseconds. Must be between 1000 and 20000 milliseconds.
func (c SendingCategory) SendInteractiveButtons(chatId, body string, buttons []InteractiveButton, options ...SendInteractiveButtonsOption) (*APIResponse, error) {
	err := ValidateChatId(chatId)
	if err != nil {
		return nil, err
	}

	err = ValidateMessageLength(body, 255)
	if err != nil {
		return nil, err
	}

	if len(buttons) == 0 {
		return nil, fmt.Errorf("buttons cannot be empty")
	} else if len(buttons) > 10 {
		return nil, fmt.Errorf("cannot create more than 10 buttons")
	}

	for _, button := range buttons {
		if len(button.ButtonId) == 0 {
			return nil, fmt.Errorf(`"buttonId" cannot be empty`)
		}
		if len(button.ButtonText) == 0 {
			return nil, fmt.Errorf(`"buttonText" cannot be empty`)
		}
		if len(button.ButtonText) > 100 {
			return nil, fmt.Errorf(`"buttonText" should not exceed 100 characters`)
		}

		switch button.Type {
		case "copy":
			if len(button.CopyCode) == 0 {
				return nil, fmt.Errorf(`"copyCode" is required for "copy" button type`)
			}
			if len(button.PhoneNumber) > 0 {
				return nil, fmt.Errorf(`"phoneNumber" is not allowed for "copy" button type`)
			}
			if len(button.URL) > 0 {
				return nil, fmt.Errorf(`"url" is not allowed for "copy" button type`)
			}
		case "call":
			if len(button.PhoneNumber) == 0 {
				return nil, fmt.Errorf(`"phoneNumber" is required for "call" button type`)
			}
			if len(button.CopyCode) > 0 {
				return nil, fmt.Errorf(`"copyCode" is not allowed for "call" button type`)
			}
			if len(button.URL) > 0 {
				return nil, fmt.Errorf(`"url" is not allowed for "call" button type`)
			}
		case "url":
			if len(button.URL) == 0 {
				return nil, fmt.Errorf(`"url" is required for "url" button type`)
			}
			if len(button.CopyCode) > 0 {
				return nil, fmt.Errorf(`"copyCode" is not allowed for "url" button type`)
			}
			if len(button.PhoneNumber) > 0 {
				return nil, fmt.Errorf(`"phoneNumber" is not allowed for "url" button type`)
			}
		default:
			return nil, fmt.Errorf(`invalid button type: "%s". Allowed types: "copy", "call", "url"`, button.Type)
		}
	}

	r := &RequestSendInteractiveButtons{
		ChatId:  chatId,
		Body:    body,
		Buttons: buttons,
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

	return c.GreenAPI.Request("POST", "sendInteractiveButtons", jsonData)
}

// ------------------------------------------------------------------ SendInteractiveButtonsReply

type InteractiveReplyButton struct {
	ButtonId   string `json:"buttonId"`
	ButtonText string `json:"buttonText"`
}

type RequestSendInteractiveButtonsReply struct {
	ChatId          string                   `json:"chatId"`
	Header          string                   `json:"header,omitempty"`
	Body            string                   `json:"body"`
	Footer          string                   `json:"footer,omitempty"`
	Buttons         []InteractiveReplyButton `json:"buttons"`
	QuotedMessageId string                   `json:"quotedMessageId,omitempty"`
	TypingTime      int                      `json:"typingTime,omitempty"`
}

type SendInteractiveButtonsReplyOption func(*RequestSendInteractiveButtonsReply) error

// Message header.
func OptionalInteractiveReplyHeader(header string) SendInteractiveButtonsReplyOption {
	return func(r *RequestSendInteractiveButtonsReply) error {
		r.Header = header
		return nil
	}
}

// Message footer.
func OptionalInteractiveReplyFooter(footer string) SendInteractiveButtonsReplyOption {
	return func(r *RequestSendInteractiveButtonsReply) error {
		r.Footer = footer
		return nil
	}
}

// If specified, the message will be sent quoting the specified chat message.
func OptionalInteractiveReplyQuotedMessageId(quotedMessageId string) SendInteractiveButtonsReplyOption {
	return func(r *RequestSendInteractiveButtonsReply) error {
		r.QuotedMessageId = quotedMessageId
		return nil
	}
}

// Duration of typing indication in milliseconds. Must be between 1000 and 20000 milliseconds.
func OptionalInteractiveReplyTypingTime(typingTime int) SendInteractiveButtonsReplyOption {
	return func(r *RequestSendInteractiveButtonsReply) error {
		if typingTime != 0 && (typingTime < 1000 || typingTime > 20000) {
			return fmt.Errorf("typingTime must be between 1000 and 20000 milliseconds, got: %d", typingTime)
		}
		r.TypingTime = typingTime
		return nil
	}
}

// Sending interactive reply buttons.
//
// https://green-api.com/en/docs/api/sending/SendInteractiveButtons/
//
// Add optional arguments by passing these functions:
//
//	OptionalInteractiveReplyHeader(header string) <- Message header.
//	OptionalInteractiveReplyFooter(footer string) <- Message footer.
//	OptionalInteractiveReplyQuotedMessageId(quotedMessageId string) <- If specified, the message will be sent quoting the specified chat message.
//	OptionalInteractiveReplyTypingTime(typingTime int) <- Duration of typing indication in milliseconds. Must be between 1000 and 20000 milliseconds.
func (c SendingCategory) SendInteractiveButtonsReply(chatId, body string, buttons []InteractiveReplyButton, options ...SendInteractiveButtonsReplyOption) (*APIResponse, error) {
	err := ValidateChatId(chatId)
	if err != nil {
		return nil, err
	}

	err = ValidateMessageLength(body, 255)
	if err != nil {
		return nil, err
	}

	if len(buttons) == 0 {
		return nil, fmt.Errorf("buttons cannot be empty")
	} else if len(buttons) > 10 {
		return nil, fmt.Errorf("cannot create more than 10 buttons")
	}

	for _, button := range buttons {
		if len(button.ButtonId) == 0 {
			return nil, fmt.Errorf(`"buttonId" cannot be empty`)
		}
		if len(button.ButtonText) == 0 {
			return nil, fmt.Errorf(`"buttonText" cannot be empty`)
		}
		if len(button.ButtonText) > 100 {
			return nil, fmt.Errorf(`"buttonText" should not exceed 100 characters`)
		}
	}

	r := &RequestSendInteractiveButtonsReply{
		ChatId:  chatId,
		Body:    body,
		Buttons: buttons,
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

	return c.GreenAPI.Request("POST", "sendInteractiveButtonsReply", jsonData)
}
