package notes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chudno/amo_crm_sdk/client"
)

func TestGetNote(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/leads/123/notes/456"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 456,
			"entity_id": 123,
			"entity_type": "leads",
			"note_type": 4,
			"text": "Тестовое примечание",
			"created_by": 789,
			"created_at": "2023-01-01T10:00:00Z",
			"updated_at": "2023-01-01T11:00:00Z"
		}`))
	}))
	defer server.Close()

	// Создаем клиент
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	note, err := GetNote(apiClient, "leads", 123, 456)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при получении примечания: %v", err)
	}

	if note.ID != 456 {
		t.Errorf("Ожидался ID примечания 456, получен %d", note.ID)
	}

	if note.EntityID != 123 {
		t.Errorf("Ожидался ID сущности 123, получен %d", note.EntityID)
	}

	if note.EntityType != "leads" {
		t.Errorf("Ожидался тип сущности 'leads', получен '%s'", note.EntityType)
	}

	if note.Text != "Тестовое примечание" {
		t.Errorf("Ожидался текст примечания 'Тестовое примечание', получен '%s'", note.Text)
	}
}

func TestCreateNote(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/contacts/789/notes"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Декодируем тело запроса
		var requestNote Note
		if err := json.NewDecoder(r.Body).Decode(&requestNote); err != nil {
			t.Fatalf("Ошибка при декодировании тела запроса: %v", err)
		}

		// Проверяем содержимое запроса
		if requestNote.Text != "Новое примечание" {
			t.Errorf("Ожидался текст примечания 'Новое примечание', получен '%s'", requestNote.Text)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 123,
			"entity_id": 789,
			"entity_type": "contacts",
			"note_type": 4,
			"text": "Новое примечание",
			"created_by": 456,
			"created_at": "2023-01-01T10:00:00Z",
			"updated_at": "2023-01-01T10:00:00Z"
		}`))
	}))
	defer server.Close()

	// Создаем клиент
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Создаем примечание для теста
	createdAt, _ := time.Parse(time.RFC3339, "2023-01-01T10:00:00Z")
	updatedAt, _ := time.Parse(time.RFC3339, "2023-01-01T10:00:00Z")
	noteToCreate := &Note{
		EntityID:   789,
		EntityType: "contacts",
		NoteType:   4,
		Text:       "Новое примечание",
		CreatedBy:  456,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}

	// Вызываем тестируемый метод
	createdNote, err := CreateNote(apiClient, "contacts", 789, noteToCreate)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при создании примечания: %v", err)
	}

	if createdNote.ID != 123 {
		t.Errorf("Ожидался ID примечания 123, получен %d", createdNote.ID)
	}

	if createdNote.Text != "Новое примечание" {
		t.Errorf("Ожидался текст примечания 'Новое примечание', получен '%s'", createdNote.Text)
	}
}

func TestUpdateNote(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "PATCH" {
			t.Errorf("Ожидался метод PATCH, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/leads/456/notes/123"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Декодируем тело запроса
		var requestNote Note
		if err := json.NewDecoder(r.Body).Decode(&requestNote); err != nil {
			t.Fatalf("Ошибка при декодировании тела запроса: %v", err)
		}

		// Проверяем содержимое запроса
		if requestNote.Text != "Обновленное примечание" {
			t.Errorf("Ожидался текст примечания 'Обновленное примечание', получен '%s'", requestNote.Text)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 123,
			"entity_id": 456,
			"entity_type": "leads",
			"note_type": 4,
			"text": "Обновленное примечание",
			"created_by": 789,
			"created_at": "2023-01-01T10:00:00Z",
			"updated_at": "2023-01-01T12:00:00Z"
		}`))
	}))
	defer server.Close()

	// Создаем клиент
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Создаем примечание для обновления
	createdAt, _ := time.Parse(time.RFC3339, "2023-01-01T10:00:00Z")
	updatedAt, _ := time.Parse(time.RFC3339, "2023-01-01T12:00:00Z")
	noteToUpdate := &Note{
		ID:         123,
		EntityID:   456,
		EntityType: "leads",
		NoteType:   4,
		Text:       "Обновленное примечание",
		CreatedBy:  789,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}

	// Вызываем тестируемый метод
	updatedNote, err := UpdateNote(apiClient, "leads", 456, noteToUpdate)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при обновлении примечания: %v", err)
	}

	if updatedNote.ID != 123 {
		t.Errorf("Ожидался ID примечания 123, получен %d", updatedNote.ID)
	}

	if updatedNote.Text != "Обновленное примечание" {
		t.Errorf("Ожидался текст примечания 'Обновленное примечание', получен '%s'", updatedNote.Text)
	}
}

func TestListNotes(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/companies/123/notes"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Проверяем параметры запроса
		query := r.URL.Query()
		if query.Get("limit") != "10" {
			t.Errorf("Ожидался параметр limit=10, получен %s", query.Get("limit"))
		}
		if query.Get("page") != "1" {
			t.Errorf("Ожидался параметр page=1, получен %s", query.Get("page"))
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"_embedded": {
				"items": [
					{
						"id": 456,
						"entity_id": 123,
						"entity_type": "companies",
						"note_type": 4,
						"text": "Примечание 1",
						"created_by": 789,
						"created_at": "2023-01-01T10:00:00Z",
						"updated_at": "2023-01-01T11:00:00Z"
					},
					{
						"id": 789,
						"entity_id": 123,
						"entity_type": "companies",
						"note_type": 4,
						"text": "Примечание 2",
						"created_by": 789,
						"created_at": "2023-01-02T10:00:00Z",
						"updated_at": "2023-01-02T11:00:00Z"
					}
				]
			}
		}`))
	}))
	defer server.Close()

	// Создаем клиент
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	notes, err := ListNotes(apiClient, "companies", 123, 10, 1)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при получении списка примечаний: %v", err)
	}

	if len(notes) != 2 {
		t.Errorf("Ожидалось 2 примечания, получено %d", len(notes))
	}

	if notes[0].ID != 456 {
		t.Errorf("Ожидался ID первого примечания 456, получен %d", notes[0].ID)
	}

	if notes[1].ID != 789 {
		t.Errorf("Ожидался ID второго примечания 789, получен %d", notes[1].ID)
	}
}

func TestDeleteNote(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "DELETE" {
			t.Errorf("Ожидался метод DELETE, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/leads/123/notes/456"
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
	err := DeleteNote(apiClient, "leads", 123, 456)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при удалении примечания: %v", err)
	}
}

func TestDeleteNoteError(t *testing.T) {
	// Создаем тестовый сервер с ошибкой
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Отправляем ответ с ошибкой
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"status": "error", "message": "Note not found"}`))
	}))
	defer server.Close()

	// Создаем клиент
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	err := DeleteNote(apiClient, "leads", 123, 456)

	// Проверяем результаты - должна быть ошибка
	if err == nil {
		t.Fatal("Ожидалась ошибка при удалении несуществующего примечания, но её не было")
	}
}
