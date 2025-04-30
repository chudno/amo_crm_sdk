package calls

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chudno/amo_crm_sdk/client"
)

func TestAddCall(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/calls"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{
			"_embedded": {
				"calls": [
					{
						"id": 123,
						"direction": "inbound",
						"status": "success",
						"responsible_user_id": 456,
						"created_by": 789,
						"updated_by": 789,
						"created_at": 1609459200,
						"updated_at": 1609459200,
						"account_id": 12345,
						"uniq": "call-uniq-123",
						"duration": 120,
						"source": "test_source",
						"call_result": "test_result",
						"phone": "+79001234567",
						"_links": {
							"self": {
								"href": "/api/v4/calls/123"
							}
						}
					}
				]
			}
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Создаем звонок для отправки
	call := &Call{
		Direction:         CallDirectionIncoming,
		Status:            CallStatusSuccess,
		ResponsibleUserID: 456,
		Duration:          120,
		Source:            "test_source",
		CallResult:        "test_result",
		Phone:             "+79001234567",
		CreatedAt:         time.Now().Unix(),
	}

	// Вызываем тестируемый метод
	createdCall, err := AddCall(apiClient, call)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при создании звонка: %v", err)
	}

	// Проверяем содержимое созданного звонка
	if createdCall.ID != 123 {
		t.Errorf("Ожидался ID 123, получен %d", createdCall.ID)
	}

	if createdCall.Direction != CallDirectionIncoming {
		t.Errorf("Ожидалось направление inbound, получено %s", createdCall.Direction)
	}

	if createdCall.Status != CallStatusSuccess {
		t.Errorf("Ожидался статус success, получен %s", createdCall.Status)
	}

	if createdCall.Phone != "+79001234567" {
		t.Errorf("Ожидался телефон +79001234567, получен %s", createdCall.Phone)
	}

	if createdCall.Duration != 120 {
		t.Errorf("Ожидалась продолжительность 120, получена %d", createdCall.Duration)
	}
}

func TestGetCalls(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/calls"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Проверяем параметры запроса
		expectedPage := "1"
		if r.URL.Query().Get("page") != expectedPage {
			t.Errorf("Ожидался параметр page=%s, получен %s", expectedPage, r.URL.Query().Get("page"))
		}

		expectedLimit := "50"
		if r.URL.Query().Get("limit") != expectedLimit {
			t.Errorf("Ожидался параметр limit=%s, получен %s", expectedLimit, r.URL.Query().Get("limit"))
		}

		// Проверяем фильтр
		expectedFilterDirection := "inbound"
		if r.URL.Query().Get("filter[direction]") != expectedFilterDirection {
			t.Errorf("Ожидался параметр filter[direction]=%s, получен %s", expectedFilterDirection, r.URL.Query().Get("filter[direction]"))
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"page": 1,
			"per_page": 50,
			"total": 2,
			"_embedded": {
				"calls": [
					{
						"id": 123,
						"direction": "inbound",
						"status": "success",
						"responsible_user_id": 456,
						"created_by": 789,
						"updated_by": 789,
						"created_at": 1609459200,
						"updated_at": 1609459200,
						"account_id": 12345,
						"uniq": "call-uniq-123",
						"duration": 120,
						"source": "test_source",
						"call_result": "test_result",
						"phone": "+79001234567",
						"_links": {
							"self": {
								"href": "/api/v4/calls/123"
							}
						}
					},
					{
						"id": 456,
						"direction": "inbound",
						"status": "missed",
						"responsible_user_id": 456,
						"created_by": 789,
						"updated_by": 789,
						"created_at": 1609459200,
						"updated_at": 1609459200,
						"account_id": 12345,
						"uniq": "call-uniq-456",
						"duration": 0,
						"source": "test_source",
						"call_result": "test_result",
						"phone": "+79001234568",
						"_links": {
							"self": {
								"href": "/api/v4/calls/456"
							}
						}
					}
				]
			}
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Создаем фильтр
	filter := map[string]string{
		"filter[direction]": string(CallDirectionIncoming),
	}

	// Вызываем тестируемый метод
	calls, err := GetCalls(apiClient, 1, 50, filter)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при получении звонков: %v", err)
	}

	if len(calls) != 2 {
		t.Fatalf("Ожидалось получение 2 звонков, получено %d", len(calls))
	}

	// Проверяем содержимое первого звонка
	if calls[0].ID != 123 {
		t.Errorf("Ожидался ID 123, получен %d", calls[0].ID)
	}

	if calls[0].Direction != CallDirectionIncoming {
		t.Errorf("Ожидалось направление inbound, получено %s", calls[0].Direction)
	}

	if calls[0].Status != CallStatusSuccess {
		t.Errorf("Ожидался статус success, получен %s", calls[0].Status)
	}

	// Проверяем содержимое второго звонка
	if calls[1].ID != 456 {
		t.Errorf("Ожидался ID 456, получен %d", calls[1].ID)
	}

	if calls[1].Status != CallStatusMissed {
		t.Errorf("Ожидался статус missed, получен %s", calls[1].Status)
	}
}

