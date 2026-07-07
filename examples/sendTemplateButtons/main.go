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

	templateButtons := []greenapi.TemplateButton{
		{
			Index: 1,
			UrlButton: &greenapi.UrlButton{
				DisplayText: "Visit website",
				Url:         "https://green-api.com",
			},
		},
		{
			Index: 2,
			CallButton: &greenapi.CallButton{
				DisplayText: "Call us",
				PhoneNumber: "11001234567",
			},
		},
		{
			Index: 3,
			QuickReplyButton: &greenapi.QuickReplyButton{
				DisplayText: "Quick reply",
				Id:          "reply1",
			},
		},
	}

	response, err := GreenAPI.Sending().SendTemplateButtons(
		"11001234567@c.us",
		"Choose an action:",
		templateButtons,
		greenapi.OptionalTemplateFooter("Powered by Green API"),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Status: %v %s \n\rResponse: %s\n\rTimestamp: %s\n\r", response.StatusCode,
		response.StatusMessage,
		response.Body,
		response.Timestamp.Format("15:04:05.000"))
}
