package mailing

import (
	"io"
	"net/http"
	"strings"
)

// MockRequest представляет информацию о запросе в мок-клиенте
type MockRequest struct {
	Method string
	URL    string
	Body   string
	Header http.Header
}

// MockResponse представляет информацию об ответе мок-клиента
type MockResponse struct {
	StatusCode int
	Body       string
	Error      error
}

// AdvancedMockClient представляет собой расширенный мок-клиент для тестирования
type AdvancedMockClient struct {
	// BaseURL базовый URL API
	BaseURL string
	// LastRequest последний выполненный запрос
	LastRequest *MockRequest
	// MockResponses предопределенные ответы для запросов
	MockResponses map[string]MockResponse
	// DefaultResponse ответ по умолчанию, если нет совпадений в MockResponses
	DefaultResponse MockResponse
}

// NewAdvancedMockClient создает новый мок-клиент с заданными параметрами
func NewAdvancedMockClient(baseURL string, defaultResponse MockResponse) *AdvancedMockClient {
	return &AdvancedMockClient{
		BaseURL:         baseURL,
		MockResponses:   make(map[string]MockResponse),
		DefaultResponse: defaultResponse,
	}
}

// DoRequest реализует метод DoRequest из интерфейса Requester
func (c *AdvancedMockClient) DoRequest(req *http.Request) (*http.Response, error) {
	// Получаем тело запроса, если оно есть
	var bodyStr string
	if req.Body != nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err == nil {
			bodyStr = string(bodyBytes)
			// Восстанавливаем тело запроса для повторного чтения
			req.Body = io.NopCloser(strings.NewReader(bodyStr))
		}
	}

	// Сохраняем информацию о запросе
	c.LastRequest = &MockRequest{
		Method: req.Method,
		URL:    req.URL.String(),
		Body:   bodyStr,
		Header: req.Header,
	}

	// Формируем ключ запроса для поиска в предопределенных ответах
	requestKey := req.Method + " " + req.URL.Path

	// Ищем ответ в предопределенных ответах
	response, ok := c.MockResponses[requestKey]
	if !ok {
		// Если ничего не нашли, используем ответ по умолчанию
		response = c.DefaultResponse
	}

	// Если в ответе есть ошибка, возвращаем её
	if response.Error != nil {
		return nil, response.Error
	}

	// Создаем и возвращаем HTTP-ответ
	httpResponse := &http.Response{
		StatusCode: response.StatusCode,
		Body:       io.NopCloser(strings.NewReader(response.Body)),
		Header:     make(http.Header),
	}
	return httpResponse, nil
}

// GetBaseURL реализует метод GetBaseURL из интерфейса Requester
func (c *AdvancedMockClient) GetBaseURL() string {
	return c.BaseURL
}

// AddMockResponse добавляет новый предопределенный ответ для заданного пути
func (c *AdvancedMockClient) AddMockResponse(method, path string, response MockResponse) {
	c.MockResponses[method+" "+path] = response
}
