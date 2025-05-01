package short_links

import (
	"net/http"
	"strings"
	"testing"
)

// createGetShortLinksSuccessMockClient создает мок-клиент для успешного запроса списка коротких ссылок
func createGetShortLinksSuccessMockClient() *AdvancedMockClient {
	return &AdvancedMockClient{
		BaseURL:        "https://example.amocrm.ru",
		ExpectedMethod: "GET",
		ExpectedURL:    "/api/v4/short_links",
		MockResponse: &MockResponse{
			StatusCode: http.StatusOK,
			Body: `{
				"_page": 1,
				"_links": {
					"self": {
						"href": "https://example.amocrm.ru/api/v4/short_links?page=1&limit=50"
					}
				},
				"_embedded": {
					"short_links": [
						{
							"id": 123,
							"url": "https://example.com/test1",
							"key": "abc123",
							"short_url": "https://amo.im/abc123",
							"account_id": 12345,
							"entity_id": 67890,
							"entity_type": "leads",
							"created_at": 1609459200,
							"created_by": 1111,
							"updated_at": 1609459300,
							"visit_count": 42,
							"last_visit_at": 1609459400
						},
						{
							"id": 456,
							"url": "https://example.com/test2",
							"key": "def456",
							"short_url": "https://amo.im/def456",
							"account_id": 12345,
							"entity_id": 54321,
							"entity_type": "contacts",
							"created_at": 1609459500,
							"created_by": 2222,
							"updated_at": 1609459600,
							"visit_count": 24,
							"last_visit_at": 1609459700
						}
					]
				}
			}`,
		},
	}
}

// createGetShortLinksEmptyMockClient создает мок-клиент для запроса пустого списка коротких ссылок
func createGetShortLinksEmptyMockClient() *AdvancedMockClient {
	return &AdvancedMockClient{
		BaseURL:        "https://example.amocrm.ru",
		ExpectedMethod: "GET",
		ExpectedURL:    "/api/v4/short_links",
		MockResponse: &MockResponse{
			StatusCode: http.StatusOK,
			Body: `{
				"_page": 1,
				"_links": {
					"self": {
						"href": "https://example.amocrm.ru/api/v4/short_links?page=1&limit=50"
					}
				},
				"_embedded": {
					"short_links": []
				}
			}`,
		},
	}
}

