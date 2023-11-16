package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/amm/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestCollectGasFeesToIncentiveModule(t *testing.T) {
	app := simapp.InitElysTestApp(initChain)
	ctx := app.BaseApp.NewContext(initChain, tmproto.Header{})

	ik, bk, amm, oracle := app.IncentiveKeeper, app.BankKeeper, app.AmmKeeper, app.OracleKeeper

	// Setup coin prices
	SetupStableCoinPrices(ctx, oracle)

	// Collect gas fees
	collectedAmt := ik.CollectGasFeesToIncentiveModule(ctx, ptypes.BaseCurrency)

	// rewards should be zero
	require.True(t, collectedAmt.IsZero())

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(1000000))
	transferAmt := sdk.NewCoin(ptypes.Elys, sdk.NewInt(100))

	// Deposit 100elys to FeeCollectorName wallet
	err := bk.SendCoinsFromAccountToModule(ctx, addr[0], authtypes.FeeCollectorName, sdk.NewCoins(transferAmt))
	require.NoError(t, err)

	// Create a pool
	// Mint 100000USDC
	usdcToken := sdk.NewCoins(sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(100000)))

	err = app.BankKeeper.MintCoins(ctx, types.ModuleName, usdcToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr[0], usdcToken)
	require.NoError(t, err)

	poolAssets := []ammtypes.PoolAsset{
		{
			Weight: sdk.NewInt(50),
			Token:  sdk.NewCoin(ptypes.Elys, sdk.NewInt(100000)),
		},
		{
			Weight: sdk.NewInt(50),
			Token:  sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(10000)),
		},
	}

	argSwapFee := sdk.MustNewDecFromStr("0.1")
	argExitFee := sdk.MustNewDecFromStr("0.1")

	poolParams := &ammtypes.PoolParams{
		SwapFee: argSwapFee,
		ExitFee: argExitFee,
	}

	msg := types.NewMsgCreatePool(
		addr[0].String(),
		poolParams,
		poolAssets,
	)

	// Create a Elys+USDC pool
	poolId, err := amm.CreatePool(ctx, msg)
	require.NoError(t, err)
	require.Equal(t, poolId, uint64(1))

	pools := amm.GetAllPool(ctx)

	// check length of pools
	require.Equal(t, len(pools), 1)

	// check block height
	require.Equal(t, int64(0), ctx.BlockHeight())

	// Collect gas fees again
	collectedAmt = ik.CollectGasFeesToIncentiveModule(ctx, ptypes.BaseCurrency)

	// check block height
	require.Equal(t, int64(0), ctx.BlockHeight())

	// It should be 9 usdc
	require.Equal(t, collectedAmt, sdk.Coins{sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(9))})
}

func TestCollectDEXRevenueToIncentiveModule(t *testing.T) {
	app := simapp.InitElysTestApp(initChain)
	ctx := app.BaseApp.NewContext(initChain, tmproto.Header{})

	ik, bk, amm, oracle := app.IncentiveKeeper, app.BankKeeper, app.AmmKeeper, app.OracleKeeper

	// Recalculate total committed info
	ik.UpdateTotalCommitmentInfo(ctx, ptypes.BaseCurrency)

	// Setup coin prices
	SetupStableCoinPrices(ctx, oracle)

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(1000000))

	// Create 2 pools

	// #######################
	// ####### POOL 1 ########
	// Mint 100000USDC
	usdcToken := sdk.NewCoins(sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(100000)))

	err := bk.MintCoins(ctx, types.ModuleName, usdcToken)
	require.NoError(t, err)
	err = bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr[0], usdcToken)
	require.NoError(t, err)

	poolAssets := []ammtypes.PoolAsset{
		{
			Weight: sdk.NewInt(50),
			Token:  sdk.NewCoin(ptypes.Elys, sdk.NewInt(100000)),
		},
		{
			Weight: sdk.NewInt(50),
			Token:  sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(10000)),
		},
	}

	argSwapFee := sdk.MustNewDecFromStr("0.1")
	argExitFee := sdk.MustNewDecFromStr("0.1")

	poolParams := &ammtypes.PoolParams{
		SwapFee: argSwapFee,
		ExitFee: argExitFee,
	}

	msg := types.NewMsgCreatePool(
		addr[0].String(),
		poolParams,
		poolAssets,
	)

	// Create a Elys+USDC pool
	poolId, err := amm.CreatePool(ctx, msg)
	require.NoError(t, err)
	require.Equal(t, poolId, uint64(1))

	// ####### POOL 2 ########
	// ATOM+USDC pool
	// Mint uusdc
	usdcToken = sdk.NewCoins(sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(200000)))

	err = app.BankKeeper.MintCoins(ctx, types.ModuleName, usdcToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr[1], usdcToken)
	require.NoError(t, err)

	// Mint uatom
	atomToken := sdk.NewCoins(sdk.NewCoin(ptypes.ATOM, sdk.NewInt(200000)))
	err = bk.MintCoins(ctx, types.ModuleName, atomToken)
	require.NoError(t, err)
	err = bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr[1], atomToken)
	require.NoError(t, err)

	poolAssets2 := []ammtypes.PoolAsset{
		{
			Weight: sdk.NewInt(50),
			Token:  sdk.NewCoin(ptypes.ATOM, sdk.NewInt(150000)),
		},
		{
			Weight: sdk.NewInt(50),
			Token:  sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(10000)),
		},
	}

	msg = types.NewMsgCreatePool(
		addr[1].String(),
		poolParams,
		poolAssets2,
	)

	// Create a ATOM+USDC pool
	poolId, err = amm.CreatePool(ctx, msg)
	require.NoError(t, err)
	require.Equal(t, poolId, uint64(2))

	pools := amm.GetAllPool(ctx)

	// check length of pools
	require.Equal(t, len(pools), 2)

	// check block height
	require.Equal(t, int64(0), ctx.BlockHeight())

	// Fill in pool #1 revenue wallet
	revenueAddress1 := ammtypes.NewPoolRevenueAddress(0)
	usdcRevToken1 := sdk.NewCoins(sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(1000)))
	err = app.BankKeeper.MintCoins(ctx, types.ModuleName, usdcRevToken1)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, revenueAddress1, usdcRevToken1)
	require.NoError(t, err)

	// Fill in pool #2 revenue wallet
	revenueAddress2 := ammtypes.NewPoolRevenueAddress(1)
	usdcRevToken2 := sdk.NewCoins(sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(3000)))
	err = app.BankKeeper.MintCoins(ctx, types.ModuleName, usdcRevToken2)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, revenueAddress2, usdcRevToken2)
	require.NoError(t, err)

	// Collect revenue
	collectedAmt, rewardForLpsAmt := ik.CollectDEXRevenue(ctx)

	// check block height
	require.Equal(t, int64(0), ctx.BlockHeight())

	// It should be 3000=1000+2000 usdc
	require.Equal(t, collectedAmt, sdk.Coins{sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(3000))})
	// It should be 1950=3000*0.65 usdc
	require.Equal(t, rewardForLpsAmt, sdk.DecCoins{sdk.NewDecCoin(ptypes.BaseCurrency, sdk.NewInt(1950))})
}
