package server

import (
	"github.com/BazaarTrade/QuoteProtoGen/pbQ"
)

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
			s.logger.Error("failed to send POBSs", "error", err)
		}

		s.logger.Info("sent POBSs", "POBSs", POBSs)
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
			s.logger.Error("failed to send precised trades", "error", err)
		}

		s.logger.Info("sent precised trades", "trades", pbQPrecisedTrades)
	}
}
