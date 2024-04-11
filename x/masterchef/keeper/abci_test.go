package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	simapp "github.com/elys-network/elys/app"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	ctypes "github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	tokenomicskeeper "github.com/elys-network/elys/x/tokenomics/keeper"
	tokenomicstypes "github.com/elys-network/elys/x/tokenomics/types"
	"github.com/stretchr/testify/require"
)

func TestABCI_EndBlocker(t *testing.T) {
	app, genAccount, _ := simapp.InitElysTestAppWithGenAccount()
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	ik := app.IncentiveKeeper

	var committed sdk.Coins
	var unclaimed sdk.Coins

	// Prepare unclaimed tokens
	uedenToken := sdk.NewCoin(ptypes.Eden, sdk.NewInt(2000))
	uedenBToken := sdk.NewCoin(ptypes.EdenB, sdk.NewInt(2000))
	unclaimed = unclaimed.Add(uedenToken, uedenBToken)

	// Mint coins
	err := app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, unclaimed)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ctypes.ModuleName, genAccount, unclaimed)
	require.NoError(t, err)

	// Add testing commitment
	simapp.AddTestCommitment(app, ctx, genAccount, committed)
	ik.EndBlocker(ctx)

	// Get elys staked
	elysStaked := ik.GetElysStaked(ctx, genAccount.String())
	require.Equal(t, elysStaked, sdk.DefaultPowerReduction)

	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	srv := tokenomicskeeper.NewMsgServerImpl(app.TokenomicsKeeper)

	expected := &tokenomicstypes.MsgCreateTimeBasedInflation{
		Authority:        authority,
		StartBlockHeight: uint64(1),
		EndBlockHeight:   uint64(6307200),
		Inflation: &tokenomicstypes.InflationEntry{
			LmRewards:         9999999,
			IcsStakingRewards: 9999999,
			CommunityFund:     9999999,
			StrategicReserve:  9999999,
			TeamTokensVested:  9999999,
		},
	}

	wctx := sdk.WrapSDKContext(ctx)
	_, err = srv.CreateTimeBasedInflation(wctx, expected)
	require.NoError(t, err)

	expected = &tokenomicstypes.MsgCreateTimeBasedInflation{
		Authority:        authority,
		StartBlockHeight: uint64(6307201),
		EndBlockHeight:   uint64(12614401),
		Inflation: &tokenomicstypes.InflationEntry{
			LmRewards:         9999999,
			IcsStakingRewards: 9999999,
			CommunityFund:     9999999,
			StrategicReserve:  9999999,
			TeamTokensVested:  9999999,
		},
	}
	_, err = srv.CreateTimeBasedInflation(wctx, expected)
	require.NoError(t, err)

	// Set tokenomics params
	listTimeBasdInflations := app.TokenomicsKeeper.GetAllTimeBasedInflation(ctx)

	// After the first year
	ctx = ctx.WithBlockHeight(1)
	paramSet := ik.ProcessUpdateIncentiveParams(ctx)
	require.Equal(t, paramSet, true)

	// Check if the params are correctly set
	params := ik.GetParams(ctx)
	require.NotNil(t, params.StakeIncentives)
	require.NotNil(t, params.LpIncentives)

	require.Equal(t, params.StakeIncentives.EdenAmountPerYear, sdk.NewInt(int64(listTimeBasdInflations[0].Inflation.IcsStakingRewards)))
	require.Equal(t, params.LpIncentives.EdenAmountPerYear, sdk.NewInt(int64(listTimeBasdInflations[0].Inflation.LmRewards)))

	// After the first year
	ctx = ctx.WithBlockHeight(6307210)

	// Incentive param should be empty
	stakerEpoch, _ := ik.IsStakerRewardsDistributionEpoch(ctx)
	params = ik.GetParams(ctx)
	require.Equal(t, stakerEpoch, false)
	require.Nil(t, params.StakeIncentives)

	// Incentive param should be empty
	lpEpoch, _ := ik.IsLPRewardsDistributionEpoch(ctx)
	params = ik.GetParams(ctx)
	require.Equal(t, lpEpoch, false)
	require.Nil(t, params.LpIncentives)

	// After reading tokenomics again
	paramSet = ik.ProcessUpdateIncentiveParams(ctx)
	require.Equal(t, paramSet, true)

	// Check params
	_, stakeIncentive := ik.IsStakerRewardsDistributionEpoch(ctx)
	params = ik.GetParams(ctx)
	require.Equal(t, stakeIncentive.EdenAmountPerYear, sdk.NewInt(int64(listTimeBasdInflations[0].Inflation.IcsStakingRewards)))

	// Check params
	_, lpIncentive := ik.IsLPRewardsDistributionEpoch(ctx)
	params = ik.GetParams(ctx)
	require.Equal(t, lpIncentive.EdenAmountPerYear, sdk.NewInt(int64(listTimeBasdInflations[0].Inflation.IcsStakingRewards)))
}

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

	err = app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, usdcToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[0], usdcToken)
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

	msg := ammtypes.NewMsgCreatePool(
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

	err := bk.MintCoins(ctx, ammtypes.ModuleName, usdcToken)
	require.NoError(t, err)
	err = bk.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[0], usdcToken)
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

	msg := ammtypes.NewMsgCreatePool(
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

	err = app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, usdcToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[1], usdcToken)
	require.NoError(t, err)

	// Mint uatom
	atomToken := sdk.NewCoins(sdk.NewCoin(ptypes.ATOM, sdk.NewInt(200000)))
	err = bk.MintCoins(ctx, ammtypes.ModuleName, atomToken)
	require.NoError(t, err)
	err = bk.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[1], atomToken)
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

	msg = ammtypes.NewMsgCreatePool(
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
	err = app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, usdcRevToken1)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, revenueAddress1, usdcRevToken1)
	require.NoError(t, err)

	// Fill in pool #2 revenue wallet
	revenueAddress2 := ammtypes.NewPoolRevenueAddress(1)
	usdcRevToken2 := sdk.NewCoins(sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(3000)))
	err = app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, usdcRevToken2)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, revenueAddress2, usdcRevToken2)
	require.NoError(t, err)

	// Collect revenue
	collectedAmt, rewardForLpsAmt, _ := ik.CollectDEXRevenue(ctx)

	// check block height
	require.Equal(t, int64(0), ctx.BlockHeight())

	// It should be 3000=1000+2000 usdc
	require.Equal(t, collectedAmt, sdk.Coins{sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(3000))})
	// It should be 1950=3000*0.65 usdc
	require.Equal(t, rewardForLpsAmt, sdk.DecCoins{sdk.NewDecCoin(ptypes.BaseCurrency, sdk.NewInt(1800))})
}
