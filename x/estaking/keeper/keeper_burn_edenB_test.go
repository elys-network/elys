package keeper_test

import (
	"testing"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	commkeeper "github.com/elys-network/elys/x/commitment/keeper"
	ctypes "github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestBurnEdenBFromElysUnstaked(t *testing.T) {
	app, genAccount, valAddr := simapp.InitElysTestAppWithGenAccount(t)
	ctx := app.BaseApp.NewContext(true)

	ek, sk := app.EstakingKeeper, app.StakingKeeper

	var committed sdk.Coins
	var unclaimed sdk.Coins

	// Prepare unclaimed tokens
	uedenToken := sdk.NewCoin(ptypes.Eden, math.NewInt(2000))
	uedenBToken := sdk.NewCoin(ptypes.EdenB, math.NewInt(20000))
	unclaimed = unclaimed.Add(uedenToken, uedenBToken)

	// Mint coins
	err := app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, unclaimed)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ctypes.ModuleName, genAccount, unclaimed)
	require.NoError(t, err)

	// Prepare committed tokens
	uedenToken = sdk.NewCoin(ptypes.Eden, math.NewInt(10000))
	uedenBToken = sdk.NewCoin(ptypes.EdenB, math.NewInt(5000))
	committed = committed.Add(uedenToken, uedenBToken)

	// Mint coins
	err = app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, committed)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ctypes.ModuleName, genAccount, committed)
	require.NoError(t, err)

	// Add testing commitment
	simapp.AddTestCommitment(app, ctx, genAccount, committed)

	// Take elys staked snapshot
	ek.TakeDelegationSnapshot(ctx, genAccount)

	// burn amount = 100000 (unbonded amt) / (1000000 (elys staked) + 10000 (Eden committed)) * (20000 EdenB + 5000 EdenB committed)
	unbondAmt, err := sk.Unbond(ctx, genAccount, valAddr, math.LegacyNewDecWithPrec(10, 2))
	require.Equal(t, unbondAmt, math.NewInt(100000))
	require.NoError(t, err)

	// Process EdenB burn operation
	ek.EndBlocker(ctx)
}

func TestBurnEdenBFromEdenUncommitted(t *testing.T) {
	app, genAccount, _ := simapp.InitElysTestAppWithGenAccount(t)
	ctx := app.BaseApp.NewContext(true)

	ek, commitmentKeeper := app.EstakingKeeper, app.CommitmentKeeper

	var committed sdk.Coins
	var unclaimed sdk.Coins

	// Prepare unclaimed tokens
	uedenToken := sdk.NewCoin(ptypes.Eden, math.NewInt(2000))
	uedenBToken := sdk.NewCoin(ptypes.EdenB, math.NewInt(20000))
	unclaimed = unclaimed.Add(uedenToken, uedenBToken)

	// Mint coins
	err := app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, unclaimed)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ctypes.ModuleName, genAccount, unclaimed)
	require.NoError(t, err)

	// Prepare committed tokens
	uedenToken = sdk.NewCoin(ptypes.Eden, math.NewInt(10000))
	uedenBToken = sdk.NewCoin(ptypes.EdenB, math.NewInt(5000))
	committed = committed.Add(uedenToken, uedenBToken)

	// Set assetprofile entry for denom
	app.AssetprofileKeeper.SetEntry(ctx, assetprofiletypes.Entry{BaseDenom: ptypes.Eden, CommitEnabled: true, WithdrawEnabled: true})

	commitment := app.CommitmentKeeper.GetCommitments(ctx, genAccount)
	commitment.Claimed = commitment.Claimed.Add(committed...)
	app.CommitmentKeeper.SetCommitments(ctx, commitment)

	msgServer := commkeeper.NewMsgServerImpl(*commitmentKeeper)
	_, err = msgServer.CommitClaimedRewards(ctx, &ctypes.MsgCommitClaimedRewards{
		Creator: genAccount.String(),
		Amount:  math.NewInt(1000),
		Denom:   ptypes.Eden,
	})
	require.NoError(t, err)

	// Track Elys staked amount
	ek.EndBlocker(ctx)

	// Uncommit tokens
	_, err = msgServer.UncommitTokens(sdk.WrapSDKContext(ctx), &ctypes.MsgUncommitTokens{
		Creator: genAccount.String(),
		Amount:  math.NewInt(1000),
		Denom:   ptypes.Eden,
	})
	require.NoError(t, err)
}
