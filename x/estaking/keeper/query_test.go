package keeper_test

import (
	"testing"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	simapp "github.com/elys-network/elys/app"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	exdistr "github.com/elys-network/elys/x/estaking/modules/distribution"
	"github.com/elys-network/elys/x/estaking/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestQueryRewards(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

	stakingKeeper := app.StakingKeeper
	estakingKeeper := app.EstakingKeeper

	// create validator with 50% commission
	validators, _ := stakingKeeper.GetAllValidators(ctx)
	require.True(t, len(validators) > 0)
	valAddr, _ := sdk.ValAddressFromBech32(validators[0].GetOperator())
	delegations, _ := stakingKeeper.GetValidatorDelegations(ctx, valAddr)
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
	err := estakingKeeper.UpdateStakersRewards(ctx)
	require.Nil(t, err)

	distrAppModule := exdistr.NewAppModule(
		app.AppCodec(), app.DistrKeeper, app.AccountKeeper,
		app.CommitmentKeeper, app.EstakingKeeper,
		&app.AssetprofileKeeper,
		authtypes.FeeCollectorName, app.GetSubspace(distrtypes.ModuleName))
	distrAppModule.AllocateTokens(ctx)

	// query rewards
	res, err := estakingKeeper.Rewards(ctx, &types.QueryRewardsRequest{
		Address: addr.String(),
	})
	require.Nil(t, err)
	require.Equal(t, res.Total.String(), "147608ueden")
}
