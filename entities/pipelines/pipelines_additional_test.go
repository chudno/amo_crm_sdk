package pipelines

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

// TestGetPipelineError проверяет обработку ошибок при получении воронки
func TestGetPipelineError(t *testing.T) {
	t.Run("Ошибка HTTP", func(t *testing.T) {
		apiClient := client.NewClient("http://non-existent-domain.example", "test_api_key")

		_, err := Get(context.Background(), apiClient, 999)

		if err == nil {
			t.Error("Ожидалась ошибка, но её не было")
		}
	})

	t.Run("Некорректный JSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(`{"id": 123, "name": "Тестовая воронка", "is_main": true, status`)); err != nil { // Некорректный JSON
				t.Fatalf("Ошибка при записи ответа: %v", err)
			}
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		_, err := Get(context.Background(), apiClient, 123)

		if err == nil {
			t.Error("Ожидалась ошибка из-за некорректного JSON, но её не было")
		}
	})
}

// TestListPipelinesError проверяет обработку ошибок при получении списка воронок
func TestListPipelinesError(t *testing.T) {
	t.Run("Ошибка HTTP", func(t *testing.T) {
		apiClient := client.NewClient("http://non-existent-domain.example", "test_api_key")

		_, err := List(context.Background(), apiClient)

		if err == nil {
			t.Error("Ожидалась ошибка, но её не было")
		}
	})

	t.Run("Пустой список воронок", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(`{"_embedded": {"pipelines": []}}`)); err != nil {
				t.Fatalf("Ошибка при записи ответа: %v", err)
			}
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		pipelines, err := List(context.Background(), apiClient)

		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получена: %v", err)
		}

		if len(pipelines) != 0 {
			t.Errorf("Ожидался пустой список воронок, получено %d элементов", len(pipelines))
		}
	})

	t.Run("Некорректный JSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(`{"_embedded": {"pipelines": [{"id": 123, "name": "Тестовая воронка"`)); err != nil { // Некорректный JSON
				t.Fatalf("Ошибка при записи ответа: %v", err)
			}
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		_, err := List(context.Background(), apiClient)

		if err == nil {
			t.Error("Ожидалась ошибка из-за некорректного JSON, но её не было")
		}
	})
}

// TestUpdatePipelineError проверяет обработку ошибок при обновлении воронки
func TestUpdatePipelineError(t *testing.T) {
	t.Run("ID не указан", func(t *testing.T) {
		apiClient := client.NewClient("http://non-existent-domain.example", "test_api_key")

		pipelineWithoutID := &Pipeline{
			Name:     "Воронка без ID",
			IsActive: true,
		}

		_, err := Update(context.Background(), apiClient, pipelineWithoutID)

		if err == nil {
			t.Error("Ожидалась ошибка, но её не было")
		}
	})

	t.Run("Ошибка HTTP", func(t *testing.T) {
		apiClient := client.NewClient("http://non-existent-domain.example", "test_api_key")

		pipelineToUpdate := &Pipeline{
			ID:       999,
			Name:     "Несуществующая воронка",
			IsActive: true,
		}

		_, err := Update(context.Background(), apiClient, pipelineToUpdate)

		if err == nil {
			t.Error("Ожидалась ошибка, но её не было")
		}
	})
}

// TestDeletePipelineError проверяет обработку ошибок при удалении воронки
func TestDeletePipelineError(t *testing.T) {
	t.Run("Ошибка HTTP", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			if _, err := w.Write([]byte(`{"error": "Forbidden to delete main pipeline"}`)); err != nil {
				t.Fatalf("Ошибка при записи ответа: %v", err)
			}
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		err := Delete(context.Background(), apiClient, 123)

		if err == nil {
			t.Error("Ожидалась ошибка, но её не было")
		}
	})

	t.Run("Некорректный код ответа", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK) // Должен быть StatusNoContent
			if _, err := w.Write([]byte(`{}`)); err != nil {
				t.Fatalf("Ошибка при записи ответа: %v", err)
			}
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		err := Delete(context.Background(), apiClient, 123)

		if err == nil {
			t.Error("Ожидалась ошибка из-за некорректного кода ответа, но её не было")
		}
	})
}

// TestGetStatusError проверяет обработку ошибок при получении статуса воронки
func TestGetStatusError(t *testing.T) {
	t.Run("Ошибка HTTP", func(t *testing.T) {
		apiClient := client.NewClient("http://non-existent-domain.example", "test_api_key")

		_, err := GetStatus(context.Background(), apiClient, 123, 999)

		if err == nil {
			t.Error("Ожидалась ошибка, но её не было")
		}
	})
}

// TestCreatePipelineError проверяет обработку ошибок при создании воронки
func TestCreatePipelineError(t *testing.T) {
	t.Run("Ошибка HTTP", func(t *testing.T) {
		apiClient := client.NewClient("http://non-existent-domain.example", "test_api_key")

		invalidPipeline := &Pipeline{
			Name: "Тестовая воронка",
		}

		_, err := Create(context.Background(), apiClient, invalidPipeline)

		if err == nil {
			t.Error("Ожидалась ошибка, но её не было")
		}
	})
}

// TestCreateStatusError проверяет обработку ошибок при создании статуса
func TestCreateStatusError(t *testing.T) {
	t.Run("Ошибка HTTP", func(t *testing.T) {
		apiClient := client.NewClient("http://non-existent-domain.example", "test_api_key")

		testStatus := &Status{
			Name:  "Тестовый статус",
			Color: "#FF0000",
		}

		_, err := CreateStatus(context.Background(), apiClient, 123, testStatus)

		if err == nil {
			t.Error("Ожидалась ошибка, но её не было")
		}
	})
}
