package exchanges

// Market represents a trading market
type Market string

// OrderType represents a type of order
type OrderType string

const (
	// MarketETH represents the ETH market
	MarketETH Market = "ETH"
	// MarketBTC represents the BTC market
	MarketBTC Market = "BTC"

	// MarketOrder represents a market order type
	MarketOrder OrderType = "MARKET"
	// LimitOrder represents a limit order type
	LimitOrder OrderType = "LIMIT"
)

// PlaceOrderRequest is a data structure for placing orders via API
type PlaceOrderRequest struct {
	UserID uint64
	Type   OrderType
	IsBid  bool
	Amount float64
	Price  float64
	Market Market
}

// PlaceOrderResponse is a response for a successful order placement
type PlaceOrderResponse struct {
	OrderID uint64
}

// Order represents a simplified order for API responses
type Order struct {
	UserID    uint64
	ID        uint64
	Amount    float64
	IsBid     bool
	Price     float64
	Timestamp int64
}

// OrderbookResponse represents an orderbook for API responses
type OrderbookResponse struct {
	TotalAsksVolume float64
	TotalBidsVolume float64
	Asks            []*Order
	Bids            []*Order
}

// GetOrdersResponse represents a response to get orders API call
type GetOrdersResponse struct {
	Asks []Order
	Bids []Order
}

// PriceResponse represents a response with a price
type PriceResponse struct {
	Price float64
}

// MatchedOrder represents a matched order for API responses
type MatchedOrder struct {
	UserID       uint64
	Price        float64
	AmountFilled float64
	ID           uint64
}