// createGetShortLinksErrorMockClient создает мок-клиент с ошибкой сервера
func createGetShortLinksErrorMockClient() *AdvancedMockClient {
	return &AdvancedMockClient{
		BaseURL:        "https://example.amocrm.ru",
		ExpectedMethod: "GET",
		ExpectedURL:    "/api/v4/short_links",
		MockResponse: &MockResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"status": "error", "message": "Internal Server Error"}`,
		},
	}
}

// createGetShortLinksWithFilterMockClient создает мок-клиент для запроса с фильтром
func createGetShortLinksWithFilterMockClient() *AdvancedMockClient {
	return &AdvancedMockClient{
		BaseURL:        "https://example.amocrm.ru",
		ExpectedMethod: "GET",
		ExpectedURL:    "/api/v4/short_links",
		MockResponse: &MockResponse{
			StatusCode: http.StatusOK,
			Body: `{
				"_page": 1,
				"_links": {
					"self": {
						"href": "https://example.amocrm.ru/api/v4/short_links?page=1&limit=50&filter[entity_type]=leads"
					}
				},
				"_embedded": {
					"short_links": [
						{
							"id": 123,
							"url": "https://example.com/test1",
							"key": "abc123",
							"short_url": "https://amo.im/abc123",
							"entity_type": "leads",
							"entity_id": 67890
						}
					]
				}
			}`,
		},
	}
}

// verifyFirstShortLink проверяет данные первой ссылки
func verifyFirstShortLink(t *testing.T, link ShortLink) {
	if link.ID != 123 ||
		link.URL != "https://example.com/test1" ||
		link.Key != "abc123" ||
		link.ShortURL != "https://amo.im/abc123" ||
		link.EntityType != "leads" ||
		link.EntityID != 67890 ||
		link.VisitCount != 42 {
		t.Errorf("Данные первой ссылки не соответствуют ожидаемым")
	}
}

// verifySecondShortLink проверяет данные второй ссылки
func verifySecondShortLink(t *testing.T, link ShortLink) {
	if link.ID != 456 ||
		link.URL != "https://example.com/test2" ||
		link.Key != "def456" ||
		link.EntityType != "contacts" {
		t.Errorf("Данные второй ссылки не соответствуют ожидаемым")
	}
}

// verifyFilteredShortLink проверяет данные отфильтрованной ссылки
func verifyFilteredShortLink(t *testing.T, link ShortLink) {
	if link.EntityType != "leads" {
		t.Errorf("Тип сущности не соответствует ожидаемому: %s", link.EntityType)
	}
}

// TestGetShortLinks проверяет получение списка коротких ссылок
func TestGetShortLinks(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := createGetShortLinksSuccessMockClient()

		// Вызываем тестируемый метод
		links, err := GetShortLinksWithRequester(mockClient, 1, 50)

		// Проверка наличия ошибки
		if err != nil {
			t.Fatalf("Ошибка при получении списка коротких ссылок: %v", err)
		}

		// Проверка количества полученных ссылок
		if len(links) != 2 {
			t.Errorf("Ожидалось получение 2 ссылок, получено %d", len(links))
		}

		// Проверка данных ссылок
		verifyFirstShortLink(t, links[0])
		verifySecondShortLink(t, links[1])
	})

	t.Run("Empty", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := createGetShortLinksEmptyMockClient()

		// Вызываем тестируемый метод
		links, err := GetShortLinksWithRequester(mockClient, 1, 50)

		// Проверка наличия ошибки
		if err != nil {
			t.Fatalf("Ошибка при получении списка коротких ссылок: %v", err)
		}

		// Проверка количества полученных ссылок
		if len(links) != 0 {
			t.Errorf("Ожидался пустой список ссылок, получено %d", len(links))
		}
	})

	t.Run("ServerError", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := createGetShortLinksErrorMockClient()

		// Вызываем тестируемый метод
		_, err := GetShortLinksWithRequester(mockClient, 1, 50)

		// Проверка наличия ошибки
		if err == nil {
			t.Fatal("Ожидалась ошибка, но её не было")
		}
	})

	t.Run("WithFilter", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := createGetShortLinksWithFilterMockClient()

		// Создаем фильтр
		filter := map[string]string{
			"filter[entity_type]": "leads",
		}

		// Вызываем тестируемый метод с фильтром
		links, err := GetShortLinksWithRequester(mockClient, 1, 50, WithFilter(filter))

		// Проверка наличия ошибки
		if err != nil {
			t.Fatalf("Ошибка при получении списка коротких ссылок: %v", err)
		}

		// Проверка количества полученных ссылок
		if len(links) != 1 {
			t.Errorf("Ожидалось получение 1 ссылки, получено %d", len(links))
		}

		// Проверка фильтра в запросе
		if !strings.Contains(mockClient.LastRequest.URL, "filter%5Bentity_type%5D=leads") {
			t.Errorf("Фильтр не был добавлен к URL запроса: %s", mockClient.LastRequest.URL)
		}

		// Проверка данных ссылки
		verifyFilteredShortLink(t, links[0])
	})
}

// TestGetShortLink проверяет получение информации о конкретной короткой ссылке
func TestGetShortLink(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// ID ссылки для теста
		linkID := 123

		// Создаем мок-клиент
		mockClient := &AdvancedMockClient{
			BaseURL:        "https://example.amocrm.ru",
			ExpectedMethod: "GET",
			ExpectedURL:    "/api/v4/short_links/123",
			MockResponse: &MockResponse{
				StatusCode: http.StatusOK,
				Body: `{
					"id": 123,
					"url": "https://example.com/test",
					"key": "abc123",
					"short_url": "https://amo.im/abc123",
					"account_id": 12345,
					"entity_id": 67890,
					"entity_type": "leads",
					"created_at": 1609459200,
					"created_by": 1111,
					"updated_at": 1609459300,
					"visit_count": 42,
					"last_visit_at": 1609459400
				}`,
			},
		}

		// Вызываем тестируемый метод
		link, err := GetShortLinkWithRequester(mockClient, linkID)

		// Проверка наличия ошибки
		if err != nil {
			t.Fatalf("Ошибка при получении короткой ссылки: %v", err)
		}

		// Проверка данных ссылки
		if link.ID != linkID ||
			link.URL != "https://example.com/test" ||
			link.Key != "abc123" ||
			link.ShortURL != "https://amo.im/abc123" ||
			link.EntityType != "leads" ||
			link.EntityID != 67890 ||
			link.VisitCount != 42 {
			t.Errorf("Данные ссылки не соответствуют ожидаемым")
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		// ID несуществующей ссылки
		linkID := 999

		// Создаем мок-клиент
		mockClient := &AdvancedMockClient{
			BaseURL:        "https://example.amocrm.ru",
			ExpectedMethod: "GET",
			ExpectedURL:    "/api/v4/short_links/999",
			MockResponse: &MockResponse{
				StatusCode: http.StatusNotFound,
				Body:       `{"status": "error", "message": "Short link not found"}`,
			},
		}

		// Вызываем тестируемый метод
		_, err := GetShortLinkWithRequester(mockClient, linkID)

		// Проверка наличия ошибки
		if err == nil {
			t.Fatal("Ожидалась ошибка, но её не было")
		}
	})
}

// TestCreateShortLink проверяет создание новой короткой ссылки
func TestCreateShortLink(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := &AdvancedMockClient{
			BaseURL:        "https://example.amocrm.ru",
			ExpectedMethod: "POST",
			ExpectedURL:    "/api/v4/short_links",
			ExpectedBody: &ShortLink{
				URL:        "https://example.com/test",
				EntityType: "leads",
				EntityID:   123,
			},
			MockResponse: &MockResponse{
				StatusCode: http.StatusCreated,
				Body: `{
					"_embedded": {
						"short_links": [
							{
								"id": 456,
								"url": "https://example.com/test",
								"key": "def456",
								"short_url": "https://amo.im/def456",
								"account_id": 12345,
								"entity_id": 123,
								"entity_type": "leads",
								"created_at": 1609459200,
								"created_by": 1111,
								"updated_at": 1609459200
							}
						]
					}
				}`,
			},
		}

		// Создаем новую ссылку
		newLink := &ShortLink{
			URL:        "https://example.com/test",
			EntityType: "leads",
			EntityID:   123,
		}

		// Вызываем тестируемый метод
		createdLink, err := CreateShortLinkWithRequester(mockClient, newLink)

		// Проверка наличия ошибки
		if err != nil {
			t.Fatalf("Ошибка при создании короткой ссылки: %v", err)
		}

		// Проверка данных созданной ссылки
		if createdLink.ID != 456 ||
			createdLink.URL != "https://example.com/test" ||
			createdLink.Key != "def456" ||
			createdLink.ShortURL != "https://amo.im/def456" ||
			createdLink.EntityType != "leads" ||
			createdLink.EntityID != 123 {
			t.Errorf("Данные созданной ссылки не соответствуют ожидаемым")
		}
	})

	t.Run("Error", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := &AdvancedMockClient{
			BaseURL:        "https://example.amocrm.ru",
			ExpectedMethod: "POST",
			ExpectedURL:    "/api/v4/short_links",
			ExpectedBody: &ShortLink{
				URL: "https://example.com/test",
			},
			MockResponse: &MockResponse{
				StatusCode: http.StatusBadRequest,
				Body:       `{"status": "error", "message": "Invalid parameters"}`,
			},
		}

		// Создаем новую ссылку с некорректными данными
		newLink := &ShortLink{
			URL: "https://example.com/test",
		}

		// Вызываем тестируемый метод
		_, err := CreateShortLinkWithRequester(mockClient, newLink)

		// Проверка наличия ошибки
		if err == nil {
			t.Fatal("Ожидалась ошибка, но её не было")
		}
	})
}

// createUpdateShortLinkSuccessMockClient создает мок-клиент для успешного обновления короткой ссылки
func createUpdateShortLinkSuccessMockClient() *AdvancedMockClient {
	return &AdvancedMockClient{
		BaseURL:        "https://example.amocrm.ru",
		ExpectedMethod: "PATCH",
		ExpectedURL:    "/api/v4/short_links/123",
		ExpectedBody: &ShortLink{
			ID:        123,
			URL:       "https://updated-example.com/test",
			ExpireAt:  1640995200, // 01.01.2022
			UTMSource: "newsletter",
		},
		MockResponse: &MockResponse{
			StatusCode: http.StatusOK,
			Body: `{
				"id": 123,
				"url": "https://updated-example.com/test",
				"key": "abc123",
				"short_url": "https://amo.im/abc123",
				"account_id": 12345,
				"entity_id": 67890,
				"entity_type": "leads",
				"created_at": 1609459200,
				"created_by": 1111,
				"updated_at": 1609545600,
				"expire_at": 1640995200,
				"utm_source": "newsletter"
			}`,
		},
	}
}

// createUpdateShortLinkErrorMockClient создает мок-клиент с ошибкой при обновлении
func createUpdateShortLinkErrorMockClient() *AdvancedMockClient {
	return &AdvancedMockClient{
		BaseURL:        "https://example.amocrm.ru",
		ExpectedMethod: "PATCH",
		ExpectedURL:    "/api/v4/short_links/0",
		MockResponse: &MockResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"status": "error", "message": "Invalid ID"}`,
		},
	}
}

