package access_rights

import (
	"fmt"
	"net/http"
	"testing"
)

// TestCreateAccessRightWithRequester проверяет функцию для создания прав доступа
func TestCreateAccessRightWithRequester(t *testing.T) {
	// Подготавливаем тестовые данные
	newRight := &AccessRight{
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
		// Создаем мок-клиент
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

		// Вызываем тестируемую функцию
		result, err := CreateAccessRightWithRequester(mockClient, newRight)

		// Проверяем результаты
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
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("POST", "/api/v4/access_rights", http.StatusBadRequest, `{"error": "Bad request"}`, nil)

		// Вызываем тестируемую функцию
		_, err := CreateAccessRightWithRequester(mockClient, newRight)

		// Проверяем результаты
		if err == nil {
			t.Fatalf("Ожидалась ошибка, но получен nil")
		}
	})
}

// TestUpdateAccessRightWithRequester проверяет функцию для обновления прав доступа
func TestUpdateAccessRightWithRequester(t *testing.T) {
	// Подготавливаем тестовые данные
	accessRightID := 456
	updateData := &AccessRight{
		ID:   accessRightID, // Добавляем ID в объект
		Name: "Обновленное право",
		Rights: Rights{
			Leads: EntityRights{
				View: true,
			},
		},
	}

	t.Run("Success", func(t *testing.T) {
		// Создаем мок-клиент
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

		// Вызываем тестируемую функцию с правильной сигнатурой
		result, err := UpdateAccessRightWithRequester(mockClient, updateData)

		// Проверяем результаты
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
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("PATCH", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusNotFound, `{"error": "Not found"}`, nil)

		// Вызываем тестируемую функцию с правильной сигнатурой
		_, err := UpdateAccessRightWithRequester(mockClient, updateData)

		// Проверяем результаты
		if err == nil {
			t.Fatalf("Ожидалась ошибка, но получен nil")
		}
	})
}

// TestDeleteAccessRightWithRequester проверяет функцию для удаления прав доступа
func TestDeleteAccessRightWithRequester(t *testing.T) {
	// ID права доступа для теста
	accessRightID := 789

	t.Run("Success", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("DELETE", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusOK, `{"success": true}`, nil)

		// Вызываем тестируемую функцию
		err := DeleteAccessRightWithRequester(mockClient, accessRightID)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при удалении права доступа: %v", err)
		}
	})

	t.Run("Error", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient()
		mockClient.AddResponse("DELETE", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusNotFound, `{"error": "Not found"}`, nil)

		// Вызываем тестируемую функцию
		err := DeleteAccessRightWithRequester(mockClient, accessRightID)

		// Проверяем результаты
		if err == nil {
			t.Fatalf("Ожидалась ошибка, но получен nil")
		}
	})
}

// TestAddUsersToAccessRightWithRequester проверяет функцию для добавления пользователей к правам доступа
func TestAddUsersToAccessRightWithRequester(t *testing.T) {
	// ID права доступа для теста
	accessRightID := 123

	// Тестовые ID пользователей
	userIDs := []int{101, 102}

	t.Run("Success", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient()
		
		// Добавляем ответ для запроса получения текущего права доступа
		mockClient.AddResponse("GET", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusOK, fmt.Sprintf(`{
			"id": %d,
			"name": "Тестовое право",
			"user_ids": [103, 104]
		}`, accessRightID), nil)
		
		// Добавляем ответ для запроса обновления права доступа
		mockClient.AddResponse("PATCH", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusOK, fmt.Sprintf(`{
			"id": %d,
			"name": "Тестовое право",
			"user_ids": [103, 104, 101, 102]
		}`, accessRightID), nil)

		// Вызываем тестируемую функцию
		result, err := AddUsersToAccessRightWithRequester(mockClient, accessRightID, userIDs)

		// Проверяем результаты
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
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient()
		
		// Добавляем ошибку для запроса получения текущего права доступа
		mockClient.AddResponse("GET", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusNotFound, `{"error": "Not found"}`, nil)

		// Вызываем тестируемую функцию
		_, err := AddUsersToAccessRightWithRequester(mockClient, accessRightID, userIDs)

		// Проверяем результаты
		if err == nil {
			t.Fatalf("Ожидалась ошибка, но получен nil")
		}
	})

	t.Run("ErrorOnUpdate", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient()
		
		// Добавляем ответ для запроса получения текущего права доступа
		mockClient.AddResponse("GET", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusOK, fmt.Sprintf(`{
			"id": %d,
			"name": "Тестовое право",
			"user_ids": [103, 104]
		}`, accessRightID), nil)
		
		// Добавляем ошибку для запроса обновления права доступа
		mockClient.AddResponse("PATCH", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusInternalServerError, `{"error": "Server error"}`, nil)

		// Вызываем тестируемую функцию
		_, err := AddUsersToAccessRightWithRequester(mockClient, accessRightID, userIDs)

		// Проверяем результаты
		if err == nil {
			t.Fatalf("Ожидалась ошибка, но получен nil")
		}
	})
}

// TestRemoveUsersFromAccessRightWithRequester проверяет функцию для удаления пользователей из прав доступа
func TestRemoveUsersFromAccessRightWithRequester(t *testing.T) {
	// ID права доступа для теста
	accessRightID := 456

	// Тестовые ID пользователей
	userIDs := []int{101, 102}

	t.Run("Success", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient()
		
		// Добавляем ответ для запроса получения текущего права доступа
		mockClient.AddResponse("GET", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusOK, fmt.Sprintf(`{
			"id": %d,
			"name": "Тестовое право",
			"user_ids": [101, 102, 103, 104]
		}`, accessRightID), nil)
		
		// Добавляем ответ для запроса обновления права доступа
		mockClient.AddResponse("PATCH", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusOK, fmt.Sprintf(`{
			"id": %d,
			"name": "Тестовое право",
			"user_ids": [103, 104]
		}`, accessRightID), nil)

		// Вызываем тестируемую функцию
		result, err := RemoveUsersFromAccessRightWithRequester(mockClient, accessRightID, userIDs)

		// Проверяем результаты
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
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient()
		
		// Добавляем ошибку для запроса получения текущего права доступа
		mockClient.AddResponse("GET", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusNotFound, `{"error": "Not found"}`, nil)

		// Вызываем тестируемую функцию
		_, err := RemoveUsersFromAccessRightWithRequester(mockClient, accessRightID, userIDs)

		// Проверяем результаты
		if err == nil {
			t.Fatalf("Ожидалась ошибка, но получен nil")
		}
	})

	t.Run("ErrorOnUpdate", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := NewAdvancedMockClient()
		
		// Добавляем ответ для запроса получения текущего права доступа
		mockClient.AddResponse("GET", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusOK, fmt.Sprintf(`{
			"id": %d,
			"name": "Тестовое право",
			"user_ids": [101, 102, 103, 104]
		}`, accessRightID), nil)
		
		// Добавляем ошибку для запроса обновления права доступа
		mockClient.AddResponse("PATCH", fmt.Sprintf("/api/v4/access_rights/%d", accessRightID), http.StatusInternalServerError, `{"error": "Server error"}`, nil)

		// Вызываем тестируемую функцию
		_, err := RemoveUsersFromAccessRightWithRequester(mockClient, accessRightID, userIDs)

		// Проверяем результаты
		if err == nil {
			t.Fatalf("Ожидалась ошибка, но получен nil")
		}
	})
}
