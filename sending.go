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

type requestSendMessage struct {
	ChatId          string `json:"chatId"`
	Message         string `json:"message"`
	QuotedMessageId string `json:"quotedMessageId"`
	LinkPreview     bool   `json:"linkPreview"`
}

type sendMessageOption func(*requestSendMessage)

func WithQuotedMessageId(quotedMessageId string) sendMessageOption {
	return func(r *requestSendMessage) {
		r.QuotedMessageId = quotedMessageId
	}
}

func WithLinkPreview(linkPreview bool) sendMessageOption {
	return func(r *requestSendMessage) {
		r.LinkPreview = linkPreview
	}
}

// https://green-api.com/en/docs/api/sending/SendMessage/
func (c SendingCategory) SendMessage(chatId, message string, options ...sendMessageOption) (any, error) {

	r := &requestSendMessage{
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

type requestSendPoll struct {
	ChatId          string       `json:"chatId"`
	Message         string       `json:"message"`
	PollOptions     []PollOption `json:"options"`
	MultipleAnswers bool         `json:"multipleAnswers"`
	QuotedMessageId string       `json:"quotedMessageId"`
}

type sendPollOption func(*requestSendPoll)

func WithMultipleAnswers(multipleAnswers bool) sendPollOption {
	return func(r *requestSendPoll) {
		r.MultipleAnswers = multipleAnswers
	}
}

func WithPollQuotedMessageId(quotedMessageId string) sendPollOption {
	return func(r *requestSendPoll) {
		r.QuotedMessageId = quotedMessageId
	}
}

func (c SendingCategory) SendPoll(chatId, message string, pollOptions []string, options ...sendPollOption) (interface{}, error) {

	r := &requestSendPoll{
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

type requestSendFileByUpload struct {
	ChatId          string `json:"chatId"`
	File            string `json:"file"`
	FileName        string `json:"fileName"`
	Caption         string `json:"caption"`
	QuotedMessageId string `json:"quotedMessageId"`
}

type sendFileByUploadOption func(*requestSendFileByUpload)

func WithCaptionSendUpload(caption string) sendFileByUploadOption {
	return func(r *requestSendFileByUpload) {
		r.Caption = caption
	}
}

func WithQuotedMessageIdSendUpload(quotedMessageId string) sendFileByUploadOption {
	return func(r *requestSendFileByUpload) {
		r.QuotedMessageId = quotedMessageId
	}
}

func (c SendingCategory) SendFileByUpload(chatId, filePath, fileName string, options ...sendFileByUploadOption) (interface{}, error) {

	r := &requestSendFileByUpload{
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

	return c.GreenAPI.Request("POST", "sendFileByUpload", payload)
}

// ------------------------------------------------------------------ SendFileByUrl block

type requestSendFileByUrl struct {
	ChatId          string `json:"chatId"`
	UrlFile         string `json:"urlFile"`
	FileName        string `json:"fileName"`
	Caption         string `json:"caption"`
	QuotedMessageId string `json:"quotedMessageId"`
}

type sendFileByUrlOption func(*requestSendFileByUrl)

func WithCaptionSendUrl(caption string) sendFileByUrlOption {
	return func(r *requestSendFileByUrl) {
		r.Caption = caption
	}
}

func WithQuotedMessageIdSendUrl(quotedMessageId string) sendFileByUrlOption {
	return func(r *requestSendFileByUrl) {
		r.QuotedMessageId = quotedMessageId
	}
}

func (c SendingCategory) SendFileByUrl(chatId, urlFile, fileName string, options ...sendFileByUrlOption) (interface{}, error) {
	r := &requestSendFileByUrl{
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

type requestUploadFile struct {
	File  []byte `json:"file"`
	Mtype string `json:"mtype"`
}

func (c SendingCategory) UploadFile(filepath string) (interface{}, error) {

	binary, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	r := &requestUploadFile{
		File:  binary,
		Mtype: mimetype.Detect(binary).String(),
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	var payload map[string]interface{}

	if err := json.Unmarshal(jsonData, &payload); err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "uploadFile", payload)
}
