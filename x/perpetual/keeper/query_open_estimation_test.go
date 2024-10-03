package keeper_test

import (
	"cosmossdk.io/math"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/perpetual/types"

	simapp "github.com/elys-network/elys/app"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenEstimation_Long5XAtom100Usdc(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true)

	mk, amm, oracle := app.PerpetualKeeper, app.AmmKeeper, app.OracleKeeper

	// Setup coin prices
	SetupStableCoinPrices(ctx, oracle)

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, math.NewInt(1000000000000))

	// Create a pool
	// Mint 100000USDC
	usdcToken := []sdk.Coin{sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000000000000))}
	// Mint 100000ATOM
	atomToken := []sdk.Coin{sdk.NewCoin(ptypes.ATOM, math.NewInt(1000000000000))}

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
			Weight: math.NewInt(50),
			Token:  sdk.NewCoin(ptypes.ATOM, math.NewInt(600000000000)),
		},
		{
			Weight: math.NewInt(50),
			Token:  sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(100000000000)),
		},
	}

	argSwapFee := math.LegacyMustNewDecFromStr("0.0")
	argExitFee := math.LegacyMustNewDecFromStr("0.0")

	poolParams := &ammtypes.PoolParams{
		UseOracle:                   true,
		ExternalLiquidityRatio:      math.LegacyNewDec(2),
		WeightBreakingFeeMultiplier: math.LegacyZeroDec(),
		WeightBreakingFeeExponent:   math.LegacyNewDecWithPrec(25, 1), // 2.5
		WeightRecoveryFeePortion:    math.LegacyNewDecWithPrec(10, 2), // 10%
		ThresholdWeightDifference:   math.LegacyZeroDec(),
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

	// call min collateral query
	res, err := mk.OpenEstimation(ctx, &types.QueryOpenEstimationRequest{
		Position:        types.Position_LONG,
		Leverage:        math.LegacyMustNewDecFromStr("5.0"),
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(100000000)),
		Discount:        math.LegacyMustNewDecFromStr("0.0"),
		TakeProfitPrice: math.LegacyMustNewDecFromStr("20.0"),
	})
	require.NoError(t, err)
	require.Equal(t, &types.QueryOpenEstimationResponse{
		Position:           types.Position_LONG,
		Leverage:           math.LegacyMustNewDecFromStr("5.0"),
		TradingAsset:       ptypes.ATOM,
		Collateral:         sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(100000000)),
		PositionSize:       sdk.NewCoin(ptypes.ATOM, math.NewInt(499136435)),
		SwapFee:            math.LegacyMustNewDecFromStr("0.001000000000000000"),
		Discount:           math.LegacyMustNewDecFromStr("0.000000000000000000"),
		OpenPrice:          math.LegacyMustNewDecFromStr("1.000000000000000000"),
		TakeProfitPrice:    math.LegacyMustNewDecFromStr("20.00000000000000000"),
		LiquidationPrice:   math.LegacyMustNewDecFromStr("0.799653976372211738"),
		InterestAmount:     types.NewParams().MinBorrowInterestAmount,
		EstimatedPnl:       math.NewInt(474136435),
		EstimatedPnlDenom:  ptypes.ATOM,
		AvailableLiquidity: sdk.NewCoin(ptypes.ATOM, math.NewInt(600000000000)),
		Slippage:           math.LegacyMustNewDecFromStr("0.000727856000000000"),
		WeightBalanceRatio: math.LegacyMustNewDecFromStr("0.000000000000000000"),
		BorrowInterestRate: math.LegacyMustNewDecFromStr("0.000000000000000000"),
		FundingRate:        math.LegacyMustNewDecFromStr("0.000000000000000000"),
		PriceImpact:        math.LegacyMustNewDecFromStr("0.000000000000000000"),
		BorrowFee:          sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(0)),
		FundingFee:         sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(0)),
	}, res)
}

