package access_rights

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

// TestSetEntityRights проверяет обновление прав для конкретной сущности
func TestSetEntityRights(t *testing.T) {
	accessRightID := 123

	entityType := EntityLead

	entityRights := EntityRights{
		View:   true,
		Edit:   true,
		Add:    true,
		Delete: true,
		Export: true,
	}

	successResponse := fmt.Sprintf(`{
		"id": %d,
		"name": "Тестовое право",
		"type": "group",
		"rights": {
			"leads": {
				"view": true,
				"edit": true,
				"add": true,
				"delete": true,
				"export": true
			}
		},
		"created_by": 789,
		"updated_by": 789,
		"created_at": 1609459200,
		"updated_at": 1609460000,
		"account_id": 12345,
		"user_ids": [101, 102]
	}`, accessRightID)

	t.Run("Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "PATCH" {
				t.Errorf("Ожидался метод PATCH, получен %s", r.Method)
			}

			expectedPath := fmt.Sprintf("/api/v4/access_rights/%d", accessRightID)
			if r.URL.Path != expectedPath {
				t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
			}

			var requestBody struct {
				Rights map[string]EntityRights `json:"rights"`
			}
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&requestBody); err != nil {
				t.Errorf("Ошибка при декодировании тела запроса: %v", err)
			}

			rights, ok := requestBody.Rights[string(entityType)]
			if !ok {
				t.Errorf("Ожидались права для сущности '%s', но они отсутствуют", entityType)
			}

			if !rights.View || !rights.Edit || !rights.Add || !rights.Delete || !rights.Export {
				t.Errorf("Некоторые права для сущности '%s' не установлены", entityType)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(successResponse))
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		updatedRight, err := SetEntityRights(context.Background(), apiClient, accessRightID, entityType, entityRights)

		if err != nil {
			t.Fatalf("Ошибка при обновлении прав для сущности: %v", err)
		}

		if updatedRight.ID != accessRightID {
			t.Errorf("Ожидался ID %d, получен %d", accessRightID, updatedRight.ID)
		}

		if !updatedRight.Rights.Leads.Delete {
			t.Errorf("Ожидалось право удаления сделок")
		}
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("PATCH", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusBadRequest, `{"error": "Invalid request"}`, nil)

		_, err := SetEntityRightsWithRequester(context.Background(), mockClient, accessRightID, entityType, entityRights)

		if err == nil {
			t.Fatalf("Ожидалась ошибка, но получен nil")
		}
	})
}

// TestAddUsersToAccessRight проверяет добавление пользователей в право доступа
func TestAddUsersToAccessRight(t *testing.T) {
	accessRightID := 123

	existingUsers := []int{101, 102}

	newUsers := []int{103, 104}

	existingRightResponse := fmt.Sprintf(`{
		"id": %d,
		"name": "Тестовое право",
		"type": "group",
		"rights": {
			"leads": {
				"view": true,
				"edit": true
			}
		},
		"user_ids": [101, 102]
	}`, accessRightID)

	updatedRightResponse := fmt.Sprintf(`{
		"id": %d,
		"name": "Тестовое право",
		"type": "group",
		"rights": {
			"leads": {
				"view": true,
				"edit": true
			}
		},
		"user_ids": [101, 102, 103, 104]
	}`, accessRightID)

	t.Run("Success", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()

		mockClient.AddResponse("GET", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusOK, existingRightResponse, nil)

		mockClient.AddResponse("PATCH", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusOK, updatedRightResponse, nil)

		updatedRight, err := AddUsersWithRequester(context.Background(), mockClient, accessRightID, newUsers)

		if err != nil {
			t.Fatalf("Ошибка при добавлении пользователей: %v", err)
		}

		if updatedRight.ID != accessRightID {
			t.Errorf("Ожидался ID %d, получен %d", accessRightID, updatedRight.ID)
		}

		if len(updatedRight.UserIDs) != 4 {
			t.Errorf("Ожидалось 4 пользователя, получено %d", len(updatedRight.UserIDs))
		}

		expectedUsers := append(existingUsers, newUsers...)
		userMap := make(map[int]bool)
		for _, id := range updatedRight.UserIDs {
			userMap[id] = true
		}

		for _, id := range expectedUsers {
			if !userMap[id] {
				t.Errorf("Пользователь с ID %d отсутствует в обновленном праве доступа", id)
			}
		}
	})

	t.Run("ErrorGettingAccessRight", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()

		mockClient.AddResponse("GET", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusNotFound, `{"error": "Access right not found"}`, nil)

		_, err := AddUsersWithRequester(context.Background(), mockClient, accessRightID, newUsers)

		if err == nil {
			t.Fatalf("Ожидалась ошибка, но получен nil")
		}
	})

	t.Run("ErrorUpdatingAccessRight", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()

		mockClient.AddResponse("GET", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusOK, existingRightResponse, nil)

		mockClient.AddResponse("PATCH", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusBadRequest, `{"error": "Invalid request"}`, nil)

		_, err := AddUsersWithRequester(context.Background(), mockClient, accessRightID, newUsers)

		if err == nil {
			t.Fatalf("Ожидалась ошибка, но получен nil")
		}
	})
}

