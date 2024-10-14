package service

import (
	"log/slog"
	"sync"
)

type Service struct {
	pricePrecisions map[string][]int32
	qtyPrecisions   map[string]int32
	mu              sync.RWMutex
	logger          *slog.Logger
}

func New(logger *slog.Logger) *Service {
	return &Service{
		pricePrecisions: map[string][]int32{"BTC_USDT": {-1}},
		qtyPrecisions:   map[string]int32{"BTC_USDT": 6},
		mu:              sync.RWMutex{},
		logger:          logger,
	}
}
