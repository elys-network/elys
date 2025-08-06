package keeper_test

import (
	"time"

	"cosmossdk.io/math"
	"github.com/elys-network/elys/v7/x/vaults/types"
)

func (suite *KeeperTestSuite) TestPoolRewardsAccum() {

	now := time.Now()
	accums := []types.PoolRewardsAccum{
		{
			PoolId:      1,
			Timestamp:   uint64(now.Unix() - 86400),
			BlockHeight: 1,
			UsdcReward:  math.LegacyNewDec(1000),
			EdenReward:  math.LegacyNewDec(1000),
		},
		{
			PoolId:      1,
			Timestamp:   uint64(now.Unix()),
			BlockHeight: 1,
			UsdcReward:  math.LegacyNewDec(2000),
			EdenReward:  math.LegacyNewDec(2000),
		},
		{
			PoolId:      2,
			Timestamp:   uint64(now.Unix() - 86400),
			BlockHeight: 1,
			UsdcReward:  math.LegacyNewDec(1000),
			EdenReward:  math.LegacyNewDec(1000),
		},
		{
			PoolId:      2,
			Timestamp:   uint64(now.Unix()),
			BlockHeight: 1,
			UsdcReward:  math.LegacyNewDec(3000),
			EdenReward:  math.LegacyNewDec(3000),
		},
	}
	for _, accum := range accums {
		suite.app.VaultsKeeper.SetPoolRewardsAccum(suite.ctx, accum)
	}

	for _, accum := range accums {
		storedAccum, err := suite.app.VaultsKeeper.GetPoolRewardsAccum(suite.ctx, accum.PoolId, accum.Timestamp)
		suite.Require().NoError(err)
		suite.Require().Equal(storedAccum, accum)
	}

	accum := suite.app.VaultsKeeper.FirstPoolRewardsAccum(suite.ctx, 1)
	suite.Require().Equal(accum, accums[0])
	accum = suite.app.VaultsKeeper.LastPoolRewardsAccum(suite.ctx, 1)
	suite.Require().Equal(accum, accums[1])

	suite.app.VaultsKeeper.DeletePoolRewardsAccum(suite.ctx, accums[0])
	accum = suite.app.VaultsKeeper.FirstPoolRewardsAccum(suite.ctx, 1)
	suite.Require().Equal(accum, accums[1])
	accum = suite.app.VaultsKeeper.LastPoolRewardsAccum(suite.ctx, 1)
	suite.Require().Equal(accum, accums[1])
}

func (suite *KeeperTestSuite) TestAddPoolRewardsAccum() {

	tests := []struct {
		name       string
		poolId     uint64
		timestamp  uint64
		height     int64
		usdcReward math.LegacyDec
		edenReward math.LegacyDec
	}{
		{
			name:       "Add rewards to new pool",
			poolId:     1,
			timestamp:  uint64(time.Now().Unix()),
			height:     100,
			usdcReward: math.LegacyNewDec(5),
			edenReward: math.LegacyNewDec(3),
		},
		{
			name:       "Add rewards to existing pool",
			poolId:     1,
			timestamp:  uint64(time.Now().Unix()) + 3600, // 1 hour later
			height:     200,
			usdcReward: math.LegacyNewDec(10),
			edenReward: math.LegacyNewDec(6),
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.app.VaultsKeeper.AddPoolRewardsAccum(suite.ctx, tt.poolId, tt.timestamp, tt.height, tt.usdcReward, tt.edenReward)

			accum, err := suite.app.VaultsKeeper.GetPoolRewardsAccum(suite.ctx, tt.poolId, tt.timestamp)
			suite.Require().NoError(err)

			suite.Require().Equal(tt.poolId, accum.PoolId)
			suite.Require().Equal(tt.timestamp, accum.Timestamp)
			suite.Require().Equal(tt.height, accum.BlockHeight)

			if tt.name == "Add rewards to new pool" {
				suite.Require().Equal(tt.usdcReward, accum.UsdcReward)
				suite.Require().Equal(tt.edenReward, accum.EdenReward)
			} else {
				// For existing pool, rewards should be cumulative
				suite.Require().Equal(math.LegacyNewDec(15), accum.UsdcReward)
				suite.Require().Equal(math.LegacyNewDec(9), accum.EdenReward)
			}
		})
	}
}
