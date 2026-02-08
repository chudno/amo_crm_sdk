//go:build integration

package integration

import (
	"testing"

	"github.com/chudno/amo_crm_sdk/entities/companies"
)

func TestIntegration_CompaniesCRUD(t *testing.T) {
	apiClient, ctx := setupClient(t)

	// CREATE
	newCompany := &companies.Company{
		Name: "Тестовая компания (integration test)",
	}

	created, err := companies.Create(ctx, apiClient, newCompany)
	if err != nil {
		t.Fatalf("CreateCompany: %v", err)
	}
	if created.ID == 0 {
		t.Fatal("ID компании не должен быть 0")
	}
	t.Logf("Создана компания: ID=%d, Name=%q", created.ID, created.Name)

	// READ
	got, err := companies.Get(ctx, apiClient, created.ID)
	if err != nil {
		t.Fatalf("GetCompany(%d): %v", created.ID, err)
	}
	if got.Name != newCompany.Name {
		t.Errorf("Ожидалось имя %q, получено %q", newCompany.Name, got.Name)
	}

	// UPDATE
	got.Name = "Обновлённая компания (integration test)"
	updated, err := companies.Update(ctx, apiClient, got)
	if err != nil {
		t.Fatalf("UpdateCompany: %v", err)
	}
	t.Logf("Обновлена компания: ID=%d, Name=%q", updated.ID, updated.Name)

	// LIST
	list, err := companies.List(ctx, apiClient, 1, 5)
	if err != nil {
		t.Fatalf("GetCompanies: %v", err)
	}
	t.Logf("GetCompanies вернул %d компаний", len(list))
}
