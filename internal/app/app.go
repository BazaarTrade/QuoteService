package app

import (
	"log/slog"
	"os"

	server "github.com/BazaarTrade/QuoteService/internal/api/gRPC"
	"github.com/BazaarTrade/QuoteService/internal/service"
)

func Run() {
	handler := slog.NewTextHandler(os.Stdout, nil)
	logger := slog.New(handler)

	logger.Info("starting aplication...")

	service := service.New(logger)

	server := server.New(service, logger)

	err := server.RunGRPCClient()
	if err != nil {
		return
	}

	go server.ReadOrderBookSnapshotStream()
	go server.ReadTrades()

	err = server.RunGRPCServer()
	if err != nil {
		return
	}
}
