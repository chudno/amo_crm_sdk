# Неразобранное (Unsorted)

Модуль для работы с неразобранными заявками в amoCRM.

## Оглавление

- [Возможности](#возможности)
- [Структуры](#структуры)
- [Типы источников и категорий](#типы-источников-и-категорий)
- [Примеры использования](#примеры-использования)
  - [Создание неразобранной заявки (сделки)](#создание-неразобранной-заявки-сделки)
  - [Создание неразобранной заявки (контакта)](#создание-неразобранной-заявки-контакта)
  - [Получение списка неразобранных заявок](#получение-списка-неразобранных-заявок)
  - [Получение сводки по неразобранным заявкам](#получение-сводки-по-неразобранным-заявкам)
  - [Принятие неразобранной заявки](#принятие-неразобранной-заявки)
  - [Отклонение неразобранной заявки](#отклонение-неразобранной-заявки)
  - [Связывание неразобранной заявки с существующими сущностями](#связывание-неразобранной-заявки-с-существующими-сущностями)

## Возможности

- Создание неразобранных заявок для сделок и контактов
- Получение списка неразобранных заявок с фильтрацией
- Получение сводки по неразобранным заявкам
- Принятие неразобранных заявок (преобразование в сделки/контакты)
- Отклонение неразобранных заявок
- Связывание неразобранных заявок с существующими контактами и компаниями

## Структуры

Основные структуры для работы с неразобранными заявками:

```go
// UnsortedLeadCreate представляет структуру для создания сделки из неразобранной заявки
type UnsortedLeadCreate struct {
    UnsortedBase
    Metadata         UnsortedMetadata  `json:"metadata,omitempty"`
    Contact          *UnsortedContact  `json:"contact,omitempty"`
    Company          *UnsortedCompany  `json:"company,omitempty"`
    LeadName         string            `json:"lead_name,omitempty"`
    StatusID         int               `json:"status_id,omitempty"`
    ResponsibleUserID int              `json:"responsible_user_id,omitempty"`
    Price            int               `json:"price,omitempty"`
    PipelineType     PipelineType      `json:"pipeline_type,omitempty"`
}

// UnsortedContactCreate представляет структуру для создания контакта из неразобранной заявки
type UnsortedContactCreate struct {
    UnsortedBase
    Metadata         UnsortedMetadata  `json:"metadata,omitempty"`
    Contact          *UnsortedContact  `json:"contact,omitempty"`
    Company          *UnsortedCompany  `json:"company,omitempty"`
    ResponsibleUserID int              `json:"responsible_user_id,omitempty"`
}
```

## Типы источников и категорий

Доступные типы источников неразобранных заявок:

```go
const (
    // SourceTypeAPI источник - API
    SourceTypeAPI SourceType = "api"
    // SourceTypeForms источник - Формы
    SourceTypeForms SourceType = "forms"
    // SourceTypeSite источник - Сайт
    SourceTypeSite SourceType = "site"
    // SourceTypeSip источник - Телефония
    SourceTypeSip SourceType = "sip"
    // SourceTypeEmail источник - Email
    SourceTypeEmail SourceType = "mail"
    // SourceTypeChats источник - Чаты
    SourceTypeChats SourceType = "chats"
)
```

Доступные категории неразобранных заявок:

```go
const (
    // CategoryTypeForms категория - Формы
    CategoryTypeForms CategoryType = "forms"
    // CategoryTypeSite категория - Сайт
    CategoryTypeSite CategoryType = "site"
    // CategoryTypeSip категория - Телефония
    CategoryTypeSip CategoryType = "sip"
    // CategoryTypeEmail категория - Email
    CategoryTypeEmail CategoryType = "mail"
    // CategoryTypeChats категория - Чаты
    CategoryTypeChats CategoryType = "chats"
)
```

## Примеры использования

### Создание неразобранной заявки (сделки)

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/unsorted"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // Создаем неразобранную заявку
    lead := &unsorted.UnsortedLeadCreate{
        UnsortedBase: unsorted.UnsortedBase{
            SourceName: "Наш сайт",
            SourceType: unsorted.SourceTypeSite,
            Category:   unsorted.CategoryTypeForms,
            PipelineID: 123,                     // ID воронки
            CreatedAt:  time.Now().Unix(),       // Время создания заявки
        },
        LeadName: "Заявка с сайта",              // Название сделки
        Price:    15000,                         // Бюджет сделки
        Contact: &unsorted.UnsortedContact{
            Name:  "Иван Иванов",                // Имя контакта
            Email: "ivan@example.com",           // Email контакта
            Phone: "+79001234567",               // Телефон контакта
        },
        Company: &unsorted.UnsortedCompany{
            Name: "ООО Рога и Копыта",           // Название компании
        },
        ResponsibleUserID: 456,                  // ID ответственного пользователя
        PipelineType:      unsorted.PipelineTypeLead, // Тип воронки
        Metadata: unsorted.UnsortedMetadata{
            IP: "192.168.1.1",                   // IP-адрес посетителя
            Form: map[string]interface{}{
                "form_name": "Форма обратной связи",
                "form_id": "contact-form-123",
                "form_page": "https://example.com/contact",
            },
        },
    }

    // Отправляем запрос на создание неразобранной заявки
    response, err := unsorted.CreateUnsortedLead(apiClient, lead)
    if err != nil {
        log.Fatalf("Ошибка при создании неразобранной заявки: %v", err)
    }

    // Выводим информацию о созданной заявке
    fmt.Printf("Неразобранная заявка успешно создана. UID: %s\n", response.UID)
}
```

### Создание неразобранной заявки (контакта)

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/unsorted"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // Создаем неразобранную заявку контакта
    contact := &unsorted.UnsortedContactCreate{
        UnsortedBase: unsorted.UnsortedBase{
            SourceName: "Чат на сайте",
            SourceType: unsorted.SourceTypeChats,
            Category:   unsorted.CategoryTypeChats,
            CreatedAt:  time.Now().Unix(),       // Время создания заявки
        },
        Contact: &unsorted.UnsortedContact{
            Name:  "Петр Петров",                // Имя контакта
            Email: "petr@example.com",           // Email контакта
            Phone: "+79001234568",               // Телефон контакта
        },
        ResponsibleUserID: 456,                  // ID ответственного пользователя
        Metadata: unsorted.UnsortedMetadata{
            Service: "LiveChat",                 // Название сервиса чата
            IP: "192.168.1.2",                   // IP-адрес посетителя
        },
    }

    // Отправляем запрос на создание неразобранной заявки
    response, err := unsorted.CreateUnsortedContact(apiClient, contact)
    if err != nil {
        log.Fatalf("Ошибка при создании неразобранного контакта: %v", err)
    }

    // Выводим информацию о созданной заявке
    fmt.Printf("Неразобранный контакт успешно создан. UID: %s\n", response.UID)
}
```

### Получение списка неразобранных заявок

Получение списка неразобранных заявок с фильтрацией:

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/unsorted"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // Параметры запроса
    page := 1
    limit := 50
    
    // Фильтр для получения заявок с определенной категорией
    filter := map[string]string{
        "filter[category]": string(unsorted.CategoryTypeSite),
    }

    // Получаем список неразобранных заявок для сделок
    items, err := unsorted.GetUnsortedLeads(apiClient, page, limit, filter)
    if err != nil {
        log.Fatalf("Ошибка при получении неразобранных заявок: %v", err)
    }

    // Выводим информацию о полученных заявках
    fmt.Printf("Получено %d неразобранных заявок\n", len(items))
    for _, item := range items {
        fmt.Printf("ID: %s, UID: %s\n", item.ID, item.UID)
        fmt.Printf("Источник: %s (%s)\n", item.SourceName, item.SourceType)
        fmt.Printf("Категория: %s\n", item.Category)
        fmt.Printf("Создана: %s\n", time.Unix(item.CreatedAt, 0).Format("2006-01-02 15:04:05"))
        
        // Выводим информацию о связанных сущностях, если они есть
        if item.Embedded != nil {
            if len(item.Embedded.Leads) > 0 {
                fmt.Printf("Связанные сделки: %d\n", len(item.Embedded.Leads))
                for _, lead := range item.Embedded.Leads {
                    fmt.Printf("  ID: %d, Название: %s\n", lead.ID, lead.Name)
                }
            }
            
            if len(item.Embedded.Contacts) > 0 {
                fmt.Printf("Связанные контакты: %d\n", len(item.Embedded.Contacts))
                for _, contact := range item.Embedded.Contacts {
                    fmt.Printf("  ID: %d, Имя: %s\n", contact.ID, contact.Name)
                }
            }
        }
        
        fmt.Println("---")
    }
}
```

### Получение сводки по неразобранным заявкам

```go
package main

import (
    "fmt"
    "log"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/unsorted"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // Получаем сводку по неразобранным заявкам
    summary, err := unsorted.GetUnsortedSummary(apiClient)
    if err != nil {
        log.Fatalf("Ошибка при получении сводки по неразобранным заявкам: %v", err)
    }

    // Выводим информацию о количестве заявок
    if total, ok := summary["total"].(map[string]interface{}); ok {
        if count, ok := total["count"].(float64); ok {
            fmt.Printf("Всего неразобранных заявок: %.0f\n", count)
        }
    }

    if accepted, ok := summary["accepted"].(map[string]interface{}); ok {
        if count, ok := accepted["count"].(float64); ok {
            fmt.Printf("Принятых заявок: %.0f\n", count)
        }
    }

    if declined, ok := summary["declined"].(map[string]interface{}); ok {
        if count, ok := declined["count"].(float64); ok {
            fmt.Printf("Отклоненных заявок: %.0f\n", count)
        }
    }

    if unprocessed, ok := summary["unprocessed"].(map[string]interface{}); ok {
        if count, ok := unprocessed["count"].(float64); ok {
            fmt.Printf("Необработанных заявок: %.0f\n", count)
        }
    }
}
```

### Принятие неразобранной заявки

Принятие неразобранной заявки (преобразование в сделку):

```go
package main

import (
    "fmt"
    "log"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/unsorted"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // UID неразобранной заявки
    unsortedUID := "unsorted-uid-123"
    
    // ID статуса, в который нужно поместить сделку
    statusID := 12345
    
    // ID ответственного пользователя
    responsibleUserID := 456

    // Принимаем неразобранную заявку и преобразуем ее в сделку
    leadID, err := unsorted.AcceptUnsortedLead(apiClient, unsortedUID, statusID, responsibleUserID)
    if err != nil {
        log.Fatalf("Ошибка при принятии неразобранной заявки: %v", err)
    }

    // Выводим информацию о созданной сделке
    fmt.Printf("Неразобранная заявка успешно принята. ID созданной сделки: %d\n", leadID)
}
```

Принятие неразобранной заявки (преобразование в контакт):

```go
package main

import (
    "fmt"
    "log"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/unsorted"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // UID неразобранной заявки
    unsortedUID := "unsorted-contact-uid-123"
    
    // ID ответственного пользователя
    responsibleUserID := 456

    // Принимаем неразобранную заявку и преобразуем ее в контакт
    contactID, err := unsorted.AcceptUnsortedContact(apiClient, unsortedUID, responsibleUserID)
    if err != nil {
        log.Fatalf("Ошибка при принятии неразобранного контакта: %v", err)
    }

    // Выводим информацию о созданном контакте
    fmt.Printf("Неразобранный контакт успешно принят. ID созданного контакта: %d\n", contactID)
}
```

### Отклонение неразобранной заявки

```go
package main

import (
    "fmt"
    "log"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/unsorted"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // UID неразобранной заявки сделки
    unsortedLeadUID := "unsorted-lead-uid-123"

    // Отклоняем неразобранную заявку сделки
    err := unsorted.DeclineUnsortedLead(apiClient, unsortedLeadUID)
    if err != nil {
        log.Fatalf("Ошибка при отклонении неразобранной заявки сделки: %v", err)
    }

    fmt.Printf("Неразобранная заявка сделки с UID %s успешно отклонена\n", unsortedLeadUID)

    // UID неразобранной заявки контакта
    unsortedContactUID := "unsorted-contact-uid-123"

    // Отклоняем неразобранную заявку контакта
    err = unsorted.DeclineUnsortedContact(apiClient, unsortedContactUID)
    if err != nil {
        log.Fatalf("Ошибка при отклонении неразобранной заявки контакта: %v", err)
    }

    fmt.Printf("Неразобранная заявка контакта с UID %s успешно отклонена\n", unsortedContactUID)
}
```

### Связывание неразобранной заявки с существующими сущностями

Связывание неразобранной заявки сделки с существующим контактом:

```go
package main

import (
    "fmt"
    "log"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/unsorted"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // UID неразобранной заявки сделки
    unsortedLeadUID := "unsorted-lead-uid-123"
    
    // ID существующего контакта
    contactID := 789

    // Связываем неразобранную заявку с контактом
    err := unsorted.LinkUnsortedLeadWithContact(apiClient, unsortedLeadUID, contactID)
    if err != nil {
        log.Fatalf("Ошибка при связывании неразобранной заявки с контактом: %v", err)
    }

    fmt.Printf("Неразобранная заявка сделки с UID %s успешно связана с контактом ID %d\n", unsortedLeadUID, contactID)
}
```

Связывание неразобранной заявки сделки с существующей компанией:

```go
package main

import (
    "fmt"
    "log"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/unsorted"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // UID неразобранной заявки сделки
    unsortedLeadUID := "unsorted-lead-uid-123"
    
    // ID существующей компании
    companyID := 456

    // Связываем неразобранную заявку с компанией
    err := unsorted.LinkUnsortedLeadWithCompany(apiClient, unsortedLeadUID, companyID)
    if err != nil {
        log.Fatalf("Ошибка при связывании неразобранной заявки с компанией: %v", err)
    }

    fmt.Printf("Неразобранная заявка сделки с UID %s успешно связана с компанией ID %d\n", unsortedLeadUID, companyID)
}
```
