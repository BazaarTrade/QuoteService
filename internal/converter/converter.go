package converter

import (
	"github.com/BazaarTrade/MatchingEngineProtoGen/pbM"
	"github.com/BazaarTrade/QuoteProtoGen/pbQ"
	"github.com/BazaarTrade/QuoteService/internal/models.go"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ModelsTickerToProtoTicker(ticker models.Ticker) *pbQ.Ticker {
	return &pbQ.Ticker{
		Pair:      ticker.Pair,
		LastPrice: ticker.LastPrice.String(),
		Change:    ticker.Change.String(),
		HighPrice: ticker.HighPrice.String(),
		LowPrice:  ticker.LowPrice.String(),
		Volume:    ticker.Volume.String(),
		Turnover:  ticker.Turnover.String(),
	}
}

func ModelsTradesToPbqTrades(trades []models.Trade) *pbQ.Trades {
	var pbQTrades = &pbQ.Trades{Trades: make([]*pbQ.Trade, 0, len(trades))}
	for _, trade := range trades {
		pbQTrades.Trades = append(pbQTrades.Trades, &pbQ.Trade{
			Pair:  trade.Pair,
			IsBid: trade.IsBid,
			Price: trade.Price.String(),
			Qty:   trade.Qty.String(),
			Time:  timestamppb.New(trade.Time),
		})
	}
	return pbQTrades
}

func PbMTradesToModelsTrades(pbMTrades *pbM.Trades) []models.Trade {
	var trades = make([]models.Trade, 0, len(pbMTrades.Trades))
	for _, pbMTrade := range pbMTrades.Trades {
		trades = append(trades, models.Trade{
			Pair:  pbMTrade.Pair,
			IsBid: pbMTrade.IsBid,
			Price: decimal.RequireFromString(pbMTrade.Price),
			Qty:   decimal.RequireFromString(pbMTrade.Qty),
			Time:  pbMTrade.Time.AsTime(),
		})
	}
	return trades
}

func ModelsCandleStickToPbQCandleStick(candleStick models.Candlestick) *pbQ.CandleStick {
	return &pbQ.CandleStick{
		ID:         int64(candleStick.ID),
		Pair:       candleStick.Pair,
		Timeframe:  candleStick.Timeframe,
		OpenTime:   timestamppb.New(candleStick.OpenTime),
		CloseTime:  timestamppb.New(candleStick.CloseTime),
		OpenPrice:  candleStick.OpenPrice.String(),
		ClosePrice: candleStick.ClosePrice.String(),
		HighPrice:  candleStick.HighPrice.String(),
		LowPrice:   candleStick.LowPrice.String(),
		Volume:     candleStick.Volume.String(),
		Turnover:   candleStick.Turnover.String(),
		IsClosed:   candleStick.IsClosed,
	}
}

func PbMOBSToModelsOBS(OBS *pbM.OrderBookSnapshot) models.OrderBookSnapshot {
	var bids = make([]models.Limit, 0, len(OBS.Bids))
	for _, bid := range OBS.Bids {
		bids = append(bids, models.Limit{
			Price: decimal.RequireFromString(bid.Price),
			Qty:   decimal.RequireFromString(bid.Qty),
		})
	}

	var asks = make([]models.Limit, 0, len(OBS.Asks))
	for _, ask := range OBS.Asks {
		asks = append(asks, models.Limit{
			Price: decimal.RequireFromString(ask.Price),
			Qty:   decimal.RequireFromString(ask.Qty),
		})
	}

	return models.OrderBookSnapshot{
		Pair:    OBS.Pair,
		Bids:    bids,
		Asks:    asks,
		BidsQty: decimal.RequireFromString(OBS.BidsQty),
		AsksQty: decimal.RequireFromString(OBS.AsksQty),
	}
}

func ModelsOBSToPbQOBS(OBS models.OrderBookSnapshot) *pbQ.OrderBookSnapshot {
	var bids = make([]*pbQ.Limit, 0, len(OBS.Bids))
	for _, bid := range OBS.Bids {
		bids = append(bids, &pbQ.Limit{
			Price: bid.Price.String(),
			Qty:   bid.Qty.String(),
		})
	}

	var asks = make([]*pbQ.Limit, 0, len(OBS.Asks))
	for _, ask := range OBS.Asks {
		asks = append(asks, &pbQ.Limit{
			Price: ask.Price.String(),
			Qty:   ask.Qty.String(),
		})
	}

	return &pbQ.OrderBookSnapshot{
		Pair:    OBS.Pair,
		Bids:    bids,
		Asks:    asks,
		BidsQty: OBS.BidsQty.String(),
		AsksQty: OBS.AsksQty.String(),
	}
}
