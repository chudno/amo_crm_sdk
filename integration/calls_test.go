//go:build integration

package integration

import (
	"testing"

	"github.com/chudno/amo_crm_sdk/entities/calls"
)

func TestIntegration_CallsCRUD(t *testing.T) {
	apiClient, ctx := setupClient(t)

	// ADD call — единственная операция, поддерживаемая API amoCRM для звонков
	newCall := &calls.Call{
		Direction: calls.CallDirectionOutgoing,
		Status:    calls.CallStatusSuccess,
		Phone:     "+79991234567",
		Duration:  120,
		Source:    "integration-test",
	}

	created, err := calls.Add(ctx, apiClient, newCall)
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if created.ID == 0 {
		t.Fatal("ID звонка не должен быть 0")
	}
	t.Logf("Создан звонок: ID=%d", created.ID)
}
