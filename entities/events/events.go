package events

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/chudno/amo_crm_sdk/client"
)

// EventType тип события.
type EventType string

const (
	// EventTypeNote тип события - Примечание
	EventTypeNote EventType = "note"
	// EventTypeCall тип события - Звонок
	EventTypeCall EventType = "call"
	// EventTypeMeeting тип события - Встреча
	EventTypeMeeting EventType = "meeting"
	// EventTypeIncomingLead тип события - Входящий лид
	EventTypeIncomingLead EventType = "incoming_lead"
	// EventTypeTaskResult тип события - Результат по задаче
	EventTypeTaskResult EventType = "task_result"
	// EventTypeMessage тип события - Сообщение
	EventTypeMessage EventType = "message"
	// EventTypeSendEmailStatus тип события - Статус отправки email
	EventTypeSendEmailStatus EventType = "send_email_status"
	// EventTypeCatalogObject тип события - Объект каталога
	EventTypeCatalogObject EventType = "catalog_object"
	// EventTypeEntityView тип события - Просмотр сущности
	EventTypeEntityView EventType = "entity_view"
	// EventTypeEntityUpdate тип события - Обновление сущности
	EventTypeEntityUpdate EventType = "entity_update"
	// EventTypeEntityStatusChange тип события - Изменение статуса сущности
	EventTypeEntityStatusChange EventType = "entity_status_change"
	// EventTypeEntityResponsibleChange тип события - Изменение ответственного сущности
	EventTypeEntityResponsibleChange EventType = "entity_responsible_change"
	// EventTypeEntityCreate тип события - Создание сущности
	EventTypeEntityCreate EventType = "entity_create"
	// EventTypeEntityDelete тип события - Удаление сущности
	EventTypeEntityDelete EventType = "entity_delete"
	// EventTypeActivityCreate тип события - Создание активности
	EventTypeActivityCreate EventType = "activity_create"
	// EventTypeActivityUpdate тип события - Обновление активности
	EventTypeActivityUpdate EventType = "activity_update"
	// EventTypeActivityStatusChange тип события - Изменение статуса активности
	EventTypeActivityStatusChange EventType = "activity_status_change"
	// EventTypeActivityDelete тип события - Удаление активности
	EventTypeActivityDelete EventType = "activity_delete"
)

// EventEntityType тип сущности события.
type EventEntityType string

const (
	// EventEntityTypeLead тип сущности события - Сделка
	EventEntityTypeLead EventEntityType = "lead"
	// EventEntityTypeContact тип сущности события - Контакт
	EventEntityTypeContact EventEntityType = "contact"
	// EventEntityTypeCompany тип сущности события - Компания
	EventEntityTypeCompany EventEntityType = "company"
	// EventEntityTypeCustomer тип сущности события - Покупатель
	EventEntityTypeCustomer EventEntityType = "customer"
	// EventEntityTypeTask тип сущности события - Задача
	EventEntityTypeTask EventEntityType = "task"
)

// Event структура события в amoCRM.
type Event struct {
	ID                 int             `json:"id,omitempty"`
	Type               EventType       `json:"type"`
	EntityID           int             `json:"entity_id"`
	EntityType         EventEntityType `json:"entity_type"`
	CreatedBy          int             `json:"created_by,omitempty"`
	AccountID          int             `json:"account_id,omitempty"`
	CreatedAt          int64           `json:"created_at,omitempty"`
	ValueAfter         json.RawMessage `json:"value_after,omitempty"`
	ValueBefore        json.RawMessage `json:"value_before,omitempty"`
	ValueBeforePretty  string          `json:"value_before_pretty,omitempty"`
	ValueAfterPretty   string          `json:"value_after_pretty,omitempty"`
	AdditionalEntities EventEntities   `json:"additional_entities,omitempty"`
	Link               string          `json:"link,omitempty"`
	Ver                string          `json:"__v,omitempty"`
	Embedded           *EventEmbedded  `json:"_embedded,omitempty"`
	Links              *EventLinks     `json:"_links,omitempty"`
}

// EventEntities дополнительные сущности события.
type EventEntities struct {
	Lead    *EventEntity `json:"lead,omitempty"`
	Contact *EventEntity `json:"contact,omitempty"`
	Company *EventEntity `json:"company,omitempty"`
	Task    *EventEntity `json:"task,omitempty"`
}

// EventEntity сущность события.
type EventEntity struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Created int64  `json:"created_at,omitempty"`
	Updated int64  `json:"updated_at,omitempty"`
}

// EventLinks ссылки события.
type EventLinks struct {
	Self struct {
		Href string `json:"href"`
	} `json:"self"`
}

// EventEmbedded вложенные данные события.
type EventEmbedded struct {
	Entity *EventEntity `json:"entity,omitempty"`
}

// GetEventsResponse структура ответа при получении списка событий.
type GetEventsResponse struct {
	Page      int            `json:"page"`
	PerPage   int            `json:"per_page"`
	Total     int            `json:"total"`
	Order     []Order        `json:"order"`
	Embedded  EventsEmbedded `json:"_embedded"`
	NextPage  string         `json:"_next_page"`
	PrevPage  string         `json:"_prev_page"`
	TotalPath string         `json:"_total_path"`
}

// EventsEmbedded вложенные данные списка событий.
type EventsEmbedded struct {
	Events []Event `json:"events"`
}

// Order структура для сортировки.
type Order struct {
	Field string `json:"field"`
	Order string `json:"order"`
}

// WithOption опции для запроса событий.
type WithOption func(map[string]string)

// WithEntity опция для получения информации о сущности события.
func WithEntity() WithOption {
	return func(params map[string]string) {
		params["with"] = "entity"
	}
}

// WithLimit опция для установки лимита результатов.
func WithLimit(limit int) WithOption {
	return func(params map[string]string) {
		params["limit"] = strconv.Itoa(limit)
	}
}

// WithPage опция для установки страницы результатов.
func WithPage(page int) WithOption {
	return func(params map[string]string) {
		params["page"] = strconv.Itoa(page)
	}
}

// WithFilter опция для фильтрации результатов.
func WithFilter(filterParams map[string]string) WithOption {
	return func(params map[string]string) {
		for key, value := range filterParams {
			params[key] = value
		}
	}
}

// WithOrder опция для сортировки результатов.
func WithOrder(field, order string) WithOption {
	return func(params map[string]string) {
		params["order["+field+"]"] = order
	}
}

// GetEvents получает список событий с возможностью фильтрации.
//
// Пример использования:
//
//	filter := map[string]string{
//		"filter[type]": string(events.EventTypeNote),
//		"filter[entity_type]": string(events.EventEntityTypeLead),
//	}
//	eventsList, err := events.GetEvents(apiClient, events.WithFilter(filter), events.WithLimit(50), events.WithPage(1))
func GetEvents(apiClient *client.Client, options ...WithOption) ([]Event, error) {
	params := make(map[string]string)

	// Применяем опции
	for _, option := range options {
		option(params)
	}

	// Формируем URL с параметрами
	url := "/api/v4/events"
	if len(params) > 0 {
		var queryParams []string
		for key, value := range params {
			queryParams = append(queryParams, fmt.Sprintf("%s=%s", key, value))
		}
		url += "?" + strings.Join(queryParams, "&")
	}

	// Формируем полный URL с базовым URL клиента
	fullURL := fmt.Sprintf("%s%s", apiClient.GetBaseURL(), url)

	// Создаем запрос
	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	// Выполняем запрос
	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	// Разбираем ответ
	var eventsResponse GetEventsResponse
	err = json.NewDecoder(resp.Body).Decode(&eventsResponse)
	if err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return eventsResponse.Embedded.Events, nil
}

// GetEvent получает информацию о конкретном событии по его ID.
//
// Пример использования:
//
//	event, err := events.GetEvent(apiClient, 123, events.WithEntity())
func GetEvent(apiClient *client.Client, eventID int, options ...WithOption) (*Event, error) {
	params := make(map[string]string)

	// Применяем опции
	for _, option := range options {
		option(params)
	}

	// Формируем URL с параметрами
	url := fmt.Sprintf("/api/v4/events/%d", eventID)
	if len(params) > 0 {
		var queryParams []string
		for key, value := range params {
			queryParams = append(queryParams, fmt.Sprintf("%s=%s", key, value))
		}
		url += "?" + strings.Join(queryParams, "&")
	}

	// Формируем полный URL с базовым URL клиента
	fullURL := fmt.Sprintf("%s%s", apiClient.GetBaseURL(), url)

	// Создаем запрос
	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	// Выполняем запрос
	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	// Разбираем ответ
	var event Event
	err = json.NewDecoder(resp.Body).Decode(&event)
	if err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return &event, nil
}
