# Sending — `client.Sending()`

All methods return `(*greenapi.APIResponse, error)`. Successful send responses have Body `{"idMessage":"..."}`.
`chatId` must end with `@c.us` or `@g.us` (SDK-validated). `typingTime` accepts 1000–20000 ms.

## SendMessage — send text

Docs: https://green-api.com/en/docs/api/sending/SendMessage/

```go
SendMessage(chatId, message string, options ...SendMessageOption)
```

`message` max 20 000 chars (SDK-validated). Options:

- `OptionalQuotedMessageId(quotedMessageId string)` — reply to a message in the same chat
- `OptionalLinkPreview(linkPreview bool)` — link preview on/off (API default true)
- `OptionalTypePreview(typePreview string)` — `"large"` or `"small"`
- `OptionalCustomPreview(customPreview CustomPreview)` — custom link preview; struct fields: `Title`, `Description`, `Link`, `UrlFile`, `JpegThumbnail` (all string; preview image ≤ 2 MB)
- `OptionalMessageTypingTime(typingTime int)` — show "typing..." before sending

## SendPoll — send a poll

Docs: https://green-api.com/en/docs/api/sending/SendPoll/

```go
SendPoll(chatId, message string, pollOptions []string, options ...SendPollOption)
```

`message` (poll question) max 255 chars; `pollOptions` must contain 2–12 unique strings (SDK-validated). Options:

- `OptionalMultipleAnswers(multipleAnswers bool)` — allow multiple answers (default false)
- `OptionalPollQuotedMessageId(quotedMessageId string)`
- `OptionalPollTypingTime(typingTime int)`

## SendFileByUpload — send a local file

Docs: https://green-api.com/en/docs/api/sending/SendFileByUpload/
Goes through the media host (`MediaURL`), multipart form-data. Max 100 MB.

```go
SendFileByUpload(chatId, filePath, fileName string, options ...SendFileByUploadOption)
```

- `filePath` — path on local disk; `fileName` — name with extension (drives file-type detection).
- `OptionalCaptionSendUpload(caption string)` — max 20 000 chars
- `OptionalQuotedMessageIdSendUpload(quotedMessageId string)`
- `OptionalUploadTypingTime(typingTime int)`
- `OptionalUploadTypingType(typingType string)` — `"typing"` or `"recording"`

## SendFileByUrl — send a file by URL

Docs: https://green-api.com/en/docs/api/sending/SendFileByUrl/

```go
SendFileByUrl(chatId, urlFile, fileName string, options ...SendFileByUrlOption)
```

`urlFile` must be a direct, publicly downloadable link. Max 100 MB. Options:

- `OptionalCaptionSendUrl(caption string)` — max 20 000 chars
- `OptionalQuotedMessageIdSendUrl(quotedMessageId string)`
- `OptionalUrlTypingTime(typingTime int)`
- `OptionalUrlTypingType(typingType string)` — `"typing"` or `"recording"`

## UploadFile — upload file to GREEN-API storage

Docs: https://green-api.com/en/docs/api/sending/UploadFile/
Uploads to the media host and returns Body `{"urlFile":"..."}` — use that URL with `SendFileByUrl`.

```go
UploadFile(filePath string)
```

## SendLocation

Docs: https://green-api.com/en/docs/api/sending/SendLocation/

```go
SendLocation(chatId string, latitude, longitude float32, options ...SendLocationOption)
```

- `OptionalNameLocation(nameLocation string)` — location name
- `OptionalAddress(address string)` — location address
- `OptionalQuotedMessageIdLocation(quotedMessageId string)`
- `OptionalLocationTypingTime(typingTime int)`

## SendContact

Docs: https://green-api.com/en/docs/api/sending/SendContact/

```go
SendContact(chatId string, contact Contact, options ...SendContactOption)
```

```go
contact := greenapi.Contact{
	PhoneContact: 79876543210,      // int, country code + number, digits only (required)
	FirstName:    "John",           // at least one name field required
	MiddleName:   "",
	LastName:     "Doe",
	Company:      "",
}
```

- `OptionalQuotedMessageIdContact(quotedMessageId string)`
- `OptionalContactTypingTime(typingTime int)`

## ForwardMessages

Docs: https://green-api.com/en/docs/api/sending/ForwardMessages/

```go
ForwardMessages(chatId, chatIdFrom string, messages []string, options ...ForwardMessagesOption)
```

`chatId` — destination; `chatIdFrom` — chat the messages are forwarded from; `messages` — slice of message ids. Response Body: `{"messages":["idMessage1", ...]}`.

- `OptionalForwardTypingTime(typingTime int)`

## SendInteractiveButtons

Docs: https://green-api.com/en/docs/api/sending/SendInteractiveButtons/

```go
SendInteractiveButtons(chatId, body string, buttons []InteractiveButton, options ...SendInteractiveButtonsOption)
```

`body` max 255 chars (SDK-validated).

```go
buttons := []greenapi.InteractiveButton{
	{Type: "reply", ButtonId: "1", ButtonText: "Yes"},
	{Type: "copy",  ButtonId: "2", ButtonText: "Copy code", CopyCode: "ABC123"},
	{Type: "call",  ButtonId: "3", ButtonText: "Call us",   PhoneNumber: "79876543210"},
	{Type: "url",   ButtonId: "4", ButtonText: "Open site", URL: "https://example.com"},
}
```

- `OptionalInteractiveHeader(header string)`
- `OptionalInteractiveFooter(footer string)`
- `OptionalInteractiveQuotedMessageId(quotedMessageId string)`
- `OptionalInteractiveTypingTime(typingTime int)`

## SendInteractiveButtonsReply

Docs: https://green-api.com/en/docs/api/sending/SendInteractiveButtonsReply/
Reply-only buttons (no type field). `body` max 255 chars.

```go
SendInteractiveButtonsReply(chatId, body string, buttons []InteractiveReplyButton, options ...SendInteractiveButtonsReplyOption)
```

```go
buttons := []greenapi.InteractiveReplyButton{
	{ButtonId: "1", ButtonText: "Yes"},
	{ButtonId: "2", ButtonText: "No"},
}
```

- `OptionalInteractiveReplyHeader(header string)`
- `OptionalInteractiveReplyFooter(footer string)`
- `OptionalInteractiveReplyQuotedMessageId(quotedMessageId string)`
- `OptionalInteractiveReplyTypingTime(typingTime int)`
