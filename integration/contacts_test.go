//go:build integration

package integration

import (
	"testing"

	"github.com/chudno/amo_crm_sdk/entities/companies"
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

	// LINK WITH COMPANY
	// Создаём компанию для привязки к контакту
	newCompany := &companies.Company{
		Name: "Компания для привязки (integration test)",
	}
	company, err := companies.Create(ctx, apiClient, newCompany)
	if err != nil {
		t.Fatalf("CreateCompany: %v", err)
	}
	if company.ID == 0 {
		t.Fatal("ID компании не должен быть 0")
	}
	t.Logf("Создана компания: ID=%d, Name=%q", company.ID, company.Name)

	// Привязываем контакт к компании
	err = contacts.LinkWithCompany(ctx, apiClient, created.ID, company.ID)
	if err != nil {
		t.Fatalf("LinkWithCompany(contactID=%d, companyID=%d): %v", created.ID, company.ID, err)
	}
	t.Logf("Контакт ID=%d привязан к компании ID=%d", created.ID, company.ID)

	// Проверяем, что связь отображается при запросе контакта с WithCompanies
	contactWithCompanies, err := contacts.Get(ctx, apiClient, created.ID, contacts.WithCompanies)
	if err != nil {
		t.Fatalf("GetContact(%d, WithCompanies): %v", created.ID, err)
	}
	if contactWithCompanies.Embedded != nil {
		for _, c := range contactWithCompanies.Embedded.Companies {
			if c.ID == company.ID {
				t.Logf("Проверено: компания ID=%d присутствует в embedded контакта", company.ID)
			}
		}
	}
}
