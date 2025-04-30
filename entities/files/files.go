// Пакет files предоставляет методы для работы с файлами в amoCRM.
package files

import (
	"bytes"
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
	ID          int       `json:"id"`
	UUID        string    `json:"uuid,omitempty"`
	EntityID    int       `json:"entity_id"`
	EntityType  EntityType `json:"entity_type"`
	CreatedBy   int       `json:"created_by,omitempty"`
	UpdatedBy   int       `json:"updated_by,omitempty"`
	CreatedAt   int64     `json:"created_at,omitempty"`
	UpdatedAt   int64     `json:"updated_at,omitempty"`
	Size        int       `json:"size,omitempty"`
	Name        string    `json:"name,omitempty"`
	Type        string    `json:"type,omitempty"`
	Version     int       `json:"version,omitempty"`
	AccountID   int       `json:"account_id,omitempty"`
	Title       string    `json:"title,omitempty"`
	URL         string    `json:"url,omitempty"`
	Download    string    `json:"download_link,omitempty"`
	Preview     string    `json:"preview,omitempty"`
	Links       FileLinks `json:"_links,omitempty"`
}

// FileLinks содержит URL-ссылки для файла
type FileLinks struct {
	Self struct {
		Href string `json:"href"`
	} `json:"self"`
	Download struct {
		Href string `json:"href"`
	} `json:"download"`
}

// FilesResponse представляет ответ API при получении списка файлов
type FilesResponse struct {
	Page     int `json:"page"`
	PerPage  int `json:"per_page"`
	Total    int `json:"total"`
	Embedded struct {
		Files []File `json:"files"`
	} `json:"_embedded"`
}

// UploadFile загружает файл в amoCRM и прикрепляет его к указанной сущности
func UploadFile(apiClient *client.Client, entityType EntityType, entityID int, filePath string) (*File, error) {
	// Открываем файл для чтения
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Получаем информацию о файле
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// Формируем URL для загрузки файла
	uploadURL := fmt.Sprintf("%s/api/v4/%s/%d/files", apiClient.GetBaseURL(), entityType, entityID)

	// Создаем буфер для multipart формы
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Добавляем файл в форму
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	// Закрываем multipart writer
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	// Создаем запрос
	req, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		return nil, err
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Выполняем запрос
	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	// Декодируем ответ
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

// UploadFileByContent загружает файл в amoCRM по содержимому и прикрепляет его к указанной сущности
func UploadFileByContent(apiClient *client.Client, entityType EntityType, entityID int, fileName string, content []byte) (*File, error) {
	// Формируем URL для загрузки файла
	uploadURL := fmt.Sprintf("%s/api/v4/%s/%d/files", apiClient.GetBaseURL(), entityType, entityID)

	// Создаем буфер для multipart формы
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Добавляем файл в форму
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	// Закрываем multipart writer
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	// Создаем запрос
	req, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		return nil, err
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Выполняем запрос
	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	// Декодируем ответ
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

// GetFiles получает список файлов, прикрепленных к сущности
func GetFiles(apiClient *client.Client, entityType EntityType, entityID int, page, limit int) ([]File, error) {
	// Формируем URL для запроса
	baseURL := fmt.Sprintf("%s/api/v4/%s/%d/files", apiClient.GetBaseURL(), entityType, entityID)

	// Добавляем параметры пагинации
	params := url.Values{}
	params.Add("page", strconv.Itoa(page))
	params.Add("limit", strconv.Itoa(limit))

	// Добавляем параметры к URL
	baseURL = baseURL + "?" + params.Encode()

	// Создаем запрос
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, err
	}

	// Выполняем запрос
	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	// Декодируем ответ
	var filesResponse FilesResponse
	if err := json.NewDecoder(resp.Body).Decode(&filesResponse); err != nil {
		return nil, err
	}

	return filesResponse.Embedded.Files, nil
}

// GetFile получает информацию о конкретном файле
func GetFile(apiClient *client.Client, entityType EntityType, entityID, fileID int) (*File, error) {
	// Формируем URL для запроса
	fileURL := fmt.Sprintf("%s/api/v4/%s/%d/files/%d", apiClient.GetBaseURL(), entityType, entityID, fileID)

	// Создаем запрос
	req, err := http.NewRequest("GET", fileURL, nil)
	if err != nil {
		return nil, err
	}

	// Выполняем запрос
	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	// Декодируем ответ
	var file File
	if err := json.NewDecoder(resp.Body).Decode(&file); err != nil {
		return nil, err
	}

	return &file, nil
}

// DeleteFile удаляет файл
func DeleteFile(apiClient *client.Client, entityType EntityType, entityID, fileID int) error {
	// Формируем URL для запроса
	deleteURL := fmt.Sprintf("%s/api/v4/%s/%d/files/%d", apiClient.GetBaseURL(), entityType, entityID, fileID)

	// Создаем запрос
	req, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		return err
	}

	// Выполняем запрос
	resp, err := apiClient.DoRequest(req)
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

// BatchDeleteFiles удаляет несколько файлов одним запросом
func BatchDeleteFiles(apiClient *client.Client, entityType EntityType, entityID int, fileIDs []int) error {
	// Формируем URL для запроса
	deleteURL := fmt.Sprintf("%s/api/v4/%s/%d/files", apiClient.GetBaseURL(), entityType, entityID)

	// Создаем список ID файлов для удаления
	idsStr := make([]string, len(fileIDs))
	for i, id := range fileIDs {
		idsStr[i] = fmt.Sprintf("%d", id)
	}

	// Добавляем параметр с ID файлов к URL
	params := url.Values{}
	params.Add("filter[id]", strings.Join(idsStr, ","))
	deleteURL = deleteURL + "?" + params.Encode()

	// Создаем запрос
	req, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		return err
	}

	// Выполняем запрос
	resp, err := apiClient.DoRequest(req)
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

// DownloadFile скачивает файл и сохраняет его по указанному пути
func DownloadFile(apiClient *client.Client, entityType EntityType, entityID, fileID int, savePath string) error {
	// Получаем информацию о файле
	file, err := GetFile(apiClient, entityType, entityID, fileID)
	if err != nil {
		return err
	}

	// Проверяем, есть ли ссылка для скачивания
	if file.Links.Download.Href == "" {
		return fmt.Errorf("ссылка для скачивания файла не найдена")
	}

	// Создаем запрос для скачивания файла
	downloadURL := fmt.Sprintf("%s%s", apiClient.GetBaseURL(), file.Links.Download.Href)
	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return err
	}

	// Выполняем запрос
	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("неожиданный статус-код при скачивании файла: %d", resp.StatusCode)
	}

	// Создаем файл для сохранения
	outFile, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Копируем содержимое ответа в файл
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// GetDownloadFileURL получает URL для скачивания файла
func GetDownloadFileURL(apiClient *client.Client, entityType EntityType, entityID, fileID int) (string, error) {
	// Получаем информацию о файле
	file, err := GetFile(apiClient, entityType, entityID, fileID)
	if err != nil {
		return "", err
	}

	// Проверяем, есть ли ссылка для скачивания
	if file.Links.Download.Href == "" {
		return "", fmt.Errorf("ссылка для скачивания файла не найдена")
	}

	// Формируем полный URL для скачивания
	downloadURL := fmt.Sprintf("%s%s", apiClient.GetBaseURL(), file.Links.Download.Href)
	return downloadURL, nil
}
