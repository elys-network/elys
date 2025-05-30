package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	simapp "github.com/elys-network/elys/v6/app"
	"github.com/elys-network/elys/v6/x/leveragelp/keeper"
	"github.com/elys-network/elys/v6/x/leveragelp/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
)

func (suite *KeeperTestSuite) TestRemove_Pool() {
	suite.ResetSuite()
	suite.SetupCoinPrices(suite.ctx)
	addresses := simapp.AddTestAddrs(suite.app, suite.ctx, 10, sdkmath.NewInt(1000000))
	asset1 := ptypes.ATOM
	asset2 := ptypes.BaseCurrency
	initializeForOpen(suite, addresses, asset1, asset2)
	testCases := []struct {
		name                 string
		input                *types.MsgRemovePool
		expectErr            bool
		expectErrMsg         string
		prerequisiteFunction func()
	}{
		{name: "not allowed invalid authority",
			input: &types.MsgRemovePool{
				Authority: addresses[0].String(),
				Id:        1,
			},
			expectErr:    true,
			expectErrMsg: "invalid authority",
			prerequisiteFunction: func() {
			},
		},
		{name: "non zero pool leveraged amount",
			input: &types.MsgRemovePool{
				Authority: "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
				Id:        1,
			},
			expectErr:    true,
			expectErrMsg: "pool leverage amount is greater than zero",
			prerequisiteFunction: func() {
				pool := types.NewPool(1, sdkmath.LegacyMustNewDecFromStr("10"))
				pool.LeveragedLpAmount = sdkmath.OneInt()
				suite.app.LeveragelpKeeper.SetPool(suite.ctx, pool)
			},
		},
		{name: "success",
			input: &types.MsgRemovePool{
				Authority: "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
				Id:        1,
			},
			expectErr:    false,
			expectErrMsg: "",
			prerequisiteFunction: func() {
				pool := types.NewPool(1, sdkmath.LegacyMustNewDecFromStr("10"))
				pool.LeveragedLpAmount = sdkmath.ZeroInt()
				suite.app.LeveragelpKeeper.SetPool(suite.ctx, pool)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			msgServer := keeper.NewMsgServerImpl(*suite.app.LeveragelpKeeper)
			_, err := msgServer.RemovePool(suite.ctx, tc.input)
			if tc.expectErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
