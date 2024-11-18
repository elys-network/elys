package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	leveragelpmodulekeeper "github.com/elys-network/elys/x/leveragelp/keeper"
	leveragelpmoduletypes "github.com/elys-network/elys/x/leveragelp/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	"github.com/elys-network/elys/x/perpetual/types"

	"github.com/cometbft/cometbft/crypto/ed25519"
	simapp "github.com/elys-network/elys/app"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenEstimation_Long5XAtom100Usdc(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)
	simapp.SetStakingParam(app, ctx)
	simapp.SetPerpetualParams(app, ctx)
	simapp.SetupAssetProfile(app, ctx)

	mk, amm, oracle := app.PerpetualKeeper, app.AmmKeeper, app.OracleKeeper

	// Setup coin prices
	SetupStableCoinPrices(ctx, oracle)
	tradingAssetPrice, err := app.PerpetualKeeper.GetAssetPrice(ctx, ptypes.ATOM)
	require.NoError(t, err)
	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, math.NewInt(1000000000000))

	// Create a pool
	// Mint 100000USDC
	usdcToken := []sdk.Coin{sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000000000000))}
	// Mint 100000ATOM
	atomToken := []sdk.Coin{sdk.NewCoin(ptypes.ATOM, math.NewInt(1000000000000))}

	err = app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, usdcToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[0], usdcToken)
	require.NoError(t, err)

	err = app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, atomToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[0], atomToken)
	require.NoError(t, err)

	poolAssets := []ammtypes.PoolAsset{
		{
			Weight:                 math.NewInt(50),
			Token:                  sdk.NewCoin(ptypes.ATOM, math.NewInt(600000000000)),
			ExternalLiquidityRatio: math.LegacyNewDec(2),
		},
		{
			Weight:                 math.NewInt(50),
			Token:                  sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(100000000000)),
			ExternalLiquidityRatio: math.LegacyNewDec(2),
		},
	}

	argSwapFee := math.LegacyMustNewDecFromStr("0.0")
	argExitFee := math.LegacyMustNewDecFromStr("0.0")

	poolParams := &ammtypes.PoolParams{
		UseOracle:                   true,
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

	enablePoolMsg := leveragelpmoduletypes.MsgAddPool{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Pool: leveragelpmoduletypes.AddPool{
			AmmPoolId:   poolId,
			LeverageMax: math.LegacyMustNewDecFromStr("10"),
		},
	}
	_, err = leveragelpmodulekeeper.NewMsgServerImpl(*app.LeveragelpKeeper).AddPool(ctx, &enablePoolMsg)
	require.NoError(t, err)

	// call min collateral query
	res, err := mk.OpenEstimation(ctx, &types.QueryOpenEstimationRequest{
		PoolId:          1,
		Position:        types.Position_LONG,
		Leverage:        math.LegacyMustNewDecFromStr("5.0"),
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(100_000_000)),
		Address:         "",
		TakeProfitPrice: tradingAssetPrice.MulInt64(3),
	})
	require.NoError(t, err)

	require.Equal(t, &types.QueryOpenEstimationResponse{
		Position:           types.Position_LONG,
		EffectiveLeverage:  math.LegacyMustNewDecFromStr("5.033154750282434265"),
		TradingAsset:       ptypes.ATOM,
		Collateral:         sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(100_000_000)),
		HourlyInterestRate: math.LegacyZeroDec(),
		PositionSize:       sdk.NewCoin(ptypes.ATOM, math.NewInt(99_835_589)),
		OpenPrice:          math.LegacyMustNewDecFromStr("5.008234087746004083"),
		TakeProfitPrice:    tradingAssetPrice.MulInt64(3),
		LiquidationPrice:   math.LegacyMustNewDecFromStr("4.106751951951723348"),
		EstimatedPnl:       sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(997_533_835)),
		AvailableLiquidity: sdk.NewCoin(ptypes.ATOM, math.NewInt(600000000000)),
		Slippage:           math.LegacyMustNewDecFromStr("0.000644750000000000"),
		BorrowInterestRate: math.LegacyMustNewDecFromStr("0.000000000000000000"),
		FundingRate:        math.LegacyMustNewDecFromStr("0.000000000000000000"),
		PriceImpact:        math.LegacyMustNewDecFromStr("-0.001646817549200817"),
		Custody:            sdk.NewCoin(ptypes.ATOM, math.NewInt(99_835_589)),
		Liabilities:        sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(400_000_000)),
		WeightBreakingFee:  math.LegacyZeroDec(),
	}, res)
}

