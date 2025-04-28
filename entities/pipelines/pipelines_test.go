package pipelines

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

func TestGetPipeline(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/leads/pipelines/123"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 123,
			"name": "Тестовая воронка",
			"sort": 1,
			"is_main": true,
			"is_active": true,
			"statuses": [
				{
					"id": 456,
					"name": "Новый",
					"sort": 1,
					"color": "#99ccff",
					"type": 1,
					"pipeline_id": 123,
					"is_editable": true
				},
				{
					"id": 789,
					"name": "В работе",
					"sort": 2,
					"color": "#ffcc66",
					"type": 2,
					"pipeline_id": 123,
					"is_editable": true
				}
			]
		}`))
	}))
	defer server.Close()

	// Создаем клиент
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	pipeline, err := GetPipeline(apiClient, 123)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при получении воронки: %v", err)
	}

	if pipeline.ID != 123 {
		t.Errorf("Ожидался ID воронки 123, получен %d", pipeline.ID)
	}

	if pipeline.Name != "Тестовая воронка" {
		t.Errorf("Ожидалось имя воронки 'Тестовая воронка', получено '%s'", pipeline.Name)
	}

	if !pipeline.IsMain {
		t.Errorf("Ожидалось, что воронка основная (IsMain=true)")
	}

	if len(pipeline.Statuses) != 2 {
		t.Errorf("Ожидалось 2 статуса, получено %d", len(pipeline.Statuses))
	} else {
		if pipeline.Statuses[0].ID != 456 {
			t.Errorf("Ожидался ID первого статуса 456, получен %d", pipeline.Statuses[0].ID)
		}
		if pipeline.Statuses[1].ID != 789 {
			t.Errorf("Ожидался ID второго статуса 789, получен %d", pipeline.Statuses[1].ID)
		}
	}
}

func TestCreatePipeline(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/leads/pipelines"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Декодируем тело запроса
		var requestPipeline Pipeline
		if err := json.NewDecoder(r.Body).Decode(&requestPipeline); err != nil {
			t.Fatalf("Ошибка при декодировании тела запроса: %v", err)
		}

		// Проверяем содержимое запроса
		if requestPipeline.Name != "Новая воронка" {
			t.Errorf("Ожидалось имя воронки 'Новая воронка', получено '%s'", requestPipeline.Name)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 456,
			"name": "Новая воронка",
			"sort": 2,
			"is_main": false,
			"is_active": true
		}`))
	}))
	defer server.Close()

	// Создаем клиент
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Создаем воронку для теста
	pipelineToCreate := &Pipeline{
		Name:     "Новая воронка",
		Sort:     2,
		IsMain:   false,
		IsActive: true,
	}

	// Вызываем тестируемый метод
	createdPipeline, err := CreatePipeline(apiClient, pipelineToCreate)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при создании воронки: %v", err)
	}

	if createdPipeline.ID != 456 {
		t.Errorf("Ожидался ID воронки 456, получен %d", createdPipeline.ID)
	}

	if createdPipeline.Name != "Новая воронка" {
		t.Errorf("Ожидалось имя воронки 'Новая воронка', получено '%s'", createdPipeline.Name)
	}

	if createdPipeline.IsMain {
		t.Errorf("Ожидалось, что воронка не основная (IsMain=false)")
	}
}

