// Пакет sources предоставляет методы для работы с источниками сделок в amoCRM.
package sources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/chudno/amo_crm_sdk/client"
)

// Requester - интерфейс для выполнения HTTP-запросов, используется для тестирования.
type Requester interface {
	DoRequest(req *http.Request) (*http.Response, error)
	GetBaseURL() string
}

// Source представляет источник сделок в amoCRM.
type Source struct {
	ID            int         `json:"id,omitempty"`
	Name          string      `json:"name"`
	Type          string      `json:"type,omitempty"`
	Default       bool        `json:"default,omitempty"`
	CreatedAt     int64       `json:"created_at,omitempty"`
	UpdatedAt     int64       `json:"updated_at,omitempty"`
	Deleted       bool        `json:"deleted,omitempty"`
	EffectiveFrom int64       `json:"effective_from,omitempty"`
	EffectiveTo   int64       `json:"effective_to,omitempty"`
	Pipeline      *Pipeline   `json:"pipeline,omitempty"`
	Services      []Service   `json:"services,omitempty"`
	External      *External   `json:"external,omitempty"`
	Params        interface{} `json:"params,omitempty"`
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
	ID             string      `json:"id,omitempty"`
	Service        string      `json:"service,omitempty"`
	ExternalParams interface{} `json:"external_params,omitempty"`
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

// GetSources получает список источников сделок с поддержкой фильтрации и пагинации.
//
// Пример использования:
//
//	filter := map[string]string{
//		"filter[name]": "Реклама",
//	}
//	sources, err := sources.GetSources(apiClient, 1, 50, sources.WithFilter(filter))
func GetSources(apiClient *client.Client, page, limit int, options ...WithOption) ([]Source, error) {
	return GetSourcesWithRequester(apiClient, page, limit, options...)
}

// GetSourcesWithRequester получает список источников с использованием интерфейса Requester.
func GetSourcesWithRequester(requester Requester, page, limit int, options ...WithOption) ([]Source, error) {
	// Формируем URL для запроса
	baseURL := fmt.Sprintf("%s/api/v4/sources", requester.GetBaseURL())

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
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	// Выполняем запрос
	resp, err := requester.DoRequest(req)
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
			Sources []Source `json:"sources"`
		} `json:"_embedded"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Embedded.Sources, nil
}

// GetSource получает информацию о конкретном источнике по ID.
//
// Пример использования:
//
//	sourceInfo, err := sources.GetSource(apiClient, 123)
func GetSource(apiClient *client.Client, id int) (*Source, error) {
	return GetSourceWithRequester(apiClient, id)
}

// GetSourceWithRequester получает информацию о конкретном источнике с использованием интерфейса Requester.
func GetSourceWithRequester(requester Requester, id int) (*Source, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/sources/%d", requester.GetBaseURL(), id)

	// Создаем запрос
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Выполняем запрос
	resp, err := requester.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	// Декодируем ответ
	var sourceInfo Source
	if err := json.NewDecoder(resp.Body).Decode(&sourceInfo); err != nil {
		return nil, err
	}

	return &sourceInfo, nil
}

// CreateSource создает новый источник сделок.
//
// Пример использования:
//
//	newSource := &sources.Source{
//		Name: "Новый источник",
//		Type: "other",
//	}
//	createdSource, err := sources.CreateSource(apiClient, newSource)
func CreateSource(apiClient *client.Client, sourceData *Source) (*Source, error) {
	return CreateSourceWithRequester(apiClient, sourceData)
}

// CreateSourceWithRequester создает новый источник с использованием интерфейса Requester.
func CreateSourceWithRequester(requester Requester, sourceData *Source) (*Source, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/sources", requester.GetBaseURL())

	// Подготавливаем данные для запроса
	data, err := json.Marshal(sourceData)
	if err != nil {
		return nil, err
	}

	// Создаем запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := requester.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	// Декодируем ответ
	var createdSource Source
	if err := json.NewDecoder(resp.Body).Decode(&createdSource); err != nil {
		return nil, err
	}

	return &createdSource, nil
}

// UpdateSource обновляет существующий источник сделок.
//
// Пример использования:
//
//	sourceUpdate := &sources.Source{
//		ID:   123,
//		Name: "Обновленный источник",
//	}
//	updatedSource, err := sources.UpdateSource(apiClient, sourceUpdate)
func UpdateSource(apiClient *client.Client, sourceData *Source) (*Source, error) {
	return UpdateSourceWithRequester(apiClient, sourceData)
}

// UpdateSourceWithRequester обновляет существующий источник с использованием интерфейса Requester.
func UpdateSourceWithRequester(requester Requester, sourceData *Source) (*Source, error) {
	if sourceData.ID == 0 {
		return nil, fmt.Errorf("ID источника не указан")
	}

	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/sources/%d", requester.GetBaseURL(), sourceData.ID)

	// Подготавливаем данные для запроса
	data, err := json.Marshal(sourceData)
	if err != nil {
		return nil, err
	}

	// Создаем запрос
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := requester.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	// Декодируем ответ
	var updatedSource Source
	if err := json.NewDecoder(resp.Body).Decode(&updatedSource); err != nil {
		return nil, err
	}

	return &updatedSource, nil
}

// DeleteSource удаляет источник по ID.
//
// Пример использования:
//
//	err := sources.DeleteSource(apiClient, 123)
func DeleteSource(apiClient *client.Client, id int) error {
	return DeleteSourceWithRequester(apiClient, id)
}

// DeleteSourceWithRequester удаляет источник с использованием интерфейса Requester.
func DeleteSourceWithRequester(requester Requester, id int) error {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/sources/%d", requester.GetBaseURL(), id)

	// Создаем запрос
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	// Выполняем запрос
	resp, err := requester.DoRequest(req)
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

// SetSourceDefault устанавливает источник как используемый по умолчанию.
//
// Пример использования:
//
//	updatedSource, err := sources.SetSourceDefault(apiClient, 123)
func SetSourceDefault(apiClient *client.Client, id int) (*Source, error) {
	return SetSourceDefaultWithRequester(apiClient, id)
}

// SetSourceDefaultWithRequester устанавливает источник как используемый по умолчанию с использованием интерфейса Requester.
func SetSourceDefaultWithRequester(requester Requester, id int) (*Source, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/sources/%d/default", requester.GetBaseURL(), id)

	// Создаем запрос
	req, err := http.NewRequest("PATCH", url, nil)
	if err != nil {
		return nil, err
	}

	// Выполняем запрос
	resp, err := requester.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	// Декодируем ответ
	var updatedSource Source
	if err := json.NewDecoder(resp.Body).Decode(&updatedSource); err != nil {
		return nil, err
	}

	return &updatedSource, nil
}

// GetSourceServices получает список сервисов, доступных для источников сделок.
//
// Пример использования:
//
//	services, err := sources.GetSourceServices(apiClient)
func GetSourceServices(apiClient *client.Client) ([]Service, error) {
	return GetSourceServicesWithRequester(apiClient)
}

// GetSourceServicesWithRequester получает список сервисов с использованием интерфейса Requester.
func GetSourceServicesWithRequester(requester Requester) ([]Service, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/sources/services", requester.GetBaseURL())

	// Создаем запрос
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Выполняем запрос
	resp, err := requester.DoRequest(req)
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
			Services []Service `json:"services"`
		} `json:"_embedded"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Embedded.Services, nil
}