func TestOpenEstimation_Long5XAtom10Atom(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)
	simapp.SetStakingParam(app, ctx)
	simapp.SetPerpetualParams(app, ctx)
	simapp.SetupAssetProfile(app, ctx)

	mk, amm, oracle := app.PerpetualKeeper, app.AmmKeeper, app.OracleKeeper

	// Setup coin prices
	SetupStableCoinPrices(ctx, oracle)
	tradingAssetPrice, err := app.PerpetualKeeper.GetAssetPrice(ctx, ptypes.ATOM)
	require.NoError(t, err)
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

	err = app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, usdcToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[0], usdcToken)
	require.NoError(t, err)

	err = app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, atomToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[0], atomToken)
	require.NoError(t, err)

	poolAssets := []ammtypes.PoolAsset{
		{
			Weight:                 math.NewInt(50),
			Token:                  sdk.NewCoin(ptypes.ATOM, math.NewInt(10000000000)),
			ExternalLiquidityRatio: math.LegacyNewDec(2),
		},
		{
			Weight:                 math.NewInt(50),
			Token:                  sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(100000000000)),
			ExternalLiquidityRatio: math.LegacyNewDec(2),
		},
	}

	argSwapFee := math.LegacyMustNewDecFromStr("0.0")
	argExitFee := math.LegacyMustNewDecFromStr("0.0")

	poolParams := &ammtypes.PoolParams{
		UseOracle:                   true,
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

	enablePoolMsg := leveragelpmoduletypes.MsgAddPool{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Pool: leveragelpmoduletypes.AddPool{
			AmmPoolId:   poolId,
			LeverageMax: math.LegacyMustNewDecFromStr("10"),
		},
	}
	_, err = leveragelpmodulekeeper.NewMsgServerImpl(*app.LeveragelpKeeper).AddPool(ctx, &enablePoolMsg)
	require.NoError(t, err)

	// check length of pools
	require.Equal(t, len(pools), 1)

	// call min collateral query
	res, err := mk.OpenEstimation(ctx, &types.QueryOpenEstimationRequest{
		PoolId:          1,
		Position:        types.Position_LONG,
		Leverage:        math.LegacyMustNewDecFromStr("5.0"),
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.ATOM, math.NewInt(10_000_000)),
		Address:         "",
		TakeProfitPrice: tradingAssetPrice.MulInt64(3),
	})
	require.NoError(t, err)
	require.Equal(t, &types.QueryOpenEstimationResponse{
		Position:           types.Position_LONG,
		EffectiveLeverage:  math.LegacyMustNewDecFromStr("5.035307678145308466"),
		TradingAsset:       ptypes.ATOM,
		Collateral:         sdk.NewCoin(ptypes.ATOM, math.NewInt(10_000_000)),
		PositionSize:       sdk.NewCoin(ptypes.ATOM, math.NewInt(50_000_000)),
		OpenPrice:          math.LegacyMustNewDecFromStr("5.008765025000000000"),
		TakeProfitPrice:    tradingAssetPrice.MulInt64(3),
		LiquidationPrice:   math.LegacyMustNewDecFromStr("4.107187320500000000"),
		EstimatedPnl:       sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(399649399)),
		HourlyInterestRate: math.LegacyZeroDec(),
		AvailableLiquidity: sdk.NewCoin(ptypes.ATOM, math.NewInt(10000000000)),
		Slippage:           math.LegacyMustNewDecFromStr("0.000751255000000000"),
		BorrowInterestRate: math.LegacyMustNewDecFromStr("0.000000000000000000"),
		FundingRate:        math.LegacyMustNewDecFromStr("0.000000000000000000"),
		PriceImpact:        math.LegacyMustNewDecFromStr("-0.001753005000000000"),
		Custody:            sdk.NewCoin(ptypes.ATOM, math.NewInt(50_000_000)),
		Liabilities:        sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(200_350_601)),
		WeightBreakingFee:  math.LegacyZeroDec(),
	}, res)
}

