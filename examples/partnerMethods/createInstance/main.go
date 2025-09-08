package main

import (
	"fmt"
	"log"

	greenapi "github.com/green-api/max-api-client-golang"
)

func main() {
	Partner := greenapi.GreenAPIPartner{
		PartnerToken: "gac.1234567891234567891234567891213456789",
		Email:        "mail@email.com",
	}

	response, err := Partner.Partner().CreateInstance(
		greenapi.OptionalName("Created by GO SDK"),
		greenapi.OptionalWebhookUrl("https://webhook.url"),
		greenapi.OptionalWebhookUrlToken("auth_token"),
		greenapi.OptionalDelaySendMessages(5000),
		greenapi.OptionalMarkIncomingMessagesRead(true),
		greenapi.OptionalMarkIncomingMessagesReadOnReply(true),
		greenapi.OptionalOutgoingWebhook(true),
		greenapi.OptionalOutgoingMessageWebhook(true),
		greenapi.OptionalOutgoingAPIMessageWebhook(true),
		greenapi.OptionalStateWebhook(true),
		greenapi.OptionalIncomingWebhook(true),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Status: %v %s \n\rResponse: %s\n\rTimestamp: %s\n\r", response.StatusCode,
		response.StatusMessage,
		response.Body,
		response.Timestamp.Format("15:04:05.000"))
}
