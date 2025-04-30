package widgets

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

// WidgetType определяет тип виджета
type WidgetType string

// Константы для типов виджетов
const (
	WidgetTypeIntercom     WidgetType = "intercom"
	WidgetTypeJivosite     WidgetType = "jivosite"
	WidgetTypeCallback     WidgetType = "callback"
	WidgetTypePipeline     WidgetType = "pipeline"
	WidgetTypeMailchimp    WidgetType = "mailchimp"
	WidgetTypeCustom       WidgetType = "custom"
	WidgetTypeGoalMeter    WidgetType = "goal_meter"
	WidgetTypeDigitalPipeline WidgetType = "digital_pipeline"
	WidgetTypeSupport      WidgetType = "support"
	WidgetTypeIpTelephony  WidgetType = "ip_telephony"
	WidgetTypePayment      WidgetType = "payment"
	WidgetTypeAmoButtons   WidgetType = "amo_buttons"
	WidgetTypeEmailSubscription WidgetType = "email_subscription"
)

// WidgetStatus определяет статус виджета
type WidgetStatus string

// Константы для статусов виджетов
const (
	WidgetStatusInstalled WidgetStatus = "installed"
	WidgetStatusDemo      WidgetStatus = "demo"
	WidgetStatusInactive  WidgetStatus = "inactive"
)

// Widget структура для работы с виджетами в amoCRM
type Widget struct {
	ID            int          `json:"id,omitempty"`
	Name          string       `json:"name,omitempty"`
	Code          string       `json:"code,omitempty"`
	Type          WidgetType   `json:"type,omitempty"`
	Status        WidgetStatus `json:"status,omitempty"`
	CreatedBy     int          `json:"created_by,omitempty"`
	UpdatedBy     int          `json:"updated_by,omitempty"`
	CreatedAt     int          `json:"created_at,omitempty"`
	UpdatedAt     int          `json:"updated_at,omitempty"`
	AccountID     int          `json:"account_id,omitempty"`
	Settings      interface{}  `json:"settings,omitempty"`
	Rights        *Rights      `json:"rights,omitempty"`
	Marketplace   *Marketplace `json:"marketplace,omitempty"`
	IsConfigured  bool         `json:"is_configured,omitempty"`
	VerifiedAt    int          `json:"verified_at,omitempty"`
	MainVersion   string       `json:"main_version,omitempty"`
	CurrentVersion string      `json:"current_version,omitempty"`
	IsDeleted     bool         `json:"is_deleted,omitempty"`
}

// Rights структура прав для виджета
type Rights struct {
	View    bool `json:"view"`
	Edit    bool `json:"edit"`
	Install bool `json:"install"`
	Delete  bool `json:"delete"`
}

// Marketplace структура рыночных данных для виджета
type Marketplace struct {
	ID          int    `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	LogoURL     string `json:"logo_url"`
	URL         string `json:"url"`
	Developer   struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"developer"`
	Categories []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"categories"`
}

// WidgetsResponse структура ответа API amoCRM для списка виджетов
type WidgetsResponse struct {
	Page       int      `json:"page"`
	PerPage    int      `json:"per_page"`
	TotalItems int      `json:"_total_items"`
	Widgets    []Widget `json:"_embedded.widgets"`
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

// WithWidgetTypes добавляет фильтрацию по типам виджетов
func WithWidgetTypes(types []WidgetType) WithOption {
	return func(params map[string]string) {
		if len(types) > 0 {
			var typesStr []string
			for _, t := range types {
				typesStr = append(typesStr, string(t))
			}
			params["filter[type]"] = strings.Join(typesStr, ",")
		}
	}
}

// GetWidgets получает список виджетов с возможностью фильтрации
//
// Пример использования:
//
//  // Фильтрация по типу
//  types := []widgets.WidgetType{widgets.WidgetTypeIntercom, widgets.WidgetTypeCallback}
//  widgetsList, err := widgets.GetWidgets(apiClient, 1, 50, widgets.WithWidgetTypes(types))
func GetWidgets(apiClient *client.Client, page, limit int, options ...WithOption) ([]Widget, error) {
	return GetWidgetsWithRequester(apiClient, page, limit, options...)
}

// GetWidgetsWithRequester получает список виджетов с использованием интерфейса Requester
func GetWidgetsWithRequester(requester Requester, page, limit int, options ...WithOption) ([]Widget, error) {
	// Формируем параметры запроса
	params := make(map[string]string)
	params["page"] = strconv.Itoa(page)
	params["limit"] = strconv.Itoa(limit)

	// Применяем опции
	for _, option := range options {
		option(params)
	}

	// Формируем URL для запроса
	url := "/api/v4/widgets"
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
	var widgetsResponse struct {
		Page    int `json:"page"`
		PerPage int `json:"per_page"`
		Embedded struct {
			Widgets []Widget `json:"widgets"`
		} `json:"_embedded"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&widgetsResponse); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return widgetsResponse.Embedded.Widgets, nil
}

