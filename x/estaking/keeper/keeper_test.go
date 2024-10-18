package keeper_test

import (
	"testing"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	simapp "github.com/elys-network/elys/app"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestEstakingExtendedFunctions(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

	stakingKeeper := app.StakingKeeper
	estakingKeeper := app.EstakingKeeper

	// create validator with 50% commission
	validators, err := stakingKeeper.GetAllValidators(ctx)
	require.NoError(t, err)
	require.True(t, len(validators) > 0)
	valAddr := validators[0].GetOperator()
	validator, err := sdk.ValAddressFromBech32(valAddr)
	require.NoError(t, err)
	delegations, err := stakingKeeper.GetValidatorDelegations(ctx, validator)
	require.NoError(t, err)
	require.True(t, len(delegations) > 0)
	addr := sdk.MustAccAddressFromBech32(delegations[0].DelegatorAddress)

	totalBonded := estakingKeeper.TotalBondedTokens(ctx)
	require.Equal(t, totalBonded.String(), "1000000")

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
	commitmentMsgServer := commitmentkeeper.NewMsgServerImpl(*app.CommitmentKeeper)
	_, err = commitmentMsgServer.CommitClaimedRewards(ctx, &commitmenttypes.MsgCommitClaimedRewards{
		Creator: addr.String(),
		Denom:   ptypes.Eden,
		Amount:  math.NewInt(1000_000),
	})
	require.Nil(t, err)

	totalBonded = estakingKeeper.TotalBondedTokens(ctx)
	require.Equal(t, totalBonded.String(), "2000000")

	edenVal := estakingKeeper.GetEdenValidator(ctx)
	require.Equal(t, edenVal.GetMoniker(), "EdenValidator")

	edenBVal := estakingKeeper.GetEdenBValidator(ctx)
	require.Equal(t, edenBVal.GetMoniker(), "EdenBValidator")

	edenValidator, err := sdk.ValAddressFromBech32(edenVal.GetOperator())
	require.NoError(t, err)
	edenBValidator, err := sdk.ValAddressFromBech32(edenBVal.GetOperator())
	require.NoError(t, err)

	edenValidatorI, _ := estakingKeeper.Validator(ctx, edenValidator)
	edenBValidatorI, _ := estakingKeeper.Validator(ctx, edenBValidator)

	require.Equal(t, edenValidatorI, edenVal)
	require.Equal(t, edenBValidatorI, edenBVal)

	edenDel, _ := estakingKeeper.Delegation(ctx, addr, edenValidator)
	require.Equal(t, edenDel.GetShares(), math.LegacyNewDec(1000_000))

	edenBDel, _ := estakingKeeper.Delegation(ctx, addr, edenBValidator)
	require.Nil(t, edenBDel)

	numDelegations := int64(0)
	estakingKeeper.IterateDelegations(ctx, addr, func(index int64, delegation stakingtypes.DelegationI) (stop bool) {
		numDelegations++
		return false
	})
	require.Equal(t, numDelegations, int64(2))

	numBondedValidators := int64(0)
	estakingKeeper.IterateBondedValidatorsByPower(ctx, func(index int64, delegation stakingtypes.ValidatorI) (stop bool) {
		numBondedValidators++
		return false
	})
	require.Equal(t, numBondedValidators, int64(2))
}
