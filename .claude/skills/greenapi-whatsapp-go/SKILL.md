---
name: greenapi-whatsapp-go
version: 1.0.0
description: Write correct Go code with the GREEN-API WhatsApp SDK (github.com/green-api/whatsapp-api-client-golang-v2). Use when sending/receiving WhatsApp messages, managing instances, groups, contacts, or statuses via GREEN-API in Go.
---

# GREEN-API WhatsApp SDK for Go (v2)

Wrapper over the GREEN-API REST API (official docs: https://green-api.com/en/docs/api/).
Module: `github.com/green-api/whatsapp-api-client-golang-v2`, package name `greenapi`. Requires Go >= 1.23.

## When to apply

Use this skill whenever writing Go code that talks to WhatsApp through GREEN-API: sending text/files/polls/locations/buttons, receiving notifications (polling or webhook settings), account/instance management, groups, contacts, statuses, message journals, queues, partner methods.

## Setup

```shell
go mod init myapp
go get github.com/green-api/whatsapp-api-client-golang-v2
```

## Client initialization

There is **no constructor function**. The client is a plain struct literal:

```go
import (
	greenapi "github.com/green-api/whatsapp-api-client-golang-v2"
)

client := greenapi.GreenAPI{
	APIURL:           "https://api.green-api.com",   // use API URL shown for your instance in console.green-api.com
	MediaURL:         "https://media.green-api.com", // media host (file upload/send-by-upload)
	IDInstance:       "1101000001",
	APITokenInstance: "d75b3a66374942c5b3c019c698abc2067e151558acbd412345",
}
```

Credentials come from the GREEN-API console (https://console.green-api.com/). The WhatsApp account must be **authorized** (QR scan in console, or `client.Account().QR()` / `GetAuthorizationCode`). Check before sending:

```go
resp, err := client.Account().GetStateInstance() // Body: {"stateInstance":"authorized"}
```

If the instance is not authorized, outgoing messages sit in the queue for up to 24 hours.

## Universal call pattern

Every API method hangs off a category accessor and returns `(*greenapi.APIResponse, error)`:

```go
client.Account()   client.Sending()   client.Receiving()  client.Groups()
client.Service()   client.Journals()  client.Queues()     client.Statuses()
client.Contacts()  client.ReadMark()
```

`APIResponse` fields: `StatusCode int`, `StatusMessage []byte`, `Body json.RawMessage`, `Timestamp time.Time`.

- `err != nil` — transport or client-side validation error.
- API-level errors arrive as `resp.StatusCode != 200` — **always check both**.
- `resp.Body` is raw JSON; unmarshal into your own struct (the SDK has no typed responses).

Optional parameters are passed as variadic `greenapi.Optional*` functions (functional options). Each method has its own set — names are method-specific (e.g. `OptionalCaptionSendUrl` vs `OptionalCaptionSendUpload`); take them from the reference files, do not guess.

## Sending a message

```go
resp, err := client.Sending().SendMessage(
	"79876543210@c.us",
	"Hello",
	greenapi.OptionalLinkPreview(false),       // optional
	greenapi.OptionalQuotedMessageId("BAE5..."), // optional
)
if err != nil {
	log.Fatal(err)
}
fmt.Printf("Status: %d, Body: %s\n", resp.StatusCode, resp.Body)
// success Body: {"idMessage":"3EB0C767D097B7C7C030"}
```

## Sending a file

```go
// By URL (file must be publicly downloadable):
resp, err := client.Sending().SendFileByUrl(
	"79876543210@c.us",
	"https://example.com/image.jpg",
	"image.jpg", // fileName with extension — determines file type
	greenapi.OptionalCaptionSendUrl("Look at this"),
)

// From local disk (goes through MediaURL host):
resp, err := client.Sending().SendFileByUpload(
	"79876543210@c.us",
	"/path/to/image.jpg", // filePath
	"image.jpg",          // fileName
	greenapi.OptionalCaptionSendUpload("Look at this"),
)
```

Max file size 100 MB; caption max 20 000 chars (per SDK validation; docs note 1024 for media captions).

## Receiving notifications (polling loop)

Official flow: `ReceiveNotification` → process → `DeleteNotification(receiptId)` → repeat. A notification is removed from the queue **only** after `DeleteNotification`; otherwise you receive it again.

```go
type Notification struct {
	ReceiptId int `json:"receiptId"`
	Body      struct {
		TypeWebhook string `json:"typeWebhook"`
		SenderData  struct {
			ChatId     string `json:"chatId"`
			Sender     string `json:"sender"`
			SenderName string `json:"senderName"`
		} `json:"senderData"`
		MessageData struct {
			TypeMessage     string `json:"typeMessage"`
			TextMessageData struct {
				TextMessage string `json:"textMessage"`
			} `json:"textMessageData"`
		} `json:"messageData"`
	} `json:"body"`
}

for {
	resp, err := client.Receiving().ReceiveNotification(
		greenapi.OptionalReceiveTimeout(20), // 5–60 s, default 5
	)
	if err != nil {
		log.Printf("receive error: %v", err)
		continue
	}
	if string(resp.Body) == "null" || len(resp.Body) == 0 {
		continue // queue empty, timeout reached
	}
	var n Notification
	if err := json.Unmarshal(resp.Body, &n); err != nil {
		log.Printf("parse error: %v", err)
		continue
	}
	// ... handle n.Body (e.g. typeWebhook == "incomingMessageReceived") ...
	client.Receiving().DeleteNotification(n.ReceiptId)
}
```

To receive via **webhooks** instead, set your server URL through settings (notifications are then POSTed to it; same JSON as the `body` above):

```go
client.Account().SetSettings(
	greenapi.OptionalWebhookUrl("https://myserver.example.com/webhook"),
	greenapi.OptionalIncomingWebhook(true),
)
```

Note: incoming-notification method names/types are listed in [references/receiving.md](references/receiving.md).

## Gotchas

- **chatId format**: personal chat `79876543210@c.us` (country code, digits only, no `+`), group `120363043968066561@g.us`. The SDK's `ValidateChatId` **rejects anything not ending in `@c.us`/`@g.us`** (so `@lid` ids cannot be passed to validating methods). Never build group ids by hand — take them from `CreateGroup`/journals/notifications.
- **CheckWhatsapp quirk**: signature is `CheckWhatsapp(phoneNumber int, options...)`; `phoneNumber` is deprecated — pass `0` and use `greenapi.OptionalChatID("79876543210@c.us")`. Passing both zero/nil returns an error.
- **Sending delays**: messages go through a FIFO queue; the interval is `delaySendMessagesMilliseconds` (min 500 ms, set via `SetSettings`, `greenapi.OptionalDelaySendMessages(5000)`). For bulk sends rely on this queue rather than client-side sleeps, and don't blast new numbers — risk of ban.
- **SetSettings applies asynchronously**: the instance restarts and settings take effect within ~5 minutes.
- Text message limit 20 000 chars; poll question / interactive-buttons body 255; text status 500 (SDK validates and returns an error before the request).
- Don't invent methods: only methods listed in the reference files exist in this SDK. Notably absent: `SendButtons`, `SendListMessage`, `SendTemplateButtons`, `GetDeviceInfo`, `ScanQRCode` websocket.

## Method reference (all implemented methods)

| Category | File |
|---|---|
| Sending: SendMessage, SendPoll, SendFileByUpload, SendFileByUrl, UploadFile, SendLocation, SendContact, ForwardMessages, SendInteractiveButtons, SendInteractiveButtonsReply | [references/sending.md](references/sending.md) |
| Receiving & Queues: ReceiveNotification, DeleteNotification, DownloadFile, ShowMessagesQueue, ClearMessagesQueue, GetWebhooksCount, ClearWebhooksQueue + webhook JSON format | [references/receiving.md](references/receiving.md) |
| Account & Partner: GetSettings, SetSettings, GetStateInstance, GetStatusInstance, Reboot, Logout, QR, GetAuthorizationCode, SetProfilePicture, GetStateInstanceHistory, UpdateApiToken, GetWaSettings; Partner: GetInstances, CreateInstance, DeleteInstanceAccount | [references/account.md](references/account.md) |
| Service, Journals, Contacts, ReadMark: CheckWhatsapp, GetAvatar, GetContacts, GetContactInfo, DeleteMessage, EditMessage, ArchiveChat, UnarchiveChat, SetDisappearingChat, SendTyping, GetChats, GetChatHistory, GetMessage, LastIncoming/OutgoingMessages, LastIncoming/OutgoingCalls, AddContact, EditContact, DeleteContact, ReadChat | [references/service-chats.md](references/service-chats.md) |
| Groups & Statuses: CreateGroup, UpdateGroupName, GetGroupData, UpdateGroupSettings, Add/RemoveGroupParticipant, SetGroupAdmin, RemoveAdmin, SetGroupPicture, LeaveGroup; SendText/Voice/MediaStatus, DeleteStatus, GetStatusStatistic, GetOutgoing/IncomingStatuses | [references/groups-statuses.md](references/groups-statuses.md) |
