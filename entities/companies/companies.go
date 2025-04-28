// Пакет companies предоставляет методы для взаимодействия с сущностями "Компании" в API amoCRM.
package companies

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/chudno/amo_crm_sdk/client"
	"github.com/chudno/amo_crm_sdk/utils/custom_fields"
)

// Company представляет собой структуру компании в amoCRM.
type Company struct {
	ID                 int                              `json:"id"`
	Name               string                           `json:"name"`
	ResponsibleUserID  int                              `json:"responsible_user_id,omitempty"`
	GroupID            int                              `json:"group_id,omitempty"`
	CreatedBy          int                              `json:"created_by,omitempty"`
	UpdatedBy          int                              `json:"updated_by,omitempty"`
	CreatedAt          int64                            `json:"created_at,omitempty"`
	UpdatedAt          int64                            `json:"updated_at,omitempty"`
	ClosestTaskAt      int64                            `json:"closest_task_at,omitempty"`
	IsDeleted          bool                             `json:"is_deleted,omitempty"`
	CustomFieldsValues []custom_fields.CustomFieldValue `json:"custom_fields_values,omitempty"`
	AccountID          int                              `json:"account_id,omitempty"`
	Tags               []Tag                            `json:"tags,omitempty"`
	Embedded           *CompanyEmbedded                 `json:"_embedded,omitempty"`
}

// Tag представляет тег компании
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// CompanyEmbedded содержит связанные с компанией сущности
type CompanyEmbedded struct {
	Contacts []Contact `json:"contacts,omitempty"`
	Tags     []Tag     `json:"tags,omitempty"`
}

// Contact представляет контакт, связанный с компанией
type Contact struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

// WithOption определяет связанные сущности, которые нужно получить вместе с компанией
type WithOption string

const (
	WithContacts WithOption = "contacts"
)

// GetCompany получает компанию по её ID.
// Параметр withOptions позволяет указать, какие связанные сущности нужно получить вместе с компанией.
func GetCompany(apiClient *client.Client, companyID int, withOptions ...WithOption) (*Company, error) {
	// Формируем базовый URL
	baseURL := fmt.Sprintf("%s/api/v4/companies/%d", apiClient.GetBaseURL(), companyID)

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

	var company Company
	if err := json.NewDecoder(resp.Body).Decode(&company); err != nil {
		return nil, err
	}

	return &company, nil
}

// CreateCompany создает новую компанию в amoCRM.
func CreateCompany(apiClient *client.Client, company *Company) (*Company, error) {
	url := apiClient.GetBaseURL() + "/api/v4/companies"
	companyJSON, err := json.Marshal(company)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(companyJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var newCompany Company
	if err := json.NewDecoder(resp.Body).Decode(&newCompany); err != nil {
		return nil, err
	}

	return &newCompany, nil
}

// UpdateCompany обновляет существующую компанию в amoCRM.
func UpdateCompany(apiClient *client.Client, company *Company) (*Company, error) {
	url := apiClient.GetBaseURL() + "/api/v4/companies/" + fmt.Sprintf("%d", company.ID)
	companyJSON, err := json.Marshal(company)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(companyJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var updatedCompany Company
	if err := json.NewDecoder(resp.Body).Decode(&updatedCompany); err != nil {
		return nil, err
	}

	return &updatedCompany, nil
}

// CompaniesResponse представляет ответ от API при получении списка компаний
type CompaniesResponse struct {
	Page     int `json:"page"`
	PerPage  int `json:"per_page"`
	Total    int `json:"total"`
	Embedded struct {
		Items []Company `json:"items"`
	} `json:"_embedded"`
}

// GetCompanies получает список компаний с возможностью фильтрации и пагинации.
// Параметр withOptions позволяет указать, какие связанные сущности нужно получить вместе с компаниями.
func GetCompanies(apiClient *client.Client, page, limit int, withOptions ...WithOption) ([]Company, error) {
	// Формируем базовый URL
	baseURL := fmt.Sprintf("%s/api/v4/companies", apiClient.GetBaseURL())

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
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var companies CompaniesResponse
	if err := json.NewDecoder(resp.Body).Decode(&companies); err != nil {
		return nil, err
	}

	return companies.Embedded.Items, nil
}
