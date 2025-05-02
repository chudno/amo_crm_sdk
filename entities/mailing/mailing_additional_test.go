package mailing

import (
	"net/http"
	"strings"
	"testing"
)

// TestGetMailingStats проверяет функцию получения статистики рассылки
func TestGetMailingStats(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Создаем мок-клиент с успешным ответом
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusOK,
			Body: `{
				"total_recipients": 100,
				"delivered": 95,
				"opened": 75,
				"clicked": 50,
				"bounced": 5,
				"unsubscribed": 2,
				"complaints": 0
			}`,
		})

		// Вызываем тестируемую функцию
		stats, err := GetMailingStatsWithRequester(mockClient, 1001)

		// Проверяем результаты
		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}

		if stats == nil {
			t.Fatal("Ожидался объект статистики, но получен nil")
		}

		// Проверяем данные статистики
		if stats.TotalRecipients != 100 {
			t.Errorf("Ожидалось TotalRecipients = 100, получено %d", stats.TotalRecipients)
		}
		if stats.Delivered != 95 {
			t.Errorf("Ожидалось Delivered = 95, получено %d", stats.Delivered)
		}
		if stats.Opened != 75 {
			t.Errorf("Ожидалось Opened = 75, получено %d", stats.Opened)
		}
		if stats.Clicked != 50 {
			t.Errorf("Ожидалось Clicked = 50, получено %d", stats.Clicked)
		}
		if stats.Bounced != 5 {
			t.Errorf("Ожидалось Bounced = 5, получено %d", stats.Bounced)
		}
		if stats.Unsubscribed != 2 {
			t.Errorf("Ожидалось Unsubscribed = 2, получено %d", stats.Unsubscribed)
		}
		if stats.Complaints != 0 {
			t.Errorf("Ожидалось Complaints = 0, получено %d", stats.Complaints)
		}

		// Проверяем, что запрос был отправлен
		if mockClient.LastRequest == nil {
			t.Fatal("Запрос не был выполнен")
		}

		// Проверяем метод запроса
		if mockClient.LastRequest.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", mockClient.LastRequest.Method)
		}

		// Проверяем URL запроса
		expectedURLPart := "/api/v4/mailings/1001/stats"
		if !strings.Contains(mockClient.LastRequest.URL, expectedURLPart) {
			t.Errorf("URL запроса не содержит ожидаемой части: %s", mockClient.LastRequest.URL)
		}
	})

	t.Run("Error", func(t *testing.T) {
		// Создаем мок-клиент с ошибкой
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusNotFound,
			Body:       `{"error": "Mailing not found"}`,
		})

		// Вызываем тестируемую функцию
		stats, err := GetMailingStatsWithRequester(mockClient, 9999)

		// Проверяем результаты
		if err == nil {
			t.Error("Ожидалась ошибка, но её нет")
		}

		if stats != nil {
			t.Errorf("Ожидался nil, но получен объект статистики: %+v", stats)
		}
	})
}

// TestAddMailingRecipients проверяет функцию добавления получателей рассылки
func TestAddMailingRecipients(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Создаем мок-клиент с успешным ответом
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusOK,
			Body:       `{"success": true}`,
		})

		// Список ID контактов для добавления
		contactIDs := []int{1001, 1002, 1003}

		// Вызываем тестируемую функцию
		err := AddMailingRecipientsWithRequester(mockClient, 1001, contactIDs)

		// Проверяем результаты
		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}

		// Проверяем, что запрос был отправлен
		if mockClient.LastRequest == nil {
			t.Fatal("Запрос не был выполнен")
		}

		// Проверяем метод запроса
		if mockClient.LastRequest.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", mockClient.LastRequest.Method)
		}

		// Проверяем URL запроса
		expectedURLPart := "/api/v4/mailings/1001/recipients"
		if !strings.Contains(mockClient.LastRequest.URL, expectedURLPart) {
			t.Errorf("URL запроса не содержит ожидаемой части: %s", mockClient.LastRequest.URL)
		}

		// Проверяем тело запроса
		expectedBodyParts := []string{`"contact_ids"`, `1001`, `1002`, `1003`}
		for _, part := range expectedBodyParts {
			if !strings.Contains(mockClient.LastRequest.Body, part) {
				t.Errorf("Тело запроса не содержит ожидаемой части '%s': %s", part, mockClient.LastRequest.Body)
			}
		}
	})

	t.Run("Error", func(t *testing.T) {
		// Создаем мок-клиент с ошибкой
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Invalid contact IDs"}`,
		})

		// Список ID контактов для добавления
		contactIDs := []int{-1, -2}

		// Вызываем тестируемую функцию
		err := AddMailingRecipientsWithRequester(mockClient, 1001, contactIDs)

		// Проверяем результаты
		if err == nil {
			t.Error("Ожидалась ошибка, но её нет")
		}
	})
}