// LinkSourceToPipeline связывает источник с воронкой.
//
// Пример использования:
//
//	linkedSource, err := sources.LinkSourceToPipeline(apiClient, 123, 456)
func LinkSourceToPipeline(apiClient *client.Client, sourceID, pipelineID int) (*Source, error) {
	return LinkSourceToPipelineWithRequester(apiClient, sourceID, pipelineID)
}

// LinkSourceToPipelineWithRequester связывает источник с воронкой с использованием интерфейса Requester.
func LinkSourceToPipelineWithRequester(requester Requester, sourceID, pipelineID int) (*Source, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/sources/%d/pipeline", requester.GetBaseURL(), sourceID)

	// Подготавливаем данные для запроса
	data, err := json.Marshal(map[string]int{
		"pipeline_id": pipelineID,
	})
	if err != nil {
		return nil, err
	}

	// Создаем запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := requester.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	// Декодируем ответ
	var linkedSource Source
	if err := json.NewDecoder(resp.Body).Decode(&linkedSource); err != nil {
		return nil, err
	}

	return &linkedSource, nil
}

// UnlinkSourceFromPipeline удаляет связь источника с воронкой.
//
// Пример использования:
//
//	unlinkedSource, err := sources.UnlinkSourceFromPipeline(apiClient, 123, 456)
func UnlinkSourceFromPipeline(apiClient *client.Client, sourceID, pipelineID int) (*Source, error) {
	return UnlinkSourceFromPipelineWithRequester(apiClient, sourceID, pipelineID)
}

// UnlinkSourceFromPipelineWithRequester удаляет связь источника с воронкой с использованием интерфейса Requester.
func UnlinkSourceFromPipelineWithRequester(requester Requester, sourceID, pipelineID int) (*Source, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/sources/%d/pipeline/%d", requester.GetBaseURL(), sourceID, pipelineID)

	// Создаем запрос
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	// Выполняем запрос
	resp, err := requester.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	// Декодируем ответ
	var unlinkedSource Source
	if err := json.NewDecoder(resp.Body).Decode(&unlinkedSource); err != nil {
		return nil, err
	}

	return &unlinkedSource, nil
}
