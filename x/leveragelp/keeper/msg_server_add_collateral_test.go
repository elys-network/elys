package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	simapp "github.com/elys-network/elys/app"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/leveragelp/keeper"
	"github.com/elys-network/elys/x/leveragelp/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	stablekeeper "github.com/elys-network/elys/x/stablestake/keeper"
	stabletypes "github.com/elys-network/elys/x/stablestake/types"
	"math"
)

func initializeForAddCollateral(suite *KeeperTestSuite, addresses []sdk.AccAddress, asset1, asset2 string, createAmmPool bool) {
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
	if createAmmPool {
		poolId, err := suite.app.AmmKeeper.CreatePool(suite.ctx, &msgCreatePool)
		if err != nil {
			panic(err)
		}

		suite.app.LeveragelpKeeper.SetPool(suite.ctx, types.NewPool(poolId))

	}
	msgBond := stabletypes.MsgBond{
		Creator: addresses[1].String(),
		Amount:  issueAmount.QuoRaw(20),
	}
	stableStakeMsgServer := stablekeeper.NewMsgServerImpl(suite.app.StablestakeKeeper)
	_, err := stableStakeMsgServer.Bond(suite.ctx, &msgBond)
	if err != nil {
		panic(err)
	}
	msgBond.Creator = addresses[2].String()
	_, err = stableStakeMsgServer.Bond(suite.ctx, &msgBond)
	if err != nil {
		panic(err)
	}
}

func openPosition(suite *KeeperTestSuite, address sdk.AccAddress, collateralAmount sdk.Int, leverage sdk.Dec) {
	msg := types.MsgOpen{
		Creator:          address.String(),
		CollateralAsset:  ptypes.BaseCurrency,
		CollateralAmount: collateralAmount,
		AmmPoolId:        1,
		Leverage:         leverage,
		StopLossPrice:    sdk.MustNewDecFromStr("50.0"),
	}
	_, err := suite.app.LeveragelpKeeper.Open(suite.ctx, &msg)
	if err != nil {
		panic(err)
	}
	return
}

func (suite *KeeperTestSuite) TestMsgServerAddCollateral() {
	addresses := simapp.AddTestAddrs(suite.app, suite.ctx, 10, sdk.NewInt(1000000))
	asset1 := ptypes.ATOM
	asset2 := ptypes.BaseCurrency
	leverage := sdk.MustNewDecFromStr("2.0")
	collateralAmount := sdk.NewInt(10000000)
	//leverageLPShares := sdk.MustNewDecFromStr("20000095238095238100").TruncateInt()
	testCases := []struct {
		name                 string
		input                *types.MsgAddCollateral
		expectErr            bool
		expectErrMsg         string
		prerequisiteFunction func()
	}{
		{"position not found",
			&types.MsgAddCollateral{
				Creator:    addresses[0].String(),
				Id:         1,
				Collateral: sdk.NewInt(0),
			},
			true,
			types.ErrPositionDoesNotExist.Error(),
			func() {
				suite.ResetSuite()
				SetupCoinPrices(suite.ctx, suite.app.OracleKeeper)
				initializeForAddCollateral(suite, addresses, asset1, asset2, false)
			},
		},
		{"pool not found",
			&types.MsgAddCollateral{
				Creator:    addresses[0].String(),
				Id:         1,
				Collateral: sdk.NewInt(0),
			},
			true,
			types.ErrPoolDoesNotExist.Error(),
			func() {
				suite.ResetSuite()
				SetupCoinPrices(suite.ctx, suite.app.OracleKeeper)
				initializeForAddCollateral(suite, addresses, asset1, asset2, true)
				openPosition(suite, addresses[0], collateralAmount, leverage)
				suite.app.LeveragelpKeeper.DeletePool(suite.ctx, 1)
			},
		},
		{"pool not enabled",
			&types.MsgAddCollateral{
				Creator:    addresses[0].String(),
				Id:         1,
				Collateral: sdk.NewInt(0),
			},
			true,
			types.ErrPositionDisabled.Error(),
			func() {
				suite.ResetSuite()
				SetupCoinPrices(suite.ctx, suite.app.OracleKeeper)
				initializeForAddCollateral(suite, addresses, asset1, asset2, true)
				openPosition(suite, addresses[0], collateralAmount, leverage)
				pool, _ := suite.app.LeveragelpKeeper.GetPool(suite.ctx, 1)
				pool.Enabled = false
				suite.app.LeveragelpKeeper.SetPool(suite.ctx, pool)
			},
		},
		{"repaying more than allowed",
			&types.MsgAddCollateral{
				Creator:    addresses[0].String(),
				Id:         1,
				Collateral: sdk.NewInt(math.MaxInt64),
			},
			true,
			types.ErrInvalidCollateral.Error(),
			func() {
				suite.ResetSuite()
				SetupCoinPrices(suite.ctx, suite.app.OracleKeeper)
				initializeForAddCollateral(suite, addresses, asset1, asset2, true)
				openPosition(suite, addresses[0], collateralAmount, leverage)
			},
		},
		{"balance too low to repay",
			&types.MsgAddCollateral{
				Creator:    addresses[0].String(),
				Id:         1,
				Collateral: sdk.NewInt(10000),
			},
			true,
			"spendable balance  is smaller",
			func() {
				suite.ResetSuite()
				SetupCoinPrices(suite.ctx, suite.app.OracleKeeper)
				initializeForAddCollateral(suite, addresses, asset1, asset2, true)
				openPosition(suite, addresses[0], collateralAmount, leverage)
				amount := suite.app.BankKeeper.GetBalance(suite.ctx, addresses[0], "uusdc")
				err := suite.app.BankKeeper.SendCoins(suite.ctx, addresses[0], addresses[1], sdk.NewCoins(amount))
				if err != nil {
					panic(err)
				}

			},
		},
		{"positive case",
			&types.MsgAddCollateral{
				Creator:    addresses[0].String(),
				Id:         1,
				Collateral: sdk.NewInt(10000),
			},
			false,
			"",
			func() {
				suite.ResetSuite()
				SetupCoinPrices(suite.ctx, suite.app.OracleKeeper)
				initializeForAddCollateral(suite, addresses, asset1, asset2, true)
				openPosition(suite, addresses[0], collateralAmount, leverage)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			msgServer := keeper.NewMsgServerImpl(suite.app.LeveragelpKeeper)
			_, err := msgServer.AddCollateral(suite.ctx, tc.input)
			if tc.expectErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
