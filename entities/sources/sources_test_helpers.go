package sources

import (
	"io"
	"net/http"
	"strings"
	"testing"
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

// createGetSourcesSuccessMockClient создает мок-клиент с успешным ответом для списка источников
func createGetSourcesSuccessMockClient() *AdvancedMockClient {
	return NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
		StatusCode: http.StatusOK,
		Body: `{
			"_embedded": {
				"sources": [
					{
						"id": 1001,
						"name": "Входящие звонки",
						"type": "calls",
						"default": true,
						"created_at": 1609459200,
						"updated_at": 1609459300
					},
					{
						"id": 1002,
						"name": "Реклама в интернете",
						"type": "advertising",
						"default": false,
						"created_at": 1609459400,
						"updated_at": 1609459500
					}
				]
			}
		}`,
	})
}

// createGetSourcesEmptyMockClient создает мок-клиент с пустым списком источников
func createGetSourcesEmptyMockClient() *AdvancedMockClient {
	return NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
		StatusCode: http.StatusOK,
		Body: `{
			"_embedded": {
				"sources": []
			}
		}`,
	})
}

// createGetSourcesErrorMockClient создает мок-клиент с ошибкой сервера
func createGetSourcesErrorMockClient() *AdvancedMockClient {
	return NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       `{"error": "Internal Server Error"}`,
	})
}

// createGetSourcesWithFilterMockClient создает мок-клиент для запроса с фильтром
func createGetSourcesWithFilterMockClient() *AdvancedMockClient {
	return NewAdvancedMockClient("https://example.amocrm.ru", MockResponse{
		StatusCode: http.StatusOK,
		Body: `{
			"_embedded": {
				"sources": [
					{
						"id": 1001,
						"name": "Входящие звонки",
						"type": "calls",
						"default": true,
						"created_at": 1609459200,
						"updated_at": 1609459300
					}
				]
			}
		}`,
	})
}

// verifySourcesList проверяет правильность полученного списка источников
func verifySourcesList(t *testing.T, sources []Source, expectedCount int) {
	if len(sources) != expectedCount {
		t.Errorf("Ожидалось %d источников, получено %d", expectedCount, len(sources))
	}

	if expectedCount > 0 {
		if sources[0].ID != 1001 || sources[0].Name != "Входящие звонки" || sources[0].Type != "calls" {
			t.Errorf("Неверные данные в первом источнике: %+v", sources[0])
		}
	}

	if expectedCount > 1 {
		if sources[1].ID != 1002 || sources[1].Name != "Реклама в интернете" || sources[1].Type != "advertising" {
			t.Errorf("Неверные данные во втором источнике: %+v", sources[1])
		}
	}
}

// verifyFilterInRequest проверяет наличие параметра фильтра в URL запроса
func verifyFilterInRequest(t *testing.T, mockClient *AdvancedMockClient, expectedFilterPart string) {
	if mockClient.LastRequest == nil {
		t.Fatal("Запрос не был выполнен")
	}

	if mockClient.LastRequest.URL == "" {
		t.Fatal("URL запроса пустой")
	}

	if !strings.Contains(mockClient.LastRequest.URL, expectedFilterPart) {
		t.Errorf("Фильтр не был добавлен к URL запроса: %s", mockClient.LastRequest.URL)
	}
}
