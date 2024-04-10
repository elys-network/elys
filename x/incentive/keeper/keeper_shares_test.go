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

func TestCalcTotalShareOfStaking(t *testing.T) {
	app := simapp.InitElysTestApp(initChain)
	ctx := app.BaseApp.NewContext(initChain, tmproto.Header{})

	ik := app.IncentiveKeeper

	// Generate 2 random accounts with 1000000uelys balanced
	addr := simapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(1000))

	var committed sdk.Coins
	var unclaimed sdk.Coins

	// Prepare unclaimed tokens
	uedenToken := sdk.NewCoin(ptypes.Eden, sdk.NewInt(1000))
	uedenBToken := sdk.NewCoin(ptypes.EdenB, sdk.NewInt(1000))
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
	uedenToken = sdk.NewCoin(ptypes.Eden, sdk.NewInt(1000))
	uedenBToken = sdk.NewCoin(ptypes.EdenB, sdk.NewInt(1000))
	committed = committed.Add(uedenToken, uedenBToken)

	// Eden
	err = app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, committed)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ctypes.ModuleName, addr[0], committed)
	require.NoError(t, err)

	err = app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, committed)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ctypes.ModuleName, addr[1], committed)
	require.NoError(t, err)

	// Add testing commitment
	simapp.AddTestCommitment(app, ctx, addr[0], committed)
	simapp.AddTestCommitment(app, ctx, addr[1], committed)

	commitment := app.CommitmentKeeper.GetCommitments(ctx, addr[0].String())

	require.Equal(t, commitment.CommittedTokens[0].Denom, ptypes.Eden)
	require.Equal(t, commitment.CommittedTokens[0].Amount, sdk.NewInt(1000))

	require.Equal(t, commitment.CommittedTokens[1].Denom, ptypes.EdenB)
	require.Equal(t, commitment.CommittedTokens[1].Amount, sdk.NewInt(1000))

	// Recalculate total committed info
	ik.UpdateTotalCommitmentInfo(ctx, ptypes.BaseCurrency)

	share1 := ik.CalcTotalShareOfStaking(sdk.ZeroInt())
	require.Equal(t, share1, sdk.ZeroDec())

	share2 := ik.CalcTotalShareOfStaking(sdk.NewInt(1004000))
	require.Equal(t, share2, sdk.NewDecWithPrec(1, 0))
}

func TestCalcDelegationAmount(t *testing.T) {
	app, genAccount, _ := simapp.InitElysTestAppWithGenAccount()
	ctx := app.BaseApp.NewContext(initChain, tmproto.Header{})

	ik := app.IncentiveKeeper
	addr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(1000))

	// Check with non-delegator
	delegatedAmount := ik.CalcDelegationAmount(ctx, addr[0].String())
	require.Equal(t, delegatedAmount, sdk.ZeroInt())

	// Check with genesis account (delegator)
	delegatedAmount = ik.CalcDelegationAmount(ctx, genAccount.String())
	require.Equal(t, delegatedAmount, sdk.DefaultPowerReduction)
}
