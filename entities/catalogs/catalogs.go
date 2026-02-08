// Package catalogs предоставляет методы для взаимодействия с сущностями "Каталоги" в API amoCRM.
package catalogs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

// List получает список каталогов с возможностью пагинации и фильтрации.
func List(ctx context.Context, apiClient *client.Client, page, limit int, filter map[string]string) ([]Catalog, error) {
	baseURL := fmt.Sprintf("%s/api/v4/catalogs", apiClient.GetBaseURL())

	params := url.Values{}
	params.Add("page", fmt.Sprintf("%d", page))
	params.Add("limit", fmt.Sprintf("%d", limit))

	for key, value := range filter {
		params.Add(key, value)
	}

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

	var catalogs CatalogsResponse
	if err := json.NewDecoder(resp.Body).Decode(&catalogs); err != nil {
		return nil, err
	}

	return catalogs.Embedded.Catalogs, nil
}

// Create создает новый каталог.
func Create(ctx context.Context, apiClient *client.Client, catalog *Catalog) (*Catalog, error) {
	apiURL := fmt.Sprintf("%s/api/v4/catalogs", apiClient.GetBaseURL())

	catalogJSON, err := json.Marshal([]*Catalog{catalog})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(catalogJSON))
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
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("неожиданный статус-код: %d, ответ: %s", resp.StatusCode, body)
	}

	var response struct {
		Embedded struct {
			Catalogs []*Catalog `json:"catalogs"`
		} `json:"_embedded"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if len(response.Embedded.Catalogs) == 0 {
		return nil, fmt.Errorf("не удалось создать каталог")
	}

	return response.Embedded.Catalogs[0], nil
}

// Get получает информацию о каталоге по его ID.
func Get(ctx context.Context, apiClient *client.Client, catalogID int) (*Catalog, error) {
	url := fmt.Sprintf("%s/api/v4/catalogs/%d", apiClient.GetBaseURL(), catalogID)

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

	var catalog Catalog
	if err := json.NewDecoder(resp.Body).Decode(&catalog); err != nil {
		return nil, err
	}

	return &catalog, nil
}

// Update обновляет информацию о каталоге по его ID.
func Update(ctx context.Context, apiClient *client.Client, catalog *Catalog) (*Catalog, error) {
	if catalog.ID == 0 {
		return nil, fmt.Errorf("ID каталога не может быть пустым")
	}

	url := fmt.Sprintf("%s/api/v4/catalogs/%d", apiClient.GetBaseURL(), catalog.ID)

	catalogJSON, err := json.Marshal(catalog)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewBuffer(catalogJSON))
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

	var updatedCatalog Catalog
	if err := json.NewDecoder(resp.Body).Decode(&updatedCatalog); err != nil {
		return nil, err
	}

	return &updatedCatalog, nil
}

// AddCustomField добавляет пользовательское поле в каталог
func AddCustomField(ctx context.Context, apiClient *client.Client, catalogID int, customField *CustomField) (*CustomField, error) {
	url := fmt.Sprintf("%s/api/v4/catalogs/%d/custom_fields", apiClient.GetBaseURL(), catalogID)

	fieldJSON, err := json.Marshal(customField)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(fieldJSON))
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

	var createdField CustomField
	if err := json.NewDecoder(resp.Body).Decode(&createdField); err != nil {
		return nil, err
	}

	return &createdField, nil
}

// ListCustomFields получает список пользовательских полей каталога
func ListCustomFields(ctx context.Context, apiClient *client.Client, catalogID int) ([]CustomField, error) {
	url := fmt.Sprintf("%s/api/v4/catalogs/%d/custom_fields", apiClient.GetBaseURL(), catalogID)

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

// GetCustomField получает информацию о пользовательском поле каталога по ID
func GetCustomField(ctx context.Context, apiClient *client.Client, catalogID, fieldID int) (*CustomField, error) {
	url := fmt.Sprintf("%s/api/v4/catalogs/%d/custom_fields/%d", apiClient.GetBaseURL(), catalogID, fieldID)

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

	var field CustomField
	if err := json.NewDecoder(resp.Body).Decode(&field); err != nil {
		return nil, err
	}

	return &field, nil
}

// UpdateCustomField обновляет пользовательское поле каталога
func UpdateCustomField(ctx context.Context, apiClient *client.Client, catalogID int, field *CustomField) (*CustomField, error) {
	if field.ID == 0 {
		return nil, fmt.Errorf("ID поля не может быть пустым")
	}

	url := fmt.Sprintf("%s/api/v4/catalogs/%d/custom_fields/%d", apiClient.GetBaseURL(), catalogID, field.ID)

	fieldJSON, err := json.Marshal(field)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewBuffer(fieldJSON))
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

	var updatedField CustomField
	if err := json.NewDecoder(resp.Body).Decode(&updatedField); err != nil {
		return nil, err
	}

	return &updatedField, nil
}

// DeleteCustomField удаляет пользовательское поле каталога
func DeleteCustomField(ctx context.Context, apiClient *client.Client, catalogID, fieldID int) error {
	url := fmt.Sprintf("%s/api/v4/catalogs/%d/custom_fields/%d", apiClient.GetBaseURL(), catalogID, fieldID)

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
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
