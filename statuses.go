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

// Status background. Default: #FFFFFF.
func OptionalBackgroundColorText(backgroundColor string) SendTextStatusOption {
	return func(r *RequestSendTextStatus) error {
		r.BackgroundColor = backgroundColor
		return nil
	}
}

// Text font. Accepts values: SERIF, SANS_SERIF, NORICAN_REGULAR, BRYNDAN_WRITE, OSWALD_HEAVY
func OptionalFont(font string) SendTextStatusOption {
	return func(r *RequestSendTextStatus) error {
		r.Font = font
		return nil
	}
}

// An array of strings with contact IDs for whom the status will be available. If the field value is empty, the status will be available to all contacts.
func OptionalParticipantsTextStatus(participants []string) SendTextStatusOption {
	return func(r *RequestSendTextStatus) error {
		for _, participant := range participants {
			err := ValidateChatId(participant)
			if err!=nil {
				return err
			}
		}

		r.Participants = participants
		return nil
	}
}

// Sending a text status.
//
// https://green-api.com/docs/api/statuses/SendTextStatus/
//
// Add optional arguments by passing these functions:
//
//  OptionalBackgroundColorText(backgroundColor string) <- Status background. Default: #FFFFFF.
//  OptionalFont(font string) <- Text font. Accepts values: SERIF, SANS_SERIF, NORICAN_REGULAR, BRYNDAN_WRITE, OSWALD_HEAVY
//  OptionalParticipantsTextStatus(participants []string) <- An array of strings with contact IDs for whom the status will be available. If the field value is empty, the status will be available to all contacts.
func (c StatusesCategory) SendTextStatus(message string, options ...SendTextStatusOption) (*APIResponse, error) {
	err := ValidateMessageLength(message, 500)
	if err != nil {
		return nil, err
	}

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

// Status background. Default: #FFFFFF.
func OptionalBackgroundColorVoice(backgroundColor string) SendVoiceStatusOption {
	return func(r *RequestSendVoiceStatus) error {
		r.BackgroundColor = backgroundColor
		return nil
	}
}

// An array of strings with contact IDs for whom the status will be available. If the field value is empty, the status will be available to all contacts.
func OptionalParticipantsVoiceStatus(participants []string) SendVoiceStatusOption {
	return func(r *RequestSendVoiceStatus) error {
		for _, participant := range participants {
			err := ValidateChatId(participant)
			if err!=nil {
				return err
			}
		}

		r.Participants = participants
		return nil
	}
}

// Sending a voice status.
// 
// https://green-api.com/en/docs/api/statuses/SendVoiceStatus/
// 
// Add optional arguments by passing these functions:
//  OptionalBackgroundColorVoice(backgroundColor string) <- Status background. Default: #FFFFFF.
//  OptionalParticipantsVoiceStatus(participants []string) <- An array of strings with contact IDs for whom the status will be available. If the field value is empty, the status will be available to all contacts.
func (c StatusesCategory) SendVoiceStatus(urlFile, fileName string, options ...SendVoiceStatusOption) (*APIResponse, error) {
	err := ValidateURL(urlFile)
	if err!=nil {
		return nil, err
	}
	
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

// Media status caption.
func OptionalCaptionMediaStatus(caption string) SendMediaStatusOption {
	return func(r *RequestSendMediaStatus) error {
		err := ValidateMessageLength(r.Caption, 1024)
		if err!=nil {
			return err
		}

		r.Caption = caption
		return nil
	}
}

// An array of strings with contact IDs for whom the status will be available. If the field value is empty, the status will be available to all contacts.
func OptionalParticipantsMediaStatus(participants []string) SendMediaStatusOption {
	return func(r *RequestSendMediaStatus) error {
		for _, participant := range participants {
			err := ValidateChatId(participant)
			if err!=nil {
				return err
			}
		}

		r.Participants = participants
		return nil
	}
}

// Sending a media status.
//
// https://green-api.com/en/docs/api/statuses/SendMediaStatus/
//
// Add optional arguments by passing these functions:
//  OptionalCaptionMediaStatus(caption string) <- Media status caption.
//  OptionalParticipantsMediaStatus(participants []string) <- An array of strings with contact IDs for whom the status will be available. If the field value is empty, the status will be available to all contacts.
func (c StatusesCategory) SendMediaStatus(urlFile, fileName string, options ...SendMediaStatusOption) (*APIResponse, error) {
	err := ValidateURL(urlFile)
	if err!=nil {
		return nil, err
	}

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

// Deleting a posted status. 
// 
// https://green-api.com/en/docs/api/statuses/DeleteStatus/
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

// Getting an array of recipients marked sent/delivered/read for a given status.
// 
// https://green-api.com/en/docs/api/statuses/GetStatusStatistic/
func (c StatusesCategory) GetStatusStatistic(idMessage string) (*APIResponse, error) {
	addUrl := fmt.Sprintf("?idMessage=%s", idMessage)

	return c.GreenAPI.Request("GET", "getStatusStatistic", nil, WithGetParams(addUrl))
}

// ------------------------------------------------------------------ GetOutgoingStatuses + GetIncomingStatuses block

type RequestGetLastStatuses struct {
	Minutes int `json:"minutes,omitempty"`
}

type GetLastStatusesOption func(*RequestGetLastStatuses) error

// Time in minutes for which the status messages should be displayed (1440 minutes by default)
func OptionalMinutesOfStatuses(minutes int) GetLastStatusesOption {
	return func(r *RequestGetLastStatuses) error {
		r.Minutes = minutes
		return nil
	}
}

// Getting the outgoing statuses of an account.
// 
// https://green-api.com/en/docs/api/statuses/GetOutgoingStatuses/
// 
// Add optional arguments by passing these functions:
//  OptionalMinutesOfStatuses(minutes int) <- Time in minutes for which the status messages should be displayed (1440 minutes by default)
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

// Getting the incoming statuses of an account.
// 
// https://green-api.com/en/docs/api/statuses/GetIncomingStatuses/
//
// Add optional arguments by passing these functions:
//  OptionalMinutesOfStatuses(minutes int) <- Time in minutes for which the status messages should be displayed (1440 minutes by default)
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
