package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/masterchef/types"
	"github.com/stretchr/testify/require"
)

func TestPool(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	pools := []types.PoolInfo{
		{
			PoolId:               1,
			RewardWallet:         ammtypes.NewPoolRevenueAddress(1).String(),
			Multiplier:           sdk.OneDec(),
			EdenApr:              sdk.OneDec(),
			DexApr:               sdk.OneDec(),
			GasApr:               sdk.OneDec(),
			ExternalIncentiveApr: sdk.OneDec(),
			ExternalRewardDenoms: []string{
				"rewardDenom1",
				"rewardDenom2",
			},
		},
		{
			PoolId:               2,
			RewardWallet:         ammtypes.NewPoolRevenueAddress(2).String(),
			Multiplier:           sdk.OneDec(),
			EdenApr:              sdk.OneDec(),
			DexApr:               sdk.OneDec(),
			GasApr:               sdk.OneDec(),
			ExternalIncentiveApr: sdk.OneDec(),
			ExternalRewardDenoms: []string{
				"rewardDenom1",
				"rewardDenom2",
			},
		},
		{
			PoolId:               3,
			RewardWallet:         ammtypes.NewPoolRevenueAddress(3).String(),
			Multiplier:           sdk.OneDec(),
			EdenApr:              sdk.OneDec(),
			DexApr:               sdk.OneDec(),
			GasApr:               sdk.OneDec(),
			ExternalIncentiveApr: sdk.OneDec(),
			ExternalRewardDenoms: []string{
				"rewardDenom1",
				"rewardDenom2",
			},
		},
	}
	for _, pool := range pools {
		app.MasterchefKeeper.SetPool(ctx, pool)
	}
	for _, pool := range pools {
		info, found := app.MasterchefKeeper.GetPool(ctx, pool.PoolId)
		require.True(t, found)
		require.Equal(t, info, pool)
	}
	poolStored := app.MasterchefKeeper.GetAllPools(ctx)
	require.Len(t, poolStored, 3)

	app.MasterchefKeeper.RemovePool(ctx, pools[0].PoolId)
	poolStored = app.MasterchefKeeper.GetAllPools(ctx)
	require.Len(t, poolStored, 2)
}

func TestUpdatePoolMultipliers(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	pools := []types.PoolInfo{
		{
			PoolId:               1,
			RewardWallet:         ammtypes.NewPoolRevenueAddress(1).String(),
			Multiplier:           sdk.OneDec(),
			EdenApr:              sdk.OneDec(),
			DexApr:               sdk.OneDec(),
			GasApr:               sdk.OneDec(),
			ExternalIncentiveApr: sdk.OneDec(),
			ExternalRewardDenoms: []string{
				"rewardDenom1",
				"rewardDenom2",
			},
		},
		{
			PoolId:               2,
			RewardWallet:         ammtypes.NewPoolRevenueAddress(2).String(),
			Multiplier:           sdk.OneDec(),
			EdenApr:              sdk.OneDec(),
			DexApr:               sdk.OneDec(),
			GasApr:               sdk.OneDec(),
			ExternalIncentiveApr: sdk.OneDec(),
			ExternalRewardDenoms: []string{
				"rewardDenom1",
				"rewardDenom2",
			},
		},
		{
			PoolId:               3,
			RewardWallet:         ammtypes.NewPoolRevenueAddress(3).String(),
			Multiplier:           sdk.OneDec(),
			EdenApr:              sdk.OneDec(),
			DexApr:               sdk.OneDec(),
			ExternalIncentiveApr: sdk.OneDec(),
			GasApr:               sdk.OneDec(),
			ExternalRewardDenoms: []string{
				"rewardDenom1",
				"rewardDenom2",
			},
		},
	}
	for _, pool := range pools {
		app.MasterchefKeeper.SetPool(ctx, pool)
	}
	for _, pool := range pools {
		info, found := app.MasterchefKeeper.GetPool(ctx, pool.PoolId)
		require.True(t, found)
		require.Equal(t, info.Multiplier, sdk.OneDec())
	}

	poolMultipliers := []types.PoolMultiplier{
		{
			PoolId:     1,
			Multiplier: sdk.OneDec().Add(sdk.OneDec()),
		}, {
			PoolId:     2,
			Multiplier: sdk.OneDec().Add(sdk.OneDec()),
		}, {
			PoolId:     3,
			Multiplier: sdk.OneDec().Add(sdk.OneDec()),
		},
	}
	success := app.MasterchefKeeper.UpdatePoolMultipliers(ctx, poolMultipliers)
	require.True(t, success)
	for _, pool := range pools {
		info, found := app.MasterchefKeeper.GetPool(ctx, pool.PoolId)
		require.True(t, found)
		require.Equal(t, info.Multiplier, sdk.OneDec().Add(sdk.OneDec()))
	}
}