// TestRemoveUsersFromAccessRight проверяет удаление пользователей из права доступа
func TestRemoveUsersFromAccessRight(t *testing.T) {
	accessRightID := 123

	usersToRemove := []int{103, 104}

	existingRightResponse := fmt.Sprintf(`{
		"id": %d,
		"name": "Тестовое право",
		"type": "group",
		"rights": {
			"leads": {
				"view": true,
				"edit": true
			}
		},
		"user_ids": [101, 102, 103, 104]
	}`, accessRightID)

	updatedRightResponse := fmt.Sprintf(`{
		"id": %d,
		"name": "Тестовое право",
		"type": "group",
		"rights": {
			"leads": {
				"view": true,
				"edit": true
			}
		},
		"user_ids": [101, 102]
	}`, accessRightID)

	t.Run("Success", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()

		mockClient.AddResponse("GET", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusOK, existingRightResponse, nil)

		mockClient.AddResponse("PATCH", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusOK, updatedRightResponse, nil)

		updatedRight, err := RemoveUsersWithRequester(context.Background(), mockClient, accessRightID, usersToRemove)

		if err != nil {
			t.Fatalf("Ошибка при удалении пользователей: %v", err)
		}

		if updatedRight.ID != accessRightID {
			t.Errorf("Ожидался ID %d, получен %d", accessRightID, updatedRight.ID)
		}

		if len(updatedRight.UserIDs) != 2 {
			t.Errorf("Ожидалось 2 пользователя, получено %d", len(updatedRight.UserIDs))
		}

		removedUserMap := make(map[int]bool)
		for _, id := range usersToRemove {
			removedUserMap[id] = true
		}

		for _, id := range updatedRight.UserIDs {
			if removedUserMap[id] {
				t.Errorf("Пользователь с ID %d должен был быть удален, но присутствует", id)
			}
		}
	})

	t.Run("ErrorGettingAccessRight", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()

		mockClient.AddResponse("GET", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusNotFound, `{"error": "Access right not found"}`, nil)

		_, err := RemoveUsersWithRequester(context.Background(), mockClient, accessRightID, usersToRemove)

		if err == nil {
			t.Fatalf("Ожидалась ошибка, но получен nil")
		}
	})

	t.Run("ErrorUpdatingAccessRight", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()

		mockClient.AddResponse("GET", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusOK, existingRightResponse, nil)

		mockClient.AddResponse("PATCH", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusBadRequest, `{"error": "Invalid request"}`, nil)

		_, err := RemoveUsersWithRequester(context.Background(), mockClient, accessRightID, usersToRemove)

		if err == nil {
			t.Fatalf("Ожидалась ошибка, но получен nil")
		}
	})
}
