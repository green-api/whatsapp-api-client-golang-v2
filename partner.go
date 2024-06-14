package greenapi

import "encoding/json"

type PartnerCategory struct {
	GreenAPIPartner GreenAPIPartnerInterface
}

// ------------------------------------------------------------------ GetInstances block

// Getting all the account instances created by the partner.
// 
// https://green-api.com/en/docs/partners/getInstances/
func (c PartnerCategory) GetInstances() (*APIResponse, error) {
	return c.GreenAPIPartner.PartnerRequest("GET", "getInstances", nil)
}

// ------------------------------------------------------------------ CreateInstance block

// Creating an instance. 
// 
// https://green-api.com/en/docs/partners/createInstance/
//
// Add optional arguments by passing these functions:
//  OptionalWebhookUrl(webhookUrl string) <- URL for sending notifications.
//  OptionalWebhookUrlToken(webhookUrlToken string) <- Token to access your notification server.
//  OptionalDelaySendMesssages(delaySendMessagesMilliseconds int) <- Message sending delay. 
//  OptionalMarkIncomingMessagesRead(markIncomingMessagesReaded bool) <- Mark incoming messages as read or not.
//  OptionalMarkIncomingMessagesReadOnReply(markIncomingMessagesReadedOnReply bool) <- Mark incoming messages as read when posting a message to the chat via API.
//  OptionalOutgoingWebhook(outgoingWebhook bool) <- Get notifications about outgoing messages sending/delivering/reading statuses.
//  OptionalOutgoingMessageWebhook(outgoingMessageWebhook bool) <- Get notifications about messages sent from the phone.
//  OptionalOutgoingAPIMessageWebhook(outgoingAPIMessageWebhook bool) <- Get notifications about messages sent from API.
//  OptionalStateWebhook(stateWebhook bool) <- Get notifications about the instance authorization state change.
//  OptionalIncomingWebhook(incomingWebhook bool) <- Get notifications about incoming messages and files.
//  OptionalDeviceWebhook(deviceWebhook bool) <- Get notifications about the device (phone) and battery level.
//  OptionalKeepOnlineStatus(keepOnlineStatus bool) <- Sets the 'Online' status for your Whatsapp account.
//  OptionalPollMessageWebhook(pollMessageWebhook bool) <- Get notifications about the creation of a poll and voting in the poll.
//  OptionalIncomingBlockWebhook(incomingBlockWebhook bool) <- Get notifications about adding a chat to the list of blocked contacts.
//  OptionalIncomingCallWebhook(incomingCallWebhook bool) <- Get notifications about incoming call statuses.
func (c PartnerCategory) CreateInstance(options ...SetSettingsOption) (*APIResponse, error) {
	r := &RequestSetSettings{}

	for _, o := range options {
		o(r)
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPIPartner.PartnerRequest("POST", "createInstance", jsonData)
}

// ------------------------------------------------------------------ DeleteInstanceAccount block

type RequestDeleteInstanceAccount struct {
	IdInstance int `json:"idInstance"`
}

// Deleting an instance.
// 
// https://green-api.com/en/docs/partners/deleteInstanceAccount/
func (c PartnerCategory) DeleteInstanceAccount(idInstance int) (*APIResponse, error) {
	r := &RequestDeleteInstanceAccount{
		IdInstance: idInstance,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return c.GreenAPIPartner.PartnerRequest("POST", "deleteInstanceAccount", jsonData)
}
