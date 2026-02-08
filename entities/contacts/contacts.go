// Package contacts предоставляет методы для взаимодействия с сущностями "Контакты" в API amoCRM.
package contacts

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/chudno/amo_crm_sdk/client"
	"github.com/chudno/amo_crm_sdk/utils/custom_fields"
)

// Contact представляет собой структуру контакта в amoCRM.
type Contact struct {
	ID                 int                   `json:"id"`
	Name               string                `json:"name"`
	FirstName          string                `json:"first_name,omitempty"`
	LastName           string                `json:"last_name,omitempty"`
	ResponsibleUserID  int                   `json:"responsible_user_id,omitempty"`
	GroupID            int                   `json:"group_id,omitempty"`
	CreatedBy          int                   `json:"created_by,omitempty"`
	UpdatedBy          int                   `json:"updated_by,omitempty"`
	CreatedAt          int64                 `json:"created_at,omitempty"`
	UpdatedAt          int64                 `json:"updated_at,omitempty"`
	ClosestTaskAt      int64                 `json:"closest_task_at,omitempty"`
	IsDeleted          bool                  `json:"is_deleted,omitempty"`
	CustomFieldsValues []custom_fields.Value `json:"custom_fields_values,omitempty"`
	AccountID          int                   `json:"account_id,omitempty"`
	Tags               []Tag                 `json:"tags,omitempty"`
	Embedded           *Embedded             `json:"_embedded,omitempty"`
}

// Tag представляет тег контакта
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Embedded содержит связанные с контактом сущности
type Embedded struct {
	Companies []Company `json:"companies,omitempty"`
	Tags      []Tag     `json:"tags,omitempty"`
}

// Company представляет компанию, связанную с контактом
type Company struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

// WithOption определяет связанные сущности, которые нужно получить вместе с контактом
type WithOption string

const (
	WithCompanies WithOption = "companies"
)

// Get получает контакт по его ID.
// Параметр withOptions позволяет указать, какие связанные сущности нужно получить вместе с контактом.
func Get(ctx context.Context, apiClient *client.Client, contactID int, withOptions ...WithOption) (*Contact, error) {
	// Формируем базовый URL
	baseURL := fmt.Sprintf("%s/api/v4/contacts/%d", apiClient.GetBaseURL(), contactID)

	// Добавляем параметры запроса, если указаны withOptions
	if len(withOptions) > 0 {
		params := url.Values{}
		var withValues []string
		for _, opt := range withOptions {
			withValues = append(withValues, string(opt))
		}
		params.Add("with", strings.Join(withValues, ","))
		baseURL = baseURL + "?" + params.Encode()
	}

	// Создаем запрос
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var contact Contact
	if err := json.NewDecoder(resp.Body).Decode(&contact); err != nil {
		return nil, err
	}

	return &contact, nil
}

// Create создает новый контакт в amoCRM.
func Create(ctx context.Context, apiClient *client.Client, contact *Contact) (*Contact, error) {
	url := apiClient.GetBaseURL() + "/api/v4/contacts"
	contactJSON, err := json.Marshal(contact)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(contactJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var newContact Contact
	if err := json.NewDecoder(resp.Body).Decode(&newContact); err != nil {
		return nil, err
	}

	return &newContact, nil
}

// DeleteResponse представляет ответ от API при удалении контактов
type DeleteResponse struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}

// ListResponse представляет ответ от API при получении списка контактов
type ListResponse struct {
	Page     int `json:"page"`
	PerPage  int `json:"per_page"`
	Total    int `json:"total"`
	Embedded struct {
		Contacts []Contact `json:"contacts"`
	} `json:"_embedded"`
}

// List получает список контактов с возможностью фильтрации и пагинации.
// Параметр withOptions позволяет указать, какие связанные сущности нужно получить вместе с контактами.
func List(ctx context.Context, apiClient *client.Client, page, limit int, withOptions ...WithOption) ([]Contact, error) {
	// Формируем базовый URL
	baseURL := fmt.Sprintf("%s/api/v4/contacts", apiClient.GetBaseURL())

	// Добавляем параметры запроса
	params := url.Values{}
	params.Add("page", fmt.Sprintf("%d", page))
	params.Add("limit", fmt.Sprintf("%d", limit))

	// Добавляем параметр with, если указаны withOptions
	if len(withOptions) > 0 {
		var withValues []string
		for _, opt := range withOptions {
			withValues = append(withValues, string(opt))
		}
		params.Add("with", strings.Join(withValues, ","))
	}

	// Добавляем параметры к URL
	baseURL = baseURL + "?" + params.Encode()

	// Создаем запрос
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var contacts ListResponse
	if err := json.NewDecoder(resp.Body).Decode(&contacts); err != nil {
		return nil, err
	}

	return contacts.Embedded.Contacts, nil
}

// LinkWithCompany связывает контакт с компанией
func LinkWithCompany(ctx context.Context, apiClient *client.Client, contactID, companyID int) error {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/contacts/%d/link", apiClient.GetBaseURL(), contactID)

	// Формируем тело запроса
	type linkRequest struct {
		To []struct {
			EntityID   int    `json:"entity_id"`
			EntityType string `json:"entity_type"`
		} `json:"to"`
	}

	reqBody := linkRequest{
		To: []struct {
			EntityID   int    `json:"entity_id"`
			EntityType string `json:"entity_type"`
		}{
			{
				EntityID:   companyID,
				EntityType: "companies",
			},
		},
	}

	// Преобразуем тело запроса в JSON
	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	// Создаем запрос
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqJSON))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
