package segments

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

// TestGetSegmentErrors проверяет обработку ошибок при получении сегмента
func TestGetSegmentErrors(t *testing.T) {
	segmentID := 123

	// Тест на ошибку сервера
	t.Run("ServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")
		_, err := GetSegment(apiClient, segmentID)

		if err == nil {
			t.Fatal("Ожидалась ошибка HTTP статуса, но получен nil")
		}
	})

	// Тест на некорректный JSON
	t.Run("InvalidJSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{invalid_json`))
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")
		_, err := GetSegment(apiClient, segmentID)

		if err == nil {
			t.Fatal("Ожидалась ошибка разбора JSON, но получен nil")
		}
	})

	// Тест на отсутствие сегмента
	t.Run("NotFound", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")
		_, err := GetSegment(apiClient, segmentID)

		if err == nil {
			t.Fatal("Ожидалась ошибка HTTP статуса, но получен nil")
		}
	})
}

// TestGetSegmentsErrors проверяет обработку ошибок при получении списка сегментов
func TestGetSegmentsErrors(t *testing.T) {
	// Тест на ошибку сервера
	t.Run("ServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")
		_, err := GetSegments(apiClient, 1, 50)

		if err == nil {
			t.Fatal("Ожидалась ошибка HTTP статуса, но получен nil")
		}
	})

	// Тест на некорректный JSON
	t.Run("InvalidJSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{invalid_json`))
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")
		_, err := GetSegments(apiClient, 1, 50)

		if err == nil {
			t.Fatal("Ожидалась ошибка разбора JSON, но получен nil")
		}
	})

	// Тест на пустой список сегментов
	t.Run("EmptySegments", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"page": 1,
				"per_page": 50,
				"_embedded": {
					"segments": []
				},
				"_links": {
					"self": {
						"href": "/api/v4/segments"
					}
				}
			}`))
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")
		segments, err := GetSegments(apiClient, 1, 50)

		if err != nil {
			t.Fatalf("Неожиданная ошибка: %v", err)
		}
		if len(segments) != 0 {
			t.Fatalf("Ожидался пустой список сегментов, получено %d", len(segments))
		}
	})
}

// TestAddSegmentErrors проверяет обработку ошибок при создании сегмента
func TestAddSegmentErrors(t *testing.T) {
	segment := &Segment{
		Name: "Тестовый сегмент",
		Type: SegmentTypeDynamic,
		Filter: &Filter{
			Logic: "and",
			Nodes: []FilterNode{
				{
					FieldCode: "email",
					Operator:  "contains",
					Value:     "example.com",
				},
			},
		},
	}

	// Тест на ошибку сервера
	t.Run("ServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")
		_, err := AddSegment(apiClient, segment)

		if err == nil {
			t.Fatal("Ожидалась ошибка HTTP статуса, но получен nil")
		}
	})

	// Тест на некорректный JSON
	t.Run("InvalidJSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{invalid_json`))
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")
		_, err := AddSegment(apiClient, segment)

		if err == nil {
			t.Fatal("Ожидалась ошибка разбора JSON, но получен nil")
		}
	})

	// Тест на пустой список сегментов в ответе
	t.Run("EmptyResponse", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{
				"_embedded": {
					"segments": []
				}
			}`))
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")
		_, err := AddSegment(apiClient, segment)

		if err == nil {
			t.Fatal("Ожидалась ошибка пустого ответа, но получен nil")
		}
	})
}

// TestUpdateSegmentErrors проверяет обработку ошибок при обновлении сегмента
func TestUpdateSegmentErrors(t *testing.T) {
	segment := &Segment{
		ID:   123,
		Name: "Обновленный сегмент",
		Type: SegmentTypeDynamic,
	}

	// Тест на ошибку сервера
	t.Run("ServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")
		_, err := UpdateSegment(apiClient, segment)

		if err == nil {
			t.Fatal("Ожидалась ошибка HTTP статуса, но получен nil")
		}
	})

	// Тест на некорректный JSON
	t.Run("InvalidJSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{invalid_json`))
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")
		_, err := UpdateSegment(apiClient, segment)

		if err == nil {
			t.Fatal("Ожидалась ошибка разбора JSON, но получен nil")
		}
	})

	// Тест на отсутствие ID
	t.Run("MissingID", func(t *testing.T) {
		apiClient := client.NewClient("http://example.com", "test_api_key")
		_, err := UpdateSegment(apiClient, &Segment{
			Name: "Сегмент без ID",
			Type: SegmentTypeDynamic,
		})

		if err == nil {
			t.Fatal("Ожидалась ошибка отсутствия ID, но получен nil")
		}
	})
}

// TestDeleteSegmentErrors проверяет обработку ошибок при удалении сегмента
func TestDeleteSegmentErrors(t *testing.T) {
	segmentID := 123

	// Тест на ошибку сервера
	t.Run("ServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")
		err := DeleteSegment(apiClient, segmentID)

		if err == nil {
			t.Fatal("Ожидалась ошибка HTTP статуса, но получен nil")
		}
	})
}

// TestAddContactsToSegmentErrors проверяет обработку ошибок при добавлении контактов в сегмент
func TestAddContactsToSegmentErrors(t *testing.T) {
	segmentID := 123
	contactIDs := []int{1001, 1002, 1003}

	// Тест на ошибку сервера
	t.Run("ServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")
		err := AddContactsToSegment(apiClient, segmentID, contactIDs)

		if err == nil {
			t.Fatal("Ожидалась ошибка HTTP статуса, но получен nil")
		}
	})
}

// TestRemoveContactsFromSegmentErrors проверяет обработку ошибок при удалении контактов из сегмента
func TestRemoveContactsFromSegmentErrors(t *testing.T) {
	segmentID := 123
	contactIDs := []int{1001, 1002, 1003}

	// Тест на ошибку сервера
	t.Run("ServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")
		err := RemoveContactsFromSegment(apiClient, segmentID, contactIDs)

		if err == nil {
			t.Fatal("Ожидалась ошибка HTTP статуса, но получен nil")
		}
	})
}

// TestGetSegmentContactsErrors проверяет обработку ошибок при получении контактов из сегмента
func TestGetSegmentContactsErrors(t *testing.T) {
	segmentID := 123

	// Тест на ошибку сервера
	t.Run("ServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")
		_, err := GetSegmentContacts(apiClient, segmentID, 1, 50)

		if err == nil {
			t.Fatal("Ожидалась ошибка HTTP статуса, но получен nil")
		}
	})

	// Тест на некорректный JSON
	t.Run("InvalidJSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{invalid_json`))
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")
		_, err := GetSegmentContacts(apiClient, segmentID, 1, 50)

		if err == nil {
			t.Fatal("Ожидалась ошибка разбора JSON, но получен nil")
		}
	})

	// Тест на пустой список контактов
	t.Run("EmptyContacts", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"page": 1,
				"per_page": 50,
				"_embedded": {
					"contacts": []
				},
				"_links": {
					"self": {
						"href": "/api/v4/segments/123/contacts"
					}
				}
			}`))
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")
		contactIDs, err := GetSegmentContacts(apiClient, segmentID, 1, 50)

		if err != nil {
			t.Fatalf("Неожиданная ошибка: %v", err)
		}
		if len(contactIDs) != 0 {
			t.Fatalf("Ожидался пустой список контактов, получено %d", len(contactIDs))
		}
	})
}
