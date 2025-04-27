package repository

import "github.com/BazaarTrade/QuoteService/internal/models"

type Repository interface {
	CreateCandleStick(models.Candlestick) (int, error)
	GetPairsParams() ([]models.PairParams, error)
	GetTickerInfo24H(string) (string, []models.TickerInfo, error)
}
