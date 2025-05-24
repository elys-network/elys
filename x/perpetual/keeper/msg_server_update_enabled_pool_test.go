package keeper_test

import (
	"slices"

	sdkmath "cosmossdk.io/math"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/v5/x/perpetual/keeper"
	"github.com/elys-network/elys/v5/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) TestMsgUpdateEnabledPools() {
	addr := suite.AddAccounts(1, nil)
	amount := sdkmath.NewInt(1000)
	_, _ = suite.ResetAndSetSuite(addr, true, amount.MulRaw(1000), sdkmath.NewInt(2))
	params := suite.app.PerpetualKeeper.GetParams(suite.ctx)
	params.LeverageMax = sdkmath.LegacyMustNewDecFromStr("10")
	err := suite.app.PerpetualKeeper.SetParams(suite.ctx, &params)
	suite.Require().NoError(err)

	testCases := []struct {
		name             string
		input            *types.MsgUpdateEnabledPools
		expectErr        bool
		expectErrMsg     string
		setup            func()
		postValidateFunc func(msg *types.MsgUpdateEnabledPools)
	}{
		{"invalid authority",
			&types.MsgUpdateEnabledPools{
				Authority: addr[0].String(),
			},
			true,
			"invalid authority",
			func() {},
			func(msg *types.MsgUpdateEnabledPools) {},
		},
		{"Happy flow",
			&types.MsgUpdateEnabledPools{
				Authority:   authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				AddPools:    []uint64{4, 6, 7},
				RemovePools: []uint64{11, 141, 960},
			},
			false,
			"",
			func() {
				params.EnabledPools = []uint64{1, 2, 3, 5, 9, 11, 141, 960}
				err = suite.app.PerpetualKeeper.SetParams(suite.ctx, &params)
				suite.Require().NoError(err)
			},
			func(msg *types.MsgUpdateEnabledPools) {
				parameters := suite.app.PerpetualKeeper.GetParams(suite.ctx)
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
			tc.setup()
			msgServer := keeper.NewMsgServerImpl(*suite.app.PerpetualKeeper)
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
