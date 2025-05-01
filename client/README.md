# Пакет client

Этот пакет предоставляет базовый HTTP-клиент для работы с API amoCRM, который используется всеми другими модулями SDK.

## Основные компоненты

### Структура Client

```go
type Client struct {
    baseURL    string
    apiKey     string
    httpClient *http.Client
}
```

Структура `Client` представляет клиент API amoCRM и содержит:
- `baseURL` - базовый URL для API amoCRM (например, https://example.amocrm.ru)
- `apiKey` - токен доступа для авторизации запросов
- `httpClient` - стандартный HTTP-клиент Go с заданным таймаутом

### Функция NewClient

```go
func NewClient(baseURL, apiKey string) *Client
```

Создает новый экземпляр клиента API amoCRM с указанным базовым URL и токеном доступа.

**Параметры:**
- `baseURL` - базовый URL для API amoCRM
- `apiKey` - токен доступа, полученный из пакета auth

**Возвращает:**
- Экземпляр `*Client`, настроенный для работы с API

### Метод DoRequest

```go
func (c *Client) DoRequest(req *http.Request) (*http.Response, error)
```

Выполняет HTTP-запрос к API amoCRM, автоматически добавляя заголовок авторизации.

**Параметры:**
- `req` - HTTP-запрос для выполнения

**Возвращает:**
- HTTP-ответ и ошибку (если есть)

### Метод GetBaseURL

```go
func (c *Client) GetBaseURL() string
```

Возвращает базовый URL API, указанный при создании клиента.

**Возвращает:**
- Строку с базовым URL API

## Примеры использования

### Создание клиента и выполнение запроса

```go
// Создание клиента с базовым URL и токеном доступа
baseURL := "https://example.amocrm.ru"
accessToken := "your-access-token"
apiClient := client.NewClient(baseURL, accessToken)

// Создание HTTP-запроса
req, err := http.NewRequest("GET", apiClient.GetBaseURL()+"/api/v4/leads", nil)
if err != nil {
    log.Fatalf("Ошибка создания запроса: %v", err)
}

// Добавление параметров запроса
q := req.URL.Query()
q.Add("limit", "50")
req.URL.RawQuery = q.Encode()

// Выполнение запроса через клиент
resp, err := apiClient.DoRequest(req)
if err != nil {
    log.Fatalf("Ошибка выполнения запроса: %v", err)
}
defer resp.Body.Close()

// Обработка ответа
if resp.StatusCode != http.StatusOK {
    log.Fatalf("Неверный статус-код: %d", resp.StatusCode)
}

// Чтение и декодирование ответа
var result struct {
    Embedded struct {
        Items []map[string]interface{} `json:"items"`
    } `json:"_embedded"`
}
if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
    log.Fatalf("Ошибка декодирования ответа: %v", err)
}

// Работа с результатом
for _, item := range result.Embedded.Items {
    fmt.Printf("ID: %v, Name: %v\n", item["id"], item["name"])
}
```

## Особенности и ограничения

- Клиент автоматически устанавливает таймаут в 30 секунд для всех запросов
- Авторизация происходит через Bearer-токен в заголовке Authorization
- Клиент не обрабатывает ошибки API, это делают методы в конкретных пакетах сущностей
