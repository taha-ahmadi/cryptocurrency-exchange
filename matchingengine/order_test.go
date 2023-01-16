package matchingengine

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOrder_IsFilled(t *testing.T) {
	tests := []struct {
		name     string
		amount   float64
		isFilled bool
	}{
		{"amount 0.0", 0.0, true},
		{"amount 10.0", 10.0, false},
		{"amount 0.0001", 0.0001, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			order := Order{Amount: test.amount}
			result := order.IsFilled()
			require.Equal(t, test.isFilled, result)
		})
	}
}
