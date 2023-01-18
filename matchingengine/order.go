package matchingengine

import (
	"fmt"
	"math/rand"
	"time"
)

type Order struct {
	ID        uint64  // The ID concept is only for external APIs
	UserID    uint64  // UserID to identify who puts the order
	Amount    float64 // Amount of our crypto
	Bid       bool    // Is this a sell or buy Order
	Limit     *Limit  // To keep track of what limit this order is set in
	Timestamp int64   // Use in64 because we will use Unix nano for Timestamp
}

type Orders []*Order

func (o Orders) Len() int           { return len(o) }
func (o Orders) Less(i, j int) bool { return o[i].Timestamp < o[j].Timestamp }
func (o Orders) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }

// NewOrder is constructor of Order struct
func NewOrder(isBid bool, amount float64, userID uint64) *Order {
	return &Order{
		ID:        uint64(rand.Intn(1000000)),
		UserID:    userID,
		Amount:    amount,
		Bid:       isBid,
		Timestamp: time.Now().UnixNano(),
	}
}

func (o *Order) String() string {
	return fmt.Sprintf("[amount: %.2f]", o.Amount)
}

func (o *Order) IsFilled() bool {
	return o.Amount == 0.0
}
