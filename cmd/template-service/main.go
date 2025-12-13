package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"template-external-api-service/internal/app"
	"template-external-api-service/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// TODO: Добавьте Swagger документацию если необходимо
// @title Template External API Service
// @version 1.0
// @description Шаблонный сервис для работы с внешним API
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig("secret_config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Configuration loaded: ENV=%s, Address=%s\n", cfg.Env, cfg.Address)

	// Инициализация логгера
	logger := setupLogger(cfg.Env)
	logger.Info("Starting application",
		slog.String("env", cfg.Env),
		slog.String("address", cfg.Address))
	logger.Debug("Debug messages enabled")

	// Создание и запуск приложения
	application := app.NewApp(logger, cfg)
	application.Run()
}

// setupLogger настраивает логгер в зависимости от окружения
func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		// Локальное окружение: текстовый формат, уровень DEBUG
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envDev:
		// Dev окружение: JSON формат, уровень DEBUG
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envProd:
		// Prod окружение: JSON формат, уровень INFO
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	default:
		// По умолчанию: текстовый формат, уровень INFO
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return logger
}
