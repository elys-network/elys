package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	simapp "github.com/elys-network/elys/app"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	stablekeeper "github.com/elys-network/elys/x/stablestake/keeper"
	stabletypes "github.com/elys-network/elys/x/stablestake/types"
	"time"
)

func initializeForClose(suite *KeeperTestSuite, addresses []sdk.AccAddress, asset1, asset2 string) {
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
	suite.app.LeveragelpKeeper.SetPool(suite.ctx, types.NewPool(poolId))
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
}

func (suite *KeeperTestSuite) TestClose() {
	suite.ResetSuite()
	suite.SetupCoinPrices(suite.ctx)
	addresses := simapp.AddTestAddrs(suite.app, suite.ctx, 10, sdk.NewInt(1000000))
	asset1 := ptypes.ATOM
	asset2 := ptypes.BaseCurrency
	initializeForClose(suite, addresses, asset1, asset2)
	leverage := sdk.MustNewDecFromStr("2.0")
	collateralAmount := sdk.NewInt(10000000)
	leverageLPShares := sdk.MustNewDecFromStr("20000095238095238100").TruncateInt()
	testCases := []struct {
		name                 string
		input                *types.MsgClose
		expectErr            bool
		expectErrMsg         string
		prerequisiteFunction func()
		postValidateFunc     func()
	}{
		{"No position to close",
			&types.MsgClose{
				Creator:  addresses[0].String(),
				Id:       1,
				LpAmount: sdk.NewInt(0),
			},
			true,
			types.ErrPositionDoesNotExist.Error(),
			func() {
			},
			func() {
			},
		},
		{"Unlock time not reached",
			&types.MsgClose{
				Creator:  addresses[0].String(),
				Id:       1,
				LpAmount: sdk.NewInt(0),
			},
			true,
			"your funds will be locked for 1 hour",
			func() {
				msg := types.MsgOpen{
					Creator:          addresses[0].String(),
					CollateralAsset:  ptypes.BaseCurrency,
					CollateralAmount: collateralAmount,
					AmmPoolId:        1,
					Leverage:         leverage,
					StopLossPrice:    sdk.MustNewDecFromStr("50.0"),
				}
				_, err := suite.app.LeveragelpKeeper.Open(suite.ctx, &msg)
				suite.Require().NoError(err)
			},
			func() {
			},
		},
		{"Repay amount is greater than exit amount",
			&types.MsgClose{
				Creator:  addresses[0].String(),
				Id:       1,
				LpAmount: leverageLPShares.QuoRaw(2),
			},
			false,
			"",
			func() {
				msg := types.MsgOpen{
					Creator:          addresses[0].String(),
					CollateralAsset:  ptypes.BaseCurrency,
					CollateralAmount: collateralAmount,
					AmmPoolId:        1,
					Leverage:         leverage,
					StopLossPrice:    sdk.MustNewDecFromStr("50.0"),
				}
				_, err := suite.app.LeveragelpKeeper.Open(suite.ctx, &msg)
				suite.Require().NoError(err)
				suite.AddBlockTime(1000000 * time.Hour)
			},
			func() {
			},
		},
		{"Invalid Leverage LP shares amount to close",
			&types.MsgClose{
				Creator:  addresses[0].String(),
				Id:       1,
				LpAmount: leverageLPShares.MulRaw(2),
			},
			true,
			types.ErrInvalidCloseSize.Error(),
			func() {
				suite.ResetSuite()
				suite.SetupCoinPrices(suite.ctx)
				initializeForClose(suite, addresses, asset1, asset2)
				msg := types.MsgOpen{
					Creator:          addresses[0].String(),
					CollateralAsset:  ptypes.BaseCurrency,
					CollateralAmount: collateralAmount,
					AmmPoolId:        1,
					Leverage:         leverage,
					StopLossPrice:    sdk.MustNewDecFromStr("50.0"),
				}
				_, err := suite.app.LeveragelpKeeper.Open(suite.ctx, &msg)
				suite.Require().NoError(err)
				suite.AddBlockTime(time.Hour)
			},
			func() {
			},
		},
		{"Position Health is lower than safety factor and closing partially, should close fully",
			&types.MsgClose{
				Creator:  addresses[0].String(),
				Id:       1,
				LpAmount: leverageLPShares.QuoRaw(2000000),
			},
			false,
			"",
			func() {
				suite.ResetSuite()
				suite.SetupCoinPrices(suite.ctx)
				initializeForClose(suite, addresses, asset1, asset2)
				msg := types.MsgOpen{
					Creator:          addresses[0].String(),
					CollateralAsset:  ptypes.BaseCurrency,
					CollateralAmount: collateralAmount,
					AmmPoolId:        1,
					Leverage:         leverage,
					StopLossPrice:    sdk.MustNewDecFromStr("50.0"),
				}
				_, err := suite.app.LeveragelpKeeper.Open(suite.ctx, &msg)
				suite.Require().NoError(err)
				suite.AddBlockTime(1000000 * time.Hour)
			},
			func() {
				_, err := suite.app.LeveragelpKeeper.GetPosition(suite.ctx, addresses[0], 1)
				suite.Require().Contains(err.Error(), "position not found")
			},
		},
		{"Position LP amount is 0",
			&types.MsgClose{
				Creator:  addresses[0].String(),
				Id:       1,
				LpAmount: leverageLPShares.QuoRaw(2000000),
			},
			false,
			"",
			func() {
				suite.ResetSuite()
				suite.SetupCoinPrices(suite.ctx)
				initializeForClose(suite, addresses, asset1, asset2)
				msg := types.MsgOpen{
					Creator:          addresses[0].String(),
					CollateralAsset:  ptypes.BaseCurrency,
					CollateralAmount: collateralAmount,
					AmmPoolId:        1,
					Leverage:         leverage,
					StopLossPrice:    sdk.MustNewDecFromStr("50.0"),
				}
				_, err := suite.app.LeveragelpKeeper.Open(suite.ctx, &msg)
				suite.Require().NoError(err)
				suite.AddBlockTime(1000000 * time.Hour)
			},
			func() {
				_, err := suite.app.LeveragelpKeeper.GetPosition(suite.ctx, addresses[0], 1)
				suite.Require().Contains(err.Error(), "position not found")
			},
		},
		{"Closing partial position",
			&types.MsgClose{
				Creator:  addresses[0].String(),
				Id:       1,
				LpAmount: leverageLPShares.QuoRaw(2),
			},
			false,
			"",
			func() {
				suite.ResetSuite()
				suite.SetupCoinPrices(suite.ctx)
				initializeForClose(suite, addresses, asset1, asset2)
				msg := types.MsgOpen{
					Creator:          addresses[0].String(),
					CollateralAsset:  ptypes.BaseCurrency,
					CollateralAmount: collateralAmount,
					AmmPoolId:        1,
					Leverage:         leverage,
					StopLossPrice:    sdk.MustNewDecFromStr("50.0"),
				}
				_, err := suite.app.LeveragelpKeeper.Open(suite.ctx, &msg)
				suite.Require().NoError(err)
				suite.AddBlockTime(time.Hour)
			},
			func() {
				position, _ := suite.app.LeveragelpKeeper.GetPosition(suite.ctx, addresses[0], 1)
				suite.Require().Equal(position.LeveragedLpAmount, leverageLPShares.QuoRaw(2))
			},
		},
		{"Closing whole position",
			&types.MsgClose{
				Creator:  addresses[0].String(),
				Id:       1,
				LpAmount: sdk.NewInt(0),
			},
			false,
			"",
			func() {
			},
			func() {
				_, err := suite.app.LeveragelpKeeper.GetPosition(suite.ctx, addresses[0], 1)
				suite.Require().Contains(err.Error(), "position not found")
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			portfolio_old, found := suite.app.TierKeeper.GetPortfolio(suite.ctx, tc.input.Creator, suite.app.TierKeeper.GetDateFromBlock(suite.ctx.BlockTime()))
			tc.prerequisiteFunction()
			_, err := suite.app.LeveragelpKeeper.Close(suite.ctx, tc.input)
			if tc.expectErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				// The new value of the portfolio after the hook is called.
				portfolio_new, _ := suite.app.TierKeeper.GetPortfolio(suite.ctx, tc.input.Creator, suite.app.TierKeeper.GetDateFromBlock(suite.ctx.BlockTime()))
				// Initially, there were no entries for the portfolio
				if !found {
					// The portfolio value changes after the hook is called.
					suite.Require().NotEqual(portfolio_old, portfolio_new)
				}
				suite.Require().NoError(err)
			}
			tc.postValidateFunc()
		})
	}
}
