package postgresPgx

import (
	"context"

	"github.com/BazaarTrade/QuoteService/internal/models.go"
)

func (p *Postgres) CreateCandleStick(candlestick models.Candlestick) (int, error) {
	var ID int64
	err := p.db.QueryRow(context.Background(), `
	INSERT INTO quote.candlestick 
	(pair, timeframe, openTime, closeTime, openPrice, closePrice, highPrice, lowPrice, volume, turnover) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id
	`, candlestick.Pair, candlestick.Timeframe, candlestick.OpenTime, candlestick.CloseTime,
		candlestick.OpenPrice.String(), candlestick.ClosePrice.String(),
		candlestick.HighPrice.String(), candlestick.LowPrice.String(),
		candlestick.Volume.String(), candlestick.Turnover.String()).Scan(&ID)

	if err != nil {
		p.logger.Error("failed to insert candlestick", "error", err)
		return 0, err
	}

	return int(ID), nil
}

func (p *Postgres) GetTickerInfo24H(pair string) (string, []models.TickerInfo, error) {
	rows, err := p.db.Query(context.Background(), `
	SELECT highPrice, lowPrice, volume, turnover
	FROM quote.candlestick 
	WHERE pair = $1 AND openTime >= NOW() - INTERVAL '24 HOURS' AND timeFrame = '1s'
	`, pair)
	if err != nil {
		p.logger.Error("failed to get ticker info", "error", err)
		return "", nil, err
	}
	defer rows.Close()

	var tickerInfo24H []models.TickerInfo
	for rows.Next() {
		var tickerInfo models.TickerInfo
		err = rows.Scan(&tickerInfo.HighPrice, &tickerInfo.LowPrice, &tickerInfo.Volume, &tickerInfo.Turnover)
		if err != nil {
			p.logger.Error("failed to scan ticker info", "error", err)
			return "", nil, err
		}

		tickerInfo24H = append(tickerInfo24H, tickerInfo)
	}

	var closePrice string
	err = p.db.QueryRow(context.Background(), `
	SELECT closePrice
	FROM quote.candlestick
	WHERE pair = $1 AND timeframe = '1s'
	ORDER BY ABS(EXTRACT(EPOCH FROM (openTime - NOW() + INTERVAL '24 hours'))) ASC
	LIMIT 1;
	`, pair).Scan(&closePrice)
	if err != nil {
		p.logger.Error("failed to get close price", "error", err)
		return "", nil, err
	}

	return closePrice, tickerInfo24H, nil
}
