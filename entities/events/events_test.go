package events

import (
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

func TestGetEvents(t *testing.T) {
	t.Run("GetEventsWithFiltersAndPagination", func(t *testing.T) {
		// Создаем тестовый сервер
		server := setupGetEventsTestServer(t, true) // true для проверки параметров запроса
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Подготавливаем фильтры и опции
		filter := map[string]string{
			"filter[type]":        string(EventTypeNote),
			"filter[entity_type]": string(EventEntityTypeLead),
		}

		// Вызываем тестируемый метод
		events, err := GetEvents(apiClient,
			WithFilter(filter),
			WithPage(2),
			WithLimit(30),
			WithOrder("created_at", "desc"),
		)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при получении событий: %v", err)
		}

		// Проверяем содержимое событий
		verifyEventsList(t, events)
	})

	t.Run("GetEventsBasic", func(t *testing.T) {
		// Создаем тестовый сервер без проверки параметров
		server := setupGetEventsTestServer(t, false)
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемый метод с базовыми параметрами
		events, err := GetEvents(apiClient, WithPage(1), WithLimit(10))

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при получении событий: %v", err)
		}

		// Проверяем содержимое событий
		verifyEventsList(t, events)
	})
}

func TestGetEvent(t *testing.T) {
	t.Run("GetEventWithEntity", func(t *testing.T) {
		// ID события для теста
		eventID := 123

		// Создаем тестовый сервер
		server := setupGetEventTestServer(t, eventID, true) // true означает с параметром with=entity
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемый метод с параметром WithEntity
		event, err := GetEvent(apiClient, eventID, WithEntity())

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при получении события: %v", err)
		}

		// Проверяем содержимое события
		verifyEventDetails(t, event, eventID, true)
	})

	t.Run("GetEventWithoutEntity", func(t *testing.T) {
		// ID события для теста
		eventID := 123

		// Создаем тестовый сервер
		server := setupGetEventTestServer(t, eventID, false) // false означает без параметра with
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемый метод без параметра WithEntity
		event, err := GetEvent(apiClient, eventID)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при получении события: %v", err)
		}

		// Проверяем содержимое события
		verifyEventDetails(t, event, eventID, false)
	})
}
