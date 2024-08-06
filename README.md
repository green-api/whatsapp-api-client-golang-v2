# whatsapp-api-client-golang-v2

- [Документация на русском языке](docs/README_RU.md).

`whatsapp-api-client-golang-v2` is a library for integration with WhatsApp messenger using the API
service [green-api.com](https://green-api.com/en/). You should get a registration token and an account ID in
your [personal cabinet](https://console.green-api.com/) to use the library. There is a free developer account tariff.

You can find the `v1` version here - https://github.com/green-api/whatsapp-api-client-golang

## API

The documentation for the REST API can be found at the [link](https://green-api.com/en/docs/api). The library is a wrapper
for the REST API, so the documentation at the link above also applies.

## Support links

[![Support](https://img.shields.io/badge/support@green--api.com-D14836?style=for-the-badge&logo=gmail&logoColor=white)](mailto:support@greenapi.com)
[![Support](https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white)](https://t.me/greenapi_support_eng_bot)
[![Support](https://img.shields.io/badge/WhatsApp-25D366?style=for-the-badge&logo=whatsapp&logoColor=white)](https://wa.me/77273122366)

## Guides & News

[![Guides](https://img.shields.io/badge/YouTube-%23FF0000.svg?style=for-the-badge&logo=YouTube&logoColor=white)](https://www.youtube.com/@greenapi-en)
[![News](https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white)](https://t.me/green_api)
[![News](https://img.shields.io/badge/WhatsApp-25D366?style=for-the-badge&logo=whatsapp&logoColor=white)](https://whatsapp.com/channel/0029VaLj6J4LNSa2B5Jx6s3h)

#### Authorization

To send a message or perform other Green API methods, the WhatsApp account in the phone app must be authorized. To
authorize the account, go to your [cabinet](https://console.green-api.com/) and scan the QR code using the WhatsApp app.

## Installation

**Make sure that you have Go installed with a version of 1.20 or newer**
```shell
go version
```

**Create a module for your project if you didn't:**

```shell
go mod init ModuleName
```

**Install the library:**

```shell
go get github.com/green-api/whatsapp-api-client-golang-v2
```

**Import:**

```go
import (
	greenapi "github.com/green-api/whatsapp-api-client-golang-v2"
)
```

## Usage and examples

**How to initialize an object:**

```go
GreenAPI := greenapi.GreenAPI{
		APIURL:           "https://api.green-api.com",
		MediaURL:         "https://media.green-api.com",
		IDInstance:       "1101000001",
		APITokenInstance: "d75b3a66374942c5b3c019c698abc2067e151558acbd412345",
	}
```

All methods of this library return two objects: `*APIResponse` and `error`. 

You can see the `APIResponse` format in the [types.go](types.go)

**How to send a message:**

Link to example: [sendMessage/main.go](examples/sendMessage/main.go)

```go
response, _ := GreenAPI.Sending().SendMessage(
		"11001234567@c.us",
		"Hello",
	)
```

**How to create a group:**

Link to example: [createGroup/main.go](examples/createGroup/main.go)

```go
response, _ := GreenAPI.Groups().CreateGroup(
		"Group Title",
		[]string{
			"11001211111@c.us",
			"11001222222@c.us",
			"11001233333@c.us",
		},
	)
```

**How to send file by upload:**

Link to example: [sendFileByUpload/main.go](examples/sendFileByUpload/main.go)

```go
response, _ := GreenAPI.Sending().SendFileByUpload(
		"11001234567@c.us",
		"C:/Users/user/Desktop/Pictures/image.png",
		"image.png",
	)
```

**How to send a message with a poll:**

Link to example: [sendPoll/main.go](examples/sendPoll/main.go)

```go
response, _ := GreenAPI.Sending().SendPoll(
		"11001234567@c.us", 
		"Choose a color:", 
		[]string{"Red", "Green", "Blue"}, 
	)
```

**How to send a text status:**

Link to example: [sendTextStatus](examples/sendTextStatus/main.go)

```go
response, _ := GreenAPI.Statuses().SendTextStatus(
		"Text of the status", 
		greenapi.OptionalFont("SERIF"),
		greenapi.OptionalBackgroundColorText("#87CEEB"),
		//greenapi.OptionalParticipantsTextStatus([]string{"1234567890@c.us", "1234567890@c.us"}),
	)
```

## Partner methods

**To use partner methods you have to initialize another object:**

```go
Partner := greenapi.GreenAPIPartner{
		PartnerToken: "gac.1234567891234567891234567891213456789",
		Email: "mail@email.com", // email is optional
	}
```

**Now you can use Partner methods as usual methods, but through the "Partner" object:**

**How to get instances:**

Link to the example: [partnerMethods/getInstances/main.go](examples/partnerMethods/getInstances/main.go)

```go
response, _ := Partner.Partner().GetInstances()
```

**How to create an instance:**

Link to the example: [partnerMethods/createInstance/main.go](examples/partnerMethods/createInstance/main.go)

```go
response, _ := Partner.Partner().CreateInstance(
		greenapi.OptionalWebhookUrl("webhook_url"),
		greenapi.OptionalWebhookUrlToken("auth_token"),
		greenapi.OptionalDelaySendMessages(5000),
		greenapi.OptionalMarkIncomingMessagesRead(true),
		greenapi.OptionalMarkIncomingMessagesReadOnReply(true),
		greenapi.OptionalOutgoingWebhook(true),
		greenapi.OptionalOutgoingMessageWebhook(true),
		greenapi.OptionalOutgoingAPIMessageWebhook(true),
		greenapi.OptionalStateWebhook(true),
		greenapi.OptionalIncomingWebhook(true),
		greenapi.OptionalDeviceWebhook(true),
		greenapi.OptionalKeepOnlineStatus(true),
		greenapi.OptionalPollMessageWebhook(true),
		greenapi.OptionalIncomingBlockWebhook(true),
		greenapi.OptionalIncomingCallWebhook(true),
	)
```

**How to delete an instance:**

Link to the example: [partnerMethods/deleteInstanceAccount/main.go](examples/partnerMethods/deleteInstanceAccount/main.go)

```go
response, _ := Partner.Partner().DeleteInstanceAccount(1101000000)
```

## Optional parameters

**Note that functions might have optional arguments, which you can pass or ignore. Optional parameters are passed as functions into the method's arguments and have similar naming format:**
```go
greenapi.Optional + name of parameter
```

**For example, in the `SetSettings` method all the arguments are optional. Here is an example of how it works:**

```go
response, _ := GreenAPI.Account().SetSettings(
        greenapi.OptionalDelaySendMessages(5000),
		greenapi.OptionalOutgoingWebhook(true),
		greenapi.OptionalIncomingWebhook(true),
		// greenapi.OptionalWebhookUrl("webhook_url"),
		// greenapi.OptionalWebhookUrlToken("auth_token"),
		// greenapi.OptionalMarkIncomingMessagesRead(true),
		// greenapi.OptionalMarkIncomingMessagesReadOnReply(true),
		// greenapi.OptionalOutgoingMessageWebhook(true),
		// greenapi.OptionalOutgoingAPIMessageWebhook(true),
		// greenapi.OptionalStateWebhook(true),
		// greenapi.OptionalDeviceWebhook(true),
		// greenapi.OptionalKeepOnlineStatus(true),
		// greenapi.OptionalPollMessageWebhook(true),
		// greenapi.OptionalIncomingBlockWebhook(true),
		// greenapi.OptionalIncomingCallWebhook(true),
	)
```

In this example, only `DelaySendMessages`, `OutgoingWebhook` and `IncomingWebhook` settings will be changed, other settings are commented so they will not be passed. However, you can uncomment any setting that you prefer. **The settings that were not used will not be affected**

One more example of using optional parameters, this time let's use `sendMessage` method:

```go
response, _ := GreenAPI.Sending().SendMessage(
		"11001234567@c.us",
		"Hello",
		greenapi.OptionalLinkPreview(false), // turns off link preview if there is any
		greenapi.OptionalQuotedMessageId("BAE59673E71FC5DB"), // quotes specified message
	)
```

## List of examples

| Description                                   | Link to example                                               |
|-----------------------------------------------|---------------------------------------------------------------|
| How to send a message                         | [sendMessage/main.go](examples/sendMessage/main.go)           |
| How to send a file by uploading from the disk | [sendFileByUpload/main.go](examples/sendFileByUpload/main.go) |
| How to upload a file to an external drive                     | [uploadFile/main.go](examples/uploadFile/main.go)       |
| How to send a poll                         | [sendPoll/main.go](examples/sendPoll/main.go)           |
| How to check if there is a WhatsApp account on the phone number         | [checkWhatsapp/main.go](examples/checkWhatsapp/main.go)                   |
| How to set instance settings             | [setSettings/main.go](examples/setSettings/main.go)                 |
| How to create a group             | [createGroup/main.go](examples/createGroup/main.go)                 |
| How to send a text status             | [sendTextStatus/main.go](examples/sendTextStatus/main.go)                 |
| How to get all instances of the account             | [partnerMethods/getInstances/main.go](examples/partnerMethods/getInstances/main.go)                 |
| How to create an instance             | [partnerMethods/createInstance/main.go](examples/partnerMethods/createInstance/main.go)                 |
| How to delete an instance            | [partnerMethods/deleteInstanceAccount/main.go](examples/partnerMethods/deleteInstanceAccount/main.go)                 |

## List of all library methods

| API method                        | Description                                                                                                               | Documentation link                                                                                          |
|-----------------------------------|---------------------------------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------|
| `Account().GetSettings`           | The method is designed to get the current settings of the account                                                         | [GetSettings](https://green-api.com/en/docs/api/account/GetSettings/)                                       |
| `Account().GetWaSettings`         | The method is designed to get information about the WhatsApp account                                                      | [GetSettings](https://green-api.com/en/docs/api/account/GetWaSettings/)                                     |
| `Account().SetSettings`           | The method is designed to set the account settings                                                                        | [SetSettings](https://green-api.com/docs/api/account/SetSettings/)                                          |
| `Account().GetStateInstance`      | The method is designed to get the state of the account                                                                    | [GetStateInstance](https://green-api.com/en/docs/api/account/GetStateInstance/)                             |
| `Account().Reboot`                | The method is designed to restart the account                                                                             | [Reboot](https://green-api.com/en/docs/api/account/Reboot/)                                                 |
| `Account().Logout`                | The method is designed to unlogin the account                                                                             | [Logout](https://green-api.com/en/docs/api/account/Logout/)                                                 |
| `Account().QR`                    | The method is designed to get a QR code                                                                                   | [QR](https://green-api.com/en/docs/api/account/QR/)                                                         |
| `Account().SetProfilePicture`     | The method is designed to set the avatar of the account                                                                   | [SetProfilePicture](https://green-api.com/en/docs/api/account/SetProfilePicture/)                           |
| `Account().GetAuthorizationCode`  | The method is designed to authorize an instance by phone number                                                           | [GetAuthorizationCode](https://green-api.com/en/docs/api/account/GetAuthorizationCode/)                     |                                  |
| `Groups().CreateGroup`            | The method is designed to create a group chat                                                                             | [CreateGroup](https://green-api.com/en/docs/api/groups/CreateGroup/)                                        |
| `Groups().UpdateGroupName`        | The method changes the name of the group chat                                                                             | [UpdateGroupName](https://green-api.com/en/docs/api/groups/UpdateGroupName/)                                |
| `Groups().GetGroupData`           | The method gets group chat data                                                                                           | [GetGroupData](https://green-api.com/en/docs/api/groups/GetGroupData/)                                      |
| `Groups().AddGroupParticipant`    | The method adds a participant to the group chat                                                                           | [AddGroupParticipant](https://green-api.com/en/docs/api/groups/AddGroupParticipant/)                        |
| `Groups().RemoveGroupParticipant` | The method removes the participant from the group chat                                                                    | [RemoveGroupParticipant](https://green-api.com/en/docs/api/groups/RemoveGroupParticipant/)                  |
| `Groups().SetGroupAdmin`          | The method designates a member of a group chat as an administrator                                                        | [SetGroupAdmin](https://green-api.com/en/docs/api/groups/SetGroupAdmin/)                                    |
| `Groups().RemoveAdmin`            | The method deprives the participant of group chat administration rights                                                   | [RemoveAdmin](https://green-api.com/en/docs/api/groups/RemoveAdmin/)                                        |
| `Groups().SetGroupPicture`        | The method sets the avatar of the group                                                                                   | [SetGroupPicture](https://green-api.com/en/docs/api/groups/SetGroupPicture/)                                |
| `Groups().LeaveGroup`             | The method logs the user of the current account out of the group chat                                                     | [LeaveGroup](https://green-api.com/en/docs/api/groups/LeaveGroup/)                                          |
| `Journals().GetChatHistory`       | The method returns the chat message history                                                                               | [GetChatHistory](https://green-api.com/en/docs/api/journals/GetChatHistory/)                                |
| `Journals().GetMessage`           | The method returns a chat message                                                                                         | [GetMessage](https://green-api.com/en/docs/api/journals/GetMessage/)                                        |
| `Journals().LastIncomingMessages` | The method returns the most recent incoming messages of the account                                                       | [LastIncomingMessages](https://green-api.com/en/docs/api/journals/LastIncomingMessages/)                    |
| `Journals().LastOutgoingMessages` | The method returns the last sent messages of the account                                                                  | [LastOutgoingMessages](https://green-api.com/en/docs/api/journals/LastOutgoingMessages/)                    |
| `Queues().ShowMessagesQueue`      | The method is designed to get the list of messages that are in the queue to be sent                                       | [ShowMessagesQueue](https://green-api.com/en/docs/api/queues/ShowMessagesQueue/)                            |
| `Queues().ClearMessagesQueue`     | The method is designed to clear the queue of messages to be sent                                                          | [ClearMessagesQueue](https://green-api.com/en/docs/api/queues/ClearMessagesQueue/)                          |
| `ReadMark().ReadChat`             | The method is designed to mark chat messages as read                                                                      | [ReadChat](https://green-api.com/en/docs/api/marks/ReadChat/)                                               |
| `Receiving().ReceiveNotification` | The method is designed to receive a single incoming notification from the notification queue                              | [ReceiveNotification](https://green-api.com/en/docs/api/receiving/technology-http-api/ReceiveNotification/) |
| `Receiving().DeleteNotification`  | The method is designed to remove an incoming notification from the notification queue                                     | [DeleteNotification](https://green-api.com/en/docs/api/receiving/technology-http-api/DeleteNotification/)   |
| `Receiving().DownloadFile`        | The method is for downloading received and sent files                                                                     | [DownloadFile](https://green-api.com/en/docs/api/receiving/files/DownloadFile/)                             |
| `Sending().SendMessage`           | The method is designed to send a text message to a personal or group chat                                                 | [SendMessage](https://green-api.com/en/docs/api/sending/SendMessage/)                                       |
| `Sending().SendFileByUpload`      | The method is designed to send a file loaded through a form (form-data)                                                   | [SendFileByUpload](https://green-api.com/en/docs/api/sending/SendFileByUpload/)                             |
| `Sending().SendFileByUrl`         | The method is designed to send a file downloaded via a link                                                               | [SendFileByUrl](https://green-api.com/en/docs/api/sending/SendFileByUrl/)                                   |
| `Sending().UploadFile`            | The method allows you to upload a file from the local file system, which can later be sent using the SendFileByUrl method | [UploadFile](https://green-api.com/en/docs/api/sending/UploadFile/)                                         |
| `Sending().SendLocation`          | The method is designed to send a geolocation message                                                                      | [SendLocation](https://green-api.com/en/docs/api/sending/SendLocation/)                                     |
| `Sending().SendContact`           | The method is for sending a message with a contact                                                                        | [SendContact](https://green-api.com/en/docs/api/sending/SendContact/)                                       |
| `Sending().ForwardMessages`       | The method is designed for forwarding messages to a personal or group chat                                                | [ForwardMessages](https://green-api.com/en/docs/api/sending/ForwardMessages/)                               |
| `Sending().SendPoll`              | The method is designed for sending messages with a poll to a private or group chat                                        | [SendPoll](https://green-api.com/en/docs/api/sending/SendPoll/)                                             |
| `Service().CheckWhatsapp`         | The method checks if there is a WhatsApp account on the phone number                                                      | [CheckWhatsapp](https://green-api.com/en/docs/api/service/CheckWhatsapp/)                                   |
| `Service().GetAvatar`             | The method returns the avatar of the correspondent or group chat                                                          | [GetAvatar](https://green-api.com/en/docs/api/service/GetAvatar/)                                           |
| `Service().GetContacts`           | The method is designed to get a list of contacts of the current account                                                   | [GetContacts](https://green-api.com/en/docs/api/service/GetContacts/)                                       |
| `Service().GetContactInfo`        | The method is designed to obtain information about the contact                                                            | [GetContactInfo](https://green-api.com/en/docs/api/service/GetContactInfo/)                                 |
| `Service().DeleteMessage`         | The method deletes the message from chat                                                                                  | [DeleteMessage](https://green-api.com/en/docs/api/service/deleteMessage/)                                   |
| `Service().ArchiveChat`           | The method archives the chat                                                                                              | [ArchiveChat](https://green-api.com/en/docs/api/service/archiveChat/)                                       |
| `Service().UnarchiveChat`         | The method unarchives the chat                                                                                            | [UnarchiveChat](https://green-api.com/en/docs/api/service/unarchiveChat/)                                   |
| `Service().SetDisappearingChat`   | The method is designed to change the settings of disappearing messages in chats                                           | [SetDisappearingChat](https://green-api.com/en/docs/api/service/SetDisappearingChat/)                       |
| `Partner().GetInstances`   | The method is for getting all the account instances created by the partner.                                           | [GetInstances](https://green-api.com/en/docs/partners/getInstances/)                       |
| `Partner().CreateInstance`   | The method is for creating an instance.                                           | [CreateInstance](https://green-api.com/en/docs/partners/createInstance/)                       |
| `Partner().DeleteInstanceAccount`   | The method is for deleting an instance.                                           | [DeleteInstanceAccount](https://green-api.com/en/docs/partners/deleteInstanceAccount/)                       |
| `Statuses().SendTextStatus`             | The method is aimed for sending a text status                                                     | [SendTextStatus](https://green-api.com/en/docs/api/statuses/SendTextStatus/)                                          |
| `Statuses().SendVoiceStatus`             | The method is aimed for sending a voice status                                                     | [SendVoiceStatus](https://green-api.com/en/docs/api/statuses/SendVoiceStatus/)                                          |
| `Statuses().SendMediaStatus`             | The method is aimed for sending a voice status                                                     | [SendMediaStatus](https://green-api.com/en/docs/api/statuses/SendMediaStatus/)                                          |      
| `Statuses().GetOutgoingStatuses`             | The method returns the outgoing statuses of the account                                                     | [GetOutgoingStatuses](https://green-api.com/en/docs/api/statuses/GetOutgoingStatuses/)                                          |      
| `Statuses().GetIncomingStatuses`             | The method returns the incoming status messages of the account                                                     | [GetIncomingStatuses](https://green-api.com/en/docs/api/statuses/GetIncomingStatuses/)                                          |      
| `Statuses().GetStatusStatistic`             | The method returns an array of recipients marked for a given status.                                                     | [GetStatusStatistic](https://green-api.com/en/docs/api/statuses/GetStatusStatistic/)                                          |      
| `Statuses().DeleteStatus`             | The method is aimed for deleting status.                                                     | [DeleteStatus](https://green-api.com/en/docs/api/statuses/DeleteStatus/)                                          |    