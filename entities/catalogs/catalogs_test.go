package catalogs

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

func TestGetCatalogs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		expectedPath := "/api/v4/catalogs"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		expectedPage := "1"
		if r.URL.Query().Get("page") != expectedPage {
			t.Errorf("Ожидался параметр page=%s, получен %s", expectedPage, r.URL.Query().Get("page"))
		}

		expectedLimit := "50"
		if r.URL.Query().Get("limit") != expectedLimit {
			t.Errorf("Ожидался параметр limit=%s, получен %s", expectedLimit, r.URL.Query().Get("limit"))
		}

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

	apiClient := client.NewClient(server.URL, "test_api_key")

	catalogs, err := List(context.Background(), apiClient, 1, 50, nil)

	if err != nil {
		t.Fatalf("Ошибка при получении каталогов: %v", err)
	}

	if len(catalogs) != 2 {
		t.Fatalf("Ожидалось получение 2 каталогов, получено %d", len(catalogs))
	}

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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		expectedPath := "/api/v4/catalogs"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"_embedded": {
				"catalogs": [
					{
						"id": 789,
						"name": "Новый каталог",
						"created_by": 456,
						"updated_by": 456,
						"created_at": 1609459200,
						"updated_at": 1609459200,
						"sort": 3,
						"type": "regular"
					}
				]
			}
		}`))
	}))
	defer server.Close()

	apiClient := client.NewClient(server.URL, "test_api_key")

	newCatalog := &Catalog{
		Name: "Новый каталог",
		Sort: 3,
		Type: "regular",
	}

	createdCatalog, err := Create(context.Background(), apiClient, newCatalog)

	if err != nil {
		t.Fatalf("Ошибка при создании каталога: %v", err)
	}

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

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		expectedPath := fmt.Sprintf("/api/v4/catalogs/%d", catalogID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

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

	apiClient := client.NewClient(server.URL, "test_api_key")

	catalog, err := Get(context.Background(), apiClient, catalogID)

	if err != nil {
		t.Fatalf("Ошибка при получении каталога: %v", err)
	}

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

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("Ожидался метод PATCH, получен %s", r.Method)
		}

		expectedPath := fmt.Sprintf("/api/v4/catalogs/%d", catalogID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

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

	apiClient := client.NewClient(server.URL, "test_api_key")

	catalogToUpdate := &Catalog{
		ID:   catalogID,
		Name: "Обновленный каталог",
		Sort: 5,
	}

	updatedCatalog, err := Update(context.Background(), apiClient, catalogToUpdate)

	if err != nil {
		t.Fatalf("Ошибка при обновлении каталога: %v", err)
	}

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

func TestAddCustomFieldToCatalog(t *testing.T) {
	catalogID := 123

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		expectedPath := fmt.Sprintf("/api/v4/catalogs/%d/custom_fields", catalogID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

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

	apiClient := client.NewClient(server.URL, "test_api_key")

	newField := &CustomField{
		Name:       "Тестовое поле",
		Type:       "text",
		IsRequired: false,
		IsMultiple: false,
		Sort:       1,
		Code:       "TEST_FIELD",
	}

	createdField, err := AddCustomField(context.Background(), apiClient, catalogID, newField)

	if err != nil {
		t.Fatalf("Ошибка при добавлении поля в каталог: %v", err)
	}

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

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		expectedPath := fmt.Sprintf("/api/v4/catalogs/%d/custom_fields", catalogID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

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

	apiClient := client.NewClient(server.URL, "test_api_key")

	fields, err := ListCustomFields(context.Background(), apiClient, catalogID)

	if err != nil {
		t.Fatalf("Ошибка при получении полей каталога: %v", err)
	}

	if len(fields) != 2 {
		t.Fatalf("Ожидалось получение 2 полей, получено %d", len(fields))
	}

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

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		expectedPath := fmt.Sprintf("/api/v4/catalogs/%d/custom_fields/%d", catalogID, fieldID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

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

	apiClient := client.NewClient(server.URL, "test_api_key")

	field, err := GetCustomField(context.Background(), apiClient, catalogID, fieldID)

	if err != nil {
		t.Fatalf("Ошибка при получении поля каталога: %v", err)
	}

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

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("Ожидался метод PATCH, получен %s", r.Method)
		}

		expectedPath := fmt.Sprintf("/api/v4/catalogs/%d/custom_fields/%d", catalogID, fieldID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

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

	apiClient := client.NewClient(server.URL, "test_api_key")

	fieldToUpdate := &CustomField{
		ID:         fieldID,
		Name:       "Обновленное поле",
		IsRequired: true,
		Sort:       3,
		Code:       "UPDATED_FIELD",
	}

	updatedField, err := UpdateCustomField(context.Background(), apiClient, catalogID, fieldToUpdate)

	if err != nil {
		t.Fatalf("Ошибка при обновлении поля каталога: %v", err)
	}

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

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Ожидался метод DELETE, получен %s", r.Method)
		}

		expectedPath := fmt.Sprintf("/api/v4/catalogs/%d/custom_fields/%d", catalogID, fieldID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	apiClient := client.NewClient(server.URL, "test_api_key")

	err := DeleteCustomField(context.Background(), apiClient, catalogID, fieldID)

	if err != nil {
		t.Fatalf("Ошибка при удалении поля каталога: %v", err)
	}
}
