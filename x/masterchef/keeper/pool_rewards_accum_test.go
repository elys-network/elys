package keeper_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/masterchef/types"
	"github.com/stretchr/testify/require"
)

func TestPoolRewardsAccum(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

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
