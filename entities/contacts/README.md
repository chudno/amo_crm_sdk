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
| `CreateContact` | Создание нового контакта |
| `GetContact` | Получение контакта по ID |
| `GetContacts` | Получение списка контактов с фильтрацией |
| `UpdateContact` | Обновление существующего контакта |
| `DeleteContact` | Удаление контакта |

## Создание контакта

```go
import (
    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/contacts"
)

// Инициализация клиента
apiClient := client.NewClient("https://your-domain.amocrm.ru", "your_access_token")

// Создание нового контакта
newContact := &contacts.Contact{
    Name: "Иван Иванов",
    ResponsibleUserID: 12345, // ID ответственного менеджера
}

// Добавление номера телефона
newContact.CustomFields = append(newContact.CustomFields, contacts.CustomField{
    FieldID: 1234, // ID поля "Телефон"
    Values: []contacts.CustomFieldValue{
        {
            Value: "+79001234567",
            Enum: "WORK", // Тип телефона (рабочий)
        },
    },
})

// Добавление email
newContact.CustomFields = append(newContact.CustomFields, contacts.CustomField{
    FieldID: 5678, // ID поля "Email"
    Values: []contacts.CustomFieldValue{
        {
            Value: "ivan@example.com",
            Enum: "WORK", // Тип email (рабочий)
        },
    },
})

// Сохранение контакта
createdContact, err := contacts.CreateContact(apiClient, newContact)
if err != nil {
    // Обработка ошибки
}
```

## Получение контакта

```go
// Получение контакта по ID
contactID := 12345
contact, err := contacts.GetContact(apiClient, contactID)
if err != nil {
    // Обработка ошибки
}
```

## Получение списка контактов

```go
// Получение первых 50 контактов
contactsList, err := contacts.GetContacts(apiClient, 1, 50)
if err != nil {
    // Обработка ошибки
}

// Получение контактов с фильтрацией
filter := map[string]string{
    "query": "Иван", // Поиск по имени
    "created_at": "1609459200", // Контакты, созданные после указанной даты (timestamp)
}
filteredContacts, err := contacts.GetContacts(apiClient, 1, 50, filter)
```

## Обновление контакта

```go
// Обновление существующего контакта
contact.Name = "Иван Петрович Иванов"

// Добавление нового номера телефона
contact.CustomFields = append(contact.CustomFields, contacts.CustomField{
    FieldID: 1234, // ID поля "Телефон"
    Values: []contacts.CustomFieldValue{
        {
            Value: "+79009876543",
            Enum: "PERSONAL", // Тип телефона (личный)
        },
    },
})

updatedContact, err := contacts.UpdateContact(apiClient, contact)
if err != nil {
    // Обработка ошибки
}
```

## Пользовательские поля

Для работы с пользовательскими полями контактов используйте структуры `CustomField` и `CustomFieldValue`:

```go
// Добавление пользовательского поля
contact.CustomFields = append(contact.CustomFields, contacts.CustomField{
    FieldID: 9876, // ID пользовательского поля
    Values: []contacts.CustomFieldValue{
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

err := leads.LinkLeadWithContact(apiClient, leadID, contactID)

// Связывание контакта с компанией
import "github.com/chudno/amo_crm_sdk/entities/companies"
err := companies.LinkCompanyWithContact(apiClient, companyID, contactID)

// Связывание контакта с лидом
import "github.com/chudno/amo_crm_sdk/entities/leads"
err := leads.LinkLeadWithContact(apiClient, leadID, contactID)
```
