package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
	leveragelpmodulekeeper "github.com/elys-network/elys/v6/x/leveragelp/keeper"
	leveragelpmoduletypes "github.com/elys-network/elys/v6/x/leveragelp/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
	"github.com/stretchr/testify/require"
)

func (suite *PerpetualKeeperTestSuite) TestOpenEstimationByFinal_Long5XAtom100Usdc() {
	app := suite.app
	ctx := suite.ctx
	suite.SetupCoinPrices()

	mk, amm := app.PerpetualKeeper, app.AmmKeeper

	// Setup coin prices
	tradingAssetPrice, _, err := app.PerpetualKeeper.GetAssetPriceAndAssetUsdcDenomRatio(ctx, ptypes.ATOM)
	require.NoError(suite.T(), err)

	// Generate 1 random account with 1000stake balanced
	addr := suite.AddAccounts(1, nil)

	// Create a pool
	// Mint 100000USDC
	usdcToken := []sdk.Coin{sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000000000000))}
	// Mint 100000ATOM
	atomToken := []sdk.Coin{sdk.NewCoin(ptypes.ATOM, math.NewInt(1000000000000))}

	err = app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, usdcToken)
	require.NoError(suite.T(), err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[0], usdcToken)
	require.NoError(suite.T(), err)

	err = app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, atomToken)
	require.NoError(suite.T(), err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[0], atomToken)
	require.NoError(suite.T(), err)

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
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), poolId, uint64(1))

	pools := amm.GetAllPool(ctx)

	// check length of pools
	require.Equal(suite.T(), len(pools), 1)

	enablePoolMsg := leveragelpmoduletypes.MsgAddPool{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Pool: leveragelpmoduletypes.AddPool{
			AmmPoolId:   poolId,
			LeverageMax: math.LegacyMustNewDecFromStr("10"),
		},
	}
	_, err = leveragelpmodulekeeper.NewMsgServerImpl(*app.LeveragelpKeeper).AddPool(ctx, &enablePoolMsg)
	require.NoError(suite.T(), err)

	// First get regular open estimation
	regularRes, err := mk.OpenEstimation(ctx, &types.QueryOpenEstimationRequest{
		PoolId:          1,
		Position:        types.Position_LONG,
		Leverage:        math.LegacyMustNewDecFromStr("5.0"),
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(100_000_000)),
		Address:         "",
		TakeProfitPrice: tradingAssetPrice.MulInt64(3),
	})
	require.NoError(suite.T(), err)

	// Now get open estimation by final
	finalRes, err := mk.OpenEstimationByFinal(ctx, &types.QueryOpenEstimationByFinalRequest{
		PoolId:          1,
		Position:        types.Position_LONG,
		Leverage:        math.LegacyMustNewDecFromStr("5.0"),
		TradingAsset:    ptypes.ATOM,
		CollateralDenom: ptypes.BaseCurrency,
		FinalAmount:     regularRes.Custody,
		Address:         "",
		TakeProfitPrice: tradingAssetPrice.MulInt64(3),
	})
	require.NoError(suite.T(), err)

	// Compare results
	require.Equal(suite.T(), regularRes.Position, finalRes.Position)
	require.Equal(suite.T(), regularRes.TradingAsset, finalRes.TradingAsset)
	require.Equal(suite.T(), regularRes.Custody, finalRes.Custody)
	// Allow small difference in custody amount due to rounding
	diff := regularRes.Collateral.Amount.Sub(finalRes.Collateral.Amount)
	if diff.IsNegative() {
		diff = diff.Neg()
	}
	require.True(suite.T(), diff.LTE(math.NewInt(7)), "collateral amount difference should be <= 7")
	diff = regularRes.Liabilities.Amount.Sub(finalRes.Liabilities.Amount)
	if diff.IsNegative() {
		diff = diff.Neg()
	}
	require.True(suite.T(), diff.LTE(math.NewInt(7)), "liabilities amount difference should be <= 2")
	priceDiff := regularRes.OpenPrice.Sub(finalRes.OpenPrice)
	if priceDiff.IsNegative() {
		priceDiff = priceDiff.Neg()
	}
	require.True(suite.T(), priceDiff.LTE(math.LegacyMustNewDecFromStr("0.0001")), "open price difference should be <= 0.0001")
	require.Equal(suite.T(), regularRes.TakeProfitPrice, finalRes.TakeProfitPrice)
	priceDiff = regularRes.LiquidationPrice.Sub(finalRes.LiquidationPrice)
	if priceDiff.IsNegative() {
		priceDiff = priceDiff.Neg()
	}
	require.True(suite.T(), priceDiff.LTE(math.LegacyMustNewDecFromStr("0.0001")), "liquidation price difference should be <= 0.0001")
}

