package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/keeper"
)

func (suite *KeeperTestSuite) TestApplyDiscount() {
	// Define test cases
	tests := []struct {
		name     string
		swapFee  types.Dec
		discount types.Dec
		wantFee  types.Dec
	}{
		{
			name:     "Zero discount",
			swapFee:  types.NewDecWithPrec(100, 2), // 1.00 as an example
			discount: types.ZeroDec(),
			wantFee:  types.NewDecWithPrec(100, 2),
		},
		{
			name:     "Positive discount with valid broker address",
			swapFee:  types.NewDecWithPrec(100, 2),
			discount: types.NewDecWithPrec(10, 2), // 0.10 (10%)
			wantFee:  types.NewDecWithPrec(90, 2), // 0.90 after discount
		},
		{
			name:     "Boundary value for discount",
			swapFee:  types.NewDecWithPrec(100, 2),
			discount: types.NewDecWithPrec(9999, 4), // 0.9999 (99.99%)
			wantFee:  types.NewDecWithPrec(1, 4),    // 0.01 after discount
		},
		{
			name:     "Discount greater than swap fee",
			swapFee:  types.NewDecWithPrec(50, 2), // 0.50
			discount: types.NewDecWithPrec(75, 2), // 0.75
			wantFee:  types.NewDecWithPrec(125, 3),
		},
		{
			name:     "Zero swap fee with valid discount",
			swapFee:  types.ZeroDec(),
			discount: types.NewDecWithPrec(10, 2),
			wantFee:  types.ZeroDec(),
		},
		{
			name:     "Large discount with valid broker address",
			swapFee:  types.NewDecWithPrec(100, 2),
			discount: types.NewDecWithPrec(9000, 4), // 0.90 (90%)
			wantFee:  types.NewDecWithPrec(10, 2),   // 0.10 after discount
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			fee := keeper.ApplyDiscount(tc.swapFee, tc.discount)
			suite.Require().Equal(tc.wantFee, fee)
		})
	}
}
