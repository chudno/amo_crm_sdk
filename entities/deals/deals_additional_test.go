package deals

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

// TestGetDealWithRelatedEntities проверяет получение сделки с связанными сущностями
func TestGetDealWithRelatedEntities(t *testing.T) {
	t.Run("С контактами и компаниями", func(t *testing.T) {
		// Создаем тестовый сервер
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Проверяем метод запроса
			if r.Method != "GET" {
				t.Errorf("Ожидался метод GET, получен %s", r.Method)
			}

			// Проверяем путь запроса
			expectedPath := "/api/v4/deals/123"
			if r.URL.Path != expectedPath {
				t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
			}

			// Проверяем параметры запроса
			query := r.URL.Query()
			with := query.Get("with")
			if with != "contacts,companies" && with != "companies,contacts" {
				t.Errorf("Ожидалось with=contacts,companies (или companies,contacts), получено %s", with)
			}

			// Отправляем ответ
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(`{
				"id": 123,
				"name": "Тестовая сделка",
				"value": 10000,
				"responsible_user_id": 456,
				"created_at": 1609459200,
				"updated_at": 1609545600,
				"_embedded": {
					"contacts": [
						{
							"id": 789,
							"name": "Иван Иванов"
						}
					],
					"companies": [
						{
							"id": 101,
							"name": "ООО Тест"
						}
					]
				}
			}`)); err != nil {
				t.Fatalf("Ошибка при записи ответа: %v", err)
			}
		}))
		defer server.Close()

		// Создаем клиент
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемый метод
		deal, err := GetDeal(apiClient, 123, WithContacts, WithCompanies)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при получении сделки: %v", err)
		}

		if deal.ID != 123 {
			t.Errorf("Ожидался ID сделки 123, получен %d", deal.ID)
		}

		if deal.Name != "Тестовая сделка" {
			t.Errorf("Ожидалось имя сделки 'Тестовая сделка', получено '%s'", deal.Name)
		}

		// Проверяем связанные сущности
		if deal.Embedded == nil {
			t.Fatalf("Ожидались связанные сущности, но Embedded == nil")
		}

		if len(deal.Embedded.Contacts) != 1 {
			t.Errorf("Ожидался 1 контакт, получено %d", len(deal.Embedded.Contacts))
		} else if deal.Embedded.Contacts[0].ID != 789 {
			t.Errorf("Ожидался ID контакта 789, получен %d", deal.Embedded.Contacts[0].ID)
		}

		if len(deal.Embedded.Companies) != 1 {
			t.Errorf("Ожидалась 1 компания, получено %d", len(deal.Embedded.Companies))
		} else if deal.Embedded.Companies[0].ID != 101 {
			t.Errorf("Ожидался ID компании 101, получен %d", deal.Embedded.Companies[0].ID)
		}
	})
}

// TestGetDealErrors проверяет обработку ошибок при получении сделки
func TestGetDealErrors(t *testing.T) {
	t.Run("Сделка не найдена", func(t *testing.T) {
		// Используем несуществующий домен для вызова ошибки
		apiClient := client.NewClient("http://non-existent-domain.example", "test_api_key")

		// Вызываем тестируемую функцию
		_, err := GetDeal(apiClient, 999) // несуществующий ID

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
			if _, err := w.Write([]byte(`{"id": 123, "name": "Тестовая сделка", "value":`)); err != nil {
				t.Fatalf("Ошибка при записи ответа: %v", err)
			}
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемую функцию
		_, err := GetDeal(apiClient, 123)

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка из-за некорректного JSON, но она не возникла")
		}
	})

	t.Run("Ошибка HTTP", func(t *testing.T) {
		// Используем несуществующий домен для вызова ошибки
		apiClient := client.NewClient("http://non-existent-domain.example", "test_api_key")

		// Вызываем тестируемую функцию
		_, err := GetDeal(apiClient, 123)

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка HTTP, но она не возникла")
		}
	})
}

// TestCreateDealErrors проверяет обработку ошибок при создании сделки
func TestCreateDealErrors(t *testing.T) {
	t.Run("Некорректные данные", func(t *testing.T) {
		// Используем несуществующий домен для вызова ошибки
		apiClient := client.NewClient("http://non-existent-domain.example", "test_api_key")

		// Создаем сделку для теста
		deal := &Deal{
			Name:  "Тестовая сделка",
			Value: 15000,
		}

		// Вызываем тестируемую функцию
		_, err := CreateDeal(apiClient, deal)

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка, но она не возникла")
		}
	})

	t.Run("Ошибка сервера", func(t *testing.T) {
		// Второй тест с несуществующим доменом для проверки обработки ошибок
		apiClient := client.NewClient("http://another-non-existent-domain.example", "test_api_key")

		// Создаем тестовую сделку
		deal := &Deal{
			Name:  "Тестовая сделка",
			Value: 10000,
		}

		// Вызываем тестируемую функцию
		_, err := CreateDeal(apiClient, deal)

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка сервера, но она не возникла")
		}
	})
}

