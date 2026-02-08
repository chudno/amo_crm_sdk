// Package access_rights предоставляет методы для работы с правами доступа в amoCRM.
package access_rights

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/chudno/amo_crm_sdk/client"
)

// Requester интерфейс для выполнения HTTP-запросов
type Requester interface {
	DoRequest(ctx context.Context, req *http.Request) (*http.Response, error)
}

// Type определяет тип доступа
type Type string

// Константы для типов доступа
const (
	TypeGroup  Type = "group"
	TypeCustom Type = "custom"
)

// EntityType определяет тип сущности для доступа
type EntityType string

// Константы для типов сущностей
const (
	EntityLead       EntityType = "leads"
	EntityContact    EntityType = "contacts"
	EntityCompany    EntityType = "companies"
	EntityTask       EntityType = "tasks"
	EntityCustomer   EntityType = "customers"
	EntityCatalog    EntityType = "catalogs"
	EntityUnsorted   EntityType = "unsorted"
	EntityWidgets    EntityType = "widgets"
	EntityMails      EntityType = "mail"
	EntityChatWidget EntityType = "chat_widget"
)

// Right структура для права доступа
type Right struct {
	ID         int         `json:"id,omitempty"`
	Name       string      `json:"name,omitempty"`
	Type       Type        `json:"type,omitempty"`
	Rights     Rights      `json:"rights,omitempty"`
	CreatedBy  int         `json:"created_by,omitempty"`
	UpdatedBy  int         `json:"updated_by,omitempty"`
	CreatedAt  int         `json:"created_at,omitempty"`
	UpdatedAt  int         `json:"updated_at,omitempty"`
	AccountID  int         `json:"account_id,omitempty"`
	UserIDs    []int       `json:"user_ids,omitempty"`
	UserGroups []UserGroup `json:"_embedded.user_groups,omitempty"`
}

// UserGroup структура для группы пользователей
type UserGroup struct {
	ID      int    `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	UserIDs []int  `json:"user_ids,omitempty"`
}

// Rights структура для прав доступа к различным сущностям
type Rights struct {
	Leads      EntityRights   `json:"leads,omitempty"`
	Contacts   EntityRights   `json:"contacts,omitempty"`
	Companies  EntityRights   `json:"companies,omitempty"`
	Tasks      EntityRights   `json:"tasks,omitempty"`
	Customers  EntityRights   `json:"customers,omitempty"`
	Catalogs   EntityRights   `json:"catalogs,omitempty"`
	Unsorted   EntityRights   `json:"unsorted,omitempty"`
	Widgets    EntityRights   `json:"widgets,omitempty"`
	Mail       EntityRights   `json:"mail,omitempty"`
	ChatWidget EntityRights   `json:"chat_widget,omitempty"`
	Settings   SettingsRights `json:"settings,omitempty"`
}

// EntityRights структура прав доступа к конкретной сущности
type EntityRights struct {
	View   bool `json:"view,omitempty"`
	Edit   bool `json:"edit,omitempty"`
	Add    bool `json:"add,omitempty"`
	Delete bool `json:"delete,omitempty"`
	Export bool `json:"export,omitempty"`
}

// SettingsRights структура прав доступа к настройкам
type SettingsRights struct {
	View bool `json:"view,omitempty"`
	Edit bool `json:"edit,omitempty"`
}

// ListResponse структура ответа API amoCRM для списка прав доступа
type ListResponse struct {
	Page       int     `json:"page"`
	PerPage    int     `json:"per_page"`
	TotalItems int     `json:"_total_items"`
	Rights     []Right `json:"_embedded.access_rights"`
}

// WithOption функциональный тип для передачи опций в методы
type WithOption func(params map[string]string)

// WithFilter добавляет фильтрацию к запросу
func WithFilter(filter map[string]string) WithOption {
	return func(params map[string]string) {
		for k, v := range filter {
			params[k] = v
		}
	}
}

// WithType добавляет фильтрацию по типу права доступа
func WithType(accessType Type) WithOption {
	return func(params map[string]string) {
		params["filter[type]"] = string(accessType)
	}
}

// List получает список прав доступа с возможностью фильтрации
//
// Пример использования:
//
//	// Фильтрация по типу
//	rights, err := access_rights.List(ctx, apiClient, 1, 50, access_rights.WithType(access_rights.TypeGroup))
func List(ctx context.Context, apiClient *client.Client, page, limit int, options ...WithOption) ([]Right, error) {
	return ListWithRequester(ctx, apiClient, page, limit, options...)
}

// ListWithRequester получает список прав доступа с использованием интерфейса Requester
func ListWithRequester(ctx context.Context, requester Requester, page, limit int, options ...WithOption) ([]Right, error) {
	params := make(map[string]string)
	params["page"] = strconv.Itoa(page)
	params["limit"] = strconv.Itoa(limit)

	for _, option := range options {
		option(params)
	}

	url := "/api/v4/access_rights"
	if len(params) > 0 {
		var queryParams []string
		for key, value := range params {
			queryParams = append(queryParams, fmt.Sprintf("%s=%s", key, value))
		}
		url += "?" + strings.Join(queryParams, "&")
	}

	baseURL := ""
	if c, ok := requester.(*client.Client); ok {
		baseURL = c.GetBaseURL()
	}

	fullURL := url
	if baseURL != "" {
		fullURL = baseURL + url
	}

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	resp, err := requester.DoRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var rightsResponse struct {
		Page     int `json:"page"`
		PerPage  int `json:"per_page"`
		Embedded struct {
			AccessRights []Right `json:"access_rights"`
		} `json:"_embedded"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&rightsResponse); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return rightsResponse.Embedded.AccessRights, nil
}

