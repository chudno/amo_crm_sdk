// Пакет deals предоставляет методы для взаимодействия с сущностями "Сделки" в API amoCRM.
package deals

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/chudno/amo_crm_sdk/client"
	"github.com/chudno/amo_crm_sdk/utils/custom_fields"
	"net/http"
	"net/url"
	"strings"
)

// Deal представляет собой структуру сделки в amoCRM.
type Deal struct {
	ID                 int                              `json:"id"`
	Name               string                           `json:"name"`
	Value              int                              `json:"value"`
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
	Embedded           *DealEmbedded                    `json:"_embedded,omitempty"`
}

// DealEmbedded содержит связанные со сделкой сущности
type DealEmbedded struct {
	Contacts  []Contact `json:"contacts,omitempty"`
	Companies []Company `json:"companies,omitempty"`
	Tags      []Tag     `json:"tags,omitempty"`
}

// Tag представляет тег сделки.
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Company представляет компанию, связанную со сделкой.
type Company struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

// Contact представляет контакт, связанный со сделкой.
type Contact struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

// WithOption определяет связанные сущности, которые нужно получить вместе со сделкой
type WithOption string

const (
	WithContacts  WithOption = "contacts"
	WithCompanies WithOption = "companies"
)

// GetDeal получает сделку по её ID.
// Параметр withOptions позволяет указать, какие связанные сущности нужно получить вместе со сделкой.
func GetDeal(apiClient *client.Client, dealID int, withOptions ...WithOption) (*Deal, error) {
	// Формируем базовый URL
	baseURL := fmt.Sprintf("%s/api/v4/deals/%d", apiClient.GetBaseURL(), dealID)

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

	var deal Deal
	if err := json.NewDecoder(resp.Body).Decode(&deal); err != nil {
		return nil, err
	}

	return &deal, nil
}

// CreateDeal создает новую сделку в amoCRM.
func CreateDeal(apiClient *client.Client, deal *Deal) (*Deal, error) {
	url := fmt.Sprintf("%s/api/v4/deals", apiClient.GetBaseURL())
	dealJSON, err := json.Marshal(deal)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dealJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var newDeal Deal
	if err := json.NewDecoder(resp.Body).Decode(&newDeal); err != nil {
		return nil, err
	}

	return &newDeal, nil
}

// UpdateDeal обновляет существующую сделку в amoCRM.
func UpdateDeal(apiClient *client.Client, deal *Deal) (*Deal, error) {
	url := fmt.Sprintf("%s/api/v4/deals/%d", apiClient.GetBaseURL(), deal.ID)
	dealJSON, err := json.Marshal(deal)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(dealJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var updatedDeal Deal
	if err := json.NewDecoder(resp.Body).Decode(&updatedDeal); err != nil {
		return nil, err
	}

	return &updatedDeal, nil
}

// DealsResponse представляет ответ от API при получении списка сделок
type DealsResponse struct {
	Page     int `json:"page"`
	PerPage  int `json:"per_page"`
	Total    int `json:"total"`
	Embedded struct {
		Items []Deal `json:"items"`
	} `json:"_embedded"`
}

// GetDeals получает список сделок с возможностью фильтрации и пагинации.
// Параметр withOptions позволяет указать, какие связанные сущности нужно получить вместе со сделками.
func GetDeals(apiClient *client.Client, page, limit int, filter map[string]string, withOptions ...WithOption) ([]Deal, error) {
	// Формируем базовый URL
	baseURL := fmt.Sprintf("%s/api/v4/deals", apiClient.GetBaseURL())

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

	// Добавляем параметры фильтрации, если они есть
	if len(filter) > 0 {
		for key, value := range filter {
			params.Add(key, value)
		}
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

	var response DealsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Embedded.Items, nil
}
