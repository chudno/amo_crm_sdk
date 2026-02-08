package sources

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

// TestGetSources проверяет функцию получения списка источников сделок
func TestGetSources(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := createGetSourcesSuccessMockClient()

		sources, err := ListWithRequester(context.Background(), mockClient, 1, 50)

		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}

		verifySourcesList(t, sources, 2)
	})

	t.Run("Empty", func(t *testing.T) {
		mockClient := createGetSourcesEmptyMockClient()

		sources, err := ListWithRequester(context.Background(), mockClient, 1, 50)

		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}

		verifySourcesList(t, sources, 0)
	})

	t.Run("ServerError", func(t *testing.T) {
		mockClient := createGetSourcesErrorMockClient()

		_, err := ListWithRequester(context.Background(), mockClient, 1, 50)

		if err == nil {
			t.Error("Ожидалась ошибка, но её нет")
		}
	})

	t.Run("WithFilter", func(t *testing.T) {
		mockClient := createGetSourcesWithFilterMockClient()

		filter := map[string]string{
			"filter[type]": "calls",
		}

		sources, err := ListWithRequester(context.Background(), mockClient, 1, 50, WithFilter(filter))

		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}

		verifySourcesList(t, sources, 1)

		verifyFilterInRequest(t, mockClient, "filter%5Btype%5D=calls")
	})
}

// TestGetSource проверяет функцию получения конкретного источника
func TestGetSource(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
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

		source, err := GetWithRequester(context.Background(), mockClient, 1001)

		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}

		if source == nil {
			t.Fatal("Ожидался объект источника, но получен nil")
		}

		if source.ID != 1001 || source.Name != "Входящие звонки" || source.Type != "calls" {
			t.Errorf("Неверные данные в источнике: %+v", source)
		}

		if mockClient.LastRequest == nil {
			t.Fatal("Запрос не был выполнен")
		}

		expectedURLPart := "/api/v4/sources/1001"
		if !strings.Contains(mockClient.LastRequest.URL, expectedURLPart) {
			t.Errorf("URL запроса не содержит ожидаемой части: %s", mockClient.LastRequest.URL)
		}
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusNotFound,
			Body:       `{"error": "Not Found"}`,
		})

		_, err := GetWithRequester(context.Background(), mockClient, 9999)

		if err == nil {
			t.Error("Ожидалась ошибка, но её нет")
		}
	})
}

// TestCreateSource проверяет функцию создания источника
func TestCreateSource(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
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

		newSource := &Source{
			Name: "Новый источник",
			Type: "other",
		}

		createdSource, err := CreateWithRequester(context.Background(), mockClient, newSource)

		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}

		if createdSource == nil {
			t.Fatal("Ожидался объект источника, но получен nil")
		}

		if createdSource.ID != 1001 || createdSource.Name != "Новый источник" || createdSource.Type != "other" {
			t.Errorf("Неверные данные в созданном источнике: %+v", createdSource)
		}

		if mockClient.LastRequest == nil {
			t.Fatal("Запрос не был выполнен")
		}

		if mockClient.LastRequest.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", mockClient.LastRequest.Method)
		}

		expectedURLPart := "/api/v4/sources"
		if !strings.Contains(mockClient.LastRequest.URL, expectedURLPart) {
			t.Errorf("URL запроса не содержит ожидаемой части: %s", mockClient.LastRequest.URL)
		}

		var requestBody map[string]any
		err = json.Unmarshal([]byte(mockClient.LastRequest.Body), &requestBody)
		if err != nil {
			t.Fatalf("Ошибка при разборе тела запроса: %v", err)
		}

		if requestBody["name"] != "Новый источник" || requestBody["type"] != "other" {
			t.Errorf("Неверные данные в теле запроса: %+v", requestBody)
		}
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Bad Request"}`,
		})

		invalidSource := &Source{
			Name: "", // Пустое имя
		}

		_, err := CreateWithRequester(context.Background(), mockClient, invalidSource)

		if err == nil {
			t.Error("Ожидалась ошибка, но её нет")
		}
	})
}

