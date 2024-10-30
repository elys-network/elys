package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/ed25519"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	simapp "github.com/elys-network/elys/app"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	leveragelpmodulekeeper "github.com/elys-network/elys/x/leveragelp/keeper"
	leveragelpmoduletypes "github.com/elys-network/elys/x/leveragelp/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func (suite *PerpetualKeeperTestSuite) TestCheckAndLiquidateUnhealthyPosition() {
	app := suite.app
	ctx := suite.ctx

	mk, amm, oracle := app.PerpetualKeeper, app.AmmKeeper, app.OracleKeeper

	// Setup coin prices
	SetupStableCoinPrices(ctx, oracle)

	// Set asset profile
	app.AssetprofileKeeper.SetEntry(ctx, assetprofiletypes.Entry{
		BaseDenom: ptypes.BaseCurrency,
		Denom:     ptypes.BaseCurrency,
		Decimals:  6,
	})
	app.AssetprofileKeeper.SetEntry(ctx, assetprofiletypes.Entry{
		BaseDenom: ptypes.ATOM,
		Denom:     ptypes.ATOM,
		Decimals:  6,
	})

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 3, sdkmath.NewInt(1000000000000))

	// Create a pool
	// Mint 100000USDC
	usdcToken := []sdk.Coin{sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(200000000000))}
	// Mint 100000ATOM
	atomToken := []sdk.Coin{sdk.NewCoin(ptypes.ATOM, sdkmath.NewInt(200000000000))}

	err := app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, usdcToken)
	suite.Require().NoError(err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[0], usdcToken)
	suite.Require().NoError(err)

	err = app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, atomToken)
	suite.Require().NoError(err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[0], atomToken)
	suite.Require().NoError(err)

	poolAssets := []ammtypes.PoolAsset{
		{
			Weight: sdkmath.NewInt(50),
			Token:  sdk.NewCoin(ptypes.ATOM, sdkmath.NewInt(10000000000)),
		},
		{
			Weight: sdkmath.NewInt(50),
			Token:  sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(100000000000)),
		},
	}

	argSwapFee := sdkmath.LegacyMustNewDecFromStr("0.0")
	argExitFee := sdkmath.LegacyMustNewDecFromStr("0.0")

	poolParams := &ammtypes.PoolParams{
		UseOracle:                   true,
		ExternalLiquidityRatio:      sdkmath.LegacyNewDec(2),
		WeightBreakingFeeMultiplier: sdkmath.LegacyZeroDec(),
		WeightBreakingFeeExponent:   sdkmath.LegacyNewDecWithPrec(25, 1), // 2.5
		WeightRecoveryFeePortion:    sdkmath.LegacyNewDecWithPrec(10, 2), // 10%
		ThresholdWeightDifference:   sdkmath.LegacyZeroDec(),
		SwapFee:                     argSwapFee,
		ExitFee:                     argExitFee,
		FeeDenom:                    ptypes.BaseCurrency,
	}

	msg := ammtypes.NewMsgCreatePool(
		addr[0].String(),
		poolParams,
		poolAssets,
	)

	// Create a ATOM+USDC pool
	poolId, err := amm.CreatePool(ctx, msg)
	suite.Require().NoError(err)
	suite.Require().Equal(poolId, uint64(1))

	pools := amm.GetAllPool(ctx)

	// check length of pools
	suite.Require().Equal(len(pools), 1)

	// check block height
	suite.Require().Equal(int64(0), ctx.BlockHeight())

	pool, found := amm.GetPool(ctx, poolId)
	suite.Require().Equal(found, true)

	poolAddress := sdk.MustAccAddressFromBech32(pool.GetAddress())
	suite.Require().NoError(err)

	err = app.BankKeeper.SendCoins(ctx, addr[0], poolAddress, sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(1000000))))
	if err != nil {
		return
	}
	// Balance check before create a perpetual position
	balances := app.BankKeeper.GetAllBalances(ctx, poolAddress)
	suite.Require().Equal(balances.AmountOf(ptypes.BaseCurrency), sdkmath.NewInt(100000000000))
	suite.Require().Equal(balances.AmountOf(ptypes.ATOM), sdkmath.NewInt(10000000000))

	// Create a perpetual position open msg
	msg2 := types.NewMsgOpen(
		addr[0].String(),
		types.Position_LONG,
		sdkmath.LegacyNewDec(5),
		1,
		ptypes.ATOM,
		sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(100000000)),
		types.TakeProfitPriceDefault,
		sdk.ZeroDec(),
	)

	params := app.PerpetualKeeper.GetParams(ctx)
	params.WhitelistingEnabled = true
	err = app.PerpetualKeeper.SetParams(ctx, &params)
	suite.Require().NoError(err)
	app.PerpetualKeeper.WhitelistAddress(ctx, sdk.MustAccAddressFromBech32(msg2.Creator))
	_, err = mk.Open(ctx, msg2, false)
	suite.Require().NoError(err)

	mtps := mk.GetAllMTPs(ctx)
	suite.Require().Equal(len(mtps), 1)

	balances = app.BankKeeper.GetAllBalances(ctx, poolAddress)
	suite.Require().Equal(balances.AmountOf(ptypes.BaseCurrency), sdkmath.NewInt(100100000000))
	suite.Require().Equal(balances.AmountOf(ptypes.ATOM), sdkmath.NewInt(10000000000))

	_, found = mk.GetPool(ctx, pool.PoolId)
	suite.Require().Equal(found, true)

	// Set params
	params = mk.GetParams(ctx)
	params.ForceCloseFundAddress = addr[1].String()
	params.IncrementalBorrowInterestPaymentFundAddress = addr[2].String()
	params.IncrementalBorrowInterestPaymentFundPercentage = sdkmath.LegacyMustNewDecFromStr("0.5")
	mk.SetParams(ctx, &params)

	mtp := mtps[0]

	perpPool, _ := mk.GetPool(ctx, pool.PoolId)

	err = mk.CheckAndLiquidateUnhealthyPosition(ctx, &mtp, perpPool, pool, ptypes.BaseCurrency, 6)
	suite.Require().NoError(err)

	// Set borrow interest rate to 100% to test liquidation
	perpPool.BorrowInterestRate = sdkmath.LegacyMustNewDecFromStr("1.0")
	mk.SetPool(ctx, perpPool)

	// Check MTP
	suite.Require().Equal(types.MTP{
		Address:                       addr[0].String(),
		CollateralAsset:               "uusdc",
		TradingAsset:                  "uatom",
		LiabilitiesAsset:              "uusdc",
		CustodyAsset:                  "uatom",
		Collateral:                    sdk.NewInt(100000000),
		Liabilities:                   sdk.NewInt(400000000),
		BorrowInterestPaidCustody:     sdk.NewInt(4998625),
		BorrowInterestUnpaidLiability: sdk.NewInt(0),
		Custody:                       sdk.NewInt(481521968),
		TakeProfitLiabilities:         sdk.NewInt(473929244),
		TakeProfitCustody:             sdk.NewInt(486520593),
		MtpHealth:                     sdk.MustNewDecFromStr("1.221533382716049383"),
		Position:                      types.Position_LONG,
		Id:                            uint64(1),
		AmmPoolId:                     uint64(1),
		TakeProfitPrice:               types.TakeProfitPriceDefault,
		TakeProfitBorrowFactor:        sdk.MustNewDecFromStr("1.0"),
		FundingFeePaidCustody:         sdk.NewInt(0),
		FundingFeeReceivedCustody:     sdk.NewInt(0),
		OpenPrice:                     sdk.MustNewDecFromStr("1.027705727555914576"),
		LastInterestCalcTime:          uint64(ctx.BlockTime().Unix()),
		LastFundingCalcTime:           uint64(ctx.BlockTime().Unix()),
		StopLossPrice:                 sdk.ZeroDec(),
	}, mtp)

	err = mk.CheckAndLiquidateUnhealthyPosition(ctx, &mtp, perpPool, pool, ptypes.BaseCurrency, 6)
	suite.Require().NoError(err)

	mtps = mk.GetAllMTPs(ctx)
	suite.Require().Equal(len(mtps), 0)
}

