package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/masterchef/types"
	"github.com/stretchr/testify/require"
)

func TestPoolRewardInfo(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	poolRewardInfos := []types.PoolRewardInfo{
		{
			PoolId:                1,
			RewardDenom:           "reward1",
			PoolAccRewardPerShare: sdk.OneDec(),
			LastUpdatedBlock:      uint64(ctx.BlockHeight()),
		},
		{
			PoolId:                1,
			RewardDenom:           "reward2",
			PoolAccRewardPerShare: sdk.OneDec(),
			LastUpdatedBlock:      uint64(ctx.BlockHeight()),
		},
		{
			PoolId:                2,
			RewardDenom:           "reward2",
			PoolAccRewardPerShare: sdk.OneDec(),
			LastUpdatedBlock:      uint64(ctx.BlockHeight()),
		},
	}
	for _, rewardInfo := range poolRewardInfos {
		err := app.MasterchefKeeper.SetPoolRewardInfo(ctx, rewardInfo)
		require.NoError(t, err)
	}
	for _, rewardInfo := range poolRewardInfos {
		info, found := app.MasterchefKeeper.GetPoolRewardInfo(ctx, rewardInfo.PoolId, rewardInfo.RewardDenom)
		require.True(t, found)
		require.Equal(t, info, rewardInfo)
	}
	rewardInfosStored := app.MasterchefKeeper.GetAllPoolRewardInfos(ctx)
	require.Len(t, rewardInfosStored, 3)

	app.MasterchefKeeper.RemovePoolRewardInfo(ctx, poolRewardInfos[0].PoolId, poolRewardInfos[0].RewardDenom)
	rewardInfosStored = app.MasterchefKeeper.GetAllPoolRewardInfos(ctx)
	require.Len(t, rewardInfosStored, 2)
}
