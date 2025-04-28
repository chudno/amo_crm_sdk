package contacts

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

func TestGetContact(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/contacts/123"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
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

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	contact, err := GetContact(apiClient, 123)

	// Проверяем результаты
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
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/contacts"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 789,
			"name": "Новый контакт",
			"responsible_user_id": 456,
			"created_at": 1609459200,
			"updated_at": 1609545600
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Создаем контакт для теста
	contactToCreate := &Contact{
		Name:              "Новый контакт",
		ResponsibleUserID: 456,
	}

	// Вызываем тестируемый метод
	createdContact, err := CreateContact(apiClient, contactToCreate)

	// Проверяем результаты
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
				// Проверяем метод запроса
				if r.Method != "GET" {
					t.Errorf("Ожидался метод GET, получен %s", r.Method)
				}

				// Проверяем URL запроса
				expectedPath := "/api/v4/contacts"
				if r.URL.Path != expectedPath {
					t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
				}

				// Проверяем параметры запроса
				query := r.URL.Query()
				if query.Get("page") != fmt.Sprintf("%d", tt.page) {
					t.Errorf("Ожидался параметр page=%d, получен %s", tt.page, query.Get("page"))
				}
				if query.Get("limit") != fmt.Sprintf("%d", tt.limit) {
					t.Errorf("Ожидался параметр limit=%d, получен %s", tt.limit, query.Get("limit"))
				}

				// Устанавливаем код ответа и тело
				w.WriteHeader(tt.responseCode)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			// Создаем клиент
			apiClient := client.NewClient(server.URL, "test_api_key")

			// Вызываем тестируемую функцию
			contacts, err := GetContacts(apiClient, tt.page, tt.limit)

			// Проверяем результаты
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
