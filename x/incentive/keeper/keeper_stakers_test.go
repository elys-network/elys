package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	ctypes "github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestCalculateRewardsForStakers(t *testing.T) {
	app := simapp.InitElysTestApp(initChain)
	ctx := app.BaseApp.NewContext(initChain, tmproto.Header{})

	ik := app.IncentiveKeeper

	// Generate 2 random accounts with 10000uelys balanced
	addr := simapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(10000))

	var committed sdk.Coins
	var unclaimed sdk.Coins

	// Prepare unclaimed tokens
	uedenToken := sdk.NewCoin(ptypes.Eden, sdk.NewInt(2000))
	uedenBToken := sdk.NewCoin(ptypes.EdenB, sdk.NewInt(2000))
	unclaimed = unclaimed.Add(uedenToken, uedenBToken)

	// Mint coins
	err := app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, unclaimed)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ctypes.ModuleName, addr[0], unclaimed)
	require.NoError(t, err)

	err = app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, unclaimed)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ctypes.ModuleName, addr[1], unclaimed)
	require.NoError(t, err)

	// Prepare committed tokens
	uedenToken = sdk.NewCoin(ptypes.Eden, sdk.NewInt(1500))
	uedenBToken = sdk.NewCoin(ptypes.EdenB, sdk.NewInt(500))
	committed = committed.Add(uedenToken, uedenBToken)

	// Mint coins
	err = app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, committed)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ctypes.ModuleName, addr[0], committed)
	require.NoError(t, err)

	err = app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, committed)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ctypes.ModuleName, addr[1], committed)
	require.NoError(t, err)

	// Add testing commitment
	simapp.AddTestCommitment(app, ctx, addr[0], committed, unclaimed)
	simapp.AddTestCommitment(app, ctx, addr[1], committed, unclaimed)

	commitment, found := app.CommitmentKeeper.GetCommitments(ctx, addr[0].String())

	require.True(t, found)
	require.Equal(t, commitment.RewardsUnclaimed[0].Denom, ptypes.Eden)
	require.Equal(t, commitment.RewardsUnclaimed[0].Amount, sdk.NewInt(2000))

	require.Equal(t, commitment.RewardsUnclaimed[1].Denom, ptypes.EdenB)
	require.Equal(t, commitment.RewardsUnclaimed[1].Amount, sdk.NewInt(2000))

	require.Equal(t, commitment.CommittedTokens[0].Denom, ptypes.Eden)
	require.Equal(t, commitment.CommittedTokens[0].Amount, sdk.NewInt(1500))

	require.Equal(t, commitment.CommittedTokens[1].Denom, ptypes.EdenB)
	require.Equal(t, commitment.CommittedTokens[1].Amount, sdk.NewInt(500))

	// Recalculate total committed info
	ik.UpdateTotalCommitmentInfo(ctx)

	totalEdenGiven := sdk.ZeroInt()
	totalRewardsGiven := sdk.ZeroInt()

	dexRevenueStakersAmt := sdk.NewDec(100000)
	edenAmountPerEpochStakers := sdk.NewInt(100000)
	// Calculate delegated amount per delegator
	delegatedAmt := sdk.NewInt(1000)
	// Calculate new unclaimed Eden tokens from Elys staked Eden & Eden boost committed, Dex rewards distribution
	newUnclaimedEdenTokens, dexRewards, _ := ik.CalculateRewardsForStakersByElysStaked(ctx, delegatedAmt, edenAmountPerEpochStakers, dexRevenueStakersAmt)
	totalEdenGiven = totalEdenGiven.Add(newUnclaimedEdenTokens)
	totalRewardsGiven = totalRewardsGiven.Add(dexRewards)

	// Calculate new unclaimed Eden tokens from Eden committed, Dex rewards distribution
	newUnclaimedEdenTokens, dexRewards = ik.CalculateRewardsForStakersByCommitted(ctx, delegatedAmt, edenAmountPerEpochStakers, dexRevenueStakersAmt)
	totalEdenGiven = totalEdenGiven.Add(newUnclaimedEdenTokens)
	totalRewardsGiven = totalRewardsGiven.Add(dexRewards)

	// Calculate new unclaimed Eden tokens from Eden boost committed, Dex rewards distribution
	newUnclaimedEdenTokens, dexRewards = ik.CalculateRewardsForStakersByCommitted(ctx, delegatedAmt, edenAmountPerEpochStakers, dexRevenueStakersAmt)
	totalEdenGiven = totalEdenGiven.Add(newUnclaimedEdenTokens)
	totalRewardsGiven = totalRewardsGiven.Add(dexRewards)

	require.Equal(t, totalEdenGiven, sdk.NewInt(291))
	require.Equal(t, totalRewardsGiven, sdk.NewInt(297))
}
