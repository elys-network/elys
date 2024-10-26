package keeper_test

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	simapp "github.com/elys-network/elys/app"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	leveragelpmodulekeeper "github.com/elys-network/elys/x/leveragelp/keeper"
	"github.com/elys-network/elys/x/leveragelp/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	stablekeeper "github.com/elys-network/elys/x/stablestake/keeper"
	stabletypes "github.com/elys-network/elys/x/stablestake/types"
)

func initializeForOpen(suite *KeeperTestSuite, addresses []sdk.AccAddress, asset1, asset2 string) {
	fee := sdk.MustNewDecFromStr("0.0002")
	issueAmount := sdk.NewInt(10_000_000_000_000_000)
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
				Token:  sdk.NewInt64Coin(asset1, 100_000_000_000_000),
				Weight: sdk.NewInt(50),
			},
			{
				Token:  sdk.NewInt64Coin(asset2, 1000_000_000_000_000),
				Weight: sdk.NewInt(50),
			},
		},
	}
	poolId, err := suite.app.AmmKeeper.CreatePool(suite.ctx, &msgCreatePool)
	if err != nil {
		panic(err)
	}
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

	addPoolMsg := types.MsgAddPool{
		Authority: authtypes.NewModuleAddress("gov").String(),
		Pool: types.AddPool{
			AmmPoolId:   poolId,
			LeverageMax: math.LegacyMustNewDecFromStr("10"),
		},
	}
	_, err = leveragelpmodulekeeper.NewMsgServerImpl(*suite.app.LeveragelpKeeper).AddPool(suite.ctx, &addPoolMsg)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) TestOpen_PoolWithBaseCurrencyAsset() {
	suite.ResetSuite()
	suite.SetupCoinPrices(suite.ctx)
	addresses := simapp.AddTestAddrs(suite.app, suite.ctx, 10, sdk.NewInt(1000000))
	asset1 := ptypes.ATOM
	asset2 := ptypes.BaseCurrency
	initializeForOpen(suite, addresses, asset1, asset2)
	testCases := []struct {
		name                 string
		input                *types.MsgOpen
		expectErr            bool
		expectErrMsg         string
		prerequisiteFunction func()
	}{
		{"failed user authorization",
			&types.MsgOpen{
				Creator:          addresses[0].String(),
				CollateralAsset:  "stake",
				CollateralAmount: sdk.NewInt(1000),
				AmmPoolId:        1,
				Leverage:         sdk.MustNewDecFromStr("10.0"),
				StopLossPrice:    sdk.MustNewDecFromStr("100.0"),
			},
			true,
			errorsmod.Wrap(types.ErrUnauthorised, "unauthorised").Error(),
			func() {
				suite.EnableWhiteListing()
			},
		},
		{"no positions allowed",
			&types.MsgOpen{
				Creator:          addresses[0].String(),
				CollateralAsset:  "stake",
				CollateralAmount: sdk.NewInt(1000000000),
				AmmPoolId:        1,
				Leverage:         sdk.MustNewDecFromStr("10.0"),
				StopLossPrice:    sdk.MustNewDecFromStr("100.0"),
			},
			true,
			errorsmod.Wrap(types.ErrMaxOpenPositions, "cannot open new positions").Error(),
			func() {
				suite.DisableWhiteListing()
				suite.SetMaxOpenPositions(0)
			},
		},
		{name: "Max positions reached",
			input: &types.MsgOpen{
				Creator:          addresses[0].String(),
				CollateralAsset:  ptypes.BaseCurrency,
				CollateralAmount: sdk.NewInt(1000),
				AmmPoolId:        1,
				Leverage:         sdk.MustNewDecFromStr("2.0"),
				StopLossPrice:    sdk.MustNewDecFromStr("100.0"),
			},
			expectErr:    true,
			expectErrMsg: types.ErrMaxOpenPositions.Wrapf("cannot open new positions").Error(),
			prerequisiteFunction: func() {
			},
		},
		{name: "Pool not found",
			input: &types.MsgOpen{
				Creator:          addresses[0].String(),
				CollateralAsset:  ptypes.BaseCurrency,
				CollateralAmount: sdk.NewInt(1000),
				AmmPoolId:        100,
				Leverage:         sdk.MustNewDecFromStr("2.0"),
				StopLossPrice:    sdk.MustNewDecFromStr("100.0"),
			},
			expectErr:    true,
			expectErrMsg: types.ErrPoolDoesNotExist.Wrapf("poolId: %d", 100).Error(),
			prerequisiteFunction: func() {
				suite.SetMaxOpenPositions(3)
			},
		},
		{name: "base currency not found",
			input: &types.MsgOpen{
				Creator:          addresses[0].String(),
				CollateralAsset:  ptypes.BaseCurrency,
				CollateralAmount: sdk.NewInt(1000),
				AmmPoolId:        2,
				Leverage:         sdk.MustNewDecFromStr("2.0"),
				StopLossPrice:    sdk.MustNewDecFromStr("100.0"),
			},
			expectErr:    true,
			expectErrMsg: "pool does not exis",
			prerequisiteFunction: func() {
				pool := types.NewPool(2, math.LegacyMustNewDecFromStr("10"))
				suite.app.LeveragelpKeeper.SetPool(suite.ctx, pool)
				suite.RemovePrices(suite.ctx, []string{"uusdc"})
				suite.SetMaxOpenPositions(20)
			},
		},
		{name: "AMM Pool not found",
			input: &types.MsgOpen{
				Creator:          addresses[0].String(),
				CollateralAsset:  ptypes.BaseCurrency,
				CollateralAmount: sdk.NewInt(1000),
				AmmPoolId:        2,
				Leverage:         sdk.MustNewDecFromStr("2.0"),
				StopLossPrice:    sdk.MustNewDecFromStr("100.0"),
			},
			expectErr:    true,
			expectErrMsg: "pool does not exis",
			prerequisiteFunction: func() {
				suite.SetupCoinPrices(suite.ctx)
			},
		},
		{name: "Pool not enabled",
			input: &types.MsgOpen{
				Creator:          addresses[0].String(),
				CollateralAsset:  ptypes.BaseCurrency,
				CollateralAmount: sdk.NewInt(1000),
				AmmPoolId:        2,
				Leverage:         sdk.MustNewDecFromStr("2.0"),
				StopLossPrice:    sdk.MustNewDecFromStr("100.0"),
			},
			expectErr:    true,
			expectErrMsg: "denom does not exist in pool",
			prerequisiteFunction: func() {
				pool := types.NewPool(2, sdk.NewDec(60))
				suite.app.LeveragelpKeeper.SetPool(suite.ctx, pool)
				amm_pool := ammtypes.Pool{PoolId: 2, Address: ammtypes.NewPoolAddress(2).String(), TotalShares: sdk.Coin{Amount: sdk.NewInt(100)}}
				suite.app.AmmKeeper.SetPool(suite.ctx, amm_pool)
			},
		},
		{"Collateral asset not equal to base currency",
			&types.MsgOpen{
				Creator:          addresses[0].String(),
				CollateralAsset:  ptypes.Elys,
				CollateralAmount: sdk.NewInt(10000000),
				AmmPoolId:        1,
				Leverage:         sdk.MustNewDecFromStr("2.0"),
				StopLossPrice:    sdk.MustNewDecFromStr("50.0"),
			},
			true,
			types.ErrOnlyBaseCurrencyAllowed.Error(),
			func() {
				suite.SetPoolThreshold(sdk.MustNewDecFromStr("0.2"))
			},
		},
		{"Base currency not found",
			&types.MsgOpen{
				Creator:          addresses[0].String(),
				CollateralAsset:  ptypes.Elys,
				CollateralAmount: sdk.NewInt(100000000),
				AmmPoolId:        1,
				Leverage:         sdk.MustNewDecFromStr("2.0"),
				StopLossPrice:    sdk.MustNewDecFromStr("50.0"),
			},
			true,
			types.ErrOnlyBaseCurrencyAllowed.Error(),
			func() {
			},
		},
		{"First open position but leads to low pool health",
			&types.MsgOpen{
				Creator:          addresses[0].String(),
				CollateralAsset:  ptypes.BaseCurrency,
				CollateralAmount: sdk.NewInt(100000000),
				AmmPoolId:        1,
				Leverage:         sdk.MustNewDecFromStr("2.0"),
				StopLossPrice:    sdk.MustNewDecFromStr("50.0"),
			},
			true,
			types.ErrInvalidPosition.Wrapf("pool health too low to open new positions").Error(),
			func() {
				suite.AddCoinPrices(suite.ctx, []string{ptypes.BaseCurrency})
				suite.SetPoolThreshold(sdk.OneDec())
			},
		},
		{"Low Balance of creator",
			&types.MsgOpen{
				Creator:          simapp.AddTestAddrs(suite.app, suite.ctx, 1, sdk.NewInt(0))[0].String(),
				CollateralAsset:  ptypes.BaseCurrency,
				CollateralAmount: sdk.NewInt(10000000),
				AmmPoolId:        1,
				Leverage:         sdk.MustNewDecFromStr("2.0"),
				StopLossPrice:    sdk.MustNewDecFromStr("50.0"),
			},
			true,
			"insufficient funds",
			func() {
				suite.SetPoolThreshold(sdk.MustNewDecFromStr("0.2"))
			},
		},
		{"Borrowing more than allowed",
			&types.MsgOpen{
				Creator:          addresses[0].String(),
				CollateralAsset:  ptypes.BaseCurrency,
				CollateralAmount: sdk.NewInt(1_000_000_000_000_000_000),
				AmmPoolId:        1,
				Leverage:         sdk.MustNewDecFromStr("2.0"),
				StopLossPrice:    sdk.MustNewDecFromStr("50.0"),
			},
			true,
			"pool is already leveraged at maximum value",
			func() {
				suite.SetPoolThreshold(sdk.MustNewDecFromStr("0.2"))
			},
		},
		{"Position safety factor too low",
			&types.MsgOpen{
				Creator:          addresses[0].String(),
				CollateralAsset:  ptypes.BaseCurrency,
				CollateralAmount: sdk.NewInt(10_000_000_000),
				AmmPoolId:        1,
				Leverage:         sdk.MustNewDecFromStr("2.0"),
				StopLossPrice:    sdk.MustNewDecFromStr("100.0"),
			},
			true,
			types.ErrPositionUnhealthy.Error(),
			func() {
				suite.SetSafetyFactor(sdk.OneDec().MulInt64(10))
			},
		},
		{"Open new Position with leverage <=1",
			&types.MsgOpen{
				Creator:          addresses[0].String(),
				CollateralAsset:  ptypes.BaseCurrency,
				CollateralAmount: sdk.NewInt(1000),
				AmmPoolId:        1,
				Leverage:         sdk.MustNewDecFromStr("0.0"),
				StopLossPrice:    sdk.MustNewDecFromStr("100.0"),
			},
			true,
			"",
			func() {
			},
		},
		{"Open Position",
			&types.MsgOpen{
				Creator:          addresses[0].String(),
				CollateralAsset:  ptypes.BaseCurrency,
				CollateralAmount: sdk.NewInt(10_000_000),
				AmmPoolId:        1,
				Leverage:         sdk.MustNewDecFromStr("2.0"),
				StopLossPrice:    sdk.MustNewDecFromStr("50.0"),
			},
			false,
			"",
			func() {
				suite.ResetSuite()
				suite.SetupCoinPrices(suite.ctx)
				initializeForOpen(suite, addresses, asset1, asset2)
				suite.SetSafetyFactor(sdk.MustNewDecFromStr("1.1"))
				suite.SetPoolThreshold(sdk.MustNewDecFromStr("0.2"))
			},
		},
		{"Add on already open position Long but with different leverage 10",
			&types.MsgOpen{
				Creator:          addresses[0].String(),
				CollateralAsset:  ptypes.BaseCurrency,
				CollateralAmount: sdk.NewInt(1000),
				AmmPoolId:        1,
				Leverage:         sdk.MustNewDecFromStr("10.0"),
				StopLossPrice:    sdk.MustNewDecFromStr("100.0"),
			},
			false,
			"",
			func() {
			},
		},
		{"Add on already open position Long but with different leverage 20",
			&types.MsgOpen{
				Creator:          addresses[0].String(),
				CollateralAsset:  ptypes.BaseCurrency,
				CollateralAmount: sdk.NewInt(1000),
				AmmPoolId:        1,
				Leverage:         sdk.MustNewDecFromStr("20.0"),
				StopLossPrice:    sdk.MustNewDecFromStr("100.0"),
			},
			false,
			"",
			func() {
			},
		},
		{"Add on already open position Long but with different leverage 1, increase position health",
			&types.MsgOpen{
				Creator:          addresses[0].String(),
				CollateralAsset:  ptypes.BaseCurrency,
				CollateralAmount: sdk.NewInt(1000000),
				AmmPoolId:        1,
				Leverage:         sdk.MustNewDecFromStr("1.0"),
				StopLossPrice:    sdk.MustNewDecFromStr("100.0"),
			},
			false,
			"",
			func() {
				suite.SetSafetyFactor(sdk.MustNewDecFromStr("2.0"))
			},
		},
		{"Add on already open position Long but with different leverage 30",
			&types.MsgOpen{
				Creator:          addresses[0].String(),
				CollateralAsset:  ptypes.BaseCurrency,
				CollateralAmount: sdk.NewInt(1000000),
				AmmPoolId:        1,
				Leverage:         sdk.MustNewDecFromStr("30.0"),
				StopLossPrice:    sdk.MustNewDecFromStr("100.0"),
			},
			true,
			types.ErrPositionUnhealthy.Error(),
			func() {
			},
		},
		{"Low pool health to open position",
			&types.MsgOpen{
				Creator:          addresses[0].String(),
				CollateralAsset:  ptypes.BaseCurrency,
				CollateralAmount: sdk.NewInt(100000000),
				AmmPoolId:        1,
				Leverage:         sdk.MustNewDecFromStr("2.0"),
				StopLossPrice:    sdk.MustNewDecFromStr("50.0"),
			},
			true,
			types.ErrInvalidPosition.Wrapf("pool health too low to open new positions").Error(),
			func() {
				suite.SetSafetyFactor(sdk.MustNewDecFromStr("1.0"))
				suite.SetPoolThreshold(sdk.OneDec())
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			portfolio_old, found := suite.app.TierKeeper.GetPortfolio(suite.ctx, sdk.MustAccAddressFromBech32(tc.input.Creator), suite.app.TierKeeper.GetDateFromContext(suite.ctx))
			_, err := suite.app.LeveragelpKeeper.Open(suite.ctx, tc.input)
			if tc.expectErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				// The new value of the portfolio after the hook is called.
				portfolio_new, _ := suite.app.TierKeeper.GetPortfolio(suite.ctx, sdk.MustAccAddressFromBech32(tc.input.Creator), suite.app.TierKeeper.GetDateFromContext(suite.ctx))
				// Initially, there were no entries for the portfolio
				if !found {
					// The portfolio value changes after the hook is called.
					suite.Require().NotEqual(portfolio_old, portfolio_new)
				}
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestOpen_PoolWithoutBaseCurrencyAsset() {
	suite.ResetSuite()
	// not adding uusdc asset info and price yet
	suite.AddCoinPrices(suite.ctx, []string{ptypes.Elys, ptypes.ATOM, "uusdt"})
	addresses := simapp.AddTestAddrs(suite.app, suite.ctx, 10, sdk.NewInt(1000000))
	asset1 := ptypes.ATOM
	asset2 := ptypes.Elys
	initializeForOpen(suite, addresses, asset1, asset2)
	testCases := []struct {
		name                 string
		input                *types.MsgOpen
		expectErr            bool
		expectErrMsg         string
		prerequisiteFunction func()
	}{
		{"Fail to do JoinPoolNoSwap",
			&types.MsgOpen{
				Creator:          addresses[0].String(),
				CollateralAsset:  ptypes.BaseCurrency,
				CollateralAmount: sdk.NewInt(10000000),
				AmmPoolId:        1,
				Leverage:         sdk.MustNewDecFromStr("2.0"),
				StopLossPrice:    sdk.MustNewDecFromStr("50.0"),
			},
			true,
			"token price not set",
			func() {
			},
		},
		{"Open Position",
			&types.MsgOpen{
				Creator:          addresses[0].String(),
				CollateralAsset:  ptypes.BaseCurrency,
				CollateralAmount: sdk.NewInt(10_000_000),
				AmmPoolId:        1,
				Leverage:         sdk.MustNewDecFromStr("2.0"),
				StopLossPrice:    sdk.MustNewDecFromStr("50.0"),
			},
			true,
			"can't find the PoolAsset",
			func() {
				suite.ResetSuite()
				suite.SetupCoinPrices(suite.ctx)
				initializeForOpen(suite, addresses, asset1, asset2)
				suite.SetSafetyFactor(sdk.MustNewDecFromStr("1.1"))
				suite.SetPoolThreshold(sdk.MustNewDecFromStr("0.2"))
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			_, err := suite.app.LeveragelpKeeper.Open(suite.ctx, tc.input)
			if tc.expectErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
