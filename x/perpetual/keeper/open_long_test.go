package keeper_test

import (
	"errors"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/perpetual/keeper"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/elys-network/elys/x/perpetual/types/mocks"
	"github.com/stretchr/testify/assert"

	simapp "github.com/elys-network/elys/app"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestOpenLong_PoolNotFound(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenDefineAssetsChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		OpenDefineAssetsChecker: mockChecker,
	}

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			Leverage:     math.LegacyNewDec(10),
			TradingAsset: "bbb",
			Collateral:   sdk.NewCoin("aaa", math.NewInt(1)),
			Position:     types.Position_LONG,
		}
		poolId = uint64(42)
	)

	// Mock behavior
	mockChecker.On("GetMaxLeverageParam", ctx).Return(msg.Leverage).Once()
	mockChecker.On("GetPool", ctx, poolId).Return(types.Pool{}, false).Once()

	_, err := k.OpenDefineAssets(ctx, poolId, msg, ptypes.BaseCurrency, false)

	// Expect an error about the pool not existing
	assert.True(t, errors.Is(err, types.ErrPoolDoesNotExist))
	mockChecker.AssertExpectations(t)
}

func TestOpenLong_PoolDisabled(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenDefineAssetsChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		OpenDefineAssetsChecker: mockChecker,
	}

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			Leverage:     math.LegacyNewDec(10),
			TradingAsset: "bbb",
			Collateral:   sdk.NewCoin("aaa", math.NewInt(1)),
			Position:     types.Position_LONG,
		}
		poolId = uint64(42)
	)

	// Mock behaviors
	mockChecker.On("GetMaxLeverageParam", ctx).Return(msg.Leverage)
	mockChecker.On("GetPool", ctx, poolId).Return(types.Pool{}, true)
	mockChecker.On("IsPoolEnabled", ctx, poolId).Return(false)

	_, err := k.OpenDefineAssets(ctx, poolId, msg, ptypes.BaseCurrency, false)

	// Expect an error about the pool being disabled
	assert.True(t, errors.Is(err, types.ErrMTPDisabled))
	mockChecker.AssertExpectations(t)
}

func TestOpenLong_InsufficientAmmPoolBalanceForCustody(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenDefineAssetsChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		OpenDefineAssetsChecker: mockChecker,
	}

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			Creator:      "",
			Leverage:     math.LegacyNewDec(10),
			Position:     types.Position_LONG,
			TradingAsset: "uatom",
			Collateral:   sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000)),
		}
		ammPool = ammtypes.Pool{
			PoolId: uint64(42),
			PoolAssets: []ammtypes.PoolAsset{
				{
					Token:  sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(10000)),
					Weight: math.NewInt(50),
				},
				{
					Token:  sdk.NewCoin("uatom", math.NewInt(10)),
					Weight: math.NewInt(50),
				},
			},
		}
		pool = types.Pool{
			AmmPoolId: ammPool.PoolId,
		}
	)
	// Mock behaviors
	mockChecker.On("GetMaxLeverageParam", ctx).Return(msg.Leverage)
	mockChecker.On("GetPool", ctx, ammPool.PoolId).Return(pool, true)
	mockChecker.On("IsPoolEnabled", ctx, ammPool.PoolId).Return(true)
	mockChecker.On("GetAmmPool", ctx, ammPool.PoolId, msg.TradingAsset).Return(ammPool, nil)

	leveragedAmtTokenIn := sdk.NewCoin(msg.Collateral.Denom, math.NewInt(10000))
	custodyAmount := math.NewInt(99)

	mockChecker.On("EstimateSwap", ctx, leveragedAmtTokenIn, msg.TradingAsset, ammPool).Return(custodyAmount, nil)

	_, err := k.OpenDefineAssets(ctx, ammPool.PoolId, msg, ptypes.BaseCurrency, false)

	// Expect an error about custody amount being too high
	assert.True(t, errors.Is(err, types.ErrCustodyTooHigh))
	mockChecker.AssertExpectations(t)
}

