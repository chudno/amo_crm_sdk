package widgets

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

// TestGetMarketplaceWidgets проверяет получение списка виджетов из маркетплейса
func TestGetMarketplaceWidgets(t *testing.T) {
	// Подготавливаем ответ для успешного сценария
	successResponse := `{
		"page": 1,
		"per_page": 50,
		"_embedded": {
			"widgets": [
				{
					"id": 123,
					"code": "intercom",
					"name": "Intercom",
					"description": "Чат для вашего сайта",
					"logo_url": "https://example.com/logo1.png",
					"url": "https://example.com/intercom",
					"installed": false,
					"developer": {
						"id": 10,
						"name": "amoCRM"
					},
					"categories": [
						{
							"id": 1,
							"name": "Коммуникации"
						}
					],
					"version": "1.0",
					"pricing": "free",
					"rating": 4.5,
					"reviews_count": 120
				},
				{
					"id": 456,
					"code": "callback",
					"name": "Callback",
					"description": "Обратный звонок для вашего сайта",
					"logo_url": "https://example.com/logo2.png",
					"url": "https://example.com/callback",
					"installed": true,
					"developer": {
						"id": 10,
						"name": "amoCRM"
					},
					"categories": [
						{
							"id": 1,
							"name": "Коммуникации"
						},
						{
							"id": 2,
							"name": "Телефония"
						}
					],
					"version": "2.1",
					"pricing": "paid",
					"rating": 4.2,
					"reviews_count": 85
				}
			]
		}
	}`

	// Ответ для ситуации, когда виджетов нет
	emptyResponse := `{
		"page": 1,
		"per_page": 50,
		"_embedded": {
			"widgets": []
		}
	}`

	// Проверяем успешный сценарий
	t.Run("Success", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("GET", "/api/v4/marketplace/widgets", http.StatusOK, successResponse, nil)

		// Создаем фильтр по категории
		categoryID := 1

		// Вызываем тестируемый метод
		widgets, err := GetMarketplaceWidgetsWithRequester(mockClient, 1, 50, WithCategory(categoryID))

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при получении виджетов из маркетплейса: %v", err)
		}

		if len(widgets) != 2 {
			t.Errorf("Ожидалось 2 виджета, получено %d", len(widgets))
		}

		if widgets[0].Code != "intercom" {
			t.Errorf("Ожидался код intercom, получен %s", widgets[0].Code)
		}

		if widgets[1].Code != "callback" {
			t.Errorf("Ожидался код callback, получен %s", widgets[1].Code)
		}

		if !widgets[1].Installed {
			t.Errorf("Ожидалось, что виджет callback установлен")
		}

		if widgets[0].Developer.Name != "amoCRM" {
			t.Errorf("Ожидался разработчик amoCRM, получен %s", widgets[0].Developer.Name)
		}
	})

	// Проверяем сценарий с пустым ответом
	t.Run("Empty", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("GET", "/api/v4/marketplace/widgets", http.StatusOK, emptyResponse, nil)

		// Вызываем тестируемый метод
		widgets, err := GetMarketplaceWidgetsWithRequester(mockClient, 1, 50)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при получении виджетов из маркетплейса: %v", err)
		}

		if len(widgets) != 0 {
			t.Errorf("Ожидалось 0 виджетов, получено %d", len(widgets))
		}
	})

	// Проверяем сценарий с ошибкой сервера
	t.Run("Error", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("GET", "/api/v4/marketplace/widgets", http.StatusInternalServerError, `{"error": "Internal Server Error"}`, nil)

		// Вызываем тестируемый метод
		_, err := GetMarketplaceWidgetsWithRequester(mockClient, 1, 50)

		// Проверяем, что есть ошибка
		if err == nil {
			t.Fatalf("Ожидалась ошибка, но её не получили")
		}
	})
}

// TestSetWidgetStatus проверяет активацию/деактивацию виджета
func TestSetWidgetStatus(t *testing.T) {
	// ID виджета для теста
	widgetID := 123

	// Статус для установки
	status := WidgetStatusInactive

	// Подготавливаем ответ для успешного сценария
	successResponse := fmt.Sprintf(`{
		"id": %d,
		"name": "Intercom",
		"code": "intercom",
		"type": "intercom",
		"status": "%s",
		"created_by": 789,
		"updated_by": 789,
		"created_at": 1609459200,
		"updated_at": 1609459200,
		"account_id": 12345,
		"is_configured": true
	}`, widgetID, status)

	// Проверяем успешный сценарий
	t.Run("Success", func(t *testing.T) {
		// Создаем тестовый сервер для проверки тела запроса
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Проверяем метод запроса
			if r.Method != "PATCH" {
				t.Errorf("Ожидался метод PATCH, получен %s", r.Method)
			}

			// Проверяем путь запроса
			expectedPath := fmt.Sprintf("/api/v4/widgets/%d", widgetID)
			if r.URL.Path != expectedPath {
				t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
			}

			// Проверяем тело запроса
			var requestBody struct {
				Status string `json:"status"`
			}
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&requestBody); err != nil {
				t.Errorf("Ошибка при декодировании тела запроса: %v", err)
			}

			if requestBody.Status != string(status) {
				t.Errorf("Ожидался статус '%s', получен '%s'", status, requestBody.Status)
			}

			// Отправляем ответ
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(successResponse))
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемый метод
		widget, err := SetWidgetStatus(apiClient, widgetID, status)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при изменении статуса виджета: %v", err)
		}

		if widget.ID != widgetID {
			t.Errorf("Ожидался ID %d, получен %d", widgetID, widget.ID)
		}

		if widget.Status != status {
			t.Errorf("Ожидался статус %s, получен %s", status, widget.Status)
		}
	})

	// Проверяем сценарий с ошибкой
	t.Run("Error", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("PATCH", fmt.Sprintf("/api/v4/widgets/%d", widgetID), http.StatusBadRequest, `{"error": "Invalid status"}`, nil)

		// Вызываем тестируемый метод
		_, err := SetWidgetStatusWithRequester(mockClient, widgetID, status)

		// Проверяем, что есть ошибка
		if err == nil {
			t.Fatalf("Ожидалась ошибка, но её не получили")
		}
	})
}

