package mailing

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"
)

// createGetMailingsSuccessMockClient создает мок-клиент с успешным ответом для списка рассылок
func createGetMailingsSuccessMockClient() *AdvancedMockClient {
	return NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
		StatusCode: http.StatusOK,
		Body: `{
			"_embedded": {
				"mailings": [
					{
						"id": 1001,
						"name": "Тестовая рассылка 1",
						"status": "active",
						"subject": "Тема рассылки 1",
						"created_at": 1609459200,
						"updated_at": 1609459300
					},
					{
						"id": 1002,
						"name": "Тестовая рассылка 2",
						"status": "scheduled",
						"subject": "Тема рассылки 2",
						"created_at": 1609459400,
						"updated_at": 1609459500
					}
				]
			}
		}`,
	})
}

// createGetMailingsEmptyMockClient создает мок-клиент с пустым списком рассылок
func createGetMailingsEmptyMockClient() *AdvancedMockClient {
	return NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
		StatusCode: http.StatusOK,
		Body: `{
			"_embedded": {
				"mailings": []
			}
		}`,
	})
}

// createGetMailingsErrorMockClient создает мок-клиент с ошибкой сервера
func createGetMailingsErrorMockClient() *AdvancedMockClient {
	return NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       `{"error": "Internal Server Error"}`,
	})
}

// createGetMailingsWithFilterMockClient создает мок-клиент для запроса с фильтром
func createGetMailingsWithFilterMockClient() *AdvancedMockClient {
	return NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
		StatusCode: http.StatusOK,
		Body: `{
			"_embedded": {
				"mailings": [
					{
						"id": 1001,
						"name": "Тестовая рассылка 1",
						"status": "active",
						"subject": "Тема рассылки 1",
						"created_at": 1609459200,
						"updated_at": 1609459300
					}
				]
			}
		}`,
	})
}

// verifyMailingsList проверяет правильность полученного списка рассылок
func verifyMailingsList(t *testing.T, mailings []Mailing, expectedCount int) {
	if len(mailings) != expectedCount {
		t.Errorf("Ожидалось %d рассылок, получено %d", expectedCount, len(mailings))
	}

	if expectedCount > 0 {
		if mailings[0].ID != 1001 || mailings[0].Name != "Тестовая рассылка 1" || mailings[0].Status != "active" {
			t.Errorf("Неверные данные в первой рассылке: %+v", mailings[0])
		}
	}

	if expectedCount > 1 {
		if mailings[1].ID != 1002 || mailings[1].Name != "Тестовая рассылка 2" || mailings[1].Status != "scheduled" {
			t.Errorf("Неверные данные во второй рассылке: %+v", mailings[1])
		}
	}
}

// verifyFilterInRequest проверяет наличие параметра фильтра в URL запроса
func verifyFilterInRequest(t *testing.T, mockClient *AdvancedMockClient, expectedFilterPart string) {
	if mockClient.LastRequest == nil {
		t.Fatal("Запрос не был выполнен")
	}

	if mockClient.LastRequest.URL == "" {
		t.Fatal("URL запроса пустой")
	}

	if !strings.Contains(mockClient.LastRequest.URL, expectedFilterPart) {
		t.Errorf("Фильтр не был добавлен к URL запроса: %s", mockClient.LastRequest.URL)
	}
}

