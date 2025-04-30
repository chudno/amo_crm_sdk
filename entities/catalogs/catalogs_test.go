package catalogs

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

func TestGetCatalogs(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/catalogs"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Проверяем параметры запроса
		expectedPage := "1"
		if r.URL.Query().Get("page") != expectedPage {
			t.Errorf("Ожидался параметр page=%s, получен %s", expectedPage, r.URL.Query().Get("page"))
		}

		expectedLimit := "50"
		if r.URL.Query().Get("limit") != expectedLimit {
			t.Errorf("Ожидался параметр limit=%s, получен %s", expectedLimit, r.URL.Query().Get("limit"))
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"page": 1,
			"per_page": 50,
			"total": 2,
			"_embedded": {
				"catalogs": [
					{
						"id": 123,
						"name": "Тестовый каталог 1",
						"created_by": 456,
						"updated_by": 456,
						"created_at": 1609459200,
						"updated_at": 1609545600,
						"sort": 1,
						"type": "regular"
					},
					{
						"id": 456,
						"name": "Тестовый каталог 2",
						"created_by": 456,
						"updated_by": 456,
						"created_at": 1609459200,
						"updated_at": 1609545600,
						"sort": 2,
						"type": "regular"
					}
				]
			}
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	catalogs, err := GetCatalogs(apiClient, 1, 50, nil)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при получении каталогов: %v", err)
	}

	if len(catalogs) != 2 {
		t.Fatalf("Ожидалось получение 2 каталогов, получено %d", len(catalogs))
	}

	// Проверяем содержимое первого каталога
	expectedCatalog1 := Catalog{
		ID:        123,
		Name:      "Тестовый каталог 1",
		CreatedBy: 456,
		UpdatedBy: 456,
		CreatedAt: 1609459200,
		UpdatedAt: 1609545600,
		Sort:      1,
		Type:      "regular",
	}
	if !reflect.DeepEqual(catalogs[0], expectedCatalog1) {
		t.Errorf("Ожидался каталог %+v, получен %+v", expectedCatalog1, catalogs[0])
	}

	// Проверяем содержимое второго каталога
	expectedCatalog2 := Catalog{
		ID:        456,
		Name:      "Тестовый каталог 2",
		CreatedBy: 456,
		UpdatedBy: 456,
		CreatedAt: 1609459200,
		UpdatedAt: 1609545600,
		Sort:      2,
		Type:      "regular",
	}
	if !reflect.DeepEqual(catalogs[1], expectedCatalog2) {
		t.Errorf("Ожидался каталог %+v, получен %+v", expectedCatalog2, catalogs[1])
	}
}

func TestCreateCatalog(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/catalogs"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{
			"id": 789,
			"name": "Новый каталог",
			"created_by": 456,
			"updated_by": 456,
			"created_at": 1609459200,
			"updated_at": 1609459200,
			"sort": 3,
			"type": "regular"
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Создаем каталог для отправки
	newCatalog := &Catalog{
		Name: "Новый каталог",
		Sort: 3,
		Type: "regular",
	}

	// Вызываем тестируемый метод
	createdCatalog, err := CreateCatalog(apiClient, newCatalog)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при создании каталога: %v", err)
	}

	// Проверяем содержимое созданного каталога
	expectedCatalog := &Catalog{
		ID:        789,
		Name:      "Новый каталог",
		CreatedBy: 456,
		UpdatedBy: 456,
		CreatedAt: 1609459200,
		UpdatedAt: 1609459200,
		Sort:      3,
		Type:      "regular",
	}
	if !reflect.DeepEqual(createdCatalog, expectedCatalog) {
		t.Errorf("Ожидался каталог %+v, получен %+v", expectedCatalog, createdCatalog)
	}
}