func TestOpenEstimation_Long5XAtom10Atom(t *testing.T) {
	app := simapp.InitElysTestApp(true)
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
	addr := simapp.AddTestAddrs(app, ctx, 1, math.NewInt(1000000000000))

	// Create a pool
	// Mint 10000USDC
	usdcToken := []sdk.Coin{sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000000000000))}
	// Mint 100000ATOM
	atomToken := []sdk.Coin{sdk.NewCoin(ptypes.ATOM, math.NewInt(1000000000000))}

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
			Weight: math.NewInt(50),
			Token:  sdk.NewCoin(ptypes.ATOM, math.NewInt(10000000000)),
		},
		{
			Weight: math.NewInt(50),
			Token:  sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(100000000000)),
		},
	}

	argSwapFee := math.LegacyMustNewDecFromStr("0.0")
	argExitFee := math.LegacyMustNewDecFromStr("0.0")

	poolParams := &ammtypes.PoolParams{
		UseOracle:                   true,
		ExternalLiquidityRatio:      math.LegacyNewDec(2),
		WeightBreakingFeeMultiplier: math.LegacyZeroDec(),
		WeightBreakingFeeExponent:   math.LegacyNewDecWithPrec(25, 1), // 2.5
		WeightRecoveryFeePortion:    math.LegacyNewDecWithPrec(10, 2), // 10%
		ThresholdWeightDifference:   math.LegacyZeroDec(),
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

	// call min collateral query
	res, err := mk.OpenEstimation(ctx, &types.QueryOpenEstimationRequest{
		Position:        types.Position_LONG,
		Leverage:        math.LegacyMustNewDecFromStr("5.0"),
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.ATOM, math.NewInt(10000000)),
		Discount:        math.LegacyMustNewDecFromStr("0.0"),
		TakeProfitPrice: math.LegacyMustNewDecFromStr("20.0"),
	})
	require.NoError(t, err)
	require.Equal(t, &types.QueryOpenEstimationResponse{
		Position:           types.Position_LONG,
		Leverage:           math.LegacyMustNewDecFromStr("5.0"),
		TradingAsset:       ptypes.ATOM,
		Collateral:         sdk.NewCoin(ptypes.ATOM, math.NewInt(10000000)),
		PositionSize:       sdk.NewCoin(ptypes.ATOM, math.NewInt(49858958)),
		SwapFee:            math.LegacyMustNewDecFromStr("0.001000000000000000"),
		Discount:           math.LegacyMustNewDecFromStr("0.000000000000000000"),
		OpenPrice:          math.LegacyMustNewDecFromStr("1.000000000000000000"),
		TakeProfitPrice:    math.LegacyMustNewDecFromStr("20.00000000000000000"),
		LiquidationPrice:   math.LegacyMustNewDecFromStr("0.799662359570370484"),
		EstimatedPnl:       math.NewInt(47361232),
		EstimatedPnlDenom:  ptypes.ATOM,
		InterestAmount:     types.NewParams().MinBorrowInterestAmount,
		AvailableLiquidity: sdk.NewCoin(ptypes.ATOM, math.NewInt(10000000000)),
		Slippage:           math.LegacyMustNewDecFromStr("0.000686040302239768"),
		WeightBalanceRatio: math.LegacyMustNewDecFromStr("0.000000000000000000"),
		BorrowInterestRate: math.LegacyMustNewDecFromStr("0.000000000000000000"),
		FundingRate:        math.LegacyMustNewDecFromStr("0.000000000000000000"),
		PriceImpact:        math.LegacyMustNewDecFromStr("0.001685356924966457"),
		BorrowFee:          sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(0)),
		FundingFee:         sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(0)),
	}, res)
}

func TestOpenEstimation_Short5XAtom10Usdc(t *testing.T) {
	app := simapp.InitElysTestApp(true)
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
	addr := simapp.AddTestAddrs(app, ctx, 1, math.NewInt(1000000000000))

	// Create a pool
	// Mint 10000USDC
	usdcToken := []sdk.Coin{sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000000000000))}
	// Mint 100000ATOM
	atomToken := []sdk.Coin{sdk.NewCoin(ptypes.ATOM, math.NewInt(1000000000000))}

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
			Weight: math.NewInt(50),
			Token:  sdk.NewCoin(ptypes.ATOM, math.NewInt(10000000000)),
		},
		{
			Weight: math.NewInt(50),
			Token:  sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(100000000000)),
		},
	}

	argSwapFee := math.LegacyMustNewDecFromStr("0.0")
	argExitFee := math.LegacyMustNewDecFromStr("0.0")

	poolParams := &ammtypes.PoolParams{
		UseOracle:                   true,
		ExternalLiquidityRatio:      math.LegacyNewDec(2),
		WeightBreakingFeeMultiplier: math.LegacyZeroDec(),
		WeightBreakingFeeExponent:   math.LegacyNewDecWithPrec(25, 1), // 2.5
		WeightRecoveryFeePortion:    math.LegacyNewDecWithPrec(10, 2), // 10%
		ThresholdWeightDifference:   math.LegacyZeroDec(),
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

	// call min collateral query
	res, err := mk.OpenEstimation(ctx, &types.QueryOpenEstimationRequest{
		Position:        types.Position_SHORT,
		Leverage:        math.LegacyMustNewDecFromStr("5.0"),
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(100000000)),
		Discount:        math.LegacyMustNewDecFromStr("0.0"),
		TakeProfitPrice: math.LegacyMustNewDecFromStr("2.0"),
	})
	require.NoError(t, err)
	require.Equal(t, &types.QueryOpenEstimationResponse{
		Position:           types.Position_SHORT,
		Leverage:           math.LegacyMustNewDecFromStr("5.0"),
		TradingAsset:       ptypes.ATOM,
		Collateral:         sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(100000000)),
		PositionSize:       sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(500000000)),
		SwapFee:            math.LegacyMustNewDecFromStr("0.001000000000000000"),
		Discount:           math.LegacyMustNewDecFromStr("0.000000000000000000"),
		OpenPrice:          math.LegacyMustNewDecFromStr("1.000000000000000000"),
		TakeProfitPrice:    math.LegacyMustNewDecFromStr("2.000000000000000000"),
		LiquidationPrice:   math.LegacyMustNewDecFromStr("1.200000000000000000"),
		EstimatedPnl:       math.NewInt(-250000000),
		EstimatedPnlDenom:  ptypes.BaseCurrency,
		InterestAmount:     types.NewParams().MinBorrowInterestAmount,
		AvailableLiquidity: sdk.NewCoin(ptypes.ATOM, math.NewInt(10000000000)),
		Slippage:           math.LegacyMustNewDecFromStr("0.006806806000000000"),
		WeightBalanceRatio: math.LegacyMustNewDecFromStr("0.000000000000000000"),
		BorrowInterestRate: math.LegacyMustNewDecFromStr("0.000000000000000000"),
		FundingRate:        math.LegacyMustNewDecFromStr("0.000000000000000000"),
		PriceImpact:        math.LegacyMustNewDecFromStr("0.007800000000000000"),
		BorrowFee:          sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(0)),
		FundingFee:         sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(0)),
	}, res)
}

