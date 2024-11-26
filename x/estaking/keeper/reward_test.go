package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	"github.com/elys-network/elys/x/estaking/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (suite *EstakingKeeperTestSuite) TestRewardDistribution() {
	testCases := []struct {
		name                 string
		prerequisiteFunction func() (addr sdk.AccAddress)
		postValidateFunction func(addr sdk.AccAddress)
	}{
		{
			"Test reward distribution",
			func() (addr sdk.AccAddress) {
				suite.ResetSuite()

				addr = suite.AddAccounts(1, nil)[0]

				return addr
			},
			func(addr sdk.AccAddress) {
				suite.SetAssetProfile()
				// Delegate elys, check commitment and elys
				validators, err := suite.app.StakingKeeper.GetAllValidators(suite.ctx)
				suite.Require().Nil(err)
				suite.Require().True(len(validators) > 0)
				valAddr := validators[0].GetOperator()

				suite.ctx = suite.ctx.WithBlockHeight(10)

				msgServerCommitment := commitmentkeeper.NewMsgServerImpl(*suite.app.CommitmentKeeper)

				_, err = msgServerCommitment.Stake(suite.ctx, &commitmenttypes.MsgStake{
					Creator:          suite.genAccount.String(),
					Amount:           math.NewInt(100000000000),
					Asset:            "uelys",
					ValidatorAddress: valAddr,
				})
				suite.Require().NoError(err)

				// Check with non-delegator
				delegatedAmount := suite.app.EstakingKeeper.CalcDelegationAmount(suite.ctx, addr)
				suite.Require().Equal(delegatedAmount, math.ZeroInt())

				// Check with genesis account (delegator)
				delegatedAmount = suite.app.EstakingKeeper.CalcDelegationAmount(suite.ctx, suite.genAccount)
				suite.Require().Equal(delegatedAmount, sdk.DefaultPowerReduction.Add(math.NewInt(100000000000)))

				params := suite.app.EstakingKeeper.GetParams(suite.ctx)
				params.StakeIncentives = &types.IncentiveInfo{
					EdenAmountPerYear: math.NewInt(10000000000000),
					BlocksDistributed: 2,
				}
				suite.app.EstakingKeeper.SetParams(suite.ctx, params)

				totalBonded, err := suite.app.EstakingKeeper.TotalBondedTokens(suite.ctx)
				suite.Require().NoError(err)
				suite.Require().Equal(totalBonded.String(), "100001000000")
				suite.Require().Equal(params.MaxEdenRewardAprStakers.String(), "0.300000000000000000")

				commitments := suite.app.CommitmentKeeper.GetAllCommitments(suite.ctx)
				suite.Require().Equal(len(commitments), 1)

				feeCollector := suite.app.AccountKeeper.GetModuleAccount(suite.ctx, authtypes.FeeCollectorName)
				suite.Require().Equal(feeCollector.GetAddress().String(), "cosmos17xpfvakm2amg962yls6f84z3kell8c5lserqta")

				// check rewards after 1000 block
				for i := 1; i <= 1000; i++ {
					_, err := suite.app.BeginBlocker(suite.ctx)
					suite.Require().NoError(err)

					suite.app.EstakingKeeper.EndBlocker(suite.ctx)
					suite.Require().NoError(err)

					suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
				}

				balances := suite.app.CommitmentKeeper.GetAllBalances(suite.ctx, feeCollector.GetAddress())
				suite.Require().Equal(balances.AmountOf(ptypes.Eden).String(), "4756")
				suite.Require().Equal(balances.AmountOf(ptypes.EdenB).String(), "15855")

				totalBonded, err = suite.app.EstakingKeeper.TotalBondedTokens(suite.ctx)
				suite.Require().NoError(err)
				suite.Require().Equal(totalBonded.String(), "100001000000")

				_, err = suite.app.EstakingKeeper.WithdrawAllRewards(suite.ctx, &types.MsgWithdrawAllRewards{
					DelegatorAddress: suite.genAccount.String(),
				})
				suite.Require().NoError(err)

				// Commit 4M eden tokens
				_, err = msgServerCommitment.CommitClaimedRewards(suite.ctx, &commitmenttypes.MsgCommitClaimedRewards{
					Creator: suite.genAccount.String(),
					Denom:   ptypes.Eden,
					Amount:  math.NewInt(4000000),
				})
				suite.Require().NoError(err)

				// Commit 14M edenB tokens
				_, err = msgServerCommitment.CommitClaimedRewards(suite.ctx, &commitmenttypes.MsgCommitClaimedRewards{
					Creator: suite.genAccount.String(),
					Denom:   ptypes.EdenB,
					Amount:  math.NewInt(14000000),
				})
				suite.Require().NoError(err)

				// Now total bonded should include ELYS + EDEN + EDENB
				totalBonded, err = suite.app.EstakingKeeper.TotalBondedTokens(suite.ctx)
				suite.Require().NoError(err)
				suite.Require().Equal(totalBonded.String(), "100019000000")

				// Exclude EDENB
				totalBondedE, err := suite.app.EstakingKeeper.TotalBondedElysEdenTokens(suite.ctx)
				suite.Require().NoError(err)
				suite.Require().Equal(totalBondedE.String(), "100005000000")

				delegatedAmount = suite.app.EstakingKeeper.CalcDelegationAmount(suite.ctx, suite.genAccount)
				suite.Require().Equal(delegatedAmount, sdk.DefaultPowerReduction.Add(math.NewInt(100000000000)))

				// Left over rewards in user address
				commitment := suite.app.CommitmentKeeper.GetCommitments(suite.ctx, suite.genAccount)
				prevEdenBalance := commitment.Claimed.AmountOf(ptypes.Eden)
				prevEdenBBalance := commitment.Claimed.AmountOf(ptypes.EdenB)

				suite.Require().Equal(commitment.Claimed.AmountOf(ptypes.Eden).String(), "423408")
				suite.Require().Equal(commitment.Claimed.AmountOf(ptypes.EdenB).String(), "746243")

				new_users := suite.AddAccounts(1, nil)
				// Stake same amount of tokens in elys from another user
				_, err = msgServerCommitment.Stake(suite.ctx, &commitmenttypes.MsgStake{
					Creator:          new_users[0].String(),
					Amount:           totalBonded,
					Asset:            "uelys",
					ValidatorAddress: valAddr,
				})
				suite.Require().NoError(err)

				for i := 1; i <= 100; i++ {
					_, err := suite.app.BeginBlocker(suite.ctx)
					suite.Require().NoError(err)

					suite.app.EstakingKeeper.EndBlocker(suite.ctx)
					suite.Require().NoError(err)

					suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
				}

				_, err = suite.app.EstakingKeeper.WithdrawAllRewards(suite.ctx, &types.MsgWithdrawAllRewards{
					DelegatorAddress: suite.genAccount.String(),
				})
				suite.Require().NoError(err)

				_, err = suite.app.EstakingKeeper.WithdrawAllRewards(suite.ctx, &types.MsgWithdrawAllRewards{
					DelegatorAddress: new_users[0].String(),
				})
				suite.Require().NoError(err)

				commitment = suite.app.CommitmentKeeper.GetCommitments(suite.ctx, suite.genAccount)

				// Should be equal as number of shares for both users are same
				// Note: Share here are of Eden + Elys + EdenB
				commitmentNew := suite.app.CommitmentKeeper.GetCommitments(suite.ctx, new_users[0])
				suite.Require().Equal(commitment.Claimed.AmountOf(ptypes.Eden).Sub(prevEdenBalance),
					commitmentNew.Claimed.AmountOf(ptypes.Eden).Add(math.NewInt(3)))
				// This number will be different as new_user doesn't have EdenB staked and has higher shares if we exclude edenb
				suite.Require().Less(commitment.Claimed.AmountOf(ptypes.EdenB).Sub(prevEdenBBalance).String(),
					commitmentNew.Claimed.AmountOf(ptypes.EdenB).String())

			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			addr := tc.prerequisiteFunction()
			tc.postValidateFunction(addr)
		})
	}
}
