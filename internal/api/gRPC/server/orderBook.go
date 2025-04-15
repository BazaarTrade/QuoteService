package server

import (
	"context"

	"github.com/BazaarTrade/QuoteProtoGen/pbQ"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) CreateOrderBook(ctx context.Context, pairParams *pbQ.PairParams) (*emptypb.Empty, error) {
	s.logger.Info("InitOrderBook request", "symbol", pairParams.Pair)

	s.service.InitPairs(pairParams.Pair)
	s.service.InitPrecisions(pairParams.Pair, pairParams.PricePrecisions, pairParams.QtyPrecision)
	s.service.InitCandleStickTimeframes(pairParams.Pair, pairParams.CandleStickTimeframes)
	s.mClient.StartStreamReaders(pairParams.Pair)

	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteOrderBook(ctx context.Context, pair *pbQ.Pair) (*emptypb.Empty, error) {
	s.logger.Info("DeleteOrderBook request", "symbol", pair.Pair)

	s.service.RemovePairs(pair.Pair)
	s.service.DeletePrecisions(pair.Pair)
	s.service.DeleteTimeframes(pair.Pair)

	return &emptypb.Empty{}, nil
}
