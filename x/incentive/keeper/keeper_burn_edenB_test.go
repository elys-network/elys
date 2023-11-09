package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	aptypes "github.com/elys-network/elys/x/assetprofile/types"
	commkeeper "github.com/elys-network/elys/x/commitment/keeper"
	ctypes "github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestBurnEdenBFromElysUnstaked(t *testing.T) {
	app, genAccount, valAddr := simapp.InitElysTestAppWithGenAccount()
	ctx := app.BaseApp.NewContext(initChain, tmproto.Header{})

	ik, sk := app.IncentiveKeeper, app.StakingKeeper

	var committed []sdk.Coins
	var unclaimed []sdk.Coins

	// Prepare unclaimed tokens
	uedenToken := sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdk.NewInt(2000)))
	uedenBToken := sdk.NewCoins(sdk.NewCoin(ptypes.EdenB, sdk.NewInt(20000)))
	unclaimed = append(unclaimed, uedenToken)
	unclaimed = append(unclaimed, uedenBToken)

	// Eden
	err := app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, uedenToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ctypes.ModuleName, genAccount, uedenToken)
	require.NoError(t, err)

	// EdenB
	err = app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, uedenBToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ctypes.ModuleName, genAccount, uedenBToken)
	require.NoError(t, err)

	// Prepare committed tokens
	uedenToken = sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdk.NewInt(10000)))
	uedenBToken = sdk.NewCoins(sdk.NewCoin(ptypes.EdenB, sdk.NewInt(5000)))
	committed = append(committed, uedenToken)
	committed = append(committed, uedenBToken)

	// Eden
	err = app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, uedenToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ctypes.ModuleName, genAccount, uedenToken)
	require.NoError(t, err)

	// EdenB
	err = app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, uedenBToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ctypes.ModuleName, genAccount, uedenBToken)
	require.NoError(t, err)

	// Add testing commitment
	simapp.AddTestCommitment(app, ctx, genAccount, committed, unclaimed)

	commitment, found := app.CommitmentKeeper.GetCommitments(ctx, genAccount.String())
	require.True(t, found)
	require.Equal(t, commitment.RewardsUnclaimed[1].Denom, ptypes.EdenB)
	require.Equal(t, commitment.RewardsUnclaimed[1].Amount, sdk.NewInt(20000))

	// Track Elys staked amount
	ik.EndBlocker(ctx)

	// burn amount = 100000 (unbonded amt) / (1000000 (elys staked) + 10000 (Eden committed)) * (20000 EdenB + 5000 EdenB committed)
	unbondAmt, err := sk.Unbond(ctx, genAccount, valAddr, sdk.NewDecWithPrec(10, 2))
	require.Equal(t, unbondAmt, sdk.NewInt(100000))

	commitment, found = app.CommitmentKeeper.GetCommitments(ctx, genAccount.String())
	require.True(t, found)

	require.Equal(t, commitment.RewardsUnclaimed[1].Denom, ptypes.EdenB)
	require.Equal(t, commitment.RewardsUnclaimed[1].Amount, sdk.NewInt(17525))
}

func TestBurnEdenBFromEdenUnclaimed(t *testing.T) {
	app, genAccount, _ := simapp.InitElysTestAppWithGenAccount()
	ctx := app.BaseApp.NewContext(initChain, tmproto.Header{})

	ik, cmk := app.IncentiveKeeper, app.CommitmentKeeper

	var committed []sdk.Coins
	var unclaimed []sdk.Coins

	// Prepare unclaimed tokens
	uedenToken := sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdk.NewInt(2000)))
	uedenBToken := sdk.NewCoins(sdk.NewCoin(ptypes.EdenB, sdk.NewInt(20000)))
	unclaimed = append(unclaimed, uedenToken)
	unclaimed = append(unclaimed, uedenBToken)

	// Eden
	err := app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, uedenToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ctypes.ModuleName, genAccount, uedenToken)
	require.NoError(t, err)

	// EdenB
	err = app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, uedenBToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ctypes.ModuleName, genAccount, uedenBToken)
	require.NoError(t, err)

	// Prepare committed tokens
	uedenToken = sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdk.NewInt(10000)))
	uedenBToken = sdk.NewCoins(sdk.NewCoin(ptypes.EdenB, sdk.NewInt(5000)))
	committed = append(committed, uedenToken)
	committed = append(committed, uedenBToken)

	// Eden
	err = app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, uedenToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ctypes.ModuleName, genAccount, uedenToken)
	require.NoError(t, err)

	// EdenB
	err = app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, uedenBToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ctypes.ModuleName, genAccount, uedenBToken)
	require.NoError(t, err)

	// Add testing commitment
	simapp.AddTestCommitment(app, ctx, genAccount, committed, unclaimed)

	commitment, found := app.CommitmentKeeper.GetCommitments(ctx, genAccount.String())
	require.True(t, found)
	require.Equal(t, commitment.RewardsUnclaimed[1].Denom, ptypes.EdenB)
	require.Equal(t, commitment.RewardsUnclaimed[1].Amount, sdk.NewInt(20000))

	// Track Elys staked amount
	ik.EndBlocker(ctx)

	// Set assetprofile entry for denom
	app.AssetprofileKeeper.SetEntry(ctx, aptypes.Entry{BaseDenom: ptypes.Eden, CommitEnabled: true})

	msg := &ctypes.MsgUncommitTokens{
		Creator: genAccount.String(),
		Amount:  sdk.NewInt(1000),
		Denom:   ptypes.Eden,
	}

	msgServer := commkeeper.NewMsgServerImpl(cmk)
	_, err = msgServer.UncommitTokens(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)

	// burn amount = 1000 (unclaimed amt) / (1000000 (elys staked) + 10000 (Eden committed)) * (20000 EdenB + 5000 EdenB committed)
	commitment, found = app.CommitmentKeeper.GetCommitments(ctx, genAccount.String())
	require.True(t, found)

	require.Equal(t, commitment.RewardsUnclaimed[1].Denom, ptypes.EdenB)
	require.Equal(t, commitment.RewardsUnclaimed[1].Amount, sdk.NewInt(19976))
}
