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
| `GetPipeline` | Получение воронки по ID |
| `GetPipelines` | Получение списка воронок |
| `CreatePipeline` | Создание новой воронки |
| `UpdatePipeline` | Обновление существующей воронки |
| `DeletePipeline` | Удаление воронки |
| `CreateStatus` | Создание нового статуса в воронке |
| `UpdateStatus` | Обновление существующего статуса |
| `DeleteStatus` | Удаление статуса |

## Получение воронки

```go
import (
    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/pipelines"
)

// Инициализация клиента
apiClient := client.NewClient("https://your-domain.amocrm.ru", "your_access_token")

// Получение воронки по ID
pipelineID := 12345
pipeline, err := pipelines.GetPipeline(apiClient, pipelineID)
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
pipelinesList, err := pipelines.GetPipelines(apiClient)
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
        IsSuccess: true, // Статус успешного завершения сделки
    },
}

// Сохранение воронки
createdPipeline, err := pipelines.CreatePipeline(apiClient, newPipeline)
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
updatedPipeline, err := pipelines.UpdatePipeline(apiClient, pipeline)
if err != nil {
    // Обработка ошибки
}
```

## Работа со статусами

```go
// Получение статуса по ID
var targetStatus *pipelines.Status
for _, status := range pipeline.Statuses {
    if status.ID == 12345 {
        targetStatus = &status
        break
    }
}

// Обновление статуса
if targetStatus != nil {
    targetStatus.Name = "Новое название статуса"
    targetStatus.Color = "#ff9900"
    
    // Обновление статуса в воронке
    updatedStatus, err := pipelines.UpdateStatus(apiClient, pipeline.ID, targetStatus)
    if err != nil {
        // Обработка ошибки
    }
}

// Создание нового статуса в существующей воронке
newStatus := &pipelines.Status{
    Name: "Новый статус",
    Sort: 45, // Позиция в воронке
    Color: "#9966ff",
}

createdStatus, err := pipelines.CreateStatus(apiClient, pipeline.ID, newStatus)
if err != nil {
    // Обработка ошибки
}

// Удаление статуса
statusID := 67890
err = pipelines.DeleteStatus(apiClient, pipeline.ID, statusID)
if err != nil {
    // Обработка ошибки
}
```

## Настройка воронок для сделок

При создании или обновлении сделки необходимо указать ID воронки и статуса:

```go
import "github.com/chudno/amo_crm_sdk/entities/deals"

// Создание сделки в определенной воронке и статусе
newDeal := &deals.Deal{
    Name: "Новая сделка",
    PipelineID: 12345, // ID воронки
    StatusID: 67890,   // ID статуса в этой воронке
}

// Создание сделки
createdDeal, err := deals.CreateDeal(apiClient, newDeal)
if err != nil {
    // Обработка ошибки
}

// Перемещение сделки в другой статус
existingDeal, err := deals.GetDeal(apiClient, dealID)
if err != nil {
    // Обработка ошибки
}

existingDeal.StatusID = 54321 // ID нового статуса
// При необходимости можно сменить и воронку
// existingDeal.PipelineID = 98765

updatedDeal, err := deals.UpdateDeal(apiClient, existingDeal)
if err != nil {
    // Обработка ошибки
}
```
