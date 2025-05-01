package unsorted

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

func TestAcceptUnsortedContact(t *testing.T) {
	// Тестовый UID неразобранной заявки
	unsortedUID := "test-unsorted-contact-uid-123"
	responsibleUserID := 456

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/contacts/unsorted/%s/accept", unsortedUID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"_links": {
				"contact": {
					"id": 789
				}
			}
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	contactID, err := AcceptUnsortedContact(apiClient, unsortedUID, responsibleUserID)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при принятии неразобранного контакта: %v", err)
	}

	if contactID != 789 {
		t.Errorf("Ожидался ID контакта 789, получен %d", contactID)
	}
}

func TestDeclineUnsortedContact(t *testing.T) {
	// Тестовый UID неразобранной заявки
	unsortedUID := "test-unsorted-contact-uid-123"

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "DELETE" {
			t.Errorf("Ожидался метод DELETE, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/contacts/unsorted/%s/decline", unsortedUID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	err := DeclineUnsortedContact(apiClient, unsortedUID)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при отклонении неразобранного контакта: %v", err)
	}
}

func TestLinkUnsortedLeadWithContact(t *testing.T) {
	// Тестовые данные
	unsortedUID := "test-unsorted-lead-uid-123"
	contactID := 456

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/leads/unsorted/%s/link", unsortedUID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	err := LinkUnsortedLeadWithContact(apiClient, unsortedUID, contactID)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при связывании неразобранной заявки с контактом: %v", err)
	}
}

func TestLinkUnsortedLeadWithCompany(t *testing.T) {
	// Тестовые данные
	unsortedUID := "test-unsorted-lead-uid-123"
	companyID := 456

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/leads/unsorted/%s/link", unsortedUID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	err := LinkUnsortedLeadWithCompany(apiClient, unsortedUID, companyID)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при связывании неразобранной заявки с компанией: %v", err)
	}
}

func TestLinkUnsortedContactWithCompany(t *testing.T) {
	// Тестовые данные
	unsortedUID := "test-unsorted-contact-uid-123"
	companyID := 456

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/contacts/unsorted/%s/link", unsortedUID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	err := LinkUnsortedContactWithCompany(apiClient, unsortedUID, companyID)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при связывании неразобранного контакта с компанией: %v", err)
	}
}
