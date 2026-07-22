# Receiving & Queues

## Receiving — `client.Receiving()`

### ReceiveNotification — poll one notification

Docs: https://green-api.com/en/docs/api/receiving/technology-http-api/ReceiveNotification/

```go
ReceiveNotification(options ...ReceiveNotificationOption)
```

- `OptionalReceiveTimeout(seconds int)` — how long to wait for a notification, 5–60 s, default 5.

Response Body: `{"receiptId": 123, "body": {...notification...}}`, or `null` when the queue is empty and the timeout elapsed — check `string(resp.Body) == "null"` before unmarshalling.

Notifications are FIFO, stored up to 24 hours. A notification stays in the queue (and is re-delivered) until you call `DeleteNotification` with its `receiptId`.

### DeleteNotification — confirm processing

Docs: https://green-api.com/en/docs/api/receiving/technology-http-api/DeleteNotification/

```go
DeleteNotification(receiptId int)
```

Response Body: `{"result": true}`.

### DownloadFile — download a file from a message

Docs: https://green-api.com/en/docs/api/receiving/files/DownloadFile/

```go
DownloadFile(chatId, idMessage string)
```

Response Body contains `{"downloadUrl": "..."}` — a link valid for downloading the file of an incoming/outgoing file message.

## Notification (webhook) JSON format

Docs: https://green-api.com/en/docs/api/receiving/notifications-format/

Same JSON arrives in the `body` field of `ReceiveNotification` and as the POST body on your webhook URL. Common fields:

```json
{
  "typeWebhook": "incomingMessageReceived",
  "instanceData": {"idInstance": 1101000001, "wid": "79876543210@c.us", "typeInstance": "whatsapp"},
  "timestamp": 1588091580,
  "idMessage": "F7AEC1B7086ECDC7E6E45923F5EDB825",
  "senderData": {"chatId": "79876543210@c.us", "sender": "79876543210@c.us", "senderName": "John"},
  "messageData": {
    "typeMessage": "textMessage",
    "textMessageData": {"textMessage": "Hello"}
  }
}
```

Main `typeWebhook` values (each has a webhook toggle in `SetSettings`, see [account.md](account.md)):

- `incomingMessageReceived` — incoming message
- `outgoingMessageReceived` — message sent from the phone
- `outgoingAPIMessageReceived` — message sent via API
- `outgoingMessageStatus` — sent/delivered/read status of outgoing messages
- `stateInstanceChanged` — instance authorization state changed
- `incomingCall` — incoming call
- `pollMessage` / poll update webhooks — poll creation and votes
- `incomingBlock` — chat added to blocked list
- `editedMessage`, `deletedMessage` — message edited / deleted
- `deviceInfo` — device/battery info
- `quotaExceeded` — plan quota exceeded

Frequent `messageData.typeMessage` values for incoming messages: `textMessage`, `extendedTextMessage` (text with link preview / quote — text is in `extendedTextMessageData.text`), `imageMessage`, `videoMessage`, `documentMessage`, `audioMessage` (file info in `fileMessageData`: `downloadUrl`, `caption`, `fileName`, `mimeType`), `locationMessage`, `contactMessage`, `pollMessage`, `reactionMessage`, `stickerMessage`, `buttonsResponseMessage` / `templateButtonsReplyMessage` / `interactiveButtonsReply` (button replies), `quotedMessage`.

Group chats: `senderData.chatId` ends with `@g.us`, `senderData.sender` is the participant who wrote (`@c.us`).

## Queues — `client.Queues()`

### ShowMessagesQueue

Docs: https://green-api.com/en/docs/api/queues/ShowMessagesQueue/

```go
ShowMessagesQueue()
```

Returns the list of messages waiting in the send queue (see sending delay in SKILL.md).

### ClearMessagesQueue

Docs: https://green-api.com/en/docs/api/queues/ClearMessagesQueue/

```go
ClearMessagesQueue()
```

Response Body: `{"isCleared": true}`.

### GetWebhooksCount

Docs: https://green-api.com/en/docs/api/queues/GetWebhooksCount/

```go
GetWebhooksCount()
```

Response Body: `{"count": 0}` — number of notifications in the incoming queue.

### ClearWebhooksQueue

Docs: https://green-api.com/en/docs/api/queues/ClearWebhooksQueue/
Rate limit: once per 60 seconds.

```go
ClearWebhooksQueue()
```

Response Body: `{"isCleared": true, "reason": ""}`.
