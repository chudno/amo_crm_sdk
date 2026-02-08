// Package unsorted предоставляет методы для взаимодействия с API неразобранных заявок в amoCRM.
package unsorted

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/chudno/amo_crm_sdk/client"
)

// SourceType определяет тип источника неразобранной заявки
type SourceType string

const (
	// SourceTypeAPI источник - API
	SourceTypeAPI SourceType = "api"
	// SourceTypeForms источник - Формы
	SourceTypeForms SourceType = "forms"
	// SourceTypeSite источник - Сайт
	SourceTypeSite SourceType = "site"
	// SourceTypeSip источник - Телефония
	SourceTypeSip SourceType = "sip"
	// SourceTypeEmail источник - Email
	SourceTypeEmail SourceType = "mail"
	// SourceTypeChats источник - Чаты
	SourceTypeChats SourceType = "chats"
)

// CategoryType определяет категорию неразобранной заявки
type CategoryType string

const (
	// CategoryTypeForms категория - Формы
	CategoryTypeForms CategoryType = "forms"
	// CategoryTypeSite категория - Сайт
	CategoryTypeSite CategoryType = "site"
	// CategoryTypeSip категория - Телефония
	CategoryTypeSip CategoryType = "sip"
	// CategoryTypeEmail категория - Email
	CategoryTypeEmail CategoryType = "mail"
	// CategoryTypeChats категория - Чаты
	CategoryTypeChats CategoryType = "chats"
)

// PipelineType определяет тип воронки для неразобранной заявки
type PipelineType string

const (
	// PipelineTypeLead тип воронки - Сделки
	PipelineTypeLead PipelineType = "lead"
	// PipelineTypeContact тип воронки - Контакты
	PipelineTypeContact PipelineType = "contact"
	// PipelineTypeCustomer тип воронки - Покупатели
	PipelineTypeCustomer PipelineType = "customer"
)

// Base базовая структура для неразобранной заявки
type Base struct {
	UID        string       `json:"uid,omitempty"`
	SourceUID  string       `json:"source_uid,omitempty"`
	CreatedAt  int64        `json:"created_at,omitempty"`
	PipelineID int          `json:"pipeline_id,omitempty"`
	SourceName string       `json:"source_name,omitempty"`
	SourceType SourceType   `json:"source_type"`
	Category   CategoryType `json:"category"`
	MetadataID int64        `json:"metadata_id,omitempty"`
	AccountID  int64        `json:"account_id,omitempty"`
}

// Contact представляет контакт в неразобранной заявке
type Contact struct {
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}

// Company представляет компанию в неразобранной заявке
type Company struct {
	Name string `json:"name"`
}

// Metadata представляет метаданные неразобранной заявки
type Metadata struct {
	IP      string `json:"ip,omitempty"`
	Form    any    `json:"form,omitempty"`
	From    string `json:"from,omitempty"`
	To      string `json:"to,omitempty"`
	Subject string `json:"subject,omitempty"`
	Thread  any    `json:"thread,omitempty"`
	Service string `json:"service,omitempty"`
}

// LeadCreate представляет структуру для создания сделки из неразобранной заявки
type LeadCreate struct {
	Base
	Metadata          Metadata     `json:"metadata,omitempty"`
	Contact           *Contact     `json:"contact,omitempty"`
	Company           *Company     `json:"company,omitempty"`
	LeadName          string       `json:"lead_name,omitempty"`
	StatusID          int          `json:"status_id,omitempty"`
	ResponsibleUserID int          `json:"responsible_user_id,omitempty"`
	Price             int          `json:"price,omitempty"`
	PipelineType      PipelineType `json:"pipeline_type,omitempty"`
}

// ContactCreate представляет структуру для создания контакта из неразобранной заявки
type ContactCreate struct {
	Base
	Metadata          Metadata `json:"metadata,omitempty"`
	Contact           *Contact `json:"contact,omitempty"`
	Company           *Company `json:"company,omitempty"`
	ResponsibleUserID int      `json:"responsible_user_id,omitempty"`
}

