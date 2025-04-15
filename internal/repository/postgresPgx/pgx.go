package postgresPgx

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewPostgres(logger *slog.Logger) (*Postgres, error) {
	conn, err := pgxpool.New(context.Background(), os.Getenv("DB_CONNECTION"))
	if err != nil {
		logger.Error("failed to create pgxPool connection", "error", err)
		return nil, err
	}

	return &Postgres{
		db:     conn,
		logger: logger,
	}, nil
}
