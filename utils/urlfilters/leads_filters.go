// Пакет urlfilters предоставляет функции для конвертации URL-фильтров из веб-интерфейса amoCRM
// в формат, используемый в SDK для запросов к API.
package urlfilters

import (
	"fmt"
	"strconv"
)

// LeadFilter представляет собой результат парсинга URL с фильтрами для лидов
type LeadFilter struct {
	// Исходный ParsedFilter с общей информацией
	*ParsedFilter
	// PageInt - номер страницы (целочисленное значение)
	PageInt int
	// LimitInt - ограничение на количество элементов (целочисленное значение)
	LimitInt int
}

// NewLeadFilterFromURL создает объект LeadFilter из URL amoCRM
func NewLeadFilterFromURL(rawURL string) (*LeadFilter, error) {
	// Парсим URL с использованием основного парсера
	parsedFilter, err := ParseLeadURL(rawURL)
	if err != nil {
		return nil, err
	}

	// Конвертируем строковые значения page и limit в целые числа
	pageInt, err := strconv.Atoi(parsedFilter.Page)
	if err != nil {
		return nil, fmt.Errorf("ошибка при конвертации page в число: %w", err)
	}

	limitInt, err := strconv.Atoi(parsedFilter.Limit)
	if err != nil {
		return nil, fmt.Errorf("ошибка при конвертации limit в число: %w", err)
	}

	return &LeadFilter{
		ParsedFilter: parsedFilter,
		PageInt:      pageInt,
		LimitInt:     limitInt,
	}, nil
}

// GetSDKFilterMap возвращает карту фильтров для использования в функции GetLeads SDK
func (lf *LeadFilter) GetSDKFilterMap() map[string]string {
	return lf.Filter
}
