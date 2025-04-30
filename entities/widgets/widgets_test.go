package widgets

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

// Используем интерфейс Requester из основного пакета

// MockResponse описывает мок-ответ для тестирования
type MockResponse struct {
	StatusCode int
	Body       string
	Headers    map[string]string
}

// MockRequest описывает мок-запрос для сопоставления
type MockRequest struct {
	Method string
	Path   string
}

// AdvancedMockClient улучшенный мок-клиент для перехвата всех видов запросов
type AdvancedMockClient struct {
	// Отображение ключа MockRequest на ответ MockResponse
	Responses map[MockRequest]MockResponse
	// Ответ по умолчанию, если запрос не найден
	DefaultResponse MockResponse
}

// NewAdvancedMockClient создает новый мок-клиент с настройками по умолчанию
func NewAdvancedMockClient() *AdvancedMockClient {
	return &AdvancedMockClient{
		Responses: make(map[MockRequest]MockResponse),
		DefaultResponse: MockResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error": "Unexpected request"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		},
	}
}

// AddResponse добавляет ответ для конкретного запроса
func (c *AdvancedMockClient) AddResponse(method, path string, statusCode int, body string, headers map[string]string) {
	if headers == nil {
		headers = map[string]string{"Content-Type": "application/json"}
	}
	c.Responses[MockRequest{Method: method, Path: path}] = MockResponse{
		StatusCode: statusCode,
		Body:       body,
		Headers:    headers,
	}
}

// DoRequest реализует интерфейс Requester
func (c *AdvancedMockClient) DoRequest(req *http.Request) (*http.Response, error) {
	// Ищем подходящий ответ для метода и пути
	resp, found := c.Responses[MockRequest{Method: req.Method, Path: req.URL.Path}]
	
	// Если не найден, возвращаем ответ по умолчанию
	if !found {
		resp = c.DefaultResponse
	}
	
	// Формируем HTTP-ответ
	response := &http.Response{
		StatusCode: resp.StatusCode,
		Body:       io.NopCloser(strings.NewReader(resp.Body)),
		Header:     make(http.Header),
		Request:    req,
	}
	
	// Добавляем заголовки
	for k, v := range resp.Headers {
		response.Header.Set(k, v)
	}
	
	return response, nil
}

// TestGetWidgets проверяет получение списка виджетов
func TestGetWidgets(t *testing.T) {
	// Подготавливаем ответ для успешного сценария
	successResponse := `{
		"page": 1,
		"per_page": 50,
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
					"is_configured": true
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
					"is_configured": true
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
		mockClient.AddResponse("GET", "/api/v4/widgets", http.StatusOK, successResponse, nil)

		// Создаем фильтр по типам виджетов
		types := []WidgetType{WidgetTypeIntercom, WidgetTypeCallback}

		// Вызываем тестируемый метод
		widgets, err := GetWidgetsWithRequester(mockClient, 1, 50, WithWidgetTypes(types))

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при получении виджетов: %v", err)
		}

		if len(widgets) != 2 {
			t.Fatalf("Ожидалось получение 2 виджетов, получено %d", len(widgets))
		}

		// Проверяем содержимое первого виджета
		if widgets[0].ID != 123 {
			t.Errorf("Ожидался ID 123, получен %d", widgets[0].ID)
		}

		if widgets[0].Name != "Intercom" {
			t.Errorf("Ожидалось имя 'Intercom', получено '%s'", widgets[0].Name)
		}

		if widgets[0].Type != WidgetTypeIntercom {
			t.Errorf("Ожидался тип 'intercom', получен '%s'", widgets[0].Type)
		}
	})

	// Проверяем сценарий с пустым списком
	t.Run("EmptyList", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("GET", "/api/v4/widgets", http.StatusOK, emptyResponse, nil)

		// Вызываем тестируемый метод
		widgets, err := GetWidgetsWithRequester(mockClient, 1, 50)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при получении виджетов: %v", err)
		}

		// Проверяем, что массив пуст
		if len(widgets) != 0 {
			t.Fatalf("Ожидался пустой массив виджетов, получено %d", len(widgets))
		}
	})

	// Проверяем сценарий с ошибкой сервера
	t.Run("ServerError", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("GET", "/api/v4/widgets", http.StatusInternalServerError, `{"error": "Internal Server Error"}`, nil)

		// Вызываем тестируемый метод
		_, err := GetWidgetsWithRequester(mockClient, 1, 50)

		// Проверяем, что есть ошибка
		if err == nil {
			t.Fatalf("Ожидалась ошибка, но её не получили")
		}
	})
}

// TestGetWidget проверяет получение информации о конкретном виджете
func TestGetWidget(t *testing.T) {
	// ID виджета для теста
	widgetID := 123

	// Подготавливаем ответ для успешного сценария
	successResponse := `{
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
		"is_configured": true,
		"settings": {
			"api_key": "test_key",
			"active": true
		}
	}`

	// Проверяем успешный сценарий
	t.Run("Success", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("GET", fmt.Sprintf("/api/v4/widgets/%d", widgetID), http.StatusOK, successResponse, nil)

		// Вызываем тестируемый метод
		widget, err := GetWidgetWithRequester(mockClient, widgetID)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при получении виджета: %v", err)
		}

		if widget.ID != widgetID {
			t.Errorf("Ожидался ID %d, получен %d", widgetID, widget.ID)
		}

		if widget.Name != "Intercom" {
			t.Errorf("Ожидалось имя 'Intercom', получено '%s'", widget.Name)
		}

		if widget.Type != WidgetTypeIntercom {
			t.Errorf("Ожидался тип 'intercom', получен '%s'", widget.Type)
		}

		if !widget.IsConfigured {
			t.Errorf("Ожидалось, что виджет настроен (is_configured=true)")
		}
	})

	// Проверяем сценарий, когда виджет не найден
	t.Run("NotFound", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("GET", fmt.Sprintf("/api/v4/widgets/%d", widgetID), http.StatusNotFound, `{"error": "Widget not found"}`, nil)

		// Вызываем тестируемый метод
		_, err := GetWidgetWithRequester(mockClient, widgetID)

		// Проверяем, что есть ошибка
		if err == nil {
			t.Fatalf("Ожидалась ошибка, но её не получили")
		}
	})
}

// TestInstallWidget проверяет установку виджета из маркетплейса
func TestInstallWidget(t *testing.T) {
	// Код виджета для установки
	widgetCode := "intercom"

	// Подготавливаем ответ для успешного сценария
	successResponse := `{
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
			if r.URL.Path != "/api/v4/widgets" {
				t.Errorf("Ожидался путь /api/v4/widgets, получен %s", r.URL.Path)
			}

			// Проверяем тело запроса
			var requestBody struct {
				Code string `json:"code"`
			}
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&requestBody); err != nil {
				t.Errorf("Ошибка при декодировании тела запроса: %v", err)
			}

			if requestBody.Code != widgetCode {
				t.Errorf("Ожидался код виджета '%s', получен '%s'", widgetCode, requestBody.Code)
			}

			// Отправляем ответ
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(successResponse))
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемый метод
		widget, err := InstallWidget(apiClient, widgetCode)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при установке виджета: %v", err)
		}

		if widget.ID != 123 {
			t.Errorf("Ожидался ID 123, получен %d", widget.ID)
		}

		if widget.Code != widgetCode {
			t.Errorf("Ожидался код '%s', получен '%s'", widgetCode, widget.Code)
		}

		if widget.Status != WidgetStatusInstalled {
			t.Errorf("Ожидался статус 'installed', получен '%s'", widget.Status)
		}
	})

	// Проверяем сценарий с ошибкой
	t.Run("Error", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("POST", "/api/v4/widgets", http.StatusBadRequest, `{"error": "Invalid widget code"}`, nil)

		// Вызываем тестируемый метод
		_, err := InstallWidgetWithRequester(mockClient, "invalid_code")

		// Проверяем, что есть ошибка
		if err == nil {
			t.Fatalf("Ожидалась ошибка, но её не получили")
		}
	})
}

