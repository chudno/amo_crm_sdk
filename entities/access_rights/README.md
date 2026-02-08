# Модуль Access Rights (Права доступа)

Модуль предназначен для управления правами доступа в аккаунте amoCRM, включая получение списка прав доступа, создание новых прав, обновление существующих, а также управление правами для конкретных сущностей и пользователей.

## Содержание

- [Типы данных](#типы-данных)
- [Методы](#методы)
  - [Получение списка прав доступа](#получение-списка-прав-доступа)
  - [Получение информации о конкретном праве доступа](#получение-информации-о-конкретном-праве-доступа)
  - [Создание нового права доступа](#создание-нового-права-доступа)
  - [Обновление права доступа](#обновление-права-доступа)
  - [Удаление права доступа](#удаление-права-доступа)
  - [Обновление прав для конкретной сущности](#обновление-прав-для-конкретной-сущности)
  - [Добавление пользователей в право доступа](#добавление-пользователей-в-право-доступа)
  - [Удаление пользователей из права доступа](#удаление-пользователей-из-права-доступа)
- [Примеры использования](#примеры-использования)

## Типы данных

### Right

Основная структура для работы с правами доступа:

```go
type Right struct {
	ID          int              `json:"id,omitempty"`
	Name        string           `json:"name,omitempty"`
	Type        Type `json:"type,omitempty"`
	Rights      Rights           `json:"rights,omitempty"`
	CreatedBy   int              `json:"created_by,omitempty"`
	UpdatedBy   int              `json:"updated_by,omitempty"`
	CreatedAt   int              `json:"created_at,omitempty"`
	UpdatedAt   int              `json:"updated_at,omitempty"`
	AccountID   int              `json:"account_id,omitempty"`
	UserIDs     []int            `json:"user_ids,omitempty"`
	UserGroups  []UserGroup      `json:"_embedded.user_groups,omitempty"`
}
```

### Type

Тип права доступа определен как строковый тип с константами:

```go
type Type string

const (
	TypeGroup  Type = "group"
	TypeCustom Type = "custom"
)
```

### EntityType

Тип сущности для настройки прав доступа:

```go
type EntityType string

const (
	EntityLead       EntityType = "leads"
	EntityContact    EntityType = "contacts"
	EntityCompany    EntityType = "companies"
	EntityTask       EntityType = "tasks"
	EntityCustomer   EntityType = "customers"
	EntityCatalog    EntityType = "catalogs"
	EntityUnsorted   EntityType = "unsorted"
	EntityWidgets    EntityType = "widgets"
	EntityMails      EntityType = "mail"
	EntityChatWidget EntityType = "chat_widget"
)
```

### Rights

Структура для прав доступа к различным сущностям:

```go
type Rights struct {
	Leads       EntityRights `json:"leads,omitempty"`
	Contacts    EntityRights `json:"contacts,omitempty"`
	Companies   EntityRights `json:"companies,omitempty"`
	Tasks       EntityRights `json:"tasks,omitempty"`
	Customers   EntityRights `json:"customers,omitempty"`
	Catalogs    EntityRights `json:"catalogs,omitempty"`
	Unsorted    EntityRights `json:"unsorted,omitempty"`
	Widgets     EntityRights `json:"widgets,omitempty"`
	Mail        EntityRights `json:"mail,omitempty"`
	ChatWidget  EntityRights `json:"chat_widget,omitempty"`
	Settings    SettingsRights `json:"settings,omitempty"`
}
```

### EntityRights

Структура прав доступа к конкретной сущности:

```go
type EntityRights struct {
	View   bool `json:"view,omitempty"`
	Edit   bool `json:"edit,omitempty"`
	Add    bool `json:"add,omitempty"`
	Delete bool `json:"delete,omitempty"`
	Export bool `json:"export,omitempty"`
}
```

### SettingsRights

Структура прав доступа к настройкам:

```go
type SettingsRights struct {
	View bool `json:"view,omitempty"`
	Edit bool `json:"edit,omitempty"`
}
```

### UserGroup

Структура для группы пользователей:

```go
type UserGroup struct {
	ID      int    `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	UserIDs []int  `json:"user_ids,omitempty"`
}
```

## Методы

### Получение списка прав доступа

Метод `List` позволяет получить список прав доступа в аккаунте с возможностью фильтрации и пагинации.

```go
func List(ctx context.Context, apiClient *client.Client, page, limit int, options ...WithOption) ([]Right, error)
```

#### Параметры

- `apiClient` - клиент API amoCRM
- `page` - номер страницы (начиная с 1)
- `limit` - количество прав доступа на странице
- `options` - опциональные параметры (фильтры)

#### Доступные опции

- `WithFilter` - общая фильтрация по параметрам
- `WithType` - фильтрация по типу права доступа

#### Пример использования

```go
// Получение всех прав доступа типа group
accessRights, err := access_rights.List(ctx, apiClient, 1, 50, access_rights.WithType(access_rights.TypeGroup))
if err != nil {
    log.Fatal(err)
}

// Вывод информации о правах доступа
for _, right := range accessRights {
    fmt.Printf("ID: %d, Название: %s, Тип: %s\n", right.ID, right.Name, right.Type)
}
```

### Получение информации о конкретном праве доступа

Метод `Get` позволяет получить подробную информацию о конкретном праве доступа по его ID.

```go
func Get(ctx context.Context, apiClient *client.Client, accessRightID int) (*Right, error)
```

#### Параметры

- `apiClient` - клиент API amoCRM
- `accessRightID` - ID права доступа

#### Пример использования

```go
// Получение информации о праве доступа
accessRight, err := access_rights.Get(ctx, apiClient, 123)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Название: %s, Тип: %s\n", accessRight.Name, accessRight.Type)
fmt.Printf("Просмотр сделок: %v, Редактирование сделок: %v\n", 
    accessRight.Rights.Leads.View, accessRight.Rights.Leads.Edit)
```

### Создание нового права доступа

Метод `Create` позволяет создать новое право доступа в аккаунте.

```go
func Create(ctx context.Context, apiClient *client.Client, accessRight *Right) (*Right, error)
```

#### Параметры

- `apiClient` - клиент API amoCRM
- `accessRight` - структура с данными нового права доступа

#### Пример использования

```go
// Создание нового права доступа
newRight := &access_rights.Right{
    Name: "Менеджеры продаж",
    Type: access_rights.TypeGroup,
    Rights: access_rights.Rights{
        Leads: access_rights.EntityRights{
            View: true,
            Edit: true,
            Add: true,
        },
        Contacts: access_rights.EntityRights{
            View: true,
            Edit: true,
            Add: true,
        },
    },
    UserIDs: []int{123, 456},
}

createdRight, err := access_rights.Create(ctx, apiClient, newRight)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Создано право доступа: ID=%d, Название=%s\n", createdRight.ID, createdRight.Name)
```

### Обновление права доступа

Метод `Update` позволяет обновить существующее право доступа.

```go
func Update(ctx context.Context, apiClient *client.Client, accessRight *Right) (*Right, error)
```

#### Параметры

- `apiClient` - клиент API amoCRM
- `accessRight` - структура с данными для обновления

#### Пример использования

```go
// Обновление существующего права доступа
updateRight := &access_rights.Right{
    ID: 123,
    Name: "Менеджеры продаж (обновлено)",
    Rights: access_rights.Rights{
        Leads: access_rights.EntityRights{
            View: true,
            Edit: true,
            Add: true,
            Delete: true, // Добавляем право на удаление
        },
    },
    UserIDs: []int{123, 456, 789}, // Добавляем нового пользователя
}

updatedRight, err := access_rights.Update(ctx, apiClient, updateRight)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Обновлено право доступа: ID=%d, Название=%s\n", updatedRight.ID, updatedRight.Name)
```

### Удаление права доступа

Метод `Delete` позволяет удалить право доступа из аккаунта.

```go
func Delete(ctx context.Context, apiClient *client.Client, accessRightID int) error
```

#### Параметры

- `apiClient` - клиент API amoCRM
- `accessRightID` - ID права доступа для удаления

#### Пример использования

```go
// Удаление права доступа
err := access_rights.Delete(ctx, apiClient, 123)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Право доступа успешно удалено")
```

### Обновление прав для конкретной сущности

Метод `SetEntityRights` позволяет обновить права доступа к конкретной сущности.

```go
func SetEntityRights(ctx context.Context, apiClient *client.Client, accessRightID int, entityType EntityType, rights EntityRights) (*Right, error)
```

#### Параметры

- `apiClient` - клиент API amoCRM
- `accessRightID` - ID права доступа
- `entityType` - тип сущности
- `rights` - новые права для сущности

#### Пример использования

```go
// Обновление прав для сделок
entityRights := access_rights.EntityRights{
    View: true,
    Edit: true,
    Add: true,
    Delete: true,
    Export: true,
}

updatedRight, err := access_rights.SetEntityRights(ctx, apiClient, 123, access_rights.EntityLead, entityRights)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Обновлены права для сделок: Просмотр=%v, Редактирование=%v, Удаление=%v\n",
    updatedRight.Rights.Leads.View, updatedRight.Rights.Leads.Edit, updatedRight.Rights.Leads.Delete)
```

### Добавление пользователей в право доступа

Метод `AddUsers` позволяет добавить пользователей в существующее право доступа.

```go
func AddUsers(ctx context.Context, apiClient *client.Client, accessRightID int, userIDs []int) (*Right, error)
```

#### Параметры

- `apiClient` - клиент API amoCRM
- `accessRightID` - ID права доступа
- `userIDs` - массив ID пользователей для добавления

#### Пример использования

```go
// Добавление пользователей в право доступа
userIDs := []int{789, 101}
updatedRight, err := access_rights.AddUsers(ctx, apiClient, 123, userIDs)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Пользователи добавлены. Всего пользователей: %d\n", len(updatedRight.UserIDs))
```

### Удаление пользователей из права доступа

Метод `RemoveUsers` позволяет удалить пользователей из существующего права доступа.

```go
func RemoveUsers(ctx context.Context, apiClient *client.Client, accessRightID int, userIDs []int) (*Right, error)
```

#### Параметры

- `apiClient` - клиент API amoCRM
- `accessRightID` - ID права доступа
- `userIDs` - массив ID пользователей для удаления

#### Пример использования

```go
// Удаление пользователей из права доступа
userIDs := []int{789, 101}
updatedRight, err := access_rights.RemoveUsers(ctx, apiClient, 123, userIDs)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Пользователи удалены. Всего пользователей: %d\n", len(updatedRight.UserIDs))
```

## Примеры использования

### Создание права доступа и настройка прав для разных сущностей

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/access_rights"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // Создаем контекст
    ctx := context.Background()

    // Создаем новое право доступа
    newRight := &access_rights.Right{
        Name: "Менеджеры продаж",
        Type: access_rights.TypeGroup,
        Rights: access_rights.Rights{
            Leads: access_rights.EntityRights{
                View: true,
                Edit: true,
                Add: true,
            },
            Contacts: access_rights.EntityRights{
                View: true,
                Edit: true,
                Add: true,
            },
        },
        UserIDs: []int{123, 456},
    }

    createdRight, err := access_rights.Create(ctx, apiClient, newRight)
    if err != nil {
        log.Fatalf("Ошибка при создании права доступа: %v", err)
    }
    
    fmt.Printf("Создано право доступа: ID=%d, Название=%s\n", createdRight.ID, createdRight.Name)
    
    // Добавляем право удаления сделок
    leadsRights := access_rights.EntityRights{
        View: true,
        Edit: true,
        Add: true,
        Delete: true,
        Export: true,
    }
    
    updatedRight, err := access_rights.SetEntityRights(ctx, apiClient, createdRight.ID, access_rights.EntityLead, leadsRights)
    if err != nil {
        log.Fatalf("Ошибка при обновлении прав для сделок: %v", err)
    }
    
    fmt.Printf("Обновлены права для сделок: Удаление=%v\n", updatedRight.Rights.Leads.Delete)
    
    // Добавляем пользователей
    userIDs := []int{789, 101}
    updatedRight, err = access_rights.AddUsers(ctx, apiClient, createdRight.ID, userIDs)
    if err != nil {
        log.Fatalf("Ошибка при добавлении пользователей: %v", err)
    }
    
    fmt.Printf("Пользователи добавлены. Всего пользователей: %d\n", len(updatedRight.UserIDs))
}
```

### Получение и фильтрация прав доступа

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/access_rights"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // Создаем контекст
    ctx := context.Background()

    // Получаем права доступа типа group
    accessRights, err := access_rights.List(ctx, apiClient, 1, 50, access_rights.WithType(access_rights.TypeGroup))
    if err != nil {
        log.Fatalf("Ошибка при получении прав доступа: %v", err)
    }

    // Выводим информацию о правах доступа
    fmt.Printf("Найдено прав доступа: %d\n", len(accessRights))
    for _, right := range accessRights {
        fmt.Printf("ID: %d, Название: %s, Тип: %s\n", right.ID, right.Name, right.Type)
        
        // Выводим информацию о правах для сделок
        if right.Rights.Leads.View {
            fmt.Printf("  Права для сделок: Просмотр=%v, Редактирование=%v, Добавление=%v, Удаление=%v\n",
                right.Rights.Leads.View, right.Rights.Leads.Edit, right.Rights.Leads.Add, right.Rights.Leads.Delete)
        }
        
        // Выводим информацию о пользователях
        if len(right.UserIDs) > 0 {
            fmt.Printf("  Пользователи: %v\n", right.UserIDs)
        }
    }
}
```

### Обновление и удаление права доступа

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/access_rights"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // Создаем контекст
    ctx := context.Background()

    // ID права доступа для обновления
    accessRightID := 123

    // Получаем информацию о праве доступа
    accessRight, err := access_rights.Get(ctx, apiClient, accessRightID)
    if err != nil {
        log.Fatalf("Ошибка при получении информации о праве доступа: %v", err)
    }
    
    fmt.Printf("Получено право доступа: ID=%d, Название=%s\n", accessRight.ID, accessRight.Name)
    
    // Обновляем право доступа
    accessRight.Name = "Обновленное название"
    accessRight.Rights.Contacts.Add = true
    accessRight.Rights.Contacts.Edit = true
    
    updatedRight, err := access_rights.Update(ctx, apiClient, accessRight)
    if err != nil {
        log.Fatalf("Ошибка при обновлении права доступа: %v", err)
    }
    
    fmt.Printf("Право доступа обновлено: ID=%d, Название=%s\n", updatedRight.ID, updatedRight.Name)
    
    // Удаляем право доступа
    err = access_rights.Delete(ctx, apiClient, accessRightID)
    if err != nil {
        log.Fatalf("Ошибка при удалении права доступа: %v", err)
    }
    
    fmt.Println("Право доступа успешно удалено")
}
```
