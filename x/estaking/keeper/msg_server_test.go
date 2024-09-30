package keeper_test

import (
	"context"
	"cosmossdk.io/math"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	simapp "github.com/elys-network/elys/app"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	"github.com/elys-network/elys/x/estaking/keeper"
	"github.com/elys-network/elys/x/estaking/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.EstakingKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}

func TestMsgServer(t *testing.T) {
	ms, ctx := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
}

func TestWithdrawElysStakingRewards(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true)

	bankKeeper := app.BankKeeper
	stakingKeeper := app.StakingKeeper
	distrKeeper := app.DistrKeeper
	estakingKeeper := app.EstakingKeeper

	// create validator with 50% commission
	validators, _ := stakingKeeper.GetAllValidators(ctx)
	require.True(t, len(validators) > 0)
	valAddr, _ := sdk.ValAddressFromBech32(validators[0].GetOperator())
	delegations, _ := stakingKeeper.GetValidatorDelegations(ctx, valAddr)
	require.True(t, len(delegations) > 0)
	addr := sdk.MustAccAddressFromBech32(delegations[0].DelegatorAddress)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	// allocate some rewards
	initial := sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction)
	tokens := sdk.DecCoins{sdk.NewDecCoin(sdk.DefaultBondDenom, initial)}

	initialCoins := sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, initial)}
	err := bankKeeper.MintCoins(ctx, minttypes.ModuleName, initialCoins)
	require.Nil(t, err)
	err = bankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, disttypes.ModuleName, initialCoins)
	require.Nil(t, err)

	distrKeeper.AllocateTokensToValidator(ctx, validators[0], tokens)

	// historical count should be 4 (initial + latest for delegation)
	require.Equal(t, uint64(4), distrKeeper.GetValidatorHistoricalReferenceCount(ctx))

	// withdraw single rewards
	msgServer := keeper.NewMsgServerImpl(estakingKeeper)
	res, err := msgServer.WithdrawElysStakingRewards(ctx, &types.MsgWithdrawElysStakingRewards{
		DelegatorAddress: addr.String(),
	})
	require.Nil(t, err)
	require.NotEmpty(t, res.Amount.String())
}

func TestWithdrawReward_NormalValidator(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true)

	bankKeeper := app.BankKeeper
	stakingKeeper := app.StakingKeeper
	distrKeeper := app.DistrKeeper
	estakingKeeper := app.EstakingKeeper

	// create validator with 50% commission
	validators, _ := stakingKeeper.GetAllValidators(ctx)
	require.True(t, len(validators) > 0)
	valAddr, _ := sdk.ValAddressFromBech32(validators[0].GetOperator())
	delegations, _ := stakingKeeper.GetValidatorDelegations(ctx, valAddr)
	require.True(t, len(delegations) > 0)
	addr := sdk.MustAccAddressFromBech32(delegations[0].DelegatorAddress)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	// allocate some rewards
	initial := sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction)
	tokens := sdk.DecCoins{sdk.NewDecCoin(sdk.DefaultBondDenom, initial)}

	initialCoins := sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, initial)}
	err := bankKeeper.MintCoins(ctx, minttypes.ModuleName, initialCoins)
	require.Nil(t, err)
	err = bankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, disttypes.ModuleName, initialCoins)
	require.Nil(t, err)

	distrKeeper.AllocateTokensToValidator(ctx, validators[0], tokens)

	// historical count should be 4 (initial + latest for delegation)
	require.Equal(t, uint64(4), distrKeeper.GetValidatorHistoricalReferenceCount(ctx))

	// withdraw single rewards
	msgServer := keeper.NewMsgServerImpl(estakingKeeper)
	res, err := msgServer.WithdrawReward(ctx, &types.MsgWithdrawReward{
		DelegatorAddress: addr.String(),
		ValidatorAddress: valAddr.String(),
	})
	require.Nil(t, err)
	require.NotEmpty(t, res.Amount.String())
}

