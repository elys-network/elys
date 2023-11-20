package keeper_test

import (
	"errors"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/margin/keeper"
	"github.com/elys-network/elys/x/margin/types"
	"github.com/elys-network/elys/x/margin/types/mocks"
	"github.com/stretchr/testify/assert"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	simapp "github.com/elys-network/elys/app"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestOpenShort_PoolNotFound(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenShortChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		OpenShortChecker: mockChecker,
	}

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			Leverage:         math.LegacyNewDec(10),
			CollateralAmount: math.NewInt(1),
			CollateralAsset:  "aaa",
			BorrowAsset:      "bbb",
		}
		poolId = uint64(42)
	)

	// Mock behavior
	mockChecker.On("GetMaxLeverageParam", ctx).Return(msg.Leverage)
	mockChecker.On("GetPool", ctx, poolId).Return(types.Pool{}, false)

	_, err := k.OpenShort(ctx, poolId, msg, ptypes.BaseCurrency)

	// Expect an error about the pool not existing
	assert.True(t, errors.Is(err, types.ErrPoolDoesNotExist))
	mockChecker.AssertExpectations(t)
}

func TestOpenShort_PoolDisabled(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenShortChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		OpenShortChecker: mockChecker,
	}

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			Leverage:         math.LegacyNewDec(10),
			CollateralAmount: math.NewInt(1),
			CollateralAsset:  "aaa",
			BorrowAsset:      "bbb",
		}
		poolId = uint64(42)
	)

	// Mock behaviors
	mockChecker.On("GetMaxLeverageParam", ctx).Return(msg.Leverage)
	mockChecker.On("GetPool", ctx, poolId).Return(types.Pool{}, true)
	mockChecker.On("IsPoolEnabled", ctx, poolId).Return(false)

	_, err := k.OpenShort(ctx, poolId, msg, ptypes.BaseCurrency)

	// Expect an error about the pool being disabled
	assert.True(t, errors.Is(err, types.ErrMTPDisabled))
	mockChecker.AssertExpectations(t)
}

func TestOpenShort_InsufficientLiabilities(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenShortChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		OpenShortChecker: mockChecker,
	}

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			Leverage:         math.LegacyNewDec(2),
			CollateralAmount: math.NewInt(1000),
			Creator:          "",
			CollateralAsset:  ptypes.BaseCurrency,
			BorrowAsset:      "uatom",
			Position:         types.Position_SHORT,
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

	// Mock the behaviors to get to the CheckMinLiabilities check
	mockChecker.On("GetMaxLeverageParam", ctx).Return(msg.Leverage)
	mockChecker.On("GetPool", ctx, ammPool.PoolId).Return(pool, true)
	mockChecker.On("IsPoolEnabled", ctx, ammPool.PoolId).Return(true)
	mockChecker.On("GetAmmPool", ctx, ammPool.PoolId, msg.BorrowAsset).Return(ammPool, nil) // Assuming a valid pool is returned

	collateralAmountDec := sdk.NewDecFromBigInt(msg.CollateralAmount.BigInt())
	leveragedAmount := sdk.NewInt(collateralAmountDec.Mul(msg.Leverage).TruncateInt().Int64())
	custodyAmtToken := sdk.NewCoin(ptypes.BaseCurrency, leveragedAmount)
	borrowingAmount := sdk.NewInt(99)

	mockChecker.On("EstimateSwapGivenOut", ctx, custodyAmtToken, msg.BorrowAsset, ammPool).Return(borrowingAmount, nil)

	// Mock the behavior where CheckMinLiabilities returns an error indicating insufficient liabilities
	liabilityError := errors.New("insufficient liabilities")
	collateralTokenAmt := sdk.NewCoin(msg.CollateralAsset, msg.CollateralAmount)

	mockChecker.On("CheckMinLiabilities", ctx, collateralTokenAmt, sdk.NewDec(1), pool, ammPool, msg.BorrowAsset).Return(liabilityError)

	_, err := k.OpenShort(ctx, ammPool.PoolId, msg, ptypes.BaseCurrency)

	// Expect the custom error indicating insufficient liabilities
	assert.True(t, errors.Is(err, liabilityError))
	mockChecker.AssertExpectations(t)
}

