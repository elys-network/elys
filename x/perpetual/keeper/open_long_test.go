package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) TestOpenLong() {
	suite.SetupCoinPrices()
	addr := suite.AddAccounts(10)
	amount := sdk.NewInt(1000)
	poolCreator := addr[0]
	positionCreator := addr[1]
	poolId := uint64(1)
	var ammPool ammtypes.Pool
	msg := &types.MsgOpen{
		Creator:         positionCreator.String(),
		Leverage:        math.LegacyNewDec(2),
		Position:        types.Position_LONG,
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, amount),
		TakeProfitPrice: sdk.ZeroDec(),
		StopLossPrice:   sdk.ZeroDec(),
	}
	testCases := []struct {
		name                 string
		expectErrMsg         string
		isBroker             bool
		prerequisiteFunction func()
		postValidateFunction func(mtp *types.MTP)
	}{
		{
			"pool not found",
			types.ErrPoolDoesNotExist.Error(),
			false,
			func() {
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"pool is disabled",
			"perpetual not enabled for pool",
			false,
			func() {
				ammPool = suite.SetAndGetAmmPool(poolCreator, poolId, true, sdk.ZeroDec(), sdk.ZeroDec(), ptypes.ATOM, amount.MulRaw(10), amount.MulRaw(10))
				pool := types.NewPool(poolId)
				err := pool.InitiatePool(suite.ctx, &ammPool)
				suite.Require().NoError(err)
				pool.Enabled = false
				suite.app.PerpetualKeeper.SetPool(suite.ctx, pool)
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"amm pool not found",
			"pool does not exist",
			false,
			func() {
				pool, found := suite.app.PerpetualKeeper.GetPool(suite.ctx, 1)
				suite.Require().True(found)
				pool.Enabled = true
				suite.app.PerpetualKeeper.SetPool(suite.ctx, pool)
				suite.app.AmmKeeper.RemovePool(suite.ctx, ammPool.PoolId)
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"collateral asset neither base currency nor present in the pool",
			"(uelys) does not exist in the pool",
			false,
			func() {
				err := suite.app.AmmKeeper.SetPool(suite.ctx, ammPool)
				suite.Require().NoError(err)
				msg.Collateral.Denom = ptypes.Elys
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"collateral is same as trading asset but pool doesn't have enough depth",
			"borrowed amount is higher than pool depth",
			false,
			func() {
				msg.Collateral.Denom = ptypes.ATOM
				params := suite.app.PerpetualKeeper.GetParams(suite.ctx)
				params.BorrowInterestRateMin = sdk.ZeroDec()
				err := suite.app.PerpetualKeeper.SetParams(suite.ctx, &params)
				suite.Require().NoError(err)
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"collateral amount is too high",
			"borrowed amount is higher than pool depth",
			false,
			func() {
				msg.Collateral.Denom = ptypes.BaseCurrency
				msg.Collateral.Amount = msg.Collateral.Amount.MulRaw(1000_000_000)
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"minimum liabilities is 0",
			"minimum borrow interest rate is zero",
			false,
			func() {
				msg.Collateral.Amount = msg.Collateral.Amount.QuoRaw(1000_000_000)
				params := suite.app.PerpetualKeeper.GetParams(suite.ctx)
				params.BorrowInterestRateMin = sdk.ZeroDec()
				err := suite.app.PerpetualKeeper.SetParams(suite.ctx, &params)
				suite.Require().NoError(err)
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"Borrow fails: lack of funds",
			"user does not have enough balance of the required coin",
			false,
			func() {
				params := suite.app.PerpetualKeeper.GetParams(suite.ctx)
				params.BorrowInterestRateMin = sdk.MustNewDecFromStr("0.12")
				err := suite.app.PerpetualKeeper.SetParams(suite.ctx, &params)
				suite.Require().NoError(err)
				err = suite.app.BankKeeper.SendCoinsFromAccountToModule(suite.ctx, positionCreator, govtypes.ModuleName, sdk.NewCoins(sdk.NewCoin(ptypes.BaseCurrency, suite.GetAccountIssueAmount())))
				suite.Require().NoError(err)
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"success: collateral USDC, trading asset ATOM, stop loss price 0, TakeProfitPrice 0",
			"",
			false,
			func() {
				err := suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, govtypes.ModuleName, positionCreator, sdk.NewCoins(sdk.NewCoin(ptypes.BaseCurrency, suite.GetAccountIssueAmount())))
				suite.Require().NoError(err)
				msg.Collateral.Denom = ptypes.BaseCurrency
				msg.Collateral.Amount = amount
				msg.TradingAsset = ptypes.ATOM
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"success: collateral ATOM, trading asset ATOM, stop loss price 0, TakeProfitPrice 0",
			"",
			false,
			func() {
				tokensIn := sdk.NewCoins(sdk.NewCoin(ptypes.ATOM, sdk.NewInt(1000_000_000)), sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(1000_000_000)))
				suite.AddLiquidity(ammPool, addr[3], tokensIn)
				msg.Creator = addr[2].String()
				msg.Collateral.Denom = ptypes.ATOM
				msg.Collateral.Amount = amount
				msg.TradingAsset = ptypes.ATOM
				msg.Leverage = sdk.OneDec().MulInt64(2)
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"success: collateral USDC, trading asset USDC, stop loss price 0, TakeProfitPrice 0",
			"",
			false,
			func() {
				msg.Creator = addr[2].String()
				msg.Collateral.Denom = ptypes.BaseCurrency
				msg.Collateral.Amount = amount
				msg.TradingAsset = ptypes.BaseCurrency
				msg.Leverage = sdk.OneDec().MulInt64(2)
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"success: collateral ATOM, trading asset USDC, stop loss price 0, TakeProfitPrice 0",
			"",
			false,
			func() {
				tokensIn := sdk.NewCoins(sdk.NewCoin(ptypes.ATOM, sdk.NewInt(1000_000_000)), sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(1000_000_000)))
				suite.AddLiquidity(ammPool, addr[5], tokensIn)
				msg.Creator = addr[4].String()
				msg.Collateral.Denom = ptypes.ATOM
				msg.Collateral.Amount = amount
				msg.TradingAsset = ptypes.BaseCurrency
				msg.Leverage = sdk.OneDec().MulInt64(2)

				params := suite.app.PerpetualKeeper.GetParams(suite.ctx)
				params.SafetyFactor = sdk.MustNewDecFromStr("0.01")
				err := suite.app.PerpetualKeeper.SetParams(suite.ctx, &params)
				suite.Require().NoError(err)
			},
			func(mtp *types.MTP) {
			},
		},
		{
			"collateral is USDC, trading asset is ATOM, amm pool has enough USDC but not enough ATOM",
			"amount too low",
			false,
			func() {
				suite.ResetSuite()
				suite.SetupCoinPrices()
				addr = suite.AddAccounts(10)
				poolCreator = addr[0]
				positionCreator = addr[1]
				ammPool = suite.SetAndGetAmmPool(poolCreator, poolId, true, sdk.ZeroDec(), sdk.ZeroDec(), ptypes.ATOM, amount.MulRaw(1000), sdk.NewInt(2))
				pool := types.NewPool(poolId)
				err := pool.InitiatePool(suite.ctx, &ammPool)
				suite.Require().NoError(err)
				pool.Enabled = true
				suite.app.PerpetualKeeper.SetPool(suite.ctx, pool)
				params := suite.app.PerpetualKeeper.GetParams(suite.ctx)
				params.BorrowInterestRateMin = sdk.MustNewDecFromStr("0.12")
				err = suite.app.PerpetualKeeper.SetParams(suite.ctx, &params)
				suite.Require().NoError(err)

				msg.Creator = positionCreator.String()
				msg.Collateral.Denom = ptypes.BaseCurrency
				msg.Collateral.Amount = amount
				msg.TradingAsset = ptypes.ATOM
			},
			func(mtp *types.MTP) {
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			err := msg.ValidateBasic()
			suite.Require().NoError(err)
			mtp, err := suite.app.PerpetualKeeper.OpenLong(suite.ctx, poolId, msg, ptypes.BaseCurrency, tc.isBroker)
			if tc.expectErrMsg != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().NoError(err)
			}
			tc.postValidateFunction(mtp)
		})
	}
}

//func TestOpenLong_InsufficientAmmPoolBalanceForCustody(t *testing.T) {
//	// Setup the mock checker
//	mockChecker := new(mocks.OpenLongChecker)
//
//	// Create an instance of Keeper with the mock checker
//	k := keeper.Keeper{
//		OpenLongChecker: mockChecker,
//	}
//
//	var (
//		ctx = sdk.Context{} // Mock or setup a context
//		msg = &types.MsgOpen{
//			Creator:      "",
//			Leverage:     math.LegacyNewDec(10),
//			Position:     types.Position_LONG,
//			TradingAsset: "uatom",
//			Collateral:   sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000)),
//		}
//		ammPool = ammtypes.Pool{
//			PoolId: uint64(42),
//			PoolAssets: []ammtypes.PoolAsset{
//				{
//					Token:  sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(10000)),
//					Weight: sdk.NewInt(50),
//				},
//				{
//					Token:  sdk.NewCoin("uatom", sdk.NewInt(10)),
//					Weight: sdk.NewInt(50),
//				},
//			},
//		}
//		pool = types.Pool{
//			AmmPoolId: ammPool.PoolId,
//		}
//	)
//	// Mock behaviors
//	mockChecker.On("GetMaxLeverageParam", ctx).Return(msg.Leverage)
//	mockChecker.On("GetPool", ctx, ammPool.PoolId).Return(pool, true)
//	mockChecker.On("IsPoolEnabled", ctx, ammPool.PoolId).Return(true)
//	mockChecker.On("GetAmmPool", ctx, ammPool.PoolId, msg.TradingAsset).Return(ammPool, nil)
//
//	eta := math.LegacyNewDec(9)
//
//	mockChecker.On("CheckMinLiabilities", ctx, msg.Collateral, eta, ammPool, msg.TradingAsset, ptypes.BaseCurrency).Return(nil)
//
//	leveragedAmtTokenIn := sdk.NewCoin(msg.Collateral.Denom, math.NewInt(10000))
//	custodyAmount := math.NewInt(99)
//
//	mockChecker.On("EstimateSwap", ctx, leveragedAmtTokenIn, msg.TradingAsset, ammPool).Return(custodyAmount, nil)
//
//	_, err := k.OpenLong(ctx, ammPool.PoolId, msg, ptypes.BaseCurrency, false)
//
//	// Expect an error about custody amount being too high
//	assert.True(t, errors.Is(err, types.ErrCustodyTooHigh))
//	mockChecker.AssertExpectations(t)
//}
//
//func TestOpenLong_ErrorsDuringOperations(t *testing.T) {
//	// Setup the mock checker
//	mockChecker := new(mocks.OpenLongChecker)
//
//	// Create an instance of Keeper with the mock checker
//	k := keeper.Keeper{
//		OpenLongChecker: mockChecker,
//	}
//
//	var (
//		ctx = sdk.Context{} // Mock or setup a context
//		msg = &types.MsgOpen{
//			Creator:         "",
//			Leverage:        math.LegacyNewDec(10),
//			Position:        types.Position_LONG,
//			TradingAsset:    "uatom",
//			Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000)),
//			TakeProfitPrice: sdk.NewDecFromInt(types.TakeProfitPriceDefault),
//		}
//		ammPool = ammtypes.Pool{
//			PoolId: uint64(42),
//			PoolAssets: []ammtypes.PoolAsset{
//				{
//					Token:  sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(10000)),
//					Weight: sdk.NewInt(50),
//				},
//				{
//					Token:  sdk.NewCoin("uatom", sdk.NewInt(10000)),
//					Weight: sdk.NewInt(50),
//				},
//			},
//		}
//		pool = types.Pool{
//			AmmPoolId: ammPool.PoolId,
//		}
//	)
//
//	// Mock behaviors
//	mockChecker.On("GetMaxLeverageParam", ctx).Return(msg.Leverage)
//	mockChecker.On("GetPool", ctx, ammPool.PoolId).Return(pool, true)
//	mockChecker.On("IsPoolEnabled", ctx, ammPool.PoolId).Return(true)
//	mockChecker.On("GetAmmPool", ctx, ammPool.PoolId, msg.TradingAsset).Return(ammPool, nil)
//
//	eta := math.LegacyNewDec(9)
//
//	mockChecker.On("CheckMinLiabilities", ctx, msg.Collateral, eta, ammPool, msg.TradingAsset, ptypes.BaseCurrency).Return(nil)
//
//	leveragedAmtTokenIn := sdk.NewCoin(msg.Collateral.Denom, math.NewInt(10000))
//	custodyAmount := math.NewInt(99)
//
//	mockChecker.On("EstimateSwap", ctx, leveragedAmtTokenIn, msg.TradingAsset, ammPool).Return(custodyAmount, nil)
//
//	mtp := types.NewMTP(msg.Creator, msg.Collateral.Denom, msg.TradingAsset, msg.Collateral.Denom, msg.TradingAsset, msg.Position, msg.Leverage, sdk.NewDecFromInt(types.TakeProfitPriceDefault), ammPool.PoolId)
//
//	borrowError := errors.New("borrow error")
//	mockChecker.On("Borrow", ctx, msg.Collateral.Amount, custodyAmount, mtp, &ammPool, &pool, eta, ptypes.BaseCurrency, false).Return(borrowError)
//
//	_, err := k.OpenLong(ctx, ammPool.PoolId, msg, ptypes.BaseCurrency, false)
//
//	// Expect the borrow error
//	assert.True(t, errors.Is(err, borrowError))
//	mockChecker.AssertExpectations(t)
//}
//
//func TestOpenLong_LeverageRatioLessThanSafetyFactor(t *testing.T) {
//	// Setup the mock checker
//	mockChecker := new(mocks.OpenLongChecker)
//
//	// Create an instance of Keeper with the mock checker
//	k := keeper.Keeper{
//		OpenLongChecker: mockChecker,
//	}
//
//	var (
//		ctx = sdk.Context{} // Mock or setup a context
//		msg = &types.MsgOpen{
//			Creator:         "",
//			Leverage:        math.LegacyNewDec(10),
//			Position:        types.Position_LONG,
//			TradingAsset:    "uatom",
//			Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000)),
//			TakeProfitPrice: sdk.NewDecFromInt(types.TakeProfitPriceDefault),
//		}
//		ammPool = ammtypes.Pool{
//			PoolId: uint64(42),
//			PoolAssets: []ammtypes.PoolAsset{
//				{
//					Token:  sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(10000)),
//					Weight: sdk.NewInt(50),
//				},
//				{
//					Token:  sdk.NewCoin("uatom", sdk.NewInt(10000)),
//					Weight: sdk.NewInt(50),
//				},
//			},
//		}
//		pool = types.Pool{
//			AmmPoolId: ammPool.PoolId,
//		}
//	)
//
//	// Mock behaviors
//	mockChecker.On("GetMaxLeverageParam", ctx).Return(msg.Leverage)
//	mockChecker.On("GetPool", ctx, ammPool.PoolId).Return(pool, true)
//	mockChecker.On("IsPoolEnabled", ctx, ammPool.PoolId).Return(true)
//	mockChecker.On("GetAmmPool", ctx, ammPool.PoolId, msg.TradingAsset).Return(ammPool, nil)
//
//	eta := math.LegacyNewDec(9)
//
//	mockChecker.On("CheckMinLiabilities", ctx, msg.Collateral, eta, ammPool, msg.TradingAsset, ptypes.BaseCurrency).Return(nil)
//
//	leveragedAmtTokenIn := sdk.NewCoin(msg.Collateral.Denom, math.NewInt(10000))
//	custodyAmount := math.NewInt(99)
//
//	mockChecker.On("EstimateSwap", ctx, leveragedAmtTokenIn, msg.TradingAsset, ammPool).Return(custodyAmount, nil)
//
//	mtp := types.NewMTP(msg.Creator, msg.Collateral.Denom, msg.TradingAsset, msg.Collateral.Denom, msg.TradingAsset, msg.Position, msg.Leverage, sdk.NewDecFromInt(types.TakeProfitPriceDefault), ammPool.PoolId)
//
//	mockChecker.On("Borrow", ctx, msg.Collateral.Amount, custodyAmount, mtp, &ammPool, &pool, eta, ptypes.BaseCurrency, false).Return(nil)
//	mockChecker.On("UpdatePoolHealth", ctx, &pool).Return(nil)
//	mockChecker.On("TakeInCustody", ctx, *mtp, &pool).Return(nil)
//
//	lr := math.LegacyNewDec(50)
//
//	mockChecker.On("UpdateMTPHealth", ctx, *mtp, ammPool, ptypes.BaseCurrency).Return(lr, nil)
//	mockChecker.On("GetSafetyFactor", ctx).Return(sdk.NewDec(100))
//
//	_, err := k.OpenLong(ctx, ammPool.PoolId, msg, ptypes.BaseCurrency, false)
//
//	// Expect an error indicating MTP is unhealthy
//	assert.True(t, errors.Is(err, types.ErrMTPUnhealthy))
//	mockChecker.AssertExpectations(t)
//}
//
//func TestOpenLong_Success(t *testing.T) {
//	// Setup the mock checker
//	mockChecker := new(mocks.OpenLongChecker)
//
//	// Create an instance of Keeper with the mock checker
//	k := keeper.Keeper{
//		OpenLongChecker: mockChecker,
//	}
//
//	var (
//		ctx = sdk.Context{} // Mock or setup a context
//		msg = &types.MsgOpen{
//			Creator:         "",
//			Leverage:        math.LegacyNewDec(10),
//			Position:        types.Position_LONG,
//			TradingAsset:    "uatom",
//			Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000)),
//			TakeProfitPrice: sdk.NewDecFromInt(types.TakeProfitPriceDefault),
//		}
//		ammPool = ammtypes.Pool{
//			PoolId: uint64(42),
//			PoolAssets: []ammtypes.PoolAsset{
//				{
//					Token:  sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(10000)),
//					Weight: sdk.NewInt(50),
//				},
//				{
//					Token:  sdk.NewCoin("uatom", sdk.NewInt(10000)),
//					Weight: sdk.NewInt(50),
//				},
//			},
//		}
//		pool = types.Pool{
//			AmmPoolId: ammPool.PoolId,
//		}
//	)
//
//	// Mock behaviors
//	mockChecker.On("GetMaxLeverageParam", ctx).Return(msg.Leverage)
//	mockChecker.On("GetPool", ctx, ammPool.PoolId).Return(pool, true)
//	mockChecker.On("IsPoolEnabled", ctx, ammPool.PoolId).Return(true)
//	mockChecker.On("GetAmmPool", ctx, ammPool.PoolId, msg.TradingAsset).Return(ammPool, nil)
//
//	eta := math.LegacyNewDec(9)
//
//	mockChecker.On("CheckMinLiabilities", ctx, msg.Collateral, eta, ammPool, msg.TradingAsset, ptypes.BaseCurrency).Return(nil)
//
//	leveragedAmtTokenIn := sdk.NewCoin(msg.Collateral.Denom, math.NewInt(10000))
//	custodyAmount := math.NewInt(99)
//
//	mockChecker.On("EstimateSwap", ctx, leveragedAmtTokenIn, msg.TradingAsset, ammPool).Return(custodyAmount, nil)
//
//	mtp := types.NewMTP(msg.Creator, msg.Collateral.Denom, msg.TradingAsset, msg.Collateral.Denom, msg.TradingAsset, msg.Position, msg.Leverage, sdk.NewDecFromInt(types.TakeProfitPriceDefault), ammPool.PoolId)
//
//	mockChecker.On("Borrow", ctx, msg.Collateral.Amount, custodyAmount, mtp, &ammPool, &pool, eta, ptypes.BaseCurrency, false).Return(nil)
//	mockChecker.On("UpdatePoolHealth", ctx, &pool).Return(nil)
//	mockChecker.On("TakeInCustody", ctx, *mtp, &pool).Return(nil)
//
//	lr := math.LegacyNewDec(50)
//
//	mockChecker.On("UpdateMTPHealth", ctx, *mtp, ammPool, ptypes.BaseCurrency).Return(lr, nil)
//
//	safetyFactor := math.LegacyNewDec(10)
//
//	mockChecker.On("GetSafetyFactor", ctx).Return(safetyFactor)
//
//	mockChecker.On("CalcMTPConsolidateCollateral", ctx, mtp, ptypes.BaseCurrency).Return(nil)
//	mockChecker.On("SetMTP", ctx, mtp).Return(nil)
//
//	_, err := k.OpenLong(ctx, ammPool.PoolId, msg, ptypes.BaseCurrency, false)
//	// Expect no error
//	assert.Nil(t, err)
//	mockChecker.AssertExpectations(t)
//}
//
//func TestOpenLong_BaseCurrency_Collateral(t *testing.T) {
//	app := simapp.InitElysTestApp(true)
//	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
//
//	mk, amm, oracle := app.PerpetualKeeper, app.AmmKeeper, app.OracleKeeper
//
//	// Setup coin prices
//	SetupStableCoinPrices(ctx, oracle)
//
//	// Set asset profile
//	app.AssetprofileKeeper.SetEntry(ctx, assetprofiletypes.Entry{
//		BaseDenom: ptypes.BaseCurrency,
//		Denom:     ptypes.BaseCurrency,
//		Decimals:  6,
//	})
//	app.AssetprofileKeeper.SetEntry(ctx, assetprofiletypes.Entry{
//		BaseDenom: ptypes.ATOM,
//		Denom:     ptypes.ATOM,
//		Decimals:  6,
//	})
//
//	// Generate 1 random account with 1000stake balanced
//	addr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(1000000000000))
//
//	// Create a pool
//	// Mint 100000USDC
//	usdcToken := []sdk.Coin{sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(200000000000))}
//	// Mint 100000ATOM
//	atomToken := []sdk.Coin{sdk.NewCoin(ptypes.ATOM, sdk.NewInt(200000000000))}
//
//	err := app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, usdcToken)
//	require.NoError(t, err)
//	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[0], usdcToken)
//	require.NoError(t, err)
//
//	err = app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, atomToken)
//	require.NoError(t, err)
//	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[0], atomToken)
//	require.NoError(t, err)
//
//	poolAssets := []ammtypes.PoolAsset{
//		{
//			Weight: sdk.NewInt(50),
//			Token:  sdk.NewCoin(ptypes.ATOM, sdk.NewInt(10000000000)),
//		},
//		{
//			Weight: sdk.NewInt(50),
//			Token:  sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(100000000000)),
//		},
//	}
//
//	argSwapFee := sdk.MustNewDecFromStr("0.0")
//	argExitFee := sdk.MustNewDecFromStr("0.0")
//
//	poolParams := &ammtypes.PoolParams{
//		SwapFee: argSwapFee,
//		ExitFee: argExitFee,
//	}
//
//	msg := ammtypes.NewMsgCreatePool(
//		addr[0].String(),
//		poolParams,
//		poolAssets,
//	)
//
//	// Create a ATOM+USDC pool
//	poolId, err := amm.CreatePool(ctx, msg)
//	require.NoError(t, err)
//	require.Equal(t, poolId, uint64(1))
//
//	pools := amm.GetAllPool(ctx)
//
//	// check length of pools
//	require.Equal(t, len(pools), 1)
//
//	// check block height
//	require.Equal(t, int64(0), ctx.BlockHeight())
//
//	pool, found := amm.GetPool(ctx, poolId)
//	require.Equal(t, found, true)
//
//	poolAddress := sdk.MustAccAddressFromBech32(pool.GetAddress())
//	require.NoError(t, err)
//
//	// Balance check before create a perpetual position
//	balances := app.BankKeeper.GetAllBalances(ctx, poolAddress)
//	require.Equal(t, balances.AmountOf(ptypes.BaseCurrency), sdk.NewInt(100000000000))
//	require.Equal(t, balances.AmountOf(ptypes.ATOM), sdk.NewInt(10000000000))
//
//	// Create a perpetual position open msg
//	msg2 := types.NewMsgOpen(
//		addr[0].String(),
//		types.Position_LONG,
//		sdk.NewDec(5),
//		ptypes.ATOM,
//		sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(100000000)),
//		sdk.NewDecFromInt(types.TakeProfitPriceDefault),
//	)
//
//	_, err = mk.Open(ctx, msg2, false)
//	require.NoError(t, err)
//
//	mtps := mk.GetAllMTPs(ctx)
//	require.Equal(t, len(mtps), 1)
//
//	balances = app.BankKeeper.GetAllBalances(ctx, poolAddress)
//	require.Equal(t, balances.AmountOf(ptypes.BaseCurrency), sdk.NewInt(100100000000))
//	require.Equal(t, balances.AmountOf(ptypes.ATOM), sdk.NewInt(10000000000))
//
//	_, found = mk.OpenLongChecker.GetPool(ctx, pool.PoolId)
//	require.Equal(t, found, true)
//
//	err = mk.InvariantCheck(ctx)
//	require.Equal(t, err, nil)
//
//	mtp := mtps[0]
//
//	// Check MTP
//	require.Equal(t, types.MTP{
//		Address:                        addr[0].String(),
//		CollateralAsset:                "uusdc",
//		TradingAsset:                   "uatom",
//		LiabilitiesAsset:               "uusdc",
//		CustodyAsset:                   "uatom",
//		Collateral:                     sdk.NewInt(100000000),
//		Liabilities:                    sdk.NewInt(400000000),
//		BorrowInterestPaidCollateral:   sdk.NewInt(0),
//		BorrowInterestPaidCustody:      sdk.NewInt(0),
//		BorrowInterestUnpaidCollateral: sdk.NewInt(0),
//		Custody:                        sdk.NewInt(49751243),
//		TakeProfitLiabilities:          sdk.NewInt(495049497),
//		TakeProfitCustody:              sdk.NewInt(49751243),
//		Leverage:                       sdk.NewDec(5),
//		MtpHealth:                      sdk.MustNewDecFromStr("1.249999982500000000"),
//		Position:                       types.Position_LONG,
//		Id:                             uint64(1),
//		AmmPoolId:                      uint64(1),
//		ConsolidateLeverage:            sdk.NewDec(4),
//		SumCollateral:                  sdk.NewInt(100000000),
//		TakeProfitPrice:                sdk.NewDecFromInt(types.TakeProfitPriceDefault),
//		TakeProfitBorrowRate:           sdk.MustNewDecFromStr("1.0"),
//		FundingFeePaidCollateral:       sdk.NewInt(0),
//		FundingFeePaidCustody:          sdk.NewInt(0),
//		FundingFeeReceivedCollateral:   sdk.NewInt(0),
//		FundingFeeReceivedCustody:      sdk.NewInt(0),
//		OpenPrice:                      sdk.MustNewDecFromStr("10.050000157785002477"),
//	}, mtp)
//}
//
//func TestOpenLong_ATOM_Collateral(t *testing.T) {
//	app := simapp.InitElysTestApp(true)
//	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
//
//	mk, amm, oracle := app.PerpetualKeeper, app.AmmKeeper, app.OracleKeeper
//
//	// Setup coin prices
//	SetupStableCoinPrices(ctx, oracle)
//
//	// Generate 1 random account with 1000stake balanced
//	addr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(1000000000000))
//
//	// Create a pool
//	// Mint 100000USDC
//	usdcToken := []sdk.Coin{sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(100000000000))}
//	// Mint 100000ATOM
//	atomToken := []sdk.Coin{sdk.NewCoin(ptypes.ATOM, sdk.NewInt(100000000000))}
//
//	err := app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, usdcToken)
//	require.NoError(t, err)
//	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[0], usdcToken)
//	require.NoError(t, err)
//
//	err = app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, atomToken)
//	require.NoError(t, err)
//	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[0], atomToken)
//	require.NoError(t, err)
//
//	poolAssets := []ammtypes.PoolAsset{
//		{
//			Weight: sdk.NewInt(50),
//			Token:  sdk.NewCoin(ptypes.ATOM, sdk.NewInt(1000000000)),
//		},
//		{
//			Weight: sdk.NewInt(50),
//			Token:  sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(10000000000)),
//		},
//	}
//
//	argSwapFee := sdk.MustNewDecFromStr("0.0")
//	argExitFee := sdk.MustNewDecFromStr("0.0")
//
//	poolParams := &ammtypes.PoolParams{
//		SwapFee: argSwapFee,
//		ExitFee: argExitFee,
//	}
//
//	msg := ammtypes.NewMsgCreatePool(
//		addr[0].String(),
//		poolParams,
//		poolAssets,
//	)
//
//	// Create a ATOM+USDC pool
//	poolId, err := amm.CreatePool(ctx, msg)
//	require.NoError(t, err)
//	require.Equal(t, poolId, uint64(1))
//
//	pools := amm.GetAllPool(ctx)
//
//	// check length of pools
//	require.Equal(t, len(pools), 1)
//
//	// check block height
//	require.Equal(t, int64(0), ctx.BlockHeight())
//
//	pool, found := amm.GetPool(ctx, poolId)
//	require.Equal(t, found, true)
//
//	poolAddress := sdk.MustAccAddressFromBech32(pool.GetAddress())
//	require.NoError(t, err)
//
//	// Balance check before create a perpetual position
//	balances := app.BankKeeper.GetAllBalances(ctx, poolAddress)
//	require.Equal(t, balances.AmountOf(ptypes.BaseCurrency), sdk.NewInt(10000000000))
//	require.Equal(t, balances.AmountOf(ptypes.ATOM), sdk.NewInt(1000000000))
//
//	// Create a perpetual position open msg
//	msg2 := types.NewMsgOpen(
//		addr[0].String(),
//		types.Position_LONG,
//		sdk.NewDec(5),
//		ptypes.ATOM,
//		sdk.NewCoin(ptypes.ATOM, sdk.NewInt(10000000)),
//		sdk.NewDecFromInt(types.TakeProfitPriceDefault),
//	)
//
//	_, err = mk.Open(ctx, msg2, false)
//	require.NoError(t, err)
//
//	mtps := mk.GetAllMTPs(ctx)
//	require.Equal(t, len(mtps), 1)
//
//	balances = app.BankKeeper.GetAllBalances(ctx, poolAddress)
//	require.Equal(t, balances.AmountOf(ptypes.BaseCurrency), sdk.NewInt(10000000000))
//	require.Equal(t, balances.AmountOf(ptypes.ATOM), sdk.NewInt(1010000000))
//
//	_, found = mk.OpenLongChecker.GetPool(ctx, pool.PoolId)
//	require.Equal(t, found, true)
//
//	err = mk.InvariantCheck(ctx)
//	require.Equal(t, err, nil)
//
//	mtp := mtps[0]
//
//	// Check MTP
//	require.Equal(t, types.MTP{
//		Address:                        addr[0].String(),
//		CollateralAsset:                "uatom",
//		TradingAsset:                   "uatom",
//		LiabilitiesAsset:               "uusdc",
//		CustodyAsset:                   "uatom",
//		Collateral:                     sdk.NewInt(10000000),
//		Liabilities:                    sdk.NewInt(416666667),
//		BorrowInterestPaidCollateral:   sdk.NewInt(0),
//		BorrowInterestPaidCustody:      sdk.NewInt(0),
//		BorrowInterestUnpaidCollateral: sdk.NewInt(0),
//		Custody:                        sdk.NewInt(50000000),
//		TakeProfitLiabilities:          sdk.NewInt(476190476),
//		TakeProfitCustody:              sdk.NewInt(50000000),
//		Leverage:                       sdk.NewDec(5),
//		MtpHealth:                      sdk.MustNewDecFromStr("1.263157894989473684"),
//		Position:                       types.Position_LONG,
//		Id:                             uint64(1),
//		AmmPoolId:                      uint64(1),
//		ConsolidateLeverage:            sdk.NewDec(4),
//		SumCollateral:                  sdk.NewInt(101010102),
//		TakeProfitPrice:                sdk.NewDecFromInt(types.TakeProfitPriceDefault),
//		TakeProfitBorrowRate:           sdk.MustNewDecFromStr("1.0"),
//		FundingFeePaidCollateral:       sdk.NewInt(0),
//		FundingFeePaidCustody:          sdk.NewInt(0),
//		FundingFeeReceivedCollateral:   sdk.NewInt(0),
//		FundingFeeReceivedCustody:      sdk.NewInt(0),
//		OpenPrice:                      sdk.MustNewDecFromStr("10.313531340000000000"),
//	}, mtp)
//}
