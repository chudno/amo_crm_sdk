# Файлы (Files)

Модуль для работы с файлами в amoCRM.

## Оглавление

- [Возможности](#возможности)
- [Типы сущностей](#типы-сущностей)
- [Структура файла](#структура-файла)
- [Примеры использования](#примеры-использования)
  - [Загрузка файла из локального файла](#загрузка-файла-из-локального-файла)
  - [Загрузка файла по содержимому](#загрузка-файла-по-содержимому)
  - [Получение списка файлов](#получение-списка-файлов)
  - [Получение информации о файле](#получение-информации-о-файле)
  - [Удаление файла](#удаление-файла)
  - [Массовое удаление файлов](#массовое-удаление-файлов)
  - [Скачивание файла](#скачивание-файла)
  - [Получение URL для скачивания файла](#получение-url-для-скачивания-файла)

## Возможности

- Загрузка файлов из локальных файлов
- Загрузка файлов из памяти (по содержимому)
- Получение списка файлов, прикрепленных к сущности
- Получение информации о конкретном файле
- Удаление файлов
- Массовое удаление файлов
- Скачивание файлов
- Получение URL для скачивания файлов

## Типы сущностей

Файлы в amoCRM можно прикреплять к различным типам сущностей:

```go
const (
    // EntityTypeLead тип сущности - Сделка
    EntityTypeLead EntityType = "leads"
    // EntityTypeContact тип сущности - Контакт
    EntityTypeContact EntityType = "contacts"
    // EntityTypeCompany тип сущности - Компания
    EntityTypeCompany EntityType = "companies"
    // EntityTypeCustomers тип сущности - Покупатель
    EntityTypeCustomers EntityType = "customers"
    // EntityTypeCatalogElements тип сущности - Элемент каталога
    EntityTypeCatalogElements EntityType = "catalog_elements"
)
```

## Структура файла

Структура `File` представляет файл в amoCRM:

```go
// File представляет структуру файла в amoCRM
type File struct {
    ID          int         `json:"id"`
    UUID        string      `json:"uuid,omitempty"`
    EntityID    int         `json:"entity_id"`
    EntityType  EntityType  `json:"entity_type"`
    CreatedBy   int         `json:"created_by,omitempty"`
    UpdatedBy   int         `json:"updated_by,omitempty"`
    CreatedAt   int64       `json:"created_at,omitempty"`
    UpdatedAt   int64       `json:"updated_at,omitempty"`
    Size        int         `json:"size,omitempty"`
    Name        string      `json:"name,omitempty"`
    Type        string      `json:"type,omitempty"`
    Version     int         `json:"version,omitempty"`
    AccountID   int         `json:"account_id,omitempty"`
    Title       string      `json:"title,omitempty"`
    URL         string      `json:"url,omitempty"`
    Download    string      `json:"download_link,omitempty"`
    Preview     string      `json:"preview,omitempty"`
    Links       FileLinks   `json:"_links,omitempty"`
}
```

## Примеры использования

### Загрузка файла из локального файла

```go
package main

import (
    "fmt"
    "log"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/files"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // Параметры загрузки
    entityType := files.EntityTypeLead // Тип сущности - сделка
    entityID := 123                    // ID сделки
    filePath := "/path/to/document.pdf" // Путь к файлу для загрузки

    // Загружаем файл
    file, err := files.UploadFile(apiClient, entityType, entityID, filePath)
    if err != nil {
        log.Fatalf("Ошибка при загрузке файла: %v", err)
    }

    // Выводим информацию о загруженном файле
    fmt.Printf("Файл успешно загружен:\n")
    fmt.Printf("ID: %d\n", file.ID)
    fmt.Printf("Название: %s\n", file.Name)
    fmt.Printf("Размер: %d байт\n", file.Size)
    fmt.Printf("URL: %s\n", file.URL)
    fmt.Printf("Ссылка для скачивания: %s\n", file.Download)
}
```

### Загрузка файла по содержимому

```go
package main

import (
    "fmt"
    "io/ioutil"
    "log"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/files"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // Параметры загрузки
    entityType := files.EntityTypeContact // Тип сущности - контакт
    entityID := 456                      // ID контакта
    fileName := "document.pdf"           // Имя файла
    
    // Читаем содержимое файла (в реальном приложении это может быть файл из памяти или другого источника)
    content, err := ioutil.ReadFile("/path/to/document.pdf")
    if err != nil {
        log.Fatalf("Ошибка при чтении файла: %v", err)
    }

    // Загружаем файл по содержимому
    file, err := files.UploadFileByContent(apiClient, entityType, entityID, fileName, content)
    if err != nil {
        log.Fatalf("Ошибка при загрузке файла: %v", err)
    }

    // Выводим информацию о загруженном файле
    fmt.Printf("Файл успешно загружен:\n")
    fmt.Printf("ID: %d\n", file.ID)
    fmt.Printf("Название: %s\n", file.Name)
    fmt.Printf("Размер: %d байт\n", file.Size)
    fmt.Printf("URL: %s\n", file.URL)
    fmt.Printf("Ссылка для скачивания: %s\n", file.Download)
}
```

### Получение списка файлов

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/files"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // Параметры запроса
    entityType := files.EntityTypeLead // Тип сущности - сделка
    entityID := 123                    // ID сделки
    page := 1                          // Номер страницы
    limit := 50                        // Лимит файлов на странице

    // Получаем список файлов
    filesList, err := files.GetFiles(apiClient, entityType, entityID, page, limit)
    if err != nil {
        log.Fatalf("Ошибка при получении списка файлов: %v", err)
    }

    // Выводим информацию о файлах
    fmt.Printf("Получено %d файлов:\n", len(filesList))
    for i, file := range filesList {
        fmt.Printf("%d. %s (ID: %d)\n", i+1, file.Name, file.ID)
        fmt.Printf("   Размер: %d байт\n", file.Size)
        fmt.Printf("   Создан: %s\n", time.Unix(file.CreatedAt, 0).Format("02.01.2006 15:04:05"))
        fmt.Printf("   URL: %s\n", file.URL)
        fmt.Printf("   Ссылка для скачивания: %s\n", file.Download)
        fmt.Println()
    }
}
```

### Получение информации о файле

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/files"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // Параметры запроса
    entityType := files.EntityTypeLead // Тип сущности - сделка
    entityID := 123                    // ID сделки
    fileID := 456                      // ID файла

    // Получаем информацию о файле
    file, err := files.GetFile(apiClient, entityType, entityID, fileID)
    if err != nil {
        log.Fatalf("Ошибка при получении информации о файле: %v", err)
    }

    // Выводим информацию о файле
    fmt.Printf("Информация о файле:\n")
    fmt.Printf("ID: %d\n", file.ID)
    fmt.Printf("Название: %s\n", file.Name)
    fmt.Printf("Размер: %d байт\n", file.Size)
    fmt.Printf("Создан: %s\n", time.Unix(file.CreatedAt, 0).Format("02.01.2006 15:04:05"))
    fmt.Printf("Обновлен: %s\n", time.Unix(file.UpdatedAt, 0).Format("02.01.2006 15:04:05"))
    fmt.Printf("URL: %s\n", file.URL)
    fmt.Printf("Ссылка для скачивания: %s\n", file.Download)
}
```

### Удаление файла

```go
package main

import (
    "fmt"
    "log"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/files"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // Параметры запроса
    entityType := files.EntityTypeLead // Тип сущности - сделка
    entityID := 123                    // ID сделки
    fileID := 456                      // ID файла для удаления

    // Удаляем файл
    err := files.DeleteFile(apiClient, entityType, entityID, fileID)
    if err != nil {
        log.Fatalf("Ошибка при удалении файла: %v", err)
    }

    fmt.Printf("Файл с ID %d успешно удален\n", fileID)
}
```

### Массовое удаление файлов

```go
package main

import (
    "fmt"
    "log"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/files"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // Параметры запроса
    entityType := files.EntityTypeLead // Тип сущности - сделка
    entityID := 123                    // ID сделки
    fileIDs := []int{456, 789, 1234}   // ID файлов для удаления

    // Удаляем файлы
    err := files.BatchDeleteFiles(apiClient, entityType, entityID, fileIDs)
    if err != nil {
        log.Fatalf("Ошибка при массовом удалении файлов: %v", err)
    }

    fmt.Printf("Файлы с ID %v успешно удалены\n", fileIDs)
}
```

### Скачивание файла

```go
package main

import (
    "fmt"
    "log"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/files"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // Параметры запроса
    entityType := files.EntityTypeLead // Тип сущности - сделка
    entityID := 123                    // ID сделки
    fileID := 456                      // ID файла для скачивания
    savePath := "/path/to/save/file.pdf" // Путь для сохранения файла

    // Скачиваем файл
    err := files.DownloadFile(apiClient, entityType, entityID, fileID, savePath)
    if err != nil {
        log.Fatalf("Ошибка при скачивании файла: %v", err)
    }

    fmt.Printf("Файл успешно скачан и сохранен в %s\n", savePath)
}
```

### Получение URL для скачивания файла

```go
package main

import (
    "fmt"
    "log"

    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/files"
)

func main() {
    // Создаем клиент API
    apiClient := client.NewClient("https://example.amocrm.ru", "TOKEN")

    // Параметры запроса
    entityType := files.EntityTypeLead // Тип сущности - сделка
    entityID := 123                    // ID сделки
    fileID := 456                      // ID файла

    // Получаем URL для скачивания файла
    downloadURL, err := files.GetDownloadFileURL(apiClient, entityType, entityID, fileID)
    if err != nil {
        log.Fatalf("Ошибка при получении URL для скачивания файла: %v", err)
    }

    fmt.Printf("URL для скачивания файла: %s\n", downloadURL)
}
```
