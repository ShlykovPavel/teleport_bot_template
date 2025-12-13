# Template External API Service

Шаблонный Go-микросервис для работы с внешним API с JWT авторизацией.

## 📋 Описание

Это готовый к использованию шаблон микросервиса на Go, который включает:

- ✅ **JWT авторизацию** для внешнего API (автоматическое обновление токенов)
- ✅ **HTTP клиент** с поддержкой всех REST методов
- ✅ **External API Service** с готовыми методами (GetDemandInfo, GetAccountInfo, UpdateDemandStatus)
- ✅ **Prometheus метрики** (http_requests_total, http_request_duration_seconds)
- ✅ **Middleware** (Metrics, CORS, Recovery)
- ✅ **Структурированное логирование** (slog)
- ✅ **Конфигурацию** через YAML + переменные окружения
- ✅ **MongoDB** подключение (опционально)
- ✅ **Graceful shutdown**
- ✅ **Docker** поддержку
- ✅ **Healthcheck** и **Metrics** endpoints

## 🏗 Архитектура

```
template-external-api-service/
├── cmd/template-service/       # Точка входа
│   └── main.go
├── internal/
│   ├── app/                    # Инициализация приложения
│   │   └── app.go
│   ├── client/                 # HTTP клиент с JWT
│   │   ├── http_client.go
│   │   ├── response.go
│   │   └── external_api_service/
│   │       ├── external_api_service.go
│   │       └── models/
│   ├── config/                 # Конфигурация
│   │   └── config.go
│   ├── server/
│   │   └── middlewares/        # Middleware
│   │       └── middlewares.go
│   ├── services/               # Бизнес-логика (TODO)
│   └── storage/
│       └── database/           # БД слой
│           ├── db.go
│           ├── db_errors/
│           └── repositories/
├── metrics/                    # Prometheus метрики
│   └── metrics.go
├── config.yaml                 # Основная конфигурация
├── secret_config.yaml          # Секретные данные (не в git!)
├── Dockerfile
├── Makefile
└── README.md
```

## 🚀 Быстрый старт

### 1. Клонирование шаблона

```bash
# Скопируйте содержимое template-external-api-service в ваш новый проект
cp -r template-external-api-service my-new-service
cd my-new-service
```

### 2. Настройка go.mod

Измените название модуля в `go.mod`:

```go
module your-service-name

go 1.24.6
```

Обновите импорты во всех файлах с `template-external-api-service` на `your-service-name`.

### 3. Настройка конфигурации

#### config.yaml
```yaml
ENV: local
address: 0.0.0.0:8080
server_timeout: 30s

# БД
dbUrl: mongodb://localhost:27017
db_name: your_db_name

# External API
bot_auth_name: your_bot_username
bot_auth_password: your_bot_password
bot_login_url: https://api.example.com/auth/login
bot_refresh_token_url: https://api.example.com/auth/refresh
external_api_base_url: https://api.example.com
```

#### secret_config.yaml (создайте файл)
```yaml
# Секретные данные - НЕ коммитьте в git!
bot_auth_name: real_bot_username
bot_auth_password: real_bot_password
```

### 4. Установка зависимостей

```bash
go mod download
```

### 5. Запуск

```bash
# Локальный запуск
make run

# Или через go
go run cmd/template-service/main.go

# Сборка и запуск бинарника
make build
./bin/template-service
```

### 6. Проверка работы

```bash
# Health check
curl http://localhost:8080/health

# Metrics
curl http://localhost:8080/metrics

# Пример запроса к external API
curl http://localhost:8080/api/v1/demands/123
curl http://localhost:8080/api/v1/accounts/456
```

## 📝 Основные возможности

### HTTP клиент с JWT авторизацией

HTTP клиент автоматически добавляет JWT токен в заголовок Authorization и обновляет его при необходимости.

```go
// Использование в вашем коде
httpClient := client.NewHTTPClient(client.ClientConfig{
    BaseURL: "https://api.example.com",
    Timeout: 30 * time.Second,
    JwtAuth: botAuth,
    Logger:  logger,
})

// Выполнение запросов
resp, err := httpClient.Get(ctx, "/api/resource")
resp, err := httpClient.Post(ctx, "/api/resource", body)
resp, err := httpClient.Put(ctx, "/api/resource/123", body)
resp, err := httpClient.Patch(ctx, "/api/resource/123", body)
resp, err := httpClient.Delete(ctx, "/api/resource/123")
```

### External API Service

Готовый сервис для работы с внешним API:

