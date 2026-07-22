# Service, Journals, Contacts, ReadMark

## Service — `client.Service()`

### CheckWhatsapp — does a number have WhatsApp

Docs: https://green-api.com/en/docs/api/service/CheckWhatsapp/

```go
CheckWhatsapp(phoneNumber int, options ...CheckWhatsappOption)
```

`phoneNumber` (int, e.g. `79876543210`) is **deprecated** — pass `0` and use:

- `OptionalChatID(chatID string)` — preferred, `"79876543210@c.us"`
- `OptionalForce(force bool)` — bypass cache (default false)

If both `phoneNumber == 0` and no `OptionalChatID`, the SDK returns an error. Response Body: `{"existsWhatsapp": true, "chatId": "79876543210@c.us"}`.

### GetAvatar

Docs: https://green-api.com/en/docs/api/service/GetAvatar/

```go
GetAvatar(chatId string) // user or group avatar; Body: {"existsWhatsapp":true,"urlAvatar":"...","reason":""}
```

### GetContacts

Docs: https://green-api.com/en/docs/api/service/GetContacts/

```go
GetContacts(options ...GetContactsOption)
```

- `OptionalGetContactsGroup(group bool)` — `true` = only groups, `false` = only personal contacts
- `OptionalGetContactsCount(count int)` — limit number of results

Returns array of `{id, name, contactName, type}` for the account's chats/contacts.

### GetContactInfo

Docs: https://green-api.com/en/docs/api/service/GetContactInfo/

```go
GetContactInfo(chatId string) // avatar, name, email, category, products, etc.
```

### DeleteMessage

Docs: https://green-api.com/en/docs/api/service/deleteMessage/

```go
DeleteMessage(chatId, idMessage string, options ...DeleteMessageOption)
```

- `OptionalOnlySenderDelete(OnlySenderDelete bool)` — delete only on sender's side (default: for everyone)

### EditMessage

Docs: https://green-api.com/en/docs/api/service/editMessage/

```go
EditMessage(chatId, idMessage, message string) // Body: {"idMessage": "..."}
```

### ArchiveChat / UnarchiveChat

Docs: https://green-api.com/en/docs/api/service/archiveChat/ , .../unarchiveChat/

```go
ArchiveChat(chatId string)
UnarchiveChat(chatId string)
```

### SetDisappearingChat

Docs: https://green-api.com/en/docs/api/service/SetDisappearingChat/

```go
SetDisappearingChat(chatId string, ephemeralExpiration int)
```

`ephemeralExpiration` seconds; allowed values: `0` (off), `86400` (24 h), `604800` (7 days), `7776000` (90 days).

### SendTyping

Docs: https://green-api.com/en/docs/api/service/SendTyping/ (typing/recording indication)

```go
SendTyping(chatId string, options ...SendTypingOption)
```

- `OptionalSendTypingTime(typingTime int)` — 1000–20000 ms
- `OptionalSendTypingType(typingType string)` — `"typing"` or `"recording"` (default `"typing"`; SDK-validated)

### GetChats

Docs: https://green-api.com/en/docs/api/service/GetChats/

```go
GetChats(options ...GetChatsOption)
```

- `OptionalChatsCount(count int)` — limit number of chats

Returns array of `{archive, id, ephemeralExpiration, ephemeralSettingTimestamp, name, type}`.

## Journals — `client.Journals()`

### GetChatHistory

Docs: https://green-api.com/en/docs/api/journals/GetChatHistory/

```go
GetChatHistory(chatId string, options ...GetChatHistoryOption)
```

- `OptionalCount(count int)` — number of messages (default 100)

### GetMessage

Docs: https://green-api.com/en/docs/api/journals/GetMessage/

```go
GetMessage(chatId, idMessage string)
```

### LastIncomingMessages / LastOutgoingMessages

Docs: https://green-api.com/en/docs/api/journals/LastIncomingMessages/ , .../LastOutgoingMessages/

```go
LastIncomingMessages(options ...LastMessagesOption)
LastOutgoingMessages(options ...LastMessagesOption)
```

- `OptionalMinutes(minutes int)` — look-back period in minutes (default 1440 = 24 h)

### LastIncomingCalls / LastOutgoingCalls

Docs: https://green-api.com/en/docs/api/journals/LastIncomingCalls/ , .../LastOutgoingCalls/

```go
LastIncomingCalls(options ...LastCallsOption)
LastOutgoingCalls(options ...LastCallsOption)
```

- `OptionalCallsMinutes(minutes int)` — look-back period in minutes (default 1440)

Returns array of call records `{type, idMessage, timestamp, typeMessage, chatId, isVideo, status, ...}`.

## Contacts — `client.Contacts()`

### AddContact

Docs: https://green-api.com/en/docs/api/service/AddContact/

```go
AddContact(chatId, firstName string, options ...AddContactOption)
```

- `OptionalAddLastName(lastName string)`
- `OptionalAddSaveInAddressbook(saveInAddressbook bool)`

### EditContact

Docs: https://green-api.com/en/docs/api/service/EditContact/

```go
EditContact(chatId, firstName string, options ...EditContactOption)
```

- `OptionalEditLastName(lastName string)`
- `OptionalEditSaveInAddressbook(saveInAddressbook bool)`

### DeleteContact

Docs: https://green-api.com/en/docs/api/service/DeleteContact/

```go
DeleteContact(chatId string)
```

## ReadMark — `client.ReadMark()`

### ReadChat — mark messages as read

Docs: https://green-api.com/en/docs/api/marks/ReadChat/

```go
ReadChat(chatId string, options ...ReadChatOption)
```

- `OptionalIdMessage(idMessage string)` — mark one specific message; without it the whole chat is marked read

Response Body: `{"setRead": true}`.
