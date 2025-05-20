package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	simapp "github.com/elys-network/elys/v4/app"
	ammtypes "github.com/elys-network/elys/v4/x/amm/types"
	"github.com/elys-network/elys/v4/x/leveragelp/keeper"
	"github.com/elys-network/elys/v4/x/leveragelp/types"
	ptypes "github.com/elys-network/elys/v4/x/parameter/types"
)

func initializeForAddPool(suite *KeeperTestSuite, addresses []sdk.AccAddress, asset1, asset2 string) {
	fee := sdkmath.LegacyMustNewDecFromStr("0.0002")
	issueAmount := sdkmath.NewInt(10_000_000_000_000)
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
	msgCreatePool := ammtypes.MsgCreatePool{
		Sender: addresses[0].String(),
		PoolParams: ammtypes.PoolParams{
			SwapFee:   fee,
			UseOracle: true,
			FeeDenom:  ptypes.Elys,
		},
		PoolAssets: []ammtypes.PoolAsset{
			{
				Token:  sdk.NewInt64Coin(asset1, 100_000_000),
				Weight: sdkmath.NewInt(50),
			},
			{
				Token:  sdk.NewInt64Coin(asset2, 1000_000_000),
				Weight: sdkmath.NewInt(50),
			},
		},
	}
	_, err := suite.app.AmmKeeper.CreatePool(suite.ctx, &msgCreatePool)
	if err != nil {
		panic(err)
	}
}

func (suite *KeeperTestSuite) TestAdd_Pool() {
	suite.ResetSuite()
	suite.SetupCoinPrices(suite.ctx)
	addresses := simapp.AddTestAddrs(suite.app, suite.ctx, 10, sdkmath.NewInt(1000000))
	asset1 := ptypes.ATOM
	asset2 := ptypes.BaseCurrency
	initializeForAddPool(suite, addresses, asset1, asset2)
	testCases := []struct {
		name                 string
		input                *types.MsgAddPool
		expectErr            bool
		expectErrMsg         string
		prerequisiteFunction func()
	}{
		{name: "not allowed",
			input: &types.MsgAddPool{
				Authority: addresses[0].String(),
				Pool: types.AddPool{
					AmmPoolId:   1,
					LeverageMax: sdkmath.LegacyMustNewDecFromStr("10"),
				},
			},
			expectErr:    true,
			expectErrMsg: "invalid authority",
			prerequisiteFunction: func() {
			},
		},
		{name: "success",
			input: &types.MsgAddPool{
				Authority: "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
				Pool: types.AddPool{
					AmmPoolId:   1,
					LeverageMax: sdkmath.LegacyMustNewDecFromStr("10"),
				},
			},
			expectErr:            false,
			expectErrMsg:         "",
			prerequisiteFunction: func() {},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			msgServer := keeper.NewMsgServerImpl(*suite.app.LeveragelpKeeper)
			_, err := msgServer.AddPool(suite.ctx, tc.input)
			if tc.expectErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
