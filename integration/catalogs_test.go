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
		t.Skipf("CreateCatalog: %v (возможно, достигнут лимит каталогов в аккаунте)", err)
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

	// ADD CUSTOM FIELD
	field, err := catalogs.AddCustomField(ctx, apiClient, created.ID, &catalogs.CustomField{
		Name: "Тестовое поле",
		Type: "text",
	})
	if err != nil {
		t.Fatalf("AddCustomField: %v", err)
	}
	if field.ID == 0 {
		t.Fatal("ID пользовательского поля не должен быть 0")
	}
	t.Logf("Создано пользовательское поле: ID=%d, Name=%q", field.ID, field.Name)

	// LIST CUSTOM FIELDS
	fields, err := catalogs.ListCustomFields(ctx, apiClient, created.ID)
	if err != nil {
		t.Fatalf("ListCustomFields: %v", err)
	}
	if len(fields) == 0 {
		t.Fatal("Список пользовательских полей не должен быть пустым после добавления")
	}
	t.Logf("ListCustomFields вернул %d полей", len(fields))

	// GET CUSTOM FIELD
	gotField, err := catalogs.GetCustomField(ctx, apiClient, created.ID, field.ID)
	if err != nil {
		t.Fatalf("GetCustomField(%d): %v", field.ID, err)
	}
	if gotField.Name != field.Name {
		t.Errorf("Ожидалось имя поля %q, получено %q", field.Name, gotField.Name)
	}
	t.Logf("GetCustomField: ID=%d, Name=%q", gotField.ID, gotField.Name)

	// UPDATE CUSTOM FIELD
	updatedField, err := catalogs.UpdateCustomField(ctx, apiClient, created.ID, &catalogs.CustomField{
		ID:   field.ID,
		Name: "Обновлённое поле",
	})
	if err != nil {
		t.Fatalf("UpdateCustomField: %v", err)
	}
	t.Logf("Обновлено пользовательское поле: ID=%d, Name=%q", updatedField.ID, updatedField.Name)

	// DELETE CUSTOM FIELD
	err = catalogs.DeleteCustomField(ctx, apiClient, created.ID, field.ID)
	if err != nil {
		t.Fatalf("DeleteCustomField(%d): %v", field.ID, err)
	}
	t.Logf("Удалено пользовательское поле: ID=%d", field.ID)
}
