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

	buttons := []greenapi.InteractiveButton{
		{
			Type:       "copy",
			ButtonId:   "1",
			ButtonText: "Copy code",
			CopyCode:   "123456",
		},
		{
			Type:        "call",
			ButtonId:    "2",
			ButtonText:  "Call support",
			PhoneNumber: "79123456789",
		},
		{
			Type:       "url",
			ButtonId:   "3",
			ButtonText: "Visit website",
			URL:        "https://green-api.com",
		},
	}

	response, err := GreenAPI.Sending().SendInteractiveButtons(
		"11001234567@c.us",
		"Please choose an action:",
		buttons,
		greenapi.OptionalInteractiveHeader("Support Options"),
		greenapi.OptionalInteractiveFooter("Select one option below"),
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
