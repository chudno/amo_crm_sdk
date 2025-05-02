package access_rights

import (
	"errors"
	"net/http"
	"testing"
)

// MockClientWithError имитирует клиент с ошибкой при выполнении запроса
type MockClientWithError struct{}

func (m *MockClientWithError) DoRequest(req *http.Request) (*http.Response, error) {
	return nil, errors.New("сетевая ошибка")
}

// TestGetAccessRightsNetworkError проверяет обработку сетевых ошибок при получении списка прав доступа
func TestGetAccessRightsNetworkError(t *testing.T) {
	// Создаём клиент с ошибкой
	mockClient := &MockClientWithError{}

	// Вызываем тестируемую функцию
	_, err := GetAccessRightsWithRequester(mockClient, 1, 50)

	// Проверяем, что вернулась ошибка
	if err == nil {
		t.Fatal("Ожидалась ошибка сети, но получен nil")
	}
}

// TestGetAccessRightNetworkError проверяет обработку сетевых ошибок при получении конкретного права доступа
func TestGetAccessRightNetworkError(t *testing.T) {
	// Создаём клиент с ошибкой
	mockClient := &MockClientWithError{}

	// Вызываем тестируемую функцию
	_, err := GetAccessRightWithRequester(mockClient, 123)

	// Проверяем, что вернулась ошибка
	if err == nil {
		t.Fatal("Ожидалась ошибка сети, но получен nil")
	}
}

// TestCreateAccessRightNetworkError проверяет обработку сетевых ошибок при создании права доступа
func TestCreateAccessRightNetworkError(t *testing.T) {
	// Создаём клиент с ошибкой
	mockClient := &MockClientWithError{}

	// Создаём тестовое право доступа
	accessRight := &AccessRight{
		Name: "Тестовое право",
		Type: TypeGroup,
		Rights: Rights{
			Leads: EntityRights{
				View: true,
				Edit: true,
			},
		},
	}

	// Вызываем тестируемую функцию
	_, err := CreateAccessRightWithRequester(mockClient, accessRight)

	// Проверяем, что вернулась ошибка
	if err == nil {
		t.Fatal("Ожидалась ошибка сети, но получен nil")
	}
}

// TestUpdateAccessRightNetworkError проверяет обработку сетевых ошибок при обновлении права доступа
func TestUpdateAccessRightNetworkError(t *testing.T) {
	// Создаём клиент с ошибкой
	mockClient := &MockClientWithError{}

	// Создаём тестовое право доступа
	accessRight := &AccessRight{
		ID:   123,
		Name: "Тестовое право",
		Type: TypeGroup,
		Rights: Rights{
			Leads: EntityRights{
				View: true,
				Edit: true,
			},
		},
	}

	// Вызываем тестируемую функцию
	_, err := UpdateAccessRightWithRequester(mockClient, accessRight)

	// Проверяем, что вернулась ошибка
	if err == nil {
		t.Fatal("Ожидалась ошибка сети, но получен nil")
	}
}

// TestDeleteAccessRightNetworkError проверяет обработку сетевых ошибок при удалении права доступа
func TestDeleteAccessRightNetworkError(t *testing.T) {
	// Создаём клиент с ошибкой
	mockClient := &MockClientWithError{}

	// Вызываем тестируемую функцию
	err := DeleteAccessRightWithRequester(mockClient, 123)

	// Проверяем, что вернулась ошибка
	if err == nil {
		t.Fatal("Ожидалась ошибка сети, но получен nil")
	}
}

// TestSetEntityRightsNetworkError проверяет обработку сетевых ошибок при установке прав для сущности
func TestSetEntityRightsNetworkError(t *testing.T) {
	// Создаём клиент с ошибкой
	mockClient := &MockClientWithError{}

	// Создаём тестовые права для сущности
	entityRights := EntityRights{
		View: true,
		Edit: true,
		Add:  true,
	}

	// Вызываем тестируемую функцию
	_, err := SetEntityRightsWithRequester(mockClient, 123, EntityLead, entityRights)

	// Проверяем, что вернулась ошибка
	if err == nil {
		t.Fatal("Ожидалась ошибка сети, но получен nil")
	}
}

// TestAddUsersToAccessRightNetworkError проверяет обработку сетевых ошибок при добавлении пользователей к праву доступа
func TestAddUsersToAccessRightNetworkError(t *testing.T) {
	// Создаём клиент с ошибкой
	mockClient := &MockClientWithError{}

	// Вызываем тестируемую функцию
	_, err := AddUsersToAccessRightWithRequester(mockClient, 123, []int{101, 102})

	// Проверяем, что вернулась ошибка
	if err == nil {
		t.Fatal("Ожидалась ошибка сети, но получен nil")
	}
}

// TestRemoveUsersFromAccessRightNetworkError проверяет обработку сетевых ошибок при удалении пользователей из права доступа
func TestRemoveUsersFromAccessRightNetworkError(t *testing.T) {
	// Создаём клиент с ошибкой
	mockClient := &MockClientWithError{}

	// Вызываем тестируемую функцию
	_, err := RemoveUsersFromAccessRightWithRequester(mockClient, 123, []int{101, 102})

	// Проверяем, что вернулась ошибка
	if err == nil {
		t.Fatal("Ожидалась ошибка сети, но получен nil")
	}
}

// TestZeroID проверяет обработку нулевого ID
func TestZeroID(t *testing.T) {
	// Создаём мок-клиент
	mockClient := &MockClientWithError{}

	// Проверяем получение права доступа с нулевым ID
	t.Run("GetAccessRightWithZeroID", func(t *testing.T) {
		_, err := GetAccessRightWithRequester(mockClient, 0)
		if err == nil {
			t.Fatal("Ожидалась ошибка при нулевом ID, но получен nil")
		}
	})

	// Проверяем удаление права доступа с нулевым ID
	t.Run("DeleteAccessRightWithZeroID", func(t *testing.T) {
		err := DeleteAccessRightWithRequester(mockClient, 0)
		if err == nil {
			t.Fatal("Ожидалась ошибка при нулевом ID, но получен nil")
		}
	})

	// Проверяем обновление права доступа с нулевым ID
	t.Run("UpdateAccessRightWithZeroID", func(t *testing.T) {
		_, err := UpdateAccessRightWithRequester(mockClient, &AccessRight{ID: 0})
		if err == nil {
			t.Fatal("Ожидалась ошибка при нулевом ID, но получен nil")
		}
	})
}
