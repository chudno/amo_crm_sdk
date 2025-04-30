// Пакет calls предоставляет методы для работы со звонками в amoCRM.
package calls

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

// CallDirection определяет направление звонка
type CallDirection string

const (
	// CallDirectionIncoming входящий звонок
	CallDirectionIncoming CallDirection = "inbound"
	// CallDirectionOutgoing исходящий звонок
	CallDirectionOutgoing CallDirection = "outbound"
)

// CallStatus определяет статус звонка
type CallStatus string

const (
	// CallStatusSuccess успешный звонок
	CallStatusSuccess CallStatus = "success"
	// CallStatusMissed пропущенный звонок
	CallStatusMissed CallStatus = "missed"
	// CallStatusVoicemail голосовая почта
	CallStatusVoicemail CallStatus = "voicemail"
	// CallStatusHungup сброшенный звонок
	CallStatusHungup CallStatus = "hung_up"
	// CallStatusBusy занято
	CallStatusBusy CallStatus = "busy"
)

// EntityType определяет тип сущности, с которой связан звонок
type EntityType string

const (
	// EntityTypeLead тип сущности - Сделка
	EntityTypeLead EntityType = "leads"
	// EntityTypeContact тип сущности - Контакт
	EntityTypeContact EntityType = "contacts"
	// EntityTypeCompany тип сущности - Компания
	EntityTypeCompany EntityType = "companies"
	// EntityTypeCustomers тип сущности - Покупатель
	EntityTypeCustomers EntityType = "customers"
)

// Call представляет структуру звонка в amoCRM.
type Call struct {
	ID                int           `json:"id,omitempty"`
	Direction         CallDirection `json:"direction"`
	Status            CallStatus    `json:"status"`
	ResponsibleUserID int           `json:"responsible_user_id,omitempty"`
	CreatedBy         int           `json:"created_by,omitempty"`
	UpdatedBy         int           `json:"updated_by,omitempty"`
	CreatedAt         int64         `json:"created_at,omitempty"`
	UpdatedAt         int64         `json:"updated_at,omitempty"`
	AccountID         int64         `json:"account_id,omitempty"`
	Uniq              string        `json:"uniq,omitempty"`
	Duration          int           `json:"duration,omitempty"`
	Source            string        `json:"source,omitempty"`
	CallResult        string        `json:"call_result,omitempty"`
	Link              string        `json:"link,omitempty"`
	ServiceCode       string        `json:"service_code,omitempty"`
	Phone             string        `json:"phone,omitempty"`
	APIID             int           `json:"api_id,omitempty"`
	ManagerName       string        `json:"manager_name,omitempty"`
	ManagerEmail      string        `json:"manager_email,omitempty"`
	ManagerPhone      string        `json:"manager_phone,omitempty"`
	ManagerICQ        string        `json:"manager_icq,omitempty"`
	ContactID         int           `json:"contact_id,omitempty"`
	LeadID            int           `json:"lead_id,omitempty"`
	CompanyID         int           `json:"company_id,omitempty"`
	SourceName        string        `json:"source_name,omitempty"`
	SourceUID         string        `json:"source_uid,omitempty"`
	IsCallbackCall    bool          `json:"is_callback_call,omitempty"`
	IsRinging         bool          `json:"is_ringing,omitempty"`
	Voice             *Voice        `json:"voice,omitempty"`
	CallStartTime     string        `json:"call_start_time,omitempty"`
	CallEndTime       string        `json:"call_end_time,omitempty"`
	Version           int           `json:"version,omitempty"`
	Embedded          *CallEmbedded `json:"_embedded,omitempty"`
	Links             *CallLinks    `json:"_links,omitempty"`
	EntityType        *EntityType   `json:"entity_type,omitempty"`
	EntityID          int           `json:"entity_id,omitempty"`
}

// Voice содержит информацию о голосовом сообщении
type Voice struct {
	URL              string `json:"url,omitempty"`
	TranscriptionURL string `json:"transcription_url,omitempty"`
}

// CallEmbedded содержит вложенные сущности
type CallEmbedded struct {
	Tags []struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Color string `json:"color,omitempty"`
	} `json:"tags,omitempty"`
}

// CallLinks содержит ссылки
type CallLinks struct {
	Self struct {
		Href string `json:"href"`
	} `json:"self"`
}

// CallsResponse представляет ответ от API при получении списка звонков
type CallsResponse struct {
	Page     int `json:"page"`
	PerPage  int `json:"per_page"`
	Total    int `json:"total"`
	Embedded struct {
		Calls []Call `json:"calls"`
	} `json:"_embedded"`
}

// WithOption представляет опцию для запроса с добавлением связанных сущностей
type WithOption string

const (
	// WithTags - получить теги звонков
	WithTags WithOption = "tags"
)

