package matchingengine

// Match struct holds information about a matched order in the matching engine.
// Ask and Bid are pointers to Order structs representing the ask and bid orders that were matched.
// AmountFilled is a float64 representing the size of the match.
// Price is also a float64 representing the price of the match.
type Match struct {
	Ask          *Order
	Bid          *Order
	AmountFilled float64
	Price        float64
}

type Matches []Match

// Limit is group of orders at the certain price level.
type Limit struct {
	Price       float64
	Orders      Orders
	TotalVolume float64
}

type Limits []*Limit

// ByBestAsk sell highest.
type ByBestAsk struct{ Limits }

func (a ByBestAsk) Len() int           { return len(a.Limits) }
func (a ByBestAsk) Less(i, j int) bool { return a.Limits[i].Price < a.Limits[j].Price }
func (a ByBestAsk) Swap(i, j int)      { a.Limits[i], a.Limits[j] = a.Limits[j], a.Limits[i] }

// ByBestBid Buy cheapest.
type ByBestBid struct{ Limits }

func (b ByBestBid) Len() int           { return len(b.Limits) }
func (b ByBestBid) Less(i, j int) bool { return b.Limits[i].Price > b.Limits[j].Price }
func (b ByBestBid) Swap(i, j int)      { b.Limits[i], b.Limits[j] = b.Limits[j], b.Limits[i] }

// NewLimit is constructor of Limit struct
func NewLimit(price float64) *Limit {
	return &Limit{
		Price:  price,
		Orders: []*Order{},
	}
}

// AddOrder will add order with certain price to Limit and increase TotalVolume of crypto in the certain price.
func (l *Limit) AddOrder(o *Order) {
	o.Limit = l
	l.Orders = append(l.Orders, o)
	l.TotalVolume += o.Amount
}

func (l *Limit) DeleteOrder(o *Order) {
	for i := 0; i < len(l.Orders); i++ {
		if l.Orders[i] == o {
			// Removed element but in efficient and unordered way.
			l.Orders[i] = l.Orders[len(l.Orders)-1]
			l.Orders = l.Orders[:len(l.Orders)-1]
		}
	}
	o.Limit = nil // we put its value to nil for garbage collector.
	l.TotalVolume -= o.Amount
}

// Fill the order with orders in the specific Limit.
func (l *Limit) Fill(o *Order) Matches {
	var matches Matches
	var ordersToDelete Orders

	for _, order := range l.Orders {
		if o.IsFilled() {
			break
		}

		match := l.fillOrder(order, o)
		matches = append(matches, match)

		l.TotalVolume -= match.AmountFilled

		if order.IsFilled() {
			ordersToDelete = append(ordersToDelete, order)
		}
	}

	// It is safer to delete orders outside the loop because we are iterating over limit(l) orders and deleting them
	// inside the loop can be dangerous.
	for _, order := range ordersToDelete {
		l.DeleteOrder(order)
	}

	return matches
}

func (l *Limit) fillOrder(a, b *Order) Match {
	var (
		bid, ask   *Order
		sizeFilled float64
	)

	// Determine which order is the bid and which is the ask.
	if a.Bid {
		bid, ask = a, b
	} else {
		bid, ask = b, a
	}

	// Check the amount of the first order and the second order.
	if a.Amount >= b.Amount {
		// If the first order is greater than the second order, subtract the second order amount from the first order.
		a.Amount -= b.Amount
		sizeFilled = b.Amount
		b.Amount = 0.0
	} else {
		// If the second order is greater than the first order, subtract the first order amount from the second order.
		b.Amount -= a.Amount
		sizeFilled = a.Amount
		a.Amount = 0.0
	}

	// Create and return a Match struct to find out matches for specific order.
	return Match{
		Ask:          ask,
		Bid:          bid,
		Price:        l.Price,
		AmountFilled: sizeFilled,
	}
}
