package service

import (
	"github.com/BazaarTrade/MatchingEngineProtoGen/pbM"
	"github.com/BazaarTrade/QuoteService/internal/models.go"
	"github.com/shopspring/decimal"
)

// make different precision variations of order book snapshot
func (s *Service) OrderBookSnapshotPreciosion(OBS models.OrderBookSnapshot) (map[int32]models.OrderBookSnapshot, error) {
	POBSs := make(map[int32]models.OrderBookSnapshot)

	for _, precision := range s.pricePrecisions[OBS.Symbol] {
		var POBS = models.OrderBookSnapshot{
			Symbol:  OBS.Symbol,
			BidsQty: OBS.BidsQty,
			AsksQty: OBS.AsksQty,
			Bids:    make([]models.Limit, 0),
			Asks:    make([]models.Limit, 0),
		}

		if len(OBS.Bids) > 0 {
			bidLimitPrice, err := decimal.NewFromString(OBS.Bids[0].Price)
			if err != nil {
				s.logger.Error("Error creating decimal from string", "error", err)
				return nil, err
			}
			precisedBidLimitPrice := bidLimitPrice.RoundFloor(precision)

			var bidLimitQty decimal.Decimal

			for i, limit := range OBS.Bids {
				limitPrice, err := decimal.NewFromString(limit.Price)
				if err != nil {
					s.logger.Error("Error creating decimal from string", "error", err)
					return nil, err
				}

				if limitPrice.GreaterThanOrEqual(precisedBidLimitPrice) {
					decimalQty, err := decimal.NewFromString(limit.Qty)
					if err != nil {
						s.logger.Error("Error creating decimal from string", "error", err)
						return nil, err
					}
					bidLimitQty = bidLimitQty.Add(decimalQty)

					if i+1 == len(OBS.Bids) {
						POBS.Bids = append(POBS.Bids, models.Limit{
							Price: precisedBidLimitPrice.String(),
							Qty:   bidLimitQty.Truncate(s.qtyPrecisions[OBS.Symbol]).String(),
						})
						break
					}
				} else {
					POBS.Bids = append(POBS.Bids, models.Limit{
						Price: precisedBidLimitPrice.String(),
						Qty:   bidLimitQty.Truncate(s.qtyPrecisions[OBS.Symbol]).String(),
					})
					precisedBidLimitPrice = limitPrice.RoundFloor(precision)
					bidLimitQty, err = decimal.NewFromString(limit.Qty)
					if err != nil {
						s.logger.Error("Error creating decimal from string", "error", err)
						return nil, err
					}

					if i+1 == len(OBS.Bids) {
						decimalQty, err := decimal.NewFromString(limit.Qty)
						if err != nil {
							s.logger.Error("Error creating decimal from string", "error", err)
							return nil, err
						}
						POBS.Bids = append(POBS.Bids, models.Limit{
							Price: precisedBidLimitPrice.String(),
							Qty:   decimalQty.Truncate(s.qtyPrecisions[OBS.Symbol]).String(),
						})
						break
					}
				}
				if len(POBS.Bids) == 30 {
					break
				}
			}
		}

		if len(OBS.Asks) > 0 {
			askLimitPrice, err := decimal.NewFromString(OBS.Asks[0].Price)
			if err != nil {
				s.logger.Error("Error creating decimal from string", "error", err)
				return nil, err
			}
			precisedAskLimitPrice := askLimitPrice.RoundCeil(precision)

			var askLimitQty decimal.Decimal

			for i, limit := range OBS.Asks {
				limitPrice, err := decimal.NewFromString(limit.Price)
				if err != nil {
					s.logger.Error("Error creating decimal from string", "error", err)
					return nil, err
				}

				if limitPrice.LessThanOrEqual(precisedAskLimitPrice) {
					decimalQty, err := decimal.NewFromString(limit.Qty)
					if err != nil {
						s.logger.Error("Error creating decimal from string", "error", err)
						return nil, err
					}
					askLimitQty = askLimitQty.Add(decimalQty)

					if i+1 == len(OBS.Asks) {
						POBS.Asks = append(POBS.Asks, models.Limit{
							Price: precisedAskLimitPrice.String(),
							Qty:   askLimitQty.Truncate(s.qtyPrecisions[OBS.Symbol]).String(),
						})
						break
					}
				} else {
					POBS.Asks = append(POBS.Asks, models.Limit{
						Price: precisedAskLimitPrice.String(),
						Qty:   askLimitQty.Truncate(s.qtyPrecisions[OBS.Symbol]).String(),
					})
					precisedAskLimitPrice = limitPrice.RoundCeil(precision)
					askLimitQty, err = decimal.NewFromString(limit.Qty)
					if err != nil {
						s.logger.Error("Error creating decimal from string", "error", err)
						return nil, err
					}

					if i+1 == len(OBS.Asks) {
						decimalQty, err := decimal.NewFromString(limit.Qty)
						if err != nil {
							s.logger.Error("Error creating decimal from string", "error", err)
							return nil, err
						}

						POBS.Asks = append(POBS.Asks, models.Limit{
							Price: precisedAskLimitPrice.String(),
							Qty:   decimalQty.Truncate(s.qtyPrecisions[OBS.Symbol]).String(),
						})
						break
					}
				}
				if len(POBS.Asks) == 30 {
					break
				}
			}
		}
		POBSs[precision] = POBS
	}
	return POBSs, nil
}

func (s *Service) TradesPrecision(trades *pbM.Trades) (*pbM.Trades, error) {
	for _, trade := range trades.Trade {
		decimalPrice, err := decimal.NewFromString(trade.Price)
		if err != nil {
			s.logger.Error("Error creating decimal from string", "error", err)
			return nil, err
		}

		decimalQty, err := decimal.NewFromString(trade.Qty)
		if err != nil {
			s.logger.Error("Error creating decimal from string", "error", err)
			return nil, err
		}

		trade.Price = decimalPrice.Truncate(s.qtyPrecisions[trades.Symbol]).String()
		trade.Qty = decimalQty.Truncate(s.pricePrecisions[trades.Symbol][0]).String()
	}
	return trades, nil
}
