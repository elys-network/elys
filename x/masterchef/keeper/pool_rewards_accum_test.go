package keeper_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/masterchef/types"
	"github.com/stretchr/testify/require"
)

func TestPoolRewardsAccum(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true)

	now := time.Now()
	accums := []types.PoolRewardsAccum{
		{
			PoolId:      1,
			Timestamp:   uint64(now.Unix() - 86400),
			BlockHeight: 1,
			DexReward:   math.LegacyNewDec(1000),
			GasReward:   math.LegacyNewDec(1000),
			EdenReward:  math.LegacyNewDec(1000),
		},
		{
			PoolId:      1,
			Timestamp:   uint64(now.Unix()),
			BlockHeight: 1,
			DexReward:   math.LegacyNewDec(2000),
			GasReward:   math.LegacyNewDec(2000),
			EdenReward:  math.LegacyNewDec(2000),
		},
		{
			PoolId:      2,
			Timestamp:   uint64(now.Unix() - 86400),
			BlockHeight: 1,
			DexReward:   math.LegacyNewDec(1000),
			GasReward:   math.LegacyNewDec(1000),
			EdenReward:  math.LegacyNewDec(1000),
		},
		{
			PoolId:      2,
			Timestamp:   uint64(now.Unix()),
			BlockHeight: 1,
			DexReward:   math.LegacyNewDec(3000),
			GasReward:   math.LegacyNewDec(3000),
			EdenReward:  math.LegacyNewDec(3000),
		},
	}
	for _, accum := range accums {
		app.MasterchefKeeper.SetPoolRewardsAccum(ctx, accum)
	}

	for _, accum := range accums {
		storedAccum, err := app.MasterchefKeeper.GetPoolRewardsAccum(ctx, accum.PoolId, accum.Timestamp)
		require.NoError(t, err)
		require.Equal(t, storedAccum, accum)
	}

	accum := app.MasterchefKeeper.FirstPoolRewardsAccum(ctx, 1)
	require.Equal(t, accum, accums[0])
	accum = app.MasterchefKeeper.LastPoolRewardsAccum(ctx, 1)
	require.Equal(t, accum, accums[1])

	app.MasterchefKeeper.DeletePoolRewardsAccum(ctx, accums[0])
	accum = app.MasterchefKeeper.FirstPoolRewardsAccum(ctx, 1)
	require.Equal(t, accum, accums[1])
	accum = app.MasterchefKeeper.LastPoolRewardsAccum(ctx, 1)
	require.Equal(t, accum, accums[1])
}

func TestAddPoolRewardsAccum(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true)
	k := app.MasterchefKeeper

	tests := []struct {
		name       string
		poolId     uint64
		timestamp  uint64
		height     int64
		dexReward  math.LegacyDec
		gasReward  math.LegacyDec
		edenReward math.LegacyDec
	}{
		{
			name:       "Add rewards to new pool",
			poolId:     1,
			timestamp:  uint64(time.Now().Unix()),
			height:     100,
			dexReward:  math.LegacyNewDec(10),
			gasReward:  math.LegacyNewDec(5),
			edenReward: math.LegacyNewDec(3),
		},
		{
			name:       "Add rewards to existing pool",
			poolId:     1,
			timestamp:  uint64(time.Now().Unix()) + 3600, // 1 hour later
			height:     200,
			dexReward:  math.LegacyNewDec(20),
			gasReward:  math.LegacyNewDec(10),
			edenReward: math.LegacyNewDec(6),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k.AddPoolRewardsAccum(ctx, tt.poolId, tt.timestamp, tt.height, tt.dexReward, tt.gasReward, tt.edenReward)

			accum, err := k.GetPoolRewardsAccum(ctx, tt.poolId, tt.timestamp)
			require.NoError(t, err)

			require.Equal(t, tt.poolId, accum.PoolId)
			require.Equal(t, tt.timestamp, accum.Timestamp)
			require.Equal(t, tt.height, accum.BlockHeight)

			if tt.name == "Add rewards to new pool" {
				require.Equal(t, tt.dexReward, accum.DexReward)
				require.Equal(t, tt.gasReward, accum.GasReward)
				require.Equal(t, tt.edenReward, accum.EdenReward)

				// Check forward
				forwardEden := k.ForwardEdenCalc(ctx, tt.poolId)
				require.Equal(t, math.LegacyZeroDec(), forwardEden)
			} else {
				// For existing pool, rewards should be cumulative
				require.Equal(t, math.LegacyNewDec(30), accum.DexReward)
				require.Equal(t, math.LegacyNewDec(15), accum.GasReward)
				require.Equal(t, math.LegacyNewDec(9), accum.EdenReward)

				// Check forward
				forwardEden := k.ForwardEdenCalc(ctx, tt.poolId)
				require.Equal(t, math.LegacyMustNewDecFromStr("21600").Mul(tt.edenReward), forwardEden)
			}
		})
	}
}
