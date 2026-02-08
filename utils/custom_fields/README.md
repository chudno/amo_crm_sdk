# Модуль Пользовательские поля

Модуль `custom_fields` предоставляет вспомогательные структуры для работы с пользовательскими полями в amoCRM. Эти структуры используются при создании и обновлении сущностей (контактов, сделок и др.) через соответствующие модули SDK.

## Содержание

- [Структуры данных](#структуры-данных)
- [Работа с пользовательскими полями в сущностях](#работа-с-пользовательскими-полями-в-сущностях)

## Структуры данных

### Value

Значение пользовательского поля, привязанное к сущности:

```go
type Value struct {
    FieldID   int          `json:"field_id"`
    FieldName string       `json:"field_name,omitempty"`
    Values    []FieldValue `json:"values"`
}
```

### FieldValue

Конкретное значение поля:

```go
type FieldValue struct {
    Value     any    `json:"value"`
    EnumID    int    `json:"enum_id,omitempty"`
    EnumCode  string `json:"enum_code,omitempty"`
    EnumValue string `json:"enum_value,omitempty"`
}
```

### Field

Структура пользовательского поля:

```go
type Field struct {
    ID           int    `json:"id"`
    Name         string `json:"name"`
    Type         string `json:"type"`
    Code         string `json:"code,omitempty"`
    Sort         int    `json:"sort,omitempty"`
    IsMultiple   bool   `json:"is_multiple,omitempty"`
    IsSystem     bool   `json:"is_system,omitempty"`
    IsEditable   bool   `json:"is_editable,omitempty"`
    IsRequired   bool   `json:"is_required,omitempty"`
    IsDeleteable bool   `json:"is_deleteable,omitempty"`
    IsVisible    bool   `json:"is_visible,omitempty"`
    Enums        []Enum `json:"enums,omitempty"`
}
```

### Enum

Вариант значения для поля типа список:

```go
type Enum struct {
    ID    int    `json:"id"`
    Value string `json:"value"`
    Sort  int    `json:"sort,omitempty"`
    Code  string `json:"code,omitempty"`
}
```

## Работа с пользовательскими полями в сущностях

Структуры этого пакета используются совместно с модулями сущностей SDK. Например, при установке значений пользовательских полей для контакта или сделки можно использовать `custom_fields.Value` и `custom_fields.FieldValue`.
