package custom_fields

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestCustomFieldValueJSON(t *testing.T) {
	customFieldValue := Value{
		FieldID:   12345,
		FieldName: "Тестовое поле",
		Values: []FieldValue{
			{
				Value: "Значение поля",
			},
		},
	}

	jsonData, err := json.Marshal(customFieldValue)
	if err != nil {
		t.Fatalf("Ошибка при маршализации в JSON: %v", err)
	}

	jsonStr := string(jsonData)
	if !contains(jsonStr, "field_id") {
		t.Errorf("JSON не содержит ключ field_id: %s", jsonStr)
	}
	if !contains(jsonStr, "field_name") {
		t.Errorf("JSON не содержит ключ field_name: %s", jsonStr)
	}
	if !contains(jsonStr, "values") {
		t.Errorf("JSON не содержит ключ values: %s", jsonStr)
	}

	var decodedValue Value
	if err := json.Unmarshal(jsonData, &decodedValue); err != nil {
		t.Fatalf("Ошибка при демаршализации из JSON: %v", err)
	}

	if decodedValue.FieldID != customFieldValue.FieldID {
		t.Errorf("FieldID не совпадает: ожидалось %d, получено %d", customFieldValue.FieldID, decodedValue.FieldID)
	}
	if decodedValue.FieldName != customFieldValue.FieldName {
		t.Errorf("FieldName не совпадает: ожидалось %s, получено %s", customFieldValue.FieldName, decodedValue.FieldName)
	}
	if len(decodedValue.Values) != len(customFieldValue.Values) {
		t.Errorf("Длина Values не совпадает: ожидалось %d, получено %d", len(customFieldValue.Values), len(decodedValue.Values))
	} else if decodedValue.Values[0].Value != customFieldValue.Values[0].Value {
		t.Errorf("Value не совпадает: ожидалось %v, получено %v", customFieldValue.Values[0].Value, decodedValue.Values[0].Value)
	}
}

func TestFieldValueJSON(t *testing.T) {
	fieldValue := FieldValue{
		Value:     "Тестовое значение",
		EnumID:    789,
		EnumCode:  "test_code",
		EnumValue: "Тестовое значение списка",
	}

	jsonData, err := json.Marshal(fieldValue)
	if err != nil {
		t.Fatalf("Ошибка при маршализации в JSON: %v", err)
	}

	jsonStr := string(jsonData)
	if !contains(jsonStr, "value") {
		t.Errorf("JSON не содержит ключ value: %s", jsonStr)
	}
	if !contains(jsonStr, "enum_id") {
		t.Errorf("JSON не содержит ключ enum_id: %s", jsonStr)
	}
	if !contains(jsonStr, "enum_code") {
		t.Errorf("JSON не содержит ключ enum_code: %s", jsonStr)
	}
	if !contains(jsonStr, "enum_value") {
		t.Errorf("JSON не содержит ключ enum_value: %s", jsonStr)
	}

	var decodedValue FieldValue
	if err := json.Unmarshal(jsonData, &decodedValue); err != nil {
		t.Fatalf("Ошибка при демаршализации из JSON: %v", err)
	}

	if decodedValue.Value != fieldValue.Value {
		t.Errorf("Value не совпадает: ожидалось %v, получено %v", fieldValue.Value, decodedValue.Value)
	}
	if decodedValue.EnumID != fieldValue.EnumID {
		t.Errorf("EnumID не совпадает: ожидалось %d, получено %d", fieldValue.EnumID, decodedValue.EnumID)
	}
	if decodedValue.EnumCode != fieldValue.EnumCode {
		t.Errorf("EnumCode не совпадает: ожидалось %s, получено %s", fieldValue.EnumCode, decodedValue.EnumCode)
	}
	if decodedValue.EnumValue != fieldValue.EnumValue {
		t.Errorf("EnumValue не совпадает: ожидалось %s, получено %s", fieldValue.EnumValue, decodedValue.EnumValue)
	}
}

func TestCustomFieldJSON(t *testing.T) {
	customField := Field{
		ID:           123,
		Name:         "Тестовое поле",
		Type:         "text",
		Code:         "test_field",
		Sort:         10,
		IsMultiple:   false,
		IsSystem:     false,
		IsEditable:   true,
		IsRequired:   false,
		IsDeleteable: true,
		IsVisible:    true,
		Enums: []Enum{
			{
				ID:    456,
				Value: "Вариант 1",
				Sort:  1,
				Code:  "option_1",
			},
		},
	}

	jsonData, err := json.Marshal(customField)
	if err != nil {
		t.Fatalf("Ошибка при маршализации в JSON: %v", err)
	}

	jsonStr := string(jsonData)
	if !contains(jsonStr, "id") {
		t.Errorf("JSON не содержит ключ id: %s", jsonStr)
	}
	if !contains(jsonStr, "name") {
		t.Errorf("JSON не содержит ключ name: %s", jsonStr)
	}
	if !contains(jsonStr, "type") {
		t.Errorf("JSON не содержит ключ type: %s", jsonStr)
	}
	if !contains(jsonStr, "enums") {
		t.Errorf("JSON не содержит ключ enums: %s", jsonStr)
	}

	var decodedField Field
	if err := json.Unmarshal(jsonData, &decodedField); err != nil {
		t.Fatalf("Ошибка при демаршализации из JSON: %v", err)
	}

	if !reflect.DeepEqual(customField, decodedField) {
		t.Errorf("Демаршализованная структура не совпадает с исходной")
		t.Errorf("Ожидалось: %+v", customField)
		t.Errorf("Получено: %+v", decodedField)
	}
}

func TestCustomFieldEnumJSON(t *testing.T) {
	customFieldEnum := Enum{
		ID:    789,
		Value: "Тестовый вариант",
		Sort:  5,
		Code:  "test_option",
	}

	jsonData, err := json.Marshal(customFieldEnum)
	if err != nil {
		t.Fatalf("Ошибка при маршализации в JSON: %v", err)
	}

	jsonStr := string(jsonData)
	if !contains(jsonStr, "id") {
		t.Errorf("JSON не содержит ключ id: %s", jsonStr)
	}
	if !contains(jsonStr, "value") {
		t.Errorf("JSON не содержит ключ value: %s", jsonStr)
	}
	if !contains(jsonStr, "sort") {
		t.Errorf("JSON не содержит ключ sort: %s", jsonStr)
	}
	if !contains(jsonStr, "code") {
		t.Errorf("JSON не содержит ключ code: %s", jsonStr)
	}

	var decodedEnum Enum
	if err := json.Unmarshal(jsonData, &decodedEnum); err != nil {
		t.Fatalf("Ошибка при демаршализации из JSON: %v", err)
	}

	if decodedEnum.ID != customFieldEnum.ID {
		t.Errorf("ID не совпадает: ожидалось %d, получено %d", customFieldEnum.ID, decodedEnum.ID)
	}
	if decodedEnum.Value != customFieldEnum.Value {
		t.Errorf("Value не совпадает: ожидалось %s, получено %s", customFieldEnum.Value, decodedEnum.Value)
	}
	if decodedEnum.Sort != customFieldEnum.Sort {
		t.Errorf("Sort не совпадает: ожидалось %d, получено %d", customFieldEnum.Sort, decodedEnum.Sort)
	}
	if decodedEnum.Code != customFieldEnum.Code {
		t.Errorf("Code не совпадает: ожидалось %s, получено %s", customFieldEnum.Code, decodedEnum.Code)
	}
}

// Вспомогательная функция для проверки наличия подстроки в строке
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
