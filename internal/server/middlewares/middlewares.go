package middlewares

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"template-external-api-service/internal/config"
	jwtTokens "template-external-api-service/internal/lib/jwt_tokens"
	"time"

	"template-external-api-service/metrics"

	"github.com/gin-gonic/gin"
)

// Middlewares структура для middleware компонентов
type Middlewares struct {
	logger  *slog.Logger
	config  *config.Config
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

func (m *Middlewares) AuthMiddleware(c *gin.Context) {
	token := c.GetHeader("Authorization")
	m.logger.Debug("AuthMiddleware", slog.String("token", token))

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(token, bearerPrefix) {
		m.logger.Warn("AuthMiddleware: Missing Bearer prefix")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	tokenString := strings.TrimPrefix(token, bearerPrefix)

	_, err := jwtTokens.VerifyToken(tokenString, m.config.BotConfig.JWTSecretKey)
	if err != nil {
		m.logger.Error("AuthMiddleware", slog.String("error", err.Error()))
		c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	c.Next()

}

// APIKeyMiddleware проверяет API ключ для webhook endpoints
func (m *Middlewares) APIKeyMiddleware(c *gin.Context) {
	apiKey := c.GetHeader("X-API-KEY")

	if apiKey == "" {
		m.logger.Warn("APIKeyMiddleware: Missing API key")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing API key"})
		return
	}

	if apiKey != m.config.BotConfig.WebHookAPIKey {
		m.logger.Warn("APIKeyMiddleware: Invalid API key", slog.String("provided_key", apiKey))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
		return
	}

	m.logger.Debug("APIKeyMiddleware: API key validated successfully")
	c.Next()
}
