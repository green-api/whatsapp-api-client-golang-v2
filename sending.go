package greenapi

import (
	"encoding/json"
	"os"

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
	LinkPreview     bool   `json:"linkPreview,omitempty"`
}

type SendMessageOption func(*RequestSendMessage) error

func OptionQuotedMessageId(quotedMessageId string) SendMessageOption {
	return func(r *RequestSendMessage) error {
		r.QuotedMessageId = quotedMessageId
		return nil
	}
}

func OptionLinkPreview(linkPreview bool) SendMessageOption {
	return func(r *RequestSendMessage) error {
		r.LinkPreview = linkPreview
		return nil
	}
}

// https://green-api.com/en/docs/api/sending/SendMessage/
func (c SendingCategory) SendMessage(chatId, message string, options ...SendMessageOption) (any, error) {

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

	var payload map[string]interface{}

	if err := json.Unmarshal(jsonData, &payload); err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "sendMessage", payload)
}

// ------------------------------------------------------------------ SendPoll block

type PollOption struct {
	OptionName string `json:"optionName"`
}

type RequestSendPoll struct {
	ChatId          string       `json:"chatId"`
	Message         string       `json:"message"`
	PollOptions     []PollOption `json:"options"`
	MultipleAnswers bool         `json:"multipleAnswers,omitempty"`
	QuotedMessageId string       `json:"quotedMessageId,omitempty"`
}

type SendPollOption func(*RequestSendPoll) error

func OptionMultipleAnswers(multipleAnswers bool) SendPollOption {
	return func(r *RequestSendPoll) error {
		r.MultipleAnswers = multipleAnswers
		return nil
	}
}

func OptionPollQuotedMessageId(quotedMessageId string) SendPollOption {
	return func(r *RequestSendPoll) error {
		r.QuotedMessageId = quotedMessageId
		return nil
	}
}

func (c SendingCategory) SendPoll(chatId, message string, pollOptions []string, options ...SendPollOption) (interface{}, error) {

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

	var payload map[string]interface{}

	if err := json.Unmarshal(jsonData, &payload); err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "sendPoll", payload)
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

func OptionCaptionSendUpload(caption string) SendFileByUploadOption {
	return func(r *RequestSendFileByUpload) error {
		r.Caption = caption
		return nil
	}
}

func OptionQuotedMessageIdSendUpload(quotedMessageId string) SendFileByUploadOption {
	return func(r *RequestSendFileByUpload) error {
		r.QuotedMessageId = quotedMessageId
		return nil
	}
}

func (c SendingCategory) SendFileByUpload(chatId, filePath, fileName string, options ...SendFileByUploadOption) (interface{}, error) {

	r := &RequestSendFileByUpload{
		ChatId:   chatId,
		FileName: fileName,
		File:     filePath,
	}

	for _, o := range options {
		o(r)
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	var payload map[string]interface{}

	if err := json.Unmarshal(jsonData, &payload); err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "sendFileByUpload", payload, WithFormData(true), WithMediaHost(true))
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

func OptionCaptionSendUrl(caption string) SendFileByUrlOption {
	return func(r *RequestSendFileByUrl) error {
		r.Caption = caption
		return nil
	}
}

func OptionQuotedMessageIdSendUrl(quotedMessageId string) SendFileByUrlOption {
	return func(r *RequestSendFileByUrl) error {
		r.QuotedMessageId = quotedMessageId
		return nil
	}
}

func (c SendingCategory) SendFileByUrl(chatId, urlFile, fileName string, options ...SendFileByUrlOption) (interface{}, error) {
	r := &RequestSendFileByUrl{
		ChatId:   chatId,
		UrlFile:  urlFile,
		FileName: fileName,
	}

	for _, o := range options {
		o(r)
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	var payload map[string]interface{}

	if err := json.Unmarshal(jsonData, &payload); err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "sendFileByUrl", payload)
}

// ------------------------------------------------------------------ UploadFile block

type RequestUploadFile struct {
	File []byte `json:"file"`
}

func (c SendingCategory) UploadFile(filepath string) (interface{}, error) {

	binary, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	r := &RequestUploadFile{
		File: binary,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	var payload map[string]interface{}

	if err := json.Unmarshal(jsonData, &payload); err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "uploadFile", payload, WithSetMimetype(mimetype.Detect(binary).String()), WithMediaHost(true))
}
