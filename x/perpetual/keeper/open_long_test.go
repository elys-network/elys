package keeper_test

import (
	"errors"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	"github.com/elys-network/elys/x/perpetual/keeper"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/elys-network/elys/x/perpetual/types/mocks"
	"github.com/stretchr/testify/assert"

	"github.com/cometbft/cometbft/crypto/ed25519"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
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
					Token:  sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(10000)),
					Weight: sdk.NewInt(50),
				},
				{
					Token:  sdk.NewCoin("uatom", sdk.NewInt(10)),
					Weight: sdk.NewInt(50),
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
			TakeProfitPrice: sdk.MustNewDecFromStr(types.TakeProfitPriceDefault),
		}
		ammPool = ammtypes.Pool{
			PoolId: uint64(42),
			PoolAssets: []ammtypes.PoolAsset{
				{
					Token:  sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(10000)),
					Weight: sdk.NewInt(50),
				},
				{
					Token:  sdk.NewCoin("uatom", sdk.NewInt(10000)),
					Weight: sdk.NewInt(50),
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

	mtp := types.NewMTP(msg.Creator, msg.Collateral.Denom, msg.TradingAsset, msg.Collateral.Denom, msg.TradingAsset, msg.Position, msg.Leverage, sdk.MustNewDecFromStr(types.TakeProfitPriceDefault), ammPool.PoolId)

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
			TakeProfitPrice: sdk.MustNewDecFromStr(types.TakeProfitPriceDefault),
		}
		ammPool = ammtypes.Pool{
			PoolId: uint64(42),
			PoolAssets: []ammtypes.PoolAsset{
				{
					Token:  sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(10000)),
					Weight: sdk.NewInt(50),
				},
				{
					Token:  sdk.NewCoin("uatom", sdk.NewInt(10000)),
					Weight: sdk.NewInt(50),
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

	mtp := types.NewMTP(msg.Creator, msg.Collateral.Denom, msg.TradingAsset, msg.Collateral.Denom, msg.TradingAsset, msg.Position, msg.Leverage, sdk.MustNewDecFromStr(types.TakeProfitPriceDefault), ammPool.PoolId)

	mockChecker.On("Borrow", ctx, msg.Collateral.Amount, custodyAmount, mtp, &ammPool, &pool, eta, ptypes.BaseCurrency, false).Return(nil)
	mockChecker.On("UpdatePoolHealth", ctx, &pool).Return(nil)
	mockChecker.On("TakeInCustody", ctx, *mtp, &pool).Return(nil)

	lr := math.LegacyNewDec(50)

	mockChecker.On("GetMTPHealth", ctx, *mtp, ammPool, ptypes.BaseCurrency).Return(lr, nil)
	mockChecker.On("GetSafetyFactor", ctx).Return(sdk.NewDec(100))

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
			TakeProfitPrice: sdk.MustNewDecFromStr(types.TakeProfitPriceDefault),
		}
		ammPool = ammtypes.Pool{
			PoolId: uint64(42),
			PoolAssets: []ammtypes.PoolAsset{
				{
					Token:  sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(10000)),
					Weight: sdk.NewInt(50),
				},
				{
					Token:  sdk.NewCoin("uatom", sdk.NewInt(10000)),
					Weight: sdk.NewInt(50),
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

	mtp := types.NewMTP(msg.Creator, msg.Collateral.Denom, msg.TradingAsset, msg.Collateral.Denom, msg.TradingAsset, msg.Position, msg.Leverage, sdk.MustNewDecFromStr(types.TakeProfitPriceDefault), ammPool.PoolId)

	mockChecker.On("Borrow", ctx, msg.Collateral.Amount, custodyAmount, mtp, &ammPool, &pool, eta, ptypes.BaseCurrency, false).Return(nil)
	mockChecker.On("UpdatePoolHealth", ctx, &pool).Return(nil)
	mockChecker.On("TakeInCustody", ctx, *mtp, &pool).Return(nil)

	lr := math.LegacyNewDec(50)

	mockChecker.On("GetMTPHealth", ctx, *mtp, ammPool, ptypes.BaseCurrency).Return(lr, nil)

	safetyFactor := math.LegacyNewDec(10)

	mockChecker.On("GetSafetyFactor", ctx).Return(safetyFactor)

	mockChecker.On("SetMTP", ctx, mtp).Return(nil)

	_, err := k.OpenDefineAssets(ctx, ammPool.PoolId, msg, ptypes.BaseCurrency, false)
	// Expect no error
	assert.Nil(t, err)
	mockChecker.AssertExpectations(t)
}

func TestOpenLong_BaseCurrency_Collateral(t *testing.T) {
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
	// Mint 100000USDC
	usdcToken := []sdk.Coin{sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(200000000000))}
	// Mint 100000ATOM
	atomToken := []sdk.Coin{sdk.NewCoin(ptypes.ATOM, sdk.NewInt(200000000000))}

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

	// check block height
	require.Equal(t, int64(0), ctx.BlockHeight())

	pool, found := amm.GetPool(ctx, poolId)
	require.Equal(t, found, true)

	poolAddress := sdk.MustAccAddressFromBech32(pool.GetAddress())
	require.NoError(t, err)

	// Balance check before create a perpetual position
	balances := app.BankKeeper.GetAllBalances(ctx, poolAddress)
	require.Equal(t, balances.AmountOf(ptypes.BaseCurrency), sdk.NewInt(100000000000))
	require.Equal(t, balances.AmountOf(ptypes.ATOM), sdk.NewInt(10000000000))

	// Create a perpetual position open msg
	msg2 := types.NewMsgOpen(
		addr[0].String(),
		types.Position_LONG,
		sdk.NewDec(5),
		ptypes.ATOM,
		sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(100000000)),
		sdk.MustNewDecFromStr(types.TakeProfitPriceDefault),
		sdk.NewDec(100),
	)

	_, err = mk.Open(ctx, msg2, false)
	require.NoError(t, err)

	mtps := mk.GetAllMTPs(ctx)
	require.Equal(t, len(mtps), 1)

	balances = app.BankKeeper.GetAllBalances(ctx, poolAddress)
	require.Equal(t, balances.AmountOf(ptypes.BaseCurrency), sdk.NewInt(100100000000))
	require.Equal(t, balances.AmountOf(ptypes.ATOM), sdk.NewInt(10000000000))

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
		Collateral:                     sdk.NewInt(100000000),
		Liabilities:                    sdk.NewInt(400000000),
		BorrowInterestPaidCollateral:   sdk.NewInt(0),
		BorrowInterestPaidCustody:      sdk.NewInt(0),
		BorrowInterestUnpaidCollateral: sdk.NewInt(0),
		Custody:                        sdk.NewInt(486520593),
		TakeProfitLiabilities:          sdk.NewInt(473929244),
		TakeProfitCustody:              sdk.NewInt(486520593),
		MtpHealth:                      sdk.MustNewDecFromStr("1.234567958024691358"),
		Position:                       types.Position_LONG,
		Id:                             uint64(1),
		AmmPoolId:                      uint64(1),
		TakeProfitPrice:                sdk.MustNewDecFromStr(types.TakeProfitPriceDefault),
		TakeProfitBorrowRate:           sdk.MustNewDecFromStr("1.0"),
		FundingFeePaidCollateral:       sdk.NewInt(0),
		FundingFeePaidCustody:          sdk.NewInt(0),
		FundingFeeReceivedCollateral:   sdk.NewInt(0),
		FundingFeeReceivedCustody:      sdk.NewInt(0),
		OpenPrice:                      sdk.MustNewDecFromStr("1.027705727555914576"),
		StopLossPrice:                  sdk.NewDec(100),
	}, mtp)
}

