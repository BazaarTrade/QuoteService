package service

import (
	"time"

	"github.com/BazaarTrade/QuoteService/internal/models"
	"github.com/shopspring/decimal"
)

func (s *Service) CandleStickTick(pair string, streamHub *StreamHub) {
	defer streamHub.wg.Done()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()

		for timeFrame, candlestick := range streamHub.candlestickByTimeframe {
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
					return
				}

				candlestick.ID = ID

				select {
				case <-streamHub.ctx.Done():
					s.logger.Info("stopped candleStick tick", "pair", pair)
					return
				case streamHub.CandlestickChan <- *candlestick:
				default:
				}

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
}

func (s *Service) CandleStickFormation(trades []models.Trade) {
	s.mu.RLock()
	hub, exists := s.Streams[trades[0].Pair]
	s.mu.RUnlock()

	if !exists {
		s.logger.Error("failed to find stream hub", "pair", trades[0].Pair)
		return
	}

	for _, trade := range trades {
		for _, candlestick := range hub.candlestickByTimeframe {
			if candlestick.HighPrice.IsZero() || trade.Price.GreaterThan(candlestick.HighPrice) {
				candlestick.HighPrice = trade.Price
			}

			if candlestick.LowPrice.IsZero() || trade.Price.LessThan(candlestick.LowPrice) {
				candlestick.LowPrice = trade.Price
			}

			candlestick.ClosePrice = trade.Price
			candlestick.Volume = candlestick.Volume.Add(trade.Qty)
			candlestick.Turnover = candlestick.Turnover.Add(trade.Price.Mul(trade.Qty))

			hub.CandlestickChan <- *candlestick
		}
	}
}
