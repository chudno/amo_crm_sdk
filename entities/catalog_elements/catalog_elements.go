// Пакет catalog_elements предоставляет методы для взаимодействия с сущностями "Элементы каталогов" в API amoCRM.
package catalog_elements

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/chudno/amo_crm_sdk/client"
)

// CatalogElement представляет собой структуру элемента каталога в amoCRM.
type CatalogElement struct {
	ID                 int                     `json:"id,omitempty"`
	Name               string                  `json:"name"`
	CreatedBy          int                     `json:"created_by,omitempty"`
	UpdatedBy          int                     `json:"updated_by,omitempty"`
	CreatedAt          int64                   `json:"created_at,omitempty"`
	UpdatedAt          int64                   `json:"updated_at,omitempty"`
	CatalogID          int                     `json:"catalog_id"`
	CustomFieldsValues []CustomFieldValue      `json:"custom_fields_values,omitempty"`
	AccountID          int                     `json:"account_id,omitempty"`
	IsDeleted          bool                    `json:"is_deleted,omitempty"`
	QuantityBounded    bool                    `json:"quantity_bounded,omitempty"`
	QuantityRaw        int                     `json:"quantity_raw,omitempty"`
	Embedded           *CatalogElementEmbedded `json:"_embedded,omitempty"`
	Links              *CatalogElementLinks    `json:"_links,omitempty"`
}

// CustomFieldValue представляет значение пользовательского поля элемента каталога
type CustomFieldValue struct {
	FieldID   int              `json:"field_id"`
	FieldName string           `json:"field_name,omitempty"`
	FieldCode string           `json:"field_code,omitempty"`
	FieldType string           `json:"field_type,omitempty"`
	Values    []FieldValueItem `json:"values"`
}

// FieldValueItem представляет значение поля
type FieldValueItem struct {
	Value  interface{} `json:"value"`
	EnumID int         `json:"enum_id,omitempty"`
	Enum   string      `json:"enum,omitempty"`
}

// CatalogElementEmbedded содержит связанные с элементом каталога сущности
type CatalogElementEmbedded struct {
	Tags []Tag `json:"tags,omitempty"`
}

// Tag представляет тег элемента каталога
type Tag struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color,omitempty"`
}

// CatalogElementLinks содержит ссылки на связанные с элементом каталога ресурсы
type CatalogElementLinks struct {
	Self struct {
		Href string `json:"href"`
	} `json:"self"`
}

// CatalogElementsResponse представляет ответ от API при получении списка элементов каталога
type CatalogElementsResponse struct {
	Page     int `json:"page"`
	PerPage  int `json:"per_page"`
	Total    int `json:"total"`
	Embedded struct {
		Elements []CatalogElement `json:"elements"`
	} `json:"_embedded"`
}

// WithOption представляет опцию для запроса с добавлением связанных сущностей
type WithOption string

const (
	// WithTags - получить теги элементов каталога
	WithTags WithOption = "tags"
	// WithFullLinkedEntities - включить информацию о связанных сущностях
	WithFullLinkedEntities WithOption = "full_linked_entities"
)

// GetCatalogElements получает список элементов каталога с возможностью пагинации и фильтрации.
func GetCatalogElements(apiClient *client.Client, catalogID, page, limit int, filter map[string]string, withOptions ...WithOption) ([]CatalogElement, error) {
	// Формируем базовый URL
	baseURL := fmt.Sprintf("%s/api/v4/catalogs/%d/elements", apiClient.GetBaseURL(), catalogID)

	// Добавляем параметры запроса
	params := url.Values{}
	params.Add("page", fmt.Sprintf("%d", page))
	params.Add("limit", fmt.Sprintf("%d", limit))

	// Добавляем фильтры, если они указаны
	if filter != nil {
		for key, value := range filter {
			params.Add(key, value)
		}
	}

	// Добавляем параметр with, если указаны withOptions
	if len(withOptions) > 0 {
		var withValues []string
		for _, opt := range withOptions {
			withValues = append(withValues, string(opt))
		}
		params.Add("with", stringsJoin(withValues, ","))
	}

	// Добавляем параметры к URL
	baseURL = baseURL + "?" + params.Encode()

	// Создаем запрос
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var elements CatalogElementsResponse
	if err := json.NewDecoder(resp.Body).Decode(&elements); err != nil {
		return nil, err
	}

	return elements.Embedded.Elements, nil
}

// stringsJoin объединяет срез строк с указанным разделителем
func stringsJoin(strings []string, sep string) string {
	if len(strings) == 0 {
		return ""
	}

	result := strings[0]
	for i := 1; i < len(strings); i++ {
		result += sep + strings[i]
	}

	return result
}

