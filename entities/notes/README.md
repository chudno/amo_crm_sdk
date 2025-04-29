# Модуль Примечания

Модуль `notes` предоставляет функциональность для работы с примечаниями в amoCRM.

## Содержание

- [Основные функции](#основные-функции)
- [Создание примечания](#создание-примечания)
- [Получение примечания](#получение-примечания)
- [Получение списка примечаний](#получение-списка-примечаний)
- [Обновление примечания](#обновление-примечания)
- [Типы примечаний](#типы-примечаний)
- [Связь с другими сущностями](#связь-с-другими-сущностями)

## Основные функции

| Функция | Описание |
|---------|----------|
| `CreateNote` | Создание нового примечания |
| `GetNote` | Получение примечания по ID |
| `GetNotes` | Получение списка примечаний с фильтрацией |
| `UpdateNote` | Обновление существующего примечания |
| `DeleteNote` | Удаление примечания |

## Создание примечания

```go
import (
    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/notes"
)

// Инициализация клиента
apiClient := client.NewClient("https://your-domain.amocrm.ru", "your_access_token")

// Создание текстового примечания для контакта
newNote := &notes.Note{
    EntityID: 12345, // ID связанной сущности (например, контакта)
    EntityType: notes.EntityTypeContact, // Тип связанной сущности
    NoteType: notes.TypeCommon, // Тип примечания - обычное примечание
    Params: notes.Params{
        Text: "Клиент заинтересован в нашем предложении",
    },
}

// Сохранение примечания
createdNote, err := notes.CreateNote(apiClient, newNote)
if err != nil {
    // Обработка ошибки
}

// Создание примечания о входящем звонке
callNote := &notes.Note{
    EntityID: 12345, // ID сущности
    EntityType: notes.EntityTypeLead, // Тип сущности - сделка
    NoteType: notes.TypeIncomingCall, // Тип примечания - входящий звонок
    Params: notes.Params{
        Text: "Клиент интересовался условиями поставки",
        Phone: "+79001234567", // Телефон клиента
        CallDate: time.Now().Unix(), // Время звонка (Unix timestamp)
        CallDuration: 300, // Длительность звонка в секундах
    },
}

// Сохранение примечания о звонке
createdCallNote, err := notes.CreateNote(apiClient, callNote)
```

## Получение примечания

```go
// Получение примечания по ID
noteID := 12345
note, err := notes.GetNote(apiClient, noteID)
if err != nil {
    // Обработка ошибки
}
```

## Получение списка примечаний

```go
// Получение всех примечаний для контакта
contactID := 12345
filter := map[string]string{
    "filter[entity_id]": fmt.Sprintf("%d", contactID),
    "filter[entity_type]": notes.EntityTypeContact,
}
notesList, err := notes.GetNotes(apiClient, 1, 50, filter)
if err != nil {
    // Обработка ошибки
}

// Получение примечаний определенного типа
typeFilter := map[string]string{
    "filter[entity_id]": fmt.Sprintf("%d", contactID),
    "filter[entity_type]": notes.EntityTypeContact,
    "filter[note_type]": fmt.Sprintf("%d", notes.TypeIncomingCall), // Только примечания о входящих звонках
}
callNotes, err := notes.GetNotes(apiClient, 1, 50, typeFilter)
```

## Обновление примечания

```go
// Обновление существующего примечания
note.Params.Text = "Клиент очень заинтересован в нашем предложении"

updatedNote, err := notes.UpdateNote(apiClient, note)
if err != nil {
    // Обработка ошибки
}
```

## Типы примечаний

Модуль `notes` предоставляет константы для типов примечаний:

| Константа | Значение | Описание |
|-----------|----------|----------|
| `notes.TypeCommon` | 4 | Обычное примечание |
| `notes.TypeIncomingCall` | 10 | Входящий звонок |
| `notes.TypeOutgoingCall` | 11 | Исходящий звонок |
| `notes.TypeService` | 25 | Сервисное примечание |
| `notes.TypeSms` | 13 | SMS-сообщение |

Также доступны константы для типов сущностей:

| Константа | Значение | Описание |
|-----------|----------|----------|
| `notes.EntityTypeContact` | "contacts" | Контакт |
| `notes.EntityTypeLead` | "leads" | Сделка |
| `notes.EntityTypeCompany` | "companies" | Компания |
| `notes.EntityTypeTask` | "tasks" | Задача |

## Связь с другими сущностями

Примечания в amoCRM всегда связаны с определенной сущностью (контактом, сделкой, компанией и т.д.). При создании примечания необходимо указать тип сущности и её ID:

```go
// Примечание для контакта
note := &notes.Note{
    EntityID: 12345,
    EntityType: notes.EntityTypeContact,
    NoteType: notes.TypeCommon,
    Params: notes.Params{
        Text: "Встреча запланирована на следующую неделю",
    },
}

// Примечание для сделки
note := &notes.Note{
    EntityID: 67890,
    EntityType: notes.EntityTypeLead,
    NoteType: notes.TypeCommon,
    Params: notes.Params{
        Text: "Клиент запросил дополнительную информацию",
    },
}
```
