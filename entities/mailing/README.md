# Рассылки (Mailing)

Модуль для работы с email-рассылками в amoCRM API.

## Содержание

- [Структуры данных](#структуры-данных)
- [Функции](#функции)
  - [Получение списка рассылок](#получение-списка-рассылок)
  - [Получение информации о конкретной рассылке](#получение-информации-о-конкретной-рассылке)
  - [Создание рассылки](#создание-рассылки)
  - [Обновление рассылки](#обновление-рассылки)
  - [Удаление рассылки](#удаление-рассылки)
  - [Изменение статуса рассылки](#изменение-статуса-рассылки)
  - [Получение шаблонов рассылок](#получение-шаблонов-рассылок)
  - [Управление получателями рассылки](#управление-получателями-рассылки)
- [Примеры использования](#примеры-использования)

## Структуры данных

### Mailing

```go
type Mailing struct {
    ID               int               `json:"id,omitempty"`
    Name             string            `json:"name"`
    Status           MailingStatus     `json:"status,omitempty"`
    Subject          string            `json:"subject"`
    Template         *Template         `json:"template,omitempty"`
    Frequency        MailingFrequency  `json:"frequency,omitempty"`
    SendAt           *time.Time        `json:"send_at,omitempty"`
    CreatedAt        int64             `json:"created_at,omitempty"`
    UpdatedAt        int64             `json:"updated_at,omitempty"`
    CreatedBy        int               `json:"created_by,omitempty"`
    UpdatedBy        int               `json:"updated_by,omitempty"`
    SegmentIDs       []int             `json:"segment_ids,omitempty"`
    SegmentFilters   []SegmentFilter   `json:"segment_filters,omitempty"`
    SelectedContacts []int             `json:"selected_contacts,omitempty"`
    ExcludedContacts []int             `json:"excluded_contacts,omitempty"`
    Stats            *MailingStats     `json:"stats,omitempty"`
    AccountID        int               `json:"account_id,omitempty"`
    FromEmail        string            `json:"from_email,omitempty"`
    FromName         string            `json:"from_name,omitempty"`
    ReplyToEmail     string            `json:"reply_to_email,omitempty"`
    Settings         map[string]string `json:"settings,omitempty"`
}
```

### MailingStatus

```go
type MailingStatus string

const (
    // MailingStatusDraft - черновик
    MailingStatusDraft MailingStatus = "draft"
    // MailingStatusScheduled - запланирована
    MailingStatusScheduled MailingStatus = "scheduled"
    // MailingStatusActive - активна
    MailingStatusActive MailingStatus = "active"
    // MailingStatusPaused - приостановлена
    MailingStatusPaused MailingStatus = "paused"
    // MailingStatusCompleted - завершена
    MailingStatusCompleted MailingStatus = "completed"
    // MailingStatusStopped - остановлена
    MailingStatusStopped MailingStatus = "stopped"
)
```

### MailingFrequency

```go
type MailingFrequency string

const (
    // MailingFrequencyOnce - однократно
    MailingFrequencyOnce MailingFrequency = "once"
    // MailingFrequencyDaily - ежедневно
    MailingFrequencyDaily MailingFrequency = "daily"
    // MailingFrequencyWeekly - еженедельно
    MailingFrequencyWeekly MailingFrequency = "weekly"
    // MailingFrequencyMonthly - ежемесячно
    MailingFrequencyMonthly MailingFrequency = "monthly"
)
```

### Template

```go
type Template struct {
    ID      int    `json:"id,omitempty"`
    Name    string `json:"name,omitempty"`
    Content string `json:"content,omitempty"`
    HTML    string `json:"html,omitempty"`
    Type    string `json:"type,omitempty"`
}
```

### MailingStats

```go
type MailingStats struct {
    TotalRecipients int `json:"total_recipients"`
    Delivered       int `json:"delivered"`
    Opened          int `json:"opened"`
    Clicked         int `json:"clicked"`
    Bounced         int `json:"bounced"`
    Unsubscribed    int `json:"unsubscribed"`
    Complaints      int `json:"complaints"`
}
```

## Функции

### Получение списка рассылок

```go
func GetMailings(apiClient *client.Client, page, limit int, options ...WithOption) ([]Mailing, error)
```

Возвращает список рассылок с поддержкой пагинации и фильтрации.

**Параметры:**
- `apiClient` - клиент API amoCRM
- `page` - номер страницы (начиная с 1)
- `limit` - количество элементов на странице
- `options` - дополнительные опции (например, фильтры)

**Фильтрация:**

Для фильтрации можно использовать функцию `WithFilter`:

```go
filter := map[string]string{
    "filter[status]": "active",
}
mailings, err := mailing.GetMailings(apiClient, 1, 50, mailing.WithFilter(filter))
```

Также доступны готовые функции для фильтрации:

```go
// Фильтрация по статусу
mailings, err := mailing.GetMailings(apiClient, 1, 50, mailing.WithStatus(mailing.MailingStatusActive))

// Фильтрация по дате создания
from := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
to := time.Now()
mailings, err := mailing.GetMailings(apiClient, 1, 50, 
    mailing.WithDateFrom(from), 
    mailing.WithDateTo(to))
```

### Получение информации о конкретной рассылке

```go
func GetMailing(apiClient *client.Client, id int) (*Mailing, error)
```

Возвращает информацию о конкретной рассылке по её ID.

**Параметры:**
- `apiClient` - клиент API amoCRM
- `id` - ID рассылки

### Создание рассылки

```go
func CreateMailing(apiClient *client.Client, mailingData *Mailing) (*Mailing, error)
```

Создаёт новую рассылку и возвращает информацию о созданной рассылке.

**Параметры:**
- `apiClient` - клиент API amoCRM
- `mailingData` - структура с данными для создания новой рассылки

### Обновление рассылки

```go
func UpdateMailing(apiClient *client.Client, mailingData *Mailing) (*Mailing, error)
```

Обновляет существующую рассылку и возвращает обновлённую информацию.

**Параметры:**
- `apiClient` - клиент API amoCRM
- `mailingData` - структура с данными для обновления рассылки (должен быть заполнен ID)

### Удаление рассылки

```go
func DeleteMailing(apiClient *client.Client, id int) error
```

Удаляет рассылку по ID.

**Параметры:**
- `apiClient` - клиент API amoCRM
- `id` - ID рассылки для удаления

### Изменение статуса рассылки

```go
func ChangeMailingStatus(apiClient *client.Client, id int, status MailingStatus) (*Mailing, error)
```

Изменяет статус рассылки.

**Параметры:**
- `apiClient` - клиент API amoCRM
- `id` - ID рассылки
- `status` - новый статус рассылки

### Получение шаблонов рассылок

```go
func GetMailingTemplates(apiClient *client.Client, page, limit int) ([]Template, error)
```

Возвращает список шаблонов рассылок.

**Параметры:**
- `apiClient` - клиент API amoCRM
- `page` - номер страницы (начиная с 1)
- `limit` - количество элементов на странице

```go
func GetMailingTemplate(apiClient *client.Client, id int) (*Template, error)
```

Возвращает информацию о конкретном шаблоне рассылки.

**Параметры:**
- `apiClient` - клиент API amoCRM
- `id` - ID шаблона рассылки

### Управление получателями рассылки

```go
func AddMailingRecipients(apiClient *client.Client, id int, contactIDs []int) error
```

Добавляет получателей в рассылку.

**Параметры:**
- `apiClient` - клиент API amoCRM
- `id` - ID рассылки
- `contactIDs` - массив ID контактов для добавления в рассылку

```go
func RemoveMailingRecipients(apiClient *client.Client, id int, contactIDs []int) error
```

Удаляет получателей из рассылки.

**Параметры:**
- `apiClient` - клиент API amoCRM
- `id` - ID рассылки
- `contactIDs` - массив ID контактов для удаления из рассылки

## Примеры использования

### Создание новой рассылки

```go
package main

import (
    "fmt"
    "log"

    "github.com/chudno/amo_crm_sdk/auth"
    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/mailing"
)

func main() {
    // Создаем клиент API
    apiClient, err := client.NewClientWithAuth(auth.NewOAuthConfig(
        "ваш_домен.amocrm.ru",
        "client_id",
        "client_secret",
        "redirect_uri",
        "code",
    ))
    if err != nil {
        log.Fatalf("Ошибка аутентификации: %v", err)
    }

    // Создаем новую рассылку
    newMailing := &mailing.Mailing{
        Name:       "Новогодняя акция",
        Subject:    "Скидки к Новому году",
        Frequency:  mailing.MailingFrequencyOnce,
        FromName:   "Компания",
        FromEmail:  "marketing@example.com",
        SegmentIDs: []int{101, 102}, // ID сегментов контактов для рассылки
    }

    // Отправляем запрос на создание
    createdMailing, err := mailing.CreateMailing(apiClient, newMailing)
    if err != nil {
        log.Fatalf("Ошибка при создании рассылки: %v", err)
    }

    // Выводим информацию о созданной рассылке
    fmt.Printf("Создана рассылка: %s (ID: %d)\n", createdMailing.Name, createdMailing.ID)
    fmt.Printf("Статус: %s\n", createdMailing.Status)
}
```

### Управление статусом рассылки

```go
package main

import (
    "fmt"
    "log"

    "github.com/chudno/amo_crm_sdk/auth"
    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/mailing"
)

func main() {
    // Создаем клиент API
    apiClient, err := client.NewClientWithAuth(auth.NewOAuthConfig(
        "ваш_домен.amocrm.ru",
        "client_id",
        "client_secret",
        "redirect_uri",
        "code",
    ))
    if err != nil {
        log.Fatalf("Ошибка аутентификации: %v", err)
    }

    // ID рассылки
    mailingID := 1001

    // Запускаем рассылку
    updatedMailing, err := mailing.ChangeMailingStatus(apiClient, mailingID, mailing.MailingStatusActive)
    if err != nil {
        log.Fatalf("Ошибка при изменении статуса рассылки: %v", err)
    }

    fmt.Printf("Статус рассылки изменен на: %s\n", updatedMailing.Status)

    // Позже приостанавливаем рассылку
    pausedMailing, err := mailing.ChangeMailingStatus(apiClient, mailingID, mailing.MailingStatusPaused)
    if err != nil {
        log.Fatalf("Ошибка при приостановке рассылки: %v", err)
    }

    fmt.Printf("Рассылка приостановлена. Текущий статус: %s\n", pausedMailing.Status)
}
```

### Получение статистики рассылок

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/chudno/amo_crm_sdk/auth"
    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/mailing"
)

func main() {
    // Создаем клиент API
    apiClient, err := client.NewClientWithAuth(auth.NewOAuthConfig(
        "ваш_домен.amocrm.ru",
        "client_id",
        "client_secret",
        "redirect_uri",
        "code",
    ))
    if err != nil {
        log.Fatalf("Ошибка аутентификации: %v", err)
    }

    // Получаем список активных рассылок
    mailings, err := mailing.GetMailings(apiClient, 1, 10, mailing.WithStatus(mailing.MailingStatusActive))
    if err != nil {
        log.Fatalf("Ошибка при получении списка рассылок: %v", err)
    }

    // Выводим информацию о рассылках и их статистике
    for _, m := range mailings {
        fmt.Printf("Рассылка: %s (ID: %d)\n", m.Name, m.ID)
        
        // Получаем детальную статистику для каждой рассылки
        stats, err := mailing.GetMailingStats(apiClient, m.ID)
        if err != nil {
            log.Printf("Ошибка при получении статистики для рассылки %d: %v", m.ID, err)
            continue
        }
        
        fmt.Printf("  Всего получателей: %d\n", stats.TotalRecipients)
        fmt.Printf("  Доставлено: %d\n", stats.Delivered)
        fmt.Printf("  Открыто: %d\n", stats.Opened)
        fmt.Printf("  Клики: %d\n", stats.Clicked)
        fmt.Printf("  Отказы доставки: %d\n", stats.Bounced)
        fmt.Printf("  Отписки: %d\n", stats.Unsubscribed)
        fmt.Printf("  Жалобы: %d\n", stats.Complaints)
        
        fmt.Println()
    }
}
```
