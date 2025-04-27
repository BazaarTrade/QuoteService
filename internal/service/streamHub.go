package service

import (
	"context"
	"errors"
	"sync"

	"github.com/BazaarTrade/QuoteService/internal/models"
)

var (
	ErrStreamHubAlreadyExists = errors.New("stream hub already exists")
	ErrStreamHubNotFound      = errors.New("failed to find stream hub")
)

type StreamHub struct {
	PrecisedOrderBookSnapshotsChan chan map[int32]models.OrderBookSnapshot
	PrecisedTradesChan             chan []models.Trade
	TickerChan                     chan models.Ticker
	CandlestickChan                chan models.Candlestick

	ticker *ticker

	ctx    context.Context
	cancel context.CancelFunc

	wg sync.WaitGroup

	orderBookPrecisions    orderBookPrecisions
	candlestickByTimeframe map[string]*models.Candlestick
}

func (s *Service) NewStreamHub(pair string, pricePrecisions []int32, qtyPrecision int32, timeFrames []string) (*StreamHub, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.Streams[pair]; exists {
		return nil, ErrStreamHubAlreadyExists
	}

	ctx, cancel := context.WithCancel(context.Background())

	streamHub := &StreamHub{
		PrecisedOrderBookSnapshotsChan: make(chan map[int32]models.OrderBookSnapshot),
		PrecisedTradesChan:             make(chan []models.Trade),
		TickerChan:                     make(chan models.Ticker),
		CandlestickChan:                make(chan models.Candlestick),
		ticker:                         &ticker{pair: pair},
		ctx:                            ctx,
		cancel:                         cancel,
		orderBookPrecisions: orderBookPrecisions{
			price: pricePrecisions,
			qty:   qtyPrecision,
		},
		candlestickByTimeframe: make(map[string]*models.Candlestick),
	}

	for _, tf := range timeFrames {
		streamHub.candlestickByTimeframe[tf] = &models.Candlestick{
			Pair:      pair,
			Timeframe: tf,
		}
	}
	s.Streams[pair] = streamHub

	streamHub.wg.Add(1)
	go s.CandleStickTick(pair, streamHub)

	streamHub.wg.Add(1)
	go s.TickerTick(pair, streamHub)

	return streamHub, nil
}

func (s *Service) DeleteStreamHubs() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for pair, streamHub := range s.Streams {
		streamHub.cancel()
		streamHub.wg.Wait()
		close(streamHub.PrecisedOrderBookSnapshotsChan)
		close(streamHub.PrecisedTradesChan)
		close(streamHub.TickerChan)
		close(streamHub.CandlestickChan)
		delete(s.Streams, pair)
		s.logger.Info("deleted stream hub", "pair", pair)
	}
}

func (s *Service) DeleteStreamHubByPair(pair string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if streamHub, exists := s.Streams[pair]; exists {
		streamHub.cancel()
		streamHub.wg.Wait()
		close(streamHub.PrecisedOrderBookSnapshotsChan)
		close(streamHub.PrecisedTradesChan)
		close(streamHub.TickerChan)
		close(streamHub.CandlestickChan)
		delete(s.Streams, pair)
		return nil
	}

	return ErrStreamHubNotFound
}
