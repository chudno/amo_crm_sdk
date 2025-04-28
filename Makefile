.PHONY: test lint fmt cyclo all

test:
	docker-compose run --rm test

lint:
	docker-compose run --rm lint

cyclo:
	docker-compose run --rm cyclo

fmt:
	docker-compose run --rm fmt

all:
	docker-compose run --rm test && docker-compose run --rm lint && docker-compose run --rm fmt && docker-compose run --rm cyclo

# Помощь
help:
	@echo "Доступные команды:"
	@echo "  make test         - Запуск тестов"
	@echo "  make lint         - Проверка кода с помощью go vet"
	@echo "  make fmt          - Форматирование кода"
	@echo "  make cyclo        - Проверка цикломатической сложности"
	@echo "  make all          - Запуск всех проверок"
	@echo "  make help         - Показать эту справку"