// Get получает информацию о конкретном праве доступа по ID
//
// Пример использования:
//
//	accessRight, err := access_rights.Get(ctx, apiClient, 123)
func Get(ctx context.Context, apiClient *client.Client, accessRightID int) (*Right, error) {
	return GetWithRequester(ctx, apiClient, accessRightID)
}

// GetWithRequester получает информацию о конкретном праве доступа по ID с использованием интерфейса Requester
func GetWithRequester(ctx context.Context, requester Requester, accessRightID int) (*Right, error) {
	url := fmt.Sprintf("/api/v4/access_rights/%d", accessRightID)

	baseURL := ""
	if c, ok := requester.(*client.Client); ok {
		baseURL = c.GetBaseURL()
	}

	fullURL := url
	if baseURL != "" {
		fullURL = baseURL + url
	}

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	resp, err := requester.DoRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var accessRight Right
	if err := json.NewDecoder(resp.Body).Decode(&accessRight); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return &accessRight, nil
}

// Create создает новое право доступа
//
// Пример использования:
//
//	newRight := &access_rights.Right{
//		Name: "Менеджеры продаж",
//		Type: access_rights.TypeGroup,
//		Rights: access_rights.Rights{
//			Leads: access_rights.EntityRights{
//				View: true,
//				Edit: true,
//				Add: true,
//			},
//		},
//		UserIDs: []int{123, 456},
//	}
//	createdRight, err := access_rights.Create(ctx, apiClient, newRight)
func Create(ctx context.Context, apiClient *client.Client, accessRight *Right) (*Right, error) {
	return CreateWithRequester(ctx, apiClient, accessRight)
}

// CreateWithRequester создает новое право доступа с использованием интерфейса Requester
func CreateWithRequester(ctx context.Context, requester Requester, accessRight *Right) (*Right, error) {
	url := "/api/v4/access_rights"

	reqBodyJSON, err := json.Marshal(accessRight)
	if err != nil {
		return nil, fmt.Errorf("ошибка при кодировании тела запроса: %w", err)
	}

	baseURL := ""
	if c, ok := requester.(*client.Client); ok {
		baseURL = c.GetBaseURL()
	}

	fullURL := url
	if baseURL != "" {
		fullURL = baseURL + url
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fullURL, strings.NewReader(string(reqBodyJSON)))
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := requester.DoRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var createdRight Right
	if err := json.NewDecoder(resp.Body).Decode(&createdRight); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return &createdRight, nil
}

