package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	simapp "github.com/elys-network/elys/v4/app"
	"github.com/elys-network/elys/v4/x/leveragelp/keeper"
	"github.com/elys-network/elys/v4/x/leveragelp/types"
)

func (suite *KeeperTestSuite) TestMsgServerUpdateParams() {
	suite.ResetSuite()
	addresses := simapp.AddTestAddrs(suite.app, suite.ctx, 10, sdkmath.NewInt(1000000))
	params := types.DefaultParams()
	testCases := []struct {
		name                 string
		input                *types.MsgUpdateParams
		expectErr            bool
		expectErrMsg         string
		prerequisiteFunction func()
		postValidateFunc     func(msg *types.MsgUpdateParams)
	}{
		{"invalid authority",
			&types.MsgUpdateParams{
				Authority: addresses[0].String(),
				Params:    &params,
			},
			true,
			"invalid authority",
			func() {
			},
			func(msg *types.MsgUpdateParams) {

			},
		},
		{"invalid params",
			&types.MsgUpdateParams{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Params:    &params,
			},
			true,
			"leverage max must be greater than 1",
			func() {
				p := &params
				p.LeverageMax = sdkmath.LegacyOneDec().MulInt64(-1)
			},
			func(msg *types.MsgUpdateParams) {

			},
		},
		{"positive case",
			&types.MsgUpdateParams{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Params:    &params,
			},
			false,
			"",
			func() {
				p := &params
				p.LeverageMax = sdkmath.LegacyMustNewDecFromStr("2.5")
			},
			func(msg *types.MsgUpdateParams) {
				parameters := suite.app.LeveragelpKeeper.GetParams(suite.ctx)
				suite.Require().Equal(sdkmath.LegacyMustNewDecFromStr("2.5"), parameters.LeverageMax)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			msgServer := keeper.NewMsgServerImpl(*suite.app.LeveragelpKeeper)
			_, err := msgServer.UpdateParams(suite.ctx, tc.input)
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
