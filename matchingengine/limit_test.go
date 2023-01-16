package matchingengine

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDeleteOrder(t *testing.T) {
	l := NewLimit(1000)
	o1 := NewOrder(true, 2)
	o2 := NewOrder(false, 3)
	l.AddOrder(o1)
	l.AddOrder(o2)

	// Test case 1: delete an existing order
	l.DeleteOrder(o1)
	require.Equal(t, 1, len(l.Orders))
	require.Equal(t, 3.0, l.TotalVolume)

	// Test case 2: delete an order that does not exist in the limit
	o3 := NewOrder(true, 2)
	l.AddOrder(o3)
	l.DeleteOrder(o3)

	require.Equal(t, 1, len(l.Orders))
	require.Equal(t, 3.0, l.TotalVolume)

	// Test case 3: delete the last order in the limit
	l.DeleteOrder(o2)
	require.Equal(t, 0, len(l.Orders))
	require.Equal(t, 0.0, l.TotalVolume)
}

func TestFill(t *testing.T) {
	l := NewLimit(1000)
	o1 := NewOrder(true, 2)
	o2 := NewOrder(false, 3)
	l.AddOrder(o1)
	l.AddOrder(o2)

	// Test case 1: fill a sell order with a buy order
	o3 := NewOrder(false, 5)
	matches := l.Fill(o3)
	require.Equal(t, 2, len(matches))
	require.Equal(t, 2.0, matches[0].AmountFilled)
	require.Equal(t, 0.0, l.TotalVolume)

	// Test case 2: fill a buy order with multiple sell orders
	o4 := NewOrder(true, 5)
	matches = l.Fill(o4)

	require.Equal(t, 0, len(matches))
}
