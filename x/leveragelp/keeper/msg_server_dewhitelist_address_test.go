package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	simapp "github.com/elys-network/elys/v5/app"
	"github.com/elys-network/elys/v5/x/leveragelp/keeper"
	"github.com/elys-network/elys/v5/x/leveragelp/types"
)

func (suite *KeeperTestSuite) TestMsgServerDewhitelistAddress() {
	suite.ResetSuite()
	addresses := simapp.AddTestAddrs(suite.app, suite.ctx, 10, sdkmath.NewInt(1000000))
	msgServer := keeper.NewMsgServerImpl(*suite.app.LeveragelpKeeper)
	msg := &types.MsgWhitelist{
		Authority:          authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		WhitelistedAddress: addresses[0].String(),
	}
	_, err := msgServer.Whitelist(suite.ctx, msg)
	suite.Require().NoError(err)
	testCases := []struct {
		name                 string
		input                *types.MsgDewhitelist
		expectErr            bool
		expectErrMsg         string
		prerequisiteFunction func()
		postValidateFunc     func(msg *types.MsgDewhitelist)
	}{
		{"invalid authority",
			&types.MsgDewhitelist{
				Authority:          addresses[0].String(),
				WhitelistedAddress: addresses[0].String(),
			},
			true,
			"invalid authority",
			func() {
			},
			func(msg *types.MsgDewhitelist) {

			},
		},
		{"positive case",
			&types.MsgDewhitelist{
				Authority:          authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				WhitelistedAddress: addresses[0].String(),
			},
			false,
			"",
			func() {
			},
			func(msg *types.MsgDewhitelist) {
				whitelisted := suite.app.LeveragelpKeeper.GetAllWhitelistedAddress(suite.ctx)
				suite.Require().NotContains(whitelisted, addresses[0])
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			_, err = msgServer.Dewhitelist(suite.ctx, tc.input)
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
