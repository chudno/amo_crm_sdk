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
- [Типы задач](#типы-задач)

## Основные функции

| Функция | Описание |
|---------|----------|
| `Create` | Создание новой задачи |
| `Get` | Получение задачи по ID |
| `List` | Получение списка задач с фильтрацией |
| `Update` | Обновление существующей задачи |
| `Complete` | Завершение задачи |
| `Delete` | Удаление задачи |

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
completionTime := time.Now().Add(24 * time.Hour) // задача на завтра
newTask := &tasks.Task{
    TaskType: tasks.TypeCall, // Тип задачи - звонок
    Text: "Перезвонить клиенту",
    CompleteTill: completionTime.Unix(), // Unix timestamp
    ResponsibleUserID: 12345, // ID ответственного менеджера
    EntityID: 67890, // ID связанной сущности (например, контакта)
    EntityType: tasks.EntityTypeContact, // Тип связанной сущности
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
tasksList, err := tasks.List(ctx, apiClient, 1, 50)
if err != nil {
    // Обработка ошибки
}

// Получение задач с фильтрацией
filter := map[string]string{
    "filter[task_type]": "1", // Фильтр по типу задачи (1 - звонок)
    "filter[is_completed]": "0", // Только незавершенные задачи
    "filter[responsible_user_id]": "12345", // Задачи конкретного менеджера
}
filteredTasks, err := tasks.List(ctx, apiClient, 1, 50, filter)
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
// Создание задачи, связанной с контактом
task := &tasks.Task{
    Text: "Позвонить контакту",
    TaskType: tasks.TypeCall,
    CompleteTill: time.Now().Add(24 * time.Hour).Unix(),
    EntityID: 67890, // ID контакта
    EntityType: tasks.EntityTypeContact, // Тип сущности - контакт
}

// Создание задачи, связанной со сделкой
task := &tasks.Task{
    Text: "Подготовить коммерческое предложение",
    TaskType: tasks.TypeTask,
    CompleteTill: time.Now().Add(24 * time.Hour).Unix(),
    EntityID: 12345, // ID сделки
    EntityType: tasks.EntityTypeLead, // Тип сущности - сделка
}
```

## Типы задач

Модуль `tasks` предоставляет константы для типов задач:

| Константа | Значение | Описание |
|-----------|----------|----------|
| `tasks.TypeCall` | 1 | Звонок |
| `tasks.TypeMeeting` | 2 | Встреча |
| `tasks.TypeTask` | 3 | Задача |

Также доступны константы для типов сущностей:

| Константа | Значение | Описание |
|-----------|----------|----------|
| `tasks.EntityTypeContact` | "contacts" | Контакт |
| `tasks.EntityTypeLead` | "leads" | Сделка |
| `tasks.EntityTypeCompany` | "companies" | Компания |
