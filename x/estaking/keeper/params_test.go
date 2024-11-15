package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/elys-network/elys/x/estaking/types"
)

func (suite *EstakingKeeperTestSuite) TestParams() {
	testCases := []struct {
		name                 string
		prerequisiteFunction func()
		postValidateFunction func()
	}{
		{
			"get params",
			func() {
				suite.ResetSuite()

				params := types.DefaultParams()

				suite.app.EstakingKeeper.SetParams(suite.ctx, params)
			},
			func() {
				params := types.DefaultParams()

				suite.Require().EqualValues(params, suite.app.EstakingKeeper.GetParams(suite.ctx))
			},
		},
		{
			"get legacy params",
			func() {
				suite.ResetSuite()
			},
			func() {
				legacyParams := types.LegacyParams{
					StakeIncentives: &types.LegacyIncentiveInfo{
						EdenAmountPerYear:      math.NewInt(0),
						DistributionStartBlock: math.NewInt(0),
						TotalBlocksPerYear:     math.NewInt(0),
						BlocksDistributed:      math.NewInt(0),
					},
					EdenCommitVal:           "",
					EdenbCommitVal:          "",
					MaxEdenRewardAprStakers: math.LegacyNewDec(0),
					EdenBoostApr:            math.LegacyNewDec(0),
					DexRewardsStakers: types.LegacyDexRewardsTracker{
						NumBlocks: math.NewInt(0),
						Amount:    math.LegacyNewDec(0),
					},
				}

				suite.app.EstakingKeeper.SetLegacyParams(suite.ctx, legacyParams)

				suite.Require().EqualValues(legacyParams, suite.app.EstakingKeeper.GetLegacyParams(suite.ctx))
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			tc.postValidateFunction()
		})
	}
}
