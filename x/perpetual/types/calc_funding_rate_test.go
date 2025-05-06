package types_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/osmosis-labs/osmosis/osmomath"
	"github.com/stretchr/testify/assert"
)

func TestCalculateFundingRate(t *testing.T) {
	baseRate := osmomath.NewBigDecWithPrec(3, 4) // 0.03%
	maxRate := osmomath.NewBigDecWithPrec(1, 3)  // 0.1%
	minRate := osmomath.NewBigDecWithPrec(-1, 3) // -0.1%

	// Test cases
	tests := []struct {
		name         string
		longAmount   math.Int
		shortAmount  math.Int
		expectedRate string
	}{
		{
			name:         "Longs More Leveraged",
			longAmount:   math.NewInt(1000000),
			shortAmount:  math.NewInt(500000),
			expectedRate: "0.000600000000000000000000000000000000",
		},
		{
			name:         "Shorts More Leveraged",
			longAmount:   math.NewInt(500000),
			shortAmount:  math.NewInt(1000000),
			expectedRate: "-0.000600000000000000000000000000000000",
		},
		{
			name:         "Balanced Leveraging",
			longAmount:   math.NewInt(750000),
			shortAmount:  math.NewInt(750000),
			expectedRate: "0.000300000000000000000000000000000000",
		},
		{
			name:         "Extreme Long Leverage",
			longAmount:   math.NewInt(2000000),
			shortAmount:  math.NewInt(500000),
			expectedRate: "0.001000000000000000000000000000000000", // Capped at maxRate
		},
		{
			name:         "Zero Short Amount",
			longAmount:   math.NewInt(1000000),
			shortAmount:  math.NewInt(0),
			expectedRate: "0.001000000000000000000000000000000000", // maxRate when short amount is zero
		},
		{
			name:         "Zero Long Amount",
			longAmount:   math.NewInt(0),
			shortAmount:  math.NewInt(1000000),
			expectedRate: "0.001000000000000000000000000000000000", // maxRate when long amount is zero
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualRate := types.CalcFundingRate(tt.longAmount, tt.shortAmount, baseRate, maxRate, minRate)
			assert.Equal(t, tt.expectedRate, actualRate.String())
		})
	}
}
