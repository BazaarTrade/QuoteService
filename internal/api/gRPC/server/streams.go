package server

import (
	"errors"

	"github.com/BazaarTrade/QuoteProtoGen/pbQ"
	"github.com/BazaarTrade/QuoteService/internal/converter"
)

func (s *Server) StreamPrecisedOrderBookSnapshots(req *pbQ.Pair, stream pbQ.Quote_StreamPrecisedOrderBookSnapshotsServer) error {
	hub, exists := s.service.Streams[req.Pair]
	if !exists {
		s.logger.Error("failed to find stream hub", "pair", req.Pair)
		return errors.New("failed to find stream hub for: " + req.Pair)
	}

	s.logger.Info("client connected to pOBSs stream", "pair", req.Pair)

	for {
		select {
		case pOBSs, ok := <-hub.PrecisedOrderBookSnapshotsChan:
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
			s.logger.Info("client disconnected from pOBSs stream", "pair", req.Pair)
			return nil
		}
	}
}

func (s *Server) StreamPrecisedTrades(req *pbQ.Pair, stream pbQ.Quote_StreamPrecisedTradesServer) error {
	hub, exists := s.service.Streams[req.Pair]
	if !exists {
		s.logger.Error("failed to find stream hub", "pair", req.Pair)
		return errors.New("failed to find stream hub for: " + req.Pair)
	}

	s.logger.Info("client connected to precised trades stream", "pair", req.Pair)

	for {
		select {
		case pbMPrecisedTrades, ok := <-hub.PrecisedTradesChan:
			if !ok {
				s.logger.Info("precised trade stream stopped manualy", "pair", req.Pair)
				return nil
			}

			if err := stream.Send(converter.ModelsTradesToPbqTrades(pbMPrecisedTrades)); err != nil {
				s.logger.Error("failed to send precised trade", "error", err)
			}

			s.logger.Debug("sent precised trade", "pair", req.Pair)

		case <-stream.Context().Done():
			s.logger.Info("client disconnected from precised trades", "pair", req.Pair)
			return nil
		}
	}
}

func (s *Server) StreamTicker(req *pbQ.Pair, stream pbQ.Quote_StreamTickerServer) error {
	hub, exists := s.service.Streams[req.Pair]
	if !exists {
		s.logger.Error("failed to find stream hub", "pair", req.Pair)
		return errors.New("failed to find stream hub for: " + req.Pair)
	}

	s.logger.Info("client connected to ticker stream", "pair", req.Pair)

	for {
		select {
		case ticker, ok := <-hub.TickerChan:
			if !ok {
				s.logger.Info("ticker stream stopped manualy", "pair", req.Pair)
				return nil
			}

			if err := stream.Send(converter.ModelsTickerToProtoTicker(ticker)); err != nil {
				s.logger.Error("failed to send ticker", "error", err)
			}

			s.logger.Debug("sent ticker", "pair", req.Pair)

		case <-stream.Context().Done():
			s.logger.Info("client disconnected from ticker stream", "pair", req.Pair)
			return nil
		}
	}
}

func (s *Server) StreamCandleStick(req *pbQ.Pair, stream pbQ.Quote_StreamCandleStickServer) error {
	hub, exists := s.service.Streams[req.Pair]
	if !exists {
		s.logger.Error("failed to find stream hub", "pair", req.Pair)
		return errors.New("failed to find stream hub for: " + req.Pair)
	}

	s.logger.Info("client connected to candlestick stream", "pair", req.Pair)

	for {
		select {
		case candlestick, ok := <-hub.CandlestickChan:
			if !ok {
				s.logger.Info("candlestick stream stopped manualy", "pair", req.Pair)
				return nil
			}

			if err := stream.Send(converter.ModelsCandleStickToPbQCandleStick(candlestick)); err != nil {
				s.logger.Error("failed to send candlestick", "error", err)
			}

			s.logger.Debug("sent candleStick", "pair", req.Pair)

		case <-stream.Context().Done():
			s.logger.Info("client disconnected from candleStick stream", "pair", req.Pair)
			return nil
		}
	}
}
