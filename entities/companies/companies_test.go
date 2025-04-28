package companies

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

func TestGetCompany(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/companies/123"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 123,
			"name": "Тестовая компания",
			"responsible_user_id": 456,
			"created_at": 1609459200,
			"updated_at": 1609545600
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	company, err := GetCompany(apiClient, 123)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при получении компании: %v", err)
	}

	if company.ID != 123 {
		t.Errorf("Ожидался ID компании 123, получен %d", company.ID)
	}

	if company.Name != "Тестовая компания" {
		t.Errorf("Ожидалось название компании 'Тестовая компания', получено '%s'", company.Name)
	}

	if company.ResponsibleUserID != 456 {
		t.Errorf("Ожидался ID ответственного пользователя 456, получен %d", company.ResponsibleUserID)
	}
}

func TestCreateCompany(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/companies"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 789,
			"name": "Новая компания",
			"responsible_user_id": 456,
			"created_at": 1609459200,
			"updated_at": 1609545600
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Создаем компанию для теста
	companyToCreate := &Company{
		Name:              "Новая компания",
		ResponsibleUserID: 456,
	}

	// Вызываем тестируемый метод
	createdCompany, err := CreateCompany(apiClient, companyToCreate)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при создании компании: %v", err)
	}

	if createdCompany.ID != 789 {
		t.Errorf("Ожидался ID компании 789, получен %d", createdCompany.ID)
	}

	if createdCompany.Name != "Новая компания" {
		t.Errorf("Ожидалось название компании 'Новая компания', получено '%s'", createdCompany.Name)
	}
}
