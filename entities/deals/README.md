# Модуль Сделки

Модуль `deals` предоставляет функциональность для работы со сделками в amoCRM.

## Содержание

- [Основные функции](#основные-функции)
- [Создание сделки](#создание-сделки)
- [Получение сделки](#получение-сделки)
- [Получение списка сделок](#получение-списка-сделок)
- [Обновление сделки](#обновление-сделки)
- [Работа со связанными сущностями](#работа-со-связанными-сущностями)
- [Пользовательские поля](#пользовательские-поля)
- [Перемещение по воронке](#перемещение-по-воронке)

## Основные функции

| Функция | Описание |
|---------|----------|
| `CreateDeal` | Создание новой сделки |
| `GetDeal` | Получение сделки по ID |
| `GetDeals` | Получение списка сделок с фильтрацией |
| `UpdateDeal` | Обновление существующей сделки |
| `DeleteDeal` | Удаление сделки |

## Создание сделки

```go
import (
    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/deals"
)

// Инициализация клиента
apiClient := client.NewClient("https://your-domain.amocrm.ru", "your_access_token")

// Создание новой сделки
newDeal := &deals.Deal{
    Name: "Продажа ПО",
    Price: 15000,
    ResponsibleUserID: 12345, // ID ответственного менеджера
    StatusID: 142,  // ID статуса сделки
    PipelineID: 3778, // ID воронки
}

// Сохранение сделки
createdDeal, err := deals.CreateDeal(apiClient, newDeal)
if err != nil {
    // Обработка ошибки
}
```

## Получение сделки

```go
// Получение сделки по ID
dealID := 12345
deal, err := deals.GetDeal(apiClient, dealID)
if err != nil {
    // Обработка ошибки
}

// Получение сделки со связанными контактами и компаниями
dealWithRelations, err := deals.GetDeal(apiClient, dealID, deals.WithContacts, deals.WithCompanies)
if err != nil {
    // Обработка ошибки
}

// Доступ к связанным контактам
if dealWithRelations.Embedded != nil {
    for _, contact := range dealWithRelations.Embedded.Contacts {
        // Работа с контактом
    }
    
    for _, company := range dealWithRelations.Embedded.Companies {
        // Работа с компанией
    }
}
```

## Получение списка сделок

```go
// Получение первых 50 сделок
dealsList, err := deals.GetDeals(apiClient, 1, 50)
if err != nil {
    // Обработка ошибки
}

// Получение сделок с фильтрацией
filter := map[string]string{
    "status_id": "142", // Фильтр по ID статуса
    "pipeline_id": "3778", // Фильтр по ID воронки
    "created_at": "1609459200", // Сделки, созданные после указанной даты (timestamp)
}
filteredDeals, err := deals.GetDeals(apiClient, 1, 50, filter)

// Получение сделок со связанными сущностями
dealsWithRelations, err := deals.GetDeals(apiClient, 1, 50, filter, deals.WithContacts, deals.WithCompanies)
```

## Обновление сделки

```go
// Обновление существующей сделки
deal.Name = "Продажа ПО - Расширенная лицензия"
deal.Price = 25000
deal.StatusID = 143 // Перемещение на следующий этап

updatedDeal, err := deals.UpdateDeal(apiClient, deal)
if err != nil {
    // Обработка ошибки
}
```

## Работа со связанными сущностями

```go
// Связывание сделки с контактом
import "github.com/chudno/amo_crm_sdk/entities/contacts"
contactID := 67890
err := deals.LinkDealWithContact(apiClient, dealID, contactID)
if err != nil {
    // Обработка ошибки
}

// Связывание сделки с компанией
import "github.com/chudno/amo_crm_sdk/entities/companies"
companyID := 54321
err = deals.LinkDealWithCompany(apiClient, dealID, companyID)
if err != nil {
    // Обработка ошибки
}
```

## Пользовательские поля

Для работы с пользовательскими полями сделок используйте структуры `CustomField` и `CustomFieldValue`:

```go
// Добавление пользовательского поля
deal.CustomFields = append(deal.CustomFields, deals.CustomField{
    FieldID: 9876, // ID пользовательского поля
    Values: []deals.CustomFieldValue{
        {
            Value: "Значение поля",
        },
    },
})
```

## Перемещение по воронке

```go
// Перемещение сделки на другой этап воронки
deal.StatusID = 143 // ID нового статуса
deal.PipelineID = 3778 // ID воронки (если меняется)

updatedDeal, err := deals.UpdateDeal(apiClient, deal)
if err != nil {
    // Обработка ошибки
}
```
