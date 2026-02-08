package access_rights

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

// TestCreateAccessRightWithRequester проверяет функцию для создания прав доступа
func TestCreateAccessRightWithRequester(t *testing.T) {
	newRight := &Right{
		Name: "Новое право доступа",
		Type: TypeGroup,
		Rights: Rights{
			Leads: EntityRights{
				View: true,
				Edit: true,
			},
		},
	}

	t.Run("Success", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("POST", "/api/v4/access_rights", http.StatusOK, `{
			"id": 789,
			"name": "Новое право доступа",
			"type": "group",
			"rights": {
				"leads": {
					"view": true,
					"edit": true
				}
			}
		}`, nil)

		result, err := CreateWithRequester(context.Background(), mockClient, newRight)

		if err != nil {
			t.Fatalf("Ошибка при создании права доступа: %v", err)
		}

		if result.ID != 789 {
			t.Errorf("Ожидался ID 789, получен %d", result.ID)
		}

		if result.Name != "Новое право доступа" {
			t.Errorf("Ожидалось имя 'Новое право доступа', получено '%s'", result.Name)
		}
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("POST", "/api/v4/access_rights", http.StatusBadRequest, `{"error": "Bad request"}`, nil)

		_, err := CreateWithRequester(context.Background(), mockClient, newRight)

		if err == nil {
			t.Fatalf("Ожидалась ошибка, но получен nil")
		}
	})
}

// TestUpdateAccessRightWithRequester проверяет функцию для обновления прав доступа
func TestUpdateAccessRightWithRequester(t *testing.T) {
	accessRightID := 456
	updateData := &Right{
		ID:   accessRightID, // Добавляем ID в объект
		Name: "Обновленное право",
		Rights: Rights{
			Leads: EntityRights{
				View: true,
			},
		},
	}

	t.Run("Success", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("PATCH", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusOK, fmt.Sprintf(`{
			"id": %d,
			"name": "Обновленное право",
			"rights": {
				"leads": {
					"view": true
				}
			}
		}`, accessRightID), nil)

		result, err := UpdateWithRequester(context.Background(), mockClient, updateData)

		if err != nil {
			t.Fatalf("Ошибка при обновлении права доступа: %v", err)
		}

		if result.ID != accessRightID {
			t.Errorf("Ожидался ID %d, получен %d", accessRightID, result.ID)
		}

		if result.Name != "Обновленное право" {
			t.Errorf("Ожидалось имя 'Обновленное право', получено '%s'", result.Name)
		}
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("PATCH", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusNotFound, `{"error": "Not found"}`, nil)

		_, err := UpdateWithRequester(context.Background(), mockClient, updateData)

		if err == nil {
			t.Fatalf("Ожидалась ошибка, но получен nil")
		}
	})
}

// TestDeleteAccessRightWithRequester проверяет функцию для удаления прав доступа
func TestDeleteAccessRightWithRequester(t *testing.T) {
	accessRightID := 789

	t.Run("Success", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("DELETE", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusOK, `{"success": true}`, nil)

		err := DeleteWithRequester(context.Background(), mockClient, accessRightID)

		if err != nil {
			t.Fatalf("Ошибка при удалении права доступа: %v", err)
		}
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("DELETE", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusNotFound, `{"error": "Not found"}`, nil)

		err := DeleteWithRequester(context.Background(), mockClient, accessRightID)

		if err == nil {
			t.Fatalf("Ожидалась ошибка, но получен nil")
		}
	})
}

