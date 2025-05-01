package files

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// createUploadTestFile создает временный тестовый файл и возвращает его путь и содержимое
func createUploadTestFile(t *testing.T, filename string, content []byte) string {
	// Создаем временный файл для тестирования
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, filename)

	// Записываем тестовые данные во временный файл
	err := os.WriteFile(tempFile, content, 0644)
	if err != nil {
		t.Fatalf("Ошибка при создании временного файла: %v", err)
	}

	return tempFile
}

// setupUploadFileTestServer создает тестовый сервер для загрузки файла
func setupUploadFileTestServer(t *testing.T, entityType EntityType, entityID int, fileName string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/%s/%d/files", entityType, entityID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Проверяем Content-Type
		contentType := r.Header.Get("Content-Type")
		if contentType == "" || len(contentType) < 10 || contentType[:10] != "multipart/" {
			t.Errorf("Ожидался Content-Type multipart, получен %s", contentType)
		}

		// Проверяем, что файл был отправлен
		err := r.ParseMultipartForm(10 << 20) // Максимальный размер формы 10 MB
		if err != nil {
			t.Errorf("Ошибка при парсинге multipart формы: %v", err)
		}

		_, header, err := r.FormFile("file")
		if err != nil {
			t.Errorf("Ошибка при получении файла из формы: %v", err)
		}

		if header.Filename != fileName {
			t.Errorf("Ожидалось имя файла %s, получено %s", fileName, header.Filename)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(fmt.Sprintf(`{
			"_embedded": {
				"id": 456,
				"uuid": "test-file-uuid-123",
				"created_by": 789,
				"updated_by": 789,
				"created_at": 1609459200,
				"updated_at": 1609459200,
				"title": "%s",
				"url": "https://example.amocrm.ru/files/%s",
				"download_link": "https://example.amocrm.ru/download/%s",
				"_links": {
					"self": {
						"href": "/api/v4/leads/%d/files/456"
					},
					"download": {
						"href": "/download/%s"
					}
				}
			}
		}`, fileName, fileName, fileName, entityID, fileName)))
	}))
}

// verifyUploadedFileResult проверяет результат загрузки файла
func verifyUploadedFileResult(t *testing.T, file *File, testContent []byte, fileName string) {
	if file.ID != 456 {
		t.Errorf("Ожидался ID файла 456, получен %d", file.ID)
	}

	if file.UUID != "test-file-uuid-123" {
		t.Errorf("Ожидался UUID test-file-uuid-123, получен %s", file.UUID)
	}

	if file.Size != len(testContent) {
		t.Errorf("Ожидался размер файла %d, получен %d", len(testContent), file.Size)
	}

	if file.Name != fileName {
		t.Errorf("Ожидалось имя файла %s, получено %s", fileName, file.Name)
	}

	expectedURL := fmt.Sprintf("https://example.amocrm.ru/files/%s", fileName)
	if file.URL != expectedURL {
		t.Errorf("Ожидался URL %s, получен %s", expectedURL, file.URL)
	}

	expectedDownload := fmt.Sprintf("https://example.amocrm.ru/download/%s", fileName)
	if file.Download != expectedDownload {
		t.Errorf("Ожидалась ссылка для скачивания %s, получена %s", expectedDownload, file.Download)
	}
}

