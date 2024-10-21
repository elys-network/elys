package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"

	"github.com/elys-network/elys/x/perpetual/types"

	simapp "github.com/elys-network/elys/app"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestCheckAndLiquidateUnhealthyPosition(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

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

	// check length of pools
	require.Equal(t, len(pools), 1)

	// check block height
	require.Equal(t, int64(0), ctx.BlockHeight())

	pool, found := amm.GetPool(ctx, poolId)
	require.Equal(t, found, true)

	poolAddress := sdk.MustAccAddressFromBech32(pool.GetAddress())
	require.NoError(t, err)

	app.BankKeeper.SendCoins(ctx, addr[0], poolAddress, sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(1000000))))
	// Balance check before create a perpetual position
	balances := app.BankKeeper.GetAllBalances(ctx, poolAddress)
	require.Equal(t, balances.AmountOf(ptypes.BaseCurrency), sdkmath.NewInt(100000000000))
	require.Equal(t, balances.AmountOf(ptypes.ATOM), sdkmath.NewInt(10000000000))

	// Create a perpetual position open msg
	msg2 := types.NewMsgOpen(
		addr[0].String(),
		types.Position_LONG,
		sdkmath.LegacyNewDec(5),
		ptypes.ATOM,
		sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(100000000)),
		sdkmath.LegacyMustNewDecFromStr(types.TakeProfitPriceDefault),
		sdkmath.LegacyZeroDec(),
	)

	_, err = mk.Open(ctx, msg2, false)
	require.NoError(t, err)

	mtps := mk.GetAllMTPs(ctx)
	require.Equal(t, len(mtps), 1)

	balances = app.BankKeeper.GetAllBalances(ctx, poolAddress)
	require.Equal(t, balances.AmountOf(ptypes.BaseCurrency), sdkmath.NewInt(100100000000))
	require.Equal(t, balances.AmountOf(ptypes.ATOM), sdkmath.NewInt(10000000000))

	_, found = mk.OpenDefineAssetsChecker.GetPool(ctx, pool.PoolId)
	require.Equal(t, found, true)

	err = mk.InvariantCheck(ctx)
	require.Equal(t, err, nil)

	// Set params
	params := mk.GetParams(ctx)
	params.FundingFeeCollectionAddress = addr[1].String()
	params.IncrementalBorrowInterestPaymentFundAddress = addr[2].String()
	params.IncrementalBorrowInterestPaymentFundPercentage = sdkmath.LegacyMustNewDecFromStr("0.5")
	mk.SetParams(ctx, &params)

	mtp := mtps[0]

	perpPool, _ := mk.GetPool(ctx, pool.PoolId)

	err = mk.CheckAndLiquidateUnhealthyPosition(ctx, &mtp, perpPool, pool, ptypes.BaseCurrency, 6)
	require.NoError(t, err)

	// Set borrow interest rate to 100% to test liquidation
	perpPool.BorrowInterestRate = sdkmath.LegacyMustNewDecFromStr("1.0")
	mk.SetPool(ctx, perpPool)

	// Check MTP
	require.Equal(t, types.MTP{
		Address:                        addr[0].String(),
		CollateralAsset:                "uusdc",
		TradingAsset:                   "uatom",
		LiabilitiesAsset:               "uusdc",
		CustodyAsset:                   "uatom",
		Collateral:                     sdkmath.NewInt(100000000),
		Liabilities:                    sdkmath.NewInt(400000000),
		BorrowInterestPaidCollateral:   sdkmath.NewInt(5000000),
		BorrowInterestPaidCustody:      sdkmath.NewInt(4998625),
		BorrowInterestUnpaidCollateral: sdkmath.NewInt(0),
		Custody:                        sdkmath.NewInt(481521968),
		TakeProfitLiabilities:          sdkmath.NewInt(473929244),
		TakeProfitCustody:              sdkmath.NewInt(486520593),
		MtpHealth:                      sdkmath.LegacyMustNewDecFromStr("1.221533382716049383"),
		Position:                       types.Position_LONG,
		Id:                             uint64(1),
		AmmPoolId:                      uint64(1),
		TakeProfitPrice:                sdkmath.LegacyMustNewDecFromStr(types.TakeProfitPriceDefault),
		TakeProfitBorrowRate:           sdkmath.LegacyMustNewDecFromStr("1.0"),
		FundingFeePaidCollateral:       sdkmath.NewInt(0),
		FundingFeePaidCustody:          sdkmath.NewInt(0),
		FundingFeeReceivedCollateral:   sdkmath.NewInt(0),
		FundingFeeReceivedCustody:      sdkmath.NewInt(0),
		OpenPrice:                      sdkmath.LegacyMustNewDecFromStr("1.027705727555914576"),
		LastInterestCalcTime:           uint64(ctx.BlockTime().Unix()),
		LastFundingCalcTime:            uint64(ctx.BlockTime().Unix()),
		StopLossPrice:                  sdkmath.LegacyZeroDec(),
	}, mtp)

	err = mk.CheckAndLiquidateUnhealthyPosition(ctx, &mtp, perpPool, pool, ptypes.BaseCurrency, 6)
	require.NoError(t, err)

	mtps = mk.GetAllMTPs(ctx)
	require.Equal(t, len(mtps), 0)
}

