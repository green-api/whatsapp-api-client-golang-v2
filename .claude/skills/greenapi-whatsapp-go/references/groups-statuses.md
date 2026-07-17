# Groups & Statuses

## Groups — `client.Groups()`

`groupId` ends with `@g.us`; participant ids end with `@c.us`. Get `groupId` from `CreateGroup`'s response, `GetContacts`, `GetChats`, or incoming notifications — never construct it manually.

### CreateGroup

Docs: https://green-api.com/en/docs/api/groups/CreateGroup/

```go
CreateGroup(groupName string, chatIds []string)
```

`chatIds` — participants, e.g. `[]string{"79876543210@c.us", "79876543211@c.us"}`. Response Body: `{"created": true, "chatId": "120363043968066561@g.us", "groupInviteLink": "https://chat.whatsapp.com/..."}`.

### UpdateGroupName

Docs: https://green-api.com/en/docs/api/groups/UpdateGroupName/

```go
UpdateGroupName(groupId, groupName string)
```

### GetGroupData

Docs: https://green-api.com/en/docs/api/groups/GetGroupData/

```go
GetGroupData(groupId string)
```

Returns subject, owner, creation time, participants (with `isAdmin` flags), invite link.

### UpdateGroupSettings

Docs: https://green-api.com/en/docs/api/groups/UpdateGroupSettings/

```go
UpdateGroupSettings(groupId string, options ...UpdateGroupSettingsOption)
```

- `OptionalAllowParticipantsEditGroupSettings(allow bool)`
- `OptionalAllowParticipantsSendMessages(allow bool)`

Response Body: `{"updateGroupSettings": true, "reason": ""}`.

### AddGroupParticipant / RemoveGroupParticipant

Docs: https://green-api.com/en/docs/api/groups/AddGroupParticipant/ , .../RemoveGroupParticipant/

```go
AddGroupParticipant(groupId, participantChatId string)
RemoveGroupParticipant(groupId, participantChatId string)
```

The instance account must be a group admin.

### SetGroupAdmin / RemoveAdmin

Docs: https://green-api.com/en/docs/api/groups/SetGroupAdmin/ , .../RemoveAdmin/

```go
SetGroupAdmin(groupId, participantChatId string)
RemoveAdmin(groupId, participantChatId string)
```

### SetGroupPicture

Docs: https://green-api.com/en/docs/api/groups/SetGroupPicture/
Note the argument order: file path first.

```go
SetGroupPicture(filepath, groupId string) // local path to a JPEG image
```

### LeaveGroup

Docs: https://green-api.com/en/docs/api/groups/LeaveGroup/

```go
LeaveGroup(groupId string)
```

## Statuses — `client.Statuses()`

WhatsApp statuses (stories). `participants` — contact ids (`@c.us`) who will see the status; empty = all contacts.

### SendTextStatus

Docs: https://green-api.com/en/docs/api/statuses/SendTextStatus/

```go
SendTextStatus(message string, options ...SendTextStatusOption)
```

`message` max 500 chars (SDK-validated). Options:

- `OptionalBackgroundColorText(backgroundColor string)` — hex color, e.g. `"#87CEEB"`
- `OptionalFont(font string)` — `SERIF`, `SANS_SERIF`, `NORICAN_REGULAR`, `BRYNDAN_WRITE`, `OSWALD_HEAVY`
- `OptionalParticipantsTextStatus(participants []string)`

### SendVoiceStatus

Docs: https://green-api.com/en/docs/api/statuses/SendVoiceStatus/

```go
SendVoiceStatus(urlFile, fileName string, options ...SendVoiceStatusOption)
```

`urlFile` — direct link to an audio file; `fileName` with extension. Options:

- `OptionalBackgroundColorVoice(backgroundColor string)`
- `OptionalParticipantsVoiceStatus(participants []string)`

### SendMediaStatus

Docs: https://green-api.com/en/docs/api/statuses/SendMediaStatus/

```go
SendMediaStatus(urlFile, fileName string, options ...SendMediaStatusOption)
```

Image or video by direct URL. Options:

- `OptionalCaptionMediaStatus(caption string)` — max 1024 chars (SDK-validated)
- `OptionalParticipantsMediaStatus(participants []string)`

### DeleteStatus

Docs: https://green-api.com/en/docs/api/statuses/DeleteStatus/

```go
DeleteStatus(idMessage string)
```

### GetStatusStatistic

Docs: https://green-api.com/en/docs/api/statuses/GetStatusStatistic/

```go
GetStatusStatistic(idMessage string) // who viewed the status
```

### GetOutgoingStatuses / GetIncomingStatuses

Docs: https://green-api.com/en/docs/api/statuses/GetOutgoingStatuses/ , .../GetIncomingStatuses/

```go
GetOutgoingStatuses(options ...GetLastStatusesOption)
GetIncomingStatuses(options ...GetLastStatusesOption)
```

- `OptionalMinutesOfStatuses(minutes int)` — look-back period in minutes (default 1440)
