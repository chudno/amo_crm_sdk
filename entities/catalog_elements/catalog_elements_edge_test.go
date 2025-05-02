package catalog_elements

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

// TestGetCatalogElementsErrors проверяет обработку ошибок при получении списка элементов каталога
func TestGetCatalogElementsErrors(t *testing.T) {
	catalogID := 123

	// Тест на ошибку сервера
	t.Run("ServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")
		_, err := GetCatalogElements(apiClient, catalogID, 1, 50, nil)

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
		_, err := GetCatalogElements(apiClient, catalogID, 1, 50, nil)

		if err == nil {
			t.Fatal("Ожидалась ошибка разбора JSON, но получен nil")
		}
	})

	// Тест на пустой список элементов
	t.Run("EmptyElements", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"page": 1,
				"per_page": 50,
				"total": 0,
				"_embedded": {
					"elements": []
				}
			}`))
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")
		elements, err := GetCatalogElements(apiClient, catalogID, 1, 50, nil)

		if err != nil {
			t.Fatalf("Неожиданная ошибка: %v", err)
		}
		if len(elements) != 0 {
			t.Fatalf("Ожидался пустой список элементов, получено %d", len(elements))
		}
	})
}

// TestCreateCatalogElementErrors проверяет обработку ошибок при создании элемента каталога
func TestCreateCatalogElementErrors(t *testing.T) {
	catalogID := 123
	element := &CatalogElement{
		Name:      "Тестовый элемент",
		CatalogID: catalogID,
	}

	// Тест на ошибку сервера
	t.Run("ServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")
		_, err := CreateCatalogElement(apiClient, catalogID, element)

		if err == nil {
			t.Fatal("Ожидалась ошибка HTTP статуса, но получен nil")
		}
	})

	// Тест на некорректный JSON в ответе
	t.Run("InvalidJSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{invalid_json`))
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")
		_, err := CreateCatalogElement(apiClient, catalogID, element)

		if err == nil {
			t.Fatal("Ожидалась ошибка разбора JSON, но получен nil")
		}
	})

	// Тест на пустой список элементов в ответе
	t.Run("EmptyResponse", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{
				"_embedded": {
					"elements": []
				}
			}`))
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")
		_, err := CreateCatalogElement(apiClient, catalogID, element)

		if err == nil {
			t.Fatal("Ожидалась ошибка пустого ответа, но получен nil")
		}
	})
}

// TestGetCatalogElementErrors проверяет обработку ошибок при получении элемента каталога
func TestGetCatalogElementErrors(t *testing.T) {
	catalogID := 123
	elementID := 456

	// Тест на ошибку сервера
	t.Run("ServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")
		_, err := GetCatalogElement(apiClient, catalogID, elementID)

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
		_, err := GetCatalogElement(apiClient, catalogID, elementID)

		if err == nil {
			t.Fatal("Ожидалась ошибка разбора JSON, но получен nil")
		}
	})

	// Тест на отсутствие элемента
	t.Run("NotFound", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")
		_, err := GetCatalogElement(apiClient, catalogID, elementID)

		if err == nil {
			t.Fatal("Ожидалась ошибка HTTP статуса, но получен nil")
		}
	})
}

// TestUpdateCatalogElementErrors проверяет обработку ошибок при обновлении элемента каталога
func TestUpdateCatalogElementErrors(t *testing.T) {
	catalogID := 123
	element := &CatalogElement{
		ID:        456,
		Name:      "Тестовый элемент",
		CatalogID: catalogID,
	}

	// Тест на ошибку сервера
	t.Run("ServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")
		_, err := UpdateCatalogElement(apiClient, catalogID, element)

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
		_, err := UpdateCatalogElement(apiClient, catalogID, element)

		if err == nil {
			t.Fatal("Ожидалась ошибка разбора JSON, но получен nil")
		}
	})

	// Тест на отсутствие ID
	t.Run("MissingID", func(t *testing.T) {
		apiClient := client.NewClient("http://example.com", "test_api_key")
		_, err := UpdateCatalogElement(apiClient, catalogID, &CatalogElement{
			Name:      "Тестовый элемент",
			CatalogID: catalogID,
		})

		if err == nil {
			t.Fatal("Ожидалась ошибка отсутствия ID, но получен nil")
		}
	})
}

// TestDeleteCatalogElementErrors проверяет обработку ошибок при удалении элемента каталога
func TestDeleteCatalogElementErrors(t *testing.T) {
	catalogID := 123
	elementID := 456

	// Тест на ошибку сервера
	t.Run("ServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")
		err := DeleteCatalogElement(apiClient, catalogID, elementID)

		if err == nil {
			t.Fatal("Ожидалась ошибка HTTP статуса, но получен nil")
		}
	})
}

// TestBatchDeleteCatalogElementsExtraTests проверяет обработку ошибок при пакетном удалении элементов каталога
func TestBatchDeleteCatalogElementsExtraTests(t *testing.T) {
	catalogID := 123
	elementIDs := []int{456, 789}

	// Тест на ошибку сервера
	t.Run("ServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")
		err := BatchDeleteCatalogElements(apiClient, catalogID, elementIDs)

		if err == nil {
			t.Fatal("Ожидалась ошибка HTTP статуса, но получен nil")
		}
	})
}

// TestLinkCatalogElementWithTagsErrors проверяет обработку ошибок при привязке тегов к элементу каталога
func TestLinkCatalogElementWithTagsErrors(t *testing.T) {
	catalogID := 123
	elementID := 456
	tags := []Tag{
		{
			Name:  "Тег 1",
			Color: "#FF0000",
		},
	}

	// Тест на ошибку сервера
	t.Run("ServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")
		err := LinkCatalogElementWithTags(apiClient, catalogID, elementID, tags)

		if err == nil {
			t.Fatal("Ожидалась ошибка HTTP статуса, но получен nil")
		}
	})
}

// TestGetCatalogElementTagsErrors проверяет обработку ошибок при получении тегов элемента каталога
func TestGetCatalogElementTagsErrors(t *testing.T) {
	catalogID := 123
	elementID := 456

	// Тест на ошибку сервера
	t.Run("ServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		apiClient := client.NewClient(server.URL, "test_api_key")
		_, err := GetCatalogElementTags(apiClient, catalogID, elementID)

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
		_, err := GetCatalogElementTags(apiClient, catalogID, elementID)

		if err == nil {
			t.Fatal("Ожидалась ошибка разбора JSON, но получен nil")
		}
	})
}
