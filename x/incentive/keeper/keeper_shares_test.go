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

func TestCalculateTotalShareOfStaking(t *testing.T) {
	app := simapp.InitElysTestApp(initChain)
	ctx := app.BaseApp.NewContext(initChain, tmproto.Header{})

	ik := app.IncentiveKeeper

	// Generate 2 random accounts with 1000000uelys balanced
	addr := simapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(1000))

	var committed []sdk.Coins
	var uncommitted []sdk.Coins

	// Prepare uncommitted tokens
	uedenToken := sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdk.NewInt(1000)))
	uedenBToken := sdk.NewCoins(sdk.NewCoin(ptypes.EdenB, sdk.NewInt(1000)))
	uncommitted = append(uncommitted, uedenToken)
	uncommitted = append(uncommitted, uedenBToken)

	// Eden
	err := app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, uedenToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ctypes.ModuleName, addr[0], uedenToken)
	require.NoError(t, err)

	err = app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, uedenToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ctypes.ModuleName, addr[1], uedenToken)
	require.NoError(t, err)

	// EdenB
	err = app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, uedenBToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ctypes.ModuleName, addr[0], uedenBToken)
	require.NoError(t, err)

	err = app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, uedenBToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ctypes.ModuleName, addr[1], uedenBToken)
	require.NoError(t, err)

	// Prepare committed tokens
	uedenToken = sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdk.NewInt(1000)))
	uedenBToken = sdk.NewCoins(sdk.NewCoin(ptypes.EdenB, sdk.NewInt(1000)))
	committed = append(committed, uedenToken)
	committed = append(committed, uedenBToken)

	// Eden
	err = app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, uedenToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ctypes.ModuleName, addr[0], uedenToken)

	err = app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, uedenToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ctypes.ModuleName, addr[1], uedenToken)

	// EdenB
	err = app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, uedenBToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ctypes.ModuleName, addr[0], uedenBToken)

	err = app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, uedenBToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ctypes.ModuleName, addr[1], uedenBToken)

	// Add testing commitment
	simapp.AddTestCommitment(app, ctx, addr[0], committed, uncommitted)
	simapp.AddTestCommitment(app, ctx, addr[1], committed, uncommitted)

	commitment, found := app.CommitmentKeeper.GetCommitments(ctx, addr[0].String())

	require.True(t, found)
	require.Equal(t, commitment.UncommittedTokens[0].Denom, ptypes.Eden)
	require.Equal(t, commitment.UncommittedTokens[0].Amount, sdk.NewInt(1000))

	require.Equal(t, commitment.UncommittedTokens[1].Denom, ptypes.EdenB)
	require.Equal(t, commitment.UncommittedTokens[1].Amount, sdk.NewInt(1000))

	require.Equal(t, commitment.CommittedTokens[0].Denom, ptypes.Eden)
	require.Equal(t, commitment.CommittedTokens[0].Amount, sdk.NewInt(1000))

	require.Equal(t, commitment.CommittedTokens[1].Denom, ptypes.EdenB)
	require.Equal(t, commitment.CommittedTokens[1].Amount, sdk.NewInt(1000))

	// Recalculate total committed info
	ik.UpdateTotalCommitmentInfo(ctx)

	share1 := ik.CalculateTotalShareOfStaking(sdk.ZeroInt())
	require.Equal(t, share1, sdk.ZeroDec())

	share2 := ik.CalculateTotalShareOfStaking(sdk.NewInt(1004000))
	require.Equal(t, share2, sdk.NewDecWithPrec(1, 0))
}

func TestCalculateDelegatedAmount(t *testing.T) {
	app, genAccount, _ := simapp.InitElysTestAppWithGenAccount()
	ctx := app.BaseApp.NewContext(initChain, tmproto.Header{})

	ik := app.IncentiveKeeper
	addr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(1000))

	// Check with non-delegator
	delegatedAmount := ik.CalculateDelegatedAmount(ctx, addr[0].String())
	require.Equal(t, delegatedAmount, sdk.ZeroInt())

	// Check with genesis account (delegator)
	delegatedAmount = ik.CalculateDelegatedAmount(ctx, genAccount.String())
	require.Equal(t, delegatedAmount, sdk.DefaultPowerReduction)
}
