package keeper_test

import (
	"cosmossdk.io/math"

	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
	"github.com/elys-network/elys/v6/x/masterchef/types"
)

func (suite *MasterchefKeeperTestSuite) TestPool() {

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
		suite.app.MasterchefKeeper.SetPoolInfo(suite.ctx, pool)
		info, found := suite.app.MasterchefKeeper.GetPoolInfo(suite.ctx, pool.PoolId)
		suite.Require().True(found)
		suite.Require().Equal(info, pool)
	}

	poolStored := suite.app.MasterchefKeeper.GetAllPoolInfos(suite.ctx)
	suite.Require().Len(poolStored, 4) // setting it 4 because PoolId = math.MaxInt16 gets initiated in EndBlock

	suite.app.MasterchefKeeper.RemovePoolInfo(suite.ctx, pools[0].PoolId)
	poolStored = suite.app.MasterchefKeeper.GetAllPoolInfos(suite.ctx)
	suite.Require().Len(poolStored, 3) // setting it 3 because PoolId = math.MaxInt16 gets initiated in EndBlock
}

func (suite *MasterchefKeeperTestSuite) TestUpdatePoolMultipliers() {

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
		suite.app.MasterchefKeeper.SetPoolInfo(suite.ctx, pool)
	}
	for _, pool := range pools {
		info, found := suite.app.MasterchefKeeper.GetPoolInfo(suite.ctx, pool.PoolId)
		suite.Require().True(found)
		suite.Require().Equal(info.Multiplier, math.LegacyOneDec())
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
	success := suite.app.MasterchefKeeper.UpdatePoolMultipliers(suite.ctx, poolMultipliers)
	suite.Require().True(success)
	for _, pool := range pools {
		info, found := suite.app.MasterchefKeeper.GetPoolInfo(suite.ctx, pool.PoolId)
		suite.Require().True(found)
		suite.Require().Equal(info.Multiplier, math.LegacyOneDec().Add(math.LegacyOneDec()))
	}
}
