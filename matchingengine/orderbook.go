package matchingengine

import (
	"sort"
)

// Orderbook contains our asks and bids; orderbook would need to be persisted in a db somehow and shared between clients
// In real world exchanges we can use distributed event stream like Apache Kafka. By doing this, we could always replay
// the orders deterministically if the exchange crashes and restore the order books to their original state.
type Orderbook struct {
	asks []*Limit // If you want to sell a crypto for a certain size of crypto and certain price, you make an ask
	bids []*Limit // If you want to buy a crypto for a certain size of crypto and certain price, you make a bid
	// We use these slices because we can use our own implementation based on ask or bid to sort them
	// And sorting maps is harder than slices and needs to do more stuffs

	// We have no convenient way to check if there is already a limit order at a certain price level
	// We should loop through each slice and check if the price is same as we want but that will take too much time
	// So we will make map that point to specific limit
	AskLimits map[float64]*Limit
	BidLimits map[float64]*Limit
}

// NewOrderbook is constructor of Orderbook struct
func NewOrderbook() *Orderbook {
	return &Orderbook{
		asks: []*Limit{},
		bids: []*Limit{},

		AskLimits: make(map[float64]*Limit),
		BidLimits: make(map[float64]*Limit),
	}
}

// PlaceMarketOrder will fill the order with orderbook asks or bids, and also checks the volume for specific order request
func (ob *Orderbook) PlaceMarketOrder(o *Order) Matches {
	var matches Matches

	if o.Bid {
		// Check if the amount of the order is greater than the total volume of the ask orders
		if o.Amount > ob.AskTotalVolume() {
			panic("there is not enough volume in the orderbook")
		}
		// Iterate through all the ask orders
		for _, ask := range ob.Asks() {
			// Fill the ask order with the market order
			asksMatches := ask.Fill(o)
			matches = append(matches, asksMatches...)

			// Check if there are no more orders in the limit. we can keep limits without any orders but we will
			// remove it because of memory efficiency
			if len(ask.Orders) == 0 {
				ob.clearLimit(true, ask)
			}
		}
	} else {
		if o.Amount > ob.BidTotalVolume() {
			panic("there is not enough volume in the orderbook")
		}
		for _, bid := range ob.Bids() {
			bidsMatches := bid.Fill(o)
			matches = append(matches, bidsMatches...)

			if len(bid.Orders) == 0 {
				ob.clearLimit(false, bid)
			}
		}

	}
	return matches
}

func (ob *Orderbook) PlaceLimitOrder(price float64, o *Order) {
	// Check if already there are asks or bids volume sitting in the order book for specific price

	var limit *Limit
	if o.Bid {
		limit = ob.BidLimits[price]
	} else {
		limit = ob.AskLimits[price]
	}

	if limit == nil {
		limit = NewLimit(price)

		if o.Bid {
			ob.bids = append(ob.bids, limit)
			ob.BidLimits[price] = limit
		} else {
			ob.asks = append(ob.asks, limit)
			ob.AskLimits[price] = limit
		}
	}

	limit.AddOrder(o)
}

func (ob *Orderbook) clearLimit(isBid bool, l *Limit) {
	if isBid {
		delete(ob.BidLimits, l.Price)
		for i := 0; i < len(ob.bids); i++ {
			if ob.bids[i] == l {
				ob.bids[i] = ob.bids[len(ob.bids)-1]
				ob.bids = ob.bids[:len(ob.bids)-1]
			}
		}
	} else {
		delete(ob.AskLimits, l.Price)
		for i := 0; i < len(ob.asks); i++ {
			if ob.asks[i] == l {
				ob.asks[i] = ob.asks[len(ob.asks)-1]
				ob.asks = ob.asks[:len(ob.asks)-1]
			}
		}
	}
}

// CancelOrder will delete the order from the limit
func (ob *Orderbook) CancelOrder(o *Order) {
	o.Limit.DeleteOrder(o)
}

// BidTotalVolume returns total volume of the asks in the market
func (ob *Orderbook) BidTotalVolume() float64 {
	totalVolume := 0.0

	for i := 0; i < len(ob.bids); i++ {
		totalVolume += ob.bids[i].TotalVolume
	}

	return totalVolume
}

// AskTotalVolume returns total volume of the asks in the market
func (ob *Orderbook) AskTotalVolume() float64 {
	totalVolume := 0.0

	for i := 0; i < len(ob.asks); i++ {
		totalVolume += ob.asks[i].TotalVolume
	}

	return totalVolume
}

// Asks sort Orderbook asks and return it
func (ob *Orderbook) Asks() []*Limit {
	sort.Sort(ByBestAsk{ob.asks})
	return ob.asks
}

// Bids sort Orderbook bids and return it
func (ob *Orderbook) Bids() []*Limit {
	sort.Sort(ByBestBid{ob.bids})
	return ob.bids
}
