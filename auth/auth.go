// Пакет auth предоставляет методы для аутентификации в API amoCRM.
//
// Этот пакет содержит все необходимые инструменты для OAuth2-авторизации,
// получения токенов доступа и их обновления. Процесс авторизации соответствует
// официальной документации amoCRM по OAuth-интеграции.
package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// TokenType определяет тип токена авторизации
type TokenType string

const (
	// TokenTypeOAuth указывает на OAuth2-токен, который необходимо обновлять через refresh_token
	TokenTypeOAuth TokenType = "oauth2"
	// TokenTypeLongLived указывает на долгоживущий токен, который не требует обновления
	TokenTypeLongLived TokenType = "long_lived"
)

// AuthResponse представляет ответ от сервера при аутентификации.
// Содержит токен доступа, токен обновления и время действия токена.
// При использовании долгоживущего токена RefreshToken будет пустым, а ExpiresIn - очень большим числом.
type AuthResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// AuthRequest представляет запрос на аутентификацию.
// Используется для формирования JSON-запроса к OAuth-серверу amoCRM.
// Различные типы grant_type позволяют получать разные типы токенов.
type AuthRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	Code         string `json:"code,omitempty"`
	RedirectURI  string `json:"redirect_uri,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// GetAccessToken получает токен доступа по коду авторизации.
//
// Параметры:
//   - baseURL: базовый URL вашего аккаунта amoCRM (например, https://example.amocrm.ru)
//   - clientID: ID клиента, полученный при регистрации интеграции
//   - clientSecret: секретный ключ, полученный при регистрации интеграции
//   - code: код авторизации, полученный после перенаправления пользователя
//   - redirectURI: URI перенаправления, указанный при регистрации интеграции
//
// Возвращает структуру AuthResponse с токенами доступа или ошибку.
func GetAccessToken(baseURL, clientID, clientSecret, code, redirectURI string) (*AuthResponse, error) {
	url := baseURL + "/oauth2/access_token"

	authReq := AuthRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		GrantType:    "authorization_code",
		Code:         code,
		RedirectURI:  redirectURI,
	}

	authJSON, err := json.Marshal(authReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(authJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, err
	}

	return &authResp, nil
}

// GetAuthURL формирует URL для авторизации в amoCRM.
//
// Параметры:
//   - baseURL: базовый URL вашего аккаунта amoCRM
//   - clientID: ID клиента, полученный при регистрации интеграции
//   - redirectURI: URI перенаправления, указанный при регистрации интеграции
//   - state: произвольная строка для проверки подлинности перенаправления
//   - mode: режим отображения ("popup" или "post_message")
//
// Возвращает полный URL для перенаправления пользователя на страницу авторизации.
func GetAuthURL(baseURL, clientID, redirectURI, state, mode string) string {
	params := url.Values{}
	params.Add("client_id", clientID)
	params.Add("mode", mode)
	params.Add("redirect_uri", redirectURI)
	params.Add("response_type", "code")
	params.Add("state", state)
	return fmt.Sprintf("%s/oauth2/access_token?%s", baseURL, params.Encode())
}

// RefreshAccessToken обновляет токен доступа по refresh токену.
//
// Параметры:
//   - baseURL: базовый URL вашего аккаунта amoCRM
//   - clientID: ID клиента, полученный при регистрации интеграции
//   - clientSecret: секретный ключ, полученный при регистрации интеграции
//   - refreshToken: токен обновления, полученный ранее от API
//
// Возвращает новую структуру AuthResponse с обновленными токенами или ошибку.
// Этот метод следует использовать, когда срок действия текущего токена истекает.
func RefreshAccessToken(baseURL, clientID, clientSecret, refreshToken string) (*AuthResponse, error) {
	url := baseURL + "/oauth2/access_token"

	authReq := AuthRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		GrantType:    "refresh_token",
		RefreshToken: refreshToken,
	}

	authJSON, err := json.Marshal(authReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(authJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, err
	}

	return &authResp, nil
}

// GetLongLivedToken получает долгоживущий токен доступа для серверных интеграций.
//
// Параметры:
//   - baseURL: базовый URL вашего аккаунта amoCRM (например, https://example.amocrm.ru)
//   - clientID: ID интеграции
//   - clientSecret: секретный ключ интеграции
//
// Возвращает структуру AuthResponse с долгоживущим токеном доступа или ошибку.
// Долгоживущие токены не требуют обновления и могут использоваться длительное время.
func GetLongLivedToken(baseURL, clientID, clientSecret string) (*AuthResponse, error) {
	url := baseURL + "/oauth2/access_token"

	authReq := AuthRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		GrantType:    "client_credentials",
	}

	authJSON, err := json.Marshal(authReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(authJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, err
	}

	return &authResp, nil
}
