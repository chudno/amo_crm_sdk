package access_rights

import (
	"context"
	"io"
	"net/http"
	"strings"
)

// MockResponse описывает мок-ответ для тестирования
type MockResponse struct {
	StatusCode int
	Body       string
	Headers    map[string]string
}

// MockRequest описывает мок-запрос для сопоставления
type MockRequest struct {
	Method string
	Path   string
}

// AdvancedMockClient улучшенный мок-клиент для перехвата всех видов запросов
type AdvancedMockClient struct {
	// Отображение ключа MockRequest на ответ MockResponse
	Responses map[MockRequest]MockResponse
	// Ответ по умолчанию, если запрос не найден
	DefaultResponse MockResponse
}

// NewAdvancedMockClient создает новый мок-клиент с настройками по умолчанию
func NewAdvancedMockClient() *AdvancedMockClient {
	return &AdvancedMockClient{
		Responses: make(map[MockRequest]MockResponse),
		DefaultResponse: MockResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error": "Unexpected request"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		},
	}
}

// AddResponse добавляет ответ для конкретного запроса
func (c *AdvancedMockClient) AddResponse(method, path string, statusCode int, body string, headers map[string]string) {
	if headers == nil {
		headers = map[string]string{"Content-Type": "application/json"}
	}
	c.Responses[MockRequest{Method: method, Path: path}] = MockResponse{
		StatusCode: statusCode,
		Body:       body,
		Headers:    headers,
	}
}

// DoRequest реализует интерфейс Requester
func (c *AdvancedMockClient) DoRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	resp, found := c.Responses[MockRequest{Method: req.Method, Path: req.URL.Path}]

	if !found {
		resp = c.DefaultResponse
	}

	response := &http.Response{
		StatusCode: resp.StatusCode,
		Body:       io.NopCloser(strings.NewReader(resp.Body)),
		Header:     make(http.Header),
		Request:    req,
	}

	for k, v := range resp.Headers {
		response.Header.Set(k, v)
	}

	return response, nil
}
