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
	// Формируем базовый URL
	baseURL := fmt.Sprintf("%s/api/v4/%s/tags", apiClient.GetBaseURL(), entityType)

	// Добавляем параметры запроса
	params := url.Values{}
	params.Add("page", fmt.Sprintf("%d", page))
	params.Add("limit", fmt.Sprintf("%d", limit))

	// Добавляем параметры к URL
	baseURL = baseURL + "?" + params.Encode()

	// Создаем запрос
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
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
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/%s/tags", apiClient.GetBaseURL(), entityType)

	// Преобразуем структуру тега в JSON
	tagJSON, err := json.Marshal(tag)
	if err != nil {
		return nil, err
	}

	// Создаем запрос
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(tagJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var tagResponse TagResponse
	if err := json.NewDecoder(resp.Body).Decode(&tagResponse); err != nil {
		return nil, err
	}

	return &tagResponse.Tag, nil
}

// CreateBatch создает несколько тегов для указанного типа сущности.
func CreateBatch(ctx context.Context, apiClient *client.Client, entityType EntityType, tags []Tag) ([]Tag, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/%s/tags", apiClient.GetBaseURL(), entityType)

	// Преобразуем структуру тегов в JSON
	tagsJSON, err := json.Marshal(tags)
	if err != nil {
		return nil, err
	}

	// Создаем запрос
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(tagsJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
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
func Get(ctx context.Context, apiClient *client.Client, entityType EntityType, tagID int) (*Tag, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/%s/tags/%d", apiClient.GetBaseURL(), entityType, tagID)

	// Создаем запрос
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Выполняем запрос
	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var tag Tag
	if err := json.NewDecoder(resp.Body).Decode(&tag); err != nil {
		return nil, err
	}

	return &tag, nil
}

// Update обновляет информацию о теге по его ID для указанного типа сущности.
func Update(ctx context.Context, apiClient *client.Client, entityType EntityType, tag *Tag) (*Tag, error) {
	if tag.ID == 0 {
		return nil, fmt.Errorf("ID тега не может быть пустым")
	}

	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/%s/tags/%d", apiClient.GetBaseURL(), entityType, tag.ID)

	// Преобразуем структуру тега в JSON
	tagJSON, err := json.Marshal(tag)
	if err != nil {
		return nil, err
	}

	// Создаем запрос
	req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewBuffer(tagJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var updatedTag Tag
	if err := json.NewDecoder(resp.Body).Decode(&updatedTag); err != nil {
		return nil, err
	}

	return &updatedTag, nil
}

// Delete удаляет тег по его ID для указанного типа сущности.
func Delete(ctx context.Context, apiClient *client.Client, entityType EntityType, tagID int) error {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/%s/tags/%d", apiClient.GetBaseURL(), entityType, tagID)

	// Создаем запрос
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return err
	}

	// Выполняем запрос
	resp, err := apiClient.DoRequest(ctx, req)
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

// LinkEntity связывает сущность с тегами
func LinkEntity(ctx context.Context, apiClient *client.Client, entityType EntityType, entityID int, tags []Tag) error {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/%s/%d/tags", apiClient.GetBaseURL(), entityType, entityID)

	// Преобразуем структуру тегов в JSON
	tagsJSON, err := json.Marshal(tags)
	if err != nil {
		return err
	}

	// Создаем запрос
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(tagsJSON))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	return nil
}

// ListForEntity получает список тегов для указанной сущности
func ListForEntity(ctx context.Context, apiClient *client.Client, entityType EntityType, entityID int) ([]Tag, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/%s/%d/tags", apiClient.GetBaseURL(), entityType, entityID)

	// Создаем запрос
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Выполняем запрос
	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
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
