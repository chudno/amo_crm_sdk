package files

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

// TestGetFileErrors проверяет обработку ошибок при получении файла
func TestGetFileErrors(t *testing.T) {
	t.Run("Файл не найден", func(t *testing.T) {
		apiClient := client.NewClient("http://non-existent-domain.example", "test_api_key")

		_, err := Get(context.Background(), apiClient, EntityTypeLead, 123, 456)

		if err == nil {
			t.Error("Ожидалась ошибка, но она не возникла")
		}
	})

	t.Run("Некорректный JSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(`{"id": 456, "invalid_json":`)); err != nil {
				t.Fatalf("Ошибка при записи ответа: %v", err)
			}
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		_, err := Get(context.Background(), apiClient, EntityTypeLead, 123, 456)

		if err == nil {
			t.Error("Ожидалась ошибка при декодировании JSON, но она не возникла")
		}
	})

	t.Run("Ошибка сервера", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte(`{"status": 500, "title": "Internal Server Error"}`)); err != nil {
				t.Fatalf("Ошибка при записи ответа: %v", err)
			}
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		_, err := Get(context.Background(), apiClient, EntityTypeLead, 123, 456)

		if err == nil {
			t.Error("Ожидалась ошибка сервера, но она не возникла")
		}
	})
}

// TestGetFilesErrors проверяет обработку ошибок при получении списка файлов
func TestGetFilesErrors(t *testing.T) {
	t.Run("Ошибка HTTP", func(t *testing.T) {
		apiClient := client.NewClient("http://non-existent-domain.example", "test_api_key")

		_, err := List(context.Background(), apiClient, EntityTypeLead, 123, 1, 50)

		if err == nil {
			t.Error("Ожидалась ошибка HTTP, но она не возникла")
		}
	})

	t.Run("Некорректный JSON-ответ", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(`{"page": 1, "_embedded": {"files": [{"id": 123, "name": "test.txt"`)); err != nil {
				t.Fatalf("Ошибка при записи ответа: %v", err)
			}
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		_, err := List(context.Background(), apiClient, EntityTypeLead, 123, 1, 50)

		if err == nil {
			t.Error("Ожидалась ошибка при декодировании JSON, но она не возникла")
		}
	})

	t.Run("Пустой список файлов", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(`{"page": 1, "per_page": 50, "total": 0, "_embedded": {"files": []}}`)); err != nil {
				t.Fatalf("Ошибка при записи ответа: %v", err)
			}
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		files, err := List(context.Background(), apiClient, EntityTypeLead, 123, 1, 50)

		if err != nil {
			t.Errorf("Не ожидалась ошибка, но получено: %v", err)
		}

		if len(files) != 0 {
			t.Errorf("Ожидался пустой список файлов, но получен список длиной %d", len(files))
		}
	})
}

// TestDeleteFileErrors проверяет обработку ошибок при удалении файла
func TestDeleteFileErrors(t *testing.T) {
	t.Run("Ошибка HTTP", func(t *testing.T) {
		apiClient := client.NewClient("http://non-existent-domain.example", "test_api_key")

		err := Delete(context.Background(), apiClient, EntityTypeLead, 123, 456)

		if err == nil {
			t.Error("Ожидалась ошибка HTTP, но она не возникла")
		}
	})

	t.Run("Ошибка сервера", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				t.Errorf("Ожидался метод DELETE, получен %s", r.Method)
			}

			expectedPath := "/api/v4/leads/123/files/456"
			if r.URL.Path != expectedPath {
				t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte(`{"status": 500, "title": "Internal Server Error"}`)); err != nil {
				t.Fatalf("Ошибка при записи ответа: %v", err)
			}
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		err := Delete(context.Background(), apiClient, EntityTypeLead, 123, 456)

		if err == nil {
			t.Error("Ожидалась ошибка сервера, но она не возникла")
		}
	})
}

// TestBatchDeleteFilesErrors проверяет обработку ошибок при массовом удалении файлов
func TestBatchDeleteFilesErrors(t *testing.T) {
	t.Run("Ошибка HTTP", func(t *testing.T) {
		apiClient := client.NewClient("http://non-existent-domain.example", "test_api_key")

		err := BatchDelete(context.Background(), apiClient, EntityTypeLead, 123, []int{456, 789})

		if err == nil {
			t.Error("Ожидалась ошибка HTTP, но она не возникла")
		}
	})

	t.Run("Ошибка сервера", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				t.Errorf("Ожидался метод DELETE, получен %s", r.Method)
			}

			if r.URL.Query().Get("filter[id]") != "456,789" {
				t.Errorf("Ожидался фильтр filter[id]=456,789, получен %s", r.URL.Query().Get("filter[id]"))
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte(`{"status": 500, "title": "Internal Server Error"}`)); err != nil {
				t.Fatalf("Ошибка при записи ответа: %v", err)
			}
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		err := BatchDelete(context.Background(), apiClient, EntityTypeLead, 123, []int{456, 789})

		if err == nil {
			t.Error("Ожидалась ошибка сервера, но она не возникла")
		}
	})
}

// TestGetDownloadFileURL проверяет работу функции получения URL для скачивания файла
func TestGetDownloadFileURL(t *testing.T) {
	t.Run("Успешное получение URL", func(t *testing.T) {
		entityType := EntityTypeLead
		entityID := 123
		fileID := 456

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			expectedPath := fmt.Sprintf("/api/v4/%s/%d/files/%d", entityType, entityID, fileID)
			if r.URL.Path != expectedPath {
				t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
			}

			w.Header().Set("Content-Type", "application/json")
			response := `{
				"id": 456,
				"uuid": "test-file-uuid-123",
				"entity_id": 123,
				"entity_type": "leads",
				"size": 1024,
				"name": "file1.txt",
				"url": "https://example.amocrm.ru/files/file1.txt",
				"download_link": "https://example.amocrm.ru/download/file1.txt",
				"_links": {
					"self": {
						"href": "/api/v4/leads/123/files/456"
					},
					"download": {
						"href": "/download/file1.txt"
					}
				}
			}`
			if _, err := w.Write([]byte(response)); err != nil {
				t.Fatalf("Ошибка при записи ответа: %v", err)
			}
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		downloadURL, err := GetDownloadURL(context.Background(), apiClient, entityType, entityID, fileID)

		if err != nil {
			t.Fatalf("Ошибка при получении URL для скачивания: %v", err)
		}

		expectedURL := server.URL + "/download/file1.txt"
		if downloadURL != expectedURL {
			t.Errorf("Ожидался URL %s, получен %s", expectedURL, downloadURL)
		}
	})

	t.Run("Отсутствие ссылки для скачивания", func(t *testing.T) {
		entityType := EntityTypeLead
		entityID := 123
		fileID := 456

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			response := `{
				"id": 456,
				"uuid": "test-file-uuid-123",
				"entity_id": 123,
				"entity_type": "leads",
				"size": 1024,
				"name": "file1.txt",
				"url": "https://example.amocrm.ru/files/file1.txt",
				"_links": {
					"self": {
						"href": "/api/v4/leads/123/files/456"
					}
				}
			}`
			if _, err := w.Write([]byte(response)); err != nil {
				t.Fatalf("Ошибка при записи ответа: %v", err)
			}
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		_, err := GetDownloadURL(context.Background(), apiClient, entityType, entityID, fileID)

		if err == nil {
			t.Error("Ожидалась ошибка из-за отсутствия ссылки для скачивания, но она не возникла")
		}
	})
}
