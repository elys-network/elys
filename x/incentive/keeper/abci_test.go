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

func TestABCI_EndBlocker(t *testing.T) {
	app, genAccount, _ := simapp.InitElysTestAppWithGenAccount()
	ctx := app.BaseApp.NewContext(initChain, tmproto.Header{})

	ik := app.IncentiveKeeper

	var committed []sdk.Coins
	var unclaimed []sdk.Coins

	// Prepare unclaimed tokens
	uedenToken := sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdk.NewInt(2000)))
	uedenBToken := sdk.NewCoins(sdk.NewCoin(ptypes.EdenB, sdk.NewInt(2000)))
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

	// Add testing commitment
	simapp.AddTestCommitment(app, ctx, genAccount, committed, unclaimed)
	// Update Elys staked amount
	ik.EndBlocker(ctx)

	// Get elys staked
	elysStaked, found := ik.GetElysStaked(ctx, genAccount.String())
	require.Equal(t, found, true)
	require.Equal(t, elysStaked.Amount, sdk.DefaultPowerReduction)
}