func TestCheckAndLiquidateStopLossPosition(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

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

	// check length of pools
	require.Equal(t, len(pools), 1)

	// check block height
	require.Equal(t, int64(0), ctx.BlockHeight())

	pool, found := amm.GetPool(ctx, poolId)
	require.Equal(t, found, true)

	poolAddress := sdk.MustAccAddressFromBech32(pool.GetAddress())
	require.NoError(t, err)

	app.BankKeeper.SendCoins(ctx, addr[0], poolAddress, sdk.NewCoins(sdk.NewCoin("stake", sdkmath.NewInt(1000000))))
	// Balance check before create a perpetual position
	balances := app.BankKeeper.GetAllBalances(ctx, poolAddress)
	require.Equal(t, balances.AmountOf(ptypes.BaseCurrency), sdkmath.NewInt(100000000000))
	require.Equal(t, balances.AmountOf(ptypes.ATOM), sdkmath.NewInt(10000000000))

	// Create a perpetual position open msg
	msg2 := types.NewMsgOpen(
		addr[0].String(),
		types.Position_LONG,
		sdkmath.LegacyNewDec(5),
		ptypes.ATOM,
		sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(100000000)),
		sdkmath.LegacyMustNewDecFromStr(types.TakeProfitPriceDefault),
		sdkmath.LegacyNewDec(2),
	)

	_, err = mk.Open(ctx, msg2, false)
	require.NoError(t, err)

	mtps := mk.GetAllMTPs(ctx)
	require.Equal(t, len(mtps), 1)

	balances = app.BankKeeper.GetAllBalances(ctx, poolAddress)
	require.Equal(t, balances.AmountOf(ptypes.BaseCurrency), sdkmath.NewInt(100100000000))
	require.Equal(t, balances.AmountOf(ptypes.ATOM), sdkmath.NewInt(10000000000))

	_, found = mk.OpenDefineAssetsChecker.GetPool(ctx, pool.PoolId)
	require.Equal(t, found, true)

	err = mk.InvariantCheck(ctx)
	require.Equal(t, err, nil)

	mtp := mtps[0]

	perpPool, _ := mk.GetPool(ctx, pool.PoolId)

	err = mk.CheckAndCloseAtStopLoss(ctx, &mtp, perpPool, pool, ptypes.BaseCurrency, 6)
	require.NoError(t, err)

	mtps = mk.GetAllMTPs(ctx)
	require.Equal(t, len(mtps), 0)
}

// TODO: Add funding rate tests