func TestCheckAndCloseAtTakeProfit(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)
	simapp.SetStakingParam(app, ctx)
	simapp.SetPerpetualParams(app, ctx)

	mk, amm, oracle := app.PerpetualKeeper, app.AmmKeeper, app.OracleKeeper

	// Setup coin prices
	SetupStableCoinPrices(ctx, oracle)

	// Set asset profile
	app.AssetprofileKeeper.SetEntry(ctx, assetprofiletypes.Entry{
		BaseDenom: ptypes.BaseCurrency,
		Denom:     ptypes.BaseCurrency,
		Decimals:  6,
	})
	app.AssetprofileKeeper.SetEntry(ctx, assetprofiletypes.Entry{
		BaseDenom: ptypes.ATOM,
		Denom:     ptypes.ATOM,
		Decimals:  6,
	})

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 3, sdkmath.NewInt(1000000000000))

	// Create a pool
	// Mint 100000USDC
	usdcToken := []sdk.Coin{sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(200000000000))}
	// Mint 100000ATOM
	atomToken := []sdk.Coin{sdk.NewCoin(ptypes.ATOM, sdkmath.NewInt(200000000000))}

	err := app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, usdcToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[0], usdcToken)
	require.NoError(t, err)

	err = app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, atomToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[0], atomToken)
	require.NoError(t, err)

	poolAssets := []ammtypes.PoolAsset{
		{
			Weight: sdkmath.NewInt(50),
			Token:  sdk.NewCoin(ptypes.ATOM, sdkmath.NewInt(10000000000)),
		},
		{
			Weight: sdkmath.NewInt(50),
			Token:  sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(100000000000)),
		},
	}

	argSwapFee := sdkmath.LegacyMustNewDecFromStr("0.0")
	argExitFee := sdkmath.LegacyMustNewDecFromStr("0.0")

	poolParams := &ammtypes.PoolParams{
		UseOracle:                   true,
		ExternalLiquidityRatio:      sdkmath.LegacyNewDec(2),
		WeightBreakingFeeMultiplier: sdkmath.LegacyZeroDec(),
		WeightBreakingFeeExponent:   sdkmath.LegacyNewDecWithPrec(25, 1), // 2.5
		WeightRecoveryFeePortion:    sdkmath.LegacyNewDecWithPrec(10, 2), // 10%
		ThresholdWeightDifference:   sdkmath.LegacyZeroDec(),
		SwapFee:                     argSwapFee,
		ExitFee:                     argExitFee,
		FeeDenom:                    ptypes.BaseCurrency,
	}

	msg := ammtypes.NewMsgCreatePool(
		addr[0].String(),
		poolParams,
		poolAssets,
	)

	// Create a ATOM+USDC pool
	poolId, err := amm.CreatePool(ctx, msg)
	require.NoError(t, err)
	require.Equal(t, poolId, uint64(1))

	pools := amm.GetAllPool(ctx)

	enablePoolMsg := leveragelpmoduletypes.MsgAddPool{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Pool: leveragelpmoduletypes.AddPool{
			poolId,
			math.LegacyMustNewDecFromStr("10"),
		},
	}
	_, err = leveragelpmodulekeeper.NewMsgServerImpl(*app.LeveragelpKeeper).AddPool(ctx, &enablePoolMsg)
	require.NoError(t, err)

	// check length of pools
	require.Equal(t, len(pools), 1)

	// check block height
	require.Equal(t, int64(0), ctx.BlockHeight())

	pool, found := amm.GetPool(ctx, poolId)
	require.Equal(t, found, true)

	poolAddress := sdk.MustAccAddressFromBech32(pool.GetAddress())
	require.NoError(t, err)

	err = app.BankKeeper.SendCoins(ctx, addr[0], poolAddress, sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(1000000))))
	if err != nil {
		return
	}
	// Balance check before create a perpetual position
	balances := app.BankKeeper.GetAllBalances(ctx, poolAddress)
	require.Equal(t, balances.AmountOf(ptypes.BaseCurrency), sdkmath.NewInt(100000000000))
	require.Equal(t, balances.AmountOf(ptypes.ATOM), sdkmath.NewInt(10000000000))

	// Create a perpetual position open msg
	msg2 := types.NewMsgOpen(
		addr[0].String(),
		types.Position_LONG,
		sdkmath.LegacyNewDec(5),
		1,
		ptypes.ATOM,
		sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(100000000)),
		sdk.MustNewDecFromStr("8"),
		sdk.ZeroDec(),
	)

	_, err = mk.Open(ctx, msg2, false)
	require.NoError(t, err)

	mtps := mk.GetAllMTPs(ctx)
	require.Equal(t, len(mtps), 1)

	balances = app.BankKeeper.GetAllBalances(ctx, poolAddress)
	require.Equal(t, balances.AmountOf(ptypes.BaseCurrency), sdkmath.NewInt(100100000000))
	require.Equal(t, balances.AmountOf(ptypes.ATOM), sdkmath.NewInt(10000000000))

	_, found = mk.GetPool(ctx, pool.PoolId)
	require.Equal(t, found, true)

	mtp := mtps[0]

	perpPool, _ := mk.GetPool(ctx, pool.PoolId)

	err = mk.CheckAndCloseAtTakeProfit(ctx, &mtp, perpPool, ptypes.BaseCurrency)
	require.Error(t, err)

	// Set price above target price
	provider := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "uatom",
		Price:     sdk.MustNewDecFromStr("8.1"),
		Source:    "uatom",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})

	err = mk.CheckAndCloseAtTakeProfit(ctx, &mtp, perpPool, ptypes.BaseCurrency)
	require.NoError(t, err)

	mtps = mk.GetAllMTPs(ctx)
	require.Equal(t, len(mtps), 0)
}

