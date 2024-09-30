package keeper_test

import (
	"cosmossdk.io/math"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	simapp "github.com/elys-network/elys/app"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/estaking/keeper"
	exdistr "github.com/elys-network/elys/x/estaking/modules/distribution"
	"github.com/elys-network/elys/x/estaking/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestUpdateStakersRewards(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true)

	stakingKeeper := app.StakingKeeper
	estakingKeeper := app.EstakingKeeper

	// create validator with 50% commission
	validators, err := stakingKeeper.GetAllValidators(ctx)
	require.True(t, len(validators) > 0)
	require.NoError(t, err)
	valAddr := validators[0].GetOperator()
	validator, err := sdk.ValAddressFromBech32(valAddr)
	require.NoError(t, err)
	delegations, err := stakingKeeper.GetValidatorDelegations(ctx, validator)
	require.NoError(t, err)
	require.True(t, len(delegations) > 0)
	addr := sdk.MustAccAddressFromBech32(delegations[0].DelegatorAddress)

	app.AssetprofileKeeper.SetEntry(ctx, assetprofiletypes.Entry{
		BaseDenom:       ptypes.Eden,
		Denom:           ptypes.Eden,
		Decimals:        6,
		CommitEnabled:   true,
		WithdrawEnabled: true,
	})
	app.AssetprofileKeeper.SetEntry(ctx, assetprofiletypes.Entry{
		BaseDenom:       ptypes.BaseCurrency,
		Denom:           ptypes.BaseCurrency,
		Decimals:        6,
		CommitEnabled:   true,
		WithdrawEnabled: true,
	})

	// next block
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	params := estakingKeeper.GetParams(ctx)
	params.StakeIncentives = &types.IncentiveInfo{
		EdenAmountPerYear: math.NewInt(1000_000_000_000_000),
		BlocksDistributed: 1,
	}
	params.MaxEdenRewardAprStakers = math.LegacyNewDec(1000_000)
	estakingKeeper.SetParams(ctx, params)

	// update staker rewards
	err = estakingKeeper.UpdateStakersRewards(ctx)
	require.Nil(t, err)

	distrAppModule := exdistr.NewAppModule(
		app.AppCodec(), app.DistrKeeper, app.AccountKeeper,
		app.CommitmentKeeper, &app.EstakingKeeper,
		&app.AssetprofileKeeper,
		authtypes.FeeCollectorName, app.GetSubspace(distrtypes.ModuleName))
	distrAppModule.AllocateTokens(ctx)

	// withdraw eden rewards
	msgServer := keeper.NewMsgServerImpl(estakingKeeper)
	res, err := msgServer.WithdrawReward(ctx, &types.MsgWithdrawReward{
		DelegatorAddress: addr.String(),
		ValidatorAddress: valAddr,
	})
	require.Nil(t, err)
	require.Equal(t, res.Amount.String(), "147608ueden")
}
