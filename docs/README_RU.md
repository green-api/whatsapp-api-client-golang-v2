# whatsapp-api-client-golang-v2

`whatsapp-api-client-golang-v2` библиотека для интеграции с мессенджером WhatsApp через API сервиса [green-api.com](https://green-api.com). Чтобы воспользоваться библиотекой, нужно получить регистрационный токен и ID аккаунта в [личном кабинете](https://console.green-api.com/). Есть бесплатный тариф аккаунта разработчика.

Вы можете найти версию `v1` по ссылке - https://github.com/green-api/whatsapp-api-client-golang

## API

Документация к REST API находится по [ссылке](https://green-api.com/docs/api). Библиотека является оберткой к REST API, поэтому документация по ссылке выше применима и к самой библиотеке.

## Поддержка

[![Support](https://img.shields.io/badge/support@green--api.com-D14836?style=for-the-badge&logo=gmail&logoColor=white)](mailto:support@greenapi.com)
[![Support](https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white)](https://t.me/greenapi_support_eng_bot)
[![Support](https://img.shields.io/badge/WhatsApp-25D366?style=for-the-badge&logo=whatsapp&logoColor=white)](https://wa.me/77273122366)

## Руководства и новости

[![Guides](https://img.shields.io/badge/YouTube-%23FF0000.svg?style=for-the-badge&logo=YouTube&logoColor=white)](https://www.youtube.com/@greenapi-en)
[![News](https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white)](https://t.me/green_api)
[![News](https://img.shields.io/badge/WhatsApp-25D366?style=for-the-badge&logo=whatsapp&logoColor=white)](https://whatsapp.com/channel/0029VaLj6J4LNSa2B5Jx6s3h)

#### Авторизация

Чтобы отправить сообщение или выполнить другие методы GREEN API, аккаунт WhatsApp в приложении телефона должен быть в авторизованном состоянии. Для авторизации аккаунта перейдите в [личный кабинет](https://console.green-api.com/) и сканируйте QR-код с использованием приложения WhatsApp.

## Установка

**Убедитесь, что у вас установлена версия Go не ниже 1.20**
```shell
go version
```

**Создайте Go модуль, если он не создан:**

```shell
go mod init ModuleName
```

**Установите библиотеку:**

```shell
go get github.com/green-api/whatsapp-api-client-golang-v2
```

**Импорт:**

```go
import (
	greenapi "github.com/green-api/whatsapp-api-client-golang-v2"
)
```

## Использование и примеры

**Как инициализировать объект:**

```go
GreenAPI := greenapi.GreenAPI{
		APIURL:           "https://api.green-api.com",
		MediaURL:         "https://media.green-api.com",
		IDInstance:       "1101000001",
		APITokenInstance: "d75b3a66374942c5b3c019c698abc2067e151558acbd412345",
	}
```

Все методы библиотеки возвращают два объекта: `*APIResponse` и `error`. 

Вы можете посмотреть формат `APIResponse` в [types.go](types.go)

**Как отправить сообщение:**

Ссылка на пример: [sendMessage/main.go](/examples/sendMessage/main.go)

```go
response, _ := GreenAPI.Sending().SendMessage(
		"11001234567@c.us",
		"Hello",
	)
```

**Как создать группу:**

Ссылка на пример: [createGroup/main.go](/examples/createGroup/main.go)

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

**Как отправить файл с диска:**

Ссылка на пример: [sendFileByUpload/main.go](/examples/sendFileByUpload/main.go)

```go
response, _ := GreenAPI.Sending().SendFileByUpload(
		"11001234567@c.us",
		"C:/Users/user/Desktop/Pictures/image.png",
		"image.png",
	)
```

**Как отправить файл по ссылке:**

Ссылка на пример: [sendFileByUrl/main.go](/examples/sendFileByUrl/main.go)

```go
response, _ := GreenAPI.Sending().SendFileByUrl(
		"11001234567@c.us",
		"urlFile",
		"fileName",
		greenapi.OptionalCaptionSendUrl("Caption"),
	)
```

**Как отправить сообщение с опросом:**

Ссылка на пример: [sendPoll/main.go](/examples/sendPoll/main.go)

```go
response, _ := GreenAPI.Sending().SendPoll(
		"11001234567@c.us", 
		"Choose a color:", 
		[]string{"Red", "Green", "Blue"}, 
	)
```

**Как отправить текстовый статус:**

Ссылка на пример: [sendTextStatus/main.go](/examples/sendTextStatus/main.go)

```go
response, _ := GreenAPI.Statuses().SendTextStatus(
		"Text of the status", 
		greenapi.OptionalFont("SERIF"),
		greenapi.OptionalBackgroundColorText("#87CEEB"),
		//greenapi.OptionalParticipantsTextStatus([]string{"1234567890@c.us", "1234567890@c.us"}),
	)
```

**Как получить входящее уведомление:**

Ссылка на пример: [receiveNotification/main.go](/examples/receiveNotification/main.go)

```go
response, _ := GreenAPI.Receiving().ReceiveNotification(
		greenapi.OptionalReceiveTimeout(5),
	)
```

## Методы партнёра

**Чтобы использовать методы партнёра, вы должны инициализировать другой объект:**

```go
Partner := greenapi.GreenAPIPartner{
		PartnerToken: "gac.1234567891234567891234567891213456789",
		Email: "mail@email.com", // поле email не обязательно 
	}
```

**Теперь вы можете использовать методы партнёра так же, как и обычные методы, но через объект "Partner":**

**Как получить все инстансы на аккаунте:**

Ссылка на пример: [partnerMethods/getInstances/main.go](/examples/partnerMethods/getInstances/main.go)

```go
response, _ := Partner.Partner().GetInstances()
```

**Как создать инстанс:**

Ссылка на пример: [partnerMethods/createInstance/main.go](/examples/partnerMethods/createInstance/main.go)

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

**Как удалить инстанс:**

Ссылка на пример: [partnerMethods/deleteInstanceAccount/main.go](/examples/partnerMethods/deleteInstanceAccount/main.go)

```go
response, _ := Partner.Partner().DeleteInstanceAccount(1101000000)
```

## Необязательные параметры

**Обратите внимание, что методы могут иметь необязательные параметры, которые вы можете передавать. Необязательные параметры передаются в аргументы методов в виде функций и имеют следующий формат:**
```go
greenapi.Optional + name of parameter
```

**К примеру в методе `SetSettings` все параметры являются опциональными. Рассмотрим пример вызова данной функции::**

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

В этом примере только настройки `DelaySendMessages`, `OutgoingWebhook` и `IncomingWebhook` будут изменены, остальные параметры закомментированы, поэтому не будут использованы. Вы можете раскомментировать любой параметр который предпочитаете. **Неиспользованные параметры никак не затронут настройки инстанса**

Ещё один пример использования опциональных параметров, в этот раз рассмотрим метод `sendMessage`:

```go
response, _ := GreenAPI.Sending().SendMessage(
		"11001234567@c.us",
		"Hello",
		greenapi.OptionalLinkPreview(false), // выключает формирование превью ссылок
		greenapi.OptionalQuotedMessageId("BAE59673E71FC5DB"), // цитирует указанное сообщение
	)
```

## Список примеров

| Описание                                   | Ссылка на пример                                               |
|-----------------------------------------------|---------------------------------------------------------------|
| Как отправить сообщение                         | [sendMessage/main.go](/examples/sendMessage/main.go)           |
| Как отправить файл с диска | [sendFileByUpload/main.go](/examples/sendFileByUpload/main.go) |
| Как отправить файл по ссылке | [sendFileByUrl/main.go](/examples/sendFileByUrl/main.go) |
| Как выгрузить файл в облачное хранилище                     | [uploadFile/main.go](/examples/uploadFile/main.go)       |
| Как отправить опрос                         | [sendPoll/main.go](/examples/sendPoll/main.go)           |
| Как проверить номер телефона на наличие аккаунта WhatsApp         | [checkWhatsapp/main.go](/examples/checkWhatsapp/main.go)                   |
| Как установить настройки инстанса             | [setSettings/main.go](/examples/setSettings/main.go)                 |
| Как создать группу             | [createGroup/main.go](/examples/createGroup/main.go)                 |
| Как отправить текстовый статус             | [sendTextStatus/main.go](/examples/sendTextStatus/main.go)                 |
| Как получить входящее уведомление | [receiveNotification/main.go](/examples/receiveNotification/main.go) |
| Как получить все инстансы на аккаунте             | [partnerMethods/getInstances/main.go](/examples/partnerMethods/getInstances/main.go)                 |
| Как создать инстанс             | [partnerMethods/createInstance/main.go](/examples/partnerMethods/createInstance/main.go)                 |
| Как удалить инстанс            | [partnerMethods/deleteInstanceAccount/main.go](/examples/partnerMethods/deleteInstanceAccount/main.go)                 |

## Список всех методов библиотеки

| API метод                        | Описание                                                                                                               | Ссылка на документацию                                                                                          |
|-----------------------------------|---------------------------------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------|
| `Account().GetSettings`           | Метод предназначен для получения текущих настроек аккаунта                                                         | [GetSettings](https://green-api.com/docs/api/account/GetSettings/)                                       |
| `Account().GetWaSettings`         | Метод предназначен для получения информации о аккаунте WhatsApp                                                      | [GetSettings](https://green-api.com/docs/api/account/GetWaSettings/)                                     |
| `Account().SetSettings`           | Метод предназначен для установки настроек аккаунта                                                                        | [SetSettings](https://green-api.com/docs/api/account/SetSettings/)                                          |
| `Account().GetStateInstance`      | Метод предназначен для получения состояния аккаунта                                                                    | [GetStateInstance](https://green-api.com/docs/api/account/GetStateInstance/)                             |
| `Account().Reboot`                | Метод предназначен для перезапуска аккаунта                                                                             | [Reboot](https://green-api.com/docs/api/account/Reboot/)                                                 |
| `Account().Logout`                | Метод предназначен для деавторизации аккаунта                                                                             | [Logout](https://green-api.com/docs/api/account/Logout/)                                                 |
| `Account().QR`                    | Метод предназначен для получения QR-кода                                                                                   | [QR](https://green-api.com/docs/api/account/QR/)                                                         |
| `Account().SetProfilePicture`     | Метод предназначен для установки аватара аккаунта                                                                   | [SetProfilePicture](https://green-api.com/docs/api/account/SetProfilePicture/)                           |
| `Account().GetAuthorizationCode`  | Метод предназначен для авторизации инстанса по номеру телефона                                                           | [GetAuthorizationCode](https://green-api.com/docs/api/account/GetAuthorizationCode/)                     |                                  |
| `Groups().CreateGroup`            | Метод предназначен для создания группового чата                                                                             | [CreateGroup](https://green-api.com/docs/api/groups/CreateGroup/)                                        |
| `Groups().UpdateGroupName`        | Метод изменяет наименование группового чата                                                                             | [UpdateGroupName](https://green-api.com/docs/api/groups/UpdateGroupName/)                                |
| `Groups().GetGroupData`           | Метод получает данные группового чата                                                                                           | [GetGroupData](https://green-api.com/docs/api/groups/GetGroupData/)                                      |
| `Groups().AddGroupParticipant`    | Метод добавляет участника в групповой чат                                                                           | [AddGroupParticipant](https://green-api.com/docs/api/groups/AddGroupParticipant/)                        |
| `Groups().RemoveGroupParticipant` | Метод удаляет участника из группового чата                                                                    | [RemoveGroupParticipant](https://green-api.com/docs/api/groups/RemoveGroupParticipant/)                  |
| `Groups().SetGroupAdmin`          | Метод назначает участника группового чата администратором                                                        | [SetGroupAdmin](https://green-api.com/docs/api/groups/SetGroupAdmin/)                                    |
| `Groups().RemoveAdmin`            | Метод лишает участника прав администрирования группового чата                                                   | [RemoveAdmin](https://green-api.com/docs/api/groups/RemoveAdmin/)                                        |
| `Groups().SetGroupPicture`        | Метод устанавливает аватар группы                                                                                   | [SetGroupPicture](https://green-api.com/docs/api/groups/SetGroupPicture/)                                |
| `Groups().LeaveGroup`             | 	Метод производит выход пользователя текущего аккаунта из группового чата                                                     | [LeaveGroup](https://green-api.com/docs/api/groups/LeaveGroup/)                                          |
| `Journals().GetChatHistory`       | Метод возвращает историю сообщений чата                                                                               | [GetChatHistory](https://green-api.com/docs/api/journals/GetChatHistory/)                                |
| `Journals().GetMessage`           | Метод возвращает сообщение чата                                                                                         | [GetMessage](https://green-api.com/docs/api/journals/GetMessage/)                                        |
| `Journals().LastIncomingMessages` | Метод возвращает крайние входящие сообщения аккаунта                                                       | [LastIncomingMessages](https://green-api.com/docs/api/journals/LastIncomingMessages/)                    |
| `Journals().LastOutgoingMessages` | Метод возвращает крайние отправленные сообщения аккаунта                                                                  | [LastOutgoingMessages](https://green-api.com/docs/api/journals/LastOutgoingMessages/)                    |
| `Queues().ShowMessagesQueue`      | Метод предназначен для получения списка сообщений, находящихся в очереди на отправку                                       | [ShowMessagesQueue](https://green-api.com/docs/api/queues/ShowMessagesQueue/)                            |
| `Queues().ClearMessagesQueue`     | Метод предназначен для очистки очереди сообщений на отправку                                                          | [ClearMessagesQueue](https://green-api.com/docs/api/queues/ClearMessagesQueue/)                          |
| `ReadMark().ReadChat`             | Метод предназначен для отметки сообщений в чате прочитанными                                                                      | [ReadChat](https://green-api.com/docs/api/marks/ReadChat/)                                               |
| `Receiving().ReceiveNotification` | Метод предназначен для получения одного входящего уведомления из очереди уведомлений                              | [ReceiveNotification](https://green-api.com/docs/api/receiving/technology-http-api/ReceiveNotification/) |
| `Receiving().DeleteNotification`  | Метод предназначен для удаления входящего уведомления из очереди уведомлений                                     | [DeleteNotification](https://green-api.com/docs/api/receiving/technology-http-api/DeleteNotification/)   |
| `Receiving().DownloadFile`        | 	Метод предназначен для скачивания принятых и отправленных файлов                                                                     | [DownloadFile](https://green-api.com/docs/api/receiving/files/DownloadFile/)                             |
| `Sending().SendMessage`           | Метод предназначен для отправки текстового сообщения в личный или групповой чат                                                 | [SendMessage](https://green-api.com/docs/api/sending/SendMessage/)                                       |
| `Sending().SendFileByUpload`      | Метод предназначен для отправки файла, загружаемого через форму (form-data)                                                   | [SendFileByUpload](https://green-api.com/docs/api/sending/SendFileByUpload/)                             |
| `Sending().SendFileByUrl`         | Метод предназначен для отправки файла, загружаемого по ссылке                                                               | [SendFileByUrl](https://green-api.com/docs/api/sending/SendFileByUrl/)                                   |
| `Sending().UploadFile`            | Метод предназначен для загрузки файла в облачное хранилище, который можно отправить методом sendFileByUrl | [UploadFile](https://green-api.com/docs/api/sending/UploadFile/)                                         |
| `Sending().SendLocation`          | Метод предназначен для отправки сообщения геолокации                                                                      | [SendLocation](https://green-api.com/docs/api/sending/SendLocation/)                                     |
| `Sending().SendContact`           | Метод предназначен для отправки сообщения с контактом                                                                        | [SendContact](https://green-api.com/docs/api/sending/SendContact/)                                       |
| `Sending().ForwardMessages`       | Метод предназначен для пересылки сообщений в личный или групповой чат                                                | [ForwardMessages](https://green-api.com/docs/api/sending/ForwardMessages/)                               |
| `Sending().SendPoll`              | Метод предназначен для отправки сообщения с опросом в личный или групповой чат                                        | [SendPoll](https://green-api.com/docs/api/sending/SendPoll/)                                             |
| `Service().CheckWhatsapp`         | Метод проверяет наличие аккаунта WhatsApp на номере телефона                                                      | [CheckWhatsapp](https://green-api.com/docs/api/service/CheckWhatsapp/)                                   |
| `Service().GetAvatar`             | Метод возвращает аватар корреспондента или группового чата	                                                          | [GetAvatar](https://green-api.com/docs/api/service/GetAvatar/)                                           |
| `Service().GetContacts`           | Метод предназначен для получения списка контактов текущего аккаунта                                                   | [GetContacts](https://green-api.com/docs/api/service/GetContacts/)                                       |
| `Service().GetContactInfo`        | Метод предназначен для получения информации о контакте                                                            | [GetContactInfo](https://green-api.com/docs/api/service/GetContactInfo/)                                 |
| `Service().DeleteMessage`         | Метод удаляет сообщение из чата                                                                                  | [DeleteMessage](https://green-api.com/docs/api/service/deleteMessage/)                                   |
| `Service().ArchiveChat`           | Метод архивирует чат                                                                                              | [ArchiveChat](https://green-api.com/docs/api/service/archiveChat/)                                       |
| `Service().UnarchiveChat`         | Метод разархивирует чат                                                                                            | [UnarchiveChat](https://green-api.com/docs/api/service/unarchiveChat/)                                   |
| `Service().SetDisappearingChat`   | Метод предназначен для изменения настроек исчезающих сообщений в чатах                                           | [SetDisappearingChat](https://green-api.com/docs/api/service/SetDisappearingChat/)                       |
| `Partner().GetInstances`   | Метод предназначен для получения всех инстансов аккаунтов созданных партнёром.                                           | [GetInstances](https://green-api.com/docs/partners/getInstances/)                       |
| `Partner().CreateInstance`   | Метод предназначен для создания инстанса от имени партнёра.                                           | [CreateInstance](https://green-api.com/docs/partners/createInstance/)                       |
| `Partner().DeleteInstanceAccount`   | Метод предназначен для удаления инстанса аккаунта партнёра.                                           | [DeleteInstanceAccount](https://green-api.com/docs/partners/deleteInstanceAccount/)                       |
| `Statuses().SendTextStatus`             | Метод предназначен для отправки текстового статуса                                                     | [SendTextStatus](https://green-api.com/docs/api/statuses/SendTextStatus/)                                          |
| `Statuses().SendVoiceStatus`             | Метод предназначен для отправки голосового статуса                                                     | [SendVoiceStatus](https://green-api.com/docs/api/statuses/SendVoiceStatus/)                                          |
| `Statuses().SendMediaStatus`             | Метод предназначен для отправки медиа-файлов                                                     | [SendMediaStatus](https://green-api.com/docs/api/statuses/SendMediaStatus/)                                          |      
| `Statuses().GetOutgoingStatuses`             | Метод возвращает крайние отправленные статусы аккаунта                                                     | [GetOutgoingStatuses](https://green-api.com/docs/api/statuses/GetOutgoingStatuses/)                                          |      
| `Statuses().GetIncomingStatuses`             | Метод возвращает крайние входящие статусы аккаунта                                                     | [GetIncomingStatuses](https://green-api.com/docs/api/statuses/GetIncomingStatuses/)                                          |      
| `Statuses().GetStatusStatistic`             | Метод возвращает массив получателей со статусами, отмеченных как отправлено/доставлено/прочитано, для данного статуса                                                     | [GetStatusStatistic](https://green-api.com/docs/api/statuses/GetStatusStatistic/)                                          |      
| `Statuses().DeleteStatus`             | Метод предназначен для удаления статуса                                                     | [DeleteStatus](https://green-api.com/docs/api/statuses/DeleteStatus/)                                          |    