package calls

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chudno/amo_crm_sdk/client"
)

func TestAddCall(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		expectedPath := "/api/v4/calls"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Ошибка при чтении тела запроса: %v", err)
		}

		var calls []map[string]any
		if err := json.Unmarshal(body, &calls); err != nil {
			t.Fatalf("Ошибка при разборе JSON: %v", err)
		}
		if len(calls) != 1 {
			t.Fatalf("Ожидался 1 элемент в массиве, получено %d", len(calls))
		}

		callStatusCode, ok := calls[0]["call_status"].(float64)
		if !ok || int(callStatusCode) != CallStatusCodeSuccess {
			t.Errorf("Ожидался call_status=%d, получено %v", CallStatusCodeSuccess, calls[0]["call_status"])
		}

		if _, hasStatus := calls[0]["status"]; hasStatus {
			t.Errorf("Поле 'status' не должно отправляться в API")
		}

		if _, hasUniq := calls[0]["uniq"]; !hasUniq {
			t.Errorf("Поле 'uniq' должно быть в запросе")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{
			"_embedded": {
				"calls": [
					{
						"id": 123,
						"direction": "inbound",
						"call_status": 4,
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

	apiClient := client.NewClient(server.URL, "test_api_key")

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

	createdCall, err := Add(context.Background(), apiClient, call)

	if err != nil {
		t.Fatalf("Ошибка при создании звонка: %v", err)
	}

	if createdCall.ID != 123 {
		t.Errorf("Ожидался ID 123, получен %d", createdCall.ID)
	}

	if createdCall.Direction != CallDirectionIncoming {
		t.Errorf("Ожидалось направление inbound, получено %s", createdCall.Direction)
	}

	if createdCall.CallStatusCode != CallStatusCodeSuccess {
		t.Errorf("Ожидался call_status=%d, получен %d", CallStatusCodeSuccess, createdCall.CallStatusCode)
	}

	if createdCall.Phone != "+79001234567" {
		t.Errorf("Ожидался телефон +79001234567, получен %s", createdCall.Phone)
	}

	if createdCall.Duration != 120 {
		t.Errorf("Ожидалась продолжительность 120, получена %d", createdCall.Duration)
	}
}

func TestAddCall_ValidationErrors(t *testing.T) {
	apiClient := client.NewClient("https://example.amocrm.ru", "test_api_key")

	// Без direction
	_, err := Add(context.Background(), apiClient, &Call{Phone: "+79001234567"})
	if err == nil {
		t.Error("Ожидалась ошибка при отсутствии direction")
	}

	// Без phone
	_, err = Add(context.Background(), apiClient, &Call{Direction: CallDirectionIncoming})
	if err == nil {
		t.Error("Ожидалась ошибка при отсутствии phone")
	}
}
