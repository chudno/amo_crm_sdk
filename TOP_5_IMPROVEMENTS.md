# Топ 5 улучшений для amoCRM SDK

Результат анализа кодовой базы SDK (версия 0.2.0).

---

## 1. Отсутствие обработки rate-limiting (HTTP 429)

**Критичность:** Высокая
**Затронутые файлы:** `client/client.go`, все файлы в `entities/`, `auth/auth.go`

### Проблема

SDK не обрабатывает HTTP 429 (Too Many Requests) от API amoCRM. При превышении лимита запросов все вызовы начинают завершаться ошибками без попытки повтора.

### Примеры

В `client/client.go:29-33` — метод `DoRequest` просто проксирует запрос без какой-либо логики повторов:

```go
func (c *Client) DoRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
    req = req.WithContext(ctx)
    req.Header.Set("Authorization", "Bearer "+c.apiKey)
    return c.httpClient.Do(req)  // Нет обработки 429
}
```

В entity-файлах ошибки выбрасываются сразу:

```go
// entities/leads/leads.go:246-248
if resp.StatusCode != http.StatusOK {
    return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
}
```

### Рекомендация

Реализовать в `client.DoRequest` автоматический retry с exponential backoff при получении HTTP 429. Читать заголовок `Retry-After` из ответа API для определения интервала ожидания. Предусмотреть настраиваемое максимальное количество повторов.

---

## 2. Отсутствие проверки HTTP-статуса в `Create`/`Update` у ключевых сущностей

**Критичность:** Высокая
**Затронутые файлы:** `entities/leads/leads.go`, `entities/companies/companies.go`

### Проблема

Функции `Create` и `Update` в модулях `leads` и `companies` не проверяют HTTP-статус ответа перед декодированием JSON. Если API вернёт ошибку (4xx/5xx), SDK попытается распарсить тело ответа как успешный результат, что приведёт к пустым или некорректным данным без информативной ошибки.

### Примеры

`entities/leads/leads.go:114-128` — `Create` не проверяет статус:

```go
resp, err := apiClient.DoRequest(ctx, req)
if err != nil {
    return nil, err
}
defer func() { _ = resp.Body.Close() }()

// Нет проверки resp.StatusCode!

var response struct {
    Embedded struct {
        Leads []*Lead `json:"leads"`
    } `json:"_embedded"`
}
if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
    return nil, err  // Ошибка парсинга вместо информативной ошибки API
}
```

`entities/leads/leads.go:157-168` — `Update` аналогично без проверки.

`entities/companies/companies.go:150-161` — `Update` без проверки.

При этом `contacts.Create` (`entities/contacts/contacts.go:117`) проверяет статус корректно:

```go
if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
    return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
}
```

### Рекомендация

Добавить проверку HTTP-статуса перед декодированием ответа во все функции `Create` и `Update`. Привести к единообразию с `contacts.Create`.

---

## 3. Ошибки API не содержат тело ответа — невозможно диагностировать проблему

**Критичность:** Средняя
**Затронутые файлы:** Большинство файлов в `entities/`, `auth/auth.go`

### Проблема

При неуспешном HTTP-статусе SDK возвращает только код статуса, отбрасывая тело ответа. API amoCRM возвращает в теле подробное описание ошибки (validation errors, описание проблемы), которое полностью теряется.

### Примеры

Типичный паттерн (19 из 20 entity-пакетов):

```go
// entities/leads/leads.go:246-248
if resp.StatusCode != http.StatusOK {
    return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
}
```

Лучший вариант, который уже реализован в некоторых файлах:

```go
// entities/catalogs/catalogs.go:136-139
if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
    body, _ := io.ReadAll(resp.Body)
    return nil, fmt.Errorf("неожиданный статус-код: %d, ответ: %s", resp.StatusCode, body)
}
```

### Рекомендация

Создать общую функцию проверки ответа в пакете `client`, которая будет читать тело при ошибочном статусе и формировать структурированную ошибку с кодом статуса, URL запроса и телом ответа. Использовать её во всех entity-пакетах.

---

## 4. Нет автоматической пагинации — пользователь должен вручную перебирать страницы

**Критичность:** Средняя
**Затронутые файлы:** Все `List`-функции в `entities/`

### Проблема

Все `List`-функции принимают `page` и `limit` и возвращают одну страницу результатов. Для получения всех записей пользователь должен вручную реализовать цикл пагинации, отслеживая текущую страницу и определяя момент окончания данных. Это повторяющийся boilerplate-код у каждого пользователя SDK.

### Примеры

`entities/leads/leads.go:213`:

```go
func List(ctx context.Context, apiClient *client.Client, page, limit int,
    filter map[string]string, withOptions ...WithOption) ([]Lead, error) {
    // Возвращает только одну страницу
}
```

`ListResponse` содержит `Total`, но пользователь должен сам вычислять число страниц:

```go
type ListResponse struct {
    Page     int `json:"page"`
    PerPage  int `json:"per_page"`
    Total    int `json:"total"`
    Embedded struct { Items []Lead `json:"leads"` } `json:"_embedded"`
}
```

### Рекомендация

Добавить функцию `ListAll` (или итератор), которая автоматически перебирает все страницы и возвращает полный список. Оставить существующую `List` для ручного контроля. Вариант с callback/channel позволит обрабатывать результаты по мере получения без загрузки всех данных в память.

---

## 5. Значительное дублирование кода между entity-пакетами

**Критичность:** Средняя
**Затронутые файлы:** Все 20 пакетов в `entities/`

### Проблема

Функции `Get`, `Create`, `Update`, `List`, `Delete` во всех entity-пакетах имеют практически идентичную структуру: marshal → NewRequest → DoRequest → проверка статуса → decode. Каждый entity-пакет реализует этот паттерн заново (~30-40 строк на функцию × 4-5 функций × 20 пакетов = ~3000+ строк дублированного кода).

Баги и несоответствия (см. пункт 2 и 3) являются прямым следствием этого дублирования — исправление вносится в одном месте, но забывается в других.

### Примеры

Сравнение `Create` в трёх пакетах — одинаковый алгоритм:

| Файл | Строки | Проверка статуса |
|------|--------|-----------------|
| `leads/leads.go:99-135` | 37 строк | Нет |
| `contacts/contacts.go:97-135` | 39 строк | Есть (200, 201) |
| `companies/companies.go:95-133` | 39 строк | Есть (200, 201) |

### Рекомендация

Вынести общую логику HTTP-запросов в пакет `client` — generic-функции (с Go 1.18+ generics) или helper-функции для типовых CRUD-операций. Это сократит объём кода, устранит несоответствия и упростит добавление новых сущностей.

---

## Сводная таблица

| # | Улучшение | Критичность | Влияние |
|---|-----------|-------------|---------|
| 1 | Rate-limiting и retry | Высокая | Стабильность в продакшне |
| 2 | Проверка HTTP-статуса в Create/Update | Высокая | Корректная обработка ошибок |
| 3 | Тело ответа в ошибках | Средняя | Диагностика проблем |
| 4 | Автоматическая пагинация | Средняя | Удобство использования |
| 5 | Устранение дублирования кода | Средняя | Поддерживаемость |
