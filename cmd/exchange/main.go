package main

import (
	"log"
	"time"

	"github.com/taha-ahmadi/cryptocurrency-exchange/cmd/exchange/server"
	"github.com/taha-ahmadi/cryptocurrency-exchange/internal/config"
	"github.com/taha-ahmadi/cryptocurrency-exchange/internal/marketmaker"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Start the server in a goroutine
	go func() {
		srv, err := server.NewServer(cfg)
		if err != nil {
			log.Fatalf("Failed to create server: %v", err)
		}

		if err := srv.Start(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for the server to start
	time.Sleep(1 * time.Second)

	// Initialize market maker client
	mm := marketmaker.NewClinet()

	// SeedMarket to add liquidity
	if err := marketmaker.SeedMarket(mm); err != nil {
		log.Fatalf("Seed Market error: %v\n", err)
	}

	// Start market making algorithm in a goroutine
	go marketmaker.MakeMarketSimple(mm)

	time.Sleep(2 * time.Second)

	// Regular users add some market orders
	marketmaker.MarketOrderPlacer(mm)

	// Keep the program running
	select {}
}
