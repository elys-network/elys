package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	simapp "github.com/elys-network/elys/v6/app"
	"github.com/elys-network/elys/v6/x/leveragelp/keeper"
	"github.com/elys-network/elys/v6/x/leveragelp/types"
)

func (suite *KeeperTestSuite) TestMsgUpdatePool() {
	suite.ResetSuite()
	addresses := simapp.AddTestAddrs(suite.app, suite.ctx, 10, sdkmath.NewInt(1000000))
	newPool := types.NewPool(1, sdkmath.LegacyMustNewDecFromStr("5.5"), sdkmath.LegacyMustNewDecFromStr("0.6"), sdkmath.LegacyMustNewDecFromStr("0.8"))
	suite.app.LeveragelpKeeper.SetPool(suite.ctx, newPool)

	testCases := []struct {
		name             string
		input            *types.MsgUpdatePool
		expectErr        bool
		expectErrMsg     string
		postValidateFunc func(msg *types.MsgUpdatePool)
	}{
		{"invalid authority",
			&types.MsgUpdatePool{
				Authority:   addresses[0].String(),
				PoolId:      1,
				LeverageMax: sdkmath.LegacyMustNewDecFromStr("4.5"),
			},
			true,
			"invalid authority",
			func(msg *types.MsgUpdatePool) {},
		},
		{"pool not found",
			&types.MsgUpdatePool{
				Authority:   authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				PoolId:      2,
				LeverageMax: sdkmath.LegacyMustNewDecFromStr("4.5"),
			},
			true,
			"pool does not exists for pool id 2",
			func(msg *types.MsgUpdatePool) {},
		},
		{"Update max leverage for pool more than max leverage allowed",
			&types.MsgUpdatePool{
				Authority:   authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				PoolId:      1,
				LeverageMax: sdkmath.LegacyMustNewDecFromStr("11.5"),
			},
			true,
			"max leverage allowed is less than the leverage max",
			func(msg *types.MsgUpdatePool) {},
		},
		{"Happy flow",
			&types.MsgUpdatePool{
				Authority:   authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				PoolId:      1,
				LeverageMax: sdkmath.LegacyMustNewDecFromStr("9"),
			},
			false,
			"",
			func(msg *types.MsgUpdatePool) {
				pool, _ := suite.app.LeveragelpKeeper.GetPool(suite.ctx, 1)
				suite.Require().Equal(sdkmath.LegacyMustNewDecFromStr("9"), pool.LeverageMax)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			msgServer := keeper.NewMsgServerImpl(*suite.app.LeveragelpKeeper)
			_, err := msgServer.UpdatePool(suite.ctx, tc.input)
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
