package service

import (
	"github.com/BazaarTrade/QuoteService/internal/models.go"
	"github.com/shopspring/decimal"
)

type orderBookPrecisions struct {
	pricePrecisions []int32
	qtyPrecision    int32
}

func (s *Service) PreciseOrderBookSnaphot(OBS models.OrderBookSnapshot) {
	preciseLimits := func(limits []models.Limit, precision int32, isBid bool) []models.Limit {
		precisedLimits := make([]models.Limit, 0)

		if len(limits) == 0 {
			return precisedLimits
		}

		var limitQty, precisedLimitPrice decimal.Decimal

		switch {
		case isBid:
			precisedLimitPrice = limits[0].Price.RoundFloor(precision)
		case !isBid:
			precisedLimitPrice = limits[0].Price.RoundCeil(precision)
		}

		for i, limit := range limits {
			if (isBid && limit.Price.GreaterThanOrEqual(precisedLimitPrice)) ||
				(!isBid && limit.Price.LessThanOrEqual(precisedLimitPrice)) {
				limitQty = limitQty.Add(limit.Qty)

				if i+1 == len(limits) {
					precisedLimits = append(precisedLimits, models.Limit{
						Price: precisedLimitPrice,
						Qty:   limitQty.Truncate(s.orderBookPrecisions[OBS.Pair].qtyPrecision),
					})
					break
				}
			} else {
				precisedLimits = append(precisedLimits, models.Limit{
					Price: precisedLimitPrice,
					Qty:   limitQty.Truncate(s.orderBookPrecisions[OBS.Pair].qtyPrecision),
				})

				if isBid {
					precisedLimitPrice = limit.Price.RoundFloor(precision)
				} else {
					precisedLimitPrice = limit.Price.RoundCeil(precision)
				}

				if i+1 == len(limits) {
					precisedLimits = append(precisedLimits, models.Limit{
						Price: precisedLimitPrice,
						Qty:   limit.Qty.Truncate(s.orderBookPrecisions[OBS.Pair].qtyPrecision),
					})
					break
				}
			}
			if len(precisedLimits) >= 30 {
				break
			}
		}
		return precisedLimits
	}

	pOBSs := make(map[int32]models.OrderBookSnapshot)

	for _, precision := range s.orderBookPrecisions[OBS.Pair].pricePrecisions {
		OBS.Bids = preciseLimits(OBS.Bids, precision, true)
		OBS.Asks = preciseLimits(OBS.Asks, precision, false)

		pOBS := models.OrderBookSnapshot{
			Pair:    OBS.Pair,
			BidsQty: OBS.BidsQty,
			AsksQty: OBS.AsksQty,
			Bids:    OBS.Bids,
			Asks:    OBS.Asks,
		}

		pOBSs[precision] = pOBS
	}

	s.PrecisedOBSs[OBS.Pair] <- pOBSs
}

func (s *Service) PreciseTrades(trades []models.Trade) {
	if len(trades) == 0 {
		s.logger.Error("trades are empty")
		return
	}

	for i := range trades {
		trades[i].Price = trades[i].Price.Truncate(6)
		trades[i].Qty = trades[i].Qty.Truncate(6)
	}

	pair := trades[0].Pair

	if ch, ok := s.PrecisedTrades[pair]; ok {
		ch <- trades
	} else {
		s.logger.Error("failed to find precised trades chan", "pair", pair)
		return
	}
}
