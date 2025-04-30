package keeper_test

import (
	"time"

	"cosmossdk.io/math"
	"github.com/elys-network/elys/x/masterchef/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (suite *MasterchefKeeperTestSuite) TestPoolRewardsAccum() {

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
		suite.app.MasterchefKeeper.SetPoolRewardsAccum(suite.ctx, accum)
	}

	for _, accum := range accums {
		storedAccum, err := suite.app.MasterchefKeeper.GetPoolRewardsAccum(suite.ctx, accum.PoolId, accum.Timestamp)
		suite.Require().NoError(err)
		suite.Require().Equal(storedAccum, accum)
	}

	accum := suite.app.MasterchefKeeper.FirstPoolRewardsAccum(suite.ctx, 1)
	suite.Require().Equal(accum, accums[0])
	accum = suite.app.MasterchefKeeper.LastPoolRewardsAccum(suite.ctx, 1)
	suite.Require().Equal(accum, accums[1])

	suite.app.MasterchefKeeper.DeletePoolRewardsAccum(suite.ctx, accums[0])
	accum = suite.app.MasterchefKeeper.FirstPoolRewardsAccum(suite.ctx, 1)
	suite.Require().Equal(accum, accums[1])
	accum = suite.app.MasterchefKeeper.LastPoolRewardsAccum(suite.ctx, 1)
	suite.Require().Equal(accum, accums[1])
}

func (suite *MasterchefKeeperTestSuite) TestAddPoolRewardsAccum() {

	tests := []struct {
		name       string
		poolId     uint64
		timestamp  uint64
		height     int64
		dexReward  osmomath.BigDec
		gasReward  osmomath.BigDec
		edenReward osmomath.BigDec
	}{
		{
			name:       "Add rewards to new pool",
			poolId:     1,
			timestamp:  uint64(time.Now().Unix()),
			height:     100,
			gasReward:  osmomath.NewBigDec(5),
			dexReward:  osmomath.NewBigDec(10),
			edenReward: osmomath.NewBigDec(3),
		},
		{
			name:       "Add rewards to existing pool",
			poolId:     1,
			timestamp:  uint64(time.Now().Unix()) + 3600, // 1 hour later
			height:     200,
			gasReward:  osmomath.NewBigDec(10),
			dexReward:  osmomath.NewBigDec(20),
			edenReward: osmomath.NewBigDec(6),
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.app.MasterchefKeeper.AddPoolRewardsAccum(suite.ctx, tt.poolId, tt.timestamp, tt.height, tt.dexReward, tt.gasReward, tt.edenReward)

			accum, err := suite.app.MasterchefKeeper.GetPoolRewardsAccum(suite.ctx, tt.poolId, tt.timestamp)
			suite.Require().NoError(err)

			suite.Require().Equal(tt.poolId, accum.PoolId)
			suite.Require().Equal(tt.timestamp, accum.Timestamp)
			suite.Require().Equal(tt.height, accum.BlockHeight)

			if tt.name == "Add rewards to new pool" {
				suite.Require().Equal(tt.dexReward, osmomath.BigDecFromDec(accum.DexReward))
				suite.Require().Equal(tt.gasReward, osmomath.BigDecFromDec(accum.GasReward))
				suite.Require().Equal(tt.edenReward, osmomath.BigDecFromDec(accum.EdenReward))

				// Check forward
				forwardEden := suite.app.MasterchefKeeper.ForwardEdenCalc(suite.ctx, tt.poolId)
				suite.Require().Equal(osmomath.ZeroBigDec(), forwardEden)
			} else {
				// For existing pool, rewards should be cumulative
				suite.Require().Equal(math.LegacyNewDec(30), accum.DexReward)
				suite.Require().Equal(math.LegacyNewDec(15), accum.GasReward)
				suite.Require().Equal(math.LegacyNewDec(9), accum.EdenReward)

				// Check forward
				forwardEden := suite.app.MasterchefKeeper.ForwardEdenCalc(suite.ctx, tt.poolId)
				suite.Require().Equal(osmomath.MustNewBigDecFromStr("21600").Mul(tt.edenReward), forwardEden)
			}
		})
	}
}