func TestOpenLong_ATOM_Collateral(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	mk, amm, oracle := app.PerpetualKeeper, app.AmmKeeper, app.OracleKeeper

	// Setup coin prices
	SetupStableCoinPrices(ctx, oracle)

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(1000000000000))

	// Create a pool
	// Mint 100000USDC
	usdcToken := []sdk.Coin{sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(100000000000))}
	// Mint 100000ATOM
	atomToken := []sdk.Coin{sdk.NewCoin(ptypes.ATOM, sdk.NewInt(100000000000))}

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
			Token:  sdk.NewCoin(ptypes.ATOM, sdk.NewInt(1000000000)),
		},
		{
			Weight: sdk.NewInt(50),
			Token:  sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(10000000000)),
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

	// check block height
	require.Equal(t, int64(0), ctx.BlockHeight())

	pool, found := amm.GetPool(ctx, poolId)
	require.Equal(t, found, true)

	poolAddress := sdk.MustAccAddressFromBech32(pool.GetAddress())
	require.NoError(t, err)

	// Balance check before create a perpetual position
	balances := app.BankKeeper.GetAllBalances(ctx, poolAddress)
	require.Equal(t, balances.AmountOf(ptypes.BaseCurrency), sdk.NewInt(10000000000))
	require.Equal(t, balances.AmountOf(ptypes.ATOM), sdk.NewInt(1000000000))

	// Create a perpetual position open msg
	msg2 := types.NewMsgOpen(
		addr[0].String(),
		types.Position_LONG,
		sdk.NewDec(5),
		ptypes.ATOM,
		sdk.NewCoin(ptypes.ATOM, sdk.NewInt(10000000)),
		sdk.MustNewDecFromStr(types.TakeProfitPriceDefault),
		sdk.NewDec(100),
	)

	_, err = mk.Open(ctx, msg2, false)
	require.NoError(t, err)

	mtps := mk.GetAllMTPs(ctx)
	require.Equal(t, len(mtps), 1)

	balances = app.BankKeeper.GetAllBalances(ctx, poolAddress)
	require.Equal(t, balances.AmountOf(ptypes.BaseCurrency), sdk.NewInt(10000000000))
	require.Equal(t, balances.AmountOf(ptypes.ATOM), sdk.NewInt(1010000000))

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
		Collateral:                     sdk.NewInt(10000000),
		Liabilities:                    sdk.NewInt(40905431),
		BorrowInterestPaidCollateral:   sdk.NewInt(0),
		BorrowInterestPaidCustody:      sdk.NewInt(0),
		BorrowInterestUnpaidCollateral: sdk.NewInt(0),
		Custody:                        sdk.NewInt(50000000),
		TakeProfitLiabilities:          sdk.NewInt(48671331),
		TakeProfitCustody:              sdk.NewInt(50000000),
		MtpHealth:                      sdk.MustNewDecFromStr("1.119902772759909084"),
		Position:                       types.Position_LONG,
		Id:                             uint64(1),
		AmmPoolId:                      uint64(1),
		TakeProfitPrice:                sdk.MustNewDecFromStr(types.TakeProfitPriceDefault),
		TakeProfitBorrowRate:           sdk.MustNewDecFromStr("1.0"),
		FundingFeePaidCollateral:       sdk.NewInt(0),
		FundingFeePaidCustody:          sdk.NewInt(0),
		FundingFeeReceivedCollateral:   sdk.NewInt(0),
		FundingFeeReceivedCustody:      sdk.NewInt(0),
		OpenPrice:                      sdk.MustNewDecFromStr("1.017016260000000000"),
		StopLossPrice:                  sdk.NewDec(100),
	}, mtp)
}

