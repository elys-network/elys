package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	simapp "github.com/elys-network/elys/v4/app"
	"github.com/elys-network/elys/v4/x/leveragelp/keeper"
	"github.com/elys-network/elys/v4/x/leveragelp/types"
)

func (suite *KeeperTestSuite) TestMsgServerWhitelistAddress() {
	suite.ResetSuite()
	addresses := simapp.AddTestAddrs(suite.app, suite.ctx, 10, sdkmath.NewInt(1000000))
	testCases := []struct {
		name                 string
		input                *types.MsgWhitelist
		expectErr            bool
		expectErrMsg         string
		prerequisiteFunction func()
		postValidateFunc     func(msg *types.MsgWhitelist)
	}{
		{"invalid authority",
			&types.MsgWhitelist{
				Authority:          addresses[0].String(),
				WhitelistedAddress: addresses[0].String(),
			},
			true,
			"invalid authority",
			func() {
			},
			func(msg *types.MsgWhitelist) {

			},
		},
		{"positive case",
			&types.MsgWhitelist{
				Authority:          authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				WhitelistedAddress: addresses[0].String(),
			},
			false,
			"",
			func() {
			},
			func(msg *types.MsgWhitelist) {
				whitelisted := suite.app.LeveragelpKeeper.GetAllWhitelistedAddress(suite.ctx)
				suite.Require().Contains(whitelisted, addresses[0])
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			msgServer := keeper.NewMsgServerImpl(*suite.app.LeveragelpKeeper)
			_, err := msgServer.Whitelist(suite.ctx, tc.input)
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
