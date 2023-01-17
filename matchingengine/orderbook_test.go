package matchingengine

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPlaceMarketOrder(t *testing.T) {
	// Initialize an Orderbook
	ob := NewOrderbook()

	// Add some asks and bids to the Orderbook
	sellOrder1 := NewOrder(false, 5.0)
	sellOrder2 := NewOrder(false, 8.0)
	buyOrder := NewOrder(true, 30)
	ob.PlaceLimitOrder(120, sellOrder1)
	ob.PlaceLimitOrder(100, sellOrder2)
	ob.PlaceLimitOrder(220, buyOrder)
	require.Equal(t, 3, len(ob.Orders))
	require.Equal(t, sellOrder1, ob.Orders[sellOrder1.ID])

	// Test case 1: Place a market buy order with amount 30
	buyMarketOrder := NewOrder(true, 10.0)
	matches := ob.PlaceMarketOrder(buyMarketOrder)

	// check if the order is filled
	if !buyMarketOrder.IsFilled() {
		t.Error("Expected buy order to be filled but it was not")
	}
	// check if the order matches are correct
	if len(matches) != 2 {
		t.Errorf("Expected 2 matches but got %d", len(matches))
	}
	if matches[0].AmountFilled != 8.0 || matches[1].AmountFilled != 2.0 {
		t.Errorf("Expected matches to be of size 10 and 20 but got %f and %f", matches[0].AmountFilled, matches[1].AmountFilled)
	}

	// Test case 2: Place a market sell order with amount 50
	sellOrder3 := NewOrder(false, 3)
	matches = ob.PlaceMarketOrder(sellOrder3)

	// check if the order is filled
	if !sellOrder3.IsFilled() {
		t.Error("Expected sell order to be filled but it was not")
	}
	// check if the order matches are correct
	if len(matches) != 1 {
		t.Errorf("Expected 2 matches but got %d", len(matches))
	}

	if matches[0].AmountFilled != 3.0 {
		t.Errorf("Expected matches to be of size 3 but got %f", matches[0].AmountFilled)
	}

	// Test case 3: Place a sell market order with amount greater than the total volume of bid Orders
	sellOrder4 := NewOrder(false, 100)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	ob.PlaceMarketOrder(sellOrder4)

	// Test case 4: Place a buy market order with amount less than the total volume of ask Orders
	// and check if the matches returned are correct
	asks1 := NewOrder(false, 50)
	ob.PlaceLimitOrder(100, asks1)
	ask2 := NewOrder(false, 30)
	ob.PlaceLimitOrder(110, ask2)
	ask3 := NewOrder(false, 20)
	ob.PlaceLimitOrder(120, ask3)
	buyOrder = NewOrder(true, 50)
	matches2 := ob.PlaceMarketOrder(buyOrder)
	if len(matches2) != 1 || matches2[0].AmountFilled != 50 {
		t.Error("Expected one match with size filled of 50, got: ", matches2)
	}
}

func TestPlaceLimitOrder(t *testing.T) {
	// Initialize an orderbook
	ob := NewOrderbook()

	// Create an order
	sellOrder := NewOrder(false, 10)

	// Place the sellOrder
	ob.PlaceLimitOrder(100, sellOrder)

	// Check if the order was added to the correct limit
	require.Equal(t, 1, len(ob.AskLimits[100].Orders))

	// Create another order with the same price
	secondSellOrder := NewOrder(false, 5)
	ob.PlaceLimitOrder(100, secondSellOrder)

	// Check if the second order was added to the same limit
	require.Equal(t, 2, len(ob.AskLimits[100].Orders))
	require.Equal(t, 1, len(ob.asks))

	// Create an order
	buyOrder := NewOrder(true, 10)

	// Place the buyOrder
	ob.PlaceLimitOrder(100, buyOrder)

	// Check if the order was added to the correct limit
	require.Equal(t, 1, len(ob.BidLimits[100].Orders))

	// Create another order with the same price
	secondBuyOrder := NewOrder(true, 5)
	ob.PlaceLimitOrder(100, secondBuyOrder)

	// Check if the second order was added to the same limit
	require.Equal(t, 2, len(ob.BidLimits[100].Orders))
	require.Equal(t, 1, len(ob.bids))
}