func TestOpenLong_Long10XAtom1000Usdc(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	mk, amm, oracle := app.PerpetualKeeper, app.AmmKeeper, app.OracleKeeper

	// Setup coin prices
	SetupStableCoinPrices(ctx, oracle)

	provider := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "USDC",
		Price:     sdk.NewDec(1),
		Source:    "elys",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
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

	// Generate 1 random account with 1stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(1_000000))

	// Create a pool
	// Mint 100_000_000USDC
	usdcToken := []sdk.Coin{sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(100_000_000_000000))}
	// Mint 100_000_000ATOM
	atomToken := []sdk.Coin{sdk.NewCoin(ptypes.ATOM, sdk.NewInt(100_000_000_000000))}

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
			Token:  sdk.NewCoin(ptypes.ATOM, sdk.NewInt(10_000_000_000000)),
		},
		{
			Weight: sdk.NewInt(50),
			Token:  sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(10_000_000_000000)),
		},
	}

	argSwapFee := sdk.MustNewDecFromStr("0.0")
	argExitFee := sdk.MustNewDecFromStr("0.0")

	poolParams := &ammtypes.PoolParams{
		UseOracle:                   true,
		ExternalLiquidityRatio:      sdk.NewDec(2),
		WeightBreakingFeeMultiplier: sdk.ZeroDec(),
		WeightBreakingFeeExponent:   sdk.NewDecWithPrec(0, 1), // 2.5
		WeightRecoveryFeePortion:    sdk.NewDecWithPrec(0, 2), // 10%
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

	// check block height
	require.Equal(t, int64(0), ctx.BlockHeight())

	pool, found := amm.GetPool(ctx, poolId)
	require.Equal(t, found, true)

	poolAddress := sdk.MustAccAddressFromBech32(pool.GetAddress())
	require.NoError(t, err)

	// Balance check before create a perpetual position
	balances := app.BankKeeper.GetAllBalances(ctx, poolAddress)
	require.Equal(t, balances.AmountOf(ptypes.BaseCurrency), sdk.NewInt(10_000_000_000000))
	require.Equal(t, balances.AmountOf(ptypes.ATOM), sdk.NewInt(10_000_000_000000))

	// Create a perpetual position open msg
	msg2 := types.NewMsgOpen(
		addr[0].String(),
		types.Position_LONG,
		sdk.NewDec(10),
		ptypes.ATOM,
		sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(1_000_000000)),
		sdk.MustNewDecFromStr("5.0"),
		sdk.ZeroDec(),
	)

	_, err = mk.Open(ctx, msg2, false)
	require.NoError(t, err)

	mtps := mk.GetAllMTPs(ctx)
	require.Equal(t, len(mtps), 1)

	balances = app.BankKeeper.GetAllBalances(ctx, poolAddress)
	require.Equal(t, balances.AmountOf(ptypes.BaseCurrency), sdk.NewInt(10_001_000_000000))
	require.Equal(t, balances.AmountOf(ptypes.ATOM), sdk.NewInt(10_000_000_000000))

	_, found = mk.OpenDefineAssetsChecker.GetPool(ctx, pool.PoolId)
	require.Equal(t, found, true)

	err = mk.InvariantCheck(ctx)
	require.Equal(t, err, nil)

	mtp := mtps[0]

	// Check MTP
	require.Equal(t, types.MTP{
		Address:                        addr[0].String(),
		CollateralAsset:                ptypes.BaseCurrency,
		TradingAsset:                   ptypes.ATOM,
		LiabilitiesAsset:               ptypes.BaseCurrency,
		CustodyAsset:                   ptypes.ATOM,
		Collateral:                     sdk.NewInt(1_000_000000),
		Liabilities:                    sdk.NewInt(9_000_000000),
		BorrowInterestPaidCollateral:   sdk.NewInt(0),
		BorrowInterestPaidCustody:      sdk.NewInt(0),
		BorrowInterestUnpaidCollateral: sdk.NewInt(0),
		Custody:                        sdk.NewInt(2_276_506970),
		TakeProfitLiabilities:          sdk.NewInt(7_898_168205),
		TakeProfitCustody:              sdk.NewInt(1_800_000000),
		MtpHealth:                      sdk.MustNewDecFromStr("1.110494168573014992"),
		Position:                       types.Position_LONG,
		Id:                             uint64(1),
		AmmPoolId:                      uint64(1),
		TakeProfitPrice:                sdk.MustNewDecFromStr("5.0"),
		TakeProfitBorrowRate:           sdk.MustNewDecFromStr("1.0"),
		FundingFeePaidCollateral:       sdk.NewInt(0),
		FundingFeePaidCustody:          sdk.NewInt(0),
		FundingFeeReceivedCollateral:   sdk.NewInt(0),
		FundingFeeReceivedCustody:      sdk.NewInt(0),
		OpenPrice:                      sdk.MustNewDecFromStr("4.392694655356139762"),
		StopLossPrice:                  sdk.ZeroDec(),
	}, mtp)

	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "USDC",
		Price:     sdk.NewDec(1),
		Source:    "elys",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "ATOM",
		Price:     sdk.MustNewDecFromStr("5.0"),
		Source:    "atom",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "uatom",
		Price:     sdk.MustNewDecFromStr("5.0"),
		Source:    "uatom",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})

	resp, _, _ := mk.GetMTPsForAddressWithPagination(ctx, addr[0], nil)
	require.Equal(t, resp[0].Pnl, sdk.NewDec(1_380_312708))
}