// TestUpdateSource проверяет функцию обновления источника
func TestUpdateSource(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
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

		sourceToUpdate := &Source{
			ID:   1001,
			Name: "Обновленный источник",
		}

		updatedSource, err := UpdateWithRequester(context.Background(), mockClient, sourceToUpdate)

		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}

		if updatedSource == nil {
			t.Fatal("Ожидался объект источника, но получен nil")
		}

		if updatedSource.ID != 1001 || updatedSource.Name != "Обновленный источник" {
			t.Errorf("Неверные данные в обновленном источнике: %+v", updatedSource)
		}

		if mockClient.LastRequest == nil {
			t.Fatal("Запрос не был выполнен")
		}

		if mockClient.LastRequest.Method != "PATCH" {
			t.Errorf("Ожидался метод PATCH, получен %s", mockClient.LastRequest.Method)
		}

		expectedURLPart := "/api/v4/sources/1001"
		if !strings.Contains(mockClient.LastRequest.URL, expectedURLPart) {
			t.Errorf("URL запроса не содержит ожидаемой части: %s", mockClient.LastRequest.URL)
		}

		var requestBody map[string]any
		err = json.Unmarshal([]byte(mockClient.LastRequest.Body), &requestBody)
		if err != nil {
			t.Fatalf("Ошибка при разборе тела запроса: %v", err)
		}

		if requestBody["name"] != "Обновленный источник" {
			t.Errorf("Неверные данные в теле запроса: %+v", requestBody)
		}
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusNotFound,
			Body:       `{"error": "Not Found"}`,
		})

		nonExistentSource := &Source{
			ID:   9999,
			Name: "Несуществующий источник",
		}

		_, err := UpdateWithRequester(context.Background(), mockClient, nonExistentSource)

		if err == nil {
			t.Error("Ожидалась ошибка, но её нет")
		}
	})
}

// TestDeleteSource проверяет функцию удаления источника
func TestDeleteSource(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusOK,
			Body:       `{}`,
		})

		err := DeleteWithRequester(context.Background(), mockClient, 1001)

		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}

		if mockClient.LastRequest == nil {
			t.Fatal("Запрос не был выполнен")
		}

		if mockClient.LastRequest.Method != "DELETE" {
			t.Errorf("Ожидался метод DELETE, получен %s", mockClient.LastRequest.Method)
		}

		expectedURLPart := "/api/v4/sources/1001"
		if !strings.Contains(mockClient.LastRequest.URL, expectedURLPart) {
			t.Errorf("URL запроса не содержит ожидаемой части: %s", mockClient.LastRequest.URL)
		}
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusNotFound,
			Body:       `{"error": "Not Found"}`,
		})

		err := DeleteWithRequester(context.Background(), mockClient, 9999)

		if err == nil {
			t.Error("Ожидалась ошибка, но её нет")
		}
	})
}

// TestSetSourceDefault проверяет функцию установки источника по умолчанию
func TestSetSourceDefault(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
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

		defaultSource, err := SetDefaultWithRequester(context.Background(), mockClient, 1001)

		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}

		if defaultSource == nil {
			t.Fatal("Ожидался объект источника, но получен nil")
		}

		if !defaultSource.Default {
			t.Error("Источник не был установлен как источник по умолчанию")
		}

		if mockClient.LastRequest == nil {
			t.Fatal("Запрос не был выполнен")
		}

		if mockClient.LastRequest.Method != "PATCH" {
			t.Errorf("Ожидался метод PATCH, получен %s", mockClient.LastRequest.Method)
		}

		expectedURLPart := "/api/v4/sources/1001/default"
		if !strings.Contains(mockClient.LastRequest.URL, expectedURLPart) {
			t.Errorf("URL запроса не содержит ожидаемой части: %s", mockClient.LastRequest.URL)
		}
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusNotFound,
			Body:       `{"error": "Not Found"}`,
		})

		_, err := SetDefaultWithRequester(context.Background(), mockClient, 9999)

		if err == nil {
			t.Error("Ожидалась ошибка, но её нет")
		}
	})
}

