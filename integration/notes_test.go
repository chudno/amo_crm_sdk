//go:build integration

package integration

import (
	"testing"

	"github.com/chudno/amo_crm_sdk/entities/leads"
	"github.com/chudno/amo_crm_sdk/entities/notes"
)

func TestIntegration_NotesCRUD(t *testing.T) {
	apiClient, ctx := setupClient(t)

	// Создаём лид, к которому привяжем примечание
	lead, err := leads.Create(ctx, apiClient, &leads.Lead{
		Name: "Лид для тестирования примечаний",
	})
	if err != nil {
		t.Fatalf("CreateLead: %v", err)
	}
	t.Cleanup(func() {
		_ = leads.Delete(ctx, apiClient, lead.ID)
	})

	// CREATE
	// NoteType 4 = обычное текстовое примечание (common)
	newNote := &notes.Note{
		NoteType: 4,
		Params: notes.NoteParams{
			Text: "Тестовое примечание (integration test)",
		},
	}

	created, err := notes.Create(ctx, apiClient, "leads", lead.ID, newNote)
	if err != nil {
		t.Fatalf("CreateNote: %v", err)
	}
	if created.ID == 0 {
		t.Fatal("ID примечания не должен быть 0")
	}
	t.Logf("Создано примечание: ID=%d", created.ID)

	// READ
	got, err := notes.Get(ctx, apiClient, "leads", lead.ID, created.ID)
	if err != nil {
		t.Fatalf("GetNote(%d): %v", created.ID, err)
	}
	t.Logf("Получено примечание: ID=%d, Type=%q", got.ID, got.NoteType)

	// LIST
	list, err := notes.List(ctx, apiClient, "leads", lead.ID, 10, 1)
	if err != nil {
		t.Fatalf("ListNotes: %v", err)
	}
	t.Logf("ListNotes вернул %d примечаний", len(list))

	// DELETE
	if err := notes.Delete(ctx, apiClient, "leads", lead.ID, created.ID); err != nil {
		t.Fatalf("Delete(%d): %v", created.ID, err)
	}
	t.Logf("Удалено примечание: ID=%d", created.ID)
}
