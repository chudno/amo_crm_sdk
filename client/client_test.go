package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		apiKey   string
		expected *Client
	}{
		{
			name:    "Создание нового клиента",
			baseURL: "https://test.amocrm.ru",
			apiKey:  "test_api_key",
			expected: &Client{
				baseURL: "https://test.amocrm.ru",
				apiKey:  "test_api_key",
				httpClient: &http.Client{
					Timeout: 30 * time.Second,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewClient(tt.baseURL, tt.apiKey)

			if got.baseURL != tt.expected.baseURL {
				t.Errorf("NewClient().baseURL = %v, want %v", got.baseURL, tt.expected.baseURL)
			}

			if got.apiKey != tt.expected.apiKey {
				t.Errorf("NewClient().apiKey = %v, want %v", got.apiKey, tt.expected.apiKey)
			}

			if got.httpClient == nil {
				t.Error("NewClient().httpClient is nil")
			}

			if got.httpClient.Timeout != tt.expected.httpClient.Timeout {
				t.Errorf("NewClient().httpClient.Timeout = %v, want %v",
					got.httpClient.Timeout, tt.expected.httpClient.Timeout)
			}
		})
	}
}

func TestGetBaseURL(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		expected string
	}{
		{
			name:     "Базовый URL возвращается корректно",
			baseURL:  "https://test.amocrm.ru",
			expected: "https://test.amocrm.ru",
		},
		{
			name:     "Пустой URL",
			baseURL:  "",
			expected: "",
		},
		{
			name:     "URL с дополнительным путем",
			baseURL:  "https://test.amocrm.ru/api/v4",
			expected: "https://test.amocrm.ru/api/v4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.baseURL, "test_api_key")

			got := client.GetBaseURL()

			if got != tt.expected {
				t.Errorf("GetBaseURL() = %v, хотим %v", got, tt.expected)
			}
		})
	}
}

func TestDoRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		expectedAuth := "Bearer test_api_key"
		if auth != expectedAuth {
			t.Errorf("Ожидался заголовок Authorization %s, получен %s", expectedAuth, auth)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":true}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test_api_key")

	req, err := http.NewRequest("GET", server.URL+"/test", nil)
	if err != nil {
		t.Fatalf("Ошибка создания запроса: %v", err)
	}

	resp, err := client.DoRequest(context.Background(), req)

	if err != nil {
		t.Fatalf("Неожиданная ошибка: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус код %d, получен %d", http.StatusOK, resp.StatusCode)
	}
}
