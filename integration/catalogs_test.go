//go:build integration

package integration

import (
	"testing"

	"github.com/chudno/amo_crm_sdk/entities/catalogs"
)

func TestIntegration_CatalogsCRUD(t *testing.T) {
	apiClient, ctx := setupClient(t)

	// CREATE
	newCatalog := &catalogs.Catalog{
		Name: "Тестовый каталог (integration test)",
		Type: "regular",
	}

	created, err := catalogs.Create(ctx, apiClient, newCatalog)
	if err != nil {
		t.Fatalf("CreateCatalog: %v", err)
	}
	if created.ID == 0 {
		t.Fatal("ID каталога не должен быть 0")
	}
	t.Logf("Создан каталог: ID=%d, Name=%q", created.ID, created.Name)

	// READ
	got, err := catalogs.Get(ctx, apiClient, created.ID)
	if err != nil {
		t.Fatalf("GetCatalog(%d): %v", created.ID, err)
	}
	t.Logf("Получен каталог: ID=%d, Name=%q", got.ID, got.Name)

	// UPDATE
	got.Name = "Обновлённый каталог (integration test)"
	updated, err := catalogs.Update(ctx, apiClient, got)
	if err != nil {
		t.Fatalf("UpdateCatalog: %v", err)
	}
	t.Logf("Обновлён каталог: ID=%d, Name=%q", updated.ID, updated.Name)

	// LIST
	list, err := catalogs.List(ctx, apiClient, 1, 50, nil)
	if err != nil {
		t.Fatalf("GetCatalogs: %v", err)
	}
	t.Logf("GetCatalogs вернул %d каталогов", len(list))

	// DELETE
	if err := catalogs.Delete(ctx, apiClient, created.ID); err != nil {
		t.Fatalf("Delete(%d): %v", created.ID, err)
	}
	t.Logf("Удалён каталог: ID=%d", created.ID)
}
