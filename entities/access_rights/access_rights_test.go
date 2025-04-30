package access_rights

import (
	"fmt"
	"net/http"
	"testing"
)

// TestGetAccessRights проверяет получение списка прав доступа
func TestGetAccessRights(t *testing.T) {
	// Подготавливаем ответ для успешного сценария
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

	// Ответ для ситуации, когда прав доступа нет
	emptyResponse := `{
		"page": 1,
		"per_page": 50,
		"_embedded": {
			"access_rights": []
		}
	}`

	// Проверяем успешный сценарий
	t.Run("Success", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("GET", "/api/v4/access_rights", http.StatusOK, successResponse, nil)

		// Создаем фильтр по типу права доступа
		accessType := TypeGroup

		// Вызываем тестируемый метод
		rights, err := GetAccessRightsWithRequester(mockClient, 1, 50, WithType(accessType))

		// Проверяем результаты
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

	// Проверяем сценарий с пустым ответом
	t.Run("Empty", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("GET", "/api/v4/access_rights", http.StatusOK, emptyResponse, nil)

		// Вызываем тестируемый метод
		rights, err := GetAccessRightsWithRequester(mockClient, 1, 50)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при получении прав доступа: %v", err)
		}

		if len(rights) != 0 {
			t.Errorf("Ожидалось 0 прав доступа, получено %d", len(rights))
		}
	})

	// Проверяем сценарий с ошибкой сервера
	t.Run("ServerError", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("GET", "/api/v4/access_rights", http.StatusInternalServerError, `{"error": "Internal server error"}`, nil)

		// Вызываем тестируемый метод
		_, err := GetAccessRightsWithRequester(mockClient, 1, 50)

		// Проверяем результаты
		if err == nil {
			t.Fatalf("Ожидалась ошибка, но получен nil")
		}
	})
}

// TestGetAccessRight проверяет получение информации о конкретном праве доступа
func TestGetAccessRight(t *testing.T) {
	// ID права доступа для теста
	accessRightID := 123

	// Подготавливаем ответ для успешного сценария
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

	// Проверяем успешный сценарий
	t.Run("Success", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("GET", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusOK, successResponse, nil)

		// Вызываем тестируемый метод
		accessRight, err := GetAccessRightWithRequester(mockClient, accessRightID)

		// Проверяем результаты
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

	// Проверяем сценарий с ошибкой
	t.Run("NotFound", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("GET", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusNotFound, `{"error": "Access right not found"}`, nil)

		// Вызываем тестируемый метод
		_, err := GetAccessRightWithRequester(mockClient, accessRightID)

		// Проверяем результаты
		if err == nil {
			t.Fatalf("Ожидалась ошибка, но получен nil")
		}
	})
}
