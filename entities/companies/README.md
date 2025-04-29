# Модуль Компании

Модуль `companies` предоставляет функциональность для работы с компаниями в amoCRM.

## Содержание

- [Основные функции](#основные-функции)
- [Создание компании](#создание-компании)
- [Получение компании](#получение-компании)
- [Получение списка компаний](#получение-списка-компаний)
- [Обновление компании](#обновление-компании)
- [Пользовательские поля](#пользовательские-поля)
- [Связывание компаний с другими сущностями](#связывание-компаний-с-другими-сущностями)

## Основные функции

| Функция | Описание |
|---------|----------|
| `CreateCompany` | Создание новой компании |
| `GetCompany` | Получение компании по ID |
| `GetCompanies` | Получение списка компаний с фильтрацией |
| `UpdateCompany` | Обновление существующей компании |
| `DeleteCompany` | Удаление компании |

## Создание компании

```go
import (
    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/companies"
)

// Инициализация клиента
apiClient := client.NewClient("https://your-domain.amocrm.ru", "your_access_token")

// Создание новой компании
newCompany := &companies.Company{
    Name: "ООО Ромашка",
    ResponsibleUserID: 12345, // ID ответственного менеджера
}

// Добавление номера телефона
newCompany.CustomFields = append(newCompany.CustomFields, companies.CustomField{
    FieldID: 1234, // ID поля "Телефон"
    Values: []companies.CustomFieldValue{
        {
            Value: "+79001234567",
            Enum: "WORK", // Тип телефона (рабочий)
        },
    },
})

// Добавление email
newCompany.CustomFields = append(newCompany.CustomFields, companies.CustomField{
    FieldID: 5678, // ID поля "Email"
    Values: []companies.CustomFieldValue{
        {
            Value: "info@romashka.ru",
            Enum: "WORK", // Тип email (рабочий)
        },
    },
})

// Сохранение компании
createdCompany, err := companies.CreateCompany(apiClient, newCompany)
if err != nil {
    // Обработка ошибки
}
```

## Получение компании

```go
// Получение компании по ID
companyID := 12345
company, err := companies.GetCompany(apiClient, companyID)
if err != nil {
    // Обработка ошибки
}
```

## Получение списка компаний

```go
// Получение первых 50 компаний
companiesList, err := companies.GetCompanies(apiClient, 1, 50)
if err != nil {
    // Обработка ошибки
}

// Получение компаний с фильтрацией
filter := map[string]string{
    "query": "Ромашка", // Поиск по названию
    "created_at": "1609459200", // Компании, созданные после указанной даты (timestamp)
}
filteredCompanies, err := companies.GetCompanies(apiClient, 1, 50, filter)
```

## Обновление компании

```go
// Обновление существующей компании
company.Name = "ООО Ромашка Технологии"

// Добавление нового номера телефона
company.CustomFields = append(company.CustomFields, companies.CustomField{
    FieldID: 1234, // ID поля "Телефон"
    Values: []companies.CustomFieldValue{
        {
            Value: "+79009876543",
            Enum: "WORK2", // Тип телефона (второй рабочий)
        },
    },
})

updatedCompany, err := companies.UpdateCompany(apiClient, company)
if err != nil {
    // Обработка ошибки
}
```

## Пользовательские поля

Для работы с пользовательскими полями компаний используйте структуры `CustomField` и `CustomFieldValue`:

```go
// Добавление пользовательского поля
company.CustomFields = append(company.CustomFields, companies.CustomField{
    FieldID: 9876, // ID пользовательского поля
    Values: []companies.CustomFieldValue{
        {
            Value: "Значение поля",
        },
    },
})
```

## Связывание компаний с другими сущностями

Для связывания компаний с другими сущностями используйте соответствующие методы из модулей leads, contacts и deals:

```go
// Связывание компании со сделкой
import "github.com/chudno/amo_crm_sdk/entities/deals"
err := deals.LinkDealWithCompany(apiClient, dealID, companyID)

// Связывание компании с контактом
import "github.com/chudno/amo_crm_sdk/entities/contacts"
err := companies.LinkCompanyWithContact(apiClient, companyID, contactID)

// Связывание компании с лидом
import "github.com/chudno/amo_crm_sdk/entities/leads"
err := leads.LinkLeadWithCompany(apiClient, leadID, companyID)
```
