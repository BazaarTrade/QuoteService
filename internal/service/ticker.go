package service

import (
	"sync"
	"time"

	"github.com/BazaarTrade/QuoteService/internal/models.go"
	"github.com/shopspring/decimal"
)

type ticker struct {
	pair      string
	trades    []models.Trade
	lastPrice decimal.Decimal
	change    decimal.Decimal
	highPrice decimal.Decimal
	lowPrice  decimal.Decimal
	volume    decimal.Decimal
	turnover  decimal.Decimal

	mu sync.Mutex
}

func (s *Service) TickerTick(pair string) {
	s.mu.RLock()
	ticker, tickerExists := s.Tickers[pair]
	tickerChan, chanrExists := s.Ticker[pair]
	s.mu.RUnlock()

	if !tickerExists {
		s.logger.Error("failed to find ticker", "pair", pair)
		return
	}

	if !chanrExists {
		s.logger.Error("failed to find ticker chan", "pair", pair)
		return
	}

	waitTime := time.Until(time.Now().Truncate(time.Second).Add(time.Second))
	time.Sleep(waitTime)

	t := time.NewTicker(time.Second)
	defer t.Stop()

	for range t.C {
		ticker.mu.Lock()
		tickerChan <- ticker.cleanTrades()
		ticker.mu.Unlock()
	}
}

func (s *Service) TickerFormation(trades []models.Trade) {
	s.mu.RLock()
	ticker, tickerExists := s.Tickers[trades[0].Pair]
	s.mu.RUnlock()

	if !tickerExists {
		s.logger.Error("failed to find ticker", "pair", trades[0].Pair)
		return
	}

	ticker.addTrades(trades)
}

func (t *ticker) addTrades(trades []models.Trade) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, trade := range trades {
		t.trades = append(t.trades, trade)

		if t.highPrice.LessThan(trade.Price) || t.highPrice.IsZero() {
			t.highPrice = trade.Price
		}

		if t.lowPrice.GreaterThan(trade.Price) || t.lowPrice.IsZero() {
			t.lowPrice = trade.Price
		}

		t.volume = t.volume.Add(trade.Qty)
		t.turnover = t.turnover.Add(trade.Price.Mul(trade.Qty))
	}

	t.lastPrice = trades[len(trades)-1].Price

	//calculate change for 24h in %
	if len(t.trades) > 1 && !t.trades[0].Price.IsZero() {
		t.change = t.lastPrice.Sub(t.trades[0].Price).Div(t.trades[0].Price).Mul(decimal.NewFromInt(100)).Truncate(2)
	}
}

func (t *ticker) cleanTrades() models.Ticker {
	var startIndex = len(t.trades)
	for i, trade := range t.trades {
		if trade.Time.After(time.Now().Add(-24 * time.Hour)) {
			if i == 0 {
				return models.Ticker{
					Pair:      t.pair,
					LastPrice: t.lastPrice,
					Change:    t.change,
					HighPrice: t.highPrice,
					LowPrice:  t.lowPrice,
					Volume:    t.volume,
					Turnover:  t.turnover,
				}
			}

			startIndex = i
			break
		}
	}

	t.trades = t.trades[startIndex:]

	if len(t.trades) == 0 {
		t.lastPrice = decimal.Zero
		t.highPrice = decimal.Zero
		t.lowPrice = decimal.Zero
		t.volume = decimal.Zero
		t.turnover = decimal.Zero
		t.change = decimal.Zero
		return models.Ticker{
			Pair:      t.pair,
			LastPrice: t.lastPrice,
			Change:    t.change,
			HighPrice: t.highPrice,
			LowPrice:  t.lowPrice,
			Volume:    t.volume,
			Turnover:  t.turnover,
		}
	}

	t.lastPrice = t.trades[len(t.trades)-1].Price
	t.highPrice = t.trades[0].Price
	t.lowPrice = t.trades[0].Price
	t.volume = decimal.Zero
	t.turnover = decimal.Zero

	for _, trade := range t.trades {
		if t.highPrice.LessThan(trade.Price) || t.highPrice.IsZero() {
			t.highPrice = trade.Price
		}

		if t.lowPrice.GreaterThan(trade.Price) || t.lowPrice.IsZero() {
			t.lowPrice = trade.Price
		}

		t.volume = t.volume.Add(trade.Qty)
		t.turnover = t.turnover.Add(trade.Price.Mul(trade.Qty))
	}

	t.change = t.lastPrice.Sub(t.trades[0].Price).Div(t.trades[0].Price).Mul(decimal.NewFromInt(100)).Truncate(2)

	return models.Ticker{
		Pair:      t.pair,
		LastPrice: t.lastPrice,
		Change:    t.change,
		HighPrice: t.highPrice,
		LowPrice:  t.lowPrice,
		Volume:    t.volume,
		Turnover:  t.turnover,
	}
}
