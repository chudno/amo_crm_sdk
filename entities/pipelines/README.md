# Модуль Воронки и статусы

Модуль `pipelines` предоставляет функциональность для работы с воронками и статусами сделок в amoCRM.

## Содержание

- [Основные функции](#основные-функции)
- [Получение воронки](#получение-воронки)
- [Получение списка воронок](#получение-списка-воронок)
- [Создание воронки](#создание-воронки)
- [Обновление воронки](#обновление-воронки)
- [Работа со статусами](#работа-со-статусами)
- [Настройка воронок для сделок](#настройка-воронок-для-сделок)

## Основные функции

| Функция | Описание |
|---------|----------|
| `Get` | Получение воронки по ID |
| `List` | Получение списка воронок |
| `Create` | Создание новой воронки |
| `Update` | Обновление существующей воронки |
| `Delete` | Удаление воронки |
| `GetStatus` | Получение статуса воронки по ID |
| `CreateStatus` | Создание нового статуса в воронке |

## Получение воронки

```go
import (
    "context"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/pipelines"
)

// Инициализация клиента
apiClient := client.NewClient("https://your-domain.amocrm.ru", "your_access_token")

// Создаем контекст
ctx := context.Background()

// Получение воронки по ID
pipelineID := 12345
pipeline, err := pipelines.Get(ctx, apiClient, pipelineID)
if err != nil {
    // Обработка ошибки
}

// Вывод информации о воронке
fmt.Printf("Название воронки: %s\n", pipeline.Name)
fmt.Printf("Количество статусов: %d\n", len(pipeline.Statuses))

// Вывод статусов воронки
for _, status := range pipeline.Statuses {
    fmt.Printf("Статус: %s (ID: %d)\n", status.Name, status.ID)
}
```

## Получение списка воронок

```go
// Получение всех воронок
pipelinesList, err := pipelines.List(ctx, apiClient)
if err != nil {
    // Обработка ошибки
}

// Вывод списка воронок
for _, pipeline := range pipelinesList {
    fmt.Printf("ID: %d, Название: %s\n", pipeline.ID, pipeline.Name)
    
    // Вывод статусов для каждой воронки
    for _, status := range pipeline.Statuses {
        fmt.Printf("  - Статус: %s (ID: %d, Сортировка: %d)\n", 
            status.Name, status.ID, status.Sort)
    }
}
```

## Создание воронки

```go
// Создание новой воронки
newPipeline := &pipelines.Pipeline{
    Name: "Новая воронка продаж",
    Sort: 100, // Порядок сортировки
    IsMain: false, // Является ли основной воронкой
}

// Добавление статусов
newPipeline.Statuses = []pipelines.Status{
    {
        Name: "Первичный контакт",
        Sort: 10,
        Color: "#99ccff", // Цвет статуса в hex-формате
    },
    {
        Name: "Переговоры",
        Sort: 20,
        Color: "#ffcc66",
    },
    {
        Name: "Коммерческое предложение",
        Sort: 30,
        Color: "#ffff99",
    },
    {
        Name: "Договор",
        Sort: 40,
        Color: "#99ff99",
    },
    {
        Name: "Успешно реализовано",
        Sort: 50,
        Color: "#00cc00",
    },
}

// Сохранение воронки
createdPipeline, err := pipelines.Create(ctx, apiClient, newPipeline)
if err != nil {
    // Обработка ошибки
}
```

## Обновление воронки

```go
// Обновление существующей воронки
pipeline.Name = "Обновленная воронка продаж"

// Добавление нового статуса
pipeline.Statuses = append(pipeline.Statuses, pipelines.Status{
    Name: "Отложенная сделка",
    Sort: 25, // Порядок сортировки между существующими статусами
    Color: "#cccccc",
})

// Сохранение изменений
updatedPipeline, err := pipelines.Update(ctx, apiClient, pipeline)
if err != nil {
    // Обработка ошибки
}
```

## Работа со статусами

```go
// Получение статуса по ID
status, err := pipelines.GetStatus(ctx, apiClient, pipeline.ID, 12345)
if err != nil {
    // Обработка ошибки
}
fmt.Printf("Статус: %s (ID: %d)\n", status.Name, status.ID)

// Создание нового статуса в существующей воронке
newStatus := &pipelines.Status{
    Name: "Новый статус",
    Sort: 45,
    Color: "#9966ff",
}

createdStatus, err := pipelines.CreateStatus(ctx, apiClient, pipeline.ID, newStatus)
if err != nil {
    // Обработка ошибки
}
```

## Настройка воронок для сделок

При создании или обновлении лида необходимо указать ID воронки и статуса:

```go
import "github.com/chudno/amo_crm_sdk/entities/leads"

newLead := &leads.Lead{
	Name:      "Новый лид",
	PipelineID: pipeline.ID,
	StatusID:  statusID,
	// Другие поля лида
}

createdLead, err := leads.Create(ctx, apiClient, newLead)
if err != nil {
    // Обработка ошибки
}

// Перемещение лида в другой статус
existingLead, err := leads.Get(ctx, apiClient, leadID)
if err != nil {
    // Обработка ошибки
}

existingLead.StatusID = 54321 // ID нового статуса
// При необходимости можно сменить и воронку
// existingLead.PipelineID = 98765

updatedLead, err := leads.Update(ctx, apiClient, existingLead)
if err != nil {
    // Обработка ошибки
}
