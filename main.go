package main

import (
	"github.com/snickers54/crypto-trading-bot/analysis"
	"github.com/snickers54/crypto-trading-bot/api"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load()
}

func main() {
	analysis.UpdateProducts()
	ordersTube := make(chan []api.MarketOrderPayload)
	go analysis.UpdateAccounts()
	go analysis.ListenDiscrepancy("ETC-BTC", ordersTube)
	go analysis.ListenDiscrepancy("ETH-BTC", ordersTube)
	go analysis.ListenDiscrepancy("LTC-BTC", ordersTube)
	go analysis.ListenDiscrepancy("BCH-BTC", ordersTube)
	go analysis.ListenDiscrepancy("ZRX-BTC", ordersTube)
	go analysis.ListenDiscrepancy("XRP-BTC", ordersTube)
	go analysis.ListenDiscrepancy("XLM-BTC", ordersTube)

	analysis.RealtimeFeed()
}
