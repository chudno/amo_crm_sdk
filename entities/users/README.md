# Модуль Пользователи

Модуль `users` предоставляет функциональность для работы с пользователями (сотрудниками) в amoCRM.

## Содержание

- [Основные функции](#основные-функции)
- [Получение пользователя](#получение-пользователя)
- [Получение списка пользователей](#получение-списка-пользователей)
- [Права доступа пользователя](#права-доступа-пользователя)
- [Назначение ответственных](#назначение-ответственных)

## Основные функции

| Функция | Описание |
|---------|----------|
| `Get` | Получение пользователя по ID |
| `List` | Получение списка пользователей с пагинацией |
| `GetCurrent` | Получение информации о текущем пользователе |

## Получение пользователя

```go
import (
    "context"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/users"
)

// Инициализация клиента
apiClient := client.NewClient("https://your-domain.amocrm.ru", "your_access_token")

// Создаем контекст
ctx := context.Background()

// Получение пользователя по ID
userID := 12345
user, err := users.Get(ctx, apiClient, userID)
if err != nil {
    // Обработка ошибки
}

// Вывод информации о пользователе
fmt.Printf("Имя: %s\n", user.Name)
fmt.Printf("Email: %s\n", user.Email)
if user.Rights.IsAdmin {
    fmt.Println("Пользователь является администратором")
}
if user.Rights.Leads != nil {
    fmt.Printf("Доступ к сделкам (просмотр): %s\n", user.Rights.Leads.View)
}
```

## Получение списка пользователей

```go
// Получение первых 50 пользователей
usersList, err := users.List(ctx, apiClient, 50, 1)
if err != nil {
    // Обработка ошибки
}

for _, user := range usersList {
    fmt.Printf("ID: %d, Имя: %s, Email: %s\n", user.ID, user.Name, user.Email)
}
```

## Права доступа пользователя

Структура `Rights` описывает права пользователя в системе amoCRM. Все поля типизированы, что позволяет
обращаться к конкретным правам без ручного разбора JSON.

### Структура `Rights`

| Поле | Тип | Описание |
|------|-----|----------|
| `Leads` | `*EntityRights` | Права на сделки |
| `Contacts` | `*EntityRights` | Права на контакты |
| `Companies` | `*EntityRights` | Права на компании |
| `Tasks` | `*EntityRights` | Права на задачи |
| `MailAccess` | `bool` | Доступ к почте |
| `CatalogAccess` | `bool` | Доступ к каталогам |
| `StatusRights` | `[]StatusRight` | Права на конкретные статусы воронок |
| `IsAdmin` | `bool` | Является ли администратором |
| `IsManager` | `bool` | Является ли менеджером |
| `IsFree` | `bool` | Свободный пользователь |
| `IsActive` | `bool` | Активен ли пользователь |
| `GroupID` | `*int` | ID группы пользователя |
| `RoleID` | `*int` | ID роли пользователя |

### Структура `EntityRights`

Описывает набор прав на действия с конкретной сущностью. Каждое поле имеет тип `AccessLevel` (псевдоним `string`).

| Поле | Описание |
|------|----------|
| `View` | Просмотр |
| `Edit` | Редактирование |
| `Add` | Создание |
| `Delete` | Удаление |
| `Export` | Экспорт |

### Константы `AccessLevel`

| Константа | Значение | Описание |
|-----------|----------|----------|
| `AccessAll` | `"A"` | Полный доступ |
| `AccessGroup` | `"G"` | Доступ в пределах группы |
| `AccessOwn` | `"M"` | Только свои |
| `AccessDeny` | `"D"` | Запрещено |

### Структура `StatusRight`

Описывает права на конкретный статус воронки.

| Поле | Тип | Описание |
|------|-----|----------|
| `EntityType` | `string` | Тип сущности (например, `"leads"`) |
| `PipelineID` | `int` | ID воронки |
| `StatusID` | `int` | ID статуса |
| `Rights` | `*EntityRights` | Права на действия в данном статусе |

### Пример работы с правами

```go
import (
    "context"
    "fmt"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/users"
)

// Инициализация клиента
apiClient := client.NewClient("https://your-domain.amocrm.ru", "your_access_token")

// Создаем контекст
ctx := context.Background()

// Получаем пользователя
user, err := users.Get(ctx, apiClient, 12345)
if err != nil {
    // Обработка ошибки
}

// Проверяем роль пользователя
if user.Rights.IsAdmin {
    fmt.Println("Пользователь является администратором")
}

// Проверяем права на сделки
if user.Rights.Leads != nil {
    fmt.Printf("Просмотр сделок: %s\n", user.Rights.Leads.View)
    fmt.Printf("Редактирование сделок: %s\n", user.Rights.Leads.Edit)

    // Сравнение с константами уровня доступа
    if user.Rights.Leads.View == users.AccessAll {
        fmt.Println("У пользователя полный доступ к просмотру сделок")
    } else if user.Rights.Leads.View == users.AccessGroup {
        fmt.Println("У пользователя доступ к сделкам своей группы")
    } else if user.Rights.Leads.View == users.AccessOwn {
        fmt.Println("У пользователя доступ только к своим сделкам")
    } else if user.Rights.Leads.View == users.AccessDeny {
        fmt.Println("Просмотр сделок запрещён")
    }
}

// Проверяем права на контакты
if user.Rights.Contacts != nil && user.Rights.Contacts.Export == users.AccessDeny {
    fmt.Println("Экспорт контактов запрещён для этого пользователя")
}

// Проверяем права на статусы воронок
for _, sr := range user.Rights.StatusRights {
    fmt.Printf("Воронка %d, статус %d: просмотр=%s\n",
        sr.PipelineID, sr.StatusID, sr.Rights.View)
}

// Проверяем группу и роль
if user.Rights.GroupID != nil {
    fmt.Printf("Группа пользователя: %d\n", *user.Rights.GroupID)
}
if user.Rights.RoleID != nil {
    fmt.Printf("Роль пользователя: %d\n", *user.Rights.RoleID)
}
```

## Назначение ответственных

Для назначения ответственного пользователя за сущность, используйте поле `ResponsibleUserID` в соответствующих структурах:

```go
// Назначение ответственного за контакт
import "github.com/chudno/amo_crm_sdk/entities/contacts"

contact := &contacts.Contact{
    Name: "Иван Иванов",
    ResponsibleUserID: 12345, // ID ответственного менеджера
}

// Назначение ответственного за лид
import "github.com/chudno/amo_crm_sdk/entities/leads"

lead := &leads.Lead{
    Name:              "Тестовый лид",
    ResponsibleUserID: 12345, // ID ответственного менеджера
}

// Назначение ответственного за задачу
import "github.com/chudno/amo_crm_sdk/entities/tasks"

task := &tasks.Task{
    Text: "Позвонить клиенту",
    ResponsibleUserID: 12345, // ID ответственного менеджера
}
```

Для смены ответственного у существующей сущности используйте функцию `Update` соответствующего модуля (например, `leads.Update`).
