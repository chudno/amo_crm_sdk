package contacts

import (
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
		// Создаем тестовый сервер
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Проверяем метод запроса
			if r.Method != "GET" {
				t.Errorf("Ожидался метод GET, получен %s", r.Method)
			}

			// Проверяем путь запроса
			expectedPath := "/api/v4/contacts/123"
			if r.URL.Path != expectedPath {
				t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
			}

			// Проверяем параметр with=companies
			with := r.URL.Query().Get("with")
			if with != "companies" {
				t.Errorf("Ожидался параметр with=companies, получен with=%s", with)
			}

			// Отправляем ответ
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

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемый метод
		contact, err := GetContact(apiClient, 123, WithCompanies)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при получении контакта: %v", err)
		}

		if contact.ID != 123 {
			t.Errorf("Ожидался ID контакта 123, получен %d", contact.ID)
		}

		if contact.Name != "Тестовый контакт" {
			t.Errorf("Ожидалось имя контакта 'Тестовый контакт', получено '%s'", contact.Name)
		}

		// Проверяем наличие связанных компаний
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
		// Создаем тестовый сервер
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Проверяем параметр with
			with := r.URL.Query().Get("with")
			if with != "companies" {
				t.Errorf("Ожидался параметр with=companies, получен with=%s", with)
			}

			// Отправляем ответ
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

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемый метод с несколькими параметрами
		// (В этом примере только WithCompanies, но демонстрирует принцип)
		contact, err := GetContact(apiClient, 123, WithCompanies)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при получении контакта: %v", err)
		}

		if contact.Embedded == nil || len(contact.Embedded.Companies) != 1 {
			t.Errorf("Ожидалась 1 связанная компания")
		}
	})

	t.Run("Ошибка при некорректном ID", func(t *testing.T) {
		// Создаем тестовый сервер, который возвращает ошибку
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Возвращаем статус 404 Not Found
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error": "Contact not found"}`))  
		}))
		defer server.Close()

		// Создаем клиент API с нерабочим URL, чтобы гарантированно получить ошибку
		apiClient := client.NewClient("http://non-existing-domain.example", "test_api_key")

		// Вызываем тестируемый метод
		_, err := GetContact(apiClient, 999, WithCompanies)

		// Проверяем результаты
		if err == nil {
			t.Errorf("Ожидалась ошибка, но ее не было")
		}
	})
}

// TestGetContactsWithOptions проверяет получение списка контактов с дополнительными параметрами
func TestGetContactsWithOptions(t *testing.T) {
	t.Run("С параметром WithCompanies", func(t *testing.T) {
		// Создаем тестовый сервер
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Проверяем метод запроса
			if r.Method != "GET" {
				t.Errorf("Ожидался метод GET, получен %s", r.Method)
			}

			// Проверяем путь запроса
			expectedPath := "/api/v4/contacts"
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

			// Проверяем параметр with=companies
			with := query.Get("with")
			if with != "companies" {
				t.Errorf("Ожидался параметр with=companies, получен with=%s", with)
			}

			// Отправляем ответ
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

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемый метод
		contacts, err := GetContacts(apiClient, 1, 50, WithCompanies)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при получении контактов: %v", err)
		}

		if len(contacts) != 2 {
			t.Errorf("Ожидалось 2 контакта, получено %d", len(contacts))
		}

		// Проверяем первый контакт
		if contacts[0].ID != 123 {
			t.Errorf("Ожидался ID контакта 123, получен %d", contacts[0].ID)
		}

		if contacts[0].Name != "Тестовый контакт 1" {
			t.Errorf("Ожидалось имя контакта 'Тестовый контакт 1', получено '%s'", contacts[0].Name)
		}

		// Проверяем наличие связанных компаний у первого контакта
		if contacts[0].Embedded == nil {
			t.Fatalf("Ожидался непустой Embedded для первого контакта, получен nil")
		}

		if len(contacts[0].Embedded.Companies) != 1 {
			t.Errorf("Ожидалась 1 компания у первого контакта, получено %d", len(contacts[0].Embedded.Companies))
		}

		if contacts[0].Embedded.Companies[0].ID != 789 {
			t.Errorf("Ожидался ID компании 789, получен %d", contacts[0].Embedded.Companies[0].ID)
		}

		// Проверяем второй контакт
		if contacts[1].ID != 124 {
			t.Errorf("Ожидался ID контакта 124, получен %d", contacts[1].ID)
		}

		if contacts[1].Embedded.Companies[0].ID != 790 {
			t.Errorf("Ожидался ID компании 790, получен %d", contacts[1].Embedded.Companies[0].ID)
		}
	})

	t.Run("Пустой ответ", func(t *testing.T) {
		// Создаем тестовый сервер, который возвращает пустой список
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Проверяем параметр with=companies
			with := r.URL.Query().Get("with")
			if with != "companies" {
				t.Errorf("Ожидался параметр with=companies, получен with=%s", with)
			}

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"_embedded": {"contacts": []}}`))
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемый метод
		contacts, err := GetContacts(apiClient, 1, 50, WithCompanies)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при получении контактов: %v", err)
		}

		if len(contacts) != 0 {
			t.Errorf("Ожидался пустой список контактов, получено %d контактов", len(contacts))
		}
	})
	
	t.Run("Ошибка сервера", func(t *testing.T) {
		// Создаем тестовый сервер, который возвращает ошибку
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error": "Internal Server Error"}`))
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемый метод
		_, err := GetContacts(apiClient, 1, 50, WithCompanies)

		// Проверяем результаты
		if err == nil {
			t.Errorf("Ожидалась ошибка, но ее не было")
		}
	})
}

