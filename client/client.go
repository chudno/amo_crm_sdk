// Пакет client предоставляет функциональность для взаимодействия с API amoCRM.
package client

import (
	"net/http"
	"time"
)

// Client - структура для создания нового amoCRM API клиента.
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient создает новый экземпляр клиента для amoCRM API.
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// doRequest выполняет HTTP-запрос к API amoCRM.
func (c *Client) DoRequest(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	return c.httpClient.Do(req)
}

// GetBaseURL возвращает базовый URL API.
func (c *Client) GetBaseURL() string {
	return c.baseURL
}
