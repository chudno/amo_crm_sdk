package urlfilters

import (
	"testing"
)

func TestNewLeadFilterFromURL(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		wantPage  int
		wantLimit int
		wantErr   bool
	}{
		{
			name:      "URL лидов с пагинацией",
			url:       "https://example.amocrm.ru/leads/list/?filter[name]=Test&page=2&limit=25",
			wantPage:  2,
			wantLimit: 25,
			wantErr:   false,
		},
		{
			name:      "URL лидов без пагинации",
			url:       "https://example.amocrm.ru/leads/list/?filter[name]=Test",
			wantPage:  1,
			wantLimit: 50,
			wantErr:   false,
		},
		{
			name:      "URL контактов (неверный тип)",
			url:       "https://example.amocrm.ru/contacts/list/?filter[name]=Test",
			wantPage:  0,
			wantLimit: 0,
			wantErr:   true,
		},
		{
			name:      "Неверный URL",
			url:       "invalid-url",
			wantPage:  0,
			wantLimit: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewLeadFilterFromURL(tt.url)

			// Проверяем ошибку
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLeadFilterFromURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Проверяем номер страницы
			if result.PageInt != tt.wantPage {
				t.Errorf("NewLeadFilterFromURL() page = %v, want %v", result.PageInt, tt.wantPage)
			}

			// Проверяем лимит
			if result.LimitInt != tt.wantLimit {
				t.Errorf("NewLeadFilterFromURL() limit = %v, want %v", result.LimitInt, tt.wantLimit)
			}
		})
	}
}

func TestLeadFilter_GetSDKFilterMap(t *testing.T) {
	t.Run("Возвращает правильную карту фильтров", func(t *testing.T) {
		url := "https://example.amocrm.ru/leads/list/?filter[name]=Test&filter[status][]=12345"
		filter, err := NewLeadFilterFromURL(url)
		if err != nil {
			t.Fatalf("Ошибка при парсинге URL: %v", err)
		}

		sdkFilter := filter.GetSDKFilterMap()

		if len(sdkFilter) != 2 {
			t.Errorf("GetSDKFilterMap() filter count = %v, want 2", len(sdkFilter))
		}

		expectedFilters := map[string]string{
			"filter[name]":     "Test",
			"filter[status][]": "12345",
		}

		for key, wantValue := range expectedFilters {
			gotValue, exists := sdkFilter[key]
			if !exists {
				t.Errorf("GetSDKFilterMap() filter missing key %v", key)
				continue
			}
			if gotValue != wantValue {
				t.Errorf("GetSDKFilterMap() filter[%v] = %v, want %v", key, gotValue, wantValue)
			}
		}
	})
}
