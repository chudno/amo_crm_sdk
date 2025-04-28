FROM golang:1.20-alpine

# Устанавливаем необходимые инструменты и зависимости
RUN apk add --no-cache git make

# Устанавливаем golangci-lint для расширенной проверки кода
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2

# Устанавливаем gocyclo для проверки цикломатической сложности
RUN go install github.com/fzipp/gocyclo/cmd/gocyclo@latest

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем только go.mod файл
COPY go.mod ./
RUN go mod download && go mod tidy

# Копируем исходный код
COPY . .

# Команда, выполняемая по умолчанию
CMD ["sh", "-c", "go test -v ./... && go vet ./... && golangci-lint run ./... && gocyclo -over 15 ."]
