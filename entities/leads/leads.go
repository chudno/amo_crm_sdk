// Пакет leads предоставляет методы для взаимодействия с сущностями "Лиды" в API amoCRM.
package leads

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/chudno/amo_crm_sdk/client"
	"github.com/chudno/amo_crm_sdk/entities/companies"
	"github.com/chudno/amo_crm_sdk/entities/contacts"
	"github.com/chudno/amo_crm_sdk/utils/custom_fields"
	"net/http"
	"net/url"
	"strings"
)

// Lead представляет собой структуру лида в amoCRM.
type Lead struct {
	ID                 int                              `json:"id"`
	Name               string                           `json:"name"`
	Price              int                              `json:"price"`
	ResponsibleUserID  int                              `json:"responsible_user_id,omitempty"`
	GroupID            int                              `json:"group_id,omitempty"`
	StatusID           int                              `json:"status_id,omitempty"`
	PipelineID         int                              `json:"pipeline_id,omitempty"`
	LossReasonID       int                              `json:"loss_reason_id,omitempty"`
	SourceID           int                              `json:"source_id,omitempty"`
	CreatedBy          int                              `json:"created_by,omitempty"`
	UpdatedBy          int                              `json:"updated_by,omitempty"`
	CreatedAt          int64                            `json:"created_at,omitempty"`
	UpdatedAt          int64                            `json:"updated_at,omitempty"`
	ClosedAt           int64                            `json:"closed_at,omitempty"`
	ClosestTaskAt      int64                            `json:"closest_task_at,omitempty"`
	IsDeleted          bool                             `json:"is_deleted,omitempty"`
	CustomFieldsValues []custom_fields.CustomFieldValue `json:"custom_fields_values,omitempty"`
	Score              int                              `json:"score,omitempty"`
	AccountID          int                              `json:"account_id,omitempty"`
	Tags               []Tag                            `json:"tags,omitempty"`
	Embedded           *LeadEmbedded                    `json:"_embedded,omitempty"`
}

// Tag представляет тег сделки
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// LeadEmbedded содержит связанные с лидом сущности
type LeadEmbedded struct {
	Contacts  []contacts.Contact  `json:"contacts,omitempty"`
	Companies []companies.Company `json:"companies,omitempty"`
	Tags      []Tag               `json:"tags,omitempty"`
}

// WithOption определяет связанные сущности, которые нужно получить вместе с лидом
type WithOption string

const (
	WithContacts  WithOption = "contacts"
	WithCompanies WithOption = "companies"
)

// GetLead получает лид по его ID.
// Параметр withOptions позволяет указать, какие связанные сущности нужно получить вместе с лидом.
func GetLead(apiClient *client.Client, leadID int, withOptions ...WithOption) (*Lead, error) {
	// Формируем базовый URL
	baseURL := fmt.Sprintf("%s/api/v4/leads/%d", apiClient.GetBaseURL(), leadID)

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
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var lead Lead
	if err := json.NewDecoder(resp.Body).Decode(&lead); err != nil {
		return nil, err
	}

	return &lead, nil
}

// CreateLead создает новый лид в amoCRM.
func CreateLead(apiClient *client.Client, lead *Lead) (*Lead, error) {
	url := fmt.Sprintf("%s/api/v4/leads", apiClient.GetBaseURL())

	leadData, err := json.Marshal([]*Lead{lead})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(leadData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response struct {
		Embedded struct {
			Leads []*Lead `json:"leads"`
		} `json:"_embedded"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if len(response.Embedded.Leads) == 0 {
		return nil, fmt.Errorf("не удалось создать лид")
	}

	return response.Embedded.Leads[0], nil
}

// UpdateLead обновляет существующий лид в amoCRM.
func UpdateLead(apiClient *client.Client, lead *Lead) (*Lead, error) {
	if lead.ID == 0 {
		return nil, fmt.Errorf("ID лида не указан")
	}

	url := fmt.Sprintf("%s/api/v4/leads/%d", apiClient.GetBaseURL(), lead.ID)

	leadData, err := json.Marshal(lead)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(leadData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var updatedLead Lead
	if err := json.NewDecoder(resp.Body).Decode(&updatedLead); err != nil {
		return nil, err
	}

	return &updatedLead, nil
}

// ListLeads получает список лидов с возможностью фильтрации и пагинации.
func ListLeads(apiClient *client.Client, limit int, page int, filter map[string]interface{}) ([]*Lead, error) {
	baseURL := fmt.Sprintf("%s/api/v4/leads", apiClient.GetBaseURL())

	// Добавляем параметры запроса
	params := url.Values{}
	params.Add("limit", fmt.Sprintf("%d", limit))
	params.Add("page", fmt.Sprintf("%d", page))

	// Если указаны фильтры, добавляем их в запрос
	if len(filter) > 0 {
		filterData, err := json.Marshal(filter)
		if err != nil {
			return nil, err
		}
		params.Add("filter", string(filterData))
	}

	url := baseURL + "?" + params.Encode()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var response struct {
		Embedded struct {
			Leads []*Lead `json:"leads"`
		} `json:"_embedded"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Embedded.Leads, nil
}

// DeleteLead удаляет лид по его ID.
func DeleteLead(apiClient *client.Client, leadID int) error {
	url := fmt.Sprintf("%s/api/v4/leads/%d", apiClient.GetBaseURL(), leadID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
