package unsorted

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

func TestAcceptUnsortedContact(t *testing.T) {
	unsortedUID := "test-unsorted-contact-uid-123"
	responsibleUserID := 456

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		expectedPath := fmt.Sprintf("/api/v4/contacts/unsorted/%s/accept", unsortedUID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

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

	apiClient := client.NewClient(server.URL, "test_api_key")

	contactID, err := AcceptContact(context.Background(), apiClient, unsortedUID, responsibleUserID)

	if err != nil {
		t.Fatalf("Ошибка при принятии неразобранного контакта: %v", err)
	}

	if contactID != 789 {
		t.Errorf("Ожидался ID контакта 789, получен %d", contactID)
	}
}

func TestDeclineUnsortedContact(t *testing.T) {
	unsortedUID := "test-unsorted-contact-uid-123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Ожидался метод DELETE, получен %s", r.Method)
		}

		expectedPath := fmt.Sprintf("/api/v4/contacts/unsorted/%s/decline", unsortedUID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	apiClient := client.NewClient(server.URL, "test_api_key")

	err := DeclineContact(context.Background(), apiClient, unsortedUID)

	if err != nil {
		t.Fatalf("Ошибка при отклонении неразобранного контакта: %v", err)
	}
}

func TestLinkUnsortedLeadWithContact(t *testing.T) {
	unsortedUID := "test-unsorted-lead-uid-123"
	contactID := 456

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		expectedPath := fmt.Sprintf("/api/v4/leads/unsorted/%s/link", unsortedUID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	apiClient := client.NewClient(server.URL, "test_api_key")

	err := LinkLeadWithContact(context.Background(), apiClient, unsortedUID, contactID)

	if err != nil {
		t.Fatalf("Ошибка при связывании неразобранной заявки с контактом: %v", err)
	}
}

func TestLinkUnsortedLeadWithCompany(t *testing.T) {
	unsortedUID := "test-unsorted-lead-uid-123"
	companyID := 456

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		expectedPath := fmt.Sprintf("/api/v4/leads/unsorted/%s/link", unsortedUID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	apiClient := client.NewClient(server.URL, "test_api_key")

	err := LinkLeadWithCompany(context.Background(), apiClient, unsortedUID, companyID)

	if err != nil {
		t.Fatalf("Ошибка при связывании неразобранной заявки с компанией: %v", err)
	}
}

func TestLinkUnsortedContactWithCompany(t *testing.T) {
	unsortedUID := "test-unsorted-contact-uid-123"
	companyID := 456

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		expectedPath := fmt.Sprintf("/api/v4/contacts/unsorted/%s/link", unsortedUID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	apiClient := client.NewClient(server.URL, "test_api_key")

	err := LinkContactWithCompany(context.Background(), apiClient, unsortedUID, companyID)

	if err != nil {
		t.Fatalf("Ошибка при связывании неразобранного контакта с компанией: %v", err)
	}
}
