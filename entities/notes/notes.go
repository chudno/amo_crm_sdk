// Package notes предоставляет методы для взаимодействия с сущностями "Примечания" в API amoCRM.
package notes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chudno/amo_crm_sdk/client"
)

// Note представляет собой структуру примечания в amoCRM.
type Note struct {
	ID         int        `json:"id,omitempty"`
	EntityID   int        `json:"entity_id,omitempty"`
	EntityType string     `json:"entity_type,omitempty"` // leads, contacts, companies, customers
	NoteType   string     `json:"note_type,omitempty"`
	Text       string     `json:"text,omitempty"`
	CreatedBy  int        `json:"created_by,omitempty"`
	CreatedAt  int64      `json:"created_at,omitempty"`
	UpdatedAt  int64      `json:"updated_at,omitempty"`
	Params     NoteParams `json:"params,omitempty"`
}

// Константы типов примечаний в amoCRM v4
const (
	TypeCommon                 = "common"
	TypeCallIn                 = "call_in"
	TypeCallOut                = "call_out"
	TypeServiceMessage         = "service_message"
	TypeIncomingChatMessage    = "incoming_chat_message"
	TypeOutgoingChatMessage    = "outgoing_chat_message"
	TypeSmsIn                  = "sms_in"
	TypeSmsOut                 = "sms_out"
	TypeExtendedServiceMessage = "extended_service_message"
	TypeAttachment             = "attachment"
	TypeAmomailMessage         = "amomail_message"
	TypeGeolocation            = "geolocation"
)

// NoteParams содержит дополнительные параметры примечания в зависимости от типа
type NoteParams struct {
	Text        string `json:"text,omitempty"`
	Service     string `json:"service,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Email       string `json:"email,omitempty"`
	Link        string `json:"link,omitempty"`
}

// Get получает примечание по его ID.
func Get(ctx context.Context, apiClient *client.Client, entityType string, entityID int, noteID int) (*Note, error) {
	url := fmt.Sprintf("%s/api/v4/%s/%d/notes/%d", apiClient.GetBaseURL(), entityType, entityID, noteID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var note Note
	if err := json.NewDecoder(resp.Body).Decode(&note); err != nil {
		return nil, err
	}

	return &note, nil
}

// Create создает новое примечание в amoCRM.
func Create(ctx context.Context, apiClient *client.Client, entityType string, entityID int, note *Note) (*Note, error) {
	apiURL := fmt.Sprintf("%s/api/v4/%s/%d/notes", apiClient.GetBaseURL(), entityType, entityID)
	noteJSON, err := json.Marshal([]*Note{note})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(noteJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var response struct {
		Embedded struct {
			Notes []*Note `json:"notes"`
		} `json:"_embedded"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if len(response.Embedded.Notes) == 0 {
		return nil, fmt.Errorf("не удалось создать примечание")
	}

	return response.Embedded.Notes[0], nil
}

// Update обновляет существующее примечание в amoCRM.
func Update(ctx context.Context, apiClient *client.Client, entityType string, entityID int, note *Note) (*Note, error) {
	url := fmt.Sprintf("%s/api/v4/%s/%d/notes/%d", apiClient.GetBaseURL(), entityType, entityID, note.ID)
	noteJSON, err := json.Marshal(note)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewBuffer(noteJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var updatedNote Note
	if err := json.NewDecoder(resp.Body).Decode(&updatedNote); err != nil {
		return nil, err
	}

	return &updatedNote, nil
}

// List получает список примечаний для указанной сущности с возможностью фильтрации и пагинации.
func List(ctx context.Context, apiClient *client.Client, entityType string, entityID int, limit int, page int) ([]Note, error) {
	url := fmt.Sprintf("%s/api/v4/%s/%d/notes?limit=%d&page=%d", apiClient.GetBaseURL(), entityType, entityID, limit, page)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var notes struct {
		Embedded struct {
			Items []Note `json:"notes"`
		} `json:"_embedded"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&notes); err != nil {
		return nil, err
	}

	return notes.Embedded.Items, nil
}
