package keeper_test

import (
	"testing"

	"cosmossdk.io/math"

	simapp "github.com/elys-network/elys/app"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/masterchef/types"
	"github.com/stretchr/testify/require"
)

func TestPool(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

	pools := []types.PoolInfo{
		{
			PoolId:               1,
			RewardWallet:         ammtypes.NewPoolRevenueAddress(1).String(),
			Multiplier:           math.LegacyOneDec(),
			EdenApr:              math.LegacyOneDec(),
			DexApr:               math.LegacyOneDec(),
			GasApr:               math.LegacyOneDec(),
			ExternalIncentiveApr: math.LegacyOneDec(),
			ExternalRewardDenoms: []string{
				"rewardDenom1",
				"rewardDenom2",
			},
		},
		{
			PoolId:               2,
			RewardWallet:         ammtypes.NewPoolRevenueAddress(2).String(),
			Multiplier:           math.LegacyOneDec(),
			EdenApr:              math.LegacyOneDec(),
			DexApr:               math.LegacyOneDec(),
			GasApr:               math.LegacyOneDec(),
			ExternalIncentiveApr: math.LegacyOneDec(),
			ExternalRewardDenoms: []string{
				"rewardDenom1",
				"rewardDenom2",
			},
		},
		{
			PoolId:               3,
			RewardWallet:         ammtypes.NewPoolRevenueAddress(3).String(),
			Multiplier:           math.LegacyOneDec(),
			EdenApr:              math.LegacyOneDec(),
			DexApr:               math.LegacyOneDec(),
			GasApr:               math.LegacyOneDec(),
			ExternalIncentiveApr: math.LegacyOneDec(),
			ExternalRewardDenoms: []string{
				"rewardDenom1",
				"rewardDenom2",
			},
		},
	}
	for _, pool := range pools {
		app.MasterchefKeeper.SetPoolInfo(ctx, pool)
	}
	for _, pool := range pools {
		info, found := app.MasterchefKeeper.GetPoolInfo(ctx, pool.PoolId)
		require.True(t, found)
		require.Equal(t, info, pool)
	}
	poolStored := app.MasterchefKeeper.GetAllPoolInfos(ctx)
	require.Len(t, poolStored, 3)

	app.MasterchefKeeper.RemovePoolInfo(ctx, pools[0].PoolId)
	poolStored = app.MasterchefKeeper.GetAllPoolInfos(ctx)
	require.Len(t, poolStored, 2)
}

func TestUpdatePoolMultipliers(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

	pools := []types.PoolInfo{
		{
			PoolId:               1,
			RewardWallet:         ammtypes.NewPoolRevenueAddress(1).String(),
			Multiplier:           math.LegacyOneDec(),
			EdenApr:              math.LegacyOneDec(),
			DexApr:               math.LegacyOneDec(),
			GasApr:               math.LegacyOneDec(),
			ExternalIncentiveApr: math.LegacyOneDec(),
			ExternalRewardDenoms: []string{
				"rewardDenom1",
				"rewardDenom2",
			},
		},
		{
			PoolId:               2,
			RewardWallet:         ammtypes.NewPoolRevenueAddress(2).String(),
			Multiplier:           math.LegacyOneDec(),
			EdenApr:              math.LegacyOneDec(),
			DexApr:               math.LegacyOneDec(),
			GasApr:               math.LegacyOneDec(),
			ExternalIncentiveApr: math.LegacyOneDec(),
			ExternalRewardDenoms: []string{
				"rewardDenom1",
				"rewardDenom2",
			},
		},
		{
			PoolId:               3,
			RewardWallet:         ammtypes.NewPoolRevenueAddress(3).String(),
			Multiplier:           math.LegacyOneDec(),
			EdenApr:              math.LegacyOneDec(),
			DexApr:               math.LegacyOneDec(),
			ExternalIncentiveApr: math.LegacyOneDec(),
			GasApr:               math.LegacyOneDec(),
			ExternalRewardDenoms: []string{
				"rewardDenom1",
				"rewardDenom2",
			},
		},
	}
	for _, pool := range pools {
		app.MasterchefKeeper.SetPoolInfo(ctx, pool)
	}
	for _, pool := range pools {
		info, found := app.MasterchefKeeper.GetPoolInfo(ctx, pool.PoolId)
		require.True(t, found)
		require.Equal(t, info.Multiplier, math.LegacyOneDec())
	}

	poolMultipliers := []types.PoolMultiplier{
		{
			PoolId:     1,
			Multiplier: math.LegacyOneDec().Add(math.LegacyOneDec()),
		}, {
			PoolId:     2,
			Multiplier: math.LegacyOneDec().Add(math.LegacyOneDec()),
		}, {
			PoolId:     3,
			Multiplier: math.LegacyOneDec().Add(math.LegacyOneDec()),
		},
	}
	success := app.MasterchefKeeper.UpdatePoolMultipliers(ctx, poolMultipliers)
	require.True(t, success)
	for _, pool := range pools {
		info, found := app.MasterchefKeeper.GetPoolInfo(ctx, pool.PoolId)
		require.True(t, found)
		require.Equal(t, info.Multiplier, math.LegacyOneDec().Add(math.LegacyOneDec()))
	}
}
