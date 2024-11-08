package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/estaking/types"
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
					Amount:  sdk.NewInt(1000),
				})

				// get elys staked
				elysStaked := suite.app.EstakingKeeper.GetElysStaked(suite.ctx, addr)
				suite.Require().Equal(sdk.NewInt(1000), elysStaked.Amount)

				return addr
			},
			func(addr sdk.AccAddress) {
				suite.app.EstakingKeeper.RemoveElysStaked(suite.ctx, addr)

				// get elys staked
				elysStaked := suite.app.EstakingKeeper.GetElysStaked(suite.ctx, addr)
				suite.Require().Equal(sdk.ZeroInt(), elysStaked.Amount)
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
					Amount:  sdk.NewInt(1000),
				})

				return addr
			},
			func(addr sdk.AccAddress) {
				list := suite.app.EstakingKeeper.GetAllElysStaked(suite.ctx)
				suite.Require().Equal(1, len(list))
				suite.Require().Equal(sdk.NewInt(1000), list[0].Amount)
			},
		},
		{
			"get all legacy elys staked",
			func() (addr sdk.AccAddress) {
				suite.ResetSuite()

				addr = suite.AddAccounts(1, nil)[0]

				// set elys staked
				suite.app.EstakingKeeper.SetLegacyElysStaked(suite.ctx, types.ElysStaked{
					Address: addr.String(),
					Amount:  sdk.NewInt(1000),
				})

				return addr
			},
			func(addr sdk.AccAddress) {
				list := suite.app.EstakingKeeper.GetAllLegacyElysStaked(suite.ctx)
				suite.Require().Equal(1, len(list))
				suite.Require().Equal(sdk.NewInt(1000), list[0].Amount)
			},
		},
		{
			"delete legacy elys staked",
			func() (addr sdk.AccAddress) {
				suite.ResetSuite()

				addr = suite.AddAccounts(1, nil)[0]

				return addr
			},
			func(addr sdk.AccAddress) {
				suite.app.EstakingKeeper.DeleteLegacyElysStaked(suite.ctx, addr.String())
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
