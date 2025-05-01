# Пакет auth

Этот пакет предоставляет функциональность для аутентификации и управления токенами в API amoCRM.

## Основные возможности

- Формирование URL для авторизации пользователя
- Получение токенов доступа по коду авторизации
- Обновление истекших токенов доступа
- Получение долгоживущих (long-lived) токенов

## Основные функции

### GetAuthURL

```go
func GetAuthURL(redirectURI, clientID string) string
```

Формирует URL для авторизации пользователя в amoCRM.

**Параметры:**
- `redirectURI` - URL, на который произойдет перенаправление после авторизации
- `clientID` - ID вашего приложения в amoCRM

**Возвращает:**
- Строку URL для перенаправления пользователя на страницу авторизации amoCRM

### GetAccessToken

```go
func GetAccessToken(baseURL, redirectURI, clientID, clientSecret, code string) (*TokenResponse, error)
```

Получает токен доступа по коду авторизации.

**Параметры:**
- `baseURL` - Базовый URL amoCRM (например, https://example.amocrm.ru)
- `redirectURI` - URL перенаправления
- `clientID` - ID приложения
- `clientSecret` - Секретный ключ приложения
- `code` - Код авторизации, полученный после успешной авторизации

**Возвращает:**
- Структуру `TokenResponse` с токенами доступа и обновления
- Ошибку, если что-то пошло не так

### RefreshAccessToken

```go
func RefreshAccessToken(baseURL, clientID, clientSecret, refreshToken string) (*TokenResponse, error)
```

Обновляет истекший токен доступа с помощью refresh-токена.

**Параметры:**
- `baseURL` - Базовый URL amoCRM
- `clientID` - ID приложения
- `clientSecret` - Секретный ключ приложения
- `refreshToken` - Токен обновления (из предыдущего ответа)

**Возвращает:**
- Обновленные токены доступа и обновления
- Ошибку, если что-то пошло не так

### GetLongLivedToken

```go
func GetLongLivedToken(baseURL, redirectURI, clientID, clientSecret string) (*TokenResponse, error)
```

Получает долгоживущий токен доступа для серверных приложений.

**Параметры:**
- `baseURL` - Базовый URL amoCRM
- `redirectURI` - URL перенаправления
- `clientID` - ID приложения
- `clientSecret` - Секретный ключ приложения

**Возвращает:**
- Долгосрочные токены доступа и обновления
- Ошибку, если что-то пошло не так

## Примеры использования

### Получение токена доступа по коду авторизации

```go
baseURL := "https://example.amocrm.ru"
redirectURI := "https://example.com/oauth2/callback"
clientID := "your-client-id"
clientSecret := "your-client-secret"
code := "auth-code-from-redirect"

tokenResponse, err := auth.GetAccessToken(baseURL, redirectURI, clientID, clientSecret, code)
if err != nil {
    log.Fatalf("Ошибка получения токена: %v", err)
}

// Теперь можно использовать tokenResponse.AccessToken для запросов к API
```

### Обновление токена доступа

```go
baseURL := "https://example.amocrm.ru"
clientID := "your-client-id"
clientSecret := "your-client-secret"
refreshToken := "your-refresh-token"

newTokens, err := auth.RefreshAccessToken(baseURL, clientID, clientSecret, refreshToken)
if err != nil {
    log.Fatalf("Ошибка обновления токена: %v", err)
}

// Сохраните новые токены для дальнейшего использования
```