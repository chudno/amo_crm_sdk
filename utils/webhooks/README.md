# Модуль Вебхуки

Модуль `webhooks` предоставляет функциональность для работы с вебхуками в amoCRM.

## Содержание

- [Основные функции](#основные-функции)
- [Создание вебхука](#создание-вебхука)
- [Получение вебхука](#получение-вебхука)
- [Получение списка вебхуков](#получение-списка-вебхуков)
- [Удаление вебхука](#удаление-вебхука)
- [Типы событий](#типы-событий)
- [Типы сущностей](#типы-сущностей)
- [Обработка вебхук-уведомлений](#обработка-вебхук-уведомлений)

## Основные функции

| Функция | Описание |
|---------|----------|
| `CreateWebhook` | Создание нового вебхука |
| `GetWebhook` | Получение вебхука по ID |
| `GetWebhooks` | Получение списка вебхуков |
| `DeleteWebhook` | Удаление вебхука |

## Создание вебхука

```go
import (
    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/utils/webhooks"
)

// Инициализация клиента
apiClient := client.NewClient("https://your-domain.amocrm.ru", "your_access_token")

// Создание нового вебхука для получения уведомлений о создании контактов
newWebhook := &webhooks.Webhook{
    Destination: "https://your-server.com/webhook-handler",
    Settings: webhooks.Settings{
        EntityType: []string{webhooks.EntityTypeContact}, // Тип сущности - контакты
        EventType: []string{webhooks.EventTypeAdd},       // Тип события - добавление
    },
}

// Сохранение вебхука
createdWebhook, err := webhooks.CreateWebhook(apiClient, newWebhook)
if err != nil {
    // Обработка ошибки
}
```

## Получение вебхука

```go
// Получение вебхука по ID
webhookID := 12345
webhook, err := webhooks.GetWebhook(apiClient, webhookID)
if err != nil {
    // Обработка ошибки
}

// Вывод информации о вебхуке
fmt.Printf("URL назначения: %s\n", webhook.Destination)
fmt.Printf("Сущности: %v\n", webhook.Settings.EntityType)
fmt.Printf("События: %v\n", webhook.Settings.EventType)
```

## Получение списка вебхуков

```go
// Получение всех вебхуков
webhooksList, err := webhooks.GetWebhooks(apiClient)
if err != nil {
    // Обработка ошибки
}

// Вывод списка вебхуков
for _, webhook := range webhooksList {
    fmt.Printf("ID: %d, URL: %s\n", webhook.ID, webhook.Destination)
    fmt.Printf("  Сущности: %v, События: %v\n", 
        webhook.Settings.EntityType, webhook.Settings.EventType)
}
```

## Удаление вебхука

```go
// Удаление вебхука по ID
webhookID := 12345
err := webhooks.DeleteWebhook(apiClient, webhookID)
if err != nil {
    // Обработка ошибки
}
```

## Типы событий

Модуль `webhooks` предоставляет константы для типов событий:

| Константа | Значение | Описание |
|-----------|----------|----------|
| `webhooks.EventTypeAdd` | "add" | Создание сущности |
| `webhooks.EventTypeUpdate` | "update" | Обновление сущности |
| `webhooks.EventTypeDelete` | "delete" | Удаление сущности |
| `webhooks.EventTypeStatusChanged` | "status_changed" | Изменение статуса сделки |
| `webhooks.EventTypeNote` | "note" | Добавление примечания |
| `webhooks.EventTypeTask` | "task" | Добавление задачи |

## Типы сущностей

Модуль также предоставляет константы для типов сущностей:

| Константа | Значение | Описание |
|-----------|----------|----------|
| `webhooks.EntityTypeContact` | "contact" | Контакт |
| `webhooks.EntityTypeCompany` | "company" | Компания |
| `webhooks.EntityTypeLead` | "lead" | Сделка |
| `webhooks.EntityTypeTask` | "task" | Задача |
| `webhooks.EntityTypeNote` | "note" | Примечание |

## Обработка вебхук-уведомлений

При создании обработчика вебхук-уведомлений на вашем сервере, вы будете получать JSON-данные от amoCRM. Вот пример обработчика на Go:

```go
package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
)

// Структура для разбора webhook-уведомлений от amoCRM
type AmoCRMWebhookPayload struct {
    Leads    WebhookLeads    `json:"leads"`
    Contacts WebhookContacts `json:"contacts"`
    // Другие типы сущностей...
}

type WebhookLeads struct {
    Add    []int `json:"add"`
    Update []int `json:"update"`
    Delete []int `json:"delete"`
    // Другие события...
}

type WebhookContacts struct {
    Add    []int `json:"add"`
    Update []int `json:"update"`
    Delete []int `json:"delete"`
    // Другие события...
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
    // Проверка метода запроса
    if r.Method != http.MethodPost {
        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
        return
    }
    
    // Чтение тела запроса
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Ошибка чтения запроса", http.StatusBadRequest)
        log.Printf("Ошибка чтения запроса: %v", err)
        return
    }
    defer r.Body.Close()
    
    // Разбор JSON-данных
    var payload AmoCRMWebhookPayload
    if err := json.Unmarshal(body, &payload); err != nil {
        http.Error(w, "Ошибка разбора JSON", http.StatusBadRequest)
        log.Printf("Ошибка разбора JSON: %v", err)
        return
    }
    
    // Обработка уведомлений о добавлении сделок
    if len(payload.Leads.Add) > 0 {
        fmt.Printf("Получено уведомление о добавлении сделок: %v\n", payload.Leads.Add)
        // Дополнительная логика обработки...
    }
    
    // Обработка уведомлений о добавлении контактов
    if len(payload.Contacts.Add) > 0 {
        fmt.Printf("Получено уведомление о добавлении контактов: %v\n", payload.Contacts.Add)
        // Дополнительная логика обработки...
    }
    
    // Успешный ответ
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}

func main() {
    http.HandleFunc("/webhook-handler", handleWebhook)
    
    fmt.Println("Сервер запущен на порту 8080...")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Ошибка запуска сервера: %v", err)
    }
}
```

Помните, что ваш сервер должен быть доступен из интернета, чтобы amoCRM мог отправлять на него уведомления. Также рекомендуется добавить проверку авторизации для вашего вебхук-обработчика.