// TestBulkInstallWidgets проверяет массовую установку виджетов
func TestBulkInstallWidgets(t *testing.T) {
	// Коды виджетов для теста
	codes := []string{"intercom", "callback"}

	// Подготавливаем ответ для успешного сценария
	successResponse := `{
		"_embedded": {
			"widgets": [
				{
					"id": 123,
					"name": "Intercom",
					"code": "intercom",
					"type": "intercom",
					"status": "installed",
					"created_by": 789,
					"updated_by": 789,
					"created_at": 1609459200,
					"updated_at": 1609459200,
					"account_id": 12345,
					"is_configured": false
				},
				{
					"id": 456,
					"name": "Callback",
					"code": "callback",
					"type": "callback",
					"status": "installed",
					"created_by": 789,
					"updated_by": 789,
					"created_at": 1609459200,
					"updated_at": 1609459200,
					"account_id": 12345,
					"is_configured": false
				}
			]
		}
	}`

	// Проверяем успешный сценарий
	t.Run("Success", func(t *testing.T) {
		// Создаем тестовый сервер для проверки тела запроса
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Проверяем метод запроса
			if r.Method != "POST" {
				t.Errorf("Ожидался метод POST, получен %s", r.Method)
			}

			// Проверяем путь запроса
			expectedPath := "/api/v4/widgets"
			if r.URL.Path != expectedPath {
				t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
			}

			// Проверяем тело запроса
			var requestBody struct {
				Codes []string `json:"codes"`
			}
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&requestBody); err != nil {
				t.Errorf("Ошибка при декодировании тела запроса: %v", err)
			}

			if len(requestBody.Codes) != len(codes) {
				t.Errorf("Ожидалось %d кодов, получено %d", len(codes), len(requestBody.Codes))
			}

			// Отправляем ответ
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(successResponse))
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемый метод
		widgets, err := BulkInstallWidgets(apiClient, codes)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при массовой установке виджетов: %v", err)
		}

		if len(widgets) != 2 {
			t.Errorf("Ожидалось 2 виджета, получено %d", len(widgets))
		}

		if widgets[0].Code != "intercom" {
			t.Errorf("Ожидался код intercom, получен %s", widgets[0].Code)
		}

		if widgets[1].Code != "callback" {
			t.Errorf("Ожидался код callback, получен %s", widgets[1].Code)
		}
	})

	// Проверяем сценарий с ошибкой
	t.Run("Error", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("POST", "/api/v4/widgets", http.StatusBadRequest, `{"error": "Invalid widget codes"}`, nil)

		// Вызываем тестируемый метод
		_, err := BulkInstallWidgetsWithRequester(mockClient, codes)

		// Проверяем, что есть ошибка
		if err == nil {
			t.Fatalf("Ожидалась ошибка, но её не получили")
		}
	})
}

// TestBulkDeleteWidgets проверяет массовое удаление виджетов
func TestBulkDeleteWidgets(t *testing.T) {
	// ID виджетов для теста
	widgetIDs := []int{123, 456}

	// Проверяем успешный сценарий
	t.Run("Success", func(t *testing.T) {
		// Создаем тестовый сервер для проверки тела запроса
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Проверяем метод запроса
			if r.Method != "DELETE" {
				t.Errorf("Ожидался метод DELETE, получен %s", r.Method)
			}

			// Проверяем путь запроса
			expectedPath := "/api/v4/widgets"
			if r.URL.Path != expectedPath {
				t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
			}

			// Проверяем тело запроса
			var requestBody struct {
				WidgetIDs []int `json:"widget_ids"`
			}
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&requestBody); err != nil {
				t.Errorf("Ошибка при декодировании тела запроса: %v", err)
			}

			if len(requestBody.WidgetIDs) != len(widgetIDs) {
				t.Errorf("Ожидалось %d ID, получено %d", len(widgetIDs), len(requestBody.WidgetIDs))
			}

			// Отправляем ответ
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемый метод
		err := BulkDeleteWidgets(apiClient, widgetIDs)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при массовом удалении виджетов: %v", err)
		}
	})

	// Проверяем сценарий с ошибкой
	t.Run("Error", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("DELETE", "/api/v4/widgets", http.StatusForbidden, `{"error": "Insufficient permissions"}`, nil)

		// Вызываем тестируемый метод
		err := BulkDeleteWidgetsWithRequester(mockClient, widgetIDs)

		// Проверяем, что есть ошибка
		if err == nil {
			t.Fatalf("Ожидалась ошибка, но её не получили")
		}
	})
}
