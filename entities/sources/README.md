# Источники сделок (Sources)

Модуль для работы с источниками сделок в amoCRM API.

## Содержание

- [Структуры данных](#структуры-данных)
- [Функции](#функции)
  - [Получение списка источников](#получение-списка-источников)
  - [Получение информации о конкретном источнике](#получение-информации-о-конкретном-источнике)
  - [Создание источника](#создание-источника)
  - [Обновление источника](#обновление-источника)
  - [Удаление источника](#удаление-источника)
  - [Установка источника по умолчанию](#установка-источника-по-умолчанию)
  - [Получение доступных сервисов](#получение-доступных-сервисов)
  - [Связывание источника с воронкой](#связывание-источника-с-воронкой)
- [Примеры использования](#примеры-использования)

## Структуры данных

### Source

Основная структура, представляющая источник сделок в amoCRM:

```go
type Source struct {
    ID            int         `json:"id,omitempty"`
    Name          string      `json:"name"`
    Type          string      `json:"type,omitempty"`
    Default       bool        `json:"default,omitempty"`
    CreatedAt     int64       `json:"created_at,omitempty"`
    UpdatedAt     int64       `json:"updated_at,omitempty"`
    Deleted       bool        `json:"deleted,omitempty"`
    EffectiveFrom int64       `json:"effective_from,omitempty"`
    EffectiveTo   int64       `json:"effective_to,omitempty"`
    Pipeline      *Pipeline   `json:"pipeline,omitempty"`
    Services      []Service   `json:"services,omitempty"`
    External      *External   `json:"external,omitempty"`
    Params        interface{} `json:"params,omitempty"`
}
```

### Pipeline

Структура, представляющая воронку, связанную с источником:

```go
type Pipeline struct {
    ID int `json:"id,omitempty"`
}
```

### Service

Структура, представляющая сервис для источника:

```go
type Service struct {
    ID   int    `json:"id,omitempty"`
    Name string `json:"name,omitempty"`
}
```

### External

Структура для внешних данных источника:

```go
type External struct {
    ID             string      `json:"id,omitempty"`
    Service        string      `json:"service,omitempty"`
    ExternalParams interface{} `json:"external_params,omitempty"`
}
```

## Функции

### Получение списка источников

```go
func GetSources(apiClient *client.Client, page, limit int, options ...WithOption) ([]Source, error)
```

Возвращает список источников сделок с поддержкой пагинации и фильтрации.

**Параметры:**
- `apiClient` - клиент API amoCRM
- `page` - номер страницы (начиная с 1)
- `limit` - количество элементов на странице
- `options` - дополнительные опции (например, фильтры)

**Фильтрация:**

Для фильтрации можно использовать функцию `WithFilter`:

```go
filter := map[string]string{
    "filter[type]": "calls",
}
sources, err := sources.GetSources(apiClient, 1, 50, sources.WithFilter(filter))
```

### Получение информации о конкретном источнике

```go
func GetSource(apiClient *client.Client, id int) (*Source, error)
```

Возвращает информацию о конкретном источнике по ID.

**Параметры:**
- `apiClient` - клиент API amoCRM
- `id` - ID источника

### Создание источника

```go
func CreateSource(apiClient *client.Client, sourceData *Source) (*Source, error)
```

Создаёт новый источник сделок и возвращает информацию о созданном источнике.

**Параметры:**
- `apiClient` - клиент API amoCRM
- `sourceData` - структура с данными для создания нового источника

### Обновление источника

```go
func UpdateSource(apiClient *client.Client, sourceData *Source) (*Source, error)
```

Обновляет существующий источник сделок и возвращает обновлённую информацию.

**Параметры:**
- `apiClient` - клиент API amoCRM
- `sourceData` - структура с данными для обновления источника (должен быть заполнен ID)

### Удаление источника

```go
func DeleteSource(apiClient *client.Client, id int) error
```

Удаляет источник по ID.

**Параметры:**
- `apiClient` - клиент API amoCRM
- `id` - ID источника для удаления

### Установка источника по умолчанию

```go
func SetSourceDefault(apiClient *client.Client, id int) (*Source, error)
```

Устанавливает источник как используемый по умолчанию.

**Параметры:**
- `apiClient` - клиент API amoCRM
- `id` - ID источника

### Получение доступных сервисов

```go
func GetSourceServices(apiClient *client.Client) ([]Service, error)
```

Возвращает список сервисов, доступных для источников сделок.

**Параметры:**
- `apiClient` - клиент API amoCRM

### Связывание источника с воронкой

```go
func LinkSourceToPipeline(apiClient *client.Client, sourceID, pipelineID int) (*Source, error)
```

Связывает источник с воронкой.

**Параметры:**
- `apiClient` - клиент API amoCRM
- `sourceID` - ID источника
- `pipelineID` - ID воронки

```go
func UnlinkSourceFromPipeline(apiClient *client.Client, sourceID, pipelineID int) (*Source, error)
```

Удаляет связь источника с воронкой.

**Параметры:**
- `apiClient` - клиент API amoCRM
- `sourceID` - ID источника
- `pipelineID` - ID воронки для удаления связи

## Примеры использования

### Получение списка источников

```go
package main

import (
    "fmt"
    "log"

    "github.com/chudno/amo_crm_sdk/auth"
    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/sources"
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

    // Получаем список источников
    sourcesList, err := sources.GetSources(apiClient, 1, 100)
    if err != nil {
        log.Fatalf("Ошибка при получении списка источников: %v", err)
    }

    // Выводим информацию о полученных источниках
    for i, source := range sourcesList {
        fmt.Printf("Источник %d: %s (ID: %d, Тип: %s)\n", i+1, source.Name, source.ID, source.Type)
        if source.Default {
            fmt.Println("  Используется по умолчанию")
        }
        
        if source.Pipeline != nil {
            fmt.Printf("  Связан с воронкой ID: %d\n", source.Pipeline.ID)
        }
        
        if len(source.Services) > 0 {
            fmt.Println("  Сервисы:")
            for _, service := range source.Services {
                fmt.Printf("    - %s (ID: %d)\n", service.Name, service.ID)
            }
        }
        
        fmt.Println()
    }
}
```

### Создание и настройка источника

```go
package main

import (
    "fmt"
    "log"

    "github.com/chudno/amo_crm_sdk/auth"
    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/sources"
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

    // Создаем новый источник
    newSource := &sources.Source{
        Name: "Сайт компании",
        Type: "form",
    }

    createdSource, err := sources.CreateSource(apiClient, newSource)
    if err != nil {
        log.Fatalf("Ошибка при создании источника: %v", err)
    }

    fmt.Printf("Создан новый источник: %s (ID: %d)\n", createdSource.Name, createdSource.ID)

    // Связываем источник с воронкой
    pipelineID := 12345 // ID вашей воронки
    linkedSource, err := sources.LinkSourceToPipeline(apiClient, createdSource.ID, pipelineID)
    if err != nil {
        log.Fatalf("Ошибка при связывании источника с воронкой: %v", err)
    }

    fmt.Printf("Источник %s связан с воронкой ID: %d\n", linkedSource.Name, linkedSource.Pipeline.ID)
    
    // Устанавливаем источник по умолчанию
    defaultSource, err := sources.SetSourceDefault(apiClient, createdSource.ID)
    if err != nil {
        log.Fatalf("Ошибка при установке источника по умолчанию: %v", err)
    }
    
    fmt.Printf("Источник %s установлен как используемый по умолчанию\n", defaultSource.Name)
}
```

### Работа с сервисами источников

```go
package main

import (
    "fmt"
    "log"

    "github.com/chudno/amo_crm_sdk/auth"
    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/sources"
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

    // Получаем список доступных сервисов
    servicesList, err := sources.GetSourceServices(apiClient)
    if err != nil {
        log.Fatalf("Ошибка при получении списка сервисов: %v", err)
    }

    fmt.Println("Доступные сервисы для источников:")
    for i, service := range servicesList {
        fmt.Printf("%d. %s (ID: %d)\n", i+1, service.Name, service.ID)
    }
    
    // Выбираем сервис для примера
    if len(servicesList) > 0 {
        fmt.Println("\nПример создания источника с выбранным сервисом:")
        
        // Создаем новый источник с выбранным сервисом
        newSource := &sources.Source{
            Name: "Интеграция с CRM",
            Type: "integration",
            Services: []sources.Service{
                {ID: servicesList[0].ID},
            },
        }
        
        createdSource, err := sources.CreateSource(apiClient, newSource)
        if err != nil {
            log.Fatalf("Ошибка при создании источника: %v", err)
        }
        
        fmt.Printf("Создан новый источник: %s (ID: %d) с сервисом %s\n", 
            createdSource.Name, 
            createdSource.ID,
            servicesList[0].Name)
    }
}
```
