package widgets

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

// TestGetMarketplaceWidgets проверяет получение списка виджетов из маркетплейса
func TestGetMarketplaceWidgets(t *testing.T) {
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

	emptyResponse := `{
		"page": 1,
		"per_page": 50,
		"_embedded": {
			"widgets": []
		}
	}`

	t.Run("Success", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("GET", "/api/v4/marketplace/widgets", http.StatusOK, successResponse, nil)

		categoryID := 1

		widgets, err := ListMarketplaceWithRequester(context.Background(), mockClient, 1, 50, WithCategory(categoryID))

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

	t.Run("Empty", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("GET", "/api/v4/marketplace/widgets", http.StatusOK, emptyResponse, nil)

		widgets, err := ListMarketplaceWithRequester(context.Background(), mockClient, 1, 50)

		if err != nil {
			t.Fatalf("Ошибка при получении виджетов из маркетплейса: %v", err)
		}

		if len(widgets) != 0 {
			t.Errorf("Ожидалось 0 виджетов, получено %d", len(widgets))
		}
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("GET", "/api/v4/marketplace/widgets", http.StatusInternalServerError, `{"error": "Internal Server Error"}`, nil)

		_, err := ListMarketplaceWithRequester(context.Background(), mockClient, 1, 50)

		if err == nil {
			t.Fatalf("Ожидалась ошибка, но её не получили")
		}
	})
}

// TestSetWidgetStatus проверяет активацию/деактивацию виджета
func TestSetWidgetStatus(t *testing.T) {
	widgetID := 123

	status := WidgetStatusInactive

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
				Status string `json:"status"`
			}
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&requestBody); err != nil {
				t.Errorf("Ошибка при декодировании тела запроса: %v", err)
			}

			if requestBody.Status != string(status) {
				t.Errorf("Ожидался статус '%s', получен '%s'", status, requestBody.Status)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(successResponse))
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		widget, err := SetStatus(context.Background(), apiClient, widgetID, status)

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

	t.Run("Error", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("PATCH", fmt.Sprintf("/api/v4/widgets/%d", widgetID), http.StatusBadRequest, `{"error": "Invalid status"}`, nil)

		_, err := SetStatusWithRequester(context.Background(), mockClient, widgetID, status)

		if err == nil {
			t.Fatalf("Ожидалась ошибка, но её не получили")
		}
	})
}

// TestBulkInstallWidgets проверяет массовую установку виджетов
func TestBulkInstallWidgets(t *testing.T) {
	codes := []string{"intercom", "callback"}

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

	t.Run("Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Errorf("Ожидался метод POST, получен %s", r.Method)
			}

			expectedPath := "/api/v4/widgets"
			if r.URL.Path != expectedPath {
				t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
			}

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

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(successResponse))
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		widgets, err := BulkInstall(context.Background(), apiClient, codes)

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

	t.Run("Error", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("POST", "/api/v4/widgets", http.StatusBadRequest, `{"error": "Invalid widget codes"}`, nil)

		_, err := BulkInstallWithRequester(context.Background(), mockClient, codes)

		if err == nil {
			t.Fatalf("Ожидалась ошибка, но её не получили")
		}
	})
}

// TestBulkDeleteWidgets проверяет массовое удаление виджетов
func TestBulkDeleteWidgets(t *testing.T) {
	widgetIDs := []int{123, 456}

	t.Run("Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "DELETE" {
				t.Errorf("Ожидался метод DELETE, получен %s", r.Method)
			}

			expectedPath := "/api/v4/widgets"
			if r.URL.Path != expectedPath {
				t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
			}

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

			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		err := BulkDelete(context.Background(), apiClient, widgetIDs)

		if err != nil {
			t.Fatalf("Ошибка при массовом удалении виджетов: %v", err)
		}
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("DELETE", "/api/v4/widgets", http.StatusForbidden, `{"error": "Insufficient permissions"}`, nil)

		err := BulkDeleteWithRequester(context.Background(), mockClient, widgetIDs)

		if err == nil {
			t.Fatalf("Ожидалась ошибка, но её не получили")
		}
	})
}
