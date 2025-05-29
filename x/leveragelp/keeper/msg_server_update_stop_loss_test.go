package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	simapp "github.com/elys-network/elys/v6/app"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
	"github.com/elys-network/elys/v6/x/leveragelp/keeper"
	"github.com/elys-network/elys/v6/x/leveragelp/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
	stablekeeper "github.com/elys-network/elys/v6/x/stablestake/keeper"
	stabletypes "github.com/elys-network/elys/v6/x/stablestake/types"
)

func initializeForUpdateStopLoss(suite *KeeperTestSuite, addresses []sdk.AccAddress, asset1, asset2 string, openPosition bool) {
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
				Token:  sdk.NewInt64Coin(asset2, 100_000_000),
				Weight: sdkmath.NewInt(50),
			},
		},
	}
	poolId, err := suite.app.AmmKeeper.CreatePool(suite.ctx, &msgCreatePool)
	suite.Require().NoError(err)
	enablePoolMsg := types.MsgAddPool{
		Authority: authtypes.NewModuleAddress("gov").String(),
		Pool: types.AddPool{
			AmmPoolId:   poolId,
			LeverageMax: sdkmath.LegacyNewDec(10),
		},
	}
	msgServer := keeper.NewMsgServerImpl(*suite.app.LeveragelpKeeper)
	_, err = msgServer.AddPool(suite.ctx, &enablePoolMsg)
	suite.Require().NoError(err)
	msgBond := stabletypes.MsgBond{
		Creator: addresses[1].String(),
		Amount:  issueAmount.QuoRaw(20),
		PoolId:  1,
	}
	stableStakeMsgServer := stablekeeper.NewMsgServerImpl(*suite.app.StablestakeKeeper)
	_, err = stableStakeMsgServer.Bond(suite.ctx, &msgBond)
	suite.Require().NoError(err)
	msgBond.Creator = addresses[2].String()
	_, err = stableStakeMsgServer.Bond(suite.ctx, &msgBond)
	suite.Require().NoError(err)

	if openPosition {
		openMsg := &types.MsgOpen{
			Creator:          addresses[0].String(),
			CollateralAsset:  "uusdc",
			CollateralAmount: sdkmath.OneInt().MulRaw(1000_0000),
			AmmPoolId:        1,
			Leverage:         sdkmath.LegacyOneDec().MulInt64(2),
			StopLossPrice:    sdkmath.LegacyZeroDec(),
		}
		_, err = suite.app.LeveragelpKeeper.Open(suite.ctx, openMsg)
		if err != nil {
			panic(err)
		}
	}
}
func (suite *KeeperTestSuite) TestUpdateStopLoss() {
	addresses := simapp.AddTestAddrs(suite.app, suite.ctx, 10, sdkmath.NewInt(1000000))
	asset1 := ptypes.ATOM
	asset2 := ptypes.BaseCurrency
	testCases := []struct {
		name                 string
		input                *types.MsgUpdateStopLoss
		expectErr            bool
		expectErrMsg         string
		prerequisiteFunction func()
		postValidateFunc     func()
	}{
		{name: "position not found",
			input: &types.MsgUpdateStopLoss{
				Creator:  addresses[0].String(),
				Position: 2,
				Price:    sdkmath.LegacyOneDec().MulInt64(10),
			},
			expectErr:    true,
			expectErrMsg: types.ErrPositionDoesNotExist.Error(),
			prerequisiteFunction: func() {
				suite.ResetSuite()
				suite.SetupCoinPrices(suite.ctx)
				initializeForUpdateStopLoss(suite, addresses, asset1, asset2, false)
			},
			postValidateFunc: func() {
			},
		},
		{name: "success",
			input: &types.MsgUpdateStopLoss{
				Creator:  addresses[0].String(),
				Position: 1,
				Price:    sdkmath.LegacyOneDec().MulInt64(10),
			},
			expectErr:    false,
			expectErrMsg: "",
			prerequisiteFunction: func() {
				suite.ResetSuite()
				suite.SetupCoinPrices(suite.ctx)
				initializeForUpdateStopLoss(suite, addresses, asset1, asset2, true)
			},
			postValidateFunc: func() {
				position, found := suite.app.LeveragelpKeeper.GetPositionWithId(suite.ctx, addresses[0], 1)
				suite.Require().True(found)
				suite.Require().Equal(position.StopLossPrice, sdkmath.LegacyOneDec().MulInt64(10))
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			msgServer := keeper.NewMsgServerImpl(*suite.app.LeveragelpKeeper)
			_, err := msgServer.UpdateStopLoss(suite.ctx, tc.input)
			if tc.expectErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().NoError(err)
			}
			tc.postValidateFunc()
		})
	}
}
