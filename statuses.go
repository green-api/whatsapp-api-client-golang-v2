package greenapi

import (
	"encoding/json"
	"fmt"
)

type StatusesCategory struct {
	GreenAPI GreenAPIInterface
}

// ------------------------------------------------------------------ SendTextStatus block

type RequestSendTextStatus struct {
	Message         string   `json:"message"`
	BackgroundColor string   `json:"backgroundColor,omitempty"`
	Font            string   `json:"font,omitempty"`
	Participants    []string `json:"participants,omitempty"`
}

type SendTextStatusOption func(*RequestSendTextStatus) error

func OptionalBackgroundColorText(backgroundColor string) SendTextStatusOption {
	return func(r *RequestSendTextStatus) error {
		r.BackgroundColor = backgroundColor
		return nil
	}
}

func OptionalFont(font string) SendTextStatusOption {
	return func(r *RequestSendTextStatus) error {
		r.Font = font
		return nil
	}
}

func OptionalParticipantsTextStatus(participants []string) SendTextStatusOption {
	return func(r *RequestSendTextStatus) error {
		r.Participants = participants
		return nil
	}
}

func (c StatusesCategory) SendTextStatus(message string, options ...SendTextStatusOption) (*APIResponse, error) {
	r := &RequestSendTextStatus{
		Message: message,
	}

	for _, o := range options {
		o(r)
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "sendTextStatus", jsonData)
}

// ------------------------------------------------------------------ SendVoiceStatus block

type RequestSendVoiceStatus struct {
	UrlFile         string   `json:"urlFile"`
	FileName        string   `json:"fileName"`
	BackgroundColor string   `json:"backgroundColor,omitempty"`
	Participants    []string `json:"participants,omitempty"`
}

type SendVoiceStatusOption func(*RequestSendVoiceStatus) error

func OptionalBackgroundColorVoice(backgroundColor string) SendVoiceStatusOption {
	return func(r *RequestSendVoiceStatus) error {
		r.BackgroundColor = backgroundColor
		return nil
	}
}

func OptionalParticipantsVoiceStatus(participants []string) SendVoiceStatusOption {
	return func(r *RequestSendVoiceStatus) error {
		r.Participants = participants
		return nil
	}
}

func (c StatusesCategory) SendVoiceStatus(urlFile, fileName string, options ...SendVoiceStatusOption) (*APIResponse, error) {
	r := &RequestSendVoiceStatus{
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

	return c.GreenAPI.Request("POST", "sendVoiceStatus", jsonData)
}

// ------------------------------------------------------------------ SendMediaStatus block

type RequestSendMediaStatus struct {
	UrlFile      string   `json:"urlFile"`
	FileName     string   `json:"fileName"`
	Caption      string   `json:"caption,omitempty"`
	Participants []string `json:"participants,omitempty"`
}

type SendMediaStatusOption func(*RequestSendMediaStatus) error

func OptionalCaptionMediaStatus(caption string) SendMediaStatusOption {
	return func(r *RequestSendMediaStatus) error {
		r.Caption = caption
		return nil
	}
}

func OptionalParticipantsMediaStatus(participants []string) SendMediaStatusOption {
	return func(r *RequestSendMediaStatus) error {
		r.Participants = participants
		return nil
	}
}

func (c StatusesCategory) SendMediaStatus(urlFile, fileName string, options ...SendMediaStatusOption) (*APIResponse, error) {
	r := &RequestSendMediaStatus{
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

	return c.GreenAPI.Request("POST", "sendMediaStatus", jsonData)
}

// ------------------------------------------------------------------ DeleteStatus block

type RequestDeleteStatus struct {
	IdMessage string `json:"idMessage"`
}

func (c StatusesCategory) DeleteStatus(idMessage string) (*APIResponse, error) {
	r := &RequestDeleteStatus{
		IdMessage: idMessage,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPI.Request("POST", "deleteStatus", jsonData)
}

// ------------------------------------------------------------------ GetStatusStatistic block

func (c StatusesCategory) GetStatusStatistic(idMessage string) (*APIResponse, error) {
	addUrl := fmt.Sprintf("?idMessage=%s", idMessage)

	return c.GreenAPI.Request("GET", "getStatusStatistic", nil, WithGetParams(addUrl))
}

// ------------------------------------------------------------------ GetOutgoingStatuses + GetIncomingStatuses block

type RequestGetLastStatuses struct {
	Minutes int `json:"minutes,omitempty"`
}

type GetLastStatusesOption func(*RequestGetLastStatuses) error

// Description
func OptionalMinutesOfStatuses(minutes int) GetLastStatusesOption {
	return func(r *RequestGetLastStatuses) error {
		r.Minutes = minutes
		return nil
	}
}

func (c StatusesCategory) GetOutgoingStatuses(options ...GetLastStatusesOption) (*APIResponse, error) {
	r := &RequestGetLastStatuses{}

	for _, o := range options {
		o(r)
	}

	var addUrl string
	if r.Minutes != 0 {
		addUrl = fmt.Sprintf("?minutes=%v", r.Minutes)
	}

	return c.GreenAPI.Request("GET", "getOutgoingStatuses", nil, WithGetParams(addUrl))
}

func (c StatusesCategory) GetIncomingStatuses(options ...GetLastStatusesOption) (*APIResponse, error) {
	r := &RequestGetLastStatuses{}

	for _, o := range options {
		o(r)
	}

	var addUrl string
	if r.Minutes != 0 {
		addUrl = fmt.Sprintf("?minutes=%v", r.Minutes)
	}

	return c.GreenAPI.Request("GET", "getIncomingStatuses", nil, WithGetParams(addUrl))
}
