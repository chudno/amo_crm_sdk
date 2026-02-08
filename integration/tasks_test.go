//go:build integration

package integration

import (
	"testing"
	"time"

	"github.com/chudno/amo_crm_sdk/entities/tasks"
)

func TestIntegration_TasksCRUD(t *testing.T) {
	apiClient, ctx := setupClient(t)

	// CREATE
	newTask := &tasks.Task{
		Text:         "Тестовая задача (integration test)",
		CompleteTill: time.Now().Add(24 * time.Hour).Unix(),
		TaskTypeID:   1, // Звонок
	}

	created, err := tasks.Create(ctx, apiClient, newTask)
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}
	if created.ID == 0 {
		t.Fatal("ID задачи не должен быть 0")
	}
	t.Logf("Создана задача: ID=%d, Text=%q", created.ID, created.Text)

	// READ
	got, err := tasks.Get(ctx, apiClient, created.ID)
	if err != nil {
		t.Fatalf("GetTask(%d): %v", created.ID, err)
	}
	t.Logf("Получена задача: ID=%d, Text=%q", got.ID, got.Text)

	// UPDATE
	got.Text = "Обновлённая задача (integration test)"
	updated, err := tasks.Update(ctx, apiClient, got)
	if err != nil {
		t.Fatalf("UpdateTask: %v", err)
	}
	t.Logf("Обновлена задача: ID=%d, Text=%q", updated.ID, updated.Text)

	// LIST
	list, err := tasks.List(ctx, apiClient, 5, 1, nil)
	if err != nil {
		t.Fatalf("ListTasks: %v", err)
	}
	t.Logf("ListTasks вернул %d задач", len(list))

	// DELETE
	if err := tasks.Delete(ctx, apiClient, created.ID); err != nil {
		t.Fatalf("Delete(%d): %v", created.ID, err)
	}
	t.Logf("Удалена задача: ID=%d", created.ID)
}
