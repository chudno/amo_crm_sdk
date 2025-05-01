package leads

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

func TestUpdateLead(t *testing.T) {
	leadID := 123

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "PATCH" {
			t.Errorf("Ожидался метод PATCH, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/leads/%d", leadID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 123,
			"name": "Обновленный лид",
			"price": 25000,
			"responsible_user_id": 456,
			"status_id": 142,
			"pipeline_id": 777,
			"created_at": 1609459200,
			"updated_at": 1609632000
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Создаем данные для обновления
	leadToUpdate := &Lead{
		ID:               leadID,
		Name:             "Обновленный лид",
		Price:            25000,
		StatusID:         142,
		PipelineID:       777,
		ResponsibleUserID: 456,
	}

	// Вызываем тестируемый метод
	updatedLead, err := UpdateLead(apiClient, leadToUpdate)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при обновлении лида: %v", err)
	}

	if updatedLead.ID != leadID {
		t.Errorf("Ожидался ID лида %d, получен %d", leadID, updatedLead.ID)
	}

	if updatedLead.Name != "Обновленный лид" {
		t.Errorf("Ожидалось имя лида 'Обновленный лид', получено '%s'", updatedLead.Name)
	}

	if updatedLead.Price != 25000 {
		t.Errorf("Ожидалась цена лида 25000, получена %d", updatedLead.Price)
	}

	if updatedLead.StatusID != 142 {
		t.Errorf("Ожидался ID статуса 142, получен %d", updatedLead.StatusID)
	}

	if updatedLead.PipelineID != 777 {
		t.Errorf("Ожидался ID воронки 777, получен %d", updatedLead.PipelineID)
	}
}

// Тест обновления лида с ошибкой - не указан ID
func TestUpdateLeadWithoutID(t *testing.T) {
	// Создаем клиент API с любым URL, так как запрос не будет отправлен
	apiClient := client.NewClient("http://localhost", "test_api_key")

	// Создаем данные для обновления без ID
	leadToUpdate := &Lead{
		Name:  "Лид без ID",
		Price: 15000,
	}

	// Вызываем тестируемый метод
	_, err := UpdateLead(apiClient, leadToUpdate)

	// Проверяем результаты - должна быть ошибка
	if err == nil {
		t.Fatalf("Ожидалась ошибка при обновлении лида без ID, но ее не было")
	}
}

func TestDeleteLead(t *testing.T) {
	leadID := 123

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "DELETE" {
			t.Errorf("Ожидался метод DELETE, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/leads/%d", leadID)
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
	err := DeleteLead(apiClient, leadID)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при удалении лида: %v", err)
	}
}

func TestGetLeadWithOptions(t *testing.T) {
	leadID := 123

	// Тестируем разные комбинации опций with
	testCases := []struct {
		name         string
		withOptions  []WithOption
		expectedPath string
	}{
		{
			name:         "С контактами",
			withOptions:  []WithOption{WithContacts},
			expectedPath: "/api/v4/leads/123?with=contacts",
		},
		{
			name:         "С компаниями",
			withOptions:  []WithOption{WithCompanies},
			expectedPath: "/api/v4/leads/123?with=companies",
		},
		{
			name:         "С контактами и компаниями",
			withOptions:  []WithOption{WithContacts, WithCompanies},
			expectedPath: "/api/v4/leads/123?with=contacts%2Ccompanies",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Создаем тестовый сервер
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Проверяем метод запроса
				if r.Method != "GET" {
					t.Errorf("Ожидался метод GET, получен %s", r.Method)
				}

				// Проверяем путь запроса с параметрами
				fullPath := r.URL.Path
				if r.URL.RawQuery != "" {
					fullPath = fullPath + "?" + r.URL.RawQuery
				}

				// Сравниваем с ожидаемым путем
				if fullPath != tc.expectedPath {
					t.Errorf("Ожидался путь %s, получен %s", tc.expectedPath, fullPath)
				}

				// Формируем ответ в зависимости от опций
				response := `{
					"id": 123,
					"name": "Тестовый лид",
					"price": 10000,
					"responsible_user_id": 456,
					"created_at": 1609459200,
					"updated_at": 1609545600`

				// Добавляем связанные сущности, если они запрошены
				var embedded bool
				for _, opt := range tc.withOptions {
					if !embedded {
						response += `, "_embedded": {`
						embedded = true
					} else {
						response += `, `
					}

					switch opt {
					case WithContacts:
						response += `"contacts": [
							{
								"id": 456,
								"name": "Иван Иванов",
								"first_name": "Иван",
								"last_name": "Иванов",
								"responsible_user_id": 789
							}
						]`
					case WithCompanies:
						response += `"companies": [
							{
								"id": 789,
								"name": "ООО Тест",
								"responsible_user_id": 789
							}
						]`
					}
				}

				if embedded {
					response += `}`
				}

				response += `}`

				// Отправляем ответ
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(response))
			}))
			defer server.Close()

			// Создаем клиент API
			apiClient := client.NewClient(server.URL, "test_api_key")

			// Вызываем тестируемый метод с опциями
			lead, err := GetLead(apiClient, leadID, tc.withOptions...)

			// Проверяем результаты
			if err != nil {
				t.Fatalf("Ошибка при получении лида с опциями: %v", err)
			}

			if lead.ID != leadID {
				t.Errorf("Ожидался ID лида %d, получен %d", leadID, lead.ID)
			}

			// Проверяем наличие связанных сущностей в ответе
			for _, opt := range tc.withOptions {
				switch opt {
				case WithContacts:
					if lead.Embedded == nil || len(lead.Embedded.Contacts) == 0 {
						t.Errorf("Ожидались связанные контакты, но их нет в ответе")
					} else if lead.Embedded.Contacts[0].ID != 456 {
						t.Errorf("Ожидался ID контакта 456, получен %d", lead.Embedded.Contacts[0].ID)
					}
				case WithCompanies:
					if lead.Embedded == nil || len(lead.Embedded.Companies) == 0 {
						t.Errorf("Ожидались связанные компании, но их нет в ответе")
					} else if lead.Embedded.Companies[0].ID != 789 {
						t.Errorf("Ожидался ID компании 789, получен %d", lead.Embedded.Companies[0].ID)
					}
				}
			}
		})
	}
}
