package main

import (
	"fmt"
	"log"

	greenapi "github.com/green-api/whatsapp-api-client-golang-v2"
)

func main() {
	GreenAPI := greenapi.GreenAPI{
		APIURL:           "https://api.green-api.com",
		MediaURL:         "https://media.green-api.com",
		IDInstance:       "1101000001",
		APITokenInstance: "d75b3a66374942c5b3c019c698abc2067e151558acbd412345",
	}

	replyButtons := []greenapi.InteractiveReplyButton{
		{
			ButtonId:   "1",
			ButtonText: "Option 1",
		},
		{
			ButtonId:   "2",
			ButtonText: "Option 2",
		},
		{
			ButtonId:   "3",
			ButtonText: "Option 3",
		},
	}

	response, err := GreenAPI.Sending().SendInteractiveButtonsReply(
		"11001234567@c.us",
		"How can we help you?",
		replyButtons,
		greenapi.OptionalInteractiveReplyHeader("Customer Support"),
		greenapi.OptionalInteractiveReplyFooter("Choose one option"),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Status: %v %s \nResponse: %s\nTimestamp: %s\n\n",
		response.StatusCode,
		response.StatusMessage,
		response.Body,
		response.Timestamp.Format("15:04:05.000"))
}
