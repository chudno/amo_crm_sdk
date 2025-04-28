package users

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

func TestGetUser(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/users/123"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 123,
			"name": "Иван Иванов",
			"email": "ivan@example.com",
			"lang": "ru",
			"is_active": true,
			"rights": {
				"leads": true,
				"contacts": true,
				"companies": true,
				"tasks": true,
				"mailbox": false,
				"catalog": false,
				"is_admin": false,
				"is_manager": true
			}
		}`))
	}))
	defer server.Close()

	// Создаем клиент
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	user, err := GetUser(apiClient, 123)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при получении пользователя: %v", err)
	}

	if user.ID != 123 {
		t.Errorf("Ожидался ID пользователя 123, получен %d", user.ID)
	}

	if user.Name != "Иван Иванов" {
		t.Errorf("Ожидалось имя пользователя 'Иван Иванов', получено '%s'", user.Name)
	}

	if user.Email != "ivan@example.com" {
		t.Errorf("Ожидался email пользователя 'ivan@example.com', получен '%s'", user.Email)
	}

	if !user.Rights.Leads {
		t.Errorf("Ожидались права на лиды (Rights.Leads=true)")
	}

	if !user.Rights.IsManager {
		t.Errorf("Ожидались права менеджера (Rights.IsManager=true)")
	}

	if user.Rights.IsAdmin {
		t.Errorf("Не ожидались права администратора (Rights.IsAdmin=false)")
	}
}

func TestGetCurrentUser(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/users/self"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 456,
			"name": "Петр Петров",
			"email": "petr@example.com",
			"lang": "ru",
			"is_active": true,
			"rights": {
				"leads": true,
				"contacts": true,
				"companies": true,
				"tasks": true,
				"mailbox": true,
				"catalog": true,
				"is_admin": true,
				"is_manager": false
			}
		}`))
	}))
	defer server.Close()

	// Создаем клиент
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	user, err := GetCurrentUser(apiClient)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при получении текущего пользователя: %v", err)
	}

	if user.ID != 456 {
		t.Errorf("Ожидался ID пользователя 456, получен %d", user.ID)
	}

	if user.Name != "Петр Петров" {
		t.Errorf("Ожидалось имя пользователя 'Петр Петров', получено '%s'", user.Name)
	}

	if user.Email != "petr@example.com" {
		t.Errorf("Ожидался email пользователя 'petr@example.com', получен '%s'", user.Email)
	}

	if !user.Rights.IsAdmin {
		t.Errorf("Ожидались права администратора (Rights.IsAdmin=true)")
	}

	if user.Rights.IsManager {
		t.Errorf("Не ожидались права менеджера (Rights.IsManager=false)")
	}
}

func TestListUsers(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/users"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Проверяем параметры запроса
		query := r.URL.Query()
		if query.Get("limit") != "50" {
			t.Errorf("Ожидался параметр limit=50, получен %s", query.Get("limit"))
		}
		if query.Get("page") != "1" {
			t.Errorf("Ожидался параметр page=1, получен %s", query.Get("page"))
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"_embedded": {
				"items": [
					{
						"id": 123,
						"name": "Иван Иванов",
						"email": "ivan@example.com",
						"lang": "ru",
						"is_active": true,
						"rights": {
							"leads": true,
							"contacts": true,
							"companies": true,
							"tasks": true,
							"mailbox": false,
							"catalog": false,
							"is_admin": false,
							"is_manager": true
						}
					},
					{
						"id": 456,
						"name": "Петр Петров",
						"email": "petr@example.com",
						"lang": "ru",
						"is_active": true,
						"rights": {
							"leads": true,
							"contacts": true,
							"companies": true,
							"tasks": true,
							"mailbox": true,
							"catalog": true,
							"is_admin": true,
							"is_manager": false
						}
					}
				]
			}
		}`))
	}))
	defer server.Close()

	// Создаем клиент
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	users, err := ListUsers(apiClient, 50, 1)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при получении списка пользователей: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("Ожидалось 2 пользователя, получено %d", len(users))
		return
	}

	if users[0].ID != 123 {
		t.Errorf("Ожидался ID первого пользователя 123, получен %d", users[0].ID)
	}

	if users[1].ID != 456 {
		t.Errorf("Ожидался ID второго пользователя 456, получен %d", users[1].ID)
	}

	if users[0].Name != "Иван Иванов" {
		t.Errorf("Ожидалось имя первого пользователя 'Иван Иванов', получено '%s'", users[0].Name)
	}

	if users[1].Name != "Петр Петров" {
		t.Errorf("Ожидалось имя второго пользователя 'Петр Петров', получено '%s'", users[1].Name)
	}
}
