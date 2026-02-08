package tasks

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

// TestUpdateTask проверяет функциональность обновления задачи
func TestUpdateTask(t *testing.T) {
	// Сценарий: успешное обновление задачи
	t.Run("Успешное обновление", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "PATCH" {
				t.Errorf("Ожидался метод PATCH, получен %s", r.Method)
			}

			expectedPath := "/api/v4/tasks/123"
			if r.URL.Path != expectedPath {
				t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("Ошибка чтения тела запроса: %v", err)
			}

			var task Task
			if err := json.Unmarshal(body, &task); err != nil {
				t.Fatalf("Ошибка декодирования тела запроса: %v", err)
			}

			if task.ID != 123 {
				t.Errorf("Ожидался ID задачи 123, получен %d", task.ID)
			}

			if task.Text != "Обновленная задача" {
				t.Errorf("Ожидался текст задачи 'Обновленная задача', получен '%s'", task.Text)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(`{
				"id": 123,
				"text": "Обновленная задача",
				"responsible_user_id": 456,
				"entity_id": 789,
				"entity_type": "leads",
				"updated_at": 1609545600,
				"complete_till": 1609632000
			}`)); err != nil {
				t.Fatalf("Ошибка при записи ответа: %v", err)
			}
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		taskToUpdate := &Task{
			ID:                123,
			Text:              "Обновленная задача",
			ResponsibleUserID: 456,
			EntityID:          789,
			EntityType:        "leads",
			CompleteTill:      1609632000,
		}

		updatedTask, err := Update(context.Background(), apiClient, taskToUpdate)

		if err != nil {
			t.Fatalf("Ошибка при обновлении задачи: %v", err)
		}

		if updatedTask.ID != 123 {
			t.Errorf("Ожидался ID задачи 123, получен %d", updatedTask.ID)
		}

		if updatedTask.Text != "Обновленная задача" {
			t.Errorf("Ожидался текст задачи 'Обновленная задача', получен '%s'", updatedTask.Text)
		}
	})

	// Сценарий: попытка обновления задачи без указания ID
	t.Run("Ошибка: ID не указан", func(t *testing.T) {
		apiClient := client.NewClient("http://example.com", "test_api_key")

		taskWithoutID := &Task{
			Text:              "Задача без ID",
			ResponsibleUserID: 456,
		}

		_, err := Update(context.Background(), apiClient, taskWithoutID)

		if err == nil {
			t.Error("Ожидалась ошибка о не указанном ID, но ее не возникло")
		}
	})

	// Сценарий: ошибка от сервера
	t.Run("Ошибка от сервера", func(t *testing.T) {
		apiClient := client.NewClient("http://non-existent-domain.example", "test_api_key")

		taskToUpdate := &Task{
			ID:   999, // несуществующий ID
			Text: "Несуществующая задача",
		}

		_, err := Update(context.Background(), apiClient, taskToUpdate)

		if err == nil {
			t.Error("Ожидалась ошибка от сервера, но ее не возникло")
		}
	})
}

// TestCompleteTask проверяет функциональность отметки задачи как выполненной
func TestCompleteTask(t *testing.T) {
	t.Run("Успешное выполнение", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "PATCH" {
				t.Errorf("Ожидался метод PATCH, получен %s", r.Method)
			}

			expectedPath := "/api/v4/tasks/123"
			if r.URL.Path != expectedPath {
				t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("Ошибка чтения тела запроса: %v", err)
			}

			var task Task
			if err := json.Unmarshal(body, &task); err != nil {
				t.Fatalf("Ошибка декодирования тела запроса: %v", err)
			}

			if task.ID != 123 {
				t.Errorf("Ожидался ID задачи 123, получен %d", task.ID)
			}

			if !task.IsCompleted {
				t.Error("Ожидалось, что задача будет отмечена как выполненная")
			}

			// Result теперь json.RawMessage, проверяем что он содержит ожидаемый текст
			var resultObj map[string]string
			if err := json.Unmarshal(task.Result, &resultObj); err != nil {
				t.Errorf("Ошибка разбора result: %v", err)
			} else if resultObj["text"] != "Задача выполнена успешно" {
				t.Errorf("Ожидался результат с текстом 'Задача выполнена успешно', получен '%s'", string(task.Result))
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(`{
				"id": 123,
				"text": "Тестовая задача",
				"responsible_user_id": 456,
				"is_completed": true,
				"result": [{"text": "Задача выполнена успешно"}],
				"updated_at": 1609545600
			}`)); err != nil {
				t.Fatalf("Ошибка при записи ответа: %v", err)
			}
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		completedTask, err := Complete(context.Background(), apiClient, 123, "Задача выполнена успешно")

		if err != nil {
			t.Fatalf("Ошибка при выполнении задачи: %v", err)
		}

		if completedTask.ID != 123 {
			t.Errorf("Ожидался ID задачи 123, получен %d", completedTask.ID)
		}

		if !completedTask.IsCompleted {
			t.Error("Ожидалось, что задача будет отмечена как выполненная")
		}

		var completedResultArr []map[string]string
		if err := json.Unmarshal(completedTask.Result, &completedResultArr); err != nil {
			t.Errorf("Ошибка разбора result: %v", err)
		} else if len(completedResultArr) == 0 || completedResultArr[0]["text"] != "Задача выполнена успешно" {
			t.Errorf("Ожидался результат с текстом 'Задача выполнена успешно', получен '%s'", string(completedTask.Result))
		}
	})

	// Сценарий: ошибка обновления при выполнении задачи
	t.Run("Ошибка при выполнении", func(t *testing.T) {
		apiClient := client.NewClient("http://non-existent-domain.example", "test_api_key")

		_, err := Complete(context.Background(), apiClient, 123, "Результат выполнения")

		if err == nil {
			t.Error("Ожидалась ошибка от сервера, но ее не возникло")
		}
	})
}
