# Элементы каталогов (catalog_elements)

Модуль для работы с элементами каталогов в amoCRM.

## Оглавление

- [Возможности](#возможности)
- [Структура элемента каталога](#структура-элемента-каталога)
- [Примеры использования](#примеры-использования)
  - [Получение списка элементов каталога](#получение-списка-элементов-каталога)
  - [Создание элемента каталога](#создание-элемента-каталога)
  - [Получение элемента каталога по ID](#получение-элемента-каталога-по-id)
  - [Обновление элемента каталога](#обновление-элемента-каталога)
  - [Удаление элемента каталога](#удаление-элемента-каталога)
  - [Работа с тегами элемента каталога](#работа-с-тегами-элемента-каталога)

## Возможности

- Получение списка элементов каталога с фильтрацией
- Создание одного или нескольких элементов каталога
- Получение конкретного элемента каталога по ID
- Обновление элемента каталога
- Массовое обновление элементов каталога
- Удаление элемента каталога
- Массовое удаление элементов каталога
- Работа с тегами элемента каталога (привязка, получение)

## Структура элемента каталога

```go
// CatalogElement представляет собой структуру элемента каталога в amoCRM.
type CatalogElement struct {
	ID                 int                       `json:"id,omitempty"`
	Name               string                    `json:"name"`
	CreatedBy          int                       `json:"created_by,omitempty"`
	UpdatedBy          int                       `json:"updated_by,omitempty"`
	CreatedAt          int64                     `json:"created_at,omitempty"`
	UpdatedAt          int64                     `json:"updated_at,omitempty"`
	CatalogID          int                       `json:"catalog_id"`
	CustomFieldsValues []CustomFieldValue        `json:"custom_fields_values,omitempty"`
	AccountID          int                       `json:"account_id,omitempty"`
	IsDeleted          bool                      `json:"is_deleted,omitempty"`
	QuantityBounded    bool                      `json:"quantity_bounded,omitempty"`
	QuantityRaw        int                       `json:"quantity_raw,omitempty"`
	Embedded           *CatalogElementEmbedded   `json:"_embedded,omitempty"`
	Links              *CatalogElementLinks      `json:"_links,omitempty"`
}
```

## Примеры использования

### Получение списка элементов каталога

Получение списка элементов каталога с возможностью фильтрации:

```go
package main

import (
	"fmt"
	"log"

	"github.com/chudno/amo_crm_sdk/client"
	"github.com/chudno/amo_crm_sdk/entities/catalog_elements"
)

func main() {
	// Создаем клиент API
	apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

	// Получаем список элементов каталога
	catalogID := 123 // ID каталога
	page := 1
	limit := 50
	filter := map[string]string{
		"filter[name]": "Продукт", // Фильтр по имени
	}

	// Получаем элементы каталога с фильтрацией и с включением тегов
	elements, err := catalog_elements.GetCatalogElements(apiClient, catalogID, page, limit, filter, catalog_elements.WithTags)
	if err != nil {
		log.Fatalf("Ошибка при получении элементов каталога: %v", err)
	}

	// Выводим результаты
	fmt.Printf("Получено %d элементов каталога\n", len(elements))
	for _, element := range elements {
		fmt.Printf("ID: %d, Название: %s\n", element.ID, element.Name)
		
		// Вывод пользовательских полей
		if len(element.CustomFieldsValues) > 0 {
			fmt.Println("Пользовательские поля:")
			for _, field := range element.CustomFieldsValues {
				if len(field.Values) > 0 {
					fmt.Printf("  %s: %v\n", field.FieldName, field.Values[0].Value)
				}
			}
		}
		
		// Вывод тегов
		if element.Embedded != nil && len(element.Embedded.Tags) > 0 {
			fmt.Println("Теги:")
			for _, tag := range element.Embedded.Tags {
				fmt.Printf("  %s (ID: %d)\n", tag.Name, tag.ID)
			}
		}
		
		fmt.Println("---")
	}
}
```

### Создание элемента каталога

Создание нового элемента каталога:

```go
package main

import (
	"fmt"
	"log"

	"github.com/chudno/amo_crm_sdk/client"
	"github.com/chudno/amo_crm_sdk/entities/catalog_elements"
)

func main() {
	// Создаем клиент API
	apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

	// ID каталога, в который добавляем элемент
	catalogID := 123

	// Создаем элемент каталога
	newElement := &catalog_elements.CatalogElement{
		Name: "Новый товар",
		CustomFieldsValues: []catalog_elements.CustomFieldValue{
			{
				FieldID: 789, // ID поля "Код товара"
				Values: []catalog_elements.FieldValueItem{
					{
						Value: "PRD-001",
					},
				},
			},
			{
				FieldID: 790, // ID поля "Цена"
				Values: []catalog_elements.FieldValueItem{
					{
						Value: 1500,
					},
				},
			},
		},
	}

	// Отправляем запрос на создание элемента
	createdElement, err := catalog_elements.CreateCatalogElement(apiClient, catalogID, newElement)
	if err != nil {
		log.Fatalf("Ошибка при создании элемента каталога: %v", err)
	}

	// Выводим информацию о созданном элементе
	fmt.Printf("Элемент каталога успешно создан. ID: %d, Название: %s\n", createdElement.ID, createdElement.Name)
}
```

Создание нескольких элементов каталога за один запрос:

```go
package main

import (
	"fmt"
	"log"

	"github.com/chudno/amo_crm_sdk/client"
	"github.com/chudno/amo_crm_sdk/entities/catalog_elements"
)

func main() {
	// Создаем клиент API
	apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

	// ID каталога, в который добавляем элементы
	catalogID := 123

	// Создаем список элементов каталога
	elements := []catalog_elements.CatalogElement{
		{
			Name: "Товар 1",
			CustomFieldsValues: []catalog_elements.CustomFieldValue{
				{
					FieldID: 789, // ID поля "Код товара"
					Values: []catalog_elements.FieldValueItem{
						{
							Value: "PRD-001",
						},
					},
				},
			},
		},
		{
			Name: "Товар 2",
			CustomFieldsValues: []catalog_elements.CustomFieldValue{
				{
					FieldID: 789, // ID поля "Код товара"
					Values: []catalog_elements.FieldValueItem{
						{
							Value: "PRD-002",
						},
					},
				},
			},
		},
	}

	// Отправляем запрос на создание элементов
	createdElements, err := catalog_elements.CreateCatalogElements(apiClient, catalogID, elements)
	if err != nil {
		log.Fatalf("Ошибка при создании элементов каталога: %v", err)
	}

	// Выводим информацию о созданных элементах
	fmt.Printf("Создано элементов каталога: %d\n", len(createdElements))
	for _, element := range createdElements {
		fmt.Printf("ID: %d, Название: %s\n", element.ID, element.Name)
	}
}
```

### Получение элемента каталога по ID

Получение информации о конкретном элементе каталога:

```go
package main

import (
	"fmt"
	"log"

	"github.com/chudno/amo_crm_sdk/client"
	"github.com/chudno/amo_crm_sdk/entities/catalog_elements"
)

func main() {
	// Создаем клиент API
	apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

	// ID каталога и элемента каталога
	catalogID := 123
	elementID := 456

	// Получаем элемент каталога с тегами
	element, err := catalog_elements.GetCatalogElement(apiClient, catalogID, elementID, catalog_elements.WithTags)
	if err != nil {
		log.Fatalf("Ошибка при получении элемента каталога: %v", err)
	}

	// Выводим информацию об элементе
	fmt.Printf("Название: %s\n", element.Name)
	fmt.Printf("ID: %d\n", element.ID)
	fmt.Printf("Создан: %d\n", element.CreatedAt)
	fmt.Printf("Обновлен: %d\n", element.UpdatedAt)

	// Вывод пользовательских полей
	if len(element.CustomFieldsValues) > 0 {
		fmt.Println("Пользовательские поля:")
		for _, field := range element.CustomFieldsValues {
			if len(field.Values) > 0 {
				fmt.Printf("  %s: %v\n", field.FieldName, field.Values[0].Value)
			}
		}
	}

	// Вывод тегов
	if element.Embedded != nil && len(element.Embedded.Tags) > 0 {
		fmt.Println("Теги:")
		for _, tag := range element.Embedded.Tags {
			fmt.Printf("  %s (ID: %d)\n", tag.Name, tag.ID)
		}
	}
}
```

### Обновление элемента каталога

Обновление существующего элемента каталога:

```go
package main

import (
	"fmt"
	"log"

	"github.com/chudno/amo_crm_sdk/client"
	"github.com/chudno/amo_crm_sdk/entities/catalog_elements"
)

func main() {
	// Создаем клиент API
	apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

	// ID каталога и элемента каталога
	catalogID := 123
	elementID := 456

	// Создаем структуру для обновления элемента
	elementToUpdate := &catalog_elements.CatalogElement{
		ID:   elementID,
		Name: "Обновленное название товара",
		CustomFieldsValues: []catalog_elements.CustomFieldValue{
			{
				FieldID: 789, // ID поля "Код товара"
				Values: []catalog_elements.FieldValueItem{
					{
						Value: "PRD-001-UPD",
					},
				},
			},
			{
				FieldID: 790, // ID поля "Цена"
				Values: []catalog_elements.FieldValueItem{
					{
						Value: 1800,
					},
				},
			},
		},
	}

	// Отправляем запрос на обновление элемента
	updatedElement, err := catalog_elements.UpdateCatalogElement(apiClient, catalogID, elementToUpdate)
	if err != nil {
		log.Fatalf("Ошибка при обновлении элемента каталога: %v", err)
	}

	// Выводим информацию об обновленном элементе
	fmt.Printf("Элемент каталога успешно обновлен. ID: %d, Название: %s\n", updatedElement.ID, updatedElement.Name)
}
```

Массовое обновление элементов каталога:

```go
package main

import (
	"fmt"
	"log"

	"github.com/chudno/amo_crm_sdk/client"
	"github.com/chudno/amo_crm_sdk/entities/catalog_elements"
)

func main() {
	// Создаем клиент API
	apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

	// ID каталога
	catalogID := 123

	// Создаем список элементов для обновления
	elementsToUpdate := []catalog_elements.CatalogElement{
		{
			ID:   456, // ID первого элемента
			Name: "Обновленный товар 1",
			CustomFieldsValues: []catalog_elements.CustomFieldValue{
				{
					FieldID: 789, // ID поля "Код товара"
					Values: []catalog_elements.FieldValueItem{
						{
							Value: "PRD-001-UPD",
						},
					},
				},
			},
		},
		{
			ID:   789, // ID второго элемента
			Name: "Обновленный товар 2",
			CustomFieldsValues: []catalog_elements.CustomFieldValue{
				{
					FieldID: 789, // ID поля "Код товара"
					Values: []catalog_elements.FieldValueItem{
						{
							Value: "PRD-002-UPD",
						},
					},
				},
			},
		},
	}

	// Отправляем запрос на массовое обновление элементов
	updatedElements, err := catalog_elements.UpdateCatalogElements(apiClient, catalogID, elementsToUpdate)
	if err != nil {
		log.Fatalf("Ошибка при массовом обновлении элементов каталога: %v", err)
	}

	// Выводим информацию об обновленных элементах
	fmt.Printf("Обновлено элементов каталога: %d\n", len(updatedElements))
	for _, element := range updatedElements {
		fmt.Printf("ID: %d, Название: %s\n", element.ID, element.Name)
	}
}
```

### Удаление элемента каталога

Удаление элемента каталога по ID:

```go
package main

import (
	"fmt"
	"log"

	"github.com/chudno/amo_crm_sdk/client"
	"github.com/chudno/amo_crm_sdk/entities/catalog_elements"
)

func main() {
	// Создаем клиент API
	apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

	// ID каталога и элемента каталога
	catalogID := 123
	elementID := 456

	// Удаляем элемент каталога
	err := catalog_elements.DeleteCatalogElement(apiClient, catalogID, elementID)
	if err != nil {
		log.Fatalf("Ошибка при удалении элемента каталога: %v", err)
	}

	fmt.Printf("Элемент каталога с ID %d успешно удален\n", elementID)
}
```

Массовое удаление элементов каталога:

```go
package main

import (
	"fmt"
	"log"

	"github.com/chudno/amo_crm_sdk/client"
	"github.com/chudno/amo_crm_sdk/entities/catalog_elements"
)

func main() {
	// Создаем клиент API
	apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

	// ID каталога
	catalogID := 123

	// ID элементов каталога для удаления
	elementIDs := []int{456, 789, 1234}

	// Удаляем несколько элементов каталога
	err := catalog_elements.BatchDeleteCatalogElements(apiClient, catalogID, elementIDs)
	if err != nil {
		log.Fatalf("Ошибка при массовом удалении элементов каталога: %v", err)
	}

	fmt.Printf("Элементы каталога с ID %v успешно удалены\n", elementIDs)
}
```

### Работа с тегами элемента каталога

Связывание элемента каталога с тегами:

```go
package main

import (
	"fmt"
	"log"

	"github.com/chudno/amo_crm_sdk/client"
	"github.com/chudno/amo_crm_sdk/entities/catalog_elements"
)

func main() {
	// Создаем клиент API
	apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

	// ID каталога и элемента каталога
	catalogID := 123
	elementID := 456

	// Создаем теги для связывания с элементом
	tags := []catalog_elements.Tag{
		{
			Name:  "Новинка",
			Color: "#FF0000",
		},
		{
			Name:  "Распродажа",
			Color: "#00FF00",
		},
	}

	// Связываем элемент с тегами
	err := catalog_elements.LinkCatalogElementWithTags(apiClient, catalogID, elementID, tags)
	if err != nil {
		log.Fatalf("Ошибка при связывании элемента каталога с тегами: %v", err)
	}

	fmt.Printf("Элемент каталога с ID %d успешно связан с тегами\n", elementID)
}
```

Получение тегов элемента каталога:

```go
package main

import (
	"fmt"
	"log"

	"github.com/chudno/amo_crm_sdk/client"
	"github.com/chudno/amo_crm_sdk/entities/catalog_elements"
)

func main() {
	// Создаем клиент API
	apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

	// ID каталога и элемента каталога
	catalogID := 123
	elementID := 456

	// Получаем теги элемента каталога
	tags, err := catalog_elements.GetCatalogElementTags(apiClient, catalogID, elementID)
	if err != nil {
		log.Fatalf("Ошибка при получении тегов элемента каталога: %v", err)
	}

	// Выводим информацию о тегах
	fmt.Printf("Получено %d тегов для элемента каталога с ID %d:\n", len(tags), elementID)
	for _, tag := range tags {
		fmt.Printf("- %s (ID: %d, Цвет: %s)\n", tag.Name, tag.ID, tag.Color)
	}
}
```
