package server

import (
	"context"

	"github.com/BazaarTrade/MatchingEngineProtoGen/pbM"
	"github.com/BazaarTrade/QuoteProtoGen/pbQ"
	"github.com/BazaarTrade/QuoteService/internal/models.go"
)

func (s *Server) ReadTrades() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := s.pbMClient.StreamTrades(ctx, &pbM.Ping{})
	if err != nil {
		s.logger.Error("Failed to connect to pbM Trades stream", "error", err)
		return
	}

	s.logger.Info("Successfully connected to pbM Trades stream")

	for {
		trades, err := stream.Recv()
		if err != nil {
			s.logger.Error("Failed to receive Trades", "error", err)
			return
		}

		s.logger.Info("Recieved trades")

		precisedTrades, err := s.Service.TradesPrecision(trades)
		if err != nil {
			continue
		}

		s.PrecisedTrades <- precisedTrades
	}
}

func (s *Server) StreamPrecisedTrades(req *pbQ.Ping, stream pbQ.Quote_StreamPrecisedTradesServer) error {
	for {
		pbMPrecisedTrades := <-s.PrecisedTrades

		var pbQPrecisedTrades = &pbQ.Trades{
			Symbol: pbMPrecisedTrades.Symbol,
		}

		for _, trade := range pbMPrecisedTrades.Trade {
			pbQPrecisedTrades.Trade = append(pbQPrecisedTrades.Trade, &pbQ.Trade{
				IsBid: trade.IsBid,
				Price: trade.Price,
				Qty:   trade.Qty,
				Time:  trade.Time,
			})
		}

		err := stream.Send(pbQPrecisedTrades)
		if err != nil {
			s.logger.Error("Failed to send precisedTrades", "error", err)
		}

		s.logger.Info("Sent trades", "trades", pbQPrecisedTrades)
	}
}
func (s *Server) ReadOrderBookSnapshotStream() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := s.pbMClient.StreamOrderBookSnapshot(ctx, &pbM.Ping{})
	if err != nil {
		s.logger.Error("Failed to connect to pbM OBS stream", "error", err)
		return
	}

	s.logger.Info("Successfully connected to pbM OBS stream")

	for {
		pbOBS, err := stream.Recv()
		if err != nil {
			s.logger.Error("Failed to receive pbOBS", "error", err)
			return
		}

		var OBS = models.OrderBookSnapshot{
			Symbol:  pbOBS.Symbol,
			BidsQty: pbOBS.BidsQty,
			AsksQty: pbOBS.AsksQty,
		}

		for _, bidLimit := range pbOBS.Bids {
			OBS.Bids = append(OBS.Bids, models.Limit{
				Price: bidLimit.Price,
				Qty:   bidLimit.Qty,
			})
		}

		for _, askLimit := range pbOBS.Asks {
			OBS.Asks = append(OBS.Asks, models.Limit{
				Price: askLimit.Price,
				Qty:   askLimit.Qty,
			})
		}

		s.logger.Info("Recieved OBS", "OBS", OBS)

		POBSs, err := s.Service.OrderBookSnapshotPreciosion(OBS)
		if err != nil {
			continue
		}

		s.POBSs <- POBSs
	}
}

func (s *Server) StreamPrecisedOrderBookSnapshot(req *pbQ.Ping, stream pbQ.Quote_StreamPrecisedOrderBookSnapshotServer) error {
	for {
		POBSs := <-s.POBSs

		var pbPOBS = pbQ.PrecisedOrderBookSnapshots{
			PrecisedOrderBookSnapshot: make(map[int32]*pbQ.OrderBookSnapshot),
		}

		for precision, OBS := range POBSs {
			var pbOBS = pbQ.OrderBookSnapshot{
				Symbol:  OBS.Symbol,
				BidsQty: OBS.BidsQty,
				AsksQty: OBS.AsksQty,
			}

			for _, bidLimit := range OBS.Bids {
				pbOBS.Bids = append(pbOBS.Bids, &pbQ.Limit{
					Price: bidLimit.Price,
					Qty:   bidLimit.Qty,
				})
			}

			for _, askLimit := range OBS.Asks {
				pbOBS.Asks = append(pbOBS.Asks, &pbQ.Limit{
					Price: askLimit.Price,
					Qty:   askLimit.Qty,
				})
			}

			pbPOBS.PrecisedOrderBookSnapshot[precision] = &pbOBS
		}

		err := stream.Send(&pbPOBS)
		if err != nil {
			s.logger.Error("Failed to send POBSs", "error", err)
		}

		s.logger.Info("Sent POBSs", "POBSs", POBSs)
	}
}
