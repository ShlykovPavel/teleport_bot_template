# Template External API Service - Контекст для GitHub Copilot

## 📋 Обзор проекта

**Template External API Service** - это шаблонный Go-микросервис для работы с внешним API с JWT авторизацией. Предназначен для быстрого создания новых сервисов, которым нужно интегрироваться с внешними API.

### Основные технологии
- **Язык**: Go 1.24.6
- **Фреймворк**: Gin (веб-сервер)
- **БД**: MongoDB (опционально)
- **Мониторинг**: Prometheus
- **JWT**: ShlykovPavel/JWTAuth (авто-обновление токенов)

## 🏗 Архитектура проекта

### Структура модулей (module: template-external-api-service)

```
cmd/template-service/       # Точка входа приложения
internal/
  app/                      # Инициализация и запуск приложения
  client/                   # HTTP клиент с JWT авторизацией
    external_api_service/   # Сервис для работы с внешним API
      models/               # DTO для API
  config/                   # Управление конфигурацией
  server/                   # HTTP обработчики
    middlewares/            # Middleware (метрики, CORS)
  services/                 # Бизнес-логика (TODO - добавить свою)
  storage/
    database/               # Подключение к БД
      db_errors/            # Обработка ошибок MongoDB
      repositories/         # Репозитории (TODO - добавить свои)
metrics/                    # Prometheus метрики
```

## 🔄 Основной флоу работы

1. **Инициализация приложения** (`app.NewApp`)
   - Настройка логгера (slog)
   - Загрузка конфигурации (yaml + env)
   - Подключение к БД (опционально)
   - Инициализация JWT авторизации для бота
   - Создание HTTP клиента с JWT
   - Инициализация External API Service
   - Настройка Gin router с middleware

2. **Обработка HTTP запросов**
   - Middleware: Metrics → CORS → Recovery
   - Роутинг через Gin
   - Примеры endpoints: `/health`, `/metrics`, `/api/v1/demands/:id`, `/api/v1/accounts/:id`

3. **Работа с External API**
   - Автоматическая JWT авторизация
   - Методы: GetDemandInfo, GetAccountInfo, UpdateDemandStatus
   - Автоматическое обновление токенов

## 📊 Модели данных

### Конфигурация (Config)

```go
type Config struct {
    // Приложение
    Env           string
    Address       string
    ServerTimeout time.Duration
    
    // База данных
    DbUrl            string
    DbName           string
    DbUser           string
    DbPassword       string
    DbMaxConnections uint64
    
    // JWT авторизация для бота
    BotAuthName        string
    BotAuthPassword    string
    BotLoginUrl        string
    BotRefreshTokenUrl string
    
    // External API
    ExternalAPIBaseURL string
}
```

### External API Models

#### DemandInfoResponse
```go
type DemandInfoResponse struct {
    ID                  int                `json:"id"`
    Subject             string             `json:"subject"`
    SectionName         string             `json:"sectionName"`
    Status              string             `json:"status"`
    AccountFirstName    string             `json:"accountFirstName"`
    AccountLastName     string             `json:"accountLastName"`
    FormCellAnswers     [][]FormCellAnswer `json:"formCellAnswers"`
    // ... другие поля
}
```

#### AccountInfoResponse
```go
type AccountInfoResponse struct {
    ID             int              `json:"id"`
    FirstName      string           `json:"firstName"`
    MiddleName     string           `json:"middleName"`
    LastName       string           `json:"lastName"`
    Email          string           `json:"email"`
    Phone          string           `json:"phone"`
    AdditionalInfo []AdditionalInfo `json:"additionalInfo"`
    // ... другие поля
}
```

## 🔌 API Endpoints

### Public endpoints
- `GET /health` - healthcheck (возвращает статус "ok" и ENV)
- `GET /metrics` - Prometheus метрики

### API v1 endpoints (примеры)
- `GET /api/v1/demands/:id` - Получение информации о заявке через External API
- `GET /api/v1/accounts/:id` - Получение информации об аккаунте через External API

### TODO: Добавьте свои endpoints здесь

## ⚙️ Конфигурация

### Порядок загрузки конфигурации
1. Читается `config.yaml` (базовая конфигурация)
2. Перезаписывается из `secret_config.yaml` (секретные данные)
3. Перезаписывается из переменных окружения (ENV)

### Основные параметры

