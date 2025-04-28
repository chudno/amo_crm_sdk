// Пакет custom_fields предоставляет структуры и методы для работы с пользовательскими полями в amoCRM.
package custom_fields

// CustomFieldValue представляет значение пользовательского поля
type CustomFieldValue struct {
	FieldID   int          `json:"field_id"`
	FieldName string       `json:"field_name,omitempty"`
	Values    []FieldValue `json:"values"`
}

// FieldValue представляет конкретное значение поля
type FieldValue struct {
	Value     interface{} `json:"value"`
	EnumID    int         `json:"enum_id,omitempty"`
	EnumCode  string      `json:"enum_code,omitempty"`
	EnumValue string      `json:"enum_value,omitempty"`
}

// CustomField представляет структуру пользовательского поля
type CustomField struct {
	ID           int               `json:"id"`
	Name         string            `json:"name"`
	Type         string            `json:"type"`
	Code         string            `json:"code,omitempty"`
	Sort         int               `json:"sort,omitempty"`
	IsMultiple   bool              `json:"is_multiple,omitempty"`
	IsSystem     bool              `json:"is_system,omitempty"`
	IsEditable   bool              `json:"is_editable,omitempty"`
	IsRequired   bool              `json:"is_required,omitempty"`
	IsDeleteable bool              `json:"is_deleteable,omitempty"`
	IsVisible    bool              `json:"is_visible,omitempty"`
	Enums        []CustomFieldEnum `json:"enums,omitempty"`
}

// CustomFieldEnum представляет вариант значения для поля типа список
type CustomFieldEnum struct {
	ID    int    `json:"id"`
	Value string `json:"value"`
	Sort  int    `json:"sort,omitempty"`
	Code  string `json:"code,omitempty"`
}
