package service

import (
	"github.com/BazaarTrade/MatchingEngineProtoGen/pbM"
	"github.com/BazaarTrade/QuoteService/internal/models.go"
	"github.com/shopspring/decimal"
)

func (s *Service) PreciseOrderBookSnaphot(OBS models.OrderBookSnapshot) (map[int32]models.OrderBookSnapshot, error) {
	POBSs := make(map[int32]models.OrderBookSnapshot)

	preciseLimits := func(limits []models.Limit, precision int32, isBid bool) ([]models.Limit, error) {
		precisedLimits := make([]models.Limit, 0)
		var limitQty decimal.Decimal

		if len(limits) == 0 {
			return precisedLimits, nil
		}

		initialPrice, err := decimal.NewFromString(limits[0].Price)
		if err != nil {
			s.logger.Error("failed to create decimal from string", "error", err)
			return nil, err
		}

		var precisedLimitPrice decimal.Decimal
		if isBid {
			precisedLimitPrice = initialPrice.RoundFloor(precision)
		} else {
			precisedLimitPrice = initialPrice.RoundCeil(precision)
		}

		for i, limit := range limits {
			limitPrice, err := decimal.NewFromString(limit.Price)
			if err != nil {
				s.logger.Error("failed to create decimal from string", "error", err)
				return nil, err
			}

			if (isBid && limitPrice.GreaterThanOrEqual(precisedLimitPrice)) ||
				(!isBid && limitPrice.LessThanOrEqual(precisedLimitPrice)) {
				decimalQty, err := decimal.NewFromString(limit.Qty)
				if err != nil {
					s.logger.Error("failed to create decimal from string", "error", err)
					return nil, err
				}
				limitQty = limitQty.Add(decimalQty)

				if i+1 == len(limits) {
					precisedLimits = append(precisedLimits, models.Limit{
						Price: precisedLimitPrice.String(),
						Qty:   limitQty.Truncate(s.qtyPrecisions[OBS.Symbol]).String(),
					})
					break
				}
			} else {
				precisedLimits = append(precisedLimits, models.Limit{
					Price: precisedLimitPrice.String(),
					Qty:   limitQty.Truncate(s.qtyPrecisions[OBS.Symbol]).String(),
				})

				if isBid {
					precisedLimitPrice = limitPrice.RoundFloor(precision)
				} else {
					precisedLimitPrice = limitPrice.RoundCeil(precision)
				}

				limitQty, err = decimal.NewFromString(limit.Qty)
				if err != nil {
					s.logger.Error("failed to create decimal from string", "error", err)
					return nil, err
				}

				if i+1 == len(limits) {
					precisedLimits = append(precisedLimits, models.Limit{
						Price: precisedLimitPrice.String(),
						Qty:   limitQty.Truncate(s.qtyPrecisions[OBS.Symbol]).String(),
					})
					break
				}
			}
			if len(precisedLimits) >= 30 {
				break
			}
		}
		return precisedLimits, nil
	}

	for _, precision := range s.pricePrecisions[OBS.Symbol] {
		bids, err := preciseLimits(OBS.Bids, precision, true)
		if err != nil {
			return nil, err
		}

		asks, err := preciseLimits(OBS.Asks, precision, false)
		if err != nil {
			return nil, err
		}

		POBS := models.OrderBookSnapshot{
			Symbol:  OBS.Symbol,
			BidsQty: OBS.BidsQty,
			AsksQty: OBS.AsksQty,
			Bids:    bids,
			Asks:    asks,
		}

		POBSs[precision] = POBS
	}

	return POBSs, nil
}

func (s *Service) TradesPrecision(trades *pbM.Trades) (*pbM.Trades, error) {
	for _, trade := range trades.Trade {
		decimalPrice, err := decimal.NewFromString(trade.Price)
		if err != nil {
			s.logger.Error("failed to create decimal from string", "error", err)
			return nil, err
		}

		decimalQty, err := decimal.NewFromString(trade.Qty)
		if err != nil {
			s.logger.Error("failed to create decimal from string", "error", err)
			return nil, err
		}

		trade.Price = decimalPrice.Truncate(s.qtyPrecisions[trades.Symbol]).String()
		trade.Qty = decimalQty.Truncate(s.pricePrecisions[trades.Symbol][0]).String()
	}
	return trades, nil
}
