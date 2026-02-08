// Package files предоставляет методы для работы с файлами в amoCRM.
package files

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/chudno/amo_crm_sdk/client"
)

// EntityType определяет тип сущности, к которой прикрепляются файлы
type EntityType string

const (
	// EntityTypeLead тип сущности - Сделка
	EntityTypeLead EntityType = "leads"
	// EntityTypeContact тип сущности - Контакт
	EntityTypeContact EntityType = "contacts"
	// EntityTypeCompany тип сущности - Компания
	EntityTypeCompany EntityType = "companies"
	// EntityTypeCustomers тип сущности - Покупатель
	EntityTypeCustomers EntityType = "customers"
	// EntityTypeCatalogElements тип сущности - Элемент каталога
	EntityTypeCatalogElements EntityType = "catalog_elements"
)

// File представляет структуру файла в amoCRM
type File struct {
	ID         int        `json:"id"`
	UUID       string     `json:"uuid,omitempty"`
	EntityID   int        `json:"entity_id"`
	EntityType EntityType `json:"entity_type"`
	CreatedBy  int        `json:"created_by,omitempty"`
	UpdatedBy  int        `json:"updated_by,omitempty"`
	CreatedAt  int64      `json:"created_at,omitempty"`
	UpdatedAt  int64      `json:"updated_at,omitempty"`
	Size       int        `json:"size,omitempty"`
	Name       string     `json:"name,omitempty"`
	Type       string     `json:"type,omitempty"`
	Version    int        `json:"version,omitempty"`
	AccountID  int        `json:"account_id,omitempty"`
	Title      string     `json:"title,omitempty"`
	URL        string     `json:"url,omitempty"`
	Download   string     `json:"download_link,omitempty"`
	Preview    string     `json:"preview,omitempty"`
	Links      Links      `json:"_links,omitempty"`
}

// Links содержит URL-ссылки для файла
type Links struct {
	Self struct {
		Href string `json:"href"`
	} `json:"self"`
	Download struct {
		Href string `json:"href"`
	} `json:"download"`
}

// ListResponse представляет ответ API при получении списка файлов
type ListResponse struct {
	Page     int `json:"page"`
	PerPage  int `json:"per_page"`
	Total    int `json:"total"`
	Embedded struct {
		Files []File `json:"files"`
	} `json:"_embedded"`
}

// Upload загружает файл в amoCRM и прикрепляет его к указанной сущности
func Upload(ctx context.Context, apiClient *client.Client, entityType EntityType, entityID int, filePath string) (*File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	uploadURL := fmt.Sprintf("%s/api/v4/%s/%d/files", apiClient.GetBaseURL(), entityType, entityID)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", uploadURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var result struct {
		File *File `json:"_embedded"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Если успешно загружено, добавляем информацию о файле
	if result.File != nil {
		result.File.Size = int(fileInfo.Size())
		result.File.Name = filepath.Base(filePath)
		result.File.EntityID = entityID
		result.File.EntityType = entityType
	}

	return result.File, nil
}

// UploadByContent загружает файл в amoCRM по содержимому и прикрепляет его к указанной сущности
func UploadByContent(ctx context.Context, apiClient *client.Client, entityType EntityType, entityID int, fileName string, content []byte) (*File, error) {
	uploadURL := fmt.Sprintf("%s/api/v4/%s/%d/files", apiClient.GetBaseURL(), entityType, entityID)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", uploadURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var result struct {
		File *File `json:"_embedded"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Если успешно загружено, добавляем информацию о файле
	if result.File != nil {
		result.File.Size = len(content)
		result.File.Name = fileName
		result.File.EntityID = entityID
		result.File.EntityType = entityType
	}

	return result.File, nil
}

// List получает список файлов, прикрепленных к сущности
func List(ctx context.Context, apiClient *client.Client, entityType EntityType, entityID int, page, limit int) ([]File, error) {
	baseURL := fmt.Sprintf("%s/api/v4/%s/%d/files", apiClient.GetBaseURL(), entityType, entityID)

	params := url.Values{}
	params.Add("page", strconv.Itoa(page))
	params.Add("limit", strconv.Itoa(limit))

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

	var filesResponse ListResponse
	if err := json.NewDecoder(resp.Body).Decode(&filesResponse); err != nil {
		return nil, err
	}

	return filesResponse.Embedded.Files, nil
}

// Get получает информацию о конкретном файле
func Get(ctx context.Context, apiClient *client.Client, entityType EntityType, entityID, fileID int) (*File, error) {
	fileURL := fmt.Sprintf("%s/api/v4/%s/%d/files/%d", apiClient.GetBaseURL(), entityType, entityID, fileID)

	req, err := http.NewRequestWithContext(ctx, "GET", fileURL, nil)
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

	var file File
	if err := json.NewDecoder(resp.Body).Decode(&file); err != nil {
		return nil, err
	}

	return &file, nil
}

// Delete удаляет файл
func Delete(ctx context.Context, apiClient *client.Client, entityType EntityType, entityID, fileID int) error {
	deleteURL := fmt.Sprintf("%s/api/v4/%s/%d/files/%d", apiClient.GetBaseURL(), entityType, entityID, fileID)

	req, err := http.NewRequestWithContext(ctx, "DELETE", deleteURL, nil)
	if err != nil {
		return err
	}

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	return nil
}

// BatchDelete удаляет несколько файлов одним запросом
func BatchDelete(ctx context.Context, apiClient *client.Client, entityType EntityType, entityID int, fileIDs []int) error {
	deleteURL := fmt.Sprintf("%s/api/v4/%s/%d/files", apiClient.GetBaseURL(), entityType, entityID)

	idsStr := make([]string, len(fileIDs))
	for i, id := range fileIDs {
		idsStr[i] = fmt.Sprintf("%d", id)
	}

	params := url.Values{}
	params.Add("filter[id]", strings.Join(idsStr, ","))
	deleteURL = deleteURL + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, "DELETE", deleteURL, nil)
	if err != nil {
		return err
	}

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	return nil
}

// Download скачивает файл и сохраняет его по указанному пути
func Download(ctx context.Context, apiClient *client.Client, entityType EntityType, entityID, fileID int, savePath string) error {
	file, err := Get(ctx, apiClient, entityType, entityID, fileID)
	if err != nil {
		return err
	}

	// Проверяем, есть ли ссылка для скачивания
	if file.Links.Download.Href == "" {
		return fmt.Errorf("ссылка для скачивания файла не найдена")
	}

	downloadURL := fmt.Sprintf("%s%s", apiClient.GetBaseURL(), file.Links.Download.Href)
	req, err := http.NewRequestWithContext(ctx, "GET", downloadURL, nil)
	if err != nil {
		return err
	}

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("неожиданный статус-код при скачивании файла: %d", resp.StatusCode)
	}

	outFile, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer func() { _ = outFile.Close() }()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// GetDownloadURL получает URL для скачивания файла
func GetDownloadURL(ctx context.Context, apiClient *client.Client, entityType EntityType, entityID, fileID int) (string, error) {
	file, err := Get(ctx, apiClient, entityType, entityID, fileID)
	if err != nil {
		return "", err
	}

	// Проверяем, есть ли ссылка для скачивания
	if file.Links.Download.Href == "" {
		return "", fmt.Errorf("ссылка для скачивания файла не найдена")
	}

	downloadURL := fmt.Sprintf("%s%s", apiClient.GetBaseURL(), file.Links.Download.Href)
	return downloadURL, nil
}
