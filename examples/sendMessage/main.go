package main

import (
	"fmt"
	"log"

	greenapi "github.com/green-api/whatsapp-api-client-golang"
)

func main() {
	GreenAPI := greenapi.GreenAPI{
		Host:             "https://api.green-api.com",
		IDInstance:       "9903222222",
		APITokenInstance: "d75b3a66374942c5b3c019c698abc2067e151558acbd412345",
		//PartnerToken:     "",
	}

	response, err := GreenAPI.Account().SetSettings(map[string]interface{}{
		"delaySendMessagesMilliseconds": 1000,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response)
}
