package tags

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

func TestGetTags(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/contacts/tags"
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
				"tags": [
					{
						"id": 123,
						"name": "Важный клиент",
						"color": "#FF0000"
					},
					{
						"id": 456,
						"name": "Потенциальный клиент",
						"color": "#00FF00"
					}
				]
			}
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	tags, err := GetTags(apiClient, EntityTypeContact, 1, 50)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при получении тегов: %v", err)
	}

	if len(tags) != 2 {
		t.Fatalf("Ожидалось получение 2 тегов, получено %d", len(tags))
	}

	// Проверяем содержимое первого тега
	expectedTag1 := Tag{
		ID:    123,
		Name:  "Важный клиент",
		Color: "#FF0000",
	}
	if !reflect.DeepEqual(tags[0], expectedTag1) {
		t.Errorf("Ожидался тег %+v, получен %+v", expectedTag1, tags[0])
	}

	// Проверяем содержимое второго тега
	expectedTag2 := Tag{
		ID:    456,
		Name:  "Потенциальный клиент",
		Color: "#00FF00",
	}
	if !reflect.DeepEqual(tags[1], expectedTag2) {
		t.Errorf("Ожидался тег %+v, получен %+v", expectedTag2, tags[1])
	}
}

func TestCreateTag(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/contacts/tags"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{
			"tag": {
				"id": 789,
				"name": "Новый тег",
				"color": "#0000FF"
			}
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Создаем тег для отправки
	newTag := &Tag{
		Name:  "Новый тег",
		Color: "#0000FF",
	}

	// Вызываем тестируемый метод
	createdTag, err := CreateTag(apiClient, EntityTypeContact, newTag)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при создании тега: %v", err)
	}

	// Проверяем содержимое созданного тега
	expectedTag := &Tag{
		ID:    789,
		Name:  "Новый тег",
		Color: "#0000FF",
	}
	if !reflect.DeepEqual(createdTag, expectedTag) {
		t.Errorf("Ожидался тег %+v, получен %+v", expectedTag, createdTag)
	}
}

func TestGetTag(t *testing.T) {
	tagID := 123

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/contacts/tags/%d", tagID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 123,
			"name": "Важный клиент",
			"color": "#FF0000"
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	tag, err := GetTag(apiClient, EntityTypeContact, tagID)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при получении тега: %v", err)
	}

	// Проверяем содержимое тега
	expectedTag := &Tag{
		ID:    123,
		Name:  "Важный клиент",
		Color: "#FF0000",
	}
	if !reflect.DeepEqual(tag, expectedTag) {
		t.Errorf("Ожидался тег %+v, получен %+v", expectedTag, tag)
	}
}

func TestUpdateTag(t *testing.T) {
	tagID := 123

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "PATCH" {
			t.Errorf("Ожидался метод PATCH, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/contacts/tags/%d", tagID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 123,
			"name": "Очень важный клиент",
			"color": "#FF0000"
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Создаем тег для обновления
	tagToUpdate := &Tag{
		ID:    tagID,
		Name:  "Очень важный клиент",
		Color: "#FF0000",
	}

	// Вызываем тестируемый метод
	updatedTag, err := UpdateTag(apiClient, EntityTypeContact, tagToUpdate)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при обновлении тега: %v", err)
	}

	// Проверяем содержимое обновленного тега
	expectedTag := &Tag{
		ID:    123,
		Name:  "Очень важный клиент",
		Color: "#FF0000",
	}
	if !reflect.DeepEqual(updatedTag, expectedTag) {
		t.Errorf("Ожидался тег %+v, получен %+v", expectedTag, updatedTag)
	}
}

func TestDeleteTag(t *testing.T) {
	tagID := 123

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "DELETE" {
			t.Errorf("Ожидался метод DELETE, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/contacts/tags/%d", tagID)
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
	err := DeleteTag(apiClient, EntityTypeContact, tagID)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при удалении тега: %v", err)
	}
}

func TestLinkEntityWithTags(t *testing.T) {
	entityID := 456

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/contacts/%d/tags", entityID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"_embedded": {
				"tags": [
					{
						"id": 123,
						"name": "Важный клиент",
						"color": "#FF0000"
					},
					{
						"id": 456,
						"name": "Потенциальный клиент",
						"color": "#00FF00"
					}
				]
			}
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Создаем теги для связывания
	tags := []Tag{
		{
			Name:  "Важный клиент",
			Color: "#FF0000",
		},
		{
			Name:  "Потенциальный клиент",
			Color: "#00FF00",
		},
	}

	// Вызываем тестируемый метод
	err := LinkEntityWithTags(apiClient, EntityTypeContact, entityID, tags)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при связывании сущности с тегами: %v", err)
	}
}