func (suite *PerpetualKeeperTestSuite) TestCheckAndLiquidateStopLossPosition() {
	suite.ResetSuite()
	app := suite.app
	ctx := suite.ctx

	mk, amm, oracle, _ := app.PerpetualKeeper, app.AmmKeeper, app.OracleKeeper, app.LeveragelpKeeper

	// Setup coin prices
	suite.SetupCoinPrices()

	// Set asset profile
	app.AssetprofileKeeper.SetEntry(ctx, assetprofiletypes.Entry{
		BaseDenom: ptypes.BaseCurrency,
		Denom:     ptypes.BaseCurrency,
		Decimals:  6,
	})
	app.AssetprofileKeeper.SetEntry(ctx, assetprofiletypes.Entry{
		BaseDenom: ptypes.ATOM,
		Denom:     ptypes.ATOM,
		Decimals:  6,
	})

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 3, sdk.NewInt(1000000000000))

	// Create a pool
	// Mint 100000USDC
	usdcToken := []sdk.Coin{sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(200000000000))}
	// Mint 100000ATOM
	atomToken := []sdk.Coin{sdk.NewCoin(ptypes.ATOM, sdk.NewInt(200000000000))}

	err := app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, usdcToken)
	suite.Require().NoError(err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[0], usdcToken)
	suite.Require().NoError(err)

	err = app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, atomToken)
	suite.Require().NoError(err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[0], atomToken)
	suite.Require().NoError(err)

	poolAssets := []ammtypes.PoolAsset{
		{
			Weight: sdk.NewInt(50),
			Token:  sdk.NewCoin(ptypes.ATOM, sdk.NewInt(10000000000)),
		},
		{
			Weight: sdk.NewInt(50),
			Token:  sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(100000000000)),
		},
	}

	argSwapFee := sdk.MustNewDecFromStr("0.0")
	argExitFee := sdk.MustNewDecFromStr("0.0")

	poolParams := &ammtypes.PoolParams{
		UseOracle:                   true,
		ExternalLiquidityRatio:      sdk.NewDec(2),
		WeightBreakingFeeMultiplier: sdk.ZeroDec(),
		WeightBreakingFeeExponent:   sdk.NewDecWithPrec(25, 1), // 2.5
		WeightRecoveryFeePortion:    sdk.NewDecWithPrec(10, 2), // 10%
		ThresholdWeightDifference:   sdk.ZeroDec(),
		SwapFee:                     argSwapFee,
		ExitFee:                     argExitFee,
		FeeDenom:                    ptypes.BaseCurrency,
	}

	msg := ammtypes.NewMsgCreatePool(
		addr[0].String(),
		poolParams,
		poolAssets,
	)

	// Create a ATOM+USDC pool
	poolId, err := amm.CreatePool(ctx, msg)
	suite.Require().NoError(err)
	suite.Require().Equal(poolId, uint64(1))

	pools := amm.GetAllPool(ctx)

	// check length of pools
	suite.Require().Equal(len(pools), 1)

	// check block height
	suite.Require().Equal(int64(0), ctx.BlockHeight())

	ammPool, found := amm.GetPool(ctx, poolId)
	suite.Require().Equal(found, true)

	poolAddress := sdk.MustAccAddressFromBech32(ammPool.GetAddress())
	suite.Require().NoError(err)

	err = app.BankKeeper.SendCoins(ctx, addr[0], poolAddress, sdk.NewCoins(sdk.NewCoin("uelys", sdk.NewInt(1000000))))
	suite.Require().NoError(err)
	// Balance check before create a perpetual position
	balances := app.BankKeeper.GetAllBalances(ctx, poolAddress)
	suite.Require().Equal(balances.AmountOf(ptypes.BaseCurrency), sdk.NewInt(100000000000))
	suite.Require().Equal(balances.AmountOf(ptypes.ATOM), sdk.NewInt(10000000000))

	enablePoolMsg := leveragelpmoduletypes.MsgAddPool{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Pool: leveragelpmoduletypes.AddPool{
			poolId,
			math.LegacyMustNewDecFromStr("10"),
		},
	}
	_, err = leveragelpmodulekeeper.NewMsgServerImpl(*app.LeveragelpKeeper).AddPool(ctx, &enablePoolMsg)
	suite.Require().NoError(err)
	tradingAssetPrice, err := app.PerpetualKeeper.GetAssetPrice(ctx, ptypes.ATOM)
	suite.Require().NoError(err)
	// Create a perpetual position open msg
	msg2 := types.NewMsgOpen(
		addr[0].String(),
		types.Position_LONG,
		sdk.NewDec(5),
		1,
		ptypes.ATOM,
		sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(100000000)),
		tradingAssetPrice.MulInt64(10),
		tradingAssetPrice.QuoInt64(2),
	)
	params := app.PerpetualKeeper.GetParams(ctx)
	params.WhitelistingEnabled = true
	err = app.PerpetualKeeper.SetParams(ctx, &params)
	suite.Require().NoError(err)
	app.PerpetualKeeper.WhitelistAddress(ctx, sdk.MustAccAddressFromBech32(msg2.Creator))
	_, err = mk.Open(ctx, msg2, false)
	suite.Require().NoError(err)

	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "ATOM",
		Price:     tradingAssetPrice.QuoInt64(4),
		Source:    "elys",
		Provider:  authtypes.NewModuleAddress("provider").String(),
		Timestamp: uint64(ctx.BlockTime().Unix() + 6),
	})

	mtps := mk.GetAllMTPs(ctx)
	suite.Require().Equal(len(mtps), 1)

	balances = app.BankKeeper.GetAllBalances(ctx, poolAddress)
	suite.Require().Equal(balances.AmountOf(ptypes.BaseCurrency), sdk.NewInt(100100000000))
	suite.Require().Equal(balances.AmountOf(ptypes.ATOM), sdk.NewInt(10000000000))

	_, found = mk.GetPool(ctx, ammPool.PoolId)
	suite.Require().Equal(found, true)

	mtp := mtps[0]

	perpPool, _ := mk.GetPool(ctx, ammPool.PoolId)

	err = mk.CheckAndCloseAtStopLoss(ctx, &mtp, perpPool, ptypes.BaseCurrency)
	suite.Require().NoError(err)

	mtps = mk.GetAllMTPs(ctx)
	suite.Require().Equal(len(mtps), 0)
}

// TODO: Add funding rate tests
