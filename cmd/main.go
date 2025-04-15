package main

import (
	"time"

	"github.com/BazaarTrade/QuoteService/internal/app"
)

func main() {
	time.Sleep(time.Second * 6)
	app.Run()
}
