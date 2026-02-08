// Package sources предоставляет методы для работы с источниками сделок в amoCRM.
package sources

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

// Source представляет источник сделок в amoCRM.
type Source struct {
	ID            int       `json:"id,omitempty"`
	Name          string    `json:"name"`
	Type          string    `json:"type,omitempty"`
	Default       bool      `json:"default,omitempty"`
	CreatedAt     int64     `json:"created_at,omitempty"`
	UpdatedAt     int64     `json:"updated_at,omitempty"`
	Deleted       bool      `json:"deleted,omitempty"`
	EffectiveFrom int64     `json:"effective_from,omitempty"`
	EffectiveTo   int64     `json:"effective_to,omitempty"`
	Pipeline      *Pipeline `json:"pipeline,omitempty"`
	Services      []Service `json:"services,omitempty"`
	External      *External `json:"external,omitempty"`
	Params        any       `json:"params,omitempty"`
}

// Pipeline представляет воронку, связанную с источником.
type Pipeline struct {
	ID int `json:"id,omitempty"`
}

// Service представляет сервис для источника.
type Service struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// External представляет внешние данные источника.
type External struct {
	ID             string `json:"id,omitempty"`
	Service        string `json:"service,omitempty"`
	ExternalParams any    `json:"external_params,omitempty"`
}

// WithOption функциональный параметр для настройки запроса.
type WithOption func(params map[string]string)

// WithFilter добавляет фильтры при получении списка источников.
func WithFilter(filter map[string]string) WithOption {
	return func(params map[string]string) {
		for k, v := range filter {
			params[k] = v
		}
	}
}

// List получает список источников сделок с поддержкой фильтрации и пагинации.
//
// Пример использования:
//
//	filter := map[string]string{
//		"filter[name]": "Реклама",
//	}
//	sources, err := sources.List(ctx, apiClient, 1, 50, sources.WithFilter(filter))
func List(ctx context.Context, apiClient *client.Client, page, limit int, options ...WithOption) ([]Source, error) {
	return ListWithRequester(ctx, apiClient, page, limit, options...)
}

