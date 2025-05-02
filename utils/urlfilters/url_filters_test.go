package urlfilters

import (
	"strings"
	"testing"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		wantEntity  string
		wantFilters map[string]string
		wantPage    string
		wantLimit   string
		wantErr     bool
	}{
		{
			name:        "Простой URL сделок",
			url:         "https://example.amocrm.ru/leads/list/",
			wantEntity:  "leads",
			wantFilters: map[string]string{},
			wantPage:    "1",
			wantLimit:   "50",
			wantErr:     false,
		},
		{
			name:        "URL сделок с фильтром по имени",
			url:         "https://example.amocrm.ru/leads/list/?filter[name]=Тестовая+сделка",
			wantEntity:  "leads",
			wantFilters: map[string]string{"filter[name]": "Тестовая сделка"},
			wantPage:    "1",
			wantLimit:   "50",
			wantErr:     false,
		},
		{
			name:        "URL сделок с фильтром по статусу",
			url:         "https://example.amocrm.ru/leads/list/?filter[status][]=10073462&filter[status][]=10073459",
			wantEntity:  "leads",
			wantFilters: map[string]string{"filter[status][]": "10073462"},
			wantPage:    "1",
			wantLimit:   "50",
			wantErr:     false,
		},
		{
			name:        "URL сделок с фильтром по ответственному и пагинацией",
			url:         "https://example.amocrm.ru/leads/list/?filter[responsible_user_id][]=9057966&page=2&limit=25",
			wantEntity:  "leads",
			wantFilters: map[string]string{"filter[responsible_user_id][]": "9057966"},
			wantPage:    "2",
			wantLimit:   "25",
			wantErr:     false,
		},
		{
			name:        "URL контактов с фильтром",
			url:         "https://example.amocrm.ru/contacts/list/?filter[name]=Иван",
			wantEntity:  "contacts",
			wantFilters: map[string]string{"filter[name]": "Иван"},
			wantPage:    "1",
			wantLimit:   "50",
			wantErr:     false,
		},
		{
			name:       "URL сделок с множественными фильтрами",
			url:        "https://example.amocrm.ru/leads/list/?filter[name]=Тест&filter[price][from]=1000&filter[price][to]=5000",
			wantEntity: "leads",
			wantFilters: map[string]string{
				"filter[name]":        "Тест",
				"filter[price][from]": "1000",
				"filter[price][to]":   "5000",
			},
			wantPage:  "1",
			wantLimit: "50",
			wantErr:   false,
		},
		{
			name:        "Неверный URL",
			url:         "invalid-url",
			wantEntity:  "",
			wantFilters: nil,
			wantPage:    "",
			wantLimit:   "",
			wantErr:     true,
		},
		{
			name:        "URL без типа сущности",
			url:         "https://example.amocrm.ru/settings/",
			wantEntity:  "",
			wantFilters: nil,
			wantPage:    "",
			wantLimit:   "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseURL(tt.url)

			// Проверяем ошибку
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			expectedEntityType := "" // Извлекаем тип сущности
			if strings.Contains(tt.url, "/leads/list") {
				expectedEntityType = "leads" // Тип для лидов
			} else if strings.Contains(tt.url, "/contacts/list") {
				expectedEntityType = "contacts"
			} else if strings.Contains(tt.url, "/companies/list") {
				expectedEntityType = "companies"
			}

			if result.EntityType != expectedEntityType {
				t.Errorf("ParseURL() entityType = %v, want %v", result.EntityType, expectedEntityType)
			}

			// Проверяем номер страницы
			if result.Page != tt.wantPage {
				t.Errorf("ParseURL() page = %v, want %v", result.Page, tt.wantPage)
			}

			// Проверяем лимит
			if result.Limit != tt.wantLimit {
				t.Errorf("ParseURL() limit = %v, want %v", result.Limit, tt.wantLimit)
			}

			// Проверяем фильтры
			if len(result.Filter) != len(tt.wantFilters) {
				t.Errorf("ParseURL() filter count = %v, want %v", len(result.Filter), len(tt.wantFilters))
			}

			for key, wantValue := range tt.wantFilters {
				gotValue, exists := result.Filter[key]
				if !exists {
					t.Errorf("ParseURL() filter missing key %v", key)
					continue
				}
				if gotValue != wantValue {
					t.Errorf("ParseURL() filter[%v] = %v, want %v", key, gotValue, wantValue)
				}
			}
		})
	}
}

func TestExtractEntityType(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "Путь лидов",
			path: "/leads/list/",
			want: "leads", // Теперь в URL и SDK используется одинаковое название
		},
		{
			name: "Путь контактов",
			path: "/contacts/list/",
			want: "contacts",
		},
		{
			name: "Путь компаний",
			path: "/companies/list/",
			want: "companies",
		},
		{
			name: "Путь с дополнительными параметрами",
			path: "/leads/list/?filter[pipeline_id]=3898873",
			want: "leads", // Лиды в URL и SDK
		},
		{
			name: "Неверный путь",
			path: "/some/invalid/path",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractEntityType(tt.path)
			if got != tt.want {
				t.Errorf("extractEntityType() = %v, want %v", got, tt.want)
			}
		})
	}
}
