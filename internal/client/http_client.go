package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	auth2 "github.com/ShlykovPavel/JWTAuth/auth"
)

// HTTPClient интерфейс для HTTP клиента
type HTTPClient interface {
	Get(ctx context.Context, url string) (*http.Response, error)
	Post(ctx context.Context, url string, body interface{}) (*http.Response, error)
	Put(ctx context.Context, url string, body interface{}) (*http.Response, error)
	Patch(ctx context.Context, url string, body interface{}) (*http.Response, error)
	Delete(ctx context.Context, url string) (*http.Response, error)
	DoRequest(ctx context.Context, method, url string, body interface{}) (*http.Response, error)
}

// Client структура HTTP клиента с JWT авторизацией
type Client struct {
	httpClient *http.Client
	jwtAuth    *auth2.JWTAuth
	logger     *slog.Logger
	baseURL    string
}

// ClientConfig конфигурация для HTTP клиента
type ClientConfig struct {
	BaseURL string
	Timeout time.Duration
	JwtAuth *auth2.JWTAuth
	Logger  *slog.Logger
}

// NewHTTPClient создает новый HTTP клиент с JWT авторизацией
func NewHTTPClient(config ClientConfig) *Client {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	httpClient := &http.Client{
		Timeout: config.Timeout,
	}

	return &Client{
		httpClient: httpClient,
		jwtAuth:    config.JwtAuth,
		logger:     config.Logger,
		baseURL:    config.BaseURL,
	}
}

// Get выполняет GET запрос
func (c *Client) Get(ctx context.Context, url string) (*http.Response, error) {
	return c.DoRequest(ctx, http.MethodGet, url, nil)
}

// Post выполняет POST запрос
func (c *Client) Post(ctx context.Context, url string, body interface{}) (*http.Response, error) {
	return c.DoRequest(ctx, http.MethodPost, url, body)
}

// Put выполняет PUT запрос
func (c *Client) Put(ctx context.Context, url string, body interface{}) (*http.Response, error) {
	return c.DoRequest(ctx, http.MethodPut, url, body)
}

// Patch выполняет PATCH запрос
func (c *Client) Patch(ctx context.Context, url string, body interface{}) (*http.Response, error) {
	return c.DoRequest(ctx, http.MethodPatch, url, body)
}

// Delete выполняет DELETE запрос
func (c *Client) Delete(ctx context.Context, url string) (*http.Response, error) {
	return c.DoRequest(ctx, http.MethodDelete, url, nil)
}

// DoRequest выполняет HTTP запрос с автоматической авторизацией
func (c *Client) DoRequest(ctx context.Context, method, url string, body interface{}) (*http.Response, error) {
	// Подготавливаем URL
	fullURL := c.buildURL(url)

	// Подготавливаем тело запроса
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			c.logger.Error("Failed to marshal request body", slog.String("error", err.Error()))
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	// Создаем HTTP запрос
	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		c.logger.Error("Failed to create HTTP request", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Добавляем JWT токен в заголовок авторизации
	if err := c.addAuthHeader(req); err != nil {
		c.logger.Error("Failed to add auth header", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to add auth header: %w", err)
	}

	// Логируем запрос
	c.logger.Debug("Making HTTP request",
		slog.String("method", method),
		slog.String("url", fullURL),
		slog.String("Auth Header", req.Header.Get("Authorization")),
		slog.String("body", fmt.Sprintf("%v", body)))

	// Выполняем запрос
	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("HTTP request failed",
			slog.String("method", method),
			slog.String("url", fullURL),
			slog.String("error", err.Error()))
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	// Логируем ответ
	c.logger.Debug("HTTP response received",
		slog.String("method", method),
		slog.String("url", fullURL),
		slog.Int("status_code", resp.StatusCode),
		slog.String("status", resp.Status))

	return resp, nil
}

// addAuthHeader добавляет JWT токен в заголовок Authorization
func (c *Client) addAuthHeader(req *http.Request) error {
	if c.jwtAuth == nil {
		c.logger.Warn("JWT auth not configured, skipping authorization header")
		return nil
	}

	token, err := c.jwtAuth.GetToken()
	if err != nil {
		return fmt.Errorf("failed to get JWT token: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	c.logger.Debug("Authorization header set", slog.String("token", token))
	return nil
}

// buildURL строит полный URL из базового URL и переданного пути
func (c *Client) buildURL(path string) string {
	if c.baseURL == "" {
		return path
	}

	// Удаляем слеш в конце baseURL если он есть
	baseURL := c.baseURL
	if len(baseURL) > 0 && baseURL[len(baseURL)-1] == '/' {
		baseURL = baseURL[:len(baseURL)-1]
	}

	// Добавляем слеш в начало path если его нет
	if len(path) > 0 && path[0] != '/' {
		path = "/" + path
	}

	return baseURL + path
}

// Close закрывает HTTP клиент (пока что просто логирует)
func (c *Client) Close() {
	c.logger.Info("HTTP client closed")
}
