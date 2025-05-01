package access_rights

import (
	"testing"
)

// TestWithFilter проверяет функциональную опцию WithFilter
func TestWithFilter(t *testing.T) {
	// Создаем тестовые данные
	filter := map[string]string{
		"filter[name]": "Тестовая группа",
		"filter[type]": "group",
	}

	// Создаем параметры запроса
	params := make(map[string]string)

	// Применяем функциональную опцию
	option := WithFilter(filter)
	option(params)

	// Проверяем результат
	if params["filter[name]"] != "Тестовая группа" {
		t.Errorf(`Ожидалось params["filter[name]"] = "Тестовая группа", получено %q`, params["filter[name]"])
	}

	if params["filter[type]"] != "group" {
		t.Errorf(`Ожидалось params["filter[type]"] = "group", получено %q`, params["filter[type]"])
	}
}

// TestWithType проверяет функциональную опцию WithType
func TestWithType(t *testing.T) {
	// Создаем параметры запроса
	params := make(map[string]string)

	// Применяем функциональную опцию для типа group
	option := WithType(TypeGroup)
	option(params)

	// Проверяем результат
	if params["filter[type]"] != string(TypeGroup) {
		t.Errorf(`Ожидалось params["filter[type]"] = %q, получено %q`, string(TypeGroup), params["filter[type]"])
	}

	// Очищаем параметры
	params = make(map[string]string)

	// Применяем функциональную опцию для типа custom
	option = WithType(TypeCustom)
	option(params)

	// Проверяем результат
	if params["filter[type]"] != string(TypeCustom) {
		t.Errorf(`Ожидалось params["filter[type]"] = %q, получено %q`, string(TypeCustom), params["filter[type]"])
	}
}
