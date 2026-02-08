# Модуль Вебхуки

Модуль `webhooks` предоставляет функциональность для работы с вебхуками в amoCRM.

## Содержание

- [Основные функции](#основные-функции)
- [Создание вебхука](#создание-вебхука)
- [Получение вебхука](#получение-вебхука)
- [Получение списка вебхуков](#получение-списка-вебхуков)
- [Обновление вебхука](#обновление-вебхука)
- [Удаление вебхука](#удаление-вебхука)
- [Константы сущностей и действий](#константы-сущностей-и-действий)
- [Обработка вебхук-уведомлений](#обработка-вебхук-уведомлений)

## Основные функции

| Функция | Описание |
|---------|----------|
| `Create` | Создание нового вебхука |
| `CreateSimple` | Создание вебхука с указанием параметров напрямую |
| `Get` | Получение вебхука по ID |
| `List` | Получение списка вебхуков |
| `Update` | Обновление существующего вебхука |
| `Delete` | Удаление вебхука |

## Создание вебхука

```go
import (
    "context"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/utils/webhooks"
)

apiClient := client.NewClient("https://your-domain.amocrm.ru", "your_access_token")
ctx := context.Background()

newWebhook := &webhooks.Webhook{
    Destination: "https://your-server.com/webhook-handler",
    Settings: &webhooks.Settings{
        Entities: []string{webhooks.EntityContact},
        Actions:  []string{webhooks.ActionAdd},
    },
}

createdWebhook, err := webhooks.Create(ctx, apiClient, newWebhook)
if err != nil {
    // Обработка ошибки
}
```

Также можно использовать упрощённый вариант:

```go
createdWebhook, err := webhooks.CreateSimple(
    ctx,
    apiClient,
    "https://your-server.com/webhook-handler",
    []string{webhooks.EntityContact, webhooks.EntityLead},
    []string{webhooks.ActionAdd, webhooks.ActionUpdate},
)
```

## Получение вебхука

```go
webhook, err := webhooks.Get(ctx, apiClient, 12345)
if err != nil {
    // Обработка ошибки
}

fmt.Printf("URL назначения: %s\n", webhook.Destination)
fmt.Printf("Сущности: %v\n", webhook.Settings.Entities)
fmt.Printf("Действия: %v\n", webhook.Settings.Actions)
```

## Получение списка вебхуков

```go
webhooksList, err := webhooks.List(ctx, apiClient, 50, 1)
if err != nil {
    // Обработка ошибки
}

for _, wh := range webhooksList {
    fmt.Printf("ID: %d, URL: %s\n", wh.ID, wh.Destination)
}
```

## Обновление вебхука

```go
webhook.Settings = &webhooks.Settings{
    Entities: []string{webhooks.EntityLead, webhooks.EntityContact},
    Actions:  []string{webhooks.ActionAdd, webhooks.ActionUpdate, webhooks.ActionDelete},
}

updatedWebhook, err := webhooks.Update(ctx, apiClient, webhook)
if err != nil {
    // Обработка ошибки
}
```

## Удаление вебхука

```go
err := webhooks.Delete(ctx, apiClient, 12345)
if err != nil {
    // Обработка ошибки
}
```

## Константы сущностей и действий

### Типы сущностей

| Константа | Значение | Описание |
|-----------|----------|----------|
| `webhooks.EntityLead` | "leads" | Сделки |
| `webhooks.EntityContact` | "contacts" | Контакты |
| `webhooks.EntityCompany` | "companies" | Компании |
| `webhooks.EntityCustomer` | "customers" | Покупатели |
| `webhooks.EntityTask` | "tasks" | Задачи |

### Типы действий

| Константа | Значение | Описание |
|-----------|----------|----------|
| `webhooks.ActionAdd` | "add" | Создание сущности |
| `webhooks.ActionUpdate` | "update" | Обновление сущности |
| `webhooks.ActionDelete` | "delete" | Удаление сущности |
| `webhooks.ActionRestore` | "restore" | Восстановление сущности |
| `webhooks.ActionStatusChange` | "status" | Изменение статуса |

## Обработка вебхук-уведомлений

При создании обработчика вебхук-уведомлений на вашем сервере, вы будете получать JSON-данные от amoCRM. Вот пример обработчика на Go:

```go
package main

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
)

type AmoCRMWebhookPayload struct {
    Leads    WebhookLeads    `json:"leads"`
    Contacts WebhookContacts `json:"contacts"`
}

type WebhookLeads struct {
    Add    []int `json:"add"`
    Update []int `json:"update"`
    Delete []int `json:"delete"`
}

type WebhookContacts struct {
    Add    []int `json:"add"`
    Update []int `json:"update"`
    Delete []int `json:"delete"`
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
        return
    }

    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Ошибка чтения запроса", http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    var payload AmoCRMWebhookPayload
    if err := json.Unmarshal(body, &payload); err != nil {
        http.Error(w, "Ошибка разбора JSON", http.StatusBadRequest)
        return
    }

    if len(payload.Leads.Add) > 0 {
        fmt.Printf("Добавлены сделки: %v\n", payload.Leads.Add)
    }

    if len(payload.Contacts.Add) > 0 {
        fmt.Printf("Добавлены контакты: %v\n", payload.Contacts.Add)
    }

    w.WriteHeader(http.StatusOK)
}

func main() {
    http.HandleFunc("/webhook-handler", handleWebhook)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```
