package catalog_elements

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

func TestGetCatalogElements(t *testing.T) {
	catalogID := 123
	
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/catalogs/%d/elements", catalogID)
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

	// Вызываем тестируемый метод
	elements, err := GetCatalogElements(apiClient, catalogID, 1, 50, nil)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при получении элементов каталога: %v", err)
	}

	if len(elements) != 2 {
		t.Fatalf("Ожидалось получение 2 элементов, получено %d", len(elements))
	}

	// Проверяем содержимое первого элемента
	expectedElement1 := CatalogElement{
		ID:        456,
		Name:      "Тестовый элемент 1",
		CreatedBy: 789,
		UpdatedBy: 789,
		CreatedAt: 1609459200,
		UpdatedAt: 1609545600,
		CatalogID: 123,
		CustomFieldsValues: []CustomFieldValue{
			{
				FieldID:   101,
				FieldName: "Код",
				FieldCode: "CODE",
				FieldType: "text",
				Values: []FieldValueItem{
					{
						Value: "EL-001",
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(elements[0], expectedElement1) {
		t.Errorf("Ожидался элемент %+v, получен %+v", expectedElement1, elements[0])
	}
}

func TestCreateCatalogElement(t *testing.T) {
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
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{
			"_embedded": {
				"elements": [
					{
						"id": 456,
						"name": "Новый элемент",
						"created_by": 789,
						"updated_by": 789,
						"created_at": 1609459200,
						"updated_at": 1609459200,
						"catalog_id": 123,
						"custom_fields_values": [
							{
								"field_id": 101,
								"field_name": "Код",
								"field_code": "CODE",
								"field_type": "text",
								"values": [
									{
										"value": "NEW-001"
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

	// Создаем элемент для отправки
	newElement := &CatalogElement{
		Name:      "Новый элемент",
		CatalogID: catalogID,
		CustomFieldsValues: []CustomFieldValue{
			{
				FieldID: 101,
				Values: []FieldValueItem{
					{
						Value: "NEW-001",
					},
				},
			},
		},
	}

	// Вызываем тестируемый метод
	createdElement, err := CreateCatalogElement(apiClient, catalogID, newElement)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при создании элемента каталога: %v", err)
	}

	// Проверяем содержимое созданного элемента
	expectedElement := &CatalogElement{
		ID:        456,
		Name:      "Новый элемент",
		CreatedBy: 789,
		UpdatedBy: 789,
		CreatedAt: 1609459200,
		UpdatedAt: 1609459200,
		CatalogID: 123,
		CustomFieldsValues: []CustomFieldValue{
			{
				FieldID:   101,
				FieldName: "Код",
				FieldCode: "CODE",
				FieldType: "text",
				Values: []FieldValueItem{
					{
						Value: "NEW-001",
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(createdElement, expectedElement) {
		t.Errorf("Ожидался элемент %+v, получен %+v", expectedElement, createdElement)
	}
}

func TestGetCatalogElement(t *testing.T) {
	catalogID := 123
	elementID := 456
	
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// Проверяем метод запроса
if r.Method != "GET" {
t.Errorf("Ожидался метод GET, получен %s", r.Method)
}

// Проверяем путь запроса
expectedPath := fmt.Sprintf("/api/v4/catalogs/%d/elements/%d", catalogID, elementID)
if r.URL.Path != expectedPath {
t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
}

// Отправляем ответ
w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 456,
			"name": "Тестовый элемент",
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
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	element, err := GetCatalogElement(apiClient, catalogID, elementID)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при получении элемента каталога: %v", err)
	}

	// Проверяем содержимое элемента
	expectedElement := &CatalogElement{
		ID:        456,
		Name:      "Тестовый элемент",
		CreatedBy: 789,
		UpdatedBy: 789,
		CreatedAt: 1609459200,
		UpdatedAt: 1609545600,
		CatalogID: 123,
		CustomFieldsValues: []CustomFieldValue{
			{
				FieldID:   101,
				FieldName: "Код",
				FieldCode: "CODE",
				FieldType: "text",
				Values: []FieldValueItem{
					{
						Value: "EL-001",
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(element, expectedElement) {
		t.Errorf("Ожидался элемент %+v, получен %+v", expectedElement, element)
	}
}

func TestUpdateCatalogElement(t *testing.T) {
	catalogID := 123
	elementID := 456
	
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// Проверяем метод запроса
if r.Method != "PATCH" {
t.Errorf("Ожидался метод PATCH, получен %s", r.Method)
}

// Проверяем путь запроса
expectedPath := fmt.Sprintf("/api/v4/catalogs/%d/elements/%d", catalogID, elementID)
if r.URL.Path != expectedPath {
t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
}

// Отправляем ответ
w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 456,
			"name": "Обновленный элемент",
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
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Создаем элемент для обновления
	elementToUpdate := &CatalogElement{
		ID:        elementID,
		Name:      "Обновленный элемент",
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
	}

	// Вызываем тестируемый метод
	updatedElement, err := UpdateCatalogElement(apiClient, catalogID, elementToUpdate)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при обновлении элемента каталога: %v", err)
	}

	// Проверяем содержимое обновленного элемента
	expectedElement := &CatalogElement{
		ID:        456,
		Name:      "Обновленный элемент",
		CreatedBy: 789,
		UpdatedBy: 789,
		CreatedAt: 1609459200,
		UpdatedAt: 1609632000,
		CatalogID: 123,
		CustomFieldsValues: []CustomFieldValue{
			{
				FieldID:   101,
				FieldName: "Код",
				FieldCode: "CODE",
				FieldType: "text",
				Values: []FieldValueItem{
					{
						Value: "EL-001-UPD",
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(updatedElement, expectedElement) {
		t.Errorf("Ожидался элемент %+v, получен %+v", expectedElement, updatedElement)
	}
}

func TestDeleteCatalogElement(t *testing.T) {
	catalogID := 123
	elementID := 456
	
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// Проверяем метод запроса
if r.Method != "DELETE" {
t.Errorf("Ожидался метод DELETE, получен %s", r.Method)
}

// Проверяем путь запроса
expectedPath := fmt.Sprintf("/api/v4/catalogs/%d/elements/%d", catalogID, elementID)
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
	err := DeleteCatalogElement(apiClient, catalogID, elementID)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при удалении элемента каталога: %v", err)
	}
}

func TestLinkCatalogElementWithTags(t *testing.T) {
	catalogID := 123
	elementID := 456
	
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// Проверяем метод запроса
if r.Method != "POST" {
t.Errorf("Ожидался метод POST, получен %s", r.Method)
}

// Проверяем путь запроса
expectedPath := fmt.Sprintf("/api/v4/catalogs/%d/elements/%d/tags", catalogID, elementID)
if r.URL.Path != expectedPath {
t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
}

// Отправляем ответ
w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"_embedded": {
				"tags": [
					{
						"id": 101,
						"name": "Тег 1",
						"color": "#FF0000"
					},
					{
						"id": 102,
						"name": "Тег 2",
						"color": "#00FF00"
					}
				]
			}
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Создаем теги для связывания
	tags := []Tag{
		{
			Name:  "Тег 1",
			Color: "#FF0000",
		},
		{
			Name:  "Тег 2",
			Color: "#00FF00",
		},
	}

	// Вызываем тестируемый метод
	err := LinkCatalogElementWithTags(apiClient, catalogID, elementID, tags)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при связывании элемента каталога с тегами: %v", err)
	}
}

func TestGetCatalogElementTags(t *testing.T) {
	catalogID := 123
	elementID := 456
	
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// Проверяем метод запроса
if r.Method != "GET" {
t.Errorf("Ожидался метод GET, получен %s", r.Method)
}

// Проверяем путь запроса
expectedPath := fmt.Sprintf("/api/v4/catalogs/%d/elements/%d/tags", catalogID, elementID)
if r.URL.Path != expectedPath {
t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
}

// Отправляем ответ
w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"_embedded": {
				"tags": [
					{
						"id": 101,
						"name": "Тег 1",
						"color": "#FF0000"
					},
					{
						"id": 102,
						"name": "Тег 2",
						"color": "#00FF00"
					}
				]
			}
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	tags, err := GetCatalogElementTags(apiClient, catalogID, elementID)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при получении тегов элемента каталога: %v", err)
	}

	if len(tags) != 2 {
		t.Fatalf("Ожидалось получение 2 тегов, получено %d", len(tags))
	}

	// Проверяем содержимое первого тега
	expectedTag1 := Tag{
		ID:    101,
		Name:  "Тег 1",
		Color: "#FF0000",
	}
	if !reflect.DeepEqual(tags[0], expectedTag1) {
		t.Errorf("Ожидался тег %+v, получен %+v", expectedTag1, tags[0])
	}

	// Проверяем содержимое второго тега
	expectedTag2 := Tag{
		ID:    102,
		Name:  "Тег 2",
		Color: "#00FF00",
	}
	if !reflect.DeepEqual(tags[1], expectedTag2) {
		t.Errorf("Ожидался тег %+v, получен %+v", expectedTag2, tags[1])
	}
}
