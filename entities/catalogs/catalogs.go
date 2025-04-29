// Пакет catalogs предоставляет методы для взаимодействия с сущностями "Каталоги" в API amoCRM.
package catalogs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/chudno/amo_crm_sdk/client"
)

// Catalog представляет собой структуру каталога в amoCRM.
type Catalog struct {
	ID                 int            `json:"id,omitempty"`
	Name               string         `json:"name"`
	CreatedBy          int            `json:"created_by,omitempty"`
	UpdatedBy          int            `json:"updated_by,omitempty"`
	CreatedAt          int64          `json:"created_at,omitempty"`
	UpdatedAt          int64          `json:"updated_at,omitempty"`
	Sort               int            `json:"sort,omitempty"`
	Type               string         `json:"type,omitempty"`
	Can                *CatalogAccess `json:"can,omitempty"`
	CustomFieldsConfig []CustomField  `json:"custom_fields_config,omitempty"`
}

// CustomField представляет пользовательское поле для каталога
type CustomField struct {
	ID           int    `json:"id,omitempty"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	IsAPIOnly    bool   `json:"is_api_only,omitempty"`
	IsRequired   bool   `json:"is_required,omitempty"`
	IsMultiple   bool   `json:"is_multiple,omitempty"`
	IsSystem     bool   `json:"is_system,omitempty"`
	Sort         int    `json:"sort,omitempty"`
	Code         string `json:"code,omitempty"`
	GroupID      int    `json:"group_id,omitempty"`
	EntityType   string `json:"entity_type,omitempty"`
	NeedsAPICode bool   `json:"needs_api_code,omitempty"`
}

// CatalogAccess представляет права доступа к каталогу
type CatalogAccess struct {
	View   bool `json:"view"`
	Edit   bool `json:"edit"`
	Add    bool `json:"add"`
	Delete bool `json:"delete"`
	Export bool `json:"export"`
}

// CatalogsResponse представляет ответ от API при получении списка каталогов
type CatalogsResponse struct {
	Page     int `json:"page"`
	PerPage  int `json:"per_page"`
	Total    int `json:"total"`
	Embedded struct {
		Catalogs []Catalog `json:"catalogs"`
	} `json:"_embedded"`
}

// CatalogType представляет типы каталогов
type CatalogType string

const (
	// CatalogTypeRegular - обычный каталог, создаваемый пользователем
	CatalogTypeRegular CatalogType = "regular"
	// CatalogTypeContacts - системный каталог для контактов
	CatalogTypeContacts CatalogType = "contacts"
	// CatalogTypeCompanies - системный каталог для компаний
	CatalogTypeCompanies CatalogType = "companies"
)

// GetCatalogs получает список каталогов с возможностью пагинации и фильтрации.
func GetCatalogs(apiClient *client.Client, page, limit int, filter map[string]string) ([]Catalog, error) {
	// Формируем базовый URL
	baseURL := fmt.Sprintf("%s/api/v4/catalogs", apiClient.GetBaseURL())

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

	var catalogs CatalogsResponse
	if err := json.NewDecoder(resp.Body).Decode(&catalogs); err != nil {
		return nil, err
	}

	return catalogs.Embedded.Catalogs, nil
}

// CreateCatalog создает новый каталог.
func CreateCatalog(apiClient *client.Client, catalog *Catalog) (*Catalog, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/catalogs", apiClient.GetBaseURL())

	// Преобразуем структуру каталога в JSON
	catalogJSON, err := json.Marshal(catalog)
	if err != nil {
		return nil, err
	}

	// Создаем запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(catalogJSON))
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

	var createdCatalog Catalog
	if err := json.NewDecoder(resp.Body).Decode(&createdCatalog); err != nil {
		return nil, err
	}

	return &createdCatalog, nil
}

// GetCatalog получает информацию о каталоге по его ID.
func GetCatalog(apiClient *client.Client, catalogID int) (*Catalog, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/catalogs/%d", apiClient.GetBaseURL(), catalogID)

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

	var catalog Catalog
	if err := json.NewDecoder(resp.Body).Decode(&catalog); err != nil {
		return nil, err
	}

	return &catalog, nil
}

// UpdateCatalog обновляет информацию о каталоге по его ID.
func UpdateCatalog(apiClient *client.Client, catalog *Catalog) (*Catalog, error) {
	if catalog.ID == 0 {
		return nil, fmt.Errorf("ID каталога не может быть пустым")
	}

	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/catalogs/%d", apiClient.GetBaseURL(), catalog.ID)

	// Преобразуем структуру каталога в JSON
	catalogJSON, err := json.Marshal(catalog)
	if err != nil {
		return nil, err
	}

	// Создаем запрос
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(catalogJSON))
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

	var updatedCatalog Catalog
	if err := json.NewDecoder(resp.Body).Decode(&updatedCatalog); err != nil {
		return nil, err
	}

	return &updatedCatalog, nil
}

// DeleteCatalog удаляет каталог по его ID.
func DeleteCatalog(apiClient *client.Client, catalogID int) error {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/catalogs/%d", apiClient.GetBaseURL(), catalogID)

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

// AddCustomFieldToCatalog добавляет пользовательское поле в каталог
func AddCustomFieldToCatalog(apiClient *client.Client, catalogID int, customField *CustomField) (*CustomField, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/catalogs/%d/custom_fields", apiClient.GetBaseURL(), catalogID)

	// Преобразуем структуру поля в JSON
	fieldJSON, err := json.Marshal(customField)
	if err != nil {
		return nil, err
	}

	// Создаем запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(fieldJSON))
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

	var createdField CustomField
	if err := json.NewDecoder(resp.Body).Decode(&createdField); err != nil {
		return nil, err
	}

	return &createdField, nil
}

// GetCatalogCustomFields получает список пользовательских полей каталога
func GetCatalogCustomFields(apiClient *client.Client, catalogID int) ([]CustomField, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/catalogs/%d/custom_fields", apiClient.GetBaseURL(), catalogID)

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

	var fieldsResponse struct {
		Embedded struct {
			CustomFields []CustomField `json:"custom_fields"`
		} `json:"_embedded"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&fieldsResponse); err != nil {
		return nil, err
	}

	return fieldsResponse.Embedded.CustomFields, nil
}

