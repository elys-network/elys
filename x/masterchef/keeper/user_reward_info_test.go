package keeper_test

import (
	"testing"

	"github.com/cometbft/cometbft/crypto/ed25519"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/masterchef/types"
	"github.com/stretchr/testify/require"
)

func TestUserRewardInfo(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	user1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	user2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	userRewardInfos := []types.UserRewardInfo{
		{
			User:          user1.String(),
			PoolId:        1,
			RewardDenom:   "reward1",
			RewardDebt:    sdk.ZeroDec(),
			RewardPending: sdk.ZeroDec(),
		},
		{
			User:          user1.String(),
			PoolId:        1,
			RewardDenom:   "reward2",
			RewardDebt:    sdk.ZeroDec(),
			RewardPending: sdk.ZeroDec(),
		},
		{
			User:          user2.String(),
			PoolId:        2,
			RewardDenom:   "reward2",
			RewardDebt:    sdk.ZeroDec(),
			RewardPending: sdk.ZeroDec(),
		},
	}
	for _, rewardInfo := range userRewardInfos {
		app.MasterchefKeeper.SetUserRewardInfo(ctx, rewardInfo)

	}
	for _, rewardInfo := range userRewardInfos {
		info, found := app.MasterchefKeeper.GetUserRewardInfo(ctx, rewardInfo.GetUserAccount(), rewardInfo.PoolId, rewardInfo.RewardDenom)
		require.True(t, found)
		require.Equal(t, info, rewardInfo)
	}
	rewardInfosStored := app.MasterchefKeeper.GetAllUserRewardInfos(ctx)
	require.Len(t, rewardInfosStored, 3)

	app.MasterchefKeeper.RemoveUserRewardInfo(ctx, userRewardInfos[0].User, userRewardInfos[0].PoolId, userRewardInfos[0].RewardDenom)
	rewardInfosStored = app.MasterchefKeeper.GetAllUserRewardInfos(ctx)
	require.Len(t, rewardInfosStored, 2)
}
