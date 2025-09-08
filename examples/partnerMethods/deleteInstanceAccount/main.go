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

	response, err := Partner.Partner().DeleteInstanceAccount(1101000000)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Status: %v %s \n\rResponse: %s\n\rTimestamp: %s\n\r", response.StatusCode,
		response.StatusMessage,
		response.Body,
		response.Timestamp.Format("15:04:05.000"))
}