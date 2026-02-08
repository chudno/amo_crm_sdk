package users

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

func TestGetUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		expectedPath := "/api/v4/users/123"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 123,
			"name": "Иван Иванов",
			"email": "ivan@example.com",
			"lang": "ru",
			"is_active": true,
			"rights": {
				"leads": {"view": "M", "edit": "M", "add": "D", "delete": "M", "export": "M"},
				"contacts": {"view": "A", "edit": "A", "add": "A", "delete": "A", "export": "A"},
				"companies": {"view": "A", "edit": "A", "add": "D", "delete": "A", "export": "A"},
				"tasks": {"edit": "A", "delete": "A"},
				"mail_access": false,
				"catalog_access": true,
				"is_admin": false,
				"is_manager": true
			}
		}`))
	}))
	defer server.Close()

	apiClient := client.NewClient(server.URL, "test_api_key")

	user, err := Get(context.Background(), apiClient, 123)

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

	if user.Rights.Leads == nil {
		t.Fatal("Ожидались права на лиды (Rights.Leads не nil)")
	}
	if user.Rights.Leads.View != AccessOwn {
		t.Errorf("Ожидался уровень доступа Leads.View=%q, получен %q", AccessOwn, user.Rights.Leads.View)
	}
	if user.Rights.Leads.Add != AccessDeny {
		t.Errorf("Ожидался уровень доступа Leads.Add=%q, получен %q", AccessDeny, user.Rights.Leads.Add)
	}

	if !user.Rights.IsManager {
		t.Errorf("Ожидались права менеджера (Rights.IsManager=true)")
	}

	if user.Rights.IsAdmin {
		t.Errorf("Не ожидались права администратора (Rights.IsAdmin=false)")
	}

	if !user.Rights.CatalogAccess {
		t.Errorf("Ожидался доступ к каталогам (CatalogAccess=true)")
	}
}

func TestGetCurrentUser(t *testing.T) {
	// Создаем тестовый сервер, который обрабатывает два маршрута:
	// 1. /api/v4/account — возвращает current_user_id
	// 2. /api/v4/users/456 — возвращает данные пользователя
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		switch r.URL.Path {
		case "/api/v4/account":
			_, _ = w.Write([]byte(`{
				"id": 1,
				"name": "Test Account",
				"subdomain": "test",
				"current_user_id": 456
			}`))
		case "/api/v4/users/456":
			_, _ = w.Write([]byte(`{
				"id": 456,
				"name": "Петр Петров",
				"email": "petr@example.com",
				"lang": "ru",
				"is_active": true,
				"rights": {
					"leads": {"view": "A", "edit": "A", "add": "A", "delete": "A", "export": "A"},
					"contacts": {"view": "A", "edit": "A", "add": "A", "delete": "A", "export": "A"},
					"companies": {"view": "A", "edit": "A", "add": "A", "delete": "A", "export": "A"},
					"tasks": {"edit": "A", "delete": "A"},
					"mail_access": true,
					"catalog_access": true,
					"is_admin": true,
					"is_manager": false
				}
			}`))
		default:
			t.Errorf("Неожиданный путь запроса: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	apiClient := client.NewClient(server.URL, "test_api_key")

	user, err := GetCurrent(context.Background(), apiClient)

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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		expectedPath := "/api/v4/users"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		query := r.URL.Query()
		if query.Get("limit") != "50" {
			t.Errorf("Ожидался параметр limit=50, получен %s", query.Get("limit"))
		}
		if query.Get("page") != "1" {
			t.Errorf("Ожидался параметр page=1, получен %s", query.Get("page"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"_embedded": {
				"users": [
					{
						"id": 123,
						"name": "Иван Иванов",
						"email": "ivan@example.com",
						"lang": "ru",
						"is_active": true,
						"rights": {
							"leads": {"view": "M", "edit": "M", "add": "D", "delete": "M", "export": "M"},
							"contacts": {"view": "A", "edit": "A", "add": "A", "delete": "A", "export": "A"},
							"companies": {"view": "A", "edit": "A", "add": "D", "delete": "A", "export": "A"},
							"tasks": {"edit": "A", "delete": "A"},
							"mail_access": false,
							"catalog_access": false,
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
							"leads": {"view": "A", "edit": "A", "add": "A", "delete": "A", "export": "A"},
							"contacts": {"view": "A", "edit": "A", "add": "A", "delete": "A", "export": "A"},
							"companies": {"view": "A", "edit": "A", "add": "A", "delete": "A", "export": "A"},
							"tasks": {"edit": "A", "delete": "A"},
							"mail_access": true,
							"catalog_access": true,
							"is_admin": true,
							"is_manager": false
						}
					}
				]
			}
		}`))
	}))
	defer server.Close()

	apiClient := client.NewClient(server.URL, "test_api_key")

	users, err := List(context.Background(), apiClient, 50, 1)

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
