package segments

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// Requester интерфейс для выполнения HTTP-запросов
type Requester interface {
	DoRequest(req *http.Request) (*http.Response, error)
}

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
func (c *AdvancedMockClient) DoRequest(req *http.Request) (*http.Response, error) {
	// Ищем подходящий ответ для метода и пути
	resp, found := c.Responses[MockRequest{Method: req.Method, Path: req.URL.Path}]
	
	// Если не найден, возвращаем ответ по умолчанию
	if !found {
		resp = c.DefaultResponse
	}
	
	// Формируем HTTP-ответ
	response := &http.Response{
		StatusCode: resp.StatusCode,
		Body:       io.NopCloser(strings.NewReader(resp.Body)),
		Header:     make(http.Header),
		Request:    req,
	}
	
	// Добавляем заголовки
	for k, v := range resp.Headers {
		response.Header.Set(k, v)
	}
	
	return response, nil
}

// GetSegmentsWithRequester получает список сегментов с использованием интерфейса Requester
// Это вспомогательная функция для тестирования
func GetSegmentsWithRequester(requester Requester, page, limit int, options ...WithOption) ([]Segment, error) {
	// Путь к API сегментов
	path := "/api/v4/segments"

	// Создаем параметры запроса
	params := make(map[string]string)

	// Добавляем пагинацию
	params["page"] = strconv.Itoa(page)
	params["limit"] = strconv.Itoa(limit)

	// Применяем дополнительные опции
	for _, option := range options {
		option(params)
	}

	// Создаем запрос
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	// Добавляем параметры запроса
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	// Отправляем запрос
	resp, err := requester.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	// Разбираем ответ
	var segmentsResponse struct {
		Embedded struct {
			Segments []Segment `json:"segments"`
		} `json:"_embedded"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&segmentsResponse); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return segmentsResponse.Embedded.Segments, nil
}