func TestOpenShort_InsufficientAmmPoolBalanceForCustody(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenShortChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		OpenShortChecker: mockChecker,
	}

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			Leverage:         math.LegacyNewDec(10),
			CollateralAmount: math.NewInt(1000),
			Creator:          "",
			CollateralAsset:  ptypes.BaseCurrency,
			BorrowAsset:      "uatom",
			Position:         types.Position_SHORT,
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
	mockChecker.On("GetAmmPool", ctx, ammPool.PoolId, msg.BorrowAsset).Return(ammPool, nil)

	collateralAmountDec := sdk.NewDecFromBigInt(msg.CollateralAmount.BigInt())
	leveragedAmount := sdk.NewInt(collateralAmountDec.Mul(msg.Leverage).TruncateInt().Int64())
	custodyAmtToken := sdk.NewCoin(ptypes.BaseCurrency, leveragedAmount)
	borrowingAmount := sdk.NewInt(99)

	mockChecker.On("EstimateSwapGivenOut", ctx, custodyAmtToken, msg.BorrowAsset, ammPool).Return(borrowingAmount, nil)

	collateralTokenAmt := sdk.NewCoin(msg.CollateralAsset, msg.CollateralAmount)
	eta := math.LegacyNewDec(9)

	mockChecker.On("CheckMinLiabilities", ctx, collateralTokenAmt, eta, pool, ammPool, msg.BorrowAsset).Return(nil)

	leveragedAmtTokenIn := sdk.NewCoin(msg.BorrowAsset, borrowingAmount)
	custodyAmount := math.NewInt(100000)

	mockChecker.On("EstimateSwap", ctx, leveragedAmtTokenIn, ptypes.BaseCurrency, ammPool).Return(custodyAmount, nil)

	_, err := k.OpenShort(ctx, ammPool.PoolId, msg, ptypes.BaseCurrency)

	// Expect an error about custody amount being too high
	assert.True(t, errors.Is(err, types.ErrCustodyTooHigh))
	mockChecker.AssertExpectations(t)
}

func TestOpenShort_ErrorsDuringOperations(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenShortChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		OpenShortChecker: mockChecker,
	}

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			Leverage:         math.LegacyNewDec(10),
			CollateralAmount: math.NewInt(1000),
			Creator:          "",
			CollateralAsset:  ptypes.BaseCurrency,
			BorrowAsset:      "uatom",
			Position:         types.Position_SHORT,
			TakeProfitPrice:  sdk.MustNewDecFromStr(types.TakeProfitPriceDefault),
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
	mockChecker.On("GetAmmPool", ctx, ammPool.PoolId, msg.BorrowAsset).Return(ammPool, nil)

	collateralAmountDec := sdk.NewDecFromBigInt(msg.CollateralAmount.BigInt())
	leveragedAmount := sdk.NewInt(collateralAmountDec.Mul(msg.Leverage).TruncateInt().Int64())
	custodyAmtToken := sdk.NewCoin(ptypes.BaseCurrency, leveragedAmount)
	borrowingAmount := sdk.NewInt(99)

	mockChecker.On("EstimateSwapGivenOut", ctx, custodyAmtToken, msg.BorrowAsset, ammPool).Return(borrowingAmount, nil)

	collateralTokenAmt := sdk.NewCoin(msg.CollateralAsset, msg.CollateralAmount)
	eta := math.LegacyNewDec(9)

	mockChecker.On("CheckMinLiabilities", ctx, collateralTokenAmt, eta, pool, ammPool, msg.BorrowAsset).Return(nil)

	leveragedAmtTokenIn := sdk.NewCoin(msg.BorrowAsset, borrowingAmount)
	custodyAmount := math.NewInt(199)

	mockChecker.On("EstimateSwap", ctx, leveragedAmtTokenIn, ptypes.BaseCurrency, ammPool).Return(custodyAmount, nil)

	mtp := types.NewMTP(msg.Creator, msg.CollateralAsset, ptypes.BaseCurrency, msg.Position, msg.Leverage, sdk.MustNewDecFromStr(types.TakeProfitPriceDefault), ammPool.PoolId)

	borrowError := errors.New("borrow error")
	mockChecker.On("Borrow", ctx, msg.CollateralAsset, ptypes.BaseCurrency, msg.CollateralAmount, custodyAmount, mtp, &ammPool, &pool, eta, ptypes.BaseCurrency).Return(borrowError)

	_, err := k.OpenShort(ctx, ammPool.PoolId, msg, ptypes.BaseCurrency)

	// Expect the borrow error
	assert.True(t, errors.Is(err, borrowError))
	mockChecker.AssertExpectations(t)
}

