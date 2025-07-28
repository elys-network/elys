package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/elys-network/elys/v7/x/estaking/keeper"
	exdistr "github.com/elys-network/elys/v7/x/estaking/modules/distribution"
	"github.com/elys-network/elys/v7/x/estaking/types"
	ptypes "github.com/elys-network/elys/v7/x/parameter/types"
	ttypes "github.com/elys-network/elys/v7/x/tokenomics/types"
)

func (suite *EstakingKeeperTestSuite) TestAbci() {
	testCases := []struct {
		name                 string
		prerequisiteFunction func() (addr sdk.AccAddress, valAddr sdk.ValAddress)
		postValidateFunction func(addr sdk.AccAddress, valAddr sdk.ValAddress)
	}{
		{
			"update stakers rewards",
			func() (addr sdk.AccAddress, valAddr sdk.ValAddress) {
				suite.ResetSuite()

				suite.SetAssetProfile()

				// create validator with 50% commission
				validators, err := suite.app.StakingKeeper.GetAllValidators(suite.ctx)
				suite.Require().Nil(err)
				suite.Require().True(len(validators) > 0)
				valAddr, err = sdk.ValAddressFromBech32(validators[0].GetOperator())
				suite.Require().Nil(err)

				delegations, err := suite.app.StakingKeeper.GetValidatorDelegations(suite.ctx, valAddr)
				suite.Require().Nil(err)
				suite.Require().True(len(delegations) > 0)
				addr = sdk.MustAccAddressFromBech32(delegations[0].DelegatorAddress)

				// next block
				suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)

				params := suite.app.EstakingKeeper.GetParams(suite.ctx)
				params.StakeIncentives = &types.IncentiveInfo{
					EdenAmountPerYear: math.NewInt(1000_000_000_000_000),
					BlocksDistributed: 1,
				}
				params.MaxEdenRewardAprStakers = math.LegacyNewDec(1000_000)
				suite.app.EstakingKeeper.SetParams(suite.ctx, params)

				return addr, valAddr
			},
			func(addr sdk.AccAddress, valAddr sdk.ValAddress) {
				// update staker rewards
				err := suite.app.EstakingKeeper.UpdateStakersRewards(suite.ctx)
				suite.Require().Nil(err)

				distrAppModule := exdistr.NewAppModule(
					suite.app.AppCodec(), suite.app.DistrKeeper, suite.app.AccountKeeper,
					suite.app.CommitmentKeeper, suite.app.EstakingKeeper,
					&suite.app.AssetprofileKeeper,
					authtypes.FeeCollectorName, suite.app.GetSubspace(distrtypes.ModuleName))
				distrAppModule.AllocateEdenUsdcTokens(suite.ctx)
				distrAppModule.AllocateEdenBTokens(suite.ctx)

				// withdraw eden rewards
				msgServer := keeper.NewMsgServerImpl(*suite.app.EstakingKeeper)
				res, err := msgServer.WithdrawReward(suite.ctx, &types.MsgWithdrawReward{
					DelegatorAddress: addr.String(),
					ValidatorAddress: valAddr.String(),
				})
				suite.Require().Nil(err)
				suite.Require().Equal(res.Amount.String(), "147608ueden")
			},
		},
		//{
		//	"update stakers rewards missing asset profile base currency",
		//	func() (addr sdk.AccAddress, valAddr sdk.ValAddress) {
		//		suite.ResetSuite()
		//
		//		suite.app.AssetprofileKeeper.RemoveEntry(suite.ctx, ptypes.BaseCurrency)
		//
		//		return sdk.AccAddress{}, sdk.ValAddress{}
		//	},
		//	func(addr sdk.AccAddress, valAddr sdk.ValAddress) {
		//		// update staker rewards
		//		err := suite.app.EstakingKeeper.UpdateStakersRewards(suite.ctx)
		//		suite.Require().Error(err)
		//	},
		//},
		{
			"process rewards distribution",
			func() (addr sdk.AccAddress, valAddr sdk.ValAddress) {
				suite.ResetSuite()

				suite.app.AssetprofileKeeper.RemoveEntry(suite.ctx, ptypes.BaseCurrency)

				// create validator with 50% commission
				validators, err := suite.app.StakingKeeper.GetAllValidators(suite.ctx)
				suite.Require().Nil(err)
				suite.Require().True(len(validators) > 0)
				valAddr, err = sdk.ValAddressFromBech32(validators[0].GetOperator())
				suite.Require().Nil(err)

				delegations, err := suite.app.StakingKeeper.GetValidatorDelegations(suite.ctx, valAddr)
				suite.Require().Nil(err)
				suite.Require().True(len(delegations) > 0)
				addr = sdk.MustAccAddressFromBech32(delegations[0].DelegatorAddress)

				// next block
				suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)

				params := suite.app.EstakingKeeper.GetParams(suite.ctx)
				params.StakeIncentives = &types.IncentiveInfo{
					EdenAmountPerYear: math.NewInt(1000_000_000_000_000),
					BlocksDistributed: 1,
				}
				params.MaxEdenRewardAprStakers = math.LegacyNewDec(1000_000)
				suite.app.EstakingKeeper.SetParams(suite.ctx, params)

				return addr, valAddr
			},
			func(addr sdk.AccAddress, valAddr sdk.ValAddress) {
				// update staker rewards
				suite.app.EstakingKeeper.ProcessRewardsDistribution(suite.ctx)
			},
		},
		{
			"process update incentive params block height out of range",
			func() (addr sdk.AccAddress, valAddr sdk.ValAddress) {
				suite.ResetSuite()

				suite.SetAssetProfile()

				// create time based inflation item
				suite.app.TokenomicsKeeper.SetTimeBasedInflation(suite.ctx, ttypes.TimeBasedInflation{
					StartBlockHeight: 1,
					EndBlockHeight:   2,
					Inflation: &ttypes.InflationEntry{
						IcsStakingRewards: 1000,
						LmRewards:         0,
						CommunityFund:     0,
						StrategicReserve:  0,
						TeamTokensVested:  0,
					},
					Description: "test",
					Authority:   suite.app.GovKeeper.GetAuthority(),
				})

				return addr, valAddr
			},
			func(addr sdk.AccAddress, valAddr sdk.ValAddress) {
				// update staker rewards
				suite.app.EstakingKeeper.ProcessUpdateIncentiveParams(suite.ctx)
			},
		},
		{
			"process update incentive params",
			func() (addr sdk.AccAddress, valAddr sdk.ValAddress) {
				suite.ResetSuite()

				suite.SetAssetProfile()

				// create time based inflation item
				suite.app.TokenomicsKeeper.SetTimeBasedInflation(suite.ctx, ttypes.TimeBasedInflation{
					StartBlockHeight: 1,
					EndBlockHeight:   100,
					Inflation: &ttypes.InflationEntry{
						IcsStakingRewards: 1000,
						LmRewards:         0,
						CommunityFund:     0,
						StrategicReserve:  0,
						TeamTokensVested:  0,
					},
					Description: "test",
					Authority:   suite.app.GovKeeper.GetAuthority(),
				})

				suite.ctx = suite.ctx.WithBlockHeight(2)

				return addr, valAddr
			},
			func(addr sdk.AccAddress, valAddr sdk.ValAddress) {
				// update staker rewards
				suite.app.EstakingKeeper.ProcessUpdateIncentiveParams(suite.ctx)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			addr, valAddr := tc.prerequisiteFunction()
			tc.postValidateFunction(addr, valAddr)
		})
	}
}
