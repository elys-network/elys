package keeper_test

import (
	"github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestApplyDiscount() {
	k, ctx := suite.app.AmmKeeper, suite.ctx
	brokerAddress := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()).String()
	params := suite.app.ParameterKeeper.GetParams(ctx)
	params.BrokerAddress = brokerAddress
	suite.app.ParameterKeeper.SetParams(ctx, params)

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
			sender:   brokerAddress,
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
			sender:   brokerAddress,
			wantFee:  types.NewDecWithPrec(1, 4), // 0.01 after discount
		},
		{
			name:     "Discount greater than swap fee",
			swapFee:  types.NewDecWithPrec(50, 2), // 0.50
			discount: types.NewDecWithPrec(75, 2), // 0.75
			sender:   brokerAddress,
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
			sender:   brokerAddress,
			wantFee:  types.ZeroDec(),
		},
		{
			name:     "Large discount with valid broker address",
			swapFee:  types.NewDecWithPrec(100, 2),
			discount: types.NewDecWithPrec(9000, 4), // 0.90 (90%)
			sender:   brokerAddress,
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
		suite.Run(tc.name, func() {
			fee, discount, err := k.ApplyDiscount(ctx, tc.swapFee, tc.discount, tc.sender)

			if tc.wantError {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.wantFee, fee)
				suite.Require().Equal(tc.discount, discount)
			}
		})
	}
}