```go
// GetDemandInfo - получение информации о заявке
demandInfo, err := externalAPIService.GetDemandInfo(ctx, "123")

// GetAccountInfo - получение информации об аккаунте
accountInfo, err := externalAPIService.GetAccountInfo(ctx, "456")

// UpdateDemandStatus - обновление статуса заявки
response, err := externalAPIService.UpdateDemandStatus(ctx, "123", "new_status_id")
```

### Prometheus метрики

Автоматический сбор метрик:
- `http_requests_total` - счётчик HTTP запросов (method, path, status)
- `http_request_duration_seconds` - длительность запросов (method, path)

Доступны на `/metrics`

### Конфигурация

Поддерживается три способа конфигурации (с приоритетом):
1. `config.yaml` - базовая конфигурация
2. `secret_config.yaml` - секретные данные (перезаписывает config.yaml)
3. Переменные окружения (перезаписывают всё)

```bash
# Пример с переменными окружения
export ENV=prod
export ADDRESS=0.0.0.0:9000
export BOT_AUTH_NAME=my_bot
export BOT_AUTH_PASSWORD=secret
go run cmd/template-service/main.go
```

## 🔧 Как расширить шаблон

### 1. Добавление нового endpoint в External API Service

**Шаг 1**: Добавьте метод в интерфейс (`internal/client/external_api_service/external_api_service.go`):

```go
type ExternalAPIServiceInterface interface {
    // ...existing methods...
    GetUserProfile(ctx context.Context, userID string) (*models.UserProfileResponse, error)
}
```

**Шаг 2**: Создайте модель ответа (`internal/client/external_api_service/models/user_profile_dto.go`):

```go
package models

type UserProfileResponse struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    // ...остальные поля
}
```

**Шаг 3**: Реализуйте метод:

```go
func (s *ExternalAPIService) GetUserProfile(ctx context.Context, userID string) (*models.UserProfileResponse, error) {
    s.logger.Info("Getting user profile", slog.String("user_id", userID))

    resp, err := s.client.Get(ctx, fmt.Sprintf("/api/users/%s/profile", userID))
    if err != nil {
        return nil, fmt.Errorf("failed to get user profile: %w", err)
    }
    defer resp.Body.Close()

    result := &models.UserProfileResponse{}
    if err = client.ParseResponse(resp, result); err != nil {
        return nil, fmt.Errorf("failed to parse response: %w", err)
    }

    return result, nil
}
```

### 2. Добавление нового HTTP endpoint

В `internal/app/app.go` добавьте новый route:

```go
v1 := router.Group("/api/v1")
{
    // Новый endpoint
    v1.GET("/users/:id/profile", func(c *gin.Context) {
        userID := c.Param("id")
        
        ctx, cancel := context.WithTimeout(c.Request.Context(), cfg.ServerTimeout)
        defer cancel()

        profile, err := externalAPIService.GetUserProfile(ctx, userID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, profile)
    })
}
```

### 3. Добавление собственного сервиса

**Шаг 1**: Создайте сервис (`internal/services/my_service.go`):

```go
package services

import (
    "context"
    "log/slog"
    "your-service-name/internal/client/external_api_service"
)

type MyService struct {
    externalAPI external_api_service.ExternalAPIServiceInterface
    logger      *slog.Logger
}

func NewMyService(externalAPI external_api_service.ExternalAPIServiceInterface, logger *slog.Logger) *MyService {
    return &MyService{
        externalAPI: externalAPI,
        logger:      logger,
    }
}

func (s *MyService) DoSomething(ctx context.Context, id string) error {
    // Ваша бизнес-логика
    info, err := s.externalAPI.GetDemandInfo(ctx, id)
    if err != nil {
        return err
    }
    
    s.logger.Info("Processing", slog.String("subject", info.Subject))
    // ...
    return nil
}
```

**Шаг 2**: Инициализируйте в `app.go`:

```go
myService := services.NewMyService(externalAPIService, logger)
```

### 4. Добавление MongoDB репозитория

**Шаг 1**: Создайте модель (`internal/storage/database/repositories/user_repository.go`):