// TestGetMailings проверяет функцию получения списка рассылок
func TestGetMailings(t *testing.T) {
	// Используем пакет time для тестов с датами
	_ = time.Now()

	t.Run("Success", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := createGetMailingsSuccessMockClient()

		// Вызываем тестируемую функцию
		mailings, err := GetMailingsWithRequester(mockClient, 1, 50)

		// Проверяем результаты
		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}

		verifyMailingsList(t, mailings, 2)
	})

	t.Run("Empty", func(t *testing.T) {
		// Создаем мок-клиент с пустым списком рассылок
		mockClient := createGetMailingsEmptyMockClient()

		// Вызываем тестируемую функцию
		mailings, err := GetMailingsWithRequester(mockClient, 1, 50)

		// Проверяем результаты
		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}

		verifyMailingsList(t, mailings, 0)
	})

	t.Run("ServerError", func(t *testing.T) {
		// Создаем мок-клиент с ошибкой сервера
		mockClient := createGetMailingsErrorMockClient()

		// Вызываем тестируемую функцию
		_, err := GetMailingsWithRequester(mockClient, 1, 50)

		// Проверяем результаты
		if err == nil {
			t.Error("Ожидалась ошибка, но её нет")
		}
	})

	t.Run("WithFilter", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := createGetMailingsWithFilterMockClient()

		// Создаем фильтр
		filter := map[string]string{
			"filter[status]": "active",
		}

		// Вызываем тестируемую функцию с фильтром
		mailings, err := GetMailingsWithRequester(mockClient, 1, 50, WithFilter(filter))

		// Проверяем результаты
		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}

		verifyMailingsList(t, mailings, 1)

		// Проверка фильтра в запросе
		expectedFilterPart := "filter%5Bstatus%5D=active"
		verifyFilterInRequest(t, mockClient, expectedFilterPart)
	})
}

// TestGetMailing проверяет функцию получения конкретной рассылки

// TestGetMailing проверяет функцию получения конкретной рассылки
func TestGetMailing(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusOK,
			Body: `{
				"id": 1001,
				"name": "Тестовая рассылка 1",
				"status": "active",
				"subject": "Тема рассылки 1",
				"created_at": 1609459200,
				"updated_at": 1609459300,
				"segment_ids": [101, 102],
				"template": {
					"id": 201,
					"name": "Шаблон 1",
					"content": "Содержимое шаблона"
				},
				"stats": {
					"total_recipients": 1000,
					"delivered": 950,
					"opened": 800,
					"clicked": 500
				}
			}`,
		})

		// Вызываем тестируемую функцию
		mailingInfo, err := GetMailingWithRequester(mockClient, 1001)

		// Проверяем результаты
		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}

		if mailingInfo == nil {
			t.Fatal("Ожидался объект рассылки, но получен nil")
		}

		// Проверяем основные поля
		if mailingInfo.ID != 1001 || mailingInfo.Name != "Тестовая рассылка 1" || mailingInfo.Status != "active" {
			t.Errorf("Неверные данные рассылки: %+v", mailingInfo)
		}

		// Проверяем вложенные объекты
		if mailingInfo.Template == nil || mailingInfo.Template.ID != 201 {
			t.Errorf("Неверные данные шаблона: %+v", mailingInfo.Template)
		}

		if mailingInfo.Stats == nil || mailingInfo.Stats.TotalRecipients != 1000 {
			t.Errorf("Неверные данные статистики: %+v", mailingInfo.Stats)
		}

		// Проверяем slices
		if len(mailingInfo.SegmentIDs) != 2 || mailingInfo.SegmentIDs[0] != 101 || mailingInfo.SegmentIDs[1] != 102 {
			t.Errorf("Неверные ID сегментов: %+v", mailingInfo.SegmentIDs)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		// Создаем мок-клиент с ошибкой Not Found
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusNotFound,
			Body:       `{"error": "Mailing not found"}`,
		})

		// Вызываем тестируемую функцию
		_, err := GetMailingWithRequester(mockClient, 9999)

		// Проверяем результаты
		if err == nil {
			t.Error("Ожидалась ошибка, но её нет")
		}
	})
}

