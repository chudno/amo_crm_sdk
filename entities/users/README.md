# Модуль Пользователи

Модуль `users` предоставляет функциональность для работы с пользователями (сотрудниками) в amoCRM.

## Содержание

- [Основные функции](#основные-функции)
- [Получение пользователя](#получение-пользователя)
- [Получение списка пользователей](#получение-списка-пользователей)
- [Поиск пользователей](#поиск-пользователей)
- [Назначение ответственных](#назначение-ответственных)

## Основные функции

| Функция | Описание |
|---------|----------|
| `GetUser` | Получение пользователя по ID |
| `GetUsers` | Получение списка пользователей с фильтрацией |
| `GetCurrentUser` | Получение информации о текущем пользователе |

## Получение пользователя

```go
import (
    "github.com/chudno/amo_crm_sdk/client"
    "github.com/chudno/amo_crm_sdk/entities/users"
)

// Инициализация клиента
apiClient := client.NewClient("https://your-domain.amocrm.ru", "your_access_token")

// Получение пользователя по ID
userID := 12345
user, err := users.GetUser(apiClient, userID)
if err != nil {
    // Обработка ошибки
}

// Вывод информации о пользователе
fmt.Printf("Имя: %s\n", user.Name)
fmt.Printf("Email: %s\n", user.Email)
fmt.Printf("Роль: %s\n", user.Rights)
```

## Получение списка пользователей

```go
// Получение всех пользователей
usersList, err := users.GetUsers(apiClient)
if err != nil {
    // Обработка ошибки
}

// Вывод списка пользователей
for _, user := range usersList {
    fmt.Printf("ID: %d, Имя: %s, Email: %s\n", user.ID, user.Name, user.Email)
}

// Фильтрация пользователей
filter := map[string]string{
    "with": "role,group",  // Получить информацию о ролях и группах
    "page_size": "50",     // Количество пользователей на страницу
}
filteredUsers, err := users.GetUsers(apiClient, filter)
```

## Поиск пользователей

```go
// Поиск пользователей по имени
query := "Иван"
filter := map[string]string{
    "query": query,
}
foundUsers, err := users.GetUsers(apiClient, filter)
if err != nil {
    // Обработка ошибки
}

// Вывод найденных пользователей
fmt.Printf("Найдено %d пользователей по запросу '%s':\n", len(foundUsers), query)
for _, user := range foundUsers {
    fmt.Printf("ID: %d, Имя: %s\n", user.ID, user.Name)
}
```

## Назначение ответственных

Для назначения ответственного пользователя за сущность, используйте поле `ResponsibleUserID` в соответствующих структурах:

```go
// Назначение ответственного за контакт
import "github.com/chudno/amo_crm_sdk/entities/contacts"

contact := &contacts.Contact{
    Name: "Иван Иванов",
    ResponsibleUserID: 12345, // ID ответственного менеджера
}

// Назначение ответственного за лид
import "github.com/chudno/amo_crm_sdk/entities/leads"

lead := &leads.Lead{
    Name:              "Тестовый лид",
    ResponsibleUserID: 12345, // ID ответственного менеджера
}

// Назначение ответственного за задачу
import "github.com/chudno/amo_crm_sdk/entities/tasks"

task := &tasks.Task{
    Text: "Позвонить клиенту",
    ResponsibleUserID: 12345, // ID ответственного менеджера
}
```

Для смены ответственного у существующей сущности:

```go
// Получение сущности
contact, err := contacts.GetContact(apiClient, contactID)
if err != nil {
    // Обработка ошибки
}

// Смена ответственного
contact.ResponsibleUserID = 67890 // ID нового ответственного

// Сохранение изменений
updatedContact, err := contacts.UpdateContact(apiClient, contact)
if err != nil {
    // Обработка ошибки
}
```
