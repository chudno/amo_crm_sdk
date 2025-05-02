package pipelines

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

// TestGetPipelineError проверяет обработку ошибок при получении воронки
func TestGetPipelineError(t *testing.T) {
	t.Run("Ошибка HTTP", func(t *testing.T) {
		// Создаем клиент API с несуществующим доменом
		apiClient := client.NewClient("http://non-existent-domain.example", "test_api_key")

		// Вызываем тестируемый метод
		_, err := GetPipeline(apiClient, 999)

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка, но её не было")
		}
	})

	t.Run("Некорректный JSON", func(t *testing.T) {
		// Создаем тестовый сервер, который вернет некорректный JSON
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": 123, "name": "Тестовая воронка", "is_main": true, status`)) // Некорректный JSON
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемый метод
		_, err := GetPipeline(apiClient, 123)

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка из-за некорректного JSON, но её не было")
		}
	})
}

// TestListPipelinesError проверяет обработку ошибок при получении списка воронок
func TestListPipelinesError(t *testing.T) {
	t.Run("Ошибка HTTP", func(t *testing.T) {
		// Создаем клиент API с несуществующим доменом
		apiClient := client.NewClient("http://non-existent-domain.example", "test_api_key")

		// Вызываем тестируемый метод
		_, err := ListPipelines(apiClient)

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка, но её не было")
		}
	})

	t.Run("Пустой список воронок", func(t *testing.T) {
		// Создаем тестовый сервер, который вернет пустой список
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"_embedded": {"items": []}}`))
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемый метод
		pipelines, err := ListPipelines(apiClient)

		// Проверяем результаты
		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}

		if len(pipelines) != 0 {
			t.Errorf("Ожидался пустой список воронок, получено %d элементов", len(pipelines))
		}
	})

	t.Run("Некорректный JSON", func(t *testing.T) {
		// Создаем тестовый сервер, который вернет некорректный JSON
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"_embedded": {"items": [{"id": 123, "name": "Тестовая воронка"`)) // Некорректный JSON
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемый метод
		_, err := ListPipelines(apiClient)

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка из-за некорректного JSON, но её не было")
		}
	})
}

// TestUpdatePipelineError проверяет обработку ошибок при обновлении воронки
func TestUpdatePipelineError(t *testing.T) {
	t.Run("ID не указан", func(t *testing.T) {
		// Создаем клиент API с несуществующим доменом
		apiClient := client.NewClient("http://non-existent-domain.example", "test_api_key")

		// Создаем воронку без ID
		pipelineWithoutID := &Pipeline{
			Name:     "Воронка без ID",
			IsActive: true,
		}

		// Вызываем тестируемый метод
		_, err := UpdatePipeline(apiClient, pipelineWithoutID)

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка, но её не было")
		}
	})

	t.Run("Ошибка HTTP", func(t *testing.T) {
		// Создаем клиент API с несуществующим доменом
		apiClient := client.NewClient("http://non-existent-domain.example", "test_api_key")

		// Создаем воронку для обновления
		pipelineToUpdate := &Pipeline{
			ID:       999,
			Name:     "Несуществующая воронка",
			IsActive: true,
		}

		// Вызываем тестируемый метод
		_, err := UpdatePipeline(apiClient, pipelineToUpdate)

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка, но её не было")
		}
	})
}

// TestDeletePipelineError проверяет обработку ошибок при удалении воронки
func TestDeletePipelineError(t *testing.T) {
	t.Run("Ошибка HTTP", func(t *testing.T) {
		// Создаем тестовый сервер, который вернет ошибку
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"error": "Forbidden to delete main pipeline"}`))
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемый метод
		err := DeletePipeline(apiClient, 123)

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка, но её не было")
		}
	})

	t.Run("Некорректный код ответа", func(t *testing.T) {
		// Создаем тестовый сервер, который вернет ответ OK вместо NoContent
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK) // Должен быть StatusNoContent
			w.Write([]byte(`{}`))
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемый метод
		err := DeletePipeline(apiClient, 123)

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка из-за некорректного кода ответа, но её не было")
		}
	})
}

// TestGetStatusError проверяет обработку ошибок при получении статуса воронки
func TestGetStatusError(t *testing.T) {
	t.Run("Ошибка HTTP", func(t *testing.T) {
		// Создаем клиент API с несуществующим доменом
		apiClient := client.NewClient("http://non-existent-domain.example", "test_api_key")

		// Вызываем тестируемый метод
		_, err := GetStatus(apiClient, 123, 999)

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка, но её не было")
		}
	})
}

// TestCreatePipelineError проверяет обработку ошибок при создании воронки
func TestCreatePipelineError(t *testing.T) {
	t.Run("Ошибка HTTP", func(t *testing.T) {
		// Создаем клиент API с несуществующим доменом
		apiClient := client.NewClient("http://non-existent-domain.example", "test_api_key")

		// Создаем воронку для теста
		invalidPipeline := &Pipeline{
			Name: "Тестовая воронка",
		}

		// Вызываем тестируемый метод
		_, err := CreatePipeline(apiClient, invalidPipeline)

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка, но её не было")
		}
	})
}

// TestCreateStatusError проверяет обработку ошибок при создании статуса
func TestCreateStatusError(t *testing.T) {
	t.Run("Ошибка HTTP", func(t *testing.T) {
		// Создаем клиент API с несуществующим доменом
		apiClient := client.NewClient("http://non-existent-domain.example", "test_api_key")

		// Создаем статус для теста
		testStatus := &Status{
			Name:  "Тестовый статус",
			Color: "#FF0000",
		}

		// Вызываем тестируемый метод
		_, err := CreateStatus(apiClient, 123, testStatus)

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка, но её не было")
		}
	})
}
