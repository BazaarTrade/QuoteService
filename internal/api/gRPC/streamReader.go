package server

import (
	"context"

	"github.com/BazaarTrade/MatchingEngineProtoGen/pbM"
	"github.com/BazaarTrade/QuoteService/internal/models.go"
)

func (s *Server) ReadOrderBookSnapshotStream() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := s.pbMClient.StreamOrderBookSnapshot(ctx, &pbM.Ping{})
	if err != nil {
		s.logger.Error("failed to connect to pbM OBS stream", "error", err)
		return
	}

	s.logger.Info("successfully connected to pbM OBS stream")

	for {
		pbOBS, err := stream.Recv()
		if err != nil {
			s.logger.Error("failed to receive pbOBS", "error", err)
			return
		}

		var OBS = models.OrderBookSnapshot{
			Symbol:  pbOBS.Symbol,
			Bids:    make([]models.Limit, len(pbOBS.Bids)),
			Asks:    make([]models.Limit, len(pbOBS.Asks)),
			BidsQty: pbOBS.BidsQty,
			AsksQty: pbOBS.AsksQty,
		}

		for i, limit := range pbOBS.Bids {
			OBS.Bids[i] = models.Limit{
				Price: limit.Price,
				Qty:   limit.Qty,
			}
		}

		for i, limit := range pbOBS.Asks {
			OBS.Asks[i] = models.Limit{
				Price: limit.Price,
				Qty:   limit.Qty,
			}
		}

		s.logger.Info("recieved OBS", "OBS", OBS)

		POBSs, err := s.Service.PreciseOrderBookSnaphot(OBS)
		if err != nil {
			continue
		}

		s.POBSs <- POBSs
	}
}

func (s *Server) ReadTrades() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := s.pbMClient.StreamTrades(ctx, &pbM.Ping{})
	if err != nil {
		s.logger.Error("failed to connect to pbM trades stream", "error", err)
		return
	}

	s.logger.Info("successfully connected to pbM trades stream")

	for {
		trades, err := stream.Recv()
		if err != nil {
			s.logger.Error("failed to receive trades", "error", err)
			return
		}

		s.logger.Info("recieved trades")

		precisedTrades, err := s.Service.TradesPrecision(trades)
		if err != nil {
			continue
		}

		s.PrecisedTrades <- precisedTrades
	}
}
