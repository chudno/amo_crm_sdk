package files

import (
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

func TestUploadFile(t *testing.T) {
	t.Run("UploadSuccess", func(t *testing.T) {
		// Создаем тестовый файл и контент
		testContent := []byte("Тестовое содержимое файла")
		tempFile := createUploadTestFile(t, "test_file.txt", testContent)

		// Подготавливаем тестовые данные
		entityType := EntityTypeLead
		entityID := 123

		// Создаем тестовый сервер
		server := setupUploadFileTestServer(t, entityType, entityID, "test_file.txt")
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
		verifyUploadedFileResult(t, file, testContent, "test_file.txt")
	})
}

func TestUploadFileByContent(t *testing.T) {
	t.Run("UploadByContentSuccess", func(t *testing.T) {
		// Тестовые данные
		testContent := []byte("Тестовое содержимое файла")
		fileName := "test_file.txt"

		// Тип сущности и ID
		entityType := EntityTypeLead
		entityID := 123

		// Создаем тестовый сервер
		server := setupUploadFileTestServer(t, entityType, entityID, fileName)
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемый метод
		file, err := UploadFileByContent(apiClient, entityType, entityID, fileName, testContent)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при загрузке файла: %v", err)
		}

		// Проверяем полученные данные файла
		verifyUploadedFileResult(t, file, testContent, fileName)
	})
}

func TestGetFiles(t *testing.T) {
	t.Run("GetFilesWithParams", func(t *testing.T) {
		// Тип сущности и ID
		entityType := EntityTypeLead
		entityID := 123

		// Создаем тестовый сервер
		server := setupGetFilesTestServer(t, entityType, entityID, true)
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемый метод с параметрами пагинации
		files, err := GetFiles(apiClient, entityType, entityID, 2, 30)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при получении списка файлов: %v", err)
		}

		// Проверяем полученный список файлов
		verifyFilesList(t, files)
	})

	t.Run("GetFilesSimple", func(t *testing.T) {
		// Тип сущности и ID
		entityType := EntityTypeLead
		entityID := 123

		// Создаем тестовый сервер без проверки параметров
		server := setupGetFilesTestServer(t, entityType, entityID, false)
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемый метод с дефолтными параметрами
		files, err := GetFiles(apiClient, entityType, entityID, 1, 50)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при получении списка файлов: %v", err)
		}

		// Проверяем полученный список файлов
		verifyFilesList(t, files)
	})
}

func TestDeleteFile(t *testing.T) {
	t.Run("DeleteSuccess", func(t *testing.T) {
		// Тип сущности и ID
		entityType := EntityTypeLead
		entityID := 123
		fileID := 456

		// Создаем тестовый сервер
		server := setupDeleteFileTestServer(t, entityType, entityID, fileID)
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемый метод
		err := DeleteFile(apiClient, entityType, entityID, fileID)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при удалении файла: %v", err)
		}
	})
}

func TestBatchDeleteFiles(t *testing.T) {
	t.Run("BatchDeleteSuccess", func(t *testing.T) {
		// Тип сущности и ID
		entityType := EntityTypeLead
		entityID := 123
		fileIDs := []int{456, 789}

		// Создаем тестовый сервер
		server := setupBatchDeleteFilesTestServer(t, entityType, entityID, fileIDs)
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемый метод
		err := BatchDeleteFiles(apiClient, entityType, entityID, fileIDs)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при массовом удалении файлов: %v", err)
		}
	})
}

func TestGetFile(t *testing.T) {
	t.Run("GetFileSuccess", func(t *testing.T) {
		// Тип сущности и ID
		entityType := EntityTypeLead
		entityID := 123
		fileID := 456

		// Создаем тестовый сервер
		server := setupGetFileTestServer(t, entityType, entityID, fileID)
		defer server.Close()

		// Создаем клиент API
		apiClient := client.NewClient(server.URL, "test_api_key")

		// Вызываем тестируемый метод
		file, err := GetFile(apiClient, entityType, entityID, fileID)

		// Проверяем результаты
		if err != nil {
			t.Fatalf("Ошибка при получении информации о файле: %v", err)
		}

		// Проверяем идентификаторы
		if file.ID != 456 {
			t.Errorf("Ожидался ID файла 456, получен %d", file.ID)
		}

		if file.UUID != "test-file-uuid-123" {
			t.Errorf("Ожидался UUID test-file-uuid-123, получен %s", file.UUID)
		}

		// Проверяем размер и имя файла
		if file.Size != 1024 {
			t.Errorf("Ожидался размер файла 1024, получен %d", file.Size)
		}

		if file.Name != "file1.txt" {
			t.Errorf("Ожидалось имя файла file1.txt, получено %s", file.Name)
		}

		// Проверяем URL и ссылку для скачивания
		if file.URL != "https://example.amocrm.ru/files/file1.txt" {
			t.Errorf("Ожидался URL https://example.amocrm.ru/files/file1.txt, получен %s", file.URL)
		}

		if file.Download != "https://example.amocrm.ru/download/file1.txt" {
			t.Errorf("Ожидалась ссылка для скачивания https://example.amocrm.ru/download/file1.txt, получена %s", file.Download)
		}
	})
}
