//go:build integration

package integration

import (
	"testing"

	"github.com/chudno/amo_crm_sdk/entities/users"
)

func TestIntegration_GetCurrentUser(t *testing.T) {
	apiClient, ctx := setupClient(t)

	user, err := users.GetCurrent(ctx, apiClient)
	if err != nil {
		t.Fatalf("GetCurrentUser: %v", err)
	}

	if user.ID == 0 {
		t.Error("ID пользователя не должен быть 0")
	}
	if user.Name == "" {
		t.Error("Имя пользователя не должно быть пустым")
	}

	t.Logf("Текущий пользователь: ID=%d, Name=%q, Email=%q", user.ID, user.Name, user.Email)
}

func TestIntegration_ListUsers(t *testing.T) {
	apiClient, ctx := setupClient(t)

	usersList, err := users.List(ctx, apiClient, 10, 1)
	if err != nil {
		t.Fatalf("ListUsers: %v", err)
	}

	if len(usersList) == 0 {
		t.Fatal("Список пользователей не должен быть пустым")
	}

	t.Logf("Получено %d пользователей", len(usersList))
	for _, u := range usersList {
		t.Logf("  ID=%d, Name=%q", u.ID, u.Name)
	}
}

func TestIntegration_GetUser(t *testing.T) {
	apiClient, ctx := setupClient(t)

	// Сначала получаем текущего пользователя, чтобы узнать валидный ID
	current, err := users.GetCurrent(ctx, apiClient)
	if err != nil {
		t.Fatalf("GetCurrentUser: %v", err)
	}

	user, err := users.Get(ctx, apiClient, current.ID)
	if err != nil {
		t.Fatalf("GetUser(%d): %v", current.ID, err)
	}

	if user.ID != current.ID {
		t.Errorf("Ожидался ID=%d, получен ID=%d", current.ID, user.ID)
	}

	t.Logf("Пользователь: ID=%d, Name=%q", user.ID, user.Name)
}
