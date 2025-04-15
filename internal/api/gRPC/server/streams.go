package server

import (
	"errors"

	"github.com/BazaarTrade/QuoteProtoGen/pbQ"
	"github.com/BazaarTrade/QuoteService/internal/converter"
)

func (s *Server) StreamPrecisedOrderBookSnapshot(req *pbQ.Pair, stream pbQ.Quote_StreamPrecisedOrderBookSnapshotServer) error {
	pOBSsChan, exists := s.service.PrecisedOBSs[req.Pair]
	if !exists {
		s.logger.Error("failed to find pOBSs chan", "pair", req.Pair)
		return errors.New("failed to pOBSs chan for: " + req.Pair)
	}

	for {
		select {
		case pOBSs, ok := <-pOBSsChan:
			if !ok {
				s.logger.Info("pOBSs stream stopped manualy", "pair", req.Pair)
				return nil
			}

			var pbPOBSs = pbQ.PrecisedOrderBookSnapshots{
				PrecisedOrderBookSnapshot: make(map[int32]*pbQ.OrderBookSnapshot),
			}

			for precision, OBS := range pOBSs {
				pbPOBSs.PrecisedOrderBookSnapshot[precision] = converter.ModelsOBSToPbQOBS(OBS)
			}

			if err := stream.Send(&pbPOBSs); err != nil {
				s.logger.Error("failed to send pOBSs", "pOBSs", pOBSs, "error", err)
			}

			s.logger.Debug("sent pOBSs", "pair", req.Pair)

		case <-stream.Context().Done():
			s.logger.Info("precised OBSs stream stopped by client", "pair", req.Pair)
			return nil
		}
	}
}

func (s *Server) StreamPrecisedTrades(req *pbQ.Pair, stream pbQ.Quote_StreamPrecisedTradesServer) error {
	precisedTradeChan, exists := s.service.PrecisedTrades[req.Pair]
	if !exists {
		s.logger.Error("failed to find precised trade chan", "pair", req.Pair)
		return errors.New("failed to precised trade chan for: " + req.Pair)
	}

	for {
		select {
		case pbMPrecisedTrades, ok := <-precisedTradeChan:
			if !ok {
				s.logger.Info("precised trade stream stopped manualy", "pair", req.Pair)
				return nil
			}

			if err := stream.Send(converter.ModelsTradesToPbqTrades(pbMPrecisedTrades)); err != nil {
				s.logger.Error("failed to send precised trade", "error", err)
			}

			s.logger.Debug("sent precised trade", "pair", req.Pair)

		case <-stream.Context().Done():
			s.logger.Info("precised trade stream stopped by client", "pair", req.Pair)
			return nil
		}
	}
}

func (s *Server) StreamTicker(req *pbQ.Pair, stream pbQ.Quote_StreamTickerServer) error {
	tickerChan, exists := s.service.Ticker[req.Pair]
	if !exists {
		s.logger.Error("failed to find ticker chan", "pair", req.Pair)
		return errors.New("failed to ticker chan for: " + req.Pair)
	}

	for {
		select {
		case ticker, ok := <-tickerChan:
			if !ok {
				s.logger.Info("ticker stream stopped manualy", "pair", req.Pair)
				return nil
			}

			if err := stream.Send(converter.ModelsTickerToProtoTicker(ticker)); err != nil {
				s.logger.Error("failed to send ticker", "error", err)
			}

			s.logger.Debug("sent ticker", "pair", req.Pair)

		case <-stream.Context().Done():
			s.logger.Info("ticker stream stopped by client", "pair", req.Pair)
			return nil
		}
	}
}

func (s *Server) StreamCandleStick(req *pbQ.Pair, stream pbQ.Quote_StreamCandleStickServer) error {
	candlestickChan, exists := s.service.Candlestick[req.Pair]
	if !exists {
		s.logger.Error("failed to find candlestick chan", "pair", req.Pair)
		return errors.New("failed to candlestick chan for: " + req.Pair)
	}

	for {
		select {
		case candlestick, ok := <-candlestickChan:
			if !ok {
				s.logger.Info("candlestick stream stopped manualy", "pair", req.Pair)
				return nil
			}

			if err := stream.Send(converter.ModelsCandleStickToPbQCandleStick(candlestick)); err != nil {
				s.logger.Error("failed to send candlestick", "error", err)
			}

			s.logger.Debug("sent candleStick", "pair", req.Pair)

		case <-stream.Context().Done():
			s.logger.Info("candlestick stream stopped by client", "pair", req.Pair)
			return nil
		}
	}
}