// TestUpdateDealErrors проверяет обработку ошибок при обновлении сделки
func TestUpdateDealErrors(t *testing.T) {
	t.Run("Сделка не найдена", func(t *testing.T) {
		// Используем несуществующий домен для вызова ошибки
		apiClient := client.NewClient("http://non-existent-domain.example", "test_api_key")

		// Создаем сделку с несуществующим ID
		dealToUpdate := &Deal{
			ID:    999,
			Name:  "Обновленная сделка",
			Value: 20000,
		}

		// Вызываем тестируемую функцию
		_, err := UpdateDeal(apiClient, dealToUpdate)

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка, но она не возникла")
		}
	})

	t.Run("Некорректный JSON-ответ", func(t *testing.T) {
		// Создаем тестовый сервер, который вернет некорректный JSON
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(`{"id": 456, "name": "Обновленная сделка", "value":`)); err != nil {
				t.Fatalf("Ошибка при записи ответа: %v", err)
			}
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Создаем сделку для обновления
		dealToUpdate := &Deal{
			ID:    456,
			Name:  "Обновленная сделка",
			Value: 25000,
		}

		// Вызываем тестируемую функцию
		_, err := UpdateDeal(apiClient, dealToUpdate)

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка из-за некорректного JSON, но она не возникла")
		}
	})

	t.Run("Ошибка HTTP", func(t *testing.T) {
		// Используем несуществующий домен для вызова ошибки
		apiClient := client.NewClient("http://non-existent-domain.example", "test_api_key")

		// Создаем сделку для обновления
		dealToUpdate := &Deal{
			ID:    456,
			Name:  "Обновленная сделка",
			Value: 25000,
		}

		// Вызываем тестируемую функцию
		_, err := UpdateDeal(apiClient, dealToUpdate)

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка HTTP, но она не возникла")
		}
	})
}

// TestGetDealsAdditional проверяет дополнительные сценарии при получении списка сделок
func TestGetDealsAdditional(t *testing.T) {
	t.Run("С фильтрацией", func(t *testing.T) {
		// Создаем тестовый сервер
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Проверяем метод запроса
			if r.Method != "GET" {
				t.Errorf("Ожидался метод GET, получен %s", r.Method)
			}

			// Проверяем URL запроса
			expectedPath := "/api/v4/deals"
			if r.URL.Path != expectedPath {
				t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
			}

			// Проверяем параметры запроса
			query := r.URL.Query()
			if query.Get("filter[status_id]") != "1" {
				t.Errorf("Ожидался параметр filter[status_id]=1, получен %s", query.Get("filter[status_id]"))
			}

			// Отправляем ответ
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(`{
				"page": 1,
				"per_page": 50,
				"total": 1,
				"_embedded": {
					"items": [
						{
							"id": 12345,
							"name": "Сделка в статусе 1",
							"value": 10000,
							"status_id": 1
						}
					]
				}
			}`)); err != nil {
				t.Fatalf("Ошибка при записи ответа: %v", err)
			}
		}))
		defer server.Close()

		// Создаем клиент
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Создаем фильтр
		filter := map[string]string{
			"filter[status_id]": "1",
		}

		// Вызываем тестируемую функцию
		deals, err := GetDeals(apiClient, 1, 50, filter)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при получении списка сделок: %v", err)
		}

		if len(deals) != 1 {
			t.Errorf("Ожидалась 1 сделка, получено %d", len(deals))
		} else {
			if deals[0].StatusID != 1 {
				t.Errorf("Ожидался статус сделки 1, получен %d", deals[0].StatusID)
			}
		}
	})

	t.Run("С опциями (контакты и компании)", func(t *testing.T) {
		// Создаем тестовый сервер
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Проверяем метод запроса
			if r.Method != "GET" {
				t.Errorf("Ожидался метод GET, получен %s", r.Method)
			}

			// Проверяем параметр with
			query := r.URL.Query()
			with := query.Get("with")
			if with != "contacts,companies" && with != "companies,contacts" {
				t.Errorf("Ожидалось with=contacts,companies (или companies,contacts), получено %s", with)
			}

			// Отправляем ответ
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(`{
				"page": 1,
				"per_page": 50,
				"total": 1,
				"_embedded": {
					"items": [
						{
							"id": 12345,
							"name": "Сделка с контактами и компаниями",
							"value": 10000,
							"_embedded": {
								"contacts": [
									{
										"id": 789,
										"name": "Иван Иванов"
									}
								],
								"companies": [
									{
										"id": 101,
										"name": "ООО Тест"
									}
								]
							}
						}
					]
				}
			}`)); err != nil {
				t.Fatalf("Ошибка при записи ответа: %v", err)
			}
		}))
		defer server.Close()

		// Создаем клиент
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемую функцию
		deals, err := GetDeals(apiClient, 1, 50, nil, WithContacts, WithCompanies)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при получении списка сделок: %v", err)
		}

		if len(deals) != 1 {
			t.Errorf("Ожидалась 1 сделка, получено %d", len(deals))
		} else {
			deal := deals[0]
			if deal.Embedded == nil {
				t.Fatalf("Ожидались связанные сущности, но Embedded == nil")
			}

			if len(deal.Embedded.Contacts) != 1 {
				t.Errorf("Ожидался 1 контакт, получено %d", len(deal.Embedded.Contacts))
			} else if deal.Embedded.Contacts[0].ID != 789 {
				t.Errorf("Ожидался ID контакта 789, получен %d", deal.Embedded.Contacts[0].ID)
			}

			if len(deal.Embedded.Companies) != 1 {
				t.Errorf("Ожидалась 1 компания, получено %d", len(deal.Embedded.Companies))
			} else if deal.Embedded.Companies[0].ID != 101 {
				t.Errorf("Ожидался ID компании 101, получен %d", deal.Embedded.Companies[0].ID)
			}
		}
	})

	t.Run("Некорректный JSON-ответ", func(t *testing.T) {
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
		_, err := GetDeals(apiClient, 1, 50, nil)

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка из-за некорректного JSON, но она не возникла")
		}
	})
}
