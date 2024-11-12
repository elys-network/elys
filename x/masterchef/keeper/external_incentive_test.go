package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	simapp "github.com/elys-network/elys/app"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	masterchefkeeper "github.com/elys-network/elys/x/masterchef/keeper"
	"github.com/elys-network/elys/x/masterchef/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	tokenomicskeeper "github.com/elys-network/elys/x/tokenomics/keeper"
	tokenomicstypes "github.com/elys-network/elys/x/tokenomics/types"
	"github.com/stretchr/testify/require"
)

func TestExternalIncentive(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

	externalIncentives := []types.ExternalIncentive{
		{
			Id:             0,
			RewardDenom:    "reward1",
			PoolId:         1,
			FromBlock:      0,
			ToBlock:        100,
			AmountPerBlock: sdkmath.OneInt(),
			Apr:            sdkmath.LegacyZeroDec(),
		},
		{
			Id:             1,
			RewardDenom:    "reward1",
			PoolId:         1,
			FromBlock:      0,
			ToBlock:        100,
			AmountPerBlock: sdkmath.OneInt(),
			Apr:            sdkmath.LegacyZeroDec(),
		},
		{
			Id:             2,
			RewardDenom:    "reward1",
			PoolId:         2,
			FromBlock:      0,
			ToBlock:        100,
			AmountPerBlock: sdkmath.OneInt(),
			Apr:            sdkmath.LegacyZeroDec(),
		},
	}
	for _, externalIncentive := range externalIncentives {
		app.MasterchefKeeper.SetExternalIncentive(ctx, externalIncentive)
	}
	for _, externalIncentive := range externalIncentives {
		info, found := app.MasterchefKeeper.GetExternalIncentive(ctx, externalIncentive.Id)
		require.True(t, found)
		require.Equal(t, info, externalIncentive)
	}
	externalIncentivesStored := app.MasterchefKeeper.GetAllExternalIncentives(ctx)
	require.Len(t, externalIncentivesStored, 3)

	app.MasterchefKeeper.RemoveExternalIncentive(ctx, externalIncentives[0].Id)
	externalIncentivesStored = app.MasterchefKeeper.GetAllExternalIncentives(ctx)
	require.Len(t, externalIncentivesStored, 2)
}

