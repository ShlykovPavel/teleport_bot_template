# 🚀 Quick Start Guide

Быстрый старт для использования шаблона Template External API Service.

## За 5 минут до запуска

### 1. Скопируйте шаблон

```bash
cp -r template-external-api-service my-new-service
cd my-new-service
```

### 2. Измените название модуля

В `go.mod`:
```go
module my-new-service  // было: template-external-api-service

go 1.24.6
// ...
```

### 3. Обновите импорты

Замените все импорты `template-external-api-service` на `my-new-service`:

```bash
# macOS/Linux
find . -type f -name "*.go" -exec sed -i '' 's/template-external-api-service/my-new-service/g' {} +

# Linux (без macOS)
find . -type f -name "*.go" -exec sed -i 's/template-external-api-service/my-new-service/g' {} +
```

### 4. Настройте конфигурацию

Создайте `secret_config.yaml`:

```yaml
# Секретные данные - НЕ коммитьте в git!
bot_auth_name: your_bot_username
bot_auth_password: your_bot_password
```

Отредактируйте `config.yaml`:

```yaml
ENV: local
address: 0.0.0.0:8080

# Укажите реальные URL вашего External API
bot_login_url: https://your-api.com/auth/login
bot_refresh_token_url: https://your-api.com/auth/refresh
external_api_base_url: https://your-api.com
```

### 5. Установите зависимости и запустите

```bash
go mod tidy
make run
```

### 6. Проверьте работу

```bash
# Health check
curl http://localhost:8080/health

# Metrics
curl http://localhost:8080/metrics

# Пример запроса (замените ID на реальный)
curl http://localhost:8080/api/v1/demands/123
```

## Что дальше?

### Добавьте свою бизнес-логику

1. **Новые методы External API** → `internal/client/external_api_service/`
2. **Свои сервисы** → `internal/services/`
3. **Свои репозитории** → `internal/storage/database/repositories/`
4. **Свои endpoints** → `internal/app/app.go`

### Полезные ссылки

- [README.md](README.md) - Полная документация
- [.github/copilot-instructions.md](.github/copilot-instructions.md) - Документация для AI

## Docker

```bash
# Сборка
docker build -t my-new-service:latest .

# Запуск
docker run -p 8080:8080 \
  -e BOT_AUTH_NAME=bot_user \
  -e BOT_AUTH_PASSWORD=secret \
  -e EXTERNAL_API_BASE_URL=https://api.example.com \
  my-new-service:latest
```

## Готово! 🎉

Ваш сервис запущен и готов к расширению!

