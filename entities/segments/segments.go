package segments

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/chudno/amo_crm_sdk/client"
)

// Segment представляет структуру сегмента в amoCRM.
type Segment struct {
	ID                 int          `json:"id,omitempty"`
	Name               string       `json:"name"`
	Color              string       `json:"color,omitempty"`
	Type               SegmentType  `json:"type,omitempty"`
	Filter             *Filter      `json:"filter,omitempty"`
	AccountID          int          `json:"account_id,omitempty"`
	CreatedBy          int          `json:"created_by,omitempty"`
	UpdatedBy          int          `json:"updated_by,omitempty"`
	CreatedAt          int64        `json:"created_at,omitempty"`
	UpdatedAt          int64        `json:"updated_at,omitempty"`
	AvailableContactsCount int      `json:"available_contacts_count,omitempty"`
	ContactsCount      int          `json:"contacts_count,omitempty"`
	IsDeleted          bool         `json:"is_deleted,omitempty"`
	Embedded          *Embedded     `json:"_embedded,omitempty"`
	Links             *Links        `json:"_links,omitempty"`
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
	Term  string        `json:"term,omitempty"`
	Logic string        `json:"logic,omitempty"`
	Nodes []FilterNode  `json:"nodes,omitempty"`
}

// FilterNode узел фильтра
type FilterNode struct {
	FieldID     int      `json:"field_id,omitempty"`
	FieldCode   string   `json:"field_code,omitempty"`
	EntityType  string   `json:"entity_type,omitempty"`
	Operator    string   `json:"operator,omitempty"`
	Value       string   `json:"value,omitempty"`
	Values     []string  `json:"values,omitempty"`
	MinValue    string   `json:"min_value,omitempty"`
	MaxValue    string   `json:"max_value,omitempty"`
	Term        string   `json:"term,omitempty"`
	Logic       string   `json:"logic,omitempty"`
	Nodes      []FilterNode `json:"nodes,omitempty"`
}

