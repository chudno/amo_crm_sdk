// Пакет users предоставляет методы для взаимодействия с сущностями "Пользователи" в API amoCRM.
package users

import (
	"encoding/json"
	"fmt"
	"github.com/chudno/amo_crm_sdk/client"
	"net/http"
)

// User представляет собой структуру пользователя в amoCRM.
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Lang     string `json:"lang"`
	Rights   Rights `json:"rights"`
	IsActive bool   `json:"is_active"`
}

// Rights представляет права пользователя в системе.
type Rights struct {
	Leads     bool `json:"leads"`
	Contacts  bool `json:"contacts"`
	Companies bool `json:"companies"`
	Tasks     bool `json:"tasks"`
	Mailbox   bool `json:"mailbox"`
	Catalog   bool `json:"catalog"`
	IsAdmin   bool `json:"is_admin"`
	IsManager bool `json:"is_manager"`
}

// GetUser получает пользователя по его ID.
func GetUser(apiClient *client.Client, userID int) (*User, error) {
	url := fmt.Sprintf("%s/api/v4/users/%d", apiClient.GetBaseURL(), userID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// GetCurrentUser получает информацию о текущем пользователе (владельце API-ключа).
func GetCurrentUser(apiClient *client.Client) (*User, error) {
	url := fmt.Sprintf("%s/api/v4/users/self", apiClient.GetBaseURL())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// ListUsers получает список пользователей с возможностью фильтрации и пагинации.
func ListUsers(apiClient *client.Client, limit int, page int) ([]User, error) {
	url := fmt.Sprintf("%s/api/v4/users?limit=%d&page=%d", apiClient.GetBaseURL(), limit, page)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var users struct {
		Embedded struct {
			Items []User `json:"items"`
		} `json:"_embedded"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}

	return users.Embedded.Items, nil
}
