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
