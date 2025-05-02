package users

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

// TestGetUserErrors проверяет обработку ошибок при получении пользователя
func TestGetUserErrors(t *testing.T) {
	t.Run("Пользователь не найден", func(t *testing.T) {
		// Создаем тестовый сервер, который вернет ошибку 404
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			if _, err := w.Write([]byte(`{"title": "Not found", "status": 404, "detail": "User not found"}`)); err != nil {
				t.Fatalf("Ошибка при записи ответа: %v", err)
			}
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемую функцию
		_, err := GetUser(apiClient, 999) // несуществующий ID

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка, но она не возникла")
		}
	})

	t.Run("Некорректный JSON", func(t *testing.T) {
		// Создаем тестовый сервер, который вернет некорректный JSON
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(`{"id": 123, "name": "Иван Иванов", "email":`)); err != nil {
				t.Fatalf("Ошибка при записи ответа: %v", err)
			}
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемую функцию
		_, err := GetUser(apiClient, 123)

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка из-за некорректного JSON, но она не возникла")
		}
	})

	t.Run("Ошибка HTTP", func(t *testing.T) {
		// Используем несуществующий домен для вызова ошибки
		apiClient := client.NewClient("http://non-existent-domain.example", "test_api_key")

		// Вызываем тестируемую функцию
		_, err := GetUser(apiClient, 123)

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка HTTP, но она не возникла")
		}
	})

	t.Run("Ошибка сервера (500)", func(t *testing.T) {
		// Создаем тестовый сервер, который вернет ошибку 500
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte(`{"title": "Internal Server Error", "status": 500}`)); err != nil {
				t.Fatalf("Ошибка при записи ответа: %v", err)
			}
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемую функцию
		_, err := GetUser(apiClient, 123)

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка из-за статуса 500, но она не возникла")
		}
	})
}

// TestGetCurrentUserErrors проверяет обработку ошибок при получении информации о текущем пользователе
func TestGetCurrentUserErrors(t *testing.T) {
	t.Run("Некорректный JSON", func(t *testing.T) {
		// Создаем тестовый сервер, который вернет некорректный JSON
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(`{"id": 456, "name": "Петр Петров", "rights": {`)); err != nil {
				t.Fatalf("Ошибка при записи ответа: %v", err)
			}
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемую функцию
		_, err := GetCurrentUser(apiClient)

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка из-за некорректного JSON, но она не возникла")
		}
	})

	t.Run("Ошибка HTTP", func(t *testing.T) {
		// Используем несуществующий домен для вызова ошибки
		apiClient := client.NewClient("http://non-existent-domain.example", "test_api_key")

		// Вызываем тестируемую функцию
		_, err := GetCurrentUser(apiClient)

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка HTTP, но она не возникла")
		}
	})

	t.Run("Ошибка авторизации (401)", func(t *testing.T) {
		// Создаем тестовый сервер, который вернет ошибку 401
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			if _, err := w.Write([]byte(`{"title": "Unauthorized", "status": 401, "detail": "Invalid API key"}`)); err != nil {
				t.Fatalf("Ошибка при записи ответа: %v", err)
			}
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "invalid_api_key")

		// Вызываем тестируемую функцию
		_, err := GetCurrentUser(apiClient)

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка из-за неверного API-ключа, но она не возникла")
		}
	})
}

// TestListUsersErrors проверяет обработку ошибок при получении списка пользователей
func TestListUsersErrors(t *testing.T) {
	t.Run("Пустой список пользователей", func(t *testing.T) {
		// Создаем тестовый сервер, который вернет пустой список
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(`{"_embedded": {"items": []}}`)); err != nil {
				t.Fatalf("Ошибка при записи ответа: %v", err)
			}
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемую функцию
		users, err := ListUsers(apiClient, 50, 1)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Не ожидалась ошибка, но получена: %v", err)
		}

		if len(users) != 0 {
			t.Errorf("Ожидался пустой список пользователей, получено %d элементов", len(users))
		}
	})

	t.Run("Некорректный JSON", func(t *testing.T) {
		// Создаем тестовый сервер, который вернет некорректный JSON
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(`{"_embedded": {"items": [{`)); err != nil {
				t.Fatalf("Ошибка при записи ответа: %v", err)
			}
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемую функцию
		_, err := ListUsers(apiClient, 50, 1)

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка из-за некорректного JSON, но она не возникла")
		}
	})

	t.Run("Ошибка HTTP", func(t *testing.T) {
		// Используем несуществующий домен для вызова ошибки
		apiClient := client.NewClient("http://non-existent-domain.example", "test_api_key")

		// Вызываем тестируемую функцию
		_, err := ListUsers(apiClient, 50, 1)

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка HTTP, но она не возникла")
		}
	})

	t.Run("Некорректные параметры страницы", func(t *testing.T) {
		// Создаем тестовый сервер, который вернет ошибку 400
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			query := r.URL.Query()
			page := query.Get("page")
			
			if page == "0" {
				// Вернуть ошибку для некорректного номера страницы
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				if _, err := w.Write([]byte(`{"title": "Bad Request", "status": 400, "detail": "Page number must be greater than 0"}`)); err != nil {
					t.Fatalf("Ошибка при записи ответа: %v", err)
				}
				return
			}
			
			// Успешный ответ для корректных параметров
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(`{"_embedded": {"items": []}}`)); err != nil {
				t.Fatalf("Ошибка при записи ответа: %v", err)
			}
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемую функцию с некорректными параметрами
		_, err := ListUsers(apiClient, 50, 0) // page=0 - некорректный параметр

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка из-за некорректных параметров, но она не возникла")
		}
	})
}
