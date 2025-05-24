package keeper_test

import (
	sdkmath "cosmossdk.io/math"

	"github.com/elys-network/elys/v5/x/masterchef/types"
)

func (suite *MasterchefKeeperTestSuite) TestPoolRewardInfo() {

	poolRewardInfos := []types.PoolRewardInfo{
		{
			PoolId:                1,
			RewardDenom:           "reward1",
			PoolAccRewardPerShare: sdkmath.LegacyOneDec(),
			LastUpdatedBlock:      uint64(suite.ctx.BlockHeight()),
		},
		{
			PoolId:                1,
			RewardDenom:           "reward2",
			PoolAccRewardPerShare: sdkmath.LegacyOneDec(),
			LastUpdatedBlock:      uint64(suite.ctx.BlockHeight()),
		},
		{
			PoolId:                2,
			RewardDenom:           "reward2",
			PoolAccRewardPerShare: sdkmath.LegacyOneDec(),
			LastUpdatedBlock:      uint64(suite.ctx.BlockHeight()),
		},
	}
	for _, rewardInfo := range poolRewardInfos {
		suite.app.MasterchefKeeper.SetPoolRewardInfo(suite.ctx, rewardInfo)
	}
	for _, rewardInfo := range poolRewardInfos {
		info, found := suite.app.MasterchefKeeper.GetPoolRewardInfo(suite.ctx, rewardInfo.PoolId, rewardInfo.RewardDenom)
		suite.Require().True(found)
		suite.Require().Equal(info, rewardInfo)
	}
	rewardInfosStored := suite.app.MasterchefKeeper.GetAllPoolRewardInfos(suite.ctx)
	suite.Require().Len(rewardInfosStored, 3)

	suite.app.MasterchefKeeper.RemovePoolRewardInfo(suite.ctx, poolRewardInfos[0].PoolId, poolRewardInfos[0].RewardDenom)
	rewardInfosStored = suite.app.MasterchefKeeper.GetAllPoolRewardInfos(suite.ctx)
	suite.Require().Len(rewardInfosStored, 2)
}
