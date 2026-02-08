# Модуль Теги

Модуль `tags` предоставляет функциональность для работы с тегами в amoCRM.

## Содержание

- [Основные функции](#основные-функции)
- [Получение тегов](#получение-тегов)
- [Создание тега](#создание-тега)
- [Получение тега по ID](#получение-тега-по-id)
- [Обновление тега](#обновление-тега)
- [Связывание сущностей с тегами](#связывание-сущностей-с-тегами)
- [Получение тегов сущности](#получение-тегов-сущности)
- [Типы сущностей](#типы-сущностей)

## Основные функции

| Функция | Описание |
|---------|----------|
| `List` | Получение списка тегов с пагинацией |
| `Create` | Создание нового тега |
| `CreateBatch` | Создание нескольких тегов |
| `Get` | Получение тега по ID |
| `Update` | Обновление тега |
| `LinkEntity` | Связывание сущности с тегами |
| `ListForEntity` | Получение тегов сущности |

## Получение тегов

```go
import (
    "context"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/tags"
)

// Инициализация клиента
apiClient := client.NewClient("https://your-domain.amocrm.ru", "your_access_token")

// Создаем контекст
ctx := context.Background()

// Получение всех тегов контактов (1-я страница, 50 элементов)
contactTags, err := tags.List(ctx, apiClient, tags.EntityTypeContact, 1, 50)
if err != nil {
    // Обработка ошибки
}

// Вывод списка тегов
for _, tag := range contactTags {
    fmt.Printf("ID: %d, Название: %s, Цвет: %s\n", tag.ID, tag.Name, tag.Color)
}

// Получение тегов лидов
leadTags, err := tags.List(ctx, apiClient, tags.EntityTypeLead, 1, 50)
if err != nil {
    // Обработка ошибки
}
```

## Создание тега

```go
// Создание нового тега для контактов
newTag := &tags.Tag{
    Name:  "Важный клиент",
    Color: "#FF0000", // Красный цвет
}

createdTag, err := tags.Create(ctx, apiClient, tags.EntityTypeContact, newTag)
if err != nil {
    // Обработка ошибки
}

fmt.Printf("Создан тег с ID: %d\n", createdTag.ID)

// Создание нескольких тегов за один запрос
newTags := []tags.Tag{
    {
        Name:  "Потенциальный клиент",
        Color: "#00FF00", // Зеленый цвет
    },
    {
        Name:  "Партнер",
        Color: "#0000FF", // Синий цвет
    },
}

createdTags, err := tags.CreateBatch(ctx, apiClient, tags.EntityTypeContact, newTags)
if err != nil {
    // Обработка ошибки
}

fmt.Printf("Создано %d новых тегов\n", len(createdTags))
```

## Получение тега по ID

```go
// Получение тега по ID
tagID := 12345
tag, err := tags.Get(ctx, apiClient, tags.EntityTypeContact, tagID)
if err != nil {
    // Обработка ошибки
}

fmt.Printf("Тег: %s (Цвет: %s)\n", tag.Name, tag.Color)
```

> **Примечание:** API amoCRM v4 не предоставляет отдельного эндпоинта для получения тега по ID. Функция `Get` использует фильтрацию списка тегов (`filter[id][]`) и возвращает первый найденный результат.

## Обновление тега

```go
// Обновление тега
tag.Name = "Очень важный клиент"
tag.Color = "#990000" // Темно-красный цвет

updatedTag, err := tags.Update(ctx, apiClient, tags.EntityTypeContact, tag)
if err != nil {
    // Обработка ошибки
}

fmt.Printf("Тег обновлен: %s\n", updatedTag.Name)
```

## Связывание сущностей с тегами

```go
// Связывание контакта с тегами
contactID := 67890
tagsToLink := []tags.Tag{
    {
        ID: 123, // Существующий тег
    },
    {
        Name: "Новый тег", // Будет создан новый тег
        Color: "#FFCC00",
    },
}

err := tags.LinkEntity(ctx, apiClient, tags.EntityTypeContact, contactID, tagsToLink)
if err != nil {
    // Обработка ошибки
}

// Связывание лида с тегами
leadID := 54321
err = tags.LinkEntity(ctx, apiClient, tags.EntityTypeLead, leadID, tagsToLink)
if err != nil {
    // Обработка ошибки
}
```

## Получение тегов сущности

```go
// Получение тегов контакта
contactID := 67890
contactTags, err := tags.ListForEntity(ctx, apiClient, tags.EntityTypeContact, contactID)
if err != nil {
    // Обработка ошибки
}

fmt.Printf("У контакта %d тегов:\n", len(contactTags))
for _, tag := range contactTags {
    fmt.Printf("- %s (Цвет: %s)\n", tag.Name, tag.Color)
}

// Получение тегов лида
leadID := 54321
leadTags, err := tags.ListForEntity(ctx, apiClient, tags.EntityTypeLead, leadID)
if err != nil {
    // Обработка ошибки
}
```

## Типы сущностей

Модуль `tags` предоставляет константы для типов сущностей:

| Константа | Значение | Описание |
|-----------|----------|----------|
| `tags.EntityTypeContact` | "contacts" | Контакты |
| `tags.EntityTypeLead` | "leads" | Сделки |
| `tags.EntityTypeCompany` | "companies" | Компании |
| `tags.EntityTypeCustomer` | "customers" | Покупатели |
