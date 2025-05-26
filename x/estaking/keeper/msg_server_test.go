package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	assetprofiletypes "github.com/elys-network/elys/v5/x/assetprofile/types"
	commitmentkeeper "github.com/elys-network/elys/v5/x/commitment/keeper"
	commitmenttypes "github.com/elys-network/elys/v5/x/commitment/types"
	"github.com/elys-network/elys/v5/x/estaking/keeper"
	"github.com/elys-network/elys/v5/x/estaking/types"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"
)

func (suite *EstakingKeeperTestSuite) TestMsgServer() {
	testCases := []struct {
		name                 string
		prerequisiteFunction func()
		postValidateFunction func()
	}{
		{
			"new msg server impl",
			func() {
				suite.ResetSuite()
			},
			func() {
				msgServer := keeper.NewMsgServerImpl(*suite.app.EstakingKeeper)

				suite.Require().NotNil(msgServer)
				suite.Require().NotNil(suite.ctx)
			},
		},
		{
			"withdraw elys staking rewards",
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

				// next block
				ctx := suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)

				// allocate some rewards
				initial := sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction)
				tokens := sdk.DecCoins{sdk.NewDecCoin(sdk.DefaultBondDenom, initial)}

				initialCoins := sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, initial)}
				err = suite.app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, initialCoins)
				suite.Require().Nil(err)
				err = suite.app.BankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, disttypes.ModuleName, initialCoins)
				suite.Require().Nil(err)

				suite.app.DistrKeeper.AllocateTokensToValidator(ctx, validators[0], tokens)

				// historical count should be 4 (initial + latest for delegation)
				suite.Require().Equal(uint64(4), suite.app.DistrKeeper.GetValidatorHistoricalReferenceCount(ctx))

				// withdraw single rewards
				msgServer := keeper.NewMsgServerImpl(*suite.app.EstakingKeeper)
				res, err := msgServer.WithdrawElysStakingRewards(ctx, &types.MsgWithdrawElysStakingRewards{
					DelegatorAddress: addr.String(),
				})
				suite.Require().Nil(err)
				suite.Require().NotEmpty(res.Amount.String())
			},
		},
		{
			"withdraw reward normal validator",
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

				// next block
				ctx := suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)

				// allocate some rewards
				initial := sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction)
				tokens := sdk.DecCoins{sdk.NewDecCoin(sdk.DefaultBondDenom, initial)}

				initialCoins := sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, initial)}
				err = suite.app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, initialCoins)
				suite.Require().Nil(err)
				err = suite.app.BankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, disttypes.ModuleName, initialCoins)
				suite.Require().Nil(err)

				suite.app.DistrKeeper.AllocateTokensToValidator(ctx, validators[0], tokens)

				// historical count should be 4 (initial + latest for delegation)
				suite.Require().Equal(uint64(4), suite.app.DistrKeeper.GetValidatorHistoricalReferenceCount(ctx))

				// withdraw single rewards
				msgServer := keeper.NewMsgServerImpl(*suite.app.EstakingKeeper)
				res, err := msgServer.WithdrawReward(ctx, &types.MsgWithdrawReward{
					DelegatorAddress: addr.String(),
					ValidatorAddress: valAddr.String(),
				})
				suite.Require().Nil(err)
				suite.Require().NotEmpty(res.Amount.String())
			},
		},
		{
			"withdraw reward eden validator",
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

				// set commitments
				commitments := suite.app.CommitmentKeeper.GetCommitments(suite.ctx, addr)
				commitments.AddClaimed(sdk.NewInt64Coin(ptypes.Eden, 1000_000))
				suite.app.CommitmentKeeper.SetCommitments(suite.ctx, commitments)
				suite.app.AssetprofileKeeper.SetEntry(suite.ctx, assetprofiletypes.Entry{
					BaseDenom:       ptypes.Eden,
					Denom:           ptypes.Eden,
					Decimals:        6,
					CommitEnabled:   true,
					WithdrawEnabled: true,
				})
				commitmentMsgServer := commitmentkeeper.NewMsgServerImpl(*suite.app.CommitmentKeeper)
				_, err = commitmentMsgServer.CommitClaimedRewards(suite.ctx, &commitmenttypes.MsgCommitClaimedRewards{
					Creator: addr.String(),
					Denom:   ptypes.Eden,
					Amount:  math.NewInt(1000_000),
				})
				suite.Require().Nil(err)

				// next block
				ctx := suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)

				// allocate some rewards
				initial := sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction)
				tokens := sdk.DecCoins{sdk.NewDecCoin(sdk.DefaultBondDenom, initial)}

				initialCoins := sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, initial)}
				err = suite.app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, initialCoins)
				suite.Require().Nil(err)
				err = suite.app.BankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, disttypes.ModuleName, initialCoins)
				suite.Require().Nil(err)

				suite.app.DistrKeeper.AllocateTokensToValidator(ctx, validators[0], tokens)

				// historical count should be 5 (initial + latest for delegation)
				suite.Require().Equal(uint64(5), suite.app.DistrKeeper.GetValidatorHistoricalReferenceCount(ctx))

				// withdraw single rewards
				msgServer := keeper.NewMsgServerImpl(*suite.app.EstakingKeeper)
				res, err := msgServer.WithdrawReward(ctx, &types.MsgWithdrawReward{
					DelegatorAddress: addr.String(),
					ValidatorAddress: suite.app.EstakingKeeper.GetParams(ctx).EdenCommitVal,
				})
				suite.Require().Nil(err)
				suite.Require().NotEmpty(res.Amount.String())
			},
		},
		{
			"withdraw reward eden b validator",
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

				// set commitments
				commitments := suite.app.CommitmentKeeper.GetCommitments(suite.ctx, addr)
				commitments.AddClaimed(sdk.NewInt64Coin(ptypes.EdenB, 1000_000))
				suite.app.CommitmentKeeper.SetCommitments(suite.ctx, commitments)
				suite.app.AssetprofileKeeper.SetEntry(suite.ctx, assetprofiletypes.Entry{
					BaseDenom:       ptypes.EdenB,
					Denom:           ptypes.EdenB,
					Decimals:        6,
					CommitEnabled:   true,
					WithdrawEnabled: true,
				})
				commitmentMsgServer := commitmentkeeper.NewMsgServerImpl(*suite.app.CommitmentKeeper)
				_, err = commitmentMsgServer.CommitClaimedRewards(suite.ctx, &commitmenttypes.MsgCommitClaimedRewards{
					Creator: addr.String(),
					Denom:   ptypes.EdenB,
					Amount:  math.NewInt(1000_000),
				})
				suite.Require().Nil(err)

				// next block
				ctx := suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)

				// allocate some rewards
				initial := sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction)
				tokens := sdk.DecCoins{sdk.NewDecCoin(sdk.DefaultBondDenom, initial)}

				initialCoins := sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, initial)}
				err = suite.app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, initialCoins)
				suite.Require().Nil(err)
				err = suite.app.BankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, disttypes.ModuleName, initialCoins)
				suite.Require().Nil(err)

				suite.app.DistrKeeper.AllocateTokensToValidator(ctx, validators[0], tokens)

				// historical count should be 5 (initial + latest for delegation)
				suite.Require().Equal(uint64(5), suite.app.DistrKeeper.GetValidatorHistoricalReferenceCount(ctx))

				// withdraw single rewards
				msgServer := keeper.NewMsgServerImpl(*suite.app.EstakingKeeper)
				res, err := msgServer.WithdrawReward(ctx, &types.MsgWithdrawReward{
					DelegatorAddress: addr.String(),
					ValidatorAddress: suite.app.EstakingKeeper.GetParams(ctx).EdenbCommitVal,
				})
				suite.Require().Nil(err)
				suite.Require().NotEmpty(res.Amount.String())
			},
		},
		{
			"update params",
			func() {
				suite.ResetSuite()
			},
			func() {
				msgServer := keeper.NewMsgServerImpl(*suite.app.EstakingKeeper)

				params := suite.app.EstakingKeeper.GetParams(suite.ctx)
				params.EdenCommitVal = "eden"
				msg := &types.MsgUpdateParams{
					Authority: suite.app.GovKeeper.GetAuthority(),
					Params:    params,
				}

				res, err := msgServer.UpdateParams(suite.ctx, msg)
				suite.Require().Nil(err)
				suite.Require().NotNil(res)
			},
		},
		{
			"update params invalid authority",
			func() {
				suite.ResetSuite()
			},
			func() {
				msgServer := keeper.NewMsgServerImpl(*suite.app.EstakingKeeper)

				params := suite.app.EstakingKeeper.GetParams(suite.ctx)
				params.EdenCommitVal = "eden"
				msg := &types.MsgUpdateParams{
					Authority: "invalid",
					Params:    params,
				}

				res, err := msgServer.UpdateParams(suite.ctx, msg)
				suite.Require().NotNil(err)
				suite.Require().Nil(res)
			},
		},
		{
			"withdraw all rewards",
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

				// next block
				ctx := suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)

				// allocate some rewards
				initial := sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction)
				tokens := sdk.DecCoins{sdk.NewDecCoin(sdk.DefaultBondDenom, initial)}

				initialCoins := sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, initial)}
				err = suite.app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, initialCoins)
				suite.Require().Nil(err)
				err = suite.app.BankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, disttypes.ModuleName, initialCoins)
				suite.Require().Nil(err)

				suite.app.DistrKeeper.AllocateTokensToValidator(ctx, validators[0], tokens)

				// historical count should be 4 (initial + latest for delegation)
				suite.Require().Equal(uint64(4), suite.app.DistrKeeper.GetValidatorHistoricalReferenceCount(ctx))

				// withdraw all rewards
				msgServer := keeper.NewMsgServerImpl(*suite.app.EstakingKeeper)
				res, err := msgServer.WithdrawAllRewards(ctx, &types.MsgWithdrawAllRewards{
					DelegatorAddress: addr.String(),
				})
				suite.Require().Nil(err)
				suite.Require().NotNil(res)
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
