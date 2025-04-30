# Звонки (Calls)

Модуль для работы со звонками в amoCRM.

## Оглавление

- [Возможности](#возможности)
- [Структура звонка](#структура-звонка)
- [Типы и статусы звонков](#типы-и-статусы-звонков)
- [Примеры использования](#примеры-использования)
  - [Добавление звонка](#добавление-звонка)
  - [Получение списка звонков](#получение-списка-звонков)
  - [Получение информации о звонке](#получение-информации-о-звонке)
  - [Обновление звонка](#обновление-звонка)
  - [Удаление звонка](#удаление-звонка)
  - [Связывание звонка с сущностью](#связывание-звонка-с-сущностью)
  - [Отвязывание звонка от сущности](#отвязывание-звонка-от-сущности)

## Возможности

- Добавление информации о звонках в amoCRM
- Получение списка звонков с фильтрацией
- Получение информации о конкретном звонке
- Обновление информации о звонке
- Удаление звонков
- Связывание звонков с сущностями (сделками, контактами, компаниями)
- Отвязывание звонков от сущностей

## Структура звонка

```go
// Call представляет структуру звонка в amoCRM.
type Call struct {
    ID                   int          `json:"id,omitempty"`
    Direction            CallDirection `json:"direction"`
    Status               CallStatus   `json:"status"`
    ResponsibleUserID    int          `json:"responsible_user_id,omitempty"`
    CreatedBy            int          `json:"created_by,omitempty"`
    UpdatedBy            int          `json:"updated_by,omitempty"`
    CreatedAt            int64        `json:"created_at,omitempty"`
    UpdatedAt            int64        `json:"updated_at,omitempty"`
    AccountID            int64        `json:"account_id,omitempty"`
    Uniq                 string       `json:"uniq,omitempty"`
    Duration             int          `json:"duration,omitempty"`
    Source               string       `json:"source,omitempty"`
    CallResult           string       `json:"call_result,omitempty"`
    Link                 string       `json:"link,omitempty"`
    ServiceCode          string       `json:"service_code,omitempty"`
    Phone                string       `json:"phone,omitempty"`
    APIID                int          `json:"api_id,omitempty"`
    ManagerName          string       `json:"manager_name,omitempty"`
    ManagerEmail         string       `json:"manager_email,omitempty"`
    ManagerPhone         string       `json:"manager_phone,omitempty"`
    ManagerICQ           string       `json:"manager_icq,omitempty"`
    ContactID            int          `json:"contact_id,omitempty"`
    LeadID               int          `json:"lead_id,omitempty"`
    CompanyID            int          `json:"company_id,omitempty"`
    SourceName           string       `json:"source_name,omitempty"`
    SourceUID            string       `json:"source_uid,omitempty"`
    IsCallbackCall       bool         `json:"is_callback_call,omitempty"`
    IsRinging            bool         `json:"is_ringing,omitempty"`
    Voice               *Voice        `json:"voice,omitempty"`
    CallStartTime        string       `json:"call_start_time,omitempty"`
    CallEndTime          string       `json:"call_end_time,omitempty"`
    Version              int          `json:"version,omitempty"`
    Embedded            *CallEmbedded `json:"_embedded,omitempty"`
    Links              *CallLinks     `json:"_links,omitempty"`
    EntityType          *EntityType   `json:"entity_type,omitempty"`
    EntityID             int          `json:"entity_id,omitempty"`
}
```

## Типы и статусы звонков

### Направления звонков (CallDirection)

```go
const (
    // CallDirectionIncoming входящий звонок
    CallDirectionIncoming CallDirection = "inbound"
    // CallDirectionOutgoing исходящий звонок
    CallDirectionOutgoing CallDirection = "outbound"
)
```

### Статусы звонков (CallStatus)

```go
const (
    // CallStatusSuccess успешный звонок
    CallStatusSuccess CallStatus = "success"
    // CallStatusMissed пропущенный звонок
    CallStatusMissed CallStatus = "missed"
    // CallStatusVoicemail голосовая почта
    CallStatusVoicemail CallStatus = "voicemail"
    // CallStatusHungup сброшенный звонок
    CallStatusHungup CallStatus = "hung_up"
    // CallStatusBusy занято
    CallStatusBusy CallStatus = "busy"
)
```

### Типы сущностей (EntityType)

```go
const (
    // EntityTypeLead тип сущности - Сделка
    EntityTypeLead EntityType = "leads"
    // EntityTypeContact тип сущности - Контакт
    EntityTypeContact EntityType = "contacts"
    // EntityTypeCompany тип сущности - Компания
    EntityTypeCompany EntityType = "companies"
    // EntityTypeCustomers тип сущности - Покупатель
    EntityTypeCustomers EntityType = "customers"
)
```

## Примеры использования

### Добавление звонка

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/calls"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // Создаем новый звонок
    call := &calls.Call{
        Direction:         calls.CallDirectionIncoming,     // Входящий звонок
        Status:            calls.CallStatusSuccess,         // Успешный звонок
        ResponsibleUserID: 123,                            // ID ответственного менеджера
        Duration:          120,                            // Продолжительность звонка в секундах
        Source:            "SIP Phone",                    // Источник звонка
        CallResult:        "Клиент интересуется услугами", // Результат звонка
        Phone:             "+79001234567",                 // Телефон клиента
        CreatedAt:         time.Now().Unix(),              // Время создания записи
    }

    // Отправляем запрос на создание звонка
    createdCall, err := calls.AddCall(apiClient, call)
    if err != nil {
        log.Fatalf("Ошибка при добавлении звонка: %v", err)
    }

    // Выводим информацию о созданном звонке
    fmt.Printf("Звонок успешно добавлен. ID: %d\n", createdCall.ID)
    fmt.Printf("Направление: %s\n", createdCall.Direction)
    fmt.Printf("Статус: %s\n", createdCall.Status)
    fmt.Printf("Продолжительность: %d секунд\n", createdCall.Duration)
    fmt.Printf("Телефон: %s\n", createdCall.Phone)
}
```

### Получение списка звонков

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/calls"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // Параметры запроса
    page := 1
    limit := 50
    
    // Фильтр для получения входящих звонков
    filter := map[string]string{
        "filter[direction]": string(calls.CallDirectionIncoming),
    }

    // Получаем список звонков с фильтрацией и с включением тегов
    callsList, err := calls.GetCalls(apiClient, page, limit, filter, calls.WithTags)
    if err != nil {
        log.Fatalf("Ошибка при получении списка звонков: %v", err)
    }

    // Выводим информацию о полученных звонках
    fmt.Printf("Получено %d звонков\n", len(callsList))
    for i, call := range callsList {
        fmt.Printf("%d. ID: %d, Направление: %s, Статус: %s\n", i+1, call.ID, call.Direction, call.Status)
        fmt.Printf("   Продолжительность: %d секунд\n", call.Duration)
        fmt.Printf("   Телефон: %s\n", call.Phone)
        fmt.Printf("   Результат: %s\n", call.CallResult)
        fmt.Printf("   Создан: %s\n", time.Unix(call.CreatedAt, 0).Format("2006-01-02 15:04:05"))
        
        // Вывод тегов, если они есть
        if call.Embedded != nil && len(call.Embedded.Tags) > 0 {
            fmt.Println("   Теги:")
            for _, tag := range call.Embedded.Tags {
                fmt.Printf("     - %s\n", tag.Name)
            }
        }
        
        fmt.Println("---")
    }
}
```

### Получение информации о звонке

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/calls"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // ID звонка
    callID := 123

    // Получаем информацию о звонке с тегами
    call, err := calls.GetCall(apiClient, callID, calls.WithTags)
    if err != nil {
        log.Fatalf("Ошибка при получении информации о звонке: %v", err)
    }

    // Выводим информацию о звонке
    fmt.Printf("Информация о звонке (ID: %d):\n", call.ID)
    fmt.Printf("Направление: %s\n", call.Direction)
    fmt.Printf("Статус: %s\n", call.Status)
    fmt.Printf("Продолжительность: %d секунд\n", call.Duration)
    fmt.Printf("Телефон: %s\n", call.Phone)
    fmt.Printf("Результат: %s\n", call.CallResult)
    fmt.Printf("Источник: %s\n", call.Source)
    fmt.Printf("Создан: %s\n", time.Unix(call.CreatedAt, 0).Format("2006-01-02 15:04:05"))
    
    // Вывод тегов, если они есть
    if call.Embedded != nil && len(call.Embedded.Tags) > 0 {
        fmt.Println("Теги:")
        for _, tag := range call.Embedded.Tags {
            fmt.Printf("  - %s\n", tag.Name)
        }
    }
}
```

### Обновление звонка

```go
package main

import (
    "fmt"
    "log"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/calls"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // ID звонка для обновления
    callID := 123

    // Создаем структуру для обновления звонка
    call := &calls.Call{
        ID:              callID,
        Duration:        180,                           // Обновляем продолжительность
        CallResult:      "Клиент согласился на встречу", // Обновляем результат звонка
        ResponsibleUserID: 456                          // Обновляем ответственного
    }

    // Отправляем запрос на обновление звонка
    updatedCall, err := calls.UpdateCall(apiClient, call)
    if err != nil {
        log.Fatalf("Ошибка при обновлении звонка: %v", err)
    }

    // Выводим информацию об обновленном звонке
    fmt.Printf("Звонок успешно обновлен. ID: %d\n", updatedCall.ID)
    fmt.Printf("Продолжительность: %d секунд\n", updatedCall.Duration)
    fmt.Printf("Результат: %s\n", updatedCall.CallResult)
    fmt.Printf("Ответственный: %d\n", updatedCall.ResponsibleUserID)
}
```

### Удаление звонка

```go
package main

import (
    "fmt"
    "log"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/calls"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // ID звонка для удаления
    callID := 123

    // Удаляем звонок
    err := calls.DeleteCall(apiClient, callID)
    if err != nil {
        log.Fatalf("Ошибка при удалении звонка: %v", err)
    }

    fmt.Printf("Звонок с ID %d успешно удален\n", callID)
}
```

### Связывание звонка с сущностью

```go
package main

import (
    "fmt"
    "log"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/calls"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // ID звонка
    callID := 123
    
    // Тип сущности и ID для связывания
    entityType := calls.EntityTypeLead // Сделка
    entityID := 456                   // ID сделки

    // Связываем звонок с сделкой
    err := calls.LinkCallWithEntity(apiClient, callID, entityType, entityID)
    if err != nil {
        log.Fatalf("Ошибка при связывании звонка с сущностью: %v", err)
    }

    fmt.Printf("Звонок с ID %d успешно связан с сделкой ID %d\n", callID, entityID)
}
```

### Отвязывание звонка от сущности

```go
package main

import (
    "fmt"
    "log"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/calls"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // ID звонка
    callID := 123
    
    // Тип сущности и ID для отвязывания
    entityType := calls.EntityTypeLead // Сделка
    entityID := 456                   // ID сделки

    // Отвязываем звонок от сделки
    err := calls.UnlinkCallFromEntity(apiClient, callID, entityType, entityID)
    if err != nil {
        log.Fatalf("Ошибка при отвязывании звонка от сущности: %v", err)
    }

    fmt.Printf("Звонок с ID %d успешно отвязан от сделки ID %d\n", callID, entityID)
}
```
