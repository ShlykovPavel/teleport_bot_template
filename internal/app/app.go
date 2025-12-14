package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"template-external-api-service/internal/storage/database/db_connections"
	"time"

	"template-external-api-service/internal/client"
	"template-external-api-service/internal/client/teleport_open_api_service"
	"template-external-api-service/internal/config"
	"template-external-api-service/internal/server/middlewares"
	"template-external-api-service/metrics"

	auth2 "github.com/ShlykovPavel/JWTAuth/auth"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// App структура приложения. Включает в себя все необходимые элементы для запуска приложения.
type App struct {
	HTTPServer             *http.Server
	logger                 *slog.Logger
	cfg                    *config.Config
	dbClient               *pgxpool.Pool
	TeleportOpenAPIService *teleport_open_api_service.TeleportOpenAPIService
}

// NewApp создаёт экземпляр приложения, инициализируя все зависимости
func NewApp(logger *slog.Logger, cfg *config.Config) *App {
	// Инициализация метрик Prometheus
	metricsInstance := metrics.InitMetrics()
	logger.Info("Prometheus metrics initialized")

	poll, err := db_connections.CreatePool(context.Background(), &cfg.DbConfig, logger)
	if err != nil {
		logger.Error("Failed to create database pool", "error", err)
		os.Exit(1)
	}

	db_connections.PostgreMonitorPool(context.Background(), poll, metricsInstance)

	// TODO: Инициализация репозиториев
	// Пример:
	// exampleRepo := repositories.NewExampleRepository(Db, logger)

	// Авторизация за бота в внешнем API
	botAuth := auth2.NewJwtAuth(
		cfg.TeleportConfig.BotLoginUrl,
		cfg.TeleportConfig.BotRefreshTokenUrl,
		cfg.TeleportConfig.BotAuthName,
		cfg.TeleportConfig.BotAuthPassword,
		10,
		logger,
	)

	if err = botAuth.Start(); err != nil {
		logger.Error("Failed to start JWT auth:", slog.Any("error", err))
	}

	token, err := botAuth.GetToken()
	if err != nil {
		logger.Error("Failed to get token:", slog.Any("error", err))
	}

	logger.Info("Successfully authenticated. Token:", slog.String("token", token[:10]+"..."))

	// Инициализация HTTP клиента с JWT авторизацией
	httpClient := client.NewHTTPClient(client.ClientConfig{
		BaseURL: cfg.TeleportConfig.TeleportAPIBaseURL,
		Timeout: cfg.ServerTimeout,
		JwtAuth: botAuth,
		Logger:  logger,
	})

	// Сервис для взаимодействия с внешним API
	teleportOpenAPIService := teleport_open_api_service.NewTeleportOpenAPIService(httpClient, logger)
	logger.Info("Teleport OpenAPI Service initialized successfully")

	// TODO: Инициализация ваших сервисов
	// Пример:
	// myService := services.NewMyService(exampleRepo, externalAPIService, logger)

	// TODO: Инициализация хендлеров
	// Пример:
	// myHandler := handlers.NewMyHandler(myService, logger)

	// Инициализация Middleware
	middleware := middlewares.NewMiddlewares(logger, cfg, metricsInstance)

	// Настройка Gin роутера
	if cfg.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.MetricsMiddleware)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"env":    cfg.Env,
		})
	})

	// Prometheus metrics endpoint
	router.GET("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})

	// TODO: Добавьте ваши API endpoints здесь
	// Пример группы API v1
	v1 := router.Group("/api/v1")
	{
		// Пример endpoint для получения информации о заявке
		v1.GET("/demands/:id", func(c *gin.Context) {
			demandID := c.Param("id")

			ctx, cancel := context.WithTimeout(c.Request.Context(), cfg.ServerTimeout)
			defer cancel()

			// Получаем информацию о заявке через External API
			demandInfo, err := teleportOpenAPIService.GetDemandInfo(ctx, demandID)
			if err != nil {
				logger.Error("Failed to get demand info",
					slog.String("demand_id", demandID),
					slog.String("error", err.Error()))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get demand info"})
				return
			}

			c.JSON(http.StatusOK, demandInfo)
		})

		// Пример endpoint для получения информации об аккаунте
		v1.GET("/accounts/:id", func(c *gin.Context) {
			accountID := c.Param("id")

			ctx, cancel := context.WithTimeout(c.Request.Context(), cfg.ServerTimeout)
			defer cancel()

			// Получаем информацию об аккаунте через External API
			accountInfo, err := teleportOpenAPIService.GetAccountInfo(ctx, accountID)
			if err != nil {
				logger.Error("Failed to get account info",
					slog.String("account_id", accountID),
					slog.String("error", err.Error()))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get account info"})
				return
			}

			c.JSON(http.StatusOK, accountInfo)
		})

		// TODO: Добавьте здесь ваши endpoints
		// Пример:
		// myGroup := v1.Group("/my-resource")
		// {
		//     myGroup.GET("/", myHandler.List)
		//     myGroup.GET("/:id", myHandler.Get)
		//     myGroup.POST("/", myHandler.Create)
		//     myGroup.PUT("/:id", myHandler.Update)
		//     myGroup.DELETE("/:id", myHandler.Delete)
		// }
	}

	srv := &http.Server{
		Addr:              cfg.Address,
		Handler:           router,
		ReadHeaderTimeout: cfg.ServerTimeout,
		WriteTimeout:      cfg.ServerTimeout,
	}

	return &App{
		cfg:                    cfg,
		logger:                 logger,
		HTTPServer:             srv,
		dbClient:               poll,
		TeleportOpenAPIService: teleportOpenAPIService,
	}
}

// Run запускает HTTP-сервер
func (a *App) Run() {
	a.logger.Info("Starting application")

	a.logger.Info("Starting HTTP server", slog.String("address", a.cfg.Address))

	go func() {
		if err := a.HTTPServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Error("Failed to start server", "error", err.Error())
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	a.logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.HTTPServer.Shutdown(ctx); err != nil {
		a.logger.Error("Server forced to shutdown", "error", err.Error())
		os.Exit(1)
	}

	a.dbClient.Close()
	a.logger.Info("Disconnected from database")

	a.logger.Info("Server stopped")
}
