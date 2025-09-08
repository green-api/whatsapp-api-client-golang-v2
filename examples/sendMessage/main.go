package main

import (
	"fmt"
	"log"

	greenapi "github.com/green-api/max-api-client-golang"
)

func main() {
	GreenAPI := greenapi.GreenAPI{
		APIURL:           "https://api.green-api.com/v3",
		MediaURL:         "https://api.green-api.com/v3",
		IDInstance:       "3100000001",
		APITokenInstance: "d75b3a66374942c5b3c019c698abc2067e151558acbd412345",
	}

	response, err := GreenAPI.Sending().SendMessage(
		"10000000",
		"Hello",
		greenapi.OptionalLinkPreview(false),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Status: %v %s \n\rResponse: %s\n\rTimestamp: %s\n\r", response.StatusCode,
		response.StatusMessage,
		response.Body,
		response.Timestamp.Format("15:04:05.000"))
}
