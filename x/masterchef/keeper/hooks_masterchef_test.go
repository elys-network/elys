package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/ed25519"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	simapp "github.com/elys-network/elys/app"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	masterchefkeeper "github.com/elys-network/elys/x/masterchef/keeper"
	"github.com/elys-network/elys/x/masterchef/types"
	oraclekeeper "github.com/elys-network/elys/x/oracle/keeper"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	tokenomicskeeper "github.com/elys-network/elys/x/tokenomics/keeper"
	tokenomicstypes "github.com/elys-network/elys/x/tokenomics/types"
	"github.com/stretchr/testify/require"
)

func SetupStableCoinPrices(ctx sdk.Context, oracle oraclekeeper.Keeper) {
	// prices set for USDT and USDC
	provider := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	oracle.SetAssetInfo(ctx, oracletypes.AssetInfo{
		Denom:   ptypes.BaseCurrency,
		Display: "USDC",
		Decimal: 6,
	})
	oracle.SetAssetInfo(ctx, oracletypes.AssetInfo{
		Denom:   "uusdt",
		Display: "USDT",
		Decimal: 6,
	})
	oracle.SetAssetInfo(ctx, oracletypes.AssetInfo{
		Denom:   ptypes.Elys,
		Display: "ELYS",
		Decimal: 6,
	})
	oracle.SetAssetInfo(ctx, oracletypes.AssetInfo{
		Denom:   ptypes.ATOM,
		Display: "ATOM",
		Decimal: 6,
	})

	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "USDC",
		Price:     sdk.NewDec(1000000),
		Source:    "elys",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "USDT",
		Price:     sdk.NewDec(1000000),
		Source:    "elys",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "ELYS",
		Price:     sdk.NewDec(100),
		Source:    "elys",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "ATOM",
		Price:     sdk.NewDec(100),
		Source:    "atom",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
}

func TestHookMasterchef(t *testing.T) {
	app, _, _ := simapp.InitElysTestAppWithGenAccount()
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

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
	}
	_, err = srv.CreateTimeBasedInflation(wctx, expected)
	require.NoError(t, err)

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(100000000))

	// Create a pool
	// Mint 100000USDC
	usdcToken := sdk.NewCoins(sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(100000000)))

	err = app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, usdcToken.MulInt(sdk.NewInt(2)))
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[0], usdcToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[1], usdcToken)
	require.NoError(t, err)

	poolAssets := []ammtypes.PoolAsset{
		{
			Weight: sdk.NewInt(50),
			Token:  sdk.NewCoin(ptypes.Elys, sdk.NewInt(100)),
		},
		{
			Weight: sdk.NewInt(50),
			Token:  sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(100)),
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

	_, err = amm.ExitPool(ctx, addr[0], pools[0].PoolId, math.NewIntWithDecimal(1, 16), sdk.NewCoins(), "")
	require.NoError(t, err)

	// new user join pool with same shares
	share := ammtypes.InitPoolSharesSupply
	require.Equal(t, mk.GetPoolTotalSupply(ctx, pools[0].PoolId), ammtypes.InitPoolSharesSupply)
	require.Equal(t, mk.GetPoolBalance(ctx, pools[0].PoolId, addr[0].String()), share)
	_, _, err = amm.JoinPoolNoSwap(ctx, addr[1], pools[0].PoolId, share, sdk.NewCoins(sdk.NewCoin(ptypes.Elys, sdk.NewInt(10000)), sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(10000))))
	require.NoError(t, err)
	require.Equal(t, mk.GetPoolTotalSupply(ctx, pools[0].PoolId), ammtypes.InitPoolSharesSupply.MulRaw(2))
	require.Equal(t, mk.GetPoolBalance(ctx, pools[0].PoolId, addr[1].String()), share)

	atomToken := sdk.NewCoins(sdk.NewCoin("uatom", math.NewIntWithDecimal(100000000, 6)))
	err = app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, atomToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[0], atomToken)
	require.NoError(t, err)
	// external reward distribute
	masterchefSrv := masterchefkeeper.NewMsgServerImpl(app.MasterchefKeeper)
	_, err = masterchefSrv.AddExternalRewardDenom(sdk.WrapSDKContext(ctx), &types.MsgAddExternalRewardDenom{
		Authority:   app.GovKeeper.GetAuthority(),
		RewardDenom: "uatom",
		MinAmount:   sdk.OneInt(),
		Supported:   true,
	})
	require.NoError(t, err)
	_, err = masterchefSrv.AddExternalIncentive(sdk.WrapSDKContext(ctx), &types.MsgAddExternalIncentive{
		Sender:         addr[0].String(),
		RewardDenom:    "uatom",
		PoolId:         pools[0].PoolId,
		AmountPerBlock: math.NewIntWithDecimal(100, 6),
		FromBlock:      0,
		ToBlock:        1000,
	})
	require.NoError(t, err)

	// check rewards after 100 block
	for i := 1; i <= 100; i++ {
		mk.EndBlocker(ctx)
		ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	}

	require.Equal(t, ctx.BlockHeight(), int64(100))
	poolRewardInfo, _ := app.MasterchefKeeper.GetPoolRewardInfo(ctx, pools[0].PoolId, "uatom")
	require.Equal(t, poolRewardInfo.LastUpdatedBlock, uint64(99))

	res, err := mk.UserPendingReward(ctx, &types.QueryUserPendingRewardRequest{
		User: addr[0].String(),
	})
	require.NoError(t, err)
	require.Equal(t, res.TotalRewards[0].Amount, math.NewIntWithDecimal(100, 6).Mul(math.NewInt(99)).Quo(math.NewInt(2)))
	res, err = mk.UserPendingReward(ctx, &types.QueryUserPendingRewardRequest{
		User: addr[1].String(),
	})
	require.NoError(t, err)
	require.Equal(t, res.TotalRewards[0].Amount, math.NewIntWithDecimal(100, 6).Mul(math.NewInt(99)).Quo(math.NewInt(2)))

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

	atomAmount := app.BankKeeper.GetBalance(ctx, addr[1], "uatom")
	require.Equal(t, atomAmount.Amount, math.NewIntWithDecimal(100, 6).Mul(math.NewInt(99)).Quo(math.NewInt(2)))

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

	// first user exit pool
	_, err = amm.ExitPool(ctx, addr[1], pools[0].PoolId, share.Quo(math.NewInt(2)), sdk.NewCoins(), "")
	require.NoError(t, err)

	// check rewards after 100 block
	for i := 1; i <= 100; i++ {
		mk.EndBlocker(ctx)
		ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	}

	require.Equal(t, ctx.BlockHeight(), int64(200))
	poolRewardInfo, _ = app.MasterchefKeeper.GetPoolRewardInfo(ctx, pools[0].PoolId, "uatom")
	require.Equal(t, poolRewardInfo.LastUpdatedBlock, uint64(199))

	res, err = mk.UserPendingReward(ctx, &types.QueryUserPendingRewardRequest{
		User: addr[0].String(),
	})
	require.NoError(t, err)
	require.Equal(t, res.TotalRewards[0].String(), "6666670000uatom")
	res, err = mk.UserPendingReward(ctx, &types.QueryUserPendingRewardRequest{
		User: addr[1].String(),
	})
	require.NoError(t, err)
	require.Equal(t, res.TotalRewards[0].String(), "3333335000uatom")

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

	atomAmount = app.BankKeeper.GetBalance(ctx, addr[1], "uatom")
	require.Equal(t, atomAmount.String(), "8283335000uatom")

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
