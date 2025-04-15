package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type OrderBookSnapshot struct {
	Pair    string
	Bids    []Limit
	Asks    []Limit
	BidsQty decimal.Decimal
	AsksQty decimal.Decimal
}

type Limit struct {
	Price decimal.Decimal
	Qty   decimal.Decimal
}

type Trade struct {
	Pair  string
	IsBid bool
	Price decimal.Decimal
	Qty   decimal.Decimal
	Time  time.Time
}

type PairParams struct {
	Pair                  string
	PricePrecisions       []int32
	QtyPrecision          int32
	CandleStickTimeframes []string
}

type Ticker struct {
	Pair      string
	LastPrice decimal.Decimal
	Change    decimal.Decimal
	HighPrice decimal.Decimal
	LowPrice  decimal.Decimal
	Volume    decimal.Decimal
	Turnover  decimal.Decimal
}

type Candlestick struct {
	ID         int
	Pair       string
	Timeframe  string
	OpenTime   time.Time
	CloseTime  time.Time
	OpenPrice  decimal.Decimal
	ClosePrice decimal.Decimal
	HighPrice  decimal.Decimal
	LowPrice   decimal.Decimal
	Volume     decimal.Decimal
	Turnover   decimal.Decimal
	IsClosed   bool
}

type TickerInfo struct {
	HighPrice string
	LowPrice  string
	Volume    string
	Turnover  string
}