func TestOpenLongConsolidate_Success(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenDefineAssetsChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		OpenDefineAssetsChecker: mockChecker,
	}

	addr := "cosmos1ysxv266l8w76lq0vy44ktzajdr9u9yhlxzlvga"

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			Creator:         addr,
			Leverage:        math.LegacyNewDec(10),
			Position:        types.Position_LONG,
			TradingAsset:    "uatom",
			Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000)),
			TakeProfitPrice: sdk.MustNewDecFromStr(types.TakeProfitPriceDefault),
		}
		ammPool = ammtypes.Pool{
			PoolId: uint64(42),
			PoolAssets: []ammtypes.PoolAsset{
				{
					Token:  sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(10000000)),
					Weight: sdk.NewInt(50),
				},
				{
					Token:  sdk.NewCoin("uatom", sdk.NewInt(10000000)),
					Weight: sdk.NewInt(50),
				},
			},
		}
		pool = types.Pool{
			AmmPoolId: ammPool.PoolId,
		}
	)

	// Mock behaviors
	mockChecker.On("GetMaxLeverageParam", ctx).Return(math.LegacyNewDec(50)).Once()
	mockChecker.On("GetPool", ctx, ammPool.PoolId).Return(pool, true).Once()
	mockChecker.On("IsPoolEnabled", ctx, ammPool.PoolId).Return(true).Once()
	mockChecker.On("GetAmmPool", ctx, ammPool.PoolId, msg.TradingAsset).Return(ammPool, nil).Once()

	eta := math.LegacyNewDec(9)

	leveragedAmtTokenIn := sdk.NewCoin(msg.Collateral.Denom, math.NewInt(10000))
	custodyAmount := math.NewInt(99)

	mockChecker.On("EstimateSwap", ctx, leveragedAmtTokenIn, msg.TradingAsset, ammPool).Return(custodyAmount, nil).Once()

	mtp := types.NewMTP(msg.Creator, msg.Collateral.Denom, msg.TradingAsset, msg.Collateral.Denom, msg.TradingAsset, msg.Position, msg.Leverage, sdk.MustNewDecFromStr(types.TakeProfitPriceDefault), ammPool.PoolId)

	mockChecker.On("Borrow", ctx, msg.Collateral.Amount, custodyAmount, mtp, &ammPool, &pool, eta, ptypes.BaseCurrency, false).Return(nil).Once()
	mockChecker.On("UpdatePoolHealth", ctx, &pool).Return(nil).Once()
	mockChecker.On("TakeInCustody", ctx, *mtp, &pool).Return(nil).Once()

	lr := math.LegacyNewDec(50)

	mockChecker.On("GetMTPHealth", ctx, *mtp, ammPool, ptypes.BaseCurrency).Return(lr, nil).Once()

	safetyFactor := math.LegacyNewDec(10)

	mockChecker.On("GetSafetyFactor", ctx).Return(safetyFactor).Once()

	mockChecker.On("SetMTP", ctx, mtp).Return(nil).Once()

	existingMtp, err := k.OpenDefineAssets(ctx, ammPool.PoolId, msg, ptypes.BaseCurrency, false)
	// Expect no error
	assert.Nil(t, err)
	mockChecker.AssertExpectations(t)

	msg = &types.MsgOpen{
		Creator:         addr,
		Leverage:        math.LegacyNewDec(12),
		Position:        types.Position_LONG,
		TradingAsset:    "uatom",
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000)),
		TakeProfitPrice: sdk.MustNewDecFromStr(types.TakeProfitPriceDefault),
	}

	// Mock behaviors
	mockChecker.On("GetMaxLeverageParam", ctx).Return(math.LegacyNewDec(50)).Once()
	mockChecker.On("GetPool", ctx, ammPool.PoolId).Return(pool, true).Once()
	mockChecker.On("IsPoolEnabled", ctx, ammPool.PoolId).Return(true).Once()
	mockChecker.On("GetAmmPool", ctx, ammPool.PoolId, msg.TradingAsset).Return(ammPool, nil).Once()

	eta = math.LegacyNewDec(11)

	leveragedAmtTokenIn = sdk.NewCoin(msg.Collateral.Denom, math.NewInt(12000))
	custodyAmount = math.NewInt(99)
	mockChecker.On("EstimateSwap", ctx, leveragedAmtTokenIn, msg.TradingAsset, ammPool).Return(custodyAmount, nil).Once()

	mtp = types.NewMTP(msg.Creator, msg.Collateral.Denom, msg.TradingAsset, msg.Collateral.Denom, msg.TradingAsset, msg.Position, msg.Leverage, sdk.MustNewDecFromStr(types.TakeProfitPriceDefault), ammPool.PoolId)
	mockChecker.On("Borrow", ctx, msg.Collateral.Amount, custodyAmount, mtp, &ammPool, &pool, eta, ptypes.BaseCurrency, false).Return(nil).Once()
	mockChecker.On("UpdatePoolHealth", ctx, &pool).Return(nil).Once()
	mockChecker.On("TakeInCustody", ctx, *mtp, &pool).Return(nil).Once()

	lr = math.LegacyNewDec(50)

	mockChecker.On("GetMTPHealth", ctx, *mtp, ammPool, ptypes.BaseCurrency).Return(lr, nil).Once()
	mockChecker.On("GetSafetyFactor", ctx).Return(safetyFactor).Once()

	mockChecker.On("SetMTP", ctx, mtp).Return(nil).Once()

	newMtp, err := k.OpenDefineAssets(ctx, ammPool.PoolId, msg, ptypes.BaseCurrency, false)

	// Expect no error
	assert.Nil(t, err)
	mockChecker.AssertExpectations(t)

	// Mock behaviors
	mockChecker.On("GetPool", ctx, ammPool.PoolId).Return(pool, true).Once()
	mockChecker.On("IsPoolEnabled", ctx, ammPool.PoolId).Return(true).Once()
	mockChecker.On("GetAmmPool", ctx, ammPool.PoolId, msg.TradingAsset).Return(ammPool, nil).Once()

	mockChecker.On("SetMTP", ctx, existingMtp).Return(nil).Once()
	mockChecker.On("DestroyMTP", ctx, newMtp.GetAccountAddress(), newMtp.Id).Return(nil).Once()

	mockChecker.On("UpdateOpenPrice", ctx, existingMtp, ammPool, ptypes.BaseCurrency).Return(nil).Once()
	mockChecker.On("EmitOpenEvent", ctx, existingMtp).Once()

	_, err = k.OpenConsolidate(ctx, existingMtp, newMtp, msg, ptypes.BaseCurrency)

	// Expect no error
	assert.Nil(t, err)
	mockChecker.AssertExpectations(t)
}

