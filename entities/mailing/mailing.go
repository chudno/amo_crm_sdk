// Package mailing предоставляет методы для работы с email-рассылками в amoCRM.
package mailing

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

// Requester - интерфейс для выполнения HTTP-запросов, используется для тестирования.
type Requester interface {
	DoRequest(ctx context.Context, req *http.Request) (*http.Response, error)
	GetBaseURL() string
}

// MailingStatus представляет статус рассылки.
type MailingStatus string

const (
	// MailingStatusDraft - черновик
	MailingStatusDraft MailingStatus = "draft"
	// MailingStatusScheduled - запланирована
	MailingStatusScheduled MailingStatus = "scheduled"
	// MailingStatusActive - активна
	MailingStatusActive MailingStatus = "active"
	// MailingStatusPaused - приостановлена
	MailingStatusPaused MailingStatus = "paused"
	// MailingStatusCompleted - завершена
	MailingStatusCompleted MailingStatus = "completed"
	// MailingStatusStopped - остановлена
	MailingStatusStopped MailingStatus = "stopped"
)

// MailingFrequency представляет частоту отправки рассылки.
type MailingFrequency string

const (
	// MailingFrequencyOnce - однократно
	MailingFrequencyOnce MailingFrequency = "once"
	// MailingFrequencyDaily - ежедневно
	MailingFrequencyDaily MailingFrequency = "daily"
	// MailingFrequencyWeekly - еженедельно
	MailingFrequencyWeekly MailingFrequency = "weekly"
	// MailingFrequencyMonthly - ежемесячно
	MailingFrequencyMonthly MailingFrequency = "monthly"
)

// Mailing представляет email-рассылку в amoCRM.
type Mailing struct {
	ID               int               `json:"id,omitempty"`
	Name             string            `json:"name"`
	Status           MailingStatus     `json:"status,omitempty"`
	Subject          string            `json:"subject"`
	Template         *Template         `json:"template,omitempty"`
	Frequency        MailingFrequency  `json:"frequency,omitempty"`
	SendAt           *time.Time        `json:"send_at,omitempty"`
	CreatedAt        int64             `json:"created_at,omitempty"`
	UpdatedAt        int64             `json:"updated_at,omitempty"`
	CreatedBy        int               `json:"created_by,omitempty"`
	UpdatedBy        int               `json:"updated_by,omitempty"`
	SegmentIDs       []int             `json:"segment_ids,omitempty"`
	SegmentFilters   []SegmentFilter   `json:"segment_filters,omitempty"`
	SelectedContacts []int             `json:"selected_contacts,omitempty"`
	ExcludedContacts []int             `json:"excluded_contacts,omitempty"`
	Stats            *Stats            `json:"stats,omitempty"`
	AccountID        int               `json:"account_id,omitempty"`
	FromEmail        string            `json:"from_email,omitempty"`
	FromName         string            `json:"from_name,omitempty"`
	ReplyToEmail     string            `json:"reply_to_email,omitempty"`
	Settings         map[string]string `json:"settings,omitempty"`
}

// Template представляет шаблон рассылки.
type Template struct {
	ID      int    `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Content string `json:"content,omitempty"`
	HTML    string `json:"html,omitempty"`
	Type    string `json:"type,omitempty"`
}

// SegmentFilter представляет фильтр для сегмента контактов.
type SegmentFilter struct {
	Type      string `json:"type"`
	Logic     string `json:"logic"`
	Condition string `json:"condition,omitempty"`
	Value     string `json:"value,omitempty"`
}

// Stats представляет статистику рассылки.
type Stats struct {
	TotalRecipients int `json:"total_recipients"`
	Delivered       int `json:"delivered"`
	Opened          int `json:"opened"`
	Clicked         int `json:"clicked"`
	Bounced         int `json:"bounced"`
	Unsubscribed    int `json:"unsubscribed"`
	Complaints      int `json:"complaints"`
}

// WithOption функциональный параметр для настройки запроса.
type WithOption func(params map[string]string)

// WithFilter добавляет фильтры при получении списка рассылок.
func WithFilter(filter map[string]string) WithOption {
	return func(params map[string]string) {
		for k, v := range filter {
			params[k] = v
		}
	}
}

// WithStatus фильтрует рассылки по статусу.
func WithStatus(status MailingStatus) WithOption {
	return func(params map[string]string) {
		params["filter[status]"] = string(status)
	}
}

// WithDateFrom фильтрует рассылки по дате создания "от".
func WithDateFrom(from time.Time) WithOption {
	return func(params map[string]string) {
		params["filter[created_at][from]"] = strconv.FormatInt(from.Unix(), 10)
	}
}

// WithDateTo фильтрует рассылки по дате создания "до".
func WithDateTo(to time.Time) WithOption {
	return func(params map[string]string) {
		params["filter[created_at][to]"] = strconv.FormatInt(to.Unix(), 10)
	}
}

// List получает список рассылок с поддержкой фильтрации и пагинации.
//
// Пример использования:
//
//	filter := map[string]string{
//		"filter[status]": "active",
//	}
//	mailings, err := mailing.List(ctx, apiClient, 1, 50, mailing.WithFilter(filter))
func List(ctx context.Context, apiClient *client.Client, page, limit int, options ...WithOption) ([]Mailing, error) {
	return ListWithRequester(ctx, apiClient, page, limit, options...)
}

// ListWithRequester получает список рассылок с использованием интерфейса Requester.
func ListWithRequester(ctx context.Context, requester Requester, page, limit int, options ...WithOption) ([]Mailing, error) {
	baseURL := fmt.Sprintf("%s/api/v4/mailings", requester.GetBaseURL())

	params := map[string]string{
		"page":  strconv.Itoa(page),
		"limit": strconv.Itoa(limit),
	}

	for _, option := range options {
		option(params)
	}

	queryParams := url.Values{}
	for key, value := range params {
		queryParams.Add(key, value)
	}
	requestURL := fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := requester.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var response struct {
		Embedded struct {
			Mailings []Mailing `json:"mailings"`
		} `json:"_embedded"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Embedded.Mailings, nil
}