// Response представляет ответ от API при работе с неразобранными заявками
type Response struct {
	Links *struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links,omitempty"`
	Embedded *struct {
		Unsorted []Item `json:"unsorted"`
	} `json:"_embedded,omitempty"`
	UID       string `json:"uid,omitempty"`
	AccountID int64  `json:"account_id,omitempty"`
}

// Item представляет элемент неразобранных заявок в списке
type Item struct {
	ID           string       `json:"id"`
	UID          string       `json:"uid"`
	SourceUID    string       `json:"source_uid,omitempty"`
	CreatedAt    int64        `json:"created_at"`
	PipelineID   int          `json:"pipeline_id,omitempty"`
	Category     CategoryType `json:"category"`
	SourceType   SourceType   `json:"source_type"`
	SourceName   string       `json:"source_name,omitempty"`
	PipelineType PipelineType `json:"pipeline_type,omitempty"`
	AccountID    int64        `json:"account_id,omitempty"`
	Embedded     *struct {
		Contacts []struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Links struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"_links"`
		} `json:"contacts,omitempty"`
		Companies []struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Links struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"_links"`
		} `json:"companies,omitempty"`
		Leads []struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Links struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"_links"`
		} `json:"leads,omitempty"`
	} `json:"_embedded,omitempty"`
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
}

// CreateLead создает неразобранную заявку с типом "Сделка"
func CreateLead(ctx context.Context, apiClient *client.Client, lead *LeadCreate) (*Response, error) {
	// Устанавливаем временную метку создания, если не указана
	if lead.CreatedAt == 0 {
		lead.CreatedAt = time.Now().Unix()
	}

	// Устанавливаем тип заявки для сделки, если не указан
	if lead.PipelineType == "" {
		lead.PipelineType = PipelineTypeLead
	}

	leadJSON, err := json.Marshal([]LeadCreate{*lead})
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/v4/leads/unsorted/api", apiClient.GetBaseURL())

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(leadJSON))
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

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

// CreateContact создает неразобранную заявку с типом "Контакт"
func CreateContact(ctx context.Context, apiClient *client.Client, contact *ContactCreate) (*Response, error) {
	// Устанавливаем временную метку создания, если не указана
	if contact.CreatedAt == 0 {
		contact.CreatedAt = time.Now().Unix()
	}

	contactJSON, err := json.Marshal([]ContactCreate{*contact})
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/v4/contacts/unsorted/api", apiClient.GetBaseURL())

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(contactJSON))
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

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

// ListLeads получает список неразобранных заявок с типом "Сделка"
func ListLeads(ctx context.Context, apiClient *client.Client, page, limit int, filter map[string]string) ([]Item, error) {
	baseURL := fmt.Sprintf("%s/api/v4/leads/unsorted", apiClient.GetBaseURL())

	params := url.Values{}
	params.Add("page", strconv.Itoa(page))
	params.Add("limit", strconv.Itoa(limit))

	for key, value := range filter {
		params.Add(key, value)
	}

	baseURL = baseURL + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var response struct {
		Embedded struct {
			Unsorted []Item `json:"unsorted"`
		} `json:"_embedded"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Embedded.Unsorted, nil
}

