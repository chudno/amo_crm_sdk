package events

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// setupGetEventsTestServer создает тестовый сервер для получения списка событий
func setupGetEventsTestServer(t *testing.T, withQueryParams bool) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/events"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Проверяем параметры запроса
		if withQueryParams {
			expectedPage := "2"
			if r.URL.Query().Get("page") != expectedPage {
				t.Errorf("Ожидался параметр page=%s, получен %s", expectedPage, r.URL.Query().Get("page"))
			}

			expectedLimit := "30"
			if r.URL.Query().Get("limit") != expectedLimit {
				t.Errorf("Ожидался параметр limit=%s, получен %s", expectedLimit, r.URL.Query().Get("limit"))
			}

			// Проверяем фильтры
			expectedFilterType := "note"
			if r.URL.Query().Get("filter[type]") != expectedFilterType {
				t.Errorf("Ожидался параметр filter[type]=%s, получен %s", expectedFilterType, r.URL.Query().Get("filter[type]"))
			}

			expectedFilterEntityType := "lead"
			if r.URL.Query().Get("filter[entity_type]") != expectedFilterEntityType {
				t.Errorf("Ожидался параметр filter[entity_type]=%s, получен %s", expectedFilterEntityType, r.URL.Query().Get("filter[entity_type]"))
			}

			// Проверяем сортировку
			expectedOrderCreatedAt := "desc"
			if r.URL.Query().Get("order[created_at]") != expectedOrderCreatedAt {
				t.Errorf("Ожидался параметр order[created_at]=%s, получен %s", expectedOrderCreatedAt, r.URL.Query().Get("order[created_at]"))
			}
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"page": 2,
			"per_page": 30,
			"total": 150,
			"order": [
				{
					"field": "created_at",
					"order": "desc"
				}
			],
			"_embedded": {
				"events": [
					{
						"id": 123,
						"type": "note",
						"entity_id": 456,
						"entity_type": "lead",
						"created_by": 789,
						"account_id": 12345,
						"created_at": 1609459200,
						"value_after": {"text": "Это тестовое примечание"},
						"value_after_pretty": "Это тестовое примечание",
						"_links": {
							"self": {
								"href": "/api/v4/events/123"
							}
						}
					},
					{
						"id": 124,
						"type": "note",
						"entity_id": 457,
						"entity_type": "lead",
						"created_by": 789,
						"account_id": 12345,
						"created_at": 1609458200,
						"value_after": {"text": "Еще одно тестовое примечание"},
						"value_after_pretty": "Еще одно тестовое примечание",
						"_links": {
							"self": {
								"href": "/api/v4/events/124"
							}
						}
					}
				]
			},
			"_next_page": "/api/v4/events?page=3&limit=30",
			"_prev_page": "/api/v4/events?page=1&limit=30",
			"_total_path": "/api/v4/events/total"
		}`))
	}))
}

// verifyEventsList проверяет результат получения списка событий
func verifyEventsList(t *testing.T, events []Event) {
	if len(events) != 2 {
		t.Fatalf("Ожидалось 2 события, получено %d", len(events))
	}

	// Проверяем содержимое первого события
	if events[0].ID != 123 {
		t.Errorf("Ожидался ID 123, получен %d", events[0].ID)
	}

	if events[0].Type != EventTypeNote {
		t.Errorf("Ожидался тип события note, получен %s", events[0].Type)
	}

	if events[0].EntityID != 456 {
		t.Errorf("Ожидался ID сущности 456, получен %d", events[0].EntityID)
	}

	if events[0].EntityType != EventEntityTypeLead {
		t.Errorf("Ожидался тип сущности lead, получен %s", events[0].EntityType)
	}

	if events[0].ValueAfterPretty != "Это тестовое примечание" {
		t.Errorf("Ожидалось примечание 'Это тестовое примечание', получено '%s'", events[0].ValueAfterPretty)
	}

	// Проверяем содержимое второго события
	if events[1].ID != 124 {
		t.Errorf("Ожидался ID 124, получен %d", events[1].ID)
	}

	if events[1].ValueAfterPretty != "Еще одно тестовое примечание" {
		t.Errorf("Ожидалось примечание 'Еще одно тестовое примечание', получено '%s'", events[1].ValueAfterPretty)
	}
}

// setupGetEventTestServer создает тестовый сервер для получения информации о событии
func setupGetEventTestServer(t *testing.T, eventID int, withEntity bool) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/events/%d", eventID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Проверяем параметр with
		if withEntity {
			expectedWith := "entity"
			if r.URL.Query().Get("with") != expectedWith {
				t.Errorf("Ожидался параметр with=%s, получен %s", expectedWith, r.URL.Query().Get("with"))
			}
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		var response string
		if withEntity {
			response = fmt.Sprintf(`{
				"id": %d,
				"type": "note",
				"entity_id": 456,
				"entity_type": "lead",
				"created_by": 789,
				"account_id": 12345,
				"created_at": 1609459200,
				"value_after": {"text": "Это тестовое примечание"},
				"value_after_pretty": "Это тестовое примечание",
				"_embedded": {
					"entity": {
						"id": 456,
						"name": "Тестовая сделка",
						"created_at": 1609450000,
						"updated_at": 1609455000
					}
				},
				"_links": {
					"self": {
						"href": "/api/v4/events/%d"
					}
				}
			}`, eventID, eventID)
		} else {
			response = fmt.Sprintf(`{
				"id": %d,
				"type": "note",
				"entity_id": 456,
				"entity_type": "lead",
				"created_by": 789,
				"account_id": 12345,
				"created_at": 1609459200,
				"value_after": {"text": "Это тестовое примечание"},
				"value_after_pretty": "Это тестовое примечание",
				"_links": {
					"self": {
						"href": "/api/v4/events/%d"
					}
				}
			}`, eventID, eventID)
		}

		_, _ = w.Write([]byte(response))
	}))
}

// verifyEventDetails проверяет детали события
func verifyEventDetails(t *testing.T, event *Event, eventID int, withEntity bool) {
	if event.ID != eventID {
		t.Errorf("Ожидался ID %d, получен %d", eventID, event.ID)
	}

	if event.Type != EventTypeNote {
		t.Errorf("Ожидался тип события note, получен %s", event.Type)
	}

	if event.EntityID != 456 {
		t.Errorf("Ожидался ID сущности 456, получен %d", event.EntityID)
	}

	if event.EntityType != EventEntityTypeLead {
		t.Errorf("Ожидался тип сущности lead, получен %s", event.EntityType)
	}

	if event.ValueAfterPretty != "Это тестовое примечание" {
		t.Errorf("Ожидалось примечание 'Это тестовое примечание', получено '%s'", event.ValueAfterPretty)
	}

	if withEntity {
		// Проверяем вложенную сущность
		if event.Embedded == nil {
			t.Fatalf("Отсутствует секция _embedded")
		}

		if event.Embedded.Entity == nil {
			t.Fatalf("Отсутствует секция _embedded.entity")
		}

		if event.Embedded.Entity.ID != 456 {
			t.Errorf("Ожидался ID сущности 456, получен %d", event.Embedded.Entity.ID)
		}

		if event.Embedded.Entity.Name != "Тестовая сделка" {
			t.Errorf("Ожидалось имя сущности 'Тестовая сделка', получено '%s'", event.Embedded.Entity.Name)
		}
	}
}
