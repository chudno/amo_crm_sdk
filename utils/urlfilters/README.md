# Утилита конвертации URL-фильтров amoCRM в SDK

Пакет `urlfilters` предоставляет инструменты для конвертации URL-адресов с фильтрами из веб-интерфейса amoCRM в формат, используемый в SDK для запросов к API.

**Важно:** В веб-интерфейсе amoCRM сделки называются "лидами" (и в URL используется путь `/leads/`). В SDK мы также используем термин "лиды" и пакет `leads`.

## Основные возможности

- Парсинг URL-адресов из веб-интерфейса amoCRM
- Извлечение параметров фильтрации, пагинации и сортировки
- Конвертация в формат, который можно напрямую использовать в SDK
- Специализированные утилиты для работы с различными типами сущностей (сделки, контакты)

## Установка

```bash
go get github.com/chudno/amo_crm_sdk/utils/urlfilters
```

## Пример использования

### Парсинг URL-фильтров для лидов (leads)

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/chudno/amo_crm_sdk/client"
	"github.com/chudno/amo_crm_sdk/entities/leads"
	"github.com/chudno/amo_crm_sdk/utils/urlfilters"
)

func main() {
	// URL из веб-интерфейса amoCRM с фильтрами
	amoURL := "https://example.amocrm.ru/leads/list/?filter[name]=Тестовый+лид&filter[status][]=10073462&page=1&limit=50"
	
	// Парсим URL и получаем фильтры для лидов
	filter, err := urlfilters.NewLeadFilterFromURL(amoURL)
	if err != nil {
		log.Fatalf("Ошибка при парсинге URL: %v", err)
	}
	
	// Создаем клиент API
	apiClient := client.NewClient("https://example.amocrm.ru", "your_access_token")
	
	// Получаем лиды с фильтрами из URL
	leadsList, err := leads.GetLeads(
		apiClient, 
		filter.PageInt, 
		filter.LimitInt, 
		filter.GetSDKFilterMap(),
	)
	if err != nil {
		log.Fatalf("Ошибка при получении лидов: %v", err)
	}
	
	fmt.Printf("Получено %d лидов\n", len(leadsList))
	for _, lead := range leadsList {
		fmt.Printf("Лид: %s (ID: %d)\n", lead.Name, lead.ID)
	}
}
```

## Доступные методы

### Общие методы

- `ParseURL(rawURL string) (*ParsedFilter, error)` - парсит любой URL amoCRM и извлекает фильтры
- `ParseLeadURL(rawURL string) (*ParsedFilter, error)` - парсит URL лидов amoCRM

### Методы для лидов (leads)

- `NewLeadFilterFromURL(rawURL string) (*LeadFilter, error)` - создает объект фильтра для лидов
- `GetSDKFilterMap()` - получает карту фильтров для использования в SDK
- `Example()` - генерирует пример кода для использования фильтра



## Структура данных

### ParsedFilter

```go
type ParsedFilter struct {
	Filter     map[string]string // карта фильтров для SDK
	EntityType string            // тип сущности (leads, contacts, etc.)
	Page       string            // номер страницы
	Limit      string            // ограничение на количество элементов
	RawQuery   string            // исходная строка запроса
}
```

### LeadFilter

```go
type LeadFilter struct {
	*ParsedFilter              // встроенный ParsedFilter
	PageInt     int            // номер страницы (целочисленное значение)
	LimitInt    int            // ограничение на количество элементов (целочисленное значение)
}
```


