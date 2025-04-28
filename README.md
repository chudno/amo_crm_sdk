# amoCRM SDK for Go

[![Go Report Card](https://goreportcard.com/badge/github.com/chudno/amo_crm_sdk)](https://goreportcard.com/report/github.com/chudno/amo_crm_sdk)
[![GoDoc](https://godoc.org/github.com/chudno/amo_crm_sdk?status.svg)](https://godoc.org/github.com/chudno/amo_crm_sdk)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Эта библиотека предоставляет SDK на языке Go для работы с API amoCRM. Она поддерживает все методы API, что даёт возможность полноценно взаимодействовать с сервисом amoCRM.

## Установка

```bash
go get github.com/chudno/amo_crm_sdk
```

## Особенности

* Поддержка OAuth 2.0 аутентификации
* **Поддержка долгоживущих токенов (Long-lived tokens)** для серверных интеграций
* Работа со всеми основными сущностями amoCRM (лиды, контакты, сделки, компании, задачи и т.д.)
* Поддержка пользовательских полей
* Поддержка вебхуков
* Полная документация на русском языке


## Запуск тестов

Для запуска тестов в SDK используется стандартный инструментарий Go. Вы можете запустить тесты следующими способами:

### Запуск всех тестов

```bash
go test ./...
```

Эта команда запустит все тесты во всех пакетах проекта.

### Запуск тестов для конкретного пакета

```bash
go test ./webhooks
```

Эта команда запустит только тесты в пакете webhooks.

### Запуск конкретного теста

```bash
go test -run TestCreateWebhook ./webhooks
```

Эта команда запустит только тест с именем TestCreateWebhook в пакете webhooks.

### Подробный вывод тестов

Для получения подробной информации о выполнении тестов используйте флаг -v:

```bash
go test -v ./...
```

### Анализ покрытия кода тестами

Для анализа покрытия кода тестами используйте флаг -cover:

```bash
go test -cover ./...
```

Для более детального анализа покрытия можно использовать:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

Это создаст HTML-отчет о покрытии кода тестами и откроет его в браузере.

## Структура проекта

- `entities/leads`: Методы обработки сущностей "Лиды".
- `entities/deals`: Методы обработки сущностей "Сделки".
- `entities/contacts`: Методы обработки сущностей "Контакты".
- `entities/companies`: Методы обработки сущностей "Компании".
- `entities/tasks`: Методы обработки сущностей "Задачи".
- `entities/notes`: Методы обработки сущностей "Примечания".
- `entities/pipelines`: Методы для работы с воронками и статусами.
- `entities/users`: Методы для работы с пользователями.
- `utils/custom_fields`: Методы обработки пользовательских полей.
- `utils/webhooks`: Методы для работы с вебхуками.
- `auth`: Методы для аутентификации в API amoCRM.
- `client`: Клиент для работы с API amoCRM.
- `utils`: Вспомогательные утилиты и общие функции.

Каждый пакет содержит функции для соответствующих методов API. Вся функциональность сопровождается комментариями на русском языке.

## Использование

```go
import (
    "github.com/chudno/amo_crm_sdk/auth"
    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/leads"
    "github.com/chudno/amo_crm_sdk/entities/deals"
    "github.com/chudno/amo_crm_sdk/entities/contacts"
    "github.com/chudno/amo_crm_sdk/entities/companies"
    "github.com/chudno/amo_crm_sdk/entities/tasks"
    "github.com/chudno/amo_crm_sdk/entities/notes"
    "github.com/chudno/amo_crm_sdk/entities/pipelines"
    "github.com/chudno/amo_crm_sdk/entities/users"
    "github.com/chudno/amo_crm_sdk/utils/webhooks"
    "github.com/chudno/amo_crm_sdk/utils/custom_fields"
)
```

## Аутентификация

Для работы с API amoCRM необходимо сначала получить токен доступа. SDK поддерживает OAuth 2.0 аутентификацию.

```go
package main

import (
    "fmt"
    "github.com/chudno/amo_crm_sdk/auth"
)

func main() {
    // Параметры для авторизации
    baseURL := "https://your_amocrm_domain.amocrm.ru"
    clientID := "ваш_client_id"
    clientSecret := "ваш_client_secret"
    redirectURI := "https://your-redirect-uri.com"
    code := "код_авторизации" // Получен после перехода пользователя по ссылке авторизации
    
    // Получение URL для авторизации пользователя
    authURL := auth.GetAuthURL(baseURL, clientID, redirectURI, "random_state", "popup")
    fmt.Println("Перейдите по ссылке для авторизации:", authURL)
    
    // После получения кода авторизации, получаем токен доступа
    authResponse, err := auth.GetAccessToken(baseURL, clientID, clientSecret, code, redirectURI)
    if err != nil {
        fmt.Println("Ошибка при получении токена:", err)
        return
    }
    
    fmt.Println("Токен доступа получен:", authResponse.AccessToken)
    fmt.Println("Токен обновления:", authResponse.RefreshToken)
    fmt.Println("Срок действия токена (в секундах):", authResponse.ExpiresIn)
    
    // Обновление токена доступа по refresh токену
    refreshedAuth, err := auth.RefreshAccessToken(baseURL, clientID, clientSecret, authResponse.RefreshToken)
    if err != nil {
        fmt.Println("Ошибка при обновлении токена:", err)
        return
    }
    
    fmt.Println("Новый токен доступа получен:", refreshedAuth.AccessToken)
}
```

## Работа с лидами

```go
package main

import (
    "fmt"
    "github.com/chudno/amo_crm_sdk/entities/leads"
    "github.com/chudno/amo_crm_sdk/client" // Подключение клиента для API запросов
)

func main() {
    apiClient := client.NewClient("https://your_amocrm_domain.amocrm.ru", "ваш_api_ключ")

    // Создать нового лида
    newLead := &leads.Lead{Name: "Новый Лид", Price: 3000}
    createdLead, err := leads.CreateLead(apiClient, newLead)
    if err != nil {
        fmt.Println("Ошибка при создании лида:", err)
    } else {
        fmt.Println("Лид создан:", createdLead)
    }

    // Получить существующего лида
    leadID := 12345 // ID существующего лида
    fetchedLead, err := leads.GetLead(apiClient, leadID)
    if err != nil {
        fmt.Println("Ошибка при получении лида:", err)
    } else {
        fmt.Println("Информация о лиде:", fetchedLead)
    }
    
    // Получить лид вместе со связанными контактами и компаниями
    leadWithRelations, err := leads.GetLead(apiClient, leadID, leads.WithContacts, leads.WithCompanies)
    if err != nil {
        fmt.Println("Ошибка при получении лида со связанными сущностями:", err)
    } else {
        fmt.Println("Информация о лиде:", leadWithRelations.Name)
        
        // Вывод информации о связанных контактах
        if leadWithRelations.Embedded != nil && len(leadWithRelations.Embedded.Contacts) > 0 {
            fmt.Println("Связанные контакты:")
            for _, contact := range leadWithRelations.Embedded.Contacts {
                fmt.Printf("ID: %d, Имя: %s\n", contact.ID, contact.Name)
            }
        }
        
        // Вывод информации о связанных компаниях
        if leadWithRelations.Embedded != nil && len(leadWithRelations.Embedded.Companies) > 0 {
            fmt.Println("Связанные компании:")
            for _, company := range leadWithRelations.Embedded.Companies {
                fmt.Printf("ID: %d, Название: %s\n", company.ID, company.Name)
            }
        }
    }
    
    // Получить список лидов со связанными сущностями
    leads, err := leads.GetLeads(apiClient, 1, 50, leads.WithContacts, leads.WithCompanies)
    if err != nil {
        fmt.Println("Ошибка при получении списка лидов:", err)
    } else {
        fmt.Printf("Получено %d лидов\n", len(leads))
        for _, lead := range leads {
            fmt.Printf("ID: %d, Название: %s\n", lead.ID, lead.Name)
            // Обработка связанных контактов и компаний...
        }
    }
    
    // Связать лид с контактом
    contactID := 67890 // ID существующего контакта
    err = leads.LinkLeadWithContact(apiClient, leadID, contactID)
    if err != nil {
        fmt.Println("Ошибка при связывании лида с контактом:", err)
    } else {
        fmt.Println("Лид успешно связан с контактом")
    }
    
    // Связать лид с компанией
    companyID := 54321 // ID существующей компании
    err = leads.LinkLeadWithCompany(apiClient, leadID, companyID)
    if err != nil {
        fmt.Println("Ошибка при связывании лида с компанией:", err)
    } else {
        fmt.Println("Лид успешно связан с компанией")
    }
    
    // Массовое удаление лидов
    leadIDs := []int{12345, 67890, 54321}
    cookie := "session_id=abcdef1234567890; user=user@example.com" // Cookie для аутентификации
    
    response, err := leads.DeleteLeads(apiClient, leadIDs, cookie)
    if err != nil {
        fmt.Println("Ошибка при удалении лидов:", err)
    } else {
        fmt.Println("Результат удаления лидов:", response.Message)
    }
}
```

## Работа с контактами

```go
package main

import (
    "fmt"
    "github.com/chudno/amo_crm_sdk/contacts"
    "github.com/chudno/amo_crm_sdk/client"
)

func main() {
    apiClient := client.NewClient("https://your_amocrm_domain.amocrm.ru", "ваш_api_ключ")
    
    // Создать новый контакт
    newContact := &contacts.Contact{Name: "Иван Иванов", Email: "ivan@example.com"}
    createdContact, err := contacts.CreateContact(apiClient, newContact)
    if err != nil {
        fmt.Println("Ошибка при создании контакта:", err)
    } else {
        fmt.Println("Контакт создан:", createdContact)
    }
    
    // Получить существующий контакт
    contactID := 12345 // ID существующего контакта
    fetchedContact, err := contacts.GetContact(apiClient, contactID)
    if err != nil {
        fmt.Println("Ошибка при получении контакта:", err)
    } else {
        fmt.Println("Информация о контакте:", fetchedContact)
    }
    
    // Массовое удаление контактов
    contactIDs := []int{12345, 67890, 54321}
    cookie := "session_id=abcdef1234567890; user=user@example.com" // Cookie для аутентификации
    
    response, err := contacts.DeleteContacts(apiClient, contactIDs, cookie)
    if err != nil {
        fmt.Println("Ошибка при удалении контактов:", err)
    } else {
        fmt.Println("Результат удаления контактов:", response.Message)
    }
}
```

## Работа с компаниями

```go
package main

import (
    "fmt"
    "github.com/chudno/amo_crm_sdk/companies"
    "github.com/chudno/amo_crm_sdk/client"
)

func main() {
    apiClient := client.NewClient("https://your_amocrm_domain.amocrm.ru", "ваш_api_ключ")
    
    // Создать новую компанию
    newCompany := &companies.Company{Name: "ООО Рога и Копыта"}
    createdCompany, err := companies.CreateCompany(apiClient, newCompany)
    if err != nil {
        fmt.Println("Ошибка при создании компании:", err)
    } else {
        fmt.Println("Компания создана:", createdCompany)
    }
    
    // Получить существующую компанию
    companyID := 12345 // ID существующей компании
    fetchedCompany, err := companies.GetCompany(apiClient, companyID)
    if err != nil {
        fmt.Println("Ошибка при получении компании:", err)
    } else {
        fmt.Println("Информация о компании:", fetchedCompany)
    }
    
    // Обновить компанию
    fetchedCompany.Name = "ООО Рога и Копыта Интернешнл"
    updatedCompany, err := companies.UpdateCompany(apiClient, fetchedCompany)
    if err != nil {
        fmt.Println("Ошибка при обновлении компании:", err)
    } else {
        fmt.Println("Компания обновлена:", updatedCompany)
    }
    
    // Получить список компаний
    companies, err := companies.ListCompanies(apiClient, 50, 1) // Лимит 50, страница 1
    if err != nil {
        fmt.Println("Ошибка при получении списка компаний:", err)
    } else {
        fmt.Printf("Получено %d компаний\n", len(companies))
        for _, company := range companies {
            fmt.Printf("ID: %d, Название: %s\n", company.ID, company.Name)
        }
    }
    
    // Массовое удаление компаний
    companyIDs := []int{12345, 67890, 54321}
    cookie := "session_id=abcdef1234567890; user=user@example.com" // Cookie для аутентификации
    
    response, err := companies.DeleteCompanies(apiClient, companyIDs, cookie)
    if err != nil {
        fmt.Println("Ошибка при удалении компаний:", err)
    } else {
        fmt.Println("Результат удаления компаний:", response.Message)
    }
}
```

## Работа с задачами

```go
package main

import (
    "fmt"
    "github.com/chudno/amo_crm_sdk/tasks"
    "github.com/chudno/amo_crm_sdk/client"
    "time"
)

func main() {
    apiClient := client.NewClient("https://your_amocrm_domain.amocrm.ru", "ваш_api_ключ")
    
    // Создать новую задачу
    completeTill := time.Now().Add(24 * time.Hour) // Срок выполнения через 24 часа
    newTask := &tasks.Task{
        Text: "Позвонить клиенту",
        CompleteTill: completeTill,
        EntityID: 12345, // ID связанной сущности (например, сделки)
        EntityType: "leads", // Тип сущности (leads, contacts, companies)
        TaskTypeID: 1, // ID типа задачи (1 - звонок, 2 - встреча и т.д.)
        ResponsibleUserID: 123, // ID ответственного пользователя
    }
    
    createdTask, err := tasks.CreateTask(apiClient, newTask)
    if err != nil {
        fmt.Println("Ошибка при создании задачи:", err)
    } else {
        fmt.Println("Задача создана:", createdTask)
    }
    
    // Получить существующую задачу
    taskID := 12345 // ID существующей задачи
    fetchedTask, err := tasks.GetTask(apiClient, taskID)
    if err != nil {
        fmt.Println("Ошибка при получении задачи:", err)
    } else {
        fmt.Println("Информация о задаче:", fetchedTask)
    }
    
    // Обновить задачу
    fetchedTask.Text = "Срочно позвонить клиенту"
    updatedTask, err := tasks.UpdateTask(apiClient, fetchedTask)
    if err != nil {
        fmt.Println("Ошибка при обновлении задачи:", err)
    } else {
        fmt.Println("Задача обновлена:", updatedTask)
    }
    
    // Отметить задачу как выполненную
    completedTask, err := tasks.CompleteTask(apiClient, taskID)
    if err != nil {
        fmt.Println("Ошибка при выполнении задачи:", err)
    } else {
        fmt.Println("Задача выполнена:", completedTask)
    }
    
    // Получить список задач с фильтрацией
    filter := map[string]string{
        "responsible_user_id": "123",
        "is_completed": "0", // Невыполненные задачи
    }
    
    tasksList, err := tasks.ListTasks(apiClient, 50, 1, filter) // Лимит 50, страница 1
    if err != nil {
        fmt.Println("Ошибка при получении списка задач:", err)
    } else {
        fmt.Printf("Получено %d задач\n", len(tasksList))
        for _, task := range tasksList {
            fmt.Printf("ID: %d, Текст: %s\n", task.ID, task.Text)
        }
    }
    
    // Массовое удаление задач
    taskIDs := []int{12345, 67890, 54321}
    cookie := "session_id=abcdef1234567890; user=user@example.com" // Cookie для аутентификации
    
    response, err := tasks.DeleteTasks(apiClient, taskIDs, cookie)
    if err != nil {
        fmt.Println("Ошибка при удалении задач:", err)
    } else {
        fmt.Println("Результат удаления задач:", response.Message)
    }
}
```

## Работа со связанными сущностями

### Получение сущностей со связанными данными

SDK поддерживает получение сущностей (сделок и лидов) вместе с их связанными контактами и компаниями. Для этого используются параметры `WithContacts` и `WithCompanies` в методах `GetDeal`, `GetDeals`, `GetLead` и `GetLeads`.

```go
package main

import (
    "fmt"
    "github.com/chudno/amo_crm_sdk/leads"
    "github.com/chudno/amo_crm_sdk/deals"
    "github.com/chudno/amo_crm_sdk/contacts"
    "github.com/chudno/amo_crm_sdk/companies"
    "github.com/chudno/amo_crm_sdk/client"
)

func main() {
    apiClient := client.NewClient("https://your_amocrm_domain.amocrm.ru", "ваш_api_ключ")
    
    // Получение лида вместе со связанными контактами и компаниями
    leadID := 12345 // ID существующего лида
    lead, err := leads.GetLead(apiClient, leadID, leads.WithContacts, leads.WithCompanies)
    if err != nil {
        fmt.Println("Ошибка при получении лида со связанными сущностями:", err)
    } else {
        fmt.Println("Информация о лиде:", lead.Name)
        
        // Вывод информации о связанных контактах
        if lead.Embedded != nil && len(lead.Embedded.Contacts) > 0 {
            fmt.Println("Связанные контакты:")
            for _, contact := range lead.Embedded.Contacts {
                fmt.Printf("ID: %d, Имя: %s\n", contact.ID, contact.Name)
            }
        }
        
        // Вывод информации о связанных компаниях
        if lead.Embedded != nil && len(lead.Embedded.Companies) > 0 {
            fmt.Println("Связанные компании:")
            for _, company := range lead.Embedded.Companies {
                fmt.Printf("ID: %d, Название: %s\n", company.ID, company.Name)
            }
        }
    }
    
    // Получение сделки вместе со связанными контактами и компаниями
    dealID := 67890 // ID существующей сделки
    deal, err := deals.GetDeal(apiClient, dealID, deals.WithContacts, deals.WithCompanies)
    if err != nil {
        fmt.Println("Ошибка при получении сделки со связанными сущностями:", err)
    } else {
        fmt.Println("Информация о сделке:", deal.Name)
        
        // Вывод информации о связанных контактах
        if deal.Embedded != nil && len(deal.Embedded.Contacts) > 0 {
            fmt.Println("Связанные контакты:")
            for _, contact := range deal.Embedded.Contacts {
                fmt.Printf("ID: %d, Имя: %s\n", contact.ID, contact.Name)
            }
        }
        
        // Вывод информации о связанных компаниях
        if deal.Embedded != nil && len(deal.Embedded.Companies) > 0 {
            fmt.Println("Связанные компании:")
            for _, company := range deal.Embedded.Companies {
                fmt.Printf("ID: %d, Название: %s\n", company.ID, company.Name)
            }
        }
    }
    
    // Получение списка сделок со связанными сущностями и фильтрацией
    filter := map[string]string{
        "status_id": "142", // Фильтр по ID статуса
        "created_at": "1609459200", // Фильтр по дате создания (timestamp)
    }
    
    dealsList, err := deals.GetDeals(apiClient, 1, 50, filter, deals.WithContacts, deals.WithCompanies)
    if err != nil {
        fmt.Println("Ошибка при получении списка сделок:", err)
    } else {
        fmt.Printf("Получено %d сделок\n", len(dealsList))
        for _, deal := range dealsList {
            fmt.Printf("ID: %d, Название: %s\n", deal.ID, deal.Name)
            // Обработка связанных контактов и компаний...
        }
    }
}
```

### Связывание сущностей

SDK предоставляет методы для связывания сделок и лидов с контактами и компаниями.

```go
package main

import (
    "fmt"
    "github.com/chudno/amo_crm_sdk/leads"
    "github.com/chudno/amo_crm_sdk/deals"
    "github.com/chudno/amo_crm_sdk/contacts"
    "github.com/chudno/amo_crm_sdk/companies"
    "github.com/chudno/amo_crm_sdk/client"
)

func main() {
    apiClient := client.NewClient("https://your_amocrm_domain.amocrm.ru", "ваш_api_ключ")
    
    // Связывание лида с контактом
    leadID := 12345 // ID существующего лида
    contactID := 67890 // ID существующего контакта
    err := leads.LinkLeadWithContact(apiClient, leadID, contactID)
    if err != nil {
        fmt.Println("Ошибка при связывании лида с контактом:", err)
    } else {
        fmt.Println("Лид успешно связан с контактом")
    }
    
    // Связывание лида с компанией
    companyID := 54321 // ID существующей компании
    err = leads.LinkLeadWithCompany(apiClient, leadID, companyID)
    if err != nil {
        fmt.Println("Ошибка при связывании лида с компанией:", err)
    } else {
        fmt.Println("Лид успешно связан с компанией")
    }
    
    // Связывание сделки с контактом
    dealID := 98765 // ID существующей сделки
    err = deals.LinkDealWithContact(apiClient, dealID, contactID)
    if err != nil {
        fmt.Println("Ошибка при связывании сделки с контактом:", err)
    } else {
        fmt.Println("Сделка успешно связана с контактом")
    }
    
    // Связывание сделки с компанией
    err = deals.LinkDealWithCompany(apiClient, dealID, companyID)
    if err != nil {
        fmt.Println("Ошибка при связывании сделки с компанией:", err)
    } else {
        fmt.Println("Сделка успешно связана с компанией")
    }
    
    // Создание лида с привязкой к контакту и компании
    newContact := &contacts.Contact{ID: contactID}
    newCompany := &companies.Company{ID: companyID}
    
    newLead := &leads.Lead{
        Name: "Новый лид с привязками",
        Price: 5000,
        Embedded: &leads.LeadEmbedded{
            Contacts: []contacts.Contact{*newContact},
            Companies: []companies.Company{*newCompany},
        },
    }
    
    createdLead, err := leads.CreateLead(apiClient, newLead)
    if err != nil {
        fmt.Println("Ошибка при создании лида с привязками:", err)
    } else {
        fmt.Println("Создан лид с привязками:", createdLead.ID)
    }
    
    // Создание сделки с привязкой к контакту и компании
    newDeal := &deals.Deal{
        Name: "Новая сделка с привязками",
        Value: 10000,
        Embedded: &deals.DealEmbedded{
            Contacts: []deals.Contact{{ID: contactID}},
            Companies: []deals.Company{{ID: companyID}},
        },
    }
    
    createdDeal, err := deals.CreateDeal(apiClient, newDeal)
    if err != nil {
        fmt.Println("Ошибка при создании сделки с привязками:", err)
    } else {
        fmt.Println("Создана сделка с привязками:", createdDeal.ID)
    }
}
```

## Работа с примечаниями

```go
package main

import (
    "fmt"
    "github.com/chudno/amo_crm_sdk/notes"
    "github.com/chudno/amo_crm_sdk/client"
)

func main() {
    apiClient := client.NewClient("https://your_amocrm_domain.amocrm.ru", "ваш_api_ключ")
    
    // Создать новое примечание для сделки
    entityType := "leads" // Тип сущности (leads, contacts, companies, customers)
    entityID := 12345 // ID сущности
    
    newNote := &notes.Note{
        EntityID: entityID,
        EntityType: entityType,
        NoteType: 4, // Тип примечания (4 - обычное примечание)
        Params: notes.NoteParams{
            Text: "Клиент запросил дополнительную информацию",
        },
    }
    
    createdNote, err := notes.CreateNote(apiClient, entityType, entityID, newNote)
    if err != nil {
        fmt.Println("Ошибка при создании примечания:", err)
    } else {
        fmt.Println("Примечание создано:", createdNote)
    }
    
    // Получить существующее примечание
    noteID := 12345 // ID существующего примечания
    fetchedNote, err := notes.GetNote(apiClient, entityType, entityID, noteID)
    if err != nil {
        fmt.Println("Ошибка при получении примечания:", err)
    } else {
        fmt.Println("Информация о примечании:", fetchedNote)
    }
    
    // Обновить примечание
    fetchedNote.Params.Text = "Клиент запросил дополнительную информацию о продукте X"
    updatedNote, err := notes.UpdateNote(apiClient, entityType, entityID, fetchedNote)
    if err != nil {
        fmt.Println("Ошибка при обновлении примечания:", err)
    } else {
        fmt.Println("Примечание обновлено:", updatedNote)
    }
    
    // Получить список примечаний для сущности
    notesList, err := notes.ListNotes(apiClient, entityType, entityID, 50, 1) // Лимит 50, страница 1
    if err != nil {
        fmt.Println("Ошибка при получении списка примечаний:", err)
    } else {
        fmt.Printf("Получено %d примечаний\n", len(notesList))
        for _, note := range notesList {
            fmt.Printf("ID: %d, Тип: %d\n", note.ID, note.NoteType)
        }
    }
    
    // Удалить примечание
    err = notes.DeleteNote(apiClient, entityType, entityID, noteID)
    if err != nil {
        fmt.Println("Ошибка при удалении примечания:", err)
    } else {
        fmt.Println("Примечание успешно удалено")
    }
}
```

## Работа с пользовательскими полями

```go
package main

import (
    "fmt"
    "github.com/chudno/amo_crm_sdk/custom_fields"
    "github.com/chudno/amo_crm_sdk/client"
)

func main() {
    apiClient := client.NewClient("https://your_amocrm_domain.amocrm.ru", "ваш_api_ключ")
    
    // Создать новое пользовательское поле для контактов
    entityType := "contacts" // Тип сущности (leads, contacts, companies, customers)
    
    newField := &custom_fields.CustomField{
        Name: "Дата рождения",
        Type: "date", // Тип поля (text, numeric, checkbox, select, multiselect, date, url, textarea, radiobutton, streetaddress, smart_address, birthday, date_time, price)
        EntityType: entityType,
        Sort: 100,
        IsRequired: false,
        IsMultiple: false,
        IsVisible: true,
    }
    
    createdField, err := custom_fields.CreateCustomField(apiClient, entityType, newField)
    if err != nil {
        fmt.Println("Ошибка при создании пользовательского поля:", err)
    } else {
        fmt.Println("Пользовательское поле создано:", createdField)
    }
    
    // Создать поле типа "список" с вариантами
    selectField := &custom_fields.CustomField{
        Name: "Источник",
        Type: "select",
        EntityType: entityType,
        Sort: 101,
        IsRequired: true,
        IsMultiple: false,
        IsVisible: true,
        Enums: []custom_fields.Enum{
            {Value: "Реклама", Sort: 1},
            {Value: "Рекомендация", Sort: 2},
            {Value: "Холодный звонок", Sort: 3},
        },
    }
    
    createdSelectField, err := custom_fields.CreateCustomField(apiClient, entityType, selectField)
    if err != nil {
        fmt.Println("Ошибка при создании поля-списка:", err)
    } else {
        fmt.Println("Поле-список создано:", createdSelectField)
    }
    
    // Получить существующее пользовательское поле
    fieldID := 12345 // ID существующего поля
    fetchedField, err := custom_fields.GetCustomField(apiClient, entityType, fieldID)
    if err != nil {
        fmt.Println("Ошибка при получении пользовательского поля:", err)
    } else {
        fmt.Println("Информация о пользовательском поле:", fetchedField)
    }
    
    // Обновить пользовательское поле
    fetchedField.Name = "Дата рождения клиента"
    updatedField, err := custom_fields.UpdateCustomField(apiClient, entityType, fetchedField)
    if err != nil {
        fmt.Println("Ошибка при обновлении пользовательского поля:", err)
    } else {
        fmt.Println("Пользовательское поле обновлено:", updatedField)
    }
    
    // Получить список пользовательских полей для сущности
    fieldsList, err := custom_fields.ListCustomFields(apiClient, entityType)
    if err != nil {
        fmt.Println("Ошибка при получении списка пользовательских полей:", err)
    } else {
        fmt.Printf("Получено %d пользовательских полей\n", len(fieldsList))
        for _, field := range fieldsList {
            fmt.Printf("ID: %d, Название: %s, Тип: %s\n", field.ID, field.Name, field.Type)
        }
    }
    
    // Удалить пользовательское поле
    err = custom_fields.DeleteCustomField(apiClient, entityType, fieldID)
    if err != nil {
        fmt.Println("Ошибка при удалении пользовательского поля:", err)
    } else {
        fmt.Println("Пользовательское поле успешно удалено")
    }
}
```

## Работа с пользователями

```go
package main

import (
    "fmt"
    "github.com/chudno/amo_crm_sdk/entities/users"
    "github.com/chudno/amo_crm_sdk/client"
)

func main() {
    apiClient := client.NewClient("https://your_amocrm_domain.amocrm.ru", "ваш_api_ключ")
    
    // Получить информацию о текущем пользователе (владельце API-ключа)
    currentUser, err := users.GetCurrentUser(apiClient)
    if err != nil {
        fmt.Println("Ошибка при получении информации о текущем пользователе:", err)
    } else {
        fmt.Println("Информация о текущем пользователе:")
        fmt.Printf("ID: %d, Имя: %s, Email: %s\n", currentUser.ID, currentUser.Name, currentUser.Email)
        fmt.Printf("Права администратора: %v\n", currentUser.Rights.IsAdmin)
    }
    
    // Получить информацию о конкретном пользователе
    userID := 12345 // ID существующего пользователя
    user, err := users.GetUser(apiClient, userID)
    if err != nil {
        fmt.Println("Ошибка при получении информации о пользователе:", err)
    } else {
        fmt.Println("Информация о пользователе:")
        fmt.Printf("ID: %d, Имя: %s, Email: %s\n", user.ID, user.Name, user.Email)
        fmt.Printf("Активен: %v\n", user.IsActive)
    }
    
    // Получить список пользователей
    usersList, err := users.ListUsers(apiClient, 50, 1) // Лимит 50, страница 1
    if err != nil {
        fmt.Println("Ошибка при получении списка пользователей:", err)
    } else {
        fmt.Printf("Получено %d пользователей\n", len(usersList))
        for _, user := range usersList {
            fmt.Printf("ID: %d, Имя: %s, Email: %s\n", user.ID, user.Name, user.Email)
        }
    }
}
```

## Работа с воронками и статусами

```go
package main

import (
    "fmt"
    "github.com/chudno/amo_crm_sdk/pipelines"
    "github.com/chudno/amo_crm_sdk/client"
)

func main() {
    apiClient := client.NewClient("https://your_amocrm_domain.amocrm.ru", "ваш_api_ключ")
    
    // Создать новую воронку
    newPipeline := &pipelines.Pipeline{
        Name: "Новая воронка продаж",
        Sort: 100,
        IsMain: false,
        IsActive: true,
    }
    
    createdPipeline, err := pipelines.CreatePipeline(apiClient, newPipeline)
    if err != nil {
        fmt.Println("Ошибка при создании воронки:", err)
    } else {
        fmt.Println("Воронка создана:", createdPipeline)
    }
    
    // Получить существующую воронку
    pipelineID := 12345 // ID существующей воронки
    fetchedPipeline, err := pipelines.GetPipeline(apiClient, pipelineID)
    if err != nil {
        fmt.Println("Ошибка при получении воронки:", err)
    } else {
        fmt.Println("Информация о воронке:", fetchedPipeline)
    }
    
    // Обновить воронку
    fetchedPipeline.Name = "Обновленная воронка продаж"
    updatedPipeline, err := pipelines.UpdatePipeline(apiClient, fetchedPipeline)
    if err != nil {
        fmt.Println("Ошибка при обновлении воронки:", err)
    } else {
        fmt.Println("Воронка обновлена:", updatedPipeline)
    }
    
    // Получить список воронок
    pipelinesList, err := pipelines.ListPipelines(apiClient)
    if err != nil {
        fmt.Println("Ошибка при получении списка воронок:", err)
    } else {
        fmt.Printf("Получено %d воронок\n", len(pipelinesList))
        for _, pipeline := range pipelinesList {
            fmt.Printf("ID: %d, Название: %s\n", pipeline.ID, pipeline.Name)
        }
    }
    
    // Создать новый статус в воронке
    newStatus := &pipelines.Status{
        Name: "Новый статус",
        Sort: 100,
        Color: "#FFFF00", // Желтый цвет
        Type: 0, // Обычный статус
    }
    
    createdStatus, err := pipelines.CreateStatus(apiClient, pipelineID, newStatus)
    if err != nil {
        fmt.Println("Ошибка при создании статуса:", err)
    } else {
        fmt.Println("Статус создан:", createdStatus)
    }
    
    // Получить существующий статус
    statusID := 12345 // ID существующего статуса
    fetchedStatus, err := pipelines.GetStatus(apiClient, pipelineID, statusID)
    if err != nil {
        fmt.Println("Ошибка при получении статуса:", err)
    } else {
        fmt.Println("Информация о статусе:", fetchedStatus)
    }
    
    // Обновить статус
    fetchedStatus.Name = "Обновленный статус"
    fetchedStatus.Color = "#00FF00" // Зеленый цвет
    updatedStatus, err := pipelines.UpdateStatus(apiClient, pipelineID, fetchedStatus)
    if err != nil {
        fmt.Println("Ошибка при обновлении статуса:", err)
    } else {
        fmt.Println("Статус обновлен:", updatedStatus)
    }
    
    // Удалить статус
    err = pipelines.DeleteStatus(apiClient, pipelineID, statusID)
    if err != nil {
        fmt.Println("Ошибка при удалении статуса:", err)
    } else {
        fmt.Println("Статус успешно удален")
    }
    
    // Удалить воронку
    err = pipelines.DeletePipeline(apiClient, pipelineID)
    if err != nil {
        fmt.Println("Ошибка при удалении воронки:", err)
    } else {
        fmt.Println("Воронка успешно удалена")
    }
}
```

## Работа с вебхуками

```go
package main

import (
    "fmt"
    "github.com/chudno/amo_crm_sdk/webhooks"
    "github.com/chudno/amo_crm_sdk/client"
)

func main() {
    apiClient := client.NewClient("https://your_amocrm_domain.amocrm.ru", "ваш_api_ключ")
    
    // Создать новый вебхук
    newWebhook := &webhooks.Webhook{
        Destination: "https://your-webhook-handler.com/endpoint",
        Settings: webhooks.Settings{
            Entities: []string{"lead", "contact", "company"},
            Actions: []string{"add", "update", "delete"},
        },
    }
    
    createdWebhook, err := webhooks.CreateWebhook(apiClient, newWebhook)
    if err != nil {
        fmt.Println("Ошибка при создании вебхука:", err)
    } else {
        fmt.Println("Вебхук создан:", createdWebhook)
    }
    
    // Получить список вебхуков
    webhooksList, err := webhooks.ListWebhooks(apiClient)
    if err != nil {
        fmt.Println("Ошибка при получении списка вебхуков:", err)
    } else {
        fmt.Printf("Получено %d вебхуков\n", len(webhooksList))
        for _, webhook := range webhooksList {
            fmt.Printf("ID: %d, URL: %s\n", webhook.ID, webhook.Destination)
        }
    }
    
    // Удалить вебхук
    webhookID := 12345 // ID существующего вебхука
    err = webhooks.DeleteWebhook(apiClient, webhookID)
    if err != nil {
        fmt.Println("Ошибка при удалении вебхука:", err)
    } else {
        fmt.Println("Вебхук успешно удален")
    }
}
## Быстрый старт

### Инициализация клиента

```go
import (
    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/auth"
)

// Инициализация с токеном доступа
apiClient := client.NewClient("https://your-domain.amocrm.ru", "your_access_token")

// ИЛИ получение токена через OAuth
baseURL := "https://your-domain.amocrm.ru"
clientID := "your_client_id"
clientSecret := "your_client_secret"
redirectURI := "your_redirect_uri"
code := "authorization_code" // полученный после авторизации пользователя

auth, err := auth.GetAccessToken(baseURL, clientID, clientSecret, code, redirectURI)
if err != nil {
    // обработка ошибки
}

apiClient := client.NewClient(baseURL, auth.AccessToken)
```

### Примеры использования

#### Контакты

```go
import "github.com/chudno/amo_crm_sdk/entities/contacts"

// Создание контакта
newContact := &contacts.Contact{
    Name: "Иван Иванов",
    ResponsibleUserID: 12345,
}
createdContact, err := contacts.CreateContact(apiClient, newContact)

// Получение контакта
contact, err := contacts.GetContact(apiClient, contactID)

// Список контактов
contactsList, err := contacts.GetContacts(apiClient, 1, 50)
```

#### Сделки

```go
import "github.com/chudno/amo_crm_sdk/entities/deals"

// Создание сделки
newDeal := &deals.Deal{
    Name: "Новая сделка",
    Value: 10000,
    StatusID: 12345,
    PipelineID: 67890,
}
createdDeal, err := deals.CreateDeal(apiClient, newDeal)

// Изменение статуса сделки
dealToUpdate := &deals.Deal{
    ID: dealID,
    StatusID: newStatusID,
}
updatedDeal, err := deals.UpdateDeal(apiClient, dealToUpdate)
```

#### Пользовательские поля

```go
import "github.com/chudno/amo_crm_sdk/utils/custom_fields"

// Добавление пользовательского поля к контакту
contact.CustomFieldsValues = []custom_fields.CustomFieldValue{
    {
        FieldID: 12345,
        Values: []custom_fields.FieldValue{
            {Value: "Значение поля"},
        },
    },
}
```

## Запуск тестов и проверок в Docker

Для удобства тестирования и обеспечения единообразной среды разработки, в проекте настроено тестирование в Docker.

### С помощью Docker Compose

```bash
# Запуск всех тестов
docker-compose run --rm test

# Запуск только линтера
docker-compose run --rm lint

# Форматирование кода
docker-compose run --rm fmt
```

### С помощью Makefile

В проекте доступен Makefile для упрощения запуска команд:

```bash
# Список всех доступных команд
make help

# Запуск тестов в Docker
make docker-test

# Запуск линтера в Docker
make docker-lint

# Запуск всех проверок в Docker
make docker-all
```

Использование Docker гарантирует, что тесты будут запущены в одинаковой среде независимо от вашей локальной конфигурации.
