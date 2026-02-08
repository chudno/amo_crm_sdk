# Модуль Примечания

Модуль `notes` предоставляет функциональность для работы с примечаниями в amoCRM.

## Содержание

- [Основные функции](#основные-функции)
- [Создание примечания](#создание-примечания)
- [Получение примечания](#получение-примечания)
- [Получение списка примечаний](#получение-списка-примечаний)
- [Обновление примечания](#обновление-примечания)
- [Типы примечаний](#типы-примечаний)

## Основные функции

| Функция | Описание |
|---------|----------|
| `Create` | Создание нового примечания |
| `Get` | Получение примечания по ID |
| `List` | Получение списка примечаний для сущности |
| `Update` | Обновление существующего примечания |

## Создание примечания

```go
import (
    "context"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/notes"
)

apiClient := client.NewClient("https://your-domain.amocrm.ru", "your_access_token")
ctx := context.Background()

newNote := &notes.Note{
    NoteType: notes.TypeCommon,
    Params: notes.NoteParams{
        Text: "Клиент заинтересован в нашем предложении",
    },
}

createdNote, err := notes.Create(ctx, apiClient, "contacts", 12345, newNote)
if err != nil {
    // Обработка ошибки
}

callNote := &notes.Note{
    NoteType: notes.TypeCallIn,
    Params: notes.NoteParams{
        Text:        "Клиент интересовался условиями поставки",
        PhoneNumber: "+79001234567",
    },
}

createdCallNote, err := notes.Create(ctx, apiClient, "leads", 12345, callNote)
```

## Получение примечания

```go
note, err := notes.Get(ctx, apiClient, "contacts", 12345, 67890)
if err != nil {
    // Обработка ошибки
}
```

## Получение списка примечаний

```go
notesList, err := notes.List(ctx, apiClient, "contacts", 12345, 50, 1)
if err != nil {
    // Обработка ошибки
}
```

## Обновление примечания

```go
note.Params.Text = "Клиент очень заинтересован в нашем предложении"

updatedNote, err := notes.Update(ctx, apiClient, "contacts", 12345, note)
if err != nil {
    // Обработка ошибки
}
```

## Типы примечаний

Модуль `notes` предоставляет константы для типов примечаний:

| Константа | Значение | Описание |
|-----------|----------|----------|
| `notes.TypeCommon` | "common" | Обычное примечание |
| `notes.TypeCallIn` | "call_in" | Входящий звонок |
| `notes.TypeCallOut` | "call_out" | Исходящий звонок |
| `notes.TypeServiceMessage` | "service_message" | Сервисное сообщение |
| `notes.TypeIncomingChatMessage` | "incoming_chat_message" | Входящее сообщение чата |
| `notes.TypeOutgoingChatMessage` | "outgoing_chat_message" | Исходящее сообщение чата |
| `notes.TypeSmsIn` | "sms_in" | Входящее SMS |
| `notes.TypeSmsOut` | "sms_out" | Исходящее SMS |
| `notes.TypeExtendedServiceMessage` | "extended_service_message" | Расширенное сервисное сообщение |
| `notes.TypeAttachment` | "attachment" | Вложение |
| `notes.TypeAmomailMessage` | "amomail_message" | Сообщение amoCRM почты |
| `notes.TypeGeolocation` | "geolocation" | Геолокация |

Тип сущности (`entityType`) передается строкой: `"leads"`, `"contacts"`, `"companies"`, `"customers"`.
