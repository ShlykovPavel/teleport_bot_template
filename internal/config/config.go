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
	Env            string         `yaml:"ENV" env:"ENV" env-default:"local"`
	Address        string         `yaml:"address" env:"ADDRESS" env-default:"0.0.0.0:8080"`
	ServerTimeout  time.Duration  `yaml:"server_timeout" env:"SERVER_TIMEOUT" env-default:"30s"`
	DbConfig       DbConfig       `yaml:"db_config" env:"DB_CONFIG"`
	TeleportConfig TeleportConfig `yaml:"teleport_config" env:"TELEPORT_CONFIG"`
	BotConfig      BotConfig      `yaml:"bot_config" env:"BOT_CONFIG"`

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
