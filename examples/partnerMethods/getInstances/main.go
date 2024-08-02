package main

import (
	"fmt"
	"log"

	greenapi "github.com/green-api/whatsapp-api-client-golang-v2"
)

func main() {
	Partner := greenapi.GreenAPIPartner{
		PartnerToken: "gac.1234567891234567891234567891213456789",
		Email: "mail@email.com",
	}

	response, err := Partner.Partner().CreateInstance(
		greenapi.OptionalWebhookUrl("webhook_url"),
		greenapi.OptionalWebhookUrlToken("auth_token"),
		greenapi.OptionalDelaySendMesssages(5000),
		greenapi.OptionalMarkIncomingMessagesRead(true),
		greenapi.OptionalMarkIncomingMessagesReadOnReply(true),
		greenapi.OptionalOutgoingWebhook(true),
		greenapi.OptionalOutgoingMessageWebhook(true),
		greenapi.OptionalOutgoingAPIMessageWebhook(true),
		greenapi.OptionalStateWebhook(true),
		greenapi.OptionalIncomingWebhook(true),
		greenapi.OptionalDeviceWebhook(true),
		greenapi.OptionalKeepOnlineStatus(true),
		greenapi.OptionalPollMessageWebhook(true),
		greenapi.OptionalIncomingBlockWebhook(true),
		greenapi.OptionalIncomingCallWebhook(true),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Status: %v %s \n\rResponse: %s\n\rTimestamp: %s\n\r", response.StatusCode,
		response.StatusMessage,
		response.Body,
		response.Timestamp.Format("15:04:05.000"))
}