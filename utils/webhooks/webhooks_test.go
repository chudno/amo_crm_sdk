package webhooks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

func TestGetWebhook(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		expectedPath := "/api/v4/webhooks/123"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 123,
			"destination": "https://example.com/webhook",
			"settings": {
				"events": ["leads", "contacts"],
				"actions": ["add", "update"]
			},
			"created_at": 1609459200,
			"updated_at": 1609545600
		}`))
	}))
	defer server.Close()

	apiClient := client.NewClient(server.URL, "test_api_key")

	webhook, err := Get(context.Background(), apiClient, 123)

	if err != nil {
		t.Fatalf("Ошибка при получении вебхука: %v", err)
	}

	if webhook.ID != 123 {
		t.Errorf("Ожидался ID вебхука 123, получен %d", webhook.ID)
	}

	if webhook.Destination != "https://example.com/webhook" {
		t.Errorf("Ожидался URL вебхука 'https://example.com/webhook', получен '%s'", webhook.Destination)
	}

	if len(webhook.Settings.Entities) != 2 {
		t.Errorf("Ожидалось 2 типа сущностей, получено %d", len(webhook.Settings.Entities))
	}

	if len(webhook.Settings.Actions) != 2 {
		t.Errorf("Ожидалось 2 типа действий, получено %d", len(webhook.Settings.Actions))
	}
}

func TestCreateWebhook(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		expectedPath := "/api/v4/webhooks"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		var webhook Webhook
		if err := json.NewDecoder(r.Body).Decode(&webhook); err != nil {
			t.Errorf("Ошибка декодирования тела запроса: %v", err)
		}

		if webhook.Destination != "https://example.com/new-webhook" {
			t.Errorf("Ожидался URL вебхука 'https://example.com/new-webhook', получен '%s'", webhook.Destination)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 456,
			"destination": "https://example.com/new-webhook",
			"settings": {
				"events": ["leads"],
				"actions": ["add"]
			},
			"created_at": 1609459200,
			"updated_at": 1609545600
		}`))
	}))
	defer server.Close()

	apiClient := client.NewClient(server.URL, "test_api_key")

	webhookToCreate := &Webhook{
		Destination: "https://example.com/new-webhook",
		Settings: &Settings{
			Entities: []string{"leads"},
			Actions:  []string{"add"},
		},
	}

	createdWebhook, err := Create(context.Background(), apiClient, webhookToCreate)

	if err != nil {
		t.Fatalf("Ошибка при создании вебхука: %v", err)
	}

	if createdWebhook.ID != 456 {
		t.Errorf("Ожидался ID вебхука 456, получен %d", createdWebhook.ID)
	}

	if createdWebhook.Destination != "https://example.com/new-webhook" {
		t.Errorf("Ожидался URL вебхука 'https://example.com/new-webhook', получен '%s'", createdWebhook.Destination)
	}

	if len(createdWebhook.Settings.Entities) != 1 || createdWebhook.Settings.Entities[0] != "leads" {
		t.Errorf("Ожидался тип сущности 'leads', получено '%v'", createdWebhook.Settings.Entities)
	}

	if len(createdWebhook.Settings.Actions) != 1 || createdWebhook.Settings.Actions[0] != "add" {
		t.Errorf("Ожидался тип действия 'add', получено '%v'", createdWebhook.Settings.Actions)
	}
}

func TestDeleteWebhook(t *testing.T) {
	tests := []struct {
		name         string
		webhookID    int
		responseCode int
		responseBody string
		expectError  bool
	}{
		{
			name:         "Успешное удаление вебхука",
			webhookID:    12345,
			responseCode: http.StatusOK,
			responseBody: `{"success":true}`,
		},
		{
			name:         "Вебхук не найден",
			webhookID:    99999,
			responseCode: http.StatusNotFound,
			responseBody: `{"error":"not_found"}`,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "DELETE" {
					t.Errorf("Ожидался метод DELETE, получен %s", r.Method)
				}

				expectedPath := fmt.Sprintf("/api/v4/webhooks/%d", tt.webhookID)
				if r.URL.Path != expectedPath {
					t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(tt.responseCode)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			apiClient := client.NewClient(server.URL, "test_api_key")

			err := Delete(context.Background(), apiClient, tt.webhookID)

			if tt.expectError && err == nil {
				t.Error("Ожидалась ошибка, но ее не было")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Неожиданная ошибка: %v", err)
			}
		})
	}
}