// Test USDC reward as external and via dex collection
func TestUSDCExternalIncentive(t *testing.T) {
	app, _, _ := simapp.InitElysTestAppWithGenAccount(t)
	ctx := app.BaseApp.NewContext(true)

	mk, amm, oracle := app.MasterchefKeeper, app.AmmKeeper, app.OracleKeeper

	// Setup coin prices
	SetupStableCoinPrices(ctx, oracle)

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
		Description: "Description",
	}

	wctx := sdk.WrapSDKContext(ctx)
	_, err := srv.CreateTimeBasedInflation(wctx, expected)
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
		Description: "Description",
	}
	_, err = srv.CreateTimeBasedInflation(wctx, expected)
	require.NoError(t, err)

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 2, sdkmath.NewInt(10000000000))

	// Create a pool
	// Mint 100000USDC
	usdcToken := sdk.NewCoins(sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(10000000000)))

	err = app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, usdcToken.MulInt(sdkmath.NewInt(2)))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[0], usdcToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[1], usdcToken)
	require.NoError(t, err)

	usdcToken = sdk.NewCoins(sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(100000000000)))
	err = app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, usdcToken.MulInt(sdkmath.NewInt(2)))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[0], usdcToken)
	require.NoError(t, err)

	poolAssets := []ammtypes.PoolAsset{
		{
			Weight: sdkmath.NewInt(50),
			Token:  sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(10000000)),
		},
		{
			Weight: sdkmath.NewInt(50),
			Token:  sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(10000000)),
		},
	}

	argSwapFee := sdkmath.LegacyMustNewDecFromStr("0.0")
	argExitFee := sdkmath.LegacyMustNewDecFromStr("0.0")

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

	_, _, err = amm.ExitPool(ctx, addr[0], pools[0].PoolId, math.NewIntWithDecimal(1, 21), sdk.NewCoins(), "", false)
	require.NoError(t, err)

	// new user join pool with same shares
	share := ammtypes.InitPoolSharesSupply.Mul(math.NewIntWithDecimal(1, 5))
	t.Log(mk.GetPoolTotalCommit(ctx, pools[0].PoolId))
	require.Equal(t, mk.GetPoolTotalCommit(ctx, pools[0].PoolId).String(), "10002000000000000000000000")
	require.Equal(t, mk.GetPoolBalance(ctx, pools[0].PoolId, addr[0]).String(), "10000000000000000000000000")
	_, _, err = amm.JoinPoolNoSwap(ctx, addr[1], pools[0].PoolId, share, sdk.NewCoins(sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(10000000)), sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(10000000))))
	require.NoError(t, err)
	require.Equal(t, mk.GetPoolTotalCommit(ctx, pools[0].PoolId).String(), "20002000000000000000000000")
	require.Equal(t, mk.GetPoolBalance(ctx, pools[0].PoolId, addr[1]), share)

	atomToken := sdk.NewCoins(sdk.NewCoin("uatom", math.NewIntWithDecimal(100000000, 6)))
	err = app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, atomToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[0], atomToken)
	require.NoError(t, err)

	masterchefSrv := masterchefkeeper.NewMsgServerImpl(app.MasterchefKeeper)
	_, err = masterchefSrv.AddExternalRewardDenom(sdk.WrapSDKContext(ctx), &types.MsgAddExternalRewardDenom{
		Authority:   app.GovKeeper.GetAuthority(),
		RewardDenom: ptypes.BaseCurrency,
		MinAmount:   sdkmath.OneInt(),
		Supported:   true,
	})
	require.NoError(t, err)
	_, err = masterchefSrv.AddExternalIncentive(sdk.WrapSDKContext(ctx), &types.MsgAddExternalIncentive{
		Sender:         addr[0].String(),
		RewardDenom:    ptypes.BaseCurrency,
		PoolId:         pools[0].PoolId,
		AmountPerBlock: math.NewIntWithDecimal(100, 6),
		FromBlock:      0,
		ToBlock:        1000,
	})
	require.NoError(t, err)

	// Fill in pool revenue wallet
	revenueAddress1 := ammtypes.NewPoolRevenueAddress(1)
	usdcRevToken1 := sdk.NewCoins(sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(100000)))
	err = app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, usdcRevToken1)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, revenueAddress1, usdcRevToken1)
	require.NoError(t, err)

	// check rewards after 100 block
	for i := 1; i <= 100; i++ {
		mk.EndBlocker(ctx)
		ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	}

	require.Equal(t, ctx.BlockHeight(), int64(100))
	poolRewardInfo, _ := app.MasterchefKeeper.GetPoolRewardInfo(ctx, pools[0].PoolId, ptypes.BaseCurrency)
	require.Equal(t, poolRewardInfo.LastUpdatedBlock, uint64(99))

	res, err := mk.UserPendingReward(ctx, &types.QueryUserPendingRewardRequest{
		User: addr[0].String(),
	})
	require.NoError(t, err)
	require.Equal(t, res.TotalRewards[0].Amount.String(), "4949545046")
	res, err = mk.UserPendingReward(ctx, &types.QueryUserPendingRewardRequest{
		User: addr[1].String(),
	})
	require.NoError(t, err)
	require.Equal(t, res.TotalRewards[0].Amount.String(), "4949545046")

	prevUSDCBal := app.BankKeeper.GetBalance(ctx, addr[1], ptypes.BaseCurrency)

	// check rewards claimed
	_, err = masterchefSrv.ClaimRewards(sdk.WrapSDKContext(ctx), &types.MsgClaimRewards{
		Sender:  addr[0].String(),
		PoolIds: []uint64{pools[0].PoolId},
	})
	require.NoError(t, err)
	_, err = masterchefSrv.ClaimRewards(sdk.WrapSDKContext(ctx), &types.MsgClaimRewards{
		Sender:  addr[1].String(),
		PoolIds: []uint64{pools[0].PoolId},
	})
	require.NoError(t, err)

	curUSDCBal := app.BankKeeper.GetBalance(ctx, addr[1], ptypes.BaseCurrency)
	amount, _ := sdkmath.NewIntFromString("4949545046")
	require.Equal(t, curUSDCBal.Amount.String(), prevUSDCBal.Amount.Add(amount).String())

	// no pending rewards
	res, err = mk.UserPendingReward(ctx, &types.QueryUserPendingRewardRequest{
		User: addr[0].String(),
	})
	require.NoError(t, err)
	require.Len(t, res.TotalRewards, 0)
	res, err = mk.UserPendingReward(ctx, &types.QueryUserPendingRewardRequest{
		User: addr[1].String(),
	})
	require.NoError(t, err)
	require.Len(t, res.TotalRewards, 0)
}
