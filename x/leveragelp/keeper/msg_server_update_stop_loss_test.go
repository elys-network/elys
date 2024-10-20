package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	simapp "github.com/elys-network/elys/app"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/leveragelp/keeper"
	"github.com/elys-network/elys/x/leveragelp/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	stablekeeper "github.com/elys-network/elys/x/stablestake/keeper"
	stabletypes "github.com/elys-network/elys/x/stablestake/types"
)

func initializeForUpdateStopLoss(suite *KeeperTestSuite, addresses []sdk.AccAddress, asset1, asset2 string, openPosition bool) {
	fee := sdk.MustNewDecFromStr("0.0002")
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
	msgCreatePool := ammtypes.MsgCreatePool{
		Sender: addresses[0].String(),
		PoolParams: &ammtypes.PoolParams{
			SwapFee:                     fee,
			ExitFee:                     fee,
			UseOracle:                   true,
			WeightBreakingFeeMultiplier: fee,
			WeightBreakingFeeExponent:   fee,
			ExternalLiquidityRatio:      fee,
			WeightRecoveryFeePortion:    fee,
			ThresholdWeightDifference:   fee,
			FeeDenom:                    ptypes.Elys,
		},
		PoolAssets: []ammtypes.PoolAsset{
			{
				Token:  sdk.NewInt64Coin(asset1, 100_000_000),
				Weight: sdk.NewInt(50),
			},
			{
				Token:  sdk.NewInt64Coin(asset2, 100_000_000),
				Weight: sdk.NewInt(50),
			},
		},
	}
	poolId, err := suite.app.AmmKeeper.CreatePool(suite.ctx, &msgCreatePool)
	if err != nil {
		panic(err)
	}
	suite.app.LeveragelpKeeper.SetPool(suite.ctx, types.NewPool(poolId, math.LegacyMustNewDecFromStr("10")))
	msgBond := stabletypes.MsgBond{
		Creator: addresses[1].String(),
		Amount:  issueAmount.QuoRaw(20),
	}
	stableStakeMsgServer := stablekeeper.NewMsgServerImpl(suite.app.StablestakeKeeper)
	_, err = stableStakeMsgServer.Bond(suite.ctx, &msgBond)
	if err != nil {
		panic(err)
	}
	msgBond.Creator = addresses[2].String()
	_, err = stableStakeMsgServer.Bond(suite.ctx, &msgBond)
	if err != nil {
		panic(err)
	}

	if openPosition {
		openMsg := &types.MsgOpen{
			Creator:          addresses[0].String(),
			CollateralAsset:  "uusdc",
			CollateralAmount: sdk.OneInt().MulRaw(1000_0000),
			AmmPoolId:        1,
			Leverage:         sdk.OneDec().MulInt64(2),
			StopLossPrice:    sdk.ZeroDec(),
		}
		_, err = suite.app.LeveragelpKeeper.Open(suite.ctx, openMsg)
		if err != nil {
			panic(err)
		}
	}
}
func (suite *KeeperTestSuite) TestUpdateStopLoss() {
	addresses := simapp.AddTestAddrs(suite.app, suite.ctx, 10, sdk.NewInt(1000000))
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
				Price:    sdk.OneDec().MulInt64(10),
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
				Price:    sdk.OneDec().MulInt64(10),
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
				suite.Require().Equal(position.StopLossPrice, sdk.OneDec().MulInt64(10))
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
