package server

import (
	"log/slog"
	"net"

	"github.com/BazaarTrade/QuoteProtoGen/pbQ"
	mClient "github.com/BazaarTrade/QuoteService/internal/api/gRPC/matchingEngineClient"
	"github.com/BazaarTrade/QuoteService/internal/service"
	"google.golang.org/grpc"
)

type Server struct {
	pbQ.UnimplementedQuoteServer
	grpcServer *grpc.Server
	mClient    *mClient.Client
	service    *service.Service
	logger     *slog.Logger
}

func New(mClient *mClient.Client, service *service.Service, logger *slog.Logger) *Server {
	return &Server{
		service: service,
		mClient: mClient,
		logger:  logger,
	}
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		s.logger.Error("failed to listen", "error", err)
		return err
	}

	grpcServer := grpc.NewServer()

	s.grpcServer = grpcServer

	pbQ.RegisterQuoteServer(grpcServer, s)
	s.logger.Info("server is listening on port 50052...")

	if err := grpcServer.Serve(lis); err != nil {
		s.logger.Error("failed to serve", "err", err)
		return err
	}
	return nil
}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
}
