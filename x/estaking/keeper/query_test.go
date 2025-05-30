package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	assetprofiletypes "github.com/elys-network/elys/v6/x/assetprofile/types"
	exdistr "github.com/elys-network/elys/v6/x/estaking/modules/distribution"
	"github.com/elys-network/elys/v6/x/estaking/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
)

func (suite *EstakingKeeperTestSuite) TestQuery() {
	testCases := []struct {
		name                 string
		prerequisiteFunction func()
		postValidateFunction func()
	}{
		{
			"query rewards",
			func() {
				suite.ResetSuite()
			},
			func() {
				// create validator with 50% commission
				validators, err := suite.app.StakingKeeper.GetAllValidators(suite.ctx)
				suite.Require().Nil(err)
				suite.Require().True(len(validators) > 0)
				valAddr, err := sdk.ValAddressFromBech32(validators[0].GetOperator())
				suite.Require().Nil(err)

				delegations, err := suite.app.StakingKeeper.GetValidatorDelegations(suite.ctx, valAddr)
				suite.Require().Nil(err)
				suite.Require().True(len(delegations) > 0)
				addr := sdk.MustAccAddressFromBech32(delegations[0].DelegatorAddress)

				suite.app.AssetprofileKeeper.SetEntry(suite.ctx, assetprofiletypes.Entry{
					BaseDenom:       ptypes.Eden,
					Denom:           ptypes.Eden,
					Decimals:        6,
					CommitEnabled:   true,
					WithdrawEnabled: true,
				})
				suite.app.AssetprofileKeeper.SetEntry(suite.ctx, assetprofiletypes.Entry{
					BaseDenom:       ptypes.BaseCurrency,
					Denom:           ptypes.BaseCurrency,
					Decimals:        6,
					CommitEnabled:   true,
					WithdrawEnabled: true,
				})

				// next block
				ctx := suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)

				params := suite.app.EstakingKeeper.GetParams(ctx)
				params.StakeIncentives = &types.IncentiveInfo{
					EdenAmountPerYear: math.NewInt(1000_000_000_000_000),
					BlocksDistributed: 1,
				}
				params.MaxEdenRewardAprStakers = math.LegacyNewDec(1000_000)
				suite.app.EstakingKeeper.SetParams(ctx, params)

				// update staker rewards
				err = suite.app.EstakingKeeper.UpdateStakersRewards(ctx)
				suite.Require().Nil(err)

				distrAppModule := exdistr.NewAppModule(
					suite.app.AppCodec(), suite.app.DistrKeeper, suite.app.AccountKeeper,
					suite.app.CommitmentKeeper, suite.app.EstakingKeeper,
					&suite.app.AssetprofileKeeper,
					authtypes.FeeCollectorName, suite.app.GetSubspace(distrtypes.ModuleName))
				distrAppModule.AllocateEdenUsdcTokens(ctx)
				distrAppModule.AllocateEdenBTokens(ctx)

				// query rewards
				res, err := suite.app.EstakingKeeper.Rewards(ctx, &types.QueryRewardsRequest{
					Address: addr.String(),
				})
				suite.Require().Nil(err)
				suite.Require().Equal(res.Total.String(), "147608ueden")
			},
		},
		{
			"query invariant",
			func() {
				suite.ResetSuite()
			},
			func() {
				// create validator with 50% commission
				validators, err := suite.app.StakingKeeper.GetAllValidators(suite.ctx)
				suite.Require().Nil(err)
				suite.Require().True(len(validators) > 0)
				valAddr, err := sdk.ValAddressFromBech32(validators[0].GetOperator())
				suite.Require().Nil(err)

				delegations, err := suite.app.StakingKeeper.GetValidatorDelegations(suite.ctx, valAddr)
				suite.Require().Nil(err)
				suite.Require().True(len(delegations) > 0)

				suite.app.AssetprofileKeeper.SetEntry(suite.ctx, assetprofiletypes.Entry{
					BaseDenom:       ptypes.Eden,
					Denom:           ptypes.Eden,
					Decimals:        6,
					CommitEnabled:   true,
					WithdrawEnabled: true,
				})
				suite.app.AssetprofileKeeper.SetEntry(suite.ctx, assetprofiletypes.Entry{
					BaseDenom:       ptypes.BaseCurrency,
					Denom:           ptypes.BaseCurrency,
					Decimals:        6,
					CommitEnabled:   true,
					WithdrawEnabled: true,
				})

				// next block
				ctx := suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)

				params := suite.app.EstakingKeeper.GetParams(ctx)
				params.StakeIncentives = &types.IncentiveInfo{
					EdenAmountPerYear: math.NewInt(1000_000_000_000_000),
					BlocksDistributed: 1,
				}
				params.MaxEdenRewardAprStakers = math.LegacyNewDec(1000_000)
				suite.app.EstakingKeeper.SetParams(ctx, params)

				// update staker rewards
				err = suite.app.EstakingKeeper.UpdateStakersRewards(ctx)
				suite.Require().Nil(err)

				distrAppModule := exdistr.NewAppModule(
					suite.app.AppCodec(), suite.app.DistrKeeper, suite.app.AccountKeeper,
					suite.app.CommitmentKeeper, suite.app.EstakingKeeper,
					&suite.app.AssetprofileKeeper,
					authtypes.FeeCollectorName, suite.app.GetSubspace(distrtypes.ModuleName))
				distrAppModule.AllocateEdenUsdcTokens(ctx)
				distrAppModule.AllocateEdenBTokens(ctx)

				// query invariant
				res, err := suite.app.EstakingKeeper.Invariant(ctx, &types.QueryInvariantRequest{})
				suite.Require().Nil(err)
				suite.Require().Equal(res.TotalBonded.String(), "1000000")
				suite.Require().Equal(res.BondedValidatorTokensSum.String(), "1000000")
			},
		},
		{
			"query invariant with invalid request",
			func() {
				suite.ResetSuite()
			},
			func() {
				// query invariant
				_, err := suite.app.EstakingKeeper.Invariant(suite.ctx, nil)
				suite.Require().Error(err)
			},
		},
		{
			"query rewards with invalid request",
			func() {
				suite.ResetSuite()
			},
			func() {
				// query invariant
				_, err := suite.app.EstakingKeeper.Rewards(suite.ctx, nil)
				suite.Require().Error(err)
			},
		},
		{
			"query rewards with empty delegator address",
			func() {
				suite.ResetSuite()
			},
			func() {
				// query invariant
				_, err := suite.app.EstakingKeeper.Rewards(suite.ctx, &types.QueryRewardsRequest{Address: ""})
				suite.Require().Error(err)
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
