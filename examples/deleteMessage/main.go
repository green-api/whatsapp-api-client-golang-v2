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

	/* Delete message only for everyone (default) */
	response, err := GreenAPI.Service().DeleteMessage(
		"11001234567@c.us",
		"BAE56DEC325DB4AB",
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Status: %v %s \n\rResponse: %s\n\rTimestamp: %s\n\r", response.StatusCode,
		response.StatusMessage,
		response.Body,
		response.Timestamp.Format("15:04:05.000"))

	/* Delete message only for for sender */
	GreenAPI.Service().DeleteMessage(
		"11001234567@c.us",
		"BAE515D2ACF72E34",
		greenapi.OptionalOnlySenderDelete(true),
	)

	/* Delete message for everyone */
	GreenAPI.Service().DeleteMessage(
		"11001234567@c.us",
		"BAE563E61BA23E8A",
		greenapi.OptionalOnlySenderDelete(false),
	)
}