func TestOpenLong_ErrorsDuringOperations(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenDefineAssetsChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		OpenDefineAssetsChecker: mockChecker,
	}

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			Creator:         "",
			Leverage:        math.LegacyNewDec(10),
			Position:        types.Position_LONG,
			TradingAsset:    "uatom",
			Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000)),
			TakeProfitPrice: math.LegacyMustNewDecFromStr(types.TakeProfitPriceDefault),
		}
		ammPool = ammtypes.Pool{
			PoolId: uint64(42),
			PoolAssets: []ammtypes.PoolAsset{
				{
					Token:  sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(10000)),
					Weight: math.NewInt(50),
				},
				{
					Token:  sdk.NewCoin("uatom", math.NewInt(10000)),
					Weight: math.NewInt(50),
				},
			},
		}
		pool = types.Pool{
			AmmPoolId: ammPool.PoolId,
		}
	)

	// Mock behaviors
	mockChecker.On("GetMaxLeverageParam", ctx).Return(msg.Leverage)
	mockChecker.On("GetPool", ctx, ammPool.PoolId).Return(pool, true)
	mockChecker.On("IsPoolEnabled", ctx, ammPool.PoolId).Return(true)
	mockChecker.On("GetAmmPool", ctx, ammPool.PoolId, msg.TradingAsset).Return(ammPool, nil)

	eta := math.LegacyNewDec(9)

	leveragedAmtTokenIn := sdk.NewCoin(msg.Collateral.Denom, math.NewInt(10000))
	custodyAmount := math.NewInt(99)

	mockChecker.On("EstimateSwap", ctx, leveragedAmtTokenIn, msg.TradingAsset, ammPool).Return(custodyAmount, nil)

	mtp := types.NewMTP(msg.Creator, msg.Collateral.Denom, msg.TradingAsset, msg.Collateral.Denom, msg.TradingAsset, msg.Position, msg.Leverage, math.LegacyMustNewDecFromStr(types.TakeProfitPriceDefault), ammPool.PoolId)

	borrowError := errors.New("borrow error")
	mockChecker.On("Borrow", ctx, msg.Collateral.Amount, custodyAmount, mtp, &ammPool, &pool, eta, ptypes.BaseCurrency, false).Return(borrowError)

	_, err := k.OpenDefineAssets(ctx, ammPool.PoolId, msg, ptypes.BaseCurrency, false)

	// Expect the borrow error
	assert.True(t, errors.Is(err, borrowError))
	mockChecker.AssertExpectations(t)
}

func TestOpenLong_LeverageRatioLessThanSafetyFactor(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenDefineAssetsChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		OpenDefineAssetsChecker: mockChecker,
	}

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			Creator:         "",
			Leverage:        math.LegacyNewDec(10),
			Position:        types.Position_LONG,
			TradingAsset:    "uatom",
			Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000)),
			TakeProfitPrice: math.LegacyMustNewDecFromStr(types.TakeProfitPriceDefault),
		}
		ammPool = ammtypes.Pool{
			PoolId: uint64(42),
			PoolAssets: []ammtypes.PoolAsset{
				{
					Token:  sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(10000)),
					Weight: math.NewInt(50),
				},
				{
					Token:  sdk.NewCoin("uatom", math.NewInt(10000)),
					Weight: math.NewInt(50),
				},
			},
		}
		pool = types.Pool{
			AmmPoolId: ammPool.PoolId,
		}
	)

	// Mock behaviors
	mockChecker.On("GetMaxLeverageParam", ctx).Return(msg.Leverage)
	mockChecker.On("GetPool", ctx, ammPool.PoolId).Return(pool, true)
	mockChecker.On("IsPoolEnabled", ctx, ammPool.PoolId).Return(true)
	mockChecker.On("GetAmmPool", ctx, ammPool.PoolId, msg.TradingAsset).Return(ammPool, nil)

	eta := math.LegacyNewDec(9)

	leveragedAmtTokenIn := sdk.NewCoin(msg.Collateral.Denom, math.NewInt(10000))
	custodyAmount := math.NewInt(99)

	mockChecker.On("EstimateSwap", ctx, leveragedAmtTokenIn, msg.TradingAsset, ammPool).Return(custodyAmount, nil)

	mtp := types.NewMTP(msg.Creator, msg.Collateral.Denom, msg.TradingAsset, msg.Collateral.Denom, msg.TradingAsset, msg.Position, msg.Leverage, math.LegacyMustNewDecFromStr(types.TakeProfitPriceDefault), ammPool.PoolId)

	mockChecker.On("Borrow", ctx, msg.Collateral.Amount, custodyAmount, mtp, &ammPool, &pool, eta, ptypes.BaseCurrency, false).Return(nil)
	mockChecker.On("UpdatePoolHealth", ctx, &pool).Return(nil)
	mockChecker.On("TakeInCustody", ctx, *mtp, &pool).Return(nil)

	lr := math.LegacyNewDec(50)

	mockChecker.On("GetMTPHealth", ctx, *mtp, ammPool, ptypes.BaseCurrency).Return(lr, nil)
	mockChecker.On("GetSafetyFactor", ctx).Return(math.LegacyNewDec(100))

	_, err := k.OpenDefineAssets(ctx, ammPool.PoolId, msg, ptypes.BaseCurrency, false)

	// Expect an error indicating MTP is unhealthy
	assert.True(t, errors.Is(err, types.ErrMTPUnhealthy))
	mockChecker.AssertExpectations(t)
}

