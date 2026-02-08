package segments

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

// TestAddSegment проверяет создание сегмента
func TestAddSegment(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		expectedPath := "/api/v4/segments"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{
			"_embedded": {
				"segments": [
					{
						"id": 123,
						"name": "Тестовый сегмент",
						"color": "#FF5555",
						"type": "dynamic",
						"created_by": 789,
						"updated_by": 789,
						"created_at": 1609459200,
						"updated_at": 1609459200,
						"account_id": 12345,
						"contacts_count": 0,
						"is_deleted": false,
						"_links": {
							"self": {
								"href": "/api/v4/segments/123"
							}
						}
					}
				]
			}
		}`))
	}))
	defer server.Close()

	apiClient := client.NewClient(server.URL, "test_api_key")

	segment := &Segment{
		Name:  "Тестовый сегмент",
		Color: "#FF5555",
		Type:  SegmentTypeDynamic,
		Filter: &Filter{
			Logic: "and",
			Nodes: []FilterNode{
				{
					FieldCode: "email",
					Operator:  "contains",
					Value:     "example.com",
				},
			},
		},
	}

	createdSegment, err := Create(context.Background(), apiClient, segment)

	if err != nil {
		t.Fatalf("Ошибка при создании сегмента: %v", err)
	}

	if createdSegment.ID != 123 {
		t.Errorf("Ожидался ID 123, получен %d", createdSegment.ID)
	}

	if createdSegment.Name != "Тестовый сегмент" {
		t.Errorf("Ожидалось имя 'Тестовый сегмент', получено '%s'", createdSegment.Name)
	}

	if createdSegment.Color != "#FF5555" {
		t.Errorf("Ожидался цвет '#FF5555', получен '%s'", createdSegment.Color)
	}

	if createdSegment.Type != SegmentTypeDynamic {
		t.Errorf("Ожидался тип сегмента 'dynamic', получен '%s'", createdSegment.Type)
	}
}

// TestGetSegment проверяет получение информации о сегменте
func TestGetSegment(t *testing.T) {
	segmentID := 123

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		expectedPath := fmt.Sprintf("/api/v4/segments/%d", segmentID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		expectedWith := "contacts"
		if r.URL.Query().Get("with") != expectedWith {
			t.Errorf("Ожидался параметр with=%s, получен %s", expectedWith, r.URL.Query().Get("with"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 123,
			"name": "Активные клиенты",
			"color": "#FF5555",
			"type": "dynamic",
			"created_by": 789,
			"updated_by": 789,
			"created_at": 1609459200,
			"updated_at": 1609459200,
			"account_id": 12345,
			"contacts_count": 42,
			"is_deleted": false,
			"_embedded": {
				"contacts": [
					{
						"id": 1001,
						"name": "Иван Иванов",
						"_links": {
							"self": {
								"href": "/api/v4/contacts/1001"
							}
						}
					},
					{
						"id": 1002,
						"name": "Петр Петров",
						"_links": {
							"self": {
								"href": "/api/v4/contacts/1002"
							}
						}
					}
				]
			},
			"_links": {
				"self": {
					"href": "/api/v4/segments/123"
				}
			}
		}`))
	}))
	defer server.Close()

	apiClient := client.NewClient(server.URL, "test_api_key")

	segment, err := Get(context.Background(), apiClient, segmentID, WithContacts())

	if err != nil {
		t.Fatalf("Ошибка при получении сегмента: %v", err)
	}

	if segment.ID != segmentID {
		t.Errorf("Ожидался ID %d, получен %d", segmentID, segment.ID)
	}

	if segment.Name != "Активные клиенты" {
		t.Errorf("Ожидалось имя 'Активные клиенты', получено '%s'", segment.Name)
	}

	if segment.ContactsCount != 42 {
		t.Errorf("Ожидалось количество контактов 42, получено %d", segment.ContactsCount)
	}

	if segment.Embedded == nil {
		t.Fatalf("Отсутствует секция _embedded")
	}

	if len(segment.Embedded.Contacts) != 2 {
		t.Fatalf("Ожидалось 2 контакта, получено %d", len(segment.Embedded.Contacts))
	}

	if segment.Embedded.Contacts[0].ID != 1001 {
		t.Errorf("Ожидался ID контакта 1001, получен %d", segment.Embedded.Contacts[0].ID)
	}

	if segment.Embedded.Contacts[0].Name != "Иван Иванов" {
		t.Errorf("Ожидалось имя контакта 'Иван Иванов', получено '%s'", segment.Embedded.Contacts[0].Name)
	}
}

