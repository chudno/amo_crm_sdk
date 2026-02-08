# Интеграционные тесты

Тесты выполняют реальные запросы к API amoCRM. Требуется действующий аккаунт и токен доступа.

## Настройка

```bash
cp .env.example .env
```

Заполните `.env`:

```
AMO_BASE_URL=https://your-subdomain.amocrm.ru
AMO_ACCESS_TOKEN=your_long_lived_token
```

## Запуск

Все тесты:

```bash
AMO_BASE_URL=https://your-subdomain.amocrm.ru AMO_ACCESS_TOKEN=your_token \
  go test -tags integration -v ./integration/...
```

Или через `.env` (если используете утилиту вроде `direnv` или `dotenv`):

```bash
source .env
go test -tags integration -v ./integration/...
```

Конкретный тест:

```bash
go test -tags integration -v -run TestLeads ./integration/...
```

## Покрытие

| Файл | Сущность |
|---|---|
| `leads_test.go` | Сделки |
| `contacts_test.go` | Контакты |
| `companies_test.go` | Компании |
| `tasks_test.go` | Задачи |
| `pipelines_test.go` | Воронки |
| `users_test.go` | Пользователи |
| `notes_test.go` | Примечания |
| `tags_test.go` | Теги |
| `catalogs_test.go` | Каталоги |
| `catalog_elements_test.go` | Элементы каталогов |
| `calls_test.go` | Звонки |
| `events_test.go` | События |
| `helpers_test.go` | Вспомогательные функции |

## Важно

- Тесты создают и удаляют сущности в вашем аккаунте amoCRM
- Без переменных окружения тесты автоматически пропускаются (`t.Skip`)
- Таймаут на каждый запрос: 15 секунд
- Файл `.env` добавлен в `.gitignore`
