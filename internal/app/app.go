package app

import (
	"log/slog"
	"os"

	mClient "github.com/BazaarTrade/QuoteService/internal/api/gRPC/matchingEngineClient"
	"github.com/BazaarTrade/QuoteService/internal/api/gRPC/server"
	"github.com/BazaarTrade/QuoteService/internal/repository/postgresPgx"
	"github.com/BazaarTrade/QuoteService/internal/service"
	"github.com/joho/godotenv"
)

func Run() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	if _, err := os.Stat("../.env"); err == nil {
		if err := godotenv.Load("../.env"); err != nil {
			logger.Error("failed to load .env file", "error", err)
			return
		}
	}

	logger.Info("starting aplication...")

	repo, err := postgresPgx.NewPostgres(logger)
	if err != nil {
		return
	}

	service := service.New(repo, logger)

	mClient := mClient.New(service, logger)
	if err := mClient.Run(); err != nil {
		return
	}

	defer mClient.CloseConnection()

	initOrderBooks(mClient, service)

	server := server.New(mClient, service, logger)
	if err := server.Run(); err != nil {
		return
	}
}

func initOrderBooks(mClient *mClient.Client, service *service.Service) error {
	pairsParams, err := service.GetPairsParams()
	if err != nil {
		return err
	}

	for _, pairParams := range pairsParams {
		service.InitPairs(pairParams.Pair)
		service.InitCandleStickTimeframes(pairParams.Pair, pairParams.CandleStickTimeframes)
		service.InitPrecisions(pairParams.Pair, pairParams.PricePrecisions, pairParams.QtyPrecision)
		mClient.StartStreamReaders(pairParams.Pair)
		go service.CandleStickTick(pairParams.Pair)
		go service.TickerTick(pairParams.Pair)
	}

	return nil
}
