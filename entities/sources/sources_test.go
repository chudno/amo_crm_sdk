package sources

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

// TestGetSources проверяет функцию получения списка источников сделок
func TestGetSources(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Создаем мок-клиент с успешным ответом
		mockClient := createGetSourcesSuccessMockClient()

		// Вызываем тестируемую функцию
		sources, err := GetSourcesWithRequester(mockClient, 1, 50)

		// Проверяем результаты
		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}

		// Проверяем полученный список источников
		verifySourcesList(t, sources, 2)
	})

	t.Run("Empty", func(t *testing.T) {
		// Создаем мок-клиент с пустым списком источников
		mockClient := createGetSourcesEmptyMockClient()

		// Вызываем тестируемую функцию
		sources, err := GetSourcesWithRequester(mockClient, 1, 50)

		// Проверяем результаты
		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}

		// Проверяем пустой список
		verifySourcesList(t, sources, 0)
	})

	t.Run("ServerError", func(t *testing.T) {
		// Создаем мок-клиент с ошибкой сервера
		mockClient := createGetSourcesErrorMockClient()

		// Вызываем тестируемую функцию
		_, err := GetSourcesWithRequester(mockClient, 1, 50)

		// Проверяем результаты
		if err == nil {
			t.Error("Ожидалась ошибка, но её нет")
		}
	})

	t.Run("WithFilter", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := createGetSourcesWithFilterMockClient()

		// Создаем фильтр
		filter := map[string]string{
			"filter[type]": "calls",
		}

		// Вызываем тестируемую функцию с фильтром
		sources, err := GetSourcesWithRequester(mockClient, 1, 50, WithFilter(filter))

		// Проверяем результаты
		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}

		// Проверяем полученный список источников
		verifySourcesList(t, sources, 1)

		// Проверяем наличие фильтра в URL запроса
		verifyFilterInRequest(t, mockClient, "filter%5Btype%5D=calls")
	})
}

// TestGetSource проверяет функцию получения конкретного источника
func TestGetSource(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusOK,
			Body: `{
				"id": 1001,
				"name": "Входящие звонки",
				"type": "calls",
				"default": true,
				"created_at": 1609459200,
				"updated_at": 1609459300
			}`,
		})

		// Вызываем тестируемую функцию
		source, err := GetSourceWithRequester(mockClient, 1001)

		// Проверяем результаты
		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}

		if source == nil {
			t.Fatal("Ожидался объект источника, но получен nil")
		}

		if source.ID != 1001 || source.Name != "Входящие звонки" || source.Type != "calls" {
			t.Errorf("Неверные данные в источнике: %+v", source)
		}

		// Проверяем, что запрос был отправлен
		if mockClient.LastRequest == nil {
			t.Fatal("Запрос не был выполнен")
		}

		// Проверяем URL запроса
		expectedURLPart := "/api/v4/sources/1001"
		if !strings.Contains(mockClient.LastRequest.URL, expectedURLPart) {
			t.Errorf("URL запроса не содержит ожидаемой части: %s", mockClient.LastRequest.URL)
		}
	})

	t.Run("Error", func(t *testing.T) {
		// Создаем мок-клиент с ошибкой
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusNotFound,
			Body:       `{"error": "Not Found"}`,
		})

		// Вызываем тестируемую функцию
		_, err := GetSourceWithRequester(mockClient, 9999)

		// Проверяем результаты
		if err == nil {
			t.Error("Ожидалась ошибка, но её нет")
		}
	})
}

// TestCreateSource проверяет функцию создания источника
func TestCreateSource(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusOK,
			Body: `{
				"id": 1001,
				"name": "Новый источник",
				"type": "other",
				"default": false,
				"created_at": 1609459200,
				"updated_at": 1609459300
			}`,
		})

		// Создаем новый источник
		newSource := &Source{
			Name: "Новый источник",
			Type: "other",
		}

		// Вызываем тестируемую функцию
		createdSource, err := CreateSourceWithRequester(mockClient, newSource)

		// Проверяем результаты
		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}

		if createdSource == nil {
			t.Fatal("Ожидался объект источника, но получен nil")
		}

		if createdSource.ID != 1001 || createdSource.Name != "Новый источник" || createdSource.Type != "other" {
			t.Errorf("Неверные данные в созданном источнике: %+v", createdSource)
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
		expectedURLPart := "/api/v4/sources"
		if !strings.Contains(mockClient.LastRequest.URL, expectedURLPart) {
			t.Errorf("URL запроса не содержит ожидаемой части: %s", mockClient.LastRequest.URL)
		}

		// Проверяем тело запроса
		var requestBody map[string]interface{}
		err = json.Unmarshal([]byte(mockClient.LastRequest.Body), &requestBody)
		if err != nil {
			t.Fatalf("Ошибка при разборе тела запроса: %v", err)
		}

		if requestBody["name"] != "Новый источник" || requestBody["type"] != "other" {
			t.Errorf("Неверные данные в теле запроса: %+v", requestBody)
		}
	})

	t.Run("Error", func(t *testing.T) {
		// Создаем мок-клиент с ошибкой
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Bad Request"}`,
		})

		// Создаем источник с невалидными данными
		invalidSource := &Source{
			Name: "", // Пустое имя
		}

		// Вызываем тестируемую функцию
		_, err := CreateSourceWithRequester(mockClient, invalidSource)

		// Проверяем результаты
		if err == nil {
			t.Error("Ожидалась ошибка, но её нет")
		}
	})
}

