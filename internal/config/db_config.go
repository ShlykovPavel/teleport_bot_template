package config

type DbConfig struct {
	// TODO: Настройки базы данных - актуализируйте под вашу БД (MongoDB, PostgreSQL, MySQL и т.д.)
	DbUrl            string `yaml:"dbUrl" env:"DB_URL" env-required:"true"`
	DbName           string `yaml:"db_name" env:"DB_NAME" env-required:"true"`
	DbUser           string `yaml:"db_user" env:"DB_USER"`
	DbPassword       string `yaml:"db_password" env:"DB_PASSWORD"`
	DbMaxConnections uint64 `yaml:"db_max_connections" env:"DB_MAX_CONNECTIONS" env-default:"100"`
}