// createUpdateShortLinkData создает данные для обновления короткой ссылки
func createUpdateShortLinkData(linkID int) *ShortLink {
	return &ShortLink{
		ID:        linkID,
		URL:       "https://updated-example.com/test",
		ExpireAt:  1640995200, // 01.01.2022
		UTMSource: "newsletter",
	}
}

// createInvalidUpdateShortLinkData создает невалидные данные для обновления
func createInvalidUpdateShortLinkData() *ShortLink {
	return &ShortLink{
		URL: "https://updated-example.com/test",
	}
}

// verifyUpdatedShortLink проверяет обновленную короткую ссылку
func verifyUpdatedShortLink(t *testing.T, updatedLink *ShortLink, linkID int) {
	if updatedLink.ID != linkID ||
		updatedLink.URL != "https://updated-example.com/test" ||
		updatedLink.ExpireAt != 1640995200 ||
		updatedLink.UTMSource != "newsletter" {
		t.Errorf("Данные обновленной ссылки не соответствуют ожидаемым")
	}
}

// TestUpdateShortLink проверяет обновление существующей короткой ссылки
func TestUpdateShortLink(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// ID ссылки для обновления
		linkID := 123

		// Создаем мок-клиент
		mockClient := createUpdateShortLinkSuccessMockClient()

		// Создаем ссылку для обновления
		updateLink := createUpdateShortLinkData(linkID)

		// Вызываем тестируемый метод
		updatedLink, err := UpdateShortLinkWithRequester(mockClient, updateLink)

		// Проверка наличия ошибки
		if err != nil {
			t.Fatalf("Ошибка при обновлении короткой ссылки: %v", err)
		}

		// Проверка данных обновленной ссылки
		verifyUpdatedShortLink(t, updatedLink, linkID)
	})

	t.Run("Error", func(t *testing.T) {
		// Создаем мок-клиент
		mockClient := createUpdateShortLinkErrorMockClient()

		// Создаем ссылку без ID
		updateLink := createInvalidUpdateShortLinkData()

		// Вызываем тестируемый метод
		_, err := UpdateShortLinkWithRequester(mockClient, updateLink)

		// Проверка наличия ошибки
		if err == nil {
			t.Fatal("Ожидалась ошибка, но её не было")
		}
	})
}

