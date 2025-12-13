package middlewares

import (
	"log/slog"
	"strconv"
	"time"

	"template-external-api-service/metrics"

	"github.com/gin-gonic/gin"
)

// Middlewares структура для middleware компонентов
type Middlewares struct {
	logger  *slog.Logger
	metrics *metrics.Metrics
}

// NewMiddlewares создает новый набор middleware
func NewMiddlewares(logger *slog.Logger, metrics *metrics.Metrics) *Middlewares {
	return &Middlewares{
		logger:  logger,
		metrics: metrics,
	}
}

// MetricsMiddleware отслеживает метрики HTTP запросов
func (m *Middlewares) MetricsMiddleware(c *gin.Context) {
	start := time.Now()
	path := c.FullPath()
	if path == "" {
		path = c.Request.URL.Path
	}
	method := c.Request.Method

	c.Next()

	duration := time.Since(start).Seconds()
	status := strconv.Itoa(c.Writer.Status())

	m.metrics.HttpRequestsTotal.WithLabelValues(method, path, status).Inc()
	m.metrics.HttpRequestDuration.WithLabelValues(method, path).Observe(duration)
}

// CORSMiddleware добавляет CORS заголовки
func (m *Middlewares) CORSMiddleware(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-API-KEY")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}

// TODO: Добавьте здесь свои middleware по необходимости
// Примеры:
// - AuthMiddleware для JWT авторизации
// - APIKeyMiddleware для проверки API ключей
// - RateLimitMiddleware для ограничения запросов
// - LoggingMiddleware для детального логирования
