package config

type TeleportConfig struct {
	// Авторизация бота в телепорте
	// Необходимо завести пользователя бота в телепорте и получить access_key и secret_key
	BotAuthName        string `yaml:"bot_auth_name" env:"BOT_AUTH_NAME" env-required:"true"`                 // access_key бота для авторизации в систему
	BotAuthPassword    string `yaml:"bot_auth_password" env:"BOT_AUTH_PASSWORD" env-required:"true"`         // secret_key бота для авторизации в систему
	BotLoginUrl        string `yaml:"bot_login_url" env:"BOT_LOGIN_URL" env-required:"true"`                 // URL для логина бота в телепорт
	BotRefreshTokenUrl string `yaml:"bot_refresh_token_url" env:"BOT_REFRESH_TOKEN_URL" env-required:"true"` // URL для обновления токена бота в телепорт

	TeleportAPIBaseURL string `yaml:"teleport_api_base_url" env:"TELEPORT_API_BASE_URL" env-required:"true"` // Базовый URL внешнего API Teleport по которому мы будем ходить в open api

}