func TestGetEntityTags(t *testing.T) {
	entityID := 456

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := fmt.Sprintf("/api/v4/contacts/%d/tags", entityID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"_embedded": {
				"tags": [
					{
						"id": 123,
						"name": "Важный клиент",
						"color": "#FF0000"
					},
					{
						"id": 456,
						"name": "Потенциальный клиент",
						"color": "#00FF00"
					}
				]
			}
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Вызываем тестируемый метод
	tags, err := GetEntityTags(apiClient, EntityTypeContact, entityID)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при получении тегов сущности: %v", err)
	}

	if len(tags) != 2 {
		t.Fatalf("Ожидалось получение 2 тегов, получено %d", len(tags))
	}

	// Проверяем содержимое первого тега
	expectedTag1 := Tag{
		ID:    123,
		Name:  "Важный клиент",
		Color: "#FF0000",
	}
	if !reflect.DeepEqual(tags[0], expectedTag1) {
		t.Errorf("Ожидался тег %+v, получен %+v", expectedTag1, tags[0])
	}

	// Проверяем содержимое второго тега
	expectedTag2 := Tag{
		ID:    456,
		Name:  "Потенциальный клиент",
		Color: "#00FF00",
	}
	if !reflect.DeepEqual(tags[1], expectedTag2) {
		t.Errorf("Ожидался тег %+v, получен %+v", expectedTag2, tags[1])
	}
}

func TestCreateTags(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		// Проверяем путь запроса
		expectedPath := "/api/v4/contacts/tags"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{
			"_embedded": {
				"tags": [
					{
						"id": 123,
						"name": "Тег 1",
						"color": "#FF0000"
					},
					{
						"id": 456,
						"name": "Тег 2",
						"color": "#00FF00"
					}
				]
			}
		}`))
	}))
	defer server.Close()

	// Создаем клиент API
	apiClient := client.NewClient(server.URL, "test_api_key")

	// Создаем теги для отправки
	tagsToCreate := []Tag{
		{
			Name:  "Тег 1",
			Color: "#FF0000",
		},
		{
			Name:  "Тег 2",
			Color: "#00FF00",
		},
	}

	// Вызываем тестируемый метод
	createdTags, err := CreateTags(apiClient, EntityTypeContact, tagsToCreate)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("Ошибка при создании тегов: %v", err)
	}

	if len(createdTags) != 2 {
		t.Fatalf("Ожидалось создание 2 тегов, создано %d", len(createdTags))
	}

	// Проверяем содержимое первого тега
	expectedTag1 := Tag{
		ID:    123,
		Name:  "Тег 1",
		Color: "#FF0000",
	}
	if !reflect.DeepEqual(createdTags[0], expectedTag1) {
		t.Errorf("Ожидался тег %+v, получен %+v", expectedTag1, createdTags[0])
	}

	// Проверяем содержимое второго тега
	expectedTag2 := Tag{
		ID:    456,
		Name:  "Тег 2",
		Color: "#00FF00",
	}
	if !reflect.DeepEqual(createdTags[1], expectedTag2) {
		t.Errorf("Ожидался тег %+v, получен %+v", expectedTag2, createdTags[1])
	}
}
