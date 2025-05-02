package files

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

// TestGetFileErrors проверяет обработку ошибок при получении файла
func TestGetFileErrors(t *testing.T) {
	t.Run("Файл не найден", func(t *testing.T) {
		// Используем несуществующий домен для вызова ошибки
		apiClient := client.NewClient("http://non-existent-domain.example", "test_api_key")

		// Вызываем тестируемую функцию
		_, err := GetFile(apiClient, EntityTypeLead, 123, 456)

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка, но она не возникла")
		}
	})

	t.Run("Некорректный JSON", func(t *testing.T) {
		// Создаем тестовый сервер, который вернет некорректный JSON
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(`{"id": 456, "invalid_json":`)); err != nil {
				t.Fatalf("Ошибка при записи ответа: %v", err)
			}
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемую функцию
		_, err := GetFile(apiClient, EntityTypeLead, 123, 456)

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка при декодировании JSON, но она не возникла")
		}
	})

	t.Run("Ошибка сервера", func(t *testing.T) {
		// Создаем тестовый сервер, который вернет ошибку сервера
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte(`{"status": 500, "title": "Internal Server Error"}`)); err != nil {
				t.Fatalf("Ошибка при записи ответа: %v", err)
			}
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемую функцию
		_, err := GetFile(apiClient, EntityTypeLead, 123, 456)

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка сервера, но она не возникла")
		}
	})
}

// TestGetFilesErrors проверяет обработку ошибок при получении списка файлов
func TestGetFilesErrors(t *testing.T) {
	t.Run("Ошибка HTTP", func(t *testing.T) {
		// Используем несуществующий домен для вызова ошибки
		apiClient := client.NewClient("http://non-existent-domain.example", "test_api_key")

		// Вызываем тестируемую функцию
		_, err := GetFiles(apiClient, EntityTypeLead, 123, 1, 50)

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка HTTP, но она не возникла")
		}
	})

	t.Run("Некорректный JSON-ответ", func(t *testing.T) {
		// Создаем тестовый сервер, который вернет некорректный JSON
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(`{"page": 1, "_embedded": {"files": [{"id": 123, "name": "test.txt"`)); err != nil {
				t.Fatalf("Ошибка при записи ответа: %v", err)
			}
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемую функцию
		_, err := GetFiles(apiClient, EntityTypeLead, 123, 1, 50)

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка при декодировании JSON, но она не возникла")
		}
	})

	t.Run("Пустой список файлов", func(t *testing.T) {
		// Создаем тестовый сервер, который вернет пустой список файлов
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(`{"page": 1, "per_page": 50, "total": 0, "_embedded": {"files": []}}`)); err != nil {
				t.Fatalf("Ошибка при записи ответа: %v", err)
			}
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемую функцию
		files, err := GetFiles(apiClient, EntityTypeLead, 123, 1, 50)

		// Проверяем, что ошибки нет и список пустой
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
		// Используем несуществующий домен для вызова ошибки
		apiClient := client.NewClient("http://non-existent-domain.example", "test_api_key")

		// Вызываем тестируемую функцию
		err := DeleteFile(apiClient, EntityTypeLead, 123, 456)

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка HTTP, но она не возникла")
		}
	})

	t.Run("Ошибка сервера", func(t *testing.T) {
		// Создаем тестовый сервер, который вернет ошибку
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Проверяем, что это запрос DELETE
			if r.Method != http.MethodDelete {
				t.Errorf("Ожидался метод DELETE, получен %s", r.Method)
			}

			// Проверяем путь запроса
			expectedPath := "/api/v4/leads/123/files/456"
			if r.URL.Path != expectedPath {
				t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
			}

			// Возвращаем ошибку сервера
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte(`{"status": 500, "title": "Internal Server Error"}`)); err != nil {
				t.Fatalf("Ошибка при записи ответа: %v", err)
			}
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемую функцию
		err := DeleteFile(apiClient, EntityTypeLead, 123, 456)

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка сервера, но она не возникла")
		}
	})
}

// TestBatchDeleteFilesErrors проверяет обработку ошибок при массовом удалении файлов
func TestBatchDeleteFilesErrors(t *testing.T) {
	t.Run("Ошибка HTTP", func(t *testing.T) {
		// Используем несуществующий домен для вызова ошибки
		apiClient := client.NewClient("http://non-existent-domain.example", "test_api_key")

		// Вызываем тестируемую функцию
		err := BatchDeleteFiles(apiClient, EntityTypeLead, 123, []int{456, 789})

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка HTTP, но она не возникла")
		}
	})

	t.Run("Ошибка сервера", func(t *testing.T) {
		// Создаем тестовый сервер, который вернет ошибку
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Проверяем, что это запрос DELETE
			if r.Method != http.MethodDelete {
				t.Errorf("Ожидался метод DELETE, получен %s", r.Method)
			}

			// Проверяем, что в URL есть фильтр по ID
			if r.URL.Query().Get("filter[id]") != "456,789" {
				t.Errorf("Ожидался фильтр filter[id]=456,789, получен %s", r.URL.Query().Get("filter[id]"))
			}

			// Возвращаем ошибку сервера
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte(`{"status": 500, "title": "Internal Server Error"}`)); err != nil {
				t.Fatalf("Ошибка при записи ответа: %v", err)
			}
		}))
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемую функцию
		err := BatchDeleteFiles(apiClient, EntityTypeLead, 123, []int{456, 789})

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка сервера, но она не возникла")
		}
	})
}

// TestGetDownloadFileURL проверяет работу функции получения URL для скачивания файла
func TestGetDownloadFileURL(t *testing.T) {
	t.Run("Успешное получение URL", func(t *testing.T) {
		// Тип сущности и ID
		entityType := EntityTypeLead
		entityID := 123
		fileID := 456

		// Создаем тестовый сервер
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Проверяем параметры запроса
			expectedPath := fmt.Sprintf("/api/v4/%s/%d/files/%d", entityType, entityID, fileID)
			if r.URL.Path != expectedPath {
				t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
			}

			// Отправляем ответ с данными файла и ссылкой для скачивания
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

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемую функцию
		downloadURL, err := GetDownloadFileURL(apiClient, entityType, entityID, fileID)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при получении URL для скачивания: %v", err)
		}

		expectedURL := server.URL + "/download/file1.txt"
		if downloadURL != expectedURL {
			t.Errorf("Ожидался URL %s, получен %s", expectedURL, downloadURL)
		}
	})

	t.Run("Отсутствие ссылки для скачивания", func(t *testing.T) {
		// Тип сущности и ID
		entityType := EntityTypeLead
		entityID := 123
		fileID := 456

		// Создаем тестовый сервер, который вернет данные файла без ссылки для скачивания
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

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемую функцию
		_, err := GetDownloadFileURL(apiClient, entityType, entityID, fileID)

		// Проверяем, что вернулась ошибка
		if err == nil {
			t.Error("Ожидалась ошибка из-за отсутствия ссылки для скачивания, но она не возникла")
		}
	})
}
