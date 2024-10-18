package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/leveragelp/keeper"
	"github.com/elys-network/elys/x/leveragelp/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func initializeForUpdatePool(suite *KeeperTestSuite, addresses []sdk.AccAddress, asset1, asset2 string) {
	issueAmount := sdk.NewInt(10_000_000_000_000)
	for _, address := range addresses {
		coins := sdk.NewCoins(
			sdk.NewCoin(ptypes.ATOM, issueAmount),
			sdk.NewCoin(ptypes.Elys, issueAmount),
			sdk.NewCoin(ptypes.BaseCurrency, issueAmount),
		)
		err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
		if err != nil {
			panic(err)
		}
		err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, address, coins)
		if err != nil {
			panic(err)
		}
	}
	suite.app.LeveragelpKeeper.SetPool(suite.ctx, types.NewPool(1))
}

func (suite *KeeperTestSuite) TestUpdate_Pool() {
	suite.ResetSuite()
	suite.SetupCoinPrices(suite.ctx)
	addresses := simapp.AddTestAddrs(suite.app, suite.ctx, 10, sdk.NewInt(1000000))
	asset1 := ptypes.ATOM
	asset2 := ptypes.BaseCurrency
	initializeForUpdatePool(suite, addresses, asset1, asset2)
	testCases := []struct {
		name                 string
		input                *types.MsgUpdatePool
		expectErr            bool
		expectErrMsg         string
		prerequisiteFunction func()
	}{
		{name: "not allowed invalid authority",
			input: &types.MsgUpdatePool{
				Authority: addresses[0].String(),
				UpdatePool: &types.UpdatePool{
					PoolId:  1,
					Enabled: false,
					Closed:  false,
				},
			},
			expectErr:    true,
			expectErrMsg: "invalid authority",
			prerequisiteFunction: func() {
			},
		},
		{name: "success",
			input: &types.MsgUpdatePool{
				Authority: "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
				UpdatePool: &types.UpdatePool{
					PoolId:  1,
					Enabled: false,
					Closed:  false,
				},
			},
			expectErr:    false,
			expectErrMsg: "",
			prerequisiteFunction: func() {
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			msgServer := keeper.NewMsgServerImpl(*suite.app.LeveragelpKeeper)
			_, err := msgServer.UpdatePool(suite.ctx, tc.input)
			if tc.expectErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
