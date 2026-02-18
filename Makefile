.PHONY: help build run test clean docker-up docker-down migrate

# Показать помощь
help:
	@echo "Available commands:"
	@echo "  make build        - Build the application"
	@echo "  make run          - Run the application"
	@echo "  make test         - Run tests"
	@echo "  make test-cover   - Run tests with coverage"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make docker-up    - Start Docker services"
	@echo "  make docker-down  - Stop Docker services"
	@echo "  make migrate      - Run database migrations"
	@echo "  make lint         - Run linter"

# Собрать приложение
build:
	@echo "Building..."
	@go build -o bin/server cmd/server/main.go

# Запустить приложение
run:
	@echo "Running..."
	@go run cmd/server/main.go

# Запустить тесты
test:
	@echo "Running tests..."
	@go test -v ./...

# Запустить тесты с покрытием
test-cover:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Очистка
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

# Запустить Docker сервисы
docker-up:
	@echo "Starting Docker services..."
	@docker-compose up -d

# Остановить Docker сервисы
docker-down:
	@echo "Stopping Docker services..."
	@docker-compose down

# Применить миграции
migrate:
	@echo "Running migrations..."
	@./scripts/migrate.sh

# Линтер
lint:
	@echo "Running linter..."
	@golangci-lint run

# Полный цикл: docker -> migrate -> run
dev: docker-up migrate run

# Установка зависимостей
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy