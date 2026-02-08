package access_rights

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

// TestGetAccessRights проверяет получение списка прав доступа
func TestGetAccessRights(t *testing.T) {
	successResponse := `{
		"page": 1,
		"per_page": 50,
		"_embedded": {
			"access_rights": [
				{
					"id": 123,
					"name": "Менеджеры продаж",
					"type": "group",
					"rights": {
						"leads": {
							"view": true,
							"edit": true,
							"add": true,
							"delete": false,
							"export": true
						},
						"contacts": {
							"view": true,
							"edit": true,
							"add": true,
							"delete": false,
							"export": true
						}
					},
					"created_by": 789,
					"updated_by": 789,
					"created_at": 1609459200,
					"updated_at": 1609459200,
					"account_id": 12345,
					"user_ids": [101, 102]
				},
				{
					"id": 456,
					"name": "Администраторы",
					"type": "group",
					"rights": {
						"leads": {
							"view": true,
							"edit": true,
							"add": true,
							"delete": true,
							"export": true
						},
						"contacts": {
							"view": true,
							"edit": true,
							"add": true,
							"delete": true,
							"export": true
						},
						"settings": {
							"view": true,
							"edit": true
						}
					},
					"created_by": 789,
					"updated_by": 789,
					"created_at": 1609459200,
					"updated_at": 1609459200,
					"account_id": 12345,
					"user_ids": [201, 202]
				}
			]
		}
	}`

	emptyResponse := `{
		"page": 1,
		"per_page": 50,
		"_embedded": {
			"access_rights": []
		}
	}`

	t.Run("Success", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("GET", "/api/v4/access_rights", http.StatusOK, successResponse, nil)

		accessType := TypeGroup

		rights, err := ListWithRequester(context.Background(), mockClient, 1, 50, WithType(accessType))

		if err != nil {
			t.Fatalf("Ошибка при получении прав доступа: %v", err)
		}

		if len(rights) != 2 {
			t.Errorf("Ожидалось 2 права доступа, получено %d", len(rights))
		}

		if rights[0].ID != 123 {
			t.Errorf("Ожидался ID 123, получен %d", rights[0].ID)
		}

		if rights[0].Name != "Менеджеры продаж" {
			t.Errorf("Ожидалось имя 'Менеджеры продаж', получено '%s'", rights[0].Name)
		}

		if rights[0].Type != TypeGroup {
			t.Errorf("Ожидался тип '%s', получен '%s'", TypeGroup, rights[0].Type)
		}

		if !rights[0].Rights.Leads.View {
			t.Errorf("Ожидалось право просмотра сделок")
		}

		if len(rights[0].UserIDs) != 2 {
			t.Errorf("Ожидалось 2 пользователя, получено %d", len(rights[0].UserIDs))
		}

		if rights[1].ID != 456 {
			t.Errorf("Ожидался ID 456, получен %d", rights[1].ID)
		}

		if rights[1].Name != "Администраторы" {
			t.Errorf("Ожидалось имя 'Администраторы', получено '%s'", rights[1].Name)
		}

		if !rights[1].Rights.Settings.View {
			t.Errorf("Ожидалось право просмотра настроек")
		}
	})

	t.Run("Empty", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("GET", "/api/v4/access_rights", http.StatusOK, emptyResponse, nil)

		rights, err := ListWithRequester(context.Background(), mockClient, 1, 50)

		if err != nil {
			t.Fatalf("Ошибка при получении прав доступа: %v", err)
		}

		if len(rights) != 0 {
			t.Errorf("Ожидалось 0 прав доступа, получено %d", len(rights))
		}
	})

	t.Run("ServerError", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("GET", "/api/v4/access_rights", http.StatusInternalServerError, `{"error": "Internal server error"}`, nil)

		_, err := ListWithRequester(context.Background(), mockClient, 1, 50)

		if err == nil {
			t.Fatalf("Ожидалась ошибка, но получен nil")
		}
	})
}

// TestGetAccessRight проверяет получение информации о конкретном праве доступа
func TestGetAccessRight(t *testing.T) {
	accessRightID := 123

	successResponse := fmt.Sprintf(`{
		"id": %d,
		"name": "Менеджеры продаж",
		"type": "group",
		"rights": {
			"leads": {
				"view": true,
				"edit": true,
				"add": true,
				"delete": false,
				"export": true
			},
			"contacts": {
				"view": true,
				"edit": true,
				"add": true,
				"delete": false,
				"export": true
			}
		},
		"created_by": 789,
		"updated_by": 789,
		"created_at": 1609459200,
		"updated_at": 1609459200,
		"account_id": 12345,
		"user_ids": [101, 102]
	}`, accessRightID)

	t.Run("Success", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("GET", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusOK, successResponse, nil)

		accessRight, err := GetWithRequester(context.Background(), mockClient, accessRightID)

		if err != nil {
			t.Fatalf("Ошибка при получении права доступа: %v", err)
		}

		if accessRight.ID != accessRightID {
			t.Errorf("Ожидался ID %d, получен %d", accessRightID, accessRight.ID)
		}

		if accessRight.Name != "Менеджеры продаж" {
			t.Errorf("Ожидалось имя 'Менеджеры продаж', получено '%s'", accessRight.Name)
		}

		if accessRight.Type != TypeGroup {
			t.Errorf("Ожидался тип '%s', получен '%s'", TypeGroup, accessRight.Type)
		}

		if !accessRight.Rights.Leads.View {
			t.Errorf("Ожидалось право просмотра сделок")
		}

		if len(accessRight.UserIDs) != 2 {
			t.Errorf("Ожидалось 2 пользователя, получено %d", len(accessRight.UserIDs))
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("GET", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusNotFound, `{"error": "Access right not found"}`, nil)

		_, err := GetWithRequester(context.Background(), mockClient, accessRightID)

		if err == nil {
			t.Fatalf("Ожидалась ошибка, но получен nil")
		}
	})
}