func TestGetCatalog(t *testing.T) {
	catalogID := 123

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/catalogs/%d", catalogID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 123,
			"name": "Тестовый каталог",
			"created_by": 456,
			"updated_by": 456,
			"created_at": 1609459200,
			"updated_at": 1609545600,
			"sort": 1,
			"type": "regular"
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	catalog, err := GetCatalog(apiClient, catalogID)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при получении каталога: %v", err)
	}

	// Проверяем содержимое каталога
	expectedCatalog := &Catalog{
		ID:        123,
		Name:      "Тестовый каталог",
		CreatedBy: 456,
		UpdatedBy: 456,
		CreatedAt: 1609459200,
		UpdatedAt: 1609545600,
		Sort:      1,
		Type:      "regular",
	}
	if !reflect.DeepEqual(catalog, expectedCatalog) {
		t.Errorf("Ожидался каталог %+v, получен %+v", expectedCatalog, catalog)
	}
}

func TestUpdateCatalog(t *testing.T) {
	catalogID := 123

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "PATCH" {
			t.Errorf("Ожидался метод PATCH, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/catalogs/%d", catalogID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 123,
			"name": "Обновленный каталог",
			"created_by": 456,
			"updated_by": 456,
			"created_at": 1609459200,
			"updated_at": 1609632000,
			"sort": 5,
			"type": "regular"
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Создаем каталог для обновления
	catalogToUpdate := &Catalog{
		ID:   catalogID,
		Name: "Обновленный каталог",
		Sort: 5,
	}

	// Вызываем тестируемый метод
	updatedCatalog, err := UpdateCatalog(apiClient, catalogToUpdate)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при обновлении каталога: %v", err)
	}

	// Проверяем содержимое обновленного каталога
	expectedCatalog := &Catalog{
		ID:        123,
		Name:      "Обновленный каталог",
		CreatedBy: 456,
		UpdatedBy: 456,
		CreatedAt: 1609459200,
		UpdatedAt: 1609632000,
		Sort:      5,
		Type:      "regular",
	}
	if !reflect.DeepEqual(updatedCatalog, expectedCatalog) {
		t.Errorf("Ожидался каталог %+v, получен %+v", expectedCatalog, updatedCatalog)
	}
}

func TestDeleteCatalog(t *testing.T) {
	catalogID := 123

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "DELETE" {
			t.Errorf("Ожидался метод DELETE, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/catalogs/%d", catalogID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	err := DeleteCatalog(apiClient, catalogID)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при удалении каталога: %v", err)
	}
}

func TestAddCustomFieldToCatalog(t *testing.T) {
	catalogID := 123

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/catalogs/%d/custom_fields", catalogID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{
			"id": 456,
			"name": "Тестовое поле",
			"type": "text",
			"is_api_only": false,
			"is_required": false,
			"is_multiple": false,
			"is_system": false,
			"sort": 1,
			"code": "TEST_FIELD"
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Создаем поле для добавления
	newField := &CustomField{
		Name:       "Тестовое поле",
		Type:       "text",
		IsRequired: false,
		IsMultiple: false,
		Sort:       1,
		Code:       "TEST_FIELD",
	}

	// Вызываем тестируемый метод
	createdField, err := AddCustomFieldToCatalog(apiClient, catalogID, newField)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при добавлении поля в каталог: %v", err)
	}

	// Проверяем содержимое созданного поля
	expectedField := &CustomField{
		ID:         456,
		Name:       "Тестовое поле",
		Type:       "text",
		IsAPIOnly:  false,
		IsRequired: false,
		IsMultiple: false,
		IsSystem:   false,
		Sort:       1,
		Code:       "TEST_FIELD",
	}
	if !reflect.DeepEqual(createdField, expectedField) {
		t.Errorf("Ожидалось поле %+v, получено %+v", expectedField, createdField)
	}
}

func TestGetCatalogCustomFields(t *testing.T) {
	catalogID := 123

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/catalogs/%d/custom_fields", catalogID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"_embedded": {
				"custom_fields": [
					{
						"id": 456,
						"name": "Тестовое поле 1",
						"type": "text",
						"is_api_only": false,
						"is_required": false,
						"is_multiple": false,
						"is_system": false,
						"sort": 1,
						"code": "TEST_FIELD_1"
					},
					{
						"id": 789,
						"name": "Тестовое поле 2",
						"type": "select",
						"is_api_only": false,
						"is_required": true,
						"is_multiple": false,
						"is_system": false,
						"sort": 2,
						"code": "TEST_FIELD_2"
					}
				]
			}
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	fields, err := GetCatalogCustomFields(apiClient, catalogID)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при получении полей каталога: %v", err)
	}

	if len(fields) != 2 {
		t.Fatalf("Ожидалось получение 2 полей, получено %d", len(fields))
	}

	// Проверяем содержимое первого поля
	expectedField1 := CustomField{
		ID:         456,
		Name:       "Тестовое поле 1",
		Type:       "text",
		IsAPIOnly:  false,
		IsRequired: false,
		IsMultiple: false,
		IsSystem:   false,
		Sort:       1,
		Code:       "TEST_FIELD_1",
	}
	if !reflect.DeepEqual(fields[0], expectedField1) {
		t.Errorf("Ожидалось поле %+v, получено %+v", expectedField1, fields[0])
	}

	// Проверяем содержимое второго поля
	expectedField2 := CustomField{
		ID:         789,
		Name:       "Тестовое поле 2",
		Type:       "select",
		IsAPIOnly:  false,
		IsRequired: true,
		IsMultiple: false,
		IsSystem:   false,
		Sort:       2,
		Code:       "TEST_FIELD_2",
	}
	if !reflect.DeepEqual(fields[1], expectedField2) {
		t.Errorf("Ожидалось поле %+v, получено %+v", expectedField2, fields[1])
	}
}

