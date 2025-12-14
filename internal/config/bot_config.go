package config

import "time"

type BotConfig struct {
	// У бота есть своя админка с JWT авторизацией
	// JWT параметры для вашего сервиса
	JWTSecretKey string        `yaml:"jwt_secret_key" env:"JWT_SECRET_KEY" env-required:"true"`
	JWTDuration  time.Duration `yaml:"jwt_duration"  env:"JWT_DURATION" env-default:"5m"`

	UserEmail    string `yaml:"user_email" env:"USER_EMAIL" env-required:"true"`       // логин пользователя для авторизации в систему
	UserPassword string `yaml:"user_password" env:"USER_PASSWORD" env-required:"true"` // пароль пользователя для авторизации в систему
	// Бот принимает запросы по веб хуку от телепорта и читает название секции от которой нужно обрабатывать заявки
	// TODO Убрать после реализации настройки секций в админке бота
	SectionName string `yaml:"section_name" env:"SECTION_NAME" env-required:"true"` // название секции, заявки от которой ждём с веб хука

	// У заявок есть статусы, которые нужно указывать при обработке заявок
	StatusWorkInProgressId string `yaml:"status_work_in_progress_id" env:"STATUS_WORK_IN_PROGRESS_ID"`
	StatusProcessedId      string `yaml:"status_processed_id" env:"STATUS_PROCESSED_ID"`

	// При настройке веб хуков в телепорте нужно указать API ключ с которым веб хук будет присылать запросы на наш бот
	// В нашем боте необходимо тоже его указать что б проверять входящие запросы
	WebHookAPIKey string `yaml:"web_hook_api_key" env:"WEB_HOOK_API_KEY" env-required:"true"` // API ключ для приёма веб хуков
}