func (suite *PerpetualKeeperTestSuite) TestOpenEstimationByFinal_Short5XAtom100Usdc() {
	app := suite.app
	ctx := suite.ctx
	suite.SetupCoinPrices()

	mk, amm := app.PerpetualKeeper, app.AmmKeeper

	// Setup coin prices
	tradingAssetPrice, _, err := app.PerpetualKeeper.GetAssetPriceAndAssetUsdcDenomRatio(ctx, ptypes.ATOM)
	require.NoError(suite.T(), err)

	// Generate 1 random account with 1000stake balanced
	addr := suite.AddAccounts(1, nil)

	// Create a pool
	// Mint 100000USDC
	usdcToken := []sdk.Coin{sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000000000000))}
	// Mint 100000ATOM
	atomToken := []sdk.Coin{sdk.NewCoin(ptypes.ATOM, math.NewInt(1000000000000))}

	err = app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, usdcToken)
	require.NoError(suite.T(), err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[0], usdcToken)
	require.NoError(suite.T(), err)

	err = app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, atomToken)
	require.NoError(suite.T(), err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[0], atomToken)
	require.NoError(suite.T(), err)

	poolAssets := []ammtypes.PoolAsset{
		{
			Weight:                 math.NewInt(50),
			Token:                  sdk.NewCoin(ptypes.ATOM, math.NewInt(6000000000000)),
			ExternalLiquidityRatio: math.LegacyNewDec(2),
		},
		{
			Weight:                 math.NewInt(50),
			Token:                  sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000000000000)),
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
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), poolId, uint64(1))

	pools := amm.GetAllPool(ctx)

	// check length of pools
	require.Equal(suite.T(), len(pools), 1)

	enablePoolMsg := leveragelpmoduletypes.MsgAddPool{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Pool: leveragelpmoduletypes.AddPool{
			AmmPoolId:   poolId,
			LeverageMax: math.LegacyMustNewDecFromStr("10"),
		},
	}
	_, err = leveragelpmodulekeeper.NewMsgServerImpl(*app.LeveragelpKeeper).AddPool(ctx, &enablePoolMsg)
	require.NoError(suite.T(), err)

	// First get regular open estimation
	regularRes, err := mk.OpenEstimation(ctx, &types.QueryOpenEstimationRequest{
		PoolId:          1,
		Position:        types.Position_SHORT,
		Leverage:        math.LegacyMustNewDecFromStr("5.0"),
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(100_000_000)),
		Address:         "",
		TakeProfitPrice: tradingAssetPrice.QuoInt64(3),
	})
	require.NoError(suite.T(), err)

	// Now get open estimation by final
	finalRes, err := mk.OpenEstimationByFinal(ctx, &types.QueryOpenEstimationByFinalRequest{
		PoolId:          1,
		Position:        types.Position_SHORT,
		Leverage:        math.LegacyMustNewDecFromStr("5.0"),
		TradingAsset:    ptypes.ATOM,
		CollateralDenom: ptypes.BaseCurrency,
		FinalAmount:     regularRes.Liabilities,
		Address:         "",
		TakeProfitPrice: tradingAssetPrice.QuoInt64(3),
	})
	require.NoError(suite.T(), err)

	// Compare results
	require.Equal(suite.T(), regularRes.Position, finalRes.Position)
	require.Equal(suite.T(), regularRes.TradingAsset, finalRes.TradingAsset)
	diff := regularRes.Collateral.Amount.Sub(finalRes.Collateral.Amount)
	if diff.IsNegative() {
		diff = diff.Neg()
	}
	require.True(suite.T(), diff.LTE(math.NewInt(25)), "collateral amount difference should be <= 25")
	require.Equal(suite.T(), regularRes.Liabilities, finalRes.Liabilities)
	diff = regularRes.Custody.Amount.Sub(finalRes.Custody.Amount)
	if diff.IsNegative() {
		diff = diff.Neg()
	}
	require.True(suite.T(), diff.LTE(math.NewInt(200)), "custody amount difference should be <= 200")
}

