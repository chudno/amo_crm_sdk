// Пакет unsorted предоставляет методы для взаимодействия с API неразобранных заявок в amoCRM.
package unsorted

import (
	"bytes"
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

// UnsortedBase базовая структура для неразобранной заявки
type UnsortedBase struct {
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

// UnsortedContact представляет контакт в неразобранной заявке
type UnsortedContact struct {
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}

// UnsortedCompany представляет компанию в неразобранной заявке
type UnsortedCompany struct {
	Name string `json:"name"`
}

// UnsortedMetadata представляет метаданные неразобранной заявки
type UnsortedMetadata struct {
	IP      string      `json:"ip,omitempty"`
	Form    interface{} `json:"form,omitempty"`
	From    string      `json:"from,omitempty"`
	To      string      `json:"to,omitempty"`
	Subject string      `json:"subject,omitempty"`
	Thread  interface{} `json:"thread,omitempty"`
	Service string      `json:"service,omitempty"`
}

// UnsortedLeadCreate представляет структуру для создания сделки из неразобранной заявки
type UnsortedLeadCreate struct {
	UnsortedBase
	Metadata          UnsortedMetadata `json:"metadata,omitempty"`
	Contact           *UnsortedContact `json:"contact,omitempty"`
	Company           *UnsortedCompany `json:"company,omitempty"`
	LeadName          string           `json:"lead_name,omitempty"`
	StatusID          int              `json:"status_id,omitempty"`
	ResponsibleUserID int              `json:"responsible_user_id,omitempty"`
	Price             int              `json:"price,omitempty"`
	PipelineType      PipelineType     `json:"pipeline_type,omitempty"`
}

// UnsortedContactCreate представляет структуру для создания контакта из неразобранной заявки
type UnsortedContactCreate struct {
	UnsortedBase
	Metadata          UnsortedMetadata `json:"metadata,omitempty"`
	Contact           *UnsortedContact `json:"contact,omitempty"`
	Company           *UnsortedCompany `json:"company,omitempty"`
	ResponsibleUserID int              `json:"responsible_user_id,omitempty"`
}

// UnsortedResponse представляет ответ от API при работе с неразобранными заявками
type UnsortedResponse struct {
	Links *struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links,omitempty"`
	Embedded *struct {
		Unsorted []UnsortedItem `json:"unsorted"`
	} `json:"_embedded,omitempty"`
	UID       string `json:"uid,omitempty"`
	AccountID int64  `json:"account_id,omitempty"`
}

// UnsortedItem представляет элемент неразобранных заявок в списке
type UnsortedItem struct {
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

// CreateUnsortedLead создает неразобранную заявку с типом "Сделка"
func CreateUnsortedLead(apiClient *client.Client, lead *UnsortedLeadCreate) (*UnsortedResponse, error) {
	// Устанавливаем временную метку создания, если не указана
	if lead.CreatedAt == 0 {
		lead.CreatedAt = time.Now().Unix()
	}

	// Устанавливаем тип заявки для сделки, если не указан
	if lead.PipelineType == "" {
		lead.PipelineType = PipelineTypeLead
	}

	// Преобразуем заявку в JSON
	leadJSON, err := json.Marshal([]UnsortedLeadCreate{*lead})
	if err != nil {
		return nil, err
	}

	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/leads/unsorted/api", apiClient.GetBaseURL())

	// Создаем запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(leadJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var response UnsortedResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

// CreateUnsortedContact создает неразобранную заявку с типом "Контакт"
func CreateUnsortedContact(apiClient *client.Client, contact *UnsortedContactCreate) (*UnsortedResponse, error) {
	// Устанавливаем временную метку создания, если не указана
	if contact.CreatedAt == 0 {
		contact.CreatedAt = time.Now().Unix()
	}

	// Преобразуем заявку в JSON
	contactJSON, err := json.Marshal([]UnsortedContactCreate{*contact})
	if err != nil {
		return nil, err
	}

	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/contacts/unsorted/api", apiClient.GetBaseURL())

	// Создаем запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(contactJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var response UnsortedResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetUnsortedLeads получает список неразобранных заявок с типом "Сделка"
func GetUnsortedLeads(apiClient *client.Client, page, limit int, filter map[string]string) ([]UnsortedItem, error) {
	// Формируем URL для запроса
	baseURL := fmt.Sprintf("%s/api/v4/leads/unsorted", apiClient.GetBaseURL())

	// Формируем параметры запроса
	params := url.Values{}
	params.Add("page", strconv.Itoa(page))
	params.Add("limit", strconv.Itoa(limit))

	// Добавляем фильтры, если они указаны
	if filter != nil {
		for key, value := range filter {
			params.Add(key, value)
		}
	}

	// Добавляем параметры к URL
	baseURL = baseURL + "?" + params.Encode()

	// Создаем запрос
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, err
	}

	// Выполняем запрос
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
			Unsorted []UnsortedItem `json:"unsorted"`
		} `json:"_embedded"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Embedded.Unsorted, nil
}

// GetUnsortedContacts получает список неразобранных заявок с типом "Контакт"
func GetUnsortedContacts(apiClient *client.Client, page, limit int, filter map[string]string) ([]UnsortedItem, error) {
	// Формируем URL для запроса
	baseURL := fmt.Sprintf("%s/api/v4/contacts/unsorted", apiClient.GetBaseURL())

	// Формируем параметры запроса
	params := url.Values{}
	params.Add("page", strconv.Itoa(page))
	params.Add("limit", strconv.Itoa(limit))

	// Добавляем фильтры, если они указаны
	if filter != nil {
		for key, value := range filter {
			params.Add(key, value)
		}
	}

	// Добавляем параметры к URL
	baseURL = baseURL + "?" + params.Encode()

	// Создаем запрос
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, err
	}

	// Выполняем запрос
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
			Unsorted []UnsortedItem `json:"unsorted"`
		} `json:"_embedded"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Embedded.Unsorted, nil
}

// GetUnsortedSummary получает сводку по неразобранным заявкам
func GetUnsortedSummary(apiClient *client.Client) (map[string]interface{}, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/unsorted/summary", apiClient.GetBaseURL())

	// Создаем запрос
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Выполняем запрос
	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}

// AcceptUnsortedLead принимает неразобранную заявку сделки
func AcceptUnsortedLead(apiClient *client.Client, unsortedUID string, statusID, responsibleUserID int) (int, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/leads/unsorted/%s/accept", apiClient.GetBaseURL(), unsortedUID)

	// Создаем структуру для запроса
	requestBody := struct {
		StatusID          int `json:"status_id"`
		ResponsibleUserID int `json:"responsible_user_id"`
	}{
		StatusID:          statusID,
		ResponsibleUserID: responsibleUserID,
	}

	// Преобразуем структуру в JSON
	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return 0, err
	}

	// Создаем запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestJSON))
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
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

// AcceptUnsortedContact принимает неразобранную заявку контакта
func AcceptUnsortedContact(apiClient *client.Client, unsortedUID string, responsibleUserID int) (int, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/contacts/unsorted/%s/accept", apiClient.GetBaseURL(), unsortedUID)

	// Создаем структуру для запроса
	requestBody := struct {
		ResponsibleUserID int `json:"responsible_user_id"`
	}{
		ResponsibleUserID: responsibleUserID,
	}

	// Преобразуем структуру в JSON
	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return 0, err
	}

	// Создаем запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestJSON))
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
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

// DeclineUnsortedLead отклоняет неразобранную заявку сделки
func DeclineUnsortedLead(apiClient *client.Client, unsortedUID string) error {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/leads/unsorted/%s/decline", apiClient.GetBaseURL(), unsortedUID)

	// Создаем запрос
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	// Выполняем запрос
	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	return nil
}

// DeclineUnsortedContact отклоняет неразобранную заявку контакта
func DeclineUnsortedContact(apiClient *client.Client, unsortedUID string) error {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/contacts/unsorted/%s/decline", apiClient.GetBaseURL(), unsortedUID)

	// Создаем запрос
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	// Выполняем запрос
	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	return nil
}

// LinkUnsortedLeadWithContact связывает неразобранную заявку сделки с контактом
func LinkUnsortedLeadWithContact(apiClient *client.Client, unsortedUID string, contactID int) error {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/leads/unsorted/%s/link", apiClient.GetBaseURL(), unsortedUID)

	// Создаем структуру для запроса
	requestBody := struct {
		ContactID int `json:"contact_id"`
	}{
		ContactID: contactID,
	}

	// Преобразуем структуру в JSON
	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	// Создаем запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestJSON))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	return nil
}

// LinkUnsortedLeadWithCompany связывает неразобранную заявку сделки с компанией
func LinkUnsortedLeadWithCompany(apiClient *client.Client, unsortedUID string, companyID int) error {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/leads/unsorted/%s/link", apiClient.GetBaseURL(), unsortedUID)

	// Создаем структуру для запроса
	requestBody := struct {
		CompanyID int `json:"company_id"`
	}{
		CompanyID: companyID,
	}

	// Преобразуем структуру в JSON
	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	// Создаем запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestJSON))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	return nil
}

// LinkUnsortedContactWithCompany связывает неразобранную заявку контакта с компанией
func LinkUnsortedContactWithCompany(apiClient *client.Client, unsortedUID string, companyID int) error {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/contacts/unsorted/%s/link", apiClient.GetBaseURL(), unsortedUID)

	// Создаем структуру для запроса
	requestBody := struct {
		CompanyID int `json:"company_id"`
	}{
		CompanyID: companyID,
	}

	// Преобразуем структуру в JSON
	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	// Создаем запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestJSON))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	return nil
}
