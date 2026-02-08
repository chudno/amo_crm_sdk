# Модуль Контакты

Модуль `contacts` предоставляет функциональность для работы с контактами в amoCRM.

## Содержание

- [Основные функции](#основные-функции)
- [Создание контакта](#создание-контакта)
- [Получение контакта](#получение-контакта)
- [Получение списка контактов](#получение-списка-контактов)
- [Пользовательские поля](#пользовательские-поля)
- [Связывание контактов с другими сущностями](#связывание-контактов-с-другими-сущностями)

## Основные функции

| Функция | Описание |
|---------|----------|
| `Create` | Создание нового контакта |
| `Get` | Получение контакта по ID |
| `List` | Получение списка контактов с пагинацией |
| `LinkWithCompany` | Привязка контакта к компании |

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

// Получение контактов со связанными компаниями
contactsWithCompanies, err := contacts.List(ctx, apiClient, 1, 50, contacts.WithCompanies)
```

## Пользовательские поля

Для работы с пользовательскими полями контактов используйте структуры из пакета `custom_fields`:

```go
import "github.com/chudno/amo_crm_sdk/utils/custom_fields"

contact.CustomFieldsValues = append(contact.CustomFieldsValues, custom_fields.Value{
    FieldID: 9876,
    Values: []custom_fields.FieldValue{
        {Value: "Значение поля"},
    },
})
```

## Связывание контактов с другими сущностями

### Привязка контакта к компании

Для привязки контакта к компании используйте функцию `LinkWithCompany`:

```go
// Привязка контакта к компании
err := contacts.LinkWithCompany(ctx, apiClient, contactID, companyID)
if err != nil {
    // Обработка ошибки
}
```

### Получение контакта с компаниями

Для получения контакта вместе с привязанными компаниями используйте опцию `WithCompanies`:

```go
// Получение контакта с привязанными компаниями
contact, err := contacts.Get(ctx, apiClient, contactID, contacts.WithCompanies)
```