// CreateCatalogElement создает новый элемент каталога.
func CreateCatalogElement(apiClient *client.Client, catalogID int, element *CatalogElement) (*CatalogElement, error) {
	// Проверяем, что указан ID каталога
	element.CatalogID = catalogID

	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/catalogs/%d/elements", apiClient.GetBaseURL(), catalogID)

	// Преобразуем структуру элемента в JSON
	elementJSON, err := json.Marshal([]*CatalogElement{element})
	if err != nil {
		return nil, err
	}

	// Создаем запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(elementJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

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

	var response struct {
		Embedded struct {
			Elements []CatalogElement `json:"elements"`
		} `json:"_embedded"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if len(response.Embedded.Elements) == 0 {
		return nil, fmt.Errorf("не удалось создать элемент каталога")
	}

	return &response.Embedded.Elements[0], nil
}

// CreateCatalogElements создает несколько элементов каталога за один запрос.
func CreateCatalogElements(apiClient *client.Client, catalogID int, elements []CatalogElement) ([]CatalogElement, error) {
	// Проверяем, что указан ID каталога для всех элементов
	for i := range elements {
		elements[i].CatalogID = catalogID
	}

	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/catalogs/%d/elements", apiClient.GetBaseURL(), catalogID)

	// Преобразуем структуры элементов в JSON
	elementsJSON, err := json.Marshal(elements)
	if err != nil {
		return nil, err
	}

	// Создаем запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(elementsJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

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

	var response struct {
		Embedded struct {
			Elements []CatalogElement `json:"elements"`
		} `json:"_embedded"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Embedded.Elements, nil
}

// GetCatalogElement получает информацию об элементе каталога по его ID.
func GetCatalogElement(apiClient *client.Client, catalogID, elementID int, withOptions ...WithOption) (*CatalogElement, error) {
	// Формируем URL для запроса
	baseURL := fmt.Sprintf("%s/api/v4/catalogs/%d/elements/%d", apiClient.GetBaseURL(), catalogID, elementID)

	// Добавляем параметр with, если указаны withOptions
	if len(withOptions) > 0 {
		params := url.Values{}
		var withValues []string
		for _, opt := range withOptions {
			withValues = append(withValues, string(opt))
		}
		params.Add("with", stringsJoin(withValues, ","))
		baseURL = baseURL + "?" + params.Encode()
	}

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

	var element CatalogElement
	if err := json.NewDecoder(resp.Body).Decode(&element); err != nil {
		return nil, err
	}

	return &element, nil
}

// UpdateCatalogElement обновляет информацию об элементе каталога по его ID.
func UpdateCatalogElement(apiClient *client.Client, catalogID int, element *CatalogElement) (*CatalogElement, error) {
	if element.ID == 0 {
		return nil, fmt.Errorf("ID элемента каталога не может быть пустым")
	}

	// Проверяем, что указан ID каталога
	element.CatalogID = catalogID

	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/catalogs/%d/elements/%d", apiClient.GetBaseURL(), catalogID, element.ID)

	// Преобразуем структуру элемента в JSON
	elementJSON, err := json.Marshal(element)
	if err != nil {
		return nil, err
	}

	// Создаем запрос
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(elementJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

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

	var updatedElement CatalogElement
	if err := json.NewDecoder(resp.Body).Decode(&updatedElement); err != nil {
		return nil, err
	}

	return &updatedElement, nil
}

// UpdateCatalogElements обновляет информацию о нескольких элементах каталога за один запрос.
func UpdateCatalogElements(apiClient *client.Client, catalogID int, elements []CatalogElement) ([]CatalogElement, error) {
	// Проверяем, что у всех элементов есть ID
	for i := range elements {
		if elements[i].ID == 0 {
			return nil, fmt.Errorf("ID элемента каталога не может быть пустым")
		}
		// Проверяем, что указан ID каталога
		elements[i].CatalogID = catalogID
	}

	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/catalogs/%d/elements", apiClient.GetBaseURL(), catalogID)

	// Преобразуем структуры элементов в JSON
	elementsJSON, err := json.Marshal(elements)
	if err != nil {
		return nil, err
	}

	// Создаем запрос
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(elementsJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

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

	var response struct {
		Embedded struct {
			Elements []CatalogElement `json:"elements"`
		} `json:"_embedded"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Embedded.Elements, nil
}

// DeleteCatalogElement удаляет элемент каталога по его ID.
func DeleteCatalogElement(apiClient *client.Client, catalogID, elementID int) error {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/catalogs/%d/elements/%d", apiClient.GetBaseURL(), catalogID, elementID)

	// Создаем запрос
	req, err := http.NewRequest("DELETE", url, nil)
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

// BatchDeleteCatalogElements удаляет несколько элементов каталога за один запрос.
func BatchDeleteCatalogElements(apiClient *client.Client, catalogID int, elementIDs []int) error {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/catalogs/%d/elements", apiClient.GetBaseURL(), catalogID)

	// Формируем тело запроса
	type deleteRequest struct {
		ID int `json:"id"`
	}

	var requests []deleteRequest
	for _, id := range elementIDs {
		requests = append(requests, deleteRequest{ID: id})
	}

	// Преобразуем запрос в JSON
	requestJSON, err := json.Marshal(requests)
	if err != nil {
		return err
	}

	// Создаем запрос
	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(requestJSON))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

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

// LinkCatalogElementWithTags связывает элемент каталога с тегами.
func LinkCatalogElementWithTags(apiClient *client.Client, catalogID, elementID int, tags []Tag) error {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/catalogs/%d/elements/%d/tags", apiClient.GetBaseURL(), catalogID, elementID)

	// Преобразуем теги в JSON
	tagsJSON, err := json.Marshal(tags)
	if err != nil {
		return err
	}

	// Создаем запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(tagsJSON))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := apiClient.DoRequest(req)
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

// GetCatalogElementTags получает теги элемента каталога.
func GetCatalogElementTags(apiClient *client.Client, catalogID, elementID int) ([]Tag, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/catalogs/%d/elements/%d/tags", apiClient.GetBaseURL(), catalogID, elementID)

	// Создаем запрос
	req, err := http.NewRequest("GET", url, nil)
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

	var response struct {
		Embedded struct {
			Tags []Tag `json:"tags"`
		} `json:"_embedded"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Embedded.Tags, nil
}
