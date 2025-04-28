package leads

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

func TestGetLead(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/leads/123"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 123,
			"name": "Тестовый лид",
			"price": 10000,
			"responsible_user_id": 456,
			"created_at": 1609459200,
			"updated_at": 1609545600
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	lead, err := GetLead(apiClient, 123)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при получении лида: %v", err)
	}

	if lead.ID != 123 {
		t.Errorf("Ожидался ID лида 123, получен %d", lead.ID)
	}

	if lead.Name != "Тестовый лид" {
		t.Errorf("Ожидалось имя лида 'Тестовый лид', получено '%s'", lead.Name)
	}

	if lead.Price != 10000 {
		t.Errorf("Ожидалась цена лида 10000, получена %d", lead.Price)
	}

	if lead.ResponsibleUserID != 456 {
		t.Errorf("Ожидался ID ответственного пользователя 456, получен %d", lead.ResponsibleUserID)
	}
}

func TestCreateLead(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/leads"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"_embedded": {
				"leads": [{
					"id": 789,
					"name": "Новый лид",
					"price": 15000,
					"responsible_user_id": 456,
					"created_at": 1609459200,
					"updated_at": 1609545600
				}]
			}
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Создаем лид для теста
	leadToCreate := &Lead{
		Name:              "Новый лид",
		Price:             15000,
		ResponsibleUserID: 456,
	}

	// Вызываем тестируемый метод
	createdLead, err := CreateLead(apiClient, leadToCreate)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при создании лида: %v", err)
	}

	if createdLead.ID != 789 {
		t.Errorf("Ожидался ID лида 789, получен %d", createdLead.ID)
	}

	if createdLead.Name != "Новый лид" {
		t.Errorf("Ожидалось имя лида 'Новый лид', получено '%s'", createdLead.Name)
	}

	if createdLead.Price != 15000 {
		t.Errorf("Ожидалась цена лида 15000, получена %d", createdLead.Price)
	}
}

func TestListLeads(t *testing.T) {
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
			name:         "Успешное получение списка лидов",
			page:         1,
			limit:        50,
			responseCode: http.StatusOK,
			responseBody: `{"_embedded":{"leads":[{"id":12345,"name":"Лид 1","price":10000},{"id":12346,"name":"Лид 2","price":20000}]}}`,
			expectedLen:  2,
		},
		{
			name:         "Пустой список лидов",
			page:         1,
			limit:        50,
			responseCode: http.StatusOK,
			responseBody: `{"_embedded":{"leads":[]}}`,
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
				expectedPath := "/api/v4/leads"
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
			leads, err := ListLeads(apiClient, tt.limit, tt.page, nil)

			// Проверяем результаты
			if tt.expectError && err == nil {
				t.Error("Ожидалась ошибка, но ее не было")
			}

			if !tt.expectError {
				if err != nil {
					t.Errorf("Неожиданная ошибка: %v", err)
				}

				if len(leads) != tt.expectedLen {
					t.Errorf("Ожидалось %d лидов, получено %d", tt.expectedLen, len(leads))
				}
			}
		})
	}
}
