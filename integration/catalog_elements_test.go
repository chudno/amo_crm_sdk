//go:build integration

package integration

import (
	"testing"

	"github.com/chudno/amo_crm_sdk/entities/catalog_elements"
	"github.com/chudno/amo_crm_sdk/entities/catalogs"
)

func TestIntegration_CatalogElementsCRUD(t *testing.T) {
	apiClient, ctx := setupClient(t)

	// Создаём каталог для тестов
	catalog, err := catalogs.Create(ctx, apiClient, &catalogs.Catalog{
		Name: "Каталог для элементов (integration test)",
		Type: "regular",
	})
	if err != nil {
		t.Skipf("CreateCatalog: %v (возможно, достигнут лимит каталогов в аккаунте)", err)
	}
	if catalog.ID == 0 {
		t.Fatal("ID каталога не должен быть 0")
	}
	t.Logf("Создан каталог: ID=%d", catalog.ID)

	// CREATE element
	elem, err := catalog_elements.Create(ctx, apiClient, catalog.ID, &catalog_elements.Element{
		Name: "Тестовый элемент",
	})
	if err != nil {
		t.Fatalf("CreateElement: %v", err)
	}
	if elem.ID == 0 {
		t.Fatal("ID элемента не должен быть 0")
	}
	t.Logf("Создан элемент: ID=%d, Name=%q", elem.ID, elem.Name)

	// GET element
	got, err := catalog_elements.Get(ctx, apiClient, catalog.ID, elem.ID)
	if err != nil {
		t.Fatalf("GetElement(%d): %v", elem.ID, err)
	}
	t.Logf("Получен элемент: ID=%d, Name=%q", got.ID, got.Name)

	// UPDATE element
	got.Name = "Обновлённый элемент"
	updated, err := catalog_elements.Update(ctx, apiClient, catalog.ID, got)
	if err != nil {
		t.Fatalf("UpdateElement: %v", err)
	}
	t.Logf("Обновлён элемент: ID=%d, Name=%q", updated.ID, updated.Name)

	// CREATE BATCH
	batchElements, err := catalog_elements.CreateBatch(ctx, apiClient, catalog.ID, []catalog_elements.Element{
		{Name: "Батч элемент 1"},
		{Name: "Батч элемент 2"},
	})
	if err != nil {
		t.Fatalf("CreateBatch: %v", err)
	}
	if len(batchElements) < 2 {
		t.Errorf("Ожидалось минимум 2 элемента, получено %d", len(batchElements))
	}
	t.Logf("Создано батчем %d элементов", len(batchElements))

	// LIST elements
	list, err := catalog_elements.List(ctx, apiClient, catalog.ID, 1, 50, nil)
	if err != nil {
		t.Fatalf("ListElements: %v", err)
	}
	if len(list) < 3 {
		t.Errorf("Ожидалось минимум 3 элемента в списке, получено %d", len(list))
	}
	t.Logf("ListElements вернул %d элементов", len(list))

	// UPDATE BATCH
	if len(batchElements) >= 2 {
		batchElements[0].Name = "Батч обновлён 1"
		batchElements[1].Name = "Батч обновлён 2"
		updatedBatch, err := catalog_elements.UpdateBatch(ctx, apiClient, catalog.ID, batchElements)
		if err != nil {
			t.Fatalf("UpdateBatch: %v", err)
		}
		t.Logf("Обновлено батчем %d элементов", len(updatedBatch))
	}

	// DELETE element
	if err := catalog_elements.Delete(ctx, apiClient, catalog.ID, elem.ID); err != nil {
		t.Logf("DeleteElement(%d): %v (может не поддерживаться)", elem.ID, err)
	} else {
		t.Logf("Удалён элемент: ID=%d", elem.ID)
	}

	// DELETE BATCH
	if len(batchElements) >= 2 {
		batchIDs := []int{batchElements[0].ID, batchElements[1].ID}
		if err := catalog_elements.DeleteBatch(ctx, apiClient, catalog.ID, batchIDs); err != nil {
			t.Logf("DeleteBatch: %v (может не поддерживаться)", err)
		} else {
			t.Logf("Удалено батчем %d элементов", len(batchIDs))
		}
	}
}
