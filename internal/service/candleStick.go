package service

import (
	"errors"
	"time"

	"github.com/BazaarTrade/QuoteService/internal/models.go"
	"github.com/shopspring/decimal"
)

func (s *Service) CandleStickTick(pair string) error {
	s.mu.RLock()
	timeFrames, exists := s.candlestickTimeframes[pair]
	s.mu.RUnlock()

	if !exists {
		s.logger.Error("timeframe not found", "pair", pair)
		return errors.New("timeframe not found")
	}

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()

		for timeFrame, candlestick := range timeFrames {
			parsedTimeframe, err := time.ParseDuration(timeFrame)
			if err != nil {
				s.logger.Error("failed to parse timeframe", "timeframe", timeFrame)
				continue
			}

			truncatedTime := now.Truncate(parsedTimeframe)

			if candlestick.OpenTime.Before(truncatedTime) {
				candlestick.CloseTime = truncatedTime
				candlestick.IsClosed = true
				ID, err := s.db.CreateCandleStick(*candlestick)
				if err != nil {
					return err
				}

				candlestick.ID = ID

				s.Candlestick[pair] <- *candlestick

				candlestick.OpenTime = truncatedTime
				candlestick.CloseTime = time.Time{}
				candlestick.OpenPrice = candlestick.ClosePrice
				candlestick.HighPrice = candlestick.ClosePrice
				candlestick.LowPrice = candlestick.ClosePrice
				candlestick.Volume = decimal.Zero
				candlestick.Turnover = decimal.Zero
				candlestick.IsClosed = false
			}
		}
	}
	return nil
}

func (s *Service) CandleStickFormation(trades []models.Trade) {
	s.mu.RLock()
	timeFrames, exists := s.candlestickTimeframes[trades[0].Pair]
	s.mu.RUnlock()
	if !exists {
		s.logger.Error("failed to find candle stick timeframes", "pair", trades[0].Pair)
		return
	}

	for _, trade := range trades {
		for _, candlestick := range timeFrames {
			if candlestick.HighPrice.IsZero() || trade.Price.GreaterThan(candlestick.HighPrice) {
				candlestick.HighPrice = trade.Price
			}

			if candlestick.LowPrice.IsZero() || trade.Price.LessThan(candlestick.LowPrice) {
				candlestick.LowPrice = trade.Price
			}

			candlestick.ClosePrice = trade.Price
			candlestick.Volume = candlestick.Volume.Add(trade.Qty)
			candlestick.Turnover = candlestick.Turnover.Add(trade.Price.Mul(trade.Qty))

			s.Candlestick[trade.Pair] <- *candlestick
		}
	}
}
