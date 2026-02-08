# Модуль Компании

Модуль `companies` предоставляет функциональность для работы с компаниями в amoCRM.

## Содержание

- [Основные функции](#основные-функции)
- [Создание компании](#создание-компании)
- [Получение компании](#получение-компании)
- [Получение списка компаний](#получение-списка-компаний)
- [Обновление компании](#обновление-компании)
- [Пользовательские поля](#пользовательские-поля)

## Основные функции

| Функция | Описание |
|---------|----------|
| `Create` | Создание новой компании |
| `Get` | Получение компании по ID |
| `List` | Получение списка компаний с пагинацией |
| `Update` | Обновление существующей компании |

## Создание компании

```go
import (
    "context"
    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/companies"
)

// Инициализация клиента
apiClient := client.NewClient("https://your-domain.amocrm.ru", "your_access_token")
ctx := context.Background()

// Создание новой компании
newCompany := &companies.Company{
    Name: "ООО Ромашка",
    ResponsibleUserID: 12345, // ID ответственного менеджера
}

// Сохранение компании
createdCompany, err := companies.Create(ctx, apiClient, newCompany)
if err != nil {
    // Обработка ошибки
}
```

## Получение компании

```go
// Получение компании по ID
companyID := 12345
company, err := companies.Get(ctx, apiClient, companyID)
if err != nil {
    // Обработка ошибки
}
```

## Получение списка компаний

```go
// Получение первых 50 компаний
companiesList, err := companies.List(ctx, apiClient, 1, 50)
if err != nil {
    // Обработка ошибки
}

// Получение компаний со связанными контактами
companiesWithContacts, err := companies.List(ctx, apiClient, 1, 50, companies.WithContacts)
```

## Обновление компании

```go
// Обновление существующей компании
company.Name = "ООО Ромашка Технологии"

updatedCompany, err := companies.Update(ctx, apiClient, company)
if err != nil {
    // Обработка ошибки
}
```

## Пользовательские поля

Для работы с пользовательскими полями компаний используйте структуры из пакета `custom_fields`:

```go
import "github.com/chudno/amo_crm_sdk/utils/custom_fields"

company.CustomFieldsValues = append(company.CustomFieldsValues, custom_fields.Value{
    FieldID: 9876,
    Values: []custom_fields.FieldValue{
        {Value: "Значение поля"},
    },
})
```