// TestCreateMailing проверяет функцию создания рассылки
func TestCreateMailing(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusCreated,
			Body: `{
				"id": 1001,
				"name": "Новая рассылка",
				"status": "draft",
				"subject": "Тема новой рассылки",
				"created_at": 1609459200,
				"updated_at": 1609459200
			}`,
		})

		// Создаем данные для рассылки
		mailingData := &Mailing{
			Name:      "Новая рассылка",
			Subject:   "Тема новой рассылки",
			Frequency: MailingFrequencyOnce,
		}

		// Вызываем тестируемую функцию
		createdMailing, err := CreateMailingWithRequester(mockClient, mailingData)

		// Проверяем результаты
		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}

		if createdMailing == nil {
			t.Fatal("Ожидался объект созданной рассылки, но получен nil")
		}

		// Проверяем основные поля
		if createdMailing.ID != 1001 || createdMailing.Name != "Новая рассылка" || createdMailing.Status != "draft" {
			t.Errorf("Неверные данные созданной рассылки: %+v", createdMailing)
		}

		// Проверяем, что запрос был отправлен
		if mockClient.LastRequest == nil {
			t.Fatal("Запрос не был выполнен")
		}

		// Проверяем метод запроса
		if mockClient.LastRequest.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", mockClient.LastRequest.Method)
		}

		// Проверяем тело запроса
		if mockClient.LastRequest.Body == "" {
			t.Fatal("Тело запроса пустое")
		}

		// Декодируем тело запроса
		var sentData Mailing
		if err := json.Unmarshal([]byte(mockClient.LastRequest.Body), &sentData); err != nil {
			t.Fatalf("Ошибка декодирования тела запроса: %v", err)
		}

		// Проверяем отправленные данные
		if sentData.Name != mailingData.Name || sentData.Subject != mailingData.Subject {
			t.Errorf("Отправленные данные не соответствуют ожидаемым: %+v", sentData)
		}
	})

	t.Run("Error", func(t *testing.T) {
		// Создаем мок-клиент с ошибкой
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Bad Request"}`,
		})

		// Создаем данные для рассылки
		mailingData := &Mailing{
			Name:      "Новая рассылка",
			Subject:   "Тема новой рассылки",
			Frequency: MailingFrequencyOnce,
		}

		// Вызываем тестируемую функцию
		_, err := CreateMailingWithRequester(mockClient, mailingData)

		// Проверяем результаты
		if err == nil {
			t.Error("Ожидалась ошибка, но её нет")
		}
	})
}

// TestUpdateMailing проверяет функцию обновления рассылки
func TestUpdateMailing(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusOK,
			Body: `{
				"id": 1001,
				"name": "Обновленная рассылка",
				"status": "draft",
				"subject": "Новая тема рассылки",
				"updated_at": 1609459300
			}`,
		})

		// Создаем данные для обновления
		mailingData := &Mailing{
			ID:      1001,
			Name:    "Обновленная рассылка",
			Subject: "Новая тема рассылки",
		}

		// Вызываем тестируемую функцию
		updatedMailing, err := UpdateMailingWithRequester(mockClient, mailingData)

		// Проверяем результаты
		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}

		if updatedMailing == nil {
			t.Fatal("Ожидался объект обновленной рассылки, но получен nil")
		}

		// Проверяем основные поля
		if updatedMailing.ID != 1001 || updatedMailing.Name != "Обновленная рассылка" || updatedMailing.Subject != "Новая тема рассылки" {
			t.Errorf("Неверные данные обновленной рассылки: %+v", updatedMailing)
		}

		// Проверяем, что запрос был отправлен
		if mockClient.LastRequest == nil {
			t.Fatal("Запрос не был выполнен")
		}

		// Проверяем метод запроса
		if mockClient.LastRequest.Method != "PATCH" {
			t.Errorf("Ожидался метод PATCH, получен %s", mockClient.LastRequest.Method)
		}

		// Проверяем URL запроса
		expectedURLPart := "/api/v4/mailings/1001"
		if !strings.Contains(mockClient.LastRequest.URL, expectedURLPart) {
			t.Errorf("URL запроса не содержит ожидаемой части: %s", mockClient.LastRequest.URL)
		}
	})

	t.Run("Error", func(t *testing.T) {
		// Создаем мок-клиент с ошибкой
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Bad Request"}`,
		})

		// Создаем данные для обновления
		mailingData := &Mailing{
			ID:      1001,
			Name:    "Обновленная рассылка",
			Subject: "Новая тема рассылки",
		}

		// Вызываем тестируемую функцию
		_, err := UpdateMailingWithRequester(mockClient, mailingData)

		// Проверяем результаты
		if err == nil {
			t.Error("Ожидалась ошибка, но её нет")
		}
	})

	t.Run("NoID", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusOK,
			Body:       `{}`,
		})

		// Создаем данные без ID
		mailingData := &Mailing{
			Name:    "Обновленная рассылка",
			Subject: "Новая тема рассылки",
		}

		// Вызываем тестируемую функцию
		_, err := UpdateMailingWithRequester(mockClient, mailingData)

		// Проверяем результаты
		if err == nil {
			t.Error("Ожидалась ошибка из-за отсутствия ID, но её нет")
		}
	})
}