// Embedded вложенные поля
type Embedded struct {
	Contacts []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
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

// SegmentsResponse ответ при получении списка сегментов
type SegmentsResponse struct {
	Page     int      `json:"page"`
	PerPage  int      `json:"per_page"`
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

// AddSegment создает новый сегмент в amoCRM.
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
//	createdSegment, err := segments.AddSegment(apiClient, segment)
func AddSegment(apiClient *client.Client, segment *Segment) (*Segment, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/segments", apiClient.GetBaseURL())

	// Сериализуем сегмент в JSON
	segmentJSON, err := json.Marshal(segment)
	if err != nil {
		return nil, fmt.Errorf("ошибка при сериализации сегмента: %w", err)
	}

	// Создаем запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(segmentJSON))
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	// Разбираем ответ
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

// GetSegments получает список сегментов с возможностью фильтрации и пагинации.
//
// Пример использования:
//
//	filter := map[string]string{
//		"filter[name]": "Активные клиенты",
//	}
//	segments, err := segments.GetSegments(apiClient, 1, 50, segments.WithFilter(filter))
func GetSegments(apiClient *client.Client, page, limit int, options ...WithOption) ([]Segment, error) {
	// Формируем параметры запроса
	params := make(map[string]string)
	params["page"] = strconv.Itoa(page)
	params["limit"] = strconv.Itoa(limit)

	// Применяем опции
	for _, option := range options {
		option(params)
	}

	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/segments", apiClient.GetBaseURL())
	if len(params) > 0 {
		var queryParams []string
		for key, value := range params {
			queryParams = append(queryParams, fmt.Sprintf("%s=%s", key, value))
		}
		url += "?" + strings.Join(queryParams, "&")
	}

	// Создаем запрос
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	// Выполняем запрос
	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	// Разбираем ответ
	var segmentsResponse SegmentsResponse
	if err := json.NewDecoder(resp.Body).Decode(&segmentsResponse); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return segmentsResponse.Embedded.Segments, nil
}

// GetSegment получает информацию о конкретном сегменте по его ID.
//
// Пример использования:
//
//	segment, err := segments.GetSegment(apiClient, 123, segments.WithContacts())
func GetSegment(apiClient *client.Client, segmentID int, options ...WithOption) (*Segment, error) {
	// Формируем параметры запроса
	params := make(map[string]string)

	// Применяем опции
	for _, option := range options {
		option(params)
	}

	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/segments/%d", apiClient.GetBaseURL(), segmentID)
	if len(params) > 0 {
		var queryParams []string
		for key, value := range params {
			queryParams = append(queryParams, fmt.Sprintf("%s=%s", key, value))
		}
		url += "?" + strings.Join(queryParams, "&")
	}

	// Создаем запрос
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	// Выполняем запрос
	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	// Разбираем ответ
	var segment Segment
	if err := json.NewDecoder(resp.Body).Decode(&segment); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return &segment, nil
}

// UpdateSegment обновляет информацию о сегменте.
//
// Пример использования:
//
//	segment := &segments.Segment{
//		ID: 123,
//		Name: "Обновленный сегмент",
//		Color: "#FF5555",
//	}
//	updatedSegment, err := segments.UpdateSegment(apiClient, segment)
func UpdateSegment(apiClient *client.Client, segment *Segment) (*Segment, error) {
	if segment.ID == 0 {
		return nil, fmt.Errorf("ID сегмента не указан")
	}

	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/segments/%d", apiClient.GetBaseURL(), segment.ID)

	// Сериализуем сегмент в JSON
	segmentJSON, err := json.Marshal(segment)
	if err != nil {
		return nil, fmt.Errorf("ошибка при сериализации сегмента: %w", err)
	}

	// Создаем запрос
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(segmentJSON))
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	// Разбираем ответ
	var updatedSegment Segment
	if err := json.NewDecoder(resp.Body).Decode(&updatedSegment); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return &updatedSegment, nil
}

// DeleteSegment удаляет сегмент по его ID.
//
// Пример использования:
//
//	err := segments.DeleteSegment(apiClient, 123)
func DeleteSegment(apiClient *client.Client, segmentID int) error {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/segments/%d", apiClient.GetBaseURL(), segmentID)

	// Создаем запрос
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	// Выполняем запрос
	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	return nil
}

// AddContactsToSegment добавляет контакты в сегмент.
//
// Пример использования:
//
//	contactIDs := []int{123, 456, 789}
//	err := segments.AddContactsToSegment(apiClient, 42, contactIDs)
func AddContactsToSegment(apiClient *client.Client, segmentID int, contactIDs []int) error {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/segments/%d/contacts", apiClient.GetBaseURL(), segmentID)

	// Создаем тело запроса
	requestBody := struct {
		Contacts []int `json:"contacts"`
	}{
		Contacts: contactIDs,
	}

	// Сериализуем тело запроса в JSON
	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("ошибка при сериализации запроса: %w", err)
	}

	// Создаем запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestJSON))
	if err != nil {
		return fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	return nil
}

// RemoveContactsFromSegment удаляет контакты из сегмента.
//
// Пример использования:
//
//	contactIDs := []int{123, 456, 789}
//	err := segments.RemoveContactsFromSegment(apiClient, 42, contactIDs)
func RemoveContactsFromSegment(apiClient *client.Client, segmentID int, contactIDs []int) error {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/segments/%d/contacts/delete", apiClient.GetBaseURL(), segmentID)

	// Создаем тело запроса
	requestBody := struct {
		Contacts []int `json:"contacts"`
	}{
		Contacts: contactIDs,
	}

	// Сериализуем тело запроса в JSON
	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("ошибка при сериализации запроса: %w", err)
	}

	// Создаем запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestJSON))
	if err != nil {
		return fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	return nil
}

// GetSegmentContacts получает список контактов в сегменте.
//
// Пример использования:
//
//	contactIDs, err := segments.GetSegmentContacts(apiClient, 42, 1, 50)
func GetSegmentContacts(apiClient *client.Client, segmentID, page, limit int) ([]int, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/segments/%d/contacts?page=%d&limit=%d", 
		apiClient.GetBaseURL(), segmentID, page, limit)

	// Создаем запрос
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	// Выполняем запрос
	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	// Разбираем ответ
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

	// Извлекаем ID контактов
	contactIDs := make([]int, len(response.Embedded.Contacts))
	for i, contact := range response.Embedded.Contacts {
		contactIDs[i] = contact.ID
	}

	return contactIDs, nil
}
