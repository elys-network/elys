package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	"github.com/elys-network/elys/x/perpetual/types"

	"github.com/cometbft/cometbft/crypto/ed25519"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	simapp "github.com/elys-network/elys/app"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenEstimation_Long5XAtom100Usdc(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	mk, amm, oracle := app.PerpetualKeeper, app.AmmKeeper, app.OracleKeeper

	// Setup coin prices
	SetupStableCoinPrices(ctx, oracle)

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(1000000000000))

	// Create a pool
	// Mint 100000USDC
	usdcToken := []sdk.Coin{sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(1000000000000))}
	// Mint 100000ATOM
	atomToken := []sdk.Coin{sdk.NewCoin(ptypes.ATOM, sdk.NewInt(1000000000000))}

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
			Weight: sdk.NewInt(50),
			Token:  sdk.NewCoin(ptypes.ATOM, sdk.NewInt(600000000000)),
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
	require.NoError(t, err)
	require.Equal(t, poolId, uint64(1))

	pools := amm.GetAllPool(ctx)

	// check length of pools
	require.Equal(t, len(pools), 1)

	// call min collateral query
	res, err := mk.OpenEstimation(ctx, &types.QueryOpenEstimationRequest{
		Position:        types.Position_LONG,
		Leverage:        sdk.MustNewDecFromStr("5.0"),
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(100000000)),
		Discount:        sdk.MustNewDecFromStr("0.0"),
		TakeProfitPrice: sdk.MustNewDecFromStr("20.0"),
	})
	require.NoError(t, err)
	require.Equal(t, &types.QueryOpenEstimationResponse{
		Position:           types.Position_LONG,
		Leverage:           sdk.MustNewDecFromStr("5.0"),
		TradingAsset:       ptypes.ATOM,
		Collateral:         sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(100000000)),
		Custody:            sdk.NewCoin(ptypes.ATOM, sdk.NewInt(499136435)),
		PositionSize:       sdk.NewCoin(ptypes.ATOM, sdk.NewInt(499136435)),
		SwapFee:            sdk.MustNewDecFromStr("0.001000000000000000"),
		Discount:           sdk.MustNewDecFromStr("0.000000000000000000"),
		OpenPrice:          sdk.MustNewDecFromStr("1.000000000000000000"),
		TakeProfitPrice:    sdk.MustNewDecFromStr("20.00000000000000000"),
		LiquidationPrice:   sdk.MustNewDecFromStr("0.799653976372211738"),
		InterestAmount:     types.NewParams().MinBorrowInterestAmount,
		EstimatedPnl:       sdk.NewInt(9_482_728700),
		EstimatedPnlDenom:  ptypes.BaseCurrency,
		AvailableLiquidity: sdk.NewCoin(ptypes.ATOM, sdk.NewInt(600000000000)),
		Slippage:           sdk.MustNewDecFromStr("0.000727856000000000"),
		WeightBalanceRatio: sdk.MustNewDecFromStr("0.000000000000000000"),
		BorrowInterestRate: sdk.MustNewDecFromStr("0.000000000000000000"),
		FundingRate:        sdk.MustNewDecFromStr("0.000000000000000000"),
		PriceImpact:        sdk.MustNewDecFromStr("0.000000000000000000"),
		BorrowFee:          sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(0)),
		FundingFee:         sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(0)),
	}, res)
}

