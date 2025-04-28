package tasks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chudno/amo_crm_sdk/client"
)

func TestGetTask(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/tasks/123"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 123,
			"text": "Тестовая задача",
			"responsible_user_id": 456,
			"entity_id": 789,
			"entity_type": "leads",
			"created_at": 1609459200,
			"updated_at": 1609545600,
			"complete_till": 1609632000
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	task, err := GetTask(apiClient, 123)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при получении задачи: %v", err)
	}

	if task.ID != 123 {
		t.Errorf("Ожидался ID задачи 123, получен %d", task.ID)
	}

	if task.Text != "Тестовая задача" {
		t.Errorf("Ожидался текст задачи 'Тестовая задача', получен '%s'", task.Text)
	}

	if task.ResponsibleUserID != 456 {
		t.Errorf("Ожидался ID ответственного пользователя 456, получен %d", task.ResponsibleUserID)
	}

	if task.EntityID != 789 {
		t.Errorf("Ожидался ID сущности 789, получен %d", task.EntityID)
	}

	if task.EntityType != "leads" {
		t.Errorf("Ожидался тип сущности 'leads', получен '%s'", task.EntityType)
	}
}

func TestCreateTask(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/tasks"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"_embedded": {
				"tasks": [{
					"id": 789,
					"text": "Новая задача",
					"responsible_user_id": 456,
					"entity_id": 123,
					"entity_type": "leads",
					"created_at": 1609459200,
					"updated_at": 1609545600,
					"complete_till": 1609632000
				}]
			}
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Создаем задачу для теста
	taskToCreate := &Task{
		Text:              "Новая задача",
		ResponsibleUserID: 456,
		EntityID:          123,
		EntityType:        "leads",
		CompleteTill:      1609632000,
	}

	// Вызываем тестируемый метод
	createdTask, err := CreateTask(apiClient, taskToCreate)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при создании задачи: %v", err)
	}

	if createdTask.ID != 789 {
		t.Errorf("Ожидался ID задачи 789, получен %d", createdTask.ID)
	}

	if createdTask.Text != "Новая задача" {
		t.Errorf("Ожидался текст задачи 'Новая задача', получен '%s'", createdTask.Text)
	}

	if createdTask.EntityID != 123 {
		t.Errorf("Ожидался ID сущности 123, получен %d", createdTask.EntityID)
	}
}

func TestListTasks(t *testing.T) {
	tests := []struct {
		name         string
		page         int
		limit        int
		responseCode int
		responseBody string
		expectError  bool
		expectedLen  int
	}{
		{
			name:         "Успешное получение списка задач",
			page:         1,
			limit:        50,
			responseCode: http.StatusOK,
			responseBody: `{"_embedded":{"tasks":[{"id":12345,"text":"Задача 1"},{"id":12346,"text":"Задача 2"}]}}`,
			expectedLen:  2,
		},
		{
			name:         "Пустой список задач",
			page:         1,
			limit:        50,
			responseCode: http.StatusOK,
			responseBody: `{"_embedded":{"tasks":[]}}`,
			expectedLen:  0,
		},
		{
			name:         "Ошибка сервера",
			page:         1,
			limit:        50,
			responseCode: http.StatusInternalServerError,
			responseBody: `{"error":"server_error"}`,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Проверяем метод запроса
				if r.Method != "GET" {
					t.Errorf("Ожидался метод GET, получен %s", r.Method)
				}

				// Проверяем URL запроса
				expectedPath := "/api/v4/tasks"
				if r.URL.Path != expectedPath {
					t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
				}

				// Проверяем параметры запроса
				query := r.URL.Query()
				if query.Get("page") != fmt.Sprintf("%d", tt.page) {
					t.Errorf("Ожидался параметр page=%d, получен %s", tt.page, query.Get("page"))
				}
				if query.Get("limit") != fmt.Sprintf("%d", tt.limit) {
					t.Errorf("Ожидался параметр limit=%d, получен %s", tt.limit, query.Get("limit"))
				}

				// Устанавливаем код ответа и тело
				w.WriteHeader(tt.responseCode)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			// Создаем клиент
			apiClient := client.NewClient(server.URL, "test_api_key")

			// Вызываем тестируемую функцию
			tasks, err := ListTasks(apiClient, tt.limit, tt.page, nil)

			// Проверяем результаты
			if tt.expectError && err == nil {
				t.Error("Ожидалась ошибка, но ее не было")
			}

			if !tt.expectError {
				if err != nil {
					t.Errorf("Неожиданная ошибка: %v", err)
				}

				if len(tasks) != tt.expectedLen {
					t.Errorf("Ожидалось %d задач, получено %d", tt.expectedLen, len(tasks))
				}
			}
		})
	}
}

func TestCreateTaskForEntity(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/tasks"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Проверяем тело запроса
		var tasks []*Task
		if err := json.NewDecoder(r.Body).Decode(&tasks); err != nil {
			t.Errorf("Ошибка декодирования тела запроса: %v", err)
		}

		if len(tasks) != 1 {
			t.Errorf("Ожидалась 1 задача в запросе, получено %d", len(tasks))
		}

		task := tasks[0]
		if task.EntityType != "leads" {
			t.Errorf("Ожидался тип сущности 'leads', получен '%s'", task.EntityType)
		}

		if task.EntityID != 123 {
			t.Errorf("Ожидался ID сущности 123, получен %d", task.EntityID)
		}

		if task.Text != "Тестовая задача для лида" {
			t.Errorf("Ожидался текст задачи 'Тестовая задача для лида', получен '%s'", task.Text)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"_embedded": {
				"tasks": [{
					"id": 456,
					"text": "Тестовая задача для лида",
					"responsible_user_id": 789,
					"entity_id": 123,
					"entity_type": "leads",
					"created_at": 1609459200,
					"updated_at": 1609545600,
					"complete_till": 1609632000
				}]
			}
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Задаем параметры для создания задачи
	entityType := "leads"
	entityID := 123
	taskTypeID := 1
	text := "Тестовая задача для лида"
	completeTill := time.Unix(1609632000, 0)
	responsibleUserID := 789

	// Вызываем тестируемый метод
	createdTask, err := CreateTaskForEntity(apiClient, entityType, entityID, taskTypeID, text, completeTill, responsibleUserID)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при создании задачи для сущности: %v", err)
	}

	if createdTask.ID != 456 {
		t.Errorf("Ожидался ID задачи 456, получен %d", createdTask.ID)
	}

	if createdTask.Text != "Тестовая задача для лида" {
		t.Errorf("Ожидался текст задачи 'Тестовая задача для лида', получен '%s'", createdTask.Text)
	}

	if createdTask.EntityID != 123 {
		t.Errorf("Ожидался ID сущности 123, получен %d", createdTask.EntityID)
	}

	if createdTask.EntityType != "leads" {
		t.Errorf("Ожидался тип сущности 'leads', получен '%s'", createdTask.EntityType)
	}
}
