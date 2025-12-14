package teleport_open_api_service

import (
	"context"
	"fmt"
	"log/slog"
	"template-external-api-service/internal/client"
	"template-external-api-service/internal/client/teleport_open_api_service/models"
)

// TeleportOpenAPIServiceInterface интерфейс для работы с внешним API
type TeleportOpenAPIServiceInterface interface {
	UpdateDemandStatus(ctx context.Context, demandID string, status string) (*DemandStatusResponse, error)
	GetDemandInfo(ctx context.Context, demandID string) (*models.DemandInfoResponse, error)
	GetAccountInfo(ctx context.Context, userID string) (*models.AccountInfoResponse, error)
}

// TeleportOpenAPIService сервис для работы с внешним API
type TeleportOpenAPIService struct {
	client client.HTTPClient
	logger *slog.Logger
}

// NewExternalAPIService создает новый сервис для работы с внешним API
func NewTeleportOpenAPIService(client client.HTTPClient, logger *slog.Logger) *TeleportOpenAPIService {
	return &TeleportOpenAPIService{
		client: client,
		logger: logger,
	}
}

// DemandStatusRequest структура для обновления статуса заявки
type DemandStatusRequest struct {
	Status string `json:"newStatusId"`
}

// DemandStatusResponse структура ответа на обновление статуса заявки
type DemandStatusResponse struct {
	ID         string `json:"id"`
	Status     string `json:"status"`
	Message    string `json:"message,omitempty"`
	StatusCode int
}

// UpdateDemandStatus обновляет статус заявки через PATCH запрос
func (s *TeleportOpenAPIService) UpdateDemandStatus(ctx context.Context, demandID string, status string) (*DemandStatusResponse, error) {
	s.logger.Info("Updating demand status",
		slog.String("demand_id", demandID),
		slog.String("status", status))

	req := DemandStatusRequest{
		Status: status,
	}

	// Отправляем PATCH запрос
	resp, err := s.client.Patch(ctx, fmt.Sprintf("/api/open/demands/%s/status", demandID), req)
	if err != nil {
		s.logger.Error("Failed to update demand status",
			slog.String("demand_id", demandID),
			slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to update demand status: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			s.logger.Warn("Failed to close response body", slog.String("error", closeErr.Error()))
		}
	}()

	result := &DemandStatusResponse{
		StatusCode: resp.StatusCode,
	}

	// Если статус код 200 - успешное обновление без тела ответа
	if resp.StatusCode == 204 {
		s.logger.Info("Successfully updated demand status",
			slog.String("demand_id", demandID),
			slog.String("status", status),
			slog.Int("status_code", resp.StatusCode))
		return result, nil
	}

	// Если статус код не 200 - парсим тело ответа для получения деталей
	if err = client.ParseResponse(resp, result); err != nil {
		s.logger.Error("Failed to parse update demand status response",
			slog.String("demand_id", demandID),
			slog.Int("status_code", resp.StatusCode),
			slog.String("error", err.Error()))
		return result, fmt.Errorf("failed to parse response: %w", err)
	}

	// Логируем результат в зависимости от статус кода
	if resp.StatusCode >= 400 {
		s.logger.Error("Failed to update demand status",
			slog.String("demand_id", demandID),
			slog.Int("status_code", resp.StatusCode),
			slog.String("message", result.Message))
	} else {
		s.logger.Info("Updated demand status with non-200 response",
			slog.String("demand_id", demandID),
			slog.Int("status_code", resp.StatusCode),
			slog.String("message", result.Message))
	}

	return result, fmt.Errorf("Error updating demand status: HTTP %d", resp.StatusCode)
}

// GetDemandInfo получает информацию о заявке
func (s *TeleportOpenAPIService) GetDemandInfo(ctx context.Context, demandID string) (*models.DemandInfoResponse, error) {
	s.logger.Info("Getting demand info", slog.String("demand_id", demandID))

	resp, err := s.client.Get(ctx, fmt.Sprintf("/api/open/demands/%s", demandID))
	if err != nil {
		s.logger.Error("Failed to get demand info",
			slog.String("demand_id", demandID),
			slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get demand info: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			s.logger.Warn("Failed to close response body", slog.String("error", closeErr.Error()))
		}
	}()

	result := &models.DemandInfoResponse{
		StatusCode: resp.StatusCode,
	}

	// Если статус код 200 - парсим тело ответа
	if resp.StatusCode == 200 {
		if err = client.ParseResponse(resp, result); err != nil {
			s.logger.Error("Failed to parse demand info response",
				slog.String("demand_id", demandID),
				slog.String("error", err.Error()))
			return result, fmt.Errorf("failed to parse response: %w", err)
		}
		s.logger.Info("Successfully retrieved demand info",
			slog.String("demand_id", demandID),
			slog.String("demand_subject", result.Subject),
			slog.String("status_name", result.StatusName),
			slog.Int("status_code", resp.StatusCode))
		return result, nil
	}

	// Если статус код не 200 - обрабатываем как ошибку
	s.logger.Error("Failed to get demand info",
		slog.String("demand_id", demandID),
		slog.Int("status_code", resp.StatusCode))

	return result, fmt.Errorf("failed to get demand info: HTTP %d", resp.StatusCode)
}

// GetAccountInfo получает информацию об аккаунте
func (s *TeleportOpenAPIService) GetAccountInfo(ctx context.Context, userId string) (*models.AccountInfoResponse, error) {
	s.logger.Info("Getting account info", slog.String("user_id", userId))

	resp, err := s.client.Get(ctx, fmt.Sprintf("/api/open/accounts/%s", userId))
	if err != nil {
		s.logger.Error("Failed to get account info",
			slog.String("user_id", userId),
			slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get account info: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			s.logger.Warn("Failed to close response body", slog.String("error", closeErr.Error()))
		}
	}()

	result := &models.AccountInfoResponse{
		StatusCode: resp.StatusCode,
	}

	// Если статус код 200 - парсим тело ответа
	if resp.StatusCode == 200 {
		if err = client.ParseResponse(resp, result); err != nil {
			s.logger.Error("Failed to parse account info response",
				slog.String("user_id", userId),
				slog.String("error", err.Error()))
			return result, fmt.Errorf("failed to parse response: %w", err)
		}
		s.logger.Info("Successfully retrieved account info",
			slog.String("user_id", userId),
			slog.String("account_first_name", result.FirstName),
			slog.String("account_last_name", result.LastName),
			slog.String("account_email", result.Email),
			slog.Int("status_code", resp.StatusCode))
		return result, nil
	}

	// Если статус код не 200 - обрабатываем как ошибку
	s.logger.Error("Failed to get account info",
		slog.String("user_id", userId),
		slog.Int("status_code", resp.StatusCode))

	return result, fmt.Errorf("failed to get account info: HTTP %d", resp.StatusCode)
}