func TestOpenEstimation_Long5XAtom10Atom(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

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
	addr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(1000000000000))

	// Create a pool
	// Mint 10000USDC
	usdcToken := []sdk.Coin{sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(1000000000000))}
	// Mint 100000ATOM
	atomToken := []sdk.Coin{sdk.NewCoin(ptypes.ATOM, sdk.NewInt(1000000000000))}

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
	require.NoError(t, err)
	require.Equal(t, poolId, uint64(1))

	pools := amm.GetAllPool(ctx)

	// check length of pools
	require.Equal(t, len(pools), 1)

	// call min collateral query
	res, err := mk.OpenEstimation(ctx, &types.QueryOpenEstimationRequest{
		Position:        types.Position_LONG,
		Leverage:        sdk.MustNewDecFromStr("5.0"),
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.ATOM, sdk.NewInt(10000000)),
		Discount:        sdk.MustNewDecFromStr("0.0"),
		TakeProfitPrice: sdk.MustNewDecFromStr("20.0"),
	})
	require.NoError(t, err)
	require.Equal(t, &types.QueryOpenEstimationResponse{
		Position:           types.Position_LONG,
		Leverage:           sdk.MustNewDecFromStr("5.0"),
		TradingAsset:       ptypes.ATOM,
		Collateral:         sdk.NewCoin(ptypes.ATOM, sdk.NewInt(10000000)),
		Custody:            sdk.NewCoin(ptypes.ATOM, sdk.NewInt(49858958)),
		PositionSize:       sdk.NewCoin(ptypes.ATOM, sdk.NewInt(49858958)),
		SwapFee:            sdk.MustNewDecFromStr("0.001000000000000000"),
		Discount:           sdk.MustNewDecFromStr("0.000000000000000000"),
		OpenPrice:          sdk.MustNewDecFromStr("1.000000000000000000"),
		TakeProfitPrice:    sdk.MustNewDecFromStr("20.00000000000000000"),
		LiquidationPrice:   sdk.MustNewDecFromStr("0.799662359570370484"),
		EstimatedPnl:       sdk.NewInt(757_224656),
		EstimatedPnlDenom:  ptypes.BaseCurrency,
		InterestAmount:     types.NewParams().MinBorrowInterestAmount,
		AvailableLiquidity: sdk.NewCoin(ptypes.ATOM, sdk.NewInt(10000000000)),
		Slippage:           sdk.MustNewDecFromStr("0.000686040302239768"),
		WeightBalanceRatio: sdk.MustNewDecFromStr("0.000000000000000000"),
		BorrowInterestRate: sdk.MustNewDecFromStr("0.000000000000000000"),
		FundingRate:        sdk.MustNewDecFromStr("0.000000000000000000"),
		PriceImpact:        sdk.MustNewDecFromStr("0.001685356924966457"),
		BorrowFee:          sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(0)),
		FundingFee:         sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(0)),
	}, res)
}

