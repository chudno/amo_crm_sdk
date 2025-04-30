# События (Events)

Модуль для работы с событиями в amoCRM. События позволяют отслеживать действия и изменения в системе.

## Оглавление

- [Возможности](#возможности)
- [Структура события](#структура-события)
- [Типы событий](#типы-событий)
- [Типы сущностей](#типы-сущностей)
- [Примеры использования](#примеры-использования)
  - [Получение списка событий](#получение-списка-событий)
  - [Фильтрация событий](#фильтрация-событий)
  - [Получение информации о конкретном событии](#получение-информации-о-конкретном-событии)
  - [Пагинация и сортировка](#пагинация-и-сортировка)

## Возможности

- Получение списка событий с поддержкой фильтрации
- Получение информации о конкретном событии
- Поддержка пагинации и сортировки результатов
- Получение связанных сущностей

## Структура события

```go
// Event структура события в amoCRM.
type Event struct {
    ID                 int             `json:"id,omitempty"`
    Type               EventType       `json:"type"`
    EntityID           int             `json:"entity_id"`
    EntityType         EventEntityType `json:"entity_type"`
    CreatedBy          int             `json:"created_by,omitempty"`
    AccountID          int             `json:"account_id,omitempty"`
    CreatedAt          int64           `json:"created_at,omitempty"`
    ValueAfter         json.RawMessage `json:"value_after,omitempty"`
    ValueBefore        json.RawMessage `json:"value_before,omitempty"`
    ValueBeforePretty  string          `json:"value_before_pretty,omitempty"`
    ValueAfterPretty   string          `json:"value_after_pretty,omitempty"`
    AdditionalEntities EventEntities   `json:"additional_entities,omitempty"`
    Link               string          `json:"link,omitempty"`
    Ver                string          `json:"__v,omitempty"`
    Embedded           *EventEmbedded  `json:"_embedded,omitempty"`
    Links              *EventLinks     `json:"_links,omitempty"`
}
```

## Типы событий

События в amoCRM могут иметь различные типы:

```go
const (
    // EventTypeNote тип события - Примечание
    EventTypeNote EventType = "note"
    // EventTypeCall тип события - Звонок
    EventTypeCall EventType = "call"
    // EventTypeMeeting тип события - Встреча
    EventTypeMeeting EventType = "meeting"
    // EventTypeIncomingLead тип события - Входящий лид
    EventTypeIncomingLead EventType = "incoming_lead"
    // EventTypeTaskResult тип события - Результат по задаче
    EventTypeTaskResult EventType = "task_result"
    // EventTypeMessage тип события - Сообщение
    EventTypeMessage EventType = "message"
    // EventTypeSendEmailStatus тип события - Статус отправки email
    EventTypeSendEmailStatus EventType = "send_email_status"
    // EventTypeCatalogObject тип события - Объект каталога
    EventTypeCatalogObject EventType = "catalog_object"
    // EventTypeEntityView тип события - Просмотр сущности
    EventTypeEntityView EventType = "entity_view"
    // EventTypeEntityUpdate тип события - Обновление сущности
    EventTypeEntityUpdate EventType = "entity_update"
    // EventTypeEntityStatusChange тип события - Изменение статуса сущности
    EventTypeEntityStatusChange EventType = "entity_status_change"
    // EventTypeEntityResponsibleChange тип события - Изменение ответственного сущности
    EventTypeEntityResponsibleChange EventType = "entity_responsible_change"
    // EventTypeEntityCreate тип события - Создание сущности
    EventTypeEntityCreate EventType = "entity_create"
    // EventTypeEntityDelete тип события - Удаление сущности
    EventTypeEntityDelete EventType = "entity_delete"
    // EventTypeActivityCreate тип события - Создание активности
    EventTypeActivityCreate EventType = "activity_create"
    // EventTypeActivityUpdate тип события - Обновление активности
    EventTypeActivityUpdate EventType = "activity_update"
    // EventTypeActivityStatusChange тип события - Изменение статуса активности
    EventTypeActivityStatusChange EventType = "activity_status_change"
    // EventTypeActivityDelete тип события - Удаление активности
    EventTypeActivityDelete EventType = "activity_delete"
)
```

## Типы сущностей

События могут быть связаны с различными типами сущностей:

```go
const (
    // EventEntityTypeLead тип сущности события - Сделка
    EventEntityTypeLead EventEntityType = "lead"
    // EventEntityTypeContact тип сущности события - Контакт
    EventEntityTypeContact EventEntityType = "contact"
    // EventEntityTypeCompany тип сущности события - Компания
    EventEntityTypeCompany EventEntityType = "company"
    // EventEntityTypeCustomer тип сущности события - Покупатель
    EventEntityTypeCustomer EventEntityType = "customer"
    // EventEntityTypeTask тип сущности события - Задача
    EventEntityTypeTask EventEntityType = "task"
)
```

## Примеры использования

### Получение списка событий

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/events"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // Получаем список событий с лимитом 50 записей
    eventsList, err := events.GetEvents(apiClient, events.WithLimit(50))
    if err != nil {
        log.Fatalf("Ошибка при получении списка событий: %v", err)
    }

    // Выводим информацию о полученных событиях
    fmt.Printf("Получено %d событий\n", len(eventsList))
    for i, event := range eventsList {
        fmt.Printf("%d. ID: %d, Тип: %s, Связанная сущность: %s (ID: %d)\n", 
            i+1, event.ID, event.Type, event.EntityType, event.EntityID)
        fmt.Printf("   Создано: %s\n", 
            time.Unix(event.CreatedAt, 0).Format("2006-01-02 15:04:05"))
        
        if event.ValueAfterPretty != "" {
            fmt.Printf("   Содержание: %s\n", event.ValueAfterPretty)
        }
        fmt.Println("---")
    }
}
```

### Фильтрация событий

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/events"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // Создаем фильтр по типу события и типу сущности
    filter := map[string]string{
        "filter[type]":        string(events.EventTypeNote),      // Примечания
        "filter[entity_type]": string(events.EventEntityTypeLead), // Связанные со сделками
    }

    // Можно добавить фильтрацию по ID сущности
    // filter["filter[entity_id]"] = "12345"
    
    // Можно добавить фильтрацию по автору
    // filter["filter[created_by]"] = "67890"
    
    // Можно добавить фильтрацию по дате создания (в формате Unix timestamp)
    // startTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
    // endTime := time.Now().Unix()
    // filter["filter[created_at][from]"] = fmt.Sprintf("%d", startTime)
    // filter["filter[created_at][to]"] = fmt.Sprintf("%d", endTime)

    // Получаем список отфильтрованных событий
    eventsList, err := events.GetEvents(apiClient, 
        events.WithFilter(filter),
        events.WithLimit(30),
    )
    if err != nil {
        log.Fatalf("Ошибка при получении списка событий: %v", err)
    }

    // Выводим информацию о полученных событиях
    fmt.Printf("Получено %d примечаний для сделок\n", len(eventsList))
    for i, event := range eventsList {
        fmt.Printf("%d. ID: %d, Сделка ID: %d\n", 
            i+1, event.ID, event.EntityID)
        fmt.Printf("   Создано: %s\n", 
            time.Unix(event.CreatedAt, 0).Format("2006-01-02 15:04:05"))
        fmt.Printf("   Текст примечания: %s\n", event.ValueAfterPretty)
        fmt.Println("---")
    }
}
```

### Получение информации о конкретном событии

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/events"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // ID события
    eventID := 12345

    // Получаем информацию о событии с информацией о связанной сущности
    event, err := events.GetEvent(apiClient, eventID, events.WithEntity())
    if err != nil {
        log.Fatalf("Ошибка при получении информации о событии: %v", err)
    }

    // Выводим информацию о событии
    fmt.Printf("Информация о событии (ID: %d):\n", event.ID)
    fmt.Printf("Тип: %s\n", event.Type)
    fmt.Printf("Связанная сущность: %s (ID: %d)\n", event.EntityType, event.EntityID)
    fmt.Printf("Создано: %s\n", time.Unix(event.CreatedAt, 0).Format("2006-01-02 15:04:05"))
    
    if event.ValueAfterPretty != "" {
        fmt.Printf("Содержание: %s\n", event.ValueAfterPretty)
    }
    
    // Выводим информацию о связанной сущности, если есть
    if event.Embedded != nil && event.Embedded.Entity != nil {
        fmt.Println("\nИнформация о связанной сущности:")
        fmt.Printf("Название: %s\n", event.Embedded.Entity.Name)
        fmt.Printf("Создано: %s\n", time.Unix(event.Embedded.Entity.Created, 0).Format("2006-01-02 15:04:05"))
        fmt.Printf("Обновлено: %s\n", time.Unix(event.Embedded.Entity.Updated, 0).Format("2006-01-02 15:04:05"))
    }
}
```

### Пагинация и сортировка

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/events"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // Получаем список событий с пагинацией и сортировкой
    page := 1
    limit := 20
    
    // Указываем страницу, лимит и сортировку по дате создания в обратном порядке (сначала новые)
    eventsList, err := events.GetEvents(apiClient, 
        events.WithPage(page),
        events.WithLimit(limit),
        events.WithOrder("created_at", "desc"),
    )
    if err != nil {
        log.Fatalf("Ошибка при получении списка событий: %v", err)
    }

    // Выводим информацию о полученных событиях
    fmt.Printf("Страница %d. Получено %d событий\n", page, len(eventsList))
    for i, event := range eventsList {
        fmt.Printf("%d. ID: %d, Тип: %s, Связанная сущность: %s (ID: %d)\n", 
            i+1, event.ID, event.Type, event.EntityType, event.EntityID)
        fmt.Printf("   Создано: %s\n", 
            time.Unix(event.CreatedAt, 0).Format("2006-01-02 15:04:05"))
        fmt.Println("---")
    }
    
    // Для получения следующей страницы просто увеличиваем номер страницы
    // nextPage := page + 1
    // eventsNextPage, err := events.GetEvents(apiClient, 
    //     events.WithPage(nextPage),
    //     events.WithLimit(limit),
    //     events.WithOrder("created_at", "desc"),
    // )
}
```
