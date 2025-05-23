package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/leveragelp/keeper"
	"github.com/elys-network/elys/x/leveragelp/types"
	"slices"
)

func (suite *KeeperTestSuite) TestMsgServerUpdateEnabledPools() {
	suite.ResetSuite()
	addresses := simapp.AddTestAddrs(suite.app, suite.ctx, 10, sdkmath.NewInt(1000000))
	params := types.DefaultParams()
	testCases := []struct {
		name                 string
		input                *types.MsgUpdateEnabledPools
		expectErr            bool
		expectErrMsg         string
		prerequisiteFunction func()
		postValidateFunc     func(msg *types.MsgUpdateEnabledPools)
	}{
		{"invalid authority",
			&types.MsgUpdateEnabledPools{
				Authority: addresses[0].String(),
			},
			true,
			"invalid authority",
			func() {
			},
			func(msg *types.MsgUpdateEnabledPools) {

			},
		},
		{"positive case",
			&types.MsgUpdateEnabledPools{
				Authority:   authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				AddPools:    []uint64{4, 6, 7},
				RemovePools: []uint64{11, 141, 960},
			},
			false,
			"",
			func() {
				p := &params
				p.EnabledPools = []uint64{1, 2, 3, 5, 9, 11, 141, 960}
				err := suite.app.LeveragelpKeeper.SetParams(suite.ctx, p)
				suite.Require().NoError(err)
			},
			func(msg *types.MsgUpdateEnabledPools) {
				parameters := suite.app.LeveragelpKeeper.GetParams(suite.ctx)
				suite.Require().Equal(8, len(parameters.EnabledPools))
				result := []uint64{1, 2, 3, 4, 5, 6, 7, 9}
				for _, p := range result {
					suite.Require().Contains(parameters.EnabledPools, p)
					suite.Require().True(slices.Contains(parameters.EnabledPools, p))
				}
				for _, p := range []uint64{11, 141, 960} {
					suite.Require().NotContains(parameters.EnabledPools, p)
					suite.Require().False(slices.Contains(parameters.EnabledPools, p))
				}
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			msgServer := keeper.NewMsgServerImpl(*suite.app.LeveragelpKeeper)
			_, err := msgServer.UpdateEnabledPools(suite.ctx, tc.input)
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