func TestOpenLong_Success(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenDefineAssetsChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		OpenDefineAssetsChecker: mockChecker,
	}

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			Creator:         "",
			Leverage:        math.LegacyNewDec(10),
			Position:        types.Position_LONG,
			TradingAsset:    "uatom",
			Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000)),
			TakeProfitPrice: math.LegacyMustNewDecFromStr(types.TakeProfitPriceDefault),
		}
		ammPool = ammtypes.Pool{
			PoolId: uint64(42),
			PoolAssets: []ammtypes.PoolAsset{
				{
					Token:  sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(10000)),
					Weight: math.NewInt(50),
				},
				{
					Token:  sdk.NewCoin("uatom", math.NewInt(10000)),
					Weight: math.NewInt(50),
				},
			},
		}
		pool = types.Pool{
			AmmPoolId: ammPool.PoolId,
		}
	)

	// Mock behaviors
	mockChecker.On("GetMaxLeverageParam", ctx).Return(msg.Leverage)
	mockChecker.On("GetPool", ctx, ammPool.PoolId).Return(pool, true)
	mockChecker.On("IsPoolEnabled", ctx, ammPool.PoolId).Return(true)
	mockChecker.On("GetAmmPool", ctx, ammPool.PoolId, msg.TradingAsset).Return(ammPool, nil)

	eta := math.LegacyNewDec(9)

	leveragedAmtTokenIn := sdk.NewCoin(msg.Collateral.Denom, math.NewInt(10000))
	custodyAmount := math.NewInt(99)

	mockChecker.On("EstimateSwap", ctx, leveragedAmtTokenIn, msg.TradingAsset, ammPool).Return(custodyAmount, nil)

	mtp := types.NewMTP(msg.Creator, msg.Collateral.Denom, msg.TradingAsset, msg.Collateral.Denom, msg.TradingAsset, msg.Position, msg.Leverage, math.LegacyMustNewDecFromStr(types.TakeProfitPriceDefault), ammPool.PoolId)

	mockChecker.On("Borrow", ctx, msg.Collateral.Amount, custodyAmount, mtp, &ammPool, &pool, eta, ptypes.BaseCurrency, false).Return(nil)
	mockChecker.On("UpdatePoolHealth", ctx, &pool).Return(nil)
	mockChecker.On("TakeInCustody", ctx, *mtp, &pool).Return(nil)

	lr := math.LegacyNewDec(50)

	mockChecker.On("GetMTPHealth", ctx, *mtp, ammPool, ptypes.BaseCurrency).Return(lr, nil)

	safetyFactor := math.LegacyNewDec(10)

	mockChecker.On("GetSafetyFactor", ctx).Return(safetyFactor)

	mockChecker.On("CalcMTPConsolidateCollateral", ctx, mtp, ptypes.BaseCurrency).Return(nil)
	mockChecker.On("SetMTP", ctx, mtp).Return(nil)

	_, err := k.OpenDefineAssets(ctx, ammPool.PoolId, msg, ptypes.BaseCurrency, false)
	// Expect no error
	assert.Nil(t, err)
	mockChecker.AssertExpectations(t)
}