func TestWithdrawReward_EdenValidator(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true)

	bankKeeper := app.BankKeeper
	stakingKeeper := app.StakingKeeper
	distrKeeper := app.DistrKeeper
	estakingKeeper := app.EstakingKeeper

	// create validator with 50% commission
	validators, _ := stakingKeeper.GetAllValidators(ctx)
	require.True(t, len(validators) > 0)
	valAddr, _ := sdk.ValAddressFromBech32(validators[0].GetOperator())
	delegations, _ := stakingKeeper.GetValidatorDelegations(ctx, valAddr)
	require.True(t, len(delegations) > 0)
	addr := sdk.MustAccAddressFromBech32(delegations[0].DelegatorAddress)

	// set commitments
	commitments := app.CommitmentKeeper.GetCommitments(ctx, addr)
	commitments.AddClaimed(sdk.NewInt64Coin(ptypes.Eden, 1000_000))
	app.CommitmentKeeper.SetCommitments(ctx, commitments)
	app.AssetprofileKeeper.SetEntry(ctx, assetprofiletypes.Entry{
		BaseDenom:       ptypes.Eden,
		Denom:           ptypes.Eden,
		Decimals:        6,
		CommitEnabled:   true,
		WithdrawEnabled: true,
	})
	commitmentMsgServer := commitmentkeeper.NewMsgServerImpl(app.CommitmentKeeper)
	_, err := commitmentMsgServer.CommitClaimedRewards(sdk.WrapSDKContext(ctx), &commitmenttypes.MsgCommitClaimedRewards{
		Creator: addr.String(),
		Denom:   ptypes.Eden,
		Amount:  math.NewInt(1000_000),
	})
	require.Nil(t, err)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	// allocate some rewards
	initial := sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction)
	tokens := sdk.DecCoins{sdk.NewDecCoin(sdk.DefaultBondDenom, initial)}

	initialCoins := sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, initial)}
	err = bankKeeper.MintCoins(ctx, minttypes.ModuleName, initialCoins)
	require.Nil(t, err)
	err = bankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, disttypes.ModuleName, initialCoins)
	require.Nil(t, err)

	distrKeeper.AllocateTokensToValidator(ctx, validators[0], tokens)

	// historical count should be 5 (initial + latest for delegation)
	require.Equal(t, uint64(5), distrKeeper.GetValidatorHistoricalReferenceCount(ctx))

	// withdraw single rewards
	msgServer := keeper.NewMsgServerImpl(estakingKeeper)
	res, err := msgServer.WithdrawReward(ctx, &types.MsgWithdrawReward{
		DelegatorAddress: addr.String(),
		ValidatorAddress: estakingKeeper.GetParams(ctx).EdenCommitVal,
	})
	require.Nil(t, err)
	require.NotEmpty(t, res.Amount.String())
}

func TestWithdrawReward_EdenBValidator(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true)

	bankKeeper := app.BankKeeper
	stakingKeeper := app.StakingKeeper
	distrKeeper := app.DistrKeeper
	estakingKeeper := app.EstakingKeeper

	// create validator with 50% commission
	validators, _ := stakingKeeper.GetAllValidators(ctx)
	require.True(t, len(validators) > 0)
	valAddr, _ := sdk.ValAddressFromBech32(validators[0].GetOperator())
	delegations, _ := stakingKeeper.GetValidatorDelegations(ctx, valAddr)
	require.True(t, len(delegations) > 0)
	addr := sdk.MustAccAddressFromBech32(delegations[0].DelegatorAddress)

	// set commitments
	commitments := app.CommitmentKeeper.GetCommitments(ctx, addr)
	commitments.AddClaimed(sdk.NewInt64Coin(ptypes.EdenB, 1000_000))
	app.CommitmentKeeper.SetCommitments(ctx, commitments)
	app.AssetprofileKeeper.SetEntry(ctx, assetprofiletypes.Entry{
		BaseDenom:       ptypes.EdenB,
		Denom:           ptypes.EdenB,
		Decimals:        6,
		CommitEnabled:   true,
		WithdrawEnabled: true,
	})
	commitmentMsgServer := commitmentkeeper.NewMsgServerImpl(app.CommitmentKeeper)
	_, err := commitmentMsgServer.CommitClaimedRewards(sdk.WrapSDKContext(ctx), &commitmenttypes.MsgCommitClaimedRewards{
		Creator: addr.String(),
		Denom:   ptypes.EdenB,
		Amount:  math.NewInt(1000_000),
	})
	require.Nil(t, err)

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	// allocate some rewards
	initial := sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction)
	tokens := sdk.DecCoins{sdk.NewDecCoin(sdk.DefaultBondDenom, initial)}

	initialCoins := sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, initial)}
	err = bankKeeper.MintCoins(ctx, minttypes.ModuleName, initialCoins)
	require.Nil(t, err)
	err = bankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, disttypes.ModuleName, initialCoins)
	require.Nil(t, err)

	distrKeeper.AllocateTokensToValidator(ctx, validators[0], tokens)

	// historical count should be 5 (initial + latest for delegation)
	require.Equal(t, uint64(5), distrKeeper.GetValidatorHistoricalReferenceCount(ctx))

	// withdraw single rewards
	msgServer := keeper.NewMsgServerImpl(estakingKeeper)
	res, err := msgServer.WithdrawReward(ctx, &types.MsgWithdrawReward{
		DelegatorAddress: addr.String(),
		ValidatorAddress: estakingKeeper.GetParams(ctx).EdenbCommitVal,
	})
	require.Nil(t, err)
	require.NotEmpty(t, res.Amount.String())
}
