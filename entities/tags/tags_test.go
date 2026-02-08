package tags

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/chudno/amo_crm_sdk/client"
)

func TestGetTags(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		expectedPath := "/api/v4/contacts/tags"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		expectedPage := "1"
		if r.URL.Query().Get("page") != expectedPage {
			t.Errorf("Ожидался параметр page=%s, получен %s", expectedPage, r.URL.Query().Get("page"))
		}

		expectedLimit := "50"
		if r.URL.Query().Get("limit") != expectedLimit {
			t.Errorf("Ожидался параметр limit=%s, получен %s", expectedLimit, r.URL.Query().Get("limit"))
		}

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

	apiClient := client.NewClient(server.URL, "test_api_key")

	tags, err := List(context.Background(), apiClient, EntityTypeContact, 1, 50)

	if err != nil {
		t.Fatalf("Ошибка при получении тегов: %v", err)
	}

	if len(tags) != 2 {
		t.Fatalf("Ожидалось получение 2 тегов, получено %d", len(tags))
	}

	expectedTag1 := Tag{
		ID:    123,
		Name:  "Важный клиент",
		Color: "#FF0000",
	}
	if !reflect.DeepEqual(tags[0], expectedTag1) {
		t.Errorf("Ожидался тег %+v, получен %+v", expectedTag1, tags[0])
	}

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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		expectedPath := "/api/v4/contacts/tags"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{
			"_embedded": {
				"tags": [
					{
						"id": 789,
						"name": "Новый тег",
						"color": "#0000FF"
					}
				]
			}
		}`))
	}))
	defer server.Close()

	apiClient := client.NewClient(server.URL, "test_api_key")

	newTag := &Tag{
		Name:  "Новый тег",
		Color: "#0000FF",
	}

	createdTag, err := Create(context.Background(), apiClient, EntityTypeContact, newTag)

	if err != nil {
		t.Fatalf("Ошибка при создании тега: %v", err)
	}

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

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		expectedPath := "/api/v4/contacts/tags"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		expectedFilterID := fmt.Sprintf("%d", tagID)
		filterIDs := r.URL.Query()["filter[id][]"]
		if len(filterIDs) == 0 || filterIDs[0] != expectedFilterID {
			t.Errorf("Ожидался параметр filter[id][]=%s, получен %v", expectedFilterID, filterIDs)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"page": 1,
			"per_page": 50,
			"total": 1,
			"_embedded": {
				"tags": [
					{
						"id": 123,
						"name": "Важный клиент",
						"color": "#FF0000"
					}
				]
			}
		}`))
	}))
	defer server.Close()

	apiClient := client.NewClient(server.URL, "test_api_key")

	tag, err := Get(context.Background(), apiClient, EntityTypeContact, tagID)

	if err != nil {
		t.Fatalf("Ошибка при получении тега: %v", err)
	}

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

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("Ожидался метод PATCH, получен %s", r.Method)
		}

		expectedPath := fmt.Sprintf("/api/v4/contacts/tags/%d", tagID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": 123,
			"name": "Очень важный клиент",
			"color": "#FF0000"
		}`))
	}))
	defer server.Close()

	apiClient := client.NewClient(server.URL, "test_api_key")

	tagToUpdate := &Tag{
		ID:    tagID,
		Name:  "Очень важный клиент",
		Color: "#FF0000",
	}

	updatedTag, err := Update(context.Background(), apiClient, EntityTypeContact, tagToUpdate)

	if err != nil {
		t.Fatalf("Ошибка при обновлении тега: %v", err)
	}

	expectedTag := &Tag{
		ID:    123,
		Name:  "Очень важный клиент",
		Color: "#FF0000",
	}
	if !reflect.DeepEqual(updatedTag, expectedTag) {
		t.Errorf("Ожидался тег %+v, получен %+v", expectedTag, updatedTag)
	}
}

func TestLinkEntityWithTags(t *testing.T) {
	entityID := 456

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		expectedPath := fmt.Sprintf("/api/v4/contacts/%d/tags", entityID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

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

	apiClient := client.NewClient(server.URL, "test_api_key")

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

	err := LinkEntity(context.Background(), apiClient, EntityTypeContact, entityID, tags)

	if err != nil {
		t.Fatalf("Ошибка при связывании сущности с тегами: %v", err)
	}
}

func TestGetEntityTags(t *testing.T) {
	entityID := 456

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Ожидался метод GET, получен %s", r.Method)
		}

		expectedPath := fmt.Sprintf("/api/v4/contacts/%d/tags", entityID)
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

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

	apiClient := client.NewClient(server.URL, "test_api_key")

	tags, err := ListForEntity(context.Background(), apiClient, EntityTypeContact, entityID)

	if err != nil {
		t.Fatalf("Ошибка при получении тегов сущности: %v", err)
	}

	if len(tags) != 2 {
		t.Fatalf("Ожидалось получение 2 тегов, получено %d", len(tags))
	}

	expectedTag1 := Tag{
		ID:    123,
		Name:  "Важный клиент",
		Color: "#FF0000",
	}
	if !reflect.DeepEqual(tags[0], expectedTag1) {
		t.Errorf("Ожидался тег %+v, получен %+v", expectedTag1, tags[0])
	}

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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Ожидался метод POST, получен %s", r.Method)
		}

		expectedPath := "/api/v4/contacts/tags"
		if r.URL.Path != expectedPath {
			t.Errorf("Ожидался путь %s, получен %s", expectedPath, r.URL.Path)
		}

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

	apiClient := client.NewClient(server.URL, "test_api_key")

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

	createdTags, err := CreateBatch(context.Background(), apiClient, EntityTypeContact, tagsToCreate)

	if err != nil {
		t.Fatalf("Ошибка при создании тегов: %v", err)
	}

	if len(createdTags) != 2 {
		t.Fatalf("Ожидалось создание 2 тегов, создано %d", len(createdTags))
	}

	expectedTag1 := Tag{
		ID:    123,
		Name:  "Тег 1",
		Color: "#FF0000",
	}
	if !reflect.DeepEqual(createdTags[0], expectedTag1) {
		t.Errorf("Ожидался тег %+v, получен %+v", expectedTag1, createdTags[0])
	}

	expectedTag2 := Tag{
		ID:    456,
		Name:  "Тег 2",
		Color: "#00FF00",
	}
	if !reflect.DeepEqual(createdTags[1], expectedTag2) {
		t.Errorf("Ожидался тег %+v, получен %+v", expectedTag2, createdTags[1])
	}
}