// GetWidget получает информацию о конкретном виджете по ID
//
// Пример использования:
//
//  widget, err := widgets.GetWidget(apiClient, 123)
func GetWidget(apiClient *client.Client, widgetID int) (*Widget, error) {
	return GetWidgetWithRequester(apiClient, widgetID)
}

// GetWidgetWithRequester получает информацию о конкретном виджете по ID с использованием интерфейса Requester
func GetWidgetWithRequester(requester Requester, widgetID int) (*Widget, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("/api/v4/widgets/%d", widgetID)

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
	var widget Widget
	if err := json.NewDecoder(resp.Body).Decode(&widget); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return &widget, nil
}

// InstallWidget устанавливает виджет из маркетплейса по его коду
//
// Пример использования:
//
//  widget, err := widgets.InstallWidget(apiClient, "intercom")
func InstallWidget(apiClient *client.Client, code string) (*Widget, error) {
	return InstallWidgetWithRequester(apiClient, code)
}

// InstallWidgetWithRequester устанавливает виджет из маркетплейса по его коду с использованием интерфейса Requester
func InstallWidgetWithRequester(requester Requester, code string) (*Widget, error) {
	// Формируем URL для запроса
	url := "/api/v4/widgets"

	// Создаем тело запроса
	reqBody := struct {
		Code string `json:"code"`
	}{
		Code: code,
	}

	// Кодируем тело запроса в JSON
	reqBodyJSON, err := json.Marshal(reqBody)
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
	var widget Widget
	if err := json.NewDecoder(resp.Body).Decode(&widget); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return &widget, nil
}

// UpdateWidgetSettings обновляет настройки виджета
//
// Пример использования:
//
//  settings := map[string]interface{}{
//		"api_key": "abc123",
//		"active": true,
//  }
//  widget, err := widgets.UpdateWidgetSettings(apiClient, 123, settings)
func UpdateWidgetSettings(apiClient *client.Client, widgetID int, settings interface{}) (*Widget, error) {
	return UpdateWidgetSettingsWithRequester(apiClient, widgetID, settings)
}

// UpdateWidgetSettingsWithRequester обновляет настройки виджета с использованием интерфейса Requester
func UpdateWidgetSettingsWithRequester(requester Requester, widgetID int, settings interface{}) (*Widget, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("/api/v4/widgets/%d", widgetID)

	// Создаем тело запроса
	reqBody := struct {
		Settings interface{} `json:"settings"`
	}{
		Settings: settings,
	}

	// Кодируем тело запроса в JSON
	reqBodyJSON, err := json.Marshal(reqBody)
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
	var widget Widget
	if err := json.NewDecoder(resp.Body).Decode(&widget); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return &widget, nil
}

// DeleteWidget удаляет виджет
//
// Пример использования:
//
//  err := widgets.DeleteWidget(apiClient, 123)
func DeleteWidget(apiClient *client.Client, widgetID int) error {
	return DeleteWidgetWithRequester(apiClient, widgetID)
}

// DeleteWidgetWithRequester удаляет виджет с использованием интерфейса Requester
func DeleteWidgetWithRequester(requester Requester, widgetID int) error {
	// Формируем URL для запроса
	url := fmt.Sprintf("/api/v4/widgets/%d", widgetID)

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

// MarketplaceWidget структура для виджета из маркетплейса
type MarketplaceWidget struct {
	ID          int    `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	LogoURL     string `json:"logo_url"`
	URL         string `json:"url"`
	Installed   bool   `json:"installed,omitempty"`
	Developer   struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"developer"`
	Categories []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"categories"`
	Version       string `json:"version,omitempty"`
	Pricing       string `json:"pricing,omitempty"`
	Rating        float64 `json:"rating,omitempty"`
	ReviewsCount  int    `json:"reviews_count,omitempty"`
}

// MarketplaceResponse структура ответа API amoCRM для списка виджетов в маркетплейсе
type MarketplaceResponse struct {
	Page       int                `json:"page"`
	PerPage    int                `json:"per_page"`
	TotalItems int                `json:"_total_items"`
	Widgets    []MarketplaceWidget `json:"_embedded.widgets"`
}

// WithCategory добавляет фильтрацию по категории виджетов
func WithCategory(categoryID int) WithOption {
	return func(params map[string]string) {
		params["filter[category]"] = strconv.Itoa(categoryID)
	}
}

// GetMarketplaceWidgets получает список доступных виджетов из маркетплейса
//
// Пример использования:
//
//  // Фильтрация по категории
//  widgetsList, err := widgets.GetMarketplaceWidgets(apiClient, 1, 50, widgets.WithCategory(123))
func GetMarketplaceWidgets(apiClient *client.Client, page, limit int, options ...WithOption) ([]MarketplaceWidget, error) {
	return GetMarketplaceWidgetsWithRequester(apiClient, page, limit, options...)
}

// GetMarketplaceWidgetsWithRequester получает список доступных виджетов из маркетплейса с использованием интерфейса Requester
func GetMarketplaceWidgetsWithRequester(requester Requester, page, limit int, options ...WithOption) ([]MarketplaceWidget, error) {
	// Формируем параметры запроса
	params := make(map[string]string)
	params["page"] = strconv.Itoa(page)
	params["limit"] = strconv.Itoa(limit)

	// Применяем опции
	for _, option := range options {
		option(params)
	}

	// Формируем URL для запроса
	url := "/api/v4/marketplace/widgets"
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

	// ��азбираем ответ
	var marketplaceResponse struct {
		Page    int `json:"page"`
		PerPage int `json:"per_page"`
		Embedded struct {
			Widgets []MarketplaceWidget `json:"widgets"`
		} `json:"_embedded"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&marketplaceResponse); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return marketplaceResponse.Embedded.Widgets, nil
}

