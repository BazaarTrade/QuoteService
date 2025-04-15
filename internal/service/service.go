package service

import (
	"log/slog"
	"sync"

	"github.com/BazaarTrade/QuoteService/internal/models.go"
	"github.com/BazaarTrade/QuoteService/internal/repository"
)

type Service struct {
	PrecisedTrades map[string]chan []models.Trade
	PrecisedOBSs   map[string]chan map[int32]models.OrderBookSnapshot
	Ticker         map[string]chan models.Ticker
	Candlestick    map[string]chan models.Candlestick

	candlestickTimeframes map[string]map[string]*models.Candlestick
	orderBookPrecisions   map[string]orderBookPrecisions

	Tickers map[string]*ticker

	db     repository.Repository
	mu     sync.RWMutex
	logger *slog.Logger
}

func New(db repository.Repository, logger *slog.Logger) *Service {
	return &Service{
		PrecisedTrades:        make(map[string]chan []models.Trade),
		Ticker:                make(map[string]chan models.Ticker),
		PrecisedOBSs:          make(map[string]chan map[int32]models.OrderBookSnapshot),
		Candlestick:           make(map[string]chan models.Candlestick),
		candlestickTimeframes: make(map[string]map[string]*models.Candlestick),
		orderBookPrecisions:   make(map[string]orderBookPrecisions),
		Tickers:               make(map[string]*ticker),
		db:                    db,
		logger:                logger,
	}
}

func (s *Service) InitCandleStickTimeframes(pair string, candlestickTimeframes []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.candlestickTimeframes[pair] = make(map[string]*models.Candlestick)
	for _, tf := range candlestickTimeframes {
		s.candlestickTimeframes[pair][tf] = &models.Candlestick{
			Pair:      pair,
			Timeframe: tf,
		}
	}
}

func (s *Service) DeleteTimeframes(pair string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.candlestickTimeframes, pair)
}

func (s *Service) InitPrecisions(pair string, pricePrecisions []int32, qtyPrecision int32) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orderBookPrecisions[pair] = orderBookPrecisions{pricePrecisions: pricePrecisions, qtyPrecision: qtyPrecision}
}

func (s *Service) DeletePrecisions(pair string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.orderBookPrecisions, pair)
}

func (s *Service) InitPairs(pair string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.PrecisedOBSs[pair] = make(chan map[int32]models.OrderBookSnapshot)
	s.PrecisedTrades[pair] = make(chan []models.Trade)
	s.Candlestick[pair] = make(chan models.Candlestick)
	s.Ticker[pair] = make(chan models.Ticker)
	s.Tickers[pair] = &ticker{pair: pair}
}

func (s *Service) RemovePairs(pair string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if ch, exists := s.PrecisedOBSs[pair]; exists {
		close(ch)
		delete(s.PrecisedOBSs, pair)
	}

	if ch, exists := s.PrecisedTrades[pair]; exists {
		close(ch)
		delete(s.PrecisedTrades, pair)
	}

	if ch, exists := s.Candlestick[pair]; exists {
		close(ch)
		delete(s.Candlestick, pair)
	}

	if ch, exists := s.Ticker[pair]; exists {
		close(ch)
		delete(s.Ticker, pair)
		delete(s.Tickers, pair)
	}
}
