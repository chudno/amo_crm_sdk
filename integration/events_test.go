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

	// GET — получение конкретного события по ID
	if len(list) > 0 {
		event, err := events.Get(ctx, apiClient, list[0].ID)
		if err != nil {
			t.Fatalf("GetEvent(%s): %v", list[0].ID, err)
		}
		if event.ID != list[0].ID {
			t.Errorf("Ожидался ID=%q, получен ID=%q", list[0].ID, event.ID)
		}
		t.Logf("Получено событие: ID=%s, Type=%q", event.ID, event.Type)
	}

	// LIST с фильтром по типу entity_create
	filteredList, err := events.List(ctx, apiClient,
		events.WithFilter(map[string]string{"filter[type]": string(events.EventTypeEntityCreate)}),
		events.WithLimit(10),
		events.WithPage(1),
	)
	if err != nil {
		t.Fatalf("GetEvents с фильтром: %v", err)
	}
	t.Logf("Получено %d событий с фильтром type=entity_create", len(filteredList))
	for _, e := range filteredList {
		t.Logf("  ID=%q, Type=%q, EntityType=%q", e.ID, e.Type, e.EntityType)
	}
}
