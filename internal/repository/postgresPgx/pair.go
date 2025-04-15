package postgresPgx

import (
	"context"

	"github.com/BazaarTrade/QuoteService/internal/models.go"
)

func (p *Postgres) GetPairsParams() ([]models.PairParams, error) {
	rows, err := p.db.Query(context.Background(), `
	SELECT pair, orderBookPricePrecisions, qtyPrecision, candleStickTimeframes
	FROM quote.pairs
	`)
	if err != nil {
		p.logger.Error("failed to select pairs", "error", err)
		return nil, err
	}
	defer rows.Close()

	var pairsParams []models.PairParams
	for rows.Next() {
		var pairParams models.PairParams
		err := rows.Scan(&pairParams.Pair, &pairParams.PricePrecisions, &pairParams.QtyPrecision, &pairParams.CandleStickTimeframes)
		if err != nil {
			p.logger.Error("failed to scan pair", "error", err)
			return nil, err
		}
		pairsParams = append(pairsParams, pairParams)
	}
	return pairsParams, nil
}
