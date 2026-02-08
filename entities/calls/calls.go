// Package calls предоставляет методы для работы со звонками в amoCRM.
package calls

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

// Числовые коды статусов звонков для API v4
const (
	// CallStatusCodeLeftMessage — оставлено сообщение
	CallStatusCodeLeftMessage = 1
	// CallStatusCodeInterrupted — разговор прерван
	CallStatusCodeInterrupted = 2
	// CallStatusCodeNoAnswer — нет ответа
	CallStatusCodeNoAnswer = 3
	// CallStatusCodeSuccess — успешный звонок
	CallStatusCodeSuccess = 4
	// CallStatusCodeWrongNumber — неверный номер
	CallStatusCodeWrongNumber = 5
	// CallStatusCodeBusy — занято
	CallStatusCodeBusy = 6
	// CallStatusCodeVoicemail — голосовая почта
	CallStatusCodeVoicemail = 7
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
	Status            CallStatus    `json:"status,omitempty"`
	CallStatusCode    int           `json:"call_status,omitempty"`
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
	Embedded          *Embedded     `json:"_embedded,omitempty"`
	Links             *Links        `json:"_links,omitempty"`
	EntityType        *EntityType   `json:"entity_type,omitempty"`
	EntityID          int           `json:"entity_id,omitempty"`
}

// Voice содержит информацию о голосовом сообщении
type Voice struct {
	URL              string `json:"url,omitempty"`
	TranscriptionURL string `json:"transcription_url,omitempty"`
}

// Embedded содержит вложенные сущности
type Embedded struct {
	Tags []struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Color string `json:"color,omitempty"`
	} `json:"tags,omitempty"`
}

// Links содержит ссылки
type Links struct {
	Self struct {
		Href string `json:"href"`
	} `json:"self"`
}

// Add добавляет новый звонок в amoCRM.
func Add(ctx context.Context, apiClient *client.Client, call *Call) (*Call, error) {
	if call.Direction == "" {
		return nil, fmt.Errorf("direction is required")
	}
	if call.Phone == "" {
		return nil, fmt.Errorf("phone is required")
	}

	if call.CallStatusCode == 0 {
		switch call.Status {
		case CallStatusSuccess:
			call.CallStatusCode = CallStatusCodeSuccess
		case CallStatusMissed:
			call.CallStatusCode = CallStatusCodeNoAnswer
		case CallStatusVoicemail:
			call.CallStatusCode = CallStatusCodeVoicemail
		case CallStatusHungup:
			call.CallStatusCode = CallStatusCodeInterrupted
		case CallStatusBusy:
			call.CallStatusCode = CallStatusCodeBusy
		}
	}

	if call.CreatedAt == 0 {
		call.CreatedAt = time.Now().Unix()
	}

	if call.Uniq == "" {
		call.Uniq = fmt.Sprintf("%d-%d", call.CreatedAt, time.Now().UnixNano())
	}

	// API не принимает строковое поле "status", только числовой "call_status"
	// Очищаем status перед отправкой (omitempty не отправит пустую строку)
	call.Status = ""

	requestURL := fmt.Sprintf("%s/api/v4/calls", apiClient.GetBaseURL())

	callJSON, err := json.Marshal([]*Call{call})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", requestURL, bytes.NewBuffer(callJSON))
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
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("неожиданный статус-код: %d, ответ: %s", resp.StatusCode, body)
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
