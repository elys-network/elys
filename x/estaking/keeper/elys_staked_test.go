package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/x/estaking/types"
)

func (suite *EstakingKeeperTestSuite) TestElysStaked() {
	testCases := []struct {
		name                 string
		prerequisiteFunction func() (addr sdk.AccAddress)
		postValidateFunction func(addr sdk.AccAddress)
	}{
		{
			"remove elys staked",
			func() (addr sdk.AccAddress) {
				suite.ResetSuite()

				addr = suite.AddAccounts(1, nil)[0]

				// set elys staked
				suite.app.EstakingKeeper.SetElysStaked(suite.ctx, types.ElysStaked{
					Address: addr.String(),
					Amount:  math.NewInt(1000),
				})

				// get elys staked
				elysStaked := suite.app.EstakingKeeper.GetElysStaked(suite.ctx, addr)
				suite.Require().Equal(math.NewInt(1000), elysStaked.Amount)

				return addr
			},
			func(addr sdk.AccAddress) {
				suite.app.EstakingKeeper.RemoveElysStaked(suite.ctx, addr)

				// get elys staked
				elysStaked := suite.app.EstakingKeeper.GetElysStaked(suite.ctx, addr)
				suite.Require().Equal(math.ZeroInt(), elysStaked.Amount)
			},
		},
		{
			"get all elys staked",
			func() (addr sdk.AccAddress) {
				suite.ResetSuite()

				addr = suite.AddAccounts(1, nil)[0]

				// set elys staked
				suite.app.EstakingKeeper.SetElysStaked(suite.ctx, types.ElysStaked{
					Address: addr.String(),
					Amount:  math.NewInt(1000),
				})

				return addr
			},
			func(addr sdk.AccAddress) {
				list := suite.app.EstakingKeeper.GetAllElysStaked(suite.ctx)
				suite.Require().Equal(2, len(list))
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