func TestCancelOrder(t *testing.T) {
	// Initialize an empty orderbook
	ob := NewOrderbook()

	// Create three test Orders and remove one of them
	o1 := NewOrder(false, 10.0)
	o2 := NewOrder(false, 15.0)
	o3 := NewOrder(false, 35.0)
	ob.PlaceLimitOrder(100.0, o1)
	ob.PlaceLimitOrder(150.0, o2)
	ob.PlaceLimitOrder(200.0, o3)

	require.Equal(t, 3, len(ob.Orders))

	// Cancel the order
	ob.CancelOrder(o1)

	// Check if the order has been deleted from the limit
	require.Equal(t, 2, len(ob.Orders))
	_, ok := ob.Orders[o1.ID]
	require.False(t, ok)
	_, ok = ob.Orders[o2.ID]
	require.True(t, ok)
}

func TestBidTotalVolume(t *testing.T) {
	// Initialize an empty orderbook
	ob := NewOrderbook()

	// Create some test limits
	o1 := NewOrder(true, 15.0)
	o2 := NewOrder(true, 5.0)
	o3 := NewOrder(true, 15.0)

	// Add the limits to the orderbook's bids
	ob.PlaceLimitOrder(100, o1)
	ob.PlaceLimitOrder(100, o2)
	ob.PlaceLimitOrder(100, o3)

	// Calculate the total volume of the bid Orders
	totalVolume := ob.BidTotalVolume()

	// Check if the total volume is equal to 30
	require.Equal(t, 35.0, totalVolume)
}

func TestAskTotalVolume(t *testing.T) {
	// Initialize an empty orderbook
	ob := NewOrderbook()

	// Create some test limits
	o1 := NewOrder(false, 10.0)
	o2 := NewOrder(false, 5.0)
	o3 := NewOrder(false, 15.0)

	// Add the limits to the orderbook's asks
	ob.PlaceLimitOrder(100, o1)
	ob.PlaceLimitOrder(100, o2)
	ob.PlaceLimitOrder(100, o3)

	// Calculate the total volume of the ask Orders
	totalVolume := ob.AskTotalVolume()

	// Check if the total volume is equal to 30
	require.Equal(t, 30.0, totalVolume)
}

func TestAsks(t *testing.T) {
	// Initialize an empty orderbook
	ob := NewOrderbook()

	// Create some test limits
	l1 := NewLimit(100.0)
	l2 := NewLimit(50.0)
	l3 := NewLimit(150.0)

	// Add the limits to the orderbook's asks
	ob.asks = append(ob.asks, l1, l2, l3)

	// Retrieve the asks from the orderbook
	asks := ob.Asks()

	// Check if the asks are sorted by price
	if asks[0] != l2 || asks[1] != l1 || asks[2] != l3 {
		t.Errorf("Asks() = %v, expected %v", asks, []*Limit{l2, l1, l3})
	}
}

func TestBids(t *testing.T) {
	// Initialize an empty orderbook
	ob := NewOrderbook()

	// Create some test limits
	l1 := NewLimit(100.0)
	l2 := NewLimit(50.0)
	l3 := NewLimit(150.0)

	// Add the limits to the orderbook's bids
	ob.bids = append(ob.bids, l1, l2, l3)

	// Retrieve the bids from the orderbook
	bids := ob.Bids()

	// Check if the bids are sorted by price
	if bids[0] != l3 || bids[1] != l1 || bids[2] != l2 {
		t.Errorf("Bids() = %v, expected %v", bids, []*Limit{l2, l1, l3})
	}
}
