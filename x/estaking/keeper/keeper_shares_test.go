package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *EstakingKeeperTestSuite) TestKeeperShares() {
	testCases := []struct {
		name                 string
		prerequisiteFunction func() (addr sdk.AccAddress)
		postValidateFunction func(addr sdk.AccAddress)
	}{
		{
			"calc delegation amount",
			func() (addr sdk.AccAddress) {
				suite.ResetSuite()

				addr = suite.AddAccounts(1, nil)[0]

				return addr
			},
			func(addr sdk.AccAddress) {
				// Check with non-delegator
				delegatedAmount := suite.app.EstakingKeeper.CalcDelegationAmount(suite.ctx, addr)
				suite.Require().Equal(delegatedAmount, math.ZeroInt())

				// Check with genesis account (delegator)
				delegatedAmount = suite.app.EstakingKeeper.CalcDelegationAmount(suite.ctx, suite.genAccount)
				suite.Require().Equal(delegatedAmount, sdk.DefaultPowerReduction)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			addr := tc.prerequisiteFunction()
			tc.postValidateFunction(addr)
		})
	}
}