// Используем интерфейс Requester из segments_test_helpers.go

// TestGetSegments проверяет получение списка сегментов
func TestGetSegments(t *testing.T) {
	successResponse := `{
		"page": 1,
		"per_page": 50,
		"_embedded": {
			"segments": [
				{
					"id": 123,
					"name": "Активные клиенты",
					"color": "#FF5555",
					"type": "dynamic",
					"created_by": 789,
					"updated_by": 789,
					"created_at": 1609459200,
					"updated_at": 1609459200,
					"account_id": 12345,
					"contacts_count": 42,
					"is_deleted": false
				},
				{
					"id": 456,
					"name": "Потенциальные клиенты",
					"color": "#55FF55",
					"type": "dynamic",
					"created_by": 789,
					"updated_by": 789,
					"created_at": 1609459200,
					"updated_at": 1609459200,
					"account_id": 12345,
					"contacts_count": 18,
					"is_deleted": false
				}
			]
		}
	}`

	emptyResponse := `{
		"page": 1,
		"per_page": 50,
		"_embedded": {
			"segments": []
		}
	}`

	t.Run("Success", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("GET", "/api/v4/segments", http.StatusOK, successResponse, nil)

		filter := map[string]string{
			"filter[name]": "Активные клиенты",
		}

		segments, err := GetSegmentsWithRequester(mockClient, 1, 50, WithFilter(filter))

		if err != nil {
			t.Fatalf("Ошибка при получении сегментов: %v", err)
		}

		if len(segments) != 2 {
			t.Fatalf("Ожидалось получение 2 сегментов, получено %d", len(segments))
		}

		if segments[0].ID != 123 {
			t.Errorf("Ожидался ID 123, получен %d", segments[0].ID)
		}

		if segments[0].Name != "Активные клиенты" {
			t.Errorf("Ожидалось имя 'Активные клиенты', получено '%s'", segments[0].Name)
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("GET", "/api/v4/segments", http.StatusOK, emptyResponse, nil)

		segments, err := GetSegmentsWithRequester(mockClient, 1, 50)

		if err != nil {
			t.Fatalf("Ошибка при получении сегментов: %v", err)
		}

		if len(segments) != 0 {
			t.Fatalf("Ожидался пустой массив сегментов, получено %d", len(segments))
		}
	})

	t.Run("ServerError", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("GET", "/api/v4/segments", http.StatusInternalServerError, `{"error": "Internal Server Error"}`, nil)

		_, err := GetSegmentsWithRequester(mockClient, 1, 50)

		if err == nil {
			t.Fatalf("Ожидалась ошибка, но её не получили")
		}
	})
}