func TestOpenShort_LeverageRatioLessThanSafetyFactor(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenShortChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		OpenShortChecker: mockChecker,
	}

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			Leverage:         math.LegacyNewDec(10),
			CollateralAmount: math.NewInt(1000),
			Creator:          "",
			CollateralAsset:  ptypes.BaseCurrency,
			BorrowAsset:      "uatom",
			Position:         types.Position_SHORT,
			TakeProfitPrice:  sdk.MustNewDecFromStr(types.TakeProfitPriceDefault),
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
	mockChecker.On("GetAmmPool", ctx, ammPool.PoolId, msg.BorrowAsset).Return(ammPool, nil)

	collateralAmountDec := sdk.NewDecFromBigInt(msg.CollateralAmount.BigInt())
	leveragedAmount := sdk.NewInt(collateralAmountDec.Mul(msg.Leverage).TruncateInt().Int64())
	custodyAmtToken := sdk.NewCoin(ptypes.BaseCurrency, leveragedAmount)
	borrowingAmount := sdk.NewInt(99)

	mockChecker.On("EstimateSwapGivenOut", ctx, custodyAmtToken, msg.BorrowAsset, ammPool).Return(borrowingAmount, nil)

	collateralTokenAmt := sdk.NewCoin(msg.CollateralAsset, msg.CollateralAmount)
	eta := math.LegacyNewDec(9)

	mockChecker.On("CheckMinLiabilities", ctx, collateralTokenAmt, eta, pool, ammPool, msg.BorrowAsset).Return(nil)

	leveragedAmtTokenIn := sdk.NewCoin(msg.BorrowAsset, borrowingAmount)
	custodyAmount := math.NewInt(199)

	mockChecker.On("EstimateSwap", ctx, leveragedAmtTokenIn, ptypes.BaseCurrency, ammPool).Return(custodyAmount, nil)

	mtp := types.NewMTP(msg.Creator, msg.CollateralAsset, ptypes.BaseCurrency, msg.Position, msg.Leverage, sdk.MustNewDecFromStr(types.TakeProfitPriceDefault), ammPool.PoolId)

	mockChecker.On("Borrow", ctx, msg.CollateralAsset, ptypes.BaseCurrency, msg.CollateralAmount, custodyAmount, mtp, &ammPool, &pool, eta, ptypes.BaseCurrency).Return(nil)
	mockChecker.On("UpdatePoolHealth", ctx, &pool).Return(nil)
	mockChecker.On("TakeInCustody", ctx, *mtp, &pool).Return(nil)

	lr := math.LegacyNewDec(50)

	mockChecker.On("UpdateMTPHealth", ctx, *mtp, ammPool, ptypes.BaseCurrency).Return(lr, nil)
	mockChecker.On("GetSafetyFactor", ctx).Return(sdk.NewDec(100))

	_, err := k.OpenShort(ctx, ammPool.PoolId, msg, ptypes.BaseCurrency)

	// Expect an error indicating MTP is unhealthy
	assert.True(t, errors.Is(err, types.ErrMTPUnhealthy))
	mockChecker.AssertExpectations(t)
}

