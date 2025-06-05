package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	leveragelpmodulekeeper "github.com/elys-network/elys/v6/x/leveragelp/keeper"
	leveragelpmoduletypes "github.com/elys-network/elys/v6/x/leveragelp/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/v6/x/assetprofile/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
	oracletypes "github.com/ojo-network/ojo/x/oracle/types"

	"github.com/cometbft/cometbft/crypto/ed25519"
	simapp "github.com/elys-network/elys/v6/app"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
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
	tradingAssetPrice, _, err := app.PerpetualKeeper.GetAssetPriceAndAssetUsdcDenomRatio(ctx, ptypes.ATOM)
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

	poolParams := ammtypes.PoolParams{
		UseOracle: true,
		SwapFee:   argSwapFee,
		FeeDenom:  ptypes.BaseCurrency,
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
		EffectiveLeverage:  math.LegacyMustNewDecFromStr("5.025957360079110062"),
		TradingAsset:       ptypes.ATOM,
		Collateral:         sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(100_000_000)),
		HourlyInterestRate: math.LegacyZeroDec(),
		PositionSize:       sdk.NewCoin(ptypes.ATOM, math.NewInt(99871050)),
		OpenPrice:          math.LegacyMustNewDecFromStr("5.006455824786061626"),
		TakeProfitPrice:    tradingAssetPrice.MulInt64(3),
		LiquidationPrice:   math.LegacyMustNewDecFromStr("4.105293776324570533"),
		EstimatedPnl:       sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(998065750)),
		AvailableLiquidity: sdk.NewCoin(ptypes.ATOM, math.NewInt(600000000000)),
		Slippage:           math.LegacyMustNewDecFromStr("0.001289500000000000"),
		BorrowInterestRate: math.LegacyMustNewDecFromStr("0.000000000000000000"),
		FundingRate:        math.LegacyMustNewDecFromStr("0.000000000000000000"),
		PriceImpact:        math.LegacyMustNewDecFromStr("-0.001291164957212325"),
		Custody:            sdk.NewCoin(ptypes.ATOM, math.NewInt(99871050)),
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
	tradingAssetPrice, _, err := app.PerpetualKeeper.GetAssetPriceAndAssetUsdcDenomRatio(ctx, ptypes.ATOM)
	require.NoError(t, err)
	// Set asset profile
	app.AssetprofileKeeper.SetEntry(ctx, assetprofiletypes.Entry{
		BaseDenom:   ptypes.ATOM,
		Denom:       ptypes.ATOM,
		Decimals:    6,
		DisplayName: "ATOM",
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
	poolParams := ammtypes.PoolParams{
		UseOracle: true,
		SwapFee:   argSwapFee,
		FeeDenom:  ptypes.BaseCurrency,
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
		EffectiveLeverage:  math.LegacyMustNewDecFromStr("5.080191844300654694"),
		TradingAsset:       ptypes.ATOM,
		Collateral:         sdk.NewCoin(ptypes.ATOM, math.NewInt(10_000_000)),
		PositionSize:       sdk.NewCoin(ptypes.ATOM, math.NewInt(50_000_000)),
		OpenPrice:          math.LegacyMustNewDecFromStr("5.019731500000000000"),
		TakeProfitPrice:    tradingAssetPrice.MulInt64(3),
		LiquidationPrice:   math.LegacyMustNewDecFromStr("4.116179830000000000"),
		EstimatedPnl:       sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(399210740)),
		HourlyInterestRate: math.LegacyZeroDec(),
		AvailableLiquidity: sdk.NewCoin(ptypes.ATOM, math.NewInt(10000000000)),
		Slippage:           math.LegacyMustNewDecFromStr("0.001502510000000000"),
		BorrowInterestRate: math.LegacyMustNewDecFromStr("0.000000000000000000"),
		FundingRate:        math.LegacyMustNewDecFromStr("0.000000000000000000"),
		PriceImpact:        math.LegacyMustNewDecFromStr("-0.003946300000000000"),
		Custody:            sdk.NewCoin(ptypes.ATOM, math.NewInt(50_000_000)),
		Liabilities:        sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(200789260)),
		WeightBreakingFee:  math.LegacyMustNewDecFromStr("0.001435619047211833"),
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
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "uatom",
		Price:     math.LegacyMustNewDecFromStr("4.39"),
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
	app.AssetprofileKeeper.SetEntry(ctx, assetprofiletypes.Entry{
		BaseDenom:   ptypes.BaseCurrency,
		Denom:       ptypes.BaseCurrency,
		Decimals:    6,
		DisplayName: "USDC",
	})
	app.AssetprofileKeeper.SetEntry(ctx, assetprofiletypes.Entry{
		BaseDenom:   ptypes.ATOM,
		Denom:       ptypes.ATOM,
		Decimals:    6,
		DisplayName: "ATOM",
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

	poolParams := ammtypes.PoolParams{
		UseOracle: true,
		SwapFee:   argSwapFee,
		FeeDenom:  ptypes.BaseCurrency,
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
	tradingAssetPrice, _, err := app.PerpetualKeeper.GetAssetPriceAndAssetUsdcDenomRatio(ctx, ptypes.ATOM)
	require.NoError(t, err)
	// _,call min collateral query	tradingAssetPrice := app.OracleKeeper.GetAssetPriceAndAssetUsdcDenomRatioFromDenom(ctx, ptypes.ATOM)
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
		EffectiveLeverage:  math.LegacyMustNewDecFromStr("13.016012937612077877"),
		TradingAsset:       ptypes.ATOM,
		Collateral:         sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1_000_000_000)),
		PositionSize:       sdk.NewCoin(ptypes.ATOM, math.NewInt(2220729049)),
		OpenPrice:          math.LegacyMustNewDecFromStr("4.503025708833333678"),
		TakeProfitPrice:    tradingAssetPrice.MulInt64(3),
		LiquidationPrice:   math.LegacyMustNewDecFromStr("4.154041216398750318"),
		HourlyInterestRate: math.LegacyZeroDec(),
		EstimatedPnl:       sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(19247001575)),
		AvailableLiquidity: sdk.NewCoin(ptypes.ATOM, math.NewInt(600_000_000000)),
		Slippage:           math.LegacyMustNewDecFromStr("0.025099947050000000"),
		BorrowInterestRate: math.LegacyMustNewDecFromStr("0.000000000000000000"),
		FundingRate:        math.LegacyMustNewDecFromStr("0.000000000000000000"),
		PriceImpact:        math.LegacyMustNewDecFromStr("-0.025746175132877831"),
		Custody:            sdk.NewCoin(ptypes.ATOM, math.NewInt(2220729049)),
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
		BaseDenom:   ptypes.ATOM,
		Denom:       ptypes.ATOM,
		Decimals:    6,
		DisplayName: "ATOM",
	})
	provider := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "ATOM",
		Price:     math.LegacyNewDec(5),
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "uatom",
		Price:     math.LegacyNewDec(5),
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

	poolParams := ammtypes.PoolParams{
		UseOracle: true,
		SwapFee:   argSwapFee,
		FeeDenom:  ptypes.BaseCurrency,
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
	tradingAssetPrice, _, err := app.PerpetualKeeper.GetAssetPriceAndAssetUsdcDenomRatio(ctx, ptypes.ATOM)
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
		EffectiveLeverage:  math.LegacyMustNewDecFromStr("4.060893172743770318"),
		TradingAsset:       ptypes.ATOM,
		Collateral:         sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(100_000_000)),
		PositionSize:       sdk.NewCoin(ptypes.ATOM, math.NewInt(80240642)),
		OpenPrice:          math.LegacyMustNewDecFromStr("4.985004980393850786"),
		TakeProfitPrice:    tradingAssetPrice.QuoInt64(3),
		LiquidationPrice:   math.LegacyMustNewDecFromStr("6.079274366333964374"),
		EstimatedPnl:       sdk.Coin{ptypes.BaseCurrency, math.NewInt(266265597)},
		HourlyInterestRate: math.LegacyZeroDec(),
		AvailableLiquidity: sdk.NewCoin(ptypes.ATOM, math.NewInt(10000000000)),
		Slippage:           math.LegacyMustNewDecFromStr("0.003008025000000000"),
		BorrowInterestRate: math.LegacyMustNewDecFromStr("0.000000000000000000"),
		FundingRate:        math.LegacyMustNewDecFromStr("0.000000000000000000"),
		PriceImpact:        math.LegacyMustNewDecFromStr("0.002999003921229843"),
		Custody:            sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(500000000)),
		Liabilities:        sdk.NewCoin(ptypes.ATOM, math.NewInt(80240642)),
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

	tradingAssetPrice, _, err := app.PerpetualKeeper.GetAssetPriceAndAssetUsdcDenomRatio(ctx, ptypes.ATOM)
	require.NoError(t, err)

	// Set asset profile
	app.AssetprofileKeeper.SetEntry(ctx, assetprofiletypes.Entry{
		BaseDenom:   ptypes.ATOM,
		Denom:       ptypes.ATOM,
		Decimals:    6,
		DisplayName: "ATOM",
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

	poolParams := ammtypes.PoolParams{
		UseOracle: true,
		SwapFee:   argSwapFee,
		FeeDenom:  ptypes.BaseCurrency,
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