// TestUpdateSource проверяет функцию обновления источника
func TestUpdateSource(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusOK,
			Body: `{
				"id": 1001,
				"name": "Обновленный источник",
				"type": "calls",
				"default": true,
				"created_at": 1609459200,
				"updated_at": 1609459400
			}`,
		})

		// Создаем источник для обновления
		sourceToUpdate := &Source{
			ID:   1001,
			Name: "Обновленный источник",
		}

		// Вызываем тестируемую функцию
		updatedSource, err := UpdateSourceWithRequester(mockClient, sourceToUpdate)

		// Проверяем результаты
		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}

		if updatedSource == nil {
			t.Fatal("Ожидался объект источника, но получен nil")
		}

		if updatedSource.ID != 1001 || updatedSource.Name != "Обновленный источник" {
			t.Errorf("Неверные данные в обновленном источнике: %+v", updatedSource)
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
		expectedURLPart := "/api/v4/sources/1001"
		if !strings.Contains(mockClient.LastRequest.URL, expectedURLPart) {
			t.Errorf("URL запроса не содержит ожидаемой части: %s", mockClient.LastRequest.URL)
		}

		// Проверяем тело запроса
		var requestBody map[string]interface{}
		err = json.Unmarshal([]byte(mockClient.LastRequest.Body), &requestBody)
		if err != nil {
			t.Fatalf("Ошибка при разборе тела запроса: %v", err)
		}

		if requestBody["name"] != "Обновленный источник" {
			t.Errorf("Неверные данные в теле запроса: %+v", requestBody)
		}
	})

	t.Run("Error", func(t *testing.T) {
		// Создаем мок-клиент с ошибкой
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusNotFound,
			Body:       `{"error": "Not Found"}`,
		})

		// Создаем источник с несуществующим ID
		nonExistentSource := &Source{
			ID:   9999,
			Name: "Несуществующий источник",
		}

		// Вызываем тестируемую функцию
		_, err := UpdateSourceWithRequester(mockClient, nonExistentSource)

		// Проверяем результаты
		if err == nil {
			t.Error("Ожидалась ошибка, но её нет")
		}
	})
}

// TestDeleteSource проверяет функцию удаления источника
func TestDeleteSource(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusOK,
			Body:       `{}`,
		})

		// Вызываем тестируемую функцию
		err := DeleteSourceWithRequester(mockClient, 1001)

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
		expectedURLPart := "/api/v4/sources/1001"
		if !strings.Contains(mockClient.LastRequest.URL, expectedURLPart) {
			t.Errorf("URL запроса не содержит ожидаемой части: %s", mockClient.LastRequest.URL)
		}
	})

	t.Run("Error", func(t *testing.T) {
		// Создаем мок-клиент с ошибкой
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusNotFound,
			Body:       `{"error": "Not Found"}`,
		})

		// Вызываем тестируемую функцию
		err := DeleteSourceWithRequester(mockClient, 9999)

		// Проверяем результаты
		if err == nil {
			t.Error("Ожидалась ошибка, но её нет")
		}
	})
}

// TestSetSourceDefault проверяет функцию установки источника по умолчанию
func TestSetSourceDefault(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusOK,
			Body: `{
				"id": 1001,
				"name": "Входящие звонки",
				"type": "calls",
				"default": true,
				"created_at": 1609459200,
				"updated_at": 1609459300
			}`,
		})

		// Вызываем тестируемую функцию
		defaultSource, err := SetSourceDefaultWithRequester(mockClient, 1001)

		// Проверяем результаты
		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}

		if defaultSource == nil {
			t.Fatal("Ожидался объект источника, но получен nil")
		}

		if !defaultSource.Default {
			t.Error("Источник не был установлен как источник по умолчанию")
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
		expectedURLPart := "/api/v4/sources/1001/default"
		if !strings.Contains(mockClient.LastRequest.URL, expectedURLPart) {
			t.Errorf("URL запроса не содержит ожидаемой части: %s", mockClient.LastRequest.URL)
		}
	})

	t.Run("Error", func(t *testing.T) {
		// Создаем мок-клиент с ошибкой
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusNotFound,
			Body:       `{"error": "Not Found"}`,
		})

		// Вызываем тестируемую функцию
		_, err := SetSourceDefaultWithRequester(mockClient, 9999)

		// Проверяем результаты
		if err == nil {
			t.Error("Ожидалась ошибка, но её нет")
		}
	})
}

