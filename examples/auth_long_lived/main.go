package main

import (
	"fmt"
	"log"

	"github.com/chudno/amo_crm_sdk/auth"
	"github.com/chudno/amo_crm_sdk/client"
	"github.com/chudno/amo_crm_sdk/entities/users"
)

func main() {
	// Параметры для получения долгоживущего токена
	baseURL := "https://your-domain.amocrm.ru"
	clientID := "ваш_client_id"         // ID интеграции
	clientSecret := "ваш_client_secret" // Секретный ключ интеграции

	// Получение долгоживущего токена
	authResponse, err := auth.GetLongLivedToken(baseURL, clientID, clientSecret)
	if err != nil {
		log.Fatalf("Ошибка при получении долгоживущего токена: %v", err)
	}

	fmt.Println("Долгоживущий токен успешно получен:")
	fmt.Printf("Access Token: %s\n", authResponse.AccessToken)
	fmt.Printf("Срок действия: %d секунд (примерно %d дней)\n",
		authResponse.ExpiresIn, authResponse.ExpiresIn/86400)

	// Создание клиента API с полученным токеном
	apiClient := client.NewClient(baseURL, authResponse.AccessToken)

	// Пример использования API с долгоживущим токеном
	currentUser, err := users.GetCurrentUser(apiClient)
	if err != nil {
		log.Fatalf("Ошибка при получении информации о текущем пользователе: %v", err)
	}

	fmt.Println("\nИнформация о текущем пользователе:")
	fmt.Printf("ID: %d\n", currentUser.ID)
	fmt.Printf("Имя: %s\n", currentUser.Name)
	fmt.Printf("Email: %s\n", currentUser.Email)

	/*
		Важно: Долгоживущие токены не требуют обновления через refresh token,
		поэтому их можно безопасно хранить и использовать длительное время
		для серверных приложений. Срок действия обычно составляет около 1 года.
	*/
}