// createDeleteShortLinkSuccessMockClient создает мок-клиент для успешного удаления короткой ссылки
func createDeleteShortLinkSuccessMockClient() *AdvancedMockClient {
	return &AdvancedMockClient{
		BaseURL:        "https://example.amocrm.ru",
		ExpectedMethod: "DELETE",
		ExpectedURL:    "/api/v4/short_links/123",
		MockResponse: &MockResponse{
			StatusCode: http.StatusNoContent,
			Body:       ``,
		},
	}
}

// createDeleteShortLinkErrorMockClient создает мок-клиент с ошибкой при удалении
func createDeleteShortLinkErrorMockClient() *AdvancedMockClient {
	return &AdvancedMockClient{
		BaseURL:        "https://example.amocrm.ru",
		ExpectedMethod: "DELETE",
		ExpectedURL:    "/api/v4/short_links/999",
		MockResponse: &MockResponse{
			StatusCode: http.StatusNotFound,
			Body:       `{"status": "error", "message": "Short link not found"}`,
		},
	}
}

// TestDeleteShortLink проверяет удаление короткой ссылки
func TestDeleteShortLink(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// ID ссылки для удаления
		linkID := 123

		// Создаем мок-клиент
		mockClient := createDeleteShortLinkSuccessMockClient()

		// Вызываем тестируемый метод
		err := DeleteShortLinkWithRequester(mockClient, linkID)

		// Проверка наличия ошибки
		if err != nil {
			t.Fatalf("Ошибка при удалении короткой ссылки: %v", err)
		}
	})

	t.Run("Error", func(t *testing.T) {
		// ID несуществующей ссылки
		linkID := 999

		// Создаем мок-клиент
		mockClient := createDeleteShortLinkErrorMockClient()

		// Вызываем тестируемый метод
		err := DeleteShortLinkWithRequester(mockClient, linkID)

		// Проверка наличия ошибки
		if err == nil {
			t.Fatal("Ожидалась ошибка, но её не было")
		}
	})
}