func TestGetCall(t *testing.T) {
	// ID звонка для теста
	callID := 123

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/calls/%d", callID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Проверяем параметр with
		expectedWith := "tags"
		if r.URL.Query().Get("with") != expectedWith {
			t.Errorf("Ожидался параметр with=%s, получен %s", expectedWith, r.URL.Query().Get("with"))
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 123,
			"direction": "inbound",
			"status": "success",
			"responsible_user_id": 456,
			"created_by": 789,
			"updated_by": 789,
			"created_at": 1609459200,
			"updated_at": 1609459200,
			"account_id": 12345,
			"uniq": "call-uniq-123",
			"duration": 120,
			"source": "test_source",
			"call_result": "test_result",
			"phone": "+79001234567",
			"_embedded": {
				"tags": [
					{
						"id": 101,
						"name": "Важно",
						"color": "#ff0000"
					},
					{
						"id": 102,
						"name": "Перезвонить",
						"color": "#00ff00"
					}
				]
			},
			"_links": {
				"self": {
					"href": "/api/v4/calls/123"
				}
			}
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод с параметром WithTags
	call, err := GetCall(apiClient, callID, WithTags)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при получении звонка: %v", err)
	}

	// Проверяем содержимое звонка
	if call.ID != callID {
		t.Errorf("Ожидался ID %d, получен %d", callID, call.ID)
	}

	if call.Direction != CallDirectionIncoming {
		t.Errorf("Ожидалось направление inbound, получено %s", call.Direction)
	}

	if call.Status != CallStatusSuccess {
		t.Errorf("Ожидался статус success, получен %s", call.Status)
	}

	if call.Phone != "+79001234567" {
		t.Errorf("Ожидался телефон +79001234567, получен %s", call.Phone)
	}

	// Проверяем теги
	if call.Embedded == nil {
		t.Fatalf("Отсутствует секция _embedded")
	}

	if len(call.Embedded.Tags) != 2 {
		t.Fatalf("Ожидалось 2 тега, получено %d", len(call.Embedded.Tags))
	}

	if call.Embedded.Tags[0].Name != "Важно" {
		t.Errorf("Ожидалось имя тега 'Важно', получено '%s'", call.Embedded.Tags[0].Name)
	}

	if call.Embedded.Tags[1].Name != "Перезвонить" {
		t.Errorf("Ожидалось имя тега 'Перезвонить', получено '%s'", call.Embedded.Tags[1].Name)
	}
}

func TestUpdateCall(t *testing.T) {
	// ID звонка для теста
	callID := 123

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "PATCH" {
			t.Errorf("Ожидался метод PATCH, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/calls/%d", callID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 123,
			"direction": "inbound",
			"status": "success",
			"responsible_user_id": 456,
			"created_by": 789,
			"updated_by": 789,
			"created_at": 1609459200,
			"updated_at": 1609545600,
			"account_id": 12345,
			"uniq": "call-uniq-123",
			"duration": 180,
			"source": "test_source_updated",
			"call_result": "test_result_updated",
			"phone": "+79001234567",
			"_links": {
				"self": {
					"href": "/api/v4/calls/123"
				}
			}
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Создаем звонок для обновления
	call := &Call{
		ID:         callID,
		Duration:   180,
		Source:     "test_source_updated",
		CallResult: "test_result_updated",
	}

	// Вызываем тестируемый метод
	updatedCall, err := UpdateCall(apiClient, call)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при обновлении звонка: %v", err)
	}

	// Проверяем содержимое обновленного звонка
	if updatedCall.ID != callID {
		t.Errorf("Ожидался ID %d, получен %d", callID, updatedCall.ID)
	}

	if updatedCall.Duration != 180 {
		t.Errorf("Ожидалась продолжительность 180, получена %d", updatedCall.Duration)
	}

	if updatedCall.Source != "test_source_updated" {
		t.Errorf("Ожидался источник test_source_updated, получен %s", updatedCall.Source)
	}

	if updatedCall.CallResult != "test_result_updated" {
		t.Errorf("Ожидался результат test_result_updated, получен %s", updatedCall.CallResult)
	}
}

func TestDeleteCall(t *testing.T) {
	// ID звонка для теста
	callID := 123

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "DELETE" {
			t.Errorf("Ожидался метод DELETE, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/calls/%d", callID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	err := DeleteCall(apiClient, callID)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при удалении звонка: %v", err)
	}
}

func TestLinkCallWithEntity(t *testing.T) {
	// ID звонка для теста
	callID := 123
	entityType := EntityTypeLead
	entityID := 456

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/calls/%d/link", callID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	err := LinkCallWithEntity(apiClient, callID, entityType, entityID)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при связывании звонка с сущностью: %v", err)
	}
}

func TestUnlinkCallFromEntity(t *testing.T) {
	// ID звонка для теста
	callID := 123
	entityType := EntityTypeLead
	entityID := 456

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/calls/%d/unlink", callID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	err := UnlinkCallFromEntity(apiClient, callID, entityType, entityID)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при отвязывании звонка от сущности: %v", err)
	}
}
