package main

import (
	"log"

	"github.com/fatih/color"
	greenapi "github.com/green-api/whatsapp-api-client-golang"
)

func main() {
	GreenAPI := greenapi.GreenAPI{
		APIURL:           "https://api.green-api.com",
		MediaURL:         "https://media.green-api.com",
		IDInstance:       "1101000001",
		APITokenInstance: "d75b3a66374942c5b3c019c698abc2067e151558acbd412345",
	}

	response, err := GreenAPI.Account().SetSettings(greenapi.OptionDeviceWebhook(false),
		greenapi.OptionDelaySendMesssages(999),
		greenapi.OptionOutgoingWebhook(true),
		greenapi.OptionIncomingWebhook(true),
	)
	if err != nil {
		log.Fatal(err)
	}

	color.Green("Status: %v %s \n\rResponse: %s\n\rTimestamp: %s", response.StatusCode,
		response.StatusMessage,
		response.Body,
		response.Timestamp.Format("15:04:05.000"))
}
