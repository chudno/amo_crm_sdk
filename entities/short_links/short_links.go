// Package short_links предоставляет методы для работы с короткими ссылками в amoCRM.
package short_links

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/chudno/amo_crm_sdk/client"
)

// Requester - интерфейс для выполнения HTTP-запросов, используется для тестирования.
type Requester interface {
	DoRequest(ctx context.Context, req *http.Request) (*http.Response, error)
	GetBaseURL() string
}

// ShortLink представляет структуру короткой ссылки в amoCRM.
type ShortLink struct {
	ID            int    `json:"id,omitempty"`
	URL           string `json:"url"`
	Key           string `json:"key,omitempty"`
	ShortURL      string `json:"short_url,omitempty"`
	AccountID     int    `json:"account_id,omitempty"`
	EntityID      int    `json:"entity_id,omitempty"`
	EntityType    string `json:"entity_type,omitempty"`
	CreatedAt     int64  `json:"created_at,omitempty"`
	CreatedBy     int    `json:"created_by,omitempty"`
	UpdatedAt     int64  `json:"updated_at,omitempty"`
	MetadataID    int    `json:"metadata_id,omitempty"`
	VisitCount    int    `json:"visit_count,omitempty"`
	LastVisitAt   int64  `json:"last_visit_at,omitempty"`
	ExpireAt      int64  `json:"expire_at,omitempty"`
	UTMSource     string `json:"utm_source,omitempty"`
	UTMMedium     string `json:"utm_medium,omitempty"`
	UTMCampaign   string `json:"utm_campaign,omitempty"`
	UTMContent    string `json:"utm_content,omitempty"`
	UTMTerm       string `json:"utm_term,omitempty"`
	UseInEmbedded bool   `json:"use_in_embedded,omitempty"`
}

// Filter представляет параметры фильтрации для списка коротких ссылок.
type Filter struct {
	EntityID   int    `json:"filter[entity_id],omitempty"`
	EntityType string `json:"filter[entity_type],omitempty"`
	CreatedBy  int    `json:"filter[created_by],omitempty"`
}

// WithOption функциональный параметр для настройки запроса.
type WithOption func(params map[string]string)

// WithFilter добавляет фильтры при получении списка коротких ссылок.
func WithFilter(filter map[string]string) WithOption {
	return func(params map[string]string) {
		for k, v := range filter {
			params[k] = v
		}
	}
}

// List получает список коротких ссылок с поддержкой фильтрации и пагинации.
//
// Пример использования:
//
//	filter := map[string]string{
//		"filter[entity_type]": "leads",
//		"filter[entity_id]": "123",
//	}
//	shortLinks, err := short_links.List(ctx, apiClient, 1, 50, short_links.WithFilter(filter))
func List(ctx context.Context, apiClient *client.Client, page, limit int, options ...WithOption) ([]ShortLink, error) {
	return ListWithRequester(ctx, apiClient, page, limit, options...)
}

// ListWithRequester получает список коротких ссылок с использованием интерфейса Requester.
func ListWithRequester(ctx context.Context, requester Requester, page, limit int, options ...WithOption) ([]ShortLink, error) {
	baseURL := fmt.Sprintf("%s/api/v4/short_links", requester.GetBaseURL())

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
			ShortLinks []ShortLink `json:"short_links"`
		} `json:"_embedded"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Embedded.ShortLinks, nil
}

// Get получает информацию о конкретной короткой ссылке по ID.
//
// Пример использования:
//
//	shortLink, err := short_links.Get(ctx, apiClient, 123)
func Get(ctx context.Context, apiClient *client.Client, id int) (*ShortLink, error) {
	return GetWithRequester(ctx, apiClient, id)
}

// GetWithRequester получает информацию о конкретной короткой ссылке с использованием интерфейса Requester.
func GetWithRequester(ctx context.Context, requester Requester, id int) (*ShortLink, error) {
	url := fmt.Sprintf("%s/api/v4/short_links/%d", requester.GetBaseURL(), id)

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

	var shortLink ShortLink
	if err := json.NewDecoder(resp.Body).Decode(&shortLink); err != nil {
		return nil, err
	}

	return &shortLink, nil
}

// Create создает новую короткую ссылку.
//
// Пример использования:
//
//	newLink := &short_links.ShortLink{
//		URL: "https://example.com",
//		EntityType: "leads",
//		EntityID: 123,
//	}
//	createdLink, err := short_links.Create(ctx, apiClient, newLink)
func Create(ctx context.Context, apiClient *client.Client, shortLink *ShortLink) (*ShortLink, error) {
	return CreateWithRequester(ctx, apiClient, shortLink)
}

// CreateWithRequester создает новую короткую ссылку с использованием интерфейса Requester.
func CreateWithRequester(ctx context.Context, requester Requester, shortLink *ShortLink) (*ShortLink, error) {
	url := fmt.Sprintf("%s/api/v4/short_links", requester.GetBaseURL())

	data, err := json.Marshal(shortLink)
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

	var response struct {
		Embedded struct {
			ShortLinks []ShortLink `json:"short_links"`
		} `json:"_embedded"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if len(response.Embedded.ShortLinks) == 0 {
		return nil, fmt.Errorf("короткая ссылка не была создана")
	}

	return &response.Embedded.ShortLinks[0], nil
}

// Update обновляет существующую короткую ссылку.
//
// Пример использования:
//
//	link := &short_links.ShortLink{
//		ID: 123,
//		URL: "https://updated-example.com",
//	}
//	updatedLink, err := short_links.Update(ctx, apiClient, link)
func Update(ctx context.Context, apiClient *client.Client, shortLink *ShortLink) (*ShortLink, error) {
	return UpdateWithRequester(ctx, apiClient, shortLink)
}

// UpdateWithRequester обновляет существующую короткую ссылку с использованием интерфейса Requester.
func UpdateWithRequester(ctx context.Context, requester Requester, shortLink *ShortLink) (*ShortLink, error) {
	if shortLink.ID == 0 {
		return nil, fmt.Errorf("ID короткой ссылки не указан")
	}

	url := fmt.Sprintf("%s/api/v4/short_links/%d", requester.GetBaseURL(), shortLink.ID)

	data, err := json.Marshal(shortLink)
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

	var updatedLink ShortLink
	if err := json.NewDecoder(resp.Body).Decode(&updatedLink); err != nil {
		return nil, err
	}

	return &updatedLink, nil
}

// Delete удаляет короткую ссылку по ID.
//
// Пример использования:
//
//	err := short_links.Delete(ctx, apiClient, 123)
func Delete(ctx context.Context, apiClient *client.Client, id int) error {
	return DeleteWithRequester(ctx, apiClient, id)
}

// DeleteWithRequester удаляет короткую ссылку с использованием интерфейса Requester.
func DeleteWithRequester(ctx context.Context, requester Requester, id int) error {
	url := fmt.Sprintf("%s/api/v4/short_links/%d", requester.GetBaseURL(), id)

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

// GetStats получает статистику использования короткой ссылки.
//
// Пример использования:
//
//	stats, err := short_links.GetStats(ctx, apiClient, 123)
func GetStats(ctx context.Context, apiClient *client.Client, id int) (*ShortLink, error) {
	return GetStatsWithRequester(ctx, apiClient, id)
}

// GetStatsWithRequester получает статистику короткой ссылки с использованием интерфейса Requester.
func GetStatsWithRequester(ctx context.Context, requester Requester, id int) (*ShortLink, error) {
	url := fmt.Sprintf("%s/api/v4/short_links/%d/statistics", requester.GetBaseURL(), id)

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

	var stats ShortLink
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, err
	}

	return &stats, nil
}
