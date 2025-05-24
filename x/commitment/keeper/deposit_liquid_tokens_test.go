package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/v5/app"
	assetprofiletypes "github.com/elys-network/elys/v5/x/assetprofile/types"
	"github.com/elys-network/elys/v5/x/commitment/types"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"
	"github.com/stretchr/testify/require"
)

// TestDepositLiquidTokensWithNoEntry tests the deposit of liquid tokens with no assetprofile entry
func TestDepositLiquidTokensWithNoEntry(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(false)
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	addr := simapp.AddTestAddrs(app, ctx, 1, sdkmath.NewInt(1000000))

	// Mint 100ueden
	edenToken := sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(100)))

	err := app.BankKeeper.MintCoins(ctx, types.ModuleName, edenToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr[0], edenToken)
	require.NoError(t, err)

	creator := addr[0]

	// Set up the commitments for the creator
	commitments := types.Commitments{
		Creator: creator.String(),
		CommittedTokens: []*types.CommittedTokens{
			{
				Denom:  ptypes.Eden,
				Amount: sdkmath.NewInt(50),
			},
		},
		Claimed: sdk.Coins{sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(150))},
	}
	keeper.SetCommitments(ctx, commitments)

	// Deposit liquid eden to become claimed state
	err = keeper.DepositLiquidTokensClaimed(ctx, ptypes.Eden, sdkmath.NewInt(100), creator)
	require.Error(t, err)
}

// TestDepositLiquidTokensWithEntryDisabled tests the deposit of liquid tokens with assetprofile entry disabled
func TestDepositLiquidTokensWithEntryDisabled(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(false)
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	addr := simapp.AddTestAddrs(app, ctx, 1, sdkmath.NewInt(1000000))

	// Mint 100ueden
	edenToken := sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(100)))

	err := app.BankKeeper.MintCoins(ctx, types.ModuleName, edenToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr[0], edenToken)
	require.NoError(t, err)

	creator := addr[0]

	// Set assetprofile entry for denom
	app.AssetprofileKeeper.SetEntry(ctx, assetprofiletypes.Entry{BaseDenom: ptypes.Eden, CommitEnabled: false})

	// Set up the commitments for the creator
	commitments := types.Commitments{
		Creator: creator.String(),
		CommittedTokens: []*types.CommittedTokens{
			{
				Denom:  ptypes.Eden,
				Amount: sdkmath.NewInt(50),
			},
		},
		Claimed: sdk.Coins{sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(150))},
	}
	keeper.SetCommitments(ctx, commitments)

	// Deposit liquid eden to become claimed state
	err = keeper.DepositLiquidTokensClaimed(ctx, ptypes.Eden, sdkmath.NewInt(100), creator)
	require.Error(t, err)
}

func TestDepositLiquidTokens(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(true)
	// Create a test context and keeper
	keeper := app.CommitmentKeeper

	simapp.SetStakingParam(app, ctx)
	simapp.SetupAssetProfile(app, ctx)
	addr := simapp.AddTestAddrs(app, ctx, 1, sdkmath.NewInt(1000000))

	// Mint 100ueden
	edenToken := sdk.NewCoins(sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(100)))

	err := app.BankKeeper.MintCoins(ctx, types.ModuleName, edenToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr[0], edenToken)
	require.NoError(t, err)

	creator := addr[0]

	// Set assetprofile entry for denom
	app.AssetprofileKeeper.SetEntry(ctx, assetprofiletypes.Entry{BaseDenom: ptypes.Eden, CommitEnabled: true})

	// Set up the commitments for the creator
	commitments := types.Commitments{
		Creator: creator.String(),
		CommittedTokens: []*types.CommittedTokens{
			{
				Denom:  ptypes.Eden,
				Amount: sdkmath.NewInt(50),
			},
		},
		Claimed: sdk.Coins{sdk.NewCoin(ptypes.Eden, sdkmath.NewInt(150))},
	}
	keeper.SetCommitments(ctx, commitments)

	// Deposit liquid eden to become claimed state
	err = keeper.DepositLiquidTokensClaimed(ctx, ptypes.Eden, sdkmath.NewInt(100), creator)
	require.NoError(t, err)

	// Check if the deposit tokens were added to commitments
	newCommitments := keeper.GetCommitments(ctx, creator)

	// Check if the claimed tokens were updated correctly
	claimed := newCommitments.GetClaimedForDenom(ptypes.Eden)
	require.Equal(t, sdkmath.NewInt(250), claimed, "claimed tokens were not updated correctly")

	// Check if the committed tokens were updated correctly
	committedToken := newCommitments.GetCommittedAmountForDenom(ptypes.Eden)
	require.Equal(t, sdkmath.NewInt(50), committedToken, "committed tokens were not updated correctly")

	edenCoin := app.BankKeeper.GetBalance(ctx, addr[0], ptypes.Eden)
	require.Equal(t, edenCoin.Amount, sdkmath.ZeroInt())
}
