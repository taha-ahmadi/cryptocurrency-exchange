package main

import (
	"log"
	"time"

	"github.com/taha-ahmadi/cryptocurrency-exchange/api"
	"github.com/taha-ahmadi/cryptocurrency-exchange/marketmaker"
)

func main() {
	//Start the server(exchange)
	go api.ServerEngine()

	mm := marketmaker.NewClinet()

	// we are going to place some orders in seedMarket() so we should wait for the server running
	time.Sleep(1 * time.Second)

	// SeedMarket to add liquidity
	if err := marketmaker.SeedMarket(mm); err != nil {
		log.Fatalf("Seed Market error: %v\n", err)
	}

	// market making algorithm
	go marketmaker.MakeMarketSimple(mm)

	time.Sleep(2 * time.Second)

	// regular users add some market order in this line
	marketmaker.MarketOrderPlacer(mm)

	select {}
}
