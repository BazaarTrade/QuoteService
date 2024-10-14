package models

import "time"

type OrderBookSnapshot struct {
	Symbol  string
	Bids    []Limit
	Asks    []Limit
	BidsQty string
	AsksQty string
}

type Limit struct {
	Price string
	Qty   string
}

type Trades struct {
	Symbol string
	Trades []Trade
}

type Trade struct {
	IsBid bool
	Price string
	Qty   string
	Time  time.Time
}
