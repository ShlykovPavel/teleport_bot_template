FROM golang:1.24-alpine AS builder

WORKDIR /build_app

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o template-service ./cmd/template-service

# Финальный образ
FROM alpine:latest AS final

# Устанавливаем CA сертификаты и timezone
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Копируем бинарник из builder
COPY --from=builder /build_app/template-service .

# Копируем конфигурационные файлы
COPY --from=builder /build_app/config.yaml .
# TODO: Не забудьте добавить secret_config.yaml при деплое или использовать переменные окружения

# TODO: Если у вас есть статические файлы, шаблоны или другие ресурсы - скопируйте их здесь
# Пример:
# COPY --from=builder /build_app/templates ./templates/

EXPOSE 8080

ENTRYPOINT ["./template-service"]