func TestOpenEstimation_Long10XAtom1000Usdc(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	mk, amm, oracle := app.PerpetualKeeper, app.AmmKeeper, app.OracleKeeper

	// Setup coin prices
	SetupStableCoinPrices(ctx, oracle)

	provider := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "ATOM",
		Price:     sdk.MustNewDecFromStr("4.39"),
		Source:    "atom",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "uatom",
		Price:     sdk.MustNewDecFromStr("4.39"),
		Source:    "uatom",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(1_000_000_000000))

	// Create a pool
	// Mint 100000USDC
	usdcToken := []sdk.Coin{sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(1_000_000_000000))}
	// Mint 100000ATOM
	atomToken := []sdk.Coin{sdk.NewCoin(ptypes.ATOM, sdk.NewInt(1_000_000_000000))}

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
			Weight: sdk.NewInt(50),
			Token:  sdk.NewCoin(ptypes.ATOM, sdk.NewInt(600_000_000000)),
		},
		{
			Weight: sdk.NewInt(50),
			Token:  sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(100_000_000000)),
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
	require.NoError(t, err)
	require.Equal(t, poolId, uint64(1))

	pools := amm.GetAllPool(ctx)

	// check length of pools
	require.Equal(t, len(pools), 1)

	// call min collateral query
	res, err := mk.OpenEstimation(ctx, &types.QueryOpenEstimationRequest{
		Position:        types.Position_LONG,
		Leverage:        sdk.MustNewDecFromStr("10.0"),
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(1_000_000000)),
		Discount:        sdk.MustNewDecFromStr("0.0"),
		TakeProfitPrice: sdk.MustNewDecFromStr("5.0"),
	})
	require.NoError(t, err)
	require.Equal(t, &types.QueryOpenEstimationResponse{
		Position:           types.Position_LONG,
		Leverage:           sdk.MustNewDecFromStr("10.0"),
		TradingAsset:       ptypes.ATOM,
		Collateral:         sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(1_000_000000)),
		Custody:            sdk.NewCoin(ptypes.ATOM, sdk.NewInt(2_247_067372)),
		PositionSize:       sdk.NewCoin(ptypes.ATOM, sdk.NewInt(2_247_067372)),
		SwapFee:            sdk.MustNewDecFromStr("0.001000000000000000"),
		Discount:           sdk.MustNewDecFromStr("0.000000000000000000"),
		OpenPrice:          sdk.MustNewDecFromStr("4.390000000000000009"),
		TakeProfitPrice:    sdk.MustNewDecFromStr("5.000000000000000000"),
		LiquidationPrice:   sdk.MustNewDecFromStr("3.944975514993148154"),
		InterestAmount:     types.NewParams().MinBorrowInterestAmount,
		EstimatedPnl:       sdk.NewInt(1_235_336860),
		EstimatedPnlDenom:  ptypes.BaseCurrency,
		AvailableLiquidity: sdk.NewCoin(ptypes.ATOM, sdk.NewInt(600_000_000000)),
		Slippage:           sdk.MustNewDecFromStr("0.012549973525000000"),
		WeightBalanceRatio: sdk.MustNewDecFromStr("0.000000000000000000"),
		BorrowInterestRate: sdk.MustNewDecFromStr("0.000000000000000000"),
		FundingRate:        sdk.MustNewDecFromStr("0.000000000000000000"),
		PriceImpact:        sdk.MustNewDecFromStr("0.000000000000000000"),
		BorrowFee:          sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(0)),
		FundingFee:         sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(0)),
	}, res)
}

func TestOpenEstimation_Short5XAtom10Usdc(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

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
		Price:     sdk.NewDec(5),
		Source:    "atom",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "uatom",
		Price:     sdk.NewDec(5),
		Source:    "uatom",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(1000000000000))

	// Create a pool
	// Mint 10000USDC
	usdcToken := []sdk.Coin{sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(1000000000000))}
	// Mint 100000ATOM
	atomToken := []sdk.Coin{sdk.NewCoin(ptypes.ATOM, sdk.NewInt(1000000000000))}

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
	require.NoError(t, err)
	require.Equal(t, poolId, uint64(1))

	pools := amm.GetAllPool(ctx)

	// check length of pools
	require.Equal(t, len(pools), 1)

	// call min collateral query
	res, err := mk.OpenEstimation(ctx, &types.QueryOpenEstimationRequest{
		Position:        types.Position_SHORT,
		Leverage:        sdk.MustNewDecFromStr("5.0"),
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(100000000)),
		Discount:        sdk.MustNewDecFromStr("0.0"),
		TakeProfitPrice: sdk.MustNewDecFromStr("2.0"),
	})
	require.NoError(t, err)
	expectedRes := &types.QueryOpenEstimationResponse{
		Position:           types.Position_SHORT,
		Leverage:           sdk.MustNewDecFromStr("5.0"),
		TradingAsset:       ptypes.ATOM,
		Collateral:         sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(100000000)),
		Custody:            sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(500000000)),
		PositionSize:       sdk.NewCoin(ptypes.ATOM, sdk.NewInt(80000000)),
		SwapFee:            sdk.MustNewDecFromStr("0.001000000000000000"),
		Discount:           sdk.MustNewDecFromStr("0.000000000000000000"),
		OpenPrice:          sdk.MustNewDecFromStr("5.000000000000000000"),
		TakeProfitPrice:    sdk.MustNewDecFromStr("2.000000000000000000"),
		LiquidationPrice:   sdk.MustNewDecFromStr("6.000000000000000000"),
		EstimatedPnl:       sdk.NewInt(-400000000),
		EstimatedPnlDenom:  ptypes.BaseCurrency,
		InterestAmount:     types.NewParams().MinBorrowInterestAmount,
		AvailableLiquidity: sdk.NewCoin(ptypes.ATOM, sdk.NewInt(10000000000)),
		Slippage:           sdk.MustNewDecFromStr("0.001868770000000000"),
		WeightBalanceRatio: sdk.MustNewDecFromStr("0.000000000000000000"),
		BorrowInterestRate: sdk.MustNewDecFromStr("0.000000000000000000"),
		FundingRate:        sdk.MustNewDecFromStr("0.000000000000000000"),
		PriceImpact:        sdk.MustNewDecFromStr("0.002866910000000000"),
		BorrowFee:          sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(0)),
		FundingFee:         sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(0)),
	}
	// print both the expectedRes and res
	t.Logf("%+v", expectedRes)
	t.Logf("%+v", res)
	require.Equal(t, expectedRes, res)
}

