package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAuthURL(t *testing.T) {
	tests := []struct {
		name        string
		baseURL     string
		clientID    string
		redirectURI string
		state       string
		mode        string
		expected    string
	}{
		{
			name:        "Базовый URL авторизации",
			baseURL:     "https://test.amocrm.ru",
			clientID:    "test_client_id",
			redirectURI: "https://test-redirect.com",
			state:       "random_state",
			mode:        "popup",
			expected:    "https://test.amocrm.ru/oauth2/access_token?client_id=test_client_id&mode=popup&redirect_uri=https%3A%2F%2Ftest-redirect.com&response_type=code&state=random_state",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetAuthURL(tt.baseURL, tt.clientID, tt.redirectURI, tt.state, tt.mode)
			if got != tt.expected {
				t.Errorf("GetAuthURL() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetAccessToken(t *testing.T) {
	tests := []struct {
		name         string
		responseCode int
		responseBody string
		expectError  bool
	}{
		{
			name:         "Успешное получение токена",
			responseCode: http.StatusOK,
			responseBody: `{"token_type":"Bearer","expires_in":86400,"access_token":"test_access_token","refresh_token":"test_refresh_token"}`},
		{
			name:         "Ошибка сервера",
			responseCode: http.StatusInternalServerError,
			responseBody: `{"error":"server_error"}`,
			expectError:  true,
		},
		{
			name:         "Неверный формат ответа",
			responseCode: http.StatusOK,
			responseBody: `{invalid_json`,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Проверяем метод запроса
				if r.Method != "POST" {
					t.Errorf("Ожидался метод POST, получен %s", r.Method)
				}

				// Проверяем заголовок Content-Type
				contentType := r.Header.Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Ожидался Content-Type application/json, получен %s", contentType)
				}

				// Проверяем тело запроса
				var authReq AuthRequest
				if err := json.NewDecoder(r.Body).Decode(&authReq); err != nil {
					t.Errorf("Ошибка декодирования тела запроса: %v", err)
				}

				// Проверяем параметры запроса
				if authReq.GrantType != "authorization_code" {
					t.Errorf("Ожидался grant_type authorization_code, получен %s", authReq.GrantType)
				}

				// Устанавливаем код ответа и тело
				w.WriteHeader(tt.responseCode)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			// Вызываем тестируемую функцию
			response, err := GetAccessToken(server.URL, "test_client_id", "test_client_secret", "test_code", "https://test-redirect.com")

			// Проверяем результаты
			if tt.expectError && err == nil {
				t.Error("Ожидалась ошибка, но ее не было")
			}

			if !tt.expectError {
				if err != nil {
					t.Errorf("Неожиданная ошибка: %v", err)
				}

				if response.AccessToken != "test_access_token" {
					t.Errorf("Ожидался access_token test_access_token, получен %s", response.AccessToken)
				}

				if response.RefreshToken != "test_refresh_token" {
					t.Errorf("Ожидался refresh_token test_refresh_token, получен %s", response.RefreshToken)
				}

				if response.ExpiresIn != 86400 {
					t.Errorf("Ожидался expires_in 86400, получен %d", response.ExpiresIn)
				}
			}
		})
	}
}

func TestRefreshAccessToken(t *testing.T) {
	tests := []struct {
		name         string
		responseCode int
		responseBody string
		expectError  bool
	}{
		{
			name:         "Успешное обновление токена",
			responseCode: http.StatusOK,
			responseBody: `{"token_type":"Bearer","expires_in":86400,"access_token":"new_access_token","refresh_token":"new_refresh_token"}`},
		{
			name:         "Ошибка сервера",
			responseCode: http.StatusInternalServerError,
			responseBody: `{"error":"server_error"}`,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Проверяем метод запроса
				if r.Method != "POST" {
					t.Errorf("Ожидался метод POST, получен %s", r.Method)
				}

				// Проверяем заголовок Content-Type
				contentType := r.Header.Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Ожидался Content-Type application/json, получен %s", contentType)
				}

				// Проверяем тело запроса
				var authReq AuthRequest
				if err := json.NewDecoder(r.Body).Decode(&authReq); err != nil {
					t.Errorf("Ошибка декодирования тела запроса: %v", err)
				}

				// Проверяем параметры запроса
				if authReq.GrantType != "refresh_token" {
					t.Errorf("Ожидался grant_type refresh_token, получен %s", authReq.GrantType)
				}

				// Устанавливаем код ответа и тело
				w.WriteHeader(tt.responseCode)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			// Вызываем тестируемую функцию
			response, err := RefreshAccessToken(server.URL, "test_client_id", "test_client_secret", "test_refresh_token")

			// Проверяем результаты
			if tt.expectError && err == nil {
				t.Error("Ожидалась ошибка, но ее не было")
			}

			if !tt.expectError {
				if err != nil {
					t.Errorf("Неожиданная ошибка: %v", err)
				}

				if response.AccessToken != "new_access_token" {
					t.Errorf("Ожидался access_token new_access_token, получен %s", response.AccessToken)
				}

				if response.RefreshToken != "new_refresh_token" {
					t.Errorf("Ожидался refresh_token new_refresh_token, получен %s", response.RefreshToken)
				}
			}
		})
	}
}

func TestGetLongLivedToken(t *testing.T) {
	tests := []struct {
		name         string
		responseCode int
		responseBody string
		expectError  bool
	}{
		{
			name:         "Успешное получение долгоживущего токена",
			responseCode: http.StatusOK,
			responseBody: `{"token_type":"Bearer","expires_in":31536000,"access_token":"long_lived_token"}`,
		},
		{
			name:         "Ошибка сервера",
			responseCode: http.StatusInternalServerError,
			responseBody: `{"error":"server_error"}`,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Проверяем метод запроса
				if r.Method != "POST" {
					t.Errorf("Ожидался метод POST, получен %s", r.Method)
				}

				// Проверяем заголовок Content-Type
				contentType := r.Header.Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Ожидался Content-Type application/json, получен %s", contentType)
				}

				// Проверяем тело запроса
				var authReq AuthRequest
				if err := json.NewDecoder(r.Body).Decode(&authReq); err != nil {
					t.Errorf("Ошибка декодирования тела запроса: %v", err)
				}

				// Проверяем параметры запроса
				if authReq.GrantType != "client_credentials" {
					t.Errorf("Ожидался grant_type client_credentials, получен %s", authReq.GrantType)
				}

				// Устанавливаем код ответа и тело
				w.WriteHeader(tt.responseCode)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			// Вызываем тестируемую функцию
			response, err := GetLongLivedToken(server.URL, "test_client_id", "test_client_secret")

			// Проверяем результаты
			if tt.expectError && err == nil {
				t.Error("Ожидалась ошибка, но ее не было")
			}

			if !tt.expectError {
				if err != nil {
					t.Errorf("Неожиданная ошибка: %v", err)
				}

				if response.AccessToken != "long_lived_token" {
					t.Errorf("Ожидался access_token long_lived_token, получен %s", response.AccessToken)
				}

				// У долгоживущего токена должно быть большое время жизни
				if response.ExpiresIn != 31536000 { // Примерно год в секундах
					t.Errorf("Ожидался expires_in 31536000, получен %d", response.ExpiresIn)
				}

				// Refresh token должен отсутствовать или быть пустым
				if response.RefreshToken != "" {
					t.Errorf("Для долгоживущего токена не должно быть refresh_token, получен %s", response.RefreshToken)
				}
			}
		})
	}
}