func TestGetCatalogCustomField(t *testing.T) {
	catalogID := 123
	fieldID := 456

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/catalogs/%d/custom_fields/%d", catalogID, fieldID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 456,
			"name": "Тестовое поле",
			"type": "text",
			"is_api_only": false,
			"is_required": false,
			"is_multiple": false,
			"is_system": false,
			"sort": 1,
			"code": "TEST_FIELD"
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	field, err := GetCatalogCustomField(apiClient, catalogID, fieldID)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при получении поля каталога: %v", err)
	}

	// Проверяем содержимое поля
	expectedField := &CustomField{
		ID:         456,
		Name:       "Тестовое поле",
		Type:       "text",
		IsAPIOnly:  false,
		IsRequired: false,
		IsMultiple: false,
		IsSystem:   false,
		Sort:       1,
		Code:       "TEST_FIELD",
	}
	if !reflect.DeepEqual(field, expectedField) {
		t.Errorf("Ожидалось поле %+v, получено %+v", expectedField, field)
	}
}

func TestUpdateCatalogCustomField(t *testing.T) {
	catalogID := 123
	fieldID := 456

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "PATCH" {
			t.Errorf("Ожидался метод PATCH, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/catalogs/%d/custom_fields/%d", catalogID, fieldID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 456,
			"name": "Обновленное поле",
			"type": "text",
			"is_api_only": false,
			"is_required": true,
			"is_multiple": false,
			"is_system": false,
			"sort": 3,
			"code": "UPDATED_FIELD"
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Создаем поле для обновления
	fieldToUpdate := &CustomField{
		ID:         fieldID,
		Name:       "Обновленное поле",
		IsRequired: true,
		Sort:       3,
		Code:       "UPDATED_FIELD",
	}

	// Вызываем тестируемый метод
	updatedField, err := UpdateCatalogCustomField(apiClient, catalogID, fieldToUpdate)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при обновлении поля каталога: %v", err)
	}

	// Проверяем содержимое обновленного поля
	expectedField := &CustomField{
		ID:         456,
		Name:       "Обновленное поле",
		Type:       "text",
		IsAPIOnly:  false,
		IsRequired: true,
		IsMultiple: false,
		IsSystem:   false,
		Sort:       3,
		Code:       "UPDATED_FIELD",
	}
	if !reflect.DeepEqual(updatedField, expectedField) {
		t.Errorf("Ожидалось поле %+v, получено %+v", expectedField, updatedField)
	}
}

func TestDeleteCatalogCustomField(t *testing.T) {
	catalogID := 123
	fieldID := 456

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "DELETE" {
			t.Errorf("Ожидался метод DELETE, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/catalogs/%d/custom_fields/%d", catalogID, fieldID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	err := DeleteCatalogCustomField(apiClient, catalogID, fieldID)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при удалении поля каталога: %v", err)
	}
}