// Get получает информацию о конкретной рассылке по ID.
//
// Пример использования:
//
//	mailingInfo, err := mailing.Get(ctx, apiClient, 123)
func Get(ctx context.Context, apiClient *client.Client, id int) (*Mailing, error) {
	return GetWithRequester(ctx, apiClient, id)
}

// GetWithRequester получает информацию о конкретной рассылке с использованием интерфейса Requester.
func GetWithRequester(ctx context.Context, requester Requester, id int) (*Mailing, error) {
	url := fmt.Sprintf("%s/api/v4/mailings/%d", requester.GetBaseURL(), id)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := requester.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var mailingInfo Mailing
	if err := json.NewDecoder(resp.Body).Decode(&mailingInfo); err != nil {
		return nil, err
	}

	return &mailingInfo, nil
}

// Create создает новую рассылку.
//
// Пример использования:
//
//	newMailing := &mailing.Mailing{
//		Name:     "Новая рассылка",
//		Subject:  "Важная информация",
//		Frequency: mailing.MailingFrequencyOnce,
//	}
//	createdMailing, err := mailing.Create(ctx, apiClient, newMailing)
func Create(ctx context.Context, apiClient *client.Client, mailingData *Mailing) (*Mailing, error) {
	return CreateWithRequester(ctx, apiClient, mailingData)
}