```yaml
# Приложение
ENV: local|dev|prod
address: 0.0.0.0:8080
server_timeout: 30s

# База данных (TODO: актуализируйте под свою БД)
dbUrl: mongodb://localhost:27017
db_name: template_service
db_max_connections: 100

# JWT авторизация для бота
bot_auth_name: bot_username
bot_auth_password: bot_password
bot_login_url: https://example.com/api/auth/login
bot_refresh_token_url: https://example.com/api/auth/refresh

# External API
external_api_base_url: https://api.example.com
```

## 🔐 Авторизация

### JWT для бота (ShlykovPavel/JWTAuth)
- Автоматическое получение токенов
- Автоматическое обновление при истечении
- Используется для всех запросов к External API
- Настраивается через `bot_auth_name`, `bot_auth_password`, `bot_login_url`, `bot_refresh_token_url`

### Как это работает
1. При старте приложения создается `JWTAuth` инстанс
2. Вызывается `botAuth.Start()` - получение первого токена
3. HTTP клиент автоматически добавляет токен в заголовок `Authorization: Bearer {token}`
4. При необходимости токен автоматически обновляется

## 📈 Метрики (Prometheus)

```go
type Metrics struct {
    HttpRequestsTotal   *prometheus.CounterVec  // method, path, status
    HttpRequestDuration *prometheus.HistogramVec // method, path
}
```

Доступны на `/metrics`. Автоматически собираются через `MetricsMiddleware`.

## 🎯 Важные паттерны и практики

### 1. Обработка ошибок
```go
if err != nil {
    logger.Error("Message", slog.String("error", err.Error()))
    c.JSON(500, gin.H{"error": "User-friendly message"})
    return
}
```

### 2. Структурированное логирование (slog)
```go
logger.Info("Message", 
    slog.String("key", value),
    slog.Int("count", 10))
```

### 3. Context с таймаутом
```go
ctx, cancel := context.WithTimeout(c.Request.Context(), cfg.ServerTimeout)
defer cancel()
```

### 4. HTTP клиент - все методы
```go
resp, err := httpClient.Get(ctx, "/api/resource")
resp, err := httpClient.Post(ctx, "/api/resource", body)
resp, err := httpClient.Put(ctx, "/api/resource/123", body)
resp, err := httpClient.Patch(ctx, "/api/resource/123", body)
resp, err := httpClient.Delete(ctx, "/api/resource/123")
```

### 5. Парсинг ответов API
```go
var result MyResponse
if err := client.ParseResponse(resp, &result); err != nil {
    return fmt.Errorf("failed to parse response: %w", err)
}
```

### 6. Repository Pattern (для БД)
```go
type MyRepository interface {
    Create(ctx context.Context, doc *MyDocument) error
    FindByID(ctx context.Context, id primitive.ObjectID) (*MyDocument, error)
    Update(ctx context.Context, doc *MyDocument) error
    Delete(ctx context.Context, id primitive.ObjectID) error
}
```

### 7. Dependency Injection
- Все зависимости передаются через конструкторы
- Инициализация в `app.NewApp()`

## 🔍 External API Service

### Интерфейс

```go
type ExternalAPIServiceInterface interface {
    UpdateDemandStatus(ctx context.Context, demandID string, status string) (*DemandStatusResponse, error)
    GetDemandInfo(ctx context.Context, demandID string) (*models.DemandInfoResponse, error)
    GetAccountInfo(ctx context.Context, userID string) (*models.AccountInfoResponse, error)
}
```

### Использование

```go
// В вашем сервисе или хендлере
demandInfo, err := externalAPIService.GetDemandInfo(ctx, "123")
if err != nil {
    return err
}

accountInfo, err := externalAPIService.GetAccountInfo(ctx, "456")
if err != nil {
    return err
}

statusResp, err := externalAPIService.UpdateDemandStatus(ctx, "123", "4242")
if err != nil {
    return err
}
```

## 💡 Как расширить шаблон

### 1. Добавление нового метода в External API Service

**Шаг 1**: Добавить в интерфейс
```go
type ExternalAPIServiceInterface interface {
    // ...existing methods...
    GetUserProfile(ctx context.Context, userID string) (*models.UserProfileResponse, error)
}
```

**Шаг 2**: Создать модель DTO
```go
// internal/client/external_api_service/models/user_profile_dto.go
type UserProfileResponse struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}
```

**Шаг 3**: Реализовать метод
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
        return nil, err
    }
    
    return result, nil
}
```

### 2. Добавление нового HTTP endpoint

В `internal/app/app.go`:

```go
v1.GET("/users/:id/profile", func(c *gin.Context) {
    userID := c.Param("id")
    
    ctx, cancel := context.WithTimeout(c.Request.Context(), cfg.ServerTimeout)
    defer cancel()
    
    profile, err := externalAPIService.GetUserProfile(ctx, userID)
    if err != nil {
        logger.Error("Failed to get user profile", 
            slog.String("user_id", userID),
            slog.String("error", err.Error()))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user profile"})
        return
    }
    
    c.JSON(http.StatusOK, profile)
})
```

### 3. Добавление собственного сервиса

**Шаг 1**: Создать сервис
```go
// internal/services/my_service.go
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

