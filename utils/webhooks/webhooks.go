// Пакет webhooks предоставляет методы для взаимодействия с вебхуками в API amoCRM.
package webhooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/chudno/amo_crm_sdk/client"
	"net/http"
	"net/url"
)

// Webhook представляет собой структуру вебхука в amoCRM.
type Webhook struct {
	ID          int              `json:"id,omitempty"`
	Destination string           `json:"destination"`
	Settings    *WebhookSettings `json:"settings,omitempty"`
	CreatedAt   int64            `json:"created_at,omitempty"`
	UpdatedAt   int64            `json:"updated_at,omitempty"`
	CreatedBy   int              `json:"created_by,omitempty"`
	AccountID   int              `json:"account_id,omitempty"`
}

// WebhookSettings содержит настройки вебхука
type WebhookSettings struct {
	Entities []string `json:"events"`
	Actions  []string `json:"actions"`
}

// Entity определяет типы сущностей для вебхуков
const (
	EntityLead     = "leads"
	EntityContact  = "contacts"
	EntityCompany  = "companies"
	EntityCustomer = "customers"
	EntityTask     = "tasks"
)

// Action определяет типы действий для вебхуков
const (
	ActionAdd          = "add"
	ActionUpdate       = "update"
	ActionDelete       = "delete"
	ActionRestore      = "restore"
	ActionStatusChange = "status"
)

// GetWebhook получает вебхук по его ID.
func GetWebhook(apiClient *client.Client, webhookID int) (*Webhook, error) {
	url := fmt.Sprintf("%s/api/v4/webhooks/%d", apiClient.GetBaseURL(), webhookID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var webhook Webhook
	if err := json.NewDecoder(resp.Body).Decode(&webhook); err != nil {
		return nil, err
	}

	return &webhook, nil
}

// CreateWebhook создает новый вебхук в amoCRM.
func CreateWebhook(apiClient *client.Client, webhook *Webhook) (*Webhook, error) {
	url := fmt.Sprintf("%s/api/v4/webhooks", apiClient.GetBaseURL())

	webhookData, err := json.Marshal(webhook)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(webhookData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var createdWebhook Webhook
	if err := json.NewDecoder(resp.Body).Decode(&createdWebhook); err != nil {
		return nil, err
	}

	return &createdWebhook, nil
}

// UpdateWebhook обновляет существующий вебхук в amoCRM.
func UpdateWebhook(apiClient *client.Client, webhook *Webhook) (*Webhook, error) {
	if webhook.ID == 0 {
		return nil, fmt.Errorf("ID вебхука не указан")
	}

	url := fmt.Sprintf("%s/api/v4/webhooks/%d", apiClient.GetBaseURL(), webhook.ID)

	webhookData, err := json.Marshal(webhook)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(webhookData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var updatedWebhook Webhook
	if err := json.NewDecoder(resp.Body).Decode(&updatedWebhook); err != nil {
		return nil, err
	}

	return &updatedWebhook, nil
}

// ListWebhooks получает список вебхуков с возможностью пагинации.
func ListWebhooks(apiClient *client.Client, limit int, page int) ([]*Webhook, error) {
	baseURL := fmt.Sprintf("%s/api/v4/webhooks", apiClient.GetBaseURL())

	// Добавляем параметры запроса
	params := url.Values{}
	params.Add("limit", fmt.Sprintf("%d", limit))
	params.Add("page", fmt.Sprintf("%d", page))

	url := baseURL + "?" + params.Encode()

	req, err := http.NewRequest("GET", url, nil)
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

	var response struct {
		Embedded struct {
			Webhooks []*Webhook `json:"webhooks"`
		} `json:"_embedded"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Embedded.Webhooks, nil
}

// DeleteWebhook удаляет вебхук по его ID.
func DeleteWebhook(apiClient *client.Client, webhookID int) error {
	url := fmt.Sprintf("%s/api/v4/webhooks/%d", apiClient.GetBaseURL(), webhookID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

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

// CreateSimpleWebhook создает новый вебхук с указанными параметрами.
func CreateSimpleWebhook(apiClient *client.Client, destination string, entities []string, actions []string) (*Webhook, error) {
	webhook := &Webhook{
		Destination: destination,
		Settings: &WebhookSettings{
			Entities: entities,
			Actions:  actions,
		},
	}

	return CreateWebhook(apiClient, webhook)
}
