// Package segments предоставляет методы для работы с сегментами покупателей в amoCRM.
package segments

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/chudno/amo_crm_sdk/client"
)

// Segment представляет структуру сегмента в amoCRM.
type Segment struct {
	ID                     int         `json:"id,omitempty"`
	Name                   string      `json:"name"`
	Color                  string      `json:"color,omitempty"`
	Type                   SegmentType `json:"type,omitempty"`
	Filter                 *Filter     `json:"filter,omitempty"`
	AccountID              int         `json:"account_id,omitempty"`
	CreatedBy              int         `json:"created_by,omitempty"`
	UpdatedBy              int         `json:"updated_by,omitempty"`
	CreatedAt              int64       `json:"created_at,omitempty"`
	UpdatedAt              int64       `json:"updated_at,omitempty"`
	AvailableContactsCount int         `json:"available_contacts_count,omitempty"`
	ContactsCount          int         `json:"contacts_count,omitempty"`
	IsDeleted              bool        `json:"is_deleted,omitempty"`
	Embedded               *Embedded   `json:"_embedded,omitempty"`
	Links                  *Links      `json:"_links,omitempty"`
}

// SegmentType тип сегмента
type SegmentType string

const (
	// SegmentTypeDisposable одноразовый сегмент
	SegmentTypeDisposable SegmentType = "disposable"
	// SegmentTypeDynamic динамический сегмент
	SegmentTypeDynamic SegmentType = "dynamic"
)

// Filter фильтр сегмента
type Filter struct {
	Term  string       `json:"term,omitempty"`
	Logic string       `json:"logic,omitempty"`
	Nodes []FilterNode `json:"nodes,omitempty"`
}

// FilterNode узел фильтра
type FilterNode struct {
	FieldID    int          `json:"field_id,omitempty"`
	FieldCode  string       `json:"field_code,omitempty"`
	EntityType string       `json:"entity_type,omitempty"`
	Operator   string       `json:"operator,omitempty"`
	Value      string       `json:"value,omitempty"`
	Values     []string     `json:"values,omitempty"`
	MinValue   string       `json:"min_value,omitempty"`
	MaxValue   string       `json:"max_value,omitempty"`
	Term       string       `json:"term,omitempty"`
	Logic      string       `json:"logic,omitempty"`
	Nodes      []FilterNode `json:"nodes,omitempty"`
}

// Embedded вложенные поля
type Embedded struct {
	Contacts []struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Links struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"_links"`
	} `json:"contacts,omitempty"`
}

// Links ссылки на объекты
type Links struct {
	Self struct {
		Href string `json:"href"`
	} `json:"self"`
}

// ListResponse ответ при получении списка сегментов
type ListResponse struct {
	Page     int `json:"page"`
	PerPage  int `json:"per_page"`
	Embedded struct {
		Segments []Segment `json:"segments"`
	} `json:"_embedded"`
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
}

// WithOption функция-опция для запросов
type WithOption func(map[string]string)

// WithContacts опция для получения контактов в сегменте
func WithContacts() WithOption {
	return func(params map[string]string) {
		params["with"] = "contacts"
	}
}

// WithPage опция для установки страницы при получении списка
func WithPage(page int) WithOption {
	return func(params map[string]string) {
		params["page"] = strconv.Itoa(page)
	}
}

// WithLimit опция для установки лимита при получении списка
func WithLimit(limit int) WithOption {
	return func(params map[string]string) {
		params["limit"] = strconv.Itoa(limit)
	}
}

// WithFilter опция для фильтрации при получении списка
func WithFilter(filter map[string]string) WithOption {
	return func(params map[string]string) {
		for k, v := range filter {
			params[k] = v
		}
	}
}

// Create создает новый сегмент в amoCRM.
//
// Пример использования:
//
//	segment := &segments.Segment{
//		Name: "Новый сегмент",
//		Type: segments.SegmentTypeDynamic,
//		Filter: &segments.Filter{
//			Logic: "and",
//			Nodes: []segments.FilterNode{
//				{
//					FieldCode: "email",
//					Operator: "contains",
//					Value: "example.com",
//				},
//			},
//		},
//	}
//	createdSegment, err := segments.Create(ctx, apiClient, segment)
func Create(ctx context.Context, apiClient *client.Client, segment *Segment) (*Segment, error) {
	url := fmt.Sprintf("%s/api/v4/segments", apiClient.GetBaseURL())

	segmentJSON, err := json.Marshal(segment)
	if err != nil {
		return nil, fmt.Errorf("ошибка при сериализации сегмента: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(segmentJSON))
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var result struct {
		Embedded struct {
			Segments []Segment `json:"segments"`
		} `json:"_embedded"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	if len(result.Embedded.Segments) == 0 {
		return nil, fmt.Errorf("не удалось создать сегмент")
	}

	return &result.Embedded.Segments[0], nil
}

