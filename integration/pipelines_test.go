//go:build integration

package integration

import (
	"testing"

	"github.com/chudno/amo_crm_sdk/entities/pipelines"
)

func TestIntegration_ListPipelines(t *testing.T) {
	apiClient, ctx := setupClient(t)

	list, err := pipelines.List(ctx, apiClient)
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	if len(list) == 0 {
		t.Fatal("Список воронок не должен быть пустым (в аккаунте всегда есть хотя бы одна воронка)")
	}

	t.Logf("Получено %d воронок", len(list))
	for _, p := range list {
		t.Logf("  ID=%d, Name=%q", p.ID, p.Name)
	}
}

func TestIntegration_GetPipeline(t *testing.T) {
	apiClient, ctx := setupClient(t)

	// Получаем список, чтобы узнать валидный ID
	list, err := pipelines.List(ctx, apiClient)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) == 0 {
		t.Fatal("Нет воронок для тестирования")
	}

	pipeline, err := pipelines.Get(ctx, apiClient, list[0].ID)
	if err != nil {
		t.Fatalf("GetPipeline(%d): %v", list[0].ID, err)
	}

	if pipeline.ID != list[0].ID {
		t.Errorf("Ожидался ID=%d, получен ID=%d", list[0].ID, pipeline.ID)
	}

	t.Logf("Воронка: ID=%d, Name=%q", pipeline.ID, pipeline.Name)
}

func TestIntegration_GetStatus(t *testing.T) {
	apiClient, ctx := setupClient(t)

	// Получаем первую воронку
	pipelinesList, err := pipelines.List(ctx, apiClient)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(pipelinesList) == 0 {
		t.Fatal("Нет воронок для тестирования")
	}

	p := pipelinesList[0]
	if len(p.Statuses) == 0 {
		t.Skip("У воронки нет статусов")
	}

	status, err := pipelines.GetStatus(ctx, apiClient, p.ID, p.Statuses[0].ID)
	if err != nil {
		t.Fatalf("GetStatus(%d, %d): %v", p.ID, p.Statuses[0].ID, err)
	}

	t.Logf("Статус: ID=%d, Name=%q", status.ID, status.Name)
}

func TestIntegration_PipelinesCRUD(t *testing.T) {
	apiClient, ctx := setupClient(t)

	// CREATE
	created, err := pipelines.Create(ctx, apiClient, &pipelines.Pipeline{
		Name:         "Тестовая воронка (integration test)",
		Sort:         100,
		IsMain:       false,
		IsUnsortedOn: true,
		Statuses: []pipelines.Status{
			{Name: "Новый", Sort: 10, Color: "#fffeb2"},
		},
	})
	if err != nil {
		t.Fatalf("CreatePipeline: %v", err)
	}
	if created.ID == 0 {
		t.Fatal("ID воронки не должен быть 0")
	}
	t.Logf("Создана воронка: ID=%d, Name=%q", created.ID, created.Name)
	// Гарантируем удаление воронки при завершении теста
	t.Cleanup(func() {
		_ = pipelines.Delete(ctx, apiClient, created.ID)
	})

	// UPDATE
	created.Name = "Обновлённая воронка (integration test)"
	updated, err := pipelines.Update(ctx, apiClient, created)
	if err != nil {
		t.Fatalf("UpdatePipeline: %v", err)
	}
	t.Logf("Обновлена воронка: ID=%d, Name=%q", updated.ID, updated.Name)

	// CREATE STATUS
	status, err := pipelines.CreateStatus(ctx, apiClient, created.ID, &pipelines.Status{
		Name:  "Тестовый статус",
		Sort:  10,
		Color: "#fffeb2",
	})
	if err != nil {
		t.Logf("CreateStatus: %v (может не поддерживаться)", err)
	} else if status.ID == 0 {
		t.Logf("CreateStatus вернул ID=0")
	} else {
		t.Logf("Создан статус: ID=%d, Name=%q", status.ID, status.Name)
	}

	// DELETE pipeline
	err = pipelines.Delete(ctx, apiClient, created.ID)
	if err != nil {
		t.Fatalf("DeletePipeline(%d): %v", created.ID, err)
	}
	t.Logf("Удалена воронка: ID=%d", created.ID)
}
