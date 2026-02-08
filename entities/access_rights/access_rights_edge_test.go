package access_rights

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

// MockClientWithError имитирует клиент с ошибкой при выполнении запроса
type MockClientWithError struct{}

func (m *MockClientWithError) DoRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	return nil, errors.New("сетевая ошибка")
}

// TestGetAccessRightsNetworkError проверяет обработку сетевых ошибок при получении списка прав доступа
func TestGetAccessRightsNetworkError(t *testing.T) {
	mockClient := &MockClientWithError{}

	_, err := ListWithRequester(context.Background(), mockClient, 1, 50)

	if err == nil {
		t.Fatal("Ожидалась ошибка сети, но получен nil")
	}
}

// TestGetAccessRightNetworkError проверяет обработку сетевых ошибок при получении конкретного права доступа
func TestGetAccessRightNetworkError(t *testing.T) {
	mockClient := &MockClientWithError{}

	_, err := GetWithRequester(context.Background(), mockClient, 123)

	if err == nil {
		t.Fatal("Ожидалась ошибка сети, но получен nil")
	}
}

// TestCreateAccessRightNetworkError проверяет обработку сетевых ошибок при создании права доступа
func TestCreateAccessRightNetworkError(t *testing.T) {
	mockClient := &MockClientWithError{}

	accessRight := &Right{
		Name: "Тестовое право",
		Type: TypeGroup,
		Rights: Rights{
			Leads: EntityRights{
				View: true,
				Edit: true,
			},
		},
	}

	_, err := CreateWithRequester(context.Background(), mockClient, accessRight)

	if err == nil {
		t.Fatal("Ожидалась ошибка сети, но получен nil")
	}
}

// TestUpdateAccessRightNetworkError проверяет обработку сетевых ошибок при обновлении права доступа
func TestUpdateAccessRightNetworkError(t *testing.T) {
	mockClient := &MockClientWithError{}

	accessRight := &Right{
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

	_, err := UpdateWithRequester(context.Background(), mockClient, accessRight)

	if err == nil {
		t.Fatal("Ожидалась ошибка сети, но получен nil")
	}
}

// TestDeleteAccessRightNetworkError проверяет обработку сетевых ошибок при удалении права доступа
func TestDeleteAccessRightNetworkError(t *testing.T) {
	mockClient := &MockClientWithError{}

	err := DeleteWithRequester(context.Background(), mockClient, 123)

	if err == nil {
		t.Fatal("Ожидалась ошибка сети, но получен nil")
	}
}

// TestSetEntityRightsNetworkError проверяет обработку сетевых ошибок при установке прав для сущности
func TestSetEntityRightsNetworkError(t *testing.T) {
	mockClient := &MockClientWithError{}

	entityRights := EntityRights{
		View: true,
		Edit: true,
		Add:  true,
	}

	_, err := SetEntityRightsWithRequester(context.Background(), mockClient, 123, EntityLead, entityRights)

	if err == nil {
		t.Fatal("Ожидалась ошибка сети, но получен nil")
	}
}

// TestAddUsersToAccessRightNetworkError проверяет обработку сетевых ошибок при добавлении пользователей к праву доступа
func TestAddUsersToAccessRightNetworkError(t *testing.T) {
	mockClient := &MockClientWithError{}

	_, err := AddUsersWithRequester(context.Background(), mockClient, 123, []int{101, 102})

	if err == nil {
		t.Fatal("Ожидалась ошибка сети, но получен nil")
	}
}

// TestRemoveUsersFromAccessRightNetworkError проверяет обработку сетевых ошибок при удалении пользователей из права доступа
func TestRemoveUsersFromAccessRightNetworkError(t *testing.T) {
	mockClient := &MockClientWithError{}

	_, err := RemoveUsersWithRequester(context.Background(), mockClient, 123, []int{101, 102})

	if err == nil {
		t.Fatal("Ожидалась ошибка сети, но получен nil")
	}
}

// TestZeroID проверяет обработку нулевого ID
func TestZeroID(t *testing.T) {
	mockClient := &MockClientWithError{}

	t.Run("GetAccessRightWithZeroID", func(t *testing.T) {
		_, err := GetWithRequester(context.Background(), mockClient, 0)
		if err == nil {
			t.Fatal("Ожидалась ошибка при нулевом ID, но получен nil")
		}
	})

	t.Run("DeleteAccessRightWithZeroID", func(t *testing.T) {
		err := DeleteWithRequester(context.Background(), mockClient, 0)
		if err == nil {
			t.Fatal("Ожидалась ошибка при нулевом ID, но получен nil")
		}
	})

	t.Run("UpdateAccessRightWithZeroID", func(t *testing.T) {
		_, err := UpdateWithRequester(context.Background(), mockClient, &Right{ID: 0})
		if err == nil {
			t.Fatal("Ожидалась ошибка при нулевом ID, но получен nil")
		}
	})
}
