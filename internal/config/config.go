package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// TODO: Актуализируйте конфигурацию под свои нужды
// Добавьте или удалите поля в зависимости от требований вашего проекта

// Config представляет конфигурацию приложения
type Config struct {
	// Основные параметры приложения
	Env           string        `yaml:"ENV" env:"ENV" env-default:"local"`
	Address       string        `yaml:"address" env:"ADDRESS" env-default:"0.0.0.0:8080"`
	ServerTimeout time.Duration `yaml:"server_timeout" env:"SERVER_TIMEOUT" env-default:"30s"`

	// TODO: Настройки базы данных - актуализируйте под вашу БД (MongoDB, PostgreSQL, MySQL и т.д.)
	DbUrl            string `yaml:"dbUrl" env:"DB_URL" env-required:"true"`
	DbName           string `yaml:"db_name" env:"DB_NAME" env-required:"true"`
	DbUser           string `yaml:"db_user" env:"DB_USER"`
	DbPassword       string `yaml:"db_password" env:"DB_PASSWORD"`
	DbMaxConnections uint64 `yaml:"db_max_connections" env:"DB_MAX_CONNECTIONS" env-default:"100"`

	// JWT авторизация для внешнего API (бот)
	BotAuthName        string `yaml:"bot_auth_name" env:"BOT_AUTH_NAME" env-required:"true"`
	BotAuthPassword    string `yaml:"bot_auth_password" env:"BOT_AUTH_PASSWORD" env-required:"true"`
	BotLoginUrl        string `yaml:"bot_login_url" env:"BOT_LOGIN_URL" env-required:"true"`
	BotRefreshTokenUrl string `yaml:"bot_refresh_token_url" env:"BOT_REFRESH_TOKEN_URL" env-required:"true"`

	// External API
	ExternalAPIBaseURL string `yaml:"external_api_base_url" env:"EXTERNAL_API_BASE_URL" env-required:"true"`

	// TODO: Добавьте здесь дополнительные параметры конфигурации для вашего сервиса
	// Примеры:
	// - API ключи для различных сервисов
	// - Параметры retry-логики
	// - Настройки таймаутов
	// - Пути к шаблонам/файлам
	// - Бизнес-параметры вашего приложения
}

// LoadConfig загружает конфигурацию из файла и переменных окружения
// Порядок загрузки:
// 1. config.yaml (базовый файл)
// 2. secret_config.yaml (секретные данные)
// 3. Переменные окружения (перезаписывают предыдущие значения)
func LoadConfig(secretFilePath string) (*Config, error) {
	var cfg Config

	// Читаем основной конфиг файл
	err := cleanenv.ReadConfig("config.yaml", &cfg)
	if err != nil {
		log.Default().Printf("Error loading config from base config file: %s", err.Error())
	}

	// Читаем конфиг файл с секретами
	err = cleanenv.ReadConfig(secretFilePath, &cfg)
	if err != nil {
		log.Default().Printf("Error loading config from secret config file: %s", err.Error())
	}

	// В конце смотрим в переменные окружения
	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Default().Printf("Error loading config from env vars: %s", err.Error())
		return nil, err
	}

	return &cfg, nil
}
