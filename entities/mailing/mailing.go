// Пакет mailing предоставляет методы для работы с email-рассылками в amoCRM.
package mailing

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

// Requester - интерфейс для выполнения HTTP-запросов, используется для тестирования.
type Requester interface {
	DoRequest(req *http.Request) (*http.Response, error)
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
	Stats            *MailingStats     `json:"stats,omitempty"`
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

// MailingStats представляет статистику рассылки.
type MailingStats struct {
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

// GetMailings получает список рассылок с поддержкой фильтрации и пагинации.
//
// Пример использования:
//
//	filter := map[string]string{
//		"filter[status]": "active",
//	}
//	mailings, err := mailing.GetMailings(apiClient, 1, 50, mailing.WithFilter(filter))
func GetMailings(apiClient *client.Client, page, limit int, options ...WithOption) ([]Mailing, error) {
	return GetMailingsWithRequester(apiClient, page, limit, options...)
}

// GetMailingsWithRequester получает список рассылок с использованием интерфейса Requester.
func GetMailingsWithRequester(requester Requester, page, limit int, options ...WithOption) ([]Mailing, error) {
	// Формируем URL для запроса
	baseURL := fmt.Sprintf("%s/api/v4/mailings", requester.GetBaseURL())

	// Формируем параметры запроса
	params := map[string]string{
		"page":  strconv.Itoa(page),
		"limit": strconv.Itoa(limit),
	}

	// Применяем опции
	for _, option := range options {
		option(params)
	}

	// Формируем URL с параметрами
	queryParams := url.Values{}
	for key, value := range params {
		queryParams.Add(key, value)
	}
	requestURL := fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())

	// Создаем запрос
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	// Выполняем запрос
	resp, err := requester.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	// Декодируем ответ
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

// GetMailing получает информацию о конкретной рассылке по ID.
//
// Пример использования:
//
//	mailingInfo, err := mailing.GetMailing(apiClient, 123)
func GetMailing(apiClient *client.Client, id int) (*Mailing, error) {
	return GetMailingWithRequester(apiClient, id)
}

// GetMailingWithRequester получает информацию о конкретной рассылке с использованием интерфейса Requester.
func GetMailingWithRequester(requester Requester, id int) (*Mailing, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/mailings/%d", requester.GetBaseURL(), id)

	// Создаем запрос
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Выполняем запрос
	resp, err := requester.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	// Декодируем ответ
	var mailingInfo Mailing
	if err := json.NewDecoder(resp.Body).Decode(&mailingInfo); err != nil {
		return nil, err
	}

	return &mailingInfo, nil
}

// CreateMailing создает новую рассылку.
//
// Пример использования:
//
//	newMailing := &mailing.Mailing{
//		Name:     "Новая рассылка",
//		Subject:  "Важная информация",
//		Frequency: mailing.MailingFrequencyOnce,
//	}
//	createdMailing, err := mailing.CreateMailing(apiClient, newMailing)
func CreateMailing(apiClient *client.Client, mailingData *Mailing) (*Mailing, error) {
	return CreateMailingWithRequester(apiClient, mailingData)
}

// CreateMailingWithRequester создает новую рассылку с использованием интерфейса Requester.
func CreateMailingWithRequester(requester Requester, mailingData *Mailing) (*Mailing, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/mailings", requester.GetBaseURL())

	// Подготавливаем данные для запроса
	data, err := json.Marshal(mailingData)
	if err != nil {
		return nil, err
	}

	// Создаем запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := requester.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	// Декодируем ответ
	var createdMailing Mailing
	if err := json.NewDecoder(resp.Body).Decode(&createdMailing); err != nil {
		return nil, err
	}

	return &createdMailing, nil
}

// UpdateMailing обновляет существующую рассылку.
//
// Пример использования:
//
//	mailingUpdate := &mailing.Mailing{
//		ID:       123,
//		Name:     "Обновленная рассылка",
//		Subject:  "Новая тема рассылки",
//	}
//	updatedMailing, err := mailing.UpdateMailing(apiClient, mailingUpdate)
func UpdateMailing(apiClient *client.Client, mailingData *Mailing) (*Mailing, error) {
	return UpdateMailingWithRequester(apiClient, mailingData)
}

// UpdateMailingWithRequester обновляет существующую рассылку с использованием интерфейса Requester.
func UpdateMailingWithRequester(requester Requester, mailingData *Mailing) (*Mailing, error) {
	if mailingData.ID == 0 {
		return nil, fmt.Errorf("ID рассылки не указан")
	}

	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/mailings/%d", requester.GetBaseURL(), mailingData.ID)

	// Подготавливаем данные для запроса
	data, err := json.Marshal(mailingData)
	if err != nil {
		return nil, err
	}

	// Создаем запрос
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := requester.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	// Декодируем ответ
	var updatedMailing Mailing
	if err := json.NewDecoder(resp.Body).Decode(&updatedMailing); err != nil {
		return nil, err
	}

	return &updatedMailing, nil
}

// DeleteMailing удаляет рассылку по ID.
//
// Пример использования:
//
//	err := mailing.DeleteMailing(apiClient, 123)
func DeleteMailing(apiClient *client.Client, id int) error {
	return DeleteMailingWithRequester(apiClient, id)
}

// DeleteMailingWithRequester удаляет рассылку с использованием интерфейса Requester.
func DeleteMailingWithRequester(requester Requester, id int) error {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/mailings/%d", requester.GetBaseURL(), id)

	// Создаем запрос
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	// Выполняем запрос
	resp, err := requester.DoRequest(req)
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

// ChangeMailingStatus изменяет статус рассылки.
//
// Пример использования:
//
//	updatedMailing, err := mailing.ChangeMailingStatus(apiClient, 123, mailing.MailingStatusPaused)
func ChangeMailingStatus(apiClient *client.Client, id int, status MailingStatus) (*Mailing, error) {
	return ChangeMailingStatusWithRequester(apiClient, id, status)
}

// ChangeMailingStatusWithRequester изменяет статус рассылки с использованием интерфейса Requester.
func ChangeMailingStatusWithRequester(requester Requester, id int, status MailingStatus) (*Mailing, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/mailings/%d/status", requester.GetBaseURL(), id)

	// Подготавливаем данные для запроса
	data, err := json.Marshal(map[string]string{
		"status": string(status),
	})
	if err != nil {
		return nil, err
	}

	// Создаем запрос
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := requester.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	// Декодируем ответ
	var updatedMailing Mailing
	if err := json.NewDecoder(resp.Body).Decode(&updatedMailing); err != nil {
		return nil, err
	}

	return &updatedMailing, nil
}

// GetMailingStats получает статистику рассылки.
//
// Пример использования:
//
//	stats, err := mailing.GetMailingStats(apiClient, 123)
func GetMailingStats(apiClient *client.Client, id int) (*MailingStats, error) {
	return GetMailingStatsWithRequester(apiClient, id)
}

// GetMailingStatsWithRequester получает статистику рассылки с использованием интерфейса Requester.
func GetMailingStatsWithRequester(requester Requester, id int) (*MailingStats, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/mailings/%d/stats", requester.GetBaseURL(), id)

	// Создаем запрос
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Выполняем запрос
	resp, err := requester.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	// Декодируем ответ
	var stats MailingStats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, err
	}

	return &stats, nil
}

// AddMailingRecipients добавляет получателей в рассылку.
//
// Пример использования:
//
//	contactIDs := []int{1001, 1002, 1003}
//	err := mailing.AddMailingRecipients(apiClient, 123, contactIDs)
func AddMailingRecipients(apiClient *client.Client, id int, contactIDs []int) error {
	return AddMailingRecipientsWithRequester(apiClient, id, contactIDs)
}

// AddMailingRecipientsWithRequester добавляет получателей в рассылку с использованием интерфейса Requester.
func AddMailingRecipientsWithRequester(requester Requester, id int, contactIDs []int) error {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/mailings/%d/recipients", requester.GetBaseURL(), id)

	// Подготавливаем данные для запроса
	data, err := json.Marshal(map[string][]int{
		"contact_ids": contactIDs,
	})
	if err != nil {
		return err
	}

	// Создаем запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := requester.DoRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	return nil
}

// RemoveMailingRecipients удаляет получателей из рассылки.
//
// Пример использования:
//
//	contactIDs := []int{1001, 1002}
//	err := mailing.RemoveMailingRecipients(apiClient, 123, contactIDs)
func RemoveMailingRecipients(apiClient *client.Client, id int, contactIDs []int) error {
	return RemoveMailingRecipientsWithRequester(apiClient, id, contactIDs)
}

// RemoveMailingRecipientsWithRequester удаляет получателей из рассылки с использованием интерфейса Requester.
func RemoveMailingRecipientsWithRequester(requester Requester, id int, contactIDs []int) error {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/mailings/%d/recipients/delete", requester.GetBaseURL(), id)

	// Подготавливаем данные для запроса
	data, err := json.Marshal(map[string][]int{
		"contact_ids": contactIDs,
	})
	if err != nil {
		return err
	}

	// Создаем запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := requester.DoRequest(req)
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

// GetMailingTemplates получает список шаблонов рассылок.
//
// Пример использования:
//
//	templates, err := mailing.GetMailingTemplates(apiClient, 1, 50)
func GetMailingTemplates(apiClient *client.Client, page, limit int) ([]Template, error) {
	return GetMailingTemplatesWithRequester(apiClient, page, limit)
}

// GetMailingTemplatesWithRequester получает список шаблонов рассылок с использованием интерфейса Requester.
func GetMailingTemplatesWithRequester(requester Requester, page, limit int) ([]Template, error) {
	// Формируем URL для запроса
	baseURL := fmt.Sprintf("%s/api/v4/mailing_templates", requester.GetBaseURL())

	// Формируем параметры запроса
	params := url.Values{}
	params.Add("page", strconv.Itoa(page))
	params.Add("limit", strconv.Itoa(limit))

	// Формируем URL с параметрами
	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Создаем запрос
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	// Выполняем запрос
	resp, err := requester.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	// Декодируем ответ
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

// GetMailingTemplate получает информацию о конкретном шаблоне рассылки.
//
// Пример использования:
//
//	template, err := mailing.GetMailingTemplate(apiClient, 123)
func GetMailingTemplate(apiClient *client.Client, id int) (*Template, error) {
	return GetMailingTemplateWithRequester(apiClient, id)
}

// GetMailingTemplateWithRequester получает информацию о конкретном шаблоне рассылки с использованием интерфейса Requester.
func GetMailingTemplateWithRequester(requester Requester, id int) (*Template, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/mailing_templates/%d", requester.GetBaseURL(), id)

	// Создаем запрос
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Выполняем запрос
	resp, err := requester.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	// Декодируем ответ
	var template Template
	if err := json.NewDecoder(resp.Body).Decode(&template); err != nil {
		return nil, err
	}

	return &template, nil
}