// TestDeleteMailing проверяет функцию удаления рассылки
func TestDeleteMailing(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusNoContent,
			Body:       ``,
		})

		// Вызываем тестируемую функцию
		err := DeleteMailingWithRequester(mockClient, 1001)

		// Проверяем результаты
		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}

		// Проверяем, что запрос был отправлен
		if mockClient.LastRequest == nil {
			t.Fatal("Запрос не был выполнен")
		}

		// Проверяем метод запроса
		if mockClient.LastRequest.Method != "DELETE" {
			t.Errorf("Ожидался метод DELETE, получен %s", mockClient.LastRequest.Method)
		}

		// Проверяем URL запроса
		expectedURLPart := "/api/v4/mailings/1001"
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
		err := DeleteMailingWithRequester(mockClient, 9999)

		// Проверяем результаты
		if err == nil {
			t.Error("Ожидалась ошибка, но её нет")
		}
	})
}

// TestChangeMailingStatus проверяет функцию изменения статуса рассылки
func TestChangeMailingStatus(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusOK,
			Body: `{
				"id": 1001,
				"name": "Тестовая рассылка",
				"status": "paused",
				"subject": "Тема рассылки",
				"updated_at": 1609459300
			}`,
		})

		// Вызываем тестируемую функцию
		updatedMailing, err := ChangeMailingStatusWithRequester(mockClient, 1001, MailingStatusPaused)

		// Проверяем результаты
		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}

		if updatedMailing == nil {
			t.Fatal("Ожидался объект рассылки, но получен nil")
		}

		// Проверяем статус
		if updatedMailing.Status != MailingStatusPaused {
			t.Errorf("Ожидался статус %s, получен %s", MailingStatusPaused, updatedMailing.Status)
		}

		// Проверяем, что запрос был отправлен
		if mockClient.LastRequest == nil {
			t.Fatal("Запрос не был выполнен")
		}

		// Проверяем метод запроса
		if mockClient.LastRequest.Method != "PATCH" {
			t.Errorf("Ожидался метод PATCH, получен %s", mockClient.LastRequest.Method)
		}

		// Проверяем URL запроса
		expectedURLPart := "/api/v4/mailings/1001/status"
		if !strings.Contains(mockClient.LastRequest.URL, expectedURLPart) {
			t.Errorf("URL запроса не содержит ожидаемой части: %s", mockClient.LastRequest.URL)
		}

		// Проверяем тело запроса
		expectedBodyPart := `"status":"paused"`
		if !strings.Contains(mockClient.LastRequest.Body, expectedBodyPart) {
			t.Errorf("Тело запроса не содержит ожидаемой части: %s", mockClient.LastRequest.Body)
		}
	})

	t.Run("Error", func(t *testing.T) {
		// Создаем мок-клиент с ошибкой
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Invalid status"}`,
		})

		// Вызываем тестируемую функцию
		_, err := ChangeMailingStatusWithRequester(mockClient, 1001, "invalid_status")

		// Проверяем результаты
		if err == nil {
			t.Error("Ожидалась ошибка, но её нет")
		}
	})
}