// TestGetSourceServices проверяет функцию получения списка сервисов для источников
func TestGetSourceServices(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
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

		services, err := ListServicesWithRequester(context.Background(), mockClient)

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

		if mockClient.LastRequest == nil {
			t.Fatal("Запрос не был выполнен")
		}

		expectedURLPart := "/api/v4/sources/services"
		if !strings.Contains(mockClient.LastRequest.URL, expectedURLPart) {
			t.Errorf("URL запроса не содержит ожидаемой части: %s", mockClient.LastRequest.URL)
		}
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error": "Internal Server Error"}`,
		})

		_, err := ListServicesWithRequester(context.Background(), mockClient)

		if err == nil {
			t.Error("Ожидалась ошибка, но её нет")
		}
	})
}

// TestLinkSourceToPipeline проверяет функцию связывания источника с воронкой
func TestLinkSourceToPipeline(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
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

		linkedSource, err := LinkToPipelineWithRequester(context.Background(), mockClient, 1001, 2001)

		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}

		if linkedSource == nil {
			t.Fatal("Ожидался объект источника, но получен nil")
		}

		if linkedSource.Pipeline == nil || linkedSource.Pipeline.ID != 2001 {
			t.Errorf("Источник не связан с воронкой 2001: %+v", linkedSource.Pipeline)
		}

		if mockClient.LastRequest == nil {
			t.Fatal("Запрос не был выполнен")
		}

		if mockClient.LastRequest.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", mockClient.LastRequest.Method)
		}

		expectedURLPart := "/api/v4/sources/1001/pipeline"
		if !strings.Contains(mockClient.LastRequest.URL, expectedURLPart) {
			t.Errorf("URL запроса не содержит ожидаемой части: %s", mockClient.LastRequest.URL)
		}

		expectedBodyPart := `"pipeline_id":2001`
		if !strings.Contains(mockClient.LastRequest.Body, expectedBodyPart) {
			t.Errorf("Тело запроса не содержит ожидаемой части: %s", mockClient.LastRequest.Body)
		}
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Bad Request"}`,
		})

		_, err := LinkToPipelineWithRequester(context.Background(), mockClient, 1001, 9999)

		if err == nil {
			t.Error("Ожидалась ошибка, но её нет")
		}
	})
}

// TestUnlinkSourceFromPipeline проверяет функцию удаления связи источника с воронкой
func TestUnlinkSourceFromPipeline(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusOK,
			Body: `{
				"id": 1001,
				"name": "Входящие звонки",
				"type": "calls"
			}`,
		})

		unlinkedSource, err := UnlinkFromPipelineWithRequester(context.Background(), mockClient, 1001, 2001)

		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}

		if unlinkedSource == nil {
			t.Fatal("Ожидался объект источника, но получен nil")
		}

		if unlinkedSource.Pipeline != nil {
			t.Errorf("Источник все еще связан с воронкой: %+v", unlinkedSource.Pipeline)
		}

		if mockClient.LastRequest == nil {
			t.Fatal("Запрос не был выполнен")
		}

		if mockClient.LastRequest.Method != "DELETE" {
			t.Errorf("Ожидался метод DELETE, получен %s", mockClient.LastRequest.Method)
		}

		expectedURLPart := "/api/v4/sources/1001/pipeline/2001"
		if !strings.Contains(mockClient.LastRequest.URL, expectedURLPart) {
			t.Errorf("URL запроса не содержит ожидаемой части: %s", mockClient.LastRequest.URL)
		}
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
			StatusCode: http.StatusNotFound,
			Body:       `{"error": "Not Found"}`,
		})

		_, err := UnlinkFromPipelineWithRequester(context.Background(), mockClient, 1001, 9999)

		if err == nil {
			t.Error("Ожидалась ошибка, но её нет")
		}
	})
}