// TestGetSourceServices проверяет функцию получения списка сервисов для источников
func TestGetSourceServices(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusOK,
			Body: `{
				"_embedded": {
					"services": [
						{
							"id": 2001,
							"name": "Телефония",
							"code": "telephony"
						},
						{
							"id": 2002,
							"name": "Email",
							"code": "email"
						}
					]
				}
			}`,
		})

		// Вызываем тестируемую функцию
		services, err := GetSourceServicesWithRequester(mockClient)

		// Проверяем результаты
		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}

		if len(services) != 2 {
			t.Errorf("Ожидалось 2 сервиса, получено %d", len(services))
		}

		if services[0].ID != 2001 || services[0].Name != "Телефония" {
			t.Errorf("Неверные данные в первом сервисе: %+v", services[0])
		}

		if services[1].ID != 2002 || services[1].Name != "Email" {
			t.Errorf("Неверные данные во втором сервисе: %+v", services[1])
		}

		// Проверяем, что запрос был отправлен
		if mockClient.LastRequest == nil {
			t.Fatal("Запрос не был выполнен")
		}

		// Проверяем URL запроса
		expectedURLPart := "/api/v4/sources/services"
		if !strings.Contains(mockClient.LastRequest.URL, expectedURLPart) {
			t.Errorf("URL запроса не содержит ожидаемой части: %s", mockClient.LastRequest.URL)
		}
	})

	t.Run("Error", func(t *testing.T) {
		// Создаем мок-клиент с ошибкой
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error": "Internal Server Error"}`,
		})

		// Вызываем тестируемую функцию
		_, err := GetSourceServicesWithRequester(mockClient)

		// Проверяем результаты
		if err == nil {
			t.Error("Ожидалась ошибка, но её нет")
		}
	})
}

// TestLinkSourceToPipeline проверяет функцию связывания источника с воронкой
func TestLinkSourceToPipeline(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusOK,
			Body: `{
				"id": 1001,
				"name": "Входящие звонки",
				"type": "calls",
				"pipeline": {
					"id": 2001
				}
			}`,
		})

		// Вызываем тестируемую функцию
		linkedSource, err := LinkSourceToPipelineWithRequester(mockClient, 1001, 2001)

		// Проверяем результаты
		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}

		if linkedSource == nil {
			t.Fatal("Ожидался объект источника, но получен nil")
		}

		// Проверяем, что источник связан с воронкой
		if linkedSource.Pipeline == nil || linkedSource.Pipeline.ID != 2001 {
			t.Errorf("Источник не связан с воронкой 2001: %+v", linkedSource.Pipeline)
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
		expectedURLPart := "/api/v4/sources/1001/pipeline"
		if !strings.Contains(mockClient.LastRequest.URL, expectedURLPart) {
			t.Errorf("URL запроса не содержит ожидаемой части: %s", mockClient.LastRequest.URL)
		}

		// Проверяем тело запроса
		expectedBodyPart := `"pipeline_id":2001`
		if !strings.Contains(mockClient.LastRequest.Body, expectedBodyPart) {
			t.Errorf("Тело запроса не содержит ожидаемой части: %s", mockClient.LastRequest.Body)
		}
	})

	t.Run("Error", func(t *testing.T) {
		// Создаем мок-клиент с ошибкой
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Bad Request"}`,
		})

		// Вызываем тестируемую функцию
		_, err := LinkSourceToPipelineWithRequester(mockClient, 1001, 9999)

		// Проверяем результаты
		if err == nil {
			t.Error("Ожидалась ошибка, но её нет")
		}
	})
}

// TestUnlinkSourceFromPipeline проверяет функцию удаления связи источника с воронкой
func TestUnlinkSourceFromPipeline(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusOK,
			Body: `{
				"id": 1001,
				"name": "Входящие звонки",
				"type": "calls"
			}`,
		})

		// Вызываем тестируемую функцию
		unlinkedSource, err := UnlinkSourceFromPipelineWithRequester(mockClient, 1001, 2001)

		// Проверяем результаты
		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}

		if unlinkedSource == nil {
			t.Fatal("Ожидался объект источника, но получен nil")
		}

		// Проверяем, что источник не связан с воронкой
		if unlinkedSource.Pipeline != nil {
			t.Errorf("Источник все еще связан с воронкой: %+v", unlinkedSource.Pipeline)
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
		expectedURLPart := "/api/v4/sources/1001/pipeline/2001"
		if !strings.Contains(mockClient.LastRequest.URL, expectedURLPart) {
			t.Errorf("URL запроса не содержит ожидаемой части: %s", mockClient.LastRequest.URL)
		}
	})

	t.Run("Error", func(t *testing.T) {
		// Создаем мок-клиент с ошибкой
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusNotFound,
			Body:       `{"error": "Not Found"}`,
		})

		// Вызываем тестируемую функцию
		_, err := UnlinkSourceFromPipelineWithRequester(mockClient, 1001, 9999)

		// Проверяем результаты
		if err == nil {
			t.Error("Ожидалась ошибка, но её нет")
		}
	})
}
