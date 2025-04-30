package files

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

func TestUploadFile(t *testing.T) {
	// Создаем временный файл для тестирования
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test_file.txt")

	// Записываем тестовые данные во временный файл
	testContent := []byte("Тестовое содержимое файла")
	err := ioutil.WriteFile(tempFile, testContent, 0644)
	if err != nil {
		t.Fatalf("Ошибка при создании временного файла: %v", err)
	}

	// Тип сущности и ID
	entityType := EntityTypeLead
	entityID := 123

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		if header.Filename != "test_file.txt" {
			t.Errorf("Ожидалось имя файла test_file.txt, получено %s", header.Filename)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{
			"_embedded": {
				"id": 456,
				"uuid": "test-file-uuid-123",
				"created_by": 789,
				"updated_by": 789,
				"created_at": 1609459200,
				"updated_at": 1609459200,
				"title": "test_file.txt",
				"url": "https://example.amocrm.ru/files/test_file.txt",
				"download_link": "https://example.amocrm.ru/download/test_file.txt",
				"_links": {
					"self": {
						"href": "/api/v4/leads/123/files/456"
					},
					"download": {
						"href": "/download/test_file.txt"
					}
				}
			}
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	file, err := UploadFile(apiClient, entityType, entityID, tempFile)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при загрузке файла: %v", err)
	}

	// Проверяем полученные данные файла
	if file.ID != 456 {
		t.Errorf("Ожидался ID файла 456, получен %d", file.ID)
	}

	if file.UUID != "test-file-uuid-123" {
		t.Errorf("Ожидался UUID test-file-uuid-123, получен %s", file.UUID)
	}

	if file.Size != len(testContent) {
		t.Errorf("Ожидался размер файла %d, получен %d", len(testContent), file.Size)
	}

	if file.Name != "test_file.txt" {
		t.Errorf("Ожидалось имя файла test_file.txt, получено %s", file.Name)
	}

	if file.URL != "https://example.amocrm.ru/files/test_file.txt" {
		t.Errorf("Ожидался URL https://example.amocrm.ru/files/test_file.txt, получен %s", file.URL)
	}

	if file.Download != "https://example.amocrm.ru/download/test_file.txt" {
		t.Errorf("Ожидалась ссылка для скачивания https://example.amocrm.ru/download/test_file.txt, получена %s", file.Download)
	}
}

func TestUploadFileByContent(t *testing.T) {
	// Тип сущности и ID
	entityType := EntityTypeContact
	entityID := 456

	// Тестовое содержимое файла
	fileName := "test_file.txt"
	fileContent := []byte("Тестовое содержимое файла")

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		_, _ = w.Write([]byte(`{
			"_embedded": {
				"id": 789,
				"uuid": "test-file-uuid-456",
				"created_by": 123,
				"updated_by": 123,
				"created_at": 1609459200,
				"updated_at": 1609459200,
				"title": "test_file.txt",
				"url": "https://example.amocrm.ru/files/test_file.txt",
				"download_link": "https://example.amocrm.ru/download/test_file.txt",
				"_links": {
					"self": {
						"href": "/api/v4/contacts/456/files/789"
					},
					"download": {
						"href": "/download/test_file.txt"
					}
				}
			}
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	file, err := UploadFileByContent(apiClient, entityType, entityID, fileName, fileContent)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при загрузке файла по содержимому: %v", err)
	}

	// Проверяем полученные данные файла
	if file.ID != 789 {
		t.Errorf("Ожидался ID файла 789, получен %d", file.ID)
	}

	if file.UUID != "test-file-uuid-456" {
		t.Errorf("Ожидался UUID test-file-uuid-456, получен %s", file.UUID)
	}

	if file.Size != len(fileContent) {
		t.Errorf("Ожидался размер файла %d, получен %d", len(fileContent), file.Size)
	}

	if file.Name != fileName {
		t.Errorf("Ожидалось имя файла %s, получено %s", fileName, file.Name)
	}

	if file.URL != "https://example.amocrm.ru/files/test_file.txt" {
		t.Errorf("Ожидался URL https://example.amocrm.ru/files/test_file.txt, получен %s", file.URL)
	}

	if file.Download != "https://example.amocrm.ru/download/test_file.txt" {
		t.Errorf("Ожидалась ссылка для скачивания https://example.amocrm.ru/download/test_file.txt, получена %s", file.Download)
	}
}

func TestGetFiles(t *testing.T) {
	// Тип сущности и ID
	entityType := EntityTypeLead
	entityID := 123

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		expectedPage := "1"
		if r.URL.Query().Get("page") != expectedPage {
			t.Errorf("Ожидался параметр page=%s, получен %s", expectedPage, r.URL.Query().Get("page"))
		}

		expectedLimit := "50"
		if r.URL.Query().Get("limit") != expectedLimit {
			t.Errorf("Ожидался параметр limit=%s, получен %s", expectedLimit, r.URL.Query().Get("limit"))
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"page": 1,
			"per_page": 50,
			"total": 2,
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
						"id": 789,
						"uuid": "test-file-uuid-456",
						"entity_id": 123,
						"entity_type": "leads",
						"created_by": 789,
						"updated_by": 789,
						"created_at": 1609459200,
						"updated_at": 1609459200,
						"size": 2048,
						"name": "file2.txt",
						"title": "file2.txt",
						"url": "https://example.amocrm.ru/files/file2.txt",
						"download_link": "https://example.amocrm.ru/download/file2.txt",
						"_links": {
							"self": {
								"href": "/api/v4/leads/123/files/789"
							},
							"download": {
								"href": "/download/file2.txt"
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

	// Вызываем тестируемый метод
	files, err := GetFiles(apiClient, entityType, entityID, 1, 50)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при получении списка файлов: %v", err)
	}

	if len(files) != 2 {
		t.Fatalf("Ожидалось получение 2 файлов, получено %d", len(files))
	}

	// Проверяем содержимое первого файла
	if files[0].ID != 456 {
		t.Errorf("Ожидался ID первого файла 456, получен %d", files[0].ID)
	}

	if files[0].Name != "file1.txt" {
		t.Errorf("Ожидалось имя первого файла file1.txt, получено %s", files[0].Name)
	}

	// Проверяем содержимое второго файла
	if files[1].ID != 789 {
		t.Errorf("Ожидался ID второго файла 789, получен %d", files[1].ID)
	}

	if files[1].Name != "file2.txt" {
		t.Errorf("Ожидалось имя второго файла file2.txt, получено %s", files[1].Name)
	}
}

func TestDeleteFile(t *testing.T) {
	// Тип сущности и ID
	entityType := EntityTypeLead
	entityID := 123
	fileID := 456

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	err := DeleteFile(apiClient, entityType, entityID, fileID)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при удалении файла: %v", err)
	}
}

func TestBatchDeleteFiles(t *testing.T) {
	// Тип сущности и ID
	entityType := EntityTypeLead
	entityID := 123
	fileIDs := []int{456, 789}

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	err := BatchDeleteFiles(apiClient, entityType, entityID, fileIDs)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при массовом удалении файлов: %v", err)
	}
}

func TestGetFile(t *testing.T) {
	// Тип сущности и ID
	entityType := EntityTypeLead
	entityID := 123
	fileID := 456

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	file, err := GetFile(apiClient, entityType, entityID, fileID)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при получении информации о файле: %v", err)
	}

	// Проверяем полученные данные файла
	if file.ID != 456 {
		t.Errorf("Ожидался ID файла 456, получен %d", file.ID)
	}

	if file.UUID != "test-file-uuid-123" {
		t.Errorf("Ожидался UUID test-file-uuid-123, получен %s", file.UUID)
	}

	if file.Size != 1024 {
		t.Errorf("Ожидался размер файла 1024, получен %d", file.Size)
	}

	if file.Name != "file1.txt" {
		t.Errorf("Ожидалось имя файла file1.txt, получено %s", file.Name)
	}

	if file.URL != "https://example.amocrm.ru/files/file1.txt" {
		t.Errorf("Ожидался URL https://example.amocrm.ru/files/file1.txt, получен %s", file.URL)
	}

	if file.Download != "https://example.amocrm.ru/download/file1.txt" {
		t.Errorf("Ожидалась ссылка для скачивания https://example.amocrm.ru/download/file1.txt, получена %s", file.Download)
	}
}
