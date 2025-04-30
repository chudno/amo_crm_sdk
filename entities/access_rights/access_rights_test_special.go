package access_rights

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

// TestSetEntityRights проверяет обновление прав для конкретной сущности
func TestSetEntityRights(t *testing.T) {
	// ID права доступа для теста
	accessRightID := 123

	// Тип сущности для теста
	entityType := EntityLead

	// Права для установки
	entityRights := EntityRights{
		View:   true,
		Edit:   true,
		Add:    true,
		Delete: true,
		Export: true,
	}

	// Подготавливаем ответ для успешного сценария
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

	// Проверяем успешный сценарий
	t.Run("Success", func(t *testing.T) {
		// Создаем тестовый сервер для проверки тела запроса
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Проверяем метод запроса
			if r.Method != "PATCH" {
				t.Errorf("Ожидался метод PATCH, получен %s", r.Method)
			}

			// Проверяем путь запроса
			expectedPath := fmt.Sprintf("/api/v4/access_rights/%d", accessRightID)
			if r.URL.Path != expectedPath {
				t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
			}

			// Проверяем тело запроса
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

			// Отправляем ответ
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(successResponse))
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемый метод
		updatedRight, err := SetEntityRights(apiClient, accessRightID, entityType, entityRights)

		// Проверяем результаты
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

	// Проверяем сценарий с ошибкой
	t.Run("Error", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("PATCH", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusBadRequest, `{"error": "Invalid request"}`, nil)

		// Вызываем тестируемый метод
		_, err := SetEntityRightsWithRequester(mockClient, accessRightID, entityType, entityRights)

		// Проверяем результаты
		if err == nil {
			t.Fatalf("Ожидалась ошибка, но получен nil")
		}
	})
}

// TestAddUsersToAccessRight проверяет добавление пользователей в право доступа
func TestAddUsersToAccessRight(t *testing.T) {
	// ID права доступа для теста
	accessRightID := 123

	// Существующие пользователи
	existingUsers := []int{101, 102}

	// Новые пользователи для добавления
	newUsers := []int{103, 104}

	// Подготавливаем ответ для получения существующего права доступа
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

	// Подготавливаем ответ для успешного сценария после добавления пользователей
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

	// Проверяем успешный сценарий
	t.Run("Success", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient()

		// Добавляем ответ для получения существующего права доступа
		mockClient.AddResponse("GET", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusOK, existingRightResponse, nil)

		// Добавляем ответ для обновления права доступа
		mockClient.AddResponse("PATCH", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusOK, updatedRightResponse, nil)

		// Вызываем тестируемый метод
		updatedRight, err := AddUsersToAccessRightWithRequester(mockClient, accessRightID, newUsers)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при добавлении пользователей: %v", err)
		}

		if updatedRight.ID != accessRightID {
			t.Errorf("Ожидался ID %d, получен %d", accessRightID, updatedRight.ID)
		}

		if len(updatedRight.UserIDs) != 4 {
			t.Errorf("Ожидалось 4 пользователя, получено %d", len(updatedRight.UserIDs))
		}

		// Проверяем наличие всех пользователей
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

	// Проверяем сценарий с ошибкой при получении существующего права доступа
	t.Run("ErrorGettingAccessRight", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient()

		// Добавляем ответ с ошибкой для получения существующего права доступа
		mockClient.AddResponse("GET", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusNotFound, `{"error": "Access right not found"}`, nil)

		// Вызываем тестируемый метод
		_, err := AddUsersToAccessRightWithRequester(mockClient, accessRightID, newUsers)

		// Проверяем результаты
		if err == nil {
			t.Fatalf("Ожидалась ошибка, но получен nil")
		}
	})

	// Проверяем сценарий с ошибкой при обновлении права доступа
	t.Run("ErrorUpdatingAccessRight", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient()

		// Добавляем ответ для получения существующего права доступа
		mockClient.AddResponse("GET", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusOK, existingRightResponse, nil)

		// Добавляем ответ с ошибкой для обновления права доступа
		mockClient.AddResponse("PATCH", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusBadRequest, `{"error": "Invalid request"}`, nil)

		// Вызываем тестируемый метод
		_, err := AddUsersToAccessRightWithRequester(mockClient, accessRightID, newUsers)

		// Проверяем результаты
		if err == nil {
			t.Fatalf("Ожидалась ошибка, но получен nil")
		}
	})
}

// TestRemoveUsersFromAccessRight проверяет удаление пользователей из права доступа
func TestRemoveUsersFromAccessRight(t *testing.T) {
	// ID права доступа для теста
	accessRightID := 123

	// Пользователи для удаления
	usersToRemove := []int{103, 104}

	// Подготавливаем ответ для получения существующего права доступа
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

	// Подготавливаем ответ для успешного сценария после удаления пользователей
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

	// Проверяем успешный сценарий
	t.Run("Success", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient()

		// Добавляем ответ для получения существующего права доступа
		mockClient.AddResponse("GET", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusOK, existingRightResponse, nil)

		// Добавляем ответ для обновления права доступа
		mockClient.AddResponse("PATCH", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusOK, updatedRightResponse, nil)

		// Вызываем тестируемый метод
		updatedRight, err := RemoveUsersFromAccessRightWithRequester(mockClient, accessRightID, usersToRemove)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при удалении пользователей: %v", err)
		}

		if updatedRight.ID != accessRightID {
			t.Errorf("Ожидался ID %d, получен %d", accessRightID, updatedRight.ID)
		}

		if len(updatedRight.UserIDs) != 2 {
			t.Errorf("Ожидалось 2 пользователя, получено %d", len(updatedRight.UserIDs))
		}

		// Проверяем отсутствие удаленных пользователей
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

	// Проверяем сценарий с ошибкой при получении существующего права доступа
	t.Run("ErrorGettingAccessRight", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient()

		// Добавляем ответ с ошибкой для получения существующего права доступа
		mockClient.AddResponse("GET", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusNotFound, `{"error": "Access right not found"}`, nil)

		// Вызываем тестируемый метод
		_, err := RemoveUsersFromAccessRightWithRequester(mockClient, accessRightID, usersToRemove)

		// Проверяем результаты
		if err == nil {
			t.Fatalf("Ожидалась ошибка, но получен nil")
		}
	})

	// Проверяем сценарий с ошибкой при обновлении права доступа
	t.Run("ErrorUpdatingAccessRight", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient()

		// Добавляем ответ для получения существующего права доступа
		mockClient.AddResponse("GET", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusOK, existingRightResponse, nil)

		// Добавляем ответ с ошибкой для обновления права доступа
		mockClient.AddResponse("PATCH", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusBadRequest, `{"error": "Invalid request"}`, nil)

		// Вызываем тестируемый метод
		_, err := RemoveUsersFromAccessRightWithRequester(mockClient, accessRightID, usersToRemove)

		// Проверяем результаты
		if err == nil {
			t.Fatalf("Ожидалась ошибка, но получен nil")
		}
	})
}
