package short_links

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// MockRequest представляет собой структуру для хранения информации о запросе.
type MockRequest struct {
	Method string
	URL    string
	Body   []byte
}

// MockResponse представляет собой структуру для хранения информации об ответе.
type MockResponse struct {
	StatusCode int
	Body       string
	Headers    map[string]string
}

// AdvancedMockClient представляет собой продвинутый мок-клиент для тестирования.
type AdvancedMockClient struct {
	BaseURL        string
	ExpectedMethod string
	ExpectedURL    string
	ExpectedBody   interface{}
	MockResponse   *MockResponse
	LastRequest    *MockRequest
}

// DoRequest выполняет запрос и возвращает ответ.
func (c *AdvancedMockClient) DoRequest(req *http.Request) (*http.Response, error) {
	// Создаем информацию о запросе
	mockReq := &MockRequest{
		Method: req.Method,
		URL:    req.URL.String(),
	}

	// Читаем тело запроса, если оно есть
	if req.Body != nil {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		mockReq.Body = body
		// Восстанавливаем тело для повторного чтения
		req.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	// Сохраняем информацию о запросе
	c.LastRequest = mockReq

	// Проверяем ожидаемый метод
	if c.ExpectedMethod != "" && c.ExpectedMethod != req.Method {
		return nil, fmt.Errorf("неожиданный метод: %s, ожидался: %s", req.Method, c.ExpectedMethod)
	}

	// Проверяем ожидаемый URL
	if c.ExpectedURL != "" {
		// Преобразуем URL в более простую форму для сравнения
		expectedURL := strings.Split(c.ExpectedURL, "?")[0]
		actualURL := strings.Split(req.URL.String(), "?")[0]

		if !strings.HasSuffix(actualURL, expectedURL) {
			return nil, fmt.Errorf("неожиданный URL: %s, ожидался URL, содержащий: %s", actualURL, expectedURL)
		}
	}

	// Проверяем ожидаемое тело запроса
	if c.ExpectedBody != nil && req.Body != nil {
		// Сериализуем ожидаемое тело для сравнения
		expectedJSON, err := json.Marshal(c.ExpectedBody)
		if err != nil {
			return nil, err
		}

		// Если тела не совпадают, возвращаем ошибку
		if !bytes.Equal(mockReq.Body, expectedJSON) {
			return nil, fmt.Errorf("неожиданное тело запроса: %s, ожидалось: %s", string(mockReq.Body), string(expectedJSON))
		}
	}

	// Если нет ответа, возвращаем ошибку
	if c.MockResponse == nil {
		return nil, fmt.Errorf("не задан мок-ответ")
	}

	// Создаем заголовки ответа
	headers := http.Header{}
	if c.MockResponse.Headers != nil {
		for k, v := range c.MockResponse.Headers {
			headers.Set(k, v)
		}
	}
	// По умолчанию устанавливаем Content-Type: application/json
	if _, ok := c.MockResponse.Headers["Content-Type"]; !ok {
		headers.Set("Content-Type", "application/json")
	}

	// Создаем ответ
	resp := &http.Response{
		StatusCode: c.MockResponse.StatusCode,
		Body:       io.NopCloser(strings.NewReader(c.MockResponse.Body)),
		Header:     headers,
	}

	return resp, nil
}

// GetBaseURL возвращает базовый URL клиента.
func (c *AdvancedMockClient) GetBaseURL() string {
	return c.BaseURL
}