func TestOpenEstimation_Long10XAtom1000Usdc(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

	mk, amm, oracle := app.PerpetualKeeper, app.AmmKeeper, app.OracleKeeper

	// Setup coin prices
	SetupStableCoinPrices(ctx, oracle)

	provider := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "ATOM",
		Price:     math.LegacyMustNewDecFromStr("4.39"),
		Source:    "atom",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "uatom",
		Price:     math.LegacyMustNewDecFromStr("4.39"),
		Source:    "uatom",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, math.NewInt(1_000_000_000000))

	// Create a pool
	// Mint 100000USDC
	usdcToken := []sdk.Coin{sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1_000_000_000000))}
	// Mint 100000ATOM
	atomToken := []sdk.Coin{sdk.NewCoin(ptypes.ATOM, math.NewInt(1_000_000_000000))}

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
			Weight:                 math.NewInt(50),
			Token:                  sdk.NewCoin(ptypes.ATOM, math.NewInt(600_000_000000)),
			ExternalLiquidityRatio: math.LegacyNewDec(2),
		},
		{
			Weight:                 math.NewInt(50),
			Token:                  sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(100_000_000000)),
			ExternalLiquidityRatio: math.LegacyNewDec(2),
		},
	}

	argSwapFee := math.LegacyMustNewDecFromStr("0.0")
	argExitFee := math.LegacyMustNewDecFromStr("0.0")

	poolParams := &ammtypes.PoolParams{
		UseOracle:                   true,
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

	enablePoolMsg := leveragelpmoduletypes.MsgAddPool{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Pool: leveragelpmoduletypes.AddPool{
			AmmPoolId:   poolId,
			LeverageMax: math.LegacyMustNewDecFromStr("10"),
		},
	}
	_, err = leveragelpmodulekeeper.NewMsgServerImpl(*app.LeveragelpKeeper).AddPool(ctx, &enablePoolMsg)
	require.NoError(t, err)

	// check length of pools
	require.Equal(t, len(pools), 1)
	tradingAssetPrice, err := app.PerpetualKeeper.GetAssetPrice(ctx, ptypes.ATOM)
	require.NoError(t, err)
	// call min collateral query	tradingAssetPrice := app.OracleKeeper.GetAssetPriceFromDenom(ctx, ptypes.ATOM)
	res, err := mk.OpenEstimation(ctx, &types.QueryOpenEstimationRequest{
		PoolId:          1,
		Position:        types.Position_LONG,
		Leverage:        math.LegacyMustNewDecFromStr("10.0"),
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1_000_000000)),
		Address:         "",
		TakeProfitPrice: tradingAssetPrice.MulInt64(3),
	})
	require.NoError(t, err)
	require.Equal(t, &types.QueryOpenEstimationResponse{
		Position:           types.Position_LONG,
		EffectiveLeverage:  math.LegacyMustNewDecFromStr("11.409127722426274478"),
		TradingAsset:       ptypes.ATOM,
		Collateral:         sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1_000_000_000)),
		PositionSize:       sdk.NewCoin(ptypes.ATOM, math.NewInt(2247067372)),
		OpenPrice:          math.LegacyMustNewDecFromStr("4.450244850068518551"),
		TakeProfitPrice:    tradingAssetPrice.MulInt64(3),
		LiquidationPrice:   math.LegacyMustNewDecFromStr("4.105350874188208363"),
		HourlyInterestRate: math.LegacyZeroDec(),
		EstimatedPnl:       sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(19593877289)),
		AvailableLiquidity: sdk.NewCoin(ptypes.ATOM, math.NewInt(600_000_000000)),
		Slippage:           math.LegacyMustNewDecFromStr("0.012549973525000000"),
		BorrowInterestRate: math.LegacyMustNewDecFromStr("0.000000000000000000"),
		FundingRate:        math.LegacyMustNewDecFromStr("0.000000000000000000"),
		PriceImpact:        math.LegacyMustNewDecFromStr("-0.013723200471188736"),
		Custody:            sdk.NewCoin(ptypes.ATOM, math.NewInt(2_247_067_372)),
		Liabilities:        sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(9000000000)),
		WeightBreakingFee:  math.LegacyZeroDec(),
	}, res)
}

