package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
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
						EdenAmountPerYear:      sdk.NewInt(0),
						DistributionStartBlock: sdk.NewInt(0),
						TotalBlocksPerYear:     sdk.NewInt(0),
						BlocksDistributed:      sdk.NewInt(0),
					},
					EdenCommitVal:           "",
					EdenbCommitVal:          "",
					MaxEdenRewardAprStakers: sdk.NewDec(0),
					EdenBoostApr:            sdk.NewDec(0),
					DexRewardsStakers: types.LegacyDexRewardsTracker{
						NumBlocks: sdk.NewInt(0),
						Amount:    sdk.NewDec(0),
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