func TestUpdatePipeline(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "PATCH" {
			t.Errorf("Ожидался метод PATCH, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/leads/pipelines/789"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Декодируем тело запроса
		var requestPipeline Pipeline
		if err := json.NewDecoder(r.Body).Decode(&requestPipeline); err != nil {
			t.Fatalf("Ошибка при декодировании тела запроса: %v", err)
		}

		// Проверяем содержимое запроса
		if requestPipeline.Name != "Обновленная воронка" {
			t.Errorf("Ожидалось имя воронки 'Обновленная воронка', получено '%s'", requestPipeline.Name)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 789,
			"name": "Обновленная воронка",
			"sort": 3,
			"is_main": true,
			"is_active": true
		}`))
	}))
	defer server.Close()

	// Создаем клиент
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Создаем воронку для обновления
	pipelineToUpdate := &Pipeline{
		ID:       789,
		Name:     "Обновленная воронка",
		Sort:     3,
		IsMain:   true,
		IsActive: true,
	}

	// Вызываем тестируемый метод
	updatedPipeline, err := UpdatePipeline(apiClient, pipelineToUpdate)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при обновлении воронки: %v", err)
	}

	if updatedPipeline.ID != 789 {
		t.Errorf("Ожидался ID воронки 789, получен %d", updatedPipeline.ID)
	}

	if updatedPipeline.Name != "Обновленная воронка" {
		t.Errorf("Ожидалось имя воронки 'Обновленная воронка', получено '%s'", updatedPipeline.Name)
	}

	if !updatedPipeline.IsMain {
		t.Errorf("Ожидалось, что воронка основная (IsMain=true)")
	}
}

func TestListPipelines(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/leads/pipelines"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"_embedded": {
				"items": [
					{
						"id": 123,
						"name": "Основная воронка",
						"sort": 1,
						"is_main": true,
						"is_active": true
					},
					{
						"id": 456,
						"name": "Дополнительная воронка",
						"sort": 2,
						"is_main": false,
						"is_active": true
					}
				]
			}
		}`))
	}))
	defer server.Close()

	// Создаем клиент
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	pipelines, err := ListPipelines(apiClient)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при получении списка воронок: %v", err)
	}

	if len(pipelines) != 2 {
		t.Errorf("Ожидалось 2 воронки, получено %d", len(pipelines))
	} else {
		if pipelines[0].ID != 123 {
			t.Errorf("Ожидался ID первой воронки 123, получен %d", pipelines[0].ID)
		}
		if pipelines[1].ID != 456 {
			t.Errorf("Ожидался ID второй воронки 456, получен %d", pipelines[1].ID)
		}
	}
}

func TestDeletePipeline(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "DELETE" {
			t.Errorf("Ожидался метод DELETE, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/leads/pipelines/123"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем успешный ответ без тела
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	// Создаем клиент
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	err := DeletePipeline(apiClient, 123)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при удалении воронки: %v", err)
	}
}

func TestGetStatus(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/leads/pipelines/123/statuses/456"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 456,
			"name": "Новый",
			"sort": 1,
			"color": "#99ccff",
			"type": 1,
			"pipeline_id": 123,
			"is_editable": true
		}`))
	}))
	defer server.Close()

	// Создаем клиент
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	status, err := GetStatus(apiClient, 123, 456)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при получении статуса: %v", err)
	}

	if status.ID != 456 {
		t.Errorf("Ожидался ID статуса 456, получен %d", status.ID)
	}

	if status.Name != "Новый" {
		t.Errorf("Ожидалось имя статуса 'Новый', получено '%s'", status.Name)
	}

	if status.PipelineID != 123 {
		t.Errorf("Ожидался ID воронки 123, получен %d", status.PipelineID)
	}

	if status.Color != "#99ccff" {
		t.Errorf("Ожидался цвет статуса '#99ccff', получен '%s'", status.Color)
	}
}

func TestCreateStatus(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/leads/pipelines/123/statuses"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Декодируем тело запроса
		var requestStatus Status
		if err := json.NewDecoder(r.Body).Decode(&requestStatus); err != nil {
			t.Fatalf("Ошибка при декодировании тела запроса: %v", err)
		}

		// Проверяем содержимое запроса
		if requestStatus.Name != "Новый статус" {
			t.Errorf("Ожидалось имя статуса 'Новый статус', получено '%s'", requestStatus.Name)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 789,
			"name": "Новый статус",
			"sort": 3,
			"color": "#ff9999",
			"type": 2,
			"pipeline_id": 123,
			"is_editable": true
		}`))
	}))
	defer server.Close()

	// Создаем клиент
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Создаем статус для теста
	statusToCreate := &Status{
		Name:  "Новый статус",
		Sort:  3,
		Color: "#ff9999",
		Type:  2,
	}

	// Вызываем тестируемый метод
	createdStatus, err := CreateStatus(apiClient, 123, statusToCreate)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при создании статуса: %v", err)
	}

	if createdStatus.ID != 789 {
		t.Errorf("Ожидался ID статуса 789, получен %d", createdStatus.ID)
	}

	if createdStatus.Name != "Новый статус" {
		t.Errorf("Ожидалось имя статуса 'Новый статус', получено '%s'", createdStatus.Name)
	}

	if createdStatus.PipelineID != 123 {
		t.Errorf("Ожидался ID воронки 123, получен %d", createdStatus.PipelineID)
	}
}
