# Сегменты (Segments)

Модуль для работы с сегментами контактов в amoCRM. Сегменты позволяют группировать контакты по определенным критериям.

## Оглавление

- [Возможности](#возможности)
- [Структура сегмента](#структура-сегмента)
- [Типы сегментов](#типы-сегментов)
- [Структура фильтров](#структура-фильтров)
- [Примеры использования](#примеры-использования)
  - [Создание сегмента](#создание-сегмента)
  - [Получение списка сегментов](#получение-списка-сегментов)
  - [Получение информации о сегменте](#получение-информации-о-сегменте)
  - [Обновление сегмента](#обновление-сегмента)
  - [Удаление сегмента](#удаление-сегмента)
  - [Добавление контактов в сегмент](#добавление-контактов-в-сегмент)
  - [Удаление контактов из сегмента](#удаление-контактов-из-сегмента)
  - [Получение контактов сегмента](#получение-контактов-сегмента)

## Возможности

- Создание новых сегментов с гибкими фильтрами
- Получение списка сегментов с фильтрацией и пагинацией
- Получение детальной информации о сегменте
- Обновление сегментов
- Удаление сегментов
- Управление контактами в сегменте (добавление и удаление)
- Получение списка контактов в сегменте

## Структура сегмента

```go
// Segment представляет структуру сегмента в amoCRM.
type Segment struct {
    ID                 int          `json:"id,omitempty"`
    Name               string       `json:"name"`
    Color              string       `json:"color,omitempty"`
    Type               SegmentType  `json:"type,omitempty"`
    Filter             *Filter      `json:"filter,omitempty"`
    AccountID          int          `json:"account_id,omitempty"`
    CreatedBy          int          `json:"created_by,omitempty"`
    UpdatedBy          int          `json:"updated_by,omitempty"`
    CreatedAt          int64        `json:"created_at,omitempty"`
    UpdatedAt          int64        `json:"updated_at,omitempty"`
    AvailableContactsCount int      `json:"available_contacts_count,omitempty"`
    ContactsCount      int          `json:"contacts_count,omitempty"`
    IsDeleted          bool         `json:"is_deleted,omitempty"`
    Embedded          *Embedded     `json:"_embedded,omitempty"`
    Links             *Links        `json:"_links,omitempty"`
}
```

## Типы сегментов

В amoCRM существует два типа сегментов:

```go
// SegmentType тип сегмента
type SegmentType string

const (
    // SegmentTypeDisposable одноразовый сегмент
    SegmentTypeDisposable SegmentType = "disposable"
    // SegmentTypeDynamic динамический сегмент
    SegmentTypeDynamic SegmentType = "dynamic"
)
```

- **Одноразовые сегменты (disposable)** - статические списки контактов, созданные вручную
- **Динамические сегменты (dynamic)** - автоматически обновляемые списки контактов, удовлетворяющих заданным критериям

## Структура фильтров

Для динамических сегментов используются фильтры, определяющие критерии отбора контактов:

```go
// Filter фильтр сегмента
type Filter struct {
    Term  string        `json:"term,omitempty"`
    Logic string        `json:"logic,omitempty"` // "and" или "or"
    Nodes []FilterNode  `json:"nodes,omitempty"`
}

// FilterNode узел фильтра
type FilterNode struct {
    FieldID     int      `json:"field_id,omitempty"`
    FieldCode   string   `json:"field_code,omitempty"`
    EntityType  string   `json:"entity_type,omitempty"`
    Operator    string   `json:"operator,omitempty"`
    Value       string   `json:"value,omitempty"`
    Values     []string  `json:"values,omitempty"`
    MinValue    string   `json:"min_value,omitempty"`
    MaxValue    string   `json:"max_value,omitempty"`
    Term        string   `json:"term,omitempty"`
    Logic       string   `json:"logic,omitempty"`
    Nodes      []FilterNode `json:"nodes,omitempty"`
}
```

## Примеры использования

### Создание сегмента

```go
package main

import (
    "fmt"
    "log"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/segments"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // Создаем новый динамический сегмент для контактов с почтой на example.com
    segment := &segments.Segment{
        Name:  "Клиенты с почтой example.com",
        Color: "#FF5555",  // Красный цвет
        Type:  segments.SegmentTypeDynamic,
        Filter: &segments.Filter{
            Logic: "and",
            Nodes: []segments.FilterNode{
                {
                    FieldCode: "email",
                    Operator:  "contains",
                    Value:     "example.com",
                },
            },
        },
    }

    // Отправляем запрос на создание сегмента
    createdSegment, err := segments.AddSegment(apiClient, segment)
    if err != nil {
        log.Fatalf("Ошибка при создании сегмента: %v", err)
    }

    // Выводим информацию о созданном сегменте
    fmt.Printf("Сегмент успешно создан. ID: %d\n", createdSegment.ID)
    fmt.Printf("Название: %s\n", createdSegment.Name)
    fmt.Printf("Тип: %s\n", createdSegment.Type)
    fmt.Printf("Цвет: %s\n", createdSegment.Color)
}
```

### Получение списка сегментов

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/segments"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // Параметры запроса
    page := 1
    limit := 50
    
    // Фильтр по названию сегмента
    filter := map[string]string{
        "filter[name]": "Клиенты",  // ищем сегменты, в названии которых есть слово "Клиенты"
    }

    // Получаем список сегментов с фильтрацией
    segmentsList, err := segments.GetSegments(apiClient, page, limit, segments.WithFilter(filter))
    if err != nil {
        log.Fatalf("Ошибка при получении списка сегментов: %v", err)
    }

    // Выводим информацию о полученных сегментах
    fmt.Printf("Получено %d сегментов\n", len(segmentsList))
    for i, segment := range segmentsList {
        fmt.Printf("%d. ID: %d, Название: %s\n", i+1, segment.ID, segment.Name)
        fmt.Printf("   Тип: %s, Цвет: %s\n", segment.Type, segment.Color)
        fmt.Printf("   Контактов в сегменте: %d\n", segment.ContactsCount)
        fmt.Printf("   Создан: %s\n", time.Unix(segment.CreatedAt, 0).Format("2006-01-02 15:04:05"))
        fmt.Println("---")
    }
}
```

### Получение информации о сегменте

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/segments"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // ID сегмента
    segmentID := 123

    // Получаем информацию о сегменте с контактами
    segment, err := segments.GetSegment(apiClient, segmentID, segments.WithContacts())
    if err != nil {
        log.Fatalf("Ошибка при получении информации о сегменте: %v", err)
    }

    // Выводим информацию о сегменте
    fmt.Printf("Информация о сегменте (ID: %d):\n", segment.ID)
    fmt.Printf("Название: %s\n", segment.Name)
    fmt.Printf("Тип: %s\n", segment.Type)
    fmt.Printf("Цвет: %s\n", segment.Color)
    fmt.Printf("Контактов в сегменте: %d\n", segment.ContactsCount)
    fmt.Printf("Создан: %s\n", time.Unix(segment.CreatedAt, 0).Format("2006-01-02 15:04:05"))
    fmt.Printf("Обновлен: %s\n", time.Unix(segment.UpdatedAt, 0).Format("2006-01-02 15:04:05"))
    
    // Вывод контактов, если они есть
    if segment.Embedded != nil && len(segment.Embedded.Contacts) > 0 {
        fmt.Println("\nКонтакты в сегменте:")
        for i, contact := range segment.Embedded.Contacts {
            fmt.Printf("%d. ID: %d, Имя: %s\n", i+1, contact.ID, contact.Name)
        }
    }
}
```

### Обновление сегмента

```go
package main

import (
    "fmt"
    "log"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/segments"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // ID сегмента для обновления
    segmentID := 123

    // Создаем структуру с обновленными данными
    segment := &segments.Segment{
        ID:    segmentID,
        Name:  "Активные клиенты (обновлено)",
        Color: "#55FF55",  // Зеленый цвет
    }

    // Отправляем запрос на обновление сегмента
    updatedSegment, err := segments.UpdateSegment(apiClient, segment)
    if err != nil {
        log.Fatalf("Ошибка при обновлении сегмента: %v", err)
    }

    // Выводим информацию об обновленном сегменте
    fmt.Printf("Сегмент успешно обновлен. ID: %d\n", updatedSegment.ID)
    fmt.Printf("Новое название: %s\n", updatedSegment.Name)
    fmt.Printf("Новый цвет: %s\n", updatedSegment.Color)
}
```

### Удаление сегмента

```go
package main

import (
    "fmt"
    "log"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/segments"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // ID сегмента для удаления
    segmentID := 123

    // Удаляем сегмент
    err := segments.DeleteSegment(apiClient, segmentID)
    if err != nil {
        log.Fatalf("Ошибка при удалении сегмента: %v", err)
    }

    fmt.Printf("Сегмент с ID %d успешно удален\n", segmentID)
}
```

### Добавление контактов в сегмент

```go
package main

import (
    "fmt"
    "log"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/segments"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // ID сегмента
    segmentID := 123
    
    // ID контактов для добавления
    contactIDs := []int{1001, 1002, 1003}

    // Добавляем контакты в сегмент
    err := segments.AddContactsToSegment(apiClient, segmentID, contactIDs)
    if err != nil {
        log.Fatalf("Ошибка при добавлении контактов в сегмент: %v", err)
    }

    fmt.Printf("Контакты (ID: %v) успешно добавлены в сегмент ID %d\n", contactIDs, segmentID)
}
```

### Удаление контактов из сегмента

```go
package main

import (
    "fmt"
    "log"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/segments"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // ID сегмента
    segmentID := 123
    
    // ID контактов для удаления
    contactIDs := []int{1001, 1002, 1003}

    // Удаляем контакты из сегмента
    err := segments.RemoveContactsFromSegment(apiClient, segmentID, contactIDs)
    if err != nil {
        log.Fatalf("Ошибка при удалении контактов из сегмента: %v", err)
    }

    fmt.Printf("Контакты (ID: %v) успешно удалены из сегмента ID %d\n", contactIDs, segmentID)
}
```

### Получение контактов сегмента

```go
package main

import (
    "fmt"
    "log"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/segments"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // ID сегмента
    segmentID := 123
    
    // Параметры пагинации
    page := 1
    limit := 50

    // Получаем ID контактов в сегменте
    contactIDs, err := segments.GetSegmentContacts(apiClient, segmentID, page, limit)
    if err != nil {
        log.Fatalf("Ошибка при получении контактов сегмента: %v", err)
    }

    // Выводим информацию о полученных контактах
    fmt.Printf("Получено %d контактов в сегменте ID %d:\n", len(contactIDs), segmentID)
    for i, contactID := range contactIDs {
        fmt.Printf("%d. Контакт ID: %d\n", i+1, contactID)
    }
    
    // Для получения полной информации о контактах необходимо использовать модуль contacts
    // и выполнить запрос для каждого contactID
}
```
