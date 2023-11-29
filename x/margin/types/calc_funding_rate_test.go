package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/elys-network/elys/x/margin/types"
	"github.com/stretchr/testify/assert"
)

func TestCalculateFundingRate(t *testing.T) {
	baseRate := sdk.NewDecWithPrec(3, 4) // 0.03%
	maxRate := sdk.NewDecWithPrec(1, 3)  // 0.1%
	minRate := sdk.NewDecWithPrec(-1, 3) // -0.1%

	// Test cases
	tests := []struct {
		name         string
		longAmount   sdk.Int
		shortAmount  sdk.Int
		expectedRate string
	}{
		{
			name:         "Longs More Leveraged",
			longAmount:   sdk.NewInt(1000000),
			shortAmount:  sdk.NewInt(500000),
			expectedRate: "0.000600000000000000",
		},
		{
			name:         "Shorts More Leveraged",
			longAmount:   sdk.NewInt(500000),
			shortAmount:  sdk.NewInt(1000000),
			expectedRate: "-0.000600000000000000",
		},
		{
			name:         "Balanced Leveraging",
			longAmount:   sdk.NewInt(750000),
			shortAmount:  sdk.NewInt(750000),
			expectedRate: "0.000300000000000000",
		},
		{
			name:         "Extreme Long Leverage",
			longAmount:   sdk.NewInt(2000000),
			shortAmount:  sdk.NewInt(500000),
			expectedRate: "0.001000000000000000", // Capped at maxRate
		},
		{
			name:         "Zero Short Amount",
			longAmount:   sdk.NewInt(1000000),
			shortAmount:  sdk.NewInt(0),
			expectedRate: "0.001000000000000000", // maxRate when short amount is zero
		},
		{
			name:         "Zero Long Amount",
			longAmount:   sdk.NewInt(0),
			shortAmount:  sdk.NewInt(1000000),
			expectedRate: "0.001000000000000000", // maxRate when long amount is zero
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualRate := types.CalcFundingRate(tt.longAmount, tt.shortAmount, baseRate, maxRate, minRate)
			assert.Equal(t, tt.expectedRate, actualRate.String())
		})
	}
}