```go
package repositories

import (
    "context"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "log/slog"
)

type UserDocument struct {
    ID        primitive.ObjectID `bson:"_id,omitempty"`
    Username  string             `bson:"username"`
    Email     string             `bson:"email"`
    CreatedAt time.Time          `bson:"created_at"`
}

type UserRepository interface {
    Create(ctx context.Context, user *UserDocument) error
    FindByID(ctx context.Context, id primitive.ObjectID) (*UserDocument, error)
}

type userRepository struct {
    collection *mongo.Collection
    logger     *slog.Logger
}

func NewUserRepository(db *mongo.Database, logger *slog.Logger) UserRepository {
    return &userRepository{
        collection: db.Collection("users"),
        logger:     logger,
    }
}

func (r *userRepository) Create(ctx context.Context, user *UserDocument) error {
    user.CreatedAt = time.Now()
    result, err := r.collection.InsertOne(ctx, user)
    if err != nil {
        return err
    }
    user.ID = result.InsertedID.(primitive.ObjectID)
    return nil
}

func (r *userRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*UserDocument, error) {
    var user UserDocument
    err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
    if err != nil {
        return nil, err
    }
    return &user, nil
}
```

**Шаг 2**: Инициализируйте в `app.go`:

```go
userRepo := repositories.NewUserRepository(Db, logger)
```

### 5. Добавление middleware

В `internal/server/middlewares/middlewares.go`:

```go
func (m *Middlewares) AuthMiddleware(c *gin.Context) {
    token := c.GetHeader("Authorization")
    
    // Ваша логика проверки токена
    if token == "" {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    
    c.Next()
}
```

Использование в `app.go`:

```go
protectedGroup := v1.Group("/protected")
protectedGroup.Use(middleware.AuthMiddleware)
{
    protectedGroup.GET("/resource", handler.GetResource)
}
```

## 🐳 Docker

### Сборка образа

```bash
make docker-build
# или
docker build -t your-service-name:latest .
```

### Запуск в Docker

```bash
# С переменными окружения
docker run -p 8080:8080 \
  -e ENV=prod \
  -e BOT_AUTH_NAME=bot_user \
  -e BOT_AUTH_PASSWORD=secret \
  -e EXTERNAL_API_BASE_URL=https://api.example.com \
  your-service-name:latest

# С файлом .env
docker run -p 8080:8080 --env-file .env your-service-name:latest
```

### Docker Compose (пример)

```yaml
version: '3.8'

services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - ENV=prod
      - DB_URL=mongodb://mongo:27017
      - DB_NAME=mydb
    depends_on:
      - mongo
  
  mongo:
    image: mongo:latest
    ports:
      - "27017:27017"
    volumes:
      - mongo_data:/data/db

volumes:
  mongo_data:
```

## 📊 Мониторинг

### Healthcheck

```bash
curl http://localhost:8080/health
```

Ответ:
```json
{
  "status": "ok",
  "env": "local"
}
```

### Prometheus метрики

```bash
curl http://localhost:8080/metrics
```

Пример метрик:
```
# HELP http_requests_total Total number of HTTP requests
# TYPE http_requests_total counter
http_requests_total{method="GET",path="/health",status="200"} 42

# HELP http_request_duration_seconds Duration of HTTP requests
# TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{method="GET",path="/health",le="0.01"} 40
```

## 🔐 Безопасность

1. **Никогда не коммитьте** `secret_config.yaml` в git
2. Добавьте `secret_config.yaml` в `.gitignore` (уже добавлен)
3. Используйте переменные окружения в production
4. Храните секреты в защищённых хранилищах (Vault, AWS Secrets Manager и т.д.)

## 🧪 Тестирование

```bash
# Запуск всех тестов
make test

# Запуск с покрытием
go test -cover ./...

# Запуск конкретного теста
go test -v ./internal/client/...
```

## 📚 Полезные команды

```bash
# Форматирование кода
make fmt

# Обновление зависимостей
make deps

# Сборка
make build

# Запуск
make run

# Очистка
make clean

# Docker сборка
make docker-build
```

## 🎓 Best Practices

1. **Структурированное логирование**: Используйте slog для логирования
   ```go
   logger.Info("Message", slog.String("key", value), slog.Int("count", 10))
   ```

2. **Context с таймаутом**: Всегда используйте context.WithTimeout
   ```go
   ctx, cancel := context.WithTimeout(c.Request.Context(), cfg.ServerTimeout)
   defer cancel()
   ```

3. **Обработка ошибок**: Всегда обрабатывайте ошибки
   ```go
   if err != nil {
       logger.Error("Operation failed", slog.String("error", err.Error()))
       return err
   }
   ```

4. **Dependency Injection**: Передавайте зависимости через конструкторы

5. **Repository Pattern**: Используйте репозитории для работы с БД

## 📖 Дополнительная документация

- [GitHub Copilot Instructions](.github/copilot-instructions.md) - детальное описание для AI ассистента
- [External API Service Models](internal/client/external_api_service/models/) - модели данных API

## 🤝 Поддержка

При возникновении вопросов или проблем, обращайтесь к документации основного проекта.

## 📄 Лицензия

MIT License

