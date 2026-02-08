//go:build integration

package integration

import (
	"testing"

	"github.com/chudno/amo_crm_sdk/entities/events"
)

func TestIntegration_GetEvents(t *testing.T) {
	apiClient, ctx := setupClient(t)

	list, err := events.List(ctx, apiClient, events.WithLimit(5), events.WithPage(1))
	if err != nil {
		t.Fatalf("GetEvents: %v", err)
	}

	t.Logf("Получено %d событий", len(list))
	for _, e := range list {
		t.Logf("  ID=%q, Type=%q, CreatedAt=%d", e.ID, e.Type, e.CreatedAt)
	}
}
