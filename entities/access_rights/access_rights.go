package access_rights

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/chudno/amo_crm_sdk/client"
)

// Requester интерфейс для выполнения HTTP-запросов
type Requester interface {
	DoRequest(req *http.Request) (*http.Response, error)
}

// AccessRightsType определяет тип доступа
type AccessRightsType string

// Константы для типов доступа
const (
	TypeGroup  AccessRightsType = "group"
	TypeCustom AccessRightsType = "custom"
)

// AccessEntityType определяет тип сущности для доступа
type AccessEntityType string

// Константы для типов сущностей
const (
	EntityLead       AccessEntityType = "leads"
	EntityContact    AccessEntityType = "contacts"
	EntityCompany    AccessEntityType = "companies"
	EntityTask       AccessEntityType = "tasks"
	EntityCustomer   AccessEntityType = "customers"
	EntityCatalog    AccessEntityType = "catalogs"
	EntityUnsorted   AccessEntityType = "unsorted"
	EntityWidgets    AccessEntityType = "widgets"
	EntityMails      AccessEntityType = "mail"
	EntityChatWidget AccessEntityType = "chat_widget"
)

// AccessRight структура для права доступа
type AccessRight struct {
	ID         int              `json:"id,omitempty"`
	Name       string           `json:"name,omitempty"`
	Type       AccessRightsType `json:"type,omitempty"`
	Rights     Rights           `json:"rights,omitempty"`
	CreatedBy  int              `json:"created_by,omitempty"`
	UpdatedBy  int              `json:"updated_by,omitempty"`
	CreatedAt  int              `json:"created_at,omitempty"`
	UpdatedAt  int              `json:"updated_at,omitempty"`
	AccountID  int              `json:"account_id,omitempty"`
	UserIDs    []int            `json:"user_ids,omitempty"`
	UserGroups []UserGroup      `json:"_embedded.user_groups,omitempty"`
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

// AccessRightsResponse структура ответа API amoCRM для списка прав доступа
type AccessRightsResponse struct {
	Page       int           `json:"page"`
	PerPage    int           `json:"per_page"`
	TotalItems int           `json:"_total_items"`
	Rights     []AccessRight `json:"_embedded.access_rights"`
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
func WithType(accessType AccessRightsType) WithOption {
	return func(params map[string]string) {
		params["filter[type]"] = string(accessType)
	}
}

// GetAccessRights получает список прав доступа с возможностью фильтрации
//
// Пример использования:
//
//	// Фильтрация по типу
//	rights, err := access_rights.GetAccessRights(apiClient, 1, 50, access_rights.WithType(access_rights.TypeGroup))
func GetAccessRights(apiClient *client.Client, page, limit int, options ...WithOption) ([]AccessRight, error) {
	return GetAccessRightsWithRequester(apiClient, page, limit, options...)
}

// GetAccessRightsWithRequester получает список прав доступа с использованием интерфейса Requester
func GetAccessRightsWithRequester(requester Requester, page, limit int, options ...WithOption) ([]AccessRight, error) {
	// Формируем параметры запроса
	params := make(map[string]string)
	params["page"] = strconv.Itoa(page)
	params["limit"] = strconv.Itoa(limit)

	// Применяем опции
	for _, option := range options {
		option(params)
	}

	// Формируем URL для запроса
	url := "/api/v4/access_rights"
	if len(params) > 0 {
		var queryParams []string
		for key, value := range params {
			queryParams = append(queryParams, fmt.Sprintf("%s=%s", key, value))
		}
		url += "?" + strings.Join(queryParams, "&")
	}

	// Проверяем, что клиент имеет метод GetBaseURL()
	baseURL := ""
	if client, ok := requester.(*client.Client); ok {
		baseURL = client.GetBaseURL()
	}

	fullURL := url
	if baseURL != "" {
		fullURL = baseURL + url
	}

	// Создаем запрос
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	// Выполняем запрос
	resp, err := requester.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	// Разбираем ответ
	var rightsResponse struct {
		Page     int `json:"page"`
		PerPage  int `json:"per_page"`
		Embedded struct {
			AccessRights []AccessRight `json:"access_rights"`
		} `json:"_embedded"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&rightsResponse); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return rightsResponse.Embedded.AccessRights, nil
}

// GetAccessRight получает информацию о конкретном праве доступа по ID
//
// Пример использования:
//
//	accessRight, err := access_rights.GetAccessRight(apiClient, 123)
func GetAccessRight(apiClient *client.Client, accessRightID int) (*AccessRight, error) {
	return GetAccessRightWithRequester(apiClient, accessRightID)
}

// GetAccessRightWithRequester получает информацию о конкретном праве доступа по ID с использованием интерфейса Requester
func GetAccessRightWithRequester(requester Requester, accessRightID int) (*AccessRight, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("/api/v4/access_rights/%d", accessRightID)

	// Проверяем, что клиент имеет метод GetBaseURL()
	baseURL := ""
	if client, ok := requester.(*client.Client); ok {
		baseURL = client.GetBaseURL()
	}

	fullURL := url
	if baseURL != "" {
		fullURL = baseURL + url
	}

	// Создаем запрос
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	// Выполняем запрос
	resp, err := requester.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	// Разбираем ответ
	var accessRight AccessRight
	if err := json.NewDecoder(resp.Body).Decode(&accessRight); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return &accessRight, nil
}

// CreateAccessRight создает новое право доступа
//
// Пример использования:
//
//	newRight := &access_rights.AccessRight{
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
//	createdRight, err := access_rights.CreateAccessRight(apiClient, newRight)
func CreateAccessRight(apiClient *client.Client, accessRight *AccessRight) (*AccessRight, error) {
	return CreateAccessRightWithRequester(apiClient, accessRight)
}

// CreateAccessRightWithRequester создает новое право доступа с использованием интерфейса Requester
func CreateAccessRightWithRequester(requester Requester, accessRight *AccessRight) (*AccessRight, error) {
	// Формируем URL для запроса
	url := "/api/v4/access_rights"

	// Кодируем тело запроса в JSON
	reqBodyJSON, err := json.Marshal(accessRight)
	if err != nil {
		return nil, fmt.Errorf("ошибка при кодировании тела запроса: %w", err)
	}

	// Проверяем, что клиент имеет метод GetBaseURL()
	baseURL := ""
	if client, ok := requester.(*client.Client); ok {
		baseURL = client.GetBaseURL()
	}

	fullURL := url
	if baseURL != "" {
		fullURL = baseURL + url
	}

	// Создаем запрос
	req, err := http.NewRequest("POST", fullURL, strings.NewReader(string(reqBodyJSON)))
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := requester.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	// Разбираем ответ
	var createdRight AccessRight
	if err := json.NewDecoder(resp.Body).Decode(&createdRight); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return &createdRight, nil
}

// UpdateAccessRight обновляет существующее право доступа
//
// Пример использования:
//
//	updateRight := &access_rights.AccessRight{
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
//	updatedRight, err := access_rights.UpdateAccessRight(apiClient, updateRight)
func UpdateAccessRight(apiClient *client.Client, accessRight *AccessRight) (*AccessRight, error) {
	return UpdateAccessRightWithRequester(apiClient, accessRight)
}

// UpdateAccessRightWithRequester обновляет существующее право доступа с использованием интерфейса Requester
func UpdateAccessRightWithRequester(requester Requester, accessRight *AccessRight) (*AccessRight, error) {
	if accessRight.ID == 0 {
		return nil, fmt.Errorf("ID права доступа не может быть пустым")
	}

	// Формируем URL для запроса
	url := fmt.Sprintf("/api/v4/access_rights/%d", accessRight.ID)

	// Кодируем тело запроса в JSON
	reqBodyJSON, err := json.Marshal(accessRight)
	if err != nil {
		return nil, fmt.Errorf("ошибка при кодировании тела запроса: %w", err)
	}

	// Проверяем, что клиент имеет метод GetBaseURL()
	baseURL := ""
	if client, ok := requester.(*client.Client); ok {
		baseURL = client.GetBaseURL()
	}

	fullURL := url
	if baseURL != "" {
		fullURL = baseURL + url
	}

	// Создаем запрос
	req, err := http.NewRequest("PATCH", fullURL, strings.NewReader(string(reqBodyJSON)))
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := requester.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	// Разбираем ответ
	var updatedRight AccessRight
	if err := json.NewDecoder(resp.Body).Decode(&updatedRight); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return &updatedRight, nil
}

// DeleteAccessRight удаляет право доступа
//
// Пример использования:
//
//	err := access_rights.DeleteAccessRight(apiClient, 123)
func DeleteAccessRight(apiClient *client.Client, accessRightID int) error {
	return DeleteAccessRightWithRequester(apiClient, accessRightID)
}

// DeleteAccessRightWithRequester удаляет право доступа с использованием интерфейса Requester
func DeleteAccessRightWithRequester(requester Requester, accessRightID int) error {
	// Формируем URL для запроса
	url := fmt.Sprintf("/api/v4/access_rights/%d", accessRightID)

	// Проверяем, что клиент имеет метод GetBaseURL()
	baseURL := ""
	if client, ok := requester.(*client.Client); ok {
		baseURL = client.GetBaseURL()
	}

	fullURL := url
	if baseURL != "" {
		fullURL = baseURL + url
	}

	// Создаем запрос
	req, err := http.NewRequest("DELETE", fullURL, nil)
	if err != nil {
		return fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	// Выполняем запрос
	resp, err := requester.DoRequest(req)
	if err != nil {
		return fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
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
//	updatedRight, err := access_rights.SetEntityRights(apiClient, 123, access_rights.EntityLead, entityRights)
func SetEntityRights(apiClient *client.Client, accessRightID int, entityType AccessEntityType, rights EntityRights) (*AccessRight, error) {
	return SetEntityRightsWithRequester(apiClient, accessRightID, entityType, rights)
}

// SetEntityRightsWithRequester обновляет права доступа к конкретной сущности с использованием интерфейса Requester
func SetEntityRightsWithRequester(requester Requester, accessRightID int, entityType AccessEntityType, rights EntityRights) (*AccessRight, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("/api/v4/access_rights/%d", accessRightID)

	// Создаем структуру для обновления прав
	updateData := struct {
		Rights map[string]EntityRights `json:"rights"`
	}{
		Rights: map[string]EntityRights{
			string(entityType): rights,
		},
	}

	// Кодируем тело запроса в JSON
	reqBodyJSON, err := json.Marshal(updateData)
	if err != nil {
		return nil, fmt.Errorf("ошибка при кодировании тела запроса: %w", err)
	}

	// Проверяем, что клиент имеет метод GetBaseURL()
	baseURL := ""
	if client, ok := requester.(*client.Client); ok {
		baseURL = client.GetBaseURL()
	}

	fullURL := url
	if baseURL != "" {
		fullURL = baseURL + url
	}

	// Создаем запрос
	req, err := http.NewRequest("PATCH", fullURL, strings.NewReader(string(reqBodyJSON)))
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := requester.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	// Разбираем ответ
	var updatedRight AccessRight
	if err := json.NewDecoder(resp.Body).Decode(&updatedRight); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return &updatedRight, nil
}

// AddUsersToAccessRight добавляет пользователей в право доступа
//
// Пример использования:
//
//	userIDs := []int{123, 456, 789}
//	updatedRight, err := access_rights.AddUsersToAccessRight(apiClient, 123, userIDs)
func AddUsersToAccessRight(apiClient *client.Client, accessRightID int, userIDs []int) (*AccessRight, error) {
	return AddUsersToAccessRightWithRequester(apiClient, accessRightID, userIDs)
}

// AddUsersToAccessRightWithRequester добавляет пользователей в право доступа с использованием интерфейса Requester
func AddUsersToAccessRightWithRequester(requester Requester, accessRightID int, userIDs []int) (*AccessRight, error) {
	// Получаем текущее право доступа
	currentRight, err := GetAccessRightWithRequester(requester, accessRightID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении права доступа: %w", err)
	}

	// Создаем новый список пользователей без дубликатов
	existingUsers := make(map[int]bool)
	for _, id := range currentRight.UserIDs {
		existingUsers[id] = true
	}

	// Добавляем новых пользователей
	for _, id := range userIDs {
		if !existingUsers[id] {
			currentRight.UserIDs = append(currentRight.UserIDs, id)
			existingUsers[id] = true
		}
	}

	// Обновляем право доступа
	updateData := struct {
		UserIDs []int `json:"user_ids"`
	}{
		UserIDs: currentRight.UserIDs,
	}

	// Формируем URL для запроса
	url := fmt.Sprintf("/api/v4/access_rights/%d", accessRightID)

	// Кодируем тело запроса в JSON
	reqBodyJSON, err := json.Marshal(updateData)
	if err != nil {
		return nil, fmt.Errorf("ошибка при кодировании тела запроса: %w", err)
	}

	// Проверяем, что клиент имеет метод GetBaseURL()
	baseURL := ""
	if client, ok := requester.(*client.Client); ok {
		baseURL = client.GetBaseURL()
	}

	fullURL := url
	if baseURL != "" {
		fullURL = baseURL + url
	}

	// Создаем запрос
	req, err := http.NewRequest("PATCH", fullURL, strings.NewReader(string(reqBodyJSON)))
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := requester.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	// Разбираем ответ
	var updatedRight AccessRight
	if err := json.NewDecoder(resp.Body).Decode(&updatedRight); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return &updatedRight, nil
}

// RemoveUsersFromAccessRight удаляет пользователей из права доступа
//
// Пример использования:
//
//	userIDs := []int{123, 456}
//	updatedRight, err := access_rights.RemoveUsersFromAccessRight(apiClient, 123, userIDs)
func RemoveUsersFromAccessRight(apiClient *client.Client, accessRightID int, userIDs []int) (*AccessRight, error) {
	return RemoveUsersFromAccessRightWithRequester(apiClient, accessRightID, userIDs)
}

// RemoveUsersFromAccessRightWithRequester удаляет пользователей из права доступа с использованием интерфейса Requester
func RemoveUsersFromAccessRightWithRequester(requester Requester, accessRightID int, userIDs []int) (*AccessRight, error) {
	// Получаем текущее право доступа
	currentRight, err := GetAccessRightWithRequester(requester, accessRightID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении права доступа: %w", err)
	}

	// Создаем map для быстрого поиска пользователей для удаления
	removeUsers := make(map[int]bool)
	for _, id := range userIDs {
		removeUsers[id] = true
	}

	// Создаем новый список пользователей без удаляемых
	newUserIDs := make([]int, 0, len(currentRight.UserIDs))
	for _, id := range currentRight.UserIDs {
		if !removeUsers[id] {
			newUserIDs = append(newUserIDs, id)
		}
	}

	// Обновляем право доступа
	updateData := struct {
		UserIDs []int `json:"user_ids"`
	}{
		UserIDs: newUserIDs,
	}

	// Формируем URL для запроса
	url := fmt.Sprintf("/api/v4/access_rights/%d", accessRightID)

	// Кодируем тело запроса в JSON
	reqBodyJSON, err := json.Marshal(updateData)
	if err != nil {
		return nil, fmt.Errorf("ошибка при кодировании тела запроса: %w", err)
	}

	// Проверяем, что клиент имеет метод GetBaseURL()
	baseURL := ""
	if client, ok := requester.(*client.Client); ok {
		baseURL = client.GetBaseURL()
	}

	fullURL := url
	if baseURL != "" {
		fullURL = baseURL + url
	}

	// Создаем запрос
	req, err := http.NewRequest("PATCH", fullURL, strings.NewReader(string(reqBodyJSON)))
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := requester.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	// Разбираем ответ
	var updatedRight AccessRight
	if err := json.NewDecoder(resp.Body).Decode(&updatedRight); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return &updatedRight, nil
}
