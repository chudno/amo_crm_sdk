// Package widgets предоставляет методы для работы с виджетами в amoCRM.
package widgets

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

// WidgetType определяет тип виджета
type WidgetType string

// Константы для типов виджетов
const (
	WidgetTypeIntercom          WidgetType = "intercom"
	WidgetTypeJivosite          WidgetType = "jivosite"
	WidgetTypeCallback          WidgetType = "callback"
	WidgetTypePipeline          WidgetType = "pipeline"
	WidgetTypeMailchimp         WidgetType = "mailchimp"
	WidgetTypeCustom            WidgetType = "custom"
	WidgetTypeGoalMeter         WidgetType = "goal_meter"
	WidgetTypeDigitalPipeline   WidgetType = "digital_pipeline"
	WidgetTypeSupport           WidgetType = "support"
	WidgetTypeIpTelephony       WidgetType = "ip_telephony"
	WidgetTypePayment           WidgetType = "payment"
	WidgetTypeAmoButtons        WidgetType = "amo_buttons"
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
	ID             int          `json:"id,omitempty"`
	Name           string       `json:"name,omitempty"`
	Code           string       `json:"code,omitempty"`
	Type           WidgetType   `json:"type,omitempty"`
	Status         WidgetStatus `json:"status,omitempty"`
	CreatedBy      int          `json:"created_by,omitempty"`
	UpdatedBy      int          `json:"updated_by,omitempty"`
	CreatedAt      int          `json:"created_at,omitempty"`
	UpdatedAt      int          `json:"updated_at,omitempty"`
	AccountID      int          `json:"account_id,omitempty"`
	Settings       any          `json:"settings,omitempty"`
	Rights         *Rights      `json:"rights,omitempty"`
	Marketplace    *Marketplace `json:"marketplace,omitempty"`
	IsConfigured   bool         `json:"is_configured,omitempty"`
	VerifiedAt     int          `json:"verified_at,omitempty"`
	MainVersion    string       `json:"main_version,omitempty"`
	CurrentVersion string       `json:"current_version,omitempty"`
	IsDeleted      bool         `json:"is_deleted,omitempty"`
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

// ListResponse структура ответа API amoCRM для списка виджетов
type ListResponse struct {
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

// List получает список виджетов с возможностью фильтрации
//
// Пример использования:
//
//	// Фильтрация по типу
//	types := []widgets.WidgetType{widgets.WidgetTypeIntercom, widgets.WidgetTypeCallback}
//	widgetsList, err := widgets.List(apiClient, 1, 50, widgets.WithWidgetTypes(types))
func List(ctx context.Context, apiClient *client.Client, page, limit int, options ...WithOption) ([]Widget, error) {
	return ListWithRequester(ctx, apiClient, page, limit, options...)
}

// ListWithRequester получает список виджетов с использованием интерфейса Requester
func ListWithRequester(ctx context.Context, requester Requester, page, limit int, options ...WithOption) ([]Widget, error) {
	params := make(map[string]string)
	params["page"] = strconv.Itoa(page)
	params["limit"] = strconv.Itoa(limit)

	for _, option := range options {
		option(params)
	}

	url := "/api/v4/widgets"
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

	var widgetsResponse struct {
		Page     int `json:"page"`
		PerPage  int `json:"per_page"`
		Embedded struct {
			Widgets []Widget `json:"widgets"`
		} `json:"_embedded"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&widgetsResponse); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return widgetsResponse.Embedded.Widgets, nil
}

// Get получает информацию о конкретном виджете по ID
//
// Пример использования:
//
//	widget, err := widgets.Get(apiClient, 123)
func Get(ctx context.Context, apiClient *client.Client, widgetID int) (*Widget, error) {
	return GetWithRequester(ctx, apiClient, widgetID)
}

// GetWithRequester получает информацию о конкретном виджете по ID с использованием интерфейса Requester
func GetWithRequester(ctx context.Context, requester Requester, widgetID int) (*Widget, error) {
	url := fmt.Sprintf("/api/v4/widgets/%d", widgetID)

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

	var widget Widget
	if err := json.NewDecoder(resp.Body).Decode(&widget); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return &widget, nil
}

// Install устанавливает виджет из маркетплейса по его коду
//
// Пример использования:
//
//	widget, err := widgets.Install(apiClient, "intercom")
func Install(ctx context.Context, apiClient *client.Client, code string) (*Widget, error) {
	return InstallWithRequester(ctx, apiClient, code)
}

// InstallWithRequester устанавливает виджет из маркетплейса по его коду с использованием интерфейса Requester
func InstallWithRequester(ctx context.Context, requester Requester, code string) (*Widget, error) {
	url := "/api/v4/widgets"

	reqBody := struct {
		Code string `json:"code"`
	}{
		Code: code,
	}

	reqBodyJSON, err := json.Marshal(reqBody)
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

	var widget Widget
	if err := json.NewDecoder(resp.Body).Decode(&widget); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return &widget, nil
}

// UpdateSettings обновляет настройки виджета
//
// Пример использования:
//
//	 settings := map[string]any{
//			"api_key": "abc123",
//			"active": true,
//	 }
//	 widget, err := widgets.UpdateSettings(apiClient, 123, settings)
func UpdateSettings(ctx context.Context, apiClient *client.Client, widgetID int, settings any) (*Widget, error) {
	return UpdateSettingsWithRequester(ctx, apiClient, widgetID, settings)
}

// UpdateSettingsWithRequester обновляет настройки виджета с использованием интерфейса Requester
func UpdateSettingsWithRequester(ctx context.Context, requester Requester, widgetID int, settings any) (*Widget, error) {
	url := fmt.Sprintf("/api/v4/widgets/%d", widgetID)

	reqBody := struct {
		Settings any `json:"settings"`
	}{
		Settings: settings,
	}

	reqBodyJSON, err := json.Marshal(reqBody)
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

	var widget Widget
	if err := json.NewDecoder(resp.Body).Decode(&widget); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return &widget, nil
}

// Delete удаляет виджет
//
// Пример использования:
//
//	err := widgets.Delete(apiClient, 123)
func Delete(ctx context.Context, apiClient *client.Client, widgetID int) error {
	return DeleteWithRequester(ctx, apiClient, widgetID)
}

// DeleteWithRequester удаляет виджет с использованием интерфейса Requester
func DeleteWithRequester(ctx context.Context, requester Requester, widgetID int) error {
	url := fmt.Sprintf("/api/v4/widgets/%d", widgetID)

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
	Version      string  `json:"version,omitempty"`
	Pricing      string  `json:"pricing,omitempty"`
	Rating       float64 `json:"rating,omitempty"`
	ReviewsCount int     `json:"reviews_count,omitempty"`
}

// MarketplaceResponse структура ответа API amoCRM для списка виджетов в маркетплейсе
type MarketplaceResponse struct {
	Page       int                 `json:"page"`
	PerPage    int                 `json:"per_page"`
	TotalItems int                 `json:"_total_items"`
	Widgets    []MarketplaceWidget `json:"_embedded.widgets"`
}

// WithCategory добавляет фильтрацию по категории виджетов
func WithCategory(categoryID int) WithOption {
	return func(params map[string]string) {
		params["filter[category]"] = strconv.Itoa(categoryID)
	}
}

// ListMarketplace получает список доступных виджетов из маркетплейса
//
// Пример использования:
//
//	// Фильтрация по категории
//	widgetsList, err := widgets.ListMarketplace(apiClient, 1, 50, widgets.WithCategory(123))
func ListMarketplace(ctx context.Context, apiClient *client.Client, page, limit int, options ...WithOption) ([]MarketplaceWidget, error) {
	return ListMarketplaceWithRequester(ctx, apiClient, page, limit, options...)
}

// ListMarketplaceWithRequester получает список доступных виджетов из маркетплейса с использованием интерфейса Requester
func ListMarketplaceWithRequester(ctx context.Context, requester Requester, page, limit int, options ...WithOption) ([]MarketplaceWidget, error) {
	params := make(map[string]string)
	params["page"] = strconv.Itoa(page)
	params["limit"] = strconv.Itoa(limit)

	for _, option := range options {
		option(params)
	}

	url := "/api/v4/marketplace/widgets"
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

	var marketplaceResponse struct {
		Page     int `json:"page"`
		PerPage  int `json:"per_page"`
		Embedded struct {
			Widgets []MarketplaceWidget `json:"widgets"`
		} `json:"_embedded"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&marketplaceResponse); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return marketplaceResponse.Embedded.Widgets, nil
}

// SetStatus активирует или деактивирует виджет
//
// Пример использования:
//
//	// Деактивация виджета
//	widget, err := widgets.SetStatus(apiClient, 123, widgets.WidgetStatusInactive)
func SetStatus(ctx context.Context, apiClient *client.Client, widgetID int, status WidgetStatus) (*Widget, error) {
	return SetStatusWithRequester(ctx, apiClient, widgetID, status)
}

// SetStatusWithRequester активирует или деактивирует виджет с использованием интерфейса Requester
func SetStatusWithRequester(ctx context.Context, requester Requester, widgetID int, status WidgetStatus) (*Widget, error) {
	url := fmt.Sprintf("/api/v4/widgets/%d", widgetID)

	reqBody := struct {
		Status WidgetStatus `json:"status"`
	}{
		Status: status,
	}

	reqBodyJSON, err := json.Marshal(reqBody)
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

	var widget Widget
	if err := json.NewDecoder(resp.Body).Decode(&widget); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return &widget, nil
}

// BulkInput входные данные для массовой установки/удаления виджетов
type BulkInput struct {
	WidgetIDs []int    `json:"widget_ids,omitempty"`
	Codes     []string `json:"codes,omitempty"`
	Settings  []any    `json:"settings,omitempty"`
}

// BulkResponse ответ при массовых операциях с виджетами
type BulkResponse struct {
	Widgets []Widget `json:"_embedded.widgets"`
}

// BulkInstall массово устанавливает виджеты по их кодам
//
// Пример использования:
//
//	codes := []string{"intercom", "callback"}
//	widgets, err := widgets.BulkInstall(apiClient, codes)
func BulkInstall(ctx context.Context, apiClient *client.Client, codes []string) ([]Widget, error) {
	return BulkInstallWithRequester(ctx, apiClient, codes)
}

// BulkInstallWithRequester массово устанавливает виджеты по их кодам с использованием интерфейса Requester
func BulkInstallWithRequester(ctx context.Context, requester Requester, codes []string) ([]Widget, error) {
	url := "/api/v4/widgets"

	reqBody := BulkInput{
		Codes: codes,
	}

	reqBodyJSON, err := json.Marshal(reqBody)
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

// BulkDelete массово удаляет виджеты по их ID
//
// Пример использования:
//
//	ids := []int{123, 456}
//	err := widgets.BulkDelete(apiClient, ids)
func BulkDelete(ctx context.Context, apiClient *client.Client, widgetIDs []int) error {
	return BulkDeleteWithRequester(ctx, apiClient, widgetIDs)
}

// BulkDeleteWithRequester массово удаляет виджеты по их ID с использованием интерфейса Requester
func BulkDeleteWithRequester(ctx context.Context, requester Requester, widgetIDs []int) error {
	url := "/api/v4/widgets"

	reqBody := BulkInput{
		WidgetIDs: widgetIDs,
	}

	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("ошибка при кодировании тела запроса: %w", err)
	}

	baseURL := ""
	if c, ok := requester.(*client.Client); ok {
		baseURL = c.GetBaseURL()
	}

	fullURL := url
	if baseURL != "" {
		fullURL = baseURL + url
	}

	req, err := http.NewRequestWithContext(ctx, "DELETE", fullURL, strings.NewReader(string(reqBodyJSON)))
	if err != nil {
		return fmt.Errorf("ошибка при создании запроса: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

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
