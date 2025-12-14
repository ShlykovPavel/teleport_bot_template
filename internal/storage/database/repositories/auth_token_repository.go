package repositories

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthTokenRepository struct {
	logger         *slog.Logger
	connectionPool *pgxpool.Pool
}
