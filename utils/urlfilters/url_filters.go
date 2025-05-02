// Пакет urlfilters предоставляет функции для конвертации URL-фильтров из веб-интерфейса amoCRM
// в формат, используемый в SDK для запросов к API.
package urlfilters

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// ParsedFilter представляет собой результат парсинга URL с фильтрами
type ParsedFilter struct {
	// Filter - карта фильтров для использования в функциях SDK
	Filter map[string]string
	// EntityType - тип сущности (leads, contacts, etc.)
	EntityType string
	// Page - номер страницы
	Page string
	// Limit - ограничение на количество элементов
	Limit string
	// RawQuery - исходная строка запроса
	RawQuery string
}

// ParseURL разбирает URL из веб-интерфейса amoCRM и возвращает структуру с фильтрами,
// которые можно использовать в SDK.
func ParseURL(rawURL string) (*ParsedFilter, error) {
	// Парсим URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка при парсинге URL: %w", err)
	}

	// Извлекаем тип сущности из пути
	entityType := extractEntityType(parsedURL.Path)
	if entityType == "" {
		return nil, fmt.Errorf("не удалось определить тип сущности из URL: %s", parsedURL.Path)
	}

	// Парсим query-параметры
	queryParams, err := url.ParseQuery(parsedURL.RawQuery)
	if err != nil {
		return nil, fmt.Errorf("ошибка при парсинге параметров запроса: %w", err)
	}

	// Создаем результат
	result := &ParsedFilter{
		Filter:     make(map[string]string),
		EntityType: entityType,
		Page:       "1",  // По умолчанию первая страница
		Limit:      "50", // По умолчанию 50 элементов
		RawQuery:   parsedURL.RawQuery,
	}

	// Извлекаем страницу и лимит, если они указаны
	if page := queryParams.Get("page"); page != "" {
		result.Page = page
	}
	if limit := queryParams.Get("limit"); limit != "" {
		result.Limit = limit
	}

	// Обрабатываем фильтры
	for key, values := range queryParams {
		if strings.HasPrefix(key, "filter") {
			// В случае с множественными значениями для одного фильтра, берем первое
			if len(values) > 0 {
				result.Filter[key] = values[0]
			}
		}
	}

	return result, nil
}

// ParseLeadURL парсит URL лидов amoCRM и возвращает фильтры
func ParseLeadURL(rawURL string) (*ParsedFilter, error) {
	// Выполняем базовый парсинг URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка при парсинге URL: %w", err)
	}

	// Вручную проверяем, что это URL лидов - ищем "/leads/list/" в пути
	reLeads := regexp.MustCompile(`/leads/list/?`)
	if !reLeads.MatchString(parsedURL.Path) {
		return nil, fmt.Errorf("УРЛ не относится к лидам")
	}

	// Теперь, когда мы точно знаем, что это URL лидов, используем обычный парсер
	parsedFilter, err := ParseURL(rawURL)
	if err != nil {
		return nil, err
	}

	// Если все прошло успешно, то entityType должен быть "leads" благодаря маппингу
	if parsedFilter.EntityType != "leads" {
		return nil, fmt.Errorf("Внутренняя ошибка: тип сущности не был распознан как leads")
	}

	return parsedFilter, nil
}

// AmoCRMEntityTypeMap содержит маппинг между типами сущностей в URL amoCRM и соответствующими модулями SDK
var AmoCRMEntityTypeMap = map[string]string{
	"leads":     "leads",     // Теперь лиды в URL и SDK называются одинаково - "leads"
	"contacts":  "contacts",  // Контакты имеют одинаковые названия в URL и SDK
	"customers": "customers", // Клиенты
	"companies": "companies", // Компании
	"catalogs":  "catalogs",  // Каталоги
	"tasks":     "tasks",     // Задачи
}

// extractEntityType извлекает тип сущности из пути URL и преобразует его в формат SDK
func extractEntityType(path string) string {
	// Ищем шаблоны типа "/leads/list/", "/contacts/list/" и т.д.
	re := regexp.MustCompile(`/([a-z_]+)/list/?`)
	matches := re.FindStringSubmatch(path)
	if len(matches) >= 2 {
		// Получаем тип сущности из URL
		entityTypeInURL := matches[1]

		// Преобразуем тип сущности из URL в тип сущности в SDK
		if sdkType, exists := AmoCRMEntityTypeMap[entityTypeInURL]; exists {
			return sdkType
		}

		// Если маппинг не найден, возвращаем оригинальное значение
		return entityTypeInURL
	}
	return ""
}