func TestOpenShort_Success(t *testing.T) {
	// Setup the mock checker
	mockChecker := new(mocks.OpenShortChecker)

	// Create an instance of Keeper with the mock checker
	k := keeper.Keeper{
		OpenShortChecker: mockChecker,
	}

	var (
		ctx = sdk.Context{} // Mock or setup a context
		msg = &types.MsgOpen{
			Leverage:         math.LegacyNewDec(10),
			CollateralAmount: math.NewInt(1000),
			Creator:          "",
			CollateralAsset:  ptypes.BaseCurrency,
			BorrowAsset:      "uatom",
			Position:         types.Position_SHORT,
			TakeProfitPrice:  sdk.MustNewDecFromStr(types.TakeProfitPriceDefault),
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
	mockChecker.On("GetAmmPool", ctx, ammPool.PoolId, msg.BorrowAsset).Return(ammPool, nil)

	collateralAmountDec := sdk.NewDecFromBigInt(msg.CollateralAmount.BigInt())
	leveragedAmount := sdk.NewInt(collateralAmountDec.Mul(msg.Leverage).TruncateInt().Int64())
	custodyAmtToken := sdk.NewCoin(ptypes.BaseCurrency, leveragedAmount)
	borrowingAmount := sdk.NewInt(99)

	mockChecker.On("EstimateSwapGivenOut", ctx, custodyAmtToken, msg.BorrowAsset, ammPool).Return(borrowingAmount, nil)

	collateralTokenAmt := sdk.NewCoin(msg.CollateralAsset, msg.CollateralAmount)
	eta := math.LegacyNewDec(9)

	mockChecker.On("CheckMinLiabilities", ctx, collateralTokenAmt, eta, pool, ammPool, msg.BorrowAsset).Return(nil)

	leveragedAmtTokenIn := sdk.NewCoin(msg.BorrowAsset, borrowingAmount)
	custodyAmount := math.NewInt(199)

	mockChecker.On("EstimateSwap", ctx, leveragedAmtTokenIn, ptypes.BaseCurrency, ammPool).Return(custodyAmount, nil)

	mtp := types.NewMTP(msg.Creator, msg.CollateralAsset, ptypes.BaseCurrency, msg.Position, msg.Leverage, sdk.MustNewDecFromStr(types.TakeProfitPriceDefault), ammPool.PoolId)

	mockChecker.On("Borrow", ctx, msg.CollateralAsset, ptypes.BaseCurrency, msg.CollateralAmount, custodyAmount, mtp, &ammPool, &pool, eta, ptypes.BaseCurrency).Return(nil)
	mockChecker.On("UpdatePoolHealth", ctx, &pool).Return(nil)
	mockChecker.On("TakeInCustody", ctx, *mtp, &pool).Return(nil)

	lr := math.LegacyNewDec(50)

	mockChecker.On("UpdateMTPHealth", ctx, *mtp, ammPool, ptypes.BaseCurrency).Return(lr, nil)

	safetyFactor := math.LegacyNewDec(10)

	mockChecker.On("GetSafetyFactor", ctx).Return(safetyFactor)

	mockChecker.On("CalcMTPConsolidateCollateral", ctx, mtp, ptypes.BaseCurrency).Return(nil)
	mockChecker.On("SetMTP", ctx, mtp).Return(nil)

	_, err := k.OpenShort(ctx, ammPool.PoolId, msg, ptypes.BaseCurrency)
	// Expect no error
	assert.Nil(t, err)
	mockChecker.AssertExpectations(t)
}

func TestOpenShort_BaseCurrency_Collateral(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	mk, amm, oracle := app.MarginKeeper, app.AmmKeeper, app.OracleKeeper

	// Setup coin prices
	SetupStableCoinPrices(ctx, oracle)

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(1000000))

	// Create a pool
	// Mint 100000USDC
	usdcToken := []sdk.Coin{sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(100000))}
	// Mint 100000ATOM
	atomToken := []sdk.Coin{sdk.NewCoin(ptypes.ATOM, sdk.NewInt(100000))}

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
			Token:  sdk.NewCoin(ptypes.ATOM, sdk.NewInt(100000)),
		},
		{
			Weight: sdk.NewInt(50),
			Token:  sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(10000)),
		},
	}

	argSwapFee := sdk.MustNewDecFromStr("0.0")
	argExitFee := sdk.MustNewDecFromStr("0.0")

	poolParams := &ammtypes.PoolParams{
		SwapFee: argSwapFee,
		ExitFee: argExitFee,
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

	// Balance check before create a margin position
	balances := app.BankKeeper.GetAllBalances(ctx, poolAddress)
	require.Equal(t, balances.AmountOf(ptypes.BaseCurrency), sdk.NewInt(10000))
	require.Equal(t, balances.AmountOf(ptypes.ATOM), sdk.NewInt(100000))

	// Create a margin position open msg
	msg2 := types.NewMsgOpen(
		addr[0].String(),
		ptypes.BaseCurrency,
		sdk.NewInt(100),
		ptypes.ATOM,
		types.Position_SHORT,
		sdk.NewDec(5),
		sdk.MustNewDecFromStr(types.TakeProfitPriceDefault),
	)

	_, err = mk.Open(ctx, msg2)
	require.NoError(t, err)

	mtps := mk.GetAllMTPs(ctx)
	require.Equal(t, len(mtps), 1)

	balances = app.BankKeeper.GetAllBalances(ctx, poolAddress)
	require.Equal(t, balances.AmountOf(ptypes.BaseCurrency), sdk.NewInt(10100))
	require.Equal(t, balances.AmountOf(ptypes.ATOM), sdk.NewInt(100000))

	_, found = mk.OpenShortChecker.GetPool(ctx, pool.PoolId)
	require.Equal(t, found, true)

	err = mk.InvariantCheck(ctx)
	require.Equal(t, err, nil)
}

