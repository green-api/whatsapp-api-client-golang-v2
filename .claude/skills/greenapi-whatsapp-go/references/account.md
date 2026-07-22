# Account & Partner

## Account — `client.Account()`

### GetSettings

Docs: https://green-api.com/en/docs/api/account/GetSettings/

```go
GetSettings()
```

Returns current instance settings (webhookUrl, delaySendMessagesMilliseconds, webhook toggles, etc.).

### SetSettings

Docs: https://green-api.com/en/docs/api/account/SetSettings/
Applies **asynchronously**: the instance reboots and settings take effect within ~5 minutes. Response Body: `{"saveSettings": true}`. Pass only the options you want to change:

```go
SetSettings(options ...SetSettingsOption)
```

- `OptionalWebhookUrl(webhookUrl string)` — URL to POST notifications to (empty string disables)
- `OptionalWebhookUrlToken(webhookUrlToken string)` — token sent in `Authorization` header to your webhook server
- `OptionalDelaySendMessages(delaySendMessagesMilliseconds uint)` — delay between queued sends, min 500 ms
- `OptionalMarkIncomingMessagesRead(markIncomingMessagesReaded bool)` — mark incoming as read
- `OptionalMarkIncomingMessagesReadOnReply(markIncomingMessagesReadedOnReply bool)` — mark as read when replying via API
- `OptionalKeepOnlineStatus(keepOnlineStatus bool)` — keep account status "Online"
- `OptionalAutoTyping(autoTyping uint)` — automatic typing indication length
- `OptionalLinkPreviewBool(linkPreview bool)` — default link preview behavior
- `OptionalEnableLidMode(enableLidMode bool)` — enable LID identifiers mode
- Webhook toggles (all `bool`): `OptionalOutgoingWebhook` (statuses of outgoing messages), `OptionalOutgoingMessageWebhook` (messages sent from phone), `OptionalOutgoingAPIMessageWebhook` (messages sent via API), `OptionalStateWebhook` (instance state changes), `OptionalIncomingWebhook` (incoming messages/files), `OptionalDeviceWebhook` (device/battery), `OptionalPollMessageWebhook` (polls and votes), `OptionalIncomingBlockWebhook` (chat blocked), `OptionalIncomingCallWebhook` (incoming calls), `OptionalEditedMessageWebhook`, `OptionalDeletedMessageWebhook`

### GetStateInstance

Docs: https://green-api.com/en/docs/api/account/GetStateInstance/

```go
GetStateInstance()
```

Response Body: `{"stateInstance": "authorized"}`. Values: `notAuthorized`, `authorized`, `blocked`, `sleepMode`, `starting`, `yellowCard`. Only `authorized` instances can send.

### GetStatusInstance

Docs: https://green-api.com/en/docs/api/account/GetStatusInstance/

```go
GetStatusInstance()
```

Socket connection status of the instance (`online`/`offline`).

### GetStateInstanceHistory

Docs: https://green-api.com/en/docs/api/account/GetStateInstanceHistory/

```go
GetStateInstanceHistory(options ...GetStateInstanceHistoryOption)
```

- `OptionalHistoryCount(count int)` — number of records (default 100)

Response: array of `{stateInstance, timestamp, phoneNumber}`.

### Reboot / Logout

Docs: https://green-api.com/en/docs/api/account/Reboot/ , https://green-api.com/en/docs/api/account/Logout/

```go
Reboot() // Body: {"isReboot": true}
Logout() // Body: {"isLogout": true} — unlinks the WhatsApp account; re-authorization needed
```

### QR — get authorization QR code

Docs: https://green-api.com/en/docs/api/account/QR/
Only works while the instance is `notAuthorized`. Response Body: `{"type":"qrCode","message":"<base64 image>"}` — decode and render; QR refreshes every ~20 s, so poll in a loop until authorized.

```go
QR()
```

### GetAuthorizationCode — authorize by phone number

Docs: https://green-api.com/en/docs/api/account/GetAuthorizationCode/

```go
GetAuthorizationCode(phoneNumber int) // e.g. 79876543210
```

Response Body: `{"status": true, "code": "..."}` — enter the code in WhatsApp ("Link with phone number").

### SetProfilePicture

Docs: https://green-api.com/en/docs/api/account/SetProfilePicture/

```go
SetProfilePicture(filepath string) // local path to a JPEG image
```

### UpdateApiToken

Docs: https://green-api.com/en/docs/api/account/UpdateApiToken/
Rotates the instance API token. Response Body: `{"apiTokenInstance": "new-token"}` — the old token stops working, update your config.

```go
UpdateApiToken()
```

### GetWaSettings

Docs: https://green-api.com/en/docs/api/account/GetWaSettings/

```go
GetWaSettings()
```

Info about the WhatsApp account on the instance: `{"avatar":"...","phone":"79876543210","stateInstance":"authorized","deviceId":"..."}`.

## Partner — separate client

Partner methods (instance management via API) use a different client and token (docs: https://green-api.com/en/docs/partners/):

```go
partner := greenapi.GreenAPIPartner{
	PartnerToken: "gac.1234567891234567891234567891213456789",
	Email:        "mail@example.com", // used by some partner operations
}
```

### GetInstances

Docs: https://green-api.com/en/docs/partners/getInstances/

```go
partner.Partner().GetInstances()
```

Array of all instances created by the partner.

### CreateInstance

Docs: https://green-api.com/en/docs/partners/createInstance/

```go
partner.Partner().CreateInstance(options ...any)
```

Accepts `OptionalName(name string)` plus any of the `SetSettings` options listed above (webhook URL, delays, webhook toggles) to pre-configure the new instance. Response contains `idInstance` and `apiTokenInstance`.

### DeleteInstanceAccount

Docs: https://green-api.com/en/docs/partners/deleteInstanceAccount/

```go
partner.Partner().DeleteInstanceAccount(idInstance uint)
```

Response Body: `{"deleteInstanceAccount": true}`.
