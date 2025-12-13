.PHONY: build run test clean docker-build docker-run

# Сборка приложения
build:
	@echo "Building application..."
	go build -o bin/template-service cmd/template-service/main.go

# Запуск приложения
run:
	@echo "Running application..."
	go run cmd/template-service/main.go

# Запуск тестов
test:
	@echo "Running tests..."
	go test -v ./...

# Очистка бинарников
clean:
	@echo "Cleaning..."
	rm -rf bin/

# Сборка Docker образа
docker-build:
	@echo "Building Docker image..."
	docker build -t template-external-api-service:latest .

# Запуск в Docker
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 --env-file .env template-external-api-service:latest

# Форматирование кода
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Проверка линтером
lint:
	@echo "Running linter..."
	golangci-lint run

# Обновление зависимостей
deps:
	@echo "Updating dependencies..."
	go mod tidy
	go mod download

# TODO: Добавьте генерацию Swagger документации если используете
# swagger:
# 	@echo "Generating Swagger documentation..."
# 	swag init -g cmd/template-service/main.go --output docs