func TestOpenShort_ATOM_Collateral(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	mk, amm, oracle := app.MarginKeeper, app.AmmKeeper, app.OracleKeeper

	// Setup coin prices
	SetupStableCoinPrices(ctx, oracle)

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(1000000))

	// Create a pool
	// Mint 100000USDC
	usdcToken := []sdk.Coin{sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(100000))}
	// Mint 100000ATOM
	atomToken := []sdk.Coin{sdk.NewCoin(ptypes.ATOM, sdk.NewInt(100000))}

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
			Token:  sdk.NewCoin(ptypes.ATOM, sdk.NewInt(1000)),
		},
		{
			Weight: sdk.NewInt(50),
			Token:  sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(10000)),
		},
	}

	argSwapFee := sdk.MustNewDecFromStr("0.0")
	argExitFee := sdk.MustNewDecFromStr("0.0")

	poolParams := &ammtypes.PoolParams{
		SwapFee: argSwapFee,
		ExitFee: argExitFee,
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

	// Balance check before create a margin position
	balances := app.BankKeeper.GetAllBalances(ctx, poolAddress)
	require.Equal(t, balances.AmountOf(ptypes.BaseCurrency), sdk.NewInt(10000))
	require.Equal(t, balances.AmountOf(ptypes.ATOM), sdk.NewInt(1000))

	// Create a margin position open msg
	msg2 := types.NewMsgOpen(
		addr[0].String(),
		ptypes.ATOM,
		sdk.NewInt(10),
		ptypes.ATOM,
		types.Position_SHORT,
		sdk.NewDec(5),
		sdk.MustNewDecFromStr(types.TakeProfitPriceDefault),
	)

	_, err = mk.Open(ctx, msg2)
	assert.True(t, errors.Is(err, sdkerrors.Wrap(types.ErrInvalidCollateralAsset, "collateral asset cannot be the same as the borrowed asset in a short position")))

	mtps := mk.GetAllMTPs(ctx)
	require.Equal(t, len(mtps), 0)

	balances = app.BankKeeper.GetAllBalances(ctx, poolAddress)
	require.Equal(t, balances.AmountOf(ptypes.BaseCurrency), sdk.NewInt(10000))
	require.Equal(t, balances.AmountOf(ptypes.ATOM), sdk.NewInt(1000))

	_, found = mk.OpenShortChecker.GetPool(ctx, pool.PoolId)
	require.Equal(t, found, false)

	err = mk.InvariantCheck(ctx)
	require.Equal(t, err, nil)
}