// Update обновляет существующее право доступа
//
// Пример использования:
//
//	updateRight := &access_rights.Right{
//		ID: 123,
//		Name: "Менеджеры продаж (обновлено)",
//		Rights: access_rights.Rights{
//			Leads: access_rights.EntityRights{
//				View: true,
//				Edit: true,
//				Add: true,
//				Delete: true,
//			},
//		},
//		UserIDs: []int{123, 456, 789},
//	}
//	updatedRight, err := access_rights.Update(ctx, apiClient, updateRight)
func Update(ctx context.Context, apiClient *client.Client, accessRight *Right) (*Right, error) {
	return UpdateWithRequester(ctx, apiClient, accessRight)
}

// UpdateWithRequester обновляет существующее право доступа с использованием интерфейса Requester
func UpdateWithRequester(ctx context.Context, requester Requester, accessRight *Right) (*Right, error) {
	if accessRight.ID == 0 {
		return nil, fmt.Errorf("ID права доступа не может быть пустым")
	}

	url := fmt.Sprintf("/api/v4/access_rights/%d", accessRight.ID)

	reqBodyJSON, err := json.Marshal(accessRight)
	if err != nil {
		return nil, fmt.Errorf("ошибка при кодировании тела запроса: %w", err)
	}

	baseURL := ""
	if c, ok := requester.(*client.Client); ok {
		baseURL = c.GetBaseURL()
	}

	fullURL := url
	if baseURL != "" {
		fullURL = baseURL + url
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", fullURL, strings.NewReader(string(reqBodyJSON)))
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := requester.DoRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var updatedRight Right
	if err := json.NewDecoder(resp.Body).Decode(&updatedRight); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return &updatedRight, nil
}

// Delete удаляет право доступа
//
// Пример использования:
//
//	err := access_rights.Delete(ctx, apiClient, 123)
func Delete(ctx context.Context, apiClient *client.Client, accessRightID int) error {
	return DeleteWithRequester(ctx, apiClient, accessRightID)
}

// DeleteWithRequester удаляет право доступа с использованием интерфейса Requester
func DeleteWithRequester(ctx context.Context, requester Requester, accessRightID int) error {
	url := fmt.Sprintf("/api/v4/access_rights/%d", accessRightID)

	baseURL := ""
	if c, ok := requester.(*client.Client); ok {
		baseURL = c.GetBaseURL()
	}

	fullURL := url
	if baseURL != "" {
		fullURL = baseURL + url
	}

	req, err := http.NewRequestWithContext(ctx, "DELETE", fullURL, nil)
	if err != nil {
		return fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	resp, err := requester.DoRequest(ctx, req)
	if err != nil {
		return fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	return nil
}

// SetEntityRights обновляет права доступа к конкретной сущности
//
// Пример использования:
//
//	entityRights := access_rights.EntityRights{
//		View: true,
//		Edit: true,
//		Add: true,
//	}
//	updatedRight, err := access_rights.SetEntityRights(ctx, apiClient, 123, access_rights.EntityLead, entityRights)
func SetEntityRights(ctx context.Context, apiClient *client.Client, accessRightID int, entityType EntityType, rights EntityRights) (*Right, error) {
	return SetEntityRightsWithRequester(ctx, apiClient, accessRightID, entityType, rights)
}

// SetEntityRightsWithRequester обновляет права доступа к конкретной сущности с использованием интерфейса Requester
func SetEntityRightsWithRequester(ctx context.Context, requester Requester, accessRightID int, entityType EntityType, rights EntityRights) (*Right, error) {
	url := fmt.Sprintf("/api/v4/access_rights/%d", accessRightID)

	updateData := struct {
		Rights map[string]EntityRights `json:"rights"`
	}{
		Rights: map[string]EntityRights{
			string(entityType): rights,
		},
	}

	reqBodyJSON, err := json.Marshal(updateData)
	if err != nil {
		return nil, fmt.Errorf("ошибка при кодировании тела запроса: %w", err)
	}

	baseURL := ""
	if c, ok := requester.(*client.Client); ok {
		baseURL = c.GetBaseURL()
	}

	fullURL := url
	if baseURL != "" {
		fullURL = baseURL + url
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", fullURL, strings.NewReader(string(reqBodyJSON)))
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := requester.DoRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var updatedRight Right
	if err := json.NewDecoder(resp.Body).Decode(&updatedRight); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return &updatedRight, nil
}

// AddUsers добавляет пользователей в право доступа
//
// Пример использования:
//
//	userIDs := []int{123, 456, 789}
//	updatedRight, err := access_rights.AddUsers(ctx, apiClient, 123, userIDs)
func AddUsers(ctx context.Context, apiClient *client.Client, accessRightID int, userIDs []int) (*Right, error) {
	return AddUsersWithRequester(ctx, apiClient, accessRightID, userIDs)
}

// AddUsersWithRequester добавляет пользователей в право доступа с использованием интерфейса Requester
func AddUsersWithRequester(ctx context.Context, requester Requester, accessRightID int, userIDs []int) (*Right, error) {
	currentRight, err := GetWithRequester(ctx, requester, accessRightID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении права доступа: %w", err)
	}

	existingUsers := make(map[int]bool)
	for _, id := range currentRight.UserIDs {
		existingUsers[id] = true
	}

	for _, id := range userIDs {
		if !existingUsers[id] {
			currentRight.UserIDs = append(currentRight.UserIDs, id)
			existingUsers[id] = true
		}
	}

	updateData := struct {
		UserIDs []int `json:"user_ids"`
	}{
		UserIDs: currentRight.UserIDs,
	}

	url := fmt.Sprintf("/api/v4/access_rights/%d", accessRightID)

	reqBodyJSON, err := json.Marshal(updateData)
	if err != nil {
		return nil, fmt.Errorf("ошибка при кодировании тела запроса: %w", err)
	}

	baseURL := ""
	if c, ok := requester.(*client.Client); ok {
		baseURL = c.GetBaseURL()
	}

	fullURL := url
	if baseURL != "" {
		fullURL = baseURL + url
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", fullURL, strings.NewReader(string(reqBodyJSON)))
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := requester.DoRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var updatedRight Right
	if err := json.NewDecoder(resp.Body).Decode(&updatedRight); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return &updatedRight, nil
}

// RemoveUsers удаляет пользователей из права доступа
//
// Пример использования:
//
//	userIDs := []int{123, 456}
//	updatedRight, err := access_rights.RemoveUsers(ctx, apiClient, 123, userIDs)
func RemoveUsers(ctx context.Context, apiClient *client.Client, accessRightID int, userIDs []int) (*Right, error) {
	return RemoveUsersWithRequester(ctx, apiClient, accessRightID, userIDs)
}

// RemoveUsersWithRequester удаляет пользователей из права доступа с использованием интерфейса Requester
func RemoveUsersWithRequester(ctx context.Context, requester Requester, accessRightID int, userIDs []int) (*Right, error) {
	currentRight, err := GetWithRequester(ctx, requester, accessRightID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении права доступа: %w", err)
	}

	removeUsers := make(map[int]bool)
	for _, id := range userIDs {
		removeUsers[id] = true
	}

	newUserIDs := make([]int, 0, len(currentRight.UserIDs))
	for _, id := range currentRight.UserIDs {
		if !removeUsers[id] {
			newUserIDs = append(newUserIDs, id)
		}
	}

	updateData := struct {
		UserIDs []int `json:"user_ids"`
	}{
		UserIDs: newUserIDs,
	}

	url := fmt.Sprintf("/api/v4/access_rights/%d", accessRightID)

	reqBodyJSON, err := json.Marshal(updateData)
	if err != nil {
		return nil, fmt.Errorf("ошибка при кодировании тела запроса: %w", err)
	}

	baseURL := ""
	if c, ok := requester.(*client.Client); ok {
		baseURL = c.GetBaseURL()
	}

	fullURL := url
	if baseURL != "" {
		fullURL = baseURL + url
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", fullURL, strings.NewReader(string(reqBodyJSON)))
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := requester.DoRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var updatedRight Right
	if err := json.NewDecoder(resp.Body).Decode(&updatedRight); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return &updatedRight, nil
}
