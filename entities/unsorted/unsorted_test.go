package unsorted

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chudno/amo_crm_sdk/client"
)

func TestCreateUnsortedLead(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/leads/unsorted/api"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{
			"uid": "test-unsorted-uid-123",
			"account_id": 12345,
			"_links": {
				"self": {
					"href": "https://example.amocrm.ru/api/v4/unsorted/1"
				}
			}
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Создаем неразобранную заявку
	now := time.Now().Unix()
	lead := &UnsortedLeadCreate{
		UnsortedBase: UnsortedBase{
			SourceName: "API Test",
			SourceType: SourceTypeAPI,
			Category:   CategoryTypeForms,
			PipelineID: 123,
			CreatedAt:  now,
		},
		LeadName: "Тестовая заявка",
		Price:    1000,
		Contact: &UnsortedContact{
			Name:  "Иван Иванов",
			Email: "ivan@example.com",
			Phone: "+79001234567",
		},
		ResponsibleUserID: 456,
		PipelineType:      PipelineTypeLead,
	}

	// Вызываем тестируемый метод
	response, err := CreateUnsortedLead(apiClient, lead)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при создании неразобранной заявки: %v", err)
	}

	if response.UID != "test-unsorted-uid-123" {
		t.Errorf("Ожидался UID test-unsorted-uid-123, получен %s", response.UID)
	}

	if response.AccountID != 12345 {
		t.Errorf("Ожидался AccountID 12345, получен %d", response.AccountID)
	}
}

func TestCreateUnsortedContact(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/contacts/unsorted/api"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{
			"uid": "test-unsorted-contact-uid-123",
			"account_id": 12345,
			"_links": {
				"self": {
					"href": "https://example.amocrm.ru/api/v4/unsorted/1"
				}
			}
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Создаем неразобранную заявку контакта
	now := time.Now().Unix()
	contact := &UnsortedContactCreate{
		UnsortedBase: UnsortedBase{
			SourceName: "API Test",
			SourceType: SourceTypeAPI,
			Category:   CategoryTypeForms,
			CreatedAt:  now,
		},
		Contact: &UnsortedContact{
			Name:  "Петр Петров",
			Email: "petr@example.com",
			Phone: "+79001234568",
		},
		ResponsibleUserID: 456,
	}

	// Вызываем тестируемый метод
	response, err := CreateUnsortedContact(apiClient, contact)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при создании неразобранного контакта: %v", err)
	}

	if response.UID != "test-unsorted-contact-uid-123" {
		t.Errorf("Ожидался UID test-unsorted-contact-uid-123, получен %s", response.UID)
	}

	if response.AccountID != 12345 {
		t.Errorf("Ожидался AccountID 12345, получен %d", response.AccountID)
	}
}