func TestListWebhooks(t *testing.T) {
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
			name:         "Успешное получение списка вебхуков",
			page:         1,
			limit:        50,
			responseCode: http.StatusOK,
			responseBody: `{"_embedded":{"webhooks":[{"id":12345,"destination":"https://example.com/webhook1"},{"id":12346,"destination":"https://example.com/webhook2"}]}}`,
			expectedLen:  2,
		},
		{
			name:         "Пустой список вебхуков",
			page:         1,
			limit:        50,
			responseCode: http.StatusOK,
			responseBody: `{"_embedded":{"webhooks":[]}}`,
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

				expectedPath := "/api/v4/webhooks"
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

			webhooks, err := List(context.Background(), apiClient, tt.limit, tt.page)

			if tt.expectError && err == nil {
				t.Error("Ожидалась ошибка, но ее не было")
			}

			if !tt.expectError {
				if err != nil {
					t.Errorf("Неожиданная ошибка: %v", err)
				}

				if len(webhooks) != tt.expectedLen {
					t.Errorf("Ожидалось %d вебхуков, получено %d", tt.expectedLen, len(webhooks))
				}
			}
		})
	}
}

// Вспомогательная функция для проверки метода и пути запроса
func validateRequestBasics(t *testing.T, r *http.Request, expectedMethod, expectedPath string) {
	if r.Method != expectedMethod {
		t.Errorf("Ожидался метод %s, получен %s", expectedMethod, r.Method)
	}

	if r.URL.Path != expectedPath {
		t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
	}
}

// Вспомогательная функция для проверки параметров вебхука
func validateWebhook(t *testing.T, webhook Webhook, expectedDestination string, expectedEntities, expectedActions []string) {
	if webhook.Destination != expectedDestination {
		t.Errorf("Ожидался URL вебхука '%s', получен '%s'", expectedDestination, webhook.Destination)
	}

	if !reflect.DeepEqual(webhook.Settings.Entities, expectedEntities) {
		t.Errorf("Ожидались типы сущностей %v, получено %v", expectedEntities, webhook.Settings.Entities)
	}

	if !reflect.DeepEqual(webhook.Settings.Actions, expectedActions) {
		t.Errorf("Ожидались типы действий %v, получено %v", expectedActions, webhook.Settings.Actions)
	}
}

func TestCreateSimpleWebhook(t *testing.T) {
	expectedDestination := "https://example.com/simple-webhook"
	expectedEntities := []string{"leads"}
	expectedActions := []string{"add"}
	expectedID := 789

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		validateRequestBasics(t, r, "POST", "/api/v4/webhooks")

		var webhook Webhook
		if err := json.NewDecoder(r.Body).Decode(&webhook); err != nil {
			t.Errorf("Ошибка декодирования тела запроса: %v", err)
		}

		validateWebhook(t, webhook, expectedDestination, expectedEntities, expectedActions)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 789,
			"destination": "https://example.com/simple-webhook",
			"settings": {
				"events": ["leads"],
				"actions": ["add"]
			},
			"created_at": 1609459200,
			"updated_at": 1609545600
		}`))
	}))
	defer server.Close()

	apiClient := client.NewClient(server.URL, "test_api_key")

	createdWebhook, err := CreateSimple(context.Background(), apiClient, expectedDestination, expectedEntities, expectedActions)

	if err != nil {
		t.Fatalf("Ошибка при создании простого вебхука: %v", err)
	}

	if createdWebhook.ID != expectedID {
		t.Errorf("Ожидался ID вебхука %d, получен %d", expectedID, createdWebhook.ID)
	}

	validateWebhook(t, *createdWebhook, expectedDestination, expectedEntities, expectedActions)
}
