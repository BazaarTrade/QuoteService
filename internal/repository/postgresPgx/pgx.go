package postgresPgx

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func New(DB_CONNECTION string, logger *slog.Logger) (*Postgres, error) {
	conn, err := pgxpool.New(context.Background(), DB_CONNECTION)
	if err != nil {
		logger.Error("failed to create pgxPool connection", "error", err)
		return nil, err
	}

	return &Postgres{
		db:     conn,
		logger: logger,
	}, nil
}

func (p *Postgres) Close() {
	p.db.Close()
}
