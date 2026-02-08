package contacts

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

// TestGetContactWithOptions проверяет получение контакта с дополнительными параметрами
func TestGetContactWithOptions(t *testing.T) {
	t.Run("С параметром WithCompanies", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				t.Errorf("Ожидался метод GET, получен %s", r.Method)
			}

			expectedPath := "/api/v4/contacts/123"
			if r.URL.Path != expectedPath {
				t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
			}

			with := r.URL.Query().Get("with")
			if with != "companies" {
				t.Errorf("Ожидался параметр with=companies, получен with=%s", with)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"id": 123,
				"name": "Тестовый контакт",
				"responsible_user_id": 456,
				"created_at": 1609459200,
				"updated_at": 1609545600,
				"_embedded": {
					"companies": [
						{
							"id": 789,
							"name": "Тестовая компания"
						}
					]
				}
			}`))
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		contact, err := Get(context.Background(), apiClient, 123, WithCompanies)

		if err != nil {
			t.Fatalf("Ошибка при получении контакта: %v", err)
		}

		if contact.ID != 123 {
			t.Errorf("Ожидался ID контакта 123, получен %d", contact.ID)
		}

		if contact.Name != "Тестовый контакт" {
			t.Errorf("Ожидалось имя контакта 'Тестовый контакт', получено '%s'", contact.Name)
		}

		if contact.Embedded == nil {
			t.Fatalf("Ожидался непустой Embedded, получен nil")
		}

		if len(contact.Embedded.Companies) != 1 {
			t.Errorf("Ожидалась 1 компания, получено %d", len(contact.Embedded.Companies))
		}

		if contact.Embedded.Companies[0].ID != 789 {
			t.Errorf("Ожидался ID компании 789, получен %d", contact.Embedded.Companies[0].ID)
		}

		if contact.Embedded.Companies[0].Name != "Тестовая компания" {
			t.Errorf("Ожидалось имя компании 'Тестовая компания', получено '%s'", contact.Embedded.Companies[0].Name)
		}
	})

	t.Run("С несколькими параметрами", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			with := r.URL.Query().Get("with")
			if with != "companies" {
				t.Errorf("Ожидался параметр with=companies, получен with=%s", with)
			}

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"id": 123,
				"name": "Тестовый контакт",
				"_embedded": {
					"companies": [
						{ "id": 789, "name": "Тестовая компания" }
					]
				}
			}`))
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		contact, err := Get(context.Background(), apiClient, 123, WithCompanies)

		if err != nil {
			t.Fatalf("Ошибка при получении контакта: %v", err)
		}

		if contact.Embedded == nil || len(contact.Embedded.Companies) != 1 {
			t.Errorf("Ожидалась 1 связанная компания")
		}
	})

	t.Run("Ошибка при некорректном ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error": "Contact not found"}`))
		}))
		defer server.Close()

		apiClient := client.NewClient("http://non-existing-domain.example", "test_api_key")

		_, err := Get(context.Background(), apiClient, 999, WithCompanies)

		if err == nil {
			t.Errorf("Ожидалась ошибка, но ее не было")
		}
	})
}

// TestGetContactsWithOptions проверяет получение списка контактов с дополнительными параметрами
func TestGetContactsWithOptions(t *testing.T) {
	t.Run("С параметром WithCompanies", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				t.Errorf("Ожидался метод GET, получен %s", r.Method)
			}

			expectedPath := "/api/v4/contacts"
			if r.URL.Path != expectedPath {
				t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
			}

			query := r.URL.Query()
			if query.Get("page") != "1" {
				t.Errorf("Ожидался параметр page=1, получен %s", query.Get("page"))
			}
			if query.Get("limit") != "50" {
				t.Errorf("Ожидался параметр limit=50, получен %s", query.Get("limit"))
			}

			with := query.Get("with")
			if with != "companies" {
				t.Errorf("Ожидался параметр with=companies, получен with=%s", with)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"page": 1,
				"per_page": 50,
				"total": 2,
				"_embedded": {
					"contacts": [
						{
							"id": 123,
							"name": "Тестовый контакт 1",
							"_embedded": {
								"companies": [
									{
										"id": 789,
										"name": "Тестовая компания 1"
									}
								]
							}
						},
						{
							"id": 124,
							"name": "Тестовый контакт 2",
							"_embedded": {
								"companies": [
									{
										"id": 790,
										"name": "Тестовая компания 2"
									}
								]
							}
						}
					]
				}
			}`))
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		contacts, err := List(context.Background(), apiClient, 1, 50, WithCompanies)

		if err != nil {
			t.Fatalf("Ошибка при получении контактов: %v", err)
		}

		if len(contacts) != 2 {
			t.Errorf("Ожидалось 2 контакта, получено %d", len(contacts))
		}

		if contacts[0].ID != 123 {
			t.Errorf("Ожидался ID контакта 123, получен %d", contacts[0].ID)
		}

		if contacts[0].Name != "Тестовый контакт 1" {
			t.Errorf("Ожидалось имя контакта 'Тестовый контакт 1', получено '%s'", contacts[0].Name)
		}

		if contacts[0].Embedded == nil {
			t.Fatalf("Ожидался непустой Embedded для первого контакта, получен nil")
		}

		if len(contacts[0].Embedded.Companies) != 1 {
			t.Errorf("Ожидалась 1 компания у первого контакта, получено %d", len(contacts[0].Embedded.Companies))
		}

		if contacts[0].Embedded.Companies[0].ID != 789 {
			t.Errorf("Ожидался ID компании 789, получен %d", contacts[0].Embedded.Companies[0].ID)
		}

		if contacts[1].ID != 124 {
			t.Errorf("Ожидался ID контакта 124, получен %d", contacts[1].ID)
		}

		if contacts[1].Embedded.Companies[0].ID != 790 {
			t.Errorf("Ожидался ID компании 790, получен %d", contacts[1].Embedded.Companies[0].ID)
		}
	})

	t.Run("Пустой ответ", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			with := r.URL.Query().Get("with")
			if with != "companies" {
				t.Errorf("Ожидался параметр with=companies, получен with=%s", with)
			}

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"_embedded": {"contacts": []}}`))
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		contacts, err := List(context.Background(), apiClient, 1, 50, WithCompanies)

		if err != nil {
			t.Fatalf("Ошибка при получении контактов: %v", err)
		}

		if len(contacts) != 0 {
			t.Errorf("Ожидался пустой список контактов, получено %d контактов", len(contacts))
		}
	})

	t.Run("Ошибка сервера", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error": "Internal Server Error"}`))
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		_, err := List(context.Background(), apiClient, 1, 50, WithCompanies)

		if err == nil {
			t.Errorf("Ожидалась ошибка, но ее не было")
		}
	})
}

