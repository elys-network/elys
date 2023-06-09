package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	ctypes "github.com/elys-network/elys/x/commitment/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestCalculateRewardsForLPs(t *testing.T) {
	app := simapp.InitElysTestApp(initChain)
	ctx := app.BaseApp.NewContext(initChain, tmproto.Header{})

	ik := app.IncentiveKeeper

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(10000))

	var committed []sdk.Coins
	var uncommitted []sdk.Coins

	// Prepare uncommitted tokens
	uedenToken := sdk.NewCoins(sdk.NewCoin("ueden", sdk.NewInt(2000)))
	uncommitted = append(uncommitted, uedenToken)

	err := app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, uedenToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ctypes.ModuleName, addr[0], uedenToken)

	// Prepare committed tokens
	uedenToken = sdk.NewCoins(sdk.NewCoin("ueden", sdk.NewInt(500)))
	lpToken1 := sdk.NewCoins(sdk.NewCoin("lp-elys-usdc", sdk.NewInt(500)))
	lpToken2 := sdk.NewCoins(sdk.NewCoin("lp-ueden-usdc", sdk.NewInt(2000)))
	committed = append(committed, uedenToken)
	committed = append(committed, lpToken1)
	committed = append(committed, lpToken2)

	err = app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, uedenToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ctypes.ModuleName, addr[0], uedenToken)
	require.NoError(t, err)

	err = app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, lpToken1)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ctypes.ModuleName, addr[0], lpToken1)
	require.NoError(t, err)

	err = app.BankKeeper.MintCoins(ctx, ctypes.ModuleName, lpToken2)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ctypes.ModuleName, addr[0], lpToken2)
	require.NoError(t, err)

	simapp.AddTestCommitment(app, ctx, addr[0], committed, uncommitted)

	commitment, found := app.CommitmentKeeper.GetCommitments(ctx, addr[0].String())

	require.True(t, found)
	require.Equal(t, commitment.UncommittedTokens[0].Denom, "ueden")
	require.Equal(t, commitment.UncommittedTokens[0].Amount, sdk.NewInt(2000))

	require.Equal(t, commitment.CommittedTokens[0].Denom, "ueden")
	require.Equal(t, commitment.CommittedTokens[0].Amount, sdk.NewInt(500))

	require.Equal(t, commitment.CommittedTokens[1].Denom, "lp-elys-usdc")
	require.Equal(t, commitment.CommittedTokens[1].Amount, sdk.NewInt(500))

	require.Equal(t, commitment.CommittedTokens[2].Denom, "lp-ueden-usdc")
	require.Equal(t, commitment.CommittedTokens[2].Amount, sdk.NewInt(2000))

	// Add dummy liquidity pool
	ik.Lpk = simapp.AddTestLiquidityPool(addr)
	require.Equal(t, ik.Lpk.CalculateTVL(), sdk.NewInt(2500))

	edenAmountPerEpochLp := sdk.NewInt(1000000)
	totalProxyTVL := ik.Lpk.CalculateProxyTVL()

	// Recalculate total committed info
	ik.UpdateTotalCommitmentInfo(ctx)

	// Calculate rewards for LPs
	newUncommittedEdenTokensLp, dexRewardsLp := ik.CalculateRewardsForLPs(ctx, totalProxyTVL, commitment, edenAmountPerEpochLp)
	require.Equal(t, newUncommittedEdenTokensLp, sdk.NewInt(999999))
	require.Equal(t, dexRewardsLp, sdk.NewInt(11000))
}