// ListContacts получает список неразобранных заявок с типом "Контакт"
func ListContacts(ctx context.Context, apiClient *client.Client, page, limit int, filter map[string]string) ([]Item, error) {
	baseURL := fmt.Sprintf("%s/api/v4/contacts/unsorted", apiClient.GetBaseURL())

	params := url.Values{}
	params.Add("page", strconv.Itoa(page))
	params.Add("limit", strconv.Itoa(limit))

	for key, value := range filter {
		params.Add(key, value)
	}

	baseURL = baseURL + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var response struct {
		Embedded struct {
			Unsorted []Item `json:"unsorted"`
		} `json:"_embedded"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Embedded.Unsorted, nil
}

// GetSummary получает сводку по неразобранным заявкам
func GetSummary(ctx context.Context, apiClient *client.Client) (map[string]any, error) {
	url := fmt.Sprintf("%s/api/v4/unsorted/summary", apiClient.GetBaseURL())

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var response map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}

// AcceptLead принимает неразобранную заявку сделки
func AcceptLead(ctx context.Context, apiClient *client.Client, unsortedUID string, statusID, responsibleUserID int) (int, error) {
	url := fmt.Sprintf("%s/api/v4/leads/unsorted/%s/accept", apiClient.GetBaseURL(), unsortedUID)

	requestBody := struct {
		StatusID          int `json:"status_id"`
		ResponsibleUserID int `json:"responsible_user_id"`
	}{
		StatusID:          statusID,
		ResponsibleUserID: responsibleUserID,
	}

	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestJSON))
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return 0, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var response struct {
		Links struct {
			Lead struct {
				ID int `json:"id"`
			} `json:"lead"`
		} `json:"_links"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, err
	}

	return response.Links.Lead.ID, nil
}

// AcceptContact принимает неразобранную заявку контакта
func AcceptContact(ctx context.Context, apiClient *client.Client, unsortedUID string, responsibleUserID int) (int, error) {
	url := fmt.Sprintf("%s/api/v4/contacts/unsorted/%s/accept", apiClient.GetBaseURL(), unsortedUID)

	requestBody := struct {
		ResponsibleUserID int `json:"responsible_user_id"`
	}{
		ResponsibleUserID: responsibleUserID,
	}

	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestJSON))
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return 0, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var response struct {
		Links struct {
			Contact struct {
				ID int `json:"id"`
			} `json:"contact"`
		} `json:"_links"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, err
	}

	return response.Links.Contact.ID, nil
}

// DeclineLead отклоняет неразобранную заявку сделки
func DeclineLead(ctx context.Context, apiClient *client.Client, unsortedUID string) error {
	url := fmt.Sprintf("%s/api/v4/leads/unsorted/%s/decline", apiClient.GetBaseURL(), unsortedUID)

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	return nil
}

// DeclineContact отклоняет неразобранную заявку контакта
func DeclineContact(ctx context.Context, apiClient *client.Client, unsortedUID string) error {
	url := fmt.Sprintf("%s/api/v4/contacts/unsorted/%s/decline", apiClient.GetBaseURL(), unsortedUID)

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	return nil
}

// LinkLeadWithContact связывает неразобранную заявку сделки с контактом
func LinkLeadWithContact(ctx context.Context, apiClient *client.Client, unsortedUID string, contactID int) error {
	url := fmt.Sprintf("%s/api/v4/leads/unsorted/%s/link", apiClient.GetBaseURL(), unsortedUID)

	requestBody := struct {
		ContactID int `json:"contact_id"`
	}{
		ContactID: contactID,
	}

	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestJSON))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	return nil
}

// LinkLeadWithCompany связывает неразобранную заявку сделки с компанией
func LinkLeadWithCompany(ctx context.Context, apiClient *client.Client, unsortedUID string, companyID int) error {
	url := fmt.Sprintf("%s/api/v4/leads/unsorted/%s/link", apiClient.GetBaseURL(), unsortedUID)

	requestBody := struct {
		CompanyID int `json:"company_id"`
	}{
		CompanyID: companyID,
	}

	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestJSON))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	return nil
}

// LinkContactWithCompany связывает неразобранную заявку контакта с компанией
func LinkContactWithCompany(ctx context.Context, apiClient *client.Client, unsortedUID string, companyID int) error {
	url := fmt.Sprintf("%s/api/v4/contacts/unsorted/%s/link", apiClient.GetBaseURL(), unsortedUID)

	requestBody := struct {
		CompanyID int `json:"company_id"`
	}{
		CompanyID: companyID,
	}

	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestJSON))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	return nil
}
