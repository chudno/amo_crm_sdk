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

// ShortLinkFilter представляет параметры фильтрации для списка коротких ссылок.
type ShortLinkFilter struct {
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

// GetShortLinks получает список коротких ссылок с поддержкой фильтрации и пагинации.
//
// Пример использования:
//
//	filter := map[string]string{
//		"filter[entity_type]": "leads",
//		"filter[entity_id]": "123",
//	}
//	shortLinks, err := short_links.GetShortLinks(ctx, apiClient, 1, 50, short_links.WithFilter(filter))
func GetShortLinks(ctx context.Context, apiClient *client.Client, page, limit int, options ...WithOption) ([]ShortLink, error) {
	return GetShortLinksWithRequester(ctx, apiClient, page, limit, options...)
}

// GetShortLinksWithRequester получает список коротких ссылок с использованием интерфейса Requester.
func GetShortLinksWithRequester(ctx context.Context, requester Requester, page, limit int, options ...WithOption) ([]ShortLink, error) {
	// Формируем URL для запроса
	baseURL := fmt.Sprintf("%s/api/v4/short_links", requester.GetBaseURL())

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
	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	// Выполняем запрос
	resp, err := requester.DoRequest(ctx, req)
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
			ShortLinks []ShortLink `json:"short_links"`
		} `json:"_embedded"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Embedded.ShortLinks, nil
}

// GetShortLink получает информацию о конкретной короткой ссылке по ID.
//
// Пример использования:
//
//	shortLink, err := short_links.GetShortLink(ctx, apiClient, 123)
func GetShortLink(ctx context.Context, apiClient *client.Client, id int) (*ShortLink, error) {
	return GetShortLinkWithRequester(ctx, apiClient, id)
}

// GetShortLinkWithRequester получает информацию о конкретной короткой ссылке с использованием интерфейса Requester.
func GetShortLinkWithRequester(ctx context.Context, requester Requester, id int) (*ShortLink, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/short_links/%d", requester.GetBaseURL(), id)

	// Создаем запрос
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Выполняем запрос
	resp, err := requester.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	// Декодируем ответ
	var shortLink ShortLink
	if err := json.NewDecoder(resp.Body).Decode(&shortLink); err != nil {
		return nil, err
	}

	return &shortLink, nil
}

// CreateShortLink создает новую короткую ссылку.
//
// Пример использования:
//
//	newLink := &short_links.ShortLink{
//		URL: "https://example.com",
//		EntityType: "leads",
//		EntityID: 123,
//	}
//	createdLink, err := short_links.CreateShortLink(ctx, apiClient, newLink)
func CreateShortLink(ctx context.Context, apiClient *client.Client, shortLink *ShortLink) (*ShortLink, error) {
	return CreateShortLinkWithRequester(ctx, apiClient, shortLink)
}

// CreateShortLinkWithRequester создает новую короткую ссылку с использованием интерфейса Requester.
func CreateShortLinkWithRequester(ctx context.Context, requester Requester, shortLink *ShortLink) (*ShortLink, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/short_links", requester.GetBaseURL())

	// Подготавливаем данные для запроса
	data, err := json.Marshal(shortLink)
	if err != nil {
		return nil, err
	}

	// Создаем запрос
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := requester.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	// Декодируем ответ
	var response struct {
		Embedded struct {
			ShortLinks []ShortLink `json:"short_links"`
		} `json:"_embedded"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	// Проверяем, что создана хотя бы одна ссылка
	if len(response.Embedded.ShortLinks) == 0 {
		return nil, fmt.Errorf("короткая ссылка не была создана")
	}

	return &response.Embedded.ShortLinks[0], nil
}

// UpdateShortLink обновляет существующую короткую ссылку.
//
// Пример использования:
//
//	link := &short_links.ShortLink{
//		ID: 123,
//		URL: "https://updated-example.com",
//	}
//	updatedLink, err := short_links.UpdateShortLink(ctx, apiClient, link)
func UpdateShortLink(ctx context.Context, apiClient *client.Client, shortLink *ShortLink) (*ShortLink, error) {
	return UpdateShortLinkWithRequester(ctx, apiClient, shortLink)
}

// UpdateShortLinkWithRequester обновляет существующую короткую ссылку с использованием интерфейса Requester.
func UpdateShortLinkWithRequester(ctx context.Context, requester Requester, shortLink *ShortLink) (*ShortLink, error) {
	if shortLink.ID == 0 {
		return nil, fmt.Errorf("ID короткой ссылки не указан")
	}

	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/short_links/%d", requester.GetBaseURL(), shortLink.ID)

	// Подготавливаем данные для запроса
	data, err := json.Marshal(shortLink)
	if err != nil {
		return nil, err
	}

	// Создаем запрос
	req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := requester.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	// Декодируем ответ
	var updatedLink ShortLink
	if err := json.NewDecoder(resp.Body).Decode(&updatedLink); err != nil {
		return nil, err
	}

	return &updatedLink, nil
}

// DeleteShortLink удаляет короткую ссылку по ID.
//
// Пример использования:
//
//	err := short_links.DeleteShortLink(ctx, apiClient, 123)
func DeleteShortLink(ctx context.Context, apiClient *client.Client, id int) error {
	return DeleteShortLinkWithRequester(ctx, apiClient, id)
}

// DeleteShortLinkWithRequester удаляет короткую ссылку с использованием интерфейса Requester.
func DeleteShortLinkWithRequester(ctx context.Context, requester Requester, id int) error {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/short_links/%d", requester.GetBaseURL(), id)

	// Создаем запрос
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return err
	}

	// Выполняем запрос
	resp, err := requester.DoRequest(ctx, req)
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

// GetShortLinkStats получает статистику использования короткой ссылки.
//
// Пример использования:
//
//	stats, err := short_links.GetShortLinkStats(ctx, apiClient, 123)
func GetShortLinkStats(ctx context.Context, apiClient *client.Client, id int) (*ShortLink, error) {
	return GetShortLinkStatsWithRequester(ctx, apiClient, id)
}

// GetShortLinkStatsWithRequester получает статистику короткой ссылки с использованием интерфейса Requester.
func GetShortLinkStatsWithRequester(ctx context.Context, requester Requester, id int) (*ShortLink, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/short_links/%d/statistics", requester.GetBaseURL(), id)

	// Создаем запрос
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Выполняем запрос
	resp, err := requester.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	// Декодируем ответ
	var stats ShortLink
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, err
	}

	return &stats, nil
}