func TestOpenEstimation_WrongAsset(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

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
	addr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(1000000000000))

	// Create a pool
	// Mint 10000USDC
	usdcToken := []sdk.Coin{sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(1000000000000))}
	// Mint ATOM
	atomToken := []sdk.Coin{sdk.NewCoin(ptypes.ATOM, sdk.NewInt(1000000000000))}

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
	require.NoError(t, err)
	require.Equal(t, poolId, uint64(1))

	pools := amm.GetAllPool(ctx)

	// check length of pools
	require.Equal(t, len(pools), 1)

	_, err = mk.OpenEstimation(ctx, &types.QueryOpenEstimationRequest{
		Position:        types.Position_LONG,
		Leverage:        sdk.MustNewDecFromStr("5.0"),
		TradingAsset:    ptypes.BaseCurrency,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(10000000)),
		Discount:        sdk.MustNewDecFromStr("0.0"),
		TakeProfitPrice: sdk.MustNewDecFromStr("20.0"),
	})

	assert.Error(t, err)
	assert.Equal(t, "invalid operation: the borrowed asset cannot be the base currency: invalid borrowing asset", err.Error())

	_, err = mk.OpenEstimation(ctx, &types.QueryOpenEstimationRequest{
		Position:        types.Position_LONG,
		Leverage:        sdk.MustNewDecFromStr("5.0"),
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.Eden, sdk.NewInt(10000000)),
		Discount:        sdk.MustNewDecFromStr("0.0"),
		TakeProfitPrice: sdk.MustNewDecFromStr("20.0"),
	})

	assert.Error(t, err)
	assert.Equal(t, "invalid collateral: collateral must either match the borrowed asset or be the base currency: invalid borrowing asset", err.Error())

	_, err = mk.OpenEstimation(ctx, &types.QueryOpenEstimationRequest{
		Position:        types.Position_SHORT,
		Leverage:        sdk.MustNewDecFromStr("5.0"),
		TradingAsset:    ptypes.BaseCurrency,
		Collateral:      sdk.NewCoin(ptypes.ATOM, sdk.NewInt(10000000)),
		Discount:        sdk.MustNewDecFromStr("0.0"),
		TakeProfitPrice: sdk.MustNewDecFromStr("20.0"),
	})

	assert.Error(t, err)
	assert.Equal(t, "borrowing not allowed: cannot take a short position against the base currency: invalid borrowing asset", err.Error())

	_, err = mk.OpenEstimation(ctx, &types.QueryOpenEstimationRequest{
		Position:        types.Position_SHORT,
		Leverage:        sdk.MustNewDecFromStr("5.0"),
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.ATOM, sdk.NewInt(10000000)),
		Discount:        sdk.MustNewDecFromStr("0.0"),
		TakeProfitPrice: sdk.MustNewDecFromStr("20.0"),
	})

	assert.Error(t, err)
	assert.Equal(t, "invalid operation: collateral asset cannot be identical to the borrowed asset for a short position: invalid collateral asset", err.Error())
}