func TestOpenEstimation_Short5XAtom10Usdc(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)
	simapp.SetStakingParam(app, ctx)
	simapp.SetPerpetualParams(app, ctx)
	simapp.SetupAssetProfile(app, ctx)

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
	provider := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "ATOM",
		Price:     math.LegacyNewDec(5),
		Source:    "atom",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "uatom",
		Price:     math.LegacyNewDec(5),
		Source:    "uatom",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
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
			Weight:                 math.NewInt(50),
			Token:                  sdk.NewCoin(ptypes.ATOM, math.NewInt(10000000000)),
			ExternalLiquidityRatio: math.LegacyNewDec(2),
		},
		{
			Weight:                 math.NewInt(50),
			Token:                  sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(100000000000)),
			ExternalLiquidityRatio: math.LegacyNewDec(2),
		},
	}

	argSwapFee := math.LegacyMustNewDecFromStr("0.0")
	argExitFee := math.LegacyMustNewDecFromStr("0.0")

	poolParams := &ammtypes.PoolParams{
		UseOracle:                   true,
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

	enablePoolMsg := leveragelpmoduletypes.MsgAddPool{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Pool: leveragelpmoduletypes.AddPool{
			AmmPoolId:   poolId,
			LeverageMax: math.LegacyMustNewDecFromStr("10"),
		},
	}
	_, err = leveragelpmodulekeeper.NewMsgServerImpl(*app.LeveragelpKeeper).AddPool(ctx, &enablePoolMsg)
	require.NoError(t, err)

	// check length of pools
	require.Equal(t, len(pools), 1)
	tradingAssetPrice, err := app.PerpetualKeeper.GetAssetPrice(ctx, ptypes.ATOM)
	require.NoError(t, err)
	// call min collateral query
	res, err := mk.OpenEstimation(ctx, &types.QueryOpenEstimationRequest{
		PoolId:          1,
		Position:        types.Position_SHORT,
		Leverage:        math.LegacyMustNewDecFromStr("4.0"),
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(100_000_000)),
		Address:         "",
		TakeProfitPrice: tradingAssetPrice.QuoInt64(3),
	})
	require.NoError(t, err)
	expectedRes := &types.QueryOpenEstimationResponse{
		Position:           types.Position_SHORT,
		EffectiveLeverage:  math.LegacyMustNewDecFromStr("4.050637948604607222"),
		TradingAsset:       ptypes.ATOM,
		Collateral:         sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(100_000_000)),
		PositionSize:       sdk.NewCoin(ptypes.ATOM, math.NewInt(80200521)),
		OpenPrice:          math.LegacyMustNewDecFromStr("4.987498771984286735"),
		TakeProfitPrice:    tradingAssetPrice.QuoInt64(3),
		LiquidationPrice:   math.LegacyMustNewDecFromStr("6.082315575590593580"),
		EstimatedPnl:       sdk.Coin{ptypes.BaseCurrency, math.NewInt(266332466)},
		HourlyInterestRate: math.LegacyZeroDec(),
		AvailableLiquidity: sdk.NewCoin(ptypes.ATOM, math.NewInt(10000000000)),
		Slippage:           math.LegacyMustNewDecFromStr("0.001504012500000000"),
		BorrowInterestRate: math.LegacyMustNewDecFromStr("0.000000000000000000"),
		FundingRate:        math.LegacyMustNewDecFromStr("0.000000000000000000"),
		PriceImpact:        math.LegacyMustNewDecFromStr("0.002500245603142653"),
		Custody:            sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(500000000)),
		Liabilities:        sdk.NewCoin(ptypes.ATOM, math.NewInt(80_200_521)),
		WeightBreakingFee:  math.LegacyZeroDec(),
	}
	require.Equal(t, expectedRes, res)
}