// GetCatalogCustomField получает информацию о пользовательском поле каталога по ID
func GetCatalogCustomField(apiClient *client.Client, catalogID, fieldID int) (*CustomField, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/catalogs/%d/custom_fields/%d", apiClient.GetBaseURL(), catalogID, fieldID)

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

	var field CustomField
	if err := json.NewDecoder(resp.Body).Decode(&field); err != nil {
		return nil, err
	}

	return &field, nil
}

// UpdateCatalogCustomField обновляет пользовательское поле каталога
func UpdateCatalogCustomField(apiClient *client.Client, catalogID int, field *CustomField) (*CustomField, error) {
	if field.ID == 0 {
		return nil, fmt.Errorf("ID поля не может быть пустым")
	}

	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/catalogs/%d/custom_fields/%d", apiClient.GetBaseURL(), catalogID, field.ID)

	// Преобразуем структуру поля в JSON
	fieldJSON, err := json.Marshal(field)
	if err != nil {
		return nil, err
	}

	// Создаем запрос
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(fieldJSON))
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

	var updatedField CustomField
	if err := json.NewDecoder(resp.Body).Decode(&updatedField); err != nil {
		return nil, err
	}

	return &updatedField, nil
}

// DeleteCatalogCustomField удаляет пользовательское поле каталога
func DeleteCatalogCustomField(apiClient *client.Client, catalogID, fieldID int) error {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/api/v4/catalogs/%d/custom_fields/%d", apiClient.GetBaseURL(), catalogID, fieldID)

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