func TestGetUnsortedLeads(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/leads/unsorted"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Проверяем параметры запроса
		expectedPage := "1"
		if r.URL.Query().Get("page") != expectedPage {
			t.Errorf("Ожидался параметр page=%s, получен %s", expectedPage, r.URL.Query().Get("page"))
		}

		expectedLimit := "50"
		if r.URL.Query().Get("limit") != expectedLimit {
			t.Errorf("Ожидался параметр limit=%s, получен %s", expectedLimit, r.URL.Query().Get("limit"))
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"_embedded": {
				"unsorted": [
					{
						"id": "unsorted-lead-1",
						"uid": "unsorted-uid-1",
						"source_uid": "src-1",
						"created_at": 1609459200,
						"pipeline_id": 123,
						"category": "forms",
						"source_type": "api",
						"source_name": "Test API",
						"pipeline_type": "lead",
						"account_id": 12345,
						"_embedded": {
							"leads": [
								{
									"id": 456,
									"name": "Тестовая сделка",
									"_links": {
										"self": {
											"href": "https://example.amocrm.ru/api/v4/leads/456"
										}
									}
								}
							]
						},
						"_links": {
							"self": {
								"href": "https://example.amocrm.ru/api/v4/leads/unsorted/unsorted-uid-1"
							}
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
	items, err := GetUnsortedLeads(apiClient, 1, 50, nil)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при получении неразобранных заявок: %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("Ожидалось получение 1 заявки, получено %d", len(items))
	}

	// Проверяем содержимое первой заявки
	expectedItem := UnsortedItem{
		ID:           "unsorted-lead-1",
		UID:          "unsorted-uid-1",
		SourceUID:    "src-1",
		CreatedAt:    1609459200,
		PipelineID:   123,
		Category:     CategoryTypeForms,
		SourceType:   SourceTypeAPI,
		SourceName:   "Test API",
		PipelineType: PipelineTypeLead,
		AccountID:    12345,
		Embedded: &struct {
			Contacts []struct {
				ID    int    `json:"id"`
				Name  string `json:"name"`
				Links struct {
					Self struct {
						Href string `json:"href"`
					} `json:"self"`
				} `json:"_links"`
			} `json:"contacts,omitempty"`
			Companies []struct {
				ID    int    `json:"id"`
				Name  string `json:"name"`
				Links struct {
					Self struct {
						Href string `json:"href"`
					} `json:"self"`
				} `json:"_links"`
			} `json:"companies,omitempty"`
			Leads []struct {
				ID    int    `json:"id"`
				Name  string `json:"name"`
				Links struct {
					Self struct {
						Href string `json:"href"`
					} `json:"self"`
				} `json:"_links"`
			} `json:"leads,omitempty"`
		}{
			Leads: []struct {
				ID    int    `json:"id"`
				Name  string `json:"name"`
				Links struct {
					Self struct {
						Href string `json:"href"`
					} `json:"self"`
				} `json:"_links"`
			}{
				{
					ID:   456,
					Name: "Тестовая сделка",
					Links: struct {
						Self struct {
							Href string `json:"href"`
						} `json:"self"`
					}{
						Self: struct {
							Href string `json:"href"`
						}{
							Href: "https://example.amocrm.ru/api/v4/leads/456",
						},
					},
				},
			},
		},
		Links: struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
		}{
			Self: struct {
				Href string `json:"href"`
			}{
				Href: "https://example.amocrm.ru/api/v4/leads/unsorted/unsorted-uid-1",
			},
		},
	}

	// Проверяем основные поля (полное сравнение сложно из-за вложенных структур)
	if items[0].ID != expectedItem.ID ||
		items[0].UID != expectedItem.UID ||
		items[0].SourceUID != expectedItem.SourceUID ||
		items[0].CreatedAt != expectedItem.CreatedAt ||
		items[0].PipelineID != expectedItem.PipelineID ||
		items[0].Category != expectedItem.Category ||
		items[0].SourceType != expectedItem.SourceType ||
		items[0].PipelineType != expectedItem.PipelineType ||
		items[0].AccountID != expectedItem.AccountID {
		t.Errorf("Полученная заявка не соответствует ожидаемой")
	}

	// Проверяем вложенные структуры
	if items[0].Embedded != nil && items[0].Embedded.Leads != nil && len(items[0].Embedded.Leads) > 0 {
		if items[0].Embedded.Leads[0].ID != 456 || items[0].Embedded.Leads[0].Name != "Тестовая сделка" {
			t.Errorf("Вложенная сделка не соответствует ожидаемой")
		}
	} else {
		t.Errorf("Вложенная сделка отсутствует")
	}
}

func TestAcceptUnsortedLead(t *testing.T) {
	// Тестовый UID неразобранной заявки
	unsortedUID := "test-unsorted-uid-123"

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/leads/unsorted/%s/accept", unsortedUID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"_links": {
				"lead": {
					"id": 789
				}
			}
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	leadID, err := AcceptUnsortedLead(apiClient, unsortedUID, 123, 456)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при принятии неразобранной заявки: %v", err)
	}

	if leadID != 789 {
		t.Errorf("Ожидался ID сделки 789, получен %d", leadID)
	}
}

func TestDeclineUnsortedLead(t *testing.T) {
	// Тестовый UID неразобранной заявки
	unsortedUID := "test-unsorted-uid-123"

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "DELETE" {
			t.Errorf("Ожидался метод DELETE, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/leads/unsorted/%s/decline", unsortedUID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	err := DeclineUnsortedLead(apiClient, unsortedUID)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при отклонении неразобранной заявки: %v", err)
	}
}

func TestGetUnsortedSummary(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/unsorted/summary"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"total": {
				"count": 10
			},
			"accepted": {
				"count": 6
			},
			"declined": {
				"count": 2
			},
			"unprocessed": {
				"count": 2
			}
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	summary, err := GetUnsortedSummary(apiClient)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при получении сводки по неразобранным заявкам: %v", err)
	}

	// Проверяем содержимое сводки
	total, ok := summary["total"].(map[string]interface{})
	if !ok {
		t.Fatalf("Не удалось получить поле total из сводки")
	}
	totalCount, ok := total["count"].(float64)
	if !ok {
		t.Fatalf("Не удалось получить поле count из total")
	}
	if totalCount != 10 {
		t.Errorf("Ожидалось значение total.count равное 10, получено %v", totalCount)
	}

	// Проверяем количество обработанных заявок
	accepted, ok := summary["accepted"].(map[string]interface{})
	if !ok {
		t.Fatalf("Не удалось получить поле accepted из сводки")
	}
	acceptedCount, ok := accepted["count"].(float64)
	if !ok {
		t.Fatalf("Не удалось получить поле count из accepted")
	}
	if acceptedCount != 6 {
		t.Errorf("Ожидалось значение accepted.count равное 6, получено %v", acceptedCount)
	}
}
