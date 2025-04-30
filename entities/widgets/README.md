# Модуль Widgets (Виджеты)

Модуль предназначен для управления виджетами в аккаунте amoCRM, включая получение списка виджетов, установку новых виджетов из маркетплейса, настройку и удаление виджетов.

## Содержание

- [Типы данных](#типы-данных)
- [Методы](#методы)
  - [Получение списка виджетов](#получение-списка-виджетов)
  - [Получение информации о виджете](#получение-информации-о-виджете)
  - [Установка виджета](#установка-виджета)
  - [Обновление настроек виджета](#обновление-настроек-виджета)
  - [Удаление виджета](#удаление-виджета)
  - [Получение списка виджетов из маркетплейса](#получение-списка-виджетов-из-маркетплейса)
  - [Изменение статуса виджета](#изменение-статуса-виджета)
  - [Массовая установка виджетов](#массовая-установка-виджетов)
  - [Массовое удаление виджетов](#массовое-удаление-виджетов)
- [Примеры использования](#примеры-использования)

## Типы данных

### Widget

Основная структура для работы с виджетами:

```go
type Widget struct {
	ID            int          `json:"id,omitempty"`
	Name          string       `json:"name,omitempty"`
	Code          string       `json:"code,omitempty"`
	Type          WidgetType   `json:"type,omitempty"`
	Status        WidgetStatus `json:"status,omitempty"`
	CreatedBy     int          `json:"created_by,omitempty"`
	UpdatedBy     int          `json:"updated_by,omitempty"`
	CreatedAt     int          `json:"created_at,omitempty"`
	UpdatedAt     int          `json:"updated_at,omitempty"`
	AccountID     int          `json:"account_id,omitempty"`
	Settings      interface{}  `json:"settings,omitempty"`
	Rights        *Rights      `json:"rights,omitempty"`
	Marketplace   *Marketplace `json:"marketplace,omitempty"`
	IsConfigured  bool         `json:"is_configured,omitempty"`
	VerifiedAt    int          `json:"verified_at,omitempty"`
	MainVersion   string       `json:"main_version,omitempty"`
	CurrentVersion string      `json:"current_version,omitempty"`
	IsDeleted     bool         `json:"is_deleted,omitempty"`
}
```

### WidgetType

Тип виджета определен как строковый тип с константами:

```go
type WidgetType string

const (
	WidgetTypeIntercom     WidgetType = "intercom"
	WidgetTypeJivosite     WidgetType = "jivosite"
	WidgetTypeCallback     WidgetType = "callback"
	WidgetTypePipeline     WidgetType = "pipeline"
	WidgetTypeMailchimp    WidgetType = "mailchimp"
	WidgetTypeCustom       WidgetType = "custom"
	WidgetTypeGoalMeter    WidgetType = "goal_meter"
	WidgetTypeDigitalPipeline WidgetType = "digital_pipeline"
	WidgetTypeSupport      WidgetType = "support"
	WidgetTypeIpTelephony  WidgetType = "ip_telephony"
	WidgetTypePayment      WidgetType = "payment"
	WidgetTypeAmoButtons   WidgetType = "amo_buttons"
	WidgetTypeEmailSubscription WidgetType = "email_subscription"
)
```

### WidgetStatus

Статус виджета определен как строковый тип с константами:

```go
type WidgetStatus string

const (
	WidgetStatusInstalled WidgetStatus = "installed"
	WidgetStatusDemo      WidgetStatus = "demo"
	WidgetStatusInactive  WidgetStatus = "inactive"
)
```

## Методы

### Получение списка виджетов

Метод `GetWidgets` позволяет получить список виджетов в аккаунте с возможностью фильтрации и пагинации.

```go
func GetWidgets(apiClient *client.Client, page, limit int, options ...WithOption) ([]Widget, error)
```

#### Параметры

- `apiClient` - клиент API amoCRM
- `page` - номер страницы (начиная с 1)
- `limit` - количество виджетов на странице
- `options` - опциональные параметры (фильтры)

#### Доступные опции

- `WithFilter(filter map[string]string)` - добавляет произвольную фильтрацию
- `WithWidgetTypes(types []WidgetType)` - фильтрация по типам виджетов

#### Пример использования

```go
// Получение всех виджетов
widgets, err := widgets.GetWidgets(apiClient, 1, 50)

// Фильтрация по типам виджетов
types := []widgets.WidgetType{widgets.WidgetTypeIntercom, widgets.WidgetTypeCallback}
widgets, err := widgets.GetWidgets(apiClient, 1, 50, widgets.WithWidgetTypes(types))
```

### Получение информации о виджете

Метод `GetWidget` позволяет получить детальную информацию о конкретном виджете по его ID.

```go
func GetWidget(apiClient *client.Client, widgetID int) (*Widget, error)
```

#### Параметры

- `apiClient` - клиент API amoCRM
- `widgetID` - ID виджета

#### Пример использования

```go
// Получение информации о виджете с ID 123
widget, err := widgets.GetWidget(apiClient, 123)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Виджет: %s (тип: %s)\n", widget.Name, widget.Type)
```

### Установка виджета

Метод `InstallWidget` позволяет установить виджет из маркетплейса amoCRM по его коду.

```go
func InstallWidget(apiClient *client.Client, code string) (*Widget, error)
```

#### Параметры

- `apiClient` - клиент API amoCRM
- `code` - код виджета в маркетплейсе

#### Пример использования

```go
// Установка виджета Intercom
widget, err := widgets.InstallWidget(apiClient, "intercom")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Установлен виджет: %s (ID: %d)\n", widget.Name, widget.ID)
```

### Обновление настроек виджета

Метод `UpdateWidgetSettings` позволяет обновить настройки виджета.

```go
func UpdateWidgetSettings(apiClient *client.Client, widgetID int, settings interface{}) (*Widget, error)
```

#### Параметры

- `apiClient` - клиент API amoCRM
- `widgetID` - ID виджета
- `settings` - настройки виджета (объект или карта)

#### Пример использования

```go
// Обновление настроек виджета
settings := map[string]interface{}{
    "api_key": "abc123",
    "active": true,
}
widget, err := widgets.UpdateWidgetSettings(apiClient, 123, settings)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Настройки виджета обновлены: %s\n", widget.Name)
```

### Удаление виджета

Метод `DeleteWidget` позволяет удалить виджет из аккаунта.

```go
func DeleteWidget(apiClient *client.Client, widgetID int) error
```

#### Параметры

- `apiClient` - клиент API amoCRM
- `widgetID` - ID виджета

#### Пример использования

```go
// Удаление виджета
err := widgets.DeleteWidget(apiClient, 123)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Виджет успешно удален")
```

### Получение списка виджетов из маркетплейса

Метод `GetMarketplaceWidgets` позволяет получить список доступных виджетов в маркетплейсе amoCRM.

```go
func GetMarketplaceWidgets(apiClient *client.Client, page, limit int, options ...WithOption) ([]MarketplaceWidget, error)
```

#### Параметры

- `apiClient` - клиент API amoCRM
- `page` - номер страницы (начиная с 1)
- `limit` - количество виджетов на странице
- `options` - опциональные параметры (фильтры)

#### Доступные опции

- `WithFilter` - общая фильтрация по параметрам
- `WithCategory` - фильтрация по категории виджетов

#### Пример использования

```go
// Получение списка виджетов из маркетплейса в категории "Коммуникации" (ID: 1)
marketWidgets, err := widgets.GetMarketplaceWidgets(apiClient, 1, 50, widgets.WithCategory(1))
if err != nil {
    log.Fatal(err)
}

// Выводим список доступных виджетов
for _, widget := range marketWidgets {
    fmt.Printf("Виджет: %s (Код: %s, Установлен: %v)\n", 
        widget.Name, widget.Code, widget.Installed)
}
```

### Изменение статуса виджета

Метод `SetWidgetStatus` позволяет изменить статус виджета.

```go
func SetWidgetStatus(apiClient *client.Client, widgetID int, status WidgetStatus) (*Widget, error)
```

#### Параметры

- `apiClient` - клиент API amoCRM
- `widgetID` - ID виджета
- `status` - новый статус виджета (installed, demo или inactive)

#### Пример использования

```go
// Деактивация виджета
widget, err := widgets.SetWidgetStatus(apiClient, 123, widgets.WidgetStatusInactive)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Статус виджета %s изменен на: %s\n", widget.Name, widget.Status)

// Активация виджета
widget, err = widgets.SetWidgetStatus(apiClient, 123, widgets.WidgetStatusInstalled)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Статус виджета %s изменен на: %s\n", widget.Name, widget.Status)
```

### Массовая установка виджетов

Метод `BulkInstallWidgets` позволяет установить несколько виджетов одновременно.

```go
func BulkInstallWidgets(apiClient *client.Client, codes []string) ([]Widget, error)
```

#### Параметры

- `apiClient` - клиент API amoCRM
- `codes` - список кодов виджетов для установки

#### Пример использования

```go
// Массовая установка нескольких виджетов
codes := []string{"intercom", "callback", "jivosite"}
installedWidgets, err := widgets.BulkInstallWidgets(apiClient, codes)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Установлено виджетов: %d\n", len(installedWidgets))
for _, widget := range installedWidgets {
    fmt.Printf("Виджет: %s (ID: %d)\n", widget.Name, widget.ID)
}
```

### Массовое удаление виджетов

Метод `BulkDeleteWidgets` позволяет удалить несколько виджетов одновременно.

```go
func BulkDeleteWidgets(apiClient *client.Client, widgetIDs []int) error
```

#### Параметры

- `apiClient` - клиент API amoCRM
- `widgetIDs` - список ID виджетов для удаления

#### Пример использования

```go
// Массовое удаление ненужных виджетов
widgetIDs := []int{123, 456, 789}
err := widgets.BulkDeleteWidgets(apiClient, widgetIDs)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Виджеты успешно удалены\n")
```

## Примеры использования

### Получение списка виджетов определенного типа

```go
package main

import (
    "fmt"
    "log"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/widgets"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // Получаем виджеты типа "callback" и "intercom"
    types := []widgets.WidgetType{widgets.WidgetTypeCallback, widgets.WidgetTypeIntercom}
    widgetsList, err := widgets.GetWidgets(apiClient, 1, 50, widgets.WithWidgetTypes(types))
    if err != nil {
        log.Fatalf("Ошибка при получении виджетов: %v", err)
    }

    // Выводим информацию о виджетах
    fmt.Printf("Найдено виджетов: %d\n", len(widgetsList))
    for _, widget := range widgetsList {
        fmt.Printf("ID: %d, Название: %s, Тип: %s, Статус: %s\n", 
            widget.ID, widget.Name, widget.Type, widget.Status)
    }
}
```

### Установка и настройка нового виджета

```go
package main

import (
    "fmt"
    "log"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/widgets"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // Устанавливаем виджет
    widget, err := widgets.InstallWidget(apiClient, "intercom")
    if err != nil {
        log.Fatalf("Ошибка при установке виджета: %v", err)
    }
    fmt.Printf("Виджет установлен: %s (ID: %d)\n", widget.Name, widget.ID)

    // Настраиваем виджет
    settings := map[string]interface{}{
        "api_key": "your-intercom-api-key",
        "active": true,
        "workspace_id": "your-workspace-id",
    }
    
    updatedWidget, err := widgets.UpdateWidgetSettings(apiClient, widget.ID, settings)
    if err != nil {
        log.Fatalf("Ошибка при настройке виджета: %v", err)
    }
    
    fmt.Printf("Виджет настроен: %s (configured: %v)\n", 
        updatedWidget.Name, updatedWidget.IsConfigured)
}
```

### Получение информации о виджете и его удаление

```go
package main

import (
    "fmt"
    "log"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/widgets"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // ID виджета
    widgetID := 123
    
    // Получаем информацию о виджете
    widget, err := widgets.GetWidget(apiClient, widgetID)
    if err != nil {
        log.Fatalf("Ошибка при получении информации о виджете: %v", err)
    }
    
    fmt.Printf("Виджет: %s (ID: %d, Тип: %s)\n", 
        widget.Name, widget.ID, widget.Type)
    
    // Удаляем виджет
    err = widgets.DeleteWidget(apiClient, widgetID)
    if err != nil {
        log.Fatalf("Ошибка при удалении виджета: %v", err)
    }
    
    fmt.Printf("Виджет %s (ID: %d) успешно удален\n", widget.Name, widget.ID)
}
```
