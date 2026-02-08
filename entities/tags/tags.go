// Package tags предоставляет методы для взаимодействия с сущностями "Теги" в API amoCRM.
package tags

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/chudno/amo_crm_sdk/client"
)

// Tag представляет собой структуру тега в amoCRM.
type Tag struct {
	ID    int    `json:"id,omitempty"`
	Name  string `json:"name"`
	Color string `json:"color,omitempty"`
}

// ListResponse представляет ответ от API при получении списка тегов
type ListResponse struct {
	Page     int `json:"page"`
	PerPage  int `json:"per_page"`
	Total    int `json:"total"`
	Embedded struct {
		Tags []Tag `json:"tags"`
	} `json:"_embedded"`
}

// TagResponse представляет ответ от API при создании тега
type TagResponse struct {
	Tag `json:"tag"`
}

// EntityType представляет тип сущности для работы с тегами
type EntityType string

const (
	// EntityTypeContact - тип сущности "Контакт"
	EntityTypeContact EntityType = "contacts"
	// EntityTypeLead - тип сущности "Сделка"
	EntityTypeLead EntityType = "leads"
	// EntityTypeCompany - тип сущности "Компания"
	EntityTypeCompany EntityType = "companies"
	// EntityTypeCustomer - тип сущности "Покупатель"
	EntityTypeCustomer EntityType = "customers"
)

// List получает список тегов с возможностью пагинации по указанному типу сущности.
func List(ctx context.Context, apiClient *client.Client, entityType EntityType, page, limit int) ([]Tag, error) {
	baseURL := fmt.Sprintf("%s/api/v4/%s/tags", apiClient.GetBaseURL(), entityType)

	params := url.Values{}
	params.Add("page", fmt.Sprintf("%d", page))
	params.Add("limit", fmt.Sprintf("%d", limit))

	baseURL = baseURL + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var tags ListResponse
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return nil, err
	}

	return tags.Embedded.Tags, nil
}

// Create создает новый тег для указанного типа сущности.
func Create(ctx context.Context, apiClient *client.Client, entityType EntityType, tag *Tag) (*Tag, error) {
	apiURL := fmt.Sprintf("%s/api/v4/%s/tags", apiClient.GetBaseURL(), entityType)

	tagJSON, err := json.Marshal([]*Tag{tag})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(tagJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var response struct {
		Embedded struct {
			Tags []*Tag `json:"tags"`
		} `json:"_embedded"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if len(response.Embedded.Tags) == 0 {
		return nil, fmt.Errorf("не удалось создать тег")
	}

	return response.Embedded.Tags[0], nil
}

// CreateBatch создает несколько тегов для указанного типа сущности.
func CreateBatch(ctx context.Context, apiClient *client.Client, entityType EntityType, tags []Tag) ([]Tag, error) {
	url := fmt.Sprintf("%s/api/v4/%s/tags", apiClient.GetBaseURL(), entityType)

	tagsJSON, err := json.Marshal(tags)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(tagsJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var tagsResponse struct {
		Embedded struct {
			Tags []Tag `json:"tags"`
		} `json:"_embedded"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tagsResponse); err != nil {
		return nil, err
	}

	return tagsResponse.Embedded.Tags, nil
}

// Get получает информацию о теге по его ID для указанного типа сущности.
// Используется endpoint списка тегов с фильтром по ID, так как API amoCRM v4
// не поддерживает получение одного тега по прямому URL.
func Get(ctx context.Context, apiClient *client.Client, entityType EntityType, tagID int) (*Tag, error) {
	requestURL := fmt.Sprintf("%s/api/v4/%s/tags?filter[id][]=%d", apiClient.GetBaseURL(), entityType, tagID)

	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var listResponse ListResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResponse); err != nil {
		return nil, err
	}

	if len(listResponse.Embedded.Tags) == 0 {
		return nil, fmt.Errorf("тег с ID %d не найден", tagID)
	}

	return &listResponse.Embedded.Tags[0], nil
}

// Update обновляет информацию о теге по его ID для указанного типа сущности.
func Update(ctx context.Context, apiClient *client.Client, entityType EntityType, tag *Tag) (*Tag, error) {
	if tag.ID == 0 {
		return nil, fmt.Errorf("ID тега не может быть пустым")
	}

	url := fmt.Sprintf("%s/api/v4/%s/tags/%d", apiClient.GetBaseURL(), entityType, tag.ID)

	tagJSON, err := json.Marshal(tag)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewBuffer(tagJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var updatedTag Tag
	if err := json.NewDecoder(resp.Body).Decode(&updatedTag); err != nil {
		return nil, err
	}

	return &updatedTag, nil
}

// LinkEntity связывает сущность с тегами
func LinkEntity(ctx context.Context, apiClient *client.Client, entityType EntityType, entityID int, tags []Tag) error {
	url := fmt.Sprintf("%s/api/v4/%s/%d/tags", apiClient.GetBaseURL(), entityType, entityID)

	tagsJSON, err := json.Marshal(tags)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(tagsJSON))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	return nil
}

// ListForEntity получает список тегов для указанной сущности
func ListForEntity(ctx context.Context, apiClient *client.Client, entityType EntityType, entityID int) ([]Tag, error) {
	url := fmt.Sprintf("%s/api/v4/%s/%d/tags", apiClient.GetBaseURL(), entityType, entityID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var tagsResponse struct {
		Embedded struct {
			Tags []Tag `json:"tags"`
		} `json:"_embedded"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tagsResponse); err != nil {
		return nil, err
	}

	return tagsResponse.Embedded.Tags, nil
}