// TestUpdateSegment проверяет обновление сегмента
func TestUpdateSegment(t *testing.T) {
	segmentID := 123

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("Ожидался метод PATCH, получен %s", r.Method)
		}

		expectedPath := fmt.Sprintf("/api/v4/segments/%d", segmentID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 123,
			"name": "Обновленный сегмент",
			"color": "#5555FF",
			"type": "dynamic",
			"created_by": 789,
			"updated_by": 789,
			"created_at": 1609459200,
			"updated_at": 1609545600,
			"account_id": 12345,
			"contacts_count": 42,
			"is_deleted": false,
			"_links": {
				"self": {
					"href": "/api/v4/segments/123"
				}
			}
		}`))
	}))
	defer server.Close()

	apiClient := client.NewClient(server.URL, "test_api_key")

	segment := &Segment{
		ID:    segmentID,
		Name:  "Обновленный сегмент",
		Color: "#5555FF",
	}

	updatedSegment, err := Update(context.Background(), apiClient, segment)

	if err != nil {
		t.Fatalf("Ошибка при обновлении сегмента: %v", err)
	}

	if updatedSegment.ID != segmentID {
		t.Errorf("Ожидался ID %d, получен %d", segmentID, updatedSegment.ID)
	}

	if updatedSegment.Name != "Обновленный сегмент" {
		t.Errorf("Ожидалось имя 'Обновленный сегмент', получено '%s'", updatedSegment.Name)
	}

	if updatedSegment.Color != "#5555FF" {
		t.Errorf("Ожидался цвет '#5555FF', получен '%s'", updatedSegment.Color)
	}
}

// TestDeleteSegment проверяет удаление сегмента
func TestDeleteSegment(t *testing.T) {
	segmentID := 123

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Ожидался метод DELETE, получен %s", r.Method)
		}

		expectedPath := fmt.Sprintf("/api/v4/segments/%d", segmentID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	apiClient := client.NewClient(server.URL, "test_api_key")

	err := Delete(context.Background(), apiClient, segmentID)

	if err != nil {
		t.Fatalf("Ошибка при удалении сегмента: %v", err)
	}
}

// TestAddContactsToSegment проверяет добавление контактов в сегмент
func TestAddContactsToSegment(t *testing.T) {
	segmentID := 123
	contactIDs := []int{1001, 1002, 1003}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		expectedPath := fmt.Sprintf("/api/v4/segments/%d/contacts", segmentID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	apiClient := client.NewClient(server.URL, "test_api_key")

	err := AddContacts(context.Background(), apiClient, segmentID, contactIDs)

	if err != nil {
		t.Fatalf("Ошибка при добавлении контактов в сегмент: %v", err)
	}
}

// TestRemoveContactsFromSegment проверяет удаление контактов из сегмента
func TestRemoveContactsFromSegment(t *testing.T) {
	segmentID := 123
	contactIDs := []int{1001, 1002, 1003}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		expectedPath := fmt.Sprintf("/api/v4/segments/%d/contacts/delete", segmentID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	apiClient := client.NewClient(server.URL, "test_api_key")

	err := RemoveContacts(context.Background(), apiClient, segmentID, contactIDs)

	if err != nil {
		t.Fatalf("Ошибка при удалении контактов из сегмента: %v", err)
	}
}

// TestGetSegmentContacts проверяет получение контактов сегмента
func TestGetSegmentContacts(t *testing.T) {
	segmentID := 123

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		expectedPath := fmt.Sprintf("/api/v4/segments/%d/contacts", segmentID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"page": 1,
			"per_page": 50,
			"_embedded": {
				"contacts": [
					{
						"id": 1001
					},
					{
						"id": 1002
					},
					{
						"id": 1003
					}
				]
			},
			"_links": {
				"self": {
					"href": "/api/v4/segments/123/contacts"
				}
			}
		}`))
	}))
	defer server.Close()

	apiClient := client.NewClient(server.URL, "test_api_key")

	contactIDs, err := ListContacts(context.Background(), apiClient, segmentID, 1, 50)

	if err != nil {
		t.Fatalf("Ошибка при получении контактов сегмента: %v", err)
	}

	if len(contactIDs) != 3 {
		t.Fatalf("Ожидалось получение 3 ID контактов, получено %d", len(contactIDs))
	}

	expectedIDs := []int{1001, 1002, 1003}
	for i, id := range expectedIDs {
		if contactIDs[i] != id {
			t.Errorf("Ожидался ID контакта %d, получен %d", id, contactIDs[i])
		}
	}
}
