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

	sections := []greenapi.ListSection{
		{
			Title: "Section 1",
			Rows: []greenapi.ListRow{
				{Title: "Row 1", RowId: "row1", Description: "First option"},
				{Title: "Row 2", RowId: "row2", Description: "Second option"},
			},
		},
		{
			Title: "Section 2",
			Rows: []greenapi.ListRow{
				{Title: "Row 3", RowId: "row3", Description: "Third option"},
			},
		},
	}

	response, err := GreenAPI.Sending().SendListMessage(
		"11001234567@c.us",
		"Please select an option from the list",
		"View options",
		sections,
		greenapi.OptionalListTitle("Our Menu"),
		greenapi.OptionalListFooter("Powered by Green API"),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Status: %v %s \n\rResponse: %s\n\rTimestamp: %s\n\r", response.StatusCode,
		response.StatusMessage,
		response.Body,
		response.Timestamp.Format("15:04:05.000"))
}
