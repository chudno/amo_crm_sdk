# Модуль Задачи

Модуль `tasks` предоставляет функциональность для работы с задачами в amoCRM.

## Содержание

- [Основные функции](#основные-функции)
- [Создание задачи](#создание-задачи)
- [Получение задачи](#получение-задачи)
- [Получение списка задач](#получение-списка-задач)
- [Обновление задачи](#обновление-задачи)
- [Завершение задачи](#завершение-задачи)
- [Связь с другими сущностями](#связь-с-другими-сущностями)
- [Константы типов сущностей](#константы-типов-сущностей)

## Основные функции

| Функция | Описание |
|---------|----------|
| `Create` | Создание новой задачи |
| `Get` | Получение задачи по ID |
| `List` | Получение списка задач с фильтрацией |
| `Update` | Обновление существующей задачи |
| `Complete` | Завершение задачи |
| `CreateForEntity` | Создание задачи, привязанной к сущности |

## Создание задачи

```go
import (
    "context"
    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/tasks"
    "time"
)

// Инициализация клиента
apiClient := client.NewClient("https://your-domain.amocrm.ru", "your_access_token")
ctx := context.Background()

// Создание новой задачи
completionTime := time.Now().Add(24 * time.Hour)
newTask := &tasks.Task{
    TaskTypeID:        1,
    Text:              "Перезвонить клиенту",
    CompleteTill:      completionTime.Unix(),
    ResponsibleUserID: 12345,
    EntityID:          67890,
    EntityType:        tasks.EntityTypeContact,
}

// Сохранение задачи
createdTask, err := tasks.Create(ctx, apiClient, newTask)
if err != nil {
    // Обработка ошибки
}
```

## Получение задачи

```go
// Получение задачи по ID
taskID := 12345
task, err := tasks.Get(ctx, apiClient, taskID)
if err != nil {
    // Обработка ошибки
}
```

## Получение списка задач

```go
// Получение первых 50 задач
tasksList, err := tasks.List(ctx, apiClient, 50, 1, nil)
if err != nil {
    // Обработка ошибки
}

// Получение задач с фильтрацией
filter := map[string]any{
    "filter[is_completed]":        0,
    "filter[responsible_user_id]": 12345,
}
filteredTasks, err := tasks.List(ctx, apiClient, 50, 1, filter)
```

## Обновление задачи

```go
// Обновление существующей задачи
task.Text = "Срочно перезвонить клиенту"

// Перенос срока задачи
newCompletionTime := time.Now().Add(12 * time.Hour)
task.CompleteTill = newCompletionTime.Unix()

updatedTask, err := tasks.Update(ctx, apiClient, task)
if err != nil {
    // Обработка ошибки
}
```

## Завершение задачи

```go
// Завершение задачи
result := "Клиент согласился на встречу" // Результат выполнения задачи
completedTask, err := tasks.Complete(ctx, apiClient, taskID, result)
if err != nil {
    // Обработка ошибки
}
```

## Связь с другими сущностями

Задачи в amoCRM всегда связаны с определенной сущностью (контактом, сделкой, компанией и т.д.). При создании задачи необходимо указать тип сущности и её ID:

```go
// Создание задачи для контакта через CreateForEntity
createdTask, err := tasks.CreateForEntity(
    ctx, apiClient,
    tasks.EntityTypeContact, 67890,
    1,
    "Позвонить контакту",
    time.Now().Add(24*time.Hour),
    12345,
)

// Создание задачи для сделки
leadTask, err := tasks.CreateForEntity(
    ctx, apiClient,
    tasks.EntityTypeLead, 12345,
    1,
    "Подготовить коммерческое предложение",
    time.Now().Add(24*time.Hour),
    12345,
)
```

## Константы типов сущностей

| Константа | Значение | Описание |
|-----------|----------|----------|
| `tasks.EntityTypeLead` | "leads" | Сделка |
| `tasks.EntityTypeContact` | "contacts" | Контакт |
| `tasks.EntityTypeCompany` | "companies" | Компания |
| `tasks.EntityTypeCustomer` | "customers" | Покупатель |
