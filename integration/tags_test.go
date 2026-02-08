//go:build integration

package integration

import (
	"testing"

	"github.com/chudno/amo_crm_sdk/entities/tags"
)

func TestIntegration_TagsCRUD(t *testing.T) {
	apiClient, ctx := setupClient(t)

	// CREATE
	newTag := &tags.Tag{
		Name: "integration-test-tag",
	}

	created, err := tags.Create(ctx, apiClient, tags.EntityTypeLead, newTag)
	if err != nil {
		t.Fatalf("CreateTag: %v", err)
	}
	if created.ID == 0 {
		t.Fatal("ID тега не должен быть 0")
	}
	t.Logf("Создан тег: ID=%d, Name=%q", created.ID, created.Name)

	// LIST
	list, err := tags.List(ctx, apiClient, tags.EntityTypeLead, 1, 50)
	if err != nil {
		t.Fatalf("GetTags: %v", err)
	}
	t.Logf("GetTags вернул %d тегов", len(list))

	found := false
	for _, tag := range list {
		if tag.ID == created.ID {
			found = true
			break
		}
	}
	if !found {
		t.Error("Созданный тег не найден в списке")
	}

	// DELETE
	if err := tags.Delete(ctx, apiClient, tags.EntityTypeLead, created.ID); err != nil {
		t.Fatalf("Delete(%d): %v", created.ID, err)
	}
	t.Logf("Удалён тег: ID=%d", created.ID)
}