// TestRemoveMailingRecipients проверяет функцию удаления получателей рассылки
func TestRemoveMailingRecipients(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Создаем мок-клиент с успешным ответом
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusOK,
			Body:       `{"success": true}`,
		})

		// Список ID контактов для удаления
		contactIDs := []int{1001, 1002}

		// Вызываем тестируемую функцию
		err := RemoveMailingRecipientsWithRequester(mockClient, 1001, contactIDs)

		// Проверяем результаты
		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}

		// Проверяем, что запрос был отправлен
		if mockClient.LastRequest == nil {
			t.Fatal("Запрос не был выполнен")
		}

		// Проверяем метод запроса
		if mockClient.LastRequest.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", mockClient.LastRequest.Method)
		}

		// Проверяем URL запроса
		expectedURLPart := "/api/v4/mailings/1001/recipients/delete"
		if !strings.Contains(mockClient.LastRequest.URL, expectedURLPart) {
			t.Errorf("URL запроса не содержит ожидаемой части: %s", mockClient.LastRequest.URL)
		}

		// Проверяем тело запроса
		expectedBodyParts := []string{`"contact_ids"`, `1001`, `1002`}
		for _, part := range expectedBodyParts {
			if !strings.Contains(mockClient.LastRequest.Body, part) {
				t.Errorf("Тело запроса не содержит ожидаемой части '%s': %s", part, mockClient.LastRequest.Body)
			}
		}
	})

	t.Run("Error", func(t *testing.T) {
		// Создаем мок-клиент с ошибкой
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Invalid contact IDs"}`,
		})

		// Список ID контактов для удаления
		contactIDs := []int{-1, -2}

		// Вызываем тестируемую функцию
		err := RemoveMailingRecipientsWithRequester(mockClient, 1001, contactIDs)

		// Проверяем результаты
		if err == nil {
			t.Error("Ожидалась ошибка, но её нет")
		}
	})
}

// TestGetMailingTemplates проверяет функцию получения списка шаблонов рассылок
func TestGetMailingTemplates(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Создаем мок-клиент с успешным ответом
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusOK,
			Body: `{
				"_embedded": {
					"templates": [
						{
							"id": 101,
							"name": "Шаблон 1",
							"content": "Содержимое шаблона 1",
							"html": "<p>Содержимое шаблона 1</p>",
							"type": "email"
						},
						{
							"id": 102,
							"name": "Шаблон 2",
							"content": "Содержимое шаблона 2",
							"html": "<p>Содержимое шаблона 2</p>",
							"type": "email"
						}
					]
				}
			}`,
		})

		// Вызываем тестируемую функцию
		templates, err := GetMailingTemplatesWithRequester(mockClient, 1, 50)

		// Проверяем результаты
		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}

		// Проверяем количество шаблонов
		expectedCount := 2
		if len(templates) != expectedCount {
			t.Errorf("Ожидалось %d шаблонов, получено %d", expectedCount, len(templates))
		}

		// Проверяем данные первого шаблона
		if templates[0].ID != 101 || templates[0].Name != "Шаблон 1" || templates[0].Type != "email" {
			t.Errorf("Неверные данные в первом шаблоне: %+v", templates[0])
		}

		// Проверяем данные второго шаблона
		if templates[1].ID != 102 || templates[1].Name != "Шаблон 2" || templates[1].Type != "email" {
			t.Errorf("Неверные данные во втором шаблоне: %+v", templates[1])
		}

		// Проверяем, что запрос был отправлен
		if mockClient.LastRequest == nil {
			t.Fatal("Запрос не был выполнен")
		}

		// Проверяем метод запроса
		if mockClient.LastRequest.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", mockClient.LastRequest.Method)
		}

		// Проверяем URL запроса
		expectedURLPart := "/api/v4/mailing_templates"
		if !strings.Contains(mockClient.LastRequest.URL, expectedURLPart) {
			t.Errorf("URL запроса не содержит ожидаемой части: %s", mockClient.LastRequest.URL)
		}

		// Проверяем параметры запроса
		expectedParams := []string{"page=1", "limit=50"}
		for _, param := range expectedParams {
			if !strings.Contains(mockClient.LastRequest.URL, param) {
				t.Errorf("URL запроса не содержит ожидаемого параметра '%s': %s", param, mockClient.LastRequest.URL)
			}
		}
	})

	t.Run("Empty", func(t *testing.T) {
		// Создаем мок-клиент с пустым ответом
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusOK,
			Body: `{
				"_embedded": {
					"templates": []
				}
			}`,
		})

		// Вызываем тестируемую функцию
		templates, err := GetMailingTemplatesWithRequester(mockClient, 1, 50)

		// Проверяем результаты
		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}

		// Проверяем, что получен пустой список
		if len(templates) != 0 {
			t.Errorf("Ожидался пустой список шаблонов, получено %d элементов", len(templates))
		}
	})

	t.Run("Error", func(t *testing.T) {
		// Создаем мок-клиент с ошибкой
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error": "Server error"}`,
		})

		// Вызываем тестируемую функцию
		templates, err := GetMailingTemplatesWithRequester(mockClient, 1, 50)

		// Проверяем результаты
		if err == nil {
			t.Error("Ожидалась ошибка, но её нет")
		}

		if templates != nil {
			t.Errorf("Ожидался nil, но получен список шаблонов: %+v", templates)
		}
	})
}