// createGetShortLinkStatsSuccessMockClient создает мок-клиент для успешного получения статистики
func createGetShortLinkStatsSuccessMockClient() *AdvancedMockClient {
	return &AdvancedMockClient{
		BaseURL:        "https://example.amocrm.ru",
		ExpectedMethod: "GET",
		ExpectedURL:    "/api/v4/short_links/123/statistics",
		MockResponse: &MockResponse{
			StatusCode: http.StatusOK,
			Body: `{
				"id": 123,
				"url": "https://example.com/test",
				"key": "abc123",
				"short_url": "https://amo.im/abc123",
				"account_id": 12345,
				"entity_id": 67890,
				"entity_type": "leads",
				"created_at": 1609459200,
				"created_by": 1111,
				"updated_at": 1609459300,
				"visit_count": 42,
				"last_visit_at": 1609459400
			}`,
		},
	}
}

// createGetShortLinkStatsErrorMockClient создает мок-клиент с ошибкой при получении статистики
func createGetShortLinkStatsErrorMockClient() *AdvancedMockClient {
	return &AdvancedMockClient{
		BaseURL:        "https://example.amocrm.ru",
		ExpectedMethod: "GET",
		ExpectedURL:    "/api/v4/short_links/999/statistics",
		MockResponse: &MockResponse{
			StatusCode: http.StatusNotFound,
			Body:       `{"status": "error", "message": "Short link not found"}`,
		},
	}
}

// verifyShortLinkStats проверяет данные статистики короткой ссылки
func verifyShortLinkStats(t *testing.T, stats *ShortLink, linkID int) {
	if stats.ID != linkID ||
		stats.VisitCount != 42 ||
		stats.LastVisitAt != 1609459400 {
		t.Errorf("Данные статистики не соответствуют ожидаемым")
	}
}

// TestGetShortLinkStats проверяет получение статистики использования короткой ссылки
func TestGetShortLinkStats(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// ID ссылки для получения статистики
		linkID := 123

		// Создаем мок-клиент
		mockClient := createGetShortLinkStatsSuccessMockClient()

		// Вызываем тестируемый метод
		stats, err := GetShortLinkStatsWithRequester(mockClient, linkID)

		// Проверка наличия ошибки
		if err != nil {
			t.Fatalf("Ошибка при получении статистики короткой ссылки: %v", err)
		}

		// Проверка данных статистики
		verifyShortLinkStats(t, stats, linkID)
	})

	t.Run("Error", func(t *testing.T) {
		// ID несуществующей ссылки
		linkID := 999

		// Создаем мок-клиент
		mockClient := createGetShortLinkStatsErrorMockClient()

		// Вызываем тестируемый метод
		_, err := GetShortLinkStatsWithRequester(mockClient, linkID)

		// Проверка наличия ошибки
		if err == nil {
			t.Fatal("Ожидалась ошибка, но её не было")
		}
	})
}