func TestOpenLongConsolidate_FiniteTakeProfitPrice_Success(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenDefineAssetsChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		OpenDefineAssetsChecker: mockChecker,
	}

	addr := "cosmos1ysxv266l8w76lq0vy44ktzajdr9u9yhlxzlvga"

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			Creator:         addr,
			Leverage:        math.LegacyNewDec(10),
			Position:        types.Position_LONG,
			TradingAsset:    "uatom",
			Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000)),
			TakeProfitPrice: sdk.MustNewDecFromStr("5.0"),
		}
		ammPool = ammtypes.Pool{
			PoolId: uint64(42),
			PoolAssets: []ammtypes.PoolAsset{
				{
					Token:  sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(10000000)),
					Weight: sdk.NewInt(50),
				},
				{
					Token:  sdk.NewCoin("uatom", sdk.NewInt(10000000)),
					Weight: sdk.NewInt(50),
				},
			},
		}
		pool = types.Pool{
			AmmPoolId: ammPool.PoolId,
		}
	)

	// Mock behaviors
	mockChecker.On("GetMaxLeverageParam", ctx).Return(math.LegacyNewDec(50)).Once()
	mockChecker.On("GetPool", ctx, ammPool.PoolId).Return(pool, true).Once()
	mockChecker.On("IsPoolEnabled", ctx, ammPool.PoolId).Return(true).Once()
	mockChecker.On("GetAmmPool", ctx, ammPool.PoolId, msg.TradingAsset).Return(ammPool, nil).Once()

	eta := math.LegacyNewDec(9)

	leveragedAmtTokenIn := sdk.NewCoin(msg.Collateral.Denom, math.NewInt(10000))
	custodyAmount := math.NewInt(99)

	mockChecker.On("EstimateSwap", ctx, leveragedAmtTokenIn, msg.TradingAsset, ammPool).Return(custodyAmount, nil).Once()

	mtp := types.NewMTP(msg.Creator, msg.Collateral.Denom, msg.TradingAsset, msg.Collateral.Denom, msg.TradingAsset, msg.Position, msg.Leverage, sdk.MustNewDecFromStr("5.0"), ammPool.PoolId)

	mockChecker.On("Borrow", ctx, msg.Collateral.Amount, custodyAmount, mtp, &ammPool, &pool, eta, ptypes.BaseCurrency, false).Return(nil).Once()
	mockChecker.On("UpdatePoolHealth", ctx, &pool).Return(nil).Once()
	mockChecker.On("TakeInCustody", ctx, *mtp, &pool).Return(nil).Once()

	lr := math.LegacyNewDec(50)

	mockChecker.On("GetMTPHealth", ctx, *mtp, ammPool, ptypes.BaseCurrency).Return(lr, nil).Once()

	safetyFactor := math.LegacyNewDec(10)

	mockChecker.On("GetSafetyFactor", ctx).Return(safetyFactor).Once()

	mockChecker.On("SetMTP", ctx, mtp).Return(nil).Once()

	existingMtp, err := k.OpenDefineAssets(ctx, ammPool.PoolId, msg, ptypes.BaseCurrency, false)
	// Expect no error
	assert.Nil(t, err)
	mockChecker.AssertExpectations(t)

	msg = &types.MsgOpen{
		Creator:         addr,
		Leverage:        math.LegacyNewDec(12),
		Position:        types.Position_LONG,
		TradingAsset:    "uatom",
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000)),
		TakeProfitPrice: sdk.MustNewDecFromStr("10.0"),
	}

	// Mock behaviors
	mockChecker.On("GetMaxLeverageParam", ctx).Return(math.LegacyNewDec(50)).Once()
	mockChecker.On("GetPool", ctx, ammPool.PoolId).Return(pool, true).Once()
	mockChecker.On("IsPoolEnabled", ctx, ammPool.PoolId).Return(true).Once()
	mockChecker.On("GetAmmPool", ctx, ammPool.PoolId, msg.TradingAsset).Return(ammPool, nil).Once()

	eta = math.LegacyNewDec(11)

	leveragedAmtTokenIn = sdk.NewCoin(msg.Collateral.Denom, math.NewInt(12000))
	custodyAmount = math.NewInt(99)
	mockChecker.On("EstimateSwap", ctx, leveragedAmtTokenIn, msg.TradingAsset, ammPool).Return(custodyAmount, nil).Once()

	mtp = types.NewMTP(msg.Creator, msg.Collateral.Denom, msg.TradingAsset, msg.Collateral.Denom, msg.TradingAsset, msg.Position, msg.Leverage, sdk.MustNewDecFromStr("10.0"), ammPool.PoolId)
	mockChecker.On("Borrow", ctx, msg.Collateral.Amount, custodyAmount, mtp, &ammPool, &pool, eta, ptypes.BaseCurrency, false).Return(nil).Once()
	mockChecker.On("UpdatePoolHealth", ctx, &pool).Return(nil).Once()
	mockChecker.On("TakeInCustody", ctx, *mtp, &pool).Return(nil).Once()

	lr = math.LegacyNewDec(50)

	mockChecker.On("GetMTPHealth", ctx, *mtp, ammPool, ptypes.BaseCurrency).Return(lr, nil).Once()
	mockChecker.On("GetSafetyFactor", ctx).Return(safetyFactor).Once()

	mockChecker.On("SetMTP", ctx, mtp).Return(nil).Once()

	newMtp, err := k.OpenDefineAssets(ctx, ammPool.PoolId, msg, ptypes.BaseCurrency, false)

	// Expect no error
	assert.Nil(t, err)
	mockChecker.AssertExpectations(t)

	// Mock behaviors
	mockChecker.On("GetPool", ctx, ammPool.PoolId).Return(pool, true).Once()
	mockChecker.On("IsPoolEnabled", ctx, ammPool.PoolId).Return(true).Once()
	mockChecker.On("GetAmmPool", ctx, ammPool.PoolId, msg.TradingAsset).Return(ammPool, nil).Once()

	mockChecker.On("SetMTP", ctx, existingMtp).Return(nil).Once()
	mockChecker.On("DestroyMTP", ctx, newMtp.GetAccountAddress(), newMtp.Id).Return(nil).Once()

	mockChecker.On("UpdateOpenPrice", ctx, existingMtp, ammPool, ptypes.BaseCurrency).Return(nil).Once()
	mockChecker.On("EmitOpenEvent", ctx, existingMtp).Once()

	existingMtp.Custody = sdk.NewInt(1000)
	newMtp.Custody = sdk.NewInt(2000)

	_, err = k.OpenConsolidate(ctx, existingMtp, newMtp, msg, ptypes.BaseCurrency)

	// Expect no error
	assert.Nil(t, err)
	mockChecker.AssertExpectations(t)
}
