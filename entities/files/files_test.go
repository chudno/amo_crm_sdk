package files

import (
	"context"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

func TestUploadFile(t *testing.T) {
	t.Run("UploadSuccess", func(t *testing.T) {
		testContent := []byte("Тестовое содержимое файла")
		tempFile := createUploadTestFile(t, "test_file.txt", testContent)

		entityType := EntityTypeLead
		entityID := 123

		server := setupUploadFileTestServer(t, entityType, entityID, "test_file.txt")
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		file, err := Upload(context.Background(), apiClient, entityType, entityID, tempFile)

		if err != nil {
			t.Fatalf("Ошибка при загрузке файла: %v", err)
		}

		verifyUploadedFileResult(t, file, testContent, "test_file.txt")
	})
}

func TestUploadFileByContent(t *testing.T) {
	t.Run("UploadByContentSuccess", func(t *testing.T) {
		testContent := []byte("Тестовое содержимое файла")
		fileName := "test_file.txt"

		entityType := EntityTypeLead
		entityID := 123

		server := setupUploadFileTestServer(t, entityType, entityID, fileName)
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		file, err := UploadByContent(context.Background(), apiClient, entityType, entityID, fileName, testContent)

		if err != nil {
			t.Fatalf("Ошибка при загрузке файла: %v", err)
		}

		verifyUploadedFileResult(t, file, testContent, fileName)
	})
}

func TestGetFiles(t *testing.T) {
	t.Run("GetFilesWithParams", func(t *testing.T) {
		entityType := EntityTypeLead
		entityID := 123

		server := setupGetFilesTestServer(t, entityType, entityID, true)
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		files, err := List(context.Background(), apiClient, entityType, entityID, 2, 30)

		if err != nil {
			t.Fatalf("Ошибка при получении списка файлов: %v", err)
		}

		verifyFilesList(t, files)
	})

	t.Run("GetFilesSimple", func(t *testing.T) {
		entityType := EntityTypeLead
		entityID := 123

		server := setupGetFilesTestServer(t, entityType, entityID, false)
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		files, err := List(context.Background(), apiClient, entityType, entityID, 1, 50)

		if err != nil {
			t.Fatalf("Ошибка при получении списка файлов: %v", err)
		}

		verifyFilesList(t, files)
	})
}

func TestDeleteFile(t *testing.T) {
	t.Run("DeleteSuccess", func(t *testing.T) {
		entityType := EntityTypeLead
		entityID := 123
		fileID := 456

		server := setupDeleteFileTestServer(t, entityType, entityID, fileID)
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		err := Delete(context.Background(), apiClient, entityType, entityID, fileID)

		if err != nil {
			t.Fatalf("Ошибка при удалении файла: %v", err)
		}
	})
}

func TestBatchDeleteFiles(t *testing.T) {
	t.Run("BatchDeleteSuccess", func(t *testing.T) {
		entityType := EntityTypeLead
		entityID := 123
		fileIDs := []int{456, 789}

		server := setupBatchDeleteFilesTestServer(t, entityType, entityID, fileIDs)
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		err := BatchDelete(context.Background(), apiClient, entityType, entityID, fileIDs)

		if err != nil {
			t.Fatalf("Ошибка при массовом удалении файлов: %v", err)
		}
	})
}

func TestGetFile(t *testing.T) {
	t.Run("GetFileSuccess", func(t *testing.T) {
		entityType := EntityTypeLead
		entityID := 123
		fileID := 456

		server := setupGetFileTestServer(t, entityType, entityID, fileID)
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")

		file, err := Get(context.Background(), apiClient, entityType, entityID, fileID)

		if err != nil {
			t.Fatalf("Ошибка при получении информации о файле: %v", err)
		}

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
	})
}
