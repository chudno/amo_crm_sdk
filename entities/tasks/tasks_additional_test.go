package tasks

import (
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
		// Создаем тестовый сервер
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Проверяем метод запроса
			if r.Method != "PATCH" {
				t.Errorf("Ожидался метод PATCH, получен %s", r.Method)
			}

			// Проверяем путь запроса
			expectedPath := "/api/v4/tasks/123"
			if r.URL.Path != expectedPath {
				t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
			}

			// Проверяем тело запроса
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

			// Отправляем ответ
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

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Создаем задачу для обновления
		taskToUpdate := &Task{
			ID:                123,
			Text:              "Обновленная задача",
			ResponsibleUserID: 456,
			EntityID:          789,
			EntityType:        "leads",
			CompleteTill:      1609632000,
		}

		// Вызываем тестируемый метод
		updatedTask, err := UpdateTask(apiClient, taskToUpdate)

		// Проверяем результаты
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
		// Создаем клиент API (не нужен реальный сервер, т.к. ошибка будет раньше)
		apiClient := client.NewClient("http://example.com", "test_api_key")

		// Создаем задачу без ID
		taskWithoutID := &Task{
			Text:              "Задача без ID",
			ResponsibleUserID: 456,
		}

		// Вызываем тестируемый метод
		_, err := UpdateTask(apiClient, taskWithoutID)

		// Проверяем, что возникла ошибка
		if err == nil {
			t.Error("Ожидалась ошибка о не указанном ID, но ее не возникло")
		}
	})

	// Сценарий: ошибка от сервера
	t.Run("Ошибка от сервера", func(t *testing.T) {
		// Используем домен, который не существует, чтобы вызвать сетевую ошибку
		apiClient := client.NewClient("http://non-existent-domain.example", "test_api_key")

		// Создаем задачу для обновления
		taskToUpdate := &Task{
			ID:   999, // несуществующий ID
			Text: "Несуществующая задача",
		}

		// Вызываем тестируемый метод
		_, err := UpdateTask(apiClient, taskToUpdate)

		// Проверяем, что возникла ошибка
		if err == nil {
			t.Error("Ожидалась ошибка от сервера, но ее не возникло")
		}
	})
}

// TestCompleteTask проверяет функциональность отметки задачи как выполненной
func TestCompleteTask(t *testing.T) {
	// Сценарий: успешное выполнение задачи
	t.Run("Успешное выполнение", func(t *testing.T) {
		// Создаем тестовый сервер
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Проверяем метод запроса
			if r.Method != "PATCH" {
				t.Errorf("Ожидался метод PATCH, получен %s", r.Method)
			}

			// Проверяем путь запроса
			expectedPath := "/api/v4/tasks/123"
			if r.URL.Path != expectedPath {
				t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
			}

			// Проверяем тело запроса
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

			if task.Result != "Задача выполнена успешно" {
				t.Errorf("Ожидался результат 'Задача выполнена успешно', получен '%s'", task.Result)
			}

			// Отправляем ответ
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(`{
				"id": 123,
				"text": "Тестовая задача",
				"responsible_user_id": 456,
				"is_completed": true,
				"result": "Задача выполнена успешно",
				"updated_at": 1609545600
			}`)); err != nil {
				t.Fatalf("Ошибка при записи ответа: %v", err)
			}
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемый метод
		completedTask, err := CompleteTask(apiClient, 123, "Задача выполнена успешно")

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при выполнении задачи: %v", err)
		}

		if completedTask.ID != 123 {
			t.Errorf("Ожидался ID задачи 123, получен %d", completedTask.ID)
		}

		if !completedTask.IsCompleted {
			t.Error("Ожидалось, что задача будет отмечена как выполненная")
		}

		if completedTask.Result != "Задача выполнена успешно" {
			t.Errorf("Ожидался результат 'Задача выполнена успешно', получен '%s'", completedTask.Result)
		}
	})

	// Сценарий: ошибка обновления при выполнении задачи
	t.Run("Ошибка при выполнении", func(t *testing.T) {
		// Используем домен, который не существует, чтобы вызвать сетевую ошибку
		apiClient := client.NewClient("http://non-existent-domain.example", "test_api_key")

		// Вызываем тестируемый метод
		_, err := CompleteTask(apiClient, 123, "Результат выполнения")

		// Проверяем, что возникла ошибка
		if err == nil {
			t.Error("Ожидалась ошибка от сервера, но ее не возникло")
		}
	})
}

// TestDeleteTask проверяет функциональность удаления задачи
func TestDeleteTask(t *testing.T) {
	// Сценарий: успешное удаление задачи
	t.Run("Успешное удаление", func(t *testing.T) {
		// Создаем тестовый сервер
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Проверяем метод запроса
			if r.Method != "DELETE" {
				t.Errorf("Ожидался метод DELETE, получен %s", r.Method)
			}

			// Проверяем путь запроса
			expectedPath := "/api/v4/tasks/123"
			if r.URL.Path != expectedPath {
				t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
			}

			// Отправляем ответ успешного удаления
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемый метод
		err := DeleteTask(apiClient, 123)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при удалении задачи: %v", err)
		}
	})

	// Сценарий: ошибка при удалении
	t.Run("Ошибка при удалении", func(t *testing.T) {
		// Используем домен, который не существует, чтобы вызвать сетевую ошибку
		apiClient := client.NewClient("http://non-existent-domain.example", "test_api_key")

		// Вызываем тестируемый метод с несуществующим ID
		err := DeleteTask(apiClient, 999)

		// Проверяем, что возникла ошибка
		if err == nil {
			t.Error("Ожидалась ошибка от сервера, но ее не возникло")
		}
	})
}