func TestOpenLong_BaseCurrency_Collateral(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)
	//simapp.SetStableStake(app, ctx)
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
	addr := simapp.AddTestAddrs(app, ctx, 1, math.NewInt(1000000000000))

	// Create a pool
	// Mint 100000USDC
	usdcToken := []sdk.Coin{sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(200000000000))}
	// Mint 100000ATOM
	atomToken := []sdk.Coin{sdk.NewCoin(ptypes.ATOM, math.NewInt(200000000000))}

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

	// check block height
	require.Equal(t, int64(0), ctx.BlockHeight())

	pool, found := amm.GetPool(ctx, poolId)
	require.Equal(t, found, true)

	poolAddress := sdk.MustAccAddressFromBech32(pool.GetAddress())
	require.NoError(t, err)

	// Balance check before create a perpetual position
	balances := app.BankKeeper.GetAllBalances(ctx, poolAddress)
	require.Equal(t, balances.AmountOf(ptypes.BaseCurrency), math.NewInt(100000000000))
	require.Equal(t, balances.AmountOf(ptypes.ATOM), math.NewInt(10000000000))

	// Create a perpetual position open msg
	msg2 := types.NewMsgOpen(
		addr[0].String(),
		types.Position_LONG,
		math.LegacyNewDec(5),
		ptypes.ATOM,
		sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(100000000)),
		math.LegacyMustNewDecFromStr(types.TakeProfitPriceDefault),
		math.LegacyNewDec(100),
	)

	_, err = mk.Open(ctx, msg2, false)
	require.NoError(t, err)

	mtps := mk.GetAllMTPs(ctx)
	require.Equal(t, len(mtps), 1)

	balances = app.BankKeeper.GetAllBalances(ctx, poolAddress)
	require.Equal(t, balances.AmountOf(ptypes.BaseCurrency), math.NewInt(100100000000))
	require.Equal(t, balances.AmountOf(ptypes.ATOM), math.NewInt(10000000000))

	_, found = mk.OpenDefineAssetsChecker.GetPool(ctx, pool.PoolId)
	require.Equal(t, found, true)

	err = mk.InvariantCheck(ctx)
	require.Equal(t, err, nil)

	mtp := mtps[0]

	// Check MTP
	require.Equal(t, types.MTP{
		Address:                        addr[0].String(),
		CollateralAsset:                "uusdc",
		TradingAsset:                   "uatom",
		LiabilitiesAsset:               "uusdc",
		CustodyAsset:                   "uatom",
		Collateral:                     math.NewInt(100000000),
		Liabilities:                    math.NewInt(400000000),
		BorrowInterestPaidCollateral:   math.NewInt(0),
		BorrowInterestPaidCustody:      math.NewInt(0),
		BorrowInterestUnpaidCollateral: math.NewInt(0),
		Custody:                        math.NewInt(486520593),
		TakeProfitLiabilities:          math.NewInt(473929244),
		TakeProfitCustody:              math.NewInt(486520593),
		MtpHealth:                      math.LegacyMustNewDecFromStr("1.234567958024691358"),
		Position:                       types.Position_LONG,
		Id:                             uint64(1),
		AmmPoolId:                      uint64(1),
		TakeProfitPrice:                math.LegacyMustNewDecFromStr(types.TakeProfitPriceDefault),
		TakeProfitBorrowRate:           math.LegacyMustNewDecFromStr("1.0"),
		FundingFeePaidCollateral:       math.NewInt(0),
		FundingFeePaidCustody:          math.NewInt(0),
		FundingFeeReceivedCollateral:   math.NewInt(0),
		FundingFeeReceivedCustody:      math.NewInt(0),
		OpenPrice:                      math.LegacyMustNewDecFromStr("1.027705727555914576"),
		StopLossPrice:                  math.LegacyNewDec(100),
	}, mtp)
}

