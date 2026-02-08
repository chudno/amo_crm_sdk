//go:build integration

package integration

import (
	"testing"

	"github.com/chudno/amo_crm_sdk/entities/leads"
)

func TestIntegration_LeadsCRUD(t *testing.T) {
	apiClient, ctx := setupClient(t)

	// CREATE
	newLead := &leads.Lead{
		Name: "Тестовый лид (integration test)",
	}

	created, err := leads.Create(ctx, apiClient, newLead)
	if err != nil {
		t.Fatalf("CreateLead: %v", err)
	}
	if created.ID == 0 {
		t.Fatal("ID созданного лида не должен быть 0")
	}
	t.Logf("Создан лид: ID=%d, Name=%q", created.ID, created.Name)

	// READ
	got, err := leads.Get(ctx, apiClient, created.ID)
	if err != nil {
		t.Fatalf("GetLead(%d): %v", created.ID, err)
	}
	if got.Name != newLead.Name {
		t.Errorf("Ожидалось имя %q, получено %q", newLead.Name, got.Name)
	}

	// UPDATE
	got.Name = "Обновлённый лид (integration test)"
	updated, err := leads.Update(ctx, apiClient, got)
	if err != nil {
		t.Fatalf("UpdateLead: %v", err)
	}
	t.Logf("Обновлён лид: ID=%d, Name=%q", updated.ID, updated.Name)

	// LIST
	list, err := leads.List(ctx, apiClient, 1, 5, nil)
	if err != nil {
		t.Fatalf("GetLeads: %v", err)
	}
	t.Logf("GetLeads вернул %d лидов", len(list))

	// DELETE
	if err := leads.Delete(ctx, apiClient, created.ID); err != nil {
		t.Fatalf("Delete(%d): %v", created.ID, err)
	}
	t.Logf("Удалён лид: ID=%d", created.ID)
}
