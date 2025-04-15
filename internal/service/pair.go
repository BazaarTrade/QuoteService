package service

import "github.com/BazaarTrade/QuoteService/internal/models.go"

func (s *Service) GetPairsParams() ([]models.PairParams, error) {
	return s.db.GetPairsParams()
}
