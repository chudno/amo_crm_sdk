# Звонки (Calls)

Модуль для работы со звонками в amoCRM.

## Возможности

- Добавление информации о звонках в amoCRM

> **Примечание:** API amoCRM v4 для звонков поддерживает только добавление (`POST /api/v4/calls`). Получение, обновление, удаление и связывание звонков через API не поддерживается.

## Структура звонка

```go
type Call struct {
    ID                int           `json:"id,omitempty"`
    Direction         CallDirection `json:"direction"`
    Status            CallStatus    `json:"status,omitempty"`
    CallStatusCode    int           `json:"call_status,omitempty"`
    ResponsibleUserID int           `json:"responsible_user_id,omitempty"`
    CreatedBy         int           `json:"created_by,omitempty"`
    UpdatedBy         int           `json:"updated_by,omitempty"`
    CreatedAt         int64         `json:"created_at,omitempty"`
    UpdatedAt         int64         `json:"updated_at,omitempty"`
    AccountID         int64         `json:"account_id,omitempty"`
    Uniq              string        `json:"uniq,omitempty"`
    Duration          int           `json:"duration,omitempty"`
    Source            string        `json:"source,omitempty"`
    CallResult        string        `json:"call_result,omitempty"`
    Phone             string        `json:"phone,omitempty"`
    // ... и другие поля
}
```

## Типы и статусы звонков

### Направления звонков (CallDirection)

```go
const (
    CallDirectionIncoming CallDirection = "inbound"  // Входящий звонок
    CallDirectionOutgoing CallDirection = "outbound"  // Исходящий звонок
)
```

### Статусы звонков (CallStatus)

```go
const (
    CallStatusSuccess   CallStatus = "success"   // Успешный звонок
    CallStatusMissed    CallStatus = "missed"     // Пропущенный звонок
    CallStatusVoicemail CallStatus = "voicemail"  // Голосовая почта
    CallStatusHungup    CallStatus = "hung_up"    // Сброшенный звонок
    CallStatusBusy      CallStatus = "busy"       // Занято
)
```

### Числовые коды статусов (CallStatusCode)

```go
const (
    CallStatusCodeLeftMessage = 1  // Оставлено сообщение
    CallStatusCodeInterrupted = 2  // Разговор прерван
    CallStatusCodeNoAnswer    = 3  // Нет ответа
    CallStatusCodeSuccess     = 4  // Успешный звонок
    CallStatusCodeWrongNumber = 5  // Неверный номер
    CallStatusCodeBusy        = 6  // Занято
    CallStatusCodeVoicemail   = 7  // Голосовая почта
)
```

> **Примечание:** Функция `Add()` автоматически преобразует строковое значение поля `Status` в числовой `CallStatusCode`, необходимый для API, если `CallStatusCode` не задан явно. Вы можете использовать либо строковый `Status` (например, `CallStatusSuccess`), либо задать `CallStatusCode` напрямую. Если указаны оба значения, приоритет имеет `CallStatusCode`.

## Примеры использования

### Добавление звонка

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/calls"
)

func main() {
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")
    ctx := context.Background()

    call := &calls.Call{
        Direction:         calls.CallDirectionIncoming,
        Status:            calls.CallStatusSuccess,
        ResponsibleUserID: 123,
        Duration:          120,
        Source:            "SIP Phone",
        CallResult:        "Клиент интересуется услугами",
        Phone:             "+79001234567",
    }

    createdCall, err := calls.Add(ctx, apiClient, call)
    if err != nil {
        log.Fatalf("Ошибка при добавлении звонка: %v", err)
    }

    fmt.Printf("Звонок добавлен. ID: %d\n", createdCall.ID)
}
```
