package companies

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

// TestUpdateCompany проверяет функцию обновления компании
func TestUpdateCompany(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "PATCH" {
			t.Errorf("Ожидался метод PATCH, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/companies/456"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 456,
			"name": "Обновленная компания",
			"responsible_user_id": 789,
			"created_at": 1609459200,
			"updated_at": 1609632000
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Создаем компанию для обновления
	companyToUpdate := &Company{
		ID:                456,
		Name:              "Обновленная компания",
		ResponsibleUserID: 789,
	}

	// Вызываем тестируемый метод
	updatedCompany, err := UpdateCompany(apiClient, companyToUpdate)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при обновлении компании: %v", err)
	}

	if updatedCompany.ID != 456 {
		t.Errorf("Ожидался ID компании 456, получен %d", updatedCompany.ID)
	}

	if updatedCompany.Name != "Обновленная компания" {
		t.Errorf("Ожидалось название компании 'Обновленная компания', получено '%s'", updatedCompany.Name)
	}

	if updatedCompany.ResponsibleUserID != 789 {
		t.Errorf("Ожидался ID ответственного пользователя 789, получен %d", updatedCompany.ResponsibleUserID)
	}
}

// TestGetCompanies проверяет функцию получения списка компаний
func TestGetCompanies(t *testing.T) {
	t.Run("Базовый запрос", func(t *testing.T) {
		// Создаем тестовый сервер
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Проверяем метод запроса
			if r.Method != "GET" {
				t.Errorf("Ожидался метод GET, получен %s", r.Method)
			}

			// Проверяем путь запроса
			expectedPath := "/api/v4/companies"
			if r.URL.Path != expectedPath {
				t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
			}

			// Проверяем параметры запроса
			query := r.URL.Query()
			if query.Get("page") != "1" {
				t.Errorf("Ожидался параметр page=1, получен %s", query.Get("page"))
			}
			if query.Get("limit") != "50" {
				t.Errorf("Ожидался параметр limit=50, получен %s", query.Get("limit"))
			}

			// Отправляем ответ
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"page": 1,
				"per_page": 50,
				"total": 2,
				"_embedded": {
					"items": [
						{
							"id": 123,
							"name": "Компания 1",
							"responsible_user_id": 456
						},
						{
							"id": 789,
							"name": "Компания 2",
							"responsible_user_id": 456
						}
					]
				}
			}`))
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемый метод
		companies, err := GetCompanies(apiClient, 1, 50)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при получении списка компаний: %v", err)
		}

		if len(companies) != 2 {
			t.Errorf("Ожидалось 2 компании, получено %d", len(companies))
		}

		// Проверяем первую компанию
		if companies[0].ID != 123 {
			t.Errorf("Ожидался ID компании 123, получен %d", companies[0].ID)
		}
		if companies[0].Name != "Компания 1" {
			t.Errorf("Ожидалось название компании 'Компания 1', получено '%s'", companies[0].Name)
		}

		// Проверяем вторую компанию
		if companies[1].ID != 789 {
			t.Errorf("Ожидался ID компании 789, получен %d", companies[1].ID)
		}
		if companies[1].Name != "Компания 2" {
			t.Errorf("Ожидалось название компании 'Компания 2', получено '%s'", companies[1].Name)
		}
	})

	t.Run("С опцией WithContacts", func(t *testing.T) {
		// Создаем тестовый сервер
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Проверяем метод запроса
			if r.Method != "GET" {
				t.Errorf("Ожидался метод GET, получен %s", r.Method)
			}

			// Проверяем путь запроса
			expectedPath := "/api/v4/companies"
			if r.URL.Path != expectedPath {
				t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
			}

			// Проверяем параметры запроса
			query := r.URL.Query()
			if query.Get("with") != "contacts" {
				t.Errorf("Ожидался параметр with=contacts, получен %s", query.Get("with"))
			}

			// Отправляем ответ с вложенными контактами
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"page": 1,
				"per_page": 50,
				"total": 1,
				"_embedded": {
					"items": [
						{
							"id": 123,
							"name": "Компания с контактами",
							"_embedded": {
								"contacts": [
									{
										"id": 456,
										"name": "Контакт 1"
									},
									{
										"id": 789,
										"name": "Контакт 2"
									}
								]
							}
						}
					]
				}
			}`))
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемый метод с опцией WithContacts
		companies, err := GetCompanies(apiClient, 1, 50, WithContacts)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при получении списка компаний: %v", err)
		}

		if len(companies) != 1 {
			t.Errorf("Ожидалась 1 компания, получено %d", len(companies))
		}

		// Проверяем компанию и ее контакты
		if companies[0].ID != 123 {
			t.Errorf("Ожидался ID компании 123, получен %d", companies[0].ID)
		}
		if companies[0].Name != "Компания с контактами" {
			t.Errorf("Ожидалось название компании 'Компания с контактами', получено '%s'", companies[0].Name)
		}

		// Проверяем наличие и содержимое контактов
		if companies[0].Embedded == nil {
			t.Errorf("Ожидалось наличие вложенных данных, но Embedded == nil")
		} else if len(companies[0].Embedded.Contacts) != 2 {
			t.Errorf("Ожидалось 2 контакта, получено %d", len(companies[0].Embedded.Contacts))
		} else {
			if companies[0].Embedded.Contacts[0].ID != 456 {
				t.Errorf("Ожидался ID контакта 456, получен %d", companies[0].Embedded.Contacts[0].ID)
			}
			if companies[0].Embedded.Contacts[0].Name != "Контакт 1" {
				t.Errorf("Ожидалось имя контакта 'Контакт 1', получено '%s'", companies[0].Embedded.Contacts[0].Name)
			}
		}
	})

	t.Run("Пустой ответ", func(t *testing.T) {
		// Создаем тестовый сервер
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Отправляем пустой список компаний
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"page": 1,
				"per_page": 50,
				"total": 0,
				"_embedded": {
					"items": []
				}
			}`))
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемый метод
		companies, err := GetCompanies(apiClient, 1, 50)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при получении пустого списка компаний: %v", err)
		}

		if len(companies) != 0 {
			t.Errorf("Ожидалось 0 компаний, получено %d", len(companies))
		}
	})

	t.Run("Ошибка сервера", func(t *testing.T) {
		// Создаем тестовый сервер, который НЕ отвечает (сетевая ошибка)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Закрываем соединение без ответа
			hj, ok := w.(http.Hijacker)
			if !ok {
				t.Fatalf("Не удалось преобразовать ResponseWriter в Hijacker")
			}
			conn, _, _ := hj.Hijack()
			conn.Close()
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемый метод
		_, err := GetCompanies(apiClient, 1, 50)

		// Проверяем результаты
		if err == nil {
			t.Errorf("Ожидалась ошибка, получен nil")
		}
	})
}

// TestGetCompanyWithOptions проверяет функцию получения компании со связанными сущностями
func TestGetCompanyWithOptions(t *testing.T) {
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

		// Проверяем параметры запроса
		query := r.URL.Query()
		if query.Get("with") != "contacts" {
			t.Errorf("Ожидался параметр with=contacts, получен %s", query.Get("with"))
		}

		// Отправляем ответ с вложенными контактами
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 123,
			"name": "Компания с контактами",
			"_embedded": {
				"contacts": [
					{
						"id": 456,
						"name": "Контакт 1"
					},
					{
						"id": 789,
						"name": "Контакт 2"
					}
				]
			}
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод с опцией WithContacts
	company, err := GetCompany(apiClient, 123, WithContacts)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при получении компании: %v", err)
	}

	if company.ID != 123 {
		t.Errorf("Ожидался ID компании 123, получен %d", company.ID)
	}

	// Проверяем наличие и содержимое контактов
	if company.Embedded == nil {
		t.Errorf("Ожидалось наличие вложенных данных, но Embedded == nil")
	} else if len(company.Embedded.Contacts) != 2 {
		t.Errorf("Ожидалось 2 контакта, получено %d", len(company.Embedded.Contacts))
	} else {
		if company.Embedded.Contacts[0].ID != 456 {
			t.Errorf("Ожидался ID контакта 456, получен %d", company.Embedded.Contacts[0].ID)
		}
		if company.Embedded.Contacts[0].Name != "Контакт 1" {
			t.Errorf("Ожидалось имя контакта 'Контакт 1', получено '%s'", company.Embedded.Contacts[0].Name)
		}
	}
}
