.PHONY: test lint fmt cyclo all coverage coverage-html

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

# Отчет о покрытии кода тестами
coverage:
	@echo "Запуск проверки покрытия кода тестами..."
	@echo "mode: atomic" > coverage.txt
	@for pkg in $$(go list ./... | grep -v -E 'examples|mocks'); do \
		go test -coverprofile=profile.out -covermode=atomic $$pkg || exit 1; \
		if [ -f profile.out ]; then \
			tail -n +2 profile.out >> coverage.txt; \
			rm profile.out; \
		fi; \
	done
	@go tool cover -func=coverage.txt
	@echo "Отчет сохранен в файл coverage.txt"

# Генерация HTML-отчета о покрытии
coverage-html: coverage
	@echo "Генерация HTML-отчета о покрытии кода тестами..."
	@go tool cover -html=coverage.txt -o coverage.html
	@echo "Отчет сохранен в файл coverage.html"

# Помощь
help:
	@echo "Доступные команды:"
	@echo "  make test         - Запуск тестов"
	@echo "  make lint         - Проверка кода с помощью go vet"
	@echo "  make fmt          - Форматирование кода"
	@echo "  make cyclo        - Проверка цикломатической сложности"
	@echo "  make all          - Запуск всех проверок"
	@echo "  make coverage     - Сформировать отчет о покрытии кода тестами"
	@echo "  make coverage-html - Сформировать HTML-отчет о покрытии кода тестами"
	@echo "  make help         - Показать эту справку"
