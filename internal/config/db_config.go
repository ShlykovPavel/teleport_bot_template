package config

import "time"

type DbConfig struct {
	// TODO: Настройки базы данных - актуализируйте под вашу БД (MongoDB, PostgreSQL, MySQL и т.д.)
	DbHost              string        `yaml:"dbHost" env:"DB_HOST" env-required:"true"`
	DbPort              string        `yaml:"db_port" env:"DB_PORT" env-required:"true"`
	DbName              string        `yaml:"db_name" env:"DB_NAME" env-required:"true"`
	DbUser              string        `yaml:"db_user" env:"DB_USER"`
	DbPassword          string        `yaml:"db_password" env:"DB_PASSWORD"`
	DbMaxConnections    int32         `yaml:"db_max_connections" env:"DB_MAX_CONNECTIONS" env-default:"10"`
	DbMinConnections    int32         `yaml:"db_min_connections" env:"DB_MIN_CONNECTIONS"`
	DbMaxConnLifetime   time.Duration `yaml:"db_max_conn_lifetime" env:"DB_MAX_CONN_LIFETIME"`
	DbMaxConnIdleTime   time.Duration `yaml:"db_max_conn_idle_time" env:"DB_MAX_CONN_IDLE_TIME"`
	DbHealthCheckPeriod time.Duration `yaml:"db_health_check_period" env:"DB_HEALTH_CHECK_PERIOD"`
}
