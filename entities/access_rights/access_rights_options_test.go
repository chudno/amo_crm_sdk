package access_rights

import (
	"testing"
)

// TestWithFilter проверяет функциональную опцию WithFilter
func TestWithFilter(t *testing.T) {
	filter := map[string]string{
		"filter[name]": "Тестовая группа",
		"filter[type]": "group",
	}

	params := make(map[string]string)

	option := WithFilter(filter)
	option(params)

	if params["filter[name]"] != "Тестовая группа" {
		t.Errorf(`Ожидалось params["filter[name]"] = "Тестовая группа", получено %q`, params["filter[name]"])
	}

	if params["filter[type]"] != "group" {
		t.Errorf(`Ожидалось params["filter[type]"] = "group", получено %q`, params["filter[type]"])
	}
}

// TestWithType проверяет функциональную опцию WithType
func TestWithType(t *testing.T) {
	params := make(map[string]string)

	option := WithType(TypeGroup)
	option(params)

	if params["filter[type]"] != string(TypeGroup) {
		t.Errorf(`Ожидалось params["filter[type]"] = %q, получено %q`, string(TypeGroup), params["filter[type]"])
	}

	params = make(map[string]string)

	option = WithType(TypeCustom)
	option(params)

	if params["filter[type]"] != string(TypeCustom) {
		t.Errorf(`Ожидалось params["filter[type]"] = %q, получено %q`, string(TypeCustom), params["filter[type]"])
	}
}
