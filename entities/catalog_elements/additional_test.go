package catalog_elements

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

func TestCreateCatalogElements(t *testing.T) {
	catalogID := 123

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/catalogs/%d/elements", catalogID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"_embedded": {
				"elements": [
					{
						"id": 456,
						"name": "Тестовый элемент 1",
						"created_by": 789,
						"updated_by": 789,
						"created_at": 1609459200,
						"updated_at": 1609545600,
						"catalog_id": 123,
						"custom_fields_values": [
							{
								"field_id": 101,
								"field_name": "Код",
								"field_code": "CODE",
								"field_type": "text",
								"values": [
									{
										"value": "EL-001"
									}
								]
							}
						]
					},
					{
						"id": 789,
						"name": "Тестовый элемент 2",
						"created_by": 789,
						"updated_by": 789,
						"created_at": 1609459200,
						"updated_at": 1609545600,
						"catalog_id": 123,
						"custom_fields_values": [
							{
								"field_id": 101,
								"field_name": "Код",
								"field_code": "CODE",
								"field_type": "text",
								"values": [
									{
										"value": "EL-002"
									}
								]
							}
						]
					}
				]
			}
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Создаем элементы для добавления
	elements := []CatalogElement{
		{
			Name:      "Тестовый элемент 1",
			CatalogID: catalogID,
			CustomFieldsValues: []CustomFieldValue{
				{
					FieldID: 101,
					Values: []FieldValueItem{
						{
							Value: "EL-001",
						},
					},
				},
			},
		},
		{
			Name:      "Тестовый элемент 2",
			CatalogID: catalogID,
			CustomFieldsValues: []CustomFieldValue{
				{
					FieldID: 101,
					Values: []FieldValueItem{
						{
							Value: "EL-002",
						},
					},
				},
			},
		},
	}

	// Вызываем тестируемый метод
	createdElements, err := CreateCatalogElements(apiClient, catalogID, elements)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при создании элементов каталога: %v", err)
	}

	if len(createdElements) != 2 {
		t.Fatalf("Ожидалось создание 2 элементов, получено %d", len(createdElements))
	}

	// Проверяем ID созданных элементов
	if createdElements[0].ID != 456 {
		t.Errorf("Ожидался ID 456 для первого элемента, получен %d", createdElements[0].ID)
	}

	if createdElements[1].ID != 789 {
		t.Errorf("Ожидался ID 789 для второго элемента, получен %d", createdElements[1].ID)
	}
}

func TestUpdateCatalogElements(t *testing.T) {
	catalogID := 123

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "PATCH" {
			t.Errorf("Ожидался метод PATCH, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/catalogs/%d/elements", catalogID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"_embedded": {
				"elements": [
					{
						"id": 456,
						"name": "Обновленный элемент 1",
						"created_by": 789,
						"updated_by": 789,
						"created_at": 1609459200,
						"updated_at": 1609632000,
						"catalog_id": 123,
						"custom_fields_values": [
							{
								"field_id": 101,
								"field_name": "Код",
								"field_code": "CODE",
								"field_type": "text",
								"values": [
									{
										"value": "EL-001-UPD"
									}
								]
							}
						]
					},
					{
						"id": 789,
						"name": "Обновленный элемент 2",
						"created_by": 789,
						"updated_by": 789,
						"created_at": 1609459200,
						"updated_at": 1609632000,
						"catalog_id": 123,
						"custom_fields_values": [
							{
								"field_id": 101,
								"field_name": "Код",
								"field_code": "CODE",
								"field_type": "text",
								"values": [
									{
										"value": "EL-002-UPD"
									}
								]
							}
						]
					}
				]
			}
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Создаем элементы для обновления
	elements := []CatalogElement{
		{
			ID:        456,
			Name:      "Обновленный элемент 1",
			CatalogID: catalogID,
			CustomFieldsValues: []CustomFieldValue{
				{
					FieldID: 101,
					Values: []FieldValueItem{
						{
							Value: "EL-001-UPD",
						},
					},
				},
			},
		},
		{
			ID:        789,
			Name:      "Обновленный элемент 2",
			CatalogID: catalogID,
			CustomFieldsValues: []CustomFieldValue{
				{
					FieldID: 101,
					Values: []FieldValueItem{
						{
							Value: "EL-002-UPD",
						},
					},
				},
			},
		},
	}

	// Вызываем тестируемый метод
	updatedElements, err := UpdateCatalogElements(apiClient, catalogID, elements)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при обновлении элементов каталога: %v", err)
	}

	if len(updatedElements) != 2 {
		t.Fatalf("Ожидалось обновление 2 элементов, получено %d", len(updatedElements))
	}

	// Проверяем имена обновленных элементов
	if updatedElements[0].Name != "Обновленный элемент 1" {
		t.Errorf("Ожидалось имя 'Обновленный элемент 1', получено '%s'", updatedElements[0].Name)
	}

	if updatedElements[1].Name != "Обновленный элемент 2" {
		t.Errorf("Ожидалось имя 'Обновленный элемент 2', получено '%s'", updatedElements[1].Name)
	}

	// Проверяем значения пользовательских полей
	if len(updatedElements[0].CustomFieldsValues) != 1 ||
		len(updatedElements[0].CustomFieldsValues[0].Values) != 1 ||
		updatedElements[0].CustomFieldsValues[0].Values[0].Value != "EL-001-UPD" {
		t.Errorf("Неверное значение пользовательского поля для первого элемента")
	}

	if len(updatedElements[1].CustomFieldsValues) != 1 ||
		len(updatedElements[1].CustomFieldsValues[0].Values) != 1 ||
		updatedElements[1].CustomFieldsValues[0].Values[0].Value != "EL-002-UPD" {
		t.Errorf("Неверное значение пользовательского поля для второго элемента")
	}
}

func TestBatchDeleteCatalogElements(t *testing.T) {
	catalogID := 123
	elementIDs := []int{456, 789}

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "DELETE" {
			t.Errorf("Ожидался метод DELETE, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/catalogs/%d/elements", catalogID)
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
	err := BatchDeleteCatalogElements(apiClient, catalogID, elementIDs)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при пакетном удалении элементов каталога: %v", err)
	}
}
