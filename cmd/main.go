package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	mClient "github.com/BazaarTrade/QuoteService/internal/api/gRPC/matchingEngineClient"
	"github.com/BazaarTrade/QuoteService/internal/api/gRPC/server"
	"github.com/BazaarTrade/QuoteService/internal/repository/postgresPgx"
	"github.com/BazaarTrade/QuoteService/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	// time.Sleep(time.Second * 6)

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	if _, err := os.Stat("../.env"); err == nil {
		if err := godotenv.Load("../.env"); err != nil {
			logger.Error("failed to load .env file", "error", err)
			return
		}
	}

	logger.Info("starting aplication...")

	DB_CONNECTION := os.Getenv("DB_CONNECTION")
	if DB_CONNECTION == "" {
		logger.Error("DB_CONNECTION environment variable is not set")
		return
	}

	repository, err := postgresPgx.New(DB_CONNECTION, logger)
	if err != nil {
		return
	}

	service := service.New(repository, logger)

	CONN_ADDR_MATCHING_ENGINE := os.Getenv("CONN_ADDR_MATCHING_ENGINE")
	if CONN_ADDR_MATCHING_ENGINE == "" {
		logger.Error("CONN_ADDR_MATCHING_ENGINE environment variable is not set")
		return
	}

	mClient := mClient.New(service, logger)
	if err := mClient.Run(CONN_ADDR_MATCHING_ENGINE); err != nil {
		return
	}

	GRPC_PORT := os.Getenv("GRPC_PORT")
	if GRPC_PORT == "" {
		logger.Error("ADDR environment variable is not set")
		return
	}

	server := server.New(mClient, service, logger)

	if err := initOrderBook(service, mClient); err != nil {
		return
	}

	go func() {
		if err := server.Run(GRPC_PORT); err != nil {
			os.Exit(1)
		}
	}()

	//Graceful shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	logger.Info("shutting down...")

	mClient.CloseConnection()
	logger.Info("closed mClient connection")

	server.Stop()
	logger.Info("stopped gRPC server")

	repository.Close()
	logger.Info("closed database connection")

	logger.Info("gracefully stopped")
}

func initOrderBook(service *service.Service, mClient *mClient.Client) error {
	pairsParams, err := service.GetPairsParams()
	if err != nil {
		return err
	}

	for _, pairParams := range pairsParams {
		service.NewStreamHub(pairParams.Pair, pairParams.PricePrecisions, pairParams.QtyPrecision, pairParams.CandleStickTimeframes)
		mClient.StartStreamReaders(pairParams.Pair)
	}

	return nil
}
