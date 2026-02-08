package companies

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

func TestGetCompany(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		expectedPath := "/api/v4/companies/123"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

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

	apiClient := client.NewClient(server.URL, "test_api_key")

	company, err := Get(context.Background(), apiClient, 123)

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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		expectedPath := "/api/v4/companies"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"_embedded": {
				"companies": [
					{
						"id": 789,
						"name": "Новая компания",
						"responsible_user_id": 456,
						"created_at": 1609459200,
						"updated_at": 1609545600
					}
				]
			}
		}`))
	}))
	defer server.Close()

	apiClient := client.NewClient(server.URL, "test_api_key")

	companyToCreate := &Company{
		Name:              "Новая компания",
		ResponsibleUserID: 456,
	}

	createdCompany, err := Create(context.Background(), apiClient, companyToCreate)

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
