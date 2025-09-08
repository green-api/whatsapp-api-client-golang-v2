package greenapi

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
)

type SendingCategory struct {
	GreenAPI GreenAPIInterface
}

// ------------------------------------------------------------------ SendMessage

type RequestSendMessage struct {
	ChatId          string `json:"chatId"`
	Message         string `json:"message"`
	QuotedMessageId string `json:"quotedMessageId,omitempty"`
	LinkPreview     *bool  `json:"linkPreview,omitempty"`
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
// https://green-api.com/v3/docs/api/sending/SendMessage/
//
// Add optional arguments by passing these functions:
//
//	OptionalQuotedMessageId(quotedMessageId string) <- Quoted message ID. If present, the message will be sent quoting the specified chat message.
//	OptionalLinkPreview(linkPreview bool) <- The parameter includes displaying a preview and a description of the link. Enabled by default.
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

// ------------------------------------------------------------------ SendFileByUpload

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

// Uploading and sending a file.
//
// https://green-api.com/v3/docs/api/sending/SendFileByUpload/
//
// Add optional arguments by passing these functions:
//
//	OptionalCaptionSendUpload(caption string) <- File caption. Caption added to video, images. The maximum field length is 20000 characters.
//	OptionalQuotedMessageIdSendUpload(quotedMessageId string) <- If specified, the message will be sent quoting the specified chat message.
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

// Sending a file by URL.
//
// https://green-api.com/v3/docs/api/sending/SendFileByUrl/
//
// Add optional arguments by passing these functions:
//
//	OptionalCaptionSendUrl(caption string) <- File caption. Caption added to video, images. The maximum field length is 20000 characters.
//	OptionalQuotedMessageIdSendUrl(quotedMessageId string) <- If specified, the message will be sent quoting the specified chat message.
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
// https://green-api.com/v3/docs/api/sending/UploadFile/
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
