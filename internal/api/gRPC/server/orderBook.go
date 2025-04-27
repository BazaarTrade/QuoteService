package server

import (
	"context"

	"github.com/BazaarTrade/QuoteProtoGen/pbQ"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) CreateOrderBook(ctx context.Context, pairParams *pbQ.PairParams) (*emptypb.Empty, error) {
	s.logger.Info("InitOrderBook request", "symbol", pairParams.Pair)

	s.service.NewStreamHub(pairParams.Pair, pairParams.PricePrecisions, pairParams.QtyPrecision, pairParams.CandleStickTimeframes)
	s.mClient.StartStreamReaders(pairParams.Pair)

	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteOrderBook(ctx context.Context, pair *pbQ.Pair) (*emptypb.Empty, error) {
	s.logger.Info("DeleteOrderBook request", "symbol", pair.Pair)

	s.mClient.StopStreamReadersByPair(pair.Pair)
	s.service.DeleteStreamHubByPair(pair.Pair)

	return &emptypb.Empty{}, nil
}