// TestAddUsersToAccessRightWithRequester проверяет функцию для добавления пользователей к правам доступа
func TestAddUsersToAccessRightWithRequester(t *testing.T) {
	accessRightID := 123

	userIDs := []int{101, 102}

	t.Run("Success", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()

		mockClient.AddResponse("GET", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusOK, fmt.Sprintf(`{
			"id": %d,
			"name": "Тестовое право",
			"user_ids": [103, 104]
		}`, accessRightID), nil)

		mockClient.AddResponse("PATCH", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusOK, fmt.Sprintf(`{
			"id": %d,
			"name": "Тестовое право",
			"user_ids": [103, 104, 101, 102]
		}`, accessRightID), nil)

		result, err := AddUsersWithRequester(context.Background(), mockClient, accessRightID, userIDs)

		if err != nil {
			t.Fatalf("Ошибка при добавлении пользователей: %v", err)
		}

		if result.ID != accessRightID {
			t.Errorf("Ожидался ID %d, получен %d", accessRightID, result.ID)
		}

		if len(result.UserIDs) != 4 {
			t.Errorf("Ожидалось 4 пользователя, получено %d", len(result.UserIDs))
		}
	})

	t.Run("ErrorOnGet", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()

		// Добавляем ошибку для запроса получения текущего права доступа
		mockClient.AddResponse("GET", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusNotFound, `{"error": "Not found"}`, nil)

		_, err := AddUsersWithRequester(context.Background(), mockClient, accessRightID, userIDs)

		if err == nil {
			t.Fatalf("Ожидалась ошибка, но получен nil")
		}
	})

	t.Run("ErrorOnUpdate", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()

		mockClient.AddResponse("GET", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusOK, fmt.Sprintf(`{
			"id": %d,
			"name": "Тестовое право",
			"user_ids": [103, 104]
		}`, accessRightID), nil)

		// Добавляем ошибку для запроса обновления права доступа
		mockClient.AddResponse("PATCH", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusInternalServerError, `{"error": "Server error"}`, nil)

		_, err := AddUsersWithRequester(context.Background(), mockClient, accessRightID, userIDs)

		if err == nil {
			t.Fatalf("Ожидалась ошибка, но получен nil")
		}
	})
}

// TestRemoveUsersFromAccessRightWithRequester проверяет функцию для удаления пользователей из прав доступа
func TestRemoveUsersFromAccessRightWithRequester(t *testing.T) {
	accessRightID := 456

	userIDs := []int{101, 102}

	t.Run("Success", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()

		mockClient.AddResponse("GET", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusOK, fmt.Sprintf(`{
			"id": %d,
			"name": "Тестовое право",
			"user_ids": [101, 102, 103, 104]
		}`, accessRightID), nil)

		mockClient.AddResponse("PATCH", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusOK, fmt.Sprintf(`{
			"id": %d,
			"name": "Тестовое право",
			"user_ids": [103, 104]
		}`, accessRightID), nil)

		result, err := RemoveUsersWithRequester(context.Background(), mockClient, accessRightID, userIDs)

		if err != nil {
			t.Fatalf("Ошибка при удалении пользователей: %v", err)
		}

		if result.ID != accessRightID {
			t.Errorf("Ожидался ID %d, получен %d", accessRightID, result.ID)
		}

		if len(result.UserIDs) != 2 {
			t.Errorf("Ожидалось 2 пользователя, получено %d", len(result.UserIDs))
		}
	})

	t.Run("ErrorOnGet", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()

		// Добавляем ошибку для запроса получения текущего права доступа
		mockClient.AddResponse("GET", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusNotFound, `{"error": "Not found"}`, nil)

		_, err := RemoveUsersWithRequester(context.Background(), mockClient, accessRightID, userIDs)

		if err == nil {
			t.Fatalf("Ожидалась ошибка, но получен nil")
		}
	})

	t.Run("ErrorOnUpdate", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()

		mockClient.AddResponse("GET", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusOK, fmt.Sprintf(`{
			"id": %d,
			"name": "Тестовое право",
			"user_ids": [101, 102, 103, 104]
		}`, accessRightID), nil)

		// Добавляем ошибку для запроса обновления права доступа
		mockClient.AddResponse("PATCH", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusInternalServerError, `{"error": "Server error"}`, nil)

		_, err := RemoveUsersWithRequester(context.Background(), mockClient, accessRightID, userIDs)

		if err == nil {
			t.Fatalf("Ожидалась ошибка, но получен nil")
		}
	})
}