func (suite *PerpetualKeeperTestSuite) TestOpenEstimationByFinal_Long5XAtom10Atom() {
	app := suite.app
	ctx := suite.ctx
	suite.SetupCoinPrices()

	mk, amm := app.PerpetualKeeper, app.AmmKeeper

	// Setup coin prices
	tradingAssetPrice, _, err := app.PerpetualKeeper.GetAssetPriceAndAssetUsdcDenomRatio(ctx, ptypes.ATOM)
	require.NoError(suite.T(), err)

	// Generate 1 random account with 1000stake balanced
	addr := suite.AddAccounts(1, nil)

	// Create a pool
	// Mint 100000USDC
	usdcToken := []sdk.Coin{sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000000000000))}
	// Mint 100000ATOM
	atomToken := []sdk.Coin{sdk.NewCoin(ptypes.ATOM, math.NewInt(1000000000000))}

	err = app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, usdcToken)
	require.NoError(suite.T(), err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[0], usdcToken)
	require.NoError(suite.T(), err)

	err = app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, atomToken)
	require.NoError(suite.T(), err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[0], atomToken)
	require.NoError(suite.T(), err)

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
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), poolId, uint64(1))

	pools := amm.GetAllPool(ctx)

	// check length of pools
	require.Equal(suite.T(), len(pools), 1)

	enablePoolMsg := leveragelpmoduletypes.MsgAddPool{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Pool: leveragelpmoduletypes.AddPool{
			AmmPoolId:   poolId,
			LeverageMax: math.LegacyMustNewDecFromStr("10"),
		},
	}
	_, err = leveragelpmodulekeeper.NewMsgServerImpl(*app.LeveragelpKeeper).AddPool(ctx, &enablePoolMsg)
	require.NoError(suite.T(), err)

	// First get regular open estimation
	regularRes, err := mk.OpenEstimation(ctx, &types.QueryOpenEstimationRequest{
		PoolId:          1,
		Position:        types.Position_LONG,
		Leverage:        math.LegacyMustNewDecFromStr("5.0"),
		TradingAsset:    ptypes.ATOM,
		Collateral:      sdk.NewCoin(ptypes.ATOM, math.NewInt(10_000_000)),
		Address:         "",
		TakeProfitPrice: tradingAssetPrice.MulInt64(3),
	})
	require.NoError(suite.T(), err)

	// Now get open estimation by final
	finalRes, err := mk.OpenEstimationByFinal(ctx, &types.QueryOpenEstimationByFinalRequest{
		PoolId:          1,
		Position:        types.Position_LONG,
		Leverage:        math.LegacyMustNewDecFromStr("5.0"),
		TradingAsset:    ptypes.ATOM,
		CollateralDenom: ptypes.ATOM,
		FinalAmount:     regularRes.Custody,
		Address:         "",
		TakeProfitPrice: tradingAssetPrice.MulInt64(3),
	})
	require.NoError(suite.T(), err)

	// Compare results
	require.Equal(suite.T(), regularRes.Position, finalRes.Position)
	require.Equal(suite.T(), regularRes.TradingAsset, finalRes.TradingAsset)
	// Allow small difference in custody amount due to rounding
	diff := regularRes.Custody.Amount.Sub(finalRes.Custody.Amount)
	if diff.IsNegative() {
		diff = diff.Neg()
	}
	require.True(suite.T(), diff.LTE(math.NewInt(2)), "custody amount difference should be <= 2")
	require.Equal(suite.T(), regularRes.Collateral, finalRes.Collateral)
	require.Equal(suite.T(), regularRes.Liabilities, finalRes.Liabilities)
	require.Equal(suite.T(), regularRes.OpenPrice, finalRes.OpenPrice)
	require.Equal(suite.T(), regularRes.TakeProfitPrice, finalRes.TakeProfitPrice)
	require.Equal(suite.T(), regularRes.LiquidationPrice, finalRes.LiquidationPrice)
	require.Equal(suite.T(), regularRes.EstimatedPnl, finalRes.EstimatedPnl)
	require.Equal(suite.T(), regularRes.AvailableLiquidity, finalRes.AvailableLiquidity)
	require.Equal(suite.T(), regularRes.Slippage, finalRes.Slippage)
	require.Equal(suite.T(), regularRes.PriceImpact, finalRes.PriceImpact)
	require.Equal(suite.T(), regularRes.BorrowInterestRate, finalRes.BorrowInterestRate)
	require.Equal(suite.T(), regularRes.FundingRate, finalRes.FundingRate)
	require.Equal(suite.T(), regularRes.WeightBreakingFee, finalRes.WeightBreakingFee)
}
