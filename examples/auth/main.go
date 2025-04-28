package main

import (
	"fmt"
	"log"

	"github.com/chudno/amo_crm_sdk/auth"
)

func main() {
	// Заменить на свои значения
	baseURL := "https://example.amocrm.ru"
	clientID := "ваш_client_id"
	clientSecret := "ваш_client_secret"
	redirectURI := "https://your-redirect-uri.com"
	code := "код_авторизации" // Получен после перехода пользователя по ссылке авторизации

	// Получение URL для авторизации пользователя
	authURL := auth.GetAuthURL(baseURL, clientID, redirectURI, "random_state", "popup")
	fmt.Println("Перейдите по ссылке для авторизации:", authURL)

	// Получение токена доступа с помощью кода авторизации
	authResponse, err := auth.GetAccessToken(baseURL, clientID, clientSecret, code, redirectURI)
	if err != nil {
		log.Fatalf("Ошибка при получении токена: %v", err)
	}

	fmt.Println("Токен доступа получен:", authResponse.AccessToken)
	fmt.Println("Токен обновления:", authResponse.RefreshToken)
	fmt.Println("Срок действия токена (в секундах):", authResponse.ExpiresIn)

	// Обновление токена доступа по refresh токену
	refreshedAuth, err := auth.RefreshAccessToken(baseURL, clientID, clientSecret, authResponse.RefreshToken)
	if err != nil {
		log.Fatalf("Ошибка при обновлении токена: %v", err)
	}

	fmt.Println("Новый токен доступа получен:", refreshedAuth.AccessToken)
}