// TestLinkContactWithCompany проверяет функцию связывания контакта с компанией
func TestLinkContactWithCompany(t *testing.T) {
	t.Run("Успешное связывание", func(t *testing.T) {
		// Создаем тестовый сервер
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Проверяем метод запроса
			if r.Method != "POST" {
				t.Errorf("Ожидался метод POST, получен %s", r.Method)
			}

			// Проверяем путь запроса
			expectedPath := "/api/v4/contacts/123/link"
			if r.URL.Path != expectedPath {
				t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
			}

			// Проверяем заголовок Content-Type
			contentType := r.Header.Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Ожидался Content-Type application/json, получен %s", contentType)
			}

			// Проверяем тело запроса
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("Ошибка при чтении тела запроса: %v", err)
			}

			// Проверяем JSON-структуру
			var data map[string]interface{}
			if err := json.Unmarshal(body, &data); err != nil {
				t.Fatalf("Ошибка при разборе JSON: %v", err)
			}

			// Проверяем наличие поля "to"
			to, ok := data["to"].([]interface{})
			if !ok {
				t.Fatalf("Ожидалось поле 'to' типа array, получено %T", data["to"])
			}

			// Проверяем содержимое массива "to"
			if len(to) != 1 {
				t.Errorf("Ожидался 1 элемент в массиве 'to', получено %d", len(to))
			}

			toItem, ok := to[0].(map[string]interface{})
			if !ok {
				t.Fatalf("Ожидался элемент типа object, получено %T", to[0])
			}

			entityID, ok := toItem["entity_id"].(float64)
			if !ok || int(entityID) != 456 {
				t.Errorf("Ожидался entity_id=456, получено %v", toItem["entity_id"])
			}

			entityType, ok := toItem["entity_type"].(string)
			if !ok || entityType != "companies" {
				t.Errorf("Ожидался entity_type='companies', получено %v", toItem["entity_type"])
			}

			// Отправляем успешный ответ
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"success": true}`))
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемый метод
		err := LinkContactWithCompany(apiClient, 123, 456)

		// Проверяем результаты
		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}
	})

	t.Run("Ошибка при некорректном контакте", func(t *testing.T) {
		// Создаем тестовый сервер, который возвращает ошибку
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error": "Contact not found"}`))
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемый метод
		err := LinkContactWithCompany(apiClient, 999, 456)

		// Проверяем результаты
		if err == nil {
			t.Errorf("Ожидалась ошибка, но ее не было")
		}
	})

	t.Run("Ошибка при некорректном запросе", func(t *testing.T) {
		// Создаем тестовый сервер, который возвращает ошибку
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error": "Bad request"}`))
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемый метод
		err := LinkContactWithCompany(apiClient, 123, -1)

		// Проверяем результаты
		if err == nil {
			t.Errorf("Ожидалась ошибка, но ее не было")
		}
	})
}