// AddCall добавляет новый звонок в amoCRM.
func AddCall(apiClient *client.Client, call *Call) (*Call, error) {
	// Проверяем обязательные поля
	if call.Direction == "" {
		return nil, fmt.Errorf("direction is required")
	}
	if call.Status == "" {
		return nil, fmt.Errorf("status is required")
	}
	if call.Phone == "" {
		return nil, fmt.Errorf("phone is required")
	}

	// Устанавливаем время создания, если не указано
	if call.CreatedAt == 0 {
		call.CreatedAt = time.Now().Unix()
	}

	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/calls", apiClient.GetBaseURL())

	// Преобразуем структуру звонка в JSON
	callJSON, err := json.Marshal([]*Call{call})
	if err != nil {
		return nil, err
	}

	// Создаем запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(callJSON))
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

	var response struct {
		Embedded struct {
			Calls []Call `json:"calls"`
		} `json:"_embedded"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if len(response.Embedded.Calls) == 0 {
		return nil, fmt.Errorf("не удалось создать звонок")
	}

	return &response.Embedded.Calls[0], nil
}

// GetCalls получает список звонков с возможностью фильтрации и пагинации.
func GetCalls(apiClient *client.Client, page, limit int, filter map[string]string, withOptions ...WithOption) ([]Call, error) {
	// Формируем URL для запроса
	baseURL := fmt.Sprintf("%s/api/v4/calls", apiClient.GetBaseURL())

	// Добавляем параметры запроса
	params := url.Values{}
	params.Add("page", strconv.Itoa(page))
	params.Add("limit", strconv.Itoa(limit))

	// Добавляем фильтры
	for key, value := range filter {
		params.Add(key, value)
	}

	// Добавляем параметр with, если указаны withOptions
	if len(withOptions) > 0 {
		var withValues []string
		for _, opt := range withOptions {
			withValues = append(withValues, string(opt))
		}
		params.Add("with", stringsJoin(withValues, ","))
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

	var callsResponse CallsResponse
	if err := json.NewDecoder(resp.Body).Decode(&callsResponse); err != nil {
		return nil, err
	}

	return callsResponse.Embedded.Calls, nil
}

// stringsJoin объединяет срез строк с указанным разделителем
func stringsJoin(strings []string, sep string) string {
	if len(strings) == 0 {
		return ""
	}

	result := strings[0]
	for i := 1; i < len(strings); i++ {
		result += sep + strings[i]
	}

	return result
}

// GetCall получает информацию о конкретном звонке по его ID.
func GetCall(apiClient *client.Client, callID int, withOptions ...WithOption) (*Call, error) {
	// Формируем URL для запроса
	baseURL := fmt.Sprintf("%s/api/v4/calls/%d", apiClient.GetBaseURL(), callID)

	// Добавляем параметр with, если указаны withOptions
	if len(withOptions) > 0 {
		params := url.Values{}
		var withValues []string
		for _, opt := range withOptions {
			withValues = append(withValues, string(opt))
		}
		params.Add("with", stringsJoin(withValues, ","))
		baseURL = baseURL + "?" + params.Encode()
	}

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

	var call Call
	if err := json.NewDecoder(resp.Body).Decode(&call); err != nil {
		return nil, err
	}

	return &call, nil
}

// UpdateCall обновляет информацию о звонке.
func UpdateCall(apiClient *client.Client, call *Call) (*Call, error) {
	if call.ID == 0 {
		return nil, fmt.Errorf("ID звонка не может быть пустым")
	}

	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/calls/%d", apiClient.GetBaseURL(), call.ID)

	// Преобразуем структуру звонка в JSON
	callJSON, err := json.Marshal(call)
	if err != nil {
		return nil, err
	}

	// Создаем запрос
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(callJSON))
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
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var updatedCall Call
	if err := json.NewDecoder(resp.Body).Decode(&updatedCall); err != nil {
		return nil, err
	}

	return &updatedCall, nil
}

// DeleteCall удаляет звонок по его ID.
func DeleteCall(apiClient *client.Client, callID int) error {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/calls/%d", apiClient.GetBaseURL(), callID)

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

// LinkCallWithEntity связывает звонок с сущностью (сделкой, контактом, компанией).
func LinkCallWithEntity(apiClient *client.Client, callID int, entityType EntityType, entityID int) error {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/calls/%d/link", apiClient.GetBaseURL(), callID)

	// Создаем структуру для запроса
	requestBody := struct {
		EntityType EntityType `json:"entity_type"`
		EntityID   int        `json:"entity_id"`
	}{
		EntityType: entityType,
		EntityID:   entityID,
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
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	return nil
}

// UnlinkCallFromEntity отвязывает звонок от сущности.
func UnlinkCallFromEntity(apiClient *client.Client, callID int, entityType EntityType, entityID int) error {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/calls/%d/unlink", apiClient.GetBaseURL(), callID)

	// Создаем структуру для запроса
	requestBody := struct {
		EntityType EntityType `json:"entity_type"`
		EntityID   int        `json:"entity_id"`
	}{
		EntityType: entityType,
		EntityID:   entityID,
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
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	return nil
}
