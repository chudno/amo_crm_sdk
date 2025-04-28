// Пакет notes предоставляет методы для взаимодействия с сущностями "Примечания" в API amoCRM.
package notes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/chudno/amo_crm_sdk/client"
	"net/http"
	"time"
)

// Note представляет собой структуру примечания в amoCRM.
type Note struct {
	ID         int        `json:"id"`
	EntityID   int        `json:"entity_id"`
	EntityType string     `json:"entity_type"` // leads, contacts, companies, customers
	NoteType   int        `json:"note_type"`
	Text       string     `json:"text,omitempty"`
	CreatedBy  int        `json:"created_by"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	Params     NoteParams `json:"params,omitempty"`
}

// NoteParams содержит дополнительные параметры примечания в зависимости от типа
type NoteParams struct {
	Text        string `json:"text,omitempty"`
	Service     string `json:"service,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Email       string `json:"email,omitempty"`
	Link        string `json:"link,omitempty"`
}

// GetNote получает примечание по его ID.
func GetNote(apiClient *client.Client, entityType string, entityID int, noteID int) (*Note, error) {
	url := fmt.Sprintf("%s/api/v4/%s/%d/notes/%d", apiClient.GetBaseURL(), entityType, entityID, noteID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var note Note
	if err := json.NewDecoder(resp.Body).Decode(&note); err != nil {
		return nil, err
	}

	return &note, nil
}

// CreateNote создает новое примечание в amoCRM.
func CreateNote(apiClient *client.Client, entityType string, entityID int, note *Note) (*Note, error) {
	url := fmt.Sprintf("%s/api/v4/%s/%d/notes", apiClient.GetBaseURL(), entityType, entityID)
	noteJSON, err := json.Marshal(note)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(noteJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var newNote Note
	if err := json.NewDecoder(resp.Body).Decode(&newNote); err != nil {
		return nil, err
	}

	return &newNote, nil
}

// UpdateNote обновляет существующее примечание в amoCRM.
func UpdateNote(apiClient *client.Client, entityType string, entityID int, note *Note) (*Note, error) {
	url := fmt.Sprintf("%s/api/v4/%s/%d/notes/%d", apiClient.GetBaseURL(), entityType, entityID, note.ID)
	noteJSON, err := json.Marshal(note)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(noteJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var updatedNote Note
	if err := json.NewDecoder(resp.Body).Decode(&updatedNote); err != nil {
		return nil, err
	}

	return &updatedNote, nil
}

// ListNotes получает список примечаний для указанной сущности с возможностью фильтрации и пагинации.
func ListNotes(apiClient *client.Client, entityType string, entityID int, limit int, page int) ([]Note, error) {
	url := fmt.Sprintf("%s/api/v4/%s/%d/notes?limit=%d&page=%d", apiClient.GetBaseURL(), entityType, entityID, limit, page)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var notes struct {
		Embedded struct {
			Items []Note `json:"items"`
		} `json:"_embedded"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&notes); err != nil {
		return nil, err
	}

	return notes.Embedded.Items, nil
}

// DeleteNote удаляет примечание по его ID.
func DeleteNote(apiClient *client.Client, entityType string, entityID int, noteID int) error {
	url := fmt.Sprintf("%s/api/v4/%s/%d/notes/%d", apiClient.GetBaseURL(), entityType, entityID, noteID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