func (s *MyService) ProcessData(ctx context.Context, id string) error {
    info, err := s.externalAPI.GetDemandInfo(ctx, id)
    if err != nil {
        return err
    }
    
    s.logger.Info("Processing", slog.String("subject", info.Subject))
    // Ваша бизнес-логика
    return nil
}
```

**Шаг 2**: Инициализировать в `app.go`
```go
myService := services.NewMyService(externalAPIService, logger)
```

### 4. Добавление MongoDB репозитория

```go
// internal/storage/database/repositories/user_repository.go
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

// Реализация методов...
```

### 5. Добавление middleware

```go
// internal/server/middlewares/middlewares.go
func (m *Middlewares) AuthMiddleware(c *gin.Context) {
    token := c.GetHeader("Authorization")
    
    if token == "" {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    
    // Проверка токена
    // ...
    
    c.Next()
}
```

Использование:
```go
protectedGroup := v1.Group("/protected")
protectedGroup.Use(middleware.AuthMiddleware)
{
    protectedGroup.GET("/resource", handler.GetResource)
}
```

## 🐛 Дебаг и траблшутинг

### Уровни логирования по окружению
- `local`: DEBUG (текстовый формат)
- `dev`: DEBUG (JSON формат)
- `prod`: INFO (JSON формат)

### Типичные проблемы

1. **Не подключается к External API**
   - Проверить `external_api_base_url` в конфиге
   - Проверить учетные данные бота
   - Проверить доступность API

2. **JWT токен не обновляется**
   - Проверить `bot_login_url` и `bot_refresh_token_url`
   - Проверить логи: "Failed to start JWT auth"

3. **БД не подключается**
   - Проверить `dbUrl` в конфиге
   - Проверить доступность MongoDB

## 📚 Полезные команды

```bash
# Сборка
make build

# Запуск
make run

# Тесты
make test

# Docker сборка
make docker-build

# Очистка
make clean

# Форматирование
make fmt

# Обновление зависимостей
make deps
```

## 🎓 Для Copilot: Ключевые моменты

При генерации кода учитывай:
- Использовать структурированное логирование (slog)
- Всегда использовать context.Context с таймаутом
- Обрабатывать ошибки MongoDB через db_errors.HandleMongoError (если используется БД)
- Использовать Gin для HTTP handlers
- Следовать паттерну репозиториев для БД операций
- Инъекция зависимостей через конструкторы
- Конфигурация через yaml + env переменные
- Prometheus метрики для HTTP запросов
- JWT авторизация через ShlykovPavel/JWTAuth (автоматическое обновление)
- Все методы HTTP клиента: Get, Post, Put, Patch, Delete
- Парсинг ответов через client.ParseResponse

## 🔄 Типичные задачи

### Как добавить новое поле в конфигурацию
1. Добавить в `config.Config`: `NewField string yaml:"new_field" env:"NEW_FIELD"`
2. Добавить в `config.yaml`: `new_field: value`
3. Использовать в коде: `cfg.NewField`

### Как добавить новый endpoint External API
1. Создать модель DTO в `models/`
2. Добавить метод в интерфейс `ExternalAPIServiceInterface`
3. Реализовать метод в `external_api_service.go`
4. Использовать в хендлере или сервисе

### Как добавить новый HTTP endpoint
1. Создать хендлер или добавить в `app.go`
2. Зарегистрировать route в `router`
3. Использовать middleware при необходимости

### Как работать с БД
1. Создать модель документа
2. Создать интерфейс репозитория
3. Реализовать методы с обработкой ошибок
4. Инициализировать в `app.NewApp()`

## TODO: Актуализация шаблона

При использовании этого шаблона:
1. Измените `module` в `go.mod` на название вашего проекта
2. Обновите все импорты с `template-external-api-service` на название вашего модуля
3. Актуализируйте `config.yaml` и `secret_config.yaml` под ваши нужды
4. Добавьте свои репозитории в `internal/storage/database/repositories/`
5. Добавьте свои сервисы в `internal/services/`
6. Добавьте свои хендлеры в `internal/server/`
7. Расширьте `ExternalAPIService` новыми методами под ваш API
8. Обновите этот файл под свой проект

