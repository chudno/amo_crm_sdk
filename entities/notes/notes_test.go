package notes

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

func TestGetNote(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		expectedPath := "/api/v4/leads/123/notes/456"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 456,
			"entity_id": 123,
			"entity_type": "leads",
			"note_type": "common",
			"text": "Тестовое примечание",
			"created_by": 789,
			"created_at": 1672567200,
			"updated_at": 1672570800
		}`))
	}))
	defer server.Close()

	apiClient := client.NewClient(server.URL, "test_api_key")

	note, err := Get(context.Background(), apiClient, "leads", 123, 456)

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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		expectedPath := "/api/v4/contacts/789/notes"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		var requestNotes []*Note
		if err := json.NewDecoder(r.Body).Decode(&requestNotes); err != nil {
			t.Fatalf("Ошибка при декодировании тела запроса: %v", err)
		}

		if len(requestNotes) == 0 {
			t.Fatalf("Ожидался массив с примечанием, получен пустой массив")
		}
		if requestNotes[0].Text != "Новое примечание" {
			t.Errorf("Ожидался текст примечания 'Новое примечание', получен '%s'", requestNotes[0].Text)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"_embedded": {
				"notes": [
					{
						"id": 123,
						"entity_id": 789,
						"entity_type": "contacts",
						"note_type": "common",
						"text": "Новое примечание",
						"created_by": 456,
						"created_at": 1672567200,
						"updated_at": 1672567200
					}
				]
			}
		}`))
	}))
	defer server.Close()

	apiClient := client.NewClient(server.URL, "test_api_key")

	noteToCreate := &Note{
		EntityID:   789,
		EntityType: "contacts",
		NoteType:   TypeCommon,
		Text:       "Новое примечание",
		CreatedBy:  456,
	}

	createdNote, err := Create(context.Background(), apiClient, "contacts", 789, noteToCreate)

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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("Ожидался метод PATCH, получен %s", r.Method)
		}

		expectedPath := "/api/v4/leads/456/notes/123"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		var requestNote Note
		if err := json.NewDecoder(r.Body).Decode(&requestNote); err != nil {
			t.Fatalf("Ошибка при декодировании тела запроса: %v", err)
		}

		if requestNote.Text != "Обновленное примечание" {
			t.Errorf("Ожидался текст примечания 'Обновленное примечание', получен '%s'", requestNote.Text)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 123,
			"entity_id": 456,
			"entity_type": "leads",
			"note_type": "common",
			"text": "Обновленное примечание",
			"created_by": 789,
			"created_at": 1672567200,
			"updated_at": 1672574400
		}`))
	}))
	defer server.Close()

	apiClient := client.NewClient(server.URL, "test_api_key")

	noteToUpdate := &Note{
		ID:         123,
		EntityID:   456,
		EntityType: "leads",
		NoteType:   TypeCommon,
		Text:       "Обновленное примечание",
		CreatedBy:  789,
	}

	updatedNote, err := Update(context.Background(), apiClient, "leads", 456, noteToUpdate)

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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		expectedPath := "/api/v4/companies/123/notes"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		query := r.URL.Query()
		if query.Get("limit") != "10" {
			t.Errorf("Ожидался параметр limit=10, получен %s", query.Get("limit"))
		}
		if query.Get("page") != "1" {
			t.Errorf("Ожидался параметр page=1, получен %s", query.Get("page"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"_embedded": {
				"notes": [
					{
						"id": 456,
						"entity_id": 123,
						"entity_type": "companies",
						"note_type": "common",
						"text": "Примечание 1",
						"created_by": 789,
						"created_at": 1672567200,
						"updated_at": 1672570800
					},
					{
						"id": 789,
						"entity_id": 123,
						"entity_type": "companies",
						"note_type": "common",
						"text": "Примечание 2",
						"created_by": 789,
						"created_at": 1672653600,
						"updated_at": 1672657200
					}
				]
			}
		}`))
	}))
	defer server.Close()

	apiClient := client.NewClient(server.URL, "test_api_key")

	notes, err := List(context.Background(), apiClient, "companies", 123, 10, 1)

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
