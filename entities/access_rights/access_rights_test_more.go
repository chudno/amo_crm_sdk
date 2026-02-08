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

// TestCreateAccessRight проверяет создание нового права доступа
func TestCreateAccessRight(t *testing.T) {
	newAccessRight := &Right{
		Name: "Тестовое право доступа",
		Type: TypeGroup,
		Rights: Rights{
			Leads: EntityRights{
				View: true,
				Edit: true,
				Add:  true,
			},
			Contacts: EntityRights{
				View: true,
				Edit: true,
				Add:  true,
			},
		},
		UserIDs: []int{101, 102},
	}

	successResponse := `{
		"id": 789,
		"name": "Тестовое право доступа",
		"type": "group",
		"rights": {
			"leads": {
				"view": true,
				"edit": true,
				"add": true,
				"delete": false,
				"export": false
			},
			"contacts": {
				"view": true,
				"edit": true,
				"add": true,
				"delete": false,
				"export": false
			}
		},
		"created_by": 123,
		"updated_by": 123,
		"created_at": 1609459200,
		"updated_at": 1609459200,
		"account_id": 12345,
		"user_ids": [101, 102]
	}`

	t.Run("Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Errorf("Ожидался метод POST, получен %s", r.Method)
			}

			expectedPath := "/api/v4/access_rights"
			if r.URL.Path != expectedPath {
				t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
			}

			var requestBody Right
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&requestBody); err != nil {
				t.Errorf("Ошибка при декодировании тела запроса: %v", err)
			}

			if requestBody.Name != newAccessRight.Name {
				t.Errorf("Ожидалось имя '%s', получено '%s'", newAccessRight.Name, requestBody.Name)
			}

			if requestBody.Type != newAccessRight.Type {
				t.Errorf("Ожидался тип '%s', получен '%s'", newAccessRight.Type, requestBody.Type)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(successResponse))
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		createdRight, err := Create(context.Background(), apiClient, newAccessRight)

		if err != nil {
			t.Fatalf("Ошибка при создании права доступа: %v", err)
		}

		if createdRight.ID != 789 {
			t.Errorf("Ожидался ID 789, получен %d", createdRight.ID)
		}

		if createdRight.Name != newAccessRight.Name {
			t.Errorf("Ожидалось имя '%s', получено '%s'", newAccessRight.Name, createdRight.Name)
		}

		if len(createdRight.UserIDs) != 2 {
			t.Errorf("Ожидалось 2 пользователя, получено %d", len(createdRight.UserIDs))
		}
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("POST", "/api/v4/access_rights", http.StatusBadRequest, `{"error": "Invalid request"}`, nil)

		_, err := CreateWithRequester(context.Background(), mockClient, newAccessRight)

		if err == nil {
			t.Fatalf("Ожидалась ошибка, но получен nil")
		}
	})
}

// TestUpdateAccessRight проверяет обновление существующего права доступа
func TestUpdateAccessRight(t *testing.T) {
	accessRightID := 123

	updateAccessRight := &Right{
		ID:   accessRightID,
		Name: "Обновленное право доступа",
		Rights: Rights{
			Leads: EntityRights{
				View:   true,
				Edit:   true,
				Add:    true,
				Delete: true,
				Export: true,
			},
		},
		UserIDs: []int{101, 102, 103},
	}

	successResponse := fmt.Sprintf(`{
		"id": %d,
		"name": "Обновленное право доступа",
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
		"user_ids": [101, 102, 103]
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

			var requestBody Right
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&requestBody); err != nil {
				t.Errorf("Ошибка при декодировании тела запроса: %v", err)
			}

			if requestBody.Name != updateAccessRight.Name {
				t.Errorf("Ожидалось имя '%s', получено '%s'", updateAccessRight.Name, requestBody.Name)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(successResponse))
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		updatedRight, err := Update(context.Background(), apiClient, updateAccessRight)

		if err != nil {
			t.Fatalf("Ошибка при обновлении права доступа: %v", err)
		}

		if updatedRight.ID != accessRightID {
			t.Errorf("Ожидался ID %d, получен %d", accessRightID, updatedRight.ID)
		}

		if updatedRight.Name != updateAccessRight.Name {
			t.Errorf("Ожидалось имя '%s', получено '%s'", updateAccessRight.Name, updatedRight.Name)
		}

		if !updatedRight.Rights.Leads.Delete {
			t.Errorf("Ожидалось право удаления сделок")
		}

		if len(updatedRight.UserIDs) != 3 {
			t.Errorf("Ожидалось 3 пользователя, получено %d", len(updatedRight.UserIDs))
		}
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("PATCH", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusBadRequest, `{"error": "Invalid request"}`, nil)

		_, err := UpdateWithRequester(context.Background(), mockClient, updateAccessRight)

		if err == nil {
			t.Fatalf("Ожидалась ошибка, но получен nil")
		}
	})

	t.Run("EmptyID", func(t *testing.T) {
		emptyIDRight := &Right{
			Name: "Тест с пустым ID",
			Rights: Rights{
				Leads: EntityRights{
					View: true,
				},
			},
		}

		mockClient := NewAdvancedMockClient()

		_, err := UpdateWithRequester(context.Background(), mockClient, emptyIDRight)

		if err == nil {
			t.Fatalf("Ожидалась ошибка, но получен nil")
		}
	})
}

// TestDeleteAccessRight проверяет удаление права доступа
func TestDeleteAccessRight(t *testing.T) {
	accessRightID := 123

	t.Run("Success", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("DELETE", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusNoContent, "", nil)

		err := DeleteWithRequester(context.Background(), mockClient, accessRightID)

		if err != nil {
			t.Fatalf("Ошибка при удалении права доступа: %v", err)
		}
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("DELETE", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusForbidden, `{"error": "Insufficient permissions"}`, nil)

		err := DeleteWithRequester(context.Background(), mockClient, accessRightID)

		if err == nil {
			t.Fatalf("Ожидалась ошибка, но получен nil")
		}
	})
}