// ListWithRequester получает список источников с использованием интерфейса Requester.
func ListWithRequester(ctx context.Context, requester Requester, page, limit int, options ...WithOption) ([]Source, error) {
	baseURL := fmt.Sprintf("%s/api/v4/sources", requester.GetBaseURL())

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
			Sources []Source `json:"sources"`
		} `json:"_embedded"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Embedded.Sources, nil
}

// Get получает информацию о конкретном источнике по ID.
//
// Пример использования:
//
//	sourceInfo, err := sources.Get(ctx, apiClient, 123)
func Get(ctx context.Context, apiClient *client.Client, id int) (*Source, error) {
	return GetWithRequester(ctx, apiClient, id)
}

// GetWithRequester получает информацию о конкретном источнике с использованием интерфейса Requester.
func GetWithRequester(ctx context.Context, requester Requester, id int) (*Source, error) {
	url := fmt.Sprintf("%s/api/v4/sources/%d", requester.GetBaseURL(), id)

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

	var sourceInfo Source
	if err := json.NewDecoder(resp.Body).Decode(&sourceInfo); err != nil {
		return nil, err
	}

	return &sourceInfo, nil
}

// Create создает новый источник сделок.
//
// Пример использования:
//
//	newSource := &sources.Source{
//		Name: "Новый источник",
//		Type: "other",
//	}
//	createdSource, err := sources.Create(ctx, apiClient, newSource)
func Create(ctx context.Context, apiClient *client.Client, sourceData *Source) (*Source, error) {
	return CreateWithRequester(ctx, apiClient, sourceData)
}

// CreateWithRequester создает новый источник с использованием интерфейса Requester.
func CreateWithRequester(ctx context.Context, requester Requester, sourceData *Source) (*Source, error) {
	url := fmt.Sprintf("%s/api/v4/sources", requester.GetBaseURL())

	data, err := json.Marshal(sourceData)
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

	var createdSource Source
	if err := json.NewDecoder(resp.Body).Decode(&createdSource); err != nil {
		return nil, err
	}

	return &createdSource, nil
}

// Update обновляет существующий источник сделок.
//
// Пример использования:
//
//	sourceUpdate := &sources.Source{
//		ID:   123,
//		Name: "Обновленный источник",
//	}
//	updatedSource, err := sources.Update(ctx, apiClient, sourceUpdate)
func Update(ctx context.Context, apiClient *client.Client, sourceData *Source) (*Source, error) {
	return UpdateWithRequester(ctx, apiClient, sourceData)
}

// UpdateWithRequester обновляет существующий источник с использованием интерфейса Requester.
func UpdateWithRequester(ctx context.Context, requester Requester, sourceData *Source) (*Source, error) {
	if sourceData.ID == 0 {
		return nil, fmt.Errorf("ID источника не указан")
	}

	url := fmt.Sprintf("%s/api/v4/sources/%d", requester.GetBaseURL(), sourceData.ID)

	data, err := json.Marshal(sourceData)
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

	var updatedSource Source
	if err := json.NewDecoder(resp.Body).Decode(&updatedSource); err != nil {
		return nil, err
	}

	return &updatedSource, nil
}

// Delete удаляет источник по ID.
//
// Пример использования:
//
//	err := sources.Delete(ctx, apiClient, 123)
func Delete(ctx context.Context, apiClient *client.Client, id int) error {
	return DeleteWithRequester(ctx, apiClient, id)
}

// DeleteWithRequester удаляет источник с использованием интерфейса Requester.
func DeleteWithRequester(ctx context.Context, requester Requester, id int) error {
	url := fmt.Sprintf("%s/api/v4/sources/%d", requester.GetBaseURL(), id)

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

// SetDefault устанавливает источник как используемый по умолчанию.
//
// Пример использования:
//
//	updatedSource, err := sources.SetDefault(ctx, apiClient, 123)
func SetDefault(ctx context.Context, apiClient *client.Client, id int) (*Source, error) {
	return SetDefaultWithRequester(ctx, apiClient, id)
}

// SetDefaultWithRequester устанавливает источник как используемый по умолчанию с использованием интерфейса Requester.
func SetDefaultWithRequester(ctx context.Context, requester Requester, id int) (*Source, error) {
	url := fmt.Sprintf("%s/api/v4/sources/%d/default", requester.GetBaseURL(), id)

	req, err := http.NewRequestWithContext(ctx, "PATCH", url, nil)
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

	var updatedSource Source
	if err := json.NewDecoder(resp.Body).Decode(&updatedSource); err != nil {
		return nil, err
	}

	return &updatedSource, nil
}

// ListServices получает список сервисов, доступных для источников сделок.
//
// Пример использования:
//
//	services, err := sources.ListServices(ctx, apiClient)
func ListServices(ctx context.Context, apiClient *client.Client) ([]Service, error) {
	return ListServicesWithRequester(ctx, apiClient)
}

// ListServicesWithRequester получает список сервисов с использованием интерфейса Requester.
func ListServicesWithRequester(ctx context.Context, requester Requester) ([]Service, error) {
	url := fmt.Sprintf("%s/api/v4/sources/services", requester.GetBaseURL())

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

	var response struct {
		Embedded struct {
			Services []Service `json:"services"`
		} `json:"_embedded"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Embedded.Services, nil
}

// LinkToPipeline связывает источник с воронкой.
//
// Пример использования:
//
//	linkedSource, err := sources.LinkToPipeline(ctx, apiClient, 123, 456)
func LinkToPipeline(ctx context.Context, apiClient *client.Client, sourceID, pipelineID int) (*Source, error) {
	return LinkToPipelineWithRequester(ctx, apiClient, sourceID, pipelineID)
}

// LinkToPipelineWithRequester связывает источник с воронкой с использованием интерфейса Requester.
func LinkToPipelineWithRequester(ctx context.Context, requester Requester, sourceID, pipelineID int) (*Source, error) {
	url := fmt.Sprintf("%s/api/v4/sources/%d/pipeline", requester.GetBaseURL(), sourceID)

	data, err := json.Marshal(map[string]int{
		"pipeline_id": pipelineID,
	})
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

	var linkedSource Source
	if err := json.NewDecoder(resp.Body).Decode(&linkedSource); err != nil {
		return nil, err
	}

	return &linkedSource, nil
}

// UnlinkFromPipeline удаляет связь источника с воронкой.
//
// Пример использования:
//
//	unlinkedSource, err := sources.UnlinkFromPipeline(ctx, apiClient, 123, 456)
func UnlinkFromPipeline(ctx context.Context, apiClient *client.Client, sourceID, pipelineID int) (*Source, error) {
	return UnlinkFromPipelineWithRequester(ctx, apiClient, sourceID, pipelineID)
}

// UnlinkFromPipelineWithRequester удаляет связь источника с воронкой с использованием интерфейса Requester.
func UnlinkFromPipelineWithRequester(ctx context.Context, requester Requester, sourceID, pipelineID int) (*Source, error) {
	url := fmt.Sprintf("%s/api/v4/sources/%d/pipeline/%d", requester.GetBaseURL(), sourceID, pipelineID)

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
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

	var unlinkedSource Source
	if err := json.NewDecoder(resp.Body).Decode(&unlinkedSource); err != nil {
		return nil, err
	}

	return &unlinkedSource, nil
}
