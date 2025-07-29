package keeper_test

import (
	"cosmossdk.io/math"
	commkeeper "github.com/elys-network/elys/v7/x/commitment/keeper"
	ctypes "github.com/elys-network/elys/v7/x/commitment/types"
	ptypes "github.com/elys-network/elys/v7/x/parameter/types"
)

func (suite *EstakingKeeperTestSuite) TestBurnEdenBFromElysUnstaked() {
	testCases := []struct {
		name                 string
		prerequisiteFunction func()
		postValidateFunction func()
	}{
		{
			"burn EdenB from Elys unstaked",
			func() {
				_ = suite.PrepareUnclaimedTokens()
				committed := suite.PrepareCommittedTokens()
				suite.AddTestCommitment(committed)
				suite.app.EstakingKeeper.TakeDelegationSnapshot(suite.ctx, suite.genAccount)
			},
			func() {},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()

			// burn amount = 100000 (unbonded amt) / (1000000 (elys staked) + 10000 (Eden committed)) * (20000 EdenB + 5000 EdenB committed)
			unbondAmt, err := suite.app.StakingKeeper.Unbond(suite.ctx, suite.genAccount, suite.valAddr, math.LegacyNewDecWithPrec(10, 2))
			suite.Require().Equal(unbondAmt, math.NewInt(100000))
			suite.Require().NoError(err)

			// Process EdenB burn operation
			suite.app.EstakingKeeper.EndBlocker(suite.ctx)

			tc.postValidateFunction()
		})
	}
}

func (suite *EstakingKeeperTestSuite) TestBurnEdenBFromEdenUncommitted() {
	testCases := []struct {
		name                 string
		prerequisiteFunction func()
		postValidateFunction func()
	}{
		{
			"burn EdenB from Eden uncommitted",
			func() {
				suite.PrepareUnclaimedTokens()
				committed := suite.PrepareCommittedTokens()
				suite.AddTestClaimed(committed)
				suite.SetAssetProfile()
			},
			func() {},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			// Track Elys staked amount
			suite.app.EstakingKeeper.EndBlocker(suite.ctx)

			msgServer := commkeeper.NewMsgServerImpl(*suite.app.CommitmentKeeper)

			_, err := msgServer.CommitClaimedRewards(suite.ctx, &ctypes.MsgCommitClaimedRewards{
				Creator: suite.genAccount.String(),
				Amount:  math.NewInt(1000),
				Denom:   ptypes.Eden,
			})
			suite.Require().NoError(err)

			// Uncommit tokens
			_, err = msgServer.UncommitTokens(suite.ctx, &ctypes.MsgUncommitTokens{
				Creator: suite.genAccount.String(),
				Amount:  math.NewInt(1000),
				Denom:   ptypes.Eden,
			})
			suite.Require().NoError(err)

			tc.postValidateFunction()
		})
	}
}