func TestOpenEstimation_WrongAsset(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)
	simapp.SetStakingParam(app, ctx)
	simapp.SetPerpetualParams(app, ctx)
	simapp.SetupAssetProfile(app, ctx)

	mk, amm, oracle := app.PerpetualKeeper, app.AmmKeeper, app.OracleKeeper

	// Setup coin prices
	SetupStableCoinPrices(ctx, oracle)

	tradingAssetPrice, err := app.PerpetualKeeper.GetAssetPrice(ctx, ptypes.ATOM)
	require.NoError(t, err)

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

	err = app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, usdcToken)
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

	enablePoolMsg := leveragelpmoduletypes.MsgAddPool{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Pool: leveragelpmoduletypes.AddPool{
			AmmPoolId:   poolId,
			LeverageMax: math.LegacyMustNewDecFromStr("10"),
		},
	}
	_, err = leveragelpmodulekeeper.NewMsgServerImpl(*app.LeveragelpKeeper).AddPool(ctx, &enablePoolMsg)
	require.NoError(t, err)

	// check length of pools
	require.Equal(t, len(pools), 1)

	_, err = mk.OpenEstimation(ctx, &types.QueryOpenEstimationRequest{
		PoolId:          1,
		Position:        types.Position_LONG,
		Leverage:        math.LegacyMustNewDecFromStr("5.0"),
		TradingAsset:    ptypes.BaseCurrency,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(10000000)),
		Address:         "",
		TakeProfitPrice: tradingAssetPrice.MulInt64(3),
	})

	assert.Error(t, err)
	assert.Equal(t, "invalid operation: the borrowed asset cannot be the base currency: invalid borrowing asset", err.Error())

	_, err = mk.OpenEstimation(ctx, &types.QueryOpenEstimationRequest{
		PoolId:          1,
		Position:        types.Position_LONG,
		Leverage:        math.LegacyMustNewDecFromStr("5.0"),
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.Eden, math.NewInt(10000000)),
		Address:         "",
		TakeProfitPrice: tradingAssetPrice.MulInt64(3),
	})

	assert.Error(t, err)
	assert.Equal(t, "invalid collateral: collateral must either match the borrowed asset or be the base currency: invalid borrowing asset", err.Error())

	_, err = mk.OpenEstimation(ctx, &types.QueryOpenEstimationRequest{
		PoolId:          1,
		Position:        types.Position_SHORT,
		Leverage:        math.LegacyMustNewDecFromStr("5.0"),
		TradingAsset:    ptypes.BaseCurrency,
		Collateral:      sdk.NewCoin(ptypes.ATOM, math.NewInt(10000000)),
		Address:         "",
		TakeProfitPrice: tradingAssetPrice.QuoInt64(3),
	})

	assert.Error(t, err)
	assert.Equal(t, "borrowing not allowed: cannot take a short position against the base currency: invalid borrowing asset", err.Error())

	_, err = mk.OpenEstimation(ctx, &types.QueryOpenEstimationRequest{
		PoolId:          1,
		Position:        types.Position_SHORT,
		Leverage:        math.LegacyMustNewDecFromStr("5.0"),
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.ATOM, math.NewInt(10000000)),
		Address:         "",
		TakeProfitPrice: tradingAssetPrice.QuoInt64(3),
	})

	assert.Error(t, err)
	assert.Equal(t, "invalid operation: collateral asset cannot be identical to the borrowed asset for a short position: invalid collateral asset", err.Error())
}
