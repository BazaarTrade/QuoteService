package server

import (
	"log/slog"
	"net"
	"os"

	"github.com/BazaarTrade/MatchingEngineProtoGen/pbM"
	"github.com/BazaarTrade/QuoteProtoGen/pbQ"
	"github.com/BazaarTrade/QuoteService/internal/models.go"
	"github.com/BazaarTrade/QuoteService/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	Service        *service.Service
	POBSs          chan map[int32]models.OrderBookSnapshot
	PrecisedTrades chan *pbM.Trades
	logger         *slog.Logger
	pbMClient      pbM.MatchingEngineClient
	pbMConn        *grpc.ClientConn
	pbQ.UnimplementedQuoteServer
}

func New(service *service.Service, logger *slog.Logger) *Server {
	return &Server{
		Service:        service,
		POBSs:          make(chan map[int32]models.OrderBookSnapshot, 50),
		PrecisedTrades: make(chan *pbM.Trades, 50),
		logger:         logger,
	}
}

func (s *Server) RunGRPCServer() error {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		s.logger.Error("failed to listen", "error", err)
		return err
	}

	grpcServer := grpc.NewServer()

	pbQ.RegisterQuoteServer(grpcServer, s)
	s.logger.Info("server is listening on port 50052...")

	if err := grpcServer.Serve(lis); err != nil {
		s.logger.Error("failed to serve", "err", err)
		return err
	}

	return nil
}

func (s *Server) RunGRPCClient() error {
	conn, err := grpc.NewClient(os.Getenv("CONN_ADDR_MATCHING_ENGINE"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		s.logger.Error("GRPC error connecting client to localhost:50051", "error", err)
		return err
	}

	s.pbMClient = pbM.NewMatchingEngineClient(conn)

	return nil
}

func (s *Server) CloseConnection() {
	if s.pbMConn != nil {
		if err := s.pbMConn.Close(); err != nil {
			s.logger.Error("failed to close Matching Engine connection", "error", err)
		}
	}
}