// TestLinkContactWithCompany проверяет функцию связывания контакта с компанией
func TestLinkContactWithCompany(t *testing.T) {
	t.Run("Успешное связывание", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Errorf("Ожидался метод POST, получен %s", r.Method)
			}

			expectedPath := "/api/v4/contacts/123/link"
			if r.URL.Path != expectedPath {
				t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
			}

			contentType := r.Header.Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Ожидался Content-Type application/json, получен %s", contentType)
			}

			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("Ошибка при чтении тела запроса: %v", err)
			}

			var data []map[string]any
			if err := json.Unmarshal(body, &data); err != nil {
				t.Fatalf("Ошибка при разборе JSON: %v", err)
			}

			if len(data) != 1 {
				t.Errorf("Ожидался 1 элемент в массиве, получено %d", len(data))
			}

			toItem := data[0]

			entityID, ok := toItem["to_entity_id"].(float64)
			if !ok || int(entityID) != 456 {
				t.Errorf("Ожидался to_entity_id=456, получено %v", toItem["to_entity_id"])
			}

			entityType, ok := toItem["to_entity_type"].(string)
			if !ok || entityType != "companies" {
				t.Errorf("Ожидался to_entity_type='companies', получено %v", toItem["to_entity_type"])
			}

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"success": true}`))
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		err := LinkWithCompany(context.Background(), apiClient, 123, 456)

		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}
	})

	t.Run("Ошибка при некорректном контакте", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error": "Contact not found"}`))
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		err := LinkWithCompany(context.Background(), apiClient, 999, 456)

		if err == nil {
			t.Errorf("Ожидалась ошибка, но ее не было")
		}
	})

	t.Run("Ошибка при некорректном запросе", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error": "Bad request"}`))
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		err := LinkWithCompany(context.Background(), apiClient, 123, -1)

		if err == nil {
			t.Errorf("Ожидалась ошибка, но ее не было")
		}
	})
}
