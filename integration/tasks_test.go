//go:build integration

package integration

import (
	"testing"
	"time"

	"github.com/chudno/amo_crm_sdk/entities/leads"
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

	// COMPLETE
	_, err = tasks.Complete(ctx, apiClient, created.ID, "Задача выполнена")
	if err != nil {
		t.Fatalf("CompleteTask(%d): %v", created.ID, err)
	}
	// API возвращает минимальные данные на PATCH — перечитываем задачу
	completed, err := tasks.Get(ctx, apiClient, created.ID)
	if err != nil {
		t.Fatalf("GetTask(%d) после Complete: %v", created.ID, err)
	}
	if !completed.IsCompleted {
		t.Errorf("Задача должна быть отмечена как выполненная (IsCompleted=true)")
	}
	t.Logf("Задача ID=%d отмечена как выполненная", completed.ID)

	// CREATE FOR ENTITY
	// Создаём лид, к которому привяжем задачу
	lead, err := leads.Create(ctx, apiClient, &leads.Lead{
		Name: "Лид для задачи (integration test)",
	})
	if err != nil {
		t.Fatalf("CreateLead: %v", err)
	}
	t.Cleanup(func() {
		_ = leads.Delete(ctx, apiClient, lead.ID)
	})
	t.Logf("Создан лид: ID=%d", lead.ID)

	entityTask, err := tasks.CreateForEntity(
		ctx, apiClient,
		tasks.EntityTypeLead, lead.ID,
		1, "Задача для лида (integration test)",
		time.Now().Add(48*time.Hour),
		0,
	)
	if err != nil {
		t.Fatalf("CreateForEntity(lead=%d): %v", lead.ID, err)
	}
	if entityTask.ID == 0 {
		t.Fatal("ID задачи для лида не должен быть 0")
	}
	t.Logf("Создана задача для лида: TaskID=%d, LeadID=%d", entityTask.ID, lead.ID)
}
