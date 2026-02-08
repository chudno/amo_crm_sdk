package contacts

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

func TestGetContact(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		expectedPath := "/api/v4/contacts/123"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 123,
			"name": "Тестовый контакт",
			"responsible_user_id": 456,
			"created_at": 1609459200,
			"updated_at": 1609545600
		}`))
	}))
	defer server.Close()

	apiClient := client.NewClient(server.URL, "test_api_key")

	contact, err := Get(context.Background(), apiClient, 123)

	if err != nil {
		t.Fatalf("Ошибка при получении контакта: %v", err)
	}

	if contact.ID != 123 {
		t.Errorf("Ожидался ID контакта 123, получен %d", contact.ID)
	}

	if contact.Name != "Тестовый контакт" {
		t.Errorf("Ожидалось имя контакта 'Тестовый контакт', получено '%s'", contact.Name)
	}

	if contact.ResponsibleUserID != 456 {
		t.Errorf("Ожидался ID ответственного пользователя 456, получен %d", contact.ResponsibleUserID)
	}
}

func TestCreateContact(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		expectedPath := "/api/v4/contacts"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"_embedded": {
				"contacts": [
					{
						"id": 789,
						"name": "Новый контакт",
						"responsible_user_id": 456,
						"created_at": 1609459200,
						"updated_at": 1609545600
					}
				]
			}
		}`))
	}))
	defer server.Close()

	apiClient := client.NewClient(server.URL, "test_api_key")

	contactToCreate := &Contact{
		Name:              "Новый контакт",
		ResponsibleUserID: 456,
	}

	createdContact, err := Create(context.Background(), apiClient, contactToCreate)

	if err != nil {
		t.Fatalf("Ошибка при создании контакта: %v", err)
	}

	if createdContact.ID != 789 {
		t.Errorf("Ожидался ID контакта 789, получен %d", createdContact.ID)
	}

	if createdContact.Name != "Новый контакт" {
		t.Errorf("Ожидалось имя контакта 'Новый контакт', получено '%s'", createdContact.Name)
	}
}

func TestListContacts(t *testing.T) {
	tests := []struct {
		name         string
		page         int
		limit        int
		responseCode int
		responseBody string
		expectError  bool
		expectedLen  int
	}{
		{
			name:         "Успешное получение списка контактов",
			page:         1,
			limit:        50,
			responseCode: http.StatusOK,
			responseBody: `{"_embedded":{"contacts":[{"id":12345,"name":"Иван Иванов"},{"id":12346,"name":"Петр Петров"}]}}`,
			expectedLen:  2,
		},
		{
			name:         "Пустой список контактов",
			page:         1,
			limit:        50,
			responseCode: http.StatusOK,
			responseBody: `{"_embedded":{"contacts":[]}}`,
			expectedLen:  0,
		},
		{
			name:         "Ошибка сервера",
			page:         1,
			limit:        50,
			responseCode: http.StatusInternalServerError,
			responseBody: `{"error":"server_error"}`,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Ожидался метод GET, получен %s", r.Method)
				}

				expectedPath := "/api/v4/contacts"
				if r.URL.Path != expectedPath {
					t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
				}

				query := r.URL.Query()
				if query.Get("page") != fmt.Sprintf("%d", tt.page) {
					t.Errorf("Ожидался параметр page=%d, получен %s", tt.page, query.Get("page"))
				}
				if query.Get("limit") != fmt.Sprintf("%d", tt.limit) {
					t.Errorf("Ожидался параметр limit=%d, получен %s", tt.limit, query.Get("limit"))
				}

				w.WriteHeader(tt.responseCode)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			apiClient := client.NewClient(server.URL, "test_api_key")

			contacts, err := List(context.Background(), apiClient, tt.page, tt.limit)

			if tt.expectError && err == nil {
				t.Error("Ожидалась ошибка, но ее не было")
			}

			if !tt.expectError {
				if err != nil {
					t.Errorf("Неожиданная ошибка: %v", err)
				}

				if len(contacts) != tt.expectedLen {
					t.Errorf("Ожидалось %d контактов, получено %d", tt.expectedLen, len(contacts))
				}
			}
		})
	}
}
