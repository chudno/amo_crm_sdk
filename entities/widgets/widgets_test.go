package widgets

import (
	"context"
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
func (c *AdvancedMockClient) DoRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	resp, found := c.Responses[MockRequest{Method: req.Method, Path: req.URL.Path}]

	if !found {
		resp = c.DefaultResponse
	}

	response := &http.Response{
		StatusCode: resp.StatusCode,
		Body:       io.NopCloser(strings.NewReader(resp.Body)),
		Header:     make(http.Header),
		Request:    req,
	}

	for k, v := range resp.Headers {
		response.Header.Set(k, v)
	}

	return response, nil
}

// TestGetWidgets проверяет получение списка виджетов
func TestGetWidgets(t *testing.T) {
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

	emptyResponse := `{
		"page": 1,
		"per_page": 50,
		"_embedded": {
			"widgets": []
		}
	}`

	t.Run("Success", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("GET", "/api/v4/widgets", http.StatusOK, successResponse, nil)

		types := []WidgetType{WidgetTypeIntercom, WidgetTypeCallback}

		widgets, err := ListWithRequester(context.Background(), mockClient, 1, 50, WithWidgetTypes(types))

		if err != nil {
			t.Fatalf("Ошибка при получении виджетов: %v", err)
		}

		if len(widgets) != 2 {
			t.Fatalf("Ожидалось получение 2 виджетов, получено %d", len(widgets))
		}

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

	t.Run("EmptyList", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("GET", "/api/v4/widgets", http.StatusOK, emptyResponse, nil)

		widgets, err := ListWithRequester(context.Background(), mockClient, 1, 50)

		if err != nil {
			t.Fatalf("Ошибка при получении виджетов: %v", err)
		}

		if len(widgets) != 0 {
			t.Fatalf("Ожидался пустой массив виджетов, получено %d", len(widgets))
		}
	})

	t.Run("ServerError", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("GET", "/api/v4/widgets", http.StatusInternalServerError, `{"error": "Internal Server Error"}`, nil)

		_, err := ListWithRequester(context.Background(), mockClient, 1, 50)

		if err == nil {
			t.Fatalf("Ожидалась ошибка, но её не получили")
		}
	})
}

// TestGetWidget проверяет получение информации о конкретном виджете
func TestGetWidget(t *testing.T) {
	widgetID := 123

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

	t.Run("Success", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("GET", fmt.Sprintf("/api/v4/widgets/%d", widgetID), http.StatusOK, successResponse, nil)

		widget, err := GetWithRequester(context.Background(), mockClient, widgetID)

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

	t.Run("NotFound", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("GET", fmt.Sprintf("/api/v4/widgets/%d", widgetID), http.StatusNotFound, `{"error": "Widget not found"}`, nil)

		_, err := GetWithRequester(context.Background(), mockClient, widgetID)

		if err == nil {
			t.Fatalf("Ожидалась ошибка, но её не получили")
		}
	})
}

// TestInstallWidget проверяет установку виджета из маркетплейса
func TestInstallWidget(t *testing.T) {
	widgetCode := "intercom"

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

	t.Run("Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Errorf("Ожидался метод POST, получен %s", r.Method)
			}

			if r.URL.Path != "/api/v4/widgets" {
				t.Errorf("Ожидался путь /api/v4/widgets, получен %s", r.URL.Path)
			}

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

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(successResponse))
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		widget, err := Install(context.Background(), apiClient, widgetCode)

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

	t.Run("Error", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("POST", "/api/v4/widgets", http.StatusBadRequest, `{"error": "Invalid widget code"}`, nil)

		_, err := InstallWithRequester(context.Background(), mockClient, "invalid_code")

		if err == nil {
			t.Fatalf("Ожидалась ошибка, но её не получили")
		}
	})
}

// TestUpdateWidgetSettings проверяет обновление настроек виджета
func TestUpdateWidgetSettings(t *testing.T) {
	widgetID := 123

	settings := map[string]any{
		"api_key": "new_key",
		"active":  true,
	}

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

	t.Run("Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "PATCH" {
				t.Errorf("Ожидался метод PATCH, получен %s", r.Method)
			}

			expectedPath := fmt.Sprintf("/api/v4/widgets/%d", widgetID)
			if r.URL.Path != expectedPath {
				t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
			}

			var requestBody struct {
				Settings map[string]any `json:"settings"`
			}
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&requestBody); err != nil {
				t.Errorf("Ошибка при декодировании тела запроса: %v", err)
			}

			if requestBody.Settings["api_key"] != settings["api_key"] {
				t.Errorf("Ожидался api_key '%s', получен '%s'", settings["api_key"], requestBody.Settings["api_key"])
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(successResponse))
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		widget, err := UpdateSettings(context.Background(), apiClient, widgetID, settings)

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

	t.Run("Error", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("PATCH", fmt.Sprintf("/api/v4/widgets/%d", widgetID), http.StatusBadRequest, `{"error": "Invalid settings"}`, nil)

		_, err := UpdateSettingsWithRequester(context.Background(), mockClient, widgetID, settings)

		if err == nil {
			t.Fatalf("Ожидалась ошибка, но её не получили")
		}
	})
}

// TestDeleteWidget проверяет удаление виджета
func TestDeleteWidget(t *testing.T) {
	widgetID := 123

	t.Run("Success", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("DELETE", fmt.Sprintf("/api/v4/widgets/%d", widgetID), http.StatusNoContent, "", nil)

		err := DeleteWithRequester(context.Background(), mockClient, widgetID)

		if err != nil {
			t.Fatalf("Ошибка при удалении виджета: %v", err)
		}
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("DELETE", fmt.Sprintf("/api/v4/widgets/%d", widgetID), http.StatusForbidden, `{"error": "Insufficient permissions"}`, nil)

		err := DeleteWithRequester(context.Background(), mockClient, widgetID)

		if err == nil {
			t.Fatalf("Ожидалась ошибка, но её не получили")
		}
	})
}
