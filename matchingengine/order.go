package matchingengine

import (
	"fmt"
	"time"
)

type Order struct {
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
func NewOrder(isBid bool, size float64) *Order {
	return &Order{
		Amount:    size,
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
