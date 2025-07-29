package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/vaults/types"
)

func (suite *KeeperTestSuite) TestUserRewardInfo() {
	// Create test accounts
	user := sdk.AccAddress([]byte("user1"))
	user2 := sdk.AccAddress([]byte("user2"))

	// Create test user reward infos
	userRewardInfos := []types.UserRewardInfo{
		{
			User:          user.String(),
			PoolId:        1,
			RewardDenom:   "reward1",
			RewardDebt:    sdkmath.LegacyOneDec(),
			RewardPending: sdkmath.LegacyOneDec(),
		},
		{
			User:          user.String(),
			PoolId:        1,
			RewardDenom:   "reward2",
			RewardDebt:    sdkmath.LegacyOneDec(),
			RewardPending: sdkmath.LegacyOneDec(),
		},
		{
			User:          user2.String(),
			PoolId:        2,
			RewardDenom:   "reward2",
			RewardDebt:    sdkmath.LegacyOneDec(),
			RewardPending: sdkmath.LegacyOneDec(),
		},
	}

	// Test SetUserRewardInfo and GetUserRewardInfo
	for _, rewardInfo := range userRewardInfos {
		suite.app.VaultsKeeper.SetUserRewardInfo(suite.ctx, rewardInfo)
		info, found := suite.app.VaultsKeeper.GetUserRewardInfo(suite.ctx, rewardInfo.GetUserAccount(), rewardInfo.PoolId, rewardInfo.RewardDenom)
		suite.Require().True(found)
		suite.Require().Equal(info, rewardInfo)
	}

	// Test GetAllUserRewardInfos
	rewardInfosStored := suite.app.VaultsKeeper.GetAllUserRewardInfos(suite.ctx)
	suite.Require().Len(rewardInfosStored, 3)

	// Test RemoveUserRewardInfo
	suite.app.VaultsKeeper.RemoveUserRewardInfo(suite.ctx, userRewardInfos[0].GetUserAccount(), userRewardInfos[0].PoolId, userRewardInfos[0].RewardDenom)
	rewardInfosStored = suite.app.VaultsKeeper.GetAllUserRewardInfos(suite.ctx)
	suite.Require().Len(rewardInfosStored, 2)

	// Verify the removed reward info is gone
	_, found := suite.app.VaultsKeeper.GetUserRewardInfo(suite.ctx, userRewardInfos[0].GetUserAccount(), userRewardInfos[0].PoolId, userRewardInfos[0].RewardDenom)
	suite.Require().False(found)

	// Test GetUserRewardInfo for non-existent reward info
	_, found = suite.app.VaultsKeeper.GetUserRewardInfo(suite.ctx, user, 999, "nonexistent")
	suite.Require().False(found)
}