// SetWidgetStatus активирует или деактивирует виджет
//
// Пример использования:
//
//  // Деактивация виджета
//  widget, err := widgets.SetWidgetStatus(apiClient, 123, widgets.WidgetStatusInactive)
func SetWidgetStatus(apiClient *client.Client, widgetID int, status WidgetStatus) (*Widget, error) {
	return SetWidgetStatusWithRequester(apiClient, widgetID, status)
}

// SetWidgetStatusWithRequester активирует или деактивирует виджет с использованием интерфейса Requester
func SetWidgetStatusWithRequester(requester Requester, widgetID int, status WidgetStatus) (*Widget, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("/api/v4/widgets/%d", widgetID)

	// Создаем тело запроса
	reqBody := struct {
		Status WidgetStatus `json:"status"`
	}{
		Status: status,
	}

	// Кодируем тело запроса в JSON
	reqBodyJSON, err := json.Marshal(reqBody)
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
	var widget Widget
	if err := json.NewDecoder(resp.Body).Decode(&widget); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return &widget, nil
}

// BulkWidgetInput входные данные для массовой установки/удаления виджетов
type BulkWidgetInput struct {
	WidgetIDs []int         `json:"widget_ids,omitempty"`
	Codes     []string      `json:"codes,omitempty"`
	Settings  []interface{} `json:"settings,omitempty"`
}

// BulkWidgetResponse ответ при массовых операциях с виджетами
type BulkWidgetResponse struct {
	Widgets []Widget `json:"_embedded.widgets"`
}

// BulkInstallWidgets массово устанавливает виджеты по их кодам
//
// Пример использования:
//
//  codes := []string{"intercom", "callback"}
//  widgets, err := widgets.BulkInstallWidgets(apiClient, codes)
func BulkInstallWidgets(apiClient *client.Client, codes []string) ([]Widget, error) {
	return BulkInstallWidgetsWithRequester(apiClient, codes)
}

// BulkInstallWidgetsWithRequester массово устанавливает виджеты по их кодам с использованием интерфейса Requester
func BulkInstallWidgetsWithRequester(requester Requester, codes []string) ([]Widget, error) {
	// Формируем URL для запроса
	url := "/api/v4/widgets"

	// Создаем тело запроса
	reqBody := BulkWidgetInput{
		Codes: codes,
	}

	// Кодируем тело запроса в JSON
	reqBodyJSON, err := json.Marshal(reqBody)
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
	var bulkResponse struct {
		Embedded struct {
			Widgets []Widget `json:"widgets"`
		} `json:"_embedded"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&bulkResponse); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return bulkResponse.Embedded.Widgets, nil
}

// BulkDeleteWidgets массово удаляет виджеты по их ID
//
// Пример использования:
//
//  ids := []int{123, 456}
//  err := widgets.BulkDeleteWidgets(apiClient, ids)
func BulkDeleteWidgets(apiClient *client.Client, widgetIDs []int) error {
	return BulkDeleteWidgetsWithRequester(apiClient, widgetIDs)
}

// BulkDeleteWidgetsWithRequester массово удаляет виджеты по их ID с использованием интерфейса Requester
func BulkDeleteWidgetsWithRequester(requester Requester, widgetIDs []int) error {
	// Формируем URL для запроса
	url := "/api/v4/widgets"

	// Создаем тело запроса
	reqBody := BulkWidgetInput{
		WidgetIDs: widgetIDs,
	}

	// Кодируем тело запроса в JSON
	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("ошибка при кодировании тела запроса: %w", err)
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
	req, err := http.NewRequest("DELETE", fullURL, strings.NewReader(string(reqBodyJSON)))
	if err != nil {
		return fmt.Errorf("ошибка при создании запроса: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

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
