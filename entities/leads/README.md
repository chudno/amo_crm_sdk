# Модуль Лиды

Модуль `leads` предоставляет функциональность для работы с лидами в amoCRM.

## Содержание

- [Основные функции](#основные-функции)
- [Создание лида](#создание-лида)
- [Получение лида](#получение-лида)
- [Получение списка лидов](#получение-списка-лидов)
- [Обновление лида](#обновление-лида)
- [Работа со связанными сущностями](#работа-со-связанными-сущностями)
- [Пользовательские поля](#пользовательские-поля)
- [Перемещение по воронке](#перемещение-по-воронке)

## Основные функции

| Функция | Описание |
|---------|----------|
| `CreateLead` | Создание нового лида |
| `GetLead` | Получение лида по ID |
| `GetLeads` | Получение списка лидов с фильтрацией |
| `UpdateLead` | Обновление существующего лида |
| `DeleteLead` | Удаление лида |

## Создание лида

```go
import (
    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/leads"
)

// Инициализация клиента
apiClient := client.NewClient("https://your-domain.amocrm.ru", "your_access_token")

// Создание нового лида
newLead := &leads.Lead{
    Name: "Заявка с сайта",
    Price: 10000,
    ResponsibleUserID: 12345, // ID ответственного менеджера
    StatusID: 142,  // ID статуса лида
    PipelineID: 3778, // ID воронки
}

// Сохранение лида
createdLead, err := leads.CreateLead(apiClient, newLead)
if err != nil {
    // Обработка ошибки
}
```

## Получение лида

```go
// Получение лида по ID
leadID := 12345
lead, err := leads.GetLead(apiClient, leadID)
if err != nil {
    // Обработка ошибки
}

// Получение лида со связанными контактами и компаниями
leadWithRelations, err := leads.GetLead(apiClient, leadID, leads.WithContacts, leads.WithCompanies)
if err != nil {
    // Обработка ошибки
}

// Доступ к связанным контактам
if leadWithRelations.Embedded != nil {
    for _, contact := range leadWithRelations.Embedded.Contacts {
        // Работа с контактом
    }
    
    for _, company := range leadWithRelations.Embedded.Companies {
        // Работа с компанией
    }
}
```

## Получение списка лидов

```go
// Получение первых 50 лидов
leadsList, err := leads.GetLeads(apiClient, 1, 50)
if err != nil {
    // Обработка ошибки
}

// Получение лидов с фильтрацией
filter := map[string]string{
    "status_id": "142", // Фильтр по ID статуса
    "pipeline_id": "3778", // Фильтр по ID воронки
    "created_at": "1609459200", // Лиды, созданные после указанной даты (timestamp)
}
filteredLeads, err := leads.GetLeads(apiClient, 1, 50, filter)

// Получение лидов со связанными сущностями
leadsWithRelations, err := leads.GetLeads(apiClient, 1, 50, filter, leads.WithContacts, leads.WithCompanies)
```

## Обновление лида

```go
// Обновление существующего лида
lead.Name = "Заявка с сайта - Уточненная"
lead.Price = 15000
lead.StatusID = 143 // Перемещение на следующий этап

updatedLead, err := leads.UpdateLead(apiClient, lead)
if err != nil {
    // Обработка ошибки
}
```

## Работа со связанными сущностями

```go
// Связывание лида с контактом
import "github.com/chudno/amo_crm_sdk/entities/contacts"
contactID := 67890
err := leads.LinkLeadWithContact(apiClient, leadID, contactID)
if err != nil {
    // Обработка ошибки
}

// Связывание лида с компанией
import "github.com/chudno/amo_crm_sdk/entities/companies"
companyID := 54321
err = leads.LinkLeadWithCompany(apiClient, leadID, companyID)
if err != nil {
    // Обработка ошибки
}
```

## Пользовательские поля

Для работы с пользовательскими полями лидов используйте структуры `CustomField` и `CustomFieldValue`:

```go
// Добавление пользовательского поля
lead.CustomFields = append(lead.CustomFields, leads.CustomField{
    FieldID: 9876, // ID пользовательского поля
    Values: []leads.CustomFieldValue{
        {
            Value: "Значение поля",
        },
    },
})
```

## Перемещение по воронке

```go
// Перемещение лида на другой этап воронки
lead.StatusID = 143 // ID нового статуса
lead.PipelineID = 3778 // ID воронки (если меняется)

updatedLead, err := leads.UpdateLead(apiClient, lead)
if err != nil {
    // Обработка ошибки
}
```