// List получает список сегментов с возможностью фильтрации и пагинации.
//
// Пример использования:
//
//	filter := map[string]string{
//		"filter[name]": "Активные клиенты",
//	}
//	segments, err := segments.List(ctx, apiClient, 1, 50, segments.WithFilter(filter))
func List(ctx context.Context, apiClient *client.Client, page, limit int, options ...WithOption) ([]Segment, error) {
	params := make(map[string]string)
	params["page"] = strconv.Itoa(page)
	params["limit"] = strconv.Itoa(limit)

	for _, option := range options {
		option(params)
	}

	url := fmt.Sprintf("%s/api/v4/segments", apiClient.GetBaseURL())
	if len(params) > 0 {
		var queryParams []string
		for key, value := range params {
			queryParams = append(queryParams, fmt.Sprintf("%s=%s", key, value))
		}
		url += "?" + strings.Join(queryParams, "&")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var segmentsResponse ListResponse
	if err := json.NewDecoder(resp.Body).Decode(&segmentsResponse); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return segmentsResponse.Embedded.Segments, nil
}

// Get получает информацию о конкретном сегменте по его ID.
//
// Пример использования:
//
//	segment, err := segments.Get(ctx, apiClient, 123, segments.WithContacts())
func Get(ctx context.Context, apiClient *client.Client, segmentID int, options ...WithOption) (*Segment, error) {
	params := make(map[string]string)

	for _, option := range options {
		option(params)
	}

	url := fmt.Sprintf("%s/api/v4/segments/%d", apiClient.GetBaseURL(), segmentID)
	if len(params) > 0 {
		var queryParams []string
		for key, value := range params {
			queryParams = append(queryParams, fmt.Sprintf("%s=%s", key, value))
		}
		url += "?" + strings.Join(queryParams, "&")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var segment Segment
	if err := json.NewDecoder(resp.Body).Decode(&segment); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return &segment, nil
}

// Update обновляет информацию о сегменте.
//
// Пример использования:
//
//	segment := &segments.Segment{
//		ID: 123,
//		Name: "Обновленный сегмент",
//		Color: "#FF5555",
//	}
//	updatedSegment, err := segments.Update(ctx, apiClient, segment)
func Update(ctx context.Context, apiClient *client.Client, segment *Segment) (*Segment, error) {
	if segment.ID == 0 {
		return nil, fmt.Errorf("ID сегмента не указан")
	}

	url := fmt.Sprintf("%s/api/v4/segments/%d", apiClient.GetBaseURL(), segment.ID)

	segmentJSON, err := json.Marshal(segment)
	if err != nil {
		return nil, fmt.Errorf("ошибка при сериализации сегмента: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewBuffer(segmentJSON))
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var updatedSegment Segment
	if err := json.NewDecoder(resp.Body).Decode(&updatedSegment); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return &updatedSegment, nil
}

// Delete удаляет сегмент по его ID.
//
// Пример использования:
//
//	err := segments.Delete(ctx, apiClient, 123)
func Delete(ctx context.Context, apiClient *client.Client, segmentID int) error {
	url := fmt.Sprintf("%s/api/v4/segments/%d", apiClient.GetBaseURL(), segmentID)

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	return nil
}

// AddContacts добавляет контакты в сегмент.
//
// Пример использования:
//
//	contactIDs := []int{123, 456, 789}
//	err := segments.AddContacts(ctx, apiClient, 42, contactIDs)
func AddContacts(ctx context.Context, apiClient *client.Client, segmentID int, contactIDs []int) error {
	url := fmt.Sprintf("%s/api/v4/segments/%d/contacts", apiClient.GetBaseURL(), segmentID)

	requestBody := struct {
		Contacts []int `json:"contacts"`
	}{
		Contacts: contactIDs,
	}

	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("ошибка при сериализации запроса: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestJSON))
	if err != nil {
		return fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	return nil
}

// RemoveContacts удаляет контакты из сегмента.
//
// Пример использования:
//
//	contactIDs := []int{123, 456, 789}
//	err := segments.RemoveContacts(ctx, apiClient, 42, contactIDs)
func RemoveContacts(ctx context.Context, apiClient *client.Client, segmentID int, contactIDs []int) error {
	url := fmt.Sprintf("%s/api/v4/segments/%d/contacts/delete", apiClient.GetBaseURL(), segmentID)

	requestBody := struct {
		Contacts []int `json:"contacts"`
	}{
		Contacts: contactIDs,
	}

	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("ошибка при сериализации запроса: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestJSON))
	if err != nil {
		return fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	return nil
}

// ListContacts получает список контактов в сегменте.
//
// Пример использования:
//
//	contactIDs, err := segments.ListContacts(ctx, apiClient, 42, 1, 50)
func ListContacts(ctx context.Context, apiClient *client.Client, segmentID, page, limit int) ([]int, error) {
	url := fmt.Sprintf("%s/api/v4/segments/%d/contacts?page=%d&limit=%d",
		apiClient.GetBaseURL(), segmentID, page, limit)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var response struct {
		Embedded struct {
			Contacts []struct {
				ID int `json:"id"`
			} `json:"contacts"`
		} `json:"_embedded"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	contactIDs := make([]int, len(response.Embedded.Contacts))
	for i, contact := range response.Embedded.Contacts {
		contactIDs[i] = contact.ID
	}

	return contactIDs, nil
}
