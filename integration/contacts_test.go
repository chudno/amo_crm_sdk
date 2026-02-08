//go:build integration

package integration

import (
	"testing"

	"github.com/chudno/amo_crm_sdk/entities/contacts"
)

func TestIntegration_ContactsCRUD(t *testing.T) {
	apiClient, ctx := setupClient(t)

	// CREATE
	newContact := &contacts.Contact{
		FirstName: "Тест",
		LastName:  "Интеграция",
	}

	created, err := contacts.Create(ctx, apiClient, newContact)
	if err != nil {
		t.Fatalf("CreateContact: %v", err)
	}
	if created.ID == 0 {
		t.Fatal("ID контакта не должен быть 0")
	}
	t.Logf("Создан контакт: ID=%d, Name=%q", created.ID, created.FirstName)

	// READ
	got, err := contacts.Get(ctx, apiClient, created.ID)
	if err != nil {
		t.Fatalf("GetContact(%d): %v", created.ID, err)
	}
	t.Logf("Получен контакт: ID=%d, Name=%q", got.ID, got.FirstName)

	// LIST
	list, err := contacts.List(ctx, apiClient, 1, 5)
	if err != nil {
		t.Fatalf("GetContacts: %v", err)
	}
	t.Logf("GetContacts вернул %d контактов", len(list))
}