// TestGetMailingTemplate проверяет функцию получения конкретного шаблона рассылки
func TestGetMailingTemplate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Создаем мок-клиент с успешным ответом
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusOK,
			Body: `{
				"id": 101,
				"name": "Шаблон рассылки",
				"content": "Содержимое шаблона",
				"html": "<p>Содержимое шаблона</p>",
				"type": "email"
			}`,
		})

		// Вызываем тестируемую функцию
		template, err := GetMailingTemplateWithRequester(mockClient, 101)

		// Проверяем результаты
		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}

		if template == nil {
			t.Fatal("Ожидался объект шаблона, но получен nil")
		}

		// Проверяем данные шаблона
		if template.ID != 101 {
			t.Errorf("Ожидался шаблон с ID = 101, получен ID = %d", template.ID)
		}
		if template.Name != "Шаблон рассылки" {
			t.Errorf("Ожидалось имя шаблона 'Шаблон рассылки', получено '%s'", template.Name)
		}
		if template.Content != "Содержимое шаблона" {
			t.Errorf("Ожидалось содержимое 'Содержимое шаблона', получено '%s'", template.Content)
		}
		if template.HTML != "<p>Содержимое шаблона</p>" {
			t.Errorf("Ожидался HTML '<p>Содержимое шаблона</p>', получен '%s'", template.HTML)
		}
		if template.Type != "email" {
			t.Errorf("Ожидался тип 'email', получен '%s'", template.Type)
		}

		// Проверяем, что запрос был отправлен
		if mockClient.LastRequest == nil {
			t.Fatal("Запрос не был выполнен")
		}

		// Проверяем метод запроса
		if mockClient.LastRequest.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", mockClient.LastRequest.Method)
		}

		// Проверяем URL запроса
		expectedURLPart := "/api/v4/mailing_templates/101"
		if !strings.Contains(mockClient.LastRequest.URL, expectedURLPart) {
			t.Errorf("URL запроса не содержит ожидаемой части: %s", mockClient.LastRequest.URL)
		}
	})

	t.Run("Error", func(t *testing.T) {
		// Создаем мок-клиент с ошибкой
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusNotFound,
			Body:       `{"error": "Template not found"}`,
		})

		// Вызываем тестируемую функцию
		template, err := GetMailingTemplateWithRequester(mockClient, 9999)

		// Проверяем результаты
		if err == nil {
			t.Error("Ожидалась ошибка, но её нет")
		}

		if template != nil {
			t.Errorf("Ожидался nil, но получен объект шаблона: %+v", template)
		}
	})
}