// CreateWithRequester создает новую рассылку с использованием интерфейса Requester.
func CreateWithRequester(ctx context.Context, requester Requester, mailingData *Mailing) (*Mailing, error) {
	url := fmt.Sprintf("%s/api/v4/mailings", requester.GetBaseURL())

	data, err := json.Marshal(mailingData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := requester.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var createdMailing Mailing
	if err := json.NewDecoder(resp.Body).Decode(&createdMailing); err != nil {
		return nil, err
	}

	return &createdMailing, nil
}

// Update обновляет существующую рассылку.
//
// Пример использования:
//
//	mailingUpdate := &mailing.Mailing{
//		ID:       123,
//		Name:     "Обновленная рассылка",
//		Subject:  "Новая тема рассылки",
//	}
//	updatedMailing, err := mailing.Update(ctx, apiClient, mailingUpdate)
func Update(ctx context.Context, apiClient *client.Client, mailingData *Mailing) (*Mailing, error) {
	return UpdateWithRequester(ctx, apiClient, mailingData)
}

// UpdateWithRequester обновляет существующую рассылку с использованием интерфейса Requester.
func UpdateWithRequester(ctx context.Context, requester Requester, mailingData *Mailing) (*Mailing, error) {
	if mailingData.ID == 0 {
		return nil, fmt.Errorf("ID рассылки не указан")
	}

	url := fmt.Sprintf("%s/api/v4/mailings/%d", requester.GetBaseURL(), mailingData.ID)

	data, err := json.Marshal(mailingData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := requester.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var updatedMailing Mailing
	if err := json.NewDecoder(resp.Body).Decode(&updatedMailing); err != nil {
		return nil, err
	}

	return &updatedMailing, nil
}

// Delete удаляет рассылку по ID.
//
// Пример использования:
//
//	err := mailing.Delete(ctx, apiClient, 123)
func Delete(ctx context.Context, apiClient *client.Client, id int) error {
	return DeleteWithRequester(ctx, apiClient, id)
}

// DeleteWithRequester удаляет рассылку с использованием интерфейса Requester.
func DeleteWithRequester(ctx context.Context, requester Requester, id int) error {
	url := fmt.Sprintf("%s/api/v4/mailings/%d", requester.GetBaseURL(), id)

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := requester.DoRequest(ctx, req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	return nil
}

// ChangeStatus изменяет статус рассылки.
//
// Пример использования:
//
//	updatedMailing, err := mailing.ChangeStatus(ctx, apiClient, 123, mailing.MailingStatusPaused)
func ChangeStatus(ctx context.Context, apiClient *client.Client, id int, status MailingStatus) (*Mailing, error) {
	return ChangeStatusWithRequester(ctx, apiClient, id, status)
}

// ChangeStatusWithRequester изменяет статус рассылки с использованием интерфейса Requester.
func ChangeStatusWithRequester(ctx context.Context, requester Requester, id int, status MailingStatus) (*Mailing, error) {
	url := fmt.Sprintf("%s/api/v4/mailings/%d/status", requester.GetBaseURL(), id)

	data, err := json.Marshal(map[string]string{
		"status": string(status),
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := requester.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var updatedMailing Mailing
	if err := json.NewDecoder(resp.Body).Decode(&updatedMailing); err != nil {
		return nil, err
	}

	return &updatedMailing, nil
}

// GetStats получает статистику рассылки.
//
// Пример использования:
//
//	stats, err := mailing.GetStats(ctx, apiClient, 123)
func GetStats(ctx context.Context, apiClient *client.Client, id int) (*Stats, error) {
	return GetStatsWithRequester(ctx, apiClient, id)
}

// GetStatsWithRequester получает статистику рассылки с использованием интерфейса Requester.
func GetStatsWithRequester(ctx context.Context, requester Requester, id int) (*Stats, error) {
	url := fmt.Sprintf("%s/api/v4/mailings/%d/stats", requester.GetBaseURL(), id)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := requester.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var stats Stats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, err
	}

	return &stats, nil
}

// AddRecipients добавляет получателей в рассылку.
//
// Пример использования:
//
//	contactIDs := []int{1001, 1002, 1003}
//	err := mailing.AddRecipients(ctx, apiClient, 123, contactIDs)
func AddRecipients(ctx context.Context, apiClient *client.Client, id int, contactIDs []int) error {
	return AddRecipientsWithRequester(ctx, apiClient, id, contactIDs)
}

// AddRecipientsWithRequester добавляет получателей в рассылку с использованием интерфейса Requester.
func AddRecipientsWithRequester(ctx context.Context, requester Requester, id int, contactIDs []int) error {
	url := fmt.Sprintf("%s/api/v4/mailings/%d/recipients", requester.GetBaseURL(), id)

	data, err := json.Marshal(map[string][]int{
		"contact_ids": contactIDs,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := requester.DoRequest(ctx, req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	return nil
}

// RemoveRecipients удаляет получателей из рассылки.
//
// Пример использования:
//
//	contactIDs := []int{1001, 1002}
//	err := mailing.RemoveRecipients(ctx, apiClient, 123, contactIDs)
func RemoveRecipients(ctx context.Context, apiClient *client.Client, id int, contactIDs []int) error {
	return RemoveRecipientsWithRequester(ctx, apiClient, id, contactIDs)
}

// RemoveRecipientsWithRequester удаляет получателей из рассылки с использованием интерфейса Requester.
func RemoveRecipientsWithRequester(ctx context.Context, requester Requester, id int, contactIDs []int) error {
	url := fmt.Sprintf("%s/api/v4/mailings/%d/recipients/delete", requester.GetBaseURL(), id)

	data, err := json.Marshal(map[string][]int{
		"contact_ids": contactIDs,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := requester.DoRequest(ctx, req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	return nil
}

// ListTemplates получает список шаблонов рассылок.
//
// Пример использования:
//
//	templates, err := mailing.ListTemplates(ctx, apiClient, 1, 50)
func ListTemplates(ctx context.Context, apiClient *client.Client, page, limit int) ([]Template, error) {
	return ListTemplatesWithRequester(ctx, apiClient, page, limit)
}

// ListTemplatesWithRequester получает список шаблонов рассылок с использованием интерфейса Requester.
func ListTemplatesWithRequester(ctx context.Context, requester Requester, page, limit int) ([]Template, error) {
	baseURL := fmt.Sprintf("%s/api/v4/mailing_templates", requester.GetBaseURL())

	params := url.Values{}
	params.Add("page", strconv.Itoa(page))
	params.Add("limit", strconv.Itoa(limit))

	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := requester.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var response struct {
		Embedded struct {
			Templates []Template `json:"templates"`
		} `json:"_embedded"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Embedded.Templates, nil
}

// GetTemplate получает информацию о конкретном шаблоне рассылки.
//
// Пример использования:
//
//	template, err := mailing.GetTemplate(ctx, apiClient, 123)
func GetTemplate(ctx context.Context, apiClient *client.Client, id int) (*Template, error) {
	return GetTemplateWithRequester(ctx, apiClient, id)
}

// GetTemplateWithRequester получает информацию о конкретном шаблоне рассылки с использованием интерфейса Requester.
func GetTemplateWithRequester(ctx context.Context, requester Requester, id int) (*Template, error) {
	url := fmt.Sprintf("%s/api/v4/mailing_templates/%d", requester.GetBaseURL(), id)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := requester.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var template Template
	if err := json.NewDecoder(resp.Body).Decode(&template); err != nil {
		return nil, err
	}

	return &template, nil
}