func TestOpenEstimation_WrongAsset(t *testing.T) {
	app := simapp.InitElysTestApp(true)
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

	// Generate 1 random account with 1000000000000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, math.NewInt(1000000000000))

	// Create a pool
	// Mint 10000USDC
	usdcToken := []sdk.Coin{sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000000000000))}
	// Mint ATOM
	atomToken := []sdk.Coin{sdk.NewCoin(ptypes.ATOM, math.NewInt(1000000000000))}

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
			Weight: math.NewInt(50),
			Token:  sdk.NewCoin(ptypes.ATOM, math.NewInt(10000000000)),
		},
		{
			Weight: math.NewInt(50),
			Token:  sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(100000000000)),
		},
	}

	argSwapFee := math.LegacyMustNewDecFromStr("0.0")
	argExitFee := math.LegacyMustNewDecFromStr("0.0")

	poolParams := &ammtypes.PoolParams{
		UseOracle:                   true,
		ExternalLiquidityRatio:      math.LegacyNewDec(2),
		WeightBreakingFeeMultiplier: math.LegacyZeroDec(),
		WeightBreakingFeeExponent:   math.LegacyNewDecWithPrec(25, 1), // 2.5
		WeightRecoveryFeePortion:    math.LegacyNewDecWithPrec(10, 2), // 10%
		ThresholdWeightDifference:   math.LegacyZeroDec(),
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

	_, err = mk.OpenEstimation(ctx, &types.QueryOpenEstimationRequest{
		Position:        types.Position_LONG,
		Leverage:        math.LegacyMustNewDecFromStr("5.0"),
		TradingAsset:    ptypes.BaseCurrency,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(10000000)),
		Discount:        math.LegacyMustNewDecFromStr("0.0"),
		TakeProfitPrice: math.LegacyMustNewDecFromStr("20.0"),
	})

	assert.Error(t, err)
	assert.Equal(t, "invalid operation: the borrowed asset cannot be the base currency: invalid borrowing asset", err.Error())

	_, err = mk.OpenEstimation(ctx, &types.QueryOpenEstimationRequest{
		Position:        types.Position_LONG,
		Leverage:        math.LegacyMustNewDecFromStr("5.0"),
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.Eden, math.NewInt(10000000)),
		Discount:        math.LegacyMustNewDecFromStr("0.0"),
		TakeProfitPrice: math.LegacyMustNewDecFromStr("20.0"),
	})

	assert.Error(t, err)
	assert.Equal(t, "invalid collateral: collateral must either match the borrowed asset or be the base currency: invalid borrowing asset", err.Error())

	_, err = mk.OpenEstimation(ctx, &types.QueryOpenEstimationRequest{
		Position:        types.Position_SHORT,
		Leverage:        math.LegacyMustNewDecFromStr("5.0"),
		TradingAsset:    ptypes.BaseCurrency,
		Collateral:      sdk.NewCoin(ptypes.ATOM, math.NewInt(10000000)),
		Discount:        math.LegacyMustNewDecFromStr("0.0"),
		TakeProfitPrice: math.LegacyMustNewDecFromStr("20.0"),
	})

	assert.Error(t, err)
	assert.Equal(t, "borrowing not allowed: cannot take a short position against the base currency: invalid borrowing asset", err.Error())

	_, err = mk.OpenEstimation(ctx, &types.QueryOpenEstimationRequest{
		Position:        types.Position_SHORT,
		Leverage:        math.LegacyMustNewDecFromStr("5.0"),
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.ATOM, math.NewInt(10000000)),
		Discount:        math.LegacyMustNewDecFromStr("0.0"),
		TakeProfitPrice: math.LegacyMustNewDecFromStr("20.0"),
	})

	assert.Error(t, err)
	assert.Equal(t, "invalid operation: collateral asset cannot be identical to the borrowed asset for a short position: invalid collateral asset", err.Error())
}