// TestUpdateWidgetSettings проверяет обновление настроек виджета
func TestUpdateWidgetSettings(t *testing.T) {
	// ID виджета для теста
	widgetID := 123

	// Настройки для обновления
	settings := map[string]interface{}{
		"api_key": "new_key",
		"active":  true,
	}

	// Подготавливаем ответ для успешного сценария
	successResponse := `{
		"id": 123,
		"name": "Intercom",
		"code": "intercom",
		"type": "intercom",
		"status": "installed",
		"created_by": 789,
		"updated_by": 789,
		"created_at": 1609459200,
		"updated_at": 1609545600,
		"account_id": 12345,
		"is_configured": true,
		"settings": {
			"api_key": "new_key",
			"active": true
		}
	}`

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
				Settings map[string]interface{} `json:"settings"`
			}
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&requestBody); err != nil {
				t.Errorf("Ошибка при декодировании тела запроса: %v", err)
			}

			if requestBody.Settings["api_key"] != settings["api_key"] {
				t.Errorf("Ожидался api_key '%s', получен '%s'", settings["api_key"], requestBody.Settings["api_key"])
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
		widget, err := UpdateWidgetSettings(apiClient, widgetID, settings)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при обновлении настроек виджета: %v", err)
		}

		if widget.ID != widgetID {
			t.Errorf("Ожидался ID %d, получен %d", widgetID, widget.ID)
		}

		if !widget.IsConfigured {
			t.Errorf("Ожидалось, что виджет настроен (is_configured=true)")
		}
	})

	// Проверяем сценарий с ошибкой
	t.Run("Error", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("PATCH", fmt.Sprintf("/api/v4/widgets/%d", widgetID), http.StatusBadRequest, `{"error": "Invalid settings"}`, nil)

		// Вызываем тестируемый метод
		_, err := UpdateWidgetSettingsWithRequester(mockClient, widgetID, settings)

		// Проверяем, что есть ошибка
		if err == nil {
			t.Fatalf("Ожидалась ошибка, но её не получили")
		}
	})
}

// TestDeleteWidget проверяет удаление виджета
func TestDeleteWidget(t *testing.T) {
	// ID виджета для теста
	widgetID := 123

	// Проверяем успешный сценарий
	t.Run("Success", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("DELETE", fmt.Sprintf("/api/v4/widgets/%d", widgetID), http.StatusNoContent, "", nil)

		// Вызываем тестируемый метод
		err := DeleteWidgetWithRequester(mockClient, widgetID)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при удалении виджета: %v", err)
		}
	})

	// Проверяем сценарий с ошибкой
	t.Run("Error", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("DELETE", fmt.Sprintf("/api/v4/widgets/%d", widgetID), http.StatusForbidden, `{"error": "Insufficient permissions"}`, nil)

		// Вызываем тестируемый метод
		err := DeleteWidgetWithRequester(mockClient, widgetID)

		// Проверяем, что есть ошибка
		if err == nil {
			t.Fatalf("Ожидалась ошибка, но её не получили")
		}
	})
}
