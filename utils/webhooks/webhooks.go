// Package webhooks предоставляет методы для взаимодействия с вебхуками в API amoCRM.
package webhooks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/chudno/amo_crm_sdk/client"
	"net/http"
	"net/url"
)

// Webhook представляет собой структуру вебхука в amoCRM.
type Webhook struct {
	ID          int       `json:"id,omitempty"`
	Destination string    `json:"destination"`
	Settings    *Settings `json:"settings,omitempty"`
	CreatedAt   int64     `json:"created_at,omitempty"`
	UpdatedAt   int64     `json:"updated_at,omitempty"`
	CreatedBy   int       `json:"created_by,omitempty"`
	AccountID   int       `json:"account_id,omitempty"`
}

// Settings содержит настройки вебхука
type Settings struct {
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

// Get получает вебхук по его ID.
func Get(ctx context.Context, apiClient *client.Client, webhookID int) (*Webhook, error) {
	url := fmt.Sprintf("%s/api/v4/webhooks/%d", apiClient.GetBaseURL(), webhookID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var webhook Webhook
	if err := json.NewDecoder(resp.Body).Decode(&webhook); err != nil {
		return nil, err
	}

	return &webhook, nil
}

// Create создает новый вебхук в amoCRM.
func Create(ctx context.Context, apiClient *client.Client, webhook *Webhook) (*Webhook, error) {
	url := fmt.Sprintf("%s/api/v4/webhooks", apiClient.GetBaseURL())

	webhookData, err := json.Marshal(webhook)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(webhookData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var createdWebhook Webhook
	if err := json.NewDecoder(resp.Body).Decode(&createdWebhook); err != nil {
		return nil, err
	}

	return &createdWebhook, nil
}

// Update обновляет существующий вебхук в amoCRM.
func Update(ctx context.Context, apiClient *client.Client, webhook *Webhook) (*Webhook, error) {
	if webhook.ID == 0 {
		return nil, fmt.Errorf("ID вебхука не указан")
	}

	url := fmt.Sprintf("%s/api/v4/webhooks/%d", apiClient.GetBaseURL(), webhook.ID)

	webhookData, err := json.Marshal(webhook)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewBuffer(webhookData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var updatedWebhook Webhook
	if err := json.NewDecoder(resp.Body).Decode(&updatedWebhook); err != nil {
		return nil, err
	}

	return &updatedWebhook, nil
}

// List получает список вебхуков с возможностью пагинации.
func List(ctx context.Context, apiClient *client.Client, limit int, page int) ([]*Webhook, error) {
	baseURL := fmt.Sprintf("%s/api/v4/webhooks", apiClient.GetBaseURL())

	params := url.Values{}
	params.Add("limit", fmt.Sprintf("%d", limit))
	params.Add("page", fmt.Sprintf("%d", page))

	url := baseURL + "?" + params.Encode()

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

// Delete удаляет вебхук по его ID.
func Delete(ctx context.Context, apiClient *client.Client, webhookID int) error {
	url := fmt.Sprintf("%s/api/v4/webhooks/%d", apiClient.GetBaseURL(), webhookID)

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

// CreateSimple создает новый вебхук с указанными параметрами.
func CreateSimple(ctx context.Context, apiClient *client.Client, destination string, entities []string, actions []string) (*Webhook, error) {
	webhook := &Webhook{
		Destination: destination,
		Settings: &Settings{
			Entities: entities,
			Actions:  actions,
		},
	}

	return Create(ctx, apiClient, webhook)
}
