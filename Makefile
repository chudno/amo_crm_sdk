.PHONY: test lint fmt docker-test docker-lint docker-fmt

# Локальные команды
test:
	go test -v ./...

lint:
	go vet ./...

fmt:
	go fmt ./...

all: fmt lint test

# Docker-команды
docker-test:
	docker-compose run --rm test

docker-lint:
	docker-compose run --rm lint

docker-fmt:
	docker-compose run --rm fmt

docker-all:
	docker-compose run --rm test && docker-compose run --rm lint && docker-compose run --rm fmt

# Помощь
help:
	@echo "Доступные команды:"
	@echo "  make test         - Запуск тестов локально"
	@echo "  make lint         - Проверка кода с помощью go vet локально"
	@echo "  make fmt          - Форматирование кода локально"
	@echo "  make all          - Запуск всех локальных проверок"
	@echo "  make docker-test  - Запуск тестов в Docker"
	@echo "  make docker-lint  - Проверка кода в Docker"
	@echo "  make docker-fmt   - Форматирование кода в Docker"
	@echo "  make docker-all   - Запуск всех проверок в Docker"
