# Модуль Пользовательские поля

Модуль `custom_fields` предоставляет функциональность для работы с пользовательскими полями в amoCRM.

## Содержание

- [Основные функции](#основные-функции)
- [Получение пользовательских полей](#получение-пользовательских-полей)
- [Создание пользовательского поля](#создание-пользовательского-поля)
- [Обновление пользовательского поля](#обновление-пользовательского-поля)
- [Типы пользовательских полей](#типы-пользовательских-полей)
- [Работа с пользовательскими полями в сущностях](#работа-с-пользовательскими-полями-в-сущностях)

## Основные функции

| Функция | Описание |
|---------|----------|
| `GetCustomFields` | Получение списка пользовательских полей для указанного типа сущности |
| `GetCustomField` | Получение пользовательского поля по ID |
| `CreateCustomField` | Создание нового пользовательского поля |
| `UpdateCustomField` | Обновление существующего пользовательского поля |
| `DeleteCustomField` | Удаление пользовательского поля |

## Получение пользовательских полей

```go
import (
    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/utils/custom_fields"
)

// Инициализация клиента
apiClient := client.NewClient("https://your-domain.amocrm.ru", "your_access_token")

// Получение всех пользовательских полей для контактов
entityType := "contacts"
fieldsList, err := custom_fields.GetCustomFields(apiClient, entityType)
if err != nil {
    // Обработка ошибки
}

// Вывод списка полей
for _, field := range fieldsList {
    fmt.Printf("ID: %d, Название: %s, Тип: %s\n", field.ID, field.Name, field.FieldType)
    
    // Вывод возможных значений для списка
    if field.Enums != nil && len(field.Enums) > 0 {
        fmt.Println("Возможные значения:")
        for _, enum := range field.Enums {
            fmt.Printf("  - ID: %d, Значение: %s\n", enum.ID, enum.Value)
        }
    }
}

// Получение пользовательского поля по ID
fieldID := 12345
field, err := custom_fields.GetCustomField(apiClient, entityType, fieldID)
if err != nil {
    // Обработка ошибки
}
```

## Создание пользовательского поля

```go
// Создание текстового поля
newTextField := &custom_fields.CustomField{
    Name: "Комментарий",
    FieldType: custom_fields.TypeText,
    EntityType: "contacts",
    Sort: 100, // Порядок сортировки
}

createdField, err := custom_fields.CreateCustomField(apiClient, newTextField)
if err != nil {
    // Обработка ошибки
}

// Создание поля-списка с предопределенными значениями
newSelectField := &custom_fields.CustomField{
    Name: "Источник клиента",
    FieldType: custom_fields.TypeSelect,
    EntityType: "contacts",
    Sort: 110,
    Enums: []custom_fields.Enum{
        {Value: "Сайт"},
        {Value: "Реклама"},
        {Value: "Рекомендация"},
        {Value: "Холодный звонок"},
        {Value: "Партнер"},
    },
}

createdSelectField, err := custom_fields.CreateCustomField(apiClient, newSelectField)
if err != nil {
    // Обработка ошибки
}

// Создание поля с телефоном (с типами телефонов)
newPhoneField := &custom_fields.CustomField{
    Name: "Телефон",
    FieldType: custom_fields.TypePhone,
    EntityType: "contacts",
    Sort: 10,
    IsMultiple: true, // Разрешить несколько телефонов
    Enums: []custom_fields.Enum{
        {Value: "Рабочий", Code: "WORK"},
        {Value: "Домашний", Code: "HOME"},
        {Value: "Мобильный", Code: "PERSONAL"},
        {Value: "Другой", Code: "OTHER"},
    },
}

createdPhoneField, err := custom_fields.CreateCustomField(apiClient, newPhoneField)
```

## Обновление пользовательского поля

```go
// Обновление существующего поля
field.Name = "Новое название поля"

// Добавление нового значения в список
if field.FieldType == custom_fields.TypeSelect {
    field.Enums = append(field.Enums, custom_fields.Enum{
        Value: "Новое значение",
    })
}

// Сохранение изменений
updatedField, err := custom_fields.UpdateCustomField(apiClient, field)
if err != nil {
    // Обработка ошибки
}
```

## Типы пользовательских полей

Модуль `custom_fields` предоставляет константы для типов пользовательских полей:

| Константа | Значение | Описание |
|-----------|----------|----------|
| `custom_fields.TypeText` | "text" | Текстовое поле |
| `custom_fields.TypeNumeric` | "numeric" | Числовое поле |
| `custom_fields.TypeCheckbox` | "checkbox" | Флажок (Да/Нет) |
| `custom_fields.TypeSelect` | "select" | Список |
| `custom_fields.TypeMultiselect` | "multiselect" | Мультисписок |
| `custom_fields.TypeDate` | "date" | Дата |
| `custom_fields.TypeURL` | "url" | Ссылка |
| `custom_fields.TypeTextarea` | "textarea" | Текстовая область |
| `custom_fields.TypePhone` | "phone" | Телефон |
| `custom_fields.TypeEmail` | "email" | Email |

## Работа с пользовательскими полями в сущностях

Когда пользовательское поле создано, вы можете добавлять, обновлять и получать его значения в различных сущностях:

```go
import (
    "github.com/chudno/amo_crm_sdk/entities/contacts"
    "github.com/chudno/amo_crm_sdk/entities/leads"
)

// Добавление значения пользовательского поля для контакта
contact := &contacts.Contact{
    Name: "Иван Иванов",
}

// Добавление текстового поля
contact.CustomFields = append(contact.CustomFields, contacts.CustomField{
    FieldID: 12345, // ID пользовательского поля
    Values: []contacts.CustomFieldValue{
        {
            Value: "Значение поля",
        },
    },
})

// Добавление поля типа "Список"
contact.CustomFields = append(contact.CustomFields, contacts.CustomField{
    FieldID: 67890, // ID поля-списка
    Values: []contacts.CustomFieldValue{
        {
            Value: "Сайт", // Значение должно соответствовать одному из предопределенных значений
            EnumID: 123, // ID значения из списка (опционально)
        },
    },
})

// Добавление телефона с типом
contact.CustomFields = append(contact.CustomFields, contacts.CustomField{
    FieldID: 54321, // ID поля типа "Телефон"
    Values: []contacts.CustomFieldValue{
        {
            Value: "+79001234567",
            Enum: "WORK", // Тип телефона (код из списка)
        },
    },
})

// Сохранение контакта с пользовательскими полями
createdContact, err := contacts.CreateContact(apiClient, contact)
```

Работа с пользовательскими полями в лидах аналогична:

```go
// Добавление значения пользовательского поля для лида
lead := &leads.Lead{
    Name: "Продажа ПО",
}

// Добавление числового поля (например, "Бюджет")
lead.CustomFields = append(lead.CustomFields, leads.CustomField{
    FieldID: 98765, // ID числового поля
    Values: []leads.CustomFieldValue{
        {
            Value: "50000", // Числовое значение передается в виде строки
        },
    },
})

// Добавление поля типа "Дата"
lead.CustomFields = append(lead.CustomFields, leads.CustomField{
    FieldID: 45678, // ID поля типа "Дата"
    Values: []leads.CustomFieldValue{
        {
            Value: "1680307200", // Unix timestamp как строка
        },
    },
})

// Сохранение лида с пользовательскими полями
createdLead, err := leads.CreateLead(apiClient, lead)
```