func TestOpenLong_ATOM_Collateral(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

	mk, amm, oracle := app.PerpetualKeeper, app.AmmKeeper, app.OracleKeeper

	// Setup coin prices
	SetupStableCoinPrices(ctx, oracle)

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, math.NewInt(1000000000000))

	// Create a pool
	// Mint 100000USDC
	usdcToken := []sdk.Coin{sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(100000000000))}
	// Mint 100000ATOM
	atomToken := []sdk.Coin{sdk.NewCoin(ptypes.ATOM, math.NewInt(100000000000))}

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
			Token:  sdk.NewCoin(ptypes.ATOM, math.NewInt(1000000000)),
		},
		{
			Weight: math.NewInt(50),
			Token:  sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(10000000000)),
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

	// check block height
	require.Equal(t, int64(0), ctx.BlockHeight())

	pool, found := amm.GetPool(ctx, poolId)
	require.Equal(t, found, true)

	poolAddress := sdk.MustAccAddressFromBech32(pool.GetAddress())
	require.NoError(t, err)

	// Balance check before create a perpetual position
	balances := app.BankKeeper.GetAllBalances(ctx, poolAddress)
	require.Equal(t, balances.AmountOf(ptypes.BaseCurrency), math.NewInt(10000000000))
	require.Equal(t, balances.AmountOf(ptypes.ATOM), math.NewInt(1000000000))

	// Create a perpetual position open msg
	msg2 := types.NewMsgOpen(
		addr[0].String(),
		types.Position_LONG,
		math.LegacyNewDec(5),
		ptypes.ATOM,
		sdk.NewCoin(ptypes.ATOM, math.NewInt(10000000)),
		math.LegacyMustNewDecFromStr(types.TakeProfitPriceDefault),
		math.LegacyNewDec(100),
	)

	_, err = mk.Open(ctx, msg2, false)
	require.NoError(t, err)

	mtps := mk.GetAllMTPs(ctx)
	require.Equal(t, len(mtps), 1)

	balances = app.BankKeeper.GetAllBalances(ctx, poolAddress)
	require.Equal(t, balances.AmountOf(ptypes.BaseCurrency), math.NewInt(10000000000))
	require.Equal(t, balances.AmountOf(ptypes.ATOM), math.NewInt(1010000000))

	_, found = mk.OpenDefineAssetsChecker.GetPool(ctx, pool.PoolId)
	require.Equal(t, found, true)

	err = mk.InvariantCheck(ctx)
	require.Equal(t, err, nil)

	mtp := mtps[0]

	// Check MTP
	require.Equal(t, types.MTP{
		Address:                        addr[0].String(),
		CollateralAsset:                "uatom",
		TradingAsset:                   "uatom",
		LiabilitiesAsset:               "uusdc",
		CustodyAsset:                   "uatom",
		Collateral:                     math.NewInt(10000000),
		Liabilities:                    math.NewInt(40905431),
		BorrowInterestPaidCollateral:   math.NewInt(0),
		BorrowInterestPaidCustody:      math.NewInt(0),
		BorrowInterestUnpaidCollateral: math.NewInt(0),
		Custody:                        math.NewInt(50000000),
		TakeProfitLiabilities:          math.NewInt(48671331),
		TakeProfitCustody:              math.NewInt(50000000),
		MtpHealth:                      math.LegacyMustNewDecFromStr("1.119902772759909084"),
		Position:                       types.Position_LONG,
		Id:                             uint64(1),
		AmmPoolId:                      uint64(1),
		TakeProfitPrice:                math.LegacyMustNewDecFromStr(types.TakeProfitPriceDefault),
		TakeProfitBorrowRate:           math.LegacyMustNewDecFromStr("1.0"),
		FundingFeePaidCollateral:       math.NewInt(0),
		FundingFeePaidCustody:          math.NewInt(0),
		FundingFeeReceivedCollateral:   math.NewInt(0),
		FundingFeeReceivedCustody:      math.NewInt(0),
		OpenPrice:                      math.LegacyMustNewDecFromStr("1.017016260000000000"),
		StopLossPrice:                  math.LegacyNewDec(100),
	}, mtp)
}

