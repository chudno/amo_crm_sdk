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
func GetAuthURL(baseURL, clientID, redirectURI, state, mode string) string
```

Формирует URL для авторизации пользователя в amoCRM.

**Параметры:**
- `baseURL` - Базовый URL вашего аккаунта amoCRM (например, https://example.amocrm.ru)
- `clientID` - ID вашего приложения в amoCRM
- `redirectURI` - URL, на который произойдет перенаправление после авторизации
- `state` - Произвольная строка для проверки подлинности перенаправления
- `mode` - Режим отображения ("popup" или "post_message")

**Возвращает:**
- Строку URL для перенаправления пользователя на страницу авторизации amoCRM

### GetAccessToken

```go
func GetAccessToken(ctx context.Context, baseURL, clientID, clientSecret, code, redirectURI string) (*AuthResponse, error)
```

Получает токен доступа по коду авторизации.

**Параметры:**
- `ctx` - Контекст для управления временем выполнения и отменой запроса
- `baseURL` - Базовый URL amoCRM (например, https://example.amocrm.ru)
- `clientID` - ID приложения
- `clientSecret` - Секретный ключ приложения
- `code` - Код авторизации, полученный после успешной авторизации
- `redirectURI` - URL перенаправления, указанный при регистрации интеграции

**Возвращает:**
- Структуру `AuthResponse` с токенами доступа и обновления
- Ошибку, если что-то пошло не так

### RefreshAccessToken

```go
func RefreshAccessToken(ctx context.Context, baseURL, clientID, clientSecret, refreshToken string) (*AuthResponse, error)
```

Обновляет истекший токен доступа с помощью refresh-токена.

**Параметры:**
- `ctx` - Контекст для управления временем выполнения и отменой запроса
- `baseURL` - Базовый URL amoCRM
- `clientID` - ID приложения
- `clientSecret` - Секретный ключ приложения
- `refreshToken` - Токен обновления (из предыдущего ответа)

**Возвращает:**
- Структуру `AuthResponse` с обновленными токенами доступа и обновления
- Ошибку, если что-то пошло не так

### GetLongLivedToken

```go
func GetLongLivedToken(ctx context.Context, baseURL, clientID, clientSecret string) (*AuthResponse, error)
```

Получает долгоживущий токен доступа для серверных интеграций.

**Параметры:**
- `ctx` - Контекст для управления временем выполнения и отменой запроса
- `baseURL` - Базовый URL amoCRM
- `clientID` - ID интеграции
- `clientSecret` - Секретный ключ интеграции

**Возвращает:**
- Структуру `AuthResponse` с долгоживущим токеном доступа
- Ошибку, если что-то пошло не так

## Примеры использования

### Получение токена доступа по коду авторизации

```go
import "context"

ctx := context.Background()
baseURL := "https://example.amocrm.ru"
clientID := "your-client-id"
clientSecret := "your-client-secret"
code := "auth-code-from-redirect"
redirectURI := "https://example.com/oauth2/callback"

authResponse, err := auth.GetAccessToken(ctx, baseURL, clientID, clientSecret, code, redirectURI)
if err != nil {
    log.Fatalf("Ошибка получения токена: %v", err)
}

// authResponse.AccessToken содержит токен для запросов к API
```

### Обновление токена доступа

```go
import "context"

ctx := context.Background()
baseURL := "https://example.amocrm.ru"
clientID := "your-client-id"
clientSecret := "your-client-secret"
refreshToken := "your-refresh-token"

authResponse, err := auth.RefreshAccessToken(ctx, baseURL, clientID, clientSecret, refreshToken)
if err != nil {
    log.Fatalf("Ошибка обновления токена: %v", err)
}

// authResponse.AccessToken и authResponse.RefreshToken содержат новые токены
```

### Получение долгоживущего токена

```go
import "context"

ctx := context.Background()
baseURL := "https://example.amocrm.ru"
clientID := "your-client-id"
clientSecret := "your-client-secret"

authResponse, err := auth.GetLongLivedToken(ctx, baseURL, clientID, clientSecret)
if err != nil {
    log.Fatalf("Ошибка получения долгоживущего токена: %v", err)
}

// authResponse.AccessToken содержит долгоживущий токен
```