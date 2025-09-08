# max-api-client-golang

- [Документация на русском языке](docs/README_RU.md).

`max-api-client-golang` is a library for integration with MAX messenger using the API
service [green-api.com](https://green-api.com/v3/). You should get a registration token and an account ID in
your [personal cabinet](https://console.green-api.com/) to use the library. There is a free developer account tariff.



## API

The documentation for the REST API can be found at the [link](https://green-api.com/v3/docs/api). The library is a wrapper
for the REST API, so the documentation at the link above also applies.

## Support links

[![Support](https://img.shields.io/badge/support@green--api.com-D14836?style=for-the-badge&logo=gmail&logoColor=white)](mailto:support@greenapi.com)
[![Support](https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white)](https://t.me/greenapi_support_eng_bot)
[![Support](https://img.shields.io/badge/MAX-25D366?style=for-the-badge&logo=whatsapp&logoColor=white)](https://wa.me/77273122366)

## Guides & News

[![Guides](https://img.shields.io/badge/YouTube-%23FF0000.svg?style=for-the-badge&logo=YouTube&logoColor=white)](https://www.youtube.com/@greenapi-en)
[![News](https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white)](https://t.me/green_api)
[![News](https://img.shields.io/badge/MAX-25D366?style=for-the-badge&logo=whatsapp&logoColor=white)](https://whatsapp.com/channel/0029VaLj6J4LNSa2B5Jx6s3h)

#### Authorization

To send a message or perform other Green API methods, the MAX account in the phone app must be authorized. To
authorize the account, go to your [cabinet](https://console.green-api.com/) and scan the QR code using the MAX app.

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
go get github.com/green-api/max-api-client-golang
```

**Import:**

```go
import (
	greenapi "github.com/green-api/max-api-client-golang"
)
```

## Usage and examples

**How to initialize an object:**

```go
GreenAPI := greenapi.GreenAPI{
		APIURL:           "https://api.green-api.com/v3",
		MediaURL:         "https://api.green-api.com/v3",
		IDInstance:       "3100000001",
		APITokenInstance: "d75b3a66374942c5b3c019c698abc2067e151558acbd412345",
	}
```

All methods of this library return two objects: `*APIResponse` and `error`. 

You can see the `APIResponse` format in the [types.go](types.go)

**How to send a message:**

Link to example: [sendMessage/main.go](examples/sendMessage/main.go)

```go
response, _ := GreenAPI.Sending().SendMessage(
		"10000000",
		"Hello",
	)
```

**How to create a group:**

Link to example: [createGroup/main.go](examples/createGroup/main.go)

```go
response, _ := GreenAPI.Groups().CreateGroup(
		"Group Title",
		[]string{
			"10000000",
			"10000001",
		},
	)
```

**How to send file by upload:**

Link to example: [sendFileByUpload/main.go](examples/sendFileByUpload/main.go)

```go
response, _ := GreenAPI.Sending().SendFileByUpload(
		"10000000",
		"C:/Users/user/Desktop/Pictures/image.png",
		"image.png",
	)
```

**How to send a file by URL:**

Link to example: [sendFileByUrl/main.go](examples/sendFileByUrl/main.go)

```go
response, _ := GreenAPI.Sending().SendFileByUrl(
		"10000000",
		"urlFile",
		"fileName",
		greenapi.OptionalCaptionSendUrl("Caption"),
	)
```

**How to receive an incoming notification:**

Link to example: [receiveNotification/main.go](examples/receiveNotification/main.go)

```go
response, _ := GreenAPI.Receiving().ReceiveNotification(
		greenapi.OptionalReceiveTimeout(5),
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
		greenapi.OptionalName("Created by GO SDK"),
		greenapi.OptionalWebhookUrl("https://webhook.url"),
		greenapi.OptionalWebhookUrlToken("auth_token"),
		greenapi.OptionalDelaySendMessages(5000),
		greenapi.OptionalMarkIncomingMessagesRead(true),
		greenapi.OptionalMarkIncomingMessagesReadOnReply(true),
		greenapi.OptionalOutgoingWebhook(true),
		greenapi.OptionalOutgoingMessageWebhook(true),
		greenapi.OptionalOutgoingAPIMessageWebhook(true),
		greenapi.OptionalStateWebhook(true),
		greenapi.OptionalIncomingWebhook(true),

	)
```

**How to delete an instance:**

Link to the example: [partnerMethods/deleteInstanceAccount/main.go](examples/partnerMethods/deleteInstanceAccount/main.go)

```go
response, _ := Partner.Partner().DeleteInstanceAccount(3100000000)
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
	)
```

In this example, only `DelaySendMessages`, `OutgoingWebhook` and `IncomingWebhook` settings will be changed, other settings are commented so they will not be passed. However, you can uncomment any setting that you prefer. **The settings that were not used will not be affected**

One more example of using optional parameters, this time let's use `sendMessage` method:

```go
response, _ := GreenAPI.Sending().SendMessage(
		"10000000",
		"Hello",
		greenapi.OptionalQuotedMessageId("BAE59673E71FC5DB"), // quotes specified message
	)
```

## List of examples

| Description                                   | Link to example                                               |
|-----------------------------------------------|---------------------------------------------------------------|
| How to send a message                         | [sendMessage/main.go](examples/sendMessage/main.go)           |
| How to send a file by uploading from the disk | [sendFileByUpload/main.go](examples/sendFileByUpload/main.go) |
| How to send a file by URL | [sendFileByUrl/main.go](examples/sendFileByUrl/main.go) |
| How to upload a file to an external drive                     | [uploadFile/main.go](examples/uploadFile/main.go)       |
| How to check if there is a MAX account on the phone number         | [CheckAccount/main.go](examples/CheckAccount/main.go)                   |
| How to set instance settings             | [setSettings/main.go](examples/setSettings/main.go)                 |
| How to create a group             | [createGroup/main.go](examples/createGroup/main.go)                 |
| How to receive an incoming notification | [receiveNotification/main.go](examples/receiveNotification/main.go) |
| How to get all instances of the account             | [partnerMethods/getInstances/main.go](examples/partnerMethods/getInstances/main.go)                 |
| How to create an instance             | [partnerMethods/createInstance/main.go](examples/partnerMethods/createInstance/main.go)                 |
| How to delete an instance            | [partnerMethods/deleteInstanceAccount/main.go](examples/partnerMethods/deleteInstanceAccount/main.go)                 |

## List of all library methods

| API method                        | Description                                                                                                               | Documentation link                                                                                          |
|-----------------------------------|---------------------------------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------|
| `Account().GetSettings`           | The method is designed to get the current settings of the account                                                         | [GetSettings](https://green-api.com/v3/docs/api/account/GetSettings/)                                       |
| `Account().GetAccountSettings`         | The method is designed to get information about the MAX account                                                      | [GetSettings](https://green-api.com/v3/docs/api/account/GetAccountSettings/)                                     |
| `Account().SetSettings`           | The method is designed to set the account settings                                                                        | [SetSettings](https://green-api.com/v3/docs/api/account/SetSettings/)                                          |
| `Account().GetStateInstance`      | The method is designed to get the state of the account                                                                    | [GetStateInstance](https://green-api.com/v3/docs/api/account/GetStateInstance/)                             |
| `Account().Reboot`                | The method is designed to restart the account                                                                             | [Reboot](https://green-api.com/v3/docs/api/account/Reboot/)                                                 |
| `Account().Logout`                | The method is designed to unlogin the account                                                                             | [Logout](https://green-api.com/v3/docs/api/account/Logout/)                                                 |
| `Account().SetProfilePicture`     | The method is designed to set the avatar of the account                                                                   | [SetProfilePicture](https://green-api.com/v3/docs/api/account/SetProfilePicture/)                           |
| `Account().StartAuthorization`  | The method is designed to complete the instance authorization process. Use the verification code received from the SMS when calling the method | [StartAuthorization](https://green-api.com/v3/docs/api/account/StartAuthorization/)                     |                                  |
| `Account().SendAuthorizationCode`  | Метод предназначен для завершения процесса авторизации инстанса. Используйте код проверки полученный из SMS при вызове метода  | [SendAuthorizationCode](https://green-api.com/v3/docs/api/account/SendAuthorizationCode/)                     |                                  |
| `Groups().CreateGroup`            | The method is designed to create a group chat                                                                             | [CreateGroup](https://green-api.com/v3/docs/api/groups/CreateGroup/)                                        |
| `Groups().UpdateGroupName`        | The method changes the name of the group chat                                                                             | [UpdateGroupName](https://green-api.com/v3/docs/api/groups/UpdateGroupName/)                                |
| `Groups().GetGroupData`           | The method gets group chat data                                                                                           | [GetGroupData](https://green-api.com/v3/docs/api/groups/GetGroupData/)                                      |
| `Groups().AddGroupParticipant`    | The method adds a participant to the group chat                                                                           | [AddGroupParticipant](https://green-api.com/v3/docs/api/groups/AddGroupParticipant/)                        |
| `Groups().RemoveGroupParticipant` | The method removes the participant from the group chat                                                                    | [RemoveGroupParticipant](https://green-api.com/v3/docs/api/groups/RemoveGroupParticipant/)                  |
| `Groups().SetGroupAdmin`          | The method designates a member of a group chat as an administrator                                                        | [SetGroupAdmin](https://green-api.com/v3/docs/api/groups/SetGroupAdmin/)                                    |
| `Groups().RemoveAdmin`            | The method deprives the participant of group chat administration rights                                                   | [RemoveAdmin](https://green-api.com/v3/docs/api/groups/RemoveAdmin/)                                        |
| `Groups().SetGroupPicture`        | The method sets the avatar of the group                                                                                   | [SetGroupPicture](https://green-api.com/v3/docs/api/groups/SetGroupPicture/)                                |
| `Groups().LeaveGroup`             | The method logs the user of the current account out of the group chat                                                     | [LeaveGroup](https://green-api.com/v3/docs/api/groups/LeaveGroup/)                                          |
| `Journals().GetChatHistory`       | The method returns the chat message history                                                                               | [GetChatHistory](https://green-api.com/v3/docs/api/journals/GetChatHistory/)                                |
| `Journals().GetMessage`           | The method returns a chat message                                                                                         | [GetMessage](https://green-api.com/v3/docs/api/journals/GetMessage/)                                        |
| `Journals().LastIncomingMessages` | The method returns the most recent incoming messages of the account                                                       | [LastIncomingMessages](https://green-api.com/v3/docs/api/journals/LastIncomingMessages/)                    |
| `Journals().LastOutgoingMessages` | The method returns the last sent messages of the account                                                                  | [LastOutgoingMessages](https://green-api.com/v3/docs/api/journals/LastOutgoingMessages/)                    |
| `Queues().ShowMessagesQueue`      | The method is designed to get the list of messages that are in the queue to be sent                                       | [ShowMessagesQueue](https://green-api.com/v3/docs/api/queues/ShowMessagesQueue/)                            |
| `Queues().ClearMessagesQueue`     | The method is designed to clear the queue of messages to be sent                                                          | [ClearMessagesQueue](https://green-api.com/v3/docs/api/queues/ClearMessagesQueue/)                          |
| `ReadMark().ReadChat`             | The method is designed to mark chat messages as read                                                                      | [ReadChat](https://green-api.com/v3/docs/api/marks/ReadChat/)                                               |
| `Receiving().ReceiveNotification` | The method is designed to receive a single incoming notification from the notification queue                              | [ReceiveNotification](https://green-api.com/v3/docs/api/receiving/technology-http-api/ReceiveNotification/) |
| `Receiving().DeleteNotification`  | The method is designed to remove an incoming notification from the notification queue                                     | [DeleteNotification](https://green-api.com/v3/docs/api/receiving/technology-http-api/DeleteNotification/)   |
| `Receiving().DownloadFile`        | The method is for downloading received and sent files                                                                     | [DownloadFile](https://green-api.com/v3/docs/api/receiving/files/DownloadFile/)                             |
| `Sending().SendMessage`           | The method is designed to send a text message to a personal or group chat                                                 | [SendMessage](https://green-api.com/v3/docs/api/sending/SendMessage/)                                       |
| `Sending().SendFileByUpload`      | The method is designed to send a file loaded through a form (form-data)                                                   | [SendFileByUpload](https://green-api.com/v3/docs/api/sending/SendFileByUpload/)                             |
| `Sending().SendFileByUrl`         | The method is designed to send a file downloaded via a link                                                               | [SendFileByUrl](https://green-api.com/v3/docs/api/sending/SendFileByUrl/)                                   |
| `Sending().UploadFile`            | The method allows you to upload a file from the local file system, which can later be sent using the SendFileByUrl method | [UploadFile](https://green-api.com/v3/docs/api/sending/UploadFile/)                                         |
| `Service().CheckAccount`         | The method checks if there is a MAX account on the phone number                                                      | [CheckAccount](https://green-api.com/v3/docs/api/service/CheckAccount/)                                   |
| `Service().GetAvatar`             | The method returns the avatar of the correspondent or group chat                                                          | [GetAvatar](https://green-api.com/v3/docs/api/service/GetAvatar/)                                           |
| `Service().GetContacts`           | The method is designed to get a list of contacts of the current account                                                   | [GetContacts](https://green-api.com/v3/docs/api/service/GetContacts/)                                       |
| `Service().GetContactInfo`        | The method is designed to obtain information about the contact                                                            | [GetContactInfo](https://green-api.com/v3/docs/api/service/GetContactInfo/)                                 |
| `Partner().GetInstances`   | The method is for getting all the account instances created by the partner.                                           | [GetInstances](https://green-api.com/v3/docs/partners/getInstances/)                       |
| `Partner().CreateInstance`   | The method is for creating an instance.                                           | [CreateInstance](https://green-api.com/v3/docs/partners/createInstance/)                       |
| `Partner().DeleteInstanceAccount`   | The method is for deleting an instance.                                           | [DeleteInstanceAccount](https://green-api.com/v3/docs/partners/deleteInstanceAccount/)                       |
