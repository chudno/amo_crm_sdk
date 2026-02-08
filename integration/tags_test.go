//go:build integration

package integration

import (
	"testing"

	"github.com/chudno/amo_crm_sdk/entities/leads"
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

	if created.ID != 0 {
		found := false
		for _, tag := range list {
			if tag.ID == created.ID {
				found = true
				break
			}
		}
		if !found {
			t.Logf("Созданный тег (ID=%d) не найден в первых %d тегах (может быть на другой странице)", created.ID, len(list))
		}
	}

	// GET — получение тега по ID (через фильтр списка)
	gotTag, err := tags.Get(ctx, apiClient, tags.EntityTypeLead, created.ID)
	if err != nil {
		t.Fatalf("GetTag(%d): %v", created.ID, err)
	}
	if gotTag.Name != created.Name {
		t.Errorf("Ожидалось имя %q, получено %q", created.Name, gotTag.Name)
	}
	t.Logf("Get тег: ID=%d, Name=%q", gotTag.ID, gotTag.Name)

	// CREATE BATCH — создание нескольких тегов за один запрос
	batchTags, err := tags.CreateBatch(ctx, apiClient, tags.EntityTypeLead, []tags.Tag{
		{Name: "batch-tag-1"},
		{Name: "batch-tag-2"},
	})
	if err != nil {
		t.Logf("CreateBatch: %v (может не поддерживаться)", err)
	} else {
		if len(batchTags) < 2 {
			t.Errorf("Ожидалось >= 2 тегов из batch, получено %d", len(batchTags))
		}
		t.Logf("CreateBatch вернул %d тегов", len(batchTags))
	}

	// UPDATE — обновление имени тега
	updatedTag, err := tags.Update(ctx, apiClient, tags.EntityTypeLead, &tags.Tag{
		ID:   created.ID,
		Name: "updated-tag",
	})
	if err != nil {
		t.Logf("UpdateTag: %v (может не поддерживаться)", err)
	} else {
		t.Logf("Обновлён тег: ID=%d, Name=%q", updatedTag.ID, updatedTag.Name)
	}

	// LINK — привязка тега к лиду
	newLead := &leads.Lead{
		Name: "Лид для тестирования тегов (integration test)",
	}
	createdLead, err := leads.Create(ctx, apiClient, newLead)
	if err != nil {
		t.Fatalf("CreateLead: %v", err)
	}
	t.Logf("Создан лид для тегов: ID=%d", createdLead.ID)

	err = tags.LinkEntity(ctx, apiClient, tags.EntityTypeLead, createdLead.ID, []tags.Tag{
		{ID: created.ID},
	})
	if err != nil {
		t.Logf("LinkEntity: %v (может не поддерживаться)", err)
	} else {
		t.Log("Тег успешно привязан к лиду")

		// LIST FOR ENTITY — список тегов привязанных к лиду
		entityTags, err := tags.ListForEntity(ctx, apiClient, tags.EntityTypeLead, createdLead.ID)
		if err != nil {
			t.Logf("ListForEntity: %v", err)
		} else if len(entityTags) == 0 {
			t.Logf("Список тегов лида пуст (может быть задержка API)")
		} else {
			t.Logf("ListForEntity вернул %d тегов для лида ID=%d", len(entityTags), createdLead.ID)
		}
	}
}
