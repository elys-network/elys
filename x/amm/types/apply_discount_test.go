package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/stretchr/testify/require"
)

func TestApplyDiscount(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		swapFee  sdk.Dec
		discount sdk.Dec
		wantFee  sdk.Dec
	}{
		{
			name:     "Zero discount",
			swapFee:  sdk.NewDecWithPrec(100, 2), // 1.00 as an example
			discount: sdk.ZeroDec(),
			wantFee:  sdk.NewDecWithPrec(100, 2),
		},
		{
			name:     "Positive discount with valid broker address",
			swapFee:  sdk.NewDecWithPrec(100, 2),
			discount: sdk.NewDecWithPrec(10, 2), // 0.10 (10%)
			wantFee:  sdk.NewDecWithPrec(90, 2), // 0.90 after discount
		},
		{
			name:     "Boundary value for discount",
			swapFee:  sdk.NewDecWithPrec(100, 2),
			discount: sdk.NewDecWithPrec(9999, 4), // 0.9999 (99.99%)
			wantFee:  sdk.NewDecWithPrec(1, 4),    // 0.01 after discount
		},
		{
			name:     "Discount greater than swap fee",
			swapFee:  sdk.NewDecWithPrec(50, 2), // 0.50
			discount: sdk.NewDecWithPrec(75, 2), // 0.75
			wantFee:  sdk.NewDecWithPrec(125, 3),
		},
		{
			name:     "Zero swap fee with valid discount",
			swapFee:  sdk.ZeroDec(),
			discount: sdk.NewDecWithPrec(10, 2),
			wantFee:  sdk.ZeroDec(),
		},
		{
			name:     "Large discount with valid broker address",
			swapFee:  sdk.NewDecWithPrec(100, 2),
			discount: sdk.NewDecWithPrec(9000, 4), // 0.90 (90%)
			wantFee:  sdk.NewDecWithPrec(10, 2),   // 0.10 after discount
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fee := types.ApplyDiscount(tc.swapFee, tc.discount)
			require.Equal(t, tc.wantFee, fee)
		})
	}
}
