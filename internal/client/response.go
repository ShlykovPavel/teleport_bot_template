package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// APIError структура для ошибок API
type APIError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Body       string `json:"body,omitempty"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error: %d %s", e.StatusCode, e.Message)
}

// ParseResponse парсит HTTP ответ в указанную структуру
func ParseResponse(resp *http.Response, target interface{}) error {

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус код
	if resp.StatusCode >= 400 {
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status),
			Body:       string(body),
		}
	}

	// Парсим JSON ответ
	if err = json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}
