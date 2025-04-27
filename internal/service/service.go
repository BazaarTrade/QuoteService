package service

import (
	"log/slog"
	"sync"

	"github.com/BazaarTrade/QuoteService/internal/repository"
)

type Service struct {
	Streams map[string]*StreamHub

	mu     sync.RWMutex
	db     repository.Repository
	logger *slog.Logger
}

func New(db repository.Repository, logger *slog.Logger) *Service {
	return &Service{
		Streams: make(map[string]*StreamHub),
		db:      db,
		logger:  logger,
	}
}