func TestOpenLongConsolidate_Success(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenDefineAssetsChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		OpenDefineAssetsChecker: mockChecker,
	}

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			Creator:         "",
			Leverage:        math.LegacyNewDec(10),
			Position:        types.Position_LONG,
			TradingAsset:    "uatom",
			Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000)),
			TakeProfitPrice: math.LegacyMustNewDecFromStr(types.TakeProfitPriceDefault),
		}
		ammPool = ammtypes.Pool{
			PoolId: uint64(42),
			PoolAssets: []ammtypes.PoolAsset{
				{
					Token:  sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(10000000)),
					Weight: math.NewInt(50),
				},
				{
					Token:  sdk.NewCoin("uatom", math.NewInt(10000000)),
					Weight: math.NewInt(50),
				},
			},
		}
		pool = types.Pool{
			AmmPoolId: ammPool.PoolId,
		}
	)

	// Mock behaviors
	mockChecker.On("GetMaxLeverageParam", ctx).Return(math.LegacyNewDec(50))
	mockChecker.On("GetPool", ctx, ammPool.PoolId).Return(pool, true)
	mockChecker.On("IsPoolEnabled", ctx, ammPool.PoolId).Return(true)
	mockChecker.On("GetAmmPool", ctx, ammPool.PoolId, msg.TradingAsset).Return(ammPool, nil)

	eta := math.LegacyNewDec(9)

	leveragedAmtTokenIn := sdk.NewCoin(msg.Collateral.Denom, math.NewInt(10000))
	custodyAmount := math.NewInt(99)

	mockChecker.On("EstimateSwap", ctx, leveragedAmtTokenIn, msg.TradingAsset, ammPool).Return(custodyAmount, nil)

	mtp := types.NewMTP(msg.Creator, msg.Collateral.Denom, msg.TradingAsset, msg.Collateral.Denom, msg.TradingAsset, msg.Position, msg.Leverage, math.LegacyMustNewDecFromStr(types.TakeProfitPriceDefault), ammPool.PoolId)

	mockChecker.On("Borrow", ctx, msg.Collateral.Amount, custodyAmount, mtp, &ammPool, &pool, eta, ptypes.BaseCurrency, false).Return(nil)
	mockChecker.On("UpdatePoolHealth", ctx, &pool).Return(nil)
	mockChecker.On("TakeInCustody", ctx, *mtp, &pool).Return(nil)

	lr := math.LegacyNewDec(50)

	mockChecker.On("GetMTPHealth", ctx, *mtp, ammPool, ptypes.BaseCurrency).Return(lr, nil)

	safetyFactor := math.LegacyNewDec(10)

	mockChecker.On("GetSafetyFactor", ctx).Return(safetyFactor)

	mockChecker.On("CalcMTPConsolidateCollateral", ctx, mtp, ptypes.BaseCurrency).Return(nil)
	mockChecker.On("SetMTP", ctx, mtp).Return(nil)

	_, err := k.OpenDefineAssets(ctx, ammPool.PoolId, msg, ptypes.BaseCurrency, false)
	// Expect no error
	assert.Nil(t, err)
	mockChecker.AssertExpectations(t)

	msg = &types.MsgOpen{
		Creator:         "",
		Leverage:        math.LegacyNewDec(12),
		Position:        types.Position_LONG,
		TradingAsset:    "uatom",
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000)),
		TakeProfitPrice: math.LegacyMustNewDecFromStr(types.TakeProfitPriceDefault),
	}

	eta = math.LegacyNewDec(11)

	leveragedAmtTokenIn = sdk.NewCoin(msg.Collateral.Denom, math.NewInt(12000))
	custodyAmount = math.NewInt(99)
	mockChecker.On("EstimateSwap", ctx, leveragedAmtTokenIn, msg.TradingAsset, ammPool).Return(custodyAmount, nil)

	mtp = types.NewMTP(msg.Creator, msg.Collateral.Denom, msg.TradingAsset, msg.Collateral.Denom, msg.TradingAsset, msg.Position, msg.Leverage, math.LegacyMustNewDecFromStr(types.TakeProfitPriceDefault), ammPool.PoolId)
	mockChecker.On("Borrow", ctx, msg.Collateral.Amount, custodyAmount, mtp, &ammPool, &pool, eta, ptypes.BaseCurrency, false).Return(nil)

	mockChecker.On("TakeInCustody", ctx, *mtp, &pool).Return(nil)

	lr = math.LegacyNewDec(50)

	mockChecker.On("GetMTPHealth", ctx, *mtp, ammPool, ptypes.BaseCurrency).Return(lr, nil)

	mockChecker.On("CalcMTPConsolidateCollateral", ctx, mtp, ptypes.BaseCurrency).Return(nil)
	mockChecker.On("SetMTP", ctx, mtp).Return(nil)

	_, err = k.OpenDefineAssets(ctx, ammPool.PoolId, msg, ptypes.BaseCurrency, false)
	// Expect no error
	assert.Nil(t, err)
	mockChecker.AssertExpectations(t)
}
