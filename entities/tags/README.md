# Модуль Теги

Модуль `tags` предоставляет функциональность для работы с тегами в amoCRM.

## Содержание

- [Основные функции](#основные-функции)
- [Получение тегов](#получение-тегов)
- [Создание тега](#создание-тега)
- [Получение тега по ID](#получение-тега-по-id)
- [Обновление тега](#обновление-тега)
- [Удаление тега](#удаление-тега)
- [Связывание сущностей с тегами](#связывание-сущностей-с-тегами)
- [Получение тегов сущности](#получение-тегов-сущности)
- [Типы сущностей](#типы-сущностей)

## Основные функции

| Функция | Описание |
|---------|----------|
| `GetTags` | Получение списка тегов с пагинацией |
| `CreateTag` | Создание нового тега |
| `CreateTags` | Создание нескольких тегов |
| `GetTag` | Получение тега по ID |
| `UpdateTag` | Обновление тега |
| `DeleteTag` | Удаление тега |
| `LinkEntityWithTags` | Связывание сущности с тегами |
| `GetEntityTags` | Получение тегов сущности |

## Получение тегов

```go
import (
    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/tags"
)

// Инициализация клиента
apiClient := client.NewClient("https://your-domain.amocrm.ru", "your_access_token")

// Получение всех тегов контактов (1-я страница, 50 элементов)
contactTags, err := tags.GetTags(apiClient, tags.EntityTypeContact, 1, 50)
if err != nil {
    // Обработка ошибки
}

// Вывод списка тегов
for _, tag := range contactTags {
    fmt.Printf("ID: %d, Название: %s, Цвет: %s\n", tag.ID, tag.Name, tag.Color)
}

// Получение тегов лидов
leadTags, err := tags.GetTags(apiClient, tags.EntityTypeLead, 1, 50)
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

createdTag, err := tags.CreateTag(apiClient, tags.EntityTypeContact, newTag)
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

createdTags, err := tags.CreateTags(apiClient, tags.EntityTypeContact, newTags)
if err != nil {
    // Обработка ошибки
}

fmt.Printf("Создано %d новых тегов\n", len(createdTags))
```

## Получение тега по ID

```go
// Получение тега по ID
tagID := 12345
tag, err := tags.GetTag(apiClient, tags.EntityTypeContact, tagID)
if err != nil {
    // Обработка ошибки
}

fmt.Printf("Тег: %s (Цвет: %s)\n", tag.Name, tag.Color)
```

## Обновление тега

```go
// Обновление тега
tag.Name = "Очень важный клиент"
tag.Color = "#990000" // Темно-красный цвет

updatedTag, err := tags.UpdateTag(apiClient, tags.EntityTypeContact, tag)
if err != nil {
    // Обработка ошибки
}

fmt.Printf("Тег обновлен: %s\n", updatedTag.Name)
```

## Удаление тега

```go
// Удаление тега
tagID := 12345
err := tags.DeleteTag(apiClient, tags.EntityTypeContact, tagID)
if err != nil {
    // Обработка ошибки
}

fmt.Println("Тег успешно удален")
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

err := tags.LinkEntityWithTags(apiClient, tags.EntityTypeContact, contactID, tagsToLink)
if err != nil {
    // Обработка ошибки
}

// Связывание лида с тегами
leadID := 54321
err = tags.LinkEntityWithTags(apiClient, tags.EntityTypeLead, leadID, tagsToLink)
if err != nil {
    // Обработка ошибки
}
```

## Получение тегов сущности

```go
// Получение тегов контакта
contactID := 67890
contactTags, err := tags.GetEntityTags(apiClient, tags.EntityTypeContact, contactID)
if err != nil {
    // Обработка ошибки
}

fmt.Printf("У контакта %d тегов:\n", len(contactTags))
for _, tag := range contactTags {
    fmt.Printf("- %s (Цвет: %s)\n", tag.Name, tag.Color)
}

// Получение тегов лида
leadID := 54321
leadTags, err := tags.GetEntityTags(apiClient, tags.EntityTypeLead, leadID)
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
