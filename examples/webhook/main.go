package main

import (
	"fmt"

	"github.com/green-api/whatsapp-api-client-golang/pkg/api"
	"github.com/green-api/whatsapp-api-client-golang/pkg/webhook"
)

func main() {
	GreenAPI := api.GreenAPI{
		IDInstance:       "1101000001",
		APITokenInstance: "d75b3a66374942c5b3c019c698abc2067e151558acbd412345",
	}

	GreenAPIWebhook := webhook.GreenAPIWebhook{
		GreenAPI: GreenAPI,
	}

	GreenAPIWebhook.Start(func(body map[string]interface{}) {
		fmt.Println(body)
	})
}
