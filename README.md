# amoCRM SDK for Go

[![Go Report Card](https://goreportcard.com/badge/github.com/chudno/amo_crm_sdk)](https://goreportcard.com/report/github.com/chudno/amo_crm_sdk)
[![GoDoc](https://godoc.org/github.com/chudno/amo_crm_sdk?status.svg)](https://godoc.org/github.com/chudno/amo_crm_sdk)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![codecov](https://codecov.io/gh/chudno/amo_crm_sdk/branch/main/graph/badge.svg)](https://codecov.io/gh/chudno/amo_crm_sdk)
[![GitHub Actions](https://github.com/chudno/amo_crm_sdk/workflows/Go/badge.svg)](https://github.com/chudno/amo_crm_sdk/actions)
[![CodeQL](https://github.com/chudno/amo_crm_sdk/workflows/CodeQL/badge.svg)](https://github.com/chudno/amo_crm_sdk/security/code-scanning)
[![Repo Size](https://img.shields.io/github/repo-size/chudno/amo_crm_sdk)](https://github.com/chudno/amo_crm_sdk)
[![Last Commit](https://img.shields.io/github/last-commit/chudno/amo_crm_sdk)](https://github.com/chudno/amo_crm_sdk/commits/main)

Эта библиотека предоставляет SDK на языке Go для работы с API amoCRM. Она позволяет полноценно взаимодействовать с amoCRM и работать со всеми типами сущностей.

## Оглавление

- [Особенности](#особенности)
- [Установка](#установка)
- [Документация](#документация)
- [Тестирование](#тестирование)

## Особенности

* Поддержка OAuth 2.0 аутентификации
* **Поддержка долгоживущих токенов (Long-lived tokens)** для серверных интеграций
* Работа со всеми основными сущностями amoCRM (лиды, контакты, сделки, компании, задачи и т.д.)
* Поддержка пользовательских полей
* Поддержка вебхуков
* Полная документация на русском языке

## Установка

```bash
go get github.com/chudno/amo_crm_sdk
```

## Документация

SDK разделен на модули, каждый из которых отвечает за работу с определенной сущностью в amoCRM:

### Основные модули

| Модуль | Описание | Документация |
|-------|-------------|--------------|
| `auth` | Аутентификация в API amoCRM | [Подробнее](./auth/README.md) |
| `client` | Клиент для работы с API | [Подробнее](./client/README.md) |

### Сущности

| Модуль | Описание | Документация |
|-------|-------------|--------------|
| `entities/leads` | Работа с лидами | [Подробнее](./entities/leads/README.md) |
| `entities/contacts` | Работа с контактами | [Подробнее](./entities/contacts/README.md) |
| `entities/companies` | Работа с компаниями | [Подробнее](./entities/companies/README.md) |
| `entities/tasks` | Работа с задачами | [Подробнее](./entities/tasks/README.md) |
| `entities/notes` | Работа с примечаниями | [Подробнее](./entities/notes/README.md) |
| `entities/pipelines` | Работа с воронками и статусами | [Подробнее](./entities/pipelines/README.md) |
| `entities/users` | Работа с пользователями | [Подробнее](./entities/users/README.md) |
| `entities/tags` | Работа с тегами | [Подробнее](./entities/tags/README.md) |
| `entities/catalogs` | Работа с каталогами | [Подробнее](./entities/catalogs/README.md) |
| `entities/catalog_elements` | Работа с элементами пользовательских каталогов | [Подробнее](./entities/catalog_elements/README.md) |
| `entities/unsorted` | Работа с неразобранными заявками | [Подробнее](./entities/unsorted/README.md) |
| `entities/files` | Работа с файлами | [Подробнее](./entities/files/README.md) |
| `entities/calls` | Работа со звонками | [Подробнее](./entities/calls/README.md) |
| `entities/events` | Работа с событиями | [Подробнее](./entities/events/README.md) |
| `entities/segments` | Работа с сегментами | [Подробнее](./entities/segments/README.md) |
| `entities/widgets` | Работа с виджетами | [Подробнее](./entities/widgets/README.md) |
| `entities/access_rights` | Работа с правами доступа | [Подробнее](./entities/access_rights/README.md) |
| `entities/short_links` | Работа с короткими ссылками | [Подробнее](./entities/short_links/README.md) |
| `entities/mailing` | Работа с email-рассылками | [Подробнее](./entities/mailing/README.md) |
| `entities/sources` | Работа с источниками лидов | [Подробнее](./entities/sources/README.md) |

### Утилиты

| Модуль | Описание | Документация |
|-------|-------------|--------------|
| `utils/custom_fields` | Работа с пользовательскими полями | [Подробнее](./utils/custom_fields/README.md) |
| `utils/webhooks` | Работа с вебхуками | [Подробнее](./utils/webhooks/README.md) |
| `utils/urlfilters` | Конвертация URL-фильтров из веб-интерфейса в SDK | [Подробнее](./utils/urlfilters/README.md) |

### Примеры использования

Примеры использования всех модулей приведены в соответствующих README.md файлах каждого модуля.

## Тестирование

Для запуска тестов в проекте используется Docker. Все команды доступны через Makefile:

```bash
# Запуск всех проверок
make all

# Только тесты
make test

# Проверка кода с помощью go vet
make lint

# Форматирование кода
make fmt

# Проверка цикломатической сложности
make cyclo
```
