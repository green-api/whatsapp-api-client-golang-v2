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

	response, err := GreenAPI.Groups().CreateGroup(
		"Group Title",
		[]string{
			"11001211111@c.us",
			"11001222222@c.us",
			"11001233333@c.us",
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Status: %v %s \n\rResponse: %s\n\rTimestamp: %s\n\r", response.StatusCode,
		response.StatusMessage,
		response.Body,
		response.Timestamp.Format("15:04:05.000"))
}
