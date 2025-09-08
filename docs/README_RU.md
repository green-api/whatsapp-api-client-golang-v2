# max-api-client-golang

`max-api-client-golang` библиотека для интеграции с мессенджером MAX через API сервиса [green-api.com](https://green-api.com). Чтобы воспользоваться библиотекой, нужно получить регистрационный токен и ID аккаунта в [личном кабинете](https://console.green-api.com/). Есть бесплатный тариф аккаунта разработчика.

## API

Документация к REST API находится по [ссылке](https://green-api.com/v3/docs/api). Библиотека является оберткой к REST API, поэтому документация по ссылке выше применима и к самой библиотеке.

## Поддержка

[![Support](https://img.shields.io/badge/support@green--api.com-D14836?style=for-the-badge&logo=gmail&logoColor=white)](mailto:support@greenapi.com)
[![Support](https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white)](https://t.me/greenapi_support_eng_bot)
[![Support](https://img.shields.io/badge/WhatsApp-25D366?style=for-the-badge&logo=whatsapp&logoColor=white)](https://wa.me/77273122366)

## Руководства и новости

[![Guides](https://img.shields.io/badge/YouTube-%23FF0000.svg?style=for-the-badge&logo=YouTube&logoColor=white)](https://www.youtube.com/@greenapi-en)
[![News](https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white)](https://t.me/green_api)
[![News](https://img.shields.io/badge/WhatsApp-25D366?style=for-the-badge&logo=whatsapp&logoColor=white)](https://whatsapp.com/channel/0029VaLj6J4LNSa2B5Jx6s3h)

#### Авторизация

Чтобы отправить сообщение или выполнить другие методы GREEN API, аккаунт MAX в приложении телефона должен быть в авторизованном состоянии. Для авторизации аккаунта перейдите в [личный кабинет](https://console.green-api.com/) и сканируйте QR-код с использованием приложения MAX.

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
go get github.com/green-api/max-api-client-golang
```

**Импорт:**

```go
import (
	greenapi "github.com/green-api/max-api-client-golang"
)
```

## Использование и примеры

**Как инициализировать объект:**

```go
GreenAPI := greenapi.GreenAPI{
		APIURL:           "https://api.green-api.com/v3",
		MediaURL:         "https://api.green-api.com/v3",
		IDInstance:       "3100000001",
		APITokenInstance: "d75b3a66374942c5b3c019c698abc2067e151558acbd412345",
	}
```

Все методы библиотеки возвращают два объекта: `*APIResponse` и `error`. 

Вы можете посмотреть формат `APIResponse` в [types.go](types.go)

**Как отправить сообщение:**

Ссылка на пример: [sendMessage/main.go](/examples/sendMessage/main.go)

```go
response, _ := GreenAPI.Sending().SendMessage(
		"10000000",
		"Hello",
	)
```

**Как создать группу:**

Ссылка на пример: [createGroup/main.go](/examples/createGroup/main.go)

```go
response, _ := GreenAPI.Groups().CreateGroup(
		"Group Title",
		[]string{
			"10000000",
			"10000001",
		},
	)
```

**Как отправить файл с диска:**

Ссылка на пример: [sendFileByUpload/main.go](/examples/sendFileByUpload/main.go)

```go
response, _ := GreenAPI.Sending().SendFileByUpload(
		"10000000",
		"C:/Users/user/Desktop/Pictures/image.png",
		"image.png",
	)
```

**Как отправить файл по ссылке:**

Ссылка на пример: [sendFileByUrl/main.go](/examples/sendFileByUrl/main.go)

```go
response, _ := GreenAPI.Sending().SendFileByUrl(
		"10000000",
		"urlFile",
		"fileName",
		greenapi.OptionalCaptionSendUrl("Caption"),
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

**Как удалить инстанс:**

Ссылка на пример: [partnerMethods/deleteInstanceAccount/main.go](/examples/partnerMethods/deleteInstanceAccount/main.go)

```go
response, _ := Partner.Partner().DeleteInstanceAccount(3100000000)
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
	)
```

В этом примере только настройки `DelaySendMessages`, `OutgoingWebhook` и `IncomingWebhook` будут изменены, остальные параметры закомментированы, поэтому не будут использованы. Вы можете раскомментировать любой параметр который предпочитаете. **Неиспользованные параметры никак не затронут настройки инстанса**

Ещё один пример использования опциональных параметров, в этот раз рассмотрим метод `sendMessage`:

```go
response, _ := GreenAPI.Sending().SendMessage(
		"10000000",
		"Hello",
		greenapi.OptionalQuotedMessageId("2712345112345"), // цитирует указанное сообщение
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
| Как проверить номер телефона на наличие аккаунта MAX         | [CheckAccount/main.go](/examples/CheckAccount/main.go)                   |
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
| `Account().GetSettings`           | Метод предназначен для получения текущих настроек аккаунта                                                         | [GetSettings](https://green-api.com/v3/docs/api/account/GetSettings/)                                       |
| `Account().GetAccountSettings`         | Метод предназначен для получения информации о аккаунте MAX                                                      | [GetSettings](https://green-api.com/v3/docs/api/account/GetAccountSettings/)                                     |
| `Account().SetSettings`           | Метод предназначен для установки настроек аккаунта                                                                        | [SetSettings](https://green-api.com/v3/docs/api/account/SetSettings/)                                          |
| `Account().GetStateInstance`      | Метод предназначен для получения состояния аккаунта                                                                    | [GetStateInstance](https://green-api.com/v3/docs/api/account/GetStateInstance/)                             |
| `Account().Reboot`                | Метод предназначен для перезапуска аккаунта                                                                             | [Reboot](https://green-api.com/v3/docs/api/account/Reboot/)                                                 |
| `Account().Logout`                | Метод предназначен для деавторизации аккаунта                                                                             | [Logout](https://green-api.com/v3/docs/api/account/Logout/)                                                 |
| `Account().QR`                    | Метод предназначен для получения QR-кода                                                                                   | [QR](https://green-api.com/v3/docs/api/account/QR/)                                                         |
| `Account().SetProfilePicture`     | Метод предназначен для установки аватара аккаунта                                                                   | [SetProfilePicture](https://green-api.com/v3/docs/api/account/SetProfilePicture/)                           |
| `Account().StartAuthorization`  | Метод предназначен для авторизации инстанса. Процесс авторизации заключается в подключении к шлюзу GREEN-API существующего аккаунта мессенджера MAX | [StartAuthorization](https://green-api.com/v3/docs/api/account/StartAuthorization/)                     |                                  |
| `Account().SendAuthorizationCode`  | Метод предназначен для завершения процесса авторизации инстанса. Используйте код проверки полученный из SMS при вызове метода  | [SendAuthorizationCode](https://green-api.com/v3/docs/api/account/SendAuthorizationCode/)                     |                                  |
| `Groups().CreateGroup`            | Метод предназначен для создания группового чата                                                                             | [CreateGroup](https://green-api.com/v3/docs/api/groups/CreateGroup/)                                        |
| `Groups().UpdateGroupName`        | Метод изменяет наименование группового чата                                                                             | [UpdateGroupName](https://green-api.com/v3/docs/api/groups/UpdateGroupName/)                                |
| `Groups().GetGroupData`           | Метод получает данные группового чата                                                                                           | [GetGroupData](https://green-api.com/v3/docs/api/groups/GetGroupData/)                                      |
| `Groups().AddGroupParticipant`    | Метод добавляет участника в групповой чат                                                                           | [AddGroupParticipant](https://green-api.com/v3/docs/api/groups/AddGroupParticipant/)                        |
| `Groups().RemoveGroupParticipant` | Метод удаляет участника из группового чата                                                                    | [RemoveGroupParticipant](https://green-api.com/v3/docs/api/groups/RemoveGroupParticipant/)                  |
| `Groups().SetGroupAdmin`          | Метод назначает участника группового чата администратором                                                        | [SetGroupAdmin](https://green-api.com/v3/docs/api/groups/SetGroupAdmin/)                                    |
| `Groups().RemoveAdmin`            | Метод лишает участника прав администрирования группового чата                                                   | [RemoveAdmin](https://green-api.com/v3/docs/api/groups/RemoveAdmin/)                                        |
| `Groups().SetGroupPicture`        | Метод устанавливает аватар группы                                                                                   | [SetGroupPicture](https://green-api.com/v3/docs/api/groups/SetGroupPicture/)                                |
| `Groups().LeaveGroup`             | 	Метод производит выход пользователя текущего аккаунта из группового чата                                                     | [LeaveGroup](https://green-api.com/v3/docs/api/groups/LeaveGroup/)                                          |
| `Journals().GetChatHistory`       | Метод возвращает историю сообщений чата                                                                               | [GetChatHistory](https://green-api.com/v3/docs/api/journals/GetChatHistory/)                                |
| `Journals().GetMessage`           | Метод возвращает сообщение чата                                                                                         | [GetMessage](https://green-api.com/v3/docs/api/journals/GetMessage/)                                        |
| `Journals().LastIncomingMessages` | Метод возвращает крайние входящие сообщения аккаунта                                                       | [LastIncomingMessages](https://green-api.com/v3/docs/api/journals/LastIncomingMessages/)                    |
| `Journals().LastOutgoingMessages` | Метод возвращает крайние отправленные сообщения аккаунта                                                                  | [LastOutgoingMessages](https://green-api.com/v3/docs/api/journals/LastOutgoingMessages/)                    |
| `Queues().ShowMessagesQueue`      | Метод предназначен для получения списка сообщений, находящихся в очереди на отправку                                       | [ShowMessagesQueue](https://green-api.com/v3/docs/api/queues/ShowMessagesQueue/)                            |
| `Queues().ClearMessagesQueue`     | Метод предназначен для очистки очереди сообщений на отправку                                                          | [ClearMessagesQueue](https://green-api.com/v3/docs/api/queues/ClearMessagesQueue/)                          |
| `ReadMark().ReadChat`             | Метод предназначен для отметки сообщений в чате прочитанными                                                                      | [ReadChat](https://green-api.com/v3/docs/api/marks/ReadChat/)                                               |
| `Receiving().ReceiveNotification` | Метод предназначен для получения одного входящего уведомления из очереди уведомлений                              | [ReceiveNotification](https://green-api.com/v3/docs/api/receiving/technology-http-api/ReceiveNotification/) |
| `Receiving().DeleteNotification`  | Метод предназначен для удаления входящего уведомления из очереди уведомлений                                     | [DeleteNotification](https://green-api.com/v3/docs/api/receiving/technology-http-api/DeleteNotification/)   |
| `Receiving().DownloadFile`        | 	Метод предназначен для скачивания принятых и отправленных файлов                                                                     | [DownloadFile](https://green-api.com/v3/docs/api/receiving/files/DownloadFile/)                             |
| `Sending().SendMessage`           | Метод предназначен для отправки текстового сообщения в личный или групповой чат                                                 | [SendMessage](https://green-api.com/v3/docs/api/sending/SendMessage/)                                       |
| `Sending().SendFileByUpload`      | Метод предназначен для отправки файла, загружаемого через форму (form-data)                                                   | [SendFileByUpload](https://green-api.com/v3/docs/api/sending/SendFileByUpload/)                             |
| `Sending().SendFileByUrl`         | Метод предназначен для отправки файла, загружаемого по ссылке                                                               | [SendFileByUrl](https://green-api.com/v3/docs/api/sending/SendFileByUrl/)                                   |
| `Sending().UploadFile`            | Метод предназначен для загрузки файла в облачное хранилище, который можно отправить методом sendFileByUrl | [UploadFile](https://green-api.com/v3/docs/api/sending/UploadFile/)                                         |
| `Service().CheckAccount`         | Метод проверяет наличие аккаунта MAX на номере телефона                                                      | [CheckAccount](https://green-api.com/v3/docs/api/service/CheckAccount/)                                   |
| `Service().GetAvatar`             | Метод возвращает аватар корреспондента или группового чата	                                                          | [GetAvatar](https://green-api.com/v3/docs/api/service/GetAvatar/)                                           |
| `Service().GetContacts`           | Метод предназначен для получения списка контактов текущего аккаунта                                                   | [GetContacts](https://green-api.com/v3/docs/api/service/GetContacts/)                                       |
| `Service().GetContactInfo`        | Метод предназначен для получения информации о контакте                                                            | [GetContactInfo](https://green-api.com/v3/docs/api/service/GetContactInfo/)                                 |
| `Partner().GetInstances`   | Метод предназначен для получения всех инстансов аккаунтов созданных партнёром.                                           | [GetInstances](https://green-api.com/v3/docs/partners/getInstances/)                       |
| `Partner().CreateInstance`   | Метод предназначен для создания инстанса от имени партнёра.                                           | [CreateInstance](https://green-api.com/v3/docs/partners/createInstance/)                       |
| `Partner().DeleteInstanceAccount`   | Метод предназначен для удаления инстанса аккаунта партнёра.                                           | [DeleteInstanceAccount](https://green-api.com/v3/docs/partners/deleteInstanceAccount/)                       |