package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/stretchr/testify/require"
)

func TestApplyDiscount(t *testing.T) {
	// Mock context and keeper
	k, ctx, _, _ := keepertest.AmmKeeper(t)

	// Define test cases
	tests := []struct {
		name      string
		swapFee   types.Dec
		discount  types.Dec
		sender    string
		wantFee   types.Dec
		wantError bool
	}{
		{
			name:     "Zero discount",
			swapFee:  types.NewDecWithPrec(100, 2), // 1.00 as an example
			discount: types.ZeroDec(),
			sender:   "testSender",
			wantFee:  types.NewDecWithPrec(100, 2),
		},
		{
			name:     "Positive discount with valid broker address",
			swapFee:  types.NewDecWithPrec(100, 2),
			discount: types.NewDecWithPrec(10, 2), // 0.10 (10%)
			sender:   k.BrokerAddress(ctx),
			wantFee:  types.NewDecWithPrec(90, 2), // 0.90 after discount
		},
		{
			name:      "Positive discount with invalid broker address",
			swapFee:   types.NewDecWithPrec(100, 2),
			discount:  types.NewDecWithPrec(10, 2),
			sender:    "invalidBroker",
			wantError: true,
		},
		{
			name:     "Boundary value for discount",
			swapFee:  types.NewDecWithPrec(100, 2),
			discount: types.NewDecWithPrec(9999, 4), // 0.9999 (99.99%)
			sender:   k.BrokerAddress(ctx),
			wantFee:  types.NewDecWithPrec(1, 4), // 0.01 after discount
		},
		{
			name:     "Discount greater than swap fee",
			swapFee:  types.NewDecWithPrec(50, 2), // 0.50
			discount: types.NewDecWithPrec(75, 2), // 0.75
			sender:   k.BrokerAddress(ctx),
			wantFee:  types.NewDecWithPrec(125, 3),
		},
		{
			name:      "Invalid swap fee",
			swapFee:   types.NewDecWithPrec(-100, 2), // -1.00 (invalid)
			discount:  types.NewDecWithPrec(10, 2),
			sender:    "testSender",
			wantError: true,
		},
		{
			name:     "Zero swap fee with valid discount",
			swapFee:  types.ZeroDec(),
			discount: types.NewDecWithPrec(10, 2),
			sender:   k.BrokerAddress(ctx),
			wantFee:  types.ZeroDec(),
		},
		{
			name:     "Large discount with valid broker address",
			swapFee:  types.NewDecWithPrec(100, 2),
			discount: types.NewDecWithPrec(9000, 4), // 0.90 (90%)
			sender:   k.BrokerAddress(ctx),
			wantFee:  types.NewDecWithPrec(10, 2), // 0.10 after discount
		},
		{
			name:      "Large discount with invalid broker address",
			swapFee:   types.NewDecWithPrec(100, 2),
			discount:  types.NewDecWithPrec(9000, 4), // 0.90
			sender:    "invalidBroker",
			wantError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fee, discount, err := k.ApplyDiscount(ctx, tc.swapFee, tc.discount, tc.sender)

			if tc.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.wantFee, fee)
				require.Equal(t, tc.discount, discount)
			}
		})
	}
}
