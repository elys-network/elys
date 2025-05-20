package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/v4/x/perpetual/keeper"
	"github.com/elys-network/elys/v4/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) TestMsgUpdateMaxLeverageForPool() {
	addr := suite.AddAccounts(1, nil)
	amount := sdkmath.NewInt(1000)
	_, _ = suite.ResetAndSetSuite(addr, true, amount.MulRaw(1000), sdkmath.NewInt(2))
	params := suite.app.PerpetualKeeper.GetParams(suite.ctx)
	params.LeverageMax = sdkmath.LegacyMustNewDecFromStr("10")
	suite.app.PerpetualKeeper.SetParams(suite.ctx, &params)

	testCases := []struct {
		name             string
		input            *types.MsgUpdateMaxLeverageForPool
		expectErr        bool
		expectErrMsg     string
		postValidateFunc func(msg *types.MsgUpdateMaxLeverageForPool)
	}{
		{"invalid authority",
			&types.MsgUpdateMaxLeverageForPool{
				Authority:   addr[0].String(),
				PoolId:      1,
				LeverageMax: sdkmath.LegacyMustNewDecFromStr("4.5"),
			},
			true,
			"invalid authority",
			func(msg *types.MsgUpdateMaxLeverageForPool) {},
		},
		{"pool not found",
			&types.MsgUpdateMaxLeverageForPool{
				Authority:   authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				PoolId:      2,
				LeverageMax: sdkmath.LegacyMustNewDecFromStr("4.5"),
			},
			true,
			"pool does not exists for pool id 2",
			func(msg *types.MsgUpdateMaxLeverageForPool) {},
		},
		{"Update max leverage for pool more than max leverage allowed",
			&types.MsgUpdateMaxLeverageForPool{
				Authority:   authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				PoolId:      1,
				LeverageMax: sdkmath.LegacyMustNewDecFromStr("11.5"),
			},
			true,
			"max leverage allowed is less than the leverage max",
			func(msg *types.MsgUpdateMaxLeverageForPool) {},
		},
		{"Happy flow",
			&types.MsgUpdateMaxLeverageForPool{
				Authority:   authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				PoolId:      1,
				LeverageMax: sdkmath.LegacyMustNewDecFromStr("9"),
			},
			false,
			"",
			func(msg *types.MsgUpdateMaxLeverageForPool) {
				pool, _ := suite.app.PerpetualKeeper.GetPool(suite.ctx, 1)
				suite.Require().Equal(sdkmath.LegacyMustNewDecFromStr("9"), pool.LeverageMax)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			msgServer := keeper.NewMsgServerImpl(*suite.app.PerpetualKeeper)
			_, err := msgServer.UpdateMaxLeverageForPool(suite.ctx, tc.input)
			if tc.expectErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().NoError(err)
			}
			tc.postValidateFunc(tc.input)
		})
	}
}
