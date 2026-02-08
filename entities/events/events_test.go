package events

import (
	"context"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

func TestGetEvents(t *testing.T) {
	t.Run("GetEventsWithFiltersAndPagination", func(t *testing.T) {
		server := setupGetEventsTestServer(t, true) // true для проверки параметров запроса
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		filter := map[string]string{
			"filter[type]":        string(EventTypeNote),
			"filter[entity_type]": string(EventEntityTypeLead),
		}

		events, err := List(context.Background(), apiClient,
			WithFilter(filter),
			WithPage(2),
			WithLimit(30),
			WithOrder("created_at", "desc"),
		)

		if err != nil {
			t.Fatalf("Ошибка при получении событий: %v", err)
		}

		verifyEventsList(t, events)
	})

	t.Run("GetEventsBasic", func(t *testing.T) {
		server := setupGetEventsTestServer(t, false)
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		events, err := List(context.Background(), apiClient, WithPage(1), WithLimit(10))

		if err != nil {
			t.Fatalf("Ошибка при получении событий: %v", err)
		}

		verifyEventsList(t, events)
	})
}

func TestGetEvent(t *testing.T) {
	t.Run("GetEventWithEntity", func(t *testing.T) {
		eventID := "123"

		server := setupGetEventTestServer(t, eventID, true) // true означает с параметром with=entity
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		event, err := Get(context.Background(), apiClient, eventID, WithEntity())

		if err != nil {
			t.Fatalf("Ошибка при получении события: %v", err)
		}

		verifyEventDetails(t, event, eventID, true)
	})

	t.Run("GetEventWithoutEntity", func(t *testing.T) {
		eventID := "123"

		server := setupGetEventTestServer(t, eventID, false) // false означает без параметра with
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		event, err := Get(context.Background(), apiClient, eventID)

		if err != nil {
			t.Fatalf("Ошибка при получении события: %v", err)
		}

		verifyEventDetails(t, event, eventID, false)
	})
}
