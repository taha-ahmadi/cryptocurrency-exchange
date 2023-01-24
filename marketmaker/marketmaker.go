package marketmaker

import (
	"fmt"
	"log"
	"math"
	"time"
)

var (
	tick = 1 * time.Second
)

func MarketOrderPlacer(c *MMClient) {
	ticker := time.NewTicker(1 * time.Second)

	for {
		trades, err := c.GetTrades("ETH")
		if err != nil {
			panic(err)
		}

		if len(trades) > 0 {
			fmt.Printf("exchange price => %.2f\n", trades[len(trades)-1].Price)
		}

		otherMarketSell := &PlaceOrderParams{
			UserID: 8,
			Bid:    false,
			Amount: 1000,
		}
		orderResp, err := c.PlaceMarketOrder(otherMarketSell)
		if err != nil {
			log.Println(orderResp.OrderID)
		}

		marketSell := &PlaceOrderParams{
			UserID: 666,
			Bid:    false,
			Amount: 100,
		}
		orderResp, err = c.PlaceMarketOrder(marketSell)
		if err != nil {
			log.Println(orderResp.OrderID)
		}

		marketBuyOrder := &PlaceOrderParams{
			UserID: 666,
			Bid:    true,
			Amount: 100,
		}
		orderResp, err = c.PlaceMarketOrder(marketBuyOrder)
		if err != nil {
			log.Println(orderResp.OrderID)
		}

		<-ticker.C
	}
}

const userID = 7

func MakeMarketSimple(c *MMClient) {
	ticker := time.NewTicker(tick)

	for {
		orders, err := c.GetOrders(userID)

		if err != nil {
			log.Println(err)
		}

		bestAsk, err := c.GetBestAsk()
		if err != nil {
			log.Println(err)
		}
		bestBid, err := c.GetBestBid()
		if err != nil {
			log.Println(err)
		}

		spread := math.Abs(bestBid - bestAsk)
		fmt.Println("exchange spread", spread)

		log.Println("len of bids", len(orders.Bids))
		// place the bid
		if len(orders.Bids) < 3 {
			bidLimit := &PlaceOrderParams{
				UserID: 7,
				Bid:    true,
				Price:  bestBid + 100,
				Amount: 1000,
			}

			bidOrderResp, err := c.PlaceLimitOrder(bidLimit)
			if err != nil {
				log.Println(bidOrderResp.OrderID)
			}
		}

		// place the ask
		if len(orders.Asks) < 3 {
			askLimit := &PlaceOrderParams{
				UserID: 7,
				Bid:    false,
				Price:  bestAsk - 100,
				Amount: 1000,
			}

			askOrderResp, err := c.PlaceLimitOrder(askLimit)
			if err != nil {
				log.Println(askOrderResp.OrderID)
			}
		}

		fmt.Println("best ask price", bestAsk)
		fmt.Println("best bid price", bestBid)

		<-ticker.C
	}
}

func SeedMarket(c *MMClient) error {
	ask := &PlaceOrderParams{
		UserID: 8,
		Bid:    false,
		Price:  10_000,
		Amount: 1_0000,
	}

	bid := &PlaceOrderParams{
		UserID: 8,
		Bid:    true,
		Price:  9_000,
		Amount: 1_0000,
	}

	_, err := c.PlaceLimitOrder(ask)
	if err != nil {
		return err
	}

	_, err = c.PlaceLimitOrder(bid)
	if err != nil {
		return err
	}

	return nil
}
