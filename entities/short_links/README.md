# Короткие ссылки (Short Links)

Модуль для работы с короткими ссылками в amoCRM API.

## Содержание

- [Структуры данных](#структуры-данных)
- [Функции](#функции)
  - [Получение списка коротких ссылок](#получение-списка-коротких-ссылок)
  - [Получение информации о конкретной ссылке](#получение-информации-о-конкретной-ссылке)
  - [Создание короткой ссылки](#создание-короткой-ссылки)
  - [Обновление короткой ссылки](#обновление-короткой-ссылки)
  - [Удаление короткой ссылки](#удаление-короткой-ссылки)
  - [Получение статистики использования ссылки](#получение-статистики-использования-ссылки)
- [Примеры использования](#примеры-использования)

## Структуры данных

### ShortLink

```go
type ShortLink struct {
    ID           int    `json:"id,omitempty"`
    URL          string `json:"url"`
    Key          string `json:"key,omitempty"`
    ShortURL     string `json:"short_url,omitempty"`
    AccountID    int    `json:"account_id,omitempty"`
    EntityID     int    `json:"entity_id,omitempty"`
    EntityType   string `json:"entity_type,omitempty"`
    CreatedAt    int64  `json:"created_at,omitempty"`
    CreatedBy    int    `json:"created_by,omitempty"`
    UpdatedAt    int64  `json:"updated_at,omitempty"`
    MetadataID   int    `json:"metadata_id,omitempty"`
    VisitCount   int    `json:"visit_count,omitempty"`
    LastVisitAt  int64  `json:"last_visit_at,omitempty"`
    ExpireAt     int64  `json:"expire_at,omitempty"`
    UTMSource    string `json:"utm_source,omitempty"`
    UTMMedium    string `json:"utm_medium,omitempty"`
    UTMCampaign  string `json:"utm_campaign,omitempty"`
    UTMContent   string `json:"utm_content,omitempty"`
    UTMTerm      string `json:"utm_term,omitempty"`
    UseInEmbedded bool   `json:"use_in_embedded,omitempty"`
}
```

## Функции

### Получение списка коротких ссылок

```go
func List(ctx context.Context, apiClient *client.Client, page, limit int, options ...WithOption) ([]ShortLink, error)
```

Возвращает список коротких ссылок с поддержкой пагинации и фильтрации.

**Параметры:**
- `apiClient` - клиент API amoCRM
- `page` - номер страницы (начиная с 1)
- `limit` - количество элементов на странице
- `options` - дополнительные опции (например, фильтры)

**Фильтрация:**

Для фильтрации можно использовать функцию `WithFilter`:

```go
filter := map[string]string{
    "filter[entity_type]": "leads",
    "filter[entity_id]": "123",
}
shortLinks, err := short_links.List(ctx, apiClient, 1, 50, short_links.WithFilter(filter))
```

### Получение информации о конкретной ссылке

```go
func Get(ctx context.Context, apiClient *client.Client, id int) (*ShortLink, error)
```

Возвращает информацию о конкретной короткой ссылке по её ID.

**Параметры:**
- `apiClient` - клиент API amoCRM
- `id` - ID короткой ссылки

### Создание короткой ссылки

```go
func Create(ctx context.Context, apiClient *client.Client, shortLink *ShortLink) (*ShortLink, error)
```

Создаёт новую короткую ссылку и возвращает информацию о созданной ссылке.

**Параметры:**
- `apiClient` - клиент API amoCRM
- `shortLink` - структура с данными для создания новой ссылки

### Обновление короткой ссылки

```go
func Update(ctx context.Context, apiClient *client.Client, shortLink *ShortLink) (*ShortLink, error)
```

Обновляет существующую короткую ссылку и возвращает обновлённую информацию.

**Параметры:**
- `apiClient` - клиент API amoCRM
- `shortLink` - структура с данными для обновления ссылки (должен быть заполнен ID)

### Удаление короткой ссылки

```go
func Delete(ctx context.Context, apiClient *client.Client, id int) error
```

Удаляет короткую ссылку по ID.

**Параметры:**
- `apiClient` - клиент API amoCRM
- `id` - ID короткой ссылки для удаления

### Получение статистики использования ссылки

```go
func GetStats(ctx context.Context, apiClient *client.Client, id int) (*ShortLink, error)
```

Возвращает статистику использования короткой ссылки.

**Параметры:**
- `apiClient` - клиент API amoCRM
- `id` - ID короткой ссылки

## Примеры использования

### Получение списка коротких ссылок

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/short_links"
)

func main() {
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")
    ctx := context.Background()

    links, err := short_links.List(ctx, apiClient, 1, 10)
    if err != nil {
        log.Fatalf("Ошибка при получении списка ссылок: %v", err)
    }

    for i, link := range links {
        fmt.Printf("Ссылка %d: %s (ID: %d)\n", i+1, link.ShortURL, link.ID)
    }
}
```

### Создание короткой ссылки

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/short_links"
)

func main() {
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")
    ctx := context.Background()

    newLink := &short_links.ShortLink{
        URL:           "https://example.com/product/123",
        EntityType:    "leads",
        EntityID:      456,
        UTMSource:     "newsletter",
        UTMMedium:     "email",
        UTMCampaign:   "spring_promo",
        UseInEmbedded: true,
    }

    createdLink, err := short_links.Create(ctx, apiClient, newLink)
    if err != nil {
        log.Fatalf("Ошибка при создании короткой ссылки: %v", err)
    }

    fmt.Printf("Создана короткая ссылка: %s\n", createdLink.ShortURL)
    fmt.Printf("ID: %d\nКлюч: %s\n", createdLink.ID, createdLink.Key)
}
```

### Получение статистики использования ссылки

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/short_links"
)

func main() {
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")
    ctx := context.Background()

    linkID := 123

    stats, err := short_links.GetStats(ctx, apiClient, linkID)
    if err != nil {
        log.Fatalf("Ошибка при получении статистики: %v", err)
    }

    fmt.Printf("Статистика для ссылки %s:\n", stats.ShortURL)
    fmt.Printf("Количество переходов: %d\n", stats.VisitCount)

    if stats.LastVisitAt > 0 {
        lastVisit := time.Unix(stats.LastVisitAt, 0)
        fmt.Printf("Последний переход: %s\n", lastVisit.Format("02.01.2006 15:04:05"))
    } else {
        fmt.Println("Ещё не было переходов по ссылке")
    }
}
```