// setupGetFilesTestServer создает тестовый сервер для получения списка файлов
func setupGetFilesTestServer(t *testing.T, entityType EntityType, entityID int, withQueryParams bool) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/%s/%d/files", entityType, entityID)
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
			"_embedded": {
				"files": [
					{
						"id": 456,
						"uuid": "test-file-uuid-123",
						"entity_id": 123,
						"entity_type": "leads",
						"created_by": 789,
						"updated_by": 789,
						"created_at": 1609459200,
						"updated_at": 1609459200,
						"size": 1024,
						"name": "file1.txt",
						"title": "file1.txt",
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
					},
					{
						"id": 457,
						"uuid": "test-file-uuid-124",
						"entity_id": 123,
						"entity_type": "leads",
						"created_by": 789,
						"updated_by": 789,
						"created_at": 1609458200,
						"updated_at": 1609458200,
						"size": 2048,
						"name": "file2.pdf",
						"title": "file2.pdf",
						"url": "https://example.amocrm.ru/files/file2.pdf",
						"download_link": "https://example.amocrm.ru/download/file2.pdf",
						"_links": {
							"self": {
								"href": "/api/v4/leads/123/files/457"
							},
							"download": {
								"href": "/download/file2.pdf"
							}
						}
					}
				]
			},
			"_next_page": "/api/v4/leads/123/files?page=3&limit=30",
			"_prev_page": "/api/v4/leads/123/files?page=1&limit=30",
			"_total_path": "/api/v4/leads/123/files/total"
		}`))
	}))
}

// verifyFilesList проверяет результат получения списка файлов
func verifyFilesList(t *testing.T, files []File) {
	if len(files) != 2 {
		t.Fatalf("Ожидалось 2 файла, получено %d", len(files))
	}

	// Проверяем содержимое первого файла
	if files[0].ID != 456 {
		t.Errorf("Ожидался ID 456, получен %d", files[0].ID)
	}

	if files[0].UUID != "test-file-uuid-123" {
		t.Errorf("Ожидался UUID test-file-uuid-123, получен %s", files[0].UUID)
	}

	if files[0].Size != 1024 {
		t.Errorf("Ожидался размер 1024, получен %d", files[0].Size)
	}

	if files[0].Name != "file1.txt" {
		t.Errorf("Ожидалось имя file1.txt, получено %s", files[0].Name)
	}

	if files[0].URL != "https://example.amocrm.ru/files/file1.txt" {
		t.Errorf("Ожидался URL https://example.amocrm.ru/files/file1.txt, получен %s", files[0].URL)
	}

	if files[0].Download != "https://example.amocrm.ru/download/file1.txt" {
		t.Errorf("Ожидалась ссылка для скачивания https://example.amocrm.ru/download/file1.txt, получена %s", files[0].Download)
	}

	// Проверяем содержимое второго файла
	if files[1].ID != 457 {
		t.Errorf("Ожидался ID 457, получен %d", files[1].ID)
	}

	if files[1].UUID != "test-file-uuid-124" {
		t.Errorf("Ожидался UUID test-file-uuid-124, получен %s", files[1].UUID)
	}

	if files[1].Name != "file2.pdf" {
		t.Errorf("Ожидалось имя file2.pdf, получено %s", files[1].Name)
	}
}

// setupGetFileTestServer создает тестовый сервер для получения информации о файле
func setupGetFileTestServer(t *testing.T, entityType EntityType, entityID, fileID int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/%s/%d/files/%d", entityType, entityID, fileID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 456,
			"uuid": "test-file-uuid-123",
			"entity_id": 123,
			"entity_type": "leads",
			"created_by": 789,
			"updated_by": 789,
			"created_at": 1609459200,
			"updated_at": 1609459200,
			"size": 1024,
			"name": "file1.txt",
			"title": "file1.txt",
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
		}`))
	}))
}

// setupDeleteFileTestServer создает тестовый сервер для удаления файла
func setupDeleteFileTestServer(t *testing.T, entityType EntityType, entityID, fileID int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "DELETE" {
			t.Errorf("Ожидался метод DELETE, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/%s/%d/files/%d", entityType, entityID, fileID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.WriteHeader(http.StatusNoContent)
	}))
}

// setupBatchDeleteFilesTestServer создает тестовый сервер для массового удаления файлов
func setupBatchDeleteFilesTestServer(t *testing.T, entityType EntityType, entityID int, fileIDs []int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "DELETE" {
			t.Errorf("Ожидался метод DELETE, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/%s/%d/files", entityType, entityID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Проверяем параметр фильтра
		expectedFilter := "456,789"
		if r.URL.Query().Get("filter[id]") != expectedFilter {
			t.Errorf("Ожидался параметр filter[id]=%s, получен %s", expectedFilter, r.URL.Query().Get("filter[id]"))
		}

		// Отправляем ответ
		w.WriteHeader(http.StatusNoContent)
	}))
}
