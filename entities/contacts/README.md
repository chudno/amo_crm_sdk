# Модуль Контакты

Модуль `contacts` предоставляет функциональность для работы с контактами в amoCRM.

## Содержание

- [Основные функции](#основные-функции)
- [Создание контакта](#создание-контакта)
- [Получение контакта](#получение-контакта)
- [Получение списка контактов](#получение-списка-контактов)
- [Обновление контакта](#обновление-контакта)
- [Пользовательские поля](#пользовательские-поля)
- [Связывание контактов с другими сущностями](#связывание-контактов-с-другими-сущностями)

## Основные функции

| Функция | Описание |
|---------|----------|
| `Create` | Создание нового контакта |
| `Get` | Получение контакта по ID |
| `List` | Получение списка контактов с фильтрацией |
| `Update` | Обновление существующего контакта |
| `Delete` | Удаление контакта |

## Создание контакта

```go
import (
    "context"
    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/contacts"
)

// Инициализация клиента
apiClient := client.NewClient("https://your-domain.amocrm.ru", "your_access_token")
ctx := context.Background()

// Создание нового контакта
newContact := &contacts.Contact{
    Name: "Иван Иванов",
    ResponsibleUserID: 12345, // ID ответственного менеджера
}

// Добавление номера телефона
newContact.CustomFields = append(newContact.CustomFields, contacts.Field{
    FieldID: 1234, // ID поля "Телефон"
    Values: []contacts.Value{
        {
            Value: "+79001234567",
            Enum: "WORK", // Тип телефона (рабочий)
        },
    },
})

// Добавление email
newContact.CustomFields = append(newContact.CustomFields, contacts.Field{
    FieldID: 5678, // ID поля "Email"
    Values: []contacts.Value{
        {
            Value: "ivan@example.com",
            Enum: "WORK", // Тип email (рабочий)
        },
    },
})

// Сохранение контакта
createdContact, err := contacts.Create(ctx, apiClient, newContact)
if err != nil {
    // Обработка ошибки
}
```

## Получение контакта

```go
// Получение контакта по ID
contactID := 12345
contact, err := contacts.Get(ctx, apiClient, contactID)
if err != nil {
    // Обработка ошибки
}
```

## Получение списка контактов

```go
// Получение первых 50 контактов
contactsList, err := contacts.List(ctx, apiClient, 1, 50)
if err != nil {
    // Обработка ошибки
}

// Получение контактов с фильтрацией
filter := map[string]string{
    "query": "Иван", // Поиск по имени
    "created_at": "1609459200", // Контакты, созданные после указанной даты (timestamp)
}
filteredContacts, err := contacts.List(ctx, apiClient, 1, 50, filter)
```

## Обновление контакта

```go
// Обновление существующего контакта
contact.Name = "Иван Петрович Иванов"

// Добавление нового номера телефона
contact.CustomFields = append(contact.CustomFields, contacts.Field{
    FieldID: 1234, // ID поля "Телефон"
    Values: []contacts.Value{
        {
            Value: "+79009876543",
            Enum: "PERSONAL", // Тип телефона (личный)
        },
    },
})

updatedContact, err := contacts.Update(ctx, apiClient, contact)
if err != nil {
    // Обработка ошибки
}
```

## Пользовательские поля

Для работы с пользовательскими полями контактов используйте структуры `Field` и `Value`:

```go
// Добавление пользовательского поля
contact.CustomFields = append(contact.CustomFields, contacts.Field{
    FieldID: 9876, // ID пользовательского поля
    Values: []contacts.Value{
        {
            Value: "Значение поля",
        },
    },
})
```

## Связывание контактов с другими сущностями

Для связывания контактов с другими сущностями используйте соответствующие методы из модулей leads и companies:

```go
// Связывание контакта со сделкой
import "github.com/chudno/amo_crm_sdk/entities/leads"

err := leads.LinkLeadWithContact(ctx, apiClient, leadID, contactID)

// Связывание контакта с компанией
import "github.com/chudno/amo_crm_sdk/entities/companies"
err := companies.LinkCompanyWithContact(ctx, apiClient, companyID, contactID)

// Связывание контакта с лидом
import "github.com/chudno/amo_crm_sdk/entities/leads"
err := leads.LinkLeadWithContact(ctx, apiClient, leadID, contactID)
```
