package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	simapp "github.com/elys-network/elys/app"
	ctypes "github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	tokenomicskeeper "github.com/elys-network/elys/x/tokenomics/keeper"
	tokenomicstypes "github.com/elys-network/elys/x/tokenomics/types"
	"github.com/stretchr/testify/require"
)

func TestABCI_EndBlocker(t *testing.T) {
	app, genAccount, _ := simapp.InitElysTestAppWithGenAccount()
	ctx := app.BaseApp.NewContext(initChain, tmproto.Header{})

	ik := app.IncentiveKeeper

	var committed sdk.Coins
	var unclaimed sdk.Coins

	// Prepare unclaimed tokens
	uedenToken := sdk.NewCoin(ptypes.Eden, sdk.NewInt(2000))
	uedenBToken := sdk.NewCoin(ptypes.EdenB, sdk.NewInt(2000))
	unclaimed = unclaimed.Add(uedenToken, uedenBToken)

	// Mint coins
	err := app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, unclaimed)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ctypes.ModuleName, genAccount, unclaimed)
	require.NoError(t, err)

	// Add testing commitment
	simapp.AddTestCommitment(app, ctx, genAccount, committed, unclaimed)
	// Update Elys staked amount
	ik.EndBlocker(ctx)

	// Get elys staked
	elysStaked, found := ik.GetElysStaked(ctx, genAccount.String())
	require.Equal(t, found, true)
	require.Equal(t, elysStaked.Amount, sdk.DefaultPowerReduction)

	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	srv := tokenomicskeeper.NewMsgServerImpl(app.TokenomicsKeeper)

	expected := &tokenomicstypes.MsgCreateTimeBasedInflation{
		Authority:        authority,
		StartBlockHeight: uint64(1),
		EndBlockHeight:   uint64(6307200),
		Inflation: &tokenomicstypes.InflationEntry{
			LmRewards:         9999999,
			IcsStakingRewards: 9999999,
			CommunityFund:     9999999,
			StrategicReserve:  9999999,
			TeamTokensVested:  9999999,
		},
	}

	wctx := sdk.WrapSDKContext(ctx)
	_, err = srv.CreateTimeBasedInflation(wctx, expected)
	require.NoError(t, err)

	expected = &tokenomicstypes.MsgCreateTimeBasedInflation{
		Authority:        authority,
		StartBlockHeight: uint64(6307201),
		EndBlockHeight:   uint64(12614401),
		Inflation: &tokenomicstypes.InflationEntry{
			LmRewards:         9999999,
			IcsStakingRewards: 9999999,
			CommunityFund:     9999999,
			StrategicReserve:  9999999,
			TeamTokensVested:  9999999,
		},
	}
	_, err = srv.CreateTimeBasedInflation(wctx, expected)
	require.NoError(t, err)

	// Set tokenomics params
	listTimeBasdInflations := app.TokenomicsKeeper.GetAllTimeBasedInflation(ctx)

	// After the first year
	ctx = ctx.WithBlockHeight(1)
	paramSet := ik.ProcessUpdateIncentiveParams(ctx)
	require.Equal(t, paramSet, true)

	// Check if the params are correctly set
	params := ik.GetParams(ctx)
	require.Equal(t, len(params.StakeIncentives), 1)
	require.Equal(t, len(params.LpIncentives), 1)

	require.Equal(t, params.StakeIncentives[0].EdenAmountPerYear, sdk.NewInt(int64(listTimeBasdInflations[0].Inflation.IcsStakingRewards)))
	require.Equal(t, params.LpIncentives[0].EdenAmountPerYear, sdk.NewInt(int64(listTimeBasdInflations[0].Inflation.LmRewards)))

	// After the first year
	ctx = ctx.WithBlockHeight(6307210)

	// Incentive param should be empty
	stakerEpoch, stakeIncentive := ik.IsStakerRewardsDistributionEpoch(ctx)
	params = ik.GetParams(ctx)
	require.Equal(t, stakerEpoch, false)
	require.Equal(t, len(params.StakeIncentives), 0)

	// Incentive param should be empty
	lpEpoch, lpIncentive := ik.IsLPRewardsDistributionEpoch(ctx)
	params = ik.GetParams(ctx)
	require.Equal(t, lpEpoch, false)
	require.Equal(t, len(params.LpIncentives), 0)

	// After reading tokenomics again
	paramSet = ik.ProcessUpdateIncentiveParams(ctx)
	require.Equal(t, paramSet, true)
	// Check params
	stakerEpoch, stakeIncentive = ik.IsStakerRewardsDistributionEpoch(ctx)
	params = ik.GetParams(ctx)
	require.Equal(t, stakeIncentive.EdenAmountPerYear, sdk.NewInt(int64(listTimeBasdInflations[0].Inflation.IcsStakingRewards)))

	// Check params
	lpEpoch, lpIncentive = ik.IsLPRewardsDistributionEpoch(ctx)
	params = ik.GetParams(ctx)
	require.Equal(t, lpIncentive.EdenAmountPerYear, sdk.NewInt(int64(listTimeBasdInflations[0].Inflation.IcsStakingRewards)))
}
