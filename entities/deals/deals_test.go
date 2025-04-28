package deals

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

func TestGetDeal(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/deals/123"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 123,
			"name": "Тестовая сделка",
			"value": 10000,
			"responsible_user_id": 456,
			"created_at": 1609459200,
			"updated_at": 1609545600
		}`))
	}))
	defer server.Close()

	// Создаем клиент
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	deal, err := GetDeal(apiClient, 123)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при получении сделки: %v", err)
	}

	if deal.ID != 123 {
		t.Errorf("Ожидался ID сделки 123, получен %d", deal.ID)
	}

	if deal.Name != "Тестовая сделка" {
		t.Errorf("Ожидалось имя сделки 'Тестовая сделка', получено '%s'", deal.Name)
	}

	if deal.Value != 10000 {
		t.Errorf("Ожидалась сумма сделки 10000, получено %d", deal.Value)
	}

	if deal.ResponsibleUserID != 456 {
		t.Errorf("Ожидался ID ответственного пользователя 456, получен %d", deal.ResponsibleUserID)
	}
}

func TestCreateDeal(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/deals"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Декодируем тело запроса
		var requestDeal Deal
		if err := json.NewDecoder(r.Body).Decode(&requestDeal); err != nil {
			t.Fatalf("Ошибка при декодировании тела запроса: %v", err)
		}

		// Проверяем содержимое запроса
		if requestDeal.Name != "Новая сделка" {
			t.Errorf("Ожидалось имя сделки 'Новая сделка', получено '%s'", requestDeal.Name)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 789,
			"name": "Новая сделка",
			"value": 15000,
			"responsible_user_id": 456,
			"created_at": 1609459200,
			"updated_at": 1609545600
		}`))
	}))
	defer server.Close()

	// Создаем клиент
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Создаем сделку для теста
	dealToCreate := &Deal{
		Name:              "Новая сделка",
		Value:             15000,
		ResponsibleUserID: 456,
	}

	// Вызываем тестируемый метод
	createdDeal, err := CreateDeal(apiClient, dealToCreate)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при создании сделки: %v", err)
	}

	if createdDeal.ID != 789 {
		t.Errorf("Ожидался ID сделки 789, получен %d", createdDeal.ID)
	}

	if createdDeal.Name != "Новая сделка" {
		t.Errorf("Ожидалось имя сделки 'Новая сделка', получено '%s'", createdDeal.Name)
	}

	if createdDeal.Value != 15000 {
		t.Errorf("Ожидалась сумма сделки 15000, получено %d", createdDeal.Value)
	}
}

func TestUpdateDeal(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "PATCH" {
			t.Errorf("Ожидался метод PATCH, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/deals/456"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Декодируем тело запроса
		var requestDeal Deal
		if err := json.NewDecoder(r.Body).Decode(&requestDeal); err != nil {
			t.Fatalf("Ошибка при декодировании тела запроса: %v", err)
		}

		// Проверяем содержимое запроса
		if requestDeal.Name != "Обновленная сделка" {
			t.Errorf("Ожидалось имя сделки 'Обновленная сделка', получено '%s'", requestDeal.Name)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 456,
			"name": "Обновленная сделка",
			"value": 25000,
			"responsible_user_id": 789,
			"created_at": 1609459200,
			"updated_at": 1609631999
		}`))
	}))
	defer server.Close()

	// Создаем клиент
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Создаем сделку для обновления
	dealToUpdate := &Deal{
		ID:                456,
		Name:              "Обновленная сделка",
		Value:             25000,
		ResponsibleUserID: 789,
	}

	// Вызываем тестируемый метод
	updatedDeal, err := UpdateDeal(apiClient, dealToUpdate)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при обновлении сделки: %v", err)
	}

	if updatedDeal.ID != 456 {
		t.Errorf("Ожидался ID сделки 456, получен %d", updatedDeal.ID)
	}

	if updatedDeal.Name != "Обновленная сделка" {
		t.Errorf("Ожидалось имя сделки 'Обновленная сделка', получено '%s'", updatedDeal.Name)
	}

	if updatedDeal.Value != 25000 {
		t.Errorf("Ожидалась сумма сделки 25000, получено %d", updatedDeal.Value)
	}
}

func TestGetDeals(t *testing.T) {
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
			name:         "Успешное получение списка сделок",
			page:         1,
			limit:        50,
			responseCode: http.StatusOK,
			responseBody: `{"_embedded":{"items":[{"id":12345,"name":"Сделка 1","value":10000},{"id":12346,"name":"Сделка 2","value":20000}]}}`,
			expectedLen:  2,
		},
		{
			name:         "Пустой список сделок",
			page:         1,
			limit:        50,
			responseCode: http.StatusOK,
			responseBody: `{"_embedded":{"items":[]}}`,
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
				expectedPath := "/api/v4/deals"
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
			deals, err := GetDeals(apiClient, tt.page, tt.limit, nil)

			// Проверяем результаты
			if tt.expectError && err == nil {
				t.Error("Ожидалась ошибка, но ее не было")
			}

			if !tt.expectError {
				if err != nil {
					t.Errorf("Неожиданная ошибка: %v", err)
				}

				if len(deals) != tt.expectedLen {
					t.Errorf("Ожидалось %d сделок, получено %d", tt.expectedLen, len(deals))
				}
			}
		})
	}
}
