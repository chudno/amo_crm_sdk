# Модуль Каталоги

Модуль `catalogs` предоставляет функциональность для работы с каталогами в amoCRM.

## Содержание

- [Основные функции](#основные-функции)
- [Получение каталогов](#получение-каталогов)
- [Создание каталога](#создание-каталога)
- [Получение каталога по ID](#получение-каталога-по-id)
- [Обновление каталога](#обновление-каталога)
- [Удаление каталога](#удаление-каталога)
- [Работа с пользовательскими полями каталога](#работа-с-пользовательскими-полями-каталога)
- [Типы каталогов](#типы-каталогов)

## Основные функции

| Функция | Описание |
|---------|----------|
| `GetCatalogs` | Получение списка каталогов с пагинацией и фильтрацией |
| `CreateCatalog` | Создание нового каталога |
| `GetCatalog` | Получение каталога по ID |
| `UpdateCatalog` | Обновление каталога |
| `DeleteCatalog` | Удаление каталога |
| `AddCustomFieldToCatalog` | Добавление пользовательского поля в каталог |
| `GetCatalogCustomFields` | Получение списка пользовательских полей каталога |
| `GetCatalogCustomField` | Получение пользовательского поля каталога по ID |
| `UpdateCatalogCustomField` | Обновление пользовательского поля каталога |
| `DeleteCatalogCustomField` | Удаление пользовательского поля каталога |

## Получение каталогов

```go
import (
    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/catalogs"
)

// Инициализация клиента
apiClient := client.NewClient("https://your-domain.amocrm.ru", "your_access_token")

// Получение каталогов (первая страница, 50 элементов)
catalogsList, err := catalogs.GetCatalogs(apiClient, 1, 50, nil)
if err != nil {
    // Обработка ошибки
}

// Вывод списка каталогов
for _, catalog := range catalogsList {
    fmt.Printf("ID: %d, Название: %s, Тип: %s\n", catalog.ID, catalog.Name, catalog.Type)
}

// Получение каталогов с фильтрацией
filter := map[string]string{
    "filter[type]": string(catalogs.CatalogTypeRegular), // Только обычные каталоги
    "filter[name]": "Товары", // Фильтр по названию
}
filteredCatalogs, err := catalogs.GetCatalogs(apiClient, 1, 50, filter)
if err != nil {
    // Обработка ошибки
}
```

## Создание каталога

```go
// Создание нового каталога
newCatalog := &catalogs.Catalog{
    Name: "Товары",
    Type: string(catalogs.CatalogTypeRegular),
    Sort: 100,
}

createdCatalog, err := catalogs.CreateCatalog(apiClient, newCatalog)
if err != nil {
    // Обработка ошибки
}

fmt.Printf("Создан каталог с ID: %d\n", createdCatalog.ID)
```

## Получение каталога по ID

```go
// Получение каталога по ID
catalogID := 12345
catalog, err := catalogs.GetCatalog(apiClient, catalogID)
if err != nil {
    // Обработка ошибки
}

fmt.Printf("Каталог: %s (ID: %d)\n", catalog.Name, catalog.ID)
fmt.Printf("Создан: %d, Обновлен: %d\n", catalog.CreatedAt, catalog.UpdatedAt)
fmt.Printf("Тип: %s\n", catalog.Type)
```

## Обновление каталога

```go
// Обновление каталога
catalog.Name = "Обновленный каталог товаров"
catalog.Sort = 50

updatedCatalog, err := catalogs.UpdateCatalog(apiClient, catalog)
if err != nil {
    // Обработка ошибки
}

fmt.Printf("Каталог обновлен: %s\n", updatedCatalog.Name)
```

## Удаление каталога

```go
// Удаление каталога
catalogID := 12345
err := catalogs.DeleteCatalog(apiClient, catalogID)
if err != nil {
    // Обработка ошибки
}

fmt.Println("Каталог успешно удален")
```

## Работа с пользовательскими полями каталога

### Добавление пользовательского поля

```go
// Добавление пользовательского поля в каталог
catalogID := 12345
newField := &catalogs.CustomField{
    Name:       "Артикул",
    Type:       "text",
    IsRequired: true,
    Sort:       10,
    Code:       "SKU",
}

createdField, err := catalogs.AddCustomFieldToCatalog(apiClient, catalogID, newField)
if err != nil {
    // Обработка ошибки
}

fmt.Printf("Создано поле с ID: %d\n", createdField.ID)

// Добавление поля-списка
selectField := &catalogs.CustomField{
    Name:       "Категория",
    Type:       "select",
    IsRequired: true,
    Sort:       20,
    Code:       "CATEGORY",
}

createdSelectField, err := catalogs.AddCustomFieldToCatalog(apiClient, catalogID, selectField)
if err != nil {
    // Обработка ошибки
}
```

### Получение полей каталога

```go
// Получение всех полей каталога
catalogID := 12345
fields, err := catalogs.GetCatalogCustomFields(apiClient, catalogID)
if err != nil {
    // Обработка ошибки
}

fmt.Printf("У каталога %d полей:\n", len(fields))
for _, field := range fields {
    fmt.Printf("- %s (ID: %d, Тип: %s, Код: %s)\n", field.Name, field.ID, field.Type, field.Code)
}

// Получение конкретного поля
fieldID := 67890
field, err := catalogs.GetCatalogCustomField(apiClient, catalogID, fieldID)
if err != nil {
    // Обработка ошибки
}

fmt.Printf("Поле: %s (Тип: %s)\n", field.Name, field.Type)
```

### Обновление поля каталога

```go
// Обновление пользовательского поля
field.Name = "Обновленное название поля"
field.IsRequired = true

updatedField, err := catalogs.UpdateCatalogCustomField(apiClient, catalogID, field)
if err != nil {
    // Обработка ошибки
}

fmt.Printf("Поле обновлено: %s\n", updatedField.Name)
```

### Удаление поля каталога

```go
// Удаление пользовательского поля
fieldID := 67890
err := catalogs.DeleteCatalogCustomField(apiClient, catalogID, fieldID)
if err != nil {
    // Обработка ошибки
}

fmt.Println("Поле успешно удалено")
```

## Типы каталогов

Модуль `catalogs` предоставляет константы для типов каталогов:

| Константа | Значение | Описание |
|-----------|----------|----------|
| `catalogs.CatalogTypeRegular` | "regular" | Обычный каталог, создаваемый пользователем |
| `catalogs.CatalogTypeContacts` | "contacts" | Системный каталог для контактов |
| `catalogs.CatalogTypeCompanies` | "companies" | Системный каталог для компаний |

## Типы пользовательских полей

При создании пользовательских полей для каталогов можно использовать следующие типы:

| Тип | Описание |
|-----|----------|
| `text` | Текстовое поле |
| `numeric` | Числовое поле |
| `checkbox` | Флажок (Да/Нет) |
| `select` | Список |
| `multiselect` | Мультисписок |
| `date` | Дата |
| `url` | Ссылка |
| `textarea` | Текстовая область |
| `radiobutton` | Переключатель |
| `street_address` | Адрес |
| `birthday` | День рождения |
| `file` | Файл |
