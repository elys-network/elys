package keeper_test

import (
	sdkmath "cosmossdk.io/math"

	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/masterchef/types"
)

func (suite *MasterchefKeeperTestSuite) TestUserRewardInfo() {

	user1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	user2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	userRewardInfos := []types.UserRewardInfo{
		{
			User:          user1.String(),
			PoolId:        1,
			RewardDenom:   "reward1",
			RewardDebt:    sdkmath.LegacyZeroDec(),
			RewardPending: sdkmath.LegacyZeroDec(),
		},
		{
			User:          user1.String(),
			PoolId:        1,
			RewardDenom:   "reward2",
			RewardDebt:    sdkmath.LegacyZeroDec(),
			RewardPending: sdkmath.LegacyZeroDec(),
		},
		{
			User:          user2.String(),
			PoolId:        2,
			RewardDenom:   "reward2",
			RewardDebt:    sdkmath.LegacyZeroDec(),
			RewardPending: sdkmath.LegacyZeroDec(),
		},
	}
	for _, rewardInfo := range userRewardInfos {
		suite.app.MasterchefKeeper.SetUserRewardInfo(suite.ctx, rewardInfo)

	}
	for _, rewardInfo := range userRewardInfos {
		info, found := suite.app.MasterchefKeeper.GetUserRewardInfo(suite.ctx, rewardInfo.GetUserAccount(), rewardInfo.PoolId, rewardInfo.RewardDenom)
		suite.Require().True(found)
		suite.Require().Equal(info, rewardInfo)
	}
	rewardInfosStored := suite.app.MasterchefKeeper.GetAllUserRewardInfos(suite.ctx)
	suite.Require().Len(rewardInfosStored, 3)

	suite.app.MasterchefKeeper.RemoveUserRewardInfo(suite.ctx, userRewardInfos[0].GetUserAccount(), userRewardInfos[0].PoolId, userRewardInfos[0].RewardDenom)
	rewardInfosStored = suite.app.MasterchefKeeper.GetAllUserRewardInfos(suite.ctx)
	suite.Require().Len(rewardInfosStored, 2)
}
