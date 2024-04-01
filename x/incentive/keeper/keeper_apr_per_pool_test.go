package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/incentive/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestAPRCalculationPerPool(t *testing.T) {
	app := simapp.InitElysTestApp(initChain)
	ctx := app.BaseApp.NewContext(initChain, tmproto.Header{})

	ik, amm, oracle := app.IncentiveKeeper, app.AmmKeeper, app.OracleKeeper

	// Setup coin prices
	SetupStableCoinPrices(ctx, oracle)

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(1000000))

	// Create a pool
	// Mint 100000USDC
	usdcToken := sdk.NewCoins(sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(100000)))

	err := app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, usdcToken)
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

	// Create a Elys+USDC pool
	poolId, err := amm.CreatePool(ctx, msg)
	require.NoError(t, err)
	require.Equal(t, poolId, uint64(1))

	pools := amm.GetAllPool(ctx)

	// check length of pools
	require.Equal(t, len(pools), 1)

	// check block height
	require.Equal(t, int64(0), ctx.BlockHeight())

	params := ik.GetParams(ctx)
	params.DistributionInterval = 10
	ik.SetParams(ctx, params)

	lpIncentive := types.IncentiveInfo{
		// reward amount in eden for 1 year
		EdenAmountPerYear: sdk.NewInt(1000000000),
		// starting block height of the distribution
		DistributionStartBlock: sdk.NewInt(1),
		// distribution duration - block number per year
		TotalBlocksPerYear: sdk.NewInt(10000),
		// maximum eden allocation per day that won't exceed 30% apr
		MaxEdenPerAllocation: sdk.NewInt(100),
		// current epoch in block number
		CurrentEpochInBlocks: sdk.NewInt(1),
	}

	ctx = ctx.WithBlockHeight(params.DistributionInterval)
	ik.UpdateLPRewardsUnclaimed(ctx, lpIncentive)

	// Get pool info from incentive param
	poolInfo, found := ik.GetPoolInfo(ctx, poolId)
	require.Equal(t, found, true)
	require.Equal(t, poolInfo.EdenApr.String(), "0.182317682317682318")

	// Get dex rewards per pool
	revenueAddress := ammtypes.NewPoolRevenueAddress(poolId)

	// Feed dex rewards
	usdcToken = sdk.NewCoins(sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(1000)))
	err = app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, usdcToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, revenueAddress, usdcToken)
	require.NoError(t, err)

	// 1 week later.
	ctx = ctx.WithBlockHeight(params.DistributionInterval * ptypes.DaysPerWeek)
	poolInfo.NumBlocks = sdk.NewInt(ctx.BlockHeight())
	ik.SetPoolInfo(ctx, poolId, poolInfo)

	ik.UpdateLPRewardsUnclaimed(ctx, lpIncentive)
	poolInfo, found = ik.GetPoolInfo(ctx, poolId)
	require.Equal(t, found, true)
	require.Equal(t, poolInfo.EdenApr.String(), "0.182317682317682318")
}
