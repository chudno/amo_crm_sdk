// Package users предоставляет методы для взаимодействия с сущностями "Пользователи" в API amoCRM.
package users

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chudno/amo_crm_sdk/client"
)

// Account представляет информацию об аккаунте amoCRM.
type Account struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Subdomain     string `json:"subdomain"`
	CurrentUserID int    `json:"current_user_id"`
}

// User представляет собой структуру пользователя в amoCRM.
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Lang     string `json:"lang"`
	Rights   Rights `json:"rights"`
	IsActive bool   `json:"is_active"`
}

// AccessLevel определяет уровень доступа к действию.
// Возможные значения: "A" (полный), "G" (группа), "M" (свои), "D" (запрещено).
type AccessLevel = string

const (
	AccessAll   AccessLevel = "A" // Полный доступ
	AccessGroup AccessLevel = "G" // Доступ в пределах группы
	AccessOwn   AccessLevel = "M" // Только свои
	AccessDeny  AccessLevel = "D" // Запрещено
)

// EntityRights представляет права на действия с сущностью.
type EntityRights struct {
	View   AccessLevel `json:"view,omitempty"`
	Edit   AccessLevel `json:"edit,omitempty"`
	Add    AccessLevel `json:"add,omitempty"`
	Delete AccessLevel `json:"delete,omitempty"`
	Export AccessLevel `json:"export,omitempty"`
}

// StatusRight представляет права на конкретный статус воронки.
type StatusRight struct {
	EntityType string        `json:"entity_type"`
	PipelineID int           `json:"pipeline_id"`
	StatusID   int           `json:"status_id"`
	Rights     *EntityRights `json:"rights,omitempty"`
}

// Rights представляет права пользователя в системе.
type Rights struct {
	Leads         *EntityRights `json:"leads,omitempty"`
	Contacts      *EntityRights `json:"contacts,omitempty"`
	Companies     *EntityRights `json:"companies,omitempty"`
	Tasks         *EntityRights `json:"tasks,omitempty"`
	MailAccess    bool          `json:"mail_access,omitempty"`
	CatalogAccess bool          `json:"catalog_access,omitempty"`
	StatusRights  []StatusRight `json:"status_rights,omitempty"`
	IsAdmin       bool          `json:"is_admin"`
	IsManager     bool          `json:"is_manager"`
	IsFree        bool          `json:"is_free,omitempty"`
	IsActive      bool          `json:"is_active,omitempty"`
	GroupID       *int          `json:"group_id,omitempty"`
	RoleID        *int          `json:"role_id,omitempty"`
}

// Get получает пользователя по его ID.
func Get(ctx context.Context, apiClient *client.Client, userID int) (*User, error) {
	url := fmt.Sprintf("%s/api/v4/users/%d", apiClient.GetBaseURL(), userID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// GetCurrent получает информацию о текущем пользователе (владельце API-ключа).
// Сначала запрашивает /api/v4/account для получения current_user_id, затем /api/v4/users/{id}.
func GetCurrent(ctx context.Context, apiClient *client.Client) (*User, error) {
	accountURL := fmt.Sprintf("%s/api/v4/account", apiClient.GetBaseURL())
	req, err := http.NewRequestWithContext(ctx, "GET", accountURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код при запросе аккаунта: %d", resp.StatusCode)
	}

	var account Account
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		return nil, err
	}

	if account.CurrentUserID == 0 {
		return nil, fmt.Errorf("не удалось определить текущего пользователя")
	}

	return Get(ctx, apiClient, account.CurrentUserID)
}

// List получает список пользователей с возможностью фильтрации и пагинации.
func List(ctx context.Context, apiClient *client.Client, limit int, page int) ([]User, error) {
	url := fmt.Sprintf("%s/api/v4/users?limit=%d&page=%d", apiClient.GetBaseURL(), limit, page)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.DoRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус-код: %d", resp.StatusCode)
	}

	var users struct {
		Embedded struct {
			Items []User `json:"users"`
		} `json:"_embedded"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}

	return users.Embedded.Items, nil
}
